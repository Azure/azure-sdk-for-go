// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/tools/generator/autorest/model"
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
}

func (ctx changelogContext) SDKRoot() string {
	return ctx.sdkRoot
}

func (ctx changelogContext) SDKCloneRoot() string {
	return ctx.clnRoot
}

func (ctx changelogContext) SpecRoot() string {
	return ctx.specRoot
}

func (ctx changelogContext) SpecCommitHash() string {
	return ctx.commitHash
}

func (ctx changelogContext) CodeGenVersion() string {
	return ctx.codeGenVer
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
	p := autorest.NewChangelogProcessorFromContext(ctx).WithLocation(metadataLocation).WithReadme(ctx.readme)
	changelogResults, err := p.Process(metadataMap)
	if err != nil {
		return nil, err
	}
	for _, result := range changelogResults {
		// we need to write the changelog file to the corresponding package here
		if err := writeChangelogFile(result); err != nil {
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
	return changelogResults, nil
}

func contains(array []autorest.ChangelogResult, item string) bool {
	for _, r := range array {
		if r.PackagePath == item {
			return true
		}
	}
	return false
}

func writeChangelogFile(result autorest.ChangelogResult) error {
	fileContent := fmt.Sprintf(`%s

%s`, result.GenerationMetadata.String(), result.Changelog.ToMarkdown())
	changelogFile, err := os.Create(filepath.Join(result.PackagePath, autorest.ChangelogFilename))
	if err != nil {
		return err
	}
	defer changelogFile.Close()
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
			// TODO -- we do have some exceptions (see file tools/pkgchk/exceptions.txt) that might need to be considered here
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
