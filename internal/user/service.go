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
	ErrEmailExists     = errors.New("email already exists")
	ErrInvalidName     = errors.New("name must not be empty")
	ErrInvalidEmail    = errors.New("invalid email")
	ErrWeakPassword    = errors.New("password must be at least 8 characters")
	ErrInvalidPassword = errors.New("password invalid")
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

func (s *service) Login(ctx context.Context, input LoginInput) (*User, error) {
	if input.Email == "" || !strings.Contains(input.Email, "@") {
		return nil, ErrInvalidEmail
	}

	if len(input.Password) < 8 {
		return nil, ErrWeakPassword
	}

	user, err := s.repo.FindByEmail(ctx, input.Email)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	checkHash, err := auth.CheckPasswordHash(input.Password, user.PasswordHash)
	if err != nil {
		return nil, fmt.Errorf("CheckPasswordHash error: %w", err)
	}

	if !checkHash {
		return nil, fmt.Errorf("")
	}

	return user, nil
}

func (s *service) GetProfile(ctx context.Context, userID uuid.UUID) (*User, error) {
	err := uuid.Validate(userID.String())
	if err != nil {
		return nil, fmt.Errorf("Invalid UUID: %w", err)
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return user, nil
}

func (s *service) UpdateProfile(ctx context.Context, userID uuid.UUID, updates UpdatesInput) (*User, error) {
	err := uuid.Validate(userID.String())
	if err != nil {
		return nil, fmt.Errorf("Invalid UUID: %w", err)
	}

	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("database error: %w", err)
	}

	if updates.Name != "" {
		user.Name = updates.Name
	}

	if len(updates.Password) < 8 {
		return nil, ErrInvalidPassword
	}

	if updates.Email == "" || !strings.Contains(updates.Email, "@") {
		return nil, ErrInvalidEmail
	}

	user.Email = updates.Email

	hash, err := auth.HashPassword(updates.Password)
	if err != nil {
		return nil, fmt.Errorf("HashPassword error: %w", err)
	}

	user.PasswordHash = hash

	err = s.repo.Update(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("error updating user: %w", err)
	}

	return user, nil
}
