// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package project

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/amrelsagaei/updock/internal/security"
)

// DeriveProjectName extracts a project name from a Docker image reference.
func DeriveProjectName(image string) string {
	name := image
	if idx := strings.LastIndex(image, "/"); idx >= 0 {
		name = image[idx+1:]
	}
	name = strings.Split(name, ":")[0]
	name = strings.Split(name, "@")[0]
	return name
}

// UniqueProjectName returns a name that doesn't collide with existing
// directories, appending -2, -3, etc. as needed.
func UniqueProjectName(root, base string) string {
	candidate := base
	for i := 2; ; i++ {
		path := filepath.Join(root, candidate)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return candidate
		}
		candidate = fmt.Sprintf("%s-%d", base, i)
	}
}

// CreateProjectDir creates the project directory structure.
func CreateProjectDir(root, name string) (string, error) {
	if err := security.ValidateProjectName(name); err != nil {
		return "", err
	}

	projectPath := filepath.Join(root, name)
	dataPath := filepath.Join(projectPath, "data")

	if err := os.MkdirAll(dataPath, 0o755); err != nil {
		return "", fmt.Errorf("creating project directory: %w", err)
	}

	return projectPath, nil
}
