// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package version

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidatePackagePath(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "test-validate-version")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	t.Run("Valid package path with go.mod and version.go", func(t *testing.T) {
		// Create a directory with go.mod and version.go files
		packageDir := filepath.Join(tempDir, "valid-package")
		err := os.MkdirAll(packageDir, 0755)
		require.NoError(t, err)

		// Create a go.mod file
		goModPath := filepath.Join(packageDir, "go.mod")
		err = os.WriteFile(goModPath, []byte("module test\n\ngo 1.20\n"), 0644)
		require.NoError(t, err)

		// Create a version.go file
		versionGoPath := filepath.Join(packageDir, "version.go")
		versionContent := `package test

const (
	moduleName    = "test"
	moduleVersion = "v1.0.0"
)
`
		err = os.WriteFile(versionGoPath, []byte(versionContent), 0644)
		require.NoError(t, err)

		err = validatePackagePath(packageDir)
		assert.NoError(t, err)
	})

	t.Run("Non-existent path", func(t *testing.T) {
		err := validatePackagePath("/non/existent/path")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist")
	})

	t.Run("Directory without go.mod file", func(t *testing.T) {
		// Create a directory without go.mod
		noGoModDir := filepath.Join(tempDir, "no-go-mod")
		err := os.MkdirAll(noGoModDir, 0755)
		require.NoError(t, err)

		// Create version.go but no go.mod
		versionGoPath := filepath.Join(noGoModDir, "version.go")
		err = os.WriteFile(versionGoPath, []byte("package test\n"), 0644)
		require.NoError(t, err)

		err = validatePackagePath(noGoModDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not contain a go.mod file")
	})

	t.Run("Directory without version.go file", func(t *testing.T) {
		// Create a directory without version.go
		noVersionDir := filepath.Join(tempDir, "no-version")
		err := os.MkdirAll(noVersionDir, 0755)
		require.NoError(t, err)

		// Create go.mod but no version.go
		goModPath := filepath.Join(noVersionDir, "go.mod")
		err = os.WriteFile(goModPath, []byte("module test\n"), 0644)
		require.NoError(t, err)

		err = validatePackagePath(noVersionDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not contain a version.go file")
	})

	t.Run("File instead of directory", func(t *testing.T) {
		// Create a file instead of directory
		filePath := filepath.Join(tempDir, "testfile.txt")
		err := os.WriteFile(filePath, []byte("test"), 0644)
		require.NoError(t, err)

		err = validatePackagePath(filePath)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not exist or is not a directory")
	})
}
