// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package security

const maskedValue = "••••••"

// MaskSecret returns the masked placeholder if secret is true,
// otherwise returns the value as-is.
func MaskSecret(value string, secret bool) string {
	if secret {
		return maskedValue
	}
	return value
}
