// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package metadata

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	"github.com/spf13/cobra"
)

// MetadataResult represents the result of a metadata operation
type MetadataResult struct {
	Success      bool     `json:"success"`
	Message      string   `json:"message"`
	PackagePath  string   `json:"package_path,omitempty"`
	CreatedFiles []string `json:"created_files,omitempty"`
	SkippedFiles []string `json:"skipped_files,omitempty"`
}

// Command returns the metadata command
func Command() *cobra.Command {
	var outputFormat string
	var verbose bool
	var packageTitle string

	metadataCmd := &cobra.Command{
		Use:   "metadata <package-path>",
		Short: "Create required metadata files for a package",
		Long: `Create required metadata files (ci.yml, README.md) for a package if they don't already exist.

This command will:
1. Create ci.yml if it doesn't exist in the package directory (for both data plane and mgmt plane packages)
2. Create README.md if it doesn't exist in the package directory (for mgmt plane packages only)

The package path should be an absolute path to a Go module (containing go.mod file).

Examples:
  # Create metadata files for an mgmt plane package
  generator metadata /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute

  # Create metadata files for a data plane package
  generator metadata /path/to/azure-sdk-for-go/sdk/messaging/azeventhubs

  # Create metadata files with a custom package title (used in README.md for mgmt plane)
  generator metadata /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --package-title "Compute"

  # Create metadata files with JSON output
  generator metadata /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute --output json`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			packagePath := args[0]

			// Validate package path
			if err := utils.ValidatePackagePath(packagePath); err != nil {
				return fmt.Errorf("package path validation error: %v", err)
			}

			// Process metadata creation
			result, err := processMetadata(packagePath, packageTitle, verbose)
			if err != nil {
				return fmt.Errorf("metadata operation failed: %v", err)
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
					fmt.Printf("✓ Metadata files created successfully!\n\n")
					fmt.Printf("Package: %s\n", result.PackagePath)
					if len(result.CreatedFiles) > 0 {
						fmt.Printf("Created files:\n")
						for _, f := range result.CreatedFiles {
							fmt.Printf("  - %s\n", f)
						}
					}
					if len(result.SkippedFiles) > 0 {
						fmt.Printf("Skipped files (already exist):\n")
						for _, f := range result.SkippedFiles {
							fmt.Printf("  - %s\n", f)
						}
					}
				} else {
					fmt.Printf("✗ Metadata creation failed: %s\n", result.Message)
					return fmt.Errorf("metadata creation failed")
				}
			}

			return nil
		},
	}

	metadataCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	metadataCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	metadataCmd.Flags().StringVar(&packageTitle, "package-title", "", "Package title used in README.md (defaults to inferred from package name)")

	return metadataCmd
}

// inferPackageTitle infers the package title from the package name
// For example: "armcompute" -> "Compute", "azblob" -> "Blob"
func inferPackageTitle(packageName string) string {
	title := packageName
	// Remove "arm" prefix for management packages
	title = strings.TrimPrefix(title, "arm")
	// Remove "az" prefix for data plane packages
	title = strings.TrimPrefix(title, "az")
	// Capitalize first letter
	if len(title) > 0 {
		title = strings.ToUpper(title[:1]) + title[1:]
	}
	return title
}

// processMetadata creates the metadata files for the given package
func processMetadata(packagePath, packageTitle string, verbose bool) (*MetadataResult, error) {
	result := &MetadataResult{
		PackagePath:  packagePath,
		CreatedFiles: []string{},
		SkippedFiles: []string{},
	}

	// Find SDK root
	sdkRoot, err := utils.FindSDKRoot(packagePath)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Failed to find SDK root: %v", err)
		return result, nil
	}

	// Get module relative path
	moduleRelativePath, err := utils.GetModuleRelativePath(packagePath, sdkRoot)
	if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Failed to get module relative path: %v", err)
		return result, nil
	}

	serviceDir := utils.GetServiceDir(moduleRelativePath)
	packageName := filepath.Base(packagePath)
	isMgmt := utils.IsMgmtPlanePackage(packagePath, sdkRoot)

	if verbose {
		log.Printf("Package path: %s", packagePath)
		log.Printf("SDK root: %s", sdkRoot)
		log.Printf("Module relative path: %s", moduleRelativePath)
		log.Printf("Service directory: %s", serviceDir)
		log.Printf("Package name: %s", packageName)
		log.Printf("Is mgmt plane: %v", isMgmt)
	}

	// Infer package title if not provided
	if packageTitle == "" {
		packageTitle = inferPackageTitle(packageName)
	}

	// Create template data
	data := map[string]string{
		"moduleRelativePath": moduleRelativePath,
		"serviceDir":         serviceDir,
		"packageName":        packageName,
		"packageTitle":       packageTitle,
	}

	templateDir := filepath.Join(sdkRoot, "eng", "tools", "generator", "template", "typespec")

	// Create ci.yml if not exists
	ciYmlPath := filepath.Join(packagePath, "ci.yml")
	if _, err := os.Stat(ciYmlPath); os.IsNotExist(err) {
		if verbose {
			log.Printf("Creating ci.yml at %s", ciYmlPath)
		}
		ciTplPath := filepath.Join(templateDir, "ci.yml.tpl")
		if err := typespec.ParseSingleTemplate(ciTplPath, ciYmlPath, data); err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Failed to create ci.yml: %v", err)
			return result, nil
		}
		result.CreatedFiles = append(result.CreatedFiles, "ci.yml")
	} else if err != nil {
		result.Success = false
		result.Message = fmt.Sprintf("Failed to check ci.yml: %v", err)
		return result, nil
	} else {
		result.SkippedFiles = append(result.SkippedFiles, "ci.yml")
		if verbose {
			log.Printf("ci.yml already exists at %s, skipping", ciYmlPath)
		}
	}

	// Create README.md if not exists (only for mgmt plane packages)
	if isMgmt {
		readmePath := filepath.Join(packagePath, "README.md")
		if _, err := os.Stat(readmePath); os.IsNotExist(err) {
			if verbose {
				log.Printf("Creating README.md at %s", readmePath)
			}
			readmeTplPath := filepath.Join(templateDir, "README.md.tpl")
			if err := typespec.ParseSingleTemplate(readmeTplPath, readmePath, data); err != nil {
				result.Success = false
				result.Message = fmt.Sprintf("Failed to create README.md: %v", err)
				return result, nil
			}
			result.CreatedFiles = append(result.CreatedFiles, "README.md")
		} else if err != nil {
			result.Success = false
			result.Message = fmt.Sprintf("Failed to check README.md: %v", err)
			return result, nil
		} else {
			result.SkippedFiles = append(result.SkippedFiles, "README.md")
			if verbose {
				log.Printf("README.md already exists at %s, skipping", readmePath)
			}
		}
	}

	result.Success = true
	result.Message = "Metadata files processed successfully"
	return result, nil
}
