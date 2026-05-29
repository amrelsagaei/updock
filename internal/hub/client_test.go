// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSearchWithMockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.URL.Path, "/v2/search/repositories") {
			http.NotFound(w, r)
			return
		}
		query := r.URL.Query().Get("query")
		resp := searchResponse{
			Count: 2,
			Results: []SearchResult{
				{RepoName: "nginx", IsOfficial: true, PullCount: 1000000, StarCount: 1000},
				{RepoName: "someuser/nginx-custom", PullCount: 500, StarCount: 5},
			},
		}
		if query != "nginx" {
			resp = searchResponse{Count: 0, Results: nil}
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestClient(server)
	results, err := client.searchURL(server.URL + "/v2/search/repositories/?query=nginx&page_size=25")
	if err != nil {
		t.Fatalf("Search() error: %v", err)
	}
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d", len(results))
	}
	if results[0].RepoName != "nginx" {
		t.Errorf("expected first result 'nginx', got %q", results[0].RepoName)
	}
}

func TestSearchRateLimited(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("RateLimit-Remaining", "0")
		w.Header().Set("Retry-After", "60")
		w.WriteHeader(http.StatusTooManyRequests)
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.searchURL(server.URL + "/v2/search/repositories/?query=nginx&page_size=25")
	if err == nil {
		t.Fatal("expected rate limit error")
	}

	rle, ok := err.(*RateLimitError)
	if !ok {
		t.Fatalf("expected *RateLimitError, got %T: %v", err, err)
	}
	if rle.RetryAfter != "60" {
		t.Errorf("expected RetryAfter=60, got %q", rle.RetryAfter)
	}
}

func TestSearchHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("internal error"))
	}))
	defer server.Close()

	client := newTestClient(server)
	_, err := client.searchURL(server.URL + "/v2/search/repositories/?query=nginx&page_size=25")
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
	if !strings.Contains(err.Error(), "500") {
		t.Errorf("error should contain status code, got: %v", err)
	}
}

func TestTagsWithMockServer(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		resp := tagsResponse{
			Count: 3,
			Results: []Tag{
				{Name: "latest"},
				{Name: "1.0.0", FullSize: 50 * 1024 * 1024},
				{Name: "2.0.0", FullSize: 60 * 1024 * 1024},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := newTestClient(server)
	tags, err := client.tagsURL(server.URL + "/v2/repositories/library/nginx/tags?page_size=100")
	if err != nil {
		t.Fatalf("Tags() error: %v", err)
	}
	if len(tags) != 3 {
		t.Fatalf("expected 3 tags, got %d", len(tags))
	}
}

func TestRateLimitErrorMessage(t *testing.T) {
	tests := []struct {
		name       string
		retryAfter string
		wantMsg    string
	}{
		{"with retry", "60", "retry after 60 seconds"},
		{"without retry", "", "rate limit exceeded"},
		{"non-numeric", "abc", "rate limit exceeded"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := &RateLimitError{RetryAfter: tt.retryAfter}
			if !strings.Contains(err.Error(), tt.wantMsg) {
				t.Errorf("error %q should contain %q", err.Error(), tt.wantMsg)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	c := NewClient()
	if c == nil {
		t.Fatal("NewClient returned nil")
	}
	if c.http == nil {
		t.Error("http client should not be nil")
	}
	if c.userAgent == "" {
		t.Error("userAgent should not be empty")
	}
}

// helper methods that accept a URL directly (for testing with mock servers)
func newTestClient(server *httptest.Server) *Client {
	return NewClientWith(server.Client())
}

func (c *Client) searchURL(url string) ([]SearchResult, error) {
	var resp searchResponse
	if err := c.getJSON(url, &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}

func (c *Client) tagsURL(url string) ([]Tag, error) {
	var resp tagsResponse
	if err := c.getJSON(url, &resp); err != nil {
		return nil, err
	}
	return resp.Results, nil
}
