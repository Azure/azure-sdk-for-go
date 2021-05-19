// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/tools/generator/validate"
	"github.com/Azure/azure-sdk-for-go/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/tools/internal/utils"
)

type changelogContext struct {
	sdkRoot string
	readme  string

	removedPackages []packageOutput

	repoContent map[string]exports.Content

	commonMetadata autorest.GenerationMetadata

	autorestArguments []model.Option
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
		// we need to write the generation metadata to the corresponding package here
		metadata := ctx.commonMetadata
		metadata.Tag = result.Tag
		options := autorest.AdditionalOptionsToString(ctx.autorestArguments)
		metadata.AdditionalProperties = autorest.GenerationMetadataAdditionalProperties{
			AdditionalOptions: strings.Join(options, " "),
		}
		metadata.AutorestCommand = fmt.Sprintf("autorest --use=%s --tag=%s --go-sdk-folder=/_/azure-sdk-for-go %s /_/azure-rest-api-specs/%s",
			metadata.CodeGenVersion, result.Tag, strings.Join(options, " "), utils.NormalizePath(ctx.readme))
		if err := WriteGenerationMetadata(result.PackageFullPath, metadata); err != nil {
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
		result, err := p.GenerateChangelog(rp.outputFolder, rp.tag)
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

// WriteGenerationMetadata writes the metadata to _meta.json file
func WriteGenerationMetadata(path string, metadata autorest.GenerationMetadata) error {
	b, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("cannot marshal metadata: %+v", err)
	}

	metadataFile, err := os.Create(filepath.Join(path, autorest.MetadataFilename))
	if err != nil {
		return err
	}
	defer metadataFile.Close()

	if _, err := metadataFile.Write(b); err != nil {
		return err
	}
	return nil
}

// WriteChangelogFile writes the changelog to CHANGELOG.md
func WriteChangelogFile(result autorest.ChangelogResult) error {
	changelogFile, err := os.Create(filepath.Join(result.PackageFullPath, autorest.ChangelogFilename))
	if err != nil {
		return err
	}
	defer changelogFile.Close()

	if _, err := changelogFile.WriteString(`# Unreleased Content

`); err != nil {
		return err
	}

	return result.Write(changelogFile)
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
