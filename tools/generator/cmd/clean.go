package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/tools/generator/changelog"
	"github.com/Azure/azure-sdk-for-go/tools/generator/utils"
)

type cleanUpContext struct {
	root        string
	readmeFiles []string
}

func (ctx *cleanUpContext) clean() ([]PackageOutput, error) {
	log.Printf("Summarying all the generation metadata in '%s'...", ctx.root)
	m, err := summaryReadmePackageOutputMap(ctx.root)
	if err != nil {
		return nil, err
	}

	var removedPackages []PackageOutput
	for _, readme := range ctx.readmeFiles {
		log.Printf("Cleaning up the packages generated from readme '%s'...", readme)
		for _, p := range m[readme] {
			if err := os.RemoveAll(p.OutputFolder); err != nil {
				return nil, fmt.Errorf("cannot remove package '%s': %+v", p.OutputFolder, err)
			}
			removedPackages = append(removedPackages, p)
		}
	}
	return removedPackages, nil
}

func summaryReadmePackageOutputMap(root string) (ReadmePackageOutputMap, error) {
	// first we list all the go sdk track 1 packages
	m, err := listGenerationMetadata(root)
	if err != nil {
		return nil, err
	}
	result := ReadmePackageOutputMap{}
	for pkg, metadata := range m {
		result.add(metadata.Readme, PackageOutput{
			Tag:          metadata.Tag,
			OutputFolder: pkg,
		})
	}
	return result, nil
}

// ReadmePackageOutputMap is a map with key of readme relative path (starts with `specification`) and values of the corresponding package path info
type ReadmePackageOutputMap map[string][]PackageOutput

// PackageOutput contains the output folder and corresponding tag
type PackageOutput struct {
	Tag          string
	OutputFolder string
}

func (m *ReadmePackageOutputMap) add(readme string, output PackageOutput) {
	if l, ok := (*m)[readme]; ok {
		(*m)[readme] = append(l, output)
	} else {
		(*m)[readme] = []PackageOutput{output}
	}
}

func listGenerationMetadata(root string) (map[string]changelog.GenerationMetadata, error) {
	pkgs, err := utils.ListTrack1SDKPackages(root)
	if err != nil {
		return nil, fmt.Errorf("failed to get track 1 package list under root '%s': %+v", root, err)
	}
	result := make(map[string]changelog.GenerationMetadata)
	for _, pkg := range pkgs {
		m, err := getGenerationMetadata(pkg)
		if err != nil {
			return nil, err
		}
		result[pkg] = *m
	}
	return result, nil
}

func getGenerationMetadata(pkg string) (*changelog.GenerationMetadata, error) {
	changelogPath := filepath.Join(pkg, ChangelogFileName)
	file, err := os.Open(changelogPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %+v", changelogPath, err)
	}
	defer file.Close()
	return changelog.Parse(file)
}

const (
	ChangelogFileName = "CHANGELOG.md"
)
