package pacman

import (
	"fmt"
	"os"
	"os/exec"

	"aurora/internal/resolve"
)

func InstallPacman(pkg *resolve.PacmanResult) error {
	fmt.Printf("\nInstalling %s from official repositories...\n", pkg.Name)
	fmt.Printf("\nRunning: sudo pacman -S %s\n", pkg.Name)

	installCmd := exec.Command("sudo", "pacman", "-S", pkg.Name)
	installCmd.Stdout = os.Stdout
	installCmd.Stderr = os.Stderr
	installCmd.Stdin = os.Stdin

	return installCmd.Run()
}
