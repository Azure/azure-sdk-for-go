// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"encoding/json"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	generated "github.com/Azure/azure-sdk-for-go/sdk/data/aztables/internal"
)

// Client represents a client to the tables service affinitized to a specific table.
type Client struct {
	client  *generated.TableClient
	service *ServiceClient
	cred    *SharedKeyCredential
	name    string
	con     *generated.Connection
}

// NewClient creates a Client struct in the context of the table specified in the serviceURL, authorizing requests with an Azure AD access token.
// The serviceURL param is expected to have the name of the table in a format similar to: "https://myAccountName.core.windows.net/<myTableName>".
// Pass in nil for options to construct the client with the default ClientOptions.
func NewClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	rawServiceURL, tableName, err := parseURL(serviceURL)
	if err != nil {
		return nil, err
	}
	s, err := NewServiceClient(rawServiceURL, cred, options)
	if err != nil {
		return nil, err
	}
	return s.NewClient(tableName), nil
}

// NewClientWithNoCredential creates a Client struct in the context of the table specified in the serviceURL.
// The serviceURL param is expected to have the name of the table in a format similar to: "https://myAccountName.core.windows.net/<myTableName>?<SAS token>".
// Pass in nil for options to construct the client with the default ClientOptions.
func NewClientWithNoCredential(serviceURL string, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	rawServiceURL, tableName, err := parseURL(serviceURL)
	if err != nil {
		return nil, err
	}
	s, err := NewServiceClientWithNoCredential(rawServiceURL, options)
	if err != nil {
		return nil, err
	}
	return s.NewClient(tableName), nil
}

// NewClientWithSharedKey creates a Client struct in the context of the table specified in the serviceURL, authorizing requests with a shared key.
// The serviceURL param is expected to have the name of the table in a format similar to: "https://myAccountName.core.windows.net/<myTableName>".
// Pass in nil for options to construct the client with the default ClientOptions.
func NewClientWithSharedKey(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	rawServiceURL, tableName, err := parseURL(serviceURL)
	if err != nil {
		return nil, err
	}
	s, err := NewServiceClientWithSharedKey(rawServiceURL, cred, options)
	if err != nil {
		return nil, err
	}
	return s.NewClient(tableName), nil
}

func parseURL(serviceURL string) (string, string, error) {
	parsedUrl, err := url.Parse(serviceURL)
	if err != nil {
		return "", "", err
	}

	tableName := parsedUrl.Path[1:]
	rawServiceURL := parsedUrl.Scheme + "://" + parsedUrl.Host
	if parsedUrl.Scheme == "" {
		rawServiceURL = parsedUrl.Host
	}
	if strings.Contains(tableName, "/") {
		splits := strings.Split(parsedUrl.Path, "/")
		tableName = splits[len(splits)-1]
		rawServiceURL += strings.Join(splits[:len(splits)-1], "/")
	}
	sas := parsedUrl.Query()
	if len(sas) > 0 {
		rawServiceURL += "/?" + sas.Encode()
	}

	return rawServiceURL, tableName, nil
}

// CreateTableResponse contains response fields for Client.Create and ServiceClient.CreateTable
type CreateTableResponse struct {
	// The name of the table.
	TableName *string `json:"TableName,omitempty"`
}

func createTableResponseFromGen(g *generated.TableClientCreateResponse) CreateTableResponse {
	return CreateTableResponse{
		TableName: g.TableName,
	}
}

// CreateTable creates the table with the tableName specified when NewClient was called. If the service returns a non-successful
// HTTP status code, the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
// NOTE: creating a table with the same name as a table that's in the process of being deleted will return an *azcore.ResponseError
// with error code TableBeingDeleted and status code http.StatusConflict.
func (t *Client) CreateTable(ctx context.Context, options *CreateTableOptions) (CreateTableResponse, error) {
	if options == nil {
		options = &CreateTableOptions{}
	}
	resp, err := t.client.Create(ctx, generated.Enum1Three0, generated.TableProperties{TableName: &t.name}, options.toGenerated(), &generated.QueryOptions{})
	if err != nil {
		return CreateTableResponse{}, err
	}
	return createTableResponseFromGen(&resp), err
}

