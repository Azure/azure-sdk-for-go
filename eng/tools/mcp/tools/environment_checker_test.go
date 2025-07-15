// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package tools

import (
	"context"
	"encoding/json"
	"os/exec"
	"strings"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
)

// Mock command execution for testing
type mockCommandExecutor struct {
	responses map[string]mockResponse
}

type mockResponse struct {
	output string
	err    error
}

func (m *mockCommandExecutor) runCommand(command string, args ...string) (string, error) {
	key := command + " " + strings.Join(args, " ")
	if response, exists := m.responses[key]; exists {
		return response.output, response.err
	}
	// Default to command not found
	return "", &exec.ExitError{}
}

// Override the global runCommand function for testing
var testCommandExecutor *mockCommandExecutor

func init() {
	// Save original function if we need to restore it later
	originalRunCommand := runCommand
	_ = originalRunCommand // Prevent unused variable warning
}

func TestCheckGoVersion(t *testing.T) {
	tests := []struct {
		name           string
		commandOutput  string
		commandError   error
		minVersion     string
		expectedStatus string
		expectedName   string
	}{
		{
			name:           "Go installed with sufficient version",
			commandOutput:  "go version go1.23.0 windows/amd64",
			commandError:   nil,
			minVersion:     "1.23",
			expectedStatus: "SUCCESS",
			expectedName:   "Go",
		},
		{
			name:           "Go installed with old version",
			commandOutput:  "go version go1.20.0 windows/amd64",
			commandError:   nil,
			minVersion:     "1.23",
			expectedStatus: "ERROR",
			expectedName:   "Go",
		},
		{
			name:           "Go not installed",
			commandOutput:  "",
			commandError:   &exec.ExitError{},
			minVersion:     "1.23",
			expectedStatus: "ERROR",
			expectedName:   "Go",
		},
		{
			name:           "Go version cannot be parsed",
			commandOutput:  "invalid version output",
			commandError:   nil,
			minVersion:     "1.23",
			expectedStatus: "ERROR",
			expectedName:   "Go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			testCommandExecutor = &mockCommandExecutor{
				responses: map[string]mockResponse{
					"go version": {
						output: tt.commandOutput,
						err:    tt.commandError,
					},
				},
			}

			// Override runCommand for this test
			originalRunCommand := runCommand
			runCommand = testCommandExecutor.runCommand
			defer func() { runCommand = originalRunCommand }()

			result := checkGoVersion(tt.minVersion)

			if result.Name != tt.expectedName {
				t.Errorf("Expected name '%s', got '%s'", tt.expectedName, result.Name)
			}

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status '%s', got '%s'", tt.expectedStatus, result.Status)
			}

			if tt.expectedStatus == "SUCCESS" && result.Version == "" {
				t.Error("Expected version to be set for successful check")
			}

			if result.Message == "" {
				t.Error("Expected message to be set")
			}
		})
	}
}

func TestCheckNodeVersion(t *testing.T) {
	tests := []struct {
		name           string
		commandOutput  string
		commandError   error
		minVersion     string
		expectedStatus string
	}{
		{
			name:           "Node.js installed with sufficient version",
			commandOutput:  "v20.5.0",
			commandError:   nil,
			minVersion:     "20.0.0",
			expectedStatus: "SUCCESS",
		},
		{
			name:           "Node.js installed with old version",
			commandOutput:  "v18.0.0",
			commandError:   nil,
			minVersion:     "20.0.0",
			expectedStatus: "ERROR",
		},
		{
			name:           "Node.js not installed",
			commandOutput:  "",
			commandError:   &exec.ExitError{},
			minVersion:     "20.0.0",
			expectedStatus: "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			testCommandExecutor = &mockCommandExecutor{
				responses: map[string]mockResponse{
					"node --version": {
						output: tt.commandOutput,
						err:    tt.commandError,
					},
				},
			}

			// Override runCommand for this test
			originalRunCommand := runCommand
			runCommand = testCommandExecutor.runCommand
			defer func() { runCommand = originalRunCommand }()

			result := checkNodeVersion(tt.minVersion)

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status '%s', got '%s'", tt.expectedStatus, result.Status)
			}

			if result.Name != "Node.js" {
				t.Errorf("Expected name 'Node.js', got '%s'", result.Name)
			}
		})
	}
}

