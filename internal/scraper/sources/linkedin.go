package sources

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/luis-octavius/cintia/internal/job"
)

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type LinkedInSource struct {
	BaseURL  string
	Keywords []string
	Location string
	Client   HTTPClient
}

func NewLinkedInSource(baseURL string, keywords []string, location string) *LinkedInSource {
	if baseURL == "" {
		baseURL = "https://www.linkedin.com/jobs/search"
	}

	return &LinkedInSource{
		BaseURL:  baseURL,
		Keywords: keywords,
		Location: location,
		Client:   &http.Client{Timeout: 15 * time.Second},
	}
}

func (s *LinkedInSource) Name() string {
	return "linkedin"
}

func (s *LinkedInSource) FetchJobs(ctx context.Context) ([]job.CreateJobInput, error) {
	jobs := make([]job.CreateJobInput, 0)
	seen := make(map[string]struct{})

	for _, keyword := range s.Keywords {
		keyword = strings.TrimSpace(keyword)
		if keyword == "" {
			continue
		}

		html, err := s.fetchHTML(ctx, s.searchLink(keyword))
		if err != nil {
			return nil, err
		}

		parsed := parseLinkedInJobs(html, s.Location)
		for _, item := range parsed {
			if _, exists := seen[item.Link]; exists {
				continue
			}

			seen[item.Link] = struct{}{}
			jobs = append(jobs, item)
		}
	}

	return jobs, nil
}

func (s *LinkedInSource) fetchHTML(ctx context.Context, pageURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pageURL, nil)
	if err != nil {
		return "", fmt.Errorf("create linkedin request: %w", err)
	}

	req.Header.Set("User-Agent", "cintia-scraper/1.0")

	res, err := s.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("linkedin request failed: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return "", fmt.Errorf("linkedin request status: %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("read linkedin body: %w", err)
	}

	return string(body), nil
}

func (s *LinkedInSource) searchLink(keyword string) string {
	query := url.Values{}
	query.Set("keywords", keyword)
	if s.Location != "" {
		query.Set("location", s.Location)
	}

	return s.BaseURL + "?" + query.Encode()
}
