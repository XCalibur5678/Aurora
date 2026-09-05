package ui

import (
	"regexp"
	"strings"
)

var validPkgNameRegex = regexp.MustCompile(`^[a-zA-Z0-9][a-zA-Z0-9@._+-]*$`)

// IsValidPkgName checks if a package name matches Arch Linux naming conventions
// and prevents command-line flag injection or directory traversal.
func IsValidPkgName(name string) bool {
	if name == "" || len(name) > 255 || strings.Contains(name, "..") {
		return false
	}
	return validPkgNameRegex.MatchString(name)
}
