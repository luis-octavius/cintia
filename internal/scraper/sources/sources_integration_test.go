package sources

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLinkedInSource_FetchJobs_ParsesHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
			<html><body>
				<a href="https://www.linkedin.com/jobs/view/12345?utm_source=test">Senior Go Engineer</a>
				<a href="https://www.linkedin.com/jobs/view/67890">Backend Engineer</a>
			</body></html>
		`))
	}))
	defer server.Close()

	source := NewLinkedInSource(server.URL, []string{"golang"}, "Remote")
	source.Client = server.Client()

	jobs, err := source.FetchJobs(context.Background())
	if err != nil {
		t.Fatalf("unexpected fetch error: %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs parsed, got %d", len(jobs))
	}

	if jobs[0].Title != "Senior Go Engineer" {
		t.Fatalf("unexpected title: %s", jobs[0].Title)
	}

	if jobs[0].Source != "linkedin" {
		t.Fatalf("unexpected source: %s", jobs[0].Source)
	}

	if jobs[0].Link != "https://www.linkedin.com/jobs/view/12345" {
		t.Fatalf("expected normalized link without tracking query, got %s", jobs[0].Link)
	}
}

func TestIndeedSource_FetchJobs_ParsesHTML(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`
			<html><body>
				<a href="https://www.indeed.com/viewjob?jk=abc123&utm_campaign=tracking">Go Developer</a>
				<a href="https://www.indeed.com/rc/clk?jk=def456">Platform Engineer</a>
			</body></html>
		`))
	}))
	defer server.Close()

	source := NewIndeedSource(server.URL, []string{"backend"}, "Sao Paulo")
	source.Client = server.Client()

	jobs, err := source.FetchJobs(context.Background())
	if err != nil {
		t.Fatalf("unexpected fetch error: %v", err)
	}

	if len(jobs) != 2 {
		t.Fatalf("expected 2 jobs parsed, got %d", len(jobs))
	}

	if jobs[0].Title != "Go Developer" {
		t.Fatalf("unexpected title: %s", jobs[0].Title)
	}

	if jobs[0].Source != "indeed" {
		t.Fatalf("unexpected source: %s", jobs[0].Source)
	}

	if jobs[0].Link != "https://www.indeed.com/viewjob?jk=abc123" {
		t.Fatalf("expected normalized indeed link, got %s", jobs[0].Link)
	}
}
