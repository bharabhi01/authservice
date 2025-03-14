package auth

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/bharabhi01/authservice/internal/user"
	"github.com/bharabhi01/authservice/pkg/jwt"
	"github.com/bharabhi01/authservice/pkg/audit"
)

type Handler struct {
	userRepo *user.Repository
	auditLogger *audit.Logger
}

func NewHandler(userRepo *user.Repository, auditLogger *audit.Logger) *Handler {
	return &Handler{
		userRepo: userRepo,
		auditLogger: auditLogger,
	}
}

func (h *Handler) Register(c *gin.Context) {
	var registration user.UserRegistration
	if err := c.ShouldBindJSON(&registration); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid registration data: " + err.Error(),
		})
		return
	}

	existingUser, err := h.userRepo.GetByUsername(registration.Username)
	if err == nil && existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error": "Username already exists",
		})
		return
	}

	newUser, err := h.userRepo.Create(&registration)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create user: " + err.Error(),
		})
		return
	}

	token, err := jwt.GenerateToken(newUser.ID, newUser.Username, newUser.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token: " + err.Error(),
		})
		return
	}

	if h.auditLogger != nil {
		details := map[string]interface{}{
			"username": newUser.Username,
			"email": newUser.Email,
		}
		h.auditLogger.LogFromGin(c, "REGISTER", "user", newUser.ID, details)
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user": newUser.ToResponse(),
		"token": token,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var login user.UserLogin
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid login data: " + err.Error(),
		})
		return
	}

	user, err := h.userRepo.GetByUsername(login.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	if !h.userRepo.VerifyPassword(user, login.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid credentials",
		})
		return
	}

	if !user.Active {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "Account is not active",
		})
		return
	}

	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to generate token: " + err.Error(),
		})
		return
	}

	if h.auditLogger != nil {
		details := map[string]interface{}{
			"username": user.Username,
		}
		h.auditLogger.LogFromGin(c, "LOGIN", "user", user.ID, details)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user": user.ToResponse(),
		"token": token,
	})
}

func (h *Handler) CurrentUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Not authenticated",
		})
		return
	}

	user, err := h.userRepo.GetByID(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to get user: " + err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": user.ToResponse(),
	})
}

