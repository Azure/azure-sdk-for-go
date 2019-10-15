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
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Pkg struct {
	// the directory where the package relative to the root dir
	Dest string
	// the ast of the package
	Package *ast.Package
}

func (p Pkg) IsARMPkg() bool {
	return strings.Index(p.Dest, "/mgmt/") > -1
}

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
