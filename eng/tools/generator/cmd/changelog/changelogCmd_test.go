// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package changelog

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCommandCreation(t *testing.T) {
	cmd := Command()
	assert.NotNil(t, cmd)
	assert.Equal(t, "changelog", cmd.Use[:9])
	assert.Contains(t, cmd.Short, "Update changelog content")
}

func TestValidatePackagePath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-validate-package")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("Valid package path with go.mod", func(t *testing.T) {
		// Create a directory with go.mod file
		packageDir := filepath.Join(tempDir, "package")
		err := os.MkdirAll(packageDir, 0755)
		require.NoError(t, err)

		// Create a go.mod file
		goModPath := filepath.Join(packageDir, "go.mod")
		err = os.WriteFile(goModPath, []byte("module test\n\ngo 1.20\n"), 0644)
		require.NoError(t, err)

		err = utils.ValidatePackagePath(packageDir)
		assert.NoError(t, err)
	})

	t.Run("Non-existent path", func(t *testing.T) {
		err := utils.ValidatePackagePath("/non/existent/path")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("Directory without go.mod file", func(t *testing.T) {
		// Create a directory without go.mod
		noGoModDir := filepath.Join(tempDir, "no-go-mod")
		err := os.MkdirAll(noGoModDir, 0755)
		require.NoError(t, err)

		err = utils.ValidatePackagePath(noGoModDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not contain a go.mod file")
	})

	t.Run("File instead of directory", func(t *testing.T) {
		// Create a file instead of directory
		filePath := filepath.Join(tempDir, "testfile.txt")
		err := os.WriteFile(filePath, []byte("test"), 0644)
		require.NoError(t, err)

		err = utils.ValidatePackagePath(filePath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist or is not a directory")
	})
}

// TestReportFileFlag_NewPackageOutput exercises the changelog command end-to-end
// in --report-file (report-only) mode against a brand-new (un-released) module.
// It verifies that:
//   - The command succeeds.
//   - A JSON report is written to the given file with the expected fields.
//   - CHANGELOG.md is NOT created in the package (report-only mode is non-mutating).
func TestReportFileFlag_NewPackageOutput(t *testing.T) {
	// Initialize a temp directory as a git repository so FindSDKRoot / OpenSDKRepository succeed.
	tempDir := t.TempDir()

	gitInit := exec.Command("git", "init")
	gitInit.Dir = tempDir
	if err := gitInit.Run(); err != nil {
		t.Skip("Git not available, skipping test")
	}

	// Create a fake SDK package layout: <root>/sdk/foo/armbar with a go.mod file.
	packagePath := filepath.Join(tempDir, "sdk", "foo", "armbar")
	require.NoError(t, os.MkdirAll(packagePath, 0755))
	require.NoError(t, os.WriteFile(
		filepath.Join(packagePath, "go.mod"),
		[]byte("module github.com/Azure/azure-sdk-for-go/sdk/foo/armbar\n\ngo 1.23\n"),
		0644,
	))

	reportFile := filepath.Join(tempDir, "sdkchange.json")

	cmd := Command()
	cmd.SetArgs([]string{packagePath, "--report-file", reportFile})
	// Silence the cobra-managed output streams; the command also writes JSON via fmt.Println,
	// which will surface in test output but is harmless.
	cmd.SetOut(os.Stderr)
	cmd.SetErr(os.Stderr)

	require.NoError(t, cmd.Execute(), "report-only changelog execution should succeed for a new package")

	// Verify the report file was written and parses as the expected schema.
	data, err := os.ReadFile(reportFile)
	require.NoError(t, err, "report file should be written when --report-file is set")
	require.NotEmpty(t, data, "report file should not be empty")

	var result ChangelogResult
	require.NoError(t, json.Unmarshal(data, &result), "report file should contain valid ChangelogResult JSON")

	assert.True(t, result.Success, "Success should be true for a new package report")
	assert.Equal(t, packagePath, result.PackagePath)
	assert.Equal(t, "New Package", result.PackageStatus)
	assert.False(t, result.HasBreakingChange, "a new package has no breaking changes")
	assert.NotEmpty(t, result.ChangelogMD, "ChangelogMD should be populated in the report")
	assert.NotEmpty(t, result.ReleaseDate, "ReleaseDate should be populated")

	// Verify CHANGELOG.md was NOT created (report-only mode must not mutate the package).
	_, err = os.Stat(filepath.Join(packagePath, utils.ChangelogFileName))
	assert.True(t, os.IsNotExist(err), "report-only mode must not create CHANGELOG.md, got err=%v", err)
}
