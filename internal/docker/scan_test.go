// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"testing"
)

func withLookPath(fn func(string) (string, error), test func()) {
	orig := lookPath
	lookPath = fn
	defer func() { lookPath = orig }()
	test()
}

func TestScanImageNoTool(t *testing.T) {
	withLookPath(func(_ string) (string, error) { return "", errNotFound }, func() {
		_, err := ScanImage("nginx", "latest")
		if err == nil {
			t.Error("expected error when no scanner is found")
		}
	})
}

func TestScanImageWithTrivyClean(t *testing.T) {
	withLookPath(func(name string) (string, error) {
		if name == "trivy" {
			return "/usr/bin/trivy", nil
		}
		return "", errNotFound
	}, func() {
		m := &mockRunner{output: []byte("Total: 0 (CRITICAL: 0, HIGH: 0)\nNo vulnerabilities")}
		withMockRunner(m, func() {
			res, err := ScanImage("nginx", "latest")
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if res.Tool != "trivy" {
				t.Errorf("expected tool 'trivy', got %q", res.Tool)
			}
		})
	})
}

func TestScanImageWithTrivyVulnerable(t *testing.T) {
	withLookPath(func(name string) (string, error) {
		if name == "trivy" {
			return "/usr/bin/trivy", nil
		}
		return "", errNotFound
	}, func() {
		m := &mockRunner{output: []byte("nginx (debian 12)\nCVE-2024-1234 CRITICAL openssl")}
		withMockRunner(m, func() {
			res, err := ScanImage("nginx", "latest")
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if !res.HasIssue {
				t.Error("expected HasIssue=true for CRITICAL finding")
			}
		})
	})
}

func TestScanImageWithDockerScout(t *testing.T) {
	withLookPath(func(name string) (string, error) {
		if name == "docker" {
			return "/usr/bin/docker", nil
		}
		return "", errNotFound
	}, func() {
		m := &mockRunner{output: []byte("0 vulnerabilities found")}
		withMockRunner(m, func() {
			res, err := ScanImage("nginx", "latest")
			if err != nil {
				t.Fatalf("error: %v", err)
			}
			if res.Tool != "docker scout" {
				t.Errorf("expected 'docker scout', got %q", res.Tool)
			}
		})
	})
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		input  string
		maxLen int
		want   string
	}{
		{"short", 10, "short"},
		{"hello world", 5, "hello..."},
		{"", 5, ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

var errNotFound = &notFoundErr{}

type notFoundErr struct{}

func (e *notFoundErr) Error() string { return "not found" }
