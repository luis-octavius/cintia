package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler interface {
	RegisterHandler(w http.ResponseWriter, r *http.Request)
	LoginHandler(w http.ResponseWriter, r *http.Request)
	GetProfileHandler(w http.ResponseWriter, r *http.Request)
	UpdateProfileHandler(w http.ResponseWriter, r *http.Request)
}

type GinHandler struct {
	service Service
}

func NewGinHandler(s Service) *GinHandler {
	return &GinHandler{
		service: s,
	}
}

func (h *GinHandler) RegisterHandler(c *gin.Context) {
	// verify HTTP method
	var req RegisterInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"details": err.Error(),
		})
		return
	}

	user, err := h.service.Register(c.Request.Context(), req)
	if err != nil {
		status := http.StatusBadRequest
		if err == ErrEmailExists {
			status = http.StatusConflict
		}

		c.JSON(status, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         user.ID,
		"name":       user.Name,
		"email":      user.Email,
		"created_at": user.CreatedAt,
	})
}
