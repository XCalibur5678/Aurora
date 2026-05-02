package aur

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
)

func InstallAUR(packageName string) error {
	usr, err := user.Current()
	if err != nil {
		return fmt.Errorf("could not determine current user: %v", err)
	}

	cacheDir := filepath.Join(usr.HomeDir, ".cache", "aurora")

	err = os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return fmt.Errorf("could not create cache directory: %v", err)
	}

	pkgDir := filepath.Join(cacheDir, packageName)

	os.RemoveAll(pkgDir)

	repoURL := fmt.Sprintf("https://aur.archlinux.org/%s.git", packageName)
	fmt.Printf("Cloning %s...\n", repoURL)
	cloneCmd := exec.Command("git", "clone", repoURL, pkgDir)
	cloneCmd.Stdout = os.Stdout
	cloneCmd.Stderr = os.Stderr
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("failed to clone: %v", err)
	}

	fmt.Println("Building package...")
	buildCmd := exec.Command("makepkg", "-si")
	buildCmd.Dir = pkgDir
	buildCmd.Stdout = os.Stdout
	buildCmd.Stderr = os.Stderr
	buildCmd.Stdin = os.Stdin

	if err := buildCmd.Run(); err != nil {
		return fmt.Errorf("makepkg failed: %v", err)
	}

	return nil
}
