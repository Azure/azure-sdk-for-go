// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package build

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestValidatePath(t *testing.T) {
	// Test with current directory (should exist)
	if err := validatePath("."); err != nil {
		t.Errorf("validatePath failed for current directory: %v", err)
	}

	// Test with non-existent directory
	if err := validatePath("/non/existent/path"); err == nil {
		t.Error("validatePath should fail for non-existent path")
	}
}

func TestBuildResultJSON(t *testing.T) {
	result := &BuildResult{
		Success: true,
		Message: "Build successful",
		Path:    "/test/path",
	}

	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Errorf("Failed to marshal BuildResult: %v", err)
	}

	var unmarshaled BuildResult
	if err := json.Unmarshal(jsonData, &unmarshaled); err != nil {
		t.Errorf("Failed to unmarshal BuildResult: %v", err)
	}

	if unmarshaled.Success != result.Success {
		t.Error("JSON marshaling/unmarshaling failed for Success field")
	}
	if unmarshaled.Message != result.Message {
		t.Error("JSON marshaling/unmarshaling failed for Message field")
	}
}

// TestBuildAndVetIntegration tests the actual build and vet functionality
// This test requires a valid Go module to work properly
func TestBuildAndVetIntegration(t *testing.T) {
	// Create a temporary directory with a simple Go file
	tempDir, err := os.MkdirTemp("", "build_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a simple go.mod file
	goModContent := `module test

go 1.21
`
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create a valid Go file
	goFileContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, world!")
}
`
	if err := os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(goFileContent), 0644); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	// Test build and vet
	result, err := buildAndVet(tempDir, false)
	if err != nil {
		t.Errorf("buildAndVet failed: %v", err)
	}

	if !result.Success {
		t.Errorf("buildAndVet should succeed for valid Go code, but got: %s", result.Message)
	}
}

// TestBuildAndVetWithErrors tests the build and vet functionality with intentional errors
func TestBuildAndVetWithErrors(t *testing.T) {
	// Create a temporary directory with an invalid Go file
	tempDir, err := os.MkdirTemp("", "build_test_errors")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a simple go.mod file
	goModContent := `module test

go 1.21
`
	if err := os.WriteFile(filepath.Join(tempDir, "go.mod"), []byte(goModContent), 0644); err != nil {
		t.Fatalf("Failed to create go.mod: %v", err)
	}

	// Create an invalid Go file (syntax error)
	goFileContent := `package main

import "fmt"

func main() {
	fmt.Println("Hello, world!"
	// Missing closing parenthesis
}
`
	if err := os.WriteFile(filepath.Join(tempDir, "main.go"), []byte(goFileContent), 0644); err != nil {
		t.Fatalf("Failed to create main.go: %v", err)
	}

	// Test build and vet
	result, err := buildAndVet(tempDir, false)
	if err != nil {
		t.Errorf("buildAndVet failed: %v", err)
	}

	if result.Success {
		t.Error("buildAndVet should fail for invalid Go code")
	}

	if result.BuildOutput == "" {
		t.Error("buildAndVet should report build output for invalid Go code")
	}
}

// TestFilterHandWrittenTests verifies that hand-written test files are excluded
// from build while generated test files remain.
func TestFilterHandWrittenTests(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "build_test_filter")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	handWritten := filepath.Join(tempDir, "manual_test.go")
	if err := os.WriteFile(handWritten, []byte("package main\n"), 0644); err != nil {
		t.Fatalf("Failed to create hand-written test: %v", err)
	}
	generated := filepath.Join(tempDir, "zz_generated_test.go")
	genContent := generatedMarker + "\n\npackage main\n"
	if err := os.WriteFile(generated, []byte(genContent), 0644); err != nil {
		t.Fatalf("Failed to create generated test: %v", err)
	}

	restore, err := filterHandWrittenTests(tempDir)
	if err != nil {
		t.Fatalf("filterHandWrittenTests failed: %v", err)
	}

	if _, err := os.Stat(handWritten); !os.IsNotExist(err) {
		t.Error("hand-written test file should be filtered out during build")
	}
	if _, err := os.Stat(generated); err != nil {
		t.Error("generated test file should remain present during build")
	}

	restore()

	if _, err := os.Stat(handWritten); err != nil {
		t.Error("hand-written test file should be restored after build")
	}
}
