// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package recipe

import (
	"testing"
)

func TestValidateValid(t *testing.T) {
	r := &Recipe{
		Meta: Meta{Name: "test"},
		Prompts: []Prompt{
			{Name: "PORT", Default: "8080"},
		},
		Services: map[string]Service{
			"app": {Image: "nginx", Ports: []string{"${PORT}:80"}},
		},
	}

	if err := r.Validate(); err != nil {
		t.Errorf("expected valid, got: %v", err)
	}
}

func TestValidateMissingName(t *testing.T) {
	r := &Recipe{
		Services: map[string]Service{"app": {Image: "nginx"}},
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing name")
	}
}

func TestValidateNoServices(t *testing.T) {
	r := &Recipe{
		Meta: Meta{Name: "test"},
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for no services")
	}
}

func TestValidateEmptyPromptName(t *testing.T) {
	r := &Recipe{
		Meta:    Meta{Name: "test"},
		Prompts: []Prompt{{Name: ""}},
		Services: map[string]Service{
			"app": {Image: "nginx"},
		},
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for empty prompt name")
	}
}

func TestValidateServiceMissingImage(t *testing.T) {
	r := &Recipe{
		Meta: Meta{Name: "test"},
		Services: map[string]Service{
			"app": {Image: ""},
		},
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for missing image")
	}
}

func TestValidateUndefinedVarInPort(t *testing.T) {
	r := &Recipe{
		Meta: Meta{Name: "test"},
		Services: map[string]Service{
			"app": {Image: "nginx", Ports: []string{"${UNDEFINED}:80"}},
		},
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for undefined var in port")
	}
}

func TestValidateUndefinedVarInEnv(t *testing.T) {
	r := &Recipe{
		Meta: Meta{Name: "test"},
		Services: map[string]Service{
			"app": {
				Image:       "nginx",
				Environment: map[string]string{"KEY": "${MISSING}"},
			},
		},
	}
	if err := r.Validate(); err == nil {
		t.Error("expected error for undefined var in env")
	}
}

func TestExtractVarRefs(t *testing.T) {
	tests := []struct {
		input string
		want  int
	}{
		{"${PORT}:80", 1},
		{"${HOST}:${PORT}", 2},
		{"no vars here", 0},
		{"${A}${B}${C}", 3},
		{"${}", 0},
		{"${UNCLOSED", 0},
		{"plain", 0},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			refs := extractVarRefs(tt.input)
			if len(refs) != tt.want {
				t.Errorf("extractVarRefs(%q) = %d refs, want %d", tt.input, len(refs), tt.want)
			}
		})
	}
}
