// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GeneratePassword creates a cryptographically strong random password.
// Returns a 32-character base64url-encoded string (~192 bits of entropy).
func GeneratePassword() (string, error) {
	b := make([]byte, 24) // 24 bytes = 192 bits -> 32 base64 chars
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating password: %w", err)
	}
	return base64.URLEncoding.EncodeToString(b), nil
}
