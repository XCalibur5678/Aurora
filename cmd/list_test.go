package cmd

import (
	"strings"
	"testing"
)

func TestList_Help(t *testing.T) {
	out := captureStdout(func() {
		_ = listCmd.Help()
	})
	if !strings.Contains(out, "apt-style") && !strings.Contains(out, "upgradable") {
		t.Errorf("expected apt-style or upgradable in list help, got: %s", out)
	}
}

func TestList_InstalledQuery(t *testing.T) {
	listOfficialFlag = false
	listAURFlag = false
	listOrphansFlag = false
	listUpgradableFlag = false

	out := captureStdout(func() {
		listPackages(listCmd, []string{"bash"})
	})

	if !strings.Contains(out, "bash") {
		t.Errorf("expected 'bash' to be found in installed packages, got: %s", out)
	}
}

func TestList_Nonexistent(t *testing.T) {
	listOfficialFlag = false
	listAURFlag = false
	listOrphansFlag = false
	listUpgradableFlag = false

	out := captureStdout(func() {
		listPackages(listCmd, []string{"definitely_nonexistent_pkg_12345"})
	})

	if !strings.Contains(out, "No matching installed packages found") {
		t.Errorf("expected no matching packages message, got: %s", out)
	}
}
