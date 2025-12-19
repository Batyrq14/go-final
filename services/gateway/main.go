package main

import (
	"context"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"qasynda/shared/pkg/auth"
	"qasynda/shared/pkg/config"
	"qasynda/shared/pkg/logger"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

func RateLimitMiddleware(rps float64, burst int) gin.HandlerFunc {
	limiter := rate.NewLimiter(rate.Limit(rps), burst)
	return func(c *gin.Context) {
		if !limiter.Allow() {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
			return
		}
		c.Next()
	}
}

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

	r.Use(RateLimitMiddleware(10, 20))

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

	protected := api.Group("/")
	protected.Use(AuthMiddleware(cfg))
	{
		protected.GET("/auth/me", handler.GetProfile)

		protected.POST("/services", handler.CreateService)
		protected.POST("/bookings", handler.CreateBooking)
		protected.GET("/bookings", handler.GetBookings)
		protected.PUT("/bookings/:id/status", handler.UpdateBookingStatus)

		protected.PUT("/providers/status", handler.UpdateProviderStatus)
		protected.GET("/providers/status", handler.GetProviderStatus)

		protected.GET("/chat/history", handler.GetChatHistory)
	}

	r.GET("/ws", func(c *gin.Context) {
		target := cfg.Services.ChatUrl
		remoteUrl, err := url.Parse(target)
		if err != nil {
			logger.Error("failed to parse chat url", err)
			c.Status(http.StatusInternalServerError)
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remoteUrl)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remoteUrl.Host
			req.URL.Scheme = remoteUrl.Scheme
			req.URL.Host = remoteUrl.Host
			req.URL.Path = "/ws"
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	})

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		logger.Info("Gateway Service starting on " + cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("failed to start gateway", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutting down Gateway Service...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Error("Gateway Service forced to shutdown", err)
	}

	logger.Info("Gateway Service exiting")
}
