// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"fmt"
	"io"

	"github.com/amrelsagaei/updock/internal/hub"
	"github.com/amrelsagaei/updock/internal/project"
	"github.com/amrelsagaei/updock/internal/prompt"
	"github.com/amrelsagaei/updock/internal/recipe"
	"github.com/amrelsagaei/updock/internal/scaffold"
	"github.com/amrelsagaei/updock/internal/security"
	"github.com/amrelsagaei/updock/internal/ui"
)

// imageSource is the subset of the Docker Hub client the create flow needs.
// *hub.Client satisfies it; tests provide a fake.
type imageSource interface {
	Search(query string, limit int) ([]hub.SearchResult, error)
	Tags(repo string, limit int) ([]hub.Tag, error)
	Inspect(image, tag string) (*hub.ImageConfig, error)
}

// createOptions carries everything the create flow needs.
type createOptions struct {
	name         string
	projectName  string
	projectsRoot string
	autoGenPass  bool
	source       imageSource
	recipes      map[string]*recipe.Recipe
	out          io.Writer
}

// createProject runs the full create pipeline and returns the project path.
// It dispatches to the recipe flow if a recipe matches, otherwise the
// single-image flow. All interactive steps go through the mockable prompt
// package functions, so this is fully testable.
func createProject(opts *createOptions) (string, error) {
	if r := recipe.Match(opts.recipes, opts.name); r != nil {
		return createFromRecipe(r, opts)
	}
	return createFromImage(opts)
}

func createFromRecipe(r *recipe.Recipe, opts *createOptions) (string, error) {
	ui.Info(opts.out, "Using recipe: %s - %s", ui.Bold(r.Meta.Name), r.Meta.Description)

	values, err := recipe.CollectValues(r, opts.autoGenPass)
	if err != nil {
		return "", err
	}

	base := opts.projectName
	if base == "" {
		base = r.Meta.Name
	}
	name := project.UniqueProjectName(opts.projectsRoot, base)

	projectPath, err := project.CreateProjectDir(opts.projectsRoot, name)
	if err != nil {
		return "", err
	}

	if err := recipe.Render(r, values, projectPath, name); err != nil {
		return "", err
	}

	ui.Success(opts.out, "Created project '%s' (%d services)", name, len(r.Services))
	return projectPath, nil
}

func createFromImage(opts *createOptions) (string, error) {
	if err := security.ValidateImageName(opts.name); err != nil {
		return "", err
	}

	ui.Step(opts.out, "Searching Docker Hub for %q...", opts.name)
	results, err := opts.source.Search(opts.name, 25)
	if err != nil {
		return "", err
	}
	if len(results) == 0 {
		return "", fmt.Errorf("no images found for %q - check the spelling or try a shorter query", opts.name)
	}

	ranked := hub.Rank(results, opts.name)
	sel, err := prompt.PickImage(ranked)
	if err != nil {
		return "", err
	}

	image := sel.Result.RepoName
	tag := "latest"
	if sel.ChooseVersion {
		tag, err = chooseVersion(opts, image)
		if err != nil {
			return "", err
		}
	}

	imgCfg := inspectOrWarn(opts, image, tag)

	cfg, err := configureImage(opts, image, tag, imgCfg)
	if err != nil {
		return "", err
	}

	base := opts.projectName
	if base == "" {
		base = project.DeriveProjectName(image)
	}
	name := project.UniqueProjectName(opts.projectsRoot, base)
	cfg.ProjectName = name

	projectPath, err := project.CreateProjectDir(opts.projectsRoot, name)
	if err != nil {
		return "", err
	}

	if err := scaffold.WriteAll(projectPath, cfg); err != nil {
		return "", err
	}

	ui.Success(opts.out, "Created project '%s' (%s:%s)", name, image, tag)
	return projectPath, nil
}

func chooseVersion(opts *createOptions, image string) (string, error) {
	tags, err := opts.source.Tags(image, 100)
	if err != nil {
		return "", err
	}
	if len(tags) == 0 {
		return "latest", nil
	}
	sorted := hub.SortTags(tags)
	picked, err := prompt.PickTag(sorted)
	if err != nil {
		return "", err
	}
	return picked.Name, nil
}

func inspectOrWarn(opts *createOptions, image, tag string) *hub.ImageConfig {
	imgCfg, err := opts.source.Inspect(image, tag)
	if err != nil {
		ui.Warn(opts.out, "Could not read image metadata (%v).", err)
		ui.Info(opts.out, "Continuing with a minimal configuration.")
		return &hub.ImageConfig{EnvDefaults: map[string]string{}}
	}
	return imgCfg
}

func configureImage(opts *createOptions, image, tag string, imgCfg *hub.ImageConfig) (*project.Config, error) {
	ports, err := prompt.ConfirmPortMappings(prompt.ProposePorts(imgCfg.ExposedPorts))
	if err != nil {
		return nil, err
	}

	envVars, err := prompt.CollectEnvVars(prompt.ClassifyEnvVars(image, imgCfg), opts.autoGenPass)
	if err != nil {
		return nil, err
	}

	volumes, err := prompt.ConfirmVolumeMappings(prompt.ProposeVolumes(imgCfg.Volumes))
	if err != nil {
		return nil, err
	}

	return &project.Config{
		Image:   image,
		Tag:     tag,
		Ports:   ports,
		Env:     envVars,
		Volumes: volumes,
		Labels:  imgCfg.Labels,
	}, nil
}
