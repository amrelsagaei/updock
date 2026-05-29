// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"github.com/amrelsagaei/updock/internal/hub"
	"github.com/amrelsagaei/updock/internal/security"
	"github.com/spf13/cobra"
)

func newSearchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search [name]",
		Short: "Search Docker Hub and show ranked results",
		Long: `Search Docker Hub for a name and print the ranked matches without
running anything. Official and popular images surface first. Use this to
browse before running 'updock <name>'.`,
		Args: cobra.ExactArgs(1),
		RunE: runSearch,
	}
}

func runSearch(cmd *cobra.Command, args []string) error {
	query := args[0]

	if err := security.ValidateImageName(query); err != nil {
		return err
	}

	client := hub.NewClient()
	results, err := client.Search(query, 25)
	if err != nil {
		return err
	}

	ranked := hub.Rank(results, query)
	hub.FormatResults(cmd.OutOrStdout(), ranked, query)
	return nil
}
