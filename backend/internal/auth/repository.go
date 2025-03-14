package auth

import (
	"database/sql"
	"time"

	"github.com/bharabhi01/authservice/pkg/database"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{
		db: database.DB,
	}
}

func (r *Repository) GetRoles() ([]Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var roles []Role

	for rows.Next() {
		var role role

		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			return nil, err
		}

		roles = append(roles, role)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

func (r * Repository) GetRoleByID(id string) (*Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE id = $1
	`

	var role Role
	err := r.db.QueryRow(query, id).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, err
	}

	return &role, nil
}

func (r *Repository) GetRoleByName(name string) (*Role, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM roles
		WHERE name = $1	
	`

	var role Role
	err := r.db.QueryRow(query, name).Scan(
		&role.ID,
		&role.Name,
		&role.Description,
		&role.CreatedAt,
		&role.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &role, nil
}

func (r *Repository) GetPermissions() ([]Permission, error) {
	query := `
		SELECT id, name, description, created_at, updated_at
		FROM permissions
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission

	for rows.Next() {
		var permission Permission
		if err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

func (r *Repository) GetRolePermissions(roleID string) ([]Permission, error) {
	query := `
		SELECT p.id, p.name, p.description, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		WHERE rp.role_id = $1
		ORDER BY p.name
	`

	rows, err := r.db.Query(query, roleID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions []Permission
	for rows.Next() {
		var permission Permission
		if err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}	

func (r *Repository) GetUserRoles(userID string) ([]Role, error) {
	// SQL query to get roles for a user
	query := `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM roles r
		JOIN user_roles ur ON r.id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY r.name
	`

	// Execute the query
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(
			&role.ID,
			&role.Name,
			&role.Description,
			&role.CreatedAt,
			&role.UpdatedAt,
		); err != nil {
			return nil, err
		}
		roles = append(roles, role)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return roles, nil
}

// GetUserPermissions retrieves all permissions for a specific user
// This function is used to determine what permissions a user has through their roles
func (r *Repository) GetUserPermissions(userID string) ([]Permission, error) {
	// SQL query to get permissions for a user through their roles
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.created_at, p.updated_at
		FROM permissions p
		JOIN role_permissions rp ON p.id = rp.permission_id
		JOIN user_roles ur ON rp.role_id = ur.role_id
		WHERE ur.user_id = $1
		ORDER BY p.name
	`

	// Execute the query
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Process the results
	var permissions []Permission
	for rows.Next() {
		var permission Permission
		if err := rows.Scan(
			&permission.ID,
			&permission.Name,
			&permission.Description,
			&permission.CreatedAt,
			&permission.UpdatedAt,
		); err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}

	// Check for errors during iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return permissions, nil
}

// AssignRoleToUser assigns a role to a user
// This function is used to give a user a new role
func (r *Repository) AssignRoleToUser(userID, roleID string) error {
	// SQL query to insert a user-role relationship
	query := `
		INSERT INTO user_roles (user_id, role_id, created_at)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`

	// Execute the query
	_, err := r.db.Exec(query, userID, roleID, time.Now())
	return err
}

// RemoveRoleFromUser removes a role from a user
// This function is used to revoke a role from a user
func (r *Repository) RemoveRoleFromUser(userID, roleID string) error {
	// SQL query to delete a user-role relationship
	query := `
		DELETE FROM user_roles
		WHERE user_id = $1 AND role_id = $2
	`

	// Execute the query
	_, err := r.db.Exec(query, userID, roleID)
	return err
}

// HasPermission checks if a user has a specific permission
// This function is used for permission-based access control
func (r *Repository) HasPermission(userID, permissionName string) (bool, error) {
	// SQL query to check if a user has a permission through any of their roles
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM permissions p
			JOIN role_permissions rp ON p.id = rp.permission_id
			JOIN user_roles ur ON rp.role_id = ur.role_id
			WHERE ur.user_id = $1 AND p.name = $2
		)
	`

	// Execute the query
	var hasPermission bool
	err := r.db.QueryRow(query, userID, permissionName).Scan(&hasPermission)
	return hasPermission, err
}