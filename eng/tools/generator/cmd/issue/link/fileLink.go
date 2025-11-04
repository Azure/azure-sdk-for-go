// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package link

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/cmd/issue/query"
	"github.com/Azure/azure-sdk-for-go/eng/tools/generator/typespec"
)

type fileLink struct {
	linkBase

	path string
}

var _ CommitHashLink = (*fileLink)(nil)

// NewFileLink parses a file link to its corresponding readme.md file link
func NewFileLink(ctx context.Context, client *query.Client, requestLink, releaseLink string) Resolver {
	return &fileLink{
		linkBase: linkBase{
			ctx:         ctx,
			client:      client,
			releaseLink: releaseLink,
			requestLink: requestLink,
		},
	}
}

// Resolve ...
func (l fileLink) Resolve() (ResolveResult, error) {
	commitRef, err := l.GetCommitHash()
	if err != nil {
		return nil, err
	}
	l.path = strings.TrimPrefix(l.GetReleaseLink(), FilePrefix+commitRef+"/")

	tspConfig, err := GetTspConfigFromPath(l.ctx, l.client, l.path)
	if err != nil {
		if !errors.Is(err, errNoTspConfig) {
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

	readme, err := GetReadmeFromPath(l.ctx, l.client, l.path)
	if err != nil {
		return nil, fmt.Errorf("cannot resolve file link '%s': %+v", l.GetReleaseLink(), err)
	}
	return getResult(readme), nil
}

// String ...
func (l fileLink) String() string {
	return l.GetReleaseLink()
}

// GetCommitHash ...
func (l fileLink) GetCommitHash() (string, error) {
	return getCommitRefFromLink(l.GetReleaseLink(), FilePrefix)
}
