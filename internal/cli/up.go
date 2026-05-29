// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"fmt"

	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/ui"
	"github.com/spf13/cobra"
)

func newUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up [number]",
		Short: "Start a project",
		Long: `Start the project with the given number (from 'updock ls'). Pulls images
if needed and brings the stack up in the background, then prints the URLs
you can open.`,
		Args: cobra.ExactArgs(1),
		RunE: runUp,
	}
}

func runUp(cmd *cobra.Command, args []string) error {
	entry, err := resolveProject(args[0])
	if err != nil {
		return err
	}

	if err := docker.RunPreflightFast(); err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	errW := cmd.ErrOrStderr()

	ui.Step(w, "Pulling images for %s...", entry.Name)
	if err := docker.Pull(entry.Path, entry.Name, w, errW); err != nil {
		return err
	}

	ui.Step(w, "Starting %s...", entry.Name)
	if err := docker.Up(entry.Path, entry.Name, w, errW); err != nil {
		return err
	}

	ui.Success(w, "%s is running.", entry.Name)
	for _, p := range entry.Metadata.Ports {
		_, _ = fmt.Fprintf(w, "    %s\n", ui.Link(fmt.Sprintf("http://localhost:%d", p.Host)))
	}
	return nil
}
