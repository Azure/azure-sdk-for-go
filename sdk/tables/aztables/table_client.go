// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
)

// A TableClient represents a URL to an Azure Storage blob; the blob may be a block blob, append blob, or page blob.
type TableClient struct {
	client  *tableClient
	service *TableServiceClient
	cred    SharedKeyCredential
	name    string
}

// NewTableClient creates a TableClient object using the specified URL and request policy pipeline.
func NewTableClient(tableName string, serviceURL string, cred azcore.Credential, options *TableClientOptions) (*TableClient, error) {
	s, err := NewTableServiceClient(serviceURL, cred, options)
	return s.GetTableClient(tableName), err
}

func (t *TableClient) Name() string {
	return t.name
}

// Create creates the table with the name specified in NewTableClient
func (t *TableClient) Create(ctx context.Context) (*TableResponseResponse, *runtime.ResponseError) {
	return t.service.Create(ctx, t.name)
}

// Delete deletes the current table
func (t *TableClient) Delete(ctx context.Context) (*TableDeleteResponse, *runtime.ResponseError) {
	return t.service.Delete(ctx, t.name)
}

// Query queries the tables using the specified QueryOptions
func (t *TableClient) Query(queryOptions QueryOptions) TableEntityQueryResponsePager {
	return &tableEntityQueryResponsePager{tableClient: t, queryOptions: &queryOptions, tableQueryOptions: &TableQueryEntitiesOptions{}}
}

func (t *TableClient) QueryAsModel(opt QueryOptions, s FromMapper) StructEntityQueryResponsePager {
	return &structQueryResponsePager{mapper: s, tableClient: t, queryOptions: &opt, tableQueryOptions: &TableQueryEntitiesOptions{}}
}

// AddEntity Creates an entity from a map value.
func (t *TableClient) AddEntity(ctx context.Context, entity map[string]interface{}) (*TableInsertEntityResponse, *runtime.ResponseError) {
	toOdataAnnotatedDictionary(&entity)
	resp, err := t.client.InsertEntity(ctx, t.name, &TableInsertEntityOptions{TableEntityProperties: &entity, ResponsePreference: ResponseFormatReturnNoContent.ToPtr()}, &QueryOptions{})
	if err == nil {
		insertResp := resp.(TableInsertEntityResponse)
		return &insertResp, nil
	} else {
		return nil, convertErr(err)
	}
}

// AddModelEntity creates an entity from an arbitrary struct value.
func (t *TableClient) AddModelEntity(ctx context.Context, entity interface{}) (*TableInsertEntityResponse, *runtime.ResponseError) {
	entmap, err := toMap(entity)
	if err != nil {
		return nil, azcore.NewResponseError(err, nil).(*runtime.ResponseError)
	}
	resp, err := t.client.InsertEntity(ctx, t.name, &TableInsertEntityOptions{TableEntityProperties: entmap, ResponsePreference: ResponseFormatReturnNoContent.ToPtr()}, &QueryOptions{})
	if err == nil {
		insertResp := resp.(TableInsertEntityResponse)
		return &insertResp, nil
	} else {
		return nil, convertErr(err)
	}
}

func (t *TableClient) DeleteEntity(ctx context.Context, partitionKey string, rowKey string, etag string) (TableDeleteEntityResponse, error) {
	return t.client.DeleteEntity(ctx, t.name, partitionKey, rowKey, etag, nil, &QueryOptions{})
}
