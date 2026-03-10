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

func search(cmd *cobra.Command, args []string) {

	//if no argument is provided, print a message and exit
	if len(args) == 0 {
		fmt.Println("Please provide a package name to search for.")
		return
	}

	packageName := args[0]
	fmt.Printf("Searching for package: %s\n", packageName)
	const aurRPCURL = "https://aur.archlinux.org/rpc/v5/search/"
	url := aurRPCURL + packageName

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Printf("Error making HTTP request: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Received non-OK HTTP status: %s\n", resp.Status)
		return
	}
	var result SearchResult

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		fmt.Printf("Error decoding JSON response: %v\n", err)
		return
	}

	//If no packages are found, print a message and exit
	if result.ResultCount == 0 {
		fmt.Println("No packages found matching the search query.")
		return
	}

	/*Print number of packages found and ask the user how many packages to display.
	If the user does not enter a number , length defaults to 10.
	Sort the packages by number of votes in descending order before displaying.*/
	fmt.Printf("Found %d package(s):\n", result.ResultCount)
	fmt.Println("Displaying the Top 10 packages based on number of votes.")

	sort.SliceStable(result.Results, func(i, j int) bool {
		return result.Results[i].NumVotes > result.Results[j].NumVotes
	})

	for _, pkg := range result.Results[:10] {
		fmt.Printf("- Name: %s Votes: %d\n", pkg.Name, pkg.NumVotes)
	}

	fmt.Print("Do you want to display more packages? (Enter a number or press Enter to skip): ")
	var displayCount int = 10
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
	fmt.Print("For more details, enter the package name. (Press Enter to skip): ")
	name := ""
	fmt.Scanln(&name)
	if name == "" {
		return
	}

	for _, pkg := range result.Results[:displayCount] {
		if pkg.Name == name {
			fmt.Println("----------Start----------")
			fmt.Printf("Name: %s\nVersion: %s\nDescription: %s\nURL: %s\nLast Modified: %s\nVotes: %d\n",
				pkg.Name, pkg.Version, pkg.Description, pkg.URL, time.Unix(pkg.LastModified, 0).Format(time.RFC1123), pkg.NumVotes)
			fmt.Println("---------End-----------")
			return
		}
	}
	fmt.Printf("Package '%s' not found in the search results.\n", name)
}

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "searches for the specified package in AUR via the AUR RPC interface",
	Long: `Searches for the specified package in AUR via the AUR RPC interface.
	Example usage:
  		aurora search <package-name>
  	This command will return a list of packages that match the search query, along with their details such as package name, version, description, and votes.
  	Currently, this command only supports a single package name as an argument, but in the future, it may be extended to support multiple package names or additional search parameters.`,
	Run: search,
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
