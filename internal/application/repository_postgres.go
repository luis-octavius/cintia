package application

import (
	"context"
	"database/sql"
	"errors"
	"strings"

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

func (r *PostgresRepository) Create(ctx context.Context, app *Application) (*Application, error) {
	dbApp, err := r.queries.CreateApplication(ctx, database.CreateApplicationParams{
		UserID: app.UserID,
		JobID:  app.JobID,
		Notes:  toNullString(app.Notes),
	})
	if err != nil {
		// Check for unique constraint violation (duplicate application)
		if strings.Contains(err.Error(), "unique_user_job") {
			return nil, ErrAlreadyExists
		}
		return nil, err
	}

	return dbAppToApp(&dbApp), nil
}

func (r *PostgresRepository) GetByID(ctx context.Context, id uuid.UUID) (*Application, error) {
	dbApp, err := r.queries.GetApplicationByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return dbAppToApp(&dbApp), nil
}

func (r *PostgresRepository) GetUserApplications(ctx context.Context, userID uuid.UUID) ([]*Application, error) {
	dbApps, err := r.queries.GetUserApplications(ctx, userID)
	if err != nil {
		return nil, err
	}

	apps := make([]*Application, len(dbApps))
	for i, dbApp := range dbApps {
		apps[i] = dbAppToApp(&dbApp)
	}

	return apps, nil
}

func (r *PostgresRepository) GetUserJobApplication(ctx context.Context, userID, jobID uuid.UUID) (*Application, error) {
	dbApp, err := r.queries.GetApplicationByUserAndJob(ctx, database.GetApplicationByUserAndJobParams{
		UserID: userID,
		JobID:  jobID,
	})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Not found is not an error for this method
		}
		return nil, err
	}

	return dbAppToApp(&dbApp), nil
}

func (r *PostgresRepository) GetJobApplications(ctx context.Context, jobID uuid.UUID) ([]*Application, error) {
	dbApps, err := r.queries.GetJobApplications(ctx, jobID)
	if err != nil {
		return nil, err
	}

	apps := make([]*Application, len(dbApps))
	for i, dbApp := range dbApps {
		apps[i] = dbAppToApp(&dbApp)
	}

	return apps, nil
}

func (r *PostgresRepository) Update(ctx context.Context, app *Application) error {
	params := database.UpdateApplicationParams{
		ID: app.ID,
	}

	if app.InterviewDate != nil {
		params.InterviewDate = sql.NullTime{Time: *app.InterviewDate, Valid: true}
	}
	if app.OfferDate != nil {
		params.OfferDate = sql.NullTime{Time: *app.OfferDate, Valid: true}
	}
	if app.Notes != "" {
		params.Notes = toNullString(app.Notes)
	}
	if app.SalaryOffer != "" {
		params.SalaryOffer = toNullString(app.SalaryOffer)
	}
	params.ReminderSent = sql.NullBool{Bool: app.ReminderSent, Valid: true}
	if app.FollowUpDate != nil {
		params.FollowUpDate = sql.NullTime{Time: *app.FollowUpDate, Valid: true}
	}

	dbApp, err := r.queries.UpdateApplication(ctx, params)
	if err != nil {
		return err
	}

	// Update the application object with returned values
	*app = *dbAppToApp(&dbApp)

	return nil
}

func (r *PostgresRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status ApplicationStatus) error {
	if !status.IsValid() {
		return ErrInvalidStatus
	}

	_, err := r.queries.UpdateApplicationStatus(ctx, database.UpdateApplicationStatusParams{
		ID:     id,
		Status: string(status),
	})
	return err
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteApplication(ctx, id)
}

// Helper functions to convert between domain and database models

func dbAppToApp(dbApp *database.Application) *Application {
	app := &Application{
		ID:           dbApp.ID,
		UserID:       dbApp.UserID,
		JobID:        dbApp.JobID,
		Status:       ApplicationStatus(dbApp.Status),
		AppliedAt:    dbApp.AppliedAt,
		UpdatedAt:    dbApp.UpdatedAt,
		Notes:        fromNullString(dbApp.Notes),
		SalaryOffer:  fromNullString(dbApp.SalaryOffer),
		ReminderSent: dbApp.ReminderSent,
	}

	if dbApp.InterviewDate.Valid {
		app.InterviewDate = &dbApp.InterviewDate.Time
	}
	if dbApp.OfferDate.Valid {
		app.OfferDate = &dbApp.OfferDate.Time
	}
	if dbApp.FollowUpDate.Valid {
		app.FollowUpDate = &dbApp.FollowUpDate.Time
	}

	return app
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
