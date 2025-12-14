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

	// CORS Config
	r.Use(func(c *gin.Context) {
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

	// Serve Frontend
	r.Static("/static", "./frontend/static")
	r.StaticFile("/", "./frontend/index.html")

	// Public Routes
	api := r.Group("/api")
	{
		api.GET("/services", handler.GetServices)
		api.GET("/providers", handler.GetProviders)

		auth := api.Group("/auth")
		{
			auth.POST("/register", handler.Register)
			auth.POST("/login", handler.Login)
		}
	}

	// Protected Routes
	protected := api.Group("/")
	protected.Use(AuthMiddleware(cfg))
	{
		protected.GET("/auth/me", handler.GetProfile) // Frontend calls /api/auth/me

		protected.POST("/services", handler.CreateService)
		protected.POST("/bookings", handler.CreateBooking)
		protected.GET("/bookings", handler.GetBookings)
		protected.PUT("/bookings/:id/status", handler.UpdateBookingStatus)

		protected.PUT("/providers/status", handler.UpdateProviderStatus)
		protected.GET("/providers/status", handler.GetProviderStatus)

		protected.GET("/chat/history", handler.GetChatHistory)
	}

	logger.Info("Gateway Service starting on " + cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Error("failed to start gateway", err)
	}
}
