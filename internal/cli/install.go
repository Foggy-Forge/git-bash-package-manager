package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Foggy-Forge/git-bash-package-manager/internal/installer"
	"github.com/Foggy-Forge/git-bash-package-manager/internal/manifest"
	"github.com/Foggy-Forge/git-bash-package-manager/internal/paths"
	"github.com/Foggy-Forge/git-bash-package-manager/internal/registry"
)

func newInstallCmd() *cobra.Command {
	var manifestFile string

	cmd := &cobra.Command{
		Use:   "install [package]",
		Short: "Install a package",
		Long: `Install a package from the registry or from a local manifest file.

Examples:
  gbpm install fzf              # Install from registry
  gbpm install --file fzf.yaml  # Install from local manifest`,
		RunE: func(cmd *cobra.Command, args []string) error {
			p := paths.NewDefault()
			statePath := filepath.Join(p.Home, "state.json")

			inst, err := installer.New(p, statePath)
			if err != nil {
				return fmt.Errorf("failed to create installer: %w", err)
			}

			// Install from file
			if manifestFile != "" {
				m, err := manifest.LoadManifest(manifestFile)
				if err != nil {
					return fmt.Errorf("failed to load manifest: %w", err)
				}

				return inst.Install(m)
			}

			// Install from registry
			if len(args) == 0 {
				return fmt.Errorf("package name or --file required")
			}

			packageName := args[0]

			// Load registry
			reg, err := registry.New(p.Registry)
			if err != nil {
				return fmt.Errorf("failed to load registry: %w", err)
			}

			// Find manifest
			manifestPath, err := reg.FindManifest(packageName)
			if err != nil {
				return fmt.Errorf("package not found: %w\n\nRun 'gbpm update' to update the package registry", err)
			}

			m, err := manifest.LoadManifest(manifestPath)
			if err != nil {
				return fmt.Errorf("failed to load manifest: %w", err)
			}

			return inst.Install(m)
		},
	}

	cmd.Flags().StringVarP(&manifestFile, "file", "f", "", "Install from a local manifest file")

	return cmd
}
