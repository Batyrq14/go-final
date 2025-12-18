package main

import (
	"net/http"
	"strconv"
	"time"

	"strings"

	"qasynda/shared/pkg/auth"
	"qasynda/shared/pkg/logger"
	"qasynda/shared/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	store     *UserStore
	jwtSecret string
}

func NewServer(store *UserStore, jwtSecret string) *Server {
	return &Server{
		store:     store,
		jwtSecret: jwtSecret,
	}
}

func (s *Server) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Basic validation to avoid storing obviously invalid data
	req.Email = strings.TrimSpace(req.Email)
	req.FullName = strings.TrimSpace(req.FullName)
	req.Phone = strings.TrimSpace(req.Phone)

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "password must be at least 8 characters"})
		return
	}
	if req.Role != "client" && req.Role != "provider" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "role must be either 'client' or 'provider'"})
		return
	}

	existing, err := s.store.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		logger.Error("failed to check existing user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if existing != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "email already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("failed to hash password", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	userID := uuid.New()
	user := &User{
		ID:           userID,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		Role:         req.Role,
		FullName:     req.FullName,
		Phone:        req.Phone,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.store.Create(c.Request.Context(), user); err != nil {
		logger.Error("failed to create user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	token, err := auth.GenerateToken(
		user.ID.String(),
		user.Email,
		user.Role,
		s.jwtSecret,
		24*time.Hour,
	)
	if err != nil {
		logger.Error("failed to generate token", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, &models.AuthResponse{
		Token: token,
		User: &models.UserResponse{
			ID:       user.ID.String(),
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
			Phone:    user.Phone,
		},
	})
}

func (s *Server) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.Email = strings.TrimSpace(req.Email)
	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email and password are required"})
		return
	}

	user, err := s.store.GetByEmail(c.Request.Context(), req.Email)
	if err != nil {
		logger.Error("failed to get user", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(
		user.ID.String(),
		user.Email,
		user.Role,
		s.jwtSecret,
		24*time.Hour,
	)
	if err != nil {
		logger.Error("failed to generate token", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, &models.AuthResponse{
		Token: token,
		User: &models.UserResponse{
			ID:       user.ID.String(),
			Email:    user.Email,
			FullName: user.FullName,
			Role:     user.Role,
			Phone:    user.Phone,
		},
	})
}

func (s *Server) ValidateToken(c *gin.Context) {
	var req models.ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	claims, err := auth.ValidateToken(req.Token, s.jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	uid, err := uuid.Parse(claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
		return
	}

	user, err := s.store.GetByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, &models.UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		Phone:    user.Phone,
	})
}

func (s *Server) GetUser(c *gin.Context) {
	usrIdStr := c.Param("id")
	if usrIdStr == "" {
		// Fallback to body query if needed or usually REST uses /users/:id
		// For internal calls we might strictly use body in some architectures but REST implies params
		// But let's support binding JSON for internal communication consistency if we want
		// ACTUALLY, internal REST usually implies URL params for GET.
		// Let's assume URL param :id
	}

	// For GetUserRequest struct compatibility, let's allow binding JSON if method is POST/PUT?
	// But GetUser is likely GET.
	// Let's use `c.Param("id")`.

	uid, err := uuid.Parse(usrIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := s.store.GetByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, &models.UserResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
		FullName: user.FullName,
		Role:     user.Role,
		Phone:    user.Phone,
	})
}

func (s *Server) ListProviders(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}
	// Protect the database from very large limits
	if limit > 100 {
		limit = 100
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	providers, err := s.store.ListProviders(limit, offset)
	if err != nil {
		logger.Error("failed to list providers", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	var respProviders []*models.ProviderResponse
	for _, p := range providers {
		respProviders = append(respProviders, &models.ProviderResponse{
			User: &models.UserResponse{
				ID:       p.ID.String(),
				Email:    p.Email,
				FullName: p.FullName,
				Role:     p.Role,
				Phone:    p.Phone,
			},
			Location:        p.Location,
			HourlyRate:      p.HourlyRate,
			ExperienceYears: int(p.ExperienceYears),
			Bio:             p.Bio,
			IsAvailable:     p.IsAvailable,
			Rating:          p.Rating,
			ProviderID:      p.ServiceProviderID.String(),
		})
	}

	c.JSON(http.StatusOK, &models.ListProvidersResponse{Providers: respProviders})
}
