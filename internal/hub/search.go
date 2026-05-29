// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"fmt"
	"net/url"
)

// SearchResult represents a single image result from Docker Hub.
type SearchResult struct {
	RepoName    string `json:"repo_name"`
	IsOfficial  bool   `json:"is_official"`
	IsAutomated bool   `json:"is_automated"`
	StarCount   int    `json:"star_count"`
	PullCount   int    `json:"pull_count"`
	Description string `json:"description"`
}

type searchResponse struct {
	Count   int            `json:"count"`
	Results []SearchResult `json:"results"`
}

// Search queries Docker Hub for images matching the given term.
func (c *Client) Search(query string, limit int) ([]SearchResult, error) {
	if limit <= 0 {
		limit = defaultLimit
	}

	u := fmt.Sprintf(
		"%s/v2/search/repositories/?query=%s&page_size=%d",
		hubBaseURL,
		url.QueryEscape(query),
		limit,
	)

	var resp searchResponse
	if err := c.getJSON(u, &resp); err != nil {
		return nil, fmt.Errorf("searching Docker Hub: %w", err)
	}

	return resp.Results, nil
}
