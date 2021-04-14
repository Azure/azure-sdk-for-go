// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/tools/pkgchk/track1"
)

// GenerationMetadata contains all the metadata that has been used when generating a track 1 package
type GenerationMetadata struct {
	AutorestVersion      string                 `json:"autorest,omitempty"`
	CommitHash           string                 `json:"commit,omitempty"`
	Readme               string                 `json:"readme,omitempty"`
	Tag                  string                 `json:"tag,omitempty"`
	CodeGenVersion       string                 `json:"use,omitempty"`
	RepositoryURL        string                 `json:"repository_url,omitempty"`
	AutorestCommand      string                 `json:"autorest_command,omitempty"`
	AdditionalProperties map[string]interface{} `json:"additional_properties,omitempty"`
}

// CollectGenerationMetadata iterates every track 1 go sdk package under root, and collect all the GenerationMetadata into a map
// using relative path of the package as keys
func CollectGenerationMetadata(root string) (map[string]GenerationMetadata, error) {
	pkgs, err := track1.List(root)
	if err != nil {
		return nil, fmt.Errorf("failed to get track 1 package list under root '%s': %+v", root, err)
	}
	result := make(map[string]GenerationMetadata)
	for _, pkg := range pkgs {
		m, err := GetGenerationMetadata(pkg)
		if err != nil {
			return nil, err
		}
		result[pkg.FullPath()] = *m
	}
	return result, nil
}

// GetGenerationMetadata gets the GenerationMetadata in one specific package
func GetGenerationMetadata(pkg track1.Package) (*GenerationMetadata, error) {
	metadataFilepath := filepath.Join(pkg.FullPath(), MetadataFilename)

	// some classical package might not have a changelog, therefore we need to identify whether the changelog file exist or not
	if _, err := os.Stat(metadataFilepath); os.IsNotExist(err) {
		log.Printf("package '%s' does not have a metadata file", pkg.Path())
		return &GenerationMetadata{}, nil
	}

	b, err := ioutil.ReadFile(metadataFilepath)
	if err != nil {
		return nil, fmt.Errorf("cannot read file %s: %+v", metadataFilepath, err)
	}

	var metadata GenerationMetadata
	if err := json.Unmarshal(b, &metadata); err != nil {
		return nil, fmt.Errorf("cannot unmarshal metadata: %+v", err)
	}

	return &metadata, nil
}

const (
	// MetadataFilename
	MetadataFilename = "_meta.json"
)
