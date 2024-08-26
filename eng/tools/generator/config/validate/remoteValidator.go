// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package validate

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/config"
	"github.com/ahmetb/go-linq/v3"
	"github.com/google/go-github/v62/github"
)

type remoteValidator struct {
	ctx    context.Context
	client *query.Client
}

func (v *remoteValidator) Validate(cfg config.Config) error {
	var errResult error
	for readme, infoMap := range cfg.Track2Requests {
		// first we validate whether the readme exists
		file, err := v.validateReadmeExistence(readme)
		if err != nil {
			errResult = errors.Join(errResult, err)
			continue // readme file does not exist, we could just skip all of the other steps of validations
		}
		// get content of the readme
		contentOfReadme, err := file.GetContent()
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("cannot get readme.md content: %+v", err))
			continue
		}
		// validate the existence of readme.go.md
		fileGo, err := v.validateReadmeExistence(getReadmeGoFromReadme(readme))
		if err != nil {
			errResult = errors.Join(errResult, err)
			continue // readme.go.md is mandatory
		}
		// get content of the readme.go.md
		contentOfReadmeGo, err := fileGo.GetContent()
		if err != nil {
			errResult = errors.Join(errResult, fmt.Errorf("cannot get readme.go.md content: %+v", err))
			continue
		}
		// get the keys from infoMap, which is the tags
		var tags []string
		linq.From(infoMap).Select(func(item interface{}) interface{} {
			return item.(config.Track2Request).PackageFlag
		}).ToSlice(&tags)
		// check the tags one by one
		if err := validateTagsInReadme([]byte(contentOfReadme), readme, tags...); err != nil {
			errResult = errors.Join(errResult, err)
		}
		// check module-name exist
		if err := validateModuleNameInReadmeGo([]byte(contentOfReadmeGo), readme); err != nil {
			errResult = errors.Join(errResult, err)
		}
	}
	return errResult
}

func (v *remoteValidator) validateReadmeExistence(readme string) (*github.RepositoryContent, error) {
	file, _, _, err := v.client.Repositories.GetContents(v.ctx, SpecOwner, SpecRepo, readme, nil)
	if err != nil {
		return nil, fmt.Errorf("cannot get readme file '%s' on remote: %+v", readme, err)
	}
	return file, nil
}

func validateTagsInReadme(content []byte, readme string, tags ...string) error {
	var notFoundTags []string
	for _, t := range tags {
		if !findTagInReadme(content, t) {
			notFoundTags = append(notFoundTags, t)
		}
	}

	if len(notFoundTags) > 0 {
		return fmt.Errorf("%d tag(s) not defined in readme.md '%s': %s", len(notFoundTags), readme, strings.Join(notFoundTags, ", "))
	}

	return nil
}

func validateModuleNameInReadmeGo(content []byte, readme string) error {
	moduleExist := regexp.MustCompile(goReadmeModuleName).Match(content)
	if !moduleExist {
		return fmt.Errorf("module name not found in readme.go.md '%s'", getReadmeGoFromReadme(readme))
	}

	return nil
}

const (
	SpecOwner = "Azure"
	SpecRepo  = "azure-rest-api-specs"
)
