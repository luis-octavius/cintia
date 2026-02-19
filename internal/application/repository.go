package application

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var (
	ErrNotFound          = errors.New("application not found")
	ErrAlreadyExists     = errors.New("application already exists")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInvalidTransition = errors.New("invalid status transition")
)

type Repository interface {
	Create(ctx context.Context, app *Application) (*Application, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Application, error)
	GetUserApplications(ctx context.Context, userID uuid.UUID) ([]*Application, error)
	GetUserJobApplication(ctx context.Context, userID, jobID uuid.UUID) (*Application, error)
	GetJobApplications(ctx context.Context, jobID uuid.UUID) ([]*Application, error)
	Update(ctx context.Context, app *Application) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status ApplicationStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}
