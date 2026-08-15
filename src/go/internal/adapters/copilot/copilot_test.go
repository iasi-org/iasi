package copilot

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testDescriptor = `schema_version: 1
id: copilot
platform: github-copilot
supports:
  instructions: true
instructions:
  general:
    type: repository
    target: .github/copilot-instructions.md
  documentation:
    type: path
    target: .github/instructions/documentation.instructions.md
    applyTo: "**/*.md"
`

func TestRunUsesInstalledContentAndGeneratesDeterministically(t *testing.T) {
	workspace := t.TempDir()
	project := filepath.Join(workspace, "project")
	writeFile(t, filepath.Join(workspace, ".iasi", "manifest.yml"), "version: 0.1.0\n")
	writeFile(t, filepath.Join(workspace, ".iasi", "adapters", "copilot", "adapter.yml"), testDescriptor)
	writeInstruction(t, filepath.Join(workspace, ".iasi", "instructions", "general", "rule.md"), "general.rule", "active", "general", "Installed rule")
	writeInstruction(t, filepath.Join(workspace, ".iasi", "instructions", "documentation", "guide.md"), "documentation.guide", "active", "documentation", "Installed guide")
	writeFile(t, filepath.Join(workspace, ".iasi", "instructions", "README.md"), "id: ignored")
	writeFile(t, filepath.Join(workspace, ".iasi", "instructions", "schema", "instructions.md"), "id: ignored")

	first, err := Run(project)
	if err != nil {
		t.Fatal(err)
	}
	generalPath := filepath.Join(project, ".github", "copilot-instructions.md")
	general, err := os.ReadFile(generalPath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(general), "Installed rule") || strings.Contains(string(general), "id: general.rule") {
		t.Fatalf("unexpected generated general instructions: %s", general)
	}
	if !strings.Contains(first, ".github/copilot-instructions.md") {
		t.Fatalf("unexpected success output: %s", first)
	}
	before, err := os.ReadFile(filepath.Join(project, ".github", "instructions", "documentation.instructions.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(before), "---\napplyTo: \"**/*.md\"\n---") || !strings.Contains(string(before), "<!-- IASI-GENERATED: copilot; version=0.1.0 -->") {
		t.Fatalf("unexpected path output: %s", before)
	}
	if _, err := Run(project); err != nil {
		t.Fatal(err)
	}
	after, err := os.ReadFile(filepath.Join(project, ".github", "instructions", "documentation.instructions.md"))
	if err != nil {
		t.Fatal(err)
	}
	if string(before) != string(after) {
		t.Fatal("repeated adaptation was not byte-for-byte identical")
	}
}

func TestRunRejectsUnknownActiveScopeWithoutChanges(t *testing.T) {
	workspace := t.TempDir()
	project := filepath.Join(workspace, "project")
	writeFile(t, filepath.Join(workspace, ".iasi", "manifest.yml"), "version: 0.1.0\n")
	writeFile(t, filepath.Join(workspace, ".iasi", "adapters", "copilot", "adapter.yml"), testDescriptor)
	writeInstruction(t, filepath.Join(workspace, ".iasi", "instructions", "unknown.md"), "unknown.rule", "active", "unknown", "Unsupported")

	if _, err := Run(project); err == nil || !strings.Contains(err.Error(), "Unsupported instruction scope") {
		t.Fatalf("expected unsupported scope error, got %v", err)
	}
	if _, err := os.Stat(filepath.Join(project, ".github")); !os.IsNotExist(err) {
		t.Fatalf("adaptation created project files after preflight failure: %v", err)
	}
}

func TestRunRejectsHumanOwnedTarget(t *testing.T) {
	workspace := t.TempDir()
	project := filepath.Join(workspace, "project")
	writeFile(t, filepath.Join(workspace, ".iasi", "manifest.yml"), "version: 0.1.0\n")
	writeFile(t, filepath.Join(workspace, ".iasi", "adapters", "copilot", "adapter.yml"), testDescriptor)
	writeInstruction(t, filepath.Join(workspace, ".iasi", "instructions", "rule.md"), "general.rule", "active", "general", "Rule")
	target := filepath.Join(project, ".github", "copilot-instructions.md")
	writeFile(t, target, "human content")

	if _, err := Run(project); err == nil {
		t.Fatal("expected collision error")
	}
	data, err := os.ReadFile(target)
	if err != nil || string(data) != "human content" {
		t.Fatalf("human-owned file changed: %s (%v)", data, err)
	}
}

func writeInstruction(t *testing.T, path, id, status, scope, body string) {
	writeFile(t, path, "---\nid: "+id+"\nstatus: "+status+"\nscope: "+scope+"\n---\n\n# "+body+"\n")
}

func writeFile(t *testing.T, path, content string) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
