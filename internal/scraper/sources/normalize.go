package sources

import (
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/luis-octavius/cintia/internal/job"
)

var (
	anchorRegex = regexp.MustCompile(`(?is)<a[^>]+href=["']([^"']+)["'][^>]*>(.*?)</a>`)
	tagRegex    = regexp.MustCompile(`(?is)<[^>]+>`)
)

func titleFromKeyword(keyword string) string {
	parts := strings.Fields(strings.TrimSpace(keyword))
	if len(parts) == 0 {
		return "Software Engineer"
	}

	for i := range parts {
		parts[i] = strings.ToUpper(parts[i][:1]) + strings.ToLower(parts[i][1:])
	}

	return strings.Join(parts, " ") + " Engineer"
}

func parseLinkedInJobs(html, location string) []job.CreateJobInput {
	jobs := make([]job.CreateJobInput, 0)
	for _, match := range anchorRegex.FindAllStringSubmatch(html, -1) {
		href := strings.TrimSpace(match[1])
		if !strings.Contains(href, "/jobs/view") {
			continue
		}

		title := cleanText(match[2])
		if title == "" {
			title = "Software Engineer"
		}

		jobs = append(jobs, job.CreateJobInput{
			Title:       title,
			Company:     "Unknown Company",
			Location:    location,
			Description: "Scraped from LinkedIn search results page.",
			Source:      "linkedin",
			Link:        NormalizeJobLink(href),
			PostedDate:  time.Now(),
		})
	}

	return jobs
}

func parseIndeedJobs(html, location string) []job.CreateJobInput {
	jobs := make([]job.CreateJobInput, 0)
	for _, match := range anchorRegex.FindAllStringSubmatch(html, -1) {
		href := strings.TrimSpace(match[1])
		if !strings.Contains(href, "/viewjob") && !strings.Contains(href, "/rc/clk") {
			continue
		}

		title := cleanText(match[2])
		if title == "" {
			title = "Software Engineer"
		}

		jobs = append(jobs, job.CreateJobInput{
			Title:       title,
			Company:     "Unknown Company",
			Location:    location,
			Description: "Scraped from Indeed search results page.",
			Source:      "indeed",
			Link:        NormalizeJobLink(href),
			PostedDate:  time.Now(),
		})
	}

	return jobs
}

func NormalizeJobLink(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return raw
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return raw
	}

	query := parsed.Query()
	for key := range query {
		k := strings.ToLower(key)
		if strings.HasPrefix(k, "utm_") || k == "trk" || k == "refid" || k == "ref" {
			query.Del(key)
		}
	}

	parsed.RawQuery = query.Encode()
	parsed.Fragment = ""
	parsed.Host = strings.ToLower(parsed.Host)

	if parsed.Path != "/" {
		parsed.Path = strings.TrimSuffix(parsed.Path, "/")
	}

	return parsed.String()
}

func cleanText(input string) string {
	stripped := tagRegex.ReplaceAllString(input, " ")
	return strings.Join(strings.Fields(stripped), " ")
}
