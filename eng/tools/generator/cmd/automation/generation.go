// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package automation

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/automation/validate"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/utils"
)

type generateContext struct {
	sdkRoot    string
	specRoot   string
	commitHash string

	repoContent map[string]exports.Content

	existingPackages  packagesForReadme
	defaultOptions    model.Options
	additionalOptions []model.Option
}

func (ctx generateContext) SDKRoot() string {
	return ctx.sdkRoot
}

func (ctx generateContext) SpecRoot() string {
	return ctx.specRoot
}

func (ctx generateContext) RepoContent() map[string]exports.Content {
	return ctx.repoContent
}

var _ autorest.GenerateContext = (*generateContext)(nil)

func (ctx generateContext) generate(readme string) ([]autorest.GenerateResult, []error) {
	absReadme := filepath.Join(ctx.specRoot, readme)
	absReadmeGo := filepath.Join(filepath.Dir(absReadme), "readme.go.md")
	log.Printf("Reading tags from readme.go.md '%s'...", absReadmeGo)
	reader, err := os.Open(absReadmeGo)
	if err != nil {
		return nil, []error{
			fmt.Errorf("cannot read from readme.go.md: %+v", err),
		}
	}
	log.Printf("Parsing tags from readme.go.md '%s'...", absReadmeGo)
	tags, err := autorest.ReadBatchTags(reader)
	if err != nil {
		return nil, []error{
			fmt.Errorf("cannot read batch tags in readme.go.md '%s': %+v", absReadmeGo, err),
		}
	}

	log.Printf("Cleaning all the packages from readme '%s'...", readme)
	removedPackages, err := clean(ctx.sdkRoot, ctx.existingPackages)
	if err != nil {
		return nil, []error{
			fmt.Errorf("cannot clean packages from readme '%s': %+v", readme, err),
		}
	}

	log.Printf("Generating the following tags: \n[%s]", strings.Join(tags, ", "))
	var packageResults []autorest.GenerateResult
	var errors []error
	for _, tag := range tags {
		result, err := ctx.generateForTag(readme, tag)
		if err != nil {
			errors = append(errors, err)
			continue
		}
		packageResults = append(packageResults, *result)
	}

	// also add the removed packages in the results if it is not regenerated
	for _, removedPackage := range removedPackages {
		if !contains(packageResults, removedPackage.packageName) {
			// this package is not regenerated, therefore it is removed
			packageResults = append(packageResults, autorest.GenerateResult{
				Package: autorest.ChangelogResult{
					Tag:             removedPackage.Tag,
					PackageName:     removedPackage.packageName,
					PackageFullPath: utils.NormalizePath(filepath.Join(ctx.sdkRoot, removedPackage.packageName)),
					Changelog: model.Changelog{
						RemovedPackage: true,
					},
				},
			})
		}
	}

	return packageResults, errors
}

func (ctx generateContext) generateForTag(readme, tag string) (*autorest.GenerateResult, error) {
	var options model.Options
	// Get the proper options to use depending on whether this tag has been already generated in the SDK or not
	if metadata, ok := ctx.existingPackages[tag]; ok {
		// this tag has been generated, use the existing parameters in its metadata
		additionalOptions, err := model.ParseOptions(strings.Split(metadata.AdditionalProperties.AdditionalOptions, " "))
		if err != nil {
			return nil, fmt.Errorf("cannot parse existing defaultOptions for readme '%s'/tag '%s': %+v", readme, tag, err)
		}
		options = ctx.defaultOptions.MergeOptions(additionalOptions.Arguments()...)
	} else {
		// this is a new tag
		options = ctx.defaultOptions.MergeOptions(ctx.additionalOptions...)
	}

	input := autorest.GenerateInput{
		Readme:     readme,
		Tag:        tag,
		CommitHash: ctx.commitHash,
		Options:    options,
	}
	validateCtx := validate.MetadataValidateContext{
		Readme:  readme,
		SDKRoot: ctx.sdkRoot,
	}
	return autorest.GeneratePackage(ctx, input, autorest.GenerateOptions{
		Stderr:            os.Stderr,
		Stdout:            os.Stdout,
		AutoRestLogPrefix: "[AUTOREST] ",
		ChangelogTitle:    "Unreleased",
		Validators: []autorest.MetadataValidateFunc{
			validateCtx.PreviewCheck,
			validateCtx.MgmtCheck,
			validateCtx.NamespaceCheck,
		},
	})
}
