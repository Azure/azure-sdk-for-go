// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

// A TableClient represents a client to the tables service affinitized to a specific table.
type TableClient struct {
	client  *generated.TableClient
	service *TableServiceClient
	cred    azcore.Credential
	name    string
}

// EntityUpdateMode specifies what type of update to do on InsertEntity or UpdateEntity. ReplaceEntity
// will replace an existing entity, MergeEntity will merge properties of the entities.
type EntityUpdateMode string

const (
	ReplaceEntity EntityUpdateMode = "replace"
	MergeEntity   EntityUpdateMode = "merge"
)

// NewTableClient creates a TableClient struct in the context of the table specified in tableName, using the specified serviceURL, credential, and options.
func NewTableClient(tableName string, serviceURL string, cred azcore.Credential, options *TableClientOptions) (*TableClient, error) {
	s, err := NewTableServiceClient(serviceURL, cred, options)
	return s.NewTableClient(tableName), err
}

// Create creates the table with the tableName specified when NewTableClient was called.
func (t *TableClient) Create(ctx context.Context, options *CreateTableOptions) (*CreateTableResponse, error) {
	resp, err := t.service.CreateTable(ctx, t.name, options)
	return createTableResponseFromGen(&resp), err
}

// Delete deletes the table with the tableName specified when NewTableClient was called.
func (t *TableClient) Delete(ctx context.Context, options *DeleteTableOptions) (*DeleteTableResponse, error) {
	resp, err := t.service.DeleteTable(ctx, t.name, options.toGenerated())
	return deleteTableResponseFromGen(&resp), err
}

// List queries the entities using the specified ListEntitiesOptions.
// ListOptions can specify the following properties to affect the query results returned:
//
// Filter: An OData filter expression that limits results to those entities that satisfy the filter expression.
// For example, the following expression would return only entities with a PartitionKey of 'foo': "PartitionKey eq 'foo'"
//
// Select: A comma delimited list of entity property names that selects which set of entity properties to return in the result set.
// For example, the following value would return results containing only the PartitionKey and RowKey properties: "PartitionKey, RowKey"
//
// Top: The maximum number of entities that will be returned per page of results.
// Note: This value does not limit the total number of results if NextPage is called on the returned Pager until it returns false.
//
// List returns a Pager, which allows iteration through each page of results. Example:
//
// options := &ListEntitiesOptions{Filter: to.StringPtr("PartitionKey eq 'pk001'"), Top: to.Int32Ptr(25), Select: to.StringPtr("PartitionKey,RowKey,Value,Price")}
// pager := client.List(options) // Pass in 'nil' if you want to return all Entities for an account.
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.Printf("The page contains %i results.\n", len(resp.TableEntityQueryResponse.Value))
// }
// err := pager.Err()
func (t *TableClient) List(listOptions *ListEntitiesOptions) TableEntityListResponsePager {
	return &tableEntityQueryResponsePager{
		tableClient:       t,
		listOptions:       listOptions,
		tableQueryOptions: &generated.TableQueryEntitiesOptions{},
	}
}

// GetEntity retrieves a specific entity from the service using the specified partitionKey and rowKey values.
func (t *TableClient) GetEntity(ctx context.Context, partitionKey string, rowKey string, options *GetEntityOptions) (ByteArrayResponse, error) {
	if options == nil {
		options = &GetEntityOptions{}
	}

	genOptions, queryOptions := options.toGenerated()
	resp, err := t.client.QueryEntityWithPartitionAndRowKey(ctx, t.name, partitionKey, rowKey, genOptions, queryOptions)
	if err != nil {
		return ByteArrayResponse{}, err
	}
	return newByteArrayResponse(resp)
}

// AddEntity adds an entity (described by a JSON byte slice) to the table. This method returns an error if an entity with
// the same PartitionKey and RowKey already exists in the table.
func (t *TableClient) AddEntity(ctx context.Context, entity []byte, options *AddEntityOptions) (AddEntityResponse, error) {
	var mapEntity map[string]interface{}
	err := json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return AddEntityResponse{}, err
	}
	resp, err := t.client.InsertEntity(ctx, t.name, &generated.TableInsertEntityOptions{TableEntityProperties: mapEntity, ResponsePreference: generated.ResponseFormatReturnNoContent.ToPtr()}, nil)
	if err != nil {
		err = checkEntityForPkRk(&mapEntity, err)
		return AddEntityResponse{}, err
	}
	return *addEntityResponseFromGenerated(&resp), err
}

