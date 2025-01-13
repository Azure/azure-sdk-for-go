// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package dirs

import (
	"os"
	"path/filepath"
)

// GetSubdirs returns all of the subdirectories under current.
// Returns an empty slice if current contains no subdirectories.
func GetSubdirs(current string) ([]string, error) {
	children, err := os.ReadDir(current)
	if err != nil {
		return nil, err
	}
	dirs := []string{}
	for _, info := range children {
		if info.IsDir() {
			dirs = append(dirs, info.Name())
		}
	}
	return dirs, nil
}

// DeleteChildDirs deletes all child directories in the specified directory.
func DeleteChildDirs(dir string) error {
	children, err := GetSubdirs(dir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, child := range children {
		childPath := filepath.Join(dir, child)
		err = os.RemoveAll(childPath)
		if err != nil {
			return err
		}
	}
	return nil
}
