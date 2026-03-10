package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type PackageInfo struct {
	Name    string `json:"Name"`
	Version string `json:"Version"`
}

type InfoResponse struct {
	ResultCount int           `json:"resultcount"`
	Results     []PackageInfo `json:"results"`
}

func getInstalledForeign() (map[string]string, error) {
	cmd := exec.Command("pacman", "-Qm")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	packages := make(map[string]string)
	lines := strings.SplitSeq(strings.TrimSpace(string(output)), "\n")
	for line := range lines {
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			packages[parts[0]] = parts[1]
		}
	}
	return packages, nil
}

// isNewer returns true if aurVer is newer than localVer using 'vercmp'
func isNewer(localVer, aurVer string) bool {
	cmd := exec.Command("vercmp", aurVer, localVer)
	output, err := cmd.Output()
	if err != nil {
		// If vercmp fails, fallback to basic inequality
		return aurVer != localVer
	}
	// vercmp returns 1 if arg1 > arg2
	return strings.TrimSpace(string(output)) == "1"
}

func update(cmd *cobra.Command, args []string) {
	fmt.Println("Checking for AUR updates...")

	installed, err := getInstalledForeign()
	if err != nil {
		fmt.Printf("Error: could not list foreign packages: %v\n", err)
		return
	}

	if len(installed) == 0 {
		fmt.Println("No AUR packages installed.")
		return
	}

	var argsQuery []string
	for name := range installed {
		argsQuery = append(argsQuery, "arg[]="+name)
	}

	url := "https://aur.archlinux.org/rpc/?v=5&type=info&" + strings.Join(argsQuery, "&")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error checking AUR: %v\n", err)
		return
	}
	defer resp.Body.Close()

	var response InfoResponse

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		fmt.Printf("Error decoding AUR response: %v\n", err)
		return
	}

	foundUpdates := false
	for _, pkg := range response.Results {
		localVer := installed[pkg.Name]
		if isNewer(localVer, pkg.Version) {
			fmt.Printf("Update available for %s: %s -> %s\n", pkg.Name, localVer, pkg.Version)
			foundUpdates = true
		}
	}

	if !foundUpdates {
		fmt.Println("All AUR packages are up to date.")
	} else {
		fmt.Print("Do you want to update these AUR packages? (y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if strings.ToLower(confirm) == "y" {
			for _, pkg := range response.Results {
				localVer := installed[pkg.Name]
				if isNewer(localVer, pkg.Version) {
					fmt.Printf("Updating %s...\n", pkg.Name)
					if err := cloneAndBuild(pkg.Name); err != nil {
						fmt.Printf("Error updating %s: %v\n", pkg.Name, err)
					}
				}
			}
		} else {
			fmt.Println("AUR updates skipped.")
		}
	}

	fmt.Println("Do you want to perform a full system upgrade? We discourage partial upgrades, but you can choose to do so if you want.")
	fmt.Println("Partial updates may cause the system to break. ")
	fmt.Print("If you choose to fully update the system, type y. This will invoke: sudo pacman -Syu. (y/N): ")

	var sysConfirm string
	fmt.Scanln(&sysConfirm)
	if strings.ToLower(sysConfirm) == "y" || strings.ToLower(sysConfirm) == "yes" {
		sysCmd := exec.Command("sudo", "pacman", "-Syu")
		sysCmd.Stdout = os.Stdout
		sysCmd.Stderr = os.Stderr
		sysCmd.Stdin = os.Stdin
		if err := sysCmd.Run(); err != nil {
			fmt.Printf("Error during system upgrade: %v\n", err)
		} else {
			fmt.Println("System upgrade complete. Please restart your system for the changes to take effect.")
		}
	}
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Checks for updates for installed AUR packages",
	Run:   update,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
