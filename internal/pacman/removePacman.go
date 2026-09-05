package pacman

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// RemovePacman removes specified packages in a single transaction via pacman.
func RemovePacman(packageNames ...string) error {
	if len(packageNames) == 0 {
		return nil
	}
	args := append([]string{"pacman", "-Rns", "--noconfirm"}, packageNames...)
	fmt.Printf("\nAurora will run: sudo %s\n", strings.Join(args, " "))
	fmt.Println("Meaning: remove specified packages, their unused dependencies, and related system configuration files.")

	removeCmd := exec.Command("sudo", args...)
	removeCmd.Stdout = os.Stdout
	removeCmd.Stderr = os.Stderr
	removeCmd.Stdin = os.Stdin

	return removeCmd.Run()
}