func TestCheckTypeSpecCompiler(t *testing.T) {
	tests := []struct {
		name           string
		commandError   error
		expectedStatus string
	}{
		{
			name:           "TypeSpec compiler installed",
			commandError:   nil,
			expectedStatus: "SUCCESS",
		},
		{
			name:           "TypeSpec compiler not installed",
			commandError:   &exec.ExitError{},
			expectedStatus: "WARNING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			testCommandExecutor = &mockCommandExecutor{
				responses: map[string]mockResponse{
					"tsp --version": {
						output: "1.0.0",
						err:    tt.commandError,
					},
				},
			}

			// Override runCommand for this test
			originalRunCommand := runCommand
			runCommand = testCommandExecutor.runCommand
			defer func() { runCommand = originalRunCommand }()

			result := checkTypeSpecCompiler()

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status '%s', got '%s'", tt.expectedStatus, result.Status)
			}

			if result.Name != "TypeSpec Compiler" {
				t.Errorf("Expected name 'TypeSpec Compiler', got '%s'", result.Name)
			}
		})
	}
}

func TestCheckTypeSpecClientCLI(t *testing.T) {
	tests := []struct {
		name           string
		commandError   error
		expectedStatus string
	}{
		{
			name:           "TypeSpec client CLI installed",
			commandError:   nil,
			expectedStatus: "SUCCESS",
		},
		{
			name:           "TypeSpec client CLI not installed",
			commandError:   &exec.ExitError{},
			expectedStatus: "WARNING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			testCommandExecutor = &mockCommandExecutor{
				responses: map[string]mockResponse{
					"tsp-client --version": {
						output: "1.0.0",
						err:    tt.commandError,
					},
				},
			}

			// Override runCommand for this test
			originalRunCommand := runCommand
			runCommand = testCommandExecutor.runCommand
			defer func() { runCommand = originalRunCommand }()

			result := checkTypeSpecClientCLI()

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status '%s', got '%s'", tt.expectedStatus, result.Status)
			}

			if result.Name != "TypeSpec Client Generator CLI" {
				t.Errorf("Expected name 'TypeSpec Client Generator CLI', got '%s'", result.Name)
			}
		})
	}
}

func TestInstallTypeSpecCompiler(t *testing.T) {
	tests := []struct {
		name           string
		commandError   error
		expectedStatus string
	}{
		{
			name:           "Installation successful",
			commandError:   nil,
			expectedStatus: "SUCCESS",
		},
		{
			name:           "Installation failed",
			commandError:   &exec.ExitError{},
			expectedStatus: "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			testCommandExecutor = &mockCommandExecutor{
				responses: map[string]mockResponse{
					"npm install -g @typespec/compiler": {
						output: "installed successfully",
						err:    tt.commandError,
					},
				},
			}

			// Override runCommand for this test
			originalRunCommand := runCommand
			runCommand = testCommandExecutor.runCommand
			defer func() { runCommand = originalRunCommand }()

			result := installTypeSpecCompiler()

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status '%s', got '%s'", tt.expectedStatus, result.Status)
			}

			if result.Name != "TypeSpec Compiler Installation" {
				t.Errorf("Expected name 'TypeSpec Compiler Installation', got '%s'", result.Name)
			}
		})
	}
}

func TestInstallTypeSpecClientCLI(t *testing.T) {
	tests := []struct {
		name           string
		commandError   error
		expectedStatus string
	}{
		{
			name:           "Installation successful",
			commandError:   nil,
			expectedStatus: "SUCCESS",
		},
		{
			name:           "Installation failed",
			commandError:   &exec.ExitError{},
			expectedStatus: "ERROR",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			testCommandExecutor = &mockCommandExecutor{
				responses: map[string]mockResponse{
					"npm install -g @azure-tools/typespec-client-generator-cli": {
						output: "installed successfully",
						err:    tt.commandError,
					},
				},
			}

			// Override runCommand for this test
			originalRunCommand := runCommand
			runCommand = testCommandExecutor.runCommand
			defer func() { runCommand = originalRunCommand }()

			result := installTypeSpecClientCLI()

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status '%s', got '%s'", tt.expectedStatus, result.Status)
			}

			if result.Name != "TypeSpec Client Generator CLI Installation" {
				t.Errorf("Expected name 'TypeSpec Client Generator CLI Installation', got '%s'", result.Name)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		minimum  string
		expected bool
	}{
		{
			name:     "Current version is higher",
			current:  "1.25.0",
			minimum:  "1.23.0",
			expected: true,
		},
		{
			name:     "Current version is equal",
			current:  "1.23.0",
			minimum:  "1.23.0",
			expected: true,
		},
		{
			name:     "Current version is lower",
			current:  "1.20.0",
			minimum:  "1.23.0",
			expected: false,
		},
		{
			name:     "Version without patch number",
			current:  "1.23",
			minimum:  "1.23.0",
			expected: true,
		},
		{
			name:     "Node.js style version",
			current:  "20.5.0",
			minimum:  "20.0.0",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := compareVersions(tt.current, tt.minimum)
			if result != tt.expected {
				t.Errorf("compareVersions(%s, %s) = %v, expected %v", tt.current, tt.minimum, result, tt.expected)
			}
		})
	}
}

