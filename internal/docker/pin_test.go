// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"testing"
)

func TestResolveDigestSuccess(t *testing.T) {
	m := &mockRunner{output: []byte("nginx@sha256:abc123def456\n")}
	withMockRunner(m, func() {
		digest, err := ResolveDigest("nginx", "latest")
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if digest != "sha256:abc123def456" {
			t.Errorf("expected 'sha256:abc123def456', got %q", digest)
		}
	})
}

func TestResolveDigestRawSha(t *testing.T) {
	m := &mockRunner{output: []byte("sha256:abc123\n")}
	withMockRunner(m, func() {
		digest, err := ResolveDigest("nginx", "latest")
		if err != nil {
			t.Fatalf("error: %v", err)
		}
		if digest != "sha256:abc123" {
			t.Errorf("expected 'sha256:abc123', got %q", digest)
		}
	})
}

func TestResolveDigestBadFormat(t *testing.T) {
	m := &mockRunner{output: []byte("not-a-digest\n")}
	withMockRunner(m, func() {
		_, err := ResolveDigest("nginx", "latest")
		if err == nil {
			t.Error("expected error for bad format")
		}
	})
}

func TestPinnedImageRef(t *testing.T) {
	ref := PinnedImageRef("nginx", "sha256:abc123")
	if ref != "nginx@sha256:abc123" {
		t.Errorf("expected 'nginx@sha256:abc123', got %q", ref)
	}
}