// Delete deletes the table with the tableName specified when NewClient was called. If the service returns a non-successful HTTP status
// code, the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
// NOTE: deleting a table can take up to 40 seconds or more to complete.  If a table with the same name is created while the delete is still
// in progress, an *azcore.ResponseError is returned with error code TableBeingDeleted and status code http.StatusConflict.
func (t *Client) Delete(ctx context.Context, options *DeleteTableOptions) (DeleteTableResponse, error) {
	return t.service.DeleteTable(ctx, t.name, options)
}

// ListEntitiesOptions contains optional parameters for Table.Query
type ListEntitiesOptions struct {
	// OData filter expression.
	Filter *string

	// Select expression using OData notation. Limits the columns on each record
	// to just those requested, e.g. "$select=PolicyAssignmentId, ResourceId".
	Select *string

	// Maximum number of records to return.
	Top *int32

	// The NextPartitionKey to start paging from
	NextPartitionKey *string

	// The NextRowKey to start paging from
	NextRowKey *string
}

func (l *ListEntitiesOptions) toQueryOptions() *generated.QueryOptions {
	if l == nil {
		return &generated.QueryOptions{}
	}

	return &generated.QueryOptions{
		Filter: l.Filter,
		Format: to.Ptr(generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata),
		Select: l.Select,
		Top:    l.Top,
	}
}

// ListEntitiesResponse contains response fields for ListEntitiesPager.NextPage
type ListEntitiesResponse struct {
	// NextPartitionKey contains the information returned from the x-ms-continuation-NextPartitionKey header response.
	NextPartitionKey *string

	// NextRowKey contains the information returned from the x-ms-continuation-NextRowKey header response.
	NextRowKey *string

	// List of table entities.
	Entities [][]byte
}

// transforms a generated query response into the ListEntitiesPaged
func newListEntitiesPage(resp generated.TableClientQueryEntitiesResponse) (ListEntitiesResponse, error) {
	marshalledValue := make([][]byte, 0)
	for _, e := range resp.TableEntityQueryResponse.Value {
		m, err := json.Marshal(e)
		if err != nil {
			return ListEntitiesResponse{}, err
		}
		marshalledValue = append(marshalledValue, m)
	}

	return ListEntitiesResponse{
		NextPartitionKey: resp.XMSContinuationNextPartitionKey,
		NextRowKey:       resp.XMSContinuationNextRowKey,
		Entities:         marshalledValue,
	}, nil
}

// NewListEntitiesPager queries the entities using the specified ListEntitiesOptions.
// ListEntitiesOptions can specify the following properties to affect the query results returned:
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
// NewListEntitiesPager returns a Pager, which allows iteration through each page of results. Use nil for listOptions if you want to use the default options.
// For more information about writing query strings, check out:
//  - API Documentation: https://docs.microsoft.com/en-us/rest/api/storageservices/querying-tables-and-entities
//  - README samples: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/data/aztables/README.md#writing-filters
func (t *Client) NewListEntitiesPager(listOptions *ListEntitiesOptions) *runtime.Pager[ListEntitiesResponse] {
	if listOptions == nil {
		listOptions = &ListEntitiesOptions{}
	}
	return runtime.NewPager(runtime.PagingHandler[ListEntitiesResponse]{
		More: func(page ListEntitiesResponse) bool {
			if page.NextPartitionKey == nil || len(*page.NextPartitionKey) == 0 || page.NextRowKey == nil || len(*page.NextRowKey) == 0 {
				return false
			}
			return true
		},
		Fetcher: func(ctx context.Context, page *ListEntitiesResponse) (ListEntitiesResponse, error) {
			var partKey *string
			var rowKey *string
			if page != nil {
				if page.NextPartitionKey != nil {
					partKey = page.NextPartitionKey
				}
				if page.NextRowKey != nil {
					rowKey = page.NextRowKey
				}
			} else {
				partKey = listOptions.NextPartitionKey
				rowKey = listOptions.NextRowKey
			}
			resp, err := t.client.QueryEntities(ctx, generated.Enum1Three0, t.name, &generated.TableClientQueryEntitiesOptions{
				NextPartitionKey: partKey,
				NextRowKey:       rowKey,
			}, listOptions.toQueryOptions())
			if err != nil {
				return ListEntitiesResponse{}, err
			}
			return newListEntitiesPage(resp)
		},
	})
}

