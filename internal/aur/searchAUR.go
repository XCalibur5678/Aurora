package aur

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/abhigyan-chatterjee/aurora/internal/resolve"
)

const aurUserAgent = "aurora/0.2"

func newAURRequest(url string) (*http.Request, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", aurUserAgent)
	return req, nil
}

func SearchAUR(packageName string) ([]resolve.AURResult, error) {
	const aurRPCURL = "https://aur.archlinux.org/rpc/v5/search/"
	url := aurRPCURL + packageName

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := newAURRequest(url)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

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

	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	if raw.ResultCount == 0 {
		return nil, nil
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

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].NumVotes > results[j].NumVotes
	})

	return results, nil
}

func SearchAURExact(packageName string) (*resolve.AURResult, error) {
	const aurInfoURL = "https://aur.archlinux.org/rpc/v5/info/"
	url := aurInfoURL + packageName

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := newAURRequest(url)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-OK HTTP status: %s", resp.Status)
	}

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

	err = json.NewDecoder(resp.Body).Decode(&raw)
	if err != nil {
		return nil, fmt.Errorf("error decoding JSON response: %v", err)
	}

	if raw.ResultCount == 0 {
		return nil, nil
	}

	r := raw.Results[0]
	return &resolve.AURResult{
		Name:         r.Name,
		Version:      r.Version,
		Description:  r.Description,
		URL:          r.URL,
		LastModified: r.LastModified,
		NumVotes:     r.NumVotes,
	}, nil
}
