// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"strings"
)

// ContainerState represents the state of a project's containers.
type ContainerState string

// Container state values.
const (
	StateRunning  ContainerState = "running"
	StateStopped  ContainerState = "stopped"
	StateCreated  ContainerState = "created"
	StateNotFound ContainerState = "not found"
)

// DetectState checks the state of containers for a project.
func DetectState(projectPath, projectName string) ContainerState {
	out, err := composeOutput(projectPath, projectName, "ps", "--format", "{{.State}}")
	if err != nil {
		return StateNotFound
	}

	lines := strings.Fields(strings.TrimSpace(string(out)))
	if len(lines) == 0 {
		return StateNotFound
	}

	allRunning := true
	for _, line := range lines {
		state := strings.ToLower(line)
		if state != "running" {
			allRunning = false
			break
		}
	}

	if allRunning {
		return StateRunning
	}

	for _, line := range lines {
		state := strings.ToLower(line)
		if state == "created" {
			return StateCreated
		}
	}

	return StateStopped
}
