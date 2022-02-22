//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfiguration

import (
	"context"

	"sdk/appconfiguration/azappconfiguration/internal/generated"
)

type GetRevisionsPager interface {
	PageResponse() GetRevisionsPage
	Err() error
	NextPage(context.Context) bool
}

type getRevisionsPager struct {
	genPager *generated.AzureAppConfigurationClientGetRevisionsPager
}

func (l *getRevisionsPager) PageResponse() GetRevisionsPage {
	return getRevisionsPageFromGenerated(l.genPager.PageResponse())
}

func (l *getRevisionsPager) Err() error {
	return l.genPager.Err()
}

func (l *getRevisionsPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

type GetRevisionsOptions struct {
}

func (c *Client) GetRevisions(selector SettingSelector, options *GetRevisionsOptions) GetRevisionsPager {
	_ = options
	return c.appConfigClient.GetRevisions(selector.toGenerated())
}
