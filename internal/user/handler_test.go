package user

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetProfileHandler_ValidUserID(t *testing.T) {
	// Setup
	userID := uuid.New()
	mockService := &mockService{
		mockGetProfile: func(ctx context.Context, id uuid.UUID) (*User, error) {
			if id == userID {
				return &User{
					ID:    userID,
					Name:  "Test User",
					Email: "test@example.com",
					Role:  "candidate",
				}, nil
			}
			return nil, ErrInvalidEmail
		},
	}

	handler := NewGinHandler(mockService)
	router := gin.New()
	router.GET("/profile", handler.GetProfileHandler)

	// Create request with userID in context
	req, _ := http.NewRequest("GET", "/profile", nil)
	w := httptest.NewRecorder()

	c, _ := gin.CreateTestContext(w)
	c.Request = req
	c.Set("userID", userID.String())

	// Execute
	handler.GetProfileHandler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response["user"])
}

func TestGetProfileHandler_NoAuthentication(t *testing.T) {
	// Setup
	mockService := &mockService{}
	handler := NewGinHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/profile", nil)

	// Execute
	handler.GetProfileHandler(c)

	// Assert
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestUpdateProfileHandler_ValidUpdate(t *testing.T) {
	// Setup
	userID := uuid.New()
	updatedUser := &User{
		ID:    userID,
		Name:  "Updated User",
		Email: "updated@example.com",
		Role:  "candidate",
	}

	mockService := &mockService{
		mockUpdateProfile: func(ctx context.Context, id uuid.UUID, updates UpdatesInput) (*User, error) {
			if id == userID {
				return updatedUser, nil
			}
			return nil, ErrInvalidEmail
		},
	}

	handler := NewGinHandler(mockService)

	// Create request body
	updateInput := UpdatesInput{
		Name: "Updated User",
	}
	body, _ := json.Marshal(updateInput)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/profile", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", userID.String())

	// Execute
	handler.UpdateProfileHandler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "profile updated successfully", response["message"])
}

func TestUpdateProfileHandler_NoFields(t *testing.T) {
	// Setup
	userID := uuid.New()
	mockService := &mockService{}
	handler := NewGinHandler(mockService)

	// Create request body with empty fields
	updateInput := UpdatesInput{}
	body, _ := json.Marshal(updateInput)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PUT", "/profile", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", userID.String())

	// Execute
	handler.UpdateProfileHandler(c)

	// Assert
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// Mock service for testing
type mockService struct {
	mockRegister       func(context.Context, RegisterInput) (*User, error)
	mockLogin          func(context.Context, LoginInput) (*LoginResponse, error)
	mockGetProfile     func(context.Context, uuid.UUID) (*User, error)
	mockUpdateProfile  func(context.Context, uuid.UUID, UpdatesInput) (*User, error)
}

func (m *mockService) Register(ctx context.Context, input RegisterInput) (*User, error) {
	if m.mockRegister != nil {
		return m.mockRegister(ctx, input)
	}
	return nil, nil
}

func (m *mockService) Login(ctx context.Context, input LoginInput) (*LoginResponse, error) {
	if m.mockLogin != nil {
		return m.mockLogin(ctx, input)
	}
	return nil, nil
}

func (m *mockService) GetProfile(ctx context.Context, userID uuid.UUID) (*User, error) {
	if m.mockGetProfile != nil {
		return m.mockGetProfile(ctx, userID)
	}
	return nil, nil
}

func (m *mockService) UpdateProfile(ctx context.Context, userID uuid.UUID, updates UpdatesInput) (*User, error) {
	if m.mockUpdateProfile != nil {
		return m.mockUpdateProfile(ctx, userID, updates)
	}
	return nil, nil
}
