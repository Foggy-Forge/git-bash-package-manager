package cli

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Foggy-Forge/git-bash-package-manager/internal/paths"
	"github.com/Foggy-Forge/git-bash-package-manager/internal/state"
)

func newListCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List installed packages",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			p := paths.NewDefault()
			statePath := filepath.Join(p.Home, "state.json")

			s, err := state.Load(statePath)
			if err != nil {
				return fmt.Errorf("failed to load state: %w", err)
			}

			if len(s.Installed) == 0 {
				fmt.Println("No packages installed.")
				return nil
			}

			fmt.Println("Installed packages:")
			for name, pkg := range s.Installed {
				fmt.Printf("  %s v%s (installed: %s)\n", 
					name, 
					pkg.Version, 
					pkg.InstalledAt.Format("2006-01-02"))
			}

			return nil
		},
	}
}
