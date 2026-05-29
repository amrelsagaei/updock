// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package project

import (
	"testing"
)

func TestNewMetadata(t *testing.T) {
	cfg := Config{
		Image:       "nginx",
		Tag:         "latest",
		Digest:      "sha256:abc123",
		ProjectName: "nginx",
		Ports: []PortMapping{
			{Host: 8080, Container: 80, Protocol: "tcp"},
		},
	}

	meta := NewMetadata(&cfg)

	if meta.Version != 1 {
		t.Errorf("expected version 1, got %d", meta.Version)
	}
	if meta.Image != "nginx" {
		t.Errorf("expected image 'nginx', got %q", meta.Image)
	}
	if meta.Tag != "latest" {
		t.Errorf("expected tag 'latest', got %q", meta.Tag)
	}
	if meta.State != "created" {
		t.Errorf("expected state 'created', got %q", meta.State)
	}
	if len(meta.Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(meta.Ports))
	}
	if meta.Ports[0].Host != 8080 || meta.Ports[0].Container != 80 {
		t.Errorf("unexpected port: %+v", meta.Ports[0])
	}
	if meta.CreatedAt.IsZero() {
		t.Error("CreatedAt should not be zero")
	}
}

func TestWriteAndReadMetadata(t *testing.T) {
	dir := t.TempDir()

	pgCfg := Config{
		Image:       "postgres",
		Tag:         "16",
		ProjectName: "postgres-dev",
		Ports: []PortMapping{
			{Host: 5432, Container: 5432},
		},
	}
	original := NewMetadata(&pgCfg)

	if err := WriteMetadata(dir, &original); err != nil {
		t.Fatalf("WriteMetadata error: %v", err)
	}

	loaded, err := ReadMetadata(dir)
	if err != nil {
		t.Fatalf("ReadMetadata error: %v", err)
	}

	if loaded.Image != original.Image {
		t.Errorf("image mismatch: got %q, want %q", loaded.Image, original.Image)
	}
	if loaded.Tag != original.Tag {
		t.Errorf("tag mismatch: got %q, want %q", loaded.Tag, original.Tag)
	}
	if loaded.ProjectName != original.ProjectName {
		t.Errorf("project name mismatch: got %q", loaded.ProjectName)
	}
	if loaded.State != "created" {
		t.Errorf("state should be 'created', got %q", loaded.State)
	}
	if len(loaded.Ports) != 1 {
		t.Fatalf("expected 1 port, got %d", len(loaded.Ports))
	}
}

func TestReadMetadataNotFound(t *testing.T) {
	_, err := ReadMetadata(t.TempDir())
	if err == nil {
		t.Error("expected error for missing updock.json")
	}
}
