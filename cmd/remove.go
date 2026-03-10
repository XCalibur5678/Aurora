package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

func listInstalled(query string) ([]string, error) {

	cmd := exec.Command("pacman", "-Qq")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	var matches []string
	for _, line := range lines {
		if strings.Contains(line, query) {
			matches = append(matches, line)
		}
	}
	return matches, nil
}

func remove(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Error: Please specify a package name to remove.")
		return
	}

	pkgQuery := args[0]

	matches, err := listInstalled(pkgQuery)
	if err != nil {
		fmt.Printf("Error checking installed packages: %v\n", err)
		return
	}

	var pkgName string
	if len(matches) == 0 {
		fmt.Printf("No installed package found matching '%s'.\n", pkgQuery)
		return
	} else if len(matches) > 1 {
		fmt.Printf("Found multiple packages matching '%s':\n", pkgQuery)
		for _, m := range matches {
			fmt.Printf("- %s\n", m)
		}
		fmt.Print("Please be more specific: ")
		fmt.Scanln(&pkgName)
		found := false
		for _, m := range matches {
			if m == pkgName {
				found = true
				break
			}
		}
		if !found {
			fmt.Println("Invalid selection.")
			return
		}
	} else {
		pkgName = matches[0]
		fmt.Printf("Found package: %s\n", pkgName)
	}

	fmt.Printf("Aurora is about to remove: %s\n", pkgName)
	fmt.Printf("This will invoke: sudo pacman -Rns %s\n", pkgName)
	fmt.Print("Are you sure? (y/N): ")

	var confirm string
	fmt.Scanln(&confirm)
	if strings.ToLower(confirm) != "y" && strings.ToLower(confirm) != "Yes" {
		fmt.Println("Removal cancelled.")
		return
	}

	removeCmd := exec.Command("sudo", "pacman", "-Rns", pkgName)
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr
	removeCmd.Stdin = os.Stdin

	if err := removeCmd.Run(); err != nil {
		fmt.Printf("Error: Removal failed: %v\n", err)
	} else {
		fmt.Printf("Successfully removed %s\n", pkgName)
	}
}

var removeCmd = &cobra.Command{
	Use:   "remove [package]",
	Short: "removes the specified package from the system",
	Run:   remove,
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
