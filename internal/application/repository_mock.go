package application

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/google/uuid"
)

type mockRepository struct {
	mu           sync.RWMutex
	applications map[uuid.UUID]*Application
}

func NewMockRepository() *mockRepository {
	return &mockRepository{
		applications: make(map[uuid.UUID]*Application)
	}
}

func (m *mockRepository) Create(ctx context.Context, app *Application) (*Application, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if app.ID == uuid.Nil {
		app.ID = uuid.New()
	}

	if _, exists := m.applications[app.ID]; exists {
		return nil, ErrAlreadyExists
	}

	now := time.Now()
	if app.AppliedAt.IsZero() {
		app.AppliedAt = now
	}
	app.UpdatedAt = now

	m.applications[app.ID] = app
	return app, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id uuid.UUID) (*Application, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	application, exists := m.applications[id]
	if !exists {
		return nil, ErrNotFound
	}

	return application, nil
}

func (m *mockRepository) GetUserApplications(ctx context.Context, userID uuid.UUID) ([]*Application, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	applications := []*Application{}

	for _, application := range m.applications {
		if application.UserID == userID {
			applications = append(applications, application)
		}
	}

	return applications, nil
}

func (m *mockRepository) GetUserJobApplication(ctx context.Context, userID, jobID uuid.UUID) (*Application, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, app := range m.applications {
		if app.UserID == userID && app.JobID == jobID {
			return app, nil
		}
	}

	return nil, ErrNotFound
}

func (m *mockRepository) GetJobApplications(ctx context.Context, jobID uuid.UUID) ([]*Application, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	applications := []*Application{}

	for _, application := range m.applications {
		if application.JobID == jobID {
			applications = append(applications, application)
		}
	}

	return applications, nil
}

func (m *mockRepository) Update(ctx context.Context, app *Application) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	application, exists := m.applications[app.id]
	if !exists {
		return ErrNotFound
	}
	 
	if app.Status != {
		if !app.Status.IsValid() {
			return ErrInvalidStatus
		}
		if !application.CanTransitionTo(app.Status) {
			return ErrInvalidTransition
		}
	}

	application.Status = app.Status 
	application.InterviewDate = app.InterviewDate
	application.OfferDate = app.OfferDate 
	application.SalaryOffer = app.SalaryOffer 
	application.ReminderSent = app.ReminderSent 
	application.FollowUpdate = app.FollowUpdate 
	application.UpdatedAt = time.Now()

	m.applications[app.ID] = application
	return nil
}

func (m *mockRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status ApplicationStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	application, exists := m.applications[id]
	if !exists {
		return ErrNotFound
	}

	if ok := status.IsValid(); !ok {
		return ErrInvalidStatus
	}

	if ok := application.CanTransitionTo(status); !ok {
		return ErrInvalidTransition
	}

	m.applications[id].Status = status
	m.applications[id].UpdatedAt = time.Now()
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.applications[id]
	if !exists {
		return ErrNotFound
	}

	delete(m.applications, id)
	return nil
}
