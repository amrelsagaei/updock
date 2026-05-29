// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package project

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

// Entry represents a discovered project in the registry.
type Entry struct {
	Number   int
	Name     string
	Path     string
	Metadata Metadata
}

// Registry scans the projects root and assigns stable numbers.
type Registry struct {
	Root    string
	entries []Entry
}

// NewRegistry scans the root directory for projects.
func NewRegistry(root string) (*Registry, error) {
	r := &Registry{Root: root}
	if err := r.scan(); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *Registry) scan() error {
	entries, err := os.ReadDir(r.Root)
	if err != nil {
		if os.IsNotExist(err) {
			r.entries = nil
			return nil
		}
		return fmt.Errorf("scanning projects root: %w", err)
	}

	var projects []Entry
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		projectPath := filepath.Join(r.Root, e.Name())
		meta, err := ReadMetadata(projectPath)
		if err != nil {
			continue
		}
		projects = append(projects, Entry{
			Name:     e.Name(),
			Path:     projectPath,
			Metadata: meta,
		})
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].Name < projects[j].Name
	})

	for i := range projects {
		projects[i].Number = i + 1
	}

	r.entries = projects
	return nil
}

// List returns all discovered projects.
func (r *Registry) List() []Entry {
	return r.entries
}

// Resolve maps a number to a project entry.
func (r *Registry) Resolve(number int) (Entry, error) {
	for i := range r.entries {
		if r.entries[i].Number == number {
			return r.entries[i], nil
		}
	}
	return Entry{}, fmt.Errorf("no project with number %d (run 'updock ls' to see available projects)", number)
}

// ResolveByName finds a project by name.
func (r *Registry) ResolveByName(name string) (Entry, error) {
	for i := range r.entries {
		if r.entries[i].Name == name {
			return r.entries[i], nil
		}
	}
	return Entry{}, fmt.Errorf("no project named %q", name)
}
