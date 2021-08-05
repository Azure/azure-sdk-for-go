// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package autorest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/tools/internal/report"
	"github.com/Azure/azure-sdk-for-go/tools/internal/utils"
)

// ChangelogContext describes all necessary data that would be needed in the processing of changelogs
type ChangelogContext interface {
	SDKRoot() string
	RepoContent() map[string]exports.Content
}

// ChangelogProcessor processes the metadata and output changelog with the desired format
type ChangelogProcessor struct {
	ctx ChangelogContext
}

// NewChangelogProcessorFromContext returns a new ChangelogProcessor
func NewChangelogProcessorFromContext(ctx ChangelogContext) *ChangelogProcessor {
	return &ChangelogProcessor{
		ctx: ctx,
	}
}

// ChangelogResult describes the result of the generated changelog for one package
type ChangelogResult struct {
	Tag             string
	PackageName     string
	PackageFullPath string
	Changelog       model.Changelog
}

// ChangelogProcessError describes the errors during the processing
type ChangelogProcessError struct {
	Errors []error
}

// Error ...
func (e *ChangelogProcessError) Error() string {
	return fmt.Sprintf("total %d error(s) during processing changelog: %+v", len(e.Errors), e.Errors)
}

type changelogErrorBuilder struct {
	errors []error
}

func (b *changelogErrorBuilder) add(err error) {
	b.errors = append(b.errors, err)
}

func (b *changelogErrorBuilder) build() error {
	if len(b.errors) == 0 {
		return nil
	}
	return &ChangelogProcessError{
		Errors: b.errors,
	}
}

// Process generates the changelogs using the input metadata map.
// Please ensure the input metadata map does not contain any package that is not under the sdk root, otherwise this might give weird results.
func (p *ChangelogProcessor) Process(metadataMap map[string]model.Metadata) ([]ChangelogResult, error) {
	builder := changelogErrorBuilder{}
	var results []ChangelogResult
	for tag, metadata := range metadataMap {
		outputFolder := filepath.Clean(metadata.PackagePath())
		// format the package before getting the changelog
		if err := FormatPackage(outputFolder); err != nil {
			builder.add(err)
			continue
		}
		result, err := p.GenerateChangelog(outputFolder, tag)
		if err != nil {
			builder.add(err)
			continue
		}
		results = append(results, *result)
	}
	return results, builder.build()
}

// GenerateChangelog generates a changelog for one package
func (p *ChangelogProcessor) GenerateChangelog(packageFullPath, tag string) (*ChangelogResult, error) {
	relativePackagePath, err := p.getRelativePackagePath(packageFullPath)
	if err != nil {
		return nil, err
	}
	lhs := p.getLHSContent(relativePackagePath)
	rhs, err := getExportsForPackage(packageFullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to get exports from package '%s' in the sdk '%s': %+v", relativePackagePath, p.ctx.SDKRoot(), err)
	}
	r, err := GetChangelogForPackage(lhs, rhs)
	if err != nil {
		return nil, fmt.Errorf("failed to generate changelog for package '%s': %+v", relativePackagePath, err)
	}
	return &ChangelogResult{
		Tag:             tag,
		PackageName:     relativePackagePath,
		PackageFullPath: packageFullPath,
		Changelog:       *r,
	}, nil
}

func (p *ChangelogProcessor) getRelativePackagePath(fullPath string) (string, error) {
	relative, err := filepath.Rel(p.ctx.SDKRoot(), fullPath)
	if err != nil {
		return "", err
	}
	return utils.NormalizePath(relative), nil
}

func (p *ChangelogProcessor) getLHSContent(relativePackagePath string) *exports.Content {
	if v, ok := p.ctx.RepoContent()[relativePackagePath]; ok {
		return &v
	}
	return nil
}

func getExportsForPackage(dir string) (*exports.Content, error) {
	// The function exports.Get does not handle the circumstance that the package does not exist
	// therefore we have to check if it exists and if not exit early to ensure we do not return an error
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, nil
	}
	exp, err := exports.Get(dir)
	if err != nil {
		return nil, err
	}
	return &exp, nil
}

// GetChangelogForPackage generates the changelog report with the given two Contents
func GetChangelogForPackage(lhs, rhs *exports.Content) (*model.Changelog, error) {
	if lhs == nil && rhs == nil {
		return nil, fmt.Errorf("this package does not exist even after the generation, this should never happen")
	}
	if lhs == nil {
		// the package does not exist before the generation: this is a new package
		return &model.Changelog{
			NewPackage: true,
		}, nil
	}
	if rhs == nil {
		// the package no longer exists after the generation: this package was removed
		return &model.Changelog{
			RemovedPackage: true,
		}, nil
	}
	// lhs and rhs are both non-nil
	p := report.Generate(*lhs, *rhs, nil)
	return &model.Changelog{
		Modified: &p,
	}, nil
}

// ToMarkdown convert this ChangelogResult to markdown string
func (r ChangelogResult) ToMarkdown() string {
	return r.Changelog.ToMarkdown()
}

// Write writes this ChangelogResult to the specified writer in markdown format
func (r ChangelogResult) Write(writer io.Writer) error {
	_, err := writer.Write([]byte(r.ToMarkdown()))
	return err
}

const (
	// ChangelogFilename ...
	ChangelogFilename = "CHANGELOG.md"
)
