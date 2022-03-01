//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azappconfig

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/appconfiguration/azappconfig/internal/generated"
)

// ListRevisionsPager is a Pager for revision list operations.
//
// NextPage should be called first. It fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on this Pager.
// If the result is false, the value of Err() will indicate if an error occurred.
type ListRevisionsPager interface {
	// PageResponse returns the current ListRevisionsPage.
	PageResponse() ListRevisionsPage

	// Err returns an error if there was an error on the last request.
	Err() error

	// NextPage returns true if there is another page of data available, false if not.
	NextPage(context.Context) bool
}

type listRevisionsPager struct {
	genPager *generated.AzureAppConfigurationClientGetRevisionsPager
}

func (p listRevisionsPager) PageResponse() ListRevisionsPage {
	return fromGeneratedGetRevisionsPage(p.genPager.PageResponse())
}

func (p listRevisionsPager) Err() error {
	return p.genPager.Err()
}

func (p listRevisionsPager) NextPage(ctx context.Context) bool {
	return p.genPager.NextPage(ctx)
}

// ListRevisionsOptions contains the optional parameters for the ListRevisions method.
type ListRevisionsOptions struct {
}

// ListRevisions retrieves the revisions of one or more configuration setting entities that match the specified setting selector.
func (c *Client) ListRevisions(selector SettingSelector, options *ListRevisionsOptions) ListRevisionsPager {
	_ = options
	return listRevisionsPager{genPager: c.appConfigClient.GetRevisions(selector.toGenerated())}
}
