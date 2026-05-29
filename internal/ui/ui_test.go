// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package ui

import (
	"bytes"
	"strings"
	"testing"
)

// Tests run with non-TTY stdout, so Lip Gloss emits no color codes and we can
// assert on plain text content.

func TestSuccess(t *testing.T) {
	var buf bytes.Buffer
	Success(&buf, "started %s", "nginx")
	out := buf.String()
	if !strings.Contains(out, "started nginx") {
		t.Errorf("got %q", out)
	}
	if !strings.Contains(out, "✓") {
		t.Errorf("expected check mark, got %q", out)
	}
}

func TestFailWarnInfoStep(t *testing.T) {
	tests := []struct {
		name string
		fn   func(*bytes.Buffer)
		want string
	}{
		{"fail", func(b *bytes.Buffer) { Fail(b, "boom") }, "✗"},
		{"warn", func(b *bytes.Buffer) { Warn(b, "careful") }, "⚠"},
		{"info", func(b *bytes.Buffer) { Info(b, "fyi") }, "fyi"},
		{"step", func(b *bytes.Buffer) { Step(b, "working") }, "working"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			tt.fn(&buf)
			if !strings.Contains(buf.String(), tt.want) {
				t.Errorf("got %q, want substring %q", buf.String(), tt.want)
			}
		})
	}
}

func TestBadge(t *testing.T) {
	tests := []struct{ kind, want string }{
		{"official", "official"},
		{"popular", "popular"},
		{"community", "community"},
		{"best match", "best match"},
	}
	for _, tt := range tests {
		t.Run(tt.kind, func(t *testing.T) {
			if !strings.Contains(Badge(tt.kind), tt.want) {
				t.Errorf("Badge(%q) = %q", tt.kind, Badge(tt.kind))
			}
		})
	}
}

func TestState(t *testing.T) {
	for _, s := range []string{"running", "stopped", "created", "not found"} {
		if !strings.Contains(State(s), s) {
			t.Errorf("State(%q) lost its text: %q", s, State(s))
		}
	}
}

func TestStyleHelpersPreserveText(t *testing.T) {
	for _, s := range []string{"hello", "updock", "http://localhost:8080"} {
		if !strings.Contains(Title(s), s) {
			t.Errorf("Title dropped text: %q", Title(s))
		}
		if !strings.Contains(Bold(s), s) {
			t.Errorf("Bold dropped text: %q", Bold(s))
		}
		if !strings.Contains(Dim(s), s) {
			t.Errorf("Dim dropped text: %q", Dim(s))
		}
		if !strings.Contains(Link(s), s) {
			t.Errorf("Link dropped text: %q", Link(s))
		}
		if !strings.Contains(Label(s), s) {
			t.Errorf("Label dropped text: %q", Label(s))
		}
	}
}

func TestTable(t *testing.T) {
	headers := []string{"#", "PROJECT", "STATE"}
	rows := [][]string{
		{"1", "nginx", "running"},
		{"2", "postgres", "stopped"},
	}
	out := Table(headers, rows, 2)

	for _, want := range []string{"PROJECT", "nginx", "running", "postgres", "stopped"} {
		if !strings.Contains(out, want) {
			t.Errorf("table missing %q in:\n%s", want, out)
		}
	}
}

func TestTableNoStateColumn(t *testing.T) {
	out := Table([]string{"A", "B"}, [][]string{{"x", "y"}}, -1)
	if !strings.Contains(out, "x") || !strings.Contains(out, "y") {
		t.Errorf("table dropped content:\n%s", out)
	}
}

func TestStateWidth(t *testing.T) {
	// Width must be the visible length, ignoring any color codes.
	if got := StateWidth("running"); got != len("running") {
		t.Errorf("StateWidth(running) = %d, want %d", got, len("running"))
	}
}
