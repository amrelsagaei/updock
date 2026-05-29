// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package scaffold

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amrelsagaei/updock/internal/project"
)

func testConfig() project.Config {
	return project.Config{
		Image:       "bkimminich/juice-shop",
		Tag:         "latest",
		ProjectName: "juice-shop",
		Ports: []project.PortMapping{
			{Host: 3000, Container: 3000, Protocol: "tcp"},
		},
		Env: []project.EnvVar{
			{Key: "NODE_ENV", Value: "production"},
			{Key: "ADMIN_PASSWORD", Value: "s3cret!", Secret: true, Required: true},
		},
		Volumes: []project.VolumeMapping{
			{HostPath: "data/app", ContainerPath: "/data"},
		},
	}
}

func TestWriteComposeFile(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig()

	if err := WriteComposeFile(dir, &cfg); err != nil {
		t.Fatalf("WriteComposeFile error: %v", err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml"))
	if err != nil {
		t.Fatalf("reading compose file: %v", err)
	}
	content := string(data)

	if !strings.Contains(content, "bkimminich/juice-shop:latest") {
		t.Error("compose should contain image:tag")
	}
	if !strings.Contains(content, "${HOST_PORT_3000:-3000}:3000") {
		t.Error("compose should reference port via env var")
	}
	if !strings.Contains(content, "env_file") {
		t.Error("compose should use env_file")
	}
	if !strings.Contains(content, "./data/app:/data") {
		t.Error("compose should contain volume mapping")
	}
	if !strings.Contains(content, "restart: unless-stopped") {
		t.Error("compose should set restart policy")
	}
}

func TestComposeNeverContainsSecrets(t *testing.T) {
	dir := t.TempDir()
	cfg := testConfig()

	if err := WriteComposeFile(dir, &cfg); err != nil {
		t.Fatal(err)
	}

	data, err := os.ReadFile(filepath.Join(dir, "docker-compose.yml"))
	if err != nil {
		t.Fatal(err)
	}
	content := string(data)

	for _, env := range cfg.Env {
		if env.Secret && strings.Contains(content, env.Value) {
			t.Errorf("compose file must never contain secret value for %s", env.Key)
		}
	}
}

func TestBuildComposeContentNoPorts(t *testing.T) {
	cfg := project.Config{
		Image: "nginx",
		Tag:   "latest",
	}
	content := buildComposeContent(&cfg)

	if strings.Contains(content, "ports:") {
		t.Error("should not have ports section when no ports defined")
	}
}

func TestBuildComposeContentNoVolumes(t *testing.T) {
	cfg := project.Config{
		Image: "nginx",
		Tag:   "latest",
	}
	content := buildComposeContent(&cfg)

	if strings.Contains(content, "volumes:") {
		t.Error("should not have volumes section when no volumes defined")
	}
}
