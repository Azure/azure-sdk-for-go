// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azinternal "github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/azblob"
)

// ListContainersIterator provides an abstraction for iterating over a collection of containers.
type ListContainersIterator struct {
	client *ServiceClient
	op     *ListContainersOptions
}

// NextPage fetches the next page.  It returns the page and nil or nil and an
// error.  If there are no more pages the error returned is azcore.IterationDone.
func (iter *ListContainersIterator) NextPage(ctx context.Context) (*ListContainersPage, error) {
	if iter.op.Marker != nil && *iter.op.Marker == "" {
		return nil, azcore.IterationDone
	}
	req := iter.client.s.ListContainersCreateRequest(iter.client.u, iter.client.p, iter.op)
	resp, err := req.Do(ctx)
	if err != nil {
		return nil, err
	}
	page, err := iter.client.s.ListContainersHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	iter.op.Marker = page.NextMarker
	return page, nil
}

// ListContainersOptions - Options for listing containers.
type ListContainersOptions = azinternal.ListContainersOptions

// ListContainersPage - An enumeration of containers
type ListContainersPage = azinternal.ListContainersPage
