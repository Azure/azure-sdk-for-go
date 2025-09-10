// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package link

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
	"github.com/ahmetb/go-linq/v3"
	"github.com/google/go-github/v62/github"
)

var (
	splitRegex = regexp.MustCompile(`[^A-Fa-f0-9]`)
)

type commitLink struct {
	linkBase

	rawLink string
}

var _ CommitHashLink = (*commitLink)(nil)

// NewCommitLink parses a commit link to its corresponding readme.md file link
func NewCommitLink(ctx context.Context, client *query.Client, requestLink, releaseLink string) Resolver {
	segments := splitRegex.Split(strings.TrimPrefix(releaseLink, CommitPrefix), -1)
	realLink := fmt.Sprintf("%s%s", CommitPrefix, segments[0])
	return &commitLink{
		linkBase: linkBase{
			ctx:         ctx,
			client:      client,
			releaseLink: realLink,
			requestLink: requestLink,
		},
		rawLink: releaseLink,
	}
}

// Resolve ...
func (l commitLink) Resolve() (ResolveResult, error) {
	hash, err := l.GetCommitHash()
	if err != nil {
		return nil, err
	}
	commit, _, err := l.client.Repositories.GetCommit(l.ctx, SpecOwner, SpecRepo, hash, nil)
	if err != nil {
		return nil, err
	}
	var filePaths []string
	linq.From(commit.Files).Select(func(item interface{}) interface{} {
		return item.(*github.CommitFile).GetFilename()
	}).ToSlice(&filePaths)

	tspConfig, err := GetTspConfigPathFromChangedFiles(l.ctx, l.client, filePaths)
	if err != nil {
		if !strings.Contains(err.Error(), "cannot get any tspconfig files from these changed files:") {
			return nil, err
		}
	} else if tspConfig != "" {
		remoteTspConfigPath := fmt.Sprintf("https://raw.githubusercontent.com/Azure/azure-rest-api-specs/main/%s", tspConfig)
		exist, err := typespec.ExistGoConfigInTspConfig(remoteTspConfigPath)
		if err != nil {
			return nil, err
		}
		if exist {
			return getResult(tspConfig), nil
		}
	}

	readme, err := GetReadmePathFromChangedFiles(l.ctx, l.client, filePaths)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve commit link '%s': %+v", l.GetReleaseLink(), err)
	}
	return getResult(readme), nil
}

// GetCommitHash ...
func (l commitLink) GetCommitHash() (string, error) {
	return getCommitRefFromLink(l.GetReleaseLink(), CommitPrefix)
}

// String ...
func (l commitLink) String() string {
	return l.GetReleaseLink()
}
