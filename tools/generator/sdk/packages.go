package sdk

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

const apiDirSuffix = "api"

// TODO -- use the functions from the sdk itself instead when there is one in sdk
func GetPackages(root string) ([]string, error) {
	var pkgDirs []string
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
			hasApiSubDirs := false
			for _, f := range fi {
				// check if this is the interfaces subdir, if it is don't count it as a subdir
				if f.IsDir() && f.Name() != filepath.Base(path)+apiDirSuffix {
					hasSubDirs = true
				}
				if f.IsDir() && f.Name() == filepath.Base(path)+apiDirSuffix {
					hasApiSubDirs = true
				}
			}
			if !hasSubDirs && hasApiSubDirs {
				pkgDirs = append(pkgDirs, path)
				// skip any dirs under us (i.e. interfaces subdir)
				return filepath.SkipDir
			}
		}
		return nil
	})
	return pkgDirs, err
}
