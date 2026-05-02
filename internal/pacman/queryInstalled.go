package pacman

import (
	"os/exec"
	"strings"
)

func SearchInstalled(query string) ([]string, error) {
	cmd := exec.Command("pacman", "-Qq")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(lines) == 1 && lines[0] == "" {
		return nil, nil
	}

	var matches []string
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), strings.ToLower(query)) {
			matches = append(matches, line)
		}
	}
	return matches, nil
}
