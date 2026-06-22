// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
)

// FindSDKRoot finds the SDK repository root by walking up the directory tree
// until it finds a .git directory or file (for worktree support)
func FindSDKRoot(packagePath string) (string, error) {
	path := packagePath
	maxLevel := 10 // Prevent infinite loops

	for i := 0; i < maxLevel; i++ {
		// Check if .git exists at this level (can be a directory or file for worktrees)
		gitPath := filepath.Join(path, ".git")
		if info, err := os.Stat(gitPath); err == nil {
			// .git can be a directory (regular repo) or a file (worktree)
			if info.IsDir() || info.Mode().IsRegular() {
				return path, nil
			}
		}

		// Move up one directory
		parent := filepath.Dir(path)
		if parent == path {
			// We've reached the root of the filesystem
			break
		}
		path = parent
	}

	return "", fmt.Errorf("could not find SDK root (no .git directory or file found) starting from '%s'", packagePath)
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

// ValidatePackagePath validates that the provided package path exists and contains a go.mod file
func ValidatePackagePath(packagePath string) error {
	if info, err := os.Stat(packagePath); err != nil || !info.IsDir() {
		return fmt.Errorf("package path '%s' does not exist or is not a directory", packagePath)
	}

	goModPath := filepath.Join(packagePath, "go.mod")
	if _, err := os.Stat(goModPath); os.IsNotExist(err) {
		return fmt.Errorf("package path '%s' does not contain a go.mod file", packagePath)
	} else if err != nil {
		return fmt.Errorf("failed to check for go.mod file: %v", err)
	}

	return nil
}

// IsMgmtPlanePackage determines if the package is a management plane package
// by checking if its path relative to sdkRoot starts with "sdk/resourcemanager/"
func IsMgmtPlanePackage(packagePath, sdkRoot string) bool {
	relPath, err := filepath.Rel(sdkRoot, packagePath)
	if err != nil {
		return false
	}
	relPath = filepath.ToSlash(relPath)
	return strings.HasPrefix(relPath, "sdk/resourcemanager/")
}

// GetModuleRelativePath returns the module path relative to the SDK root
// (e.g., "sdk/resourcemanager/compute/armcompute")
func GetModuleRelativePath(packagePath, sdkRoot string) (string, error) {
	relPath, err := filepath.Rel(sdkRoot, packagePath)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path: %v", err)
	}
	return filepath.ToSlash(relPath), nil
}

// GetServiceDir returns the service directory for the pipeline
// by stripping the "sdk/" prefix (e.g., "resourcemanager/compute/armcompute")
func GetServiceDir(moduleRelativePath string) string {
	return strings.TrimPrefix(moduleRelativePath, "sdk/")
}