func TestEnvironmentCheckerHandler(t *testing.T) {
	tests := []struct {
		name                string
		mockResponses       map[string]mockResponse
		expectedSuccess     bool
		expectedFailedCount int
	}{
		{
			name: "All environment checks satisfied",
			mockResponses: map[string]mockResponse{
				"go version":           {output: "go version go1.23.0 windows/amd64", err: nil},
				"node --version":       {output: "v20.5.0", err: nil},
				"tsp --version":        {output: "1.0.0", err: nil},
				"tsp-client --version": {output: "1.0.0", err: nil},
			},
			expectedSuccess:     true,
			expectedFailedCount: 0,
		},
		{
			name: "Missing Node.js",
			mockResponses: map[string]mockResponse{
				"go version":           {output: "go version go1.23.0 windows/amd64", err: nil},
				"node --version":       {output: "", err: &exec.ExitError{}},
				"tsp --version":        {output: "1.0.0", err: nil},
				"tsp-client --version": {output: "1.0.0", err: nil},
			},
			expectedSuccess:     false,
			expectedFailedCount: 1,
		},
		{
			name: "Missing TypeSpec tools but installation succeeds",
			mockResponses: map[string]mockResponse{
				"go version":                        {output: "go version go1.23.0 windows/amd64", err: nil},
				"node --version":                    {output: "v20.5.0", err: nil},
				"tsp --version":                     {output: "", err: &exec.ExitError{}},
				"npm install -g @typespec/compiler": {output: "installed", err: nil},
				"tsp-client --version":              {output: "", err: &exec.ExitError{}},
				"npm install -g @azure-tools/typespec-client-generator-cli": {output: "installed", err: nil},
			},
			expectedSuccess:     true,
			expectedFailedCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock
			testCommandExecutor = &mockCommandExecutor{
				responses: tt.mockResponses,
			}

			// Track call order for re-checks after installation
			callCount := make(map[string]int)
			originalRunCommand := runCommand
			runCommand = func(command string, args ...string) (string, error) {
				key := command + " " + strings.Join(args, " ")
				callCount[key]++

				// For re-checks after installation, return success
				if callCount[key] > 1 && (key == "tsp --version" || key == "tsp-client --version") {
					return "1.0.0", nil
				}

				return testCommandExecutor.runCommand(command, args...)
			}
			defer func() { runCommand = originalRunCommand }()

			// Create a mock request
			request := mcp.CallToolRequest{}

			result, err := EnvironmentCheckerHandler(context.Background(), request)
			if err != nil {
				t.Fatalf("Handler returned error: %v", err)
			}

			// Parse the JSON result
			var prereqResult EnvironmentCheckResult
			if len(result.Content) > 0 {
				if textContent, ok := result.Content[0].(mcp.TextContent); ok {
					err = json.Unmarshal([]byte(textContent.Text), &prereqResult)
					if err != nil {
						t.Fatalf("Failed to unmarshal result: %v", err)
					}
				} else {
					// Try to get the text directly from the result
					resultText := ""
					switch content := result.Content[0].(type) {
					case mcp.TextContent:
						resultText = content.Text
					case *mcp.TextContent:
						resultText = content.Text
					default:
						t.Logf("Content type: %T", content)
						t.Fatal("Unexpected content type in result")
					}

					err = json.Unmarshal([]byte(resultText), &prereqResult)
					if err != nil {
						t.Fatalf("Failed to unmarshal result: %v", err)
					}
				}
			} else {
				t.Fatal("Expected content in result")
			}

			if prereqResult.Success != tt.expectedSuccess {
				t.Errorf("Expected success %v, got %v", tt.expectedSuccess, prereqResult.Success)
			}

			if len(prereqResult.Failed) != tt.expectedFailedCount {
				t.Errorf("Expected %d failed environment checks, got %d: %v", tt.expectedFailedCount, len(prereqResult.Failed), prereqResult.Failed)
			}

			if len(prereqResult.Checks) == 0 {
				t.Error("Expected at least one check to be performed")
			}

			// Verify summary is set
			if prereqResult.Summary == "" {
				t.Error("Expected summary to be set")
			}
		})
	}
}
