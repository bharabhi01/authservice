package auth

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// RoleHandler handles role and permission management requests
// It contains dependencies needed for role operations
type RoleHandler struct {
	authRepo *Repository
}

// NewRoleHandler creates a new role handler
// This function initializes the handler with required dependencies
func NewRoleHandler(authRepo *Repository) *RoleHandler {
	return &RoleHandler{
		authRepo: authRepo,
	}
}

// GetRoles returns all roles
// This endpoint lists all available roles
func (h *RoleHandler) GetRoles(c *gin.Context) {
	// Get all roles from the repository
	roles, err := h.authRepo.GetRoles()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get roles: " + err.Error(),
		})
		return
	}

	// Convert to response format
	var response []RoleResponse
	for _, role := range roles {
		// Get permissions for this role
		permissions, err := h.authRepo.GetRolePermissions(role.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get role permissions: " + err.Error(),
			})
			return
		}

		// Extract permission names
		permissionNames := make([]string, len(permissions))
		for i, p := range permissions {
			permissionNames[i] = p.Name
		}

		// Add to response
		response = append(response, RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissionNames,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": response,
	})
}

// GetPermissions returns all permissions
// This endpoint lists all available permissions
func (h *RoleHandler) GetPermissions(c *gin.Context) {
	// Get all permissions from the repository
	permissions, err := h.authRepo.GetPermissions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get permissions: " + err.Error(),
		})
		return
	}

	// Convert to response format
	var response []PermissionResponse
	for _, permission := range permissions {
		response = append(response, PermissionResponse{
			ID:          permission.ID,
			Name:        permission.Name,
			Description: permission.Description,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"permissions": response,
	})
}

// GetUserRoles returns all roles for a user
// This endpoint lists the roles assigned to a specific user
func (h *RoleHandler) GetUserRoles(c *gin.Context) {
	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Get user roles from the repository
	roles, err := h.authRepo.GetUserRoles(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user roles: " + err.Error(),
		})
		return
	}

	// Convert to response format
	var response []RoleResponse
	for _, role := range roles {
		// Get permissions for this role
		permissions, err := h.authRepo.GetRolePermissions(role.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to get role permissions: " + err.Error(),
			})
			return
		}

		// Extract permission names
		permissionNames := make([]string, len(permissions))
		for i, p := range permissions {
			permissionNames[i] = p.Name
		}

		// Add to response
		response = append(response, RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			Permissions: permissionNames,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"roles": response,
	})
}

// AssignRoleToUser assigns a role to a user
// This endpoint gives a user a new role
func (h *RoleHandler) AssignRoleToUser(c *gin.Context) {
	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Parse request body
	var request struct {
		RoleID string `json:"role_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request: " + err.Error(),
		})
		return
	}

	// Assign role to user
	err := h.authRepo.AssignRoleToUser(userID, request.RoleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to assign role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role assigned successfully",
	})
}

// RemoveRoleFromUser removes a role from a user
// This endpoint revokes a role from a user
func (h *RoleHandler) RemoveRoleFromUser(c *gin.Context) {
	// Get user ID and role ID from URL parameters
	userID := c.Param("id")
	roleID := c.Param("roleId")
	
	if userID == "" || roleID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID and Role ID are required",
		})
		return
	}

	// Remove role from user
	err := h.authRepo.RemoveRoleFromUser(userID, roleID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to remove role: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Role removed successfully",
	})
}

// CheckPermission checks if a user has a specific permission
// This endpoint is used for permission-based access control
func (h *RoleHandler) CheckPermission(c *gin.Context) {
	// Get user ID from URL parameter
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "User ID is required",
		})
		return
	}

	// Get permission name from query parameter
	permissionName := c.Query("permission")
	if permissionName == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Permission name is required",
		})
		return
	}

	// Check if user has the permission
	hasPermission, err := h.authRepo.HasPermission(userID, permissionName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check permission: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_permission": hasPermission,
	})
}