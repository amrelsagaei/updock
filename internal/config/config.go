// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the global updock settings loaded from TOML.
type Config struct {
	ProjectsRoot          string `mapstructure:"projects_root"`
	BrowserCommand        string `mapstructure:"browser_command"`
	AutoGeneratePasswords bool   `mapstructure:"auto_generate_passwords"`
	DefaultRegistry       string `mapstructure:"default_registry"`
	ScanBeforeRun         bool   `mapstructure:"scan_before_run"`
}

// DefaultConfig returns the configuration with all defaults applied.
func DefaultConfig() Config {
	return Config{
		ProjectsRoot:          defaultProjectsRoot(),
		BrowserCommand:        defaultBrowser(),
		AutoGeneratePasswords: true,
		DefaultRegistry:       "docker.io",
		ScanBeforeRun:         false,
	}
}

// Load reads the config file from disk and merges with defaults.
func Load() (Config, error) {
	cfg := DefaultConfig()

	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(configDir())

	v.SetDefault("projects_root", cfg.ProjectsRoot)
	v.SetDefault("browser_command", cfg.BrowserCommand)
	v.SetDefault("auto_generate_passwords", cfg.AutoGeneratePasswords)
	v.SetDefault("default_registry", cfg.DefaultRegistry)
	v.SetDefault("scan_before_run", cfg.ScanBeforeRun)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return cfg, fmt.Errorf("reading config: %w", err)
		}
	}

	if err := v.Unmarshal(&cfg); err != nil {
		return cfg, fmt.Errorf("parsing config: %w", err)
	}

	if err := validate(cfg); err != nil {
		return cfg, err
	}

	cfg.ProjectsRoot = expandHome(cfg.ProjectsRoot)

	return cfg, nil
}

func validate(cfg Config) error {
	if cfg.ProjectsRoot == "" {
		return fmt.Errorf("config: projects_root cannot be empty")
	}
	if cfg.DefaultRegistry == "" {
		return fmt.Errorf("config: default_registry cannot be empty")
	}
	return nil
}

func configDir() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "updock")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", ".config", "updock")
	}
	return filepath.Join(home, ".config", "updock")
}

func defaultProjectsRoot() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", "updock")
	}
	return filepath.Join(home, "updock")
}

func defaultBrowser() string {
	if runtime.GOOS == "darwin" {
		return "open"
	}
	return "xdg-open"
}

func expandHome(path string) string {
	if !strings.HasPrefix(path, "~/") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return filepath.Join(home, path[2:])
}
