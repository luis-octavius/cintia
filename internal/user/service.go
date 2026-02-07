package user

/*
In this file we have the business rules
So the methods of service are intended to specify how we handle
constraint or rules to validate a user login, or a registering user, for example
*/

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/luis-octavius/cintia/internal/auth"
)

var (
	ErrEmailExists  = errors.New("email already exists")
	ErrInvalidName  = errors.New("name must not be empty")
	ErrInvalidEmail = errors.New("invalid email")
	ErrWeakPassword = errors.New("password must be at least 8 characters")
	InvalidName     = errors.New("invalid name")
)

var ctx = context.Background()

type Service interface {
	Register(ctx context.Context, input RegisterInput) (*User, error)
	Login(ctx context.Context, input LoginInput) (*User, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*User, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, updates UpdatesInput) (*User, error)
}

type service struct {
	repo      Repository
	jwtSecret string
	// in the future, it is possible to add logger, metrics, etc. here
}

func NewService(r Repository, s string) *service {
	return &service{
		repo:      r,
		jwtSecret: s,
	}
}

func (s *service) Register(ctx context.Context, input RegisterInput) (*User, error) {
	if input.Name == "" {
		return nil, ErrInvalidName
	}

	if input.Email == "" || !strings.Contains(input.Email, "@") {
		return nil, ErrInvalidEmail
	}

	if len(input.Password) < 8 {
		return nil, ErrWeakPassword
	}

	existing, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if existing != nil {
		return nil, ErrEmailExists
	}

	hash, err := auth.HashPassword(input.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user := &User{
		ID:           uuid.New(),
		Name:         input.Name,
		Email:        input.Email,
		PasswordHash: hash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	createdUser, err := s.repo.Create(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}
