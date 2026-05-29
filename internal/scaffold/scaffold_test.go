// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package scaffold

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/amrelsagaei/updock/internal/project"
)

func TestWriteAll(t *testing.T) {
	dir := t.TempDir()
	cfg := project.Config{
		Image:       "nginx",
		Tag:         "latest",
		ProjectName: "nginx",
		Ports: []project.PortMapping{
			{Host: 80, Container: 80, Protocol: "tcp"},
		},
		Env: []project.EnvVar{
			{Key: "PORT", Value: "80"},
		},
	}

	if err := WriteAll(dir, &cfg); err != nil {
		t.Fatalf("WriteAll error: %v", err)
	}

	expectedFiles := []string{
		"docker-compose.yml",
		".env",
		".gitignore",
		"updock.json",
	}

	for _, name := range expectedFiles {
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("expected file %q to exist", name)
		}
	}
}

func TestWriteAllRoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfg := project.Config{
		Image:       "postgres",
		Tag:         "16",
		Digest:      "sha256:abc123",
		ProjectName: "postgres-dev",
		Ports: []project.PortMapping{
			{Host: 5432, Container: 5432, Protocol: "tcp"},
		},
		Env: []project.EnvVar{
			{Key: "POSTGRES_PASSWORD", Value: "test123", Secret: true, Required: true},
			{Key: "POSTGRES_USER", Value: "admin"},
		},
	}

	if err := WriteAll(dir, &cfg); err != nil {
		t.Fatal(err)
	}

	meta, err := project.ReadMetadata(dir)
	if err != nil {
		t.Fatalf("ReadMetadata error: %v", err)
	}

	if meta.Image != cfg.Image {
		t.Errorf("round-trip image: got %q, want %q", meta.Image, cfg.Image)
	}
	if meta.Tag != cfg.Tag {
		t.Errorf("round-trip tag: got %q, want %q", meta.Tag, cfg.Tag)
	}
	if meta.Digest != cfg.Digest {
		t.Errorf("round-trip digest: got %q, want %q", meta.Digest, cfg.Digest)
	}
	if meta.ProjectName != cfg.ProjectName {
		t.Errorf("round-trip name: got %q, want %q", meta.ProjectName, cfg.ProjectName)
	}
	if len(meta.Ports) != 1 || meta.Ports[0].Host != 5432 {
		t.Errorf("round-trip ports: got %+v", meta.Ports)
	}
}
