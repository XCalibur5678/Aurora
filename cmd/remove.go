package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func remove(cmd *cobra.Command, args []string) {
	fmt.Println("remove called")
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the specified package from the system ",
	Run:   remove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
