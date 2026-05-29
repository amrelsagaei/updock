// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"fmt"
	"io"

	"github.com/amrelsagaei/updock/internal/ui"
)

// FormatResults writes a human-readable ranked results list to w.
func FormatResults(w io.Writer, results []RankedResult, query string) {
	if len(results) == 0 {
		_, _ = fmt.Fprintf(w, "  No results found for %q. Check the spelling or try a shorter query.\n", query)
		return
	}

	_, _ = fmt.Fprintf(w, "\n  %s\n\n", ui.Title(fmt.Sprintf("Results for %q", query)))

	for i, r := range results {
		num := ui.Dim(fmt.Sprintf("%d)", i+1))
		name := ui.Bold(r.RepoName)
		pulls := ui.Dim(fmt.Sprintf("%s pulls", formatPulls(r.PullCount)))
		badge := ui.Badge(r.Badge)

		line := fmt.Sprintf("  %s %s  %s  %s", num, name, badge, pulls)
		if i == 0 {
			line += "  " + ui.Badge("best match")
		}
		_, _ = fmt.Fprintln(w, line)
	}
	_, _ = fmt.Fprintln(w)
}

func formatPulls(n int) string {
	switch {
	case n >= 1_000_000_000:
		return fmt.Sprintf("%.1fB+", float64(n)/1_000_000_000)
	case n >= 1_000_000:
		return fmt.Sprintf("%.1fM+", float64(n)/1_000_000)
	case n >= 1_000:
		return fmt.Sprintf("%.1fK+", float64(n)/1_000)
	case n > 0:
		return fmt.Sprintf("%d", n)
	default:
		return "-"
	}
}
