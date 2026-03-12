package job

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/luis-octavius/cintia/internal/database"
)

type PostgresRepository struct {
	queries *database.Queries
}

func NewPostgresRepository(db *sql.DB) Repository {
	return &PostgresRepository{
		queries: database.New(db),
	}
}

func (r *PostgresRepository) Create(ctx context.Context, job *Job) (*Job, error) {
	dbJob, err := r.queries.CreateJob(ctx, database.CreateJobParams{
		Title:        job.Title,
		Company:      job.Company,
		Location:     job.Location,
		Description:  job.Description,
		SalaryRange:  toNullString(job.SalaryRange),
		Requirements: toNullString(job.Requirements),
		Source:       job.Source,
		Link:         job.Link,
		PostedDate:   job.PostedDate,
	})
	if err != nil {
		return nil, err
	}

	return dbJobToJob(&dbJob), nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Job, error) {
	dbJob, err := r.queries.GetJobByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return dbJobToJob(&dbJob), nil
}

func (r *PostgresRepository) GetByLink(ctx context.Context, link string) (*Job, error) {
	dbJob, err := r.queries.GetJobByLink(ctx, link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found is not an error for this method
		}
		return nil, err
	}

	return dbJobToJob(&dbJob), nil
}

func (r *PostgresRepository) Search(ctx context.Context, filters JobFilters) ([]*Job, error) {
	// Calculate offset for pagination
	limit := filters.Limit
	if limit == 0 {
		limit = 20 // default
	}
	offset := (filters.Page - 1) * limit
	if offset < 0 {
		offset = 0
	}

	dbJobs, err := r.queries.ListJobs(ctx, database.ListJobsParams{
		Limit:    int32(limit),
		Offset:   int32(offset),
		Title:    toNullString(filters.Title),
		Company:  toNullString(filters.Company),
		Location: toNullString(filters.Location),
		Source:   toNullString(filters.Source),
		IsActive: toNullBool(filters.IsActive),
	})
	if err != nil {
		return nil, err
	}

	jobs := make([]*Job, len(dbJobs))
	for i, dbJob := range dbJobs {
		jobs[i] = dbJobToJob(&dbJob)
	}

	return jobs, nil
}

func (r *PostgresRepository) Count(ctx context.Context, filters JobFilters) (int, error) {
	count, err := r.queries.CountJobs(ctx, database.CountJobsParams{
		Title:    toNullString(filters.Title),
		Company:  toNullString(filters.Company),
		Location: toNullString(filters.Location),
		Source:   toNullString(filters.Source),
		IsActive: toNullBool(filters.IsActive),
	})
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

func (r *PostgresRepository) Update(ctx context.Context, job *Job) error {
	params := database.UpdateJobParams{
		ID: job.ID,
	}

	if job.Title != "" {
		params.Title = toNullString(job.Title)
	}
	if job.Company != "" {
		params.Company = toNullString(job.Company)
	}
	if job.Location != "" {
		params.Location = toNullString(job.Location)
	}
	if job.Description != "" {
		params.Description = toNullString(job.Description)
	}
	if job.SalaryRange != "" {
		params.SalaryRange = toNullString(job.SalaryRange)
	}
	if job.Requirements != "" {
		params.Requirements = toNullString(job.Requirements)
	}
	if job.Source != "" {
		params.Source = toNullString(job.Source)
	}
	if job.Link != "" {
		params.Link = toNullString(job.Link)
	}
	params.IsActive = sql.NullBool{Bool: job.IsActive, Valid: true}
	if !job.PostedDate.IsZero() {
		params.PostedDate = sql.NullTime{Time: job.PostedDate, Valid: true}
	}

	dbJob, err := r.queries.UpdateJob(ctx, params)
	if err != nil {
		return err
	}

	// Update the job object with returned values
	*job = *dbJobToJob(&dbJob)

	return nil
}

func (r *PostgresRepository) MarkJobAsInactive(ctx context.Context, id uuid.UUID) error {
	_, err := r.queries.UpdateJob(ctx, database.UpdateJobParams{
		ID:       id,
		IsActive: sql.NullBool{Bool: false, Valid: true},
	})
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteJob(ctx, id)
}

// Helper functions to convert between domain and database models

func dbJobToJob(dbJob *database.Job) *Job {
	return &Job{
		ID:           dbJob.ID,
		Title:        dbJob.Title,
		Company:      dbJob.Company,
		Location:     dbJob.Location,
		Description:  dbJob.Description,
		SalaryRange:  fromNullString(dbJob.SalaryRange),
		Requirements: fromNullString(dbJob.Requirements),
		Source:       dbJob.Source,
		Link:         dbJob.Link,
		PostedDate:   dbJob.PostedDate,
		ScrapedAt:    dbJob.ScrapedAt,
		IsActive:     dbJob.IsActive,
		CreatedAt:    dbJob.CreatedAt,
		UpdatedAt:    dbJob.UpdatedAt,
	}
}

func toNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

func fromNullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}

func toNullBool(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}
