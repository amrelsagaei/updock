// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package security

import (
	"strings"
	"testing"
)

func TestValidateImageName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"simple", "nginx", false},
		{"with namespace", "library/nginx", false},
		{"with dots", "my.registry/my.image", false},
		{"with hyphens", "bkimminich/juice-shop", false},
		{"with underscores", "my_image", false},
		{"empty", "", true},
		{"uppercase", "Nginx", true},
		{"with colon", "nginx:latest", true},
		{"with spaces", "my image", true},
		{"shell injection semicolon", "nginx;rm -rf /", true},
		{"shell injection pipe", "nginx|cat /etc/passwd", true},
		{"shell injection backtick", "nginx`whoami`", true},
		{"starts with dot", ".nginx", true},
		{"starts with hyphen", "-nginx", true},
		{"double slash", "a//b", true},
		{"too long", strings.Repeat("a", 256), true},
		{"max length", strings.Repeat("a", 255), false},
		{"unicode", "nginx🐋", true},
		{"deep namespace", "a/b/c/d", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateImageName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateImageName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
			if err != nil {
				var ve *ValidationError
				if ok := errorAs(err, &ve); !ok {
					t.Errorf("expected *ValidationError, got %T", err)
				}
			}
		})
	}
}

func TestValidateTag(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"latest", "latest", false},
		{"semver", "1.10.0", false},
		{"with hyphen", "v1.0-alpine", false},
		{"with underscore", "my_tag", false},
		{"empty", "", true},
		{"with colon", "tag:1", true},
		{"with slash", "tag/1", true},
		{"too long", strings.Repeat("a", 129), true},
		{"max length", strings.Repeat("a", 128), false},
		{"starts with dot", ".hidden", true},
		{"starts with hyphen", "-tag", true},
		{"shell injection", "latest;rm -rf /", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTag(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateTag(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateProjectName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"simple", "myproject", false},
		{"with hyphens", "juice-shop", false},
		{"with dots", "my.project", false},
		{"with digits", "project123", false},
		{"starts with digit", "1project", false},
		{"empty", "", true},
		{"starts with dot", ".hidden", true},
		{"starts with hyphen", "-bad", true},
		{"with spaces", "my project", true},
		{"with slash", "a/b", true},
		{"too long", strings.Repeat("a", 129), true},
		{"shell injection", "proj;rm -rf /", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateProjectName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateProjectName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidatePort(t *testing.T) {
	tests := []struct {
		name    string
		port    int
		wantErr bool
	}{
		{"min valid", 1, false},
		{"common port", 8080, false},
		{"max valid", 65535, false},
		{"zero", 0, true},
		{"negative", -1, true},
		{"too high", 65536, true},
		{"very high", 100000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidatePort(tt.port)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidatePort(%d) error = %v, wantErr %v", tt.port, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEnvVarName(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"simple", "MY_VAR", false},
		{"with digits", "VAR_123", false},
		{"starts with underscore", "_PRIVATE", false},
		{"lowercase", "my_var", false},
		{"empty", "", true},
		{"starts with digit", "1VAR", true},
		{"with spaces", "MY VAR", true},
		{"with equals", "MY=VAR", true},
		{"with hyphen", "MY-VAR", true},
		{"too long", strings.Repeat("A", 257), true},
		{"max length", strings.Repeat("A", 256), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnvVarName(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEnvVarName(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidateEnvVarValue(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"normal", "hello world", false},
		{"with newline", "line1\nline2", false},
		{"with tab", "col1\tcol2", false},
		{"empty", "", false},
		{"with null", "before\x00after", true},
		{"with bell", "ding\x07dong", true},
		{"with escape", "esc\x1bcode", true},
		{"too long", strings.Repeat("a", 32769), true},
		{"max length", strings.Repeat("a", 32768), false},
		{"unicode", "日本語", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateEnvVarValue(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateEnvVarValue(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
			}
		})
	}
}

func TestValidationErrorMessage(t *testing.T) {
	err := &ValidationError{
		Field:   "image name",
		Value:   "bad;input",
		Message: "contains invalid characters",
	}
	got := err.Error()
	if !strings.Contains(got, "image name") {
		t.Errorf("error message should contain field name, got %q", got)
	}
	if !strings.Contains(got, "contains invalid characters") {
		t.Errorf("error message should contain reason, got %q", got)
	}
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{"short", "hello", 10, "hello"},
		{"truncate", "hello world", 5, "hello"},
		{"control chars", "a\x00b\x07c", 10, "abc"},
		{"empty", "", 10, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sanitize(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("sanitize(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

// errorAs is a simple helper to avoid importing errors package in tests.
func errorAs(err error, target interface{}) bool {
	ve, ok := err.(*ValidationError)
	if ok {
		*target.(**ValidationError) = ve
	}
	return ok
}
