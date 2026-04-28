package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var search = func(cmd *cobra.Command, args []string) {
	//if no argument is provided, print a message and exit
	if len(args) == 0 {
		fmt.Println("Please provide a package name to search for.")
		return
	}
	packageName := args[0]
	_, err := searchPackage(packageName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

var searchCmd = &cobra.Command{
	Use:   "search [package]",
	Short: "searches for the specified package in AUR via the AUR RPC interface",
	Run:   search,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
