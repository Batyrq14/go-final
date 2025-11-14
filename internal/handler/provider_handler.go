package handler

import (
	"net/http"
	"qasynda/internal/models"
	"qasynda/internal/service"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ProviderHandler struct {
	providerService service.ProviderService
	validator       *validator.Validate
}

func NewProviderHandler(providerService service.ProviderService) *ProviderHandler {
	return &ProviderHandler{
		providerService: providerService,
		validator:       validator.New(),
	}
}

func (h *ProviderHandler) List(c *gin.Context) {
	filter := &models.ProviderFilter{
		Page:  1,
		Limit: 10,
	}

	// Parse query parameters
	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page > 0 {
			filter.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 && limit <= 100 {
			filter.Limit = limit
		}
	}

	if serviceIDStr := c.Query("service"); serviceIDStr != "" {
		if serviceID, err := uuid.Parse(serviceIDStr); err == nil {
			filter.ServiceID = &serviceID
		}
	}

	if city := c.Query("city"); city != "" {
		filter.City = &city
	}

	if minRatingStr := c.Query("min_rating"); minRatingStr != "" {
		if minRating, err := strconv.ParseFloat(minRatingStr, 64); err == nil {
			filter.MinRating = &minRating
		}
	}

	if availableStr := c.Query("available"); availableStr != "" {
		available := availableStr == "true"
		filter.IsAvailable = &available
	}

	if sortBy := c.Query("sort"); sortBy != "" {
		filter.SortBy = sortBy
	}

	if sortOrder := c.Query("order"); sortOrder != "" {
		filter.SortOrder = sortOrder
	}

	resp, err := h.providerService.List(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list providers"})
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ProviderHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	provider, err := h.providerService.GetByID(c.Request.Context(), id)
	if err != nil {
		if err == service.ErrProviderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get provider"})
		return
	}

	c.JSON(http.StatusOK, provider)
}

func (h *ProviderHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req models.UpdateProviderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	provider, err := h.providerService.Update(c.Request.Context(), id, userUUID, &req)
	if err != nil {
		if err == service.ErrProviderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update provider"})
		return
	}

	c.JSON(http.StatusOK, provider)
}

func (h *ProviderHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userUUID, err := uuid.Parse(userID.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.providerService.Delete(c.Request.Context(), id, userUUID); err != nil {
		if err == service.ErrProviderNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == service.ErrUnauthorized {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete provider"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "provider deleted"})
}

