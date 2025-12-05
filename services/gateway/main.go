package main

import (
	"net/http"
	"strings"

	"qasynda/shared/pkg/auth"
	"qasynda/shared/pkg/config"
	"qasynda/shared/pkg/logger"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			return
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)
		claims, err := auth.ValidateToken(tokenString, cfg.JWTSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func main() {
	logger.Init()
	cfg := config.Load()

	clients := InitClients(cfg)
	handler := NewHandler(clients)

	r := gin.Default()

	// Public Routes
	api := r.Group("/api")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
	}

	// Protected Routes
	protected := api.Group("/")
	protected.Use(AuthMiddleware(cfg))
	{
		protected.GET("/profile", handler.GetProfile)

		protected.POST("/services", handler.CreateService)
		protected.GET("/services", handler.GetServices)
		protected.POST("/bookings", handler.CreateBooking)

		protected.GET("/chat/history", handler.GetChatHistory)
	}

	logger.Info("Gateway Service starting on " + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Error("failed to start gateway", err)
	}
}
