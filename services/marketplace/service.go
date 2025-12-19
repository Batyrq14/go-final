package main

import (
	"net/http"
	"time"

	"qasynda/shared/pkg/logger"
	"qasynda/shared/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	store IStore
}

func NewServer(store IStore) *Server {
	return &Server{store: store}
}

func (s *Server) CreateService(c *gin.Context) {
	var req models.CreateServiceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := uuid.New()
	service := &Service{
		ID:          id,
		Name:        req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.store.CreateService(c.Request.Context(), service); err != nil {
		logger.Error("failed to create service", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, &models.ServiceResponse{
		ID:          service.ID.String(),
		Title:       service.Name,
		Description: service.Description,
	})
}

func (s *Server) GetServices(c *gin.Context) {

	services, err := s.store.ListServices(c.Request.Context())
	if err != nil {
		logger.Error("failed to list services", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	var respServices []*models.ServiceResponse
	for _, svc := range services {
		respServices = append(respServices, &models.ServiceResponse{
			ID:          svc.ID.String(),
			Title:       svc.Name,
			Description: svc.Description,
		})
	}

	c.JSON(http.StatusOK, &models.GetServicesResponse{
		Services: respServices,
	})
}

func (s *Server) CreateBooking(c *gin.Context) {
	var req models.CreateBookingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	serviceID, err := uuid.Parse(req.ServiceID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid service id"})
		return
	}
	clientID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	scheduledTime, err := time.Parse(time.RFC3339, req.ScheduledTime)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scheduled time format (use ISO8601/RFC3339)"})
		return
	}

	providerID, err := uuid.Parse(req.ProviderID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid provider id"})
		return
	}

	id := uuid.New()
	booking := &Booking{
		ID:            id,
		ClientID:      clientID,
		ProviderID:    providerID,
		ServiceID:     serviceID,
		ScheduledDate: scheduledTime,
		Status:        "pending",
		DurationHours: 1.0,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := s.store.CreateBooking(c.Request.Context(), booking); err != nil {
		logger.Error("failed to create booking", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, &models.BookingResponse{
		ID:     booking.ID.String(),
		Status: booking.Status,
	})
}

func (s *Server) ListBookings(c *gin.Context) {

	userID := c.Query("user_id")
	role := c.Query("role")

	if userID == "" {

	}

	bookings, err := s.store.ListBookings(c.Request.Context(), userID, role)
	if err != nil {
		logger.Error("failed to list bookings", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	var respBookings []*models.BookingDetails
	for _, b := range bookings {
		respBookings = append(respBookings, &models.BookingDetails{
			ID:             b.ID.String(),
			ServiceID:      b.ServiceID.String(),
			ClientID:       b.ClientID.String(),
			ProviderID:     b.ProviderID.String(),
			Status:         b.Status,
			ScheduledTime:  b.ScheduledDate.Format(time.RFC3339),
			ServiceTitle:   "Service " + b.ServiceID.String(),
			OtherPartyName: "User " + b.ClientID.String(),
		})
	}

	c.JSON(http.StatusOK, &models.ListBookingsResponse{
		Bookings: respBookings,
	})
}

func (s *Server) UpdateBookingStatus(c *gin.Context) {

	bookingID := c.Param("id")
	var req models.UpdateBookingStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := s.store.UpdateBookingStatus(c.Request.Context(), bookingID, req.Status)
	if err != nil {
		logger.Error("failed to update booking status", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, &models.BookingResponse{
		ID:     bookingID,
		Status: req.Status,
	})
}
