package automation

import "github.com/Azure/azure-sdk-for-go/tools/generator/autorest"

// existingPackageMap is a map with readme relative path as keys
type existingPackageMap map[string]packagesForReadme

// packagesForReadme is a map with tag as keys
type packagesForReadme map[string]existingGenerationMetadata

type existingGenerationMetadata struct {
	autorest.GenerationMetadata
	packageFullPath string
}

func (m existingPackageMap) add(path string, metadata autorest.GenerationMetadata) {
	if _, ok := m[metadata.RelativeReadme()]; !ok {
		m[metadata.RelativeReadme()] = packagesForReadme{}
	}

	m[metadata.RelativeReadme()].add(path, metadata)
}

func (m packagesForReadme) add(path string, metadata autorest.GenerationMetadata) {
	m[metadata.Tag] = existingGenerationMetadata{
		GenerationMetadata: metadata,
		packageFullPath:    path,
	}
}
