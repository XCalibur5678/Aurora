package ui

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
)

// SelectionItem represents a selectable package, topic, or candidate item in an interactive prompt.
type SelectionItem struct {
	Index       int
	Name        string
	FullName    string
	SourceLabel string
	Version     string
	Description string
}

// ConfirmAction prompts the user for a yes/no confirmation.
// If autoConfirm is true (e.g., from --yes / -y), it auto-confirms without blocking.
func ConfirmAction(reader *bufio.Reader, prompt string, autoConfirm bool) bool {
	if autoConfirm {
		fmt.Printf("%s (y/N): y [auto-confirmed]\n", prompt)
		return true
	}

	fmt.Printf("%s (y/N): ", prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// PromptSelection presents a numbered picker and parses either a numerical index or item name.
// Returns the slice index and pointer to the selected item, or (-1, nil) on skip/invalid input.
func PromptSelection(reader *bufio.Reader, prompt string, items []SelectionItem) (int, *SelectionItem) {
	if len(items) == 0 {
		return -1, nil
	}

	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	if err != nil {
		return -1, nil
	}
	input = strings.TrimSpace(input)
	if input == "" {
		return -1, nil
	}

	// 1. Try numeric choice (1-based Index matching)
	if num, err := strconv.Atoi(input); err == nil {
		for i, item := range items {
			if item.Index == num {
				return i, &items[i]
			}
		}
	}

	// 2. Try exact/case-insensitive name match or full name match
	for i, item := range items {
		if strings.EqualFold(item.Name, input) || (item.FullName != "" && strings.EqualFold(item.FullName, input)) {
			return i, &items[i]
		}
	}

	return -1, nil
}
