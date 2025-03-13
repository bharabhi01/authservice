package main

import (
	"github.com/gin-gonic/gin"
	"github.com/bharabhi01/authservice/internal/auth"
	"github.com/bharabhi01/authservice/internal/middleware"
	"github.com/bharabhi01/authservice/internal/user"
	"github.com/bharabhi01/authservice/pkg/config"     
	"github.com/bharabhi01/authservice/pkg/jwt" 
)

func setupRoutes(router *gin.Engine, cfg *config.Config) {
	jwt.Init(cfg.JWTSecret, cfg.JWTExpirationHours)
	
	userRepo := user.NewRepository()
	authHandler := auth.NewHandler(userRepo)

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
			users.GET("/userinfo", authHandler.GetUserInfo)
		}

		// admin := protected.Group("/admin")
		// admin.Use(middleware.RoleMiddleware("admin"))
		// {

		// }
	}
}