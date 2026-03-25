// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package metadata

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// repoRoot returns the absolute path of the repository root.
func repoRoot(t *testing.T) string {
	t.Helper()
	_, filename, _, ok := runtime.Caller(0)
	require.True(t, ok, "Failed to get caller info")
	dir := filepath.Dir(filename)
	for i := 0; i < 10; i++ {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("Could not find repo root")
	return ""
}

// copyTemplate copies a template file from the repo's template/typespec directory into destDir.
func copyTemplate(t *testing.T, root, name, destDir string) {
	t.Helper()
	src := filepath.Join(root, "eng", "tools", "generator", "template", "typespec", name)
	data, err := os.ReadFile(src)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(destDir, name), data, 0644))
}

func TestInferPackageTitle(t *testing.T) {
	tests := []struct {
		name        string
		packageName string
		expected    string
	}{
		{name: "ArmPrefix", packageName: "armcompute", expected: "Compute"},
		{name: "AzPrefix", packageName: "azblob", expected: "Blob"},
		{name: "NoPrefix", packageName: "storage", expected: "Storage"},
		{name: "Empty", packageName: "", expected: ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := inferPackageTitle(tt.packageName)
			require.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessMetadata_MgmtPlaneNewPackage(t *testing.T) {
	root := repoRoot(t)

	// Create a mock SDK root with .git directory
	sdkRoot := t.TempDir()
	err := os.MkdirAll(filepath.Join(sdkRoot, ".git"), 0755)
	require.NoError(t, err)

	// Copy real templates into the mock SDK root
	templateDir := filepath.Join(sdkRoot, "eng", "tools", "generator", "template", "typespec")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)
	copyTemplate(t, root, "ci.yml.tpl", templateDir)
	copyTemplate(t, root, "README.md.tpl", templateDir)

	// Create the package directory
	packageDir := filepath.Join(sdkRoot, "sdk", "resourcemanager", "compute", "armcompute")
	err = os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(packageDir, "go.mod"), []byte("module github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute"), 0644)
	require.NoError(t, err)

	// Process metadata
	result, err := processMetadata(packageDir, "Compute", false)
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Contains(t, result.CreatedFiles, "ci.yml")
	require.Contains(t, result.CreatedFiles, "README.md")

	// Verify ci.yml was created
	ciContent, err := os.ReadFile(filepath.Join(packageDir, "ci.yml"))
	require.NoError(t, err)
	require.Contains(t, string(ciContent), "sdk/resourcemanager/compute/armcompute/")
	require.Contains(t, string(ciContent), "ServiceDirectory: 'resourcemanager/compute/armcompute'")

	// Verify README.md was created
	readmeContent, err := os.ReadFile(filepath.Join(packageDir, "README.md"))
	require.NoError(t, err)
	require.Contains(t, string(readmeContent), "Azure Compute Module for Go")
	require.Contains(t, string(readmeContent), "armcompute")
}

func TestProcessMetadata_DataPlaneNewPackage(t *testing.T) {
	root := repoRoot(t)

	// Create a mock SDK root with .git directory
	sdkRoot := t.TempDir()
	err := os.MkdirAll(filepath.Join(sdkRoot, ".git"), 0755)
	require.NoError(t, err)

	// Copy real templates into the mock SDK root
	templateDir := filepath.Join(sdkRoot, "eng", "tools", "generator", "template", "typespec")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)
	copyTemplate(t, root, "ci.yml.tpl", templateDir)

	// Create a three-layer data plane package directory (no README.md should be created)
	packageDir := filepath.Join(sdkRoot, "sdk", "monitor", "ingestion", "azlogs")
	err = os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(packageDir, "go.mod"), []byte("module github.com/Azure/azure-sdk-for-go/sdk/monitor/ingestion/azlogs"), 0644)
	require.NoError(t, err)

	// Process metadata
	result, err := processMetadata(packageDir, "", false)
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Contains(t, result.CreatedFiles, "ci.yml")
	// README.md should NOT be created for data plane packages
	for _, f := range result.CreatedFiles {
		require.NotEqual(t, "README.md", f)
	}

	// Verify ci.yml was created
	ciContent, err := os.ReadFile(filepath.Join(packageDir, "ci.yml"))
	require.NoError(t, err)
	require.Contains(t, string(ciContent), "sdk/monitor/ingestion/azlogs/")
	require.Contains(t, string(ciContent), "ServiceDirectory: 'monitor/ingestion/azlogs'")
}

func TestProcessMetadata_ExistingFiles(t *testing.T) {
	root := repoRoot(t)

	// Create a mock SDK root with .git directory
	sdkRoot := t.TempDir()
	err := os.MkdirAll(filepath.Join(sdkRoot, ".git"), 0755)
	require.NoError(t, err)

	// Copy real templates into the mock SDK root
	templateDir := filepath.Join(sdkRoot, "eng", "tools", "generator", "template", "typespec")
	err = os.MkdirAll(templateDir, 0755)
	require.NoError(t, err)
	copyTemplate(t, root, "ci.yml.tpl", templateDir)
	copyTemplate(t, root, "README.md.tpl", templateDir)

	// Create the package directory with existing files
	packageDir := filepath.Join(sdkRoot, "sdk", "resourcemanager", "compute", "armcompute")
	err = os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(packageDir, "go.mod"), []byte("module test"), 0644)
	require.NoError(t, err)

	existingCIContent := "# existing ci.yml"
	err = os.WriteFile(filepath.Join(packageDir, "ci.yml"), []byte(existingCIContent), 0644)
	require.NoError(t, err)

	existingReadmeContent := "# existing README"
	err = os.WriteFile(filepath.Join(packageDir, "README.md"), []byte(existingReadmeContent), 0644)
	require.NoError(t, err)

	// Process metadata - should skip existing files
	result, err := processMetadata(packageDir, "", false)
	require.NoError(t, err)
	require.True(t, result.Success)
	require.Empty(t, result.CreatedFiles)
	require.Contains(t, result.SkippedFiles, "ci.yml")
	require.Contains(t, result.SkippedFiles, "README.md")

	// Verify existing files were not modified
	ciContent, err := os.ReadFile(filepath.Join(packageDir, "ci.yml"))
	require.NoError(t, err)
	require.Equal(t, existingCIContent, string(ciContent))

	readmeContent, err := os.ReadFile(filepath.Join(packageDir, "README.md"))
	require.NoError(t, err)
	require.Equal(t, existingReadmeContent, string(readmeContent))
}

func TestProcessMetadata_NoSDKRoot(t *testing.T) {
	tmpDir := t.TempDir()
	packageDir := filepath.Join(tmpDir, "pkg")
	err := os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)
	err = os.WriteFile(filepath.Join(packageDir, "go.mod"), []byte("module test"), 0644)
	require.NoError(t, err)

	result, err := processMetadata(packageDir, "", false)
	require.NoError(t, err)
	require.False(t, result.Success)
	require.Contains(t, result.Message, "Failed to find SDK root")
}
