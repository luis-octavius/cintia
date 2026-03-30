package scraper

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/luis-octavius/cintia/internal/job"
)

type JobService interface {
	CreateJob(ctx context.Context, input job.CreateJobInput) (*job.Job, error)
}

type Source interface {
	Name() string
	FetchJobs(ctx context.Context) ([]job.CreateJobInput, error)
}

type RunStats struct {
	TotalFetched  int
	TotalCreated  int
	TotalSkipped  int
	SourceResults map[string]SourceStats
}

type SourceStats struct {
	Fetched int
	Created int
	Skipped int
	Failed  bool
	Error   string
}

type Scheduler struct {
	service  JobService
	sources  []Source
	logger   *log.Logger
	interval time.Duration
}

func NewScheduler(service JobService, sources []Source, interval time.Duration, logger *log.Logger) *Scheduler {
	if logger == nil {
		logger = log.Default()
	}

	if interval <= 0 {
		interval = 30 * time.Minute
	}

	return &Scheduler{
		service:  service,
		sources:  sources,
		logger:   logger,
		interval: interval,
	}
}

func (s *Scheduler) RunOnce(ctx context.Context) RunStats {
	stats := RunStats{
		SourceResults: make(map[string]SourceStats),
	}

	for _, source := range s.sources {
		name := source.Name()
		sourceStats := SourceStats{}

		jobs, err := source.FetchJobs(ctx)
		if err != nil {
			sourceStats.Failed = true
			sourceStats.Error = err.Error()
			stats.SourceResults[name] = sourceStats
			s.logger.Printf("scraper source %s failed: %v", name, err)
			continue
		}

		sourceStats.Fetched = len(jobs)
		stats.TotalFetched += len(jobs)

		for _, jobInput := range jobs {
			_, err := s.service.CreateJob(ctx, jobInput)
			if err != nil {
				if errors.Is(err, job.ErrDuplicateJob) {
					sourceStats.Skipped++
					stats.TotalSkipped++
					continue
				}

				sourceStats.Skipped++
				stats.TotalSkipped++
				s.logger.Printf("failed creating job from source %s (link: %s): %v", name, jobInput.Link, err)
				continue
			}

			sourceStats.Created++
			stats.TotalCreated++
		}

		stats.SourceResults[name] = sourceStats
	}

	return stats
}

func (s *Scheduler) Run(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	s.logger.Printf("scraper scheduler started with interval: %s", s.interval)

	stats := s.RunOnce(ctx)
	s.logger.Printf("initial scraper run: fetched=%d created=%d skipped=%d", stats.TotalFetched, stats.TotalCreated, stats.TotalSkipped)

	for {
		select {
		case <-ctx.Done():
			s.logger.Println("scraper scheduler stopped")
			return
		case <-ticker.C:
			stats = s.RunOnce(ctx)
			s.logger.Printf("scraper run finished: fetched=%d created=%d skipped=%d", stats.TotalFetched, stats.TotalCreated, stats.TotalSkipped)
		}
	}
}
