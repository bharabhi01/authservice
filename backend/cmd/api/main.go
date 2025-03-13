package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/bharabhi01/authservice/pkg/database"
	"github.com/bharabhi01/authservice/pkg/config" 
)

func main() {
	// Load environment variables
	if err := config.Load(); err != nil {
		log.Println("Warning: .env file not found, using system env")
	}

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize the database connection 
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.CloseDB()	

	router := gin.Default()

	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	setupRoutes(router, cfg)

	// HTTP server with Gin router
	server := &http.Server {
		Addr : ":" + cfg.Port,
		Handler : router,
		ReadTimeout : cfg.ReadTimeout,
		WriteTimeout : cfg.WriteTimeout,
	}

	// Start the server in a goroutine so that it doesn't block the main function
	go func() {
		log.Printf("Server starting on port %s in %s mode\n", cfg.Port, cfg.Env)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for an interrupt signal to gracefully shut down the server
	quit := make (chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<- quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting")
}