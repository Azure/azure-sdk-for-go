// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package utils

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
)

// FindSDKRoot finds the SDK repository root by walking up the directory tree
// until it finds a .git directory
func FindSDKRoot(packagePath string) (string, error) {
	path := packagePath
	maxLevel := 10 // Prevent infinite loops

	for i := 0; i < maxLevel; i++ {
		// Check if .git directory exists at this level
		gitPath := filepath.Join(path, ".git")
		if info, err := os.Stat(gitPath); err == nil && info.IsDir() {
			return path, nil
		}

		// Move up one directory
		parent := filepath.Dir(path)
		if parent == path {
			// We've reached the root of the filesystem
			break
		}
		path = parent
	}

	return "", fmt.Errorf("could not find SDK root (no .git directory found) starting from '%s'", packagePath)
}

// GetRelativePath returns the path of 'path' relative to the SDK repository root
func GetRelativePath(path string, sdkRepo repo.SDKRepository) (string, error) {
	relativePath, err := filepath.Rel(sdkRepo.Root(), path)
	if err != nil {
		return "", err
	}
	relativePath = filepath.ToSlash(relativePath)
	return relativePath, nil
}
