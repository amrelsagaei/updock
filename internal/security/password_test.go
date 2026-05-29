// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package security

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	pw, err := GeneratePassword()
	if err != nil {
		t.Fatalf("GeneratePassword() error: %v", err)
	}

	if len(pw) < 32 {
		t.Errorf("password should be at least 32 chars, got %d", len(pw))
	}
}

func TestGeneratePasswordUniqueness(t *testing.T) {
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		pw, err := GeneratePassword()
		if err != nil {
			t.Fatalf("iteration %d: %v", i, err)
		}
		if seen[pw] {
			t.Fatalf("duplicate password generated at iteration %d", i)
		}
		seen[pw] = true
	}
}

func TestGeneratePasswordLength(t *testing.T) {
	pw, err := GeneratePassword()
	if err != nil {
		t.Fatal(err)
	}
	if len(pw) != 32 {
		t.Errorf("expected exactly 32 chars, got %d", len(pw))
	}
}
