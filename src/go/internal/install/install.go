package install

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"iasi-cli/internal/manifest"
)

var categories = []string{"instructions", "commands", "skills", "mcp"}

func Run(workspace, sourceRoot string) (string, error) {
	target := filepath.Join(workspace, ".iasi")
	if _, err := os.Lstat(target); err == nil {
		return "", fmt.Errorf("IASI is already installed in this workspace: %s", target)
	} else if !os.IsNotExist(err) {
		return "", fmt.Errorf("check installation path: %w", err)
	}

	temporary := filepath.Join(workspace, fmt.Sprintf(".iasi.tmp-%d", time.Now().UnixNano()))
	if err := os.Mkdir(temporary, 0o755); err != nil {
		return "", fmt.Errorf("create installation directory: %w", err)
	}
	defer os.RemoveAll(temporary)

	for _, category := range categories {
		destination := filepath.Join(temporary, category)
		if err := os.MkdirAll(destination, 0o755); err != nil {
			return "", fmt.Errorf("create %s directory: %w", category, err)
		}
		source := filepath.Join(sourceRoot, "agentics", category)
		if err := copyDirectory(source, destination); err != nil {
			return "", fmt.Errorf("copy %s: %w", category, err)
		}
	}
	if err := manifest.Write(filepath.Join(temporary, "manifest.yml")); err != nil {
		return "", err
	}
	if err := os.Rename(temporary, target); err != nil {
		return "", fmt.Errorf("complete installation: %w", err)
	}
	return target, nil
}

func copyDirectory(source, destination string) error {
	info, err := os.Stat(source)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return fmt.Errorf("source is not a directory: %s", source)
	}
	return filepath.Walk(source, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		relative, err := filepath.Rel(source, path)
		if err != nil {
			return err
		}
		if relative == "." {
			return nil
		}
		target := filepath.Join(destination, relative)
		if info.IsDir() {
			return os.MkdirAll(target, info.Mode().Perm())
		}
		return copyFile(path, target, info.Mode().Perm())
	})
}

func copyFile(source, destination string, mode os.FileMode) error {
	input, err := os.Open(source)
	if err != nil {
		return err
	}
	defer input.Close()
	output, err := os.OpenFile(destination, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer output.Close()
	_, err = io.Copy(output, input)
	return err
}
