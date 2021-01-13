package utils

import (
	"io/ioutil"
	"os"
	"path/filepath"
)

// ListTrack1SDKPackages will return all the full path of track 1 SDK packages in a slice
// a track 1 SDK package is a golang package that contains a subdirectory `{packageName}api`
// and a file `version.go`
func ListTrack1SDKPackages(root string) ([]string, error) {
	var results []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			// check if this is a leaf dir
			fi, err := ioutil.ReadDir(path)
			if err != nil {
				return err
			}
			base := filepath.Base(path)
			hasSubDir := false
			hasAPISubDirs := false
			hasVersionFile := false
			for _, f := range fi {
				hasSubDir = hasSubDir || isSubDir(base, f)
				hasAPISubDirs = hasAPISubDirs || isAPIDir(base, f)
				hasVersionFile = hasVersionFile || isVersionFile(f)
			}
			if !hasSubDir && hasAPISubDirs && hasVersionFile {
				results = append(results, path)
				// skip all subdirectories
				return filepath.SkipDir
			}
		}
		return nil
	})
	return results, err
}

func isSubDir(base string, fi os.FileInfo) bool {
	return fi.IsDir() && fi.Name() != base+apiDirSuffix
}

func isAPIDir(base string, fi os.FileInfo) bool {
	return fi.IsDir() && fi.Name() == base+apiDirSuffix
}

func isVersionFile(fi os.FileInfo) bool {
	return !fi.IsDir() && fi.Name() == versionFileName
}

const (
	apiDirSuffix    = "api"
	versionFileName = "version.go"
)
