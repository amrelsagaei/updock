// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

// Package ui provides styled terminal output built on Lip Gloss.
// Colors degrade automatically on non-TTY output and when NO_COLOR is set,
// so piped output stays plain and tests see unstyled text.
package ui

import (
	"fmt"
	"io"

	"github.com/charmbracelet/lipgloss"
)

var (
	green  = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	red    = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	yellow = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	blue   = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	purple = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	gray   = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))

	boldStyle  = lipgloss.NewStyle().Bold(true)
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("14"))
	linkStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("14")).Underline(true)
)

// Success prints a green check line.
func Success(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, "  %s %s\n", green.Render("✓"), fmt.Sprintf(format, args...))
}

// Fail prints a red cross line.
func Fail(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, "  %s %s\n", red.Render("✗"), fmt.Sprintf(format, args...))
}

// Warn prints a yellow warning line.
func Warn(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, "  %s %s\n", yellow.Render("⚠"), fmt.Sprintf(format, args...))
}

// Info prints a plain indented line.
func Info(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, "  %s\n", fmt.Sprintf(format, args...))
}

// Step prints a dimmed in-progress line (e.g. "Pulling images...").
func Step(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, "  %s\n", gray.Render(fmt.Sprintf(format, args...)))
}

// Title renders a bold accent title string.
func Title(s string) string { return titleStyle.Render(s) }

// Bold renders bold text.
func Bold(s string) string { return boldStyle.Render(s) }

// Dim renders muted text.
func Dim(s string) string { return gray.Render(s) }

// Link renders an underlined accent URL.
func Link(s string) string { return linkStyle.Render(s) }

// Label renders a dimmed field label for detail views.
func Label(s string) string { return gray.Render(s) }

// Badge renders a colored badge for an image's trust level.
func Badge(kind string) string {
	switch kind {
	case "official":
		return blue.Render("official")
	case "popular":
		return purple.Render("popular")
	case "best match":
		return green.Render("[best match]")
	default:
		return gray.Render("community")
	}
}

// State renders a container state with a status color.
func State(s string) string {
	switch s {
	case "running":
		return green.Render(s)
	case "stopped", "exited", "created":
		return yellow.Render(s)
	default:
		return gray.Render(s)
	}
}

// StateWidth returns the visible (unstyled) width of a state string, so
// callers can pad columns correctly regardless of color codes.
func StateWidth(s string) int { return lipgloss.Width(State(s)) }
