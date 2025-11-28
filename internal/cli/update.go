package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Foggy-Forge/git-bash-package-manager/internal/paths"
	"github.com/Foggy-Forge/git-bash-package-manager/internal/registry"
)

func newUpdateCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "update",
		Short: "Update the package registry",
		Long:  "Clone or update the package registry from GitHub.",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := paths.NewDefault()

			reg, err := registry.New(p.Registry)
			if err != nil {
				return fmt.Errorf("failed to create registry: %w", err)
			}

			return reg.Pull()
		},
	}
}
