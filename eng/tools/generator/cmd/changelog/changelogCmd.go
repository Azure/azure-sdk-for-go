// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package changelog

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/version"
	"github.com/spf13/cobra"
)

// ChangelogResult represents the result of a changelog operation
type ChangelogResult struct {
	Success           bool   `json:"success"`
	Message           string `json:"message"`
	PackagePath       string `json:"package_path,omitempty"`
	PackageStatus     string `json:"package_status,omitempty"`
	Version           string `json:"version,omitempty"`
	ReleaseDate       string `json:"release_date,omitempty"`
	HasBreakingChange bool   `json:"hasBreakingChange"`
	ChangelogMD       string `json:"changelog_md,omitempty"`
}

// Command returns the changelog command
func Command() *cobra.Command {
	var outputFormat string
	var verbose bool
	var reportFile string

	changelogCmd := &cobra.Command{
		Use:   "changelog <package-path>",
		Short: "Update changelog content for a package",
		Long: `Update changelog content for a package by determining the package status and generating appropriate changelog entries.

This command will:
1. Determine the package status (new package, existing package with new preview version, existing package with new stable version)
2. For new packages: generate changelog according to the template
3. For existing packages: compare current package exports with previous released version and calculate the changelog
4. Put the changelog entry into CHANGELOG.md file, replacing existing version entry if it exists

The package path should be an absolute path to a Go module (containing go.mod file).

When --report-file is provided, the command runs in report-only mode: it computes the
SDK changes (including breaking change detection) and writes a JSON report to the given
file path. CHANGELOG.md is NOT modified in this mode.

Examples:
  # Update changelog for an existing package
  generator changelog /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute

  # Update changelog with verbose output
  generator changelog /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --verbose

  # Generate changelog with JSON output
  generator changelog /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --output json

  # Report-only mode: compute SDK changes and write JSON report without updating CHANGELOG.md
  generator changelog /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --report-file /path/to/sdkchange.json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packagePath := args[0]

			// Validate package path
			if err := utils.ValidatePackagePath(packagePath); err != nil {
				return fmt.Errorf("package path validation error: %v", err)
			}

			// Get SDK root path
			sdkRoot, err := utils.FindSDKRoot(packagePath)
			if err != nil {
				return fmt.Errorf("failed to find SDK root: %v", err)
			}

			reportOnly := reportFile != ""

			// Process changelog
			result, err := processChangelog(sdkRoot, packagePath, verbose, reportOnly)
			if err != nil {
				return fmt.Errorf("changelog operation failed: %v", err)
			}

			// In report-only mode, write JSON report to file and emit JSON to stdout.
			if reportOnly {
				jsonResult, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal result: %v", err)
				}
				fmt.Println(string(jsonResult))
				if err := os.WriteFile(reportFile, jsonResult, 0644); err != nil {
					return fmt.Errorf("failed to write report file: %v", err)
				}
				return nil
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
					fmt.Printf("✓ Changelog updated successfully!\n\n")
					fmt.Printf("Package: %s\n", result.PackagePath)
					fmt.Printf("Status: %s\n", result.PackageStatus)
					if result.Version != "" {
						fmt.Printf("Version: %s\n", result.Version)
					}
					fmt.Printf("Release Date: %s\n", result.ReleaseDate)
					if verbose && result.ChangelogMD != "" {
						fmt.Printf("\nGenerated Changelog:\n%s\n", result.ChangelogMD)
					}
				} else {
					fmt.Printf("✗ Changelog update failed: %s\n", result.Message)
					return fmt.Errorf("changelog update failed")
				}
			}

			return nil
		},
	}

	changelogCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	changelogCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	changelogCmd.Flags().StringVar(&reportFile, "report-file", "", "Write SDK change report as JSON to the given file path. When set, CHANGELOG.md is not modified.")

	return changelogCmd
}

// processChangelog processes the changelog for the given package.
// When reportOnly is true, CHANGELOG.md is not modified; only the change report
// (HasBreakingChange, ChangelogMD) is populated.
func processChangelog(sdkRoot, packagePath string, verbose, reportOnly bool) (*ChangelogResult, error) {
	result := &ChangelogResult{
		PackagePath: packagePath,
		ReleaseDate: time.Now().Format("2006-01-02"),
	}

	// Initialize SDK repo
	sdkRepo, err := common.GetSDKRepo(sdkRoot, "")
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Failed to initialize SDK repository: %v", err)
		return result, nil
	}

	if verbose {
		log.Printf("Processing package at: %s", packagePath)
	}

	// Determine package status
	status, err := changelog.DetermineModuleStatus(packagePath, sdkRepo)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Failed to determine package status: %v", err)
		return result, nil
	}

	switch status {
	case utils.PackageStatusNew:
		result.PackageStatus = "New Package"
		err = handleNewPackage(packagePath, sdkRepo, result, verbose, reportOnly)
	case utils.PackageStatusExisting:
		result.PackageStatus = "Existing Package"
		err = handleExistingPackage(packagePath, sdkRepo, result, verbose, reportOnly)
	default:
		result.Success = false
		result.Message = fmt.Sprintf("Unknown package status: %v", status)
		return result, nil
	}

	if err != nil {
		result.Success = false
		result.Message = err.Error()
		return result, nil
	}

	result.Success = true
	result.Message = "Changelog updated successfully"
	return result, nil
}

func handleNewPackage(modulePath string, sdkRepo repo.SDKRepository, result *ChangelogResult, verbose, reportOnly bool) error {
	if verbose {
		log.Printf("Handling new package...")
	}

	result.Version = "0.1.0"
	result.ChangelogMD = "New package"
	result.HasBreakingChange = false

	if reportOnly {
		return nil
	}

	if err := changelog.CreateNewChangelog(modulePath, sdkRepo, result.Version, result.ReleaseDate); err != nil {
		return fmt.Errorf("failed to create new package changelog: %v", err)
	}

	if verbose {
		log.Printf("Created new package changelog for %s", modulePath)
	}

	return nil
}

func handleExistingPackage(modulePath string, sdkRepo repo.SDKRepository, result *ChangelogResult, verbose, reportOnly bool) error {
	if verbose {
		log.Printf("Handling existing package...")
	}

	// Generate changelog using the new GenerateChangelog function
	if verbose {
		log.Printf("Generating changelog...")
	}

	// Generate changelog
	isCurrentPreview, err := version.IsCurrentPreviewVersion(modulePath, sdkRepo, nil)
	if err != nil {
		return fmt.Errorf("failed to determine if current version is preview: %v", err)
	}
	changelogResult, err := changelog.GenerateChangelog(modulePath, sdkRepo, isCurrentPreview)
	if err != nil {
		return fmt.Errorf("failed to generate changelog: %v", err)
	}

	result.HasBreakingChange = changelogResult.ChangelogData.HasBreakingChanges()

	if reportOnly {
		// In report-only mode, do not modify CHANGELOG.md or compute a new version.
		result.ChangelogMD = changelogResult.ChangelogData.ToCompactMarkdown()
		return nil
	}

	// Calculate new version
	newVersion, _, err := version.CalculateNewVersion(changelogResult.ChangelogData, changelogResult.PreviousVersion, isCurrentPreview)
	if err != nil {
		return fmt.Errorf("failed to calculate new version: %v", err)
	}

	if verbose {
		log.Printf("Calculated new version: %s", newVersion.String())
	}

	result.Version = newVersion.String()

	// Update changelog file
	changelogMd, err := changelog.AddChangelogToFileWithReplacement(changelogResult.ChangelogData, newVersion, modulePath, result.ReleaseDate)
	if err != nil {
		return fmt.Errorf("failed to update changelog file: %v", err)
	}

	result.ChangelogMD = changelogMd

	if verbose {
		log.Printf("Successfully updated changelog for version %s", newVersion.String())
	}

	return nil
}
