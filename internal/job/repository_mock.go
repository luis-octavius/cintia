package job

import (
	"context"
	"strings"
	"sync"

	"github.com/google/uuid"
)

type mockRepository struct {
	mu    sync.RWMutex
	jobs  map[uuid.UUID]*Job
	links map[string]*Job
}

func NewMockRepository() Repository {
	return &mockRepository{
		jobs:  map[uuid.UUID]*Job{},
		links: map[string]*Job{},
	}
}

func (m *mockRepository) Create(ctx context.Context, job *Job) (*Job, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.jobs[job.ID] = job
	m.links[job.Link] = job

	return job, nil
}

func (m *mockRepository) GetByLink(ctx context.Context, link string) (*Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	job, exists := m.links[link]
	if !exists {
		return nil, ErrNotFound
	}

	return job, nil
}

func (m *mockRepository) GetByID(ctx context.Context, id uuid.UUID) (*Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	job, exists := m.jobs[id]
	if !exists {
		return nil, ErrNotFound
	}

	return job, nil
}

func (m *mockRepository) Search(ctx context.Context, filters JobFilters) ([]*Job, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var jobs []*Job

	for _, job := range m.jobs {
		if filters.Limit != 0 && len(jobs) >= filters.Limit {
			break
		}

		if m.matchesFilters(job, filters) {
			jobs = append(jobs, job)
		}
	}

	return jobs, nil
}

func (m *mockRepository) Count(ctx context.Context, filters JobFilters) (int, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var jobs []*Job

	for _, job := range m.jobs {
		if filters.Limit > 0 && len(jobs) >= filters.Limit {
			break
		}

		if m.matchesFilters(job, filters) {
			jobs = append(jobs, job)
		}
	}

	return len(jobs), nil
}

func (m *mockRepository) Update(ctx context.Context, job *Job) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	existing, exists := m.jobs[job.ID]
	if !exists {
		return ErrNotFound
	}

	// if links has change, remove old entry and adds new
	if existing.Link != job.Link {
		delete(m.links, existing.Link)
		m.links[job.Link] = job
	}

	m.jobs[job.ID] = job
	return nil
}

func (m *mockRepository) MarkAsInactive(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	job, exists := m.jobs[id]
	if !exists {
		return ErrNotFound
	}

	job.IsActive = false
	return nil
}

func (m *mockRepository) Delete(ctx context.Context, id uuid.UUID) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, exists := m.jobs[id]
	if !exists {
		return ErrNotFound
	}

	delete(m.jobs, id)
	return nil
}

func (m *mockRepository) matchesFilters(job *Job, filters JobFilters) bool {
	if filters.Title != "" && !strings.Contains(strings.ToLower(job.Title), strings.ToLower(filters.Title)) {
		return false
	}

	if filters.Company != "" && !strings.Contains(strings.ToLower(job.Company), strings.ToLower(filters.Company)) {
		return false
	}

	if filters.Location != "" && !strings.Contains(strings.ToLower(job.Location), strings.ToLower(filters.Location)) {
		return false
	}

	if filters.Source != "" && job.Source != filters.Source {
		return false
	}

	return true
}
