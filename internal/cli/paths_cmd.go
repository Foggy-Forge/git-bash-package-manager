package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/Foggy-Forge/git-bash-package-manager/internal/paths"
)

func newPathsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "paths",
		Short: "Show gbpm paths",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := paths.NewDefault()
			fmt.Println("GBPM_HOME:", p.Home)
			fmt.Println("GBPM_BIN:", p.Bin)
			fmt.Println("GBPM_CACHE:", p.Cache)
			fmt.Println("GBPM_REGISTRY:", p.Registry)
			return nil
		},
	}
}
