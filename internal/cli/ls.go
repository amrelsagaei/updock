// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/amrelsagaei/updock/internal/config"
	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/project"
	"github.com/amrelsagaei/updock/internal/ui"
)

// stateDetector resolves a project's live container state. Overridable in tests.
var stateDetector = func(projectPath, projectName string) string {
	return string(docker.DetectState(projectPath, projectName))
}

func newLsCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "ls",
		Aliases: []string{"list"},
		Short:   "List all projects with state",
		Long: `List every updock project with its number, image, live state, and ports.
Use the number with any other command, e.g. 'updock up 2'.`,
		Args: cobra.NoArgs,
		RunE: runLs,
	}
}

func runLs(cmd *cobra.Command, _ []string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	reg, err := project.NewRegistry(cfg.ProjectsRoot)
	if err != nil {
		return err
	}

	entries := reg.List()
	w := cmd.OutOrStdout()

	if len(entries) == 0 {
		_, _ = fmt.Fprintln(w, "  No projects yet. Run 'updock <name>' to create one.")
		return nil
	}

	rows := make([][]string, 0, len(entries))
	for i := range entries {
		e := &entries[i]
		rows = append(rows, []string{
			fmt.Sprintf("%d", e.Number),
			e.Name,
			fmt.Sprintf("%s:%s", e.Metadata.Image, e.Metadata.Tag),
			stateDetector(e.Path, e.Name),
			formatPorts(e.Metadata.Ports),
		})
	}

	headers := []string{"#", "PROJECT", "IMAGE", "STATE", "PORTS"}
	const stateColumn = 3
	_, _ = fmt.Fprintln(w)
	_, _ = fmt.Fprintln(w, ui.Table(headers, rows, stateColumn))
	return nil
}

func formatPorts(ports []project.PortMetadata) string {
	if len(ports) == 0 {
		return "-"
	}
	result := ""
	for i, p := range ports {
		if i > 0 {
			result += ", "
		}
		result += fmt.Sprintf("%d->%d", p.Container, p.Host)
	}
	return result
}
