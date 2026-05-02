package pacman

import (
	"fmt"
	"os"
	"os/exec"
)

func RemovePacman(packageName string) error {
	fmt.Printf("\nAurora will run: sudo pacman -Rns %s\n", packageName)
	fmt.Println("Meaning: remove the package, its unused dependencies, and related system configuration files.")

	removeCmd := exec.Command("sudo", "pacman", "-Rns", packageName)
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr
	removeCmd.Stdin = os.Stdin

	return removeCmd.Run()
}
