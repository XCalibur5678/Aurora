package ui

import (
	"bufio"
	"strings"
	"testing"
)

func TestIsValidPkgName(t *testing.T) {
	valid := []string{
		"neovim",
		"git",
		"gtk+3",
		"libsigc++",
		"gcc@11",
		"clang@14",
		"python-pip",
		"lib32-glibc",
		"a.b-c_d+e@f",
	}

	for _, name := range valid {
		if !IsValidPkgName(name) {
			t.Errorf("expected %q to be a valid package name, but got false", name)
		}
	}

	invalid := []string{
		"",
		".",
		"..",
		"-Syu",
		"-foo",
		".hidden",
		"foo..bar",
		"../escaped",
		"foo bar",
		"foo;rm",
		"foo/bar",
		"foo*bar",
		"foo$bar",
		strings.Repeat("a", 256),
	}

	for _, name := range invalid {
		if IsValidPkgName(name) {
			t.Errorf("expected %q to be an invalid package name, but got true", name)
		}
	}
}

func TestPromptSelection(t *testing.T) {
	items := []SelectionItem{
		{Index: 1, Name: "neovim", FullName: "extra/neovim", SourceLabel: "official"},
		{Index: 2, Name: "neovim-git", SourceLabel: "AUR"},
		{Index: 3, Name: "gtk+3", SourceLabel: "official"},
	}

	// Test selection by number
	r := bufio.NewReader(strings.NewReader("2\n"))
	idx, sel := PromptSelection(r, "prompt", items)
	if idx != 1 || sel == nil || sel.Name != "neovim-git" {
		t.Fatalf("expected item index 1 (neovim-git), got idx=%d, sel=%+v", idx, sel)
	}

	// Test selection by exact name
	r = bufio.NewReader(strings.NewReader("neovim\n"))
	idx, sel = PromptSelection(r, "prompt", items)
	if idx != 0 || sel == nil || sel.Name != "neovim" {
		t.Fatalf("expected item index 0 (neovim), got idx=%d, sel=%+v", idx, sel)
	}

	// Test selection by FullName (e.g. repo/name)
	r = bufio.NewReader(strings.NewReader("extra/neovim\n"))
	idx, sel = PromptSelection(r, "prompt", items)
	if idx != 0 || sel == nil || sel.Name != "neovim" {
		t.Fatalf("expected item index 0 (extra/neovim), got idx=%d, sel=%+v", idx, sel)
	}

	// Test selection by case-insensitive name
	r = bufio.NewReader(strings.NewReader("NEOVIM-GIT\n"))
	idx, sel = PromptSelection(r, "prompt", items)
	if idx != 1 || sel == nil || sel.Name != "neovim-git" {
		t.Fatalf("expected item index 1 (neovim-git), got idx=%d, sel=%+v", idx, sel)
	}

	// Test empty / skip
	r = bufio.NewReader(strings.NewReader("\n"))
	idx, sel = PromptSelection(r, "prompt", items)
	if idx != -1 || sel != nil {
		t.Fatalf("expected (-1, nil) on empty input, got idx=%d, sel=%+v", idx, sel)
	}

	// Test invalid input
	r = bufio.NewReader(strings.NewReader("nonexistent\n"))
	idx, sel = PromptSelection(r, "prompt", items)
	if idx != -1 || sel != nil {
		t.Fatalf("expected (-1, nil) on invalid input, got idx=%d, sel=%+v", idx, sel)
	}

	// Test empty list
	r = bufio.NewReader(strings.NewReader("1\n"))
	idx, sel = PromptSelection(r, "prompt", nil)
	if idx != -1 || sel != nil {
		t.Fatalf("expected (-1, nil) on empty items, got idx=%d, sel=%+v", idx, sel)
	}
}

func TestConfirmAction(t *testing.T) {
	// Auto-confirm
	r := bufio.NewReader(strings.NewReader(""))
	if !ConfirmAction(r, "Proceed?", true) {
		t.Errorf("expected auto-confirmation when autoConfirm is true")
	}

	// Input = "y"
	r = bufio.NewReader(strings.NewReader("y\n"))
	if !ConfirmAction(r, "Proceed?", false) {
		t.Errorf("expected confirmation on 'y'")
	}

	// Input = "yes"
	r = bufio.NewReader(strings.NewReader("yes\n"))
	if !ConfirmAction(r, "Proceed?", false) {
		t.Errorf("expected confirmation on 'yes'")
	}

	// Input = "n"
	r = bufio.NewReader(strings.NewReader("n\n"))
	if ConfirmAction(r, "Proceed?", false) {
		t.Errorf("expected rejection on 'n'")
	}

	// Input = empty
	r = bufio.NewReader(strings.NewReader("\n"))
	if ConfirmAction(r, "Proceed?", false) {
		t.Errorf("expected rejection on empty input")
	}
}
