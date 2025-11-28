package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print gbpm version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("gbpm version", version)
		},
	}
}
