package pacman

import (
	"strings"
	"testing"
)

func TestSearchPacmanExact_Found(t *testing.T) {
	pkg, err := SearchPacmanExact("bash")
	if err != nil {
		t.Fatalf("unexpected error searching for bash: %v", err)
	}
	if pkg == nil {
		t.Fatalf("expected bash package to be found, got nil")
	}
	if !strings.EqualFold(pkg.Name, "bash") {
		t.Errorf("expected package name 'bash', got %q", pkg.Name)
	}
	if pkg.Repository != "core" {
		t.Errorf("expected repository 'core', got %q", pkg.Repository)
	}
}

func TestSearchPacmanExact_NotFound(t *testing.T) {
	pkg, err := SearchPacmanExact("definitely_nonexistent_package_12345")
	if err != nil {
		t.Fatalf("expected no error for missing package, got: %v", err)
	}
	if pkg != nil {
		t.Fatalf("expected nil for missing package, got: %+v", pkg)
	}
}

func TestSearchPacman_NotFound(t *testing.T) {
	results, err := SearchPacman("definitely_nonexistent_package_12345")
	if err != nil {
		t.Fatalf("expected no error for missing package search, got: %v", err)
	}
	if len(results) != 0 {
		t.Fatalf("expected 0 results, got: %d", len(results))
	}
}

func TestSearchInstalled(t *testing.T) {
	matches, err := SearchInstalled("bash")
	if err != nil {
		t.Fatalf("unexpected error searching installed packages: %v", err)
	}
	if len(matches) == 0 {
		t.Fatalf("expected at least one installed match for 'bash', got 0")
	}

	found := false
	for _, m := range matches {
		if m == "bash" {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("expected exact match 'bash' in installed matches: %v", matches)
	}
}

func TestParsePacmanSearchOutput_UnmatchedHeaderRegression(t *testing.T) {
	// Concrete regression:
	// pacman -Ss matches packages on name OR description.
	// When an unmatched header appears after a matching header (which had no description),
	// the parser must not assign the unmatched package's description to the previous matching result.
	rawOutput := `core/match-one 1.0.0-1
extra/other-pkg 2.5.0-1
    This is the description of other package
extra/match-two 3.0.0-1
    This is the description of match two
`
	results := ParsePacmanSearchOutput(rawOutput, "match")
	if len(results) != 2 {
		t.Fatalf("expected 2 results, got %d: %+v", len(results), results)
	}

	// Results should be sorted by name: match-one, then match-two
	if results[0].Name != "match-one" {
		t.Errorf("expected first result to be match-one, got %q", results[0].Name)
	}
	if results[0].Description != "" {
		t.Errorf("expected match-one description to be empty, got %q (unmatched pkg description was wrongly assigned!)", results[0].Description)
	}

	if results[1].Name != "match-two" {
		t.Errorf("expected second result to be match-two, got %q", results[1].Name)
	}
	if results[1].Description != "This is the description of match two" {
		t.Errorf("expected match-two description to be 'This is the description of match two', got %q", results[1].Description)
	}
}

func TestParsePacmanSearchOutput_SortingAndTieBreak(t *testing.T) {
	rawOutput := `extra/zeta 1.0-1
    Zeta description
core/alpha 1.0-1
    Alpha core description
extra/alpha 1.0-1
    Alpha extra description
`
	results := ParsePacmanSearchOutput(rawOutput, "a")
	// "zeta" contains "a", "alpha" contains "a"
	if len(results) != 3 {
		t.Fatalf("expected 3 results, got %d", len(results))
	}

	// Expected order: core/alpha, extra/alpha, extra/zeta
	if results[0].Name != "alpha" || results[0].Repository != "core" {
		t.Errorf("expected results[0] to be core/alpha, got %s/%s", results[0].Repository, results[0].Name)
	}
	if results[1].Name != "alpha" || results[1].Repository != "extra" {
		t.Errorf("expected results[1] to be extra/alpha, got %s/%s", results[1].Repository, results[1].Name)
	}
	if results[2].Name != "zeta" || results[2].Repository != "extra" {
		t.Errorf("expected results[2] to be extra/zeta, got %s/%s", results[2].Repository, results[2].Name)
	}
}

func TestParsePacmanSearchOutput_EmptyAndMalformed(t *testing.T) {
	if res := ParsePacmanSearchOutput("", "bash"); res != nil {
		t.Errorf("expected nil for empty output, got %+v", res)
	}
	if res := ParsePacmanSearchOutput("   \n\n  \t\n", "bash"); res != nil {
		t.Errorf("expected nil for whitespace output, got %+v", res)
	}
	if res := ParsePacmanSearchOutput("not-a-valid-line\n   indented orphan", "bash"); res != nil {
		t.Errorf("expected nil for invalid output, got %+v", res)
	}
}

func TestParsePacmanExactOutput(t *testing.T) {
	raw := `Repository      : core
Name            : bash
Version         : 5.2.026-2
Description     : The GNU Bourne Again shell
Architecture    : x86_64
`
	pkg := ParsePacmanExactOutput(raw, "bash")
	if pkg == nil {
		t.Fatalf("expected non-nil package")
	}
	if pkg.Name != "bash" {
		t.Errorf("expected name 'bash', got %q", pkg.Name)
	}
	if pkg.Repository != "core" {
		t.Errorf("expected repository 'core', got %q", pkg.Repository)
	}
	if pkg.Version != "5.2.026-2" {
		t.Errorf("expected version '5.2.026-2', got %q", pkg.Version)
	}
	if pkg.Description != "The GNU Bourne Again shell" {
		t.Errorf("expected description 'The GNU Bourne Again shell', got %q", pkg.Description)
	}

	// Case-insensitive query
	pkgCase := ParsePacmanExactOutput(raw, "BASH")
	if pkgCase == nil || pkgCase.Name != "bash" {
		t.Errorf("expected case-insensitive match for BASH")
	}

	// Non-matching query
	pkgOther := ParsePacmanExactOutput(raw, "zsh")
	if pkgOther != nil {
		t.Errorf("expected nil for non-matching package query, got %+v", pkgOther)
	}
}
