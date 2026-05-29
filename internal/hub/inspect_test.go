// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package hub

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseConfigBlob(t *testing.T) {
	blob := imageConfigBlob{}
	blob.Config.ExposedPorts = map[string]struct{}{
		"80/tcp":  {},
		"443/tcp": {},
		"8080":    {},
	}
	blob.Config.Env = []string{
		"PATH=/usr/local/bin:/usr/bin",
		"POSTGRES_PASSWORD=",
		"POSTGRES_USER=postgres",
	}
	blob.Config.Labels = map[string]string{
		"maintainer": "test",
	}
	blob.Config.Volumes = map[string]struct{}{
		"/var/lib/data": {},
		"/tmp":          {},
	}
	blob.Config.Cmd = []string{"postgres"}
	blob.Config.Entrypoint = []string{"docker-entrypoint.sh"}

	cfg := parseConfigBlob(&blob)

	if len(cfg.ExposedPorts) != 3 {
		t.Errorf("expected 3 ports, got %d", len(cfg.ExposedPorts))
	}
	if cfg.ExposedPorts[0] != 80 {
		t.Errorf("ports should be sorted, first should be 80, got %d", cfg.ExposedPorts[0])
	}

	if cfg.EnvDefaults["POSTGRES_USER"] != "postgres" {
		t.Errorf("expected POSTGRES_USER=postgres, got %q", cfg.EnvDefaults["POSTGRES_USER"])
	}
	if val, ok := cfg.EnvDefaults["POSTGRES_PASSWORD"]; !ok || val != "" {
		t.Errorf("expected POSTGRES_PASSWORD with empty value, got %q (ok=%v)", val, ok)
	}

	if len(cfg.Volumes) != 2 {
		t.Errorf("expected 2 volumes, got %d", len(cfg.Volumes))
	}

	if cfg.Labels["maintainer"] != "test" {
		t.Errorf("expected label maintainer=test, got %q", cfg.Labels["maintainer"])
	}

	if len(cfg.Cmd) != 1 || cfg.Cmd[0] != "postgres" {
		t.Errorf("unexpected Cmd: %v", cfg.Cmd)
	}
	if len(cfg.Entrypoint) != 1 || cfg.Entrypoint[0] != "docker-entrypoint.sh" {
		t.Errorf("unexpected Entrypoint: %v", cfg.Entrypoint)
	}
}

func TestParseConfigBlobEmpty(t *testing.T) {
	empty := imageConfigBlob{}
	cfg := parseConfigBlob(&empty)

	if len(cfg.ExposedPorts) != 0 {
		t.Errorf("expected 0 ports, got %d", len(cfg.ExposedPorts))
	}
	if len(cfg.EnvDefaults) != 0 {
		t.Errorf("expected 0 env defaults, got %d", len(cfg.EnvDefaults))
	}
	if cfg.Labels == nil {
		t.Error("Labels should be initialized, not nil")
	}
}

func TestNormalizeRepo(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"nginx", "library/nginx"},
		{"bitnami/nginx", "bitnami/nginx"},
		{"a/b/c", "a/b/c"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := normalizeRepo(tt.input)
			if got != tt.want {
				t.Errorf("normalizeRepo(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestInspectWithMockRegistry(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(tokenResponse{Token: "test-token"})
	})

	mux.HandleFunc("/v2/library/nginx/manifests/latest", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
		_ = json.NewEncoder(w).Encode(manifest{
			SchemaVersion: 2,
			Config: struct {
				Digest string `json:"digest"`
			}{Digest: "sha256:abc123"},
		})
	})

	mux.HandleFunc("/v2/library/nginx/blobs/sha256:abc123", func(w http.ResponseWriter, _ *http.Request) {
		blob := imageConfigBlob{}
		blob.Config.ExposedPorts = map[string]struct{}{"80/tcp": {}}
		blob.Config.Env = []string{"NGINX_VERSION=1.25.0", "PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin"}
		_ = json.NewEncoder(w).Encode(blob)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClientWith(server.Client())
	cfg, err := client.inspectWith("nginx", "latest", server.URL+"/token", server.URL)
	if err != nil {
		t.Fatalf("Inspect error: %v", err)
	}

	if len(cfg.ExposedPorts) != 1 || cfg.ExposedPorts[0] != 80 {
		t.Errorf("expected port 80, got %v", cfg.ExposedPorts)
	}
	if cfg.EnvDefaults["NGINX_VERSION"] != "1.25.0" {
		t.Errorf("expected NGINX_VERSION=1.25.0, got %q", cfg.EnvDefaults["NGINX_VERSION"])
	}
}

func TestInspectManifestListFallback(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/token", func(w http.ResponseWriter, _ *http.Request) {
		_ = json.NewEncoder(w).Encode(tokenResponse{Token: "test-token"})
	})

	mux.HandleFunc("/v2/library/nginx/manifests/latest", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.list.v2+json")
		_ = json.NewEncoder(w).Encode(manifestList{
			SchemaVersion: 2,
			Manifests: []struct {
				Digest    string `json:"digest"`
				MediaType string `json:"mediaType"`
				Platform  struct {
					Architecture string `json:"architecture"`
					OS           string `json:"os"`
				} `json:"platform"`
			}{
				{
					Digest:    "sha256:plat1",
					MediaType: "application/vnd.docker.distribution.manifest.v2+json",
					Platform: struct {
						Architecture string `json:"architecture"`
						OS           string `json:"os"`
					}{Architecture: "amd64", OS: "linux"},
				},
			},
		})
	})

	mux.HandleFunc("/v2/library/nginx/manifests/sha256:plat1", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.docker.distribution.manifest.v2+json")
		_ = json.NewEncoder(w).Encode(manifest{
			SchemaVersion: 2,
			Config: struct {
				Digest string `json:"digest"`
			}{Digest: "sha256:cfg1"},
		})
	})

	mux.HandleFunc("/v2/library/nginx/blobs/sha256:cfg1", func(w http.ResponseWriter, _ *http.Request) {
		blob := imageConfigBlob{}
		blob.Config.ExposedPorts = map[string]struct{}{"80/tcp": {}}
		_ = json.NewEncoder(w).Encode(blob)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	client := NewClientWith(server.Client())
	cfg, err := client.inspectWith("nginx", "latest", server.URL+"/token", server.URL)
	if err != nil {
		t.Fatalf("Inspect error: %v", err)
	}

	if len(cfg.ExposedPorts) != 1 {
		t.Errorf("expected 1 port from manifest list, got %d", len(cfg.ExposedPorts))
	}
}

func TestInspectTokenError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer server.Close()

	client := NewClientWith(server.Client())
	_, err := client.inspectWith("nginx", "latest", server.URL+"/token", server.URL)
	if err == nil {
		t.Fatal("expected error for auth failure")
	}
}
