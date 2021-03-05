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

type Package struct {
	root    string
	dir     string
	pkgName string
}

func (p Package) Root() string {
	return p.root
}

func (p Package) Path() string {
	path, _ := filepath.Rel(p.root, p.dir)
	return strings.ReplaceAll(path, "\\", "/")
}

func (p Package) FullPath() string {
	return p.dir
}

func (p Package) Name() string {
	return p.pkgName
}

func (p Package) IsARMPackage() bool {
	return strings.Index(p.Path(), "/mgmt/") > -1
}

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
