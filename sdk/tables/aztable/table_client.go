// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// A TableClient represents a client to the tables service affinitized to a specific table.
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
//     fmt.Sprintf("The page contains %i results", len(resp.TableEntityQueryResponse.Value))
// }
// err := pager.Err()
func (t *TableClient) Query(queryOptions *QueryOptions) TableEntityQueryResponsePager {
	return &tableEntityQueryResponsePager{tableClient: t, queryOptions: queryOptions, tableQueryOptions: &TableQueryEntitiesOptions{}}
}

// GetEntity retrieves a specific entity from the service using the specified partitionKey and rowKey values.
func (t *TableClient) GetEntity(ctx context.Context, partitionKey string, rowKey string) (ByteArrayResponse, error) {
	resp, err := t.client.QueryEntityWithPartitionAndRowKey(ctx, t.Name, partitionKey, rowKey, &TableQueryEntityWithPartitionAndRowKeyOptions{}, &QueryOptions{})
	if err != nil {
		return ByteArrayResponse{}, err
	}
	// err = castAndRemoveAnnotations(&resp.Value)
	return newByteArrayResponse(resp)
}

// AddEntity adds an entity from an arbitrary interface value to the table.
// An entity must have at least a PartitionKey and RowKey property.
func (t *TableClient) AddEntity(ctx context.Context, entity []byte) (interface{}, error) {
	var mapEntity map[string]interface{}
	err := json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return entity, err
	}
	resp, err := t.client.InsertEntity(ctx, t.Name, &TableInsertEntityOptions{TableEntityProperties: mapEntity, ResponsePreference: ResponseFormatReturnNoContent.ToPtr()}, &QueryOptions{})
	if err == nil {
		insertResp := resp.(TableInsertEntityResponse)
		return insertResp, nil
	} else {
		err = checkEntityForPkRk(&mapEntity, err)
		return TableInsertEntityResponse{}, err
	}
}

// DeleteEntity deletes the entity with the specified partitionKey and rowKey from the table.
func (t *TableClient) DeleteEntity(ctx context.Context, partitionKey string, rowKey string, etag *string) (TableDeleteEntityResponse, error) {
	if etag == nil {
		nilEtag := "*"
		etag = &nilEtag
	}
	return t.client.DeleteEntity(ctx, t.Name, partitionKey, rowKey, *etag, nil, &QueryOptions{})
}

// UpdateEntity updates the specified table entity if it exists.
// If updateMode is Replace, the entity will be replaced. This is the only way to remove properties from an existing entity.
// If updateMode is Merge, the property values present in the specified entity will be merged with the existing entity. Properties not specified in the merge will be unaffected.
// The specified etag value will be used for optimistic concurrency. If the etag does not match the value of the entity in the table, the operation will fail.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
func (t *TableClient) UpdateEntity(ctx context.Context, entity []byte, etag *string, updateMode TableUpdateMode) (interface{}, error) {
	var ifMatch string = "*"
	if etag != nil {
		ifMatch = *etag
	}

	var mapEntity map[string]interface{}
	err := json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return entity, err
	}

	pk, _ := mapEntity[partitionKey]
	partKey := pk.(string)

	rk, _ := mapEntity[rowKey]
	rowkey := rk.(string)

	switch updateMode {
	case Merge:
		return t.client.MergeEntity(ctx, t.Name, partKey, rowkey, &TableMergeEntityOptions{IfMatch: &ifMatch, TableEntityProperties: mapEntity}, &QueryOptions{})
	case Replace:
		return t.client.UpdateEntity(ctx, t.Name, partKey, rowkey, &TableUpdateEntityOptions{IfMatch: &ifMatch, TableEntityProperties: mapEntity}, &QueryOptions{})
	}
	return nil, errors.New("Invalid TableUpdateMode")
}

// UpsertEntity replaces the specified table entity if it exists or creates the entity if it does not exist.
// If the entity exists and updateMode is Merge, the property values present in the specified entity will be merged with the existing entity rather than replaced.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
func (t *TableClient) UpsertEntity(ctx context.Context, entity []byte, updateMode TableUpdateMode) (interface{}, error) {
	var mapEntity map[string]interface{}
	err := json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return entity, err
	}

	pk, _ := mapEntity[partitionKey]
	partKey := pk.(string)

	rk, _ := mapEntity[rowKey]
	rowkey := rk.(string)

	switch updateMode {
	case Merge:
		return t.client.MergeEntity(ctx, t.Name, partKey, rowkey, &TableMergeEntityOptions{TableEntityProperties: mapEntity}, &QueryOptions{})
	case Replace:
		return t.client.UpdateEntity(ctx, t.Name, partKey, rowkey, &TableUpdateEntityOptions{TableEntityProperties: mapEntity}, &QueryOptions{})
	}
	return nil, errors.New("Invalid TableUpdateMode")
}

type TableAccessPolicy struct {
	Start      time.Time
	Expiry     time.Time
	Permission string
}

// GetTableAccessPolicy retrieves details about any stored access policies specified on the table that may be used with Shared Access Signatures
func (t *TableClient) GetTableAccessPolicy(ctx context.Context) (SignedIdentifierArrayResponse, error) {
	accessPolicies, err := t.client.GetAccessPolicy(ctx, t.Name, nil)
	return accessPolicies, err
}
