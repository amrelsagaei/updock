// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package recipe

import (
	"fmt"
	"strings"
)

// Recipe is the top-level structure of a recipe YAML file.
type Recipe struct {
	Meta     Meta                `yaml:"meta"`
	Prompts  []Prompt            `yaml:"prompts"`
	Services map[string]Service  `yaml:"services"`
	Volumes  map[string]struct{} `yaml:"volumes"`
}

// Meta holds recipe metadata.
type Meta struct {
	Name        string   `yaml:"name"`
	Description string   `yaml:"description"`
	Version     string   `yaml:"version"`
	Author      string   `yaml:"author"`
	Tags        []string `yaml:"tags"`
}

// Prompt defines a user-facing question in the recipe.
type Prompt struct {
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Default     string `yaml:"default"`
	Required    bool   `yaml:"required"`
	Secret      bool   `yaml:"secret"`
	Generate    string `yaml:"generate"`
	Type        string `yaml:"type"`
}

// Service defines a service in the recipe.
type Service struct {
	Image       string            `yaml:"image"`
	DefaultTag  string            `yaml:"default_tag"`
	Ports       []string          `yaml:"ports"`
	Environment map[string]string `yaml:"environment"`
	DependsOn   []string          `yaml:"depends_on"`
	Volumes     []string          `yaml:"volumes"`
	Restart     string            `yaml:"restart"`
}

// Validate checks that a recipe has all required fields and that all
// ${VAR} references resolve to a defined prompt.
func (r *Recipe) Validate() error {
	if r.Meta.Name == "" {
		return fmt.Errorf("recipe: meta.name is required")
	}
	if len(r.Services) == 0 {
		return fmt.Errorf("recipe %q: must have at least one service", r.Meta.Name)
	}

	promptNames := make(map[string]bool)
	for _, p := range r.Prompts {
		if p.Name == "" {
			return fmt.Errorf("recipe %q: prompt name cannot be empty", r.Meta.Name)
		}
		promptNames[p.Name] = true
	}

	for svcName := range r.Services {
		svc := r.Services[svcName]
		if svc.Image == "" {
			return fmt.Errorf("recipe %q: service %q must have an image", r.Meta.Name, svcName)
		}

		for _, port := range svc.Ports {
			for _, ref := range extractVarRefs(port) {
				if !promptNames[ref] {
					return fmt.Errorf("recipe %q: service %q port references undefined variable ${%s}", r.Meta.Name, svcName, ref)
				}
			}
		}

		for key, val := range svc.Environment {
			for _, ref := range extractVarRefs(val) {
				if !promptNames[ref] {
					return fmt.Errorf("recipe %q: service %q env %s references undefined variable ${%s}", r.Meta.Name, svcName, key, ref)
				}
			}
		}
	}

	return nil
}

func extractVarRefs(s string) []string {
	var refs []string
	for {
		start := strings.Index(s, "${")
		if start < 0 {
			break
		}
		end := strings.Index(s[start:], "}")
		if end < 0 {
			break
		}
		ref := s[start+2 : start+end]
		if ref != "" {
			refs = append(refs, ref)
		}
		s = s[start+end+1:]
	}
	return refs
}
