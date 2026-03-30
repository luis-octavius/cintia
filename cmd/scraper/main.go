package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/luis-octavius/cintia/internal/database"
	"github.com/luis-octavius/cintia/internal/job"
	"github.com/luis-octavius/cintia/internal/scraper"
	"github.com/luis-octavius/cintia/internal/scraper/sources"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Println("no .env file found")
	}

	interval := parseDuration(getEnv("SCRAPER_INTERVAL", "30m"), 30*time.Minute)
	runOnce := strings.EqualFold(getEnv("SCRAPER_ONCE", "false"), "true")
	keywords := parseKeywords(getEnv("SCRAPER_KEYWORDS", "golang,backend"))
	location := getEnv("SCRAPER_LOCATION", "")

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
		log.Fatal("failed to connect to database: ", err)
	}
	defer db.Close()

	jobRepo := job.NewPostgresRepository(db)
	jobService := job.NewService(jobRepo)

	jobSources := []scraper.Source{
		sources.NewLinkedInSource("", keywords, location),
		sources.NewIndeedSource("", keywords, location),
	}

	scheduler := scraper.NewScheduler(jobService, jobSources, interval, log.Default())

	if runOnce {
		stats := scheduler.RunOnce(context.Background())
		log.Printf("scraper run once finished: fetched=%d created=%d skipped=%d", stats.TotalFetched, stats.TotalCreated, stats.TotalSkipped)
		return
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	scheduler.Run(ctx)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func parseKeywords(raw string) []string {
	parts := strings.Split(raw, ",")
	keywords := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		keywords = append(keywords, part)
	}

	if len(keywords) == 0 {
		return []string{"golang", "backend"}
	}

	return keywords
}

func parseDuration(raw string, fallback time.Duration) time.Duration {
	parsed, err := time.ParseDuration(raw)
	if err != nil {
		log.Printf("invalid SCRAPER_INTERVAL %q, fallback to %s", raw, fallback)
		return fallback
	}

	if parsed <= 0 {
		return fallback
	}

	return parsed
}
