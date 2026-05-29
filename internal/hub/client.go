// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/amrelsagaei/updock/internal/version"
)

const (
	hubBaseURL   = "https://hub.docker.com"
	authURL      = "https://auth.docker.io/token"
	registryURL  = "https://registry-1.docker.io"
	defaultLimit = 25
)

// Client talks to Docker Hub and the Docker registry.
type Client struct {
	http      *http.Client
	userAgent string
}

// NewClient creates a Hub client with sensible defaults.
func NewClient() *Client {
	return &Client{
		http: &http.Client{
			Timeout: 10 * time.Second,
		},
		userAgent: "updock/" + version.Version,
	}
}

// NewClientWith creates a Hub client backed by the given http.Client.
func NewClientWith(h *http.Client) *Client {
	return &Client{
		http:      h,
		userAgent: "updock/" + version.Version,
	}
}

func (c *Client) get(url string) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)
	return c.http.Do(req)
}

func (c *Client) getJSON(url string, target any) error {
	resp, err := c.get(url)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode == http.StatusTooManyRequests {
		remaining := resp.Header.Get("RateLimit-Remaining")
		retry := resp.Header.Get("Retry-After")
		return &RateLimitError{Remaining: remaining, RetryAfter: retry}
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decoding response: %w", err)
	}
	return nil
}

// RateLimitError is returned when Docker Hub rate-limits a request.
type RateLimitError struct {
	Remaining  string
	RetryAfter string
}

func (e *RateLimitError) Error() string {
	msg := "Docker Hub rate limit exceeded"
	if e.RetryAfter != "" {
		if secs, err := strconv.Atoi(e.RetryAfter); err == nil {
			msg += fmt.Sprintf(", retry after %d seconds", secs)
		}
	}
	return msg
}
