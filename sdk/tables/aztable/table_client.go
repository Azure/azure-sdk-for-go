// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// A TableClient represents a client to the tables service affinitized to a specific table.
type TableClient struct {
	client  *tableClient
	service *TableServiceClient
	cred    azcore.Credential
	Name    string
}

type TableUpdateMode string

const (
	Replace TableUpdateMode = "replace"
	Merge   TableUpdateMode = "merge"
)

// NewTableClient creates a TableClient struct in the context of the table specified in tableName, using the specified serviceURL, credential, and options.
func NewTableClient(tableName string, serviceURL string, cred azcore.Credential, options *TableClientOptions) (*TableClient, error) {
	s, err := NewTableServiceClient(serviceURL, cred, options)
	return s.NewTableClient(tableName), err
}

// Create creates the table with the tableName specified when NewTableClient was called.
func (t *TableClient) Create(ctx context.Context) (TableResponseResponse, error) {
	return t.service.Create(ctx, t.Name)
}

// Delete deletes the table with the tableName specified when NewTableClient was called.
func (t *TableClient) Delete(ctx context.Context) (TableDeleteResponse, error) {
	return t.service.Delete(ctx, t.Name)
}

// Query queries the tables using the specified QueryOptions.
// QueryOptions can specify the following properties to affect the query results returned:
//
// Filter: An Odata filter expression that limits results to those entities that satisfy the filter expression.
// For example, the following expression would return only entities with a PartitionKey of 'foo': "PartitionKey eq 'foo'"
//
// Select: A comma delimited list of entity property names that selects which set of entity properties to return in the result set.
// For example, the following value would return results containing only the PartitionKey and RowKey properties: "PartitionKey, RowKey"
//
// Top: The maximum number of entities that will be returned per page of results.
// Note: This value does not limit the total number of results if NextPage is called on the returned Pager until it returns false.
//
// Query returns a Pager, which allows iteration through each page of results. Example:
//
// pager := client.Query(nil)
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.sprintf("The page contains %i results", len(resp.TableEntityQueryResponse.Value))
// }
// err := pager.Err()
func (t *TableClient) Query(queryOptions *QueryOptions) TableEntityQueryResponsePager {
	if queryOptions == nil {
		queryOptions = &QueryOptions{}
	}
	return &tableEntityQueryResponsePager{tableClient: t, queryOptions: queryOptions, tableQueryOptions: &TableQueryEntitiesOptions{}}
}

// GetEntity retrieves a specific entity from the service using the specified partitionKey and rowKey values.
func (t *TableClient) GetEntity(ctx context.Context, partitionKey string, rowKey string) (MapOfInterfaceResponse, error) {
	resp, err := t.client.QueryEntityWithPartitionAndRowKey(ctx, t.Name, partitionKey, rowKey, &TableQueryEntityWithPartitionAndRowKeyOptions{}, nil)
	if err != nil {
		return resp, err
	}
	err = castAndRemoveAnnotations(&resp.Value)
	return resp, err
}

// AddEntity adds an entity from an arbitrary interface value to the table.
// An entity must have at least a PartitionKey and RowKey property.
func (t *TableClient) AddEntity(ctx context.Context, entity interface{}) (TableInsertEntityResponse, error) {
	entmap, err := toMap(entity)
	if err != nil {
		return TableInsertEntityResponse{}, azcore.NewResponseError(err, nil)
	}
	resp, err := t.client.InsertEntity(ctx, t.Name, &TableInsertEntityOptions{TableEntityProperties: *entmap, ResponsePreference: ResponseFormatReturnNoContent.ToPtr()}, nil)
	if err == nil {
		insertResp := resp.(TableInsertEntityResponse)
		return insertResp, nil
	} else {
		err = checkEntityForPkRk(entmap, err)
		return TableInsertEntityResponse{}, err
	}
}

// DeleteEntity deletes the entity with the specified partitionKey and rowKey from the table.
func (t *TableClient) DeleteEntity(ctx context.Context, partitionKey string, rowKey string, etag string) (TableDeleteEntityResponse, error) {
	return t.client.DeleteEntity(ctx, t.Name, partitionKey, rowKey, etag, nil, nil)
}

// UpdateEntity updates the specified table entity if it exists.
// If updateMode is Replace, the entity will be replaced. This is the only way to remove properties from an existing entity.
// If updateMode is Merge, the property values present in the specified entity will be merged with the existing entity. Properties not specified in the merge will be unaffected.
// The specified etag value will be used for optimistic concurrency. If the etag does not match the value of the entity in the table, the operation will fail.
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
		return t.client.MergeEntity(ctx, t.Name, pk, rk, &TableMergeEntityOptions{IfMatch: &ifMatch, TableEntityProperties: entity}, nil)
	case Replace:
		return t.client.UpdateEntity(ctx, t.Name, pk, rk, &TableUpdateEntityOptions{IfMatch: &ifMatch, TableEntityProperties: entity}, nil)
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
		return t.client.MergeEntity(ctx, t.Name, pk, rk, &TableMergeEntityOptions{TableEntityProperties: entity}, nil)
	case Replace:
		return t.client.UpdateEntity(ctx, t.Name, pk, rk, &TableUpdateEntityOptions{TableEntityProperties: entity}, nil)
	}
	return nil, errors.New("Invalid TableUpdateMode")
}
