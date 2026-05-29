// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package security

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ValidationError describes an input that failed validation.
type ValidationError struct {
	Field   string
	Value   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid %s: %s", e.Field, e.Message)
}

const (
	maxImageNameLen = 255
	maxTagLen       = 128
	maxProjectLen   = 128
	maxEnvNameLen   = 256
	maxEnvValueLen  = 32768
	minPort         = 1
	maxPort         = 65535
)

var (
	imageNameRe  = regexp.MustCompile(`^[a-z0-9]+(?:[._-][a-z0-9]+)*(?:/[a-z0-9]+(?:[._-][a-z0-9]+)*)*$`)
	tagRe        = regexp.MustCompile(`^[a-zA-Z0-9_][a-zA-Z0-9_.\-]{0,127}$`)
	projectRe    = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9._\-]*$`)
	envVarNameRe = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)
)

func sanitize(s string, maxLen int) string {
	if len(s) > maxLen {
		s = s[:maxLen]
	}
	var b strings.Builder
	for _, r := range s {
		if unicode.IsPrint(r) {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ValidateImageName checks a Docker image name against the allowed pattern.
func ValidateImageName(name string) error {
	if name == "" {
		return &ValidationError{Field: "image name", Value: "", Message: "cannot be empty"}
	}
	if len(name) > maxImageNameLen {
		return &ValidationError{
			Field:   "image name",
			Value:   sanitize(name, 50),
			Message: fmt.Sprintf("exceeds maximum length of %d characters", maxImageNameLen),
		}
	}
	if !imageNameRe.MatchString(name) {
		return &ValidationError{
			Field:   "image name",
			Value:   sanitize(name, 50),
			Message: "must contain only lowercase letters, digits, and separators (., -, _)",
		}
	}
	return nil
}

// ValidateTag checks a Docker image tag against the allowed pattern.
func ValidateTag(tag string) error {
	if tag == "" {
		return &ValidationError{Field: "tag", Value: "", Message: "cannot be empty"}
	}
	if len(tag) > maxTagLen {
		return &ValidationError{
			Field:   "tag",
			Value:   sanitize(tag, 50),
			Message: fmt.Sprintf("exceeds maximum length of %d characters", maxTagLen),
		}
	}
	if !tagRe.MatchString(tag) {
		return &ValidationError{
			Field:   "tag",
			Value:   sanitize(tag, 50),
			Message: "must contain only letters, digits, underscores, dots, and hyphens",
		}
	}
	return nil
}

// ValidateProjectName checks a project directory name against the allowed pattern.
func ValidateProjectName(name string) error {
	if name == "" {
		return &ValidationError{Field: "project name", Value: "", Message: "cannot be empty"}
	}
	if len(name) > maxProjectLen {
		return &ValidationError{
			Field:   "project name",
			Value:   sanitize(name, 50),
			Message: fmt.Sprintf("exceeds maximum length of %d characters", maxProjectLen),
		}
	}
	if !projectRe.MatchString(name) {
		return &ValidationError{
			Field:   "project name",
			Value:   sanitize(name, 50),
			Message: "must start with a letter or digit and contain only letters, digits, dots, hyphens, and underscores",
		}
	}
	return nil
}

// ValidatePort checks that a port number is within the valid TCP/UDP range.
func ValidatePort(port int) error {
	if port < minPort || port > maxPort {
		return &ValidationError{
			Field:   "port",
			Value:   fmt.Sprintf("%d", port),
			Message: fmt.Sprintf("must be between %d and %d", minPort, maxPort),
		}
	}
	return nil
}

// ValidateEnvVarName checks an environment variable name against POSIX rules.
func ValidateEnvVarName(name string) error {
	if name == "" {
		return &ValidationError{Field: "env var name", Value: "", Message: "cannot be empty"}
	}
	if len(name) > maxEnvNameLen {
		return &ValidationError{
			Field:   "env var name",
			Value:   sanitize(name, 50),
			Message: fmt.Sprintf("exceeds maximum length of %d characters", maxEnvNameLen),
		}
	}
	if !envVarNameRe.MatchString(name) {
		return &ValidationError{
			Field:   "env var name",
			Value:   sanitize(name, 50),
			Message: "must start with a letter or underscore and contain only letters, digits, and underscores",
		}
	}
	return nil
}

// ValidateEnvVarValue rejects null bytes and disallowed control characters.
func ValidateEnvVarValue(value string) error {
	if len(value) > maxEnvValueLen {
		return &ValidationError{
			Field:   "env var value",
			Value:   sanitize(value, 50),
			Message: fmt.Sprintf("exceeds maximum length of %d characters", maxEnvValueLen),
		}
	}
	for i, r := range value {
		if r == 0 {
			return &ValidationError{
				Field:   "env var value",
				Value:   sanitize(value, 50),
				Message: fmt.Sprintf("contains null byte at position %d", i),
			}
		}
		if unicode.IsControl(r) && r != '\n' && r != '\t' {
			return &ValidationError{
				Field:   "env var value",
				Value:   sanitize(value, 50),
				Message: fmt.Sprintf("contains disallowed control character at position %d", i),
			}
		}
	}
	return nil
}
