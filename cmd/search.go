package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"aurora/internal/aur"
	"aurora/internal/pacman"
	"aurora/internal/resolve"

	"github.com/spf13/cobra"
)

func search(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Println("Please provide a package name to search for.")
		return
	}

	packageName := args[0]
	reader := bufio.NewReader(os.Stdin)

	pacmanResults, pacmanErr := pacman.SearchPacman(packageName)
	if pacmanErr != nil {
		fmt.Printf("Error searching official repositories: %v\n", pacmanErr)
	}

	aurResults := []resolve.AURResult{}
	hasAUR := false

	if len(pacmanResults) > 0 {
		displayPacmanResults(pacmanResults)
		fmt.Println("\nOfficial repositories are the recommended source as packages are prebuilt and maintained as part of the normal Arch package flow.")
		fmt.Print("\nWould you also like to search the AUR? (y/N): ")

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

	fmt.Print("\nWould you like to see details for a specific package? (Enter package name, or press Enter to skip): ")
	detailName, _ := reader.ReadString('\n')
	detailName = strings.TrimSpace(detailName)

	if detailName == "" {
		return
	}

	if pkg, found := lookupPacmanPackage(pacmanResults, detailName); found {
		displayPacmanDetail(pkg)
		return
	}
	if hasAUR {
		if pkg, found := lookupAURPackage(aurResults, detailName); found {
			displayAURDetail(pkg)
			return
		}
	}

	fmt.Printf("Could not find package \"%s\" in the search results.\n", detailName)
}

func lookupPacmanPackage(results []resolve.PacmanResult, name string) (resolve.PacmanResult, bool) {
	resultMap := make(map[string]resolve.PacmanResult, len(results))
	for _, r := range results {
		resultMap[r.Name] = r
	}
	pkg, exists := resultMap[name]
	return pkg, exists
}

func lookupAURPackage(results []resolve.AURResult, name string) (resolve.AURResult, bool) {
	resultMap := make(map[string]resolve.AURResult, len(results))
	for _, r := range results {
		resultMap[r.Name] = r
	}
	pkg, exists := resultMap[name]
	return pkg, exists
}

func displayPacmanResults(results []resolve.PacmanResult) {
	fmt.Printf("\nResults for \"%s\"\n", "search")
	fmt.Println("\nOfficial repositories")
	for _, r := range results {
		fmt.Printf("  %-25s %s\n", r.Repository+"/"+r.Name, r.Version)
		if r.Description != "" {
			fmt.Printf("    %s\n", r.Description)
		}
	}
}

func displayAURResults(results []resolve.AURResult) {
	fmt.Println("\nAUR")
	for i, r := range results {
		if i >= 10 {
			break
		}
		fmt.Printf("  %-25s %s\n", r.Name, r.Version)
		fmt.Printf("    %s\n", r.Description)
		fmt.Printf("    Votes: %d\n", r.NumVotes)
	}
}

func displayPacmanDetail(pkg resolve.PacmanResult) {
	fmt.Printf("\n--- Package Details ---\n")
	fmt.Printf("Name        : %s\n", pkg.Name)
	fmt.Printf("Version     : %s\n", pkg.Version)
	fmt.Printf("Repository  : %s\n", pkg.Repository)
	fmt.Printf("Description : %s\n", pkg.Description)
	fmt.Printf("Source      : official repositories\n")
	fmt.Println("\nRecommended: install via official repositories using pacman.")
}

func displayAURDetail(pkg resolve.AURResult) {
	fmt.Printf("\n--- Package Details ---\n")
	fmt.Printf("Name        : %s\n", pkg.Name)
	fmt.Printf("Version     : %s\n", pkg.Version)
	fmt.Printf("Description : %s\n", pkg.Description)
	fmt.Printf("URL         : %s\n", pkg.URL)
	fmt.Printf("Votes       : %d\n", pkg.NumVotes)
	fmt.Printf("Source      : AUR\n")
	fmt.Println("\nNote: installing from AUR means Aurora will clone the build files and build locally using makepkg.")
}

var searchCmd = &cobra.Command{
	Use:   "search [package]",
	Short: "Search for packages in official repositories and the AUR",
	Run:   search,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
