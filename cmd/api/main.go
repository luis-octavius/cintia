package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luis-octavius/cintia/internal/middleware"
	"github.com/luis-octavius/cintia/internal/user"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("no .env file found")
	}

	secret := os.Getenv("SECRET")
	if secret == "" {
		secret = "dev-secret-key"
		log.Println("secret in .env not set, fallback...")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	repo := user.NewMockRepository()

	service := user.NewService(repo, secret)

	handler := user.NewGinHandler(service)

	api := r.Group("/api")
	{
		api.POST("/users/register", handler.RegisterHandler)
		api.POST("/users/login", handler.LoginHandler)

		protected := api.Group("/users")
		protected.Use(middleware.AuthMiddleware(secret))
		{
			protected.GET("/me", handler.GetProfileHandler)
			protected.PUT("/me", handler.UpdateProfileHandler)
		}
	}

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "cintia",
		})
	})

	log.Println("server starting on", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("failed to start server:", err)
	}
}
