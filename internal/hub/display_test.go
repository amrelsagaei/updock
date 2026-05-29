// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"bytes"
	"strings"
	"testing"
)

func TestFormatResultsWithResults(t *testing.T) {
	results := []RankedResult{
		{SearchResult: SearchResult{RepoName: "nginx", PullCount: 1000000000}, Badge: "official", Score: 10000},
		{SearchResult: SearchResult{RepoName: "bitnami/nginx", PullCount: 5000000}, Badge: "popular", Score: 5000},
	}

	var buf bytes.Buffer
	FormatResults(&buf, results, "nginx")
	out := buf.String()

	if !strings.Contains(out, "nginx") {
		t.Error("output should contain 'nginx'")
	}
	if !strings.Contains(out, "[best match]") {
		t.Error("first result should have [best match]")
	}
	if !strings.Contains(out, "1)") {
		t.Error("should have numbered results")
	}
	if !strings.Contains(out, "official") {
		t.Error("should show badge")
	}
}

func TestFormatResultsEmpty(t *testing.T) {
	var buf bytes.Buffer
	FormatResults(&buf, nil, "nonexistent")
	out := buf.String()

	if !strings.Contains(out, "No results") {
		t.Error("empty results should show 'No results' message")
	}
	if !strings.Contains(out, "nonexistent") {
		t.Error("should include the query in the message")
	}
}

func TestFormatPulls(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"billions", 2000000000, "2.0B+"},
		{"millions", 5500000, "5.5M+"},
		{"thousands", 1500, "1.5K+"},
		{"small", 42, "42"},
		{"zero", 0, "-"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatPulls(tt.n)
			if got != tt.want {
				t.Errorf("formatPulls(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}
