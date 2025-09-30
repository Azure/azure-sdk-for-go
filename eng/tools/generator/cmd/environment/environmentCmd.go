// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package environment

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// Configuration constants
const (
	MinNodeVersion = "20.0.0"
)

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

// Command returns the environment command
func Command() *cobra.Command {
	var outputFormat string
	var autoInstall bool = true

	envCmd := &cobra.Command{
		Use:   "environment",
		Short: "Check and validate environment prerequisites for Azure Go SDK generation",
		Long: `This command checks and validates environment prerequisites for Azure Go SDK generation.
It verifies the installation and versions of required tools including:
- Node.js (minimum version 20.0.0)
- TypeSpec compiler
- TypeSpec client generator CLI
- GitHub CLI and authentication

The command automatically installs missing TypeSpec tools by default. Use --auto-install=false to disable this behavior.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			result := &EnvironmentCheckResult{
				Success:   true,
				Checks:    []EnvironmentChecker{},
				Failed:    []string{},
				Installed: []string{},
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
			if (tspCheck.Status == "ERROR" || tspCheck.Status == "WARNING") && autoInstall {
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
			} else if tspCheck.Status == "ERROR" || tspCheck.Status == "WARNING" {
				result.Success = false
				result.Failed = append(result.Failed, "TypeSpec compiler")
			}

			// Check TypeSpec client generator CLI
			tspClientCheck := checkTypeSpecClientCLI()
			result.Checks = append(result.Checks, tspClientCheck)
			if (tspClientCheck.Status == "ERROR" || tspClientCheck.Status == "WARNING") && autoInstall {
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
			} else if tspClientCheck.Status == "ERROR" || tspClientCheck.Status == "WARNING" {
				result.Success = false
				result.Failed = append(result.Failed, "TypeSpec client generator CLI")
			}

			// Check GitHub CLI installation
			ghCheck := checkGitHubCLI()
			result.Checks = append(result.Checks, ghCheck)
			if ghCheck.Status == "ERROR" {
				result.Success = false
				result.Failed = append(result.Failed, "GitHub CLI (gh)")
			}

			// Check GitHub CLI authentication
			ghAuthCheck := checkGitHubCLIAuth()
			result.Checks = append(result.Checks, ghAuthCheck)
			if ghAuthCheck.Status == "ERROR" {
				result.Success = false
				result.Failed = append(result.Failed, "GitHub CLI authentication")
			}

			// Generate summary
			if result.Success {
				result.Summary = "All environment checks are satisfied! ✓"
			} else {
				result.Summary = fmt.Sprintf("Missing environment: %s", strings.Join(result.Failed, ", "))
			}

			// Output result
			switch outputFormat {
			case "json":
				jsonResult, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal result: %v", err)
				}
				fmt.Println(string(jsonResult))
			default:
				// Human-readable output
				fmt.Println(result.Summary)
				fmt.Println()
				for _, check := range result.Checks {
					status := getStatusSymbol(check.Status)
					fmt.Printf("%s %s: %s\n", status, check.Name, check.Message)
					if check.InstallHint != "" && check.Status != "SUCCESS" {
						fmt.Printf("   Hint: %s\n", check.InstallHint)
					}
				}
				if len(result.Installed) > 0 {
					fmt.Printf("\n✓ Automatically installed: %s\n", strings.Join(result.Installed, ", "))
				}
			}

			// Return exit code 1 if checks failed
			if !result.Success {
				return fmt.Errorf("environment checks failed")
			}

			return nil
		},
	}

	envCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	envCmd.Flags().BoolVar(&autoInstall, "auto-install", true, "Automatically install missing TypeSpec tools")

	return envCmd
}

func getStatusSymbol(status string) string {
	switch status {
	case "SUCCESS":
		return "✓"
	case "WARNING":
		return "⚠"
	case "ERROR":
		return "✗"
	default:
		return "?"
	}
}

// Helper functions

// runCommand executes a shell command and returns the output
var runCommand = func(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	return string(output), err
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
			InstallHint: "Will be installed automatically, or run: npm install -g @typespec/compiler",
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
			InstallHint: "Will be installed automatically, or run: npm install -g @azure-tools/typespec-client-generator-cli",
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

// checkGitHubCLI checks if GitHub CLI is installed and available
func checkGitHubCLI() EnvironmentChecker {
	output, err := runCommand("gh", "--version")
	if err != nil {
		return EnvironmentChecker{
			Name:        "GitHub CLI",
			Status:      "ERROR",
			Message:     "GitHub CLI (gh) is not installed or not in PATH",
			InstallHint: "Install from https://cli.github.com/ or run: winget install GitHub.CLI",
		}
	}

	// Extract version from output (e.g., "gh version 2.40.1 (2023-12-13)")
	versionPattern := regexp.MustCompile(`gh version (\d+\.\d+\.\d+)`)
	matches := versionPattern.FindStringSubmatch(output)

	var version string
	if len(matches) > 1 {
		version = matches[1]
	} else {
		version = "unknown"
	}

	return EnvironmentChecker{
		Name:    "GitHub CLI",
		Status:  "SUCCESS",
		Version: version,
		Message: fmt.Sprintf("GitHub CLI %s is installed ✓", version),
	}
}

// checkGitHubCLIAuth checks if GitHub CLI is authenticated
func checkGitHubCLIAuth() EnvironmentChecker {
	output, err := runCommand("gh", "auth", "status")
	if err != nil {
		return EnvironmentChecker{
			Name:        "GitHub CLI Authentication",
			Status:      "ERROR",
			Message:     "GitHub CLI is not authenticated",
			InstallHint: "Run: gh auth login",
		}
	}

	// Check if output indicates successful authentication
	if strings.Contains(output, "Logged in to github.com") || strings.Contains(output, "✓") {
		return EnvironmentChecker{
			Name:    "GitHub CLI Authentication",
			Status:  "SUCCESS",
			Message: "GitHub CLI is authenticated ✓",
		}
	}

	return EnvironmentChecker{
		Name:        "GitHub CLI Authentication",
		Status:      "ERROR",
		Message:     "GitHub CLI authentication status unclear",
		InstallHint: "Run: gh auth login",
	}
}
