package source

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFindConfiguredSource(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "agentics"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("IASI_SOURCE_ROOT", root)

	found, err := New().Find()
	if err != nil || found != root {
		t.Fatalf("expected %s, got %s, %v", root, found, err)
	}
}
