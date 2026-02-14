package job

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler interface {
	CreateJobHandler(c *gin.Context)
	SearchJobsHandler(c *gin.Context)
	GetJobHandler(c *gin.Context)
	ToggleJobStatusHandler(c *gin.Context)
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
	filters := JobFilters{}

	if title := c.Query("title"); title != "" {
		filters.Title = title
	}

	if company := c.Query("company"); company != "" {
		filters.Company = company
	}

	if location := c.Query("location"); location != "" {
		filters.Location = location
	}

	if source := c.Query("source"); source != "" {
		filters.Source = source
	}

	if isActive := c.Query("is_active"); isActive != "" {
		active := isActive == "true"
		filters.IsActive = &active
	}

	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		page = 1
	}
	filters.Page = page

	limit, err := strconv.Atoi(c.DefaultQuery("limit", "20"))
	if err != nil || limit < 1 || limit > 100 {
		limit = 20
	}
	filters.Limit = limit

	response, err := h.service.SearchJobs(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to search jobs",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

func (h *GinHandler) GetJobHandler(c *gin.Context) {
	jobID := c.Param("jobID")

	parsedID, err := uuid.Parse(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	job, err := h.service.GetJob(c.Request.Context(), parsedID)
	if err != nil {
		status := http.StatusBadRequest
		if err == ErrJobNotFound {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{
			"error":   err.Error(),
			"message": "failed to get job",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"job": gin.H{	
			"id":           job.ID,
			"title":        job.Title,
			"company":      job.Company,
			"source":       job.Source,
			"posted_date":  job.PostedDate,
			"link":         job.Link,
			"created_at":   job.CreatedAt,
			"location": 	  job.Location, 
			"description":  job.Description, 
			"salary_range": job.SalaryRange,
			"is_active":    job.IsActive,
		}
	})
}

func (h *GinHandler) ToggleJobStatusHandler(c *gin.Context) {
	jobID := c.Param("jobID")

	parsedID, err := uuid.Parse(jobID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err = h.service.MarkJobAsInactive(c.Request.Context(), parsedID)
	if err != nil {
		status := http.StatusBadRequest 
		if error.Is(err, ErrNotFound) {
			status = http.StatusNotFound 
		}
		c.JSON(status, gin.H{
			"error": err.Error(),
			"message": "failed to mark job as inactive",
		})
		return 
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "job marked as inactive successfully",
	})
}
