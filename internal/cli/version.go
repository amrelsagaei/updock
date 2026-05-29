// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/amrelsagaei/updock/internal/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version and build info",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, _ []string) {
			w := cmd.OutOrStdout()
			_, _ = fmt.Fprintf(w, "updock %s\n", version.Version)
			_, _ = fmt.Fprintf(w, "  commit:  %s\n", version.Commit)
			_, _ = fmt.Fprintf(w, "  built:   %s\n", version.Date)
			_, _ = fmt.Fprintf(w, "  go:      %s\n", runtime.Version())
			_, _ = fmt.Fprintf(w, "  os/arch: %s/%s\n", runtime.GOOS, runtime.GOARCH)
		},
	}
}
