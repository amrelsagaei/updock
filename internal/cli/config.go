// Copyright (c) 2026 Amr Elsagaei. All rights reserved.

package cli

import (
	"os"
	"path/filepath"
	"time"

	"github.com/amrelsagaei/updock/internal/docker"
	"github.com/amrelsagaei/updock/internal/project"
	"github.com/amrelsagaei/updock/internal/prompt"
	"github.com/amrelsagaei/updock/internal/scaffold"
	"github.com/amrelsagaei/updock/internal/ui"
	"github.com/spf13/cobra"
)

func newConfigCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "config [number]",
		Short: "Re-run the config prompts for a project",
		Long: `Re-run the configuration prompts for the project with the given number
(from 'updock ls'). Current values are pre-filled - press enter to keep
them. Regenerates the .env file and offers to restart.`,
		Args: cobra.ExactArgs(1),
		RunE: runConfig,
	}
}

func runConfig(cmd *cobra.Command, args []string) error {
	entry, err := resolveProject(args[0])
	if err != nil {
		return err
	}

	cfg, err := loadProjectConfig(entry.Path)
	if err != nil {
		return err
	}

	w := cmd.OutOrStdout()
	ui.Info(w, "Reconfiguring '%s'. Press enter to keep current values.", entry.Name)

	if len(cfg.Env) > 0 {
		cfg.Env, err = prompt.EditEnvVars(cfg.Env)
		if err != nil {
			return err
		}
	}

	if len(cfg.Ports) > 0 {
		cfg.Ports, err = prompt.ConfirmPortMappings(cfg.Ports)
		if err != nil {
			return err
		}
	}

	// Regenerating .env is sufficient: the compose file references every
	// value via ${VAR} and ${HOST_PORT_n}, so changes apply without touching it.
	if err := scaffold.WriteEnvFile(entry.Path, cfg); err != nil {
		return err
	}

	if err := updatePortMetadata(entry.Path, cfg); err != nil {
		return err
	}

	ui.Success(w, "Configuration saved.")

	restart, err := prompt.Confirm("Restart to apply changes?")
	if err != nil {
		return err
	}
	if restart {
		if err := docker.Restart(entry.Path, entry.Name); err != nil {
			return err
		}
		ui.Success(w, "%s restarted.", entry.Name)
	}

	return nil
}

func updatePortMetadata(projectPath string, cfg *project.Config) error {
	meta, err := project.ReadMetadata(projectPath)
	if err != nil {
		return err
	}
	ports := make([]project.PortMetadata, len(cfg.Ports))
	for i, p := range cfg.Ports {
		ports[i] = project.PortMetadata{Host: p.Host, Container: p.Container}
	}
	meta.Ports = ports
	meta.UpdatedAt = time.Now().UTC()
	return project.WriteMetadata(projectPath, &meta)
}

// userRecipesDir returns the directory where user-supplied recipes live.
func userRecipesDir() string {
	if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
		return filepath.Join(xdg, "updock", "recipes")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".local", "share", "updock", "recipes")
}
