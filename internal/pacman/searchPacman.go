package pacman

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"github.com/abhigyan-chatterjee/aurora/internal/resolve"
)

var (
	pacmanSearchHeaderRe = regexp.MustCompile(`^([^/]+)/([^\s]+)\s+(.+)$`)
	pacmanExactFieldRe   = regexp.MustCompile(`^([^:]+)\s*:\s*(.+)$`)
)

// ParsePacmanSearchOutput parses the standard output of 'pacman -Ss' into structured results.
// It filters results to those where the package name contains the query (case-insensitively).
// Package descriptions appearing on subsequent indented lines are correctly associated with
// the current matching package, and are cleared if an unmatched package header is encountered.
func ParsePacmanSearchOutput(output string, packageName string) []resolve.PacmanResult {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return nil
	}

	var results []resolve.PacmanResult
	currentIdx := -1

	for _, line := range lines {
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			if currentIdx >= 0 && results[currentIdx].Description == "" {
				results[currentIdx].Description = strings.TrimSpace(line)
			}
			continue
		}

		matches := pacmanSearchHeaderRe.FindStringSubmatch(line)
		if matches == nil {
			currentIdx = -1
			continue
		}

		pkgName := matches[2]

		if !strings.Contains(strings.ToLower(pkgName), strings.ToLower(packageName)) {
			currentIdx = -1
			continue
		}

		results = append(results, resolve.PacmanResult{
			Repository: matches[1],
			Name:       pkgName,
			Version:    matches[3],
		})
		currentIdx = len(results) - 1
	}

	if len(results) == 0 {
		return nil
	}

	sort.SliceStable(results, func(i, j int) bool {
		if results[i].Name != results[j].Name {
			return results[i].Name < results[j].Name
		}
		return results[i].Repository < results[j].Repository
	})

	return results
}

func SearchPacman(packageName string) ([]resolve.PacmanResult, error) {
	cmd := exec.Command("pacman", "-Ss", packageName)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("error running pacman -Ss: %v", err)
	}

	return ParsePacmanSearchOutput(string(output), packageName), nil
}

// ParsePacmanExactOutput parses the standard output of 'pacman -Si' into a PacmanResult.
func ParsePacmanExactOutput(output string, packageName string) *resolve.PacmanResult {
	lines := strings.Split(strings.TrimSpace(output), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return nil
	}

	var pkg resolve.PacmanResult
	for _, line := range lines {
		if line == "" && pkg.Name != "" {
			break
		}
		matches := pacmanExactFieldRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}
		field := strings.TrimSpace(matches[1])
		value := strings.TrimSpace(matches[2])

		switch field {
		case "Name":
			pkg.Name = value
		case "Version":
			pkg.Version = value
		case "Description":
			pkg.Description = value
		case "Repository":
			pkg.Repository = value
		}
	}

	if pkg.Name == "" || !strings.EqualFold(pkg.Name, packageName) {
		return nil
	}

	return &pkg
}

func SearchPacmanExact(packageName string) (*resolve.PacmanResult, error) {
	cmd := exec.Command("pacman", "-Si", packageName)
	output, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil, nil
		}
		return nil, fmt.Errorf("error running pacman -Si: %v", err)
	}

	return ParsePacmanExactOutput(string(output), packageName), nil
}
