package status

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"iasi-cli/internal/manifest"
)

var ErrNotInstalled = errors.New("IASI is not installed for this location")

type Result struct {
	Path, Type, Version string
	Counts              map[string]int
}

func Find(start string) (Result, error) {
	for current := filepath.Clean(start); ; current = filepath.Dir(current) {
		installation := filepath.Join(current, ".iasi")
		if info, err := os.Stat(installation); err == nil && info.IsDir() {
			version, err := manifest.ReadVersion(filepath.Join(installation, "manifest.yml"))
			if err != nil {
				return Result{}, err
			}
			counts := make(map[string]int)
			for _, category := range []string{"instructions", "commands", "skills", "mcp"} {
				counts[category] = countFiles(filepath.Join(installation, category))
			}
			return Result{Path: installation, Type: "workspace", Version: version, Counts: counts}, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
	}
	return Result{}, ErrNotInstalled
}

func Format(result Result) string {
	return fmt.Sprintf("IASI\n\nType    : %s\nPath    : %s\nVersion : %s\n\nInstructions : %d\nCommands     : %d\nSkills       : %d\nMCP          : %d\n", result.Type, result.Path, result.Version, result.Counts["instructions"], result.Counts["commands"], result.Counts["skills"], result.Counts["mcp"])
}

func countFiles(root string) int {
	count := 0
	filepath.Walk(root, func(_ string, info os.FileInfo, err error) error {
		if err == nil && info.Mode().IsRegular() {
			count++
		}
		return nil
	})
	return count
}
