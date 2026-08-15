package copilot

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
	"iasi-cli/internal/status"
)

const schemaVersion = 1
const ownershipMarker = "IASI-GENERATED: copilot"
const manifestRelativePath = ".github/.iasi/copilot-manifest.yml"

type descriptor struct {
	SchemaVersion int                    `yaml:"schema_version"`
	ID            string                 `yaml:"id"`
	Platform      string                 `yaml:"platform"`
	Supports      map[string]bool        `yaml:"supports"`
	Instructions  map[string]instruction `yaml:"instructions"`
}

type instruction struct {
	Type    string `yaml:"type"`
	Target  string `yaml:"target"`
	ApplyTo string `yaml:"applyTo"`
}

type metadata struct {
	ID     string `yaml:"id"`
	Status string `yaml:"status"`
	Scope  string `yaml:"scope"`
}

type candidate struct {
	ID, Scope, Body, Path string
	Status                string
}

type output struct {
	Path string
	Data []byte
}

type generatedManifest struct {
	SchemaVersion int      `yaml:"schema_version"`
	Adapter       string   `yaml:"adapter"`
	IASIVersion   string   `yaml:"iasi_version"`
	Generated     []string `yaml:"generated"`
}

func Run(project string) (string, error) {
	installation, err := status.Find(project)
	if err != nil {
		return "", err
	}
	adapterRoot := filepath.Join(installation.Path, "adapters", "copilot")
	desc, err := loadDescriptor(filepath.Join(adapterRoot, "adapter.yml"))
	if err != nil {
		return "", err
	}
	if !desc.Supports["instructions"] {
		return "", errors.New("Copilot adapter does not support instructions")
	}

	candidates, err := discover(filepath.Join(installation.Path, "instructions"), desc)
	if err != nil {
		return "", err
	}
	previous, manifestExists, err := loadManifest(filepath.Join(project, manifestRelativePath))
	if err != nil {
		return "", err
	}

	outputs, err := buildOutputs(candidates, desc, installation.InstalledVersion)
	if err != nil {
		return "", err
	}
	manifestData, err := marshalManifest(installation.InstalledVersion, outputs)
	if err != nil {
		return "", err
	}
	manifestOutput := output{Path: manifestRelativePath, Data: manifestData}

	stale := []string{}
	if manifestExists {
		stale = stalePaths(previous, outputPaths(outputs))
	}
	if err := preflight(project, outputs, manifestOutput, stale); err != nil {
		return "", err
	}
	if err := commit(project, outputs, manifestOutput, stale); err != nil {
		return "", err
	}

	all := outputPaths(outputs)
	return formatSuccess(installation.Path, project, all), nil
}

func loadDescriptor(path string) (descriptor, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return descriptor{}, errors.New("Copilot adapter is not available in this IASI installation")
		}
		return descriptor{}, fmt.Errorf("read Copilot adapter: %w", err)
	}
	var result descriptor
	if err := yaml.Unmarshal(data, &result); err != nil {
		return descriptor{}, fmt.Errorf("invalid Copilot adapter descriptor: %w", err)
	}
	if result.SchemaVersion != schemaVersion || result.ID != "copilot" || result.Platform != "github-copilot" {
		return descriptor{}, errors.New("invalid Copilot adapter descriptor")
	}
	if len(result.Instructions) == 0 {
		return descriptor{}, errors.New("Copilot adapter has no instruction mappings")
	}
	for scope, mapping := range result.Instructions {
		if mapping.Type != "repository" && mapping.Type != "path" {
			return descriptor{}, fmt.Errorf("invalid Copilot mapping type for scope: %s", scope)
		}
		if err := validateTarget(mapping.Target); err != nil {
			return descriptor{}, err
		}
		if mapping.Type == "path" && mapping.ApplyTo == "" {
			return descriptor{}, fmt.Errorf("missing applyTo for scope: %s", scope)
		}
	}
	return result, nil
}

