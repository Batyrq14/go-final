package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"qasynda/shared/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) CreateService(ctx context.Context, service *Service) error {
	args := m.Called(ctx, service)
	return args.Error(0)
}

func (m *MockStore) ListServices(ctx context.Context) ([]*Service, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*Service), args.Error(1)
}

func (m *MockStore) CreateBooking(ctx context.Context, booking *Booking) error {
	args := m.Called(ctx, booking)
	return args.Error(0)
}

func (m *MockStore) ListBookings(ctx context.Context, userID string, role string) ([]*Booking, error) {
	args := m.Called(ctx, userID, role)
	return args.Get(0).([]*Booking), args.Error(1)
}

func (m *MockStore) UpdateBookingStatus(ctx context.Context, bookingID string, status string) error {
	args := m.Called(ctx, bookingID, status)
	return args.Error(0)
}

func (m *MockStore) GetBooking(ctx context.Context, bookingID string) (*Booking, error) {
	args := m.Called(ctx, bookingID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*Booking), args.Error(1)
}

func TestCreateBooking(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStore := new(MockStore)
	server := NewServer(mockStore)

	r := gin.Default()
	r.POST("/bookings", server.CreateBooking)

	req := models.CreateBookingRequest{
		ServiceID:     uuid.New().String(),
		UserID:        uuid.New().String(),
		ProviderID:    uuid.New().String(),
		ScheduledTime: "2023-12-25T10:00:00Z",
	}

	mockStore.On("CreateBooking", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/bookings", bytes.NewBuffer(body))
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.BookingResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, "pending", resp.Status)
}

func TestUpdateBookingStatus(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStore := new(MockStore)
	server := NewServer(mockStore)

	r := gin.Default()
	r.PUT("/bookings/:id/status", server.UpdateBookingStatus)

	bookingID := uuid.New().String()
	req := models.UpdateBookingStatusRequest{
		BookingID: bookingID,
		Status:    "accepted",
		UserID:    "provider-user-1",
	}

	mockStore.On("GetBooking", mock.Anything, bookingID).Return(&Booking{
		ID:         uuid.MustParse(bookingID),
		ProviderID: uuid.New(),
	}, nil)

	mockStore.On("UpdateBookingStatus", mock.Anything, bookingID, "accepted").Return(nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("PUT", "/bookings/"+bookingID+"/status", bytes.NewBuffer(body))
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)
}
