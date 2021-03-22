// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	//"net/url"
	//"time"
)

// A TableClient represents a URL to an Azure Storage blob; the blob may be a block blob, append blob, or page blob.
type TableClient struct {
	client *tableClient
	cred   SharedKeyCredential
}

// NewTableClient creates a TableClient object using the specified URL and request policy pipeline.
func NewTableClient(serviceURL string, cred azcore.Credential, options *ClientOptions) (TableClient, error) {
	con := newConnection(serviceURL, cred, options.getConnectionOptions())

	c, _ := cred.(*SharedKeyCredential)

	return TableClient{client: &tableClient{con}, cred: *c}, nil
}

// Create
func (t TableClient) Create(ctx context.Context, name string) (TableResponseResponse, error) {
	resp, err := t.client.Create(ctx, TableProperties{&name}, nil, nil)
	if resp == nil {
		return TableResponseResponse{}, err
	} else {
		return resp.(TableResponseResponse), err
	}
}

// Delete
func (t TableClient) Delete(ctx context.Context, name string) (TableDeleteResponse, error){
	return t.client.Delete(ctx, name, nil)
}