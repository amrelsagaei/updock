// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package recipe

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadAllEmbedded(t *testing.T) {
	recipes, err := LoadAll("")
	if err != nil {
		t.Fatalf("LoadAll error: %v", err)
	}

	expected := []string{"wordpress", "ghost", "nextcloud", "gitea", "n8n", "plausible"}
	for _, name := range expected {
		if _, ok := recipes[name]; !ok {
			t.Errorf("expected embedded recipe %q to be loaded", name)
		}
	}
}

func TestAllEmbeddedRecipesValidate(t *testing.T) {
	recipes, err := LoadAll("")
	if err != nil {
		t.Fatalf("LoadAll error: %v", err)
	}

	for name, r := range recipes {
		t.Run(name, func(t *testing.T) {
			if err := r.Validate(); err != nil {
				t.Errorf("recipe %q failed validation: %v", name, err)
			}
		})
	}
}

func TestLoadAllUserOverride(t *testing.T) {
	userDir := t.TempDir()

	custom := `
meta:
  name: wordpress
  description: Custom WordPress
  version: "2.0"
  author: me
services:
  wp:
    image: wordpress
    default_tag: custom
`
	if err := os.WriteFile(filepath.Join(userDir, "wordpress.yaml"), []byte(custom), 0o644); err != nil {
		t.Fatal(err)
	}

	recipes, err := LoadAll(userDir)
	if err != nil {
		t.Fatalf("LoadAll error: %v", err)
	}

	r := recipes["wordpress"]
	if r == nil {
		t.Fatal("expected wordpress recipe")
	}
	if r.Meta.Version != "2.0" {
		t.Errorf("user recipe should override embedded, got version %q", r.Meta.Version)
	}
}

func TestLoadAllUserDirNotExist(t *testing.T) {
	recipes, err := LoadAll("/nonexistent/path")
	if err != nil {
		t.Fatalf("should not error for missing user dir: %v", err)
	}
	if len(recipes) == 0 {
		t.Error("should still have embedded recipes")
	}
}

func TestParseInvalidYAML(t *testing.T) {
	_, err := Parse([]byte("invalid: [[[yaml"))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseValidationFailure(t *testing.T) {
	_, err := Parse([]byte("meta:\n  name: \"\"\nservices: {}"))
	if err == nil {
		t.Error("expected validation error")
	}
}
