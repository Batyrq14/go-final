package main

import (
	"context"
	"net/http"
	"strconv"

	"qasynda/shared/pkg/logger"
	"qasynda/shared/pkg/models"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	clients *Clients
}

func NewHandler(clients *Clients) *Handler {
	return &Handler{clients: clients}
}

func (h *Handler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	res, err := h.clients.User.Register(ctx, &req)
	if err != nil {
		logger.Error("gateway: register failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	res, err := h.clients.User.Login(ctx, &req)
	if err != nil {
		logger.Error("gateway: login failed", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	ctx := c.Request.Context()
	res, err := h.clients.User.GetUser(ctx, &models.GetUserRequest{UserID: userID})
	if err != nil {
		logger.Error("gateway: get profile failed", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProviders(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	ctx := c.Request.Context()
	res, err := h.clients.User.ListProviders(ctx, &models.ListProvidersRequest{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateService(c *gin.Context) {
	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserID = c.GetString("user_id")

	ctx := c.Request.Context()
	res, err := h.clients.Marketplace.CreateService(ctx, &req)
	if err != nil {
		logger.Error("gateway: create service failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create service"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetServices(c *gin.Context) {
	category := c.Query("category")
	ctx := c.Request.Context()
	res, err := h.clients.Marketplace.GetServices(ctx, &models.GetServicesRequest{Category: category})
	if err != nil {
		logger.Error("gateway: get services failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load services"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateBooking(c *gin.Context) {
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserID = c.GetString("user_id")

	ctx := c.Request.Context()
	res, err := h.clients.Marketplace.CreateBooking(ctx, &req)
	if err != nil {
		logger.Error("gateway: create booking failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create booking"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetBookings(c *gin.Context) {
	userId := c.GetString("user_id")
	role := c.GetString("role")

	ctx := c.Request.Context()
	res, err := h.clients.Marketplace.ListBookings(ctx, &models.ListBookingsRequest{
		UserID: userId,
		Role:   role,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) UpdateBookingStatus(c *gin.Context) {
	bookingID := c.Param("id")
	var req models.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.BookingID = bookingID
	req.UserID = c.GetString("user_id")

	ctx := c.Request.Context()
	res, err := h.clients.Marketplace.UpdateBookingStatus(ctx, &req)
	if err != nil {
		logger.Error("gateway: update booking status failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update booking status"})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetChatHistory(c *gin.Context) {
	otherUserID := c.Query("other_user_id")
	if otherUserID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "other_user_id is required"})
		return
	}
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	req := &models.GetHistoryRequest{
		UserID1: c.GetString("user_id"),
		UserID2: otherUserID,
		Limit:   limit,
		Offset:  offset,
	}

	ctx := c.Request.Context()
	res, err := h.clients.Chat.GetHistory(ctx, req)
	if err != nil {
		logger.Error("gateway: get chat history failed", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load chat history"})
		return
	}

	c.JSON(http.StatusOK, res)
}

// HealthCheck returns the health status of the gateway service
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "gateway",
	})
}
