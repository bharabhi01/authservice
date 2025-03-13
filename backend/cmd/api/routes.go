package main

import (
	"github.com/gin-gonic/gin"
	"github.com/bharabhi01/authservice/internal/auth"
	"github.com/bharabhi01/authservice/internal/middleware"
	"github.com/bharabhi01/authservice/internal/user"
)

func setupRoutes(router *gin.Engine) {
	userRepo := user.NewRepository()
	authHandler := auth.NewHandler(userRepo)

	public := router.Group("/api/v1")
	{
		public.http.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status": "ok",
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