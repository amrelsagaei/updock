// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package project

import (
	"path/filepath"
	"testing"
)

func setupTestRegistry(t *testing.T) (string, *Registry) {
	t.Helper()
	root := t.TempDir()

	projects := []Config{
		{Image: "nginx", Tag: "latest", ProjectName: "nginx", Ports: []PortMapping{{Host: 80, Container: 80}}},
		{Image: "postgres", Tag: "16", ProjectName: "postgres-dev", Ports: []PortMapping{{Host: 5432, Container: 5432}}},
		{Image: "bkimminich/juice-shop", Tag: "latest", ProjectName: "juice-shop"},
	}

	for i := range projects {
		cfg := &projects[i]
		path, err := CreateProjectDir(root, cfg.ProjectName)
		if err != nil {
			t.Fatal(err)
		}
		meta := NewMetadata(cfg)
		if err := WriteMetadata(path, &meta); err != nil {
			t.Fatal(err)
		}
	}

	reg, err := NewRegistry(root)
	if err != nil {
		t.Fatal(err)
	}
	return root, reg
}

func TestRegistryList(t *testing.T) {
	_, reg := setupTestRegistry(t)

	entries := reg.List()
	if len(entries) != 3 {
		t.Fatalf("expected 3 projects, got %d", len(entries))
	}

	if entries[0].Name != "juice-shop" {
		t.Errorf("first project should be 'juice-shop' (alphabetical), got %q", entries[0].Name)
	}
	if entries[0].Number != 1 {
		t.Errorf("first project should be #1, got %d", entries[0].Number)
	}
	if entries[1].Name != "nginx" {
		t.Errorf("second should be 'nginx', got %q", entries[1].Name)
	}
	if entries[2].Name != "postgres-dev" {
		t.Errorf("third should be 'postgres-dev', got %q", entries[2].Name)
	}
}

func TestRegistryResolve(t *testing.T) {
	_, reg := setupTestRegistry(t)

	entry, err := reg.Resolve(2)
	if err != nil {
		t.Fatalf("Resolve(2) error: %v", err)
	}
	if entry.Name != "nginx" {
		t.Errorf("expected 'nginx', got %q", entry.Name)
	}
	if entry.Metadata.Image != "nginx" {
		t.Errorf("expected image 'nginx', got %q", entry.Metadata.Image)
	}
}

func TestRegistryResolveNotFound(t *testing.T) {
	_, reg := setupTestRegistry(t)

	_, err := reg.Resolve(99)
	if err == nil {
		t.Error("expected error for nonexistent number")
	}
}

func TestRegistryResolveByName(t *testing.T) {
	_, reg := setupTestRegistry(t)

	entry, err := reg.ResolveByName("postgres-dev")
	if err != nil {
		t.Fatalf("ResolveByName error: %v", err)
	}
	if entry.Metadata.Tag != "16" {
		t.Errorf("expected tag '16', got %q", entry.Metadata.Tag)
	}
}

func TestRegistryResolveByNameNotFound(t *testing.T) {
	_, reg := setupTestRegistry(t)

	_, err := reg.ResolveByName("nonexistent")
	if err == nil {
		t.Error("expected error")
	}
}

func TestRegistryEmptyRoot(t *testing.T) {
	root := t.TempDir()
	reg, err := NewRegistry(root)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(reg.List()) != 0 {
		t.Error("expected empty list")
	}
}

func TestRegistryNonexistentRoot(t *testing.T) {
	reg, err := NewRegistry(filepath.Join(t.TempDir(), "nonexistent"))
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if len(reg.List()) != 0 {
		t.Error("expected empty list for nonexistent root")
	}
}

func TestRegistryIgnoresNonProjectDirs(t *testing.T) {
	root := t.TempDir()

	if _, err := CreateProjectDir(root, "valid-project"); err != nil {
		t.Fatal(err)
	}
	validCfg := Config{Image: "nginx", Tag: "latest", ProjectName: "valid-project"}
	meta := NewMetadata(&validCfg)
	if err := WriteMetadata(filepath.Join(root, "valid-project"), &meta); err != nil {
		t.Fatal(err)
	}

	if _, err := CreateProjectDir(root, "no-metadata"); err != nil {
		t.Fatal(err)
	}

	reg, err := NewRegistry(root)
	if err != nil {
		t.Fatal(err)
	}

	entries := reg.List()
	if len(entries) != 1 {
		t.Errorf("expected 1 project (ignoring dir without updock.json), got %d", len(entries))
	}
}
