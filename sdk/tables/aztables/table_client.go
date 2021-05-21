// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/runtime"
)

// A TableClient represents a URL to an Azure Storage blob; the blob may be a block blob, append blob, or page blob.
type TableClient struct {
	client  *tableClient
	service *TableServiceClient
	cred    SharedKeyCredential
	Name    string
}

type TableUpdateMode string

const (
	Replace TableUpdateMode = "replace"
	Merge   TableUpdateMode = "merge"
)

// NewTableClient creates a TableClient object using the specified URL and request policy pipeline.
func NewTableClient(tableName string, serviceURL string, cred azcore.Credential, options *TableClientOptions) (*TableClient, error) {
	s, err := NewTableServiceClient(serviceURL, cred, options)
	return s.GetTableClient(tableName), err
}

// Create creates the table with the name specified in NewTableClient
func (t *TableClient) Create(ctx context.Context) (*TableResponseResponse, *runtime.ResponseError) {
	return t.service.Create(ctx, t.Name)
}

// Delete deletes the current table
func (t *TableClient) Delete(ctx context.Context) (*TableDeleteResponse, *runtime.ResponseError) {
	return t.service.Delete(ctx, t.Name)
}

// Query queries the tables using the specified QueryOptions
func (t *TableClient) Query(queryOptions QueryOptions) TableEntityQueryResponsePager {
	return &tableEntityQueryResponsePager{tableClient: t, queryOptions: &queryOptions, tableQueryOptions: &TableQueryEntitiesOptions{}}
}

// QueryAsModel queries the table using the specified QueryOptions and attempts to serialize the response as the supplied interface type
func (t *TableClient) QueryAsModel(opt QueryOptions, s FromMapper) StructEntityQueryResponsePager {
	return &structQueryResponsePager{mapper: s, tableClient: t, queryOptions: &opt, tableQueryOptions: &TableQueryEntitiesOptions{}}
}

func (t *TableClient) GetEntity(ctx context.Context, partitionKey string, rowKey string) (MapOfInterfaceResponse, error) {
	resp, err := t.client.QueryEntityWithPartitionAndRowKey(ctx, t.Name, partitionKey, rowKey, &TableQueryEntityWithPartitionAndRowKeyOptions{}, &QueryOptions{})
	if err != nil {
		return resp, err
	}
	castAndRemoveAnnotations(&resp.Value)
	return resp, err
}

// AddEntity Creates an entity from a map value.
func (t *TableClient) AddEntity(ctx context.Context, entity map[string]interface{}) (*TableInsertEntityResponse, *runtime.ResponseError) {
	toOdataAnnotatedDictionary(&entity)
	resp, err := t.client.InsertEntity(ctx, t.Name, &TableInsertEntityOptions{TableEntityProperties: entity, ResponsePreference: ResponseFormatReturnNoContent.ToPtr()}, &QueryOptions{})
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
	resp, err := t.client.InsertEntity(ctx, t.Name, &TableInsertEntityOptions{TableEntityProperties: *entmap, ResponsePreference: ResponseFormatReturnNoContent.ToPtr()}, &QueryOptions{})
	if err == nil {
		insertResp := resp.(TableInsertEntityResponse)
		return &insertResp, nil
	} else {
		return nil, convertErr(err)
	}
}

func (t *TableClient) DeleteEntity(ctx context.Context, partitionKey string, rowKey string, etag string) (TableDeleteEntityResponse, error) {
	return t.client.DeleteEntity(ctx, t.Name, partitionKey, rowKey, etag, nil, &QueryOptions{})
}

// UpdateEntity updates the specified table entity if it exists.
// If updateMode is Replace, the entity will be replaced.
// If updateMode is Merge, the property values present in the specified entity will be merged with the existing entity.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
func (t *TableClient) UpdateEntity(ctx context.Context, entity map[string]interface{}, etag *string, updateMode TableUpdateMode) (interface{}, error) {
	pk := entity[partitionKey].(string)
	rk := entity[rowKey].(string)
	var ifMatch string = "*"
	if etag != nil {
		ifMatch = *etag
	}
	switch updateMode {
	case Merge:
		return t.client.MergeEntity(ctx, t.Name, pk, rk, &TableMergeEntityOptions{IfMatch: &ifMatch, TableEntityProperties: entity}, &QueryOptions{})
	case Replace:
		return t.client.UpdateEntity(ctx, t.Name, pk, rk, &TableUpdateEntityOptions{IfMatch: &ifMatch, TableEntityProperties: entity}, &QueryOptions{})
	}
	return nil, errors.New("Invalid TableUpdateMode")
}

// UpsertEntity replaces the specified table entity if it exists or creates the entity if it does not exist.
// If the entity exists and updateMode is Merge, the property values present in the specified entity will be merged with the existing entity rather than replaced.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
func (t *TableClient) UpsertEntity(ctx context.Context, entity map[string]interface{}, updateMode TableUpdateMode) (interface{}, error) {
	pk := entity[partitionKey].(string)
	rk := entity[rowKey].(string)

	switch updateMode {
	case Merge:
		return t.client.MergeEntity(ctx, t.Name, pk, rk, &TableMergeEntityOptions{TableEntityProperties: entity}, &QueryOptions{})
	case Replace:
		return t.client.UpdateEntity(ctx, t.Name, pk, rk, &TableUpdateEntityOptions{TableEntityProperties: entity}, &QueryOptions{})
	}
	return nil, errors.New("Invalid TableUpdateMode")
}
