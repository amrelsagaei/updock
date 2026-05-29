// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"context"
	"fmt"
	"net"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
	"time"
)

// CheckStatus represents the outcome of a preflight check.
type CheckStatus int

// Preflight check statuses.
const (
	StatusPass CheckStatus = iota
	StatusWarn
	StatusFail
)

// CheckResult holds the outcome of a single preflight check.
type CheckResult struct {
	Name    string
	Status  CheckStatus
	Message string
	Fix     string
}

// Env abstracts system queries so preflight checks are testable.
type Env struct {
	LookPath      func(file string) (string, error)
	RunCommand    func(name string, args ...string) error
	CommandOutput func(ctx context.Context, name string, args ...string) ([]byte, error)
	DialSocket    func(network, addr string, timeout time.Duration) (net.Conn, error)
	CurrentUser   func() (*user.User, error)
	UserGroupIDs  func(u *user.User) ([]string, error)
	LookupGroupID func(gid string) (*user.Group, error)
	RuntimeGOOS   string
}

// DefaultEnv returns an Env backed by the real operating system.
func DefaultEnv() Env {
	return Env{
		LookPath: exec.LookPath,
		RunCommand: func(name string, args ...string) error {
			return exec.Command(name, args...).Run()
		},
		CommandOutput: func(ctx context.Context, name string, args ...string) ([]byte, error) {
			return exec.CommandContext(ctx, name, args...).Output()
		},
		DialSocket:    net.DialTimeout,
		CurrentUser:   user.Current,
		UserGroupIDs:  func(u *user.User) ([]string, error) { return u.GroupIds() },
		LookupGroupID: user.LookupGroupId,
		RuntimeGOOS:   runtime.GOOS,
	}
}

// RunPreflight executes all environment checks and returns their results.
func RunPreflight() []CheckResult {
	return RunPreflightWith(DefaultEnv())
}

// RunPreflightWith executes all checks using the given environment.
func RunPreflightWith(env Env) []CheckResult {
	return []CheckResult{
		checkDockerInstalled(env),
		checkComposeAvailable(env),
		checkDaemonRunning(env),
		checkSocketAccessible(env),
	}
}

// RunPreflightFast runs all checks and returns the first failure as an error.
func RunPreflightFast() error {
	for _, r := range RunPreflight() {
		if r.Status == StatusFail {
			return fmt.Errorf("%s: %s", r.Name, r.Message)
		}
	}
	return nil
}

func checkDockerInstalled(env Env) CheckResult {
	_, err := env.LookPath("docker")
	if err != nil {
		return CheckResult{
			Name:    "Docker installed",
			Status:  StatusFail,
			Message: "docker is not on your PATH",
			Fix:     "Install Docker: https://docs.docker.com/get-docker/",
		}
	}
	return CheckResult{
		Name:   "Docker installed",
		Status: StatusPass,
	}
}

func checkComposeAvailable(env Env) CheckResult {
	if err := env.RunCommand("docker", "compose", "version"); err == nil {
		return CheckResult{
			Name:   "Docker Compose available",
			Status: StatusPass,
		}
	}

	if _, err := env.LookPath("docker-compose"); err == nil {
		return CheckResult{
			Name:   "Docker Compose available",
			Status: StatusPass,
		}
	}

	return CheckResult{
		Name:    "Docker Compose available",
		Status:  StatusFail,
		Message: "neither 'docker compose' plugin nor 'docker-compose' found",
		Fix:     "Install the Docker Compose plugin: https://docs.docker.com/compose/install/",
	}
}

func checkDaemonRunning(env Env) CheckResult {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	out, err := env.CommandOutput(ctx, "docker", "info", "--format", "{{.ServerVersion}}")
	if err != nil {
		return CheckResult{
			Name:    "Docker daemon running",
			Status:  StatusFail,
			Message: "cannot connect to the Docker daemon",
			Fix:     daemonFix(env.RuntimeGOOS),
		}
	}

	ver := strings.TrimSpace(string(out))
	if ver == "" {
		return CheckResult{
			Name:    "Docker daemon running",
			Status:  StatusFail,
			Message: "Docker daemon did not return a version",
			Fix:     daemonFix(env.RuntimeGOOS),
		}
	}

	return CheckResult{
		Name:   "Docker daemon running",
		Status: StatusPass,
	}
}

func daemonFix(goos string) string {
	if goos == "linux" {
		return "Start the daemon: sudo systemctl start docker"
	}
	return "Start Docker Desktop or the Docker daemon"
}

func checkSocketAccessible(env Env) CheckResult {
	if env.RuntimeGOOS == "windows" {
		return CheckResult{
			Name:   "Docker socket accessible",
			Status: StatusPass,
		}
	}

	socketPath := "/var/run/docker.sock"
	conn, err := env.DialSocket("unix", socketPath, 2*time.Second)
	if err != nil {
		msg := fmt.Sprintf("cannot reach %s", socketPath)

		u, uErr := env.CurrentUser()
		if uErr == nil {
			groups, gErr := env.UserGroupIDs(u)
			inDockerGroup := false
			if gErr == nil {
				for _, gid := range groups {
					g, lErr := env.LookupGroupID(gid)
					if lErr == nil && g.Name == "docker" {
						inDockerGroup = true
						break
					}
				}
			}
			if !inDockerGroup {
				return CheckResult{
					Name:    "Docker socket accessible",
					Status:  StatusFail,
					Message: msg,
					Fix: fmt.Sprintf(
						"Add yourself to the docker group: sudo usermod -aG docker %s && newgrp docker\n"+
							"    Warning: the docker group grants root-equivalent access to your machine.",
						u.Username,
					),
				}
			}
		}

		return CheckResult{
			Name:    "Docker socket accessible",
			Status:  StatusFail,
			Message: msg,
			Fix:     "Ensure the Docker daemon is running and you have permission to access the socket.",
		}
	}
	_ = conn.Close()

	return CheckResult{
		Name:   "Docker socket accessible",
		Status: StatusPass,
	}
}
