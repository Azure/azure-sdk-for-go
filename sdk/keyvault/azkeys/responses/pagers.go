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

// More returns true if there are more pages to retrieve.
func (l *ListKeysPager) More() bool {
	return l.genPager.More()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *ListKeysPager) NextPage(ctx context.Context) (ListKeysPage, error) {
	page, err := l.genPager.NextPage(ctx)
	if err != nil {
		return ListKeysPage{}, err
	}
	return listKeysPageFromGenerated(page), nil
}

// ListDeletedKeys is the interface for the Client.ListDeletedKeys operation
type ListDeletedKeysPager struct {
	genPager *generated.KeyVaultClientGetDeletedKeysPager
}

func NewListDeletedKeysPager(p *generated.KeyVaultClientGetDeletedKeysPager) *ListDeletedKeysPager {
	return &ListDeletedKeysPager{genPager: p}
}

// More returns true if there are more pages to retrieve.
func (l *ListDeletedKeysPager) More() bool {
	return l.genPager.More()
}

// NextPage fetches the next page of results.
func (l *ListDeletedKeysPager) NextPage(ctx context.Context) (ListDeletedKeysPage, error) {
	page, err := l.genPager.NextPage(ctx)
	if err != nil {
		return ListDeletedKeysPage{}, err
	}

	var values []*models.DeletedKeyItem
	for _, d := range page.Value {
		values = append(values, convert.DeletedKeyItemFromGenerated(d))
	}

	return ListDeletedKeysPage{
		NextLink:    page.NextLink,
		DeletedKeys: values,
	}, nil
}

// ListKeyVersionsPager is a Pager for Client.ListKeyVersions results
type ListKeyVersionsPager struct {
	genPager *generated.KeyVaultClientGetKeyVersionsPager
}

func NewListKeyVersionsPager(p *generated.KeyVaultClientGetKeyVersionsPager) *ListKeyVersionsPager {
	return &ListKeyVersionsPager{genPager: p}
}

// More returns true if there are more pages to retrieve.
func (l *ListKeyVersionsPager) More() bool {
	return l.genPager.More()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *ListKeyVersionsPager) NextPage(ctx context.Context) (ListKeyVersionsPage, error) {
	page, err := l.genPager.NextPage(ctx)
	if err != nil {
		return ListKeyVersionsPage{}, err
	}
	return listKeyVersionsPageFromGenerated(page), nil
}

// convert internal Response to ListKeysPage
func listKeysPageFromGenerated(i generated.KeyVaultClientGetKeysResponse) ListKeysPage {
	var keys []*models.KeyItem
	for _, k := range i.Value {
		keys = append(keys, convert.KeyItemFromGenerated(k))
	}
	return ListKeysPage{
		NextLink: i.NextLink,
		Keys:     keys,
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
		NextLink: i.NextLink,
		Keys:     keys,
	}
}
