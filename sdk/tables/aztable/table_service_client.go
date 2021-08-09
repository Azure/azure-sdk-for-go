// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztable

import (
	"context"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

const (
	LegacyCosmosTableDomain = ".table.cosmosdb."
	CosmosTableDomain       = ".table.cosmos."
)

// A TableServiceClient represents a client to the table service. It can be used to query the available tables, add/remove tables, and various other service level operations.
type TableServiceClient struct {
	client  *tableClient
	service *serviceClient
	cred    azcore.Credential
}

// NewTableServiceClient creates a TableServiceClient struct using the specified serviceURL, credential, and options.
func NewTableServiceClient(serviceURL string, cred azcore.Credential, options *TableClientOptions) (*TableServiceClient, error) {
	if options == nil {
		options = &TableClientOptions{}
	}
	conOptions := options.getConnectionOptions()
	if isCosmosEndpoint(serviceURL) {
		conOptions.PerCallPolicies = []azcore.Policy{CosmosPatchTransformPolicy{}}
	}
	conOptions.PerCallPolicies = append(conOptions.PerCallPolicies, cred.AuthenticationPolicy(azcore.AuthenticationPolicyOptions{Options: azcore.TokenRequestOptions{Scopes: options.Scopes}}))
	for _, p := range options.PerCallOptions {
		conOptions.PerCallPolicies = append(conOptions.PerCallPolicies, p)
	}
	con := newConnection(serviceURL, conOptions)
	return &TableServiceClient{client: &tableClient{con}, service: &serviceClient{con}, cred: cred}, nil
}

// NewTableClient returns a pointer to a TableClient affinitzed to the specified table name and initialized with the same serviceURL and credentials as this TableServiceClient
func (t *TableServiceClient) NewTableClient(tableName string) *TableClient {
	return &TableClient{client: t.client, cred: t.cred, Name: tableName, service: t}
}

// Create creates a table with the specified name.
func (t *TableServiceClient) CreateTable(ctx context.Context, name string) (TableResponseResponse, error) {
	resp, err := t.client.Create(ctx, TableProperties{&name}, new(TableCreateOptions), new(QueryOptions))
	if err == nil {
		tableResp := resp.(TableResponseResponse)
		return tableResp, nil
	}
	return TableResponseResponse{}, err
}

// Delete deletes a table by name.
func (t *TableServiceClient) DeleteTable(ctx context.Context, name string, options *TableDeleteOptions) (TableDeleteResponse, error) {
	if options == nil {
		options = &TableDeleteOptions{}
	}
	return t.client.Delete(ctx, name, options)
}

// List queries the existing tables using the specified ListOptions.
// ListOptions can specify the following properties to affect the query results returned:
//
// Filter: An Odata filter expression that limits results to those tables that satisfy the filter expression.
// For example, the following expression would return only tables with a TableName of 'foo': "TableName eq 'foo'"
//
// Top: The maximum number of tables that will be returned per page of results.
// Note: This value does not limit the total number of results if NextPage is called on the returned Pager until it returns false.
//
// List returns a Pager, which allows iteration through each page of results. Example:
//
// options := &ListOptions{Filter: to.StringPtr("PartitionKey eq 'pk001'"), Top: to.Int32Ptr(25)}
// pager := client.List(options) // Pass in 'nil' if you want to return all Tables for an account.
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.Printf("The page contains %i results.\n", len(resp.TableQueryResponse.Value))
// }
// err := pager.Err()
func (t *TableServiceClient) ListTables(listOptions *ListOptions) TableQueryResponsePager {
	return &tableQueryResponsePager{
		client:            t.client,
		queryOptions:      listOptions,
		tableQueryOptions: new(TableQueryOptions),
	}
}

// GetStatistics retrieves all the statistics for an account with Geo-redundancy established.
//
// response, err := client.GetStatistics(context.Background, nil)
// handle(err)
// fmt.Println("Status: ", response.StorageServiceStats.GeoReplication.Status)
// fmt.Println(Last Sync Time: ", response.StorageServiceStats.GeoReplication.LastSyncTime)
func (t *TableServiceClient) GetStatistics(ctx context.Context, options *ServiceGetStatisticsOptions) (TableServiceStatsResponse, error) {
	if options == nil {
		options = &ServiceGetStatisticsOptions{}
	}
	return t.service.GetStatistics(ctx, options)
}

// GetProperties retrieves the properties for an account including the metrics, logging, and cors rules established.
//
// response, err := client.GetProperties(context.Background, nil)
// handle(err)
// fmt.Println(resopnse.StorageServiceStats.Cors)
// fmt.Println(resopnse.StorageServiceStats.HourMetrics)
// fmt.Println(resopnse.StorageServiceStats.Logging)
// fmt.Println(resopnse.StorageServiceStats.MinuteMetrics)
func (t *TableServiceClient) GetProperties(ctx context.Context, options *ServiceGetPropertiesOptions) (TableServicePropertiesResponse, error) {
	if options == nil {
		options = &ServiceGetPropertiesOptions{}
	}
	return t.service.GetProperties(ctx, options)
}

// SetProperties allows the user to set cors , metrics, and logging rules for the account.
//
// Cors: A slice of CorsRules.
//
// HoursMetrics: A summary of request statistics grouped in hourly aggregatess for tables
//
// HoursMetrics: A summary of request statistics grouped in minute aggregates for tables
//
// Logging: Azure Analytics logging settings
//
//
// logging := Logging{
// 		Read:    to.BoolPtr(true),
// 		Write:   to.BoolPtr(true),
// 		Delete:  to.BoolPtr(true),
// 		Version: to.StringPtr("1.0"),
// 		RetentionPolicy: &RetentionPolicy{
// 			Enabled: to.BoolPtr(true),
// 		Days:    to.Int32Ptr(5),
// 		},
// }
// props := TableServiceProperties{Logging: &logging}
// resp, err := context.client.SetProperties(ctx, props, nil)
// handle(err)
func (t *TableServiceClient) SetProperties(ctx context.Context, properties TableServiceProperties, options *ServiceSetPropertiesOptions) (ServiceSetPropertiesResponse, error) {
	if options == nil {
		options = &ServiceSetPropertiesOptions{}
	}
	return t.service.SetProperties(ctx, properties, options)
}

func (s TableServiceClient) CanGetAccountSASToken() bool {
	return s.cred != nil
}

// GetAccountSASToken is a convenience method for generating a SAS token for the currently pointed at account.
// It can only be used if the supplied azcore.Credential during creation was a SharedKeyCredential.
// This validity can be checked with CanGetAccountSASToken().
func (t TableServiceClient) GetAccountSASToken(resources AccountSASResourceTypes, permissions AccountSASPermissions, start, expiry time.Time) (SASQueryParameters, error) {
	return AccountSASSignatureValues{
		Version: SASVersion,

		Permissions:   permissions.String(),
		Services:      "t",
		ResourceTypes: resources.String(),

		StartTime:  start.UTC(),
		ExpiryTime: expiry.UTC(),
	}.NewSASQueryParameters(t.cred.(*SharedKeyCredential))
}

func isCosmosEndpoint(url string) bool {
	isCosmosEmulator := strings.Contains(url, "localhost") && strings.Contains(url, "8902")
	return isCosmosEmulator || strings.Contains(url, CosmosTableDomain) || strings.Contains(url, LegacyCosmosTableDomain)
}
