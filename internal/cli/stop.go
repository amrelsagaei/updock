// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/ui"
	"github.com/spf13/cobra"
)

func newStopCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "stop [number]",
		Short: "Stop containers, keep data",
		Long: `Stop the containers for the project with the given number (from
'updock ls'). The project and its data are kept; start it again with
'updock up <n>'.`,
		Args: cobra.ExactArgs(1),
		RunE: runStop,
	}
}

func runStop(cmd *cobra.Command, args []string) error {
	entry, err := resolveProject(args[0])
	if err != nil {
		return err
	}

	if err := docker.Stop(entry.Path, entry.Name); err != nil {
		return err
	}

	ui.Success(cmd.OutOrStdout(), "%s stopped.", entry.Name)
	return nil
}
