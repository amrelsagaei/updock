// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package project

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDeriveProjectName(t *testing.T) {
	tests := []struct {
		image string
		want  string
	}{
		{"nginx", "nginx"},
		{"bkimminich/juice-shop", "juice-shop"},
		{"library/postgres", "postgres"},
		{"nginx:latest", "nginx"},
		{"nginx:1.25@sha256:abc", "nginx"},
		{"registry.example.com/org/app", "app"},
	}

	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			got := DeriveProjectName(tt.image)
			if got != tt.want {
				t.Errorf("DeriveProjectName(%q) = %q, want %q", tt.image, got, tt.want)
			}
		})
	}
}

func TestUniqueProjectName(t *testing.T) {
	root := t.TempDir()

	got := UniqueProjectName(root, "myapp")
	if got != "myapp" {
		t.Errorf("expected 'myapp', got %q", got)
	}

	if err := os.MkdirAll(filepath.Join(root, "myapp"), 0o755); err != nil {
		t.Fatal(err)
	}
	got = UniqueProjectName(root, "myapp")
	if got != "myapp-2" {
		t.Errorf("expected 'myapp-2', got %q", got)
	}

	if err := os.MkdirAll(filepath.Join(root, "myapp-2"), 0o755); err != nil {
		t.Fatal(err)
	}
	got = UniqueProjectName(root, "myapp")
	if got != "myapp-3" {
		t.Errorf("expected 'myapp-3', got %q", got)
	}
}

func TestCreateProjectDir(t *testing.T) {
	root := t.TempDir()

	path, err := CreateProjectDir(root, "test-project")
	if err != nil {
		t.Fatalf("CreateProjectDir error: %v", err)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Error("project directory should exist")
	}
	if _, err := os.Stat(filepath.Join(path, "data")); os.IsNotExist(err) {
		t.Error("data directory should exist")
	}
}

func TestCreateProjectDirInvalidName(t *testing.T) {
	root := t.TempDir()

	_, err := CreateProjectDir(root, ".hidden")
	if err == nil {
		t.Error("expected error for invalid project name")
	}
}
