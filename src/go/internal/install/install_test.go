package install

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunCopiesSourceAndCreatesEmptyCategories(t *testing.T) {
	sourceRoot := t.TempDir()
	workspace := t.TempDir()
	writeTestFile(t, filepath.Join(sourceRoot, "agentics", "instructions", "general", "behavior.md"), "be clear")
	writeTestFile(t, filepath.Join(sourceRoot, "agentics", "instructions", "documentation", "guide.md"), "guide")

	path, err := Run(workspace, os.DirFS(sourceRoot), "0.2.0")
	if err != nil {
		t.Fatal(err)
	}
	for _, category := range categories {
		if info, err := os.Stat(filepath.Join(path, category)); err != nil || !info.IsDir() {
			t.Fatalf("category %s was not created", category)
		}
	}
	manifest, err := os.ReadFile(filepath.Join(path, "manifest.yml"))
	if err != nil || !strings.Contains(string(manifest), "version: 0.2.0") {
		t.Fatalf("manifest was not created correctly: %v", err)
	}
	content, err := os.ReadFile(filepath.Join(path, "instructions", "general", "behavior.md"))
	if err != nil || string(content) != "be clear" {
		t.Fatalf("instruction was not copied: %v", err)
	}
}

func TestRunRejectsExistingInstallation(t *testing.T) {
	workspace := t.TempDir()
	if err := os.Mkdir(filepath.Join(workspace, ".iasi"), 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(workspace, os.DirFS(t.TempDir()), "0.2.0"); err == nil || !strings.Contains(err.Error(), "already installed") {
		t.Fatalf("expected existing installation error, got %v", err)
	}
}

func writeTestFile(t *testing.T, path, content string) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
