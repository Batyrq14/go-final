package handler

import (
	"errors"
	"net/http"
	"qasynda/internal/models"
	"qasynda/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type ServiceHandler struct {
	serviceService service.ServiceService
	validator      *validator.Validate
}

func NewServiceHandler(serviceService service.ServiceService) *ServiceHandler {
	return &ServiceHandler{
		serviceService: serviceService,
		validator:      validator.New(),
	}
}

func (h *ServiceHandler) List(c *gin.Context) {
	services, err := h.serviceService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to list services"})
		return
	}

	c.JSON(http.StatusOK, services)
}

func (h *ServiceHandler) Create(c *gin.Context) {
	userRole, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// Only admins can create services
	if userRole.(string) != string(models.RoleAdmin) {
		c.JSON(http.StatusForbidden, gin.H{"error": "only admins can create services"})
		return
	}

	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.validator.Struct(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	service, err := h.serviceService.Create(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create service"})
		return
	}

	c.JSON(http.StatusCreated, service)
}

func (h *ServiceHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}

	svc, err := h.serviceService.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, service.ErrServiceNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get service"})
		return
	}
	c.JSON(http.StatusOK, svc)
}

