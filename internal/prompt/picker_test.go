// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package prompt

import (
	"fmt"
	"testing"

	"github.com/amrelsagaei/updock/internal/hub"
	"github.com/charmbracelet/huh"
)

func mockSelect(returnValue string) func(string, []huh.Option[string]) (string, error) {
	return func(_ string, _ []huh.Option[string]) (string, error) {
		return returnValue, nil
	}
}

func mockSelectError() func(string, []huh.Option[string]) (string, error) {
	return func(_ string, _ []huh.Option[string]) (string, error) {
		return "", fmt.Errorf("user cancelled")
	}
}

func testResults() []hub.RankedResult {
	return []hub.RankedResult{
		{SearchResult: hub.SearchResult{RepoName: "nginx", PullCount: 1000000}, Badge: "official"},
		{SearchResult: hub.SearchResult{RepoName: "bitnami/nginx", PullCount: 500000}, Badge: "popular"},
	}
}

func TestPickImageBestMatch(t *testing.T) {
	orig := SelectFunc
	defer func() { SelectFunc = orig }()
	SelectFunc = mockSelect("nginx")

	sel, err := PickImage(testResults())
	if err != nil {
		t.Fatalf("PickImage error: %v", err)
	}
	if sel.Result.RepoName != "nginx" {
		t.Errorf("expected nginx, got %q", sel.Result.RepoName)
	}
	if sel.ChooseVersion {
		t.Error("should not be choose version")
	}
}

func TestPickImageChooseVersion(t *testing.T) {
	orig := SelectFunc
	defer func() { SelectFunc = orig }()
	SelectFunc = mockSelect(chooseVersionSentinel)

	sel, err := PickImage(testResults())
	if err != nil {
		t.Fatalf("PickImage error: %v", err)
	}
	if !sel.ChooseVersion {
		t.Error("expected ChooseVersion=true")
	}
	if sel.Result.RepoName != "nginx" {
		t.Errorf("choose version should return first result, got %q", sel.Result.RepoName)
	}
}

func TestPickImageSecondResult(t *testing.T) {
	orig := SelectFunc
	defer func() { SelectFunc = orig }()
	SelectFunc = mockSelect("bitnami/nginx")

	sel, err := PickImage(testResults())
	if err != nil {
		t.Fatalf("PickImage error: %v", err)
	}
	if sel.Result.RepoName != "bitnami/nginx" {
		t.Errorf("expected bitnami/nginx, got %q", sel.Result.RepoName)
	}
}

func TestPickImageError(t *testing.T) {
	orig := SelectFunc
	defer func() { SelectFunc = orig }()
	SelectFunc = mockSelectError()

	_, err := PickImage(testResults())
	if err == nil {
		t.Error("expected error")
	}
}

func TestPickImageEmpty(t *testing.T) {
	_, err := PickImage(nil)
	if err == nil {
		t.Error("expected error for empty results")
	}
}

func TestPickImageUnknownSelection(t *testing.T) {
	orig := SelectFunc
	defer func() { SelectFunc = orig }()
	SelectFunc = mockSelect("unknown-repo")

	sel, err := PickImage(testResults())
	if err != nil {
		t.Fatalf("PickImage error: %v", err)
	}
	if sel.Result.RepoName != "nginx" {
		t.Errorf("unknown selection should fall back to first result, got %q", sel.Result.RepoName)
	}
}

func TestPickTag(t *testing.T) {
	orig := SelectFunc
	defer func() { SelectFunc = orig }()
	SelectFunc = mockSelect("1.0.0")

	tags := []hub.Tag{
		{Name: "latest"},
		{Name: "1.0.0", FullSize: 50 * 1024 * 1024},
	}

	tag, err := PickTag(tags)
	if err != nil {
		t.Fatalf("PickTag error: %v", err)
	}
	if tag.Name != "1.0.0" {
		t.Errorf("expected 1.0.0, got %q", tag.Name)
	}
}

func TestPickTagError(t *testing.T) {
	orig := SelectFunc
	defer func() { SelectFunc = orig }()
	SelectFunc = mockSelectError()

	_, err := PickTag([]hub.Tag{{Name: "latest"}})
	if err == nil {
		t.Error("expected error")
	}
}

