// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/prompt"
	"github.com/amrelsagaei/updock/internal/ui"
)

func newRmCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rm [number]",
		Short: "Delete a project (confirms first)",
		Long: `Delete the project with the given number (from 'updock ls'): removes the
containers and volumes and deletes the project folder. Asks for confirmation
first unless you pass --yes. This cannot be undone.`,
		Args: cobra.ExactArgs(1),
		RunE: runRm,
	}
	cmd.Flags().BoolP("yes", "y", false, "Skip confirmation")
	return cmd
}

func runRm(cmd *cobra.Command, args []string) error {
	entry, err := resolveProject(args[0])
	if err != nil {
		return err
	}

	yes, _ := cmd.Flags().GetBool("yes")
	if !yes {
		confirmed, err := prompt.Confirm(
			fmt.Sprintf("Delete project '%s' and all its data? This cannot be undone.", entry.Name),
		)
		if err != nil {
			return err
		}
		if !confirmed {
			ui.Info(cmd.OutOrStdout(), "Cancelled.")
			return nil
		}
	}

	_ = docker.DownWithVolumes(entry.Path, entry.Name)

	if err := os.RemoveAll(entry.Path); err != nil {
		return fmt.Errorf("deleting project directory: %w", err)
	}

	ui.Success(cmd.OutOrStdout(), "Project '%s' deleted.", entry.Name)
	return nil
}
