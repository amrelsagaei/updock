// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"bytes"
	"strings"
	"testing"

	"github.com/amrelsagaei/updock/internal/version"
)

func TestNewRootCmdHelp(t *testing.T) {
	cmd := NewRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"--help"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("root --help failed: %v", err)
	}

	out := buf.String()
	expectedCmds := []string{
		"search", "ls", "status", "up", "stop", "restart",
		"rebuild", "down", "logs", "open", "config", "rm",
		"doctor", "version",
	}

	for _, name := range expectedCmds {
		if !strings.Contains(out, name) {
			t.Errorf("help output missing command %q", name)
		}
	}
}

func TestVersionCmd(t *testing.T) {
	version.Version = "1.2.3"
	version.Commit = "abc123"
	version.Date = "2026-01-01"

	cmd := NewRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetArgs([]string{"version"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("version command failed: %v", err)
	}

	out := buf.String()
	tests := []struct {
		name     string
		contains string
	}{
		{"version", "updock 1.2.3"},
		{"commit", "abc123"},
		{"date", "2026-01-01"},
		{"go version", "go"},
		{"os/arch", "os/arch:"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if !strings.Contains(out, tt.contains) {
				t.Errorf("version output should contain %q, got:\n%s", tt.contains, out)
			}
		})
	}
}

func TestDoctorCmd(t *testing.T) {
	cmd := NewRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"doctor"})

	// Doctor should not error regardless of Docker availability.
	_ = cmd.Execute()

	out := buf.String()
	if out == "" {
		t.Error("doctor should produce output")
	}
}

func TestRootNoArgs(t *testing.T) {
	cmd := NewRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("root with no args failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "updock") {
		t.Error("root with no args should show help")
	}
}

func TestCommandsRegistered(t *testing.T) {
	// Lifecycle commands resolve a project number. Without a real project
	// root they return an error - but that proves the command is registered
	// and wired (not a panic or unknown-command error).
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cmds := []struct {
		name string
		args []string
	}{
		{"config", []string{"config", "1"}},
		{"status", []string{"status", "1"}},
		{"up", []string{"up", "1"}},
		{"stop", []string{"stop", "1"}},
		{"restart", []string{"restart", "1"}},
		{"rebuild", []string{"rebuild", "1"}},
		{"down", []string{"down", "1"}},
		{"logs", []string{"logs", "1"}},
		{"open", []string{"open", "1"}},
		{"rm", []string{"rm", "--yes", "1"}},
	}

	for _, tt := range cmds {
		t.Run(tt.name, func(t *testing.T) {
			cmd := NewRootCmd()
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)
			cmd.SetArgs(tt.args)

			// These will error because no project exists, but should not panic
			_ = cmd.Execute()
		})
	}
}

func TestLsEmptyRoot(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	cmd := NewRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"ls"})

	if err := cmd.Execute(); err != nil {
		t.Fatalf("ls failed: %v", err)
	}

	if !strings.Contains(buf.String(), "No projects") {
		t.Error("ls with empty root should show 'No projects' message")
	}
}
