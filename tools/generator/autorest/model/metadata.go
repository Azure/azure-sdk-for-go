// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package model

import (
	"encoding/json"
	"io"
	"io/ioutil"
)

// Metadata ...
type Metadata interface {
	// SwaggerFiles returns the related swagger files in this tag
	SwaggerFiles() []string
	// PackagePath returns the output package path of this tag
	PackagePath() string
	// Namespace returns the namespace of this tag
	Namespace() string
}

// NewMetadataFrom reads a new Metadata from a io.Reader
func NewMetadataFrom(reader io.Reader) (Metadata, error) {
	b, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	var result localMetadata
	if err := json.Unmarshal(b, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

type localMetadata struct {
	// InputFiles ...
	InputFiles []string `json:"inputFiles"`
	// OutputFolder ...
	OutputFolder string `json:"outputFolder"`
	// LocalNamespace ...
	LocalNamespace string `json:"namespace"`
}

// SwaggerFiles ...
func (m localMetadata) SwaggerFiles() []string {
	return m.InputFiles
}

// PackagePath ...
func (m localMetadata) PackagePath() string {
	return m.OutputFolder
}

// Namespace ...
func (m localMetadata) Namespace() string {
	return m.LocalNamespace
}

// String ...
func (m localMetadata) String() string {
	b, _ := json.MarshalIndent(m, "", "  ")
	return string(b)
}
