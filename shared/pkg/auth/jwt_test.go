package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret"
	userID := "user-123"
	email := "test@example.com"
	role := "admin"
	expiration := 1 * time.Hour

	token, err := GenerateToken(userID, email, role, secret, expiration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := ValidateToken(token, secret)
	assert.NoError(t, err)
	assert.Equal(t, userID, claims.UserID)
	assert.Equal(t, email, claims.Email)
	assert.Equal(t, role, claims.Role)
}

func TestInvalidToken(t *testing.T) {
	secret := "test-secret"
	token := "invalid-token-string"

	claims, err := ValidateToken(token, secret)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestExpiredToken(t *testing.T) {
	secret := "test-secret"
	userID := "user-123"
	email := "test@example.com"
	role := "admin"
	expiration := -1 * time.Hour

	token, err := GenerateToken(userID, email, role, secret, expiration)
	assert.NoError(t, err)

	claims, err := ValidateToken(token, secret)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}

func TestWrongSecret(t *testing.T) {
	secret := "test-secret"
	wrongSecret := "wrong-secret"
	userID := "user-123"
	email := "test@example.com"
	role := "admin"
	expiration := 1 * time.Hour

	token, err := GenerateToken(userID, email, role, secret, expiration)
	assert.NoError(t, err)

	claims, err := ValidateToken(token, wrongSecret)
	assert.Error(t, err)
	assert.Nil(t, claims)
	assert.Equal(t, ErrInvalidToken, err)
}
