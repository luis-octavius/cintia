package user

import (
	"context"

	"github.com/google/uuid"
)

type mockRepository struct {
	users map[string]*User
}

func (m *mockRepository) Create(ctx context.Context, user *User) (*User, error) {
	user.ID = uuid.New()
	m.users[user.Email] = user
	return user, nil
}
