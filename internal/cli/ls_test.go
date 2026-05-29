// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amrelsagaei/updock/internal/project"
	"github.com/amrelsagaei/updock/internal/scaffold"
)

// writeProjectsRootConfig writes a config.toml under XDG_CONFIG_HOME pointing
// projects_root at the given directory. Requires XDG_CONFIG_HOME to be set.
func writeProjectsRootConfig(t *testing.T, root string) {
	t.Helper()
	dir := filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "updock")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	content := "projects_root = \"" + root + "\"\n"
	if err := os.WriteFile(filepath.Join(dir, "config.toml"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestLsShowsLiveState(t *testing.T) {
	root := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())

	// Point the config's projects_root at our temp dir via a config file.
	writeProjectsRootConfig(t, root)

	// Scaffold one project.
	projectPath, err := project.CreateProjectDir(root, "nginx")
	if err != nil {
		t.Fatal(err)
	}
	cfg := &project.Config{Image: "nginx", Tag: "latest", ProjectName: "nginx", Ports: []project.PortMapping{{Host: 80, Container: 80}}}
	if err := scaffold.WriteAll(projectPath, cfg); err != nil {
		t.Fatal(err)
	}

	// Inject a fake state detector returning "running" - proves ls uses
	// live detection, not the cached "created" metadata value.
	orig := stateDetector
	t.Cleanup(func() { stateDetector = orig })
	stateDetector = func(_, _ string) string { return "running" }

	cmd := NewRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs([]string{"ls"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("ls failed: %v", err)
	}

	out := buf.String()
	if !strings.Contains(out, "running") {
		t.Errorf("ls should show live 'running' state, got:\n%s", out)
	}
	if strings.Contains(out, "created") {
		t.Error("ls should NOT show the stale cached 'created' state")
	}
}
