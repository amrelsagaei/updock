// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package project

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const metadataFile = "updock.json"

// Metadata is the JSON structure stored in updock.json.
type Metadata struct {
	Version     int            `json:"version"`
	Image       string         `json:"image"`
	Tag         string         `json:"tag"`
	Digest      string         `json:"digest,omitempty"`
	ProjectName string         `json:"project_name"`
	Ports       []PortMetadata `json:"ports"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	State       string         `json:"state"`
}

// PortMetadata is the JSON representation of a port mapping.
type PortMetadata struct {
	Host      int `json:"host"`
	Container int `json:"container"`
}

// NewMetadata creates metadata from a project Config.
func NewMetadata(cfg *Config) Metadata {
	now := time.Now().UTC()
	ports := make([]PortMetadata, len(cfg.Ports))
	for i, p := range cfg.Ports {
		ports[i] = PortMetadata{Host: p.Host, Container: p.Container}
	}

	return Metadata{
		Version:     1,
		Image:       cfg.Image,
		Tag:         cfg.Tag,
		Digest:      cfg.Digest,
		ProjectName: cfg.ProjectName,
		Ports:       ports,
		CreatedAt:   now,
		UpdatedAt:   now,
		State:       "created",
	}
}

// WriteMetadata writes updock.json to the project directory.
func WriteMetadata(projectPath string, meta *Metadata) error {
	data, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling metadata: %w", err)
	}

	path := filepath.Join(projectPath, metadataFile)
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("writing metadata: %w", err)
	}
	return nil
}

// ReadMetadata reads updock.json from a project directory.
func ReadMetadata(projectPath string) (Metadata, error) {
	path := filepath.Join(projectPath, metadataFile)
	data, err := os.ReadFile(path)
	if err != nil {
		return Metadata{}, fmt.Errorf("reading metadata: %w", err)
	}

	var meta Metadata
	if err := json.Unmarshal(data, &meta); err != nil {
		return Metadata{}, fmt.Errorf("parsing metadata: %w", err)
	}
	return meta, nil
}
