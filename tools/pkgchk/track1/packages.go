// Copyright 2021 Microsoft Corporation
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package track1

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Package represents a track 1 SDK package
type Package struct {
	root    string
	dir     string
	pkgName string
}

// Root ...
func (p Package) Root() string {
	return p.root
}

// Path ...
func (p Package) Path() string {
	path, _ := filepath.Rel(p.root, p.dir)
	return strings.ReplaceAll(path, "\\", "/")
}

// FullPath ...
func (p Package) FullPath() string {
	return p.dir
}

// Name ...
func (p Package) Name() string {
	return p.pkgName
}

// IsARMPackage ...
func (p Package) IsARMPackage() bool {
	return strings.Index(p.Path(), "/mgmt/") > -1
}

// List lists all the track 1 SDK packages under the root directory
func List(root string) ([]Package, error) {
	var results []Package
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// check if leaf dir
			fi, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}
			hasSubDirs := false
			interfacesDir := false
			for _, f := range fi {
				if f.IsDir() {
					hasSubDirs = true
					break
				}
				if f.Name() == "interfaces.go" {
					interfacesDir = true
				}
			}
			if !hasSubDirs {
				fs := token.NewFileSet()
				// with interfaces codegen the majority of leaf directories are now the
				// *api packages. when this is the case parse from the parent directory.
				if interfacesDir {
					path = filepath.Dir(path)
				}
				packages, err := parser.ParseDir(fs, path, func(fi os.FileInfo) bool {
					return fi.Name() != "interfaces.go"
				}, parser.PackageClauseOnly)
				if err != nil {
					return err
				}
				if len(packages) < 1 {
					return fmt.Errorf("did not find any packages which is unexpected")
				}
				if len(packages) > 1 {
					return fmt.Errorf("found more than one package which is unexpected")
				}
				pkgName := ""
				for _, pkg := range packages {
					pkgName = pkg.Name
				}
				// normalize the separator
				results = append(results, Package{
					root:    root,
					dir:     strings.ReplaceAll(path, "\\", "/"),
					pkgName: pkgName,
				})
			}
		}
		return nil
	})
	return results, err
}

// VerifyWithDefaultVerifiers verifies packages under the root directory with the given exceptions
func VerifyWithDefaultVerifiers(root string, exceptions map[string]bool) error {
	pkgs, err := List(root)
	if err != nil {
		return fmt.Errorf("failed to get packages: %+v", err)
	}
	verifier := GetDefaultVerifier()
	count := 0
	for _, pkg := range pkgs {
		for _, err := range verifier.Verify(pkg) {
			if _, ok := exceptions[err.Error()]; !ok {
				fmt.Fprintln(os.Stderr, err)
				count++
			}
		}
	}
	if count > 0 {
		return fmt.Errorf("found %d errors", count)
	}
	return nil
}
