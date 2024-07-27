// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package ioext

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// CopyDir recursively copies the specified source directory tree to the
// specified destination.  The destination directory must not exist.  Any
// symlinks under src are ignored.
func CopyDir(src, dst string) error {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	// verify that src is a directory
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !srcInfo.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	// now verify that dst doesn't exist
	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	if err == nil {
		return fmt.Errorf("destination directory already exists")
	}

	err = os.MkdirAll(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	// get the collection of directory entries under src.
	// for each entry build its corresponding path under dst.
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		// skip symlinks
		if entry.Type()&os.ModeSymlink != 0 {
			continue
		}

		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(srcPath, dstPath, true)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// CopyFile copies the specified source file to the specified destination file.
// Specify true for overwrite to overwrite the destination file if it already exits.
func CopyFile(src, dst string, overwrite bool) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	if !overwrite {
		// check if the file exists, if it does then return an error
		_, err := os.Stat(dst)
		if err != nil && !os.IsNotExist(err) {
			return errors.New("won't overwrite destination file")
		}
	}

	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// flush the buffer
	err = dstFile.Sync()
	if err != nil {
		return err
	}

	// copy file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	err = os.Chmod(dst, srcInfo.Mode())
	if err != nil {
		return err
	}

	return nil
}
