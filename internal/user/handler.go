package user

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler interface {
	RegisterHandler(c *gin.Context)
	LoginHandler(c *gin.Context)
	GetProfileHandler(c *gin.Context)
	UpdateProfileHandler(c *gin.Context)
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

func (h *GinHandler) LoginHandler(c *gin.Context) {
	var req LoginInput

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "invalid request format",
			"details": err.Error(),
		})
		return
	}

	log.Println("Email: ", req.Email)
	log.Println("Password: ", req.Password)

	response, err := h.service.Login(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user": gin.H{
			"id":    response.User.ID,
			"name":  response.User.Name,
			"email": response.User.Email,
			"role":  response.User.Role,
		},
		"token": response.Token,
	})
}

func (h *GinHandler) GetProfileHandler(c *gin.Context) {
	userIDStr, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	user, err := h.service.GetProfile(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": err.Error(),
		})

		return
	}

	c.JSON(http.StatusFound, gin.H{
		"user": gin.H{
			"id":         user.ID,
			"name":       user.Name,
			"email":      user.Email,
			"role":       user.Role,
			"created_at": user.CreatedAt,
			"updated_at": user.UpdatedAt,
		},
	})
}
