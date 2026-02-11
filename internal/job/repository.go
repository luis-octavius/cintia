package job

import (
	"context"
	"errors"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("not found")

type Repository interface {
	Create(ctx context.Context, job *Job) (*Job, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Job, error)
	GetByLink(ctx context.Context, link string) (*Job, error)
	Search(ctx context.Context, filters JobFilters) ([]*Job, error)
	Count(ctx context.Context, filters JobFilters) (int, error)
	Update(ctx context.Context, job *Job) error
	MarkAsInactive(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}
