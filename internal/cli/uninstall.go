package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Foggy-Forge/git-bash-package-manager/internal/installer"
	"github.com/Foggy-Forge/git-bash-package-manager/internal/paths"
)

func newUninstallCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "uninstall <package>",
		Short: "Uninstall a package",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packageName := args[0]

			p := paths.NewDefault()
			statePath := filepath.Join(p.Home, "state.json")

			inst, err := installer.New(p, statePath)
			if err != nil {
				return fmt.Errorf("failed to create installer: %w", err)
			}

			return inst.Uninstall(packageName)
		},
	}
}
