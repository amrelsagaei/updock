// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package scaffold

import (
	"fmt"

	"github.com/amrelsagaei/updock/internal/project"
)

// WriteAll generates all project files: docker-compose.yml, .env, .gitignore, updock.json.
func WriteAll(projectPath string, cfg *project.Config) error {
	if err := WriteComposeFile(projectPath, cfg); err != nil {
		return fmt.Errorf("scaffold: %w", err)
	}

	if err := WriteEnvFile(projectPath, cfg); err != nil {
		return fmt.Errorf("scaffold: %w", err)
	}

	if err := WriteGitignore(projectPath); err != nil {
		return fmt.Errorf("scaffold: %w", err)
	}

	meta := project.NewMetadata(cfg)
	if err := project.WriteMetadata(projectPath, &meta); err != nil {
		return fmt.Errorf("scaffold: %w", err)
	}

	return nil
}
