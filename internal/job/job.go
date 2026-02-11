package job

import (
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID           uuid.UUID `json:"id"`
	Title        string    `json:"title"`
	Company      string    `json:"company"`
	Location     string    `json:"location"`
	Description  string    `json:"description"`
	SalaryRange  string    `json:"salary_range,omitempty"`
	Requirements string    `json:"requirements,omitempty"`
	Source       string    `json:"source"`
	Link         string    `json:"link"`
	PostedDate   time.Time `json:"posted_date"`
	ScrapedAt    time.Time `json:"scraped_at"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateJobInput struct {
	Title        string    `json:"title"`
	Company      string    `json:"company"`
	Location     string    `json:"location"`
	Description  string    `json:"description"`
	SalaryRange  string    `json:"salary_range,omitempty"`
	Requirements string    `json:"requirements,omitempty"`
	Source       string    `json:"source"`
	Link         string    `json:"link"`
	PostedDate   time.Time `json:"posted_date"`
}

type JobFilters struct {
	Title    string `json:"title,omitempty"`
	Company  string `json:"company,omitempty"`
	Location string `json:"location,omitempty"`
	IsActive *bool  `json:"is_active,omitempty"`
	Source   string `json:"source,omitempty"`
	Page     int    `json:"page,omitempty"`
	Limit    int    `json:"limit,omitempty"`
}

type JobsResponse struct {
	Jobs       []*Job `json:"jobs"`
	Total      int    `json:"total"`
	Page       int    `json:"page"`
	TotalPages int    `json:"total_pages"`
	HasMore    bool   `json:"has_more"`
}

type UpdateJobInput struct {
	Title   string `json:"title"`
	Company string `json:"company"`

	IsActive bool `json:"is_active,omitempty"`
}
