// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package recipe

import "strings"

// Match finds a recipe by name. Returns nil if no recipe matches.
func Match(recipes map[string]*Recipe, name string) *Recipe {
	name = strings.ToLower(name)

	if r, ok := recipes[name]; ok {
		return r
	}

	if idx := strings.LastIndex(name, "/"); idx >= 0 {
		base := name[idx+1:]
		if r, ok := recipes[base]; ok {
			return r
		}
	}

	return nil
}
