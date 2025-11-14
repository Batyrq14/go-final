package handler

import (
	"net/http"
	"qasynda/internal/models"
	"qasynda/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ReviewHandler struct {
	reviewService service.ReviewService
	validator     *validator.Validate
}

func NewReviewHandler(reviewService service.ReviewService) *ReviewHandler {
	return &ReviewHandler{
		reviewService: reviewService,
		validator:     validator.New(),
	}
}

func (h *ReviewHandler) Create(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Only clients can create reviews
	if userRole.(string) != string(models.RoleClient) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only clients can create reviews"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req models.CreateReviewRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	review, err := h.reviewService.Create(c.Request.Context(), userUUID, &req)
	if err != nil {
		if err == service.ErrUnauthorized || err == service.ErrBookingNotCompleted || err == service.ErrReviewExists {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create review"})
		return
	}

	c.JSON(http.StatusCreated, review)
}

func (h *ReviewHandler) GetByProviderID(c *gin.Context) {
	providerID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}

	reviews, err := h.reviewService.GetByProviderID(c.Request.Context(), providerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get reviews"})
		return
	}

	c.JSON(http.StatusOK, reviews)
}

