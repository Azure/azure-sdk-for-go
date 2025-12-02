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
	generated "github.com/Azure/azure-sdk-for-go/sdk/data/aztables/internal"
)

// Client represents a client to the tables service affinitized to a specific table.
type Client struct {
	client  *generated.TableClient
	service *ServiceClient
	cred    *SharedKeyCredential
	name    string
}

// ClientOptions contains the optional parameters for client constructors.
type ClientOptions struct {
	azcore.ClientOptions
}

// NewClient creates a Client struct in the context of the table specified in the serviceURL, authorizing requests with an Azure AD access token.
// The serviceURL param is expected to have the name of the table in a format similar to: "https://myAccountName.table.core.windows.net/<myTableName>".
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
// The serviceURL param is expected to have the name of the table in a format similar to: "https://myAccountName.table.core.windows.net/<myTableName>?<SAS token>".
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
// The serviceURL param is expected to have the name of the table in a format similar to: "https://myAccountName.table.core.windows.net/<myTableName>".
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

// CreateTable creates the table with the tableName specified when NewClient was called. If the service returns a non-successful
// HTTP status code, the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
// NOTE: creating a table with the same name as a table that's in the process of being deleted will return an *azcore.ResponseError
// with error code TableBeingDeleted and status code http.StatusConflict.
func (t *Client) CreateTable(ctx context.Context, options *CreateTableOptions) (CreateTableResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.CreateTable", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &CreateTableOptions{}
	}
	resp, err := t.client.Create(ctx, generated.TableProperties{TableName: &t.name}, options.toGenerated(), &generated.QueryOptions{})
	if err != nil {
		return CreateTableResponse{}, err
	}
	return CreateTableResponse{
		TableName: resp.TableName,
	}, nil
}

// Delete deletes the table with the tableName specified when NewClient was called. If the service returns a non-successful HTTP status
// code, the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
// NOTE: deleting a table can take up to 40 seconds or more to complete.  If a table with the same name is created while the delete is still
// in progress, an *azcore.ResponseError is returned with error code TableBeingDeleted and status code http.StatusConflict.
func (t *Client) Delete(ctx context.Context, options *DeleteTableOptions) (DeleteTableResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.Delete", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	resp, err := t.service.DeleteTable(ctx, t.name, options)
	return resp, err
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
//   - API Documentation: https://learn.microsoft.com/rest/api/storageservices/querying-tables-and-entities
//   - README samples: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/data/aztables/README.md#writing-filters
func (t *Client) NewListEntitiesPager(listOptions *ListEntitiesOptions) *runtime.Pager[ListEntitiesResponse] {
	if listOptions == nil {
		listOptions = &ListEntitiesOptions{}
	}
	return runtime.NewPager(runtime.PagingHandler[ListEntitiesResponse]{
		More: func(page ListEntitiesResponse) bool {
			// if there are no continuation header values, there are no more pages
			// https://learn.microsoft.com/rest/api/storageservices/Query-Timeout-and-Pagination
			return (page.NextPartitionKey != nil && len(*page.NextPartitionKey) > 0) || (page.NextRowKey != nil && len(*page.NextRowKey) > 0)
		},
		Fetcher: func(ctx context.Context, page *ListEntitiesResponse) (ListEntitiesResponse, error) {
			var partKey *string
			var rowKey *string
			if page != nil {
				partKey = page.NextPartitionKey
				rowKey = page.NextRowKey
			} else {
				partKey = listOptions.NextPartitionKey
				rowKey = listOptions.NextRowKey
			}
			resp, err := t.client.QueryEntities(ctx, t.name, &generated.TableClientQueryEntitiesOptions{
				NextPartitionKey: partKey,
				NextRowKey:       rowKey,
			}, listOptions.toQueryOptions())
			if err != nil {
				return ListEntitiesResponse{}, err
			}

			var marshalledValue [][]byte
			if len(resp.Value) > 0 {
				marshalledValue = make([][]byte, len(resp.Value))
				for i := range resp.Value {
					m, err := json.Marshal(resp.Value[i])
					if err != nil {
						return ListEntitiesResponse{}, err
					}
					marshalledValue[i] = m
				}
			}

			return ListEntitiesResponse{
				NextPartitionKey: resp.XMSContinuationNextPartitionKey,
				NextRowKey:       resp.XMSContinuationNextRowKey,
				Entities:         marshalledValue,
			}, nil
		},
		Tracer: t.client.Tracer(),
	})
}

// GetEntity retrieves a specific entity from the service using the specified partitionKey and rowKey values. If
// no entity is available it returns an error. If the service returns a non-successful HTTP status code, the function
// returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *Client) GetEntity(ctx context.Context, partitionKey string, rowKey string, options *GetEntityOptions) (GetEntityResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.GetEntity", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &GetEntityOptions{}
	}

	resp, err := t.client.QueryEntityWithPartitionAndRowKey(ctx, t.name, prepareKey(partitionKey), prepareKey(rowKey), nil, &generated.QueryOptions{
		Format: options.Format,
	})
	if err != nil {
		return GetEntityResponse{}, err
	}
	marshalledValue, err := json.Marshal(resp.Value)
	if err != nil {
		return GetEntityResponse{}, err
	}

	var ETag azcore.ETag
	if resp.ETag != nil {
		ETag = azcore.ETag(*resp.ETag)
	}
	return GetEntityResponse{
		ETag:  ETag,
		Value: marshalledValue,
	}, nil
}

