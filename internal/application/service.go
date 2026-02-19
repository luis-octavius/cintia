package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/luis-octavius/cintia/internal/job"
	"github.com/luis-octavius/cintia/internal/user"
)

var (
	ErrAlreadyApplied      = errors.New("user already applied to this job")
	ErrJobNotFound         = errors.New("job not found")
	ErrJobInactive         = errors.New("job inactive")
	ErrUserNotFound        = errors.New("user not found")
	ErrApplicationNotFound = errors.New("application not found")
)

type Service interface {
	CreateApplication(ctx context.Context, userID uuid.UUID, input CreateApplicationInput) (*Application, error)
	GetApplicationByID(ctx context.Context, id uuid.UUID) (*Application, error)
	GetUserApplications(ctx context.Context, userID uuid.UUID) ([]*Application, error)
	GetJobApplications(ctx context.Context, jobID uuid.UUID) ([]*Application, error)
	UpdateApplication(ctx context.Context, id uuid.UUID, updates UpdateApplicationInput) error
	UpdateApplicationStatus(ctx context.Context, id uuid.UUID, status ApplicationStatus) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type service struct {
	repo        Repository
	jobService  job.Service
	userService user.Service
}

func NewService(repo Repository, jobService job.Service, userService user.Service) Service {
	return &service{
		repo:        repo,
		jobService:  jobService,
		userService: userService,
	}
}

func (s *service) CreateApplication(ctx context.Context, userID uuid.UUID, input CreateApplicationInput) (*Application, error) {
	if input.JobID == uuid.Nil {
		return nil, errors.New("job_id is required")
	}

	existing, _ := s.repo.GetUserJobApplication(ctx, userID, input.JobID)
	if existing != nil {
		return nil, ErrAlreadyApplied
	}

	job, err := s.jobService.GetJob(ctx, input.JobID)
	if err != nil {
		return nil, ErrJobNotFound
	}

	// business rules of job
	if !job.IsActive {
		return nil, ErrJobInactive
	}

	app := &Application{
		ID:        uuid.New(),
		UserID:    userID,
		JobID:     input.JobID,
		Status:    StatusApplied,
		AppliedAt: time.Now(),
		UpdatedAt: time.Now(),
		Notes:     input.Notes,
	}

	return s.repo.Create(ctx, app)
}

func (s *service) GetApplicationByID(ctx context.Context, id uuid.UUID) (*Application, error) {
	application, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrApplicationNotFound
		}
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	return application, nil
}

func (s *service) GetUserApplications(ctx context.Context, userID uuid.UUID) ([]*Application, error) {
	// sees if user exists
	_, err := s.userService.GetProfile(ctx, userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	applications, err := s.repo.GetUserApplications(ctx, userID)
	if err != nil {
		return nil, err
	}

	return applications, nil
}

func (s *service) GetJobApplications(ctx context.Context, jobID uuid.UUID) ([]*Application, error) {
	_, err := s.jobService.GetJob(ctx, jobID)
	if err != nil {
		return nil, ErrJobNotFound
	}

	applications, err := s.repo.GetJobApplications(ctx, jobID)

	return applications, nil
}

func (s *service) UpdateApplication(ctx context.Context, id uuid.UUID, updates UpdateApplicationInput) error {
	application, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return ErrApplicationNotFound
	}

	// validate updates
	updated := false

	if updates.InterviewDate != nil {
		application.InterviewDate = updates.InterviewDate
		updated = true
	}

	if updates.OfferDate != nil {
		application.OfferDate = updates.OfferDate
		updated = true
	}

	if updates.Notes != "" {
		application.Notes = updates.Notes
		updated = true
	}

	if updates.SalaryOffer != "" {
		application.SalaryOffer = updates.SalaryOffer
		updated = true
	}

	if updates.ReminderSent {
		application.ReminderSent = updates.ReminderSent
		updated = true
	}

	if updates.FollowUpDate != nil {
		application.FollowUpDate = updates.FollowUpDate
		updated = true
	}

	if !updated {
		return errors.New("nothing to update")
	}

	application.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, application); err != nil {
		return fmt.Errorf("failed to update application: %w", err)
	}

	return nil
}

func (s *service) UpdateApplicationStatus(ctx context.Context, id uuid.UUID, status ApplicationStatus) error {
	app, err := s.GetApplicationByID(ctx, id)
	if err != nil {
		return ErrApplicationNotFound
	}

	if !app.CanTransitionTo(status) {
		return fmt.Errorf("cannot transition from %s to %s", app.Status, status)
	}

	if !status.IsValid() {
		return ErrInvalidStatus
	}

	return s.repo.UpdateStatus(ctx, id, status)
}

func (s *service) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := s.GetApplicationByID(ctx, id)
	if err != nil {
		return ErrApplicationNotFound
	}

	err = s.repo.Delete(ctx, id)
	if err != nil {
		return fmt.Errorf("error deleting application: %v", err)
	}

	return nil
}
