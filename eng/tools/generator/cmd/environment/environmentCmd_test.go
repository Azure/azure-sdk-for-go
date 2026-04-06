// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package environment

import (
	"testing"
)

// mockCommand allows us to override the runCommand function for testing
type mockCommand struct {
	output string
	err    error
}

func (m *mockCommand) run(command string, args ...string) (string, error) {
	return m.output, m.err
}

func TestCheckNodeVersion(t *testing.T) {
	tests := []struct {
		name           string
		commandOutput  string
		commandErr     error
		minVersion     string
		expectedStatus string
		expectedName   string
	}{
		{
			name:           "Valid Node.js version above minimum",
			commandOutput:  "v20.1.0",
			commandErr:     nil,
			minVersion:     "20.0.0",
			expectedStatus: "SUCCESS",
			expectedName:   "Node.js",
		},
		{
			name:           "Node.js version below minimum",
			commandOutput:  "v18.0.0",
			commandErr:     nil,
			minVersion:     "20.0.0",
			expectedStatus: "ERROR",
			expectedName:   "Node.js",
		},
		{
			name:           "Node.js not installed",
			commandOutput:  "",
			commandErr:     &mockError{msg: "command not found"},
			minVersion:     "20.0.0",
			expectedStatus: "ERROR",
			expectedName:   "Node.js",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original runCommand
			originalRunCommand := runCommand
			defer func() {
				runCommand = originalRunCommand
			}()

			// Override runCommand for this test
			runCommand = func(command string, args ...string) (string, error) {
				return tt.commandOutput, tt.commandErr
			}

			result := checkNodeVersion(tt.minVersion)

			if result.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, result.Name)
			}

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		current  string
		minimum  string
		expected bool
	}{
		{"1.24", "1.24", true},
		{"1.25", "1.24", true},
		{"1.23", "1.24", false},
		{"2.0", "1.24", true},
		{"20.1.0", "20.0.0", true},
		{"20.0.0", "20.0.0", true},
		{"19.9.0", "20.0.0", false},
	}

	for _, tt := range tests {
		t.Run(tt.current+"_vs_"+tt.minimum, func(t *testing.T) {
			result := compareVersions(tt.current, tt.minimum)
			if result != tt.expected {
				t.Errorf("compareVersions(%s, %s) = %v, expected %v", tt.current, tt.minimum, result, tt.expected)
			}
		})
	}
}

func TestCheckTypeSpecCompiler(t *testing.T) {
	tests := []struct {
		name           string
		commandOutput  string
		commandErr     error
		expectedStatus string
		expectedName   string
	}{
		{
			name:           "TypeSpec compiler installed",
			commandOutput:  "0.50.0",
			commandErr:     nil,
			expectedStatus: "SUCCESS",
			expectedName:   "TypeSpec Compiler",
		},
		{
			name:           "TypeSpec compiler not installed",
			commandOutput:  "",
			commandErr:     &mockError{msg: "command not found"},
			expectedStatus: "WARNING",
			expectedName:   "TypeSpec Compiler",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original runCommand
			originalRunCommand := runCommand
			defer func() {
				runCommand = originalRunCommand
			}()

			// Override runCommand for this test
			runCommand = func(command string, args ...string) (string, error) {
				return tt.commandOutput, tt.commandErr
			}

			result := checkTypeSpecCompiler()

			if result.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, result.Name)
			}

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}
		})
	}
}

func TestCheckTypeSpecClientCLI(t *testing.T) {
	tests := []struct {
		name           string
		commandOutput  string
		commandErr     error
		expectedStatus string
		expectedName   string
	}{
		{
			name:           "TypeSpec client CLI installed",
			commandOutput:  "1.0.0",
			commandErr:     nil,
			expectedStatus: "SUCCESS",
			expectedName:   "TypeSpec Client Generator CLI",
		},
		{
			name:           "TypeSpec client CLI not installed",
			commandOutput:  "",
			commandErr:     &mockError{msg: "command not found"},
			expectedStatus: "WARNING",
			expectedName:   "TypeSpec Client Generator CLI",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original runCommand
			originalRunCommand := runCommand
			defer func() {
				runCommand = originalRunCommand
			}()

			// Override runCommand for this test
			runCommand = func(command string, args ...string) (string, error) {
				return tt.commandOutput, tt.commandErr
			}

			result := checkTypeSpecClientCLI()

			if result.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, result.Name)
			}

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}
		})
	}
}

func TestInstallTypeSpecCompiler(t *testing.T) {
	tests := []struct {
		name           string
		commandOutput  string
		commandErr     error
		expectedStatus string
		expectedName   string
	}{
		{
			name:           "Installation successful",
			commandOutput:  "installed successfully",
			commandErr:     nil,
			expectedStatus: "SUCCESS",
			expectedName:   "TypeSpec Compiler Installation",
		},
		{
			name:           "Installation failed",
			commandOutput:  "",
			commandErr:     &mockError{msg: "npm install failed"},
			expectedStatus: "ERROR",
			expectedName:   "TypeSpec Compiler Installation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original runCommand
			originalRunCommand := runCommand
			defer func() {
				runCommand = originalRunCommand
			}()

			// Override runCommand for this test
			runCommand = func(command string, args ...string) (string, error) {
				return tt.commandOutput, tt.commandErr
			}

			result := installTypeSpecCompiler()

			if result.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, result.Name)
			}

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}
		})
	}
}

