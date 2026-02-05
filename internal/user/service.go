package user

import "github.com/luis-octavius/cintia/internal/auth"

var (
	EmailNotExists = errors.New("email not exists")
	InvalidEmail = errors.New("invalid email")
	InvalidPassword = errors.New("invalid password")
	InvalidName = errors.New("invalid name")
)

import (
	"context"

	"github.com/google/uuid"
)

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

func (s *service) Register(ctx context.Context, input RegisterInput) (*User, error) {
	switch {
		case input.Email == "":
			return nil, InvalidEmail
		case input.Password == "":
			return nil, InvalidPassword
		case input.Name == "":
			return nil, InvalidName
		}

		hash, err := auth.HashPassword(input.Password)
		if err != nil {
			return nil, err
		}

		user := &User{
			ID:        uuid.New(),
			Email:     input.Email,
			Password:  hash,
			Name:      input.Name,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		return s.repo.Create(ctx, user)
}
