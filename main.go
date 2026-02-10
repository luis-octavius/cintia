package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/luis-octavius/cintia/internal/database"
)

type apiConfig struct {
	Port    string
	Queries *database.Queries
}

func main() {
	godotenv.Load()
	mux := http.NewServeMux()

	port := os.Getenv("PORT")
	dbUrl := os.Getenv("DB_URL")

	db, err := sql.Open("postgres", dbUrl)
	if err != nil {
		log.Fatal("Error opening database: ", err)
	}

	queries := database.New(db)

	cfg := apiConfig{
		Port:    ":" + port,
		Queries: queries,
	}

	fmt.Println("Config: ", cfg)

	server := &http.Server{
		Handler: mux,
		Addr:    cfg.Port,
	}

	fmt.Printf("Server listening on port %v", port)

	mux.Handle("POST /api/users", cfg.handlerCreateUser())

	err = server.ListenAndServe()
	if err != nil {
		log.Fatalf("Error opening server: %v", err)
	}
}
