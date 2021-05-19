// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
)

// A TableServiceClient represents a URL to an Azure Storage blob; the blob may be a block blob, append blob, or page blob.
type TableServiceClient struct {
	client  *tableClient
	service *serviceClient
	cred    SharedKeyCredential
}

// NewTableServiceClient creates a TableClient object using the specified URL and request policy pipeline.
func NewTableServiceClient(serviceURL string, cred azcore.Credential, options *TableClientOptions) (*TableServiceClient, error) {
	con := newConnection(serviceURL, cred, options.getConnectionOptions())
	c, _ := cred.(*SharedKeyCredential)
	return &TableServiceClient{client: &tableClient{con}, service: &serviceClient{con}, cred: *c}, nil
}

// Gets a TableClient affinitzed to the specified table name and initialized with the same serviceURL and credentials as this TableServiceClient
func (t *TableServiceClient) GetTableClient(tableName string) *TableClient {
	return &TableClient{client: t.client, cred: t.cred, name: tableName, service: t}
}

// Creates a table with the specified name
func (t *TableServiceClient) Create(ctx context.Context, name string) (*TableResponseResponse, *runtime.ResponseError) {
	var r *TableResponseResponse = nil
	resp, err := t.client.Create(ctx, TableProperties{&name}, new(TableCreateOptions), new(QueryOptions))
	if err == nil {
		tableResp := resp.(TableResponseResponse)
		r = &tableResp
	}
	return r, convertErr(err)
}

// Deletes a table by name
func (t *TableServiceClient) Delete(ctx context.Context, name string) (*TableDeleteResponse, *runtime.ResponseError) {
	resp, err := t.client.Delete(ctx, name, nil)
	return &resp, convertErr(err)
}

// Queries the tables using the specified QueryOptions
func (t *TableServiceClient) QueryTables(queryOptions QueryOptions) TableQueryResponsePager {
	return &tableQueryResponsePager{client: t.client, queryOptions: &queryOptions, tableQueryOptions: new(TableQueryOptions)}
}

func convertErr(err error) *runtime.ResponseError {
	var e *runtime.ResponseError
	if err == nil || !errors.As(err, &e) {
		return nil
	} else {
		return e
	}
}
