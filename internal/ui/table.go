// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package ui

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

// Table renders headers and rows as a rounded bordered table.
// If stateCol >= 0, cells in that column are colored by container state.
func Table(headers []string, rows [][]string, stateCol int) string {
	t := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(gray).
		Headers(headers...).
		Rows(rows...).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == table.HeaderRow {
				return lipgloss.NewStyle().Bold(true).Padding(0, 1)
			}
			if col == stateCol && row >= 0 && row < len(rows) && col < len(rows[row]) {
				return stateStyleFor(rows[row][col]).Padding(0, 1)
			}
			return lipgloss.NewStyle().Padding(0, 1)
		})
	return t.String()
}

func stateStyleFor(s string) lipgloss.Style {
	switch s {
	case "running":
		return green
	case "stopped", "exited", "created":
		return yellow
	default:
		return gray
	}
}