// AddEntity adds an entity (described by a byte slice) to the table. This method returns an error if an entity with
// the same PartitionKey and RowKey already exists in the table. If the supplied entity does not contain both a PartitionKey
// and a RowKey an error will be returned. If the service returns a non-successful HTTP status code, the function returns
// an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *Client) AddEntity(ctx context.Context, entity []byte, options *AddEntityOptions) (AddEntityResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.AddEntity", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	var mapEntity map[string]any
	err = json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return AddEntityResponse{}, err
	}

	if options == nil {
		options = &AddEntityOptions{}
	}

	resp, err := t.client.InsertEntity(ctx, t.name, &generated.TableClientInsertEntityOptions{TableEntityProperties: mapEntity}, &generated.QueryOptions{
		Format: options.Format,
	})
	if err != nil {
		err = checkEntityForPkRk(&mapEntity, err)
		return AddEntityResponse{}, err
	}
	marshalledValue, err := json.Marshal(resp.Value)
	if err != nil {
		return AddEntityResponse{}, err
	}

	var ETag azcore.ETag
	if resp.ETag != nil {
		ETag = azcore.ETag(*resp.ETag)
	}
	return AddEntityResponse{
		ETag:  ETag,
		Value: marshalledValue,
	}, nil
}

// DeleteEntity deletes the entity with the specified partitionKey and rowKey from the table. If the service returns a non-successful HTTP
// status code, the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *Client) DeleteEntity(ctx context.Context, partitionKey string, rowKey string, options *DeleteEntityOptions) (DeleteEntityResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.DeleteEntity", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &DeleteEntityOptions{}
	}
	if options.IfMatch == nil {
		nilEtag := azcore.ETag("*")
		options.IfMatch = &nilEtag
	}
	_, err = t.client.DeleteEntity(ctx, t.name, prepareKey(partitionKey), prepareKey(rowKey), string(*options.IfMatch), options.toGenerated(), &generated.QueryOptions{})
	return DeleteEntityResponse{}, err
}

