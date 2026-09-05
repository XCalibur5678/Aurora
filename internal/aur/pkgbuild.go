package aur

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// PKGBUILDInfo holds extracted high-level metadata from a PKGBUILD.
type PKGBUILDInfo struct {
	PackageName string
	Version     string
	Release     string
	Description string
	URL         string
	Maintainer  string
	Depends     []string
	MakeDepends []string
	Sources     []string
	RawContent  string
}

// FetchPKGBUILD retrieves the raw PKGBUILD file from the AUR cgit interface.
func FetchPKGBUILD(packageName string) (string, error) {
	if !isValidPackageName(packageName) {
		return "", fmt.Errorf("invalid package name: %s", packageName)
	}

	cgitURL := fmt.Sprintf("https://aur.archlinux.org/cgit/aur.git/plain/PKGBUILD?h=%s", url.QueryEscape(packageName))

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := newAURRequest(cgitURL)
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error fetching PKGBUILD: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("PKGBUILD not found for package %q", packageName)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received HTTP status %s when fetching PKGBUILD", resp.Status)
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading PKGBUILD response: %v", err)
	}

	content := string(bodyBytes)
	if strings.TrimSpace(content) == "" {
		return "", fmt.Errorf("PKGBUILD for %q is empty", packageName)
	}

	return content, nil
}

// ParsePKGBUILD extracts high-level metadata from raw PKGBUILD text.
func ParsePKGBUILD(content string) PKGBUILDInfo {
	info := PKGBUILDInfo{RawContent: content}

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "# Maintainer:") {
			info.Maintainer = strings.TrimSpace(strings.TrimPrefix(trimmed, "# Maintainer:"))
		}
	}

	info.PackageName = extractSimpleVar(content, "pkgname")
	info.Version = extractSimpleVar(content, "pkgver")
	info.Release = extractSimpleVar(content, "pkgrel")
	info.Description = extractSimpleVar(content, "pkgdesc")
	info.URL = extractSimpleVar(content, "url")
	info.Depends = extractArrayVar(content, "depends")
	info.MakeDepends = extractArrayVar(content, "makedepends")
	info.Sources = extractArrayVar(content, "source")

	return info
}

func extractSimpleVar(content, varName string) string {
	re := regexp.MustCompile(`(?m)^\s*` + regexp.QuoteMeta(varName) + `=(?:["']([^"']+)["']|([^#\s]+))`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		if matches[1] != "" {
			return matches[1]
		}
		if len(matches) > 2 {
			return matches[2]
		}
	}
	return ""
}

func extractArrayVar(content, varName string) []string {
	re := regexp.MustCompile(`(?s)` + regexp.QuoteMeta(varName) + `=\((.*?)\)`)
	matches := re.FindStringSubmatch(content)
	if len(matches) < 2 {
		return nil
	}

	rawInside := matches[1]
	itemRe := regexp.MustCompile(`['"]([^'"]+)['"]|(\S+)`)
	itemMatches := itemRe.FindAllStringSubmatch(rawInside, -1)

	var results []string
	for _, m := range itemMatches {
		val := m[1]
		if val == "" {
			val = m[2]
		}
		val = strings.TrimSpace(val)
		if val != "" && !strings.HasPrefix(val, "#") {
			results = append(results, val)
		}
	}
	return results
}
