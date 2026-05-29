// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"context"
	"fmt"
	"net"
	"os/user"
	"testing"
	"time"
)

func passingEnv() Env {
	return Env{
		LookPath:   func(_ string) (string, error) { return "/usr/bin/docker", nil },
		RunCommand: func(_ string, _ ...string) error { return nil },
		CommandOutput: func(_ context.Context, _ string, _ ...string) ([]byte, error) {
			return []byte("27.0.0"), nil
		},
		DialSocket: func(_, _ string, _ time.Duration) (net.Conn, error) {
			return &net.UnixConn{}, nil
		},
		CurrentUser:   func() (*user.User, error) { return &user.User{Username: "test"}, nil },
		UserGroupIDs:  func(_ *user.User) ([]string, error) { return []string{"1000"}, nil },
		LookupGroupID: func(_ string) (*user.Group, error) { return &user.Group{Name: "docker"}, nil },
		RuntimeGOOS:   "linux",
	}
}

func failingEnv() Env {
	return Env{
		LookPath:   func(_ string) (string, error) { return "", fmt.Errorf("not found") },
		RunCommand: func(_ string, _ ...string) error { return fmt.Errorf("not found") },
		CommandOutput: func(_ context.Context, _ string, _ ...string) ([]byte, error) {
			return nil, fmt.Errorf("cannot connect")
		},
		DialSocket: func(_, _ string, _ time.Duration) (net.Conn, error) {
			return nil, fmt.Errorf("permission denied")
		},
		CurrentUser:   func() (*user.User, error) { return &user.User{Username: "test"}, nil },
		UserGroupIDs:  func(_ *user.User) ([]string, error) { return []string{"1000"}, nil },
		LookupGroupID: func(_ string) (*user.Group, error) { return &user.Group{Name: "staff"}, nil },
		RuntimeGOOS:   "linux",
	}
}

func TestRunPreflightAllPass(t *testing.T) {
	results := RunPreflightWith(passingEnv())
	if len(results) != 4 {
		t.Fatalf("expected 4 checks, got %d", len(results))
	}
	for _, r := range results {
		if r.Status != StatusPass {
			t.Errorf("check %q should pass, got status %d: %s", r.Name, r.Status, r.Message)
		}
	}
}

func TestRunPreflightAllFail(t *testing.T) {
	results := RunPreflightWith(failingEnv())
	for _, r := range results {
		if r.Status != StatusFail {
			t.Errorf("check %q should fail, got status %d", r.Name, r.Status)
		}
		if r.Message == "" {
			t.Errorf("check %q should have a message", r.Name)
		}
		if r.Fix == "" {
			t.Errorf("check %q should have a fix", r.Name)
		}
	}
}

func TestCheckDockerInstalled(t *testing.T) {
	tests := []struct {
		name   string
		found  bool
		expect CheckStatus
	}{
		{"installed", true, StatusPass},
		{"not installed", false, StatusFail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := passingEnv()
			if !tt.found {
				env.LookPath = func(_ string) (string, error) { return "", fmt.Errorf("not found") }
			}
			r := checkDockerInstalled(env)
			if r.Status != tt.expect {
				t.Errorf("expected status %d, got %d", tt.expect, r.Status)
			}
		})
	}
}

func TestCheckComposeAvailable(t *testing.T) {
	tests := []struct {
		name   string
		plugin bool
		legacy bool
		expect CheckStatus
	}{
		{"plugin available", true, false, StatusPass},
		{"legacy available", false, true, StatusPass},
		{"neither available", false, false, StatusFail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := passingEnv()
			if !tt.plugin {
				env.RunCommand = func(_ string, _ ...string) error { return fmt.Errorf("not found") }
			}
			if !tt.legacy {
				env.LookPath = func(f string) (string, error) {
					if f == "docker-compose" {
						return "", fmt.Errorf("not found")
					}
					return "/usr/bin/" + f, nil
				}
			}
			r := checkComposeAvailable(env)
			if r.Status != tt.expect {
				t.Errorf("expected status %d, got %d: %s", tt.expect, r.Status, r.Message)
			}
		})
	}
}

