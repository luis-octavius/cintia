package user

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

var (
	ErrInvalidUser = errors.New("user does not exist")
)

type mockRepository struct {
	mu    sync.RWMutex
	users map[string]*User
}

func NewMockRepository() *mockRepository {
	return &mockRepository{
		users: make(map[string]*User)
	}
}

func (m *mockRepository) Create(ctx context.Context, user *User) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if user.ID == uuid.Nil {
		user.ID = uuid.New()
	}

	m.users[user.Email] = user
	return user, nil
}

func (m *mockRepository) FindByEmail(ctx context.Context, email string) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
 
	user, exists = m.users[email]
	if !exists {
		return nil, nil
	}

	return user, nil
}

func (m *mockRepository) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, user := range m.users {
		if user.ID == id {
			return user, nil
		}
	}

	return nil, fmt.Errorf("user not found")
}

func (m *mockRepository) Update(ctx context.Context, user *User) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	elem, exists := m.users[user.Email]
	if !exists {
		return ErrInvalidUser
	}

	m.users[user.Email] = user
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, user := range m.users {
		if user.ID == id {
			delete(m.users, user)
			return nil
		}
	}

	return ErrInvalidUser
}
