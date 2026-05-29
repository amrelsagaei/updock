// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/project"
	"github.com/amrelsagaei/updock/internal/security"
	"github.com/amrelsagaei/updock/internal/ui"
)

func newStatusCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "status [number]",
		Short: "Show details for one project",
		Long: `Show full details for the project with the given number (from 'updock ls'):
image, live state, port mappings, environment variables (secrets masked),
and the project path.`,
		Args: cobra.ExactArgs(1),
		RunE: runStatus,
	}
}

func runStatus(cmd *cobra.Command, args []string) error {
	entry, err := resolveProject(args[0])
	if err != nil {
		return err
	}

	state := docker.DetectState(entry.Path, entry.Name)

	cfg, err := loadProjectConfig(entry.Path)
	if err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	_, _ = fmt.Fprintf(w, "\n  %s %s %s\n", ui.Label("Project: "), ui.Bold(entry.Name), ui.Dim(fmt.Sprintf("(#%d)", entry.Number)))
	_, _ = fmt.Fprintf(w, "  %s %s:%s\n", ui.Label("Image:   "), entry.Metadata.Image, entry.Metadata.Tag)
	_, _ = fmt.Fprintf(w, "  %s %s\n", ui.Label("State:   "), ui.State(string(state)))

	if len(entry.Metadata.Ports) > 0 {
		_, _ = fmt.Fprintf(w, "  %s ", ui.Label("Ports:   "))
		for i, p := range entry.Metadata.Ports {
			if i > 0 {
				_, _ = fmt.Fprint(w, "           ")
			}
			_, _ = fmt.Fprintf(w, "%d -> %s\n", p.Container, ui.Link(fmt.Sprintf("localhost:%d", p.Host)))
		}
	}

	if len(cfg.Env) > 0 {
		_, _ = fmt.Fprintf(w, "  %s ", ui.Label("Env:     "))
		for i, e := range cfg.Env {
			if i > 0 {
				_, _ = fmt.Fprint(w, "           ")
			}
			_, _ = fmt.Fprintf(w, "%s=%s\n", e.Key, security.MaskSecret(e.Value, e.Secret))
		}
	}

	_, _ = fmt.Fprintf(w, "  %s %s\n", ui.Label("Created: "), entry.Metadata.CreatedAt.Format("2006-01-02 15:04"))
	_, _ = fmt.Fprintf(w, "  %s %s\n\n", ui.Label("Path:    "), ui.Dim(entry.Path))
	return nil
}

func loadProjectConfig(projectPath string) (*project.Config, error) {
	meta, err := project.ReadMetadata(projectPath)
	if err != nil {
		return nil, err
	}

	envVars := readEnvFile(projectPath)

	ports := make([]project.PortMapping, len(meta.Ports))
	for i, p := range meta.Ports {
		ports[i] = project.PortMapping{Host: p.Host, Container: p.Container}
	}

	return &project.Config{
		Image:       meta.Image,
		Tag:         meta.Tag,
		ProjectName: meta.ProjectName,
		Ports:       ports,
		Env:         envVars,
	}, nil
}

func readEnvFile(projectPath string) []project.EnvVar {
	data, err := os.ReadFile(filepath.Join(projectPath, ".env"))
	if err != nil {
		return nil
	}

	var vars []project.EnvVar
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if k, v, ok := strings.Cut(line, "="); ok {
			// HOST_PORT_* entries are managed through Ports, not as env vars.
			if strings.HasPrefix(k, "HOST_PORT_") {
				continue
			}
			vars = append(vars, project.EnvVar{
				Key:    k,
				Value:  v,
				Secret: docker.IsSecretKey(k),
			})
		}
	}
	return vars
}