func TestCheckDaemonRunning(t *testing.T) {
	tests := []struct {
		name   string
		output []byte
		err    error
		expect CheckStatus
	}{
		{"running", []byte("27.0.0"), nil, StatusPass},
		{"error", nil, fmt.Errorf("fail"), StatusFail},
		{"empty version", []byte(""), nil, StatusFail},
		{"whitespace only", []byte("  \n"), nil, StatusFail},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			env := passingEnv()
			env.CommandOutput = func(_ context.Context, _ string, _ ...string) ([]byte, error) {
				return tt.output, tt.err
			}
			r := checkDaemonRunning(env)
			if r.Status != tt.expect {
				t.Errorf("expected status %d, got %d: %s", tt.expect, r.Status, r.Message)
			}
		})
	}
}

func TestCheckSocketAccessible(t *testing.T) {
	t.Run("windows skips", func(t *testing.T) {
		env := passingEnv()
		env.RuntimeGOOS = "windows"
		r := checkSocketAccessible(env)
		if r.Status != StatusPass {
			t.Errorf("Windows should always pass socket check")
		}
	})

	t.Run("socket reachable", func(t *testing.T) {
		r := checkSocketAccessible(passingEnv())
		if r.Status != StatusPass {
			t.Errorf("expected pass, got %d", r.Status)
		}
	})

	t.Run("socket unreachable not in docker group", func(t *testing.T) {
		env := failingEnv()
		r := checkSocketAccessible(env)
		if r.Status != StatusFail {
			t.Errorf("expected fail, got %d", r.Status)
		}
		if r.Fix == "" {
			t.Error("should suggest adding to docker group")
		}
	})

	t.Run("socket unreachable user error", func(t *testing.T) {
		env := failingEnv()
		env.CurrentUser = func() (*user.User, error) { return nil, fmt.Errorf("no user") }
		r := checkSocketAccessible(env)
		if r.Status != StatusFail {
			t.Errorf("expected fail, got %d", r.Status)
		}
	})

	t.Run("socket unreachable in docker group", func(t *testing.T) {
		env := failingEnv()
		env.LookupGroupID = func(_ string) (*user.Group, error) {
			return &user.Group{Name: "docker"}, nil
		}
		r := checkSocketAccessible(env)
		if r.Status != StatusFail {
			t.Errorf("expected fail, got %d", r.Status)
		}
		if r.Fix != "Ensure the Docker daemon is running and you have permission to access the socket." {
			t.Errorf("unexpected fix: %q", r.Fix)
		}
	})
}

func TestDaemonFix(t *testing.T) {
	tests := []struct {
		goos string
		want string
	}{
		{"linux", "Start the daemon: sudo systemctl start docker"},
		{"darwin", "Start Docker Desktop or the Docker daemon"},
		{"windows", "Start Docker Desktop or the Docker daemon"},
	}
	for _, tt := range tests {
		t.Run(tt.goos, func(t *testing.T) {
			got := daemonFix(tt.goos)
			if got != tt.want {
				t.Errorf("daemonFix(%q) = %q, want %q", tt.goos, got, tt.want)
			}
		})
	}
}

func TestRunPreflightFastWithMock(t *testing.T) {
	t.Run("all pass", func(t *testing.T) {
		// Use real RunPreflightFast which calls DefaultEnv - just verify no panic
		_ = RunPreflightFast()
	})
}

func TestRunPreflightReturnsCorrectNames(t *testing.T) {
	results := RunPreflightWith(passingEnv())
	expectedNames := []string{
		"Docker installed",
		"Docker Compose available",
		"Docker daemon running",
		"Docker socket accessible",
	}
	for i, name := range expectedNames {
		if results[i].Name != name {
			t.Errorf("check %d: expected name %q, got %q", i, name, results[i].Name)
		}
	}
}
