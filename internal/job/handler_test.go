package job

import (
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

func TestGetJobHandler_ValidID(t *testing.T) {
	// Setup
	jobID := uuid.New()
	mockJob := &Job{
		ID:        jobID,
		Title:     "Software Engineer",
		Company:   "Tech Corp",
		Location:  "San Francisco",
		IsActive:  true,
		CreatedAt: time.Now(),
	}

	mockService := &mockJobService{
		mockGetJob: func(ctx context.Context, id uuid.UUID) (*Job, error) {
			if id == jobID {
				return mockJob, nil
			}
			return nil, ErrJobNotFound
		},
	}

	handler := NewGinHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/jobs/"+jobID.String(), nil)
	c.Params = gin.Params{{Key: "jobID", Value: jobID.String()}}

	// Execute
	handler.GetJobHandler(c)

	// Assert
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.NotNil(t, response["job"])
}

func TestGetJobHandler_NotFound(t *testing.T) {
	// Setup
	jobID := uuid.New()
	mockService := &mockJobService{
		mockGetJob: func(ctx context.Context, id uuid.UUID) (*Job, error) {
			return nil, ErrJobNotFound
		},
	}

	handler := NewGinHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/jobs/"+jobID.String(), nil)
	c.Params = gin.Params{{Key: "jobID", Value: jobID.String()}}

	// Execute
	handler.GetJobHandler(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestToggleJobStatusHandler_Success(t *testing.T) {
	// Setup
	jobID := uuid.New()
	mockService := &mockJobService{
		mockMarkJobAsInactive: func(ctx context.Context, id uuid.UUID) error {
			if id == jobID {
				return nil
			}
			return ErrNotFound
		},
	}

	handler := NewGinHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PATCH", "/jobs/"+jobID.String()+"/status", nil)
	c.Params = gin.Params{{Key: "jobID", Value: jobID.String()}}

	// Execute
	handler.ToggleJobStatusHandler(c)

	// Assert - should return 200 OK, not 202 Accepted
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "job marked as inactive successfully", response["message"])
}

func TestToggleJobStatusHandler_NotFound(t *testing.T) {
	// Setup
	jobID := uuid.New()
	mockService := &mockJobService{
		mockMarkJobAsInactive: func(ctx context.Context, id uuid.UUID) error {
			return ErrNotFound
		},
	}

	handler := NewGinHandler(mockService)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("PATCH", "/jobs/"+jobID.String()+"/status", nil)
	c.Params = gin.Params{{Key: "jobID", Value: jobID.String()}}

	// Execute
	handler.ToggleJobStatusHandler(c)

	// Assert
	assert.Equal(t, http.StatusNotFound, w.Code)
}

// Mock service for testing
type mockJobService struct {
	mockCreateJob          func(context.Context, CreateJobInput) (*Job, error)
	mockSearchJobs         func(context.Context, JobFilters) (*JobsResponse, error)
	mockGetJob             func(context.Context, uuid.UUID) (*Job, error)
	mockUpdateJob          func(context.Context, uuid.UUID, UpdateJobInput) (*Job, error)
	mockMarkJobAsInactive  func(context.Context, uuid.UUID) error
}

func (m *mockJobService) CreateJob(ctx context.Context, input CreateJobInput) (*Job, error) {
	if m.mockCreateJob != nil {
		return m.mockCreateJob(ctx, input)
	}
	return nil, nil
}

func (m *mockJobService) SearchJobs(ctx context.Context, filters JobFilters) (*JobsResponse, error) {
	if m.mockSearchJobs != nil {
		return m.mockSearchJobs(ctx, filters)
	}
	return nil, nil
}

func (m *mockJobService) GetJob(ctx context.Context, id uuid.UUID) (*Job, error) {
	if m.mockGetJob != nil {
		return m.mockGetJob(ctx, id)
	}
	return nil, nil
}

func (m *mockJobService) UpdateJob(ctx context.Context, id uuid.UUID, updates UpdateJobInput) (*Job, error) {
	if m.mockUpdateJob != nil {
		return m.mockUpdateJob(ctx, id, updates)
	}
	return nil, nil
}

func (m *mockJobService) MarkJobAsInactive(ctx context.Context, id uuid.UUID) error {
	if m.mockMarkJobAsInactive != nil {
		return m.mockMarkJobAsInactive(ctx, id)
	}
	return nil
}
