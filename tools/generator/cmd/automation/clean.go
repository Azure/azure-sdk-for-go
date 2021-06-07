// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package automation

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

func clean(sdkRoot string, packages packagesForReadme) ([]existingGenerationMetadata, error) {
	var removedPackages []existingGenerationMetadata

	for _, metadata := range packages {
		log.Printf("Cleaning up pakcage '%s'...", metadata.packageName)
		packageFullPath := filepath.Join(sdkRoot, metadata.packageName)
		if err := os.RemoveAll(packageFullPath); err != nil {
			return nil, fmt.Errorf("cannot remove package '%s': %+v", metadata.packageName, err)
		}

		removedPackages = append(removedPackages, metadata)

		// recursively remove all its parent if this directory is empty after the deletion
		if err := removeEmptyParents(filepath.Dir(packageFullPath)); err != nil {
			return nil, err
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
