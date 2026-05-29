// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/ui"
	"github.com/spf13/cobra"
)

func newRebuildCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rebuild [number]",
		Short: "Rebuild and recreate containers",
		Long: `Rebuild and recreate the containers for the project with the given number
(from 'updock ls'). Use this after changing the image or configuration.`,
		Args: cobra.ExactArgs(1),
		RunE: runRebuild,
	}
}

func runRebuild(cmd *cobra.Command, args []string) error {
	entry, err := resolveProject(args[0])
	if err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	ui.Step(w, "Rebuilding %s...", entry.Name)

	if err := docker.Rebuild(entry.Path, entry.Name, w, cmd.ErrOrStderr()); err != nil {
		return err
	}

	ui.Success(w, "%s rebuilt and running.", entry.Name)
	return nil
}
