// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
)

// Configuration constants
const (
	MinGoVersion   = "1.23"
	MinNodeVersion = "20.0.0"
)

// EnvironmentCheckerTool creates and returns the check-sdk-generation-environment tool
func EnvironmentCheckerTool() mcp.Tool {
	return mcp.NewTool("check-sdk-generation-environment",
		mcp.WithDescription("Checks and validates environment prerequisites for Azure Go SDK generation. Automatically installs missing tools."),
	)
}

// EnvironmentChecker represents the result of a single environment check
type EnvironmentChecker struct {
	Name        string `json:"name"`
	Status      string `json:"status"` // "SUCCESS", "ERROR", "WARNING"
	Version     string `json:"version,omitempty"`
	Message     string `json:"message"`
	InstallHint string `json:"install_hint,omitempty"`
}

// EnvironmentCheckResult represents the overall result of all checks
type EnvironmentCheckResult struct {
	Success   bool                 `json:"success"`
	Summary   string               `json:"summary"`
	Checks    []EnvironmentChecker `json:"checks"`
	Failed    []string             `json:"failed,omitempty"`
	Installed []string             `json:"installed,omitempty"`
}

// EnvironmentCheckerHandler handles the check-sdk-generation-environment tool requests
func EnvironmentCheckerHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	result := &EnvironmentCheckResult{
		Success:   true,
		Checks:    []EnvironmentChecker{},
		Failed:    []string{},
		Installed: []string{},
	}

	// Check Go version
	goCheck := checkGoVersion(MinGoVersion)
	result.Checks = append(result.Checks, goCheck)
	if goCheck.Status == "ERROR" {
		result.Success = false
		result.Failed = append(result.Failed, fmt.Sprintf("Go %s or later", MinGoVersion))
	}

	// Check Node.js version
	nodeCheck := checkNodeVersion(MinNodeVersion)
	result.Checks = append(result.Checks, nodeCheck)
	if nodeCheck.Status == "ERROR" {
		result.Success = false
		result.Failed = append(result.Failed, fmt.Sprintf("Node.js %s or later", MinNodeVersion))
	}

	// Check TypeSpec compiler
	tspCheck := checkTypeSpecCompiler()
	result.Checks = append(result.Checks, tspCheck)
	if tspCheck.Status == "ERROR" || tspCheck.Status == "WARNING" {
		installResult := installTypeSpecCompiler()
		result.Checks = append(result.Checks, installResult)
		if installResult.Status == "SUCCESS" {
			result.Installed = append(result.Installed, "TypeSpec compiler")
			// Re-check after installation
			tspRecheck := checkTypeSpecCompiler()
			result.Checks = append(result.Checks, tspRecheck)
			if tspRecheck.Status == "ERROR" {
				result.Success = false
				result.Failed = append(result.Failed, "TypeSpec compiler (failed to install properly)")
			}
		} else {
			result.Success = false
			result.Failed = append(result.Failed, "TypeSpec compiler (failed to install)")
		}
	}

	// Check TypeSpec client generator CLI
	tspClientCheck := checkTypeSpecClientCLI()
	result.Checks = append(result.Checks, tspClientCheck)
	if tspClientCheck.Status == "ERROR" || tspClientCheck.Status == "WARNING" {
		installResult := installTypeSpecClientCLI()
		result.Checks = append(result.Checks, installResult)
		if installResult.Status == "SUCCESS" {
			result.Installed = append(result.Installed, "TypeSpec client generator CLI")
			// Re-check after installation
			tspClientRecheck := checkTypeSpecClientCLI()
			result.Checks = append(result.Checks, tspClientRecheck)
			if tspClientRecheck.Status == "ERROR" {
				result.Success = false
				result.Failed = append(result.Failed, "TypeSpec client generator CLI (failed to install properly)")
			}
		} else {
			result.Success = false
			result.Failed = append(result.Failed, "TypeSpec client generator CLI (failed to install)")
		}
	}

	// Generate summary
	if result.Success {
		result.Summary = "All environment checks are satisfied! ✓"
	} else {
		result.Summary = fmt.Sprintf("Missing environment: %s", strings.Join(result.Failed, ", "))
	}

	// Return result as JSON
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to marshal result: %v", err)), nil
	}

	return mcp.NewToolResultText(string(jsonResult)), nil
}

// Helper functions