// DeleteEntity deletes the entity with the specified partitionKey and rowKey from the table.
func (t *TableClient) DeleteEntity(ctx context.Context, partitionKey string, rowKey string, etag *string, options *DeleteEntityOptions) (*DeleteEntityResponse, error) {
	if etag == nil {
		nilEtag := "*"
		etag = &nilEtag
	}
	resp, err := t.client.DeleteEntity(ctx, t.name, partitionKey, rowKey, *etag, options.toGenerated(), &generated.QueryOptions{})
	return deleteEntityResponseFromGenerated(&resp), err
}

// UpdateEntity updates the specified table entity if it exists.
// If updateMode is Replace, the entity will be replaced. This is the only way to remove properties from an existing entity.
// If updateMode is Merge, the property values present in the specified entity will be merged with the existing entity. Properties not specified in the merge will be unaffected.
// The specified etag value will be used for optimistic concurrency. If the etag does not match the value of the entity in the table, the operation will fail.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
func (t *TableClient) UpdateEntity(ctx context.Context, entity []byte, etag *string, updateMode EntityUpdateMode, options *UpdateEntityOptions) (*UpdateEntityResponse, error) {
	var ifMatch string = "*"
	if etag != nil {
		ifMatch = *etag
	}

	if options == nil {
		options = &UpdateEntityOptions{}
	}
	options.IfMatch = &ifMatch

	var mapEntity map[string]interface{}
	err := json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return nil, err
	}

	pk, okPk := mapEntity[partitionKey]
	partKey := pk.(string)

	rk, okRk := mapEntity[rowKey]
	rowkey := rk.(string)

	switch updateMode {
	case MergeEntity:
		resp, err := t.client.MergeEntity(ctx, t.name, partKey, rowkey, options.toGeneratedMergeEntity(mapEntity), &generated.QueryOptions{})
		return updateEntityResponseFromMergeGenerated(&resp), err
	case ReplaceEntity:
		resp, err := t.client.UpdateEntity(ctx, t.name, partKey, rowkey, options.toGeneratedUpdateEntity(mapEntity), &generated.QueryOptions{})
		return updateEntityResponseFromUpdateGenerated(&resp), err
	}
	if !okPk || !okRk {
		return nil, errPartitionKeyRowKeyError
	}
	return nil, errInvalidUpdateMode
}

// InsertEntity inserts an entity if it does not already exist in the table. If the entity does exist, the entity is
// replaced or merged as specified the updateMode parameter. If the entity exists and updateMode is Merge, the property
// values present in the specified entity will be merged with the existing entity rather than replaced.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
func (t *TableClient) InsertEntity(ctx context.Context, entity []byte, updateMode EntityUpdateMode, options *InsertEntityOptions) (*InsertEntityResponse, error) {
	var mapEntity map[string]interface{}
	err := json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return nil, err
	}

	pk, okPk := mapEntity[partitionKey]
	partKey := pk.(string)

	rk, okRk := mapEntity[rowKey]
	rowkey := rk.(string)

	switch updateMode {
	case MergeEntity:
		resp, err := t.client.MergeEntity(ctx, t.name, partKey, rowkey, &generated.TableMergeEntityOptions{TableEntityProperties: mapEntity}, &generated.QueryOptions{})
		return insertEntityFromGeneratedMerge(&resp), err
	case ReplaceEntity:
		resp, err := t.client.UpdateEntity(ctx, t.name, partKey, rowkey, &generated.TableUpdateEntityOptions{TableEntityProperties: mapEntity}, &generated.QueryOptions{})
		return insertEntityFromGeneratedUpdate(&resp), err
	}
	if !okPk || !okRk {
		return nil, errPartitionKeyRowKeyError
	}
	return nil, errInvalidUpdateMode
}

// GetAccessPolicy retrieves details about any stored access policies specified on the table that may be used with the Shared Access Signature
func (t *TableClient) GetAccessPolicy(ctx context.Context, options *GetAccessPolicyOptions) (*GetAccessPolicyResponse, error) {
	resp, err := t.client.GetAccessPolicy(ctx, t.name, options.toGenerated())
	return getAccessPolicyResponseFromGenerated(&resp), err
}

// SetAccessPolicy sets stored access policies for the table that may be used with SharedAccessSignature
func (t *TableClient) SetAccessPolicy(ctx context.Context, options *SetAccessPolicyOptions) (*SetAccessPolicyResponse, error) {
	response, err := t.client.SetAccessPolicy(ctx, t.name, options.toGenerated())
	if len(options.TableACL) > 5 {
		err = errTooManyAccessPoliciesError
	}
	return setAccessPolicyResponseFromGenerated(&response), err
}
