package manifest

import (
	"fmt"
	"os"
	"strings"
)

const Version = "0.1.0"

const content = "version: 0.1.0\nprofile: workspace\n\ninstalled:\n  instructions: all\n  commands: all\n  skills: all\n  mcp: all\n"

func Write(path string) error {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		return fmt.Errorf("write manifest: %w", err)
	}
	return nil
}

func ReadVersion(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read manifest: %w", err)
	}
	for _, line := range strings.Split(string(data), "\n") {
		key, value, ok := strings.Cut(line, ":")
		if ok && strings.TrimSpace(key) == "version" {
			return strings.TrimSpace(value), nil
		}
	}
	return "", fmt.Errorf("manifest has no version")
}
