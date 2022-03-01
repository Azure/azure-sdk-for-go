//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfig/internal/generated"
)

// GetRevisionsPager is a Pager for revision list operations.
//
// NextPage should be called first. It fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on this Pager.
// If the result is false, the value of Err() will indicate if an error occurred.
type GetRevisionsPager interface {
	// PageResponse returns the current GetRevisionsPage.
	PageResponse() GetRevisionsPage

	// Err returns an error if there was an error on the last request.
	Err() error

	// NextPage returns true if there is another page of data available, false if not.
	NextPage(context.Context) bool
}

type getRevisionsPager struct {
	genPager *generated.AzureAppConfigurationClientGetRevisionsPager
}

func (p getRevisionsPager) PageResponse() GetRevisionsPage {
	return fromGeneratedGetRevisionsPage(p.genPager.PageResponse())
}

func (p getRevisionsPager) Err() error {
	return p.genPager.Err()
}

func (p getRevisionsPager) NextPage(ctx context.Context) bool {
	return p.genPager.NextPage(ctx)
}

// GetRevisionsOptions contains the optional parameters for the GetRevisions method.
type GetRevisionsOptions struct {
}

// GetRevisions retrieves the revisions of one or more configuration setting entities that match the specified setting selector.
func (c *Client) GetRevisions(selector SettingSelector, options *GetRevisionsOptions) GetRevisionsPager {
	_ = options
	return getRevisionsPager{genPager: c.appConfigClient.GetRevisions(selector.toGenerated())}
}
