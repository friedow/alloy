package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Get the directory where this script is located
	scriptDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	// Apply replaces to go.mod
	if err := applyReplacesToGoMod(scriptDir); err != nil {
		log.Fatalf("Failed to apply replaces to go.mod: %v", err)
	}

	// Apply replaces to builder-config.yaml
	if err := applyReplacesToBuilderConfig(scriptDir); err != nil {
		log.Fatalf("Failed to apply replaces to builder-config.yaml: %v", err)
	}

	// Apply replaces to extension/alloyengine/go.mod
	if err := applyReplacesToExtensionGoMod(scriptDir); err != nil {
		log.Fatalf("Failed to apply replaces to extension/alloyengine/go.mod: %v", err)
	}

	log.Println("✅ Successfully applied replaces to all target files")
}

// applyReplacesToGoMod removes existing generated replaces and appends new ones to go.mod
func applyReplacesToGoMod(scriptDir string) error {
	goModPath := filepath.Join(scriptDir, "..", "..", "go.mod")
	replacesPath := filepath.Join(scriptDir, "go-mod-replaces.txt")

	// Read the generated replaces file
	replacesContent, err := os.ReadFile(replacesPath)
	if err != nil {
		return fmt.Errorf("read replaces file: %w", err)
	}

	// Process go.mod: remove old replaces and append new ones
	content, err := os.ReadFile(goModPath)
	if err != nil {
		return fmt.Errorf("read go.mod: %w", err)
	}

	// Remove content between markers (including the full marker text with " - DO NOT EDIT")
	newContent := removeBetweenMarkers(string(content), "// BEGIN GENERATED REPLACES - DO NOT EDIT", "// END GENERATED REPLACES")

	// Append new replaces
	newContent = strings.TrimRight(newContent, "\n") + "\n" + string(replacesContent)

	// Write back to file
	if err := os.WriteFile(goModPath, []byte(newContent), 0o644); err != nil {
		return fmt.Errorf("write go.mod: %w", err)
	}

	log.Printf("✅ Updated %s", goModPath)
	return nil
}

// applyReplacesToBuilderConfig removes existing generated replaces and appends new ones to builder-config.yaml
func applyReplacesToBuilderConfig(scriptDir string) error {
	builderConfigPath := filepath.Join(scriptDir, "..", "..", "collector", "builder-config.yaml")
	replacesPath := filepath.Join(scriptDir, "builder-config-replaces.txt")

	// Read the generated replaces file
	replacesContent, err := os.ReadFile(replacesPath)
	if err != nil {
		return fmt.Errorf("read replaces file: %w", err)
	}

	// Process builder-config.yaml: remove old replaces and append new ones
	content, err := os.ReadFile(builderConfigPath)
	if err != nil {
		return fmt.Errorf("read builder-config.yaml: %w", err)
	}

	// Remove content between markers (with 2-space indentation, including the full marker text)
	newContent := removeBetweenMarkers(string(content), "  # BEGIN GENERATED REPLACES - DO NOT EDIT", "  # END GENERATED REPLACES")

	// Append new replaces
	newContent = strings.TrimRight(newContent, "\n") + "\n" + string(replacesContent)

	// Write back to file
	if err := os.WriteFile(builderConfigPath, []byte(newContent), 0o644); err != nil {
		return fmt.Errorf("write builder-config.yaml: %w", err)
	}

	log.Printf("✅ Updated %s", builderConfigPath)
	return nil
}

// applyReplacesToExtensionGoMod removes existing generated replaces and appends new ones to extension/alloyengine/go.mod
func applyReplacesToExtensionGoMod(scriptDir string) error {
	extensionGoModPath := filepath.Join(scriptDir, "..", "..", "extension", "alloyengine", "go.mod")
	replacesPath := filepath.Join(scriptDir, "extension-go-mod-replaces.txt")

	// Read the generated replaces file
	replacesContent, err := os.ReadFile(replacesPath)
	if err != nil {
		return fmt.Errorf("read replaces file: %w", err)
	}

	// Process extension/alloyengine/go.mod: remove old replaces and append new ones
	content, err := os.ReadFile(extensionGoModPath)
	if err != nil {
		return fmt.Errorf("read extension/alloyengine/go.mod: %w", err)
	}

	// Remove content between markers (including the full marker text with " - DO NOT EDIT")
	newContent := removeBetweenMarkers(string(content), "// BEGIN GENERATED REPLACES - DO NOT EDIT", "// END GENERATED REPLACES")

	// Append new replaces
	newContent = strings.TrimRight(newContent, "\n") + "\n" + string(replacesContent)

	// Write back to file
	if err := os.WriteFile(extensionGoModPath, []byte(newContent), 0o644); err != nil {
		return fmt.Errorf("write extension/alloyengine/go.mod: %w", err)
	}

	log.Printf("✅ Updated %s", extensionGoModPath)
	return nil
}

// removeBetweenMarkers removes all lines between (and including) the start and end markers
// Markers are matched by checking if the trimmed line starts with the trimmed marker text
func removeBetweenMarkers(content, startMarker, endMarker string) string {
	var result strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(content))
	inMarkerBlock := false

	// Normalize markers for comparison (trim whitespace)
	startMarkerTrimmed := strings.TrimSpace(startMarker)
	endMarkerTrimmed := strings.TrimSpace(endMarker)

	for scanner.Scan() {
		line := scanner.Text()
		trimmedLine := strings.TrimSpace(line)

		// Check if we hit the start marker (line starts with the marker text)
		// This handles cases where the marker might have additional text like " - DO NOT EDIT"
		if !inMarkerBlock && strings.HasPrefix(trimmedLine, startMarkerTrimmed) {
			inMarkerBlock = true
			continue // Skip the start marker line
		}

		// Check if we hit the end marker (exact match for end marker)
		if inMarkerBlock && trimmedLine == endMarkerTrimmed {
			inMarkerBlock = false
			continue // Skip the end marker line
		}

		// If we're not in a marker block, keep the line
		if !inMarkerBlock {
			result.WriteString(line)
			result.WriteString("\n")
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Warning: error scanning content: %v", err)
	}

	return result.String()
}
