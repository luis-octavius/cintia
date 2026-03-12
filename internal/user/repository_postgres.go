package user

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

func (r *PostgresRepository) Create(ctx context.Context, user *User) (*User, error) {
	dbUser, err := r.queries.CreateUser(ctx, database.CreateUserParams{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
	})
	if err != nil {
		return nil, err
	}

	return &User{
		ID:           dbUser.ID,
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}, nil
}

func (r *PostgresRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &User{
		ID:           dbUser.ID,
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}, nil
}

func (r *PostgresRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	dbUser, err := r.queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // User not found
		}
		return nil, err
	}

	return &User{
		ID:           dbUser.ID,
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Role:         dbUser.Role,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
	}, nil
}

func (r *PostgresRepository) Update(ctx context.Context, user *User) error {
	// Build nullable parameters for partial updates
	params := database.UpdateUserParams{
		ID: user.ID,
	}

	if user.Name != "" {
		params.Name = sql.NullString{String: user.Name, Valid: true}
	}
	if user.Email != "" {
		params.Email = sql.NullString{String: user.Email, Valid: true}
	}
	if user.PasswordHash != "" {
		params.PasswordHash = sql.NullString{String: user.PasswordHash, Valid: true}
	}

	dbUser, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return err
	}

	// Update the user object with returned values
	user.Name = dbUser.Name
	user.Email = dbUser.Email
	user.PasswordHash = dbUser.PasswordHash
	user.UpdatedAt = dbUser.UpdatedAt

	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id uuid.UUID) error {
	return r.queries.DeleteUser(ctx, id)
}