// GetEntityOptions contains optional parameters for Client.GetEntity
type GetEntityOptions struct {
	// placeholder for future optional parameters
}

func (g *GetEntityOptions) toGenerated() (*generated.TableClientQueryEntityWithPartitionAndRowKeyOptions, *generated.QueryOptions) {
	return &generated.TableClientQueryEntityWithPartitionAndRowKeyOptions{}, &generated.QueryOptions{Format: to.Ptr(generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata)}
}

// GetEntityResponse contains response fields for Client.GetEntity
type GetEntityResponse struct {
	// ETag contains the information returned from the ETag header response.
	ETag azcore.ETag

	// The properties of the table entity.
	Value []byte
}

// newGetEntityResponse transforms a generated response to the GetEntityResponse type
func newGetEntityResponse(g generated.TableClientQueryEntityWithPartitionAndRowKeyResponse) (GetEntityResponse, error) {
	marshalledValue, err := json.Marshal(g.Value)
	if err != nil {
		return GetEntityResponse{}, err
	}

	var ETag azcore.ETag
	if g.ETag != nil {
		ETag = azcore.ETag(*g.ETag)
	}
	return GetEntityResponse{
		ETag:  ETag,
		Value: marshalledValue,
	}, nil
}

// GetEntity retrieves a specific entity from the service using the specified partitionKey and rowKey values. If
// no entity is available it returns an error. If the service returns a non-successful HTTP status code, the function
// returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *Client) GetEntity(ctx context.Context, partitionKey string, rowKey string, options *GetEntityOptions) (GetEntityResponse, error) {
	if options == nil {
		options = &GetEntityOptions{}
	}

	genOptions, queryOptions := options.toGenerated()
	resp, err := t.client.QueryEntityWithPartitionAndRowKey(ctx, generated.Enum1Three0, t.name, partitionKey, rowKey, genOptions, queryOptions)
	if err != nil {
		return GetEntityResponse{}, err
	}
	return newGetEntityResponse(resp)
}

// AddEntityOptions contains optional parameters for Client.AddEntity
type AddEntityOptions struct {
	// placeholder for future optional parameters
}

// AddEntityResponse contains response fields for Client.AddEntityResponse
type AddEntityResponse struct {
	ETag azcore.ETag
}

func addEntityResponseFromGenerated(g *generated.TableClientInsertEntityResponse) AddEntityResponse {
	if g == nil {
		return AddEntityResponse{}
	}

	var ETag azcore.ETag
	if g.ETag != nil {
		ETag = azcore.ETag(*g.ETag)
	}
	return AddEntityResponse{
		ETag: ETag,
	}
}

// AddEntity adds an entity (described by a byte slice) to the table. This method returns an error if an entity with
// the same PartitionKey and RowKey already exists in the table. If the supplied entity does not contain both a PartitionKey
// and a RowKey an error will be returned. If the service returns a non-successful HTTP status code, the function returns
// an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *Client) AddEntity(ctx context.Context, entity []byte, options *AddEntityOptions) (AddEntityResponse, error) {
	var mapEntity map[string]interface{}
	err := json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return AddEntityResponse{}, err
	}
	resp, err := t.client.InsertEntity(ctx, generated.Enum1Three0, t.name, &generated.TableClientInsertEntityOptions{TableEntityProperties: mapEntity, ResponsePreference: to.Ptr(generated.ResponseFormatReturnNoContent)}, nil)
	if err != nil {
		err = checkEntityForPkRk(&mapEntity, err)
		return AddEntityResponse{}, err
	}
	return addEntityResponseFromGenerated(&resp), err
}

// DeleteEntityOptions contains optional parameters for Client.DeleteEntity
type DeleteEntityOptions struct {
	IfMatch *azcore.ETag
}

