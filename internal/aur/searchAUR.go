package AUR

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"
)

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

func searchAUR(packageName string) (*SearchResult, int, error) {
	const aurRPCURL = "https://aur.archlinux.org/rpc/v5/search/"
	url := aurRPCURL + packageName

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)

	if err != nil {
		fmt.Printf("Error making HTTP request: %v\n", err)
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, 0, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}
	var result SearchResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, 0, fmt.Errorf("Error decoding JSON response: %v\n", err)
	}

	if result.ResultCount == 0 {
		return nil, 0, nil
	}

	sort.SliceStable(result.Results, func(i, j int) bool {
		return result.Results[i].NumVotes > result.Results[j].NumVotes
	})
	return &result, len(result.Results), nil

}
