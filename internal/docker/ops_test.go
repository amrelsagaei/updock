// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

type mockRunner struct {
	lastDir  string
	lastName string
	lastArgs []string
	err      error
	output   []byte
}

func (m *mockRunner) Run(dir string, _, _ io.Writer, name string, args ...string) error {
	m.lastDir = dir
	m.lastName = name
	m.lastArgs = args
	return m.err
}

func (m *mockRunner) Output(dir, name string, args ...string) ([]byte, error) {
	m.lastDir = dir
	m.lastName = name
	m.lastArgs = args
	return m.output, m.err
}

func withMockRunner(m *mockRunner, fn func()) {
	orig := CommandRunner
	CommandRunner = m
	defer func() { CommandRunner = orig }()
	fn()
}

func TestUpUsesArgArrays(t *testing.T) {
	m := &mockRunner{}
	withMockRunner(m, func() {
		var stdout, stderr bytes.Buffer
		err := Up("/tmp/project", "myproject", &stdout, &stderr)
		if err != nil {
			t.Fatal(err)
		}

		if m.lastName != "docker" {
			t.Errorf("expected 'docker', got %q", m.lastName)
		}
		args := strings.Join(m.lastArgs, " ")
		if !strings.Contains(args, "compose") {
			t.Error("should use 'compose' subcommand")
		}
		if !strings.Contains(args, "--project-name myproject") {
			t.Error("should set --project-name")
		}
		if !strings.Contains(args, "up -d") {
			t.Error("should pass 'up -d'")
		}
		if m.lastDir != "/tmp/project" {
			t.Errorf("should set working dir, got %q", m.lastDir)
		}
	})
}

func TestStopUsesArgArrays(t *testing.T) {
	m := &mockRunner{}
	withMockRunner(m, func() {
		if err := Stop("/tmp/p", "test"); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(strings.Join(m.lastArgs, " "), "stop") {
			t.Error("should pass 'stop'")
		}
	})
}

func TestRestartUsesArgArrays(t *testing.T) {
	m := &mockRunner{}
	withMockRunner(m, func() {
		if err := Restart("/tmp/p", "test"); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(strings.Join(m.lastArgs, " "), "restart") {
			t.Error("should pass 'restart'")
		}
	})
}

func TestRebuildUsesArgArrays(t *testing.T) {
	m := &mockRunner{}
	withMockRunner(m, func() {
		var stdout, stderr bytes.Buffer
		if err := Rebuild("/tmp/p", "test", &stdout, &stderr); err != nil {
			t.Fatal(err)
		}
		args := strings.Join(m.lastArgs, " ")
		if !strings.Contains(args, "--build") || !strings.Contains(args, "--force-recreate") {
			t.Errorf("should pass --build --force-recreate, got: %s", args)
		}
	})
}

func TestDownUsesArgArrays(t *testing.T) {
	m := &mockRunner{}
	withMockRunner(m, func() {
		if err := Down("/tmp/p", "test"); err != nil {
			t.Fatal(err)
		}
		args := strings.Join(m.lastArgs, " ")
		if !strings.Contains(args, "down") {
			t.Error("should pass 'down'")
		}
		if strings.Contains(args, "-v") {
			t.Error("Down should NOT pass -v")
		}
	})
}

func TestDownWithVolumesUsesArgArrays(t *testing.T) {
	m := &mockRunner{}
	withMockRunner(m, func() {
		if err := DownWithVolumes("/tmp/p", "test"); err != nil {
			t.Fatal(err)
		}
		args := strings.Join(m.lastArgs, " ")
		if !strings.Contains(args, "down -v") {
			t.Errorf("should pass 'down -v', got: %s", args)
		}
	})
}

func TestPullUsesArgArrays(t *testing.T) {
	m := &mockRunner{}
	withMockRunner(m, func() {
		var stdout, stderr bytes.Buffer
		if err := Pull("/tmp/p", "test", &stdout, &stderr); err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(strings.Join(m.lastArgs, " "), "pull") {
			t.Error("should pass 'pull'")
		}
	})
}

func TestLogsUsesArgArrays(t *testing.T) {
	m := &mockRunner{}
	withMockRunner(m, func() {
		var stdout, stderr bytes.Buffer
		if err := Logs("/tmp/p", "test", false, &stdout, &stderr); err != nil {
			t.Fatal(err)
		}
		args := strings.Join(m.lastArgs, " ")
		if !strings.Contains(args, "logs") {
			t.Error("should pass 'logs'")
		}
		if strings.Contains(args, "-f") {
			t.Error("should NOT pass -f when follow=false")
		}
	})
}

func TestLogsFollowUsesArgArrays(t *testing.T) {
	m := &mockRunner{}
	withMockRunner(m, func() {
		var stdout, stderr bytes.Buffer
		if err := Logs("/tmp/p", "test", true, &stdout, &stderr); err != nil {
			t.Fatal(err)
		}
		args := strings.Join(m.lastArgs, " ")
		if !strings.Contains(args, "-f") {
			t.Error("should pass -f when follow=true")
		}
	})
}

func TestUpErrorPropagates(t *testing.T) {
	m := &mockRunner{err: fmt.Errorf("docker failed")}
	withMockRunner(m, func() {
		var stdout, stderr bytes.Buffer
		err := Up("/tmp/p", "test", &stdout, &stderr)
		if err == nil {
			t.Error("expected error")
		}
		if !strings.Contains(err.Error(), "starting project") {
			t.Errorf("error should wrap context, got: %v", err)
		}
	})
}

func TestNoShellInterpolation(t *testing.T) {
	m := &mockRunner{}
	withMockRunner(m, func() {
		var stdout, stderr bytes.Buffer
		_ = Up("/tmp/p", "evil;rm -rf /", &stdout, &stderr)

		if m.lastName != "docker" {
			t.Error("command name should always be 'docker', never a shell")
		}
		// The project name with shell metacharacters is passed as a single
		// exec.Command argument - not parsed by a shell - so injection is impossible.
		foundProjectName := false
		for _, arg := range m.lastArgs {
			if arg == "evil;rm -rf /" {
				foundProjectName = true
			}
		}
		if !foundProjectName {
			t.Error("project name with shell chars should be passed as a single argument")
		}
	})
}

func TestOpenBrowser(t *testing.T) {
	orig := browserStarter
	defer func() { browserStarter = orig }()

	var gotCmd, gotURL string
	browserStarter = func(browserCmd, url string) error {
		gotCmd = browserCmd
		gotURL = url
		return nil
	}

	if err := OpenBrowser("http://localhost:8080", "firefox"); err != nil {
		t.Fatal(err)
	}
	if gotCmd != "firefox" {
		t.Errorf("expected 'firefox', got %q", gotCmd)
	}
	if gotURL != "http://localhost:8080" {
		t.Errorf("expected URL passed through, got %q", gotURL)
	}
}

func TestOpenBrowserDefaultCmd(t *testing.T) {
	orig := browserStarter
	defer func() { browserStarter = orig }()

	var gotCmd string
	browserStarter = func(browserCmd, _ string) error {
		gotCmd = browserCmd
		return nil
	}

	if err := OpenBrowser("http://localhost", ""); err != nil {
		t.Fatal(err)
	}
	if gotCmd != "open" {
		t.Errorf("empty browser cmd should default to 'open', got %q", gotCmd)
	}
}
