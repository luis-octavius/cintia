package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luis-octavius/cintia/internal/middleware"
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
		api.POST("/users/register", handler.RegisterHandler)
		api.POST("/users/login", handler.LoginHandler)

		protected := api.Group("/users")
		protected.Use(middleware.AuthMiddleware("jwt-secret-key"))
		{
			protected.GET("/me", handler.GetProfileHandler)
			protected.GET("/me", handler.UpdateProfileHandler)
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
