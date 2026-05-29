// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"testing"
)

func TestRankOfficialFirst(t *testing.T) {
	results := []SearchResult{
		{RepoName: "someuser/nginx-fork", IsOfficial: false, PullCount: 100, StarCount: 5},
		{RepoName: "nginx", IsOfficial: true, PullCount: 1000000, StarCount: 1000},
		{RepoName: "another/nginx-custom", IsOfficial: false, PullCount: 500, StarCount: 10},
	}

	ranked := Rank(results, "nginx")

	if ranked[0].RepoName != "nginx" {
		t.Errorf("expected official 'nginx' first, got %q", ranked[0].RepoName)
	}
	if ranked[0].Badge != "official" {
		t.Errorf("expected badge 'official', got %q", ranked[0].Badge)
	}
}

func TestRankJuiceShop(t *testing.T) {
	results := []SearchResult{
		{RepoName: "someuser/juice-shop-fork", PullCount: 50, StarCount: 1},
		{RepoName: "bkimminich/juice-shop", PullCount: 50000000, StarCount: 500},
		{RepoName: "bkimminich/juice-shop-ctf", PullCount: 10000, StarCount: 20},
	}

	ranked := Rank(results, "juice-shop")

	if ranked[0].RepoName != "bkimminich/juice-shop" {
		t.Errorf("expected 'bkimminich/juice-shop' first, got %q", ranked[0].RepoName)
	}
}

func TestRankEmptyResults(t *testing.T) {
	ranked := Rank(nil, "anything")
	if len(ranked) != 0 {
		t.Errorf("expected empty results, got %d", len(ranked))
	}
}

func TestRankScoreDecreasing(t *testing.T) {
	results := []SearchResult{
		{RepoName: "a", PullCount: 1},
		{RepoName: "b", PullCount: 1000000},
		{RepoName: "c", PullCount: 100},
	}

	ranked := Rank(results, "test")
	for i := 1; i < len(ranked); i++ {
		if ranked[i].Score > ranked[i-1].Score {
			t.Errorf("results not sorted: index %d (score %.2f) > index %d (score %.2f)",
				i, ranked[i].Score, i-1, ranked[i-1].Score)
		}
	}
}

func TestBadge(t *testing.T) {
	tests := []struct {
		name      string
		result    SearchResult
		wantBadge string
	}{
		{"official", SearchResult{IsOfficial: true}, "official"},
		{"popular", SearchResult{PullCount: 2000000}, "popular"},
		{"community", SearchResult{PullCount: 100}, "community"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := badge(tt.result)
			if got != tt.wantBadge {
				t.Errorf("badge() = %q, want %q", got, tt.wantBadge)
			}
		})
	}
}

func TestNameSimilarity(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		query    string
		wantMin  float64
	}{
		{"exact match", "nginx", "nginx", 1.0},
		{"exact with namespace", "library/nginx", "nginx", 1.0},
		{"substring", "nginx-proxy", "nginx", 0.5},
		{"completely different", "postgres", "redis", 0.0},
		{"empty query", "nginx", "", 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := nameSimilarity(tt.repoName, tt.query)
			if got < tt.wantMin {
				t.Errorf("nameSimilarity(%q, %q) = %.2f, want >= %.2f",
					tt.repoName, tt.query, got, tt.wantMin)
			}
		})
	}
}

func TestLevenshtein(t *testing.T) {
	tests := []struct {
		a, b string
		want int
	}{
		{"", "", 0},
		{"a", "", 1},
		{"", "b", 1},
		{"kitten", "sitting", 3},
		{"nginx", "nginx", 0},
		{"nginx", "nginz", 1},
	}

	for _, tt := range tests {
		t.Run(tt.a+"_"+tt.b, func(t *testing.T) {
			got := levenshtein(tt.a, tt.b)
			if got != tt.want {
				t.Errorf("levenshtein(%q, %q) = %d, want %d", tt.a, tt.b, got, tt.want)
			}
		})
	}
}
