// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package changelog

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	internalutils "github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
	"github.com/Masterminds/semver"
	"github.com/stretchr/testify/require"
)

func TestTypeToAny(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/toany")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/toany")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := getChangelog(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	FilterChangelog(changelog, TypeToAnyFilter)

	excepted := "### Breaking Changes\n\n- Type of `Client.M` has been changed from `map[string]string` to `map[string]any`\n\n### Features Added\n\n- Type of `Client.A` has been changed from `*int` to `any`\n"
	require.Equal(t, excepted, changelog.ToCompactMarkdown())
}

func TestFuncParameterChange(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/parameter")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/parameter")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := getChangelog(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	FilterChangelog(changelog, FuncFilter)

	excepted := "### Breaking Changes\n\n- Function `*Client.AfterAny` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, serviceName string, value interface{}, option ClientOption)` to `(ctx context.Context, resourceGroupName string, serviceName string, value any, option Option)`\n- Function `*Client.BeforeAny` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, serviceName string, value interface{}, option ClientOption)` to `(ctx context.Context, resourceGroupName string, serviceName any, value any, option ClientOption)`\n- Function `*Client.OnlyToAny` parameter(s) have been changed from `(ctx context.Context, resourceGroupName string, serviceName string, value interface{}, option ClientOption)` to `(ctx context.Context, resourceGroupName string, serviceName string, value any, option ClientOption)`\n"
	require.Equal(t, excepted, changelog.ToCompactMarkdown())
}

func TestFuncParameterOrderChange(t *testing.T) {
	oldExport, err := exports.Get("./testdata/old/paramorder")
	if err != nil {
		t.Fatal(err)
	}

	newExport, err := exports.Get("./testdata/new/paramorder")
	if err != nil {
		t.Fatal(err)
	}

	changelog, err := getChangelog(&oldExport, &newExport)
	if err != nil {
		t.Fatal(err)
	}

	FilterChangelog(changelog, FuncFilter)

	// Expected: Functions with parameter changes (type, name, or order) should be detected
	// Now the output includes parameter names to make it clear what changed
	// - NewListByServicePager: serviceName and resourceGroupName swapped (order change)
	// - OrderChanged: resourceGroupName, serviceName, subscriptionID -> serviceName, subscriptionID, resourceGroupName (order change)
	// - DifferentNames: oldName, newName -> firstName, lastName (name change)
	// - NoChange: should not appear (no change)
	excepted := "### Breaking Changes\n\n- Function `*AllPoliciesClient.DifferentNames` parameter(s) have been changed from `(oldName string, newName string)` to `(firstName string, lastName string)`\n- Function `*AllPoliciesClient.NewListByServicePager` parameter(s) have been changed from `(resourceGroupName string, serviceName string, options *AllPoliciesClientListByServiceOptions)` to `(serviceName string, resourceGroupName string, options *AllPoliciesClientListByServiceOptions)`\n- Function `*AllPoliciesClient.OrderChanged` parameter(s) have been changed from `(resourceGroupName string, serviceName string, subscriptionID string)` to `(serviceName string, subscriptionID string, resourceGroupName string)`\n"
	require.Equal(t, excepted, changelog.ToCompactMarkdown())
}

func TestGetAllVersionTags(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	log.Printf("Using current directory as SDK root: %s", cwd)

	// create sdk repo ref
	sdkRepo, err := repo.OpenSDKRepository(internalutils.NormalizePath(cwd))
	if err != nil {
		t.Fatal(err)
	}

	tags, err := GetAllVersionTags("refs/tags/sdk/azidentity", sdkRepo)
	if err != nil {
		t.Fatal(err)
	}
	expected := "refs/tags/sdk/azidentity/v0.1.0"
	require.Contains(t, tags, expected)
	expected = "refs/tags/sdk/azidentity/v1.10.0"
	require.Contains(t, tags, expected)
	require.GreaterOrEqual(t, len(tags), 69)

	if tags, err = GetAllVersionTags("sdk/resourcemanager/network/armnetwork", sdkRepo); err != nil {
		t.Fatal(err)
	}
	expected = "refs/tags/sdk/resourcemanager/network/armnetwork/v0.1.0"
	require.Contains(t, tags, expected)
	expected = "refs/tags/sdk/resourcemanager/network/armnetwork/v7.0.0"
	require.Contains(t, tags, expected)
	require.GreaterOrEqual(t, len(tags), 30)
}