func (d *DeleteEntityOptions) toGenerated() *generated.TableClientDeleteEntityOptions {
	return &generated.TableClientDeleteEntityOptions{}
}

// DeleteEntityResponse contains response fields for Client.DeleteEntity
type DeleteEntityResponse struct {
	// placeholder for future optional response fields
}

func deleteEntityResponseFromGenerated(g *generated.TableClientDeleteEntityResponse) DeleteEntityResponse {
	return DeleteEntityResponse{}
}

// DeleteEntity deletes the entity with the specified partitionKey and rowKey from the table. If the service returns a non-successful HTTP
// status code, the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *Client) DeleteEntity(ctx context.Context, partitionKey string, rowKey string, options *DeleteEntityOptions) (DeleteEntityResponse, error) {
	if options == nil {
		options = &DeleteEntityOptions{}
	}
	if options.IfMatch == nil {
		nilEtag := azcore.ETag("*")
		options.IfMatch = &nilEtag
	}
	resp, err := t.client.DeleteEntity(ctx, generated.Enum1Three0, t.name, partitionKey, rowKey, string(*options.IfMatch), options.toGenerated(), &generated.QueryOptions{})
	if err != nil {
		return DeleteEntityResponse{}, err
	}
	return deleteEntityResponseFromGenerated(&resp), err
}

// UpdateEntityOptions contains optional parameters for Client.UpdateEntity
type UpdateEntityOptions struct {
	IfMatch    *azcore.ETag
	UpdateMode UpdateMode
}

func (u *UpdateEntityOptions) toGeneratedMergeEntity(m map[string]interface{}) *generated.TableClientMergeEntityOptions {
	if u == nil {
		return &generated.TableClientMergeEntityOptions{}
	}
	return &generated.TableClientMergeEntityOptions{
		IfMatch:               to.Ptr(string(*u.IfMatch)),
		TableEntityProperties: m,
	}
}

func (u *UpdateEntityOptions) toGeneratedUpdateEntity(m map[string]interface{}) *generated.TableClientUpdateEntityOptions {
	if u == nil {
		return &generated.TableClientUpdateEntityOptions{}
	}
	return &generated.TableClientUpdateEntityOptions{
		IfMatch:               to.Ptr(string(*u.IfMatch)),
		TableEntityProperties: m,
	}
}

// UpdateEntityResponse contains response fields for Client.UpdateEntity
type UpdateEntityResponse struct {
	ETag azcore.ETag
}

func updateEntityResponseFromMergeGenerated(g *generated.TableClientMergeEntityResponse) UpdateEntityResponse {
	if g == nil {
		return UpdateEntityResponse{}
	}

	var ETag azcore.ETag
	if g.ETag != nil {
		ETag = azcore.ETag(*g.ETag)
	}
	return UpdateEntityResponse{
		ETag: ETag,
	}
}

func updateEntityResponseFromUpdateGenerated(g *generated.TableClientUpdateEntityResponse) UpdateEntityResponse {
	if g == nil {
		return UpdateEntityResponse{}
	}

	var ETag azcore.ETag
	if g.ETag != nil {
		ETag = azcore.ETag(*g.ETag)
	}
	return UpdateEntityResponse{
		ETag: ETag,
	}
}

// UpdateEntity updates the specified table entity if it exists.
// If updateMode is Replace, the entity will be replaced. This is the only way to remove properties from an existing entity.
// If updateMode is Merge, the property values present in the specified entity will be merged with the existing entity. Properties not specified in the merge will be unaffected.
// The specified etag value will be used for optimistic concurrency. If the etag does not match the value of the entity in the table, the operation will fail.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
// If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *Client) UpdateEntity(ctx context.Context, entity []byte, options *UpdateEntityOptions) (UpdateEntityResponse, error) {
	if options == nil {
		options = &UpdateEntityOptions{
			UpdateMode: UpdateModeMerge,
		}
	}

	if options.IfMatch == nil {
		star := azcore.ETag("*")
		options.IfMatch = &star
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

	switch options.UpdateMode {
	case UpdateModeMerge:
		resp, err := t.client.MergeEntity(
			ctx,
			generated.Enum1Three0,
			t.name,
			partKey,
			rowkey,
			options.toGeneratedMergeEntity(mapEntity),
			&generated.QueryOptions{},
		)
		if err != nil {
			return UpdateEntityResponse{}, err
		}
		return updateEntityResponseFromMergeGenerated(&resp), err
	case UpdateModeReplace:
		resp, err := t.client.UpdateEntity(
			ctx,
			generated.Enum1Three0,
			t.name,
			partKey,
			rowkey,
			options.toGeneratedUpdateEntity(mapEntity),
			&generated.QueryOptions{},
		)
		if err != nil {
			return UpdateEntityResponse{}, err
		}
		return updateEntityResponseFromUpdateGenerated(&resp), err
	}
	if pk == "" || rk == "" {
		return UpdateEntityResponse{}, errPartitionKeyRowKeyError
	}
	return UpdateEntityResponse{}, errInvalidUpdateMode
}

