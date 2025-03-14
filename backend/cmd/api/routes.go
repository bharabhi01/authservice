package main

import (
	"github.com/gin-gonic/gin"
	"github.com/bharabhi01/authservice/internal/auth"
	"github.com/bharabhi01/authservice/internal/middleware"
	"github.com/bharabhi01/authservice/internal/user"
	"github.com/bharabhi01/authservice/pkg/config"     
	"github.com/bharabhi01/authservice/pkg/jwt" 
	"github.com/bharabhi01/authservice/pkg/audit"
	auditHandler "github.com/bharabhi01/authservice/internal/audit"
)

func setupRoutes(router *gin.Engine, cfg *config.Config) {
	jwt.Init(cfg.JWTSecret, cfg.JWTExpirationHours)
	
	userRepo := user.NewRepository()
	authRepo := auth.NewRepository()
	auditLogger := audit.NewLogger()

	auditHandlerInstance := auditHandler.NewHandler(auditLogger)
	authHandler := auth.NewHandler(userRepo, auditLogger)
	roleHandler := auth.NewRoleHandler(authRepo)

	router.Use(middleware.AuditMiddleware(auditLogger))

	public := router.Group("/api/v1")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
				"env": cfg.Env,
			})
		})

		auth := public.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}
	}

	protected := router.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware())
	{
		users := protected.Group("/users")
		{
			users.GET("/userinfo", authHandler.CurrentUserInfo)

			users.GET("/:id/roles", roleHandler.GetUserRoles)
			users.POST("/:id/roles", roleHandler.AssignRoleToUser)
			users.DELETE("/:id/roles/:roleId", roleHandler.RemoveRoleFromUser)

			users.GET("/:id/permissions/check", roleHandler.CheckPermission)
		}

		roles := protected.Group("/roles")
		roles.Use(middleware.RoleMiddleware("admin")) 
		{
			roles.GET("", roleHandler.GetRoles)
		}

		permissions := protected.Group("/permissions")
		permissions.Use(middleware.RoleMiddleware("admin")) 
		{
			permissions.GET("", roleHandler.GetPermissions)
		}

		// Add audit logs endpoint
		logs := protected.Group("/audit")
		logs.Use(middleware.RoleMiddleware("admin"))
		{
			logs.GET("/logs", auditHandlerInstance.GetLogs)
		}
	}
}