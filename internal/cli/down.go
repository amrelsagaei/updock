// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/ui"
	"github.com/spf13/cobra"
)

func newDownCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "down [number]",
		Short: "Remove containers, keep folder and data",
		Long: `Remove the containers for the project with the given number (from
'updock ls'). The project folder and the data/ directory are preserved, so
'updock up <n>' brings everything back.`,
		Args: cobra.ExactArgs(1),
		RunE: runDown,
	}
}

func runDown(cmd *cobra.Command, args []string) error {
	entry, err := resolveProject(args[0])
	if err != nil {
		return err
	}

	if err := docker.Down(entry.Path, entry.Name); err != nil {
		return err
	}

	ui.Success(cmd.OutOrStdout(), "%s containers removed. Project folder and data preserved.", entry.Name)
	return nil
}
