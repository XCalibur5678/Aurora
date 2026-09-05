package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"

	"github.com/abhigyan-chatterjee/aurora/internal/aur"
	"github.com/abhigyan-chatterjee/aurora/internal/pacman"
	"github.com/abhigyan-chatterjee/aurora/internal/ui"

	"github.com/spf13/cobra"
)

type outdatedAURPkg struct {
	name     string
	localVer string
	aurVer   string
}

func update(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Checking for updates...")

	var outdated []outdatedAURPkg

	installed, err := pacman.GetForeignPackages()
	if err != nil {
		fmt.Printf("Warning: error checking foreign packages: %v\n", err)
	} else if len(installed) > 0 {
		var pkgNames []string
		for name := range installed {
			pkgNames = append(pkgNames, name)
		}
		sort.Strings(pkgNames)

		aurResults, err := aur.GetAURInfoBatch(pkgNames)
		if err != nil {
			fmt.Printf("Warning: error checking AUR updates: %v\n", err)
		} else {
			for _, pkg := range aurResults {
				localVer := installed[pkg.Name]
				if pacman.IsNewerVersion(localVer, pkg.Version) {
					outdated = append(outdated, outdatedAURPkg{
						name:     pkg.Name,
						localVer: localVer,
						aurVer:   pkg.Version,
					})
				}
			}
			sort.Slice(outdated, func(i, j int) bool {
				return outdated[i].name < outdated[j].name
			})
		}
	}

	fmt.Println("\nAurora update plan:")
	if len(outdated) > 0 {
		fmt.Printf("  AUR packages to upgrade (%d):\n", len(outdated))
		for _, o := range outdated {
			fmt.Printf("    - %s (%s -> %s)\n", o.name, o.localVer, o.aurVer)
		}
	} else {
		fmt.Println("  AUR packages: all up to date")
	}

	fmt.Println("  System upgrade: refresh package databases and upgrade official packages")

	if !ui.ConfirmAction(reader, "\nProceed with update?", yesFlag) {
		fmt.Println("Update cancelled.")
		return
	}

	if len(outdated) > 0 {
		fmt.Println("\n--- Upgrading AUR Packages ---")
		for _, o := range outdated {
			fmt.Printf("\nUpdating %s...\n", o.name)
			err := aur.InstallAUR(o.name)
			if err != nil {
				fmt.Printf("Error updating %s: %v\n", o.name, err)
			}
		}
	}

	fmt.Println("\n--- Upgrading System Packages ---")
	err = pacman.SystemUpdate()
	if err != nil {
		fmt.Printf("Error during system upgrade: %v\n", err)
	} else {
		fmt.Println("System upgrade complete.")
	}
}

var updateCmd = &cobra.Command{
	Use:     "update",
	Aliases: []string{"upgrade"},
	Short:   "Checks for AUR updates and performs a full system upgrade",
	Long: `Checks for AUR updates and performs a full system upgrade using 'sudo pacman -Syu --noconfirm'.

Underlying pacman command flags:
  -S          : synchronize/install packages from repositories
  -y          : refresh package databases
  -u          : upgrade installed packages
  --noconfirm : skip pacman's per-package confirmation prompts`,
	Run: update,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