// UpdateEntity updates the specified table entity if it exists.
// If updateMode is Replace, the entity will be replaced. This is the only way to remove properties from an existing entity.
// If updateMode is Merge, the property values present in the specified entity will be merged with the existing entity. Properties not specified in the merge will be unaffected.
// The specified etag value will be used for optimistic concurrency. If the etag does not match the value of the entity in the table, the operation will fail.
// The response type will be TableEntityMergeResponse if updateMode is Merge and TableEntityUpdateResponse if updateMode is Replace.
// If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *Client) UpdateEntity(ctx context.Context, entity []byte, options *UpdateEntityOptions) (UpdateEntityResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.UpdateEntity", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &UpdateEntityOptions{
			UpdateMode: UpdateModeMerge,
		}
	}

	if options.IfMatch == nil {
		star := azcore.ETag("*")
		options.IfMatch = &star
	}

	var mapEntity map[string]any
	err = json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return UpdateEntityResponse{}, err
	}

	pk := mapEntity[partitionKey]
	partKey := pk.(string)

	rk := mapEntity[rowKey]
	rowkey := rk.(string)

	switch options.UpdateMode {
	case UpdateModeMerge:
		var resp generated.TableClientMergeEntityResponse
		resp, err = t.client.MergeEntity(
			ctx,
			t.name,
			prepareKey(partKey),
			prepareKey(rowkey),
			options.toGeneratedMergeEntity(mapEntity),
			&generated.QueryOptions{},
		)
		if err != nil {
			return UpdateEntityResponse{}, err
		}
		var ETag azcore.ETag
		if resp.ETag != nil {
			ETag = azcore.ETag(*resp.ETag)
		}
		return UpdateEntityResponse{
			ETag: ETag,
		}, nil
	case UpdateModeReplace:
		var resp generated.TableClientUpdateEntityResponse
		resp, err = t.client.UpdateEntity(
			ctx,
			t.name,
			prepareKey(partKey),
			prepareKey(rowkey),
			options.toGeneratedUpdateEntity(mapEntity),
			&generated.QueryOptions{},
		)
		if err != nil {
			return UpdateEntityResponse{}, err
		}
		var ETag azcore.ETag
		if resp.ETag != nil {
			ETag = azcore.ETag(*resp.ETag)
		}
		return UpdateEntityResponse{
			ETag: ETag,
		}, nil
	}
	if pk == "" || rk == "" {
		err = errPartitionKeyRowKeyError
	} else {
		err = errInvalidUpdateMode
	}
	return UpdateEntityResponse{}, err
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
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.UpsertEntity", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &UpsertEntityOptions{
			UpdateMode: UpdateModeMerge,
		}
	}
	var mapEntity map[string]any
	err = json.Unmarshal(entity, &mapEntity)
	if err != nil {
		return UpsertEntityResponse{}, err
	}

	pk := mapEntity[partitionKey]
	partKey := pk.(string)

	rk := mapEntity[rowKey]
	rowkey := rk.(string)

	switch options.UpdateMode {
	case UpdateModeMerge:
		var resp generated.TableClientMergeEntityResponse
		resp, err = t.client.MergeEntity(
			ctx,
			t.name,
			prepareKey(partKey),
			prepareKey(rowkey),
			&generated.TableClientMergeEntityOptions{TableEntityProperties: mapEntity},
			&generated.QueryOptions{},
		)
		if err != nil {
			return UpsertEntityResponse{}, err
		}
		return insertEntityFromGeneratedMerge(&resp), err
	case UpdateModeReplace:
		var resp generated.TableClientUpdateEntityResponse
		resp, err = t.client.UpdateEntity(
			ctx,
			t.name,
			prepareKey(partKey),
			prepareKey(rowkey),
			&generated.TableClientUpdateEntityOptions{TableEntityProperties: mapEntity},
			&generated.QueryOptions{},
		)
		if err != nil {
			return UpsertEntityResponse{}, err
		}
		return insertEntityFromGeneratedUpdate(&resp), err
	}
	if pk == "" || rk == "" {
		err = errPartitionKeyRowKeyError
	} else {
		err = errInvalidUpdateMode
	}
	return UpsertEntityResponse{}, err
}

// GetAccessPolicy retrieves details about any stored access policies specified on the table that may be used with the Shared Access Signature.
// If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *Client) GetAccessPolicy(ctx context.Context, options *GetAccessPolicyOptions) (GetAccessPolicyResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.GetAccessPolicy", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	resp, err := t.client.GetAccessPolicy(ctx, t.name, options.toGenerated())
	if err != nil {
		return GetAccessPolicyResponse{}, err
	}
	if len(resp.SignedIdentifiers) == 0 {
		return GetAccessPolicyResponse{}, nil
	}
	sis := make([]*SignedIdentifier, len(resp.SignedIdentifiers))
	for i := range resp.SignedIdentifiers {
		sis[i] = fromGeneratedSignedIdentifier(resp.SignedIdentifiers[i])
	}
	return GetAccessPolicyResponse{
		SignedIdentifiers: sis,
	}, nil
}

// SetAccessPolicy sets stored access policies for the table that may be used with SharedAccessSignature.
// If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *Client) SetAccessPolicy(ctx context.Context, options *SetAccessPolicyOptions) (SetAccessPolicyResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "Client.SetAccessPolicy", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &SetAccessPolicyOptions{}
	}
	_, err = t.client.SetAccessPolicy(ctx, t.name, options.toGenerated())
	if err != nil && len(options.TableACL) > 5 {
		err = errTooManyAccessPoliciesError
	}
	return SetAccessPolicyResponse{}, err
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

	serviceURL := t.client.Endpoint()
	if !strings.Contains(serviceURL, "/") {
		serviceURL += "/"
	}
	serviceURL += t.name + "?" + qps
	return serviceURL, nil
}
