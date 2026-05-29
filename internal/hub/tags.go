// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"fmt"
	"sort"
	"strings"

	"github.com/Masterminds/semver/v3"
)

// Tag represents a single image tag from Docker Hub.
type Tag struct {
	Name        string `json:"name"`
	LastUpdated string `json:"last_updated"`
	FullSize    int64  `json:"full_size"`
	Digest      string `json:"digest"`
}

type tagsResponse struct {
	Count   int   `json:"count"`
	Results []Tag `json:"results"`
}

// Tags fetches the tag list for a repository from Docker Hub.
func (c *Client) Tags(repo string, limit int) ([]Tag, error) {
	if limit <= 0 {
		limit = 100
	}

	namespace, name := splitRepo(repo)
	u := fmt.Sprintf(
		"%s/v2/repositories/%s/%s/tags?page_size=%d&ordering=last_updated",
		hubBaseURL,
		namespace,
		name,
		limit,
	)

	var resp tagsResponse
	if err := c.getJSON(u, &resp); err != nil {
		return nil, fmt.Errorf("fetching tags: %w", err)
	}

	return resp.Results, nil
}

func splitRepo(repo string) (namespace, name string) {
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) == 1 {
		return "library", parts[0]
	}
	return parts[0], parts[1]
}

// SortTags sorts tags with latest pinned first, then semver descending,
// then non-semver alphabetically.
func SortTags(tags []Tag) []Tag {
	type classified struct {
		tag Tag
		ver *semver.Version
	}

	var latest *Tag
	var withSemver []classified
	var withoutSemver []classified

	for i := range tags {
		t := tags[i]
		if t.Name == "latest" {
			latest = &t
			continue
		}

		v, err := semver.NewVersion(t.Name)
		if err == nil {
			withSemver = append(withSemver, classified{tag: t, ver: v})
		} else {
			withoutSemver = append(withoutSemver, classified{tag: t})
		}
	}

	sort.Slice(withSemver, func(i, j int) bool {
		return withSemver[i].ver.GreaterThan(withSemver[j].ver)
	})

	sort.Slice(withoutSemver, func(i, j int) bool {
		return withoutSemver[i].tag.Name < withoutSemver[j].tag.Name
	})

	result := make([]Tag, 0, len(tags))
	if latest != nil {
		result = append(result, *latest)
	}
	for _, c := range withSemver {
		result = append(result, c.tag)
	}
	for _, c := range withoutSemver {
		result = append(result, c.tag)
	}

	return result
}
