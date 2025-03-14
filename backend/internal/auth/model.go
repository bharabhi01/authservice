package auth

type Role struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Permission struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RolePermission struct {
	RoleID string `json:"role_id"`
	PermissionID string `json:"permission_id"`
	CreatedAt time.Time `json:"created_at"`
}

type UserRole struct {
	UserID string `json:"user_id"`
	RoleID string `json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

type RoleResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
	Permissions []Permission `json:"permissions"`
}

type PermissionResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Description string `json:"description"`
}