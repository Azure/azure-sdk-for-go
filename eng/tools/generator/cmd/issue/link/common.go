// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package link

import (
	"context"
	"fmt"
	"log"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/ahmetb/go-linq/v3"
)

func init() {
	specCache = make(map[Readme]bool)
	specSubCache = make(map[Readme]Readme)

	tspCache = make(map[Readme]bool)
	tspSubCache = make(map[Readme]Readme)
}

const (
	mgmtSegment = "resource-manager"
)

// Readme represents a readme filepath
type Readme string

// IsMgmt returns true when the readme belongs to a mgmt plane package
func (r Readme) IsMgmt() bool {
	return strings.Contains(string(r), mgmtSegment)
}

func (r Readme) IsTsp() bool {
	return strings.Contains(string(r), "tspconfig.yaml")
}

// GetReadmePathFromChangedFiles ...
func GetReadmePathFromChangedFiles(ctx context.Context, client *query.Client, files []string) (Readme, error) {
	// find readme files one by one
	readmeFiles := make(map[Readme]bool)
	for _, file := range files {
		readme, err := GetReadmeFromPath(ctx, client, file)
		if err != nil {
			log.Printf("Changed file '%s' does not belong to any RP, ignoring", file)
			continue
		}
		readmeFiles[readme] = true
	}
	if len(readmeFiles) > 1 {
		return "", fmt.Errorf("cannot determine which RP to release because we have the following readme files involved: %+v", getMapKeys(readmeFiles))
	}
	if len(readmeFiles) == 0 {
		return "", fmt.Errorf("cannot get any readme files from these changed files: [%s]", strings.Join(files, ", "))
	}
	// we only have one readme file
	return getMapKeys(readmeFiles)[0], nil
}

func getMapKeys(m map[Readme]bool) []Readme {
	var result []Readme
	linq.From(m).Select(func(kv interface{}) interface{} {
		return kv.(linq.KeyValue).Key
	}).ToSlice(&result)
	return result
}

var (
	specCache map[Readme]bool
	tspCache  map[Readme]bool

	specSubCache map[Readme]Readme
	tspSubCache  map[Readme]Readme
	tempCache    map[Readme]Readme
)

// GetReadmeFromPath ...
func GetReadmeFromPath(ctx context.Context, client *query.Client, path string) (Readme, error) {
	// we do not need to determine this is a path of file or directory
	// we could always assume this is a directory path even if it is a file path - just waste a try in the first attempt if it is a filepath
	tempCache = make(map[Readme]Readme)
	isTypeSpec := ctx.Value("TypeSpec")
	if isTypeSpec != nil {
		if isTypeSpec.(bool) {
			return getTspConfigFromDirectoryPath(ctx, client, path)
		}
	}

	return getReadmeFromDirectoryPath(ctx, client, path)
}

func getReadmeFromDirectoryPath(ctx context.Context, client *query.Client, dir string) (Readme, error) {
	if len(dir) == 0 || dir == "." {
		return "", fmt.Errorf("cannot determine the readme.md path")
	}
	file := tryReadmePath(dir)
	// use cache to short cut
	if y, ok := specCache[file]; ok {
		if y {
			return file, nil
		} else {
			return getReadmeFromDirectoryPath(ctx, client, filepathDir(dir))
		}
	}

	if y, ok := specSubCache[file]; ok {
		if y != "" {
			appendMap(specSubCache, setMapValue(tempCache, y))
			return y, nil
		}
	}

	_, _, resp, err := client.Repositories.GetContents(ctx, SpecOwner, SpecRepo, string(file), nil)
	if err == nil {
		specCache[file] = true
		appendMap(specSubCache, setMapValue(tempCache, file))
		return file, nil
	}
	if resp != nil && resp.StatusCode == 404 {
		tempCache[file] = ""
		return getReadmeFromDirectoryPath(ctx, client, filepathDir(dir))
	}
	return "", err
}

func tryReadmePath(base string) Readme {
	return Readme(fmt.Sprintf("%s/readme.md", base))
}

func filepathDir(path string) string {
	return strings.ReplaceAll(filepath.Dir(path), "\\", "/")
}

func getTspConfigFromDirectoryPath(ctx context.Context, client *query.Client, dir string) (Readme, error) {
	if len(dir) == 0 || dir == "." {
		return "", fmt.Errorf("cannot determine the readme.md path")
	}

	file := tryTspConfigPath(dir)
	// use cache to short cut
	if y, ok := tspCache[file]; ok {
		if y {
			return file, nil
		} else {
			return getTspConfigFromDirectoryPath(ctx, client, filepathDir(dir))
		}
	}

	if y, ok := tspSubCache[file]; ok {
		if y != "" {
			appendMap(tspSubCache, setMapValue(tempCache, y))
			return y, nil
		}
	}

	_, _, resp, err := client.Repositories.GetContents(ctx, SpecOwner, SpecRepo, string(file), nil)
	if err == nil {
		tspCache[file] = true
		appendMap(tspSubCache, setMapValue(tempCache, file))
		return file, nil
	}
	if resp != nil && resp.StatusCode == 404 {
		tempCache[file] = ""
		return getTspConfigFromDirectoryPath(ctx, client, filepathDir(dir))
	}
	return "", err
}

func tryTspConfigPath(base string) Readme {
	return Readme(fmt.Sprintf("%s/tspconfig.yaml", base))
}

func setMapValue(cache map[Readme]Readme, value Readme) map[Readme]Readme {
	newCache := make(map[Readme]Readme, len(cache))
	for k, _ := range cache {
		newCache[k] = value
	}

	return newCache
}

func appendMap(dst, src map[Readme]Readme) {
	for k, v := range src {
		dst[k] = v
	}
}
