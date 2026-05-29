// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package prompt

import (
	"fmt"

	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/hub"
	"github.com/amrelsagaei/updock/internal/project"
	"github.com/amrelsagaei/updock/internal/security"
)

// ClassifyEnvVars analyzes image config and hints to produce a categorized
// list of env vars with defaults, required flags, and secret detection.
func ClassifyEnvVars(image string, imgCfg *hub.ImageConfig) []project.EnvVar {
	vars := make([]project.EnvVar, 0)
	seen := make(map[string]bool)

	hint, hasHint := docker.GetHint(image)

	if hasHint {
		for _, key := range hint.Required {
			seen[key] = true
			vars = append(vars, project.EnvVar{
				Key:      key,
				Value:    imgCfg.EnvDefaults[key],
				Secret:   docker.IsSecretKey(key),
				Required: true,
			})
		}
		for _, key := range hint.Common {
			if seen[key] {
				continue
			}
			seen[key] = true
			vars = append(vars, project.EnvVar{
				Key:    key,
				Value:  imgCfg.EnvDefaults[key],
				Secret: docker.IsSecretKey(key),
			})
		}
	}

	for key, val := range imgCfg.EnvDefaults {
		if seen[key] {
			continue
		}
		if isBoringEnvVar(key) {
			continue
		}
		secret := docker.IsSecretKey(key)
		required := secret && val == ""
		vars = append(vars, project.EnvVar{
			Key:      key,
			Value:    val,
			Secret:   secret,
			Required: required,
		})
	}

	return vars
}

func isBoringEnvVar(key string) bool {
	boring := []string{"PATH", "GOPATH", "HOME", "LANG", "LC_ALL", "TERM", "HOSTNAME"}
	for _, b := range boring {
		if key == b {
			return true
		}
	}
	return false
}

// CollectEnvVars runs interactive prompts for each env var, generating
// passwords for required secrets.
func CollectEnvVars(vars []project.EnvVar, autoGenerate bool) ([]project.EnvVar, error) {
	result := make([]project.EnvVar, len(vars))

	for i, v := range vars {
		result[i] = v

		if v.Required && v.Secret && autoGenerate && v.Value == "" {
			pw, err := security.GeneratePassword()
			if err != nil {
				return nil, fmt.Errorf("generating password for %s: %w", v.Key, err)
			}
			result[i].Value = pw
			continue
		}

		if v.Required && v.Value == "" {
			title := fmt.Sprintf("%s (required)", v.Key)
			val, err := InputFunc(title, "", "")
			if err != nil {
				return nil, fmt.Errorf("env var %s: %w", v.Key, err)
			}
			if val == "" {
				return nil, fmt.Errorf("%s is required but was left empty", v.Key)
			}
			result[i].Value = val
			continue
		}

		if v.Value == "" {
			title := v.Key
			val, err := InputFunc(title, "(press enter to skip)", "")
			if err != nil {
				return nil, fmt.Errorf("env var %s: %w", v.Key, err)
			}
			result[i].Value = val
		}
	}

	return result, nil
}

// EditEnvVars re-prompts each env var with its current value pre-filled,
// letting the user change values during reconfiguration.
func EditEnvVars(vars []project.EnvVar) ([]project.EnvVar, error) {
	result := make([]project.EnvVar, len(vars))

	for i, v := range vars {
		result[i] = v

		title := v.Key
		if v.Secret {
			title = v.Key + " (secret)"
		}

		val, err := InputFunc(title, v.Value, v.Value)
		if err != nil {
			return nil, fmt.Errorf("env var %s: %w", v.Key, err)
		}
		result[i].Value = val
	}

	return result, nil
}
