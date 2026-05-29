// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"math"
	"sort"
	"strings"
)

// RankedResult is a SearchResult augmented with a ranking score.
type RankedResult struct {
	SearchResult
	Score float64
	Badge string
}

// Rank scores and sorts search results for the given query.
func Rank(results []SearchResult, query string) []RankedResult {
	ranked := make([]RankedResult, len(results))
	for i, r := range results {
		ranked[i] = RankedResult{
			SearchResult: r,
			Score:        score(r, query),
			Badge:        badge(r),
		}
	}

	sort.Slice(ranked, func(i, j int) bool {
		return ranked[i].Score > ranked[j].Score
	})

	return ranked
}

func score(r SearchResult, query string) float64 {
	s := 0.0

	if r.IsOfficial {
		s += 10000
	}

	if r.PullCount > 0 {
		s += math.Log10(float64(r.PullCount)) * 100
	}
	if r.StarCount > 0 {
		s += math.Log10(float64(r.StarCount)) * 50
	}

	s += nameSimilarity(r.RepoName, query) * 200

	return s
}

func badge(r SearchResult) string {
	if r.IsOfficial {
		return "official"
	}
	if r.PullCount > 1_000_000 {
		return "popular"
	}
	return "community"
}

func nameSimilarity(repoName, query string) float64 {
	name := repoName
	if idx := strings.LastIndex(repoName, "/"); idx >= 0 {
		name = repoName[idx+1:]
	}

	name = strings.ToLower(name)
	query = strings.ToLower(query)

	if name == query {
		return 1.0
	}

	if strings.Contains(name, query) || strings.Contains(query, name) {
		longer := len(name)
		if len(query) > longer {
			longer = len(query)
		}
		shorter := len(name)
		if len(query) < shorter {
			shorter = len(query)
		}
		return 0.5 + 0.4*float64(shorter)/float64(longer)
	}

	dist := levenshtein(name, query)
	maxLen := len(name)
	if len(query) > maxLen {
		maxLen = len(query)
	}
	if maxLen == 0 {
		return 0
	}

	return math.Max(0, 1.0-float64(dist)/float64(maxLen))
}

func levenshtein(a, b string) int {
	la, lb := len(a), len(b)
	if la == 0 {
		return lb
	}
	if lb == 0 {
		return la
	}

	prev := make([]int, lb+1)
	curr := make([]int, lb+1)

	for j := 0; j <= lb; j++ {
		prev[j] = j
	}

	for i := 1; i <= la; i++ {
		curr[0] = i
		for j := 1; j <= lb; j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = minInt(
				prev[j]+1,
				curr[j-1]+1,
				prev[j-1]+cost,
			)
		}
		prev, curr = curr, prev
	}

	return prev[lb]
}

func minInt(vals ...int) int {
	m := vals[0]
	for _, v := range vals[1:] {
		if v < m {
			m = v
		}
	}
	return m
}
