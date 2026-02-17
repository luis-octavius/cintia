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

func (s *service) SearchJobs(ctx context.Context, filters JobFilters) (*JobsResponse, error) {
	if filters.Limit <= 0 {
		filters.Limit = 20
	}

	if filters.Page <= 0 {
		filters.Page = 1
	}

	// search jobs
	jobs, err := s.repo.Search(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search jobs: %w", err)
	}

	// count total (for pagination)
	total, err := s.repo.Count(ctx, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to count jobs: %w", err)
	}

	totalPages := (total + filters.Limit - 1) / filters.Limit

	response := &JobsResponse{
		Jobs:       jobs,
		Total:      total,
		Page:       filters.Page,
		TotalPages: totalPages,
		HasMore:    filters.Page < totalPages,
	}

	return response, nil
}

func (s *service) GetJob(ctx context.Context, id uuid.UUID) (*Job, error) {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrJobNotFound
		}
		return nil, fmt.Errorf("failed to get job: %w", err)
	}

	return job, nil
}

func (s *service) UpdateJob(ctx context.Context, id uuid.UUID, updates UpdateJobInput) (*Job, error) {
	if updates.Link != "" && !strings.HasPrefix(updates.Link, "http") {
		return nil, fmt.Errorf("link must be a valid url")
	}

	if updates.Source != "" {
		validSources := map[string]bool{
			"linkedin": true,
			"indeed":   true,
			"manual":   true,
		}
		if !validSources[updates.Source] {
			return nil, ErrInvalidSource
		}
	}

	if updates.PostedDate != nil && updates.PostedDate.After(time.Now()) {
		return nil, fmt.Errorf("posted date cannot be in the future")
	}

	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to find job: %w", err)
	}

	// if none were updated, stays false
	updated := false

	// update only provided fields
	if updates.Title != "" {
		job.Title = updates.Title
		updated = true
	}

	if updates.Company != "" {
		job.Company = updates.Company
		updated = true
	}

	if updates.Location != "" {
		job.Location = updates.Company
		updated = true
	}

	if updates.Description != "" {
		job.Description = updates.Description
		updated = true
	}

	if updates.SalaryRange != "" {
		job.SalaryRange = updates.SalaryRange
		updated = true
	}

	if updates.Requirements != "" {
		job.Requirements = updates.Requirements
		updated = true
	}

	if updates.Source != "" {
		job.Source = updates.Source
		updated = true
	}

	if updates.Link != "" && updates.Link != job.Link {
		existing, err := s.repo.GetByLink(ctx, updates.Link)
		if err != nil && !errors.Is(err, ErrNotFound) {
			return nil, fmt.Errorf("database error: %w", err)
		}

		if existing != nil && existing.ID != job.ID {
			return nil, ErrDuplicateJob
		}

		job.Link = updates.Link
		updated = true
	}

	if updates.IsActive != nil {
		job.IsActive = *updates.IsActive
		updated = true
	}

	if updates.PostedDate != nil {
		job.PostedDate = *updates.PostedDate
		updated = true
	}

	if !updated {
		return nil, errors.New("no fields to update")
	}

	job.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to update job: %w", err)
	}

	return job, nil
}

func (s *service) MarkJobAsInactive(ctx context.Context, id uuid.UUID) error {
	job, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return fmt.Errorf("database error: %w", err)
		}
		return err
	}

	job.UpdatedAt = time.Now()

	err = s.repo.MarkJobAsInactive(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to mark job as inactive: %w", err)
	}

	return nil
}
