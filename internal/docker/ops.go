// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

// CommandRunner executes shell commands. Replaceable in tests.
var CommandRunner Runner = defaultRunner{}

// Runner abstracts command execution for testability.
type Runner interface {
	Run(dir string, stdout, stderr io.Writer, name string, args ...string) error
	Output(dir, name string, args ...string) ([]byte, error)
}

type defaultRunner struct{}

func (defaultRunner) Run(dir string, stdout, stderr io.Writer, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	cmd.Stdout = stdout
	cmd.Stderr = stderr
	return cmd.Run()
}

func (defaultRunner) Output(dir, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Output()
}

func compose(projectPath, projectName string, stdout, stderr io.Writer, args ...string) error {
	fullArgs := append([]string{"compose", "--project-name", projectName}, args...)
	return CommandRunner.Run(projectPath, stdout, stderr, "docker", fullArgs...)
}

func composeQuiet(projectPath, projectName string, args ...string) error {
	return compose(projectPath, projectName, io.Discard, io.Discard, args...)
}

func composeOutput(projectPath, projectName string, args ...string) ([]byte, error) {
	fullArgs := append([]string{"compose", "--project-name", projectName}, args...)
	return CommandRunner.Output(projectPath, "docker", fullArgs...)
}

// Up starts the project in the background.
func Up(projectPath, projectName string, stdout, stderr io.Writer) error {
	if err := compose(projectPath, projectName, stdout, stderr, "up", "-d"); err != nil {
		return fmt.Errorf("starting project: %w", err)
	}
	return nil
}

// Pull pulls images for the project.
func Pull(projectPath, projectName string, stdout, stderr io.Writer) error {
	if err := compose(projectPath, projectName, stdout, stderr, "pull"); err != nil {
		return fmt.Errorf("pulling images: %w", err)
	}
	return nil
}

// Stop stops containers without removing them.
func Stop(projectPath, projectName string) error {
	if err := composeQuiet(projectPath, projectName, "stop"); err != nil {
		return fmt.Errorf("stopping project: %w", err)
	}
	return nil
}

// Restart restarts all containers.
func Restart(projectPath, projectName string) error {
	if err := composeQuiet(projectPath, projectName, "restart"); err != nil {
		return fmt.Errorf("restarting project: %w", err)
	}
	return nil
}

// Rebuild recreates containers with a fresh build.
func Rebuild(projectPath, projectName string, stdout, stderr io.Writer) error {
	if err := compose(projectPath, projectName, stdout, stderr, "up", "-d", "--build", "--force-recreate"); err != nil {
		return fmt.Errorf("rebuilding project: %w", err)
	}
	return nil
}

// Down removes containers but preserves volumes and project folder.
func Down(projectPath, projectName string) error {
	if err := composeQuiet(projectPath, projectName, "down"); err != nil {
		return fmt.Errorf("removing containers: %w", err)
	}
	return nil
}

// DownWithVolumes removes containers and named volumes.
func DownWithVolumes(projectPath, projectName string) error {
	if err := composeQuiet(projectPath, projectName, "down", "-v"); err != nil {
		return fmt.Errorf("removing containers and volumes: %w", err)
	}
	return nil
}

// Logs streams container logs. If follow is true, it tails.
func Logs(projectPath, projectName string, follow bool, stdout, stderr io.Writer) error {
	args := []string{"logs"}
	if follow {
		args = append(args, "-f")
	}
	if err := compose(projectPath, projectName, stdout, stderr, args...); err != nil {
		return fmt.Errorf("fetching logs: %w", err)
	}
	return nil
}

// browserStarter launches the browser command. Overridable in tests.
var browserStarter = func(browserCmd, url string) error {
	cmd := exec.Command(browserCmd, url)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Start()
}

// OpenBrowser opens a URL in the user's default browser.
func OpenBrowser(url, browserCmd string) error {
	if browserCmd == "" {
		browserCmd = "open"
	}
	return browserStarter(browserCmd, url)
}
