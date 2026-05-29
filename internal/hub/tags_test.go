// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"testing"
)

func TestSortTagsSemver(t *testing.T) {
	tags := []Tag{
		{Name: "1.9.0"},
		{Name: "1.10.0"},
		{Name: "1.2.0"},
		{Name: "latest"},
		{Name: "2.0.0"},
	}

	sorted := SortTags(tags)

	if sorted[0].Name != "latest" {
		t.Errorf("expected 'latest' first, got %q", sorted[0].Name)
	}
	if sorted[1].Name != "2.0.0" {
		t.Errorf("expected '2.0.0' second, got %q", sorted[1].Name)
	}
	if sorted[2].Name != "1.10.0" {
		t.Errorf("expected '1.10.0' third (above 1.9.0), got %q", sorted[2].Name)
	}
	if sorted[3].Name != "1.9.0" {
		t.Errorf("expected '1.9.0' fourth, got %q", sorted[3].Name)
	}
	if sorted[4].Name != "1.2.0" {
		t.Errorf("expected '1.2.0' fifth, got %q", sorted[4].Name)
	}
}

func TestSortTagsMixed(t *testing.T) {
	tags := []Tag{
		{Name: "alpine"},
		{Name: "1.0.0"},
		{Name: "latest"},
		{Name: "bookworm"},
		{Name: "2.0.0"},
	}

	sorted := SortTags(tags)

	if sorted[0].Name != "latest" {
		t.Errorf("expected 'latest' first, got %q", sorted[0].Name)
	}

	// semver tags next (2.0.0, 1.0.0)
	if sorted[1].Name != "2.0.0" {
		t.Errorf("expected '2.0.0' after latest, got %q", sorted[1].Name)
	}
	if sorted[2].Name != "1.0.0" {
		t.Errorf("expected '1.0.0' after 2.0.0, got %q", sorted[2].Name)
	}

	// non-semver alphabetical (alpine, bookworm)
	if sorted[3].Name != "alpine" {
		t.Errorf("expected 'alpine' before 'bookworm', got %q", sorted[3].Name)
	}
	if sorted[4].Name != "bookworm" {
		t.Errorf("expected 'bookworm' last, got %q", sorted[4].Name)
	}
}

func TestSortTagsNoLatest(t *testing.T) {
	tags := []Tag{
		{Name: "1.0.0"},
		{Name: "2.0.0"},
	}

	sorted := SortTags(tags)

	if len(sorted) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(sorted))
	}
	if sorted[0].Name != "2.0.0" {
		t.Errorf("expected '2.0.0' first, got %q", sorted[0].Name)
	}
}

func TestSortTagsEmpty(t *testing.T) {
	sorted := SortTags(nil)
	if len(sorted) != 0 {
		t.Errorf("expected empty, got %d", len(sorted))
	}
}

func TestSortTagsOnlyLatest(t *testing.T) {
	sorted := SortTags([]Tag{{Name: "latest"}})
	if len(sorted) != 1 || sorted[0].Name != "latest" {
		t.Error("expected single 'latest' tag")
	}
}

func TestSplitRepo(t *testing.T) {
	tests := []struct {
		input    string
		wantNS   string
		wantName string
	}{
		{"nginx", "library", "nginx"},
		{"bitnami/nginx", "bitnami", "nginx"},
		{"a/b/c", "a", "b/c"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ns, name := splitRepo(tt.input)
			if ns != tt.wantNS || name != tt.wantName {
				t.Errorf("splitRepo(%q) = (%q, %q), want (%q, %q)",
					tt.input, ns, name, tt.wantNS, tt.wantName)
			}
		})
	}
}
