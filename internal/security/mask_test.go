// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package security

import (
	"strings"
	"testing"
)

func TestMaskSecret(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		secret   bool
		wantMask bool
	}{
		{"secret masked", "my-password", true, true},
		{"non-secret shown", "public-value", false, false},
		{"empty secret", "", true, true},
		{"empty non-secret", "", false, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MaskSecret(tt.value, tt.secret)
			if tt.wantMask {
				if !strings.Contains(got, "•") {
					t.Errorf("expected masked output, got %q", got)
				}
				if strings.Contains(got, tt.value) && tt.value != "" {
					t.Errorf("masked output should not contain original value %q", tt.value)
				}
			} else if got != tt.value {
				t.Errorf("expected %q, got %q", tt.value, got)
			}
		})
	}
}

func TestMaskSecretNeverLeaks(t *testing.T) {
	secrets := []string{"P@ssw0rd!", "super-secret-token", "api-key-12345"}
	for _, secret := range secrets {
		masked := MaskSecret(secret, true)
		if strings.Contains(masked, secret) {
			t.Errorf("masked output leaked secret: %q", secret)
		}
	}
}