func TestPickTagEmpty(t *testing.T) {
	_, err := PickTag(nil)
	if err == nil {
		t.Error("expected error for empty tags")
	}
}

func TestPickTagUnknownFallback(t *testing.T) {
	orig := SelectFunc
	defer func() { SelectFunc = orig }()
	SelectFunc = mockSelect("nonexistent")

	tags := []hub.Tag{{Name: "latest"}}
	tag, err := PickTag(tags)
	if err != nil {
		t.Fatalf("PickTag error: %v", err)
	}
	if tag.Name != "latest" {
		t.Errorf("should fall back to first tag, got %q", tag.Name)
	}
}

func TestBuildImageOptions(t *testing.T) {
	opts := BuildImageOptions(testResults())

	// best match + choose version + 1 remaining = 3
	if len(opts) != 3 {
		t.Fatalf("expected 3 options, got %d", len(opts))
	}
}

func TestBuildImageOptionsSingle(t *testing.T) {
	results := []hub.RankedResult{
		{SearchResult: hub.SearchResult{RepoName: "nginx"}, Badge: "official"},
	}
	opts := BuildImageOptions(results)
	if len(opts) != 2 {
		t.Fatalf("expected 2 options for single result, got %d", len(opts))
	}
}

func TestBuildTagOptions(t *testing.T) {
	tags := []hub.Tag{
		{Name: "latest"},
		{Name: "1.0.0", FullSize: 50 * 1024 * 1024},
	}
	opts := BuildTagOptions(tags)
	if len(opts) != 2 {
		t.Fatalf("expected 2 options, got %d", len(opts))
	}
}

func TestResolveImageSelection(t *testing.T) {
	results := testResults()

	tests := []struct {
		name          string
		selected      string
		wantRepo      string
		wantChooseVer bool
	}{
		{"exact match", "nginx", "nginx", false},
		{"choose version", chooseVersionSentinel, "nginx", true},
		{"second result", "bitnami/nginx", "bitnami/nginx", false},
		{"fallback", "unknown", "nginx", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel := ResolveImageSelection(tt.selected, results)
			if sel.Result.RepoName != tt.wantRepo {
				t.Errorf("got repo %q, want %q", sel.Result.RepoName, tt.wantRepo)
			}
			if sel.ChooseVersion != tt.wantChooseVer {
				t.Errorf("got ChooseVersion=%v, want %v", sel.ChooseVersion, tt.wantChooseVer)
			}
		})
	}
}

func TestResolveTagSelection(t *testing.T) {
	tags := []hub.Tag{{Name: "latest"}, {Name: "1.0.0"}}

	tests := []struct {
		name     string
		selected string
		wantTag  string
	}{
		{"exact", "1.0.0", "1.0.0"},
		{"latest", "latest", "latest"},
		{"fallback", "nonexistent", "latest"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tag := ResolveTagSelection(tt.selected, tags)
			if tag.Name != tt.wantTag {
				t.Errorf("got %q, want %q", tag.Name, tt.wantTag)
			}
		})
	}
}

func TestFormatCount(t *testing.T) {
	tests := []struct {
		name string
		n    int
		want string
	}{
		{"billions", 1500000000, "1.5B"},
		{"millions", 2300000, "2.3M"},
		{"thousands", 4700, "4.7K"},
		{"small", 42, "42"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatCount(tt.n)
			if got != tt.want {
				t.Errorf("FormatCount(%d) = %q, want %q", tt.n, got, tt.want)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name string
		b    int64
		want string
	}{
		{"gigabytes", 2 * 1024 * 1024 * 1024, "2.0 GB"},
		{"megabytes", 50 * 1024 * 1024, "50.0 MB"},
		{"kilobytes", 512 * 1024, "512.0 KB"},
		{"bytes", 100, "100 B"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatBytes(tt.b)
			if got != tt.want {
				t.Errorf("FormatBytes(%d) = %q, want %q", tt.b, got, tt.want)
			}
		})
	}
}