func TestGetExportsFromTag(t *testing.T) {
	// Test cases with different package paths and versions
	testCases := []struct {
		name        string
		packagePath string
		tag         string
	}{
		{
			name:        "IoT Firmware Defense v2.0.0-beta.1",
			packagePath: "sdk/resourcemanager/iotfirmwaredefense/armiotfirmwaredefense",
			tag:         "sdk/resourcemanager/iotfirmwaredefense/armiotfirmwaredefense/v2.0.0-beta.1",
		},
		{
			name:        "Compute Schedule v1.1.0",
			packagePath: "sdk/resourcemanager/computeschedule/armcomputeschedule",
			tag:         "sdk/resourcemanager/computeschedule/armcomputeschedule/v1.1.0",
		},
		{
			name:        "Edge Order v1.2.0",
			packagePath: "sdk\\resourcemanager\\edgeorder\\armedgeorder",
			tag:         "sdk/resourcemanager/edgeorder/armedgeorder/v1.2.0",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Call GetExportsFromTag
			exports, err := getExportsFromTag(tc.packagePath, tc.tag)

			// Should not error for valid tags
			require.NoError(t, err)
			require.NotNil(t, exports)

			// Log some information about the exports for debugging
			log.Printf("Test %s: Found exports for package %s at tag %s", tc.name, tc.packagePath, tc.tag)

			require.True(t, len(exports.Funcs) > 0 || len(exports.Structs) > 0,
				"Expected to find some exports in package")
		})
	}
}

func TestGetPreviousVersionTag(t *testing.T) {
	tags := []string{"v1.2.0", "v1.1.0-beta.1", "v1.0.0"}

	// Test preview - should return latest
	result := getPreviousVersionTag(true, tags)
	require.Equal(t, "v1.2.0", result)

	// Test stable - should return latest stable
	result = getPreviousVersionTag(false, tags)
	require.Equal(t, "v1.2.0", result)

	tags = []string{"v1.1.0-beta.2", "v1.1.0-beta.1", "v1.0.0"}

	// Test preview - should return latest beta
	result = getPreviousVersionTag(true, tags)
	require.Equal(t, "v1.1.0-beta.2", result)

	// Test stable - should return latest stable
	result = getPreviousVersionTag(false, tags)
	require.Equal(t, "v1.0.0", result)
}

func TestUpdateLatestChangelogVersion(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "test-update-changelog")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	testCases := []struct {
		name            string
		initialContent  string
		newVersion      string
		releaseDate     string
		expectedContent string
		expectError     bool
	}{
		{
			name: "Update unreleased version",
			initialContent: `# Release History

## 1.0.0 (Unreleased)

### Features Added

- Added new feature X
- Added new feature Y

## 0.9.0 (2025-01-01)

### Features Added

- Initial release
`,
			newVersion:  "1.0.0",
			releaseDate: "2026-01-06",
			expectedContent: `# Release History

## 1.0.0 (2026-01-06)

### Features Added

- Added new feature X
- Added new feature Y

## 0.9.0 (2025-01-01)

### Features Added

- Initial release
`,
			expectError: false,
		},
		{
			name: "Update with different version",
			initialContent: `# Release History

## 0.5.0 (Unreleased)

### Bugs Fixed

- Fixed bug A

## 0.4.0 (2025-12-01)

### Features Added

- Feature B
`,
			newVersion:  "0.5.0",
			releaseDate: "2026-01-06",
			expectedContent: `# Release History

## 0.5.0 (2026-01-06)

### Bugs Fixed

- Fixed bug A

## 0.4.0 (2025-12-01)

### Features Added

- Feature B
`,
			expectError: false,
		},
		{
			name: "No version entry",
			initialContent: `# Release History

No versions yet.
`,
			newVersion:  "1.0.0",
			releaseDate: "2026-01-06",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a temporary file with the initial content
			changelogPath := filepath.Join(tempDir, "CHANGELOG.md")
			err := os.WriteFile(changelogPath, []byte(tc.initialContent), 0644)
			if err != nil {
				t.Fatal(err)
			}

			// Call the function
			version, err := semver.NewVersion(tc.newVersion)
			if err != nil {
				t.Fatal(err)
			}

			err = UpdateLatestChangelogVersion(tempDir, version, tc.releaseDate)

			if tc.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)

				// Read the updated content
				updatedContent, err := os.ReadFile(changelogPath)
				if err != nil {
					t.Fatal(err)
				}

				require.Equal(t, tc.expectedContent, string(updatedContent))
			}
		})
	}
}
