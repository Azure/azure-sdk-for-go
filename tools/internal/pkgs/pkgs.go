// Copyright 2018 Microsoft Corporation
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

package pkgs

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Pkg provides a package of a service
type Pkg struct {
	// the directory where the package relative to the root dir
	Dest string
	// the ast of the package
	Package *ast.Package
}

// IsARMPkg returns if this package is ARM package (management plane)
func (p Pkg) IsARMPkg() bool {
	return strings.Index(p.Dest, "/mgmt/") > -1
}

// GetAPIVersion returns the api version of this package
func (p Pkg) GetAPIVersion() (string, error) {
	dest := p.Dest
	if p.IsARMPkg() {
		// management-plane
		regex := regexp.MustCompile(`mgmt/(.+)/`)
		groups := regex.FindStringSubmatch(dest)
		if len(groups) < 2 {
			return "", fmt.Errorf("cannot find api version in %s", dest)
		}
		versionString := groups[1]
		return versionString, nil
	}
	// data-plane
	regex := regexp.MustCompile(`/(\d{4}-\d{2}.*|v?\d+(\.\d+)?)/`)
	groups := regex.FindStringSubmatch(dest)
	if len(groups) < 2 {
		return "", fmt.Errorf("cannot find api version in %s", dest)
	}
	versionString := groups[1]
	if versionString == "" {
		return "", fmt.Errorf("does not find api version in data plane package %s", dest)
	}
	return versionString, nil
}

// GetPkgs returns every package under the rootDir. Package for interfaces will be ignored.
func GetPkgs(rootDir string) ([]Pkg, error) {
	var pkgs []Pkg
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
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
				packages, err := parser.ParseDir(fs, path, func(fileInfo os.FileInfo) bool {
					return fileInfo.Name() != "interfaces.go"
				}, parser.PackageClauseOnly)
				if err != nil {
					return nil
				}
				if len(packages) < 1 {
					return errors.New("didn't find any packages which is unexpected")
				}
				if len(packages) > 1 {
					return errors.New("found more than one package which is unexpected")
				}
				var p *ast.Package
				for _, pkgs := range packages {
					p = pkgs
				}
				// normalize directory separator to '/' character
				pkgs = append(pkgs, Pkg{
					Dest:    strings.ReplaceAll(path[len(rootDir):], "\\", "/"),
					Package: p,
				})
			}
		}
		return nil
	})
	return pkgs, err
}
