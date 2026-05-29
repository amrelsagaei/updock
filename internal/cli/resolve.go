// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"fmt"
	"strconv"

	"github.com/amrelsagaei/updock/internal/config"
	"github.com/amrelsagaei/updock/internal/project"
)

func resolveProject(numberStr string) (project.Entry, error) {
	num, err := strconv.Atoi(numberStr)
	if err != nil {
		return project.Entry{}, fmt.Errorf("invalid project number %q - use a number from 'updock ls'", numberStr)
	}

	cfg, err := config.Load()
	if err != nil {
		return project.Entry{}, err
	}

	reg, err := project.NewRegistry(cfg.ProjectsRoot)
	if err != nil {
		return project.Entry{}, err
	}

	return reg.Resolve(num)
}
