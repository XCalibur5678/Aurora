package aur

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"aurora/internal/resolve"
)

func GetAURInfoBatch(packageNames []string) ([]resolve.AURResult, error) {
	if len(packageNames) == 0 {
		return nil, nil
	}

	var argsQuery []string
	for _, name := range packageNames {
		argsQuery = append(argsQuery, "arg[]="+name)
	}

	url := "https://aur.archlinux.org/rpc/?v=5&type=info&" + strings.Join(argsQuery, "&")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error checking AUR: %v", err)
	}
	defer resp.Body.Close()

	var raw struct {
		ResultCount int `json:"resultcount"`
		Results     []struct {
			Name         string `json:"Name"`
			Version      string `json:"Version"`
			Description  string `json:"Description"`
			URL          string `json:"URL"`
			LastModified int64  `json:"LastModified"`
			NumVotes     int    `json:"NumVotes"`
		} `json:"results"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, fmt.Errorf("error decoding AUR response: %v", err)
	}

	var results []resolve.AURResult
	for _, r := range raw.Results {
		results = append(results, resolve.AURResult{
			Name:         r.Name,
			Version:      r.Version,
			Description:  r.Description,
			URL:          r.URL,
			LastModified: r.LastModified,
			NumVotes:     r.NumVotes,
		})
	}

	return results, nil
}
