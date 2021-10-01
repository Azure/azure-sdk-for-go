// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package automation

import "github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"

// existingPackageMap is a map with readme relative path as keys
type existingPackageMap map[string]packagesForReadme

// packagesForReadme is a map with tag as keys
type packagesForReadme map[string]existingGenerationMetadata

type existingGenerationMetadata struct {
	autorest.GenerationMetadata
	// packageName is the relative path of a package to the root of the SDK
	packageName string
}

func (m existingPackageMap) add(relPath string, metadata autorest.GenerationMetadata) {
	if _, ok := m[metadata.RelativeReadme()]; !ok {
		m[metadata.RelativeReadme()] = packagesForReadme{}
	}

	m[metadata.RelativeReadme()].add(relPath, metadata)
}

func (m packagesForReadme) add(relPath string, metadata autorest.GenerationMetadata) {
	m[metadata.Tag] = existingGenerationMetadata{
		GenerationMetadata: metadata,
		packageName:        relPath,
	}
}
