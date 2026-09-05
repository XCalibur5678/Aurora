package pacman

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/abhigyan-chatterjee/aurora/internal/resolve"
)

// InstallPacman installs a single official package via pacman.
func InstallPacman(pkg *resolve.PacmanResult) error {
	return InstallPacmanBatch([]string{pkg.Name})
}

// InstallPacmanBatch installs multiple official packages in a single pacman transaction.
func InstallPacmanBatch(packageNames []string) error {
	if len(packageNames) == 0 {
		return nil
	}

	fmt.Printf("\nInstalling from official repositories: %s\n", strings.Join(packageNames, ", "))
	args := append([]string{"pacman", "-S", "--noconfirm"}, packageNames...)
	fmt.Printf("Running: sudo %s\n", strings.Join(args, " "))

	installCmd := exec.Command("sudo", args...)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Stdin = os.Stdin

	return installCmd.Run()
}
