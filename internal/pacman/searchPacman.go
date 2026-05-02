package pacman

import (
	"fmt"
	"os/exec"
	"regexp"
	"sort"
	"strings"

	"aurora/internal/resolve"
)

func SearchPacman(packageName string) ([]resolve.PacmanResult, error) {
	cmd := exec.Command("pacman", "-Ss", packageName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error running pacman -Ss: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 || (len(lines) == 1 && lines[0] == "") {
		return nil, nil
	}

	firstLineRe := regexp.MustCompile(`^([^/]+)/([^\s]+)\s+(.+)$`)

	var results []resolve.PacmanResult
	var currentPkg *resolve.PacmanResult

	for _, line := range lines {
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			if currentPkg != nil && currentPkg.Description == "" {
				currentPkg.Description = strings.TrimSpace(line)
			}
			continue
		}

		matches := firstLineRe.FindStringSubmatch(line)
		if matches == nil {
			continue
		}

		results = append(results, resolve.PacmanResult{
			Repository: matches[1],
			Name:       matches[2],
			Version:    matches[3],
		})
		currentPkg = &results[len(results)-1]
	}

	if len(results) == 0 {
		return nil, nil
	}

	sort.SliceStable(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results, nil
}

func SearchPacmanExact(packageName string) (*resolve.PacmanResult, error) {
	cmd := exec.Command("pacman", "-Si", packageName)
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("package not found in official repositories: %s", packageName)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 0 {
		return nil, nil
	}

	fieldRe := regexp.MustCompile(`^([^:]+)\s*:\s*(.+)$`)

	var pkg resolve.PacmanResult
	for _, line := range lines {
		matches := fieldRe.FindStringSubmatch(line)
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

	if pkg.Name == "" {
		return nil, nil
	}

	return &pkg, nil
}
