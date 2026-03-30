package sources

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/luis-octavius/cintia/internal/job"
)

type IndeedSource struct {
	BaseURL  string
	Keywords []string
	Location string
}

func NewIndeedSource(baseURL string, keywords []string, location string) *IndeedSource {
	if baseURL == "" {
		baseURL = "https://www.indeed.com/jobs"
	}

	return &IndeedSource{
		BaseURL:  baseURL,
		Keywords: keywords,
		Location: location,
	}
}

func (s *IndeedSource) Name() string {
	return "indeed"
}

func (s *IndeedSource) FetchJobs(ctx context.Context) ([]job.CreateJobInput, error) {
	_ = ctx

	// MVP: produce deterministic seed jobs from configured keywords so scheduler
	// and persistence flow can be validated without external scraping fragility.
	jobs := make([]job.CreateJobInput, 0, len(s.Keywords))
	for _, keyword := range s.Keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			continue
		}

		jobs = append(jobs, job.CreateJobInput{
			Title:       titleFromKeyword(keyword),
			Company:     "Indeed Listing",
			Location:    s.Location,
			Description: "Auto-collected placeholder job from Indeed source configuration.",
			Source:      "indeed",
			Link:        s.searchLink(keyword),
			PostedDate:  time.Now(),
		})
	}

	return jobs, nil
}

func (s *IndeedSource) searchLink(keyword string) string {
	query := url.Values{}
	query.Set("q", keyword)
	if s.Location != "" {
		query.Set("l", s.Location)
	}

	return s.BaseURL + "?" + query.Encode()
}
