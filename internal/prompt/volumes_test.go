// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package prompt

import (
	"testing"
)

func TestProposeVolumes(t *testing.T) {
	vols := ProposeVolumes([]string{"/var/lib/postgresql/data", "/tmp/uploads"})

	if len(vols) != 2 {
		t.Fatalf("expected 2 mappings, got %d", len(vols))
	}
	if vols[0].HostPath != "data/data" {
		t.Errorf("expected host path 'data/data', got %q", vols[0].HostPath)
	}
	if vols[0].ContainerPath != "/var/lib/postgresql/data" {
		t.Errorf("expected container path preserved, got %q", vols[0].ContainerPath)
	}
	if vols[1].HostPath != "data/uploads" {
		t.Errorf("expected 'data/uploads', got %q", vols[1].HostPath)
	}
}

func TestProposeVolumesEmpty(t *testing.T) {
	vols := ProposeVolumes(nil)
	if len(vols) != 0 {
		t.Errorf("expected 0, got %d", len(vols))
	}
}

func TestProposeVolumesRoot(t *testing.T) {
	vols := ProposeVolumes([]string{"/"})
	if vols[0].HostPath != "data/data" {
		t.Errorf("root volume should use 'data/data', got %q", vols[0].HostPath)
	}
}

func TestVolumeName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/var/lib/data", "data"},
		{"/tmp", "tmp"},
		{"/", "data"},
		{"", "data"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := volumeName(tt.path)
			if got != tt.want {
				t.Errorf("volumeName(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestConfirmVolumeMappings(t *testing.T) {
	orig := InputFunc
	defer func() { InputFunc = orig }()
	InputFunc = func(_, _, value string) (string, error) {
		return value, nil
	}

	proposed := ProposeVolumes([]string{"/data"})
	result, err := ConfirmVolumeMappings(proposed)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result[0].HostPath != "data/data" {
		t.Errorf("expected default 'data/data', got %q", result[0].HostPath)
	}
}

func TestConfirmVolumeMappingsOverride(t *testing.T) {
	orig := InputFunc
	defer func() { InputFunc = orig }()
	InputFunc = func(_, _, _ string) (string, error) {
		return "/custom/path", nil
	}

	proposed := ProposeVolumes([]string{"/data"})
	result, err := ConfirmVolumeMappings(proposed)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result[0].HostPath != "/custom/path" {
		t.Errorf("expected '/custom/path', got %q", result[0].HostPath)
	}
}

func TestConfirmVolumeMappingsEmptyFallback(t *testing.T) {
	orig := InputFunc
	defer func() { InputFunc = orig }()
	InputFunc = func(_, _, _ string) (string, error) {
		return "  ", nil
	}

	proposed := ProposeVolumes([]string{"/data"})
	result, err := ConfirmVolumeMappings(proposed)
	if err != nil {
		t.Fatalf("error: %v", err)
	}
	if result[0].HostPath != "data/data" {
		t.Errorf("empty input should fall back to default, got %q", result[0].HostPath)
	}
}
