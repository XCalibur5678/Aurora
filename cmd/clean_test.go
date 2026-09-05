package cmd

import (
	"strings"
	"testing"
)

func TestClean_Help(t *testing.T) {
	out := captureStdout(func() {
		_ = cleanCmd.Help()
	})
	if !strings.Contains(out, "clean") {
		t.Errorf("expected clean in help, got: %s", out)
	}
}

func TestClean_DryRun(t *testing.T) {
	cleanDryRunFlag = true
	cleanCacheFlag = true
	cleanPacmanFlag = false
	cleanAllFlag = false

	out := captureStdout(func() {
		clean(cleanCmd, []string{})
	})

	if !strings.Contains(out, "Dry Run") {
		t.Errorf("expected Dry Run in output, got: %s", out)
	}
	if !strings.Contains(out, "AUR Build Cache") {
		t.Errorf("expected AUR Build Cache in dry run output, got: %s", out)
	}
}

func TestHumanSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{500, "500 B"},
		{1024, "1.00 KB"},
		{1048576, "1.00 MB"},
		{1073741824, "1.00 GB"},
	}

	for _, tc := range tests {
		got := humanSize(tc.bytes)
		if got != tc.expected {
			t.Errorf("humanSize(%d) = %q, expected %q", tc.bytes, got, tc.expected)
		}
	}
}
