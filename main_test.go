package main

import (
	"os"
	"path/filepath"
	"testing"
)

// Test for createDirectoryStructure
func TestCreateDirectoryStructure(t *testing.T) {
	programType := "testType"
	programName := "testProgram"
	domain := "testDomain"

	expectedPath := filepath.Join(programType, programName, domain)

	// Call the function
	actualPath := createDirectoryStructure(programType, programName, domain)

	// Check if the directory was created correctly
	if actualPath != expectedPath {
		t.Errorf("Expected path %s, but got %s", expectedPath, actualPath)
	}

	// Clean up the created directory
	defer os.RemoveAll(filepath.Join(programType))
}

// Test for addProgram
func TestAddProgram(t *testing.T) {
	programType := "testType"
	programName := "testProgram"
	domains := []string{"domain1.com", "domain2.com"}

	// Call the function
	directories := addProgram(programType, programName, domains)

	// Check if the directories were created correctly
	for i, domain := range domains {
		expectedPath := filepath.Join(programType, programName, domain)
		if directories[i] != expectedPath {
			t.Errorf("Expected path %s, but got %s", expectedPath, directories[i])
		}
	}

	// Clean up the created directories
	defer os.RemoveAll(filepath.Join(programType))
}

// Test for loadCommandConfig
func TestLoadCommandConfig(t *testing.T) {
	// Create a temporary YAML file for testing
	tempFile, err := os.CreateTemp("", "config.yml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up

	// Write test content to the temp file
	testContent := "command: \"echo %s\"\n"
	if _, err := tempFile.WriteString(testContent); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Close the file
	tempFile.Close()

	// Load the command from the temp file
	command, err := loadCommandConfig(tempFile.Name())
	if err != nil {
		t.Fatalf("Error loading command configuration: %v", err)
	}

	// Check if the command is as expected
	expectedCommand := "echo %s"
	if command != expectedCommand {
		t.Errorf("Expected command %s, but got %s", expectedCommand, command)
	}

	defer os.Remove(tempFile.Name()) // Add this to clean up the temporary file
}