func discover(root string, desc descriptor) ([]candidate, error) {
	var candidates []candidate
	ids := map[string]string{}
	err := filepath.Walk(root, func(path string, info os.FileInfo, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if info.IsDir() {
			if path != root && filepath.Base(path) == "schema" {
				return filepath.SkipDir
			}
			return nil
		}
		base := filepath.Base(path)
		if !strings.HasSuffix(strings.ToLower(base), ".md") || strings.HasPrefix(strings.ToLower(base), "readme") {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		meta, body, err := parseInstruction(data)
		if err != nil {
			return fmt.Errorf("invalid instruction %s: %w", path, err)
		}
		if meta.ID == "" || meta.Status == "" || meta.Scope == "" {
			return fmt.Errorf("instruction metadata is incomplete: %s", path)
		}
		if meta.Status != "active" && meta.Status != "draft" && meta.Status != "deprecated" {
			return fmt.Errorf("invalid instruction status %q: %s", meta.Status, path)
		}
		if old, exists := ids[meta.ID]; exists {
			return fmt.Errorf("duplicate instruction ID %q in %s and %s", meta.ID, old, path)
		}
		ids[meta.ID] = path
		if meta.Status == "active" {
			if _, ok := desc.Instructions[meta.Scope]; !ok {
				return fmt.Errorf("Unsupported instruction scope for Copilot adapter: %s", meta.Scope)
			}
		}
		candidates = append(candidates, candidate{ID: meta.ID, Scope: meta.Scope, Status: meta.Status, Body: body, Path: path})
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(candidates, func(i, j int) bool { return candidates[i].ID < candidates[j].ID })
	return candidates, nil
}

func parseInstruction(data []byte) (metadata, string, error) {
	text := strings.ReplaceAll(strings.TrimPrefix(string(data), "\ufeff"), "\r\n", "\n")
	if !strings.HasPrefix(text, "---\n") {
		return metadata{}, "", errors.New("missing YAML front matter")
	}
	end := strings.Index(text[4:], "\n---")
	if end < 0 {
		return metadata{}, "", errors.New("unterminated YAML front matter")
	}
	end += 4
	var result metadata
	if err := yaml.Unmarshal([]byte(text[4:end]), &result); err != nil {
		return metadata{}, "", fmt.Errorf("invalid YAML front matter: %w", err)
	}
	body := text[end+4:]
	body = strings.TrimPrefix(body, "\n")
	return result, body, nil
}

func buildOutputs(candidates []candidate, desc descriptor, version string) ([]output, error) {
	groups := map[string][]candidate{}
	for _, candidate := range candidates {
		if candidate.Status == "active" {
			groups[candidate.Scope] = append(groups[candidate.Scope], candidate)
		}
	}
	var outputs []output
	for scope, items := range groups {
		mapping := desc.Instructions[scope]
		sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
		var builder strings.Builder
		if mapping.Type == "path" {
			builder.WriteString("---\napplyTo: \"")
			builder.WriteString(mapping.ApplyTo)
			builder.WriteString("\"\n---\n\n")
		}
		builder.WriteString("<!-- IASI-GENERATED: copilot; version=")
		builder.WriteString(version)
		builder.WriteString(" -->\nGenerated from IASI. Do not edit this file manually.\n\n")
		for _, item := range items {
			builder.WriteString("## IASI: ")
			builder.WriteString(item.ID)
			builder.WriteString("\n\n")
			builder.WriteString(item.Body)
			if !strings.HasSuffix(item.Body, "\n") {
				builder.WriteByte('\n')
			}
			builder.WriteByte('\n')
		}
		outputs = append(outputs, output{Path: mapping.Target, Data: []byte(builder.String())})
	}
	sort.Slice(outputs, func(i, j int) bool { return outputs[i].Path < outputs[j].Path })
	return outputs, nil
}

func loadManifest(path string) (generatedManifest, bool, error) {
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return generatedManifest{}, false, nil
	}
	if err != nil {
		return generatedManifest{}, false, err
	}
	if !strings.Contains(string(data), "# IASI-GENERATED") {
		return generatedManifest{}, true, errors.New("invalid Copilot manifest ownership")
	}
	var result generatedManifest
	if err := yaml.Unmarshal(data, &result); err != nil {
		return generatedManifest{}, true, fmt.Errorf("invalid Copilot manifest: %w", err)
	}
	if result.SchemaVersion != schemaVersion || result.Adapter != "copilot" {
		return generatedManifest{}, true, errors.New("invalid Copilot manifest")
	}
	for _, path := range result.Generated {
		if err := validateTarget(path); err != nil {
			return generatedManifest{}, true, fmt.Errorf("invalid Copilot manifest target: %w", err)
		}
	}
	return result, true, nil
}

func marshalManifest(version string, outputs []output) ([]byte, error) {
	paths := outputPaths(outputs)
	data, err := yaml.Marshal(generatedManifest{SchemaVersion: schemaVersion, Adapter: "copilot", IASIVersion: version, Generated: paths})
	if err != nil {
		return nil, err
	}
	return append([]byte("# IASI-GENERATED: copilot\n"), data...), nil
}

func preflight(project string, outputs []output, manifest output, stale []string) error {
	seen := map[string]bool{}
	for _, item := range append(outputs, manifest) {
		if seen[item.Path] {
			return fmt.Errorf("duplicate output target: %s", item.Path)
		}
		seen[item.Path] = true
		path := filepath.Join(project, filepath.FromSlash(item.Path))
		if data, err := os.ReadFile(path); err == nil && !strings.Contains(string(data), ownershipMarker) {
			return fmt.Errorf("Cannot generate Copilot instructions because this file already exists and is not managed by IASI: %s", item.Path)
		} else if err != nil && !os.IsNotExist(err) {
			return err
		}
	}
	for _, relative := range stale {
		if seen[relative] {
			continue
		}
		path := filepath.Join(project, filepath.FromSlash(relative))
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		if !strings.Contains(string(data), ownershipMarker) {
			return fmt.Errorf("stale Copilot output is not owned by IASI: %s", relative)
		}
	}
	return nil
}

func commit(project string, outputs []output, manifest output, stale []string) error {
	temp, err := os.MkdirTemp("", "iasi-copilot-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(temp)
	all := append(outputs, manifest)
	for _, item := range all {
		path := filepath.Join(temp, filepath.FromSlash(item.Path))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(path, item.Data, 0o644); err != nil {
			return err
		}
	}
	affected := map[string]string{}
	backupDir := filepath.Join(temp, "backup")
	if err := os.MkdirAll(backupDir, 0o755); err != nil {
		return err
	}
	for _, item := range all {
		realPath := filepath.Join(project, filepath.FromSlash(item.Path))
		if _, err := os.Stat(realPath); err == nil {
			backup := filepath.Join(backupDir, fmt.Sprintf("%d", len(affected)))
			if err := os.Rename(realPath, backup); err != nil {
				return rollback(affected, nil, err)
			}
			affected[realPath] = backup
		}
	}
	for _, relative := range stale {
		realPath := filepath.Join(project, filepath.FromSlash(relative))
		if _, err := os.Stat(realPath); err == nil {
			backup := filepath.Join(backupDir, fmt.Sprintf("%d", len(affected)))
			if err := os.Rename(realPath, backup); err != nil {
				return rollback(affected, nil, err)
			}
			affected[realPath] = backup
		}
	}
	created := []string{}
	for _, item := range all {
		realPath := filepath.Join(project, filepath.FromSlash(item.Path))
		if err := os.MkdirAll(filepath.Dir(realPath), 0o755); err != nil {
			return rollback(affected, created, err)
		}
		staged := filepath.Join(temp, filepath.FromSlash(item.Path))
		if err := os.Rename(staged, realPath); err != nil {
			return rollback(affected, created, err)
		}
		created = append(created, realPath)
	}
	return nil
}

func rollback(affected map[string]string, created []string, original error) error {
	for _, path := range created {
		_ = os.Remove(path)
	}
	for path, backup := range affected {
		_ = os.MkdirAll(filepath.Dir(path), 0o755)
		_ = os.Rename(backup, path)
	}
	return original
}

func stalePaths(previous generatedManifest, current []string) []string {
	known := map[string]bool{}
	for _, path := range current {
		known[path] = true
	}
	var stale []string
	for _, path := range previous.Generated {
		if !known[path] {
			stale = append(stale, path)
		}
	}
	return stale
}

func outputPaths(outputs []output) []string {
	paths := make([]string, 0, len(outputs))
	for _, item := range outputs {
		paths = append(paths, item.Path)
	}
	sort.Strings(paths)
	return paths
}

func validateTarget(target string) error {
	clean := filepath.ToSlash(filepath.Clean(target))
	if filepath.IsAbs(target) || clean == ".github" || !strings.HasPrefix(clean, ".github/") || strings.Contains(clean, "../") {
		return fmt.Errorf("invalid Copilot target: %s", target)
	}
	return nil
}

func formatSuccess(source, project string, paths []string) string {
	var builder strings.Builder
	builder.WriteString("IASI Copilot adapter\n\nSource : ")
	builder.WriteString(source)
	builder.WriteString("\nTarget : ")
	builder.WriteString(project)
	builder.WriteString("\n\nGenerated:\n")
	for _, path := range paths {
		builder.WriteString("  ")
		builder.WriteString(path)
		builder.WriteByte('\n')
	}
	return builder.String()
}
