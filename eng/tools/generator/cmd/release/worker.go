// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package release

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/autorest/model"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/repo"
	"github.com/Azure/azure-sdk-for-go/eng/tools/internal/exports"
)

type generateContext struct {
	sdkRepo            repo.SDKRepository
	specRepo           repo.SpecRepository
	readme             string
	specLastCommitHash string

	defaultOptions    model.Options
	additionalOptions []model.Option

	repoContent map[string]exports.Content

	skipProfiles bool
}

func (ctx generateContext) SDKRoot() string {
	return ctx.sdkRepo.Root()
}

func (ctx generateContext) SpecRoot() string {
	return ctx.specRepo.Root()
}

func (ctx generateContext) RepoContent() map[string]exports.Content {
	return ctx.repoContent
}

func (ctx *generateContext) generate(tag string, infoList []config.ReleaseRequestInfo) (*GenerateResult, error) {
	metadataOutputRoot := filepath.Join(os.TempDir(), fmt.Sprintf("release-metadata-%v", time.Now().Unix()))
	defer os.RemoveAll(metadataOutputRoot)

	var options model.Options
	// determine whether this is a new package or not
	if m, ok := repo.ContainsPackage(ctx.SDKRoot(), ctx.readme, tag); ok {
		log.Printf("Task (readme %s / tag %s) is an existing package, using the options in the metadata...", ctx.readme, tag)
		options = ctx.defaultOptions.(model.Options).MergeOptions(autorest.GetAdditionalOptions(m).Arguments()...)
	} else {
		log.Printf("Task (readme %s / tag %s) is a new package, appending the additional options to the default options...", ctx.readme, tag)
		options = ctx.defaultOptions.(model.Options).MergeOptions(ctx.additionalOptions...)
	}

	log.Printf("Generating %s from %v...", tag, infoList)
	// iterate over the tags in one request
	// Generate code
	input := autorest.GenerateInput{
		Readme: ctx.readme,
		Tag:    tag,
		//SDKRoot:    ctx.SDKRoot(),
		CommitHash: ctx.specLastCommitHash,
		Options:    options,
	}
	r, err := autorest.GeneratePackage(ctx, input, autorest.GenerateOptions{
		MetadataOutputRoot: metadataOutputRoot,
		Stderr:             os.Stderr,
		Stdout:             os.Stderr, // we redirect all the output of autorest to stderr, so that the stdout will only contain the proper output
	})
	if err != nil {
		return nil, fmt.Errorf("cannot generate readme '%s', tag '%s': %+v", ctx.readme, tag, err)
	}

	// regenerate the profiles
	if err := ctx.regenerateProfiles(); err != nil {
		return nil, err
	}

	// commit the content
	if err := ctx.commit(tag); err != nil {
		return nil, err
	}

	// get last commit
	ref, err := ctx.sdkRepo.Head()
	if err != nil {
		return nil, err
	}

	return &GenerateResult{
		Package:    r.Package,
		Readme:     ctx.readme,
		Tag:        tag,
		CommitHash: ref.Hash().String(),
		Info:       infoList,
	}, nil
}

func (ctx *generateContext) regenerateProfiles() error {
	if ctx.skipProfiles {
		return nil
	}
	return autorest.RegenerateProfiles(ctx.SDKRoot())
}

func (ctx *generateContext) commit(tag string) error {
	if err := ctx.sdkRepo.Add("profiles"); err != nil {
		return fmt.Errorf("failed to add 'profiles': %+v", err)
	}

	if err := ctx.sdkRepo.Add("services"); err != nil {
		return fmt.Errorf("failed to add 'services': %+v", err)
	}

	message := fmt.Sprintf("Generated from %s tag %s (commit hash: %s)", ctx.readme, tag, ctx.specLastCommitHash)
	if err := ctx.sdkRepo.Commit(message); err != nil {
		if repo.IsNothingToCommit(err) {
			log.Printf("There is nothing to commit. Message: %s", message)
			return nil
		}
		return fmt.Errorf("failed to commit changes: %+v", err)
	}

	return nil
}
