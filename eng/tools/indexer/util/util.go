// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package util

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"regexp"
	"strings"
)

// PackageSet is a collection of Go packages.
// Key is the full package import path, value indicates if the package has been indexed.
type PackageSet map[string]bool

// GetIndexedPackages returns the set of packages that have already been indexed.
// It finds all entries matching the regex `github\.com/Azure/azure-sdk-for-go/services/.*?"`.
func GetIndexedPackages(content io.Reader) (PackageSet, error) {
	body, err := ioutil.ReadAll(content)
	if err != nil {
		return nil, err
	}

	if len(body) < 1 {
		return nil, errors.New("did't receive a response body when lookinig for indexed packages")
	}

	// scrape the content to create the package list
	pkgs := PackageSet{}
	regex := regexp.MustCompile(`github\.com/Azure/azure-sdk-for-go/services/.*?"`)
	finds := regex.FindAllString(string(body), -1)

	for _, find := range finds {
		// strip of the trailing "
		pkg := find[:len(find)-1]
		pkgs[pkg] = true
	}
	return pkgs, nil
}

// GetPackagesForIndexing returns the set of packages, calculated from the specified directory, to be indexed.
// Each directory entry is converted to a complete package path, e.g. "github.com/Azure/azure-sdk-for-go/services/foo/...".
func GetPackagesForIndexing(dir string) (PackageSet, error) {
	leafDirs := []string{}
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// check if leaf dir
			fi, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}
			hasSubDirs := false
			for _, f := range fi {
				if f.IsDir() {
					hasSubDirs = true
					break
				}
			}
			if !hasSubDirs {
				leafDirs = append(leafDirs, path)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	// dirs will look something like "D:\work\src\github.com\Azure\azure-sdk-for-go\services\..."
	// strip off the stuff before the github.com and change the whacks so it looks like a package import

	pkgs := PackageSet{}
	for _, dir := range leafDirs {
		i := strings.Index(dir, "github.com")
		if i < 0 {
			return nil, fmt.Errorf("didn't find github.com in directory '%s'", dir)
		}
		pkg := strings.Replace(dir[i:], "\\", "/", -1)
		pkgs[pkg] = false
	}
	return pkgs, nil
}
