package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/luis-octavius/cintia/internal/application"
	"github.com/luis-octavius/cintia/internal/database"
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

	// Database connection
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", ""),
		DBName:   getEnv("DB_NAME", "cintia"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatal("failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("database connected successfully")

	r := gin.Default()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// Initialize repositories with real database
	repoUser := user.NewPostgresRepository(db)
	serviceUser := user.NewService(repoUser, secret)
	handlerUser := user.NewGinHandler(serviceUser)

	repoJob := job.NewPostgresRepository(db)
	serviceJob := job.NewService(repoJob)
	handlerJob := job.NewGinHandler(serviceJob)

	repoApp := application.NewPostgresRepository(db)
	serviceApp := application.NewService(repoApp, serviceJob, serviceUser)
	handlerApp := application.NewGinHandler(serviceApp)

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
			jobs.GET("/:jobID", handlerJob.GetJobHandler)
			jobs.GET("/:jobID/applications", handlerApp.GetJobApplicationsHandler)
			jobs.Use(middleware.AuthMiddleware(secret))
			{
				jobs.POST("/", handlerJob.CreateJobHandler)
				jobs.PATCH("/:jobID", handlerJob.ToggleJobStatusHandler)
			}

		}

		applications := api.Group("/applications")
		{
			applications.Use(middleware.AuthMiddleware(secret))
			{
				applications.POST("/", handlerApp.CreateApplicationHandler)
				applications.GET("/", handlerApp.GetUserApplicationsHandler)
				applications.GET("/:id", handlerApp.GetApplicationHandler)
				applications.PUT("/:id", handlerApp.UpdateApplicationHandler)
				applications.PATCH("/:id/status", handlerApp.UpdateStatusHandler)
				applications.DELETE("/:id", handlerApp.DeleteApplicationHandler)
			}
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

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
