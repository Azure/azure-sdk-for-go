// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

func parseConnectionString(connStr string) (string, azcore.Credential, error) {
	var serviceURL string
	var cred azcore.Credential
	splitString := strings.Split(connStr, ";")
	var pairs [][]string
	for _, splitPair := range splitString {
		temp := strings.Split(splitPair, "=")
		if len(temp) != 2 {
			return serviceURL, cred, errors.New("Connection string is either blank or malformed")
		}
		pairs = append(pairs, temp)
	}

	pairsMap := make(map[string]string)
	for _, pair := range pairs {
		pairsMap[pair[0]] = pair[1]
	}

	var accountName string
	var accountKey string
	if value, ok := pairsMap["accountname"]; ok {
		accountName = value
	}
	if value, ok := pairsMap["accountkey"]; ok {
		accountKey = value
	}

	if accountName == "" || accountKey == "" {
		// Try sharedaccesssignature
		sharedAccessSignature, ok := pairsMap["sharedaccesssignature"]
		if !ok {
			return serviceURL, cred, errors.New("Connection string missing required connection details")
		}
		cred = azcore.SharedAccessSignature(sharedAccessSignature)
	}

	cred = &SharedKeyCredential{
		accountName: accountName,
		accountKey:  accountKey,
	}

	primary, okPrimary := pairsMap["tableendpoint"]
	secondary, okSecondary := pairsMap["tablesecondaryendpoint"]
	if !okPrimary {
		if okSecondary {
			return serviceURL, cred, errors.New("Connection string specifies only secondary connection")
		}
		if endpointsProtocol, ok := pairsMap["defaultendpointsprotocol"]; ok {
			if accountName, ok := pairsMap["accountname"]; ok {
				if endpointSuffix, ok := pairsMap["endpointsuffix"]; ok {
					primary = fmt.Sprintf("%v://%v.table.%v", endpointsProtocol, accountName, endpointSuffix)
					secondary = fmt.Sprintf("%v-secondary.table.%v", accountName, endpointSuffix)
					okPrimary = true
					okSecondary = true
				}
			}
		}
	}

	if !okPrimary {

	}

	if serviceURL, ok = pairsMap["tableendpoint"]; !ok {
		return serviceURL, cred, errors.New("Connection string does not specify")
	}

	return serviceURL, cred, nil
}

func NewTableClientFromConnectionString(tableName string, connectionString string, options *TableClientOptions) (*TableClient, error) {
	endpoint, credential := parseConnectionString(connectionString)
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
// pager := client.Query(QueryOptions{})
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.sprintf("The page contains %i results", len(resp.TableEntityQueryResponse.Value))
// }
// err := pager.Err()
func (t *TableClient) Query(queryOptions QueryOptions) TableEntityQueryResponsePager {
	return &tableEntityQueryResponsePager{tableClient: t, queryOptions: &queryOptions, tableQueryOptions: &TableQueryEntitiesOptions{}}
}

// GetEntity retrieves a specific entity from the service using the specified partitionKey and rowKey values.
func (t *TableClient) GetEntity(ctx context.Context, partitionKey string, rowKey string) (MapOfInterfaceResponse, error) {
	resp, err := t.client.QueryEntityWithPartitionAndRowKey(ctx, t.Name, partitionKey, rowKey, &TableQueryEntityWithPartitionAndRowKeyOptions{}, &QueryOptions{})
	if err != nil {
		return resp, err
	}
	castAndRemoveAnnotations(&resp.Value)
	return resp, err
}

// AddEntity adds an entity from an arbitrary interface value to the table.
// An entity must have at least a PartitionKey and RowKey property.
func (t *TableClient) AddEntity(ctx context.Context, entity interface{}) (TableInsertEntityResponse, error) {
	entmap, err := toMap(entity)
	if err != nil {
		return TableInsertEntityResponse{}, azcore.NewResponseError(err, nil)
	}
	resp, err := t.client.InsertEntity(ctx, t.Name, &TableInsertEntityOptions{TableEntityProperties: *entmap, ResponsePreference: ResponseFormatReturnNoContent.ToPtr()}, &QueryOptions{})
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
	return t.client.DeleteEntity(ctx, t.Name, partitionKey, rowKey, etag, nil, &QueryOptions{})
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
