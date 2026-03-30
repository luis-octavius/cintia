package sources

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/luis-octavius/cintia/internal/job"
)

type LinkedInSource struct {
	BaseURL  string
	Keywords []string
	Location string
}

func NewLinkedInSource(baseURL string, keywords []string, location string) *LinkedInSource {
	if baseURL == "" {
		baseURL = "https://www.linkedin.com/jobs/search"
	}

	return &LinkedInSource{
		BaseURL:  baseURL,
		Keywords: keywords,
		Location: location,
	}
}

func (s *LinkedInSource) Name() string {
	return "linkedin"
}

func (s *LinkedInSource) FetchJobs(ctx context.Context) ([]job.CreateJobInput, error) {
	_ = ctx

	jobs := make([]job.CreateJobInput, 0, len(s.Keywords))
	for _, keyword := range s.Keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			continue
		}

		jobs = append(jobs, job.CreateJobInput{
			Title:       titleFromKeyword(keyword),
			Company:     "LinkedIn Listing",
			Location:    s.Location,
			Description: "Auto-collected placeholder job from LinkedIn source configuration.",
			Source:      "linkedin",
			Link:        s.searchLink(keyword),
			PostedDate:  time.Now(),
		})
	}

	return jobs, nil
}

func (s *LinkedInSource) searchLink(keyword string) string {
	query := url.Values{}
	query.Set("keywords", keyword)
	if s.Location != "" {
		query.Set("location", s.Location)
	}

	return s.BaseURL + "?" + query.Encode()
}
