// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package prompt

import (
	"fmt"
	"net"
	"strconv"

	"github.com/charmbracelet/huh"

	"github.com/amrelsagaei/updock/internal/project"
	"github.com/amrelsagaei/updock/internal/security"
)

// PortChecker tests whether a port is available. Replaceable in tests.
var PortChecker = isPortFree

func isPortFree(port int) bool {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return false
	}
	_ = ln.Close()
	return true
}

// FindFreePort returns the given port if free, or the next free port after it.
func FindFreePort(start int) int {
	for p := start; p <= 65535; p++ {
		if PortChecker(p) {
			return p
		}
	}
	return start
}

// ProposePorts builds port mappings from exposed container ports,
// checking for conflicts and suggesting free alternatives.
func ProposePorts(exposedPorts []int) []project.PortMapping {
	mappings := make([]project.PortMapping, 0, len(exposedPorts))
	usedHost := make(map[int]bool)

	for _, containerPort := range exposedPorts {
		hostPort := containerPort
		if usedHost[hostPort] || !PortChecker(hostPort) {
			hostPort = FindFreePort(containerPort + 1)
		}
		usedHost[hostPort] = true
		mappings = append(mappings, project.PortMapping{
			Host:      hostPort,
			Container: containerPort,
			Protocol:  "tcp",
		})
	}
	return mappings
}

// InputFunc runs an interactive text input. Replaceable in tests.
var InputFunc = runHuhInput

func runHuhInput(title, placeholder, value string) (string, error) {
	result := value
	err := huh.NewInput().
		Title(title).
		Placeholder(placeholder).
		Value(&result).
		Run()
	return result, err
}

// ConfirmPortMappings prompts the user to accept or override port mappings.
func ConfirmPortMappings(proposed []project.PortMapping) ([]project.PortMapping, error) {
	result := make([]project.PortMapping, len(proposed))

	for i, m := range proposed {
		title := fmt.Sprintf("Host port for container port %d", m.Container)
		defaultVal := strconv.Itoa(m.Host)

		if m.Host != m.Container {
			title = fmt.Sprintf("Host port for container port %d (port %d is in use)", m.Container, m.Container)
		}

		val, err := InputFunc(title, defaultVal, defaultVal)
		if err != nil {
			return nil, fmt.Errorf("port configuration: %w", err)
		}

		hostPort, err := strconv.Atoi(val)
		if err != nil {
			hostPort = m.Host
		}
		if portErr := security.ValidatePort(hostPort); portErr != nil {
			hostPort = m.Host
		}

		result[i] = project.PortMapping{
			Host:      hostPort,
			Container: m.Container,
			Protocol:  m.Protocol,
		}
	}

	return result, nil
}
