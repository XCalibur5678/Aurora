package cmd

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"

	"aurora/internal/aur"
	"aurora/internal/pacman"
	"aurora/internal/resolve"

	"github.com/spf13/cobra"
)

func isValidPkgName(name string) bool {
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-\._]+$`, name)
	return matched
}

func install(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a package name to install.")
		return
	}

	packageName := args[0]

	if !isValidPkgName(packageName) {
		fmt.Println("Error: Invalid package name format.")
		return
	}

	reader := bufio.NewReader(os.Stdin)

	pacmanPkg, _ := pacman.SearchPacmanExact(packageName)
	if pacmanPkg != nil {
		fmt.Printf("\nFound \"%s\" in official repositories.\n", packageName)
		displayPacmanSummary(pacmanPkg)
		fmt.Println("\nOfficial repositories are the recommended source — packages are prebuilt and maintained as part of the normal Arch package flow.")
		fmt.Print("\nProceed with install? (y/N): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			err := pacman.InstallPacman(pacmanPkg)
			if err != nil {
				fmt.Printf("Error: installation failed: %v\n", err)
			}
		}
		return
	}

	aurPkg, _ := aur.SearchAURExact(packageName)
	if aurPkg != nil {
		fmt.Printf("\n\"%s\" not found in official repositories.\n", packageName)
		fmt.Println("Querying AUR...")
		fmt.Printf("\nFound \"%s\" in the AUR.\n", packageName)
		displayAURSummary(aurPkg)
		fmt.Print("\nProceed with install? (y/N): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			fmt.Printf("\nInstalling %s from AUR...\n", aurPkg.Name)
			err := aur.InstallAUR(aurPkg.Name)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
		return
	}

	fmt.Printf("\n\"%s\" was not found in official repositories or the AUR.\n", packageName)
	fmt.Print("Would you like to search for similar packages? (y/N): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" {
		return
	}

	pacmanResults, pacmanErr := pacman.SearchPacman(packageName)
	if pacmanErr != nil {
		fmt.Printf("Error searching official repositories: %v\n", pacmanErr)
	}

	aurResults := []resolve.AURResult{}
	hasAUR := false

	if len(pacmanResults) > 0 {
		displayPacmanResults(pacmanResults)
		fmt.Println("\nOfficial repositories are the recommended source — packages are prebuilt and maintained as part of the normal Arch package flow.")
		fmt.Print("\nSearch the AUR as well? (y/N): ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			var aurErr error
			aurResults, aurErr = aur.SearchAUR(packageName)
			if aurErr != nil {
				fmt.Printf("Error searching AUR: %v\n", aurErr)
			}
			if len(aurResults) > 0 {
				displayAURResults(aurResults)
				hasAUR = true
			} else {
				fmt.Println("\nNo results found in the AUR.")
			}
		}
	} else {
		var aurErr error
		aurResults, aurErr = aur.SearchAUR(packageName)
		if aurErr != nil {
			fmt.Printf("Error searching AUR: %v\n", aurErr)
		}
		if len(aurResults) > 0 {
			displayAURResults(aurResults)
			hasAUR = true
		} else {
			fmt.Printf("\nNo results found for \"%s\" in official repositories or the AUR.\n", packageName)
			return
		}
	}

	fmt.Print("\nEnter the package name you would like to install (or press Enter to skip): ")
	detailName, _ := reader.ReadString('\n')
	detailName = strings.TrimSpace(detailName)

	if detailName == "" {
		return
	}

	if pkg, found := lookupPacmanPackage(pacmanResults, detailName); found {
		displayPacmanSummary(&pkg)
		fmt.Print("\nProceed with install? (y/N): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))
		if input == "y" || input == "yes" {
			err := pacman.InstallPacman(&pkg)
			if err != nil {
				fmt.Printf("Error: installation failed: %v\n", err)
			}
		}
		return
	}

	if hasAUR {
		if pkg, found := lookupAURPackage(aurResults, detailName); found {
			displayAURSummary(&pkg)
			fmt.Print("\nProceed with install? (y/N): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input == "y" || input == "yes" {
				fmt.Printf("\nInstalling %s from AUR...\n", pkg.Name)
				err := aur.InstallAUR(pkg.Name)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
				}
			}
			return
		}
	}

	fmt.Printf("Could not find package \"%s\" in the search results.\n", detailName)
}

func displayPacmanSummary(pkg *resolve.PacmanResult) {
	fmt.Printf("\n  %s/%s  %s\n", pkg.Repository, pkg.Name, pkg.Version)
	if pkg.Description != "" {
		fmt.Printf("  %s\n", pkg.Description)
	}
}

func displayAURSummary(pkg *resolve.AURResult) {
	fmt.Printf("\n  %s  %s\n", pkg.Name, pkg.Version)
	if pkg.Description != "" {
		fmt.Printf("  %s\n", pkg.Description)
	}
	fmt.Printf("  Votes: %d\n", pkg.NumVotes)
	fmt.Println("\nNote: installing from AUR means Aurora will clone the build files and build locally using makepkg.")
}

var installCmd = &cobra.Command{
	Use:   "install [package]",
	Short: "Installs the specified package from official repositories or AUR",
	Run:   install,
}

func init() {
	rootCmd.AddCommand(installCmd)
}
