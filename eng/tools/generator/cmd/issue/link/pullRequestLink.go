// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package link

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
	"github.com/ahmetb/go-linq/v3"
	"github.com/google/go-github/v62/github"
)

type pullRequestLink struct {
	linkBase
}

// NewPullRequestLink parses a pull request link to its corresponding readme.md file link
func NewPullRequestLink(ctx context.Context, client *query.Client, requestLink, releaseLink string) Resolver {
	return &pullRequestLink{
		linkBase: linkBase{
			ctx:         ctx,
			client:      client,
			releaseLink: releaseLink,
			requestLink: requestLink,
		},
	}
}

// Resolve ...
func (l pullRequestLink) Resolve() (ResolveResult, error) {
	n, err := l.getPullRequestNumber()
	if err != nil {
		return nil, fmt.Errorf("cannot resolve pull request number from '%s'", l)
	}
	files, err := l.listChangedFiles(SpecOwner, SpecRepo, n)
	if err != nil {
		return nil, err
	}
	var filePaths []string
	linq.From(files).Select(func(item interface{}) interface{} {
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
		return nil, fmt.Errorf("cannot resolve pull request link '%s': %+v", l.GetReleaseLink(), err)
	}
	// we need to check if the associated PR has been merged
	merged, err := l.checkStatus(SpecOwner, SpecRepo, n)
	if err != nil {
		return nil, err
	}
	if !merged {
		return result{
			readme: readme,
			code:   CodePRNotMerged,
		}, nil
	}
	return getResult(readme), nil
}

// String ...
func (l pullRequestLink) String() string {
	return l.GetReleaseLink()
}

// getPullRequestNumber returns the PR number from a PR link which should be in this form: {number}(/something)?
func (l pullRequestLink) getPullRequestNumber() (int, error) {
	segments := strings.Split(strings.TrimPrefix(l.GetReleaseLink(), PullRequestPrefix), "/")
	return strconv.Atoi(segments[0])
}

func (l pullRequestLink) checkStatus(owner, repo string, number int) (bool, error) {
	merged, _, err := l.client.PullRequests.IsMerged(l.ctx, owner, repo, number)
	if err != nil {
		return false, err
	}
	return merged, nil
}

func (l pullRequestLink) listChangedFiles(owner, repo string, number int) ([]*github.CommitFile, error) {
	opt := &github.ListOptions{
		PerPage: 10,
	}
	var files []*github.CommitFile
	for {
		f, resp, err := l.client.PullRequests.ListFiles(l.ctx, owner, repo, number, opt)
		if err != nil {
			return nil, err
		}
		files = append(files, f...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}
	return files, nil
}
