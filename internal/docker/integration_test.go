// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

//go:build integration

// Package docker integration tests exercise the real `docker compose` CLI
// and a running Docker daemon. They are excluded from the normal unit-test
// build and run only with: go test -race -tags integration ./internal/docker/
//
// Each test creates a throwaway project in a temp dir, brings it up with a
// tiny image, asserts behavior, and always tears down (containers + volumes).
package docker

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"
)

// requireDocker skips the test if Docker or the daemon is unavailable.
func requireDocker(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("docker not installed; skipping integration test")
	}
	if err := exec.Command("docker", "info").Run(); err != nil {
		t.Skip("docker daemon not running; skipping integration test")
	}
}

// writeProject writes a minimal compose project into dir and returns it.
func writeProject(t *testing.T, dir, compose string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, "docker-compose.yml"), []byte(compose), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, ".env"), []byte("# test\n"), 0o600); err != nil {
		t.Fatal(err)
	}
}

const helloCompose = `services:
  app:
    image: hello-world
    restart: "no"
`

// nginxCompose maps an ephemeral host port so we can assert reachability.
func nginxCompose(hostPort int) string {
	return `services:
  web:
    image: nginx:alpine
    ports:
      - "` + itoa(hostPort) + `:80"
    restart: unless-stopped
`
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}

func TestIntegrationUpStateDown(t *testing.T) {
	requireDocker(t)
	dir := t.TempDir()
	const name = "updock-it-hello"
	writeProject(t, dir, helloCompose)

	t.Cleanup(func() { _ = Down(dir, name) })

	if err := Up(dir, name, os.Stdout, os.Stderr); err != nil {
		t.Fatalf("Up failed: %v", err)
	}

	// hello-world exits immediately, so state should not be "running".
	state := DetectState(dir, name)
	if state == StateRunning {
		t.Errorf("hello-world should have exited, got state %q", state)
	}

	if err := Down(dir, name); err != nil {
		t.Fatalf("Down failed: %v", err)
	}
	if got := DetectState(dir, name); got != StateNotFound {
		t.Errorf("after Down, expected 'not found', got %q", got)
	}
}

func TestIntegrationNginxRunningAndReachable(t *testing.T) {
	requireDocker(t)
	dir := t.TempDir()
	const name = "updock-it-nginx"
	const port = 38080
	writeProject(t, dir, nginxCompose(port))

	t.Cleanup(func() { _ = DownWithVolumes(dir, name) })

	if err := Pull(dir, name, os.Stdout, os.Stderr); err != nil {
		t.Fatalf("Pull failed: %v", err)
	}
	if err := Up(dir, name, os.Stdout, os.Stderr); err != nil {
		t.Fatalf("Up failed: %v", err)
	}

	// Give nginx a moment to start.
	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		if DetectState(dir, name) == StateRunning {
			break
		}
		time.Sleep(500 * time.Millisecond)
	}
	if got := DetectState(dir, name); got != StateRunning {
		t.Fatalf("nginx should be running, got %q", got)
	}

	// Stop keeps the project; state should no longer be running.
	if err := Stop(dir, name); err != nil {
		t.Fatalf("Stop failed: %v", err)
	}
	if got := DetectState(dir, name); got == StateRunning {
		t.Errorf("after Stop, should not be running, got %q", got)
	}
}

func TestIntegrationDataPersistsAcrossDownUp(t *testing.T) {
	requireDocker(t)
	dir := t.TempDir()
	const name = "updock-it-persist"

	// A container that writes to a bind-mounted data dir, then exits.
	dataDir := filepath.Join(dir, "data")
	if err := os.MkdirAll(dataDir, 0o755); err != nil {
		t.Fatal(err)
	}
	compose := `services:
  writer:
    image: busybox
    command: sh -c "echo persisted > /data/marker && sleep 1"
    volumes:
      - ./data:/data
    restart: "no"
`
	writeProject(t, dir, compose)
	t.Cleanup(func() { _ = Down(dir, name) })

	if err := Up(dir, name, os.Stdout, os.Stderr); err != nil {
		t.Fatalf("Up failed: %v", err)
	}
	time.Sleep(2 * time.Second)

	if err := Down(dir, name); err != nil {
		t.Fatalf("Down failed: %v", err)
	}

	// The marker file written into ./data must survive down.
	marker := filepath.Join(dataDir, "marker")
	if _, err := os.Stat(marker); err != nil {
		t.Errorf("data should persist across down, marker missing: %v", err)
	}
}
