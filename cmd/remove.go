package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/abhigyan-chatterjee/aurora/internal/pacman"

	"github.com/spf13/cobra"
)

func remove(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please specify a package name to remove.")
		return
	}

	pkgQuery := args[0]
	reader := bufio.NewReader(os.Stdin)

	matches, err := pacman.SearchInstalled(pkgQuery)
	if err != nil {
		fmt.Printf("Error checking installed packages: %v\n", err)
		return
	}

	var pkgName string

	if len(matches) == 0 {
		fmt.Printf("No installed package found matching \"%s\".\n", pkgQuery)
		return
	} else if len(matches) > 1 {
		fmt.Printf("Found multiple packages matching \"%s\":\n", pkgQuery)
		for _, m := range matches {
			fmt.Printf("  - %s\n", m)
		}
		fmt.Print("\nEnter the exact package name to remove: ")
		input, _ := reader.ReadString('\n')
		pkgName = strings.TrimSpace(input)

		found := false
		for _, m := range matches {
			if m == pkgName {
				found = true
				break
			}
		}
		if !found {
			fmt.Printf("\"%s\" does not match any of the listed packages.\n", pkgName)
			return
		}
	} else {
		pkgName = matches[0]
	}

	fmt.Printf("\nAurora is about to remove: %s\n", pkgName)
	fmt.Print("Are you sure? (y/N): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" {
		fmt.Println("Removal cancelled.")
		return
	}

	err = pacman.RemovePacman(pkgName)
	if err != nil {
		fmt.Printf("Error: removal failed: %v\n", err)
	} else {
		fmt.Printf("Successfully removed %s.\n", pkgName)
	}
}

var removeCmd = &cobra.Command{
	Use:   "remove [package]",
	Short: "Removes the specified package from the system",
	Run:   remove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
