package main

import (
	"context"
	"net/http"
	"strconv"

	pb "qasynda/shared/proto"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	clients *Clients
}

func NewHandler(clients *Clients) *Handler {
	return &Handler{clients: clients}
}

// Auth Handlers
func (h *Handler) Register(c *gin.Context) {
	var req pb.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.clients.User.Register(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) Login(c *gin.Context) {
	var req pb.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	res, err := h.clients.User.Login(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetProfile(c *gin.Context) {
	userID := c.GetString("user_id")
	res, err := h.clients.User.GetUser(context.Background(), &pb.GetUserRequest{UserId: userID})
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// Marketplace Handlers
func (h *Handler) CreateService(c *gin.Context) {
	var req pb.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Inject user_id from auth if needed, but proto has it in request.
	// Ideally we override it with trusted claim.
	req.UserId = c.GetString("user_id")

	res, err := h.clients.Marketplace.CreateService(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) GetServices(c *gin.Context) {
	// Check query params if any
	category := c.Query("category")
	res, err := h.clients.Marketplace.GetServices(context.Background(), &pb.GetServicesRequest{Category: category})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *Handler) CreateBooking(c *gin.Context) {
	var req pb.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.UserId = c.GetString("user_id")

	res, err := h.clients.Marketplace.CreateBooking(context.Background(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// Chat Handlers
func (h *Handler) GetChatHistory(c *gin.Context) {
	otherUserID := c.Query("other_user_id")
	limitStr := c.DefaultQuery("limit", "20")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, _ := strconv.Atoi(limitStr)
	offset, _ := strconv.Atoi(offsetStr)

	req := &pb.GetHistoryRequest{
		UserId_1: c.GetString("user_id"),
		UserId_2: otherUserID,
		Limit:    int32(limit),
		Offset:   int32(offset),
	}

	res, err := h.clients.Chat.GetHistory(context.Background(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
