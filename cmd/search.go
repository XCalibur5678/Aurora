package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/spf13/cobra"
)

// Define a struct to hold the JSON response from the AUR RPC interface
type SearchResult struct {
	Type        string `json:"type"`
	ResultCount int    `json:"resultcount"`
	Results     []struct {
		Name         string `json:"Name"`
		Version      string `json:"Version"`
		Description  string `json:"Description"`
		URL          string `json:"URL"`
		LastModified int64  `json:"LastModified"`
		NumVotes     int    `json:"NumVotes"`
	} `json:"results"`
}

func searchPackage(packageName string) (string, error) {
	fmt.Printf("Searching for package: %s\n", packageName)
	const aurRPCURL = "https://aur.archlinux.org/rpc/v5/search/"
	url := aurRPCURL + packageName

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)

	if err != nil {
		fmt.Printf("Error making HTTP request: %v\n", err)
		return err.Error(), err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}
	var result SearchResult

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", fmt.Errorf("Error decoding JSON response: %v\n", err)
	}

	//If no packages are found, print a message and exit
	if result.ResultCount == 0 {
		fmt.Println("No packages found matching the search query.")
		return "", nil
	}

	/*Print number of packages found and ask the user how many packages to display.
	If the user does not enter a number , length defaults to 10.
	Sort the packages by number of votes in descending order before displaying.*/
	fmt.Printf("Found %d package(s):\n", result.ResultCount)
	fmt.Println("Displaying the Top 10 packages based on number of votes.")

	sort.SliceStable(result.Results, func(i, j int) bool {
		return result.Results[i].NumVotes > result.Results[j].NumVotes
	})

	top_result_length := len(result.Results)
	if top_result_length < 10 {
		top_result_length = top_result_length
	} else {
		top_result_length = 10
	}

	for _, pkg := range result.Results[:10] {
		fmt.Printf("- Name: %s Votes: %d\n", pkg.Name, pkg.NumVotes)
	}

	fmt.Print("Do you want to display more packages? (Enter a number or press Enter to skip): ")

	fmt.Scanln(&displayCount)
	if displayCount != 10 {
		if displayCount > 0 && displayCount <= result.ResultCount {
			for _, pkg := range result.Results[:displayCount] {
				fmt.Printf("- Name: %s Votes: %d\n", pkg.Name, pkg.NumVotes)
			}
		} else {
			fmt.Println("Invalid input.")
		}
	}
	fmt.Print("Enter your package name to install. If you want to quit, press Enter: ")
	name := ""
	fmt.Scanln(&name)
	if name == "" {
		return "", nil
	}

	for _, pkg := range result.Results[:displayCount] {
		if pkg.Name == name {
			fmt.Println("----------Start----------")
			fmt.Printf("Name: %s\nVersion: %s\nDescription: %s\nURL: %s\nLast Modified: %s\nVotes: %d\n",
				pkg.Name, pkg.Version, pkg.Description, pkg.URL, time.Unix(pkg.LastModified, 0).Format(time.RFC1123), pkg.NumVotes)
			fmt.Println("---------End-----------")
			return name, nil
		}
	}
	fmt.Printf("Package '%s' not found in the search results.\n", name)
	return name, nil

}

var search = func(cmd *cobra.Command, args []string) {
	//if no argument is provided, print a message and exit
	if len(args) == 0 {
		fmt.Println("Please provide a package name to search for.")
		return
	}
	packageName := args[0]
	_, err := searchPackage(packageName)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

var searchCmd = &cobra.Command{
	Use:   "search [package]",
	Short: "searches for the specified package in AUR via the AUR RPC interface",
	Run:   search,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
