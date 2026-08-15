package source

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type Locator struct {
	Executable string
	WorkingDir string
}

func New() Locator {
	workingDir, _ := os.Getwd()
	executable, _ := os.Executable()
	return Locator{Executable: executable, WorkingDir: workingDir}
}

func (l Locator) Find() (string, error) {
	if configured := os.Getenv("IASI_SOURCE_ROOT"); configured != "" {
		return validate(configured)
	}
	for _, start := range []string{l.WorkingDir, filepath.Dir(l.Executable)} {
		if start == "" {
			continue
		}
		for current := filepath.Clean(start); ; current = filepath.Dir(current) {
			if root, err := validate(current); err == nil {
				return root, nil
			}
			parent := filepath.Dir(current)
			if parent == current {
				break
			}
		}
	}
	return "", errors.New("could not locate IASI source directory (agentics)")
}

func validate(root string) (string, error) {
	root = filepath.Clean(root)
	info, err := os.Stat(filepath.Join(root, "agentics"))
	if err != nil || !info.IsDir() {
		return "", fmt.Errorf("%s does not contain agentics", root)
	}
	return root, nil
}
