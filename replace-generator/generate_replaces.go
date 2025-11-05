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
	Comment     string `yaml:"comment"`
	Dependency  string `yaml:"dependency"`
	Replacement string `yaml:"replacement"`
}

type ProjectReplaces struct {
	Replaces []ReplaceEntry `yaml:"replaces"`
}

func main() {
	var projectReplacesPath string
	if len(os.Args) < 2 {
		// Default to local file when invoked via `go generate`
		projectReplacesPath = "project-replaces.yaml"
	} else {
		projectReplacesPath = os.Args[1]
	}

	// Determine base directory for templates relative to project-replaces.yaml
	absReplacesPath, err := filepath.Abs(projectReplacesPath)
	if err != nil {
		log.Fatalf("Failed to resolve path to project-replaces.yaml: %v", err)
	}
	baseDir := filepath.Dir(absReplacesPath)

	// Read and parse project-replaces.yaml
	var projectReplaces ProjectReplaces
	data, err := os.ReadFile(absReplacesPath)
	if err != nil {
		log.Fatalf("Failed to read project-replaces.yaml: %v", err)
	}

	if err := yaml.Unmarshal(data, &projectReplaces); err != nil {
		log.Fatalf("Failed to parse project-replaces.yaml: %v", err)
	}

	// Normalize comments (handle multi-line comments with >)
	normalizeComments(projectReplaces.Replaces)

	// Generate YAML output using templates in the same directory
	yamlTemplatePath := filepath.Join(baseDir, "replaces_yaml.tpl")
	yamlOutputPath := "builder-config-replaces.txt"
	if err := generateOutput(yamlTemplatePath, yamlOutputPath, projectReplaces.Replaces); err != nil {
		log.Fatalf("Failed to generate YAML output: %v", err)
	}
	log.Printf("✅ Generated YAML replaces: %s", yamlOutputPath)

	// Generate go.mod output
	gomodTemplatePath := filepath.Join(baseDir, "replaces_gomod.tpl")
	gomodOutputPath := "go-mod-replaces.txt"
	if err := generateOutput(gomodTemplatePath, gomodOutputPath, projectReplaces.Replaces); err != nil {
		log.Fatalf("Failed to generate go.mod output: %v", err)
	}
	log.Printf("✅ Generated go.mod replaces: %s", gomodOutputPath)
}

func normalizeComments(entries []ReplaceEntry) {
	for i := range entries {
		// YAML > syntax creates multi-line strings with newlines
		// Replace newlines with spaces and trim
		entries[i].Comment = strings.ReplaceAll(entries[i].Comment, "\n", " ")
		entries[i].Comment = strings.TrimSpace(entries[i].Comment)
	}
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
