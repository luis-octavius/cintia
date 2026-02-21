package application

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler interface {
	CreateApplicationHandler(c *gin.Context)
	GetUserApplicationsHandler(c *gin.Context)
	GetApplicationHandler(c *gin.Context)
	GetJobApplicationsHandler(c *gin.Context)
	UpdateApplicationHandler(c *gin.Context)
	UpdateStatusHandler(c *gin.Context)
	DeleteApplicationHandler(c *gin.Context)
}

type GinHandler struct {
	service Service
}

func NewGinHandler(service Service) *GinHandler {
	return &GinHandler{service: service}
}

// 1. POST /api/applications - create new application
func (h *GinHandler) CreateApplicationHandler(c *gin.Context) {
	userIDStr, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id format",
		})
		return
	}

	var req CreateApplicationInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	application, err := h.service.CreateApplication(c.Request.Context(), userID, req)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrAlreadyApplied) {
			status = http.StatusConflict
		}

		if errors.Is(err, ErrJobNotFound) {
			status = http.StatusNotFound
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "application created successfully",
		"application": gin.H{
			"id":         application.ID,
			"job_id":     application.JobID,
			"status":     application.Status,
			"notes":      application.Notes,
			"applied_at": application.AppliedAt,
			"updated_at": application.UpdatedAt,
		},
	})
}

// 2. GET /api/applications - list user applications
func (h *GinHandler) GetUserApplicationsHandler(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id format",
		})
		return
	}

	applications, err := h.service.GetUserApplications(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := make([]gin.H, 0)

	for _, app := range applications {
		response = append(response, gin.H{
			"id":             app.ID,
			"job_id":         app.JobID,
			"status":         app.Status,
			"notes":          app.Notes,
			"applied_at":     app.AppliedAt,
			"updated_at":     app.UpdatedAt,
			"interview_date": app.InterviewDate,
			"offer_date":     app.OfferDate,
			"follow_up_date": app.FollowUpDate,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"applications": response,
		"total":        len(response),
	})
}

// 3. GET /api/applications/:id - details of an application
func (h *GinHandler) GetApplicationHandler(c *gin.Context) {
}

// 4. GET /api/jobs/:jobID/applications - sees who applied
func (h *GinHandler) GetJobApplicationsHandler(c *gin.Context) {
}

// 5. PUT /api/applications/:id - update application
func (h *GinHandler) UpdateApplicationHandler(c *gin.Context) {
}

// 6; PATCH /api/applications/:id/status - update status
func (h *GinHandler) UpdateStatusHandler(c *gin.Context) {
}

// 7. DELETE /api/applications/:id - delete application
func (h *GinHandler) DeleteApplicationHandler(c *gin.Context) {
}
