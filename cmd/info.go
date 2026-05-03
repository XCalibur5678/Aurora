package cmd

import (
	"fmt"

	"github.com/abhigyan-chatterjee/aurora/internal/aur"
	"github.com/abhigyan-chatterjee/aurora/internal/pacman"

	"github.com/spf13/cobra"
)

func info(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a package name.")
		return
	}

	packageName := args[0]

	pacmanPkg, _ := pacman.SearchPacmanExact(packageName)
	if pacmanPkg != nil {
		fmt.Printf("\n--- Package Details ---\n")
		fmt.Printf("Name        : %s\n", pacmanPkg.Name)
		fmt.Printf("Version     : %s\n", pacmanPkg.Version)
		fmt.Printf("Repository  : %s\n", pacmanPkg.Repository)
		fmt.Printf("Description : %s\n", pacmanPkg.Description)
		fmt.Printf("Source      : official repositories\n")
		return
	}

	aurPkg, _ := aur.SearchAURExact(packageName)
	if aurPkg != nil {
		fmt.Printf("\n--- Package Details ---\n")
		fmt.Printf("Name        : %s\n", aurPkg.Name)
		fmt.Printf("Version     : %s\n", aurPkg.Version)
		fmt.Printf("Description : %s\n", aurPkg.Description)
		fmt.Printf("URL         : %s\n", aurPkg.URL)
		fmt.Printf("Votes       : %d\n", aurPkg.NumVotes)
		fmt.Printf("Source      : AUR\n")
		fmt.Println("\nNote: installing from AUR means Aurora will clone the build files and build locally using makepkg.")
		return
	}

	fmt.Printf("Package \"%s\" was not found in official repositories or the AUR.\n", packageName)
}

var infoCmd = &cobra.Command{
	Use:   "info [package]",
	Short: "Show detailed information about a package",
	Run:   info,
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
