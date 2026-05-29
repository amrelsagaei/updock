// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package recipe

import (
	"testing"
)

func testRecipes() map[string]*Recipe {
	return map[string]*Recipe{
		"wordpress": {Meta: Meta{Name: "wordpress"}},
		"ghost":     {Meta: Meta{Name: "ghost"}},
	}
}

func TestMatchExact(t *testing.T) {
	r := Match(testRecipes(), "wordpress")
	if r == nil || r.Meta.Name != "wordpress" {
		t.Error("expected to match wordpress")
	}
}

func TestMatchCaseInsensitive(t *testing.T) {
	r := Match(testRecipes(), "WordPress")
	if r == nil {
		t.Error("expected case-insensitive match")
	}
}

func TestMatchWithNamespace(t *testing.T) {
	r := Match(testRecipes(), "library/wordpress")
	if r == nil {
		t.Error("expected to match by base name after stripping namespace")
	}
}

func TestMatchNotFound(t *testing.T) {
	r := Match(testRecipes(), "redis")
	if r != nil {
		t.Error("expected nil for unmatched name")
	}
}

func TestMatchEmptyMap(t *testing.T) {
	r := Match(nil, "anything")
	if r != nil {
		t.Error("expected nil for nil map")
	}
}
