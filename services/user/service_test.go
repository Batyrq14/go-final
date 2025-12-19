package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"qasynda/shared/pkg/auth"
	"qasynda/shared/pkg/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockStore struct {
	mock.Mock
}

func (m *MockStore) Create(ctx context.Context, user *User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *MockStore) GetByEmail(ctx context.Context, email string) (*User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockStore) GetByID(ctx context.Context, id uuid.UUID) (*User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockStore) ListProviders(limit, offset int) ([]*DetailedProvider, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]*DetailedProvider), args.Error(1)
}

func (m *MockStore) UpdateProviderStatus(ctx context.Context, userID uuid.UUID, isAvailable bool) error {
	args := m.Called(ctx, userID, isAvailable)
	return args.Error(0)
}

func (m *MockStore) GetProviderStatus(ctx context.Context, userID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID)
	return args.Bool(0), args.Error(1)
}

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStore := new(MockStore)
	server := NewServer(mockStore, "secret")

	r := gin.Default()
	r.POST("/register", server.Register)

	req := models.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
		Role:     "client",
	}

	mockStore.On("GetByEmail", mock.Anything, req.Email).Return(nil, nil)
	mockStore.On("Create", mock.Anything, mock.Anything).Return(nil)

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(body))
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, req.Email, resp.User.Email)
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStore := new(MockStore)
	server := NewServer(mockStore, "secret")

	r := gin.Default()
	r.POST("/login", server.Login)

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	uid := uuid.New()
	user := &User{
		ID:           uid,
		Email:        "test@example.com",
		PasswordHash: string(hashedPassword),
		Role:         "client",
	}

	mockStore.On("GetByEmail", mock.Anything, user.Email).Return(user, nil)

	req := models.LoginRequest{
		Email:    user.Email,
		Password: password,
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(body))
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.AuthResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.Token)
	assert.Equal(t, user.Email, resp.User.Email)
}

func TestValidateToken(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockStore := new(MockStore)
	server := NewServer(mockStore, "secret")

	r := gin.Default()
	r.POST("/validate", server.ValidateToken)

	uid := uuid.New()
	user := &User{
		ID:    uid,
		Email: "test@example.com",
		Role:  "client",
	}

	token, _ := auth.GenerateToken(uid.String(), user.Email, user.Role, "secret", 24*time.Hour)

	mockStore.On("GetByID", mock.Anything, uid).Return(user, nil)

	req := models.ValidateTokenRequest{
		Token: token,
	}

	body, _ := json.Marshal(req)
	w := httptest.NewRecorder()
	httpReq, _ := http.NewRequest("POST", "/validate", bytes.NewBuffer(body))
	r.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var resp models.UserResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	assert.NoError(t, err)
	assert.Equal(t, user.Email, resp.Email)
	assert.Equal(t, uid.String(), resp.ID)
}
