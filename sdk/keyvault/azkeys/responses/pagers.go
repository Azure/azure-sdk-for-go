//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package responses

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/convert"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/models"
)

// ListKeysPager is a Pager for the Client.ListSecrets operation
type ListKeysPager struct {
	genPager *generated.KeyVaultClientGetKeysPager
}

func NewListKeysPager(p *generated.KeyVaultClientGetKeysPager) *ListKeysPager {
	return &ListKeysPager{genPager: p}
}

// PageResponse returns the results from the page most recently fetched from the service
func (l *ListKeysPager) PageResponse() ListKeysPage {
	return listKeysPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil
func (l *ListKeysPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *ListKeysPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListDeletedKeys is the interface for the Client.ListDeletedKeys operation
type ListDeletedKeysPager struct {
	genPager *generated.KeyVaultClientGetDeletedKeysPager
}

func NewListDeletedKeysPager(p *generated.KeyVaultClientGetDeletedKeysPager) *ListDeletedKeysPager {
	return &ListDeletedKeysPager{genPager: p}
}

// PageResponse returns the current page of results
func (l *ListDeletedKeysPager) PageResponse() ListDeletedKeysPage {
	resp := l.genPager.PageResponse()

	var values []*models.DeletedKeyItem
	for _, d := range resp.Value {
		values = append(values, convert.DeletedKeyItemFromGenerated(d))
	}

	return ListDeletedKeysPage{
		RawResponse: resp.RawResponse,
		NextLink:    resp.NextLink,
		DeletedKeys: values,
	}
}

// Err returns an error if the last operation resulted in an error.
func (l *ListDeletedKeysPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next page of results.
func (l *ListDeletedKeysPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListKeyVersionsPager is a Pager for Client.ListKeyVersions results
type ListKeyVersionsPager struct {
	genPager *generated.KeyVaultClientGetKeyVersionsPager
}

func NewListKeyVersionsPager(p *generated.KeyVaultClientGetKeyVersionsPager) *ListKeyVersionsPager {
	return &ListKeyVersionsPager{genPager: p}
}

// PageResponse returns the results from the page most recently fetched from the service.
func (l *ListKeyVersionsPager) PageResponse() ListKeyVersionsPage {
	return listKeyVersionsPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (l *ListKeyVersionsPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *ListKeyVersionsPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// convert internal Response to ListKeysPage
func listKeysPageFromGenerated(i generated.KeyVaultClientGetKeysResponse) ListKeysPage {
	var keys []*models.KeyItem
	for _, k := range i.Value {
		keys = append(keys, convert.KeyItemFromGenerated(k))
	}
	return ListKeysPage{
		RawResponse: i.RawResponse,
		NextLink:    i.NextLink,
		Keys:        keys,
	}
}

// create ListKeysPage from generated pager
func listKeyVersionsPageFromGenerated(i generated.KeyVaultClientGetKeyVersionsResponse) ListKeyVersionsPage {
	var keys []models.KeyItem
	for _, s := range i.Value {
		if s != nil {
			keys = append(keys, *convert.KeyItemFromGenerated(s))
		}
	}
	return ListKeyVersionsPage{
		RawResponse: i.RawResponse,
		NextLink:    i.NextLink,
		Keys:        keys,
	}
}
