package aur

import (
	"strings"
	"testing"
)

func TestParsePKGBUILD(t *testing.T) {
	raw := `# Maintainer: Jane Doe <jane@example.com>
pkgname=sample-pkg
pkgver=1.2.3
pkgrel=1
pkgdesc="A sample package for testing"
arch=('x86_64')
url="https://example.com/sample"
license=('MIT')
depends=('glibc' 'openssl>=1.1')
makedepends=(
  'go>=1.20'
  'git'
)
source=("sample-1.2.3.tar.gz::https://example.com/download.tar.gz")
sha256sums=('abc1234567890abcdef')

build() {
  echo "building"
}
`

	info := ParsePKGBUILD(raw)

	if info.PackageName != "sample-pkg" {
		t.Errorf("expected pkgname 'sample-pkg', got %q", info.PackageName)
	}
	if info.Version != "1.2.3" {
		t.Errorf("expected pkgver '1.2.3', got %q", info.Version)
	}
	if info.Release != "1" {
		t.Errorf("expected pkgrel '1', got %q", info.Release)
	}
	if info.Description != "A sample package for testing" {
		t.Errorf("expected description 'A sample package for testing', got %q", info.Description)
	}
	if info.URL != "https://example.com/sample" {
		t.Errorf("expected URL 'https://example.com/sample', got %q", info.URL)
	}
	if !strings.Contains(info.Maintainer, "Jane Doe") {
		t.Errorf("expected Maintainer containing 'Jane Doe', got %q", info.Maintainer)
	}

	if len(info.Depends) != 2 || info.Depends[0] != "glibc" || info.Depends[1] != "openssl>=1.1" {
		t.Errorf("expected depends ['glibc', 'openssl>=1.1'], got %v", info.Depends)
	}

	if len(info.MakeDepends) != 2 || info.MakeDepends[0] != "go>=1.20" || info.MakeDepends[1] != "git" {
		t.Errorf("expected makedepends ['go>=1.20', 'git'], got %v", info.MakeDepends)
	}

	if len(info.Sources) != 1 || !strings.Contains(info.Sources[0], "sample-1.2.3.tar.gz") {
		t.Errorf("expected source to contain tar.gz url, got %v", info.Sources)
	}
}

func TestFetchPKGBUILD_InvalidName(t *testing.T) {
	_, err := FetchPKGBUILD("../invalid")
	if err == nil {
		t.Errorf("expected error for invalid package name, got nil")
	}
}
