// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package prompt

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/amrelsagaei/updock/internal/project"
)

// ProposeVolumes creates default volume mappings from image config,
// mounting to data/<name>/ inside the project folder.
func ProposeVolumes(containerVolumes []string) []project.VolumeMapping {
	mappings := make([]project.VolumeMapping, 0, len(containerVolumes))
	for _, vol := range containerVolumes {
		name := volumeName(vol)
		mappings = append(mappings, project.VolumeMapping{
			HostPath:      filepath.Join("data", name),
			ContainerPath: vol,
		})
	}
	return mappings
}

func volumeName(containerPath string) string {
	name := filepath.Base(containerPath)
	if name == "/" || name == "." || name == "" {
		return "data"
	}
	return name
}

// ConfirmVolumeMappings prompts the user to accept or edit volume mappings.
func ConfirmVolumeMappings(proposed []project.VolumeMapping) ([]project.VolumeMapping, error) {
	result := make([]project.VolumeMapping, len(proposed))

	for i, m := range proposed {
		title := fmt.Sprintf("Host path for %s", m.ContainerPath)
		val, err := InputFunc(title, m.HostPath, m.HostPath)
		if err != nil {
			return nil, fmt.Errorf("volume configuration: %w", err)
		}

		hostPath := strings.TrimSpace(val)
		if hostPath == "" {
			hostPath = m.HostPath
		}

		result[i] = project.VolumeMapping{
			HostPath:      hostPath,
			ContainerPath: m.ContainerPath,
		}
	}

	return result, nil
}
