package application

import (
	"time"

	"github.com/google/uuid"
)

type ApplicationStatus string

const (
	StatusApplied      ApplicationStatus = "applied"
	StatusInterviewing ApplicationStatus = "interviewing"
	StatusOffer        ApplicationStatus = "offer"
	StatusRejected     ApplicationStatus = "rejected"
	StatusAccepted     ApplicationStatus = "accepted"
)

type Application struct {
	ID            uuid.UUID         `json:"id"`
	UserID        uuid.UUID         `json:"user_id"`
	JobID         uuid.UUID         `json:"job_id"`
	Status        ApplicationStatus `json:"status"`
	AppliedAt     time.Time         `json:"applied_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	InterviewDate *time.Time        `json:"interview_date,omitempty"` // pointer, because it can be nil
	OfferDate     *time.Time        `json:"offer_date,omitempty"`
	Notes         string            `json:"notes,omitempty"`
	SalaryOffer   string            `json:"salary_offer,omitempty"`
	ReminderSent  bool              `json:"reminder_sent,omitempty"`
	FollowUpDate  *time.Time        `json:"follow_up_date,omitempty"`
}

type CreateApplicationInput struct {
	JobID uuid.UUID `json:"job_id" binding:"required"`
	Notes string    `json:"notes,omitempty"`
}

type UpdateApplicationInput struct {
	InterviewDate *time.Time `json:"interview_date,omitempty"`
	OfferDate     *time.Time `json:"offer_date,omitempty"`
	Notes         string     `json:"notes,omitempty"`
	SalaryOffer   string     `json:"salary_offer,omitempty"`
	ReminderSent  bool       `json:"reminder_sent,omitempty"`
	FollowUpDate  *time.Time `json:"follow_up_date,omitempty"`
}

// IsValid validate the application status
func (s ApplicationStatus) IsValid() bool {
	switch s {
	case StatusApplied, StatusInterviewing, StatusOffer, StatusRejected, StatusAccepted:
		return true
	}
	return false
}

// CanTransitionTo checks if a status of an Application can transition to
// a new state based on the business rules
func (a *Application) CanTransitionTo(newStatus ApplicationStatus) bool {
	transitions := map[ApplicationStatus][]ApplicationStatus{
		StatusApplied:      {StatusInterviewing, StatusRejected},
		StatusInterviewing: {StatusOffer, StatusRejected},
		StatusOffer:        {StatusAccepted, StatusRejected},
		StatusAccepted:     {}, // final state
		StatusRejected:     {}, // final state
	}

	allowed, exists := transitions[a.Status]
	if !exists {
		return false
	}

	for _, status := range allowed {
		if status == newStatus {
			return true
		}
	}
	return false
}
