// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package changelog

import (
	"os"
	"path/filepath"
	"testing"

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

		err = validatePackagePath(noGoModDir)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "does not contain a go.mod file")
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
