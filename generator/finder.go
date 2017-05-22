package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/marstr/collection"
)

// SwaggerFinder will enumerate all Swaggers in a particular directory.
type SwaggerFinder struct {
	Root string
}

// NewSwaggerFinder creates a new instance of SwaggerFinder which will search the
// directory `targetDir`.
func NewSwaggerFinder(targetDir string) (constructed *SwaggerFinder, err error) {
	constructed = &SwaggerFinder{
		Root: targetDir,
	}
	return
}

// Enumerate lists out all instances of Swagger files in the SwaggerFinder's `Root`.
func (f *SwaggerFinder) Enumerate() collection.Enumerator {
	retval := make(chan interface{})

	go func() {
		defer close(retval)

		seen := map[string][]string{}

		seenContains := func(needle Swagger) bool {
			if previouslySeen, ok := seen[needle.Info.Title]; ok {
				for _, version := range previouslySeen {
					if version == needle.Info.Version {
						return true
					}
				}
			}
			return false
		}

		filepath.Walk(f.Root, func(path string, info os.FileInfo, err error) (result error) {
			if err != nil {
				return
			}

			if strings.ToLower(filepath.Ext(path)) == ".json" {
				var contents []byte
				if temp, err := ioutil.ReadFile(path); err == nil {
					contents = temp
				} else {
					return
				}

				var manifest Swagger
				if err := json.Unmarshal(contents, &manifest); err != nil {
					return
				}
				manifest.Path = strings.TrimPrefix(path, f.Root)
				manifest.Path = strings.TrimPrefix(manifest.Path, "/")

				title := manifest.Info.Title

				if title == "" {
					return
				}

				if seenContains(manifest) {
					return
				} else if versions, ok := seen[manifest.Info.Title]; ok {
					seen[manifest.Info.Title] = append(versions, manifest.Info.Version)
				} else {
					seen[manifest.Info.Title] = []string{manifest.Info.Version}
				}

				retval <- manifest
			}
			return
		})
	}()

	return retval
}