func TestInstallTypeSpecClientCLI(t *testing.T) {
	tests := []struct {
		name           string
		commandOutput  string
		commandErr     error
		expectedStatus string
		expectedName   string
	}{
		{
			name:           "Installation successful",
			commandOutput:  "installed successfully",
			commandErr:     nil,
			expectedStatus: "SUCCESS",
			expectedName:   "TypeSpec Client Generator CLI Installation",
		},
		{
			name:           "Installation failed",
			commandOutput:  "",
			commandErr:     &mockError{msg: "npm install failed"},
			expectedStatus: "ERROR",
			expectedName:   "TypeSpec Client Generator CLI Installation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original runCommand
			originalRunCommand := runCommand
			defer func() {
				runCommand = originalRunCommand
			}()

			// Override runCommand for this test
			runCommand = func(command string, args ...string) (string, error) {
				return tt.commandOutput, tt.commandErr
			}

			result := installTypeSpecClientCLI()

			if result.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, result.Name)
			}

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}
		})
	}
}

func TestCheckGitHubCLI(t *testing.T) {
	tests := []struct {
		name            string
		commandOutput   string
		commandErr      error
		expectedStatus  string
		expectedName    string
		expectedVersion string
	}{
		{
			name:            "GitHub CLI installed",
			commandOutput:   "gh version 2.40.1 (2023-12-13)\nhttps://github.com/cli/cli/releases/tag/v2.40.1",
			commandErr:      nil,
			expectedStatus:  "SUCCESS",
			expectedName:    "GitHub CLI",
			expectedVersion: "2.40.1",
		},
		{
			name:            "GitHub CLI not installed",
			commandOutput:   "",
			commandErr:      &mockError{msg: "command not found"},
			expectedStatus:  "ERROR",
			expectedName:    "GitHub CLI",
			expectedVersion: "",
		},
		{
			name:            "GitHub CLI version format unexpected",
			commandOutput:   "gh version unknown",
			commandErr:      nil,
			expectedStatus:  "SUCCESS",
			expectedName:    "GitHub CLI",
			expectedVersion: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original runCommand
			originalRunCommand := runCommand
			defer func() {
				runCommand = originalRunCommand
			}()

			// Override runCommand for this test
			runCommand = func(command string, args ...string) (string, error) {
				return tt.commandOutput, tt.commandErr
			}

			result := checkGitHubCLI()

			if result.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, result.Name)
			}

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}

			if result.Version != tt.expectedVersion {
				t.Errorf("Expected version %s, got %s", tt.expectedVersion, result.Version)
			}
		})
	}
}

func TestCheckGitHubCLIAuth(t *testing.T) {
	tests := []struct {
		name           string
		commandOutput  string
		commandErr     error
		expectedStatus string
		expectedName   string
	}{
		{
			name:           "GitHub CLI authenticated",
			commandOutput:  "github.com\n  ✓ Logged in to github.com as user (oauth_token)",
			commandErr:     nil,
			expectedStatus: "SUCCESS",
			expectedName:   "GitHub CLI Authentication",
		},
		{
			name:           "GitHub CLI not authenticated",
			commandOutput:  "",
			commandErr:     &mockError{msg: "not authenticated"},
			expectedStatus: "ERROR",
			expectedName:   "GitHub CLI Authentication",
		},
		{
			name:           "GitHub CLI auth status unclear",
			commandOutput:  "some unexpected output",
			commandErr:     nil,
			expectedStatus: "ERROR",
			expectedName:   "GitHub CLI Authentication",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Save original runCommand
			originalRunCommand := runCommand
			defer func() {
				runCommand = originalRunCommand
			}()

			// Override runCommand for this test
			runCommand = func(command string, args ...string) (string, error) {
				return tt.commandOutput, tt.commandErr
			}

			result := checkGitHubCLIAuth()

			if result.Name != tt.expectedName {
				t.Errorf("Expected name %s, got %s", tt.expectedName, result.Name)
			}

			if result.Status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, result.Status)
			}
		})
	}
}

func TestGetStatusSymbol(t *testing.T) {
	tests := []struct {
		status   string
		expected string
	}{
		{"SUCCESS", "✓"},
		{"WARNING", "⚠"},
		{"ERROR", "✗"},
		{"UNKNOWN", "?"},
	}

	for _, tt := range tests {
		t.Run(tt.status, func(t *testing.T) {
			result := getStatusSymbol(tt.status)
			if result != tt.expected {
				t.Errorf("getStatusSymbol(%s) = %s, expected %s", tt.status, result, tt.expected)
			}
		})
	}
}

// mockError implements the error interface for testing
type mockError struct {
	msg string
}

func (e *mockError) Error() string {
	return e.msg
}
