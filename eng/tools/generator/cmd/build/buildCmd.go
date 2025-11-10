// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package build

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

// BuildResult represents the result of a build operation
type BuildResult struct {
	Success     bool   `json:"success"`
	Message     string `json:"message"`
	Path        string `json:"path,omitempty"`
	BuildOutput string `json:"build_output,omitempty"`
	VetOutput   string `json:"vet_output,omitempty"`
}

// Command returns the build command
func Command() *cobra.Command {
	var outputFormat string
	var verbose bool

	buildCmd := &cobra.Command{
		Use:   "build <folder-path>",
		Short: "Build and vet Go packages in the specified folder",
		Long: `Build and vet Go packages in the specified folder using 'go build' and 'go vet'.

This command will:
1. Run 'go build' to compile the Go packages in the specified folder
2. Run 'go vet' to check for common Go programming errors
3. Report any issues found during build or vet process

The command recursively processes all Go packages found in the specified folder.

Examples:
  # Build and vet packages in a specific folder
  generator build /path/to/generated/sdk/folder

  # Build with verbose output
  generator build /path/to/generated/sdk/folder --verbose

  # Build with JSON output
  generator build /path/to/generated/sdk/folder --output json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			folderPath := args[0]

			// Validate that the path exists and is a directory
			if err := validatePath(folderPath); err != nil {
				return fmt.Errorf("path validation error: %v", err)
			}

			// Perform build and vet
			result, err := buildAndVet(folderPath)
			if err != nil {
				return fmt.Errorf("build operation failed: %v", err)
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
				if result.Success {
					fmt.Printf("✓ Build and vet completed successfully!\n\n")
					fmt.Printf("Path: %s\n", result.Path)
					if verbose && result.BuildOutput != "" {
						fmt.Printf("\nBuild Output:\n%s\n", result.BuildOutput)
					}
					if verbose && result.VetOutput != "" {
						fmt.Printf("\nVet Output:\n%s\n", result.VetOutput)
					}
				} else {
					fmt.Printf("✗ Build and vet failed: %s\n", result.Message)
					if result.BuildOutput != "" {
						fmt.Printf("\nBuild Output:\n%s\n", result.BuildOutput)
					}
					if result.VetOutput != "" {
						fmt.Printf("\nVet Output:\n%s\n", result.VetOutput)
					}
					return fmt.Errorf("build or vet failed")
				}
			}

			return nil
		},
	}

	buildCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	buildCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	return buildCmd
}

// validatePath validates that the provided path exists and is a directory
func validatePath(path string) error {
	if info, err := os.Stat(path); err != nil || !info.IsDir() {
		return fmt.Errorf("path '%s' does not exist or is not a directory", path)
	}
	return nil
}

// buildAndVet performs go build and go vet operations on the specified folder
func buildAndVet(folderPath string) (*BuildResult, error) {
	result := &BuildResult{
		Path: folderPath,
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(folderPath)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Failed to get absolute path: %v", err)
		return result, nil
	}
	result.Path = absPath

	// Run go build
	buildSuccess, buildOutput := runGoBuild(absPath)
	result.BuildOutput = buildOutput

	// Run go vet
	vetSuccess, vetOutput := runGoVet(absPath)
	result.VetOutput = vetOutput

	// Determine overall success
	result.Success = buildSuccess && vetSuccess

	if !result.Success {
		if !buildSuccess && !vetSuccess {
			result.Message = "Both build and vet failed"
		} else if !buildSuccess {
			result.Message = "Build failed"
		} else {
			result.Message = "Vet failed"
		}
	} else {
		result.Message = "Build and vet completed successfully"
	}

	return result, nil
}

// runGoBuild executes go build and returns success status and output
func runGoBuild(path string) (bool, string) {
	cmd := exec.Command("go", "build", "./...")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		return false, outputStr
	}

	return true, outputStr
}

// runGoVet executes go vet and returns success status and output
func runGoVet(path string) (bool, string) {
	cmd := exec.Command("go", "vet", "./...")
	cmd.Dir = path

	output, err := cmd.CombinedOutput()
	outputStr := string(output)

	if err != nil {
		return false, outputStr
	}

	return true, outputStr
}
