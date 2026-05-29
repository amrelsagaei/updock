// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/spf13/cobra"
)

func newLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs [number]",
		Short: "Show or follow logs",
		Long: `Show the logs for the project with the given number (from 'updock ls').
Pass -f to follow the output live.`,
		Args: cobra.ExactArgs(1),
		RunE: runLogs,
	}
	cmd.Flags().BoolP("follow", "f", false, "Follow log output")
	return cmd
}

func runLogs(cmd *cobra.Command, args []string) error {
	entry, err := resolveProject(args[0])
	if err != nil {
		return err
	}

	follow, _ := cmd.Flags().GetBool("follow")
	return docker.Logs(entry.Path, entry.Name, follow, cmd.OutOrStdout(), cmd.ErrOrStderr())
}
