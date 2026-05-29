// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amrelsagaei/updock/internal/project"
	"github.com/amrelsagaei/updock/internal/prompt"
	"github.com/amrelsagaei/updock/internal/scaffold"
)

func setupConfigProject(t *testing.T) (root, projectPath string) {
	t.Helper()
	root = t.TempDir()
	projectPath = filepath.Join(root, "myapp")
	if err := os.MkdirAll(filepath.Join(projectPath, "data"), 0o755); err != nil {
		t.Fatal(err)
	}

	cfg := &project.Config{
		Image:       "postgres",
		Tag:         "16",
		ProjectName: "myapp",
		Ports:       []project.PortMapping{{Host: 5432, Container: 5432, Protocol: "tcp"}},
		Env: []project.EnvVar{
			{Key: "POSTGRES_PASSWORD", Value: "oldpass", Secret: true, Required: true},
			{Key: "POSTGRES_USER", Value: "admin"},
		},
	}
	if err := scaffold.WriteAll(projectPath, cfg); err != nil {
		t.Fatal(err)
	}
	return root, projectPath
}

func TestRunConfigEditsEnv(t *testing.T) {
	_, projectPath := setupConfigProject(t)

	origIn := prompt.InputFunc
	t.Cleanup(func() { prompt.InputFunc = origIn })
	prompt.InputFunc = func(_, _, value string) (string, error) {
		if value == "oldpass" {
			return "newpass", nil
		}
		return value, nil
	}

	cfg, err := loadProjectConfig(projectPath)
	if err != nil {
		t.Fatal(err)
	}

	cfg.Env, err = prompt.EditEnvVars(cfg.Env)
	if err != nil {
		t.Fatal(err)
	}
	if err := scaffold.WriteEnvFile(projectPath, cfg); err != nil {
		t.Fatal(err)
	}

	data, _ := os.ReadFile(filepath.Join(projectPath, ".env"))
	content := string(data)
	if !strings.Contains(content, "POSTGRES_PASSWORD=newpass") {
		t.Error("env should be updated to new password")
	}
	if strings.Contains(content, "oldpass") {
		t.Error("old password should be replaced")
	}
}

func TestLoadProjectConfigFiltersHostPort(t *testing.T) {
	_, projectPath := setupConfigProject(t)

	cfg, err := loadProjectConfig(projectPath)
	if err != nil {
		t.Fatal(err)
	}

	for _, e := range cfg.Env {
		if strings.HasPrefix(e.Key, "HOST_PORT_") {
			t.Errorf("HOST_PORT_* should be filtered from env vars, found %q", e.Key)
		}
	}

	// ports should come from metadata, not env
	if len(cfg.Ports) != 1 || cfg.Ports[0].Host != 5432 {
		t.Errorf("expected 1 port 5432, got %+v", cfg.Ports)
	}
}

func TestEditEnvVars(t *testing.T) {
	origIn := prompt.InputFunc
	t.Cleanup(func() { prompt.InputFunc = origIn })
	prompt.InputFunc = func(_, _, value string) (string, error) {
		return value + "-edited", nil
	}

	vars := []project.EnvVar{
		{Key: "FOO", Value: "bar"},
		{Key: "SECRET_KEY", Value: "xyz", Secret: true},
	}

	result, err := prompt.EditEnvVars(vars)
	if err != nil {
		t.Fatal(err)
	}
	if result[0].Value != "bar-edited" {
		t.Errorf("expected 'bar-edited', got %q", result[0].Value)
	}
	if result[1].Value != "xyz-edited" {
		t.Errorf("expected 'xyz-edited', got %q", result[1].Value)
	}
}

func TestUpdatePortMetadata(t *testing.T) {
	_, projectPath := setupConfigProject(t)

	cfg := &project.Config{
		Ports: []project.PortMapping{{Host: 9999, Container: 5432}},
	}
	if err := updatePortMetadata(projectPath, cfg); err != nil {
		t.Fatal(err)
	}

	meta, err := project.ReadMetadata(projectPath)
	if err != nil {
		t.Fatal(err)
	}
	if len(meta.Ports) != 1 || meta.Ports[0].Host != 9999 {
		t.Errorf("metadata ports should be updated to 9999, got %+v", meta.Ports)
	}
}