// UpsertEntityOptions contains optional parameters for Client.InsertEntity
type UpsertEntityOptions struct {
	// ETag is the optional etag for the Table
	ETag azcore.ETag

	// UpdateMode is the desired mode for the Update. Use UpdateModeReplace to replace fields on
	// the entity, use UpdateModeMerge to merge fields of the entity.
	UpdateMode UpdateMode
}

// UpsertEntityResponse contains response fields for Client.InsertEntity
type UpsertEntityResponse struct {
	ETag azcore.ETag
}

func insertEntityFromGeneratedMerge(g *generated.TableClientMergeEntityResponse) UpsertEntityResponse {
	if g == nil {
		return UpsertEntityResponse{}
	}

	var ETag azcore.ETag
	if g.ETag != nil {
		ETag = azcore.ETag(*g.ETag)
	}
	return UpsertEntityResponse{
		ETag: ETag,
	}
}

func insertEntityFromGeneratedUpdate(g *generated.TableClientUpdateEntityResponse) UpsertEntityResponse {
	if g == nil {
		return UpsertEntityResponse{}
	}

	var ETag azcore.ETag
	if g.ETag != nil {
		ETag = azcore.ETag(*g.ETag)
	}
	return UpsertEntityResponse{
		ETag: ETag,
	}
}

// UpsertEntity inserts an entity if it does not already exist in the table. If the entity does exist, the entity is
// replaced or merged as specified the updateMode parameter. If the entity exists and updateMode is Merge, the property
// values present in the specified entity will be merged with the existing entity rather than replaced.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
// If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *Client) UpsertEntity(ctx context.Context, entity []byte, options *UpsertEntityOptions) (UpsertEntityResponse, error) {
	if options == nil {
		options = &UpsertEntityOptions{
			UpdateMode: UpdateModeMerge,
		}
	}
	var mapEntity map[string]interface{}
	err := json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return UpsertEntityResponse{}, err
	}

	pk := mapEntity[partitionKey]
	partKey := pk.(string)

	rk := mapEntity[rowKey]
	rowkey := rk.(string)

	switch options.UpdateMode {
	case UpdateModeMerge:
		resp, err := t.client.MergeEntity(
			ctx,
			generated.Enum1Three0,
			t.name,
			partKey,
			rowkey,
			&generated.TableClientMergeEntityOptions{TableEntityProperties: mapEntity},
			&generated.QueryOptions{},
		)
		if err != nil {
			return UpsertEntityResponse{}, err
		}
		return insertEntityFromGeneratedMerge(&resp), err
	case UpdateModeReplace:
		resp, err := t.client.UpdateEntity(
			ctx,
			generated.Enum1Three0,
			t.name,
			partKey,
			rowkey,
			&generated.TableClientUpdateEntityOptions{TableEntityProperties: mapEntity},
			&generated.QueryOptions{},
		)
		if err != nil {
			return UpsertEntityResponse{}, err
		}
		return insertEntityFromGeneratedUpdate(&resp), err
	}
	if pk == "" || rk == "" {
		return UpsertEntityResponse{}, errPartitionKeyRowKeyError
	}
	return UpsertEntityResponse{}, errInvalidUpdateMode
}

