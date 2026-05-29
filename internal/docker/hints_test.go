// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"testing"
)

func TestGetHint(t *testing.T) {
	tests := []struct {
		image   string
		wantOK  bool
		wantReq int
	}{
		{"postgres", true, 1},
		{"library/postgres", true, 1},
		{"bitnami/mysql", true, 1},
		{"mysql", true, 1},
		{"redis", true, 0},
		{"wordpress", true, 1},
		{"mariadb", true, 1},
		{"mongo", true, 0},
		{"gitea", true, 0},
		{"nextcloud", true, 0},
		{"unknown-image", false, 0},
		{"POSTGRES", true, 1},
	}

	for _, tt := range tests {
		t.Run(tt.image, func(t *testing.T) {
			hint, ok := GetHint(tt.image)
			if ok != tt.wantOK {
				t.Errorf("GetHint(%q) ok = %v, want %v", tt.image, ok, tt.wantOK)
			}
			if ok && len(hint.Required) != tt.wantReq {
				t.Errorf("GetHint(%q) required = %d, want %d", tt.image, len(hint.Required), tt.wantReq)
			}
		})
	}
}

func TestIsSecretKey(t *testing.T) {
	tests := []struct {
		key  string
		want bool
	}{
		{"POSTGRES_PASSWORD", true},
		{"MYSQL_ROOT_PASSWORD", true},
		{"API_SECRET", true},
		{"AUTH_TOKEN", true},
		{"AWS_ACCESS_KEY", true},
		{"DB_CREDENTIALS", true},
		{"DATABASE_URL", false},
		{"POSTGRES_USER", false},
		{"NODE_ENV", false},
		{"PORT", false},
		{"password", true},
	}

	for _, tt := range tests {
		t.Run(tt.key, func(t *testing.T) {
			got := IsSecretKey(tt.key)
			if got != tt.want {
				t.Errorf("IsSecretKey(%q) = %v, want %v", tt.key, got, tt.want)
			}
		})
	}
}

func TestHintDescriptions(t *testing.T) {
	hint, ok := GetHint("postgres")
	if !ok {
		t.Fatal("expected postgres hint")
	}
	if hint.Description["POSTGRES_PASSWORD"] == "" {
		t.Error("expected description for POSTGRES_PASSWORD")
	}
}
