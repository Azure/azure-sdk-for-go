// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
)

// MetadataProcessError ...
type MetadataProcessError struct {
	MetadataLocation string
	Errors           []error
}

// Error ...
func (e *MetadataProcessError) Error() string {
	return fmt.Sprintf("total %d error(s) during processing metadata %s: %+v", len(e.Errors), e.MetadataLocation, e.Errors)
}

func (e *MetadataProcessError) add(err error) {
	e.Errors = append(e.Errors, err)
}

// MetadataProcessor processes the metadata
type MetadataProcessor struct {
	metadataOutputFolder string
}

// NewMetadataProcessorFromLocation creates a new MetadataProcessor using the metadata output folder location
func NewMetadataProcessorFromLocation(metadataOutput string) *MetadataProcessor {
	return &MetadataProcessor{
		metadataOutputFolder: metadataOutput,
	}
}

// Process returns the metadata result: a map from tag to Metadata, and an error if there is anything that could not be processed.
// the error returned must be of type *MetadataProcessError
func (p MetadataProcessor) Process() (map[string]model.Metadata, error) {
	fi, err := ioutil.ReadDir(p.metadataOutputFolder)
	if err != nil {
		return nil, &MetadataProcessError{
			MetadataLocation: p.metadataOutputFolder,
			Errors:           []error{err},
		}
	}
	result := make(map[string]model.Metadata)
	metadataErr := &MetadataProcessError{
		MetadataLocation: p.metadataOutputFolder,
	}
	for _, f := range fi {
		// a metadata output must be a json file
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".json") {
			continue
		}
		file, err := os.Open(filepath.Join(p.metadataOutputFolder, f.Name()))
		if err != nil {
			metadataErr.add(err)
			continue
		}
		metadata, err := model.NewMetadataFrom(file)
		if err != nil {
			metadataErr.add(err)
			continue
		}
		tag := strings.TrimSuffix(f.Name(), ".json")
		result[tag] = metadata
	}

	if len(metadataErr.Errors) != 0 {
		return result, metadataErr
	}
	return result, nil
}
