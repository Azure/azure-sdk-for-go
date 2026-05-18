package sdkchange

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/v2/common"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/version"
	"github.com/spf13/cobra"
)

// ChangelogResult represents the result of a changelog operation
type SdkChangeResult struct {
	Success           bool   `json:"success"`
	Message           string `json:"message"`
	PackagePath       string `json:"package_path,omitempty"`
	PackageStatus     string `json:"package_status,omitempty"`
	HasBreakingChange bool   `json:"hasBreakingChange"`
	ChangelogMD       string `json:"changelog_md,omitempty"`
}

// Command returns the changelog command
func Command() *cobra.Command {
	var outputFormat string
	var verbose bool

	changelogCmd := &cobra.Command{
		Use:   "sdkchange <package-path> <outputjson-file-path>",
		Short: "Get sdk changes for a existing package",
		Long: `Get sdk changes for a existing package.

This command will:
1. Determine the package status (new package, existing package with new preview version, existing package with new stable version)
2. For existing packages: compare current package exports with previous released version and calculate the change

The package path should be an absolute path to a Go module (containing go.mod file).

Examples:
  # get sdk change for an existing package
  generator sdkchange /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute /path/to/outputJsonFile

  # get sdk change with verbose output
  generator sdkchange /path/to/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute /path/to/outputJsonFile --verbose`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			packagePath := args[0]

			outputJsonFile := args[1]

			// Validate package path
			if err := utils.ValidatePackagePath(packagePath); err != nil {
				return fmt.Errorf("package path validation error: %v", err)
			}

			result := &SdkChangeResult{
				PackagePath: packagePath,
			}
			// Get SDK root path
			sdkRoot, err := utils.FindSDKRoot(packagePath)
			if err != nil {
				return fmt.Errorf("failed to find SDK root: %v", err)
			}

			// Initialize SDK repo
			sdkRepo, err := common.GetSDKRepo(sdkRoot, "")
			if err != nil {
				result.Success = false
				result.Message = fmt.Sprintf("Failed to initialize SDK repository: %v", err)
				return err
			}

			// Determine package status
			status, err := changelog.DetermineModuleStatus(packagePath, sdkRepo)
			if err != nil {
				result.Success = false
				result.Message = fmt.Sprintf("Failed to determine package status: %v", err)
				return err
			}

			if status != utils.PackageStatusExisting {
				result.Success = true
				result.HasBreakingChange = false
			}

			// Generate changelog
			isCurrentPreview, err := version.IsCurrentPreviewVersion(packagePath, sdkRepo, nil)
			if err != nil {
				return fmt.Errorf("failed to determine if current version is preview: %v", err)
			}
			changelogResult, err := changelog.GenerateChangelog(packagePath, sdkRepo, isCurrentPreview)
			if err != nil {
				return fmt.Errorf("failed to generate changelog: %v", err)
			}

			result.ChangelogMD = changelogResult.ChangelogData.ToCompactMarkdown()

			result.HasBreakingChange = changelogResult.ChangelogData.HasBreakingChanges()

			result.Success = true

			// Output result
			jsonResult, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				return fmt.Errorf("failed to marshal result: %v", err)
			}
			fmt.Println(string(jsonResult))
			// path := filepath.Join(packagePath, utils.SDKChangeJsonFile)
			if err = os.WriteFile(outputJsonFile, []byte(string(jsonResult)), 0644); err != nil {
				return err
			}

			return nil
		},
	}

	changelogCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", "Output format (text|json)")
	changelogCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	return changelogCmd
}
