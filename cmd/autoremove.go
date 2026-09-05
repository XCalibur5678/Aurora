package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/abhigyan-chatterjee/aurora/internal/ui"
	"github.com/spf13/cobra"
)

// GetOrphanPackages returns package names identified as unneeded orphan dependencies by pacman.
func GetOrphanPackages() ([]string, error) {
	cmd := exec.Command("pacman", "-Qtdq")
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var orphans []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			orphans = append(orphans, trimmed)
		}
	}
	return orphans, nil
}

func autoremove(cmd *cobra.Command, args []string) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Checking for unused orphan dependencies...")
	orphans, err := GetOrphanPackages()
	if err != nil {
		fmt.Printf("Error checking for orphan packages: %v\n", err)
		return
	}

	if len(orphans) == 0 {
		fmt.Println("No orphan packages found. Your system is clean!")
		return
	}

	fmt.Printf("\nFound %d orphan package(s) no longer needed by any installed software:\n", len(orphans))
	for _, p := range orphans {
		fmt.Printf("  - %s\n", p)
	}

	if !conciseFlag {
		fmt.Println("\nAurora will run: sudo pacman -Rns --noconfirm <packages>")
		fmt.Println("Meaning: safely remove unneeded dependencies and their unused configuration files.")
	}

	if !ui.ConfirmAction(reader, "\nProceed with autoremove?", yesFlag) {
		fmt.Println("Autoremove cancelled.")
		return
	}

	removeArgs := append([]string{"pacman", "-Rns", "--noconfirm"}, orphans...)
	execCmd := exec.Command("sudo", removeArgs...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Stdin = os.Stdin

	if err := execCmd.Run(); err != nil {
		fmt.Printf("Error during autoremove: %v\n", err)
		return
	}

	fmt.Println("\nAutoremove completed successfully. System cleaned.")
}

var autoremoveCmd = &cobra.Command{
	Use:     "autoremove",
	Aliases: []string{"orphans"},
	Short:   "Removes unused orphan packages and leftover dependencies",
	Long: `Autoremove scans the system for orphan dependencies (packages installed as dependencies
for software that has since been uninstalled) and safely removes them with 'sudo pacman -Rns'.

This command mimics 'apt autoremove' to keep your Arch system clean and clutter-free.`,
	Run: autoremove,
}

func init() {
	rootCmd.AddCommand(autoremoveCmd)
}
