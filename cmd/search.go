package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func search(cmd *cobra.Command, args []string) {
	//fmt.Println("search called")

}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "searches for the specified package in AUR via the AUR RPC interface",
	Long: `Searches for the specified package in AUR via the AUR RPC interface.
Example usage:
  aurora search <package-name>
  This command will return a list of packages that match the search query, along with their details such as package name, version, description, and popularity.
  Currently, this command only supports a single package name as an argument, but in the future, it may be extended to support multiple package names or additional search parameters.`,
	Run: search,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
