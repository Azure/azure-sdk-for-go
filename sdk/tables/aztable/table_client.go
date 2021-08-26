// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"encoding/json"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	generated "github.com/Azure/azure-sdk-for-go/sdk/tables/aztable/internal"
)

// A Client represents a client to the tables service affinitized to a specific table.
type Client struct {
	client  *generated.TableClient
	service *ServiceClient
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

// NewClient creates a Client struct in the context of the table specified in tableName, using the specified serviceURL, credential, and options.
func NewClient(serviceURL string, tableName string, cred azcore.Credential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	s, err := NewServiceClient(serviceURL, cred, options)
	return s.NewClient(tableName), err
}

// Create creates the table with the tableName specified when NewClient was called.
func (t *Client) Create(ctx context.Context, options *CreateTableOptions) (CreateTableResponse, error) {
	return t.service.CreateTable(ctx, t.name, options)
}

// Delete deletes the table with the tableName specified when NewClient was called.
func (t *Client) Delete(ctx context.Context, options *DeleteTableOptions) (DeleteTableResponse, error) {
	return t.service.DeleteTable(ctx, t.name, options)
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
//     fmt.Printf("The page contains %i results.\n", len(resp.Value))
// }
// err := pager.Err()
func (t *Client) List(listOptions *ListEntitiesOptions) ListEntitiesPager {
	return &tableEntityQueryResponsePager{
		tableClient:       t,
		listOptions:       listOptions,
		tableQueryOptions: &generated.TableQueryEntitiesOptions{},
	}
}

// GetEntity retrieves a specific entity from the service using the specified partitionKey and rowKey values. If no entity is available it returns an error
func (t *Client) GetEntity(ctx context.Context, partitionKey string, rowKey string, options *GetEntityOptions) (GetEntityResponse, error) {
	if options == nil {
		options = &GetEntityOptions{}
	}

	genOptions, queryOptions := options.toGenerated()
	resp, err := t.client.QueryEntityWithPartitionAndRowKey(ctx, t.name, partitionKey, rowKey, genOptions, queryOptions)
	if err != nil {
		return GetEntityResponse{}, err
	}
	return newGetEntityResponse(resp)
}

// AddEntity adds an entity (described by a byte slice) to the table. This method returns an error if an entity with
// the same PartitionKey and RowKey already exists in the table. If the supplied entity does not contain both a PartitionKey
// and a RowKey an error will be returned.
func (t *Client) AddEntity(ctx context.Context, entity []byte, options *AddEntityOptions) (AddEntityResponse, error) {
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
	return addEntityResponseFromGenerated(&resp), err
}

// DeleteEntity deletes the entity with the specified partitionKey and rowKey from the table.
func (t *Client) DeleteEntity(ctx context.Context, partitionKey string, rowKey string, options *DeleteEntityOptions) (DeleteEntityResponse, error) {
	if options == nil {
		options = &DeleteEntityOptions{}
	}
	if options.ETag == nil {
		nilEtag := azcore.ETag("*")
		options.ETag = &nilEtag
	}
	resp, err := t.client.DeleteEntity(ctx, t.name, partitionKey, rowKey, string(*options.ETag), options.toGenerated(), &generated.QueryOptions{})
	return deleteEntityResponseFromGenerated(&resp), err
}

// UpdateEntity updates the specified table entity if it exists.
// If updateMode is Replace, the entity will be replaced. This is the only way to remove properties from an existing entity.
// If updateMode is Merge, the property values present in the specified entity will be merged with the existing entity. Properties not specified in the merge will be unaffected.
// The specified etag value will be used for optimistic concurrency. If the etag does not match the value of the entity in the table, the operation will fail.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
func (t *Client) UpdateEntity(ctx context.Context, entity []byte, updateMode EntityUpdateMode, options *UpdateEntityOptions) (UpdateEntityResponse, error) {
	if options == nil {
		options = &UpdateEntityOptions{}
	}

	if options.ETag == nil {
		star := azcore.ETag("*")
		options.ETag = &star
	}

	var mapEntity map[string]interface{}
	err := json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return UpdateEntityResponse{}, err
	}

	pk := mapEntity[partitionKey]
	partKey := pk.(string)

	rk := mapEntity[rowKey]
	rowkey := rk.(string)

	switch updateMode {
	case MergeEntity:
		resp, err := t.client.MergeEntity(ctx, t.name, partKey, rowkey, options.toGeneratedMergeEntity(mapEntity), &generated.QueryOptions{})
		return updateEntityResponseFromMergeGenerated(&resp), err
	case ReplaceEntity:
		resp, err := t.client.UpdateEntity(ctx, t.name, partKey, rowkey, options.toGeneratedUpdateEntity(mapEntity), &generated.QueryOptions{})
		return updateEntityResponseFromUpdateGenerated(&resp), err
	}
	if pk == "" || rk == "" {
		return UpdateEntityResponse{}, errPartitionKeyRowKeyError
	}
	return UpdateEntityResponse{}, errInvalidUpdateMode
}

// InsertEntity inserts an entity if it does not already exist in the table. If the entity does exist, the entity is
// replaced or merged as specified the updateMode parameter. If the entity exists and updateMode is Merge, the property
// values present in the specified entity will be merged with the existing entity rather than replaced.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
func (t *Client) InsertEntity(ctx context.Context, entity []byte, updateMode EntityUpdateMode, options *InsertEntityOptions) (InsertEntityResponse, error) {
	var mapEntity map[string]interface{}
	err := json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return InsertEntityResponse{}, err
	}

	pk := mapEntity[partitionKey]
	partKey := pk.(string)

	rk := mapEntity[rowKey]
	rowkey := rk.(string)

	// TODO: Fix the options on merge/update
	switch updateMode {
	case MergeEntity:
		resp, err := t.client.MergeEntity(ctx, t.name, partKey, rowkey, &generated.TableMergeEntityOptions{TableEntityProperties: mapEntity}, &generated.QueryOptions{})
		return insertEntityFromGeneratedMerge(&resp), err
	case ReplaceEntity:
		resp, err := t.client.UpdateEntity(ctx, t.name, partKey, rowkey, &generated.TableUpdateEntityOptions{TableEntityProperties: mapEntity}, &generated.QueryOptions{})
		return insertEntityFromGeneratedUpdate(&resp), err
	}
	if pk == "" || rk == "" {
		return InsertEntityResponse{}, errPartitionKeyRowKeyError
	}
	return InsertEntityResponse{}, errInvalidUpdateMode
}

// GetAccessPolicy retrieves details about any stored access policies specified on the table that may be used with the Shared Access Signature
func (t *Client) GetAccessPolicy(ctx context.Context, options *GetAccessPolicyOptions) (GetAccessPolicyResponse, error) {
	resp, err := t.client.GetAccessPolicy(ctx, t.name, options.toGenerated())
	return getAccessPolicyResponseFromGenerated(&resp), err
}

// SetAccessPolicy sets stored access policies for the table that may be used with SharedAccessSignature
func (t *Client) SetAccessPolicy(ctx context.Context, options *SetAccessPolicyOptions) (SetAccessPolicyResponse, error) {
	response, err := t.client.SetAccessPolicy(ctx, t.name, options.toGenerated())
	if len(options.TableACL) > 5 {
		err = errTooManyAccessPoliciesError
	}
	return setAccessPolicyResponseFromGenerated(&response), err
}
