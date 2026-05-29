// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"fmt"

	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/ui"
	"github.com/spf13/cobra"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Run all environment checks",
		Long:  "Verifies that Docker, Docker Compose, and the Docker daemon are properly set up.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			w := cmd.OutOrStdout()
			results := docker.RunPreflight()
			allPassed := true

			for _, r := range results {
				switch r.Status {
				case docker.StatusPass:
					ui.Success(w, "%s", r.Name)
				case docker.StatusWarn:
					ui.Warn(w, "%s - %s", r.Name, r.Message)
					if r.Fix != "" {
						_, _ = fmt.Fprintf(w, "    %s %s\n", ui.Label("Fix:"), r.Fix)
					}
				case docker.StatusFail:
					ui.Fail(w, "%s - %s", r.Name, r.Message)
					if r.Fix != "" {
						_, _ = fmt.Fprintf(w, "    %s %s\n", ui.Label("Fix:"), r.Fix)
					}
					allPassed = false
				}
			}

			if allPassed {
				ui.Success(w, "All checks passed. updock is ready.")
			} else {
				ui.Warn(w, "Some checks failed. Fix the issues above before using updock.")
			}
			return nil
		},
	}
}
