package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luis-octavius/cintia/internal/job"
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

	repoUser := user.NewMockRepository()

	serviceUser := user.NewService(repoUser, secret)

	handlerUser := user.NewGinHandler(serviceUser)

	repoJob := job.NewMockRepository()

	serviceJob := job.NewService(repoJob)

	handlerJob := job.NewGinHandler(serviceJob)

	api := r.Group("/api")
	{
		// users routes
		users := api.Group("/users")
		{
			users.POST("/register", handlerUser.RegisterHandler)
			users.POST("/login", handlerUser.LoginHandler)
			users.Use(middleware.AuthMiddleware(secret))
			{
				users.GET("/me", handlerUser.GetProfileHandler)
				users.PUT("/me", handlerUser.UpdateProfileHandler)
			}
		}

		// jobs routes
		jobs := api.Group("/jobs")
		{
			jobs.GET("/", handlerJob.SearchJobsHandler)
			jobs.POST("/", handlerJob.CreateJobHandler)
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
