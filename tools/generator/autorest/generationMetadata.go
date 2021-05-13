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
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/internal/packages/track1"
)

// GenerationMetadata contains all the metadata that has been used when generating a track 1 package
type GenerationMetadata struct {
	// AutorestVersion is the version of autorest.core
	AutorestVersion string `json:"autorest,omitempty"`
	// CommitHash is the commit hash of azure-rest-api-specs from which this SDK package is generated
	CommitHash string `json:"commit,omitempty"`
	// Readme is the normalized path of the readme file from which this SDK package is generated. It should be in this pattern: /_/azure-rest-api-specs/{relative_path}
	Readme string `json:"readme,omitempty"`
	// Tag is the tag from which this SDK package is generated
	Tag string `json:"tag,omitempty"`
	// CodeGenVersion is the version of autorest.go using when this package is generated
	CodeGenVersion string `json:"use,omitempty"`
	// RepositoryURL is the URL of the azure-rest-api-specs. This should always be a constant "https://github.com/Azure/azure-rest-api-specs.git"
	RepositoryURL string `json:"repository_url,omitempty"`
	// AutorestCommand is the full command that generates this package
	AutorestCommand string `json:"autorest_command,omitempty"`
	// AdditionalProperties is a map of addition information in this metadata
	AdditionalProperties GenerationMetadataAdditionalProperties `json:"additional_properties,omitempty"`
}

// GenerationMetadataAdditionalProperties contains all the additional options other than go-sdk-foler, tag, multiapi, use or the readme path
type GenerationMetadataAdditionalProperties struct {
	AdditionalOptions string `json:"additional_options,omitempty"`
}

// RelativeReadme returns the relative readme path
func (meta *GenerationMetadata) RelativeReadme() string {
	return strings.TrimPrefix(meta.Readme, NormalizedSpecRoot)
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
		if m != nil {
			result[pkg.FullPath()] = *m
		}
	}
	return result, nil
}

// GetGenerationMetadata gets the GenerationMetadata in one specific package
func GetGenerationMetadata(pkg track1.Package) (*GenerationMetadata, error) {
	metadataFilepath := filepath.Join(pkg.FullPath(), MetadataFilename)

	// some classical package might not have a changelog, therefore we need to identify whether the changelog file exist or not
	if _, err := os.Stat(metadataFilepath); os.IsNotExist(err) {
		log.Printf("package '%s' does not have a metadata file", pkg.Path())
		return nil, nil
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

// AdditionalOptions removes flags that may change over scenarios
func AdditionalOptions(arguments []model.Option) []model.Option {
	var transformed []model.Option
	for _, argument := range arguments {
		switch o := argument.(type) {
		case model.ArgumentOption: // omit the readme path argument
			continue
		case model.FlagOption:
			if o.Flag() == "multiapi" { // omit the multiapi flag or use
				continue
			}
		case model.KeyValueOption:
			// omit go-sdk-folder, use, tag and metadata-output-folder
			if o.Key() == "go-sdk-folder" || o.Key() == "use" || o.Key() == "tag" || o.Key() == "metadata-output-folder" {
				continue
			}
		}
		transformed = append(transformed, argument)
	}
	return transformed
}

// AdditionalOptionsToString removes flags that may change over scenarios and cast them to strings
func AdditionalOptionsToString(arguments []model.Option) []string {
	transformed := AdditionalOptions(arguments)
	result := make([]string, len(transformed))
	for i, o := range transformed {
		result[i] = o.Format()
	}
	return result
}

const (
	// NormalizedSpecRoot this is the prefix for readme
	NormalizedSpecRoot = "/_/azure-rest-api-specs/"

	// NormalizedSDKRoot this is the prefix for readme
	NormalizedSDKRoot = "/_/azure-sdk-for-go/"

	// MetadataFilename ...
	MetadataFilename = "_meta.json"
)
