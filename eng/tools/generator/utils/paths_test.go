// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindSDKRoot_RegularGitRepo(t *testing.T) {
	// Create a temporary directory structure simulating a regular git repo
	tmpDir := t.TempDir()

	// Create .git directory (regular repo)
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	require.NoError(t, err)

	// Create nested package directory
	packageDir := filepath.Join(tmpDir, "sdk", "resourcemanager", "compute")
	err = os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)

	// Test finding root from package directory
	root, err := FindSDKRoot(packageDir)
	require.NoError(t, err)
	require.Equal(t, tmpDir, root)

	// Test finding root from the root itself
	root, err = FindSDKRoot(tmpDir)
	require.NoError(t, err)
	require.Equal(t, tmpDir, root)
}

func TestFindSDKRoot_GitWorktree(t *testing.T) {
	// Create a temporary directory structure simulating a git worktree
	tmpDir := t.TempDir()

	// Create .git file (worktree style) with content pointing to main repo
	gitFile := filepath.Join(tmpDir, ".git")
	// In a real worktree, this file contains something like:
	// gitdir: /path/to/main/repo/.git/worktrees/worktree-name
	err := os.WriteFile(gitFile, []byte("gitdir: /some/path/.git/worktrees/myworktree"), 0644)
	require.NoError(t, err)

	// Create nested package directory
	packageDir := filepath.Join(tmpDir, "sdk", "resourcemanager", "compute")
	err = os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)

	// Test finding root from package directory
	root, err := FindSDKRoot(packageDir)
	require.NoError(t, err)
	require.Equal(t, tmpDir, root)

	// Test finding root from the root itself
	root, err = FindSDKRoot(tmpDir)
	require.NoError(t, err)
	require.Equal(t, tmpDir, root)
}

func TestFindSDKRoot_NoGitDirectory(t *testing.T) {
	// Create a temporary directory without .git
	tmpDir := t.TempDir()

	// Create nested directory
	packageDir := filepath.Join(tmpDir, "sdk", "resourcemanager", "compute")
	err := os.MkdirAll(packageDir, 0755)
	require.NoError(t, err)

	// Test that it returns an error when no .git is found
	_, err = FindSDKRoot(packageDir)
	require.Error(t, err)
	require.Contains(t, err.Error(), "could not find SDK root")
}

func TestFindSDKRoot_DeepNesting(t *testing.T) {
	// Create a temporary directory with deep nesting
	tmpDir := t.TempDir()

	// Create .git directory
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	require.NoError(t, err)

	// Create deeply nested directory (within maxLevel)
	deepDir := filepath.Join(tmpDir, "a", "b", "c", "d", "e", "f", "g", "h")
	err = os.MkdirAll(deepDir, 0755)
	require.NoError(t, err)

	// Test finding root from deep directory
	root, err := FindSDKRoot(deepDir)
	require.NoError(t, err)
	require.Equal(t, tmpDir, root)
}

func TestFindSDKRoot_MaxLevelExceeded(t *testing.T) {
	// Create a temporary directory with nesting exceeding maxLevel
	tmpDir := t.TempDir()

	// Create .git directory at root
	gitDir := filepath.Join(tmpDir, ".git")
	err := os.Mkdir(gitDir, 0755)
	require.NoError(t, err)

	// Create directory nested deeper than maxLevel (10)
	// We need 11 levels deep from a starting point that doesn't contain .git
	deepDir := filepath.Join(tmpDir, "1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11")
	err = os.MkdirAll(deepDir, 0755)
	require.NoError(t, err)

	// Starting from level 11, we should still find root within 10 iterations
	// since tmpDir/.git is 11 levels up
	_, err = FindSDKRoot(deepDir)
	// This should fail because maxLevel is 10 and we need 11 levels
	require.Error(t, err)
	require.Contains(t, err.Error(), "could not find SDK root")
}
