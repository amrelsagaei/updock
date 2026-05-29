// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/project"
	"github.com/amrelsagaei/updock/internal/prompt"
	"github.com/amrelsagaei/updock/internal/scaffold"
)

// fakeRunner implements docker.Runner without touching a real daemon.
type fakeRunner struct {
	runErr    error
	output    []byte
	outputErr error
}

func (f *fakeRunner) Run(_ string, _, _ io.Writer, _ string, _ ...string) error {
	return f.runErr
}

func (f *fakeRunner) Output(_, _ string, _ ...string) ([]byte, error) {
	return f.output, f.outputErr
}

// setupLifecycleProject points config at a temp root with one scaffolded
// project named "nginx" (number 1), and installs a fake docker runner.
func setupLifecycleProject(t *testing.T) {
	t.Helper()
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	root := t.TempDir()
	writeProjectsRootConfig(t, root)

	projectPath, err := project.CreateProjectDir(root, "nginx")
	if err != nil {
		t.Fatal(err)
	}
	cfg := &project.Config{
		Image: "nginx", Tag: "latest", ProjectName: "nginx",
		Ports: []project.PortMapping{{Host: 8080, Container: 80}},
		Env:   []project.EnvVar{{Key: "FOO", Value: "bar"}},
	}
	if err := scaffold.WriteAll(projectPath, cfg); err != nil {
		t.Fatal(err)
	}

	origRunner := docker.CommandRunner
	docker.CommandRunner = &fakeRunner{output: []byte("running\n")}
	t.Cleanup(func() { docker.CommandRunner = origRunner })

	// Skip the real preflight (no daemon in CI) for the root-level up path.
	// Lifecycle subcommands don't call preflight, so most tests don't need this.
}

func runCmd(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := NewRootCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestCmdStop(t *testing.T) {
	setupLifecycleProject(t)
	out, err := runCmd(t, "stop", "1")
	if err != nil {
		t.Fatalf("stop failed: %v", err)
	}
	if !strings.Contains(out, "stopped") {
		t.Errorf("expected confirmation, got: %s", out)
	}
}

func TestCmdRestart(t *testing.T) {
	setupLifecycleProject(t)
	out, err := runCmd(t, "restart", "1")
	if err != nil {
		t.Fatalf("restart failed: %v", err)
	}
	if !strings.Contains(out, "restarted") {
		t.Errorf("expected confirmation, got: %s", out)
	}
}

func TestCmdDown(t *testing.T) {
	setupLifecycleProject(t)
	out, err := runCmd(t, "down", "1")
	if err != nil {
		t.Fatalf("down failed: %v", err)
	}
	if !strings.Contains(out, "removed") {
		t.Errorf("expected confirmation, got: %s", out)
	}
}

func TestCmdRebuild(t *testing.T) {
	setupLifecycleProject(t)
	out, err := runCmd(t, "rebuild", "1")
	if err != nil {
		t.Fatalf("rebuild failed: %v", err)
	}
	if !strings.Contains(out, "rebuilt") {
		t.Errorf("expected confirmation, got: %s", out)
	}
}

func TestCmdLogs(t *testing.T) {
	setupLifecycleProject(t)
	if _, err := runCmd(t, "logs", "1"); err != nil {
		t.Fatalf("logs failed: %v", err)
	}
}

func TestCmdStatus(t *testing.T) {
	setupLifecycleProject(t)
	out, err := runCmd(t, "status", "1")
	if err != nil {
		t.Fatalf("status failed: %v", err)
	}
	if !strings.Contains(out, "nginx") || !strings.Contains(out, "Image:") {
		t.Errorf("status should show project details, got: %s", out)
	}
}

func TestCmdStatusMasksSecrets(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	root := t.TempDir()
	writeProjectsRootConfig(t, root)

	projectPath, err := project.CreateProjectDir(root, "db")
	if err != nil {
		t.Fatal(err)
	}
	cfg := &project.Config{
		Image: "postgres", Tag: "16", ProjectName: "db",
		Env: []project.EnvVar{{Key: "POSTGRES_PASSWORD", Value: "topsecret123", Secret: true}},
	}
	if err := scaffold.WriteAll(projectPath, cfg); err != nil {
		t.Fatal(err)
	}
	origRunner := docker.CommandRunner
	docker.CommandRunner = &fakeRunner{output: []byte("running\n")}
	t.Cleanup(func() { docker.CommandRunner = origRunner })

	out, err := runCmd(t, "status", "1")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(out, "topsecret123") {
		t.Error("status must mask secret values")
	}
	if !strings.Contains(out, "••••••") {
		t.Error("status should show masked placeholder")
	}
}

func TestCmdOpenNoPorts(t *testing.T) {
	t.Setenv("XDG_CONFIG_HOME", t.TempDir())
	root := t.TempDir()
	writeProjectsRootConfig(t, root)

	projectPath, err := project.CreateProjectDir(root, "noports")
	if err != nil {
		t.Fatal(err)
	}
	cfg := &project.Config{Image: "busybox", Tag: "latest", ProjectName: "noports"}
	if err := scaffold.WriteAll(projectPath, cfg); err != nil {
		t.Fatal(err)
	}

	_, err = runCmd(t, "open", "1")
	if err == nil {
		t.Error("open should error when project has no mapped ports")
	}
}

func TestCmdRmWithYes(t *testing.T) {
	setupLifecycleProject(t)

	out, err := runCmd(t, "rm", "--yes", "1")
	if err != nil {
		t.Fatalf("rm failed: %v", err)
	}
	if !strings.Contains(out, "deleted") {
		t.Errorf("expected deletion confirmation, got: %s", out)
	}
}

func TestCmdRmCancelled(t *testing.T) {
	setupLifecycleProject(t)

	origConfirm := prompt.ConfirmFunc
	prompt.ConfirmFunc = func(_ string) (bool, error) { return false, nil }
	t.Cleanup(func() { prompt.ConfirmFunc = origConfirm })

	out, err := runCmd(t, "rm", "1")
	if err != nil {
		t.Fatalf("rm (cancelled) should not error: %v", err)
	}
	if !strings.Contains(out, "Cancelled") {
		t.Errorf("expected cancellation message, got: %s", out)
	}
}

func TestCmdResolveInvalidNumber(t *testing.T) {
	setupLifecycleProject(t)
	_, err := runCmd(t, "stop", "notanumber")
	if err == nil {
		t.Error("non-numeric project arg should error")
	}
}

func TestCmdResolveOutOfRange(t *testing.T) {
	setupLifecycleProject(t)
	_, err := runCmd(t, "stop", "99")
	if err == nil {
		t.Error("out-of-range project number should error")
	}
}
