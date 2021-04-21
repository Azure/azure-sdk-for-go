// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
)

type cleanUpContext struct {
	root        string
	readmeFiles []string
}

// clean cleans all the packages related with the readme files. Return a map using readme filepath as keys, corresponds to the cleaned packages
func (ctx *cleanUpContext) clean() (readmePackageOutputMap, error) {
	log.Printf("Summarying all the generation metadata in '%s'...", ctx.root)
	m, err := summaryReadmePackageOutputMap(ctx.root)
	if err != nil {
		return nil, err
	}

	removedPackages := readmePackageOutputMap{}
	for _, readme := range ctx.readmeFiles {
		log.Printf("Cleaning up the packages generated from readme '%s'...", readme)
		for _, p := range m[readme] {
			if err := os.RemoveAll(p.outputFolder); err != nil {
				return nil, fmt.Errorf("cannot remove package '%s': %+v", p.outputFolder, err)
			}
			removedPackages.add(readme, p)

			// recursively remove all its parent if this directory is empty after the deletion
			if err := removeEmptyParents(filepath.Dir(p.outputFolder)); err != nil {
				return nil, err
			}
		}
	}
	return removedPackages, nil
}

func removeEmptyParents(parent string) error {
	fi, err := ioutil.ReadDir(parent)
	if err != nil {
		return err
	}
	if len(fi) == 0 {
		if err := os.RemoveAll(parent); err != nil {
			return err
		}
		return removeEmptyParents(filepath.Dir(parent))
	}
	return nil
}

func summaryReadmePackageOutputMap(root string) (readmePackageOutputMap, error) {
	// first we list all the go sdk track 1 packages
	m, err := autorest.CollectGenerationMetadata(root)
	if err != nil {
		return nil, err
	}
	result := readmePackageOutputMap{}
	for pkg, metadata := range m {
		result.add(metadata.RelativeReadme(), packageOutput{
			tag:          metadata.Tag,
			outputFolder: pkg,
		})
	}
	return result, nil
}

type readmePackageOutputMap map[string][]packageOutput

type packageOutput struct {
	tag          string
	outputFolder string
}

func (m *readmePackageOutputMap) add(readme string, output packageOutput) {
	if l, ok := (*m)[readme]; ok {
		(*m)[readme] = append(l, output)
	} else {
		(*m)[readme] = []packageOutput{output}
	}
}

func (m *readmePackageOutputMap) packages() []packageOutput {
	var results []packageOutput
	for _, l := range *m {
		results = append(results, l...)
	}
	return results
}
