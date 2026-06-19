package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/abhigyan-chatterjee/aurora/internal/aur"
	"github.com/abhigyan-chatterjee/aurora/internal/pacman"

	"github.com/spf13/cobra"
)

func update(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Checking for AUR updates...")

	installed, err := pacman.GetForeignPackages()
	if err != nil {
		fmt.Printf("Error checking foreign packages: %v\n", err)
		return
	}

	if len(installed) == 0 {
		fmt.Println("No AUR packages installed.")
	} else {
		var pkgNames []string
		for name := range installed {
			pkgNames = append(pkgNames, name)
		}

		aurResults, err := aur.GetAURInfoBatch(pkgNames)
		if err != nil {
			fmt.Printf("Error checking AUR: %v\n", err)
			return
		}

		var outdated []struct {
			name     string
			localVer string
			aurVer   string
		}

		for _, pkg := range aurResults {
			localVer := installed[pkg.Name]
			if pacman.IsNewerVersion(localVer, pkg.Version) {
				outdated = append(outdated, struct {
					name     string
					localVer string
					aurVer   string
				}{
					name:     pkg.Name,
					localVer: localVer,
					aurVer:   pkg.Version,
				})
			}
		}

		if len(outdated) == 0 {
			fmt.Println("All AUR packages are up to date.")
		} else {
			fmt.Println("\nThe following AUR packages have updates available:")
			for _, o := range outdated {
				fmt.Printf("  %s  %s -> %s\n", o.name, o.localVer, o.aurVer)
			}

			fmt.Print("\nUpdate these AUR packages? (y/N): ")
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(strings.ToLower(input))
			if input == "y" || input == "yes" {
				for _, o := range outdated {
					fmt.Printf("\nUpdating %s...\n", o.name)
					err := aur.InstallAUR(o.name)
					if err != nil {
						fmt.Printf("Error updating %s: %v\n", o.name, err)
					}
				}
			} else {
				fmt.Println("AUR updates skipped.")
			}
		}
	}

	fmt.Println("\nAurora is preparing a full system update.")
	fmt.Println("\nThis will:")
	fmt.Println("  1. Refresh package databases")
	fmt.Println("  2. Upgrade installed official packages")

	fmt.Println("\nUnderlying command: sudo pacman -Syu")
	fmt.Println("\nMeaning:")
	fmt.Println("  -S : synchronize/install packages from repositories")
	fmt.Println("  -y : refresh package databases")
	fmt.Println("  -u : upgrade installed packages")

	fmt.Println("\nNote: partial upgrades are strongly discouraged on Arch-based systems.")
	fmt.Print("\nProceed with full system upgrade? (y/N): ")

	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	if input != "y" && input != "yes" {
		fmt.Println("System upgrade skipped.")
		return
	}

	err = pacman.SystemUpdate()
	if err != nil {
		fmt.Printf("Error during system upgrade: %v\n", err)
	} else {
		fmt.Println("System upgrade complete.")
	}
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Checks for AUR updates and performs a full system upgrade",
	Run:   update,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
