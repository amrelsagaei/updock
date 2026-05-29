// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package security

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestNoShellCommandPatterns(t *testing.T) {
	root := findProjectRoot(t)

	forbidden := []string{
		`exec.Command("sh"`,
		`exec.Command("bash"`,
		`exec.Command("/bin/sh"`,
		`exec.Command("/bin/bash"`,
	}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && (info.Name() == ".git" || info.Name() == "testdata") {
			return filepath.SkipDir
		}
		if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
			return nil
		}

		data, readErr := os.ReadFile(path)
		if readErr != nil {
			return readErr
		}
		content := string(data)

		for _, pattern := range forbidden {
			if strings.Contains(content, pattern) {
				t.Errorf("SECURITY: %s contains forbidden pattern %q - use exec.Command with arg arrays, never shell invocation", path, pattern)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatalf("walk error: %v", err)
	}
}

func TestNoSecretsInComposeTemplate(t *testing.T) {
	root := findProjectRoot(t)
	composePath := filepath.Join(root, "internal", "scaffold", "compose.go")

	data, err := os.ReadFile(composePath)
	if err != nil {
		t.Skipf("compose.go not found: %v", err)
	}
	content := string(data)

	if strings.Contains(content, "e.Value") && !strings.Contains(content, "${") {
		t.Error("SECURITY: compose.go should never inline env values - use ${VAR} references only")
	}

	if strings.Contains(content, ".Secret") {
		t.Error("SECURITY: compose.go should not reference Secret field - it should never touch secret values")
	}
}

func findProjectRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root (go.mod)")
		}
		dir = parent
	}
}
