package scraper

import (
	"context"
	"errors"
	"log"
	"testing"
	"time"

	"github.com/luis-octavius/cintia/internal/job"
)

func TestRunOnce_AggregatesStats(t *testing.T) {
	service := &mockJobService{}
	sources := []Source{
		mockSource{name: "linkedin", jobs: []job.CreateJobInput{{Title: "Go Dev", Company: "A", Source: "linkedin", Link: "https://x/1", PostedDate: time.Now()}}},
		mockSource{name: "indeed", jobs: []job.CreateJobInput{{Title: "Backend Dev", Company: "B", Source: "indeed", Link: "https://x/2", PostedDate: time.Now()}}},
	}

	scheduler := NewScheduler(service, sources, time.Minute, log.Default())
	stats := scheduler.RunOnce(context.Background())

	if stats.TotalFetched != 2 {
		t.Fatalf("expected 2 fetched jobs, got %d", stats.TotalFetched)
	}

	if stats.TotalCreated != 2 {
		t.Fatalf("expected 2 created jobs, got %d", stats.TotalCreated)
	}

	if stats.TotalSkipped != 0 {
		t.Fatalf("expected 0 skipped jobs, got %d", stats.TotalSkipped)
	}
}

func TestRunOnce_SkipsDuplicates(t *testing.T) {
	service := &mockJobService{createErr: job.ErrDuplicateJob}
	sources := []Source{
		mockSource{name: "linkedin", jobs: []job.CreateJobInput{{Title: "Go Dev", Company: "A", Source: "linkedin", Link: "https://x/1", PostedDate: time.Now()}}},
	}

	scheduler := NewScheduler(service, sources, time.Minute, log.Default())
	stats := scheduler.RunOnce(context.Background())

	if stats.TotalFetched != 1 {
		t.Fatalf("expected 1 fetched job, got %d", stats.TotalFetched)
	}

	if stats.TotalCreated != 0 {
		t.Fatalf("expected 0 created jobs, got %d", stats.TotalCreated)
	}

	if stats.TotalSkipped != 1 {
		t.Fatalf("expected 1 skipped job, got %d", stats.TotalSkipped)
	}
}

func TestRunOnce_RecordsSourceErrors(t *testing.T) {
	service := &mockJobService{}
	sources := []Source{
		mockSource{name: "linkedin", err: errors.New("source failed")},
	}

	scheduler := NewScheduler(service, sources, time.Minute, log.Default())
	stats := scheduler.RunOnce(context.Background())

	result, ok := stats.SourceResults["linkedin"]
	if !ok {
		t.Fatal("expected source result for linkedin")
	}

	if !result.Failed {
		t.Fatal("expected source to be marked as failed")
	}
}

type mockJobService struct {
	createErr error
}

func (m *mockJobService) CreateJob(ctx context.Context, input job.CreateJobInput) (*job.Job, error) {
	if m.createErr != nil {
		return nil, m.createErr
	}

	return &job.Job{Title: input.Title, Company: input.Company}, nil
}

type mockSource struct {
	name string
	jobs []job.CreateJobInput
	err  error
}

func (m mockSource) Name() string {
	return m.name
}

func (m mockSource) FetchJobs(ctx context.Context) ([]job.CreateJobInput, error) {
	_ = ctx
	if m.err != nil {
		return nil, m.err
	}

	return m.jobs, nil
}
