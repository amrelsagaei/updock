// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package version

import "testing"

func TestDefaults(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{"Version", Version},
		{"Commit", Commit},
		{"Date", Date},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.value == "" {
				t.Errorf("%s should have a non-empty default (got empty) so 'updock version' never prints a blank field", tt.name)
			}
		})
	}
}
