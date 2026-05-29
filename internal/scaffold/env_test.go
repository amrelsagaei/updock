// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amrelsagaei/updock/internal/project"
)

func TestWriteEnvFile(t *testing.T) {
	dir := t.TempDir()
	cfg := project.Config{
		Env: []project.EnvVar{
			{Key: "POSTGRES_PASSWORD", Value: "mysecret", Secret: true},
			{Key: "POSTGRES_USER", Value: "postgres"},
		},
		Ports: []project.PortMapping{
			{Host: 5432, Container: 5432},
		},
	}

	if err := WriteEnvFile(dir, &cfg); err != nil {
		t.Fatalf("WriteEnvFile error: %v", err)
	}

	path := filepath.Join(dir, ".env")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	if !strings.Contains(content, "POSTGRES_PASSWORD=mysecret") {
		t.Error("env file should contain password value")
	}
	if !strings.Contains(content, "POSTGRES_USER=postgres") {
		t.Error("env file should contain user value")
	}
	if !strings.Contains(content, "HOST_PORT_5432=5432") {
		t.Error("env file should contain port variable")
	}
}

func TestEnvFilePermissions(t *testing.T) {
	dir := t.TempDir()
	cfg := project.Config{
		Env: []project.EnvVar{
			{Key: "SECRET", Value: "val", Secret: true},
		},
	}

	if err := WriteEnvFile(dir, &cfg); err != nil {
		t.Fatal(err)
	}

	path := filepath.Join(dir, ".env")
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}

	perm := info.Mode().Perm()
	if perm != 0o600 {
		t.Errorf("expected permissions 0600, got %04o", perm)
	}
}

func TestWriteGitignore(t *testing.T) {
	dir := t.TempDir()

	if err := WriteGitignore(dir); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, ".gitignore"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	if !strings.Contains(content, ".env") {
		t.Error("gitignore should contain .env")
	}
	if !strings.Contains(content, "data/") {
		t.Error("gitignore should contain data/")
	}
}
