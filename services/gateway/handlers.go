package main

import (
	"context"
	"net/http"
	"strconv"

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

	res, err := h.clients.User.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	res, err := h.clients.User.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	res, err := h.clients.User.GetUser(c.Request.Context(), &models.GetUserRequest{UserID: userID})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProviders(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	res, err := h.clients.User.ListProviders(c.Request.Context(), &models.ListProvidersRequest{
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

	res, err := h.clients.Marketplace.CreateService(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetServices(c *gin.Context) {
	category := c.Query("category")
	res, err := h.clients.Marketplace.GetServices(c.Request.Context(), &models.GetServicesRequest{Category: category})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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

	res, err := h.clients.Marketplace.CreateBooking(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetBookings(c *gin.Context) {
	userId := c.GetString("user_id")
	role := c.GetString("role")

	res, err := h.clients.Marketplace.ListBookings(c.Request.Context(), &models.ListBookingsRequest{
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

	res, err := h.clients.Marketplace.UpdateBookingStatus(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetChatHistory(c *gin.Context) {
	otherUserID := c.Query("other_user_id")
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

	res, err := h.clients.Chat.GetHistory(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
