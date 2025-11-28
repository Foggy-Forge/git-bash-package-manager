package cli

import (
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
)

func Execute() error {
	return newRootCmd().Execute()
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gbpm",
		Short: "gbpm is a lightweight package manager for Git Bash",
		Long:  "gbpm installs and manages CLI tools and scripts for Git Bash on Windows.",
	}

	cmd.AddCommand(
		newVersionCmd(),
		newDoctorCmd(),
		newPathsCmd(),
		newInstallCmd(),
		newUninstallCmd(),
		newListCmd(),
		newUpdateCmd(),
		newUpgradeCmd(),
	)

	return cmd
}
