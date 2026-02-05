package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	router := gin.Default()

	router.GET("/api/v1/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Hello!"})
	})

	router.Run(":8080")
}
