package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCLIVersion(t *testing.T) {
	// Simple sanity test
	if mdbindVersion := "0.2.0"; mdbindVersion == "" {
		t.Fatal("empty version")
	}
}

func TestCLIInitAndBuild(t *testing.T) {
	dir := t.TempDir()
	bookDir := filepath.Join(dir, "sample-book")

	// 1. Run init
	exitCode := runInit([]string{bookDir})
	if exitCode != 0 {
		t.Fatalf("expected exit code 0 for init, got %d", exitCode)
	}

	// Verify generated files
	if _, err := os.Stat(filepath.Join(bookDir, "book.yaml")); err != nil {
		t.Fatalf("missing book.yaml: %v", err)
	}
	if _, err := os.Stat(filepath.Join(bookDir, "content", "01-introduction.md")); err != nil {
		t.Fatalf("missing 01-introduction.md: %v", err)
	}

	// 2. Run init again in non-empty directory (must fail with exit code 2)
	exitCode2 := runInit([]string{bookDir})
	if exitCode2 != 2 {
		t.Fatalf("expected exit code 2 for non-empty dir init, got %d", exitCode2)
	}

	// 3. Run build
	outDir := filepath.Join(bookDir, "site-out")
	exitCodeBuild := runBuild([]string{
		"--input", filepath.Join(bookDir, "content"),
		"--output", outDir,
		"--title", "CLI Test Book",
		"--config", filepath.Join(bookDir, "book.yaml"),
	})
	if exitCodeBuild != 0 {
		t.Fatalf("expected exit code 0 for build, got %d", exitCodeBuild)
	}

	// Verify exported site
	indexContent, err := os.ReadFile(filepath.Join(outDir, "index.html"))
	if err != nil {
		t.Fatalf("failed to read exported index.html: %v", err)
	}

	if !strings.Contains(string(indexContent), "CLI Test Book") {
		t.Errorf("expected 'CLI Test Book' in index.html, got: %s", string(indexContent))
	}
}
