// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/apidiff/exports"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/utils"
	"github.com/Azure/azure-sdk-for-go/tools/generator/validate"
)

type changelogContext struct {
	sdkRoot         string
	clnRoot         string
	specRoot        string
	commitHash      string
	codeGenVer      string
	readme          string
	removedPackages []packageOutput

	repoContent map[string]exports.Content
}

func (ctx changelogContext) SDKRoot() string {
	return ctx.sdkRoot
}

func (ctx changelogContext) RepoContent() map[string]exports.Content {
	return ctx.repoContent
}

func (ctx changelogContext) process(metadataLocation string) ([]autorest.ChangelogResult, error) {
	// get the metadata
	m := autorest.NewMetadataProcessorFromLocation(metadataLocation)
	metadataMap, err := m.Process()
	if err != nil {
		return nil, err
	}
	// validate the metadata output-folder
	if err := ctx.validateMetadata(metadataMap); err != nil {
		return nil, err
	}
	// generate the changelogs
	p := autorest.NewChangelogProcessorFromContext(ctx)
	changelogResults, err := p.Process(metadataMap)
	if err != nil {
		return nil, err
	}
	for _, result := range changelogResults {
		// we need to write the changelog file to the corresponding package here
		if err := WriteChangelogFile(result); err != nil {
			return nil, err
		}
	}
	// iterate over the removed packages, generate changelogs for them as well
	var removedResults []autorest.ChangelogResult
	for _, rp := range ctx.removedPackages {
		if contains(changelogResults, rp.outputFolder) {
			// this package has been regenerated
			continue
		}
		result, err := p.GenerateChangelog(rp.outputFolder)
		if err != nil {
			return nil, err
		}
		removedResults = append(removedResults, *result)
	}
	changelogResults = append(changelogResults, removedResults...)

	// omit the packages not in services directory
	var results []autorest.ChangelogResult
	for _, result := range changelogResults {
		if strings.HasPrefix(result.PackageName, "services/") {
			results = append(results, result)
		}
	}

	return results, nil
}

func contains(array []autorest.ChangelogResult, item string) bool {
	for _, r := range array {
		if utils.NormalizePath(r.PackageFullPath) == utils.NormalizePath(item) {
			return true
		}
	}
	return false
}

func WriteChangelogFile(result autorest.ChangelogResult) error {
	fileContent := result.Write()
	
	changelogFile, err := os.Create(filepath.Join(result.PackageFullPath, autorest.ChangelogFilename))
	if err != nil {
		return err
	}
	defer changelogFile.Close()

	fileContent = fmt.Sprintf(`# Unreleased content

%s`, fileContent)

	if _, err := changelogFile.WriteString(fileContent); err != nil {
		return err
	}
	return nil
}

func (ctx changelogContext) validateMetadata(metadataMap map[string]model.Metadata) error {
	builder := validationErrorBuilder{
		readme: ctx.readme,
	}
	validateContext := validate.MetadataValidateContext{
		Readme:  ctx.readme,
		SDKRoot: ctx.sdkRoot,
		Validators: []validate.MetadataValidateFunc{
			validate.PreviewCheck,
			validate.MgmtCheck,
			validate.NamespaceCheck,
		},
	}
	for tag, metadata := range metadataMap {
		// validate the output-folder, etc
		if errors := validateContext.Validate(tag, metadata); len(errors) != 0 {
			builder.addMultiple(errors)
			continue
		}
	}
	return builder.build()
}

type validationErrorBuilder struct {
	readme string
	errors []error
}

func (b *validationErrorBuilder) addMultiple(errors []error) {
	b.errors = append(b.errors, errors...)
}

func (b *validationErrorBuilder) add(err error) {
	b.errors = append(b.errors, err)
}

func (b *validationErrorBuilder) build() error {
	if len(b.errors) == 0 {
		return nil
	}
	var messages []string
	for _, e := range b.errors {
		messages = append(messages, e.Error())
	}
	return fmt.Errorf("validation failed in readme '%s' with %d error(s): \n%s", b.readme, len(b.errors), strings.Join(messages, "\n"))
}
