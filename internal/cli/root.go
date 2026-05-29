// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"fmt"
	"path/filepath"

	"github.com/amrelsagaei/updock/internal/config"
	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/hub"
	"github.com/amrelsagaei/updock/internal/prompt"
	"github.com/amrelsagaei/updock/internal/recipe"
	"github.com/amrelsagaei/updock/internal/ui"
	"github.com/spf13/cobra"
)

// NewRootCmd creates the top-level updock command with all subcommands.
func NewRootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "updock [name]",
		Short: "Run any Docker app from one word",
		Long: `updock finds the image, lets you pick a version, asks you the
few things that actually matter, writes the Compose file and .env
for you, and brings it up. After that you control everything by number.`,
		Args:                  cobra.ArbitraryArgs,
		DisableFlagsInUseLine: true,
		RunE:                  runRoot,
	}

	root.Flags().String("name", "", "Override the project name")

	root.AddCommand(
		newSearchCmd(),
		newLsCmd(),
		newStatusCmd(),
		newUpCmd(),
		newStopCmd(),
		newRestartCmd(),
		newRebuildCmd(),
		newDownCmd(),
		newLogsCmd(),
		newOpenCmd(),
		newConfigCmd(),
		newRmCmd(),
		newDoctorCmd(),
		newVersionCmd(),
	)

	return root
}

func runRoot(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	name := args[0]
	projectNameOverride, _ := cmd.Flags().GetString("name")
	w := cmd.OutOrStdout()

	if err := docker.RunPreflightFast(); err != nil {
		return fmt.Errorf("%w\n  Run 'updock doctor' for the full check and how to fix it", err)
	}

	cfg, err := config.Load()
	if err != nil {
		return err
	}

	recipes, err := recipe.LoadAll(userRecipesDir())
	if err != nil {
		return err
	}

	projectPath, err := createProject(&createOptions{
		name:         name,
		projectName:  projectNameOverride,
		projectsRoot: cfg.ProjectsRoot,
		autoGenPass:  cfg.AutoGeneratePasswords,
		source:       hub.NewClient(),
		recipes:      recipes,
		out:          w,
	})
	if err != nil {
		return err
	}

	start, err := prompt.Confirm("Start it now?")
	if err != nil {
		return err
	}
	if !start {
		ui.Info(w, "Created at %s", ui.Dim(projectPath))
		ui.Info(w, "Run 'updock ls', then 'updock up <n>' when ready.")
		return nil
	}

	projectName := filepath.Base(projectPath)
	ui.Step(w, "Pulling images...")
	if err := docker.Pull(projectPath, projectName, w, cmd.ErrOrStderr()); err != nil {
		return err
	}
	if err := docker.Up(projectPath, projectName, w, cmd.ErrOrStderr()); err != nil {
		return err
	}

	ui.Success(w, "%s is running.", projectName)
	return nil
}
