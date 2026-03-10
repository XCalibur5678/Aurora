package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"regexp"

	"github.com/spf13/cobra"
)

func isValidPkgName(name string) bool {
	// Basic AUR package naming convention: alphanumeric, dash, dot, underscore
	matched, _ := regexp.MatchString(`^[a-zA-Z0-9\-\._]+$`, name)
	return matched
}

func install(cmd *cobra.Command, args []string) {
	//if no argument is provided, print a message and exit
	if len(args) == 0 {
		fmt.Println("Please provide a package name to search for.")
		return
	}

	packageName := args[0]

	if !isValidPkgName(packageName) {
		fmt.Println("Error: Invalid package name format.")
		return
	}

	packageName, _ = searchPackage(packageName)

	if packageName == "" {
		return
	}

	err := cloneAndBuild(packageName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

func cloneAndBuild(packageName string) error {
	usr, _ := user.Current()
	cacheDir := filepath.Join(usr.HomeDir, ".cache", "aurora")

	// Create cache directory
	err := os.MkdirAll(cacheDir, 0755)
	if err != nil {
		return fmt.Errorf("could not create cache directory: %v", err)
	}

	pkgDir := filepath.Join(cacheDir, packageName)

	// Clean up old directory if it exists
	os.RemoveAll(pkgDir)

	// Clone the repo
	repoURL := fmt.Sprintf("https://aur.archlinux.org/%s.git", packageName)
	fmt.Printf("Cloning %s...\n", repoURL)
	cloneCmd := exec.Command("git", "clone", repoURL, pkgDir)
	if err := cloneCmd.Run(); err != nil {
		return fmt.Errorf("failed to clone: %v", err)
	}

	// Build and install
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

var installCmd = &cobra.Command{
	Use:   "install [package]",
	Short: "installs the specified package on the system",
	Run:   install,
}

func init() {
	rootCmd.AddCommand(installCmd)
}
