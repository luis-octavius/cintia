package application

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateApplicationHandler_Success(t *testing.T) {
	// Setup
	userID := uuid.New()
	jobID := uuid.New()
	appID := uuid.New()

	mockService := &mockApplicationService{
		mockCreateApplication: func(ctx context.Context, uID uuid.UUID, input CreateApplicationInput) (*Application, error) {
			return &Application{
				ID:        appID,
				UserID:    uID,
				JobID:     input.JobID,
				Status:    StatusApplied,
				AppliedAt: time.Now(),
				UpdatedAt: time.Now(),
				Notes:     input.Notes,
			}, nil
		},
	}

	handler := NewGinHandler(mockService)

	// Create request
	reqBody := CreateApplicationInput{
		JobID: jobID,
		Notes: "Interested in this role",
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/applications", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", userID.String())

	// Execute
	handler.CreateApplicationHandler(c)

	// Assert
	assert.Equal(t, http.StatusCreated, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "application created successfully", response["message"])
}

func TestCreateApplicationHandler_AlreadyApplied(t *testing.T) {
	// Setup
	userID := uuid.New()
	jobID := uuid.New()

	mockService := &mockApplicationService{
		mockCreateApplication: func(ctx context.Context, uID uuid.UUID, input CreateApplicationInput) (*Application, error) {
			return nil, ErrAlreadyApplied
		},
	}

	handler := NewGinHandler(mockService)

	reqBody := CreateApplicationInput{
		JobID: jobID,
	}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/applications", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Set("userID", userID.String())

	// Execute
	handler.CreateApplicationHandler(c)

	// Assert - should return 409 Conflict
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestGetApplicationHandler_Authorized(t *testing.T) {
	// Setup
	userID := uuid.New()
	appID := uuid.New()

	application := &Application{
		ID:        appID,
		UserID:    userID,
		JobID:     uuid.New(),
		Status:    StatusApplied,
		AppliedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService := &mockApplicationService{
		mockGetApplicationByID: func(ctx context.Context, id uuid.UUID) (*Application, error) {
			if id == appID {
				return application, nil
			}
			return nil, ErrApplicationNotFound
		},
	}

	handler := NewGinHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/applications/"+appID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: appID.String()}}
	c.Set("userID", userID.String())

	// Execute
	handler.GetApplicationHandler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestGetApplicationHandler_Unauthorized(t *testing.T) {
	// Setup
	userID := uuid.New()
	otherUserID := uuid.New()
	appID := uuid.New()

	application := &Application{
		ID:        appID,
		UserID:    otherUserID, // Different user
		JobID:     uuid.New(),
		Status:    StatusApplied,
		AppliedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService := &mockApplicationService{
		mockGetApplicationByID: func(ctx context.Context, id uuid.UUID) (*Application, error) {
			if id == appID {
				return application, nil
			}
			return nil, ErrApplicationNotFound
		},
	}

	handler := NewGinHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/applications/"+appID.String(), nil)
	c.Params = gin.Params{{Key: "id", Value: appID.String()}}
	c.Set("userID", userID.String())

	// Execute
	handler.GetApplicationHandler(c)

	// Assert - should return 403 Forbidden
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateStatusHandler_ValidTransition(t *testing.T) {
	// Setup
	userID := uuid.New()
	appID := uuid.New()

	application := &Application{
		ID:        appID,
		UserID:    userID,
		JobID:     uuid.New(),
		Status:    StatusApplied,
		AppliedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService := &mockApplicationService{
		mockGetApplicationByID: func(ctx context.Context, id uuid.UUID) (*Application, error) {
			return application, nil
		},
		mockUpdateApplicationStatus: func(ctx context.Context, id uuid.UUID, status ApplicationStatus) error {
			return nil
		},
	}

	handler := NewGinHandler(mockService)

	reqBody := map[string]string{"status": string(StatusInterviewing)}
	body, _ := json.Marshal(reqBody)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PATCH", "/applications/"+appID.String()+"/status", bytes.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "id", Value: appID.String()}}
	c.Set("userID", userID.String())

	// Execute
	handler.UpdateStatusHandler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "status updated successfully", response["message"])
}

// Mock service for testing
type mockApplicationService struct {
	mockCreateApplication          func(context.Context, uuid.UUID, CreateApplicationInput) (*Application, error)
	mockGetApplicationByID          func(context.Context, uuid.UUID) (*Application, error)
	mockGetUserApplications         func(context.Context, uuid.UUID) ([]*Application, error)
	mockGetJobApplications          func(context.Context, uuid.UUID) ([]*Application, error)
	mockUpdateApplication           func(context.Context, uuid.UUID, UpdateApplicationInput) error
	mockUpdateApplicationStatus     func(context.Context, uuid.UUID, ApplicationStatus) error
	mockDelete                      func(context.Context, uuid.UUID) error
}

func (m *mockApplicationService) CreateApplication(ctx context.Context, userID uuid.UUID, input CreateApplicationInput) (*Application, error) {
	if m.mockCreateApplication != nil {
		return m.mockCreateApplication(ctx, userID, input)
	}
	return nil, nil
}

func (m *mockApplicationService) GetApplicationByID(ctx context.Context, id uuid.UUID) (*Application, error) {
	if m.mockGetApplicationByID != nil {
		return m.mockGetApplicationByID(ctx, id)
	}
	return nil, nil
}

func (m *mockApplicationService) GetUserApplications(ctx context.Context, userID uuid.UUID) ([]*Application, error) {
	if m.mockGetUserApplications != nil {
		return m.mockGetUserApplications(ctx, userID)
	}
	return nil, nil
}

func (m *mockApplicationService) GetJobApplications(ctx context.Context, jobID uuid.UUID) ([]*Application, error) {
	if m.mockGetJobApplications != nil {
		return m.mockGetJobApplications(ctx, jobID)
	}
	return nil, nil
}

func (m *mockApplicationService) UpdateApplication(ctx context.Context, id uuid.UUID, updates UpdateApplicationInput) error {
	if m.mockUpdateApplication != nil {
		return m.mockUpdateApplication(ctx, id, updates)
	}
	return nil
}

func (m *mockApplicationService) UpdateApplicationStatus(ctx context.Context, id uuid.UUID, status ApplicationStatus) error {
	if m.mockUpdateApplicationStatus != nil {
		return m.mockUpdateApplicationStatus(ctx, id, status)
	}
	return nil
}

func (m *mockApplicationService) Delete(ctx context.Context, id uuid.UUID) error {
	if m.mockDelete != nil {
		return m.mockDelete(ctx, id)
	}
	return nil
}
