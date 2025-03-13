package user

import (
	"database/sql"
	"errors"
	"time"

	"github.com/bharabhi01/authservice/pkg/database"
	"golang.org/x/crypto/bcrypt"
)

type Repository struct {
	db *sql.DB
}

func NewRepository() *Repository {
	return &Repository{
		db: database.DB,
	}
}

func (r *Repository) Create(user *UserRegistration) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()

	newUser := &User{
		Username: user.Username,
		Email: user.Email,
		PasswordHash: string(hashedPassword),
		FirstName: user.FirstName,
		LastName: user.LastName,
		Role: "user",
		Active: true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	query := `
		INSERT INTO users (username, email, password_hash, first_name, last_name, role, active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		RETURNING id, username, email, password_hash, first_name, last_name, role, active, created_at, updated_at
	`
	err := r.db.QueryRow (
		query,
		newUser.Username,
		newUser.Email,
		newUser.PasswordHash,
		newUser.FirstName,
		newUser.LastName,
		newUser.Role,
		newUser.Active,
		newUser.CreatedAt,
		newUser.UpdatedAt,
	).Scan(
		&newUser.ID,
		&newUser.Username,
		&newUser.Email,
		&newUser.PasswordHash,
		&newUser.FirstName,
		&newUser.LastName,
		&newUser.Role,
		&newUser.Active,
		&newUser.CreatedAt,
		&newUser.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return newUser, nil
}

func (r *Repository) GetByUsername(username string) (*User, error) {
	user := &User{}

	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, active, created_at, updated_at
		FROM users
		WHERE username = $1
	`

	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) GetByID(id string) (*User, error) {
	user := &User{}

	query := `
		SELECT id, username, email, password_hash, first_name, last_name, role, active, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.FirstName,
		&user.LastName,
		&user.Role,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

func (r *Repository) VerifyPassword(user *User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	return err == nil
}