// runCommand is a variable that can be overridden for testing
var runCommand = func(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

func checkGoVersion(minVersion string) EnvironmentChecker {
	output, err := runCommand("go", "version")
	if err != nil {
		return EnvironmentChecker{
			Name:        "Go",
			Status:      "ERROR",
			Message:     "Go is not installed or not in PATH",
			InstallHint: "Install from https://golang.org/dl/",
		}
	}

	// Extract version from "go version go1.23.0 ..."
	re := regexp.MustCompile(`go version go(\d+\.\d+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return EnvironmentChecker{
			Name:    "Go",
			Status:  "ERROR",
			Message: "Could not parse Go version",
		}
	}

	version := matches[1]
	if compareVersions(version, minVersion) {
		return EnvironmentChecker{
			Name:    "Go",
			Status:  "SUCCESS",
			Version: version,
			Message: fmt.Sprintf("Go version %s is installed ✓", version),
		}
	}

	return EnvironmentChecker{
		Name:        "Go",
		Status:      "ERROR",
		Version:     version,
		Message:     fmt.Sprintf("Go version %s is too old. Minimum required: %s", version, minVersion),
		InstallHint: "Update from https://golang.org/dl/",
	}
}

func checkNodeVersion(minVersion string) EnvironmentChecker {
	output, err := runCommand("node", "--version")
	if err != nil {
		return EnvironmentChecker{
			Name:        "Node.js",
			Status:      "ERROR",
			Message:     "Node.js is not installed or not in PATH",
			InstallHint: "Install from https://nodejs.org/",
		}
	}

	// Extract version from "v20.0.0"
	re := regexp.MustCompile(`v(\d+\.\d+\.\d+)`)
	matches := re.FindStringSubmatch(strings.TrimSpace(output))
	if len(matches) < 2 {
		return EnvironmentChecker{
			Name:    "Node.js",
			Status:  "ERROR",
			Message: "Could not parse Node.js version",
		}
	}

	version := matches[1]
	if compareVersions(version, minVersion) {
		return EnvironmentChecker{
			Name:    "Node.js",
			Status:  "SUCCESS",
			Version: version,
			Message: fmt.Sprintf("Node.js version %s is installed ✓", version),
		}
	}

	return EnvironmentChecker{
		Name:        "Node.js",
		Status:      "ERROR",
		Version:     version,
		Message:     fmt.Sprintf("Node.js version %s is too old. Minimum required: %s", version, minVersion),
		InstallHint: "Update from https://nodejs.org/",
	}
}

func checkTypeSpecCompiler() EnvironmentChecker {
	_, err := runCommand("tsp", "--version")
	if err != nil {
		return EnvironmentChecker{
			Name:        "TypeSpec Compiler",
			Status:      "WARNING",
			Message:     "TypeSpec compiler is not installed",
			InstallHint: "Will be installed automatically",
		}
	}

	return EnvironmentChecker{
		Name:    "TypeSpec Compiler",
		Status:  "SUCCESS",
		Message: "TypeSpec compiler is installed ✓",
	}
}

func checkTypeSpecClientCLI() EnvironmentChecker {
	_, err := runCommand("tsp-client", "--version")
	if err != nil {
		return EnvironmentChecker{
			Name:        "TypeSpec Client Generator CLI",
			Status:      "WARNING",
			Message:     "TypeSpec client generator CLI is not installed",
			InstallHint: "Will be installed automatically",
		}
	}

	return EnvironmentChecker{
		Name:    "TypeSpec Client Generator CLI",
		Status:  "SUCCESS",
		Message: "TypeSpec client generator CLI is installed ✓",
	}
}

func installTypeSpecCompiler() EnvironmentChecker {
	_, err := runCommand("npm", "install", "-g", "@typespec/compiler")
	if err != nil {
		return EnvironmentChecker{
			Name:    "TypeSpec Compiler Installation",
			Status:  "ERROR",
			Message: fmt.Sprintf("Failed to install TypeSpec compiler: %v", err),
		}
	}

	return EnvironmentChecker{
		Name:    "TypeSpec Compiler Installation",
		Status:  "SUCCESS",
		Message: "TypeSpec compiler installed successfully ✓",
	}
}

func installTypeSpecClientCLI() EnvironmentChecker {
	_, err := runCommand("npm", "install", "-g", "@azure-tools/typespec-client-generator-cli")
	if err != nil {
		return EnvironmentChecker{
			Name:    "TypeSpec Client Generator CLI Installation",
			Status:  "ERROR",
			Message: fmt.Sprintf("Failed to install TypeSpec client generator CLI: %v", err),
		}
	}

	return EnvironmentChecker{
		Name:    "TypeSpec Client Generator CLI Installation",
		Status:  "SUCCESS",
		Message: "TypeSpec client generator CLI installed successfully ✓",
	}
}

func compareVersions(current, minimum string) bool {
	currentParts := strings.Split(current, ".")
	minimumParts := strings.Split(minimum, ".")

	// Ensure both have at least 2 parts (major.minor)
	for len(currentParts) < 2 {
		currentParts = append(currentParts, "0")
	}
	for len(minimumParts) < 2 {
		minimumParts = append(minimumParts, "0")
	}

	for i := 0; i < len(minimumParts) && i < len(currentParts); i++ {
		currentNum, err1 := strconv.Atoi(currentParts[i])
		minimumNum, err2 := strconv.Atoi(minimumParts[i])

		if err1 != nil || err2 != nil {
			// Fallback to string comparison
			return current >= minimum
		}

		if currentNum > minimumNum {
			return true
		} else if currentNum < minimumNum {
			return false
		}
		// Continue to next part if equal
	}

	return true // All parts are equal or current has more parts
}
