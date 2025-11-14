package tests

import (
	"context"
	"qasynda/internal/config"
	"qasynda/internal/models"
	"qasynda/internal/service"
	"testing"
	"time"

	"github.com/google/uuid"
)

// Mock repository for testing
type mockUserRepository struct {
	users map[string]*models.User
}

func (m *mockUserRepository) Create(ctx context.Context, user *models.User) error {
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}
	return nil, service.ErrUserNotFound
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, service.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) Update(ctx context.Context, user *models.User) error {
	if _, exists := m.users[user.Email]; !exists {
		return service.ErrUserNotFound
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	for email, user := range m.users {
		if user.ID == id {
			delete(m.users, email)
			return nil
		}
	}
	return service.ErrUserNotFound
}

type mockProviderRepository struct{}

func (m *mockProviderRepository) Create(ctx context.Context, provider *models.ServiceProvider) error {
	return nil
}

func (m *mockProviderRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.ServiceProvider, error) {
	return nil, nil
}

func (m *mockProviderRepository) GetByUserID(ctx context.Context, userID uuid.UUID) (*models.ServiceProvider, error) {
	return nil, nil
}

func (m *mockProviderRepository) Update(ctx context.Context, provider *models.ServiceProvider) error {
	return nil
}

func (m *mockProviderRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return nil
}

func (m *mockProviderRepository) List(ctx context.Context, filter *models.ProviderFilter) ([]*models.ServiceProvider, int, error) {
	return nil, 0, nil
}

func (m *mockProviderRepository) AddService(ctx context.Context, providerID, serviceID uuid.UUID) error {
	return nil
}

func (m *mockProviderRepository) RemoveService(ctx context.Context, providerID, serviceID uuid.UUID) error {
	return nil
}

func (m *mockProviderRepository) GetProviderServices(ctx context.Context, providerID uuid.UUID) ([]*models.Service, error) {
	return nil, nil
}

func TestAuthService_Register(t *testing.T) {
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret:     "test-secret",
			Expiration: 24 * time.Hour,
		},
	}

	userRepo := &mockUserRepository{users: make(map[string]*models.User)}
	providerRepo := &mockProviderRepository{}
	authService := service.NewAuthService(userRepo, providerRepo, cfg)

	tests := []struct {
		name    string
		req     *models.RegisterRequest
		wantErr bool
		errType error
	}{
		{
			name: "valid registration",
			req: &models.RegisterRequest{
				Email:    "newuser@example.com",
				Password: "password123",
				Role:     models.RoleClient,
				FullName: "New User",
				Phone:    "1234567890",
			},
			wantErr: false,
		},
		{
			name: "duplicate email",
			req: &models.RegisterRequest{
				Email:    "existing@example.com",
				Password: "password123",
				Role:     models.RoleClient,
				FullName: "Existing User",
				Phone:    "1234567890",
			},
			wantErr: true,
			errType: service.ErrEmailExists,
		},
	}

	// Setup: create existing user
	existingUser := &models.User{
		ID:        uuid.New(),
		Email:     "existing@example.com",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userRepo.users["existing@example.com"] = existingUser

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			resp, err := authService.Register(ctx, tt.req)

			if (err != nil) != tt.wantErr {
				t.Errorf("AuthService.Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Error("AuthService.Register() returned nil response")
				}
				if resp.Token == "" {
					t.Error("AuthService.Register() returned empty token")
				}
			} else if err != tt.errType {
				t.Errorf("AuthService.Register() error = %v, want %v", err, tt.errType)
			}
		})
	}
}

