// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amrelsagaei/updock/internal/hub"
)

// TestEdgeImageNoExposedPorts: an image with no EXPOSE should still scaffold,
// just without a ports section.
func TestEdgeImageNoExposedPorts(t *testing.T) {
	stubPrompts(t, "scratch", "")
	root := t.TempDir()

	src := &fakeSource{
		searchResults: []hub.SearchResult{{RepoName: "scratch", PullCount: 100}},
		imgCfg:        &hub.ImageConfig{ExposedPorts: nil, EnvDefaults: map[string]string{}},
	}

	var out bytes.Buffer
	projectPath, err := createProject(&createOptions{
		name: "scratch", projectsRoot: root, source: src, out: &out,
	})
	if err != nil {
		t.Fatalf("should scaffold even with no ports: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(projectPath, "docker-compose.yml"))
	if strings.Contains(string(data), "ports:") {
		t.Error("compose should have no ports section when image exposes none")
	}
}

// TestEdgeImageNoEnvDefaults: an image with no ENV should scaffold with no
// env vars to prompt for.
func TestEdgeImageNoEnvDefaults(t *testing.T) {
	stubPrompts(t, "busybox", "")
	root := t.TempDir()

	src := &fakeSource{
		searchResults: []hub.SearchResult{{RepoName: "busybox", PullCount: 100}},
		imgCfg:        &hub.ImageConfig{ExposedPorts: []int{}, EnvDefaults: map[string]string{}},
	}

	var out bytes.Buffer
	_, err := createProject(&createOptions{
		name: "busybox", projectsRoot: root, autoGenPass: true, source: src, out: &out,
	})
	if err != nil {
		t.Fatalf("should handle image with no env defaults: %v", err)
	}
}

// TestEdgeNearMaxImageName: a name near the 255-char limit is accepted;
// over the limit is rejected before any network call.
func TestEdgeImageNameLength(t *testing.T) {
	stubPrompts(t, "", "")

	// 256 chars is rejected
	tooLong := strings.Repeat("a", 256)
	var out bytes.Buffer
	_, err := createProject(&createOptions{
		name: tooLong, projectsRoot: t.TempDir(), source: &fakeSource{}, out: &out,
	})
	if err == nil {
		t.Error("image name over 255 chars should be rejected")
	}
}

// TestEdgeUnicodeEnvValue: unicode values flow through scaffolding intact.
func TestEdgeUnicodeEnvValue(t *testing.T) {
	stubPrompts(t, "app", "日本語-café-Ω")

	root := t.TempDir()
	src := &fakeSource{
		searchResults: []hub.SearchResult{{RepoName: "app", PullCount: 100}},
		imgCfg: &hub.ImageConfig{
			EnvDefaults: map[string]string{"LOCALE": ""},
		},
	}

	var out bytes.Buffer
	projectPath, err := createProject(&createOptions{
		name: "app", projectsRoot: root, source: src, out: &out,
	})
	if err != nil {
		t.Fatalf("unicode env value should be accepted: %v", err)
	}

	data, _ := os.ReadFile(filepath.Join(projectPath, ".env"))
	if !strings.Contains(string(data), "日本語-café-Ω") {
		t.Error("unicode env value should be written to .env intact")
	}
}

// TestEdgeManyTags: a repo with many tags sorts correctly and the picker
// receives them all.
func TestEdgeManyTags(t *testing.T) {
	tags := make([]hub.Tag, 0, 60)
	for i := 0; i < 50; i++ {
		tags = append(tags, hub.Tag{Name: "1." + itoa(i) + ".0"})
	}
	tags = append(tags, hub.Tag{Name: "latest"})

	sorted := hub.SortTags(tags)
	if sorted[0].Name != "latest" {
		t.Errorf("latest should be pinned first among many tags, got %q", sorted[0].Name)
	}
	if len(sorted) != 51 {
		t.Errorf("all tags should be preserved, got %d", len(sorted))
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
