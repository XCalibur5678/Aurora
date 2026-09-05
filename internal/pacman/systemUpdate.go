package pacman

import (
	"fmt"
	"os"
	"os/exec"
)

func SystemUpdate() error {
	fmt.Println("\nRunning system upgrade: sudo pacman -Syu --noconfirm")

	sysCmd := exec.Command("sudo", "pacman", "-Syu", "--noconfirm")
	sysCmd.Stdout = os.Stdout
	sysCmd.Stderr = os.Stderr
	sysCmd.Stdin = os.Stdin

	return sysCmd.Run()
}
