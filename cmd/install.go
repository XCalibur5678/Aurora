package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func install(cmd *cobra.Command, args []string) {
	fmt.Println("install called")
}

var installCmd = &cobra.Command{
	Use:   "install",
	Short: "installs the specified package on the system",
	Run:   install,
}

func init() {
	rootCmd.AddCommand(installCmd)
}
