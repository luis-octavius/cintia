package job

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	CreateJobHandler(c *gin.Context)
	SearchJobsHandler(c *gin.Context)
	GetJobHandler(c *gin.Context)
	ToggleJobsStatusHandler(c *gin.Context)
}

type GinHandler struct {
	service Service
}

func NewGinHandler(service Service) *GinHandler {
	return &GinHandler{service: service}
}

func (h *GinHandler) CreateJobHandler(c *gin.Context) {
	var req CreateJobInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"details": err.Error(),
		})
		return
	}

	job, err := h.service.CreateJob(c.Request.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		switch {
		case errors.Is(err, ErrDuplicateJob):
			status = http.StatusConflict
		case errors.Is(err, ErrInvalidSource),
			errors.Is(err, ErrMissingTitle),
			errors.Is(err, ErrMissingCompany):
			status = http.StatusBadRequest
		default:
			status = http.StatusInternalServerError

		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "job created successfully",
		"job": gin.H{
			"id":          job.ID,
			"title":       job.Title,
			"company":     job.Company,
			"source":      job.Source,
			"posted_date": job.PostedDate,
			"link":        job.Link,
			"created_at":  job.CreatedAt,
		},
	})
}

func (h *GinHandler) SearchJobsHandler(c *gin.Context) {
}

func (h *GinHandler) GetJobHandler(c *gin.Context) {
}

func (h *GinHandler) ToggleJobsStatusHandler(c *gin.Context) {
}
