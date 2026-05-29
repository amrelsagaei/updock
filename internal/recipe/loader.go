// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package recipe

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"go.yaml.in/yaml/v3"
)

//go:embed recipes/*.yaml
var embeddedRecipes embed.FS

// LoadAll loads all recipes: embedded first, then user overrides.
func LoadAll(userRecipesDir string) (map[string]*Recipe, error) {
	recipes := make(map[string]*Recipe)

	entries, err := embeddedRecipes.ReadDir("recipes")
	if err == nil {
		for _, e := range entries {
			if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
				continue
			}
			data, readErr := embeddedRecipes.ReadFile("recipes/" + e.Name())
			if readErr != nil {
				continue
			}
			r, parseErr := Parse(data)
			if parseErr != nil {
				continue
			}
			recipes[r.Meta.Name] = r
		}
	}

	if userRecipesDir != "" {
		userEntries, dirErr := os.ReadDir(userRecipesDir)
		if dirErr == nil {
			for _, e := range userEntries {
				if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
					continue
				}
				data, readErr := os.ReadFile(filepath.Join(userRecipesDir, e.Name()))
				if readErr != nil {
					continue
				}
				r, parseErr := Parse(data)
				if parseErr != nil {
					continue
				}
				recipes[r.Meta.Name] = r
			}
		}
	}

	return recipes, nil
}

// Parse parses a recipe from YAML bytes.
func Parse(data []byte) (*Recipe, error) {
	var r Recipe
	if err := yaml.Unmarshal(data, &r); err != nil {
		return nil, fmt.Errorf("parsing recipe: %w", err)
	}
	if err := r.Validate(); err != nil {
		return nil, err
	}
	return &r, nil
}
