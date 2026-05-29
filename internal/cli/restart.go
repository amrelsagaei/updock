// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"github.com/spf13/cobra"

	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/ui"
)

func newRestartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "restart [number]",
		Short: "Restart a project",
		Long: `Restart the containers for the project with the given number (from
'updock ls'). Use this to apply changes that only need a fresh start.`,
		Args: cobra.ExactArgs(1),
		RunE: runRestart,
	}
}

func runRestart(cmd *cobra.Command, args []string) error {
	entry, err := resolveProject(args[0])
	if err != nil {
		return err
	}

	if err := docker.Restart(entry.Path, entry.Name); err != nil {
		return err
	}

	ui.Success(cmd.OutOrStdout(), "%s restarted.", entry.Name)
	return nil
}
