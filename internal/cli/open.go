// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"fmt"

	"github.com/amrelsagaei/updock/internal/config"
	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/ui"
	"github.com/spf13/cobra"
)

func newOpenCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "open [number]",
		Short: "Open the mapped port in a browser",
		Long: `Open the first mapped port of the project with the given number (from
'updock ls') in your default browser. Set 'browser_command' in the config
to use a specific browser.`,
		Args: cobra.ExactArgs(1),
		RunE: runOpen,
	}
}

func runOpen(cmd *cobra.Command, args []string) error {
	entry, err := resolveProject(args[0])
	if err != nil {
		return err
	}

	if len(entry.Metadata.Ports) == 0 {
		return fmt.Errorf("project %q has no mapped ports", entry.Name)
	}

	port := entry.Metadata.Ports[0].Host
	url := fmt.Sprintf("http://localhost:%d", port)

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ui.Info(cmd.OutOrStdout(), "Opening %s", ui.Link(url))
	return docker.OpenBrowser(url, cfg.BrowserCommand)
}
