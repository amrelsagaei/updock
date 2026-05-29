// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.ProjectsRoot == "" {
		t.Error("ProjectsRoot should not be empty")
	}
	if cfg.BrowserCommand == "" {
		t.Error("BrowserCommand should not be empty")
	}
	if !cfg.AutoGeneratePasswords {
		t.Error("AutoGeneratePasswords should default to true")
	}
	if cfg.DefaultRegistry != "docker.io" {
		t.Errorf("DefaultRegistry should be 'docker.io', got %q", cfg.DefaultRegistry)
	}
	if cfg.ScanBeforeRun {
		t.Error("ScanBeforeRun should default to false")
	}
}

func TestLoadFallsBackToDefaults(t *testing.T) {
	// Point config to a nonexistent directory so no file is found.
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	defaults := DefaultConfig()
	if cfg.DefaultRegistry != defaults.DefaultRegistry {
		t.Errorf("expected default registry %q, got %q", defaults.DefaultRegistry, cfg.DefaultRegistry)
	}
	if cfg.AutoGeneratePasswords != defaults.AutoGeneratePasswords {
		t.Errorf("expected auto_generate_passwords=%v, got %v", defaults.AutoGeneratePasswords, cfg.AutoGeneratePasswords)
	}
}

func TestLoadFromTOML(t *testing.T) {
	dir := t.TempDir()
	updockDir := filepath.Join(dir, "updock")
	if err := os.MkdirAll(updockDir, 0o755); err != nil {
		t.Fatal(err)
	}

	toml := `
projects_root = "/tmp/my-projects"
browser_command = "firefox"
auto_generate_passwords = false
default_registry = "ghcr.io"
scan_before_run = true
`
	if err := os.WriteFile(filepath.Join(updockDir, "config.toml"), []byte(toml), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("XDG_CONFIG_HOME", dir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error: %v", err)
	}

	if cfg.ProjectsRoot != "/tmp/my-projects" {
		t.Errorf("ProjectsRoot: got %q", cfg.ProjectsRoot)
	}
	if cfg.BrowserCommand != "firefox" {
		t.Errorf("BrowserCommand: got %q", cfg.BrowserCommand)
	}
	if cfg.AutoGeneratePasswords {
		t.Error("AutoGeneratePasswords should be false")
	}
	if cfg.DefaultRegistry != "ghcr.io" {
		t.Errorf("DefaultRegistry: got %q", cfg.DefaultRegistry)
	}
	if !cfg.ScanBeforeRun {
		t.Error("ScanBeforeRun should be true")
	}
}

func TestValidateRejectsEmptyProjectsRoot(t *testing.T) {
	cfg := DefaultConfig()
	cfg.ProjectsRoot = ""

	err := validate(cfg)
	if err == nil {
		t.Error("expected error for empty projects_root")
	}
}

func TestValidateRejectsEmptyRegistry(t *testing.T) {
	cfg := DefaultConfig()
	cfg.DefaultRegistry = ""

	err := validate(cfg)
	if err == nil {
		t.Error("expected error for empty default_registry")
	}
}

func TestLoadBadTOML(t *testing.T) {
	dir := t.TempDir()
	updockDir := filepath.Join(dir, "updock")
	if err := os.MkdirAll(updockDir, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(updockDir, "config.toml"), []byte("invalid[[[toml"), 0o644); err != nil {
		t.Fatal(err)
	}

	t.Setenv("XDG_CONFIG_HOME", dir)

	_, err := Load()
	if err == nil {
		t.Error("expected error for invalid TOML")
	}
}

func TestConfigDir(t *testing.T) {
	t.Run("with XDG", func(t *testing.T) {
		custom := filepath.Join(t.TempDir(), "config")
		t.Setenv("XDG_CONFIG_HOME", custom)
		dir := configDir()
		want := filepath.Join(custom, "updock")
		if dir != want {
			t.Errorf("got %q, want %q", dir, want)
		}
	})

	t.Run("without XDG", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		dir := configDir()
		if dir == "" {
			t.Error("configDir should not be empty")
		}
	})
}

func TestDefaultBrowser(t *testing.T) {
	b := defaultBrowser()
	if b == "" {
		t.Error("defaultBrowser should not be empty")
	}
}

func TestDefaultProjectsRoot(t *testing.T) {
	root := defaultProjectsRoot()
	if root == "" {
		t.Error("defaultProjectsRoot should not be empty")
	}
}

func TestExpandHome(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot determine home dir")
	}

	tests := []struct {
		name string
		path string
		want string
	}{
		{"tilde path", "~/projects", filepath.Join(home, "projects")},
		{"absolute path", "/tmp/foo", "/tmp/foo"},
		{"relative path", "foo/bar", "foo/bar"},
		{"just tilde slash", "~/", home},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := expandHome(tt.path)
			if got != tt.want {
				t.Errorf("expandHome(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
