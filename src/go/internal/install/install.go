package install

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"iasi-cli/internal/manifest"
)

var categories = []string{"instructions", "commands", "skills", "mcp", "adapters"}

func Run(workspace string, methodology fs.FS, version string) (string, error) {
	target := filepath.Join(workspace, ".iasi")
	if _, err := os.Lstat(target); err == nil {
		return "", fmt.Errorf("IASI is already installed in this workspace: %s", target)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("check installation path: %w", err)
	}

	temporary, err := os.MkdirTemp(workspace, ".iasi.tmp-")
	if err != nil {
		return "", fmt.Errorf("create installation directory: %w", err)
	}
	defer os.RemoveAll(temporary)

	for _, category := range categories {
		destination := filepath.Join(temporary, category)
		if err := os.MkdirAll(destination, 0o755); err != nil {
			return "", fmt.Errorf("create %s directory: %w", category, err)
		}
		if err := copyDirectory(methodology, filepath.ToSlash(filepath.Join("agentics", category)), destination); err != nil {
			return "", fmt.Errorf("copy %s: %w", category, err)
		}
	}
	if err := manifest.Write(filepath.Join(temporary, "manifest.yml"), version); err != nil {
		return "", err
	}
	if err := os.Rename(temporary, target); err != nil {
		return "", fmt.Errorf("complete installation: %w", err)
	}
	return target, nil
}

func copyDirectory(source fs.FS, directory, destination string) error {
	if _, err := fs.Stat(source, directory); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	return fs.WalkDir(source, directory, func(path string, info fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(directory, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}
		target := filepath.Join(destination, relative)
		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		data, err := fs.ReadFile(source, path)
		if err != nil {
			return err
		}
		return os.WriteFile(target, data, 0o644)
	})
}