// GetAccessPolicyOptions contains optional parameters for Client.GetAccessPolicy
type GetAccessPolicyOptions struct {
	// placeholder for future optional parameters
}

func (g *GetAccessPolicyOptions) toGenerated() *generated.TableClientGetAccessPolicyOptions {
	return &generated.TableClientGetAccessPolicyOptions{}
}

// GetAccessPolicyResponse contains response fields for Client.GetAccessPolicy
type GetAccessPolicyResponse struct {
	SignedIdentifiers []*SignedIdentifier
}

func getAccessPolicyResponseFromGenerated(g *generated.TableClientGetAccessPolicyResponse) GetAccessPolicyResponse {
	if g == nil {
		return GetAccessPolicyResponse{}
	}

	var sis []*SignedIdentifier
	for _, s := range g.SignedIdentifiers {
		sis = append(sis, fromGeneratedSignedIdentifier(s))
	}
	return GetAccessPolicyResponse{
		SignedIdentifiers: sis,
	}
}

// GetAccessPolicy retrieves details about any stored access policies specified on the table that may be used with the Shared Access Signature.
// If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *Client) GetAccessPolicy(ctx context.Context, options *GetAccessPolicyOptions) (GetAccessPolicyResponse, error) {
	resp, err := t.client.GetAccessPolicy(ctx, t.name, generated.Enum4ACL, options.toGenerated())
	if err != nil {
		return GetAccessPolicyResponse{}, err
	}
	return getAccessPolicyResponseFromGenerated(&resp), err
}

// SetAccessPolicyOptions contains optional parameters for Client.SetAccessPolicy
type SetAccessPolicyOptions struct {
	TableACL []*SignedIdentifier
}

// SetAccessPolicyResponse contains response fields for Client.SetAccessPolicy
type SetAccessPolicyResponse struct {
	// placeholder for future optional parameters
}

func setAccessPolicyResponseFromGenerated(g *generated.TableClientSetAccessPolicyResponse) SetAccessPolicyResponse {
	return SetAccessPolicyResponse{}
}

func (s *SetAccessPolicyOptions) toGenerated() *generated.TableClientSetAccessPolicyOptions {
	var sis []*generated.SignedIdentifier
	for _, t := range s.TableACL {
		sis = append(sis, toGeneratedSignedIdentifier(t))
	}
	return &generated.TableClientSetAccessPolicyOptions{
		TableACL: sis,
	}
}

// SetAccessPolicy sets stored access policies for the table that may be used with SharedAccessSignature.
// If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *Client) SetAccessPolicy(ctx context.Context, options *SetAccessPolicyOptions) (SetAccessPolicyResponse, error) {
	if options == nil {
		options = &SetAccessPolicyOptions{}
	}
	response, err := t.client.SetAccessPolicy(ctx, t.name, generated.Enum4ACL, options.toGenerated())
	if err != nil && len(options.TableACL) > 5 {
		err = errTooManyAccessPoliciesError
	}
	if err != nil {
		return SetAccessPolicyResponse{}, err
	}
	return setAccessPolicyResponseFromGenerated(&response), err
}

// GetTableSASURL is a convenience method for generating a SAS token for a specific table.
// It can only be used by clients created by NewClientWithSharedKey().
func (t Client) GetTableSASURL(permissions SASPermissions, start time.Time, expiry time.Time) (string, error) {
	if t.cred == nil {
		return "", errors.New("SAS can only be signed with a SharedKeyCredential")
	}
	qps, err := SASSignatureValues{
		TableName:         t.name,
		Permissions:       permissions.String(),
		StartTime:         start,
		ExpiryTime:        expiry,
		StartPartitionKey: permissions.StartPartitionKey,
		StartRowKey:       permissions.StartRowKey,
		EndPartitionKey:   permissions.EndPartitionKey,
		EndRowKey:         permissions.EndRowKey,
	}.Sign(t.cred)
	if err != nil {
		return "", err
	}

	serviceURL := t.con.Endpoint()
	if !strings.Contains(serviceURL, "/") {
		serviceURL += "/"
	}
	serviceURL += t.name + "?" + qps
	return serviceURL, nil
}
