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
	appIDStr := c.Param("id")
	if appIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "application id is required",
		})
		return
	}

	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid application id format",
		})
		return
	}

	app, err := h.service.GetApplicationByID(c.Request.Context(), appID)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, ErrApplicationNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	// verify if user has permission
	userIDStr, _ := c.Get("userID")
	userID, _ := uuid.Parse(userIDStr.(string))

	if app.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you don't have permission to view this application",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"application": gin.H{
			"id":             app.ID,
			"job_id":         app.JobID,
			"status":         app.Status,
			"notes":          app.Notes,
			"applied_at":     app.AppliedAt,
			"updated_at":     app.UpdatedAt,
			"interview_date": app.InterviewDate,
			"offer_date":     app.OfferDate,
			"follow_up_date": app.FollowUpDate,
		},
	})
}

// 4. GET /api/jobs/:jobID/applications - sees who applied
func (h *GinHandler) GetJobApplicationsHandler(c *gin.Context) {
	// verify if user is admin
	userRole, exists := c.Get("userRole")
	if !exists || userRole != "admin" {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "admin access required",
		})
		return
	}

	jobIDStr := c.Param("jobID")
	if jobIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "job id is required",
		})
		return
	}

	jobID, err := uuid.Parse(jobIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid job id format",
		})
		return
	}

	jobApplications, err := h.service.GetJobApplications(c.Request.Context(), jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	response := make([]gin.H, 0)
	for _, app := range jobApplications {
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
	})
}

// 5. PUT /api/applications/:id - update application
func (h *GinHandler) UpdateApplicationHandler(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid job id format",
		})
		return
	}

	userIDVal, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "user not authenticated",
		})
		return
	}
	userIDStr, ok := userIDVal.(string)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid user id in context",
		})
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id format",
		})
		return
	}

	app, err := h.service.GetApplicationByID(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "application not found",
		})
		return
	}

	if app.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you can only update your own applications",
		})
		return
	}

	var req UpdateApplicationInput
	if err = c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	err = h.service.UpdateApplication(c.Request.Context(), appID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error updating application",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "application updated successfully",
	})
}

// 6; PATCH /api/applications/:id/status - update status
func (h *GinHandler) UpdateStatusHandler(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid application id format",
		})
		return
	}

	userIDVal, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id format",
		})
	}

	app, err := h.service.GetApplicationByID(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "application not found",
		})
		return
	}

	if app.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you can only update your own applications",
		})
		return
	}

	var input struct {
		Status ApplicationStatus `json:"status" binding:"required"`
	}

	if err = c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "status is required",
		})
		return
	}

	err = h.service.UpdateApplicationStatus(c.Request.Context(), app.ID, input.Status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "status updated successfully",
	})
}

// 7. DELETE /api/applications/:id - delete application
func (h *GinHandler) DeleteApplicationHandler(c *gin.Context) {
	appIDStr := c.Param("id")
	appID, err := uuid.Parse(appIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid application id format",
		})
		return
	}

	userIDVal, _ := c.Get("userID")
	userID, err := uuid.Parse(userIDVal.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid user id format",
		})
		return
	}

	app, err := h.service.GetApplicationByID(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "application not found",
		})
		return
	}

	if app.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "you can only delete your own applications",
		})
		return
	}

	err = h.service.Delete(c.Request.Context(), appID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "error deleting application",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "application deleted successfully",
	})
}
