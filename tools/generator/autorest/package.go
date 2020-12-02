package autorest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ChangedPackagesMap is a wrapper of a map of packages. The key is package names, and the value is the changed file list.
type ChangedPackagesMap map[string][]string

func (c *ChangedPackagesMap) addFileToPackage(pkg, file string) {
	pkg = strings.ReplaceAll(pkg, "\\", "/")
	if _, ok := (*c)[pkg]; !ok {
		(*c)[pkg] = []string{}
	}
	(*c)[pkg] = append((*c)[pkg], file)
}

func (c *ChangedPackagesMap) String() string {
	var r []string
	for k, v := range *c {
		r = append(r, fmt.Sprintf("%s: %+v", k, v))
	}
	return strings.Join(r, "\n")
}

// GetChangedPackages get the go SDK packages map from the given changed file list.
// the map returned has the package full path as key, and the changed files in the package as the value.
// This function identify the package by checking if a directory has both a `version.go` file and a `client.go` file.
func GetChangedPackages(changedFiles []string) (ChangedPackagesMap, error) {
	changedFiles, err := ExpandChangedDirectories(changedFiles)
	if err != nil {
		return nil, err
	}
	r := ChangedPackagesMap{}
	for _, file := range changedFiles {
		fi, err := os.Stat(file)
		isDir := false
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if fi != nil {
			isDir = fi.IsDir()
		}
		path := file
		if !isDir {
			path = filepath.Dir(file)
		}
		if IsValidPackage(path) {
			r.addFileToPackage(path, file)
		}
	}
	return r, nil
}

// ExpandChangedDirectories expands every directory listed in the array to all its file
func ExpandChangedDirectories(changedFiles []string) ([]string, error) {
	var result []string
	for _, path := range changedFiles {
		fi, err := os.Stat(path)
		isDir := false
		if err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		if fi != nil {
			isDir = fi.IsDir()
		}
		if isDir {
			siblings, err := getAllFiles(path)
			if err != nil {
				return nil, err
			}
			result = append(result, siblings...)
		} else {
			result = append(result, path)
		}
	}

	return result, nil
}

func getAllFiles(root string) ([]string, error) {
	var siblings []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			siblings = append(siblings, strings.ReplaceAll(path, "\\", "/"))
		}
		return nil
	})
	return siblings, err
}

const (
	clientGo  = "client.go"
	versionGo = "version.go"
)

// IsValidPackage returns true when the given directory is a valid azure-sdk-for-go package,
// otherwise returns false (including the directory does not exist)
// The criteria of a valid azure-sdk-for-go package is that the directory must have a `client.go` and a `version.go` file
func IsValidPackage(dir string) bool {
	client := filepath.Join(dir, clientGo)
	version := filepath.Join(dir, versionGo)
	// both the above files must exist to return true
	if _, err := os.Stat(client); os.IsNotExist(err) {
		return false
	}
	if _, err := os.Stat(version); os.IsNotExist(err) {
		return false
	}
	return true
}
