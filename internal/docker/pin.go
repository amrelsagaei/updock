// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package docker

import (
	"fmt"
	"strings"
)

// ResolveDigest fetches the digest for an image:tag from the registry.
// Returns the digest string (e.g., "sha256:abc123...").
func ResolveDigest(image, tag string) (string, error) {
	out, err := CommandRunner.Output("", "docker", "inspect", "--format", "{{index .RepoDigests 0}}", image+":"+tag)
	if err != nil {
		out, err = CommandRunner.Output("", "docker", "manifest", "inspect", "--verbose", image+":"+tag)
		if err != nil {
			return "", fmt.Errorf("resolving digest for %s:%s: %w", image, tag, err)
		}
	}

	digest := strings.TrimSpace(string(out))
	if idx := strings.Index(digest, "@"); idx >= 0 {
		digest = digest[idx+1:]
	}

	if !strings.HasPrefix(digest, "sha256:") {
		return "", fmt.Errorf("unexpected digest format: %q", digest)
	}

	return digest, nil
}

// PinnedImageRef returns an image reference pinned to a specific digest.
func PinnedImageRef(image, digest string) string {
	return fmt.Sprintf("%s@%s", image, digest)
}
