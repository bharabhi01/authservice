package audit 

import (
	"context"
	"database/sql"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/bharabhi01/authservice/pkg/database" 
)

type LogEntry struct {
	ID string `json:"id"`
	UserID string `json:"user_id"`
	Action string `json:"action"`
	ResourceType string `json:"resource_type"`
	ResourceID string `json:"resource_id"`
	IPAddress string `json:"ip_address"`
	UserAgent string `json:"user_agent"`
	Details json.RawMessage `json:"details"`
	CreatedAt time.Time `json:"created_at"`
}

type Logger struct {
	db *sql.DB
}

func NewLogger() *Logger {
	return &Logger{
		db: database.DB,
	}
}

func (l *Logger) Log(ctx context.Context, entry *LogEntry) error {
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}

	query := `
		INSERT INTO audit_logs (user_id, action, resource_type, resource_id, ip_address, user_agent, details, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id
	`

	err := l.db.QueryRowContext(
		ctx,
		query,
		entry.UserID,
		entry.Action,
		entry.ResourceType,
		entry.ResourceID,
		entry.IPAddress,
		entry.UserAgent,
		entry.Details,
		entry.CreatedAt,
	).Scan(&entry.ID)

	return err
}

func (l *Logger) LogFromGin(c *gin.Context, action, resourceType, resourceID string, details interface{}) error {
	userID, _ := c.Get("user_id")
	userIDStr, ok := userID.(string)
	if !ok {
		userIDStr = ""
	}

	var detailsJSON json.RawMessage
	if details != nil {
		detailsBytes, err := json.Marshal(details)
		if err != nil {
			return err
		}
		detailsJSON = detailsBytes
	}

	entry := &LogEntry{
		UserID: userIDStr,
		Action: action,
		ResourceType: resourceType,
		ResourceID: resourceID,
		IPAddress: c.ClientIP(),
		UserAgent: c.Request.UserAgent(),
		Details: detailsJSON,
		CreatedAt: time.Now(),
	}

	return l.Log(c.Request.Context(), entry)
}

func (l *Logger) GetLogs(ctx context.Context, userID, action, resourceType string, limit, offset int) ([]LogEntry, error) {
	query := `
		SELECT id, user_id, action, resource_type, resource_id, ip_address, user_agent, details, created_at
		FROM audit_logs
		WHERE 1=1
	`

	args := []interface{}{}
	argCount := 1

	if userID != "" {
		query += ` AND user_id = $` + string(argCount)
		args = append(args, userID)
		argCount++
	}

	if action != "" {
		query += ` AND action = $` + string(argCount)
		args = append(args, action)
		argCount++
	}

	if resourceType != "" {
		query += ` AND resource_type = $` + string(argCount)
		args = append(args, resourceType)
		argCount++
	}

	query += ` ORDER BY created_at DESC LIMIT $` + string(argCount) + ` OFFSET $` + string(argCount + 1)
	args = append(args, limit, offset)

	rows, err := l.db.QueryRowContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []LogEntry
	for rows.Next() {
		var log LogEntry
		if err := rows.Scan(
			&log.ID,
			&log.UserID,
			&log.Action,
			&log.ResourceType,
			&log.ResourceID,
			&log.IPAddress,
			&log.UserAgent,
			&log.Details,
			&log.CreatedAt,
		); err != nil {
			return nil, err
		}
		logs = append(logs, log)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return logs, nil
}