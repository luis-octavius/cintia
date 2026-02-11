package job

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

var (
	ErrJobNotFound    = errors.New("job not found")
	ErrDuplicateJob   = errors.New("job with this link already exists")
	ErrInvalidSource  = errors.New("invalid job source")
	ErrMissingTitle   = errors.New("job title is required")
	ErrMissingCompany = errors.New("company name is required")
	ErrFuturePostDate = errors.New("post date cannot be in the future")
)

type Service interface {
	CreateJob(ctx context.Context, input CreateJobInput) (*Job, error)
	GetJob(ctx context.Context, id uuid.UUID) (*Job, error)
	SearchJobs(ctx context.Context, filters JobFilters) (*JobsResponse, error)
	UpdateJob(ctx context.Context, id uuid.UUID, updates UpdateJobInput) (*Job, error)
	MarkJobAsInactive(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateJob(ctx context.Context, input CreateJobInput) (*Job, error) {
	if input.Title == "" {
		return nil, ErrMissingTitle
	}

	if input.Company == "" {
		return nil, ErrMissingCompany
	}

	// validate source
	validSources := map[string]bool{
		"linkedin": true,
		"indeed":   true,
		"manual":   true,
	}

	if !validSources[input.Source] {
		return nil, ErrInvalidSource
	}

	if input.Link != "" && !strings.HasPrefix(input.Link, "http") {
		return nil, errors.New("link must be a valid URL")
	}

	// search database with job link, if link exists returns an error
	existingJob, err := s.repo.GetByLink(ctx, input.Link)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("database error: %w", err)
		}
	}

	if existingJob != nil {
		return nil, ErrDuplicateJob
	}

	// validate if posted date is after right now
	postedDate := input.PostedDate
	if postedDate.IsZero() {
		postedDate = time.Now()
	} else if postedDate.After(time.Now()) {
		postedDate = time.Now()
	}

	job := &Job{
		ID:           uuid.New(),
		Title:        input.Title,
		Company:      input.Company,
		Location:     input.Location,
		Description:  input.Description,
		SalaryRange:  input.SalaryRange,
		Requirements: input.Requirements,
		Link:         input.Link,
		PostedDate:   postedDate,
		ScrapedAt:    time.Now(),
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdJob, err := s.repo.Create(ctx, job)
	if err != nil {
		return nil, fmt.Errorf("Error creating job: %w", err)
	}

	return createdJob, nil
}
