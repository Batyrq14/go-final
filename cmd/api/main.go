package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"qasynda/internal/config"
	"qasynda/internal/handler"
	"qasynda/internal/middleware"
	"qasynda/internal/repository"
	"qasynda/internal/service"
	"qasynda/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Set Gin mode
	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db, err := sqlx.Connect("postgres", cfg.Database.URL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Configure database connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	providerRepo := repository.NewProviderRepository(db)
	serviceRepo := repository.NewServiceRepository(db)
	bookingRepo := repository.NewBookingRepository(db)
	reviewRepo := repository.NewReviewRepository(db)

	// Initialize notification worker
	notificationWorker := worker.NewNotificationWorker(5)
	notificationWorker.Start()
	defer notificationWorker.Shutdown()

	_ = worker.NewNotificationService(notificationWorker) // Notification service for future use

	// Initialize services
	authService := service.NewAuthService(userRepo, providerRepo, cfg)
	providerService := service.NewProviderService(providerRepo, userRepo, serviceRepo)
	serviceService := service.NewServiceService(serviceRepo)
	bookingService := service.NewBookingService(bookingRepo, providerRepo, serviceRepo, userRepo)
	reviewService := service.NewReviewService(reviewRepo, bookingRepo, providerRepo, userRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	providerHandler := handler.NewProviderHandler(providerService)
	serviceHandler := handler.NewServiceHandler(serviceService)
	bookingHandler := handler.NewBookingHandler(bookingService)
	reviewHandler := handler.NewReviewHandler(reviewService)

	// Setup router
	router := setupRouter(cfg, authHandler, providerHandler, serviceHandler, bookingHandler, reviewHandler)

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// Wait for interrupt signal for graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func setupRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	providerHandler *handler.ProviderHandler,
	serviceHandler *handler.ServiceHandler,
	bookingHandler *handler.BookingHandler,
	reviewHandler *handler.ReviewHandler,
) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.CORS())
	router.Use(middleware.RequestLogging())
	router.Use(middleware.RateLimit())

	// Serve static files and HTML
	// Try multiple paths to work in both local and Docker environments
	webPaths := []string{"./web", "web", "/root/web"}
	var webPath string
	for _, path := range webPaths {
		if _, err := os.Stat(path + "/index.html"); err == nil {
			webPath = path
			break
		}
	}
	
	if webPath != "" {
		router.Static("/static", webPath+"/static")
		router.GET("/", func(c *gin.Context) {
			c.File(webPath + "/index.html")
		})
	} else {
		log.Println("Warning: web directory not found, HTML frontend will not be available")
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API routes
	api := router.Group("/api")
	{
		// Auth routes (public)
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.GET("/me", middleware.AuthMiddleware(cfg), authHandler.GetMe)
		}

		// Service routes (public)
		services := api.Group("/services")
		{
			services.GET("", serviceHandler.List)
			services.GET("/:id", serviceHandler.GetByID)
			services.POST("", middleware.AuthMiddleware(cfg), middleware.RoleMiddleware("admin"), serviceHandler.Create)
		}

		// Provider routes
		providers := api.Group("/providers")
		{
			providers.GET("", providerHandler.List)        // Public
			providers.GET("/:id", providerHandler.GetByID) // Public
			providers.PUT("/:id", middleware.AuthMiddleware(cfg), providerHandler.Update)
			providers.DELETE("/:id", middleware.AuthMiddleware(cfg), providerHandler.Delete)
		}

		// Booking routes (protected)
		bookings := api.Group("/bookings")
		bookings.Use(middleware.AuthMiddleware(cfg))
		{
			bookings.POST("", bookingHandler.Create)
			bookings.GET("", bookingHandler.List)
			bookings.GET("/:id", bookingHandler.GetByID)
			bookings.PATCH("/:id/status", bookingHandler.UpdateStatus)
		}

		// Review routes
		reviews := api.Group("/reviews")
		{
			reviews.GET("/providers/:id", reviewHandler.GetByProviderID) // Public
			reviews.POST("", middleware.AuthMiddleware(cfg), reviewHandler.Create)
		}
	}

	return router
}
