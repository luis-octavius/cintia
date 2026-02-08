package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luis-octavius/cintia/internal/user"
)

func main() {
	godotenv.Load()

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	repo := user.NewMockRepository()

	service := user.NewService(repo, "jwt-secret-key")

	handler := user.NewGinHandler(service)

	api := r.Group("/api")
	{
		users := api.Group("/users")
		{
			users.POST("/register", handler.RegisterHandler)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "cintia",
		})
	})

	log.Println("server starting on :8080")
	if err := r.Run(":8080"); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
