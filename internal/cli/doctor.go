package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/Foggy-Forge/git-bash-package-manager/internal/paths"
)

func newDoctorCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check environment and configuration",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Running gbpm diagnostics...")

			home, err := os.UserHomeDir()
			if err != nil {
				return fmt.Errorf("cannot determine home directory: %w", err)
			}
			fmt.Println("Home:", home)

			p := paths.NewDefault()
			fmt.Println("GBPM_HOME:", p.Home)
			fmt.Println("GBPM_BIN:", p.Bin)
			fmt.Println("GBPM_CACHE:", p.Cache)
			fmt.Println("GBPM_REGISTRY:", p.Registry)

			pathEnv := os.Getenv("PATH")
			if !isInPath(pathEnv, p.Bin) {
				fmt.Println()
				fmt.Println("WARNING: gbpm bin directory is not in PATH.")
				fmt.Println("Add the following line to your ~/.bashrc or ~/.profile:")
				fmt.Printf("    export PATH=\"%s:$PATH\"\n", p.Bin)
			} else {
				fmt.Println("OK: gbpm bin directory is in PATH.")
			}

			return nil
		},
	}
}

func isInPath(pathEnv, dir string) bool {
	for _, p := range filepath.SplitList(pathEnv) {
		if filepath.Clean(p) == filepath.Clean(dir) {
			return true
		}
	}
	return false
}
