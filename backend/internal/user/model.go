package user 

import (
	"time"
)

// Users table in the database
type User struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	PasswordHash string `json:"-"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Role string `json:"role"`
	Active bool `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// For input validation during registration
type UserRegistration struct {
	Username string `json:"username binding:"required, min=3,max=50`
	Email string `json:"email" binding:"required,email`
	Password string `json:"password" binding:"required,min=8`
	FirstName string `json:"first_name`
	LastName string `json:"last_name`
}

// For input validation during login
type UserLogin struct {
	Username string `json:"username" binding:"required`
	Password string `json:"password" binding:"required`
}

// Omits sensitive information like password hash
type UserResponse struct {
	ID string `json:"id"`
	Username string `json:"username"`
	Email string `json:"email"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Role string `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// Convert User to UserResponse to prevent leaking sensitive information
func (u *User) ToResponse() UserResponse {
	return UserResponse {
		ID: u.ID,
		Username: u.Username,
		Email: u.Email,
		FirstName: u.FirstName,
		LastName: u.LastName,
		Role: u.Role,
		CreatedAt: u.CreatedAt,
	}
}



