package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"gopkg.in/yaml.v3"
)

type ReplaceEntry struct {
	Comment     string   `yaml:"comment"`
	Dependency  string   `yaml:"dependency"`
	Replacement string   `yaml:"replacement"`
	Scope       []string `yaml:"scope"`
}

type ProjectReplaces struct {
	Replaces []ReplaceEntry `yaml:"replaces"`
}

func main() {
	var projectReplacesPath string
	if len(os.Args) < 2 {
		// Default to local file when invoked via `go generate`
		projectReplacesPath = "dependency-replacements.yaml"
	} else {
		projectReplacesPath = os.Args[1]
	}

	// Determine base directory for templates relative to dependency-replacements.yaml
	absReplacesPath, err := filepath.Abs(projectReplacesPath)
	if err != nil {
		log.Fatalf("Failed to resolve path to dependency-replacements.yaml: %v", err)
	}
	baseDir := filepath.Dir(absReplacesPath)

	// Read and parse dependency-replacements.yaml
	var projectReplaces ProjectReplaces
	data, err := os.ReadFile(absReplacesPath)
	if err != nil {
		log.Fatalf("Failed to read dependency-replacements.yaml: %v", err)
	}

	if err := yaml.Unmarshal(data, &projectReplaces); err != nil {
		log.Fatalf("Failed to parse dependency-replacements.yaml: %v", err)
	}

	// Normalize comments (handle multi-line comments with >)
	normalizeComments(projectReplaces.Replaces)

	// Set default scope for entries that don't have it
	setDefaultScope(projectReplaces.Replaces)

	// Filter replaces by scope and generate outputs
	alloyReplaces := filterByScope(projectReplaces.Replaces, "alloy")
	collectorReplaces := filterByScope(projectReplaces.Replaces, "collector")
	extensionReplaces := filterByScope(projectReplaces.Replaces, "extension")

	// Generate alloy go.mod output (using shared go.mod template)
	gomodTemplatePath := filepath.Join(baseDir, "replaces_alloy_gomod.tpl")
	alloyOutputPath := "go-mod-replaces.txt"
	if err := generateOutput(gomodTemplatePath, alloyOutputPath, alloyReplaces); err != nil {
		log.Fatalf("Failed to generate alloy go.mod output: %v", err)
	}
	log.Printf("✅ Generated alloy go.mod replaces: %s", alloyOutputPath)

	// Generate collector YAML output
	collectorTemplatePath := filepath.Join(baseDir, "replaces_yaml.tpl")
	collectorOutputPath := "builder-config-replaces.txt"
	if err := generateOutput(collectorTemplatePath, collectorOutputPath, collectorReplaces); err != nil {
		log.Fatalf("Failed to generate collector YAML output: %v", err)
	}
	log.Printf("✅ Generated collector YAML replaces: %s", collectorOutputPath)

	// Generate extension go.mod output (using shared go.mod template)
	extensionOutputPath := "extension-go-mod-replaces.txt"
	if err := generateOutput(gomodTemplatePath, extensionOutputPath, extensionReplaces); err != nil {
		log.Fatalf("Failed to generate extension go.mod output: %v", err)
	}
	log.Printf("✅ Generated extension go.mod replaces: %s", extensionOutputPath)
}

func normalizeComments(entries []ReplaceEntry) {
	for i := range entries {
		// YAML > syntax creates multi-line strings with newlines
		// Replace newlines with spaces and trim
		entries[i].Comment = strings.ReplaceAll(entries[i].Comment, "\n", " ")
		entries[i].Comment = strings.TrimSpace(entries[i].Comment)
	}
}

// setDefaultScope sets default scope to [alloy, collector, extension] if scope is empty
func setDefaultScope(entries []ReplaceEntry) {
	for i := range entries {
		if len(entries[i].Scope) == 0 {
			entries[i].Scope = []string{"alloy", "collector", "extension"}
		}
	}
}

// filterByScope returns only entries that have the specified scope
func filterByScope(entries []ReplaceEntry, scope string) []ReplaceEntry {
	var filtered []ReplaceEntry
	for _, entry := range entries {
		for _, s := range entry.Scope {
			if s == scope {
				filtered = append(filtered, entry)
				break
			}
		}
	}
	return filtered
}

func generateOutput(templatePath, outputPath string, data []ReplaceEntry) error {
	// Read template
	tmplContent, err := os.ReadFile(templatePath)
	if err != nil {
		return fmt.Errorf("read template %s: %w", templatePath, err)
	}

	// Parse template
	tmpl, err := template.New(filepath.Base(templatePath)).Parse(string(tmplContent))
	if err != nil {
		return fmt.Errorf("parse template: %w", err)
	}

	// Create output file
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}
	defer outputFile.Close()

	// Execute template
	if err := tmpl.Execute(outputFile, data); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}
