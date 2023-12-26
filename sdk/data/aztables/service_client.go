// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	generated "github.com/Azure/azure-sdk-for-go/sdk/data/aztables/internal"
)

// ServiceClient represents a client to the table service. It can be used to query
// the available tables, create/delete tables, and various other service level operations.
type ServiceClient struct {
	client  *generated.TableClient
	service *generated.ServiceClient
	cred    *SharedKeyCredential
}

// NewServiceClient creates a ServiceClient struct using the specified serviceURL, credential, and options.
// Pass in nil for options to construct the client with the default ClientOptions.
func NewServiceClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*ServiceClient, error) {
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{runtime.NewBearerTokenPolicy(cred, []string{"https://storage.azure.com/.default"}, nil)},
	}
	client, err := newClient(serviceURL, plOpts, options)
	if err != nil {
		return nil, err
	}
	return &ServiceClient{
		client:  generated.NewTableClient(serviceURL, client),
		service: generated.NewServiceClient(serviceURL, client),
	}, nil
}

// NewServiceClientWithNoCredential creates a ServiceClient struct using the specified serviceURL and options.
// Call this method when serviceURL contains a SAS token.
// Pass in nil for options to construct the client with the default ClientOptions.
func NewServiceClientWithNoCredential(serviceURL string, options *ClientOptions) (*ServiceClient, error) {
	client, err := newClient(serviceURL, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	return &ServiceClient{
		client:  generated.NewTableClient(serviceURL, client),
		service: generated.NewServiceClient(serviceURL, client),
	}, nil
}

// NewServiceClientWithSharedKey creates a ServiceClient struct using the specified serviceURL, credential, and options.
// Pass in nil for options to construct the client with the default ClientOptions.
func NewServiceClientWithSharedKey(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*ServiceClient, error) {
	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{newSharedKeyCredPolicy(cred)},
	}
	client, err := newClient(serviceURL, plOpts, options)
	if err != nil {
		return nil, err
	}
	return &ServiceClient{
		client:  generated.NewTableClient(serviceURL, client),
		service: generated.NewServiceClient(serviceURL, client),
		cred:    cred,
	}, nil
}

func newClient(serviceURL string, plOpts runtime.PipelineOptions, options *ClientOptions) (*azcore.Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	if isCosmosEndpoint(serviceURL) {
		plOpts.PerCall = append(plOpts.PerCall, cosmosPatchTransformPolicy{})
	}
	plOpts.Tracing.Namespace = "Microsoft.Tables"
	return azcore.NewClient(generated.ModuleName, generated.Version, plOpts, &options.ClientOptions)
}

// NewClient returns a pointer to a Client affinitized to the specified table name and initialized with the same serviceURL and credentials as this ServiceClient
func (t *ServiceClient) NewClient(tableName string) *Client {
	return &Client{
		client:  t.client,
		name:    tableName,
		service: t,
		cred:    t.cred,
	}
}

// CreateTableOptions contains optional parameters for Client.Create and ServiceClient.CreateTable
type CreateTableOptions struct {
	// placeholder for future optional parameters
}

func (c *CreateTableOptions) toGenerated() *generated.TableClientCreateOptions {
	return &generated.TableClientCreateOptions{}
}

// CreateTable creates a table with the specified name. If the service returns a non-successful HTTP status code,
// the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *ServiceClient) CreateTable(ctx context.Context, name string, options *CreateTableOptions) (CreateTableResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "ServiceClient.CreateTable", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &CreateTableOptions{}
	}
	resp, err := t.client.Create(ctx, generated.TableProperties{TableName: &name}, options.toGenerated(), &generated.QueryOptions{})
	if err != nil {
		return CreateTableResponse{}, err
	}
	return CreateTableResponse{
		TableName: resp.TableName,
	}, err
}

// DeleteTableOptions contains optional parameters for Client.Delete and ServiceClient.DeleteTable
type DeleteTableOptions struct {
	// placeholder for future optional parameters
}

func (c *DeleteTableOptions) toGenerated() *generated.TableClientDeleteOptions {
	return &generated.TableClientDeleteOptions{}
}

// DeleteTableResponse contains response fields for ServiceClient.DeleteTable and Client.Delete
type DeleteTableResponse struct {
	// placeholder for future optional response fields
}

func deleteTableResponseFromGen(g generated.TableClientDeleteResponse) DeleteTableResponse {
	return DeleteTableResponse{}
}

// DeleteTable deletes a table by name. If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *ServiceClient) DeleteTable(ctx context.Context, name string, options *DeleteTableOptions) (DeleteTableResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "ServiceClient.DeleteTable", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	resp, err := t.client.Delete(ctx, name, options.toGenerated())
	if err != nil {
		return DeleteTableResponse{}, err
	}
	return deleteTableResponseFromGen(resp), err
}

// ListTablesOptions contains optional parameters for ServiceClient.QueryTables
type ListTablesOptions struct {
	// OData filter expression.
	Filter *string

	// Select expression using OData notation. Limits the columns on each record to just those requested, e.g. "$select=PolicyAssignmentId, ResourceId".
	Select *string

	// Maximum number of records to return.
	Top *int32

	// NextTableName is the continuation token for the next table to page from
	NextTableName *string
}

func (l *ListTablesOptions) toQueryOptions() *generated.QueryOptions {
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

// ListTablesResponse contains response fields for ListTablesPager.NextPage
type ListTablesResponse struct {
	// NextTableName contains the information returned from the x-ms-continuation-NextTableName header response.
	NextTableName *string

	// List of tables.
	Tables []*TableProperties `json:"value,omitempty"`
}

func fromGeneratedTableQueryResponseEnvelope(g generated.TableClientQueryResponse) ListTablesResponse {
	var value []*TableProperties

	for _, v := range g.Value {
		value = append(value, fromGeneratedTableResponseProperties(v))
	}

	return ListTablesResponse{
		NextTableName: g.XMSContinuationNextTableName,
		Tables:        value,
	}
}

// TableProperties contains the properties for a single Table
type TableProperties struct {
	// The name of the table.
	Name *string `json:"TableName,omitempty"`
}

// Convets a generated TableResponseProperties to a ResponseProperties
func fromGeneratedTableResponseProperties(g *generated.TableResponseProperties) *TableProperties {
	if g == nil {
		return nil
	}

	return &TableProperties{
		Name: g.TableName,
	}
}

// NewListTablesPager queries the existing tables using the specified ListTablesOptions.
// listOptions can specify the following properties to affect the query results returned:
//
// Filter: An OData filter expression that limits results to those tables that satisfy the filter expression.
// For example, the following expression would return only tables with a TableName of 'foo': "TableName eq 'foo'"
//
// Top: The maximum number of tables that will be returned per page of results.
// Note: This value does not limit the total number of results if NextPage is called on the returned Pager until it returns false.
//
// NewListTablesPager returns a Pager, which allows iteration through each page of results. Specify nil for listOptions if you want to use the default options.
// For more information about writing query strings, check out:
//   - API Documentation: https://learn.microsoft.com/rest/api/storageservices/querying-tables-and-entities
//   - README samples: https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/data/aztables/README.md#writing-filters
func (t *ServiceClient) NewListTablesPager(listOptions *ListTablesOptions) *runtime.Pager[ListTablesResponse] {
	if listOptions == nil {
		listOptions = &ListTablesOptions{}
	}
	return runtime.NewPager(runtime.PagingHandler[ListTablesResponse]{
		More: func(page ListTablesResponse) bool {
			if page.NextTableName == nil || len(*page.NextTableName) == 0 {
				return false
			}
			return true
		},
		Fetcher: func(ctx context.Context, page *ListTablesResponse) (ListTablesResponse, error) {
			var tableName *string
			if page != nil {
				if page.NextTableName != nil {
					tableName = page.NextTableName
				}
			} else {
				tableName = listOptions.NextTableName
			}
			resp, err := t.client.Query(
				ctx,
				&generated.TableClientQueryOptions{NextTableName: tableName},
				listOptions.toQueryOptions())
			if err != nil {
				return ListTablesResponse{}, err
			}
			return fromGeneratedTableQueryResponseEnvelope(resp), nil
		},
		Tracer: t.client.Tracer(),
	})
}

// GetStatisticsOptions contains optional parameters for ServiceClient.GetStatistics
type GetStatisticsOptions struct {
	// placeholder for future optional parameters
}

// GetStatisticsResponse contains response fields for Client.GetStatistics
type GetStatisticsResponse struct {
	GeoReplication *GeoReplication `xml:"GeoReplication"`
}

func getStatisticsResponseFromGenerated(g *generated.ServiceClientGetStatisticsResponse) GetStatisticsResponse {
	return GetStatisticsResponse{
		GeoReplication: fromGeneratedGeoReplication(g.GeoReplication),
	}
}

func (g *GetStatisticsOptions) toGenerated() *generated.ServiceClientGetStatisticsOptions {
	return &generated.ServiceClientGetStatisticsOptions{}
}

// GetStatistics retrieves all the statistics for an account with Geo-redundancy established. If the service returns a non-successful
// HTTP status code, the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *ServiceClient) GetStatistics(ctx context.Context, options *GetStatisticsOptions) (GetStatisticsResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "ServiceClient.GetStatistics", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &GetStatisticsOptions{}
	}
	resp, err := t.service.GetStatistics(ctx, options.toGenerated())
	if err != nil {
		return GetStatisticsResponse{}, err
	}
	return getStatisticsResponseFromGenerated(&resp), err
}

// GetPropertiesOptions contains optional parameters for Client.GetProperties
type GetPropertiesOptions struct {
	// placeholder for future optional parameters
}

func (g *GetPropertiesOptions) toGenerated() *generated.ServiceClientGetPropertiesOptions {
	return &generated.ServiceClientGetPropertiesOptions{}
}

// GetPropertiesResponse contains response fields for Client.GetProperties
type GetPropertiesResponse struct {
	ServiceProperties
}

func getPropertiesResponseFromGenerated(g *generated.ServiceClientGetPropertiesResponse) GetPropertiesResponse {
	var cors []*CorsRule
	for _, c := range g.Cors {
		cors = append(cors, fromGeneratedCors(c))
	}
	return GetPropertiesResponse{
		ServiceProperties: ServiceProperties{
			Cors:          cors,
			HourMetrics:   fromGeneratedMetrics(g.HourMetrics),
			Logging:       fromGeneratedLogging(g.Logging),
			MinuteMetrics: fromGeneratedMetrics(g.MinuteMetrics),
		},
	}
}

// GetProperties retrieves the properties for an account including the metrics, logging, and cors rules established.
// If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *ServiceClient) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "ServiceClient.GetProperties", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &GetPropertiesOptions{}
	}
	resp, err := t.service.GetProperties(ctx, options.toGenerated())
	if err != nil {
		return GetPropertiesResponse{}, err
	}
	return getPropertiesResponseFromGenerated(&resp), err
}

// SetPropertiesOptions contains optional parameters for Client.SetProperties
type SetPropertiesOptions struct {
	// placeholder for future optional parameters
}

func (s *SetPropertiesOptions) toGenerated() *generated.ServiceClientSetPropertiesOptions {
	return &generated.ServiceClientSetPropertiesOptions{}
}

// SetPropertiesResponse contains response fields for Client.SetProperties
type SetPropertiesResponse struct {
	// placeholder for future response fields
}

func setPropertiesResponseFromGenerated(g *generated.ServiceClientSetPropertiesResponse) SetPropertiesResponse {
	return SetPropertiesResponse{}
}

// SetProperties allows the user to set cors, metrics, and logging rules for the account.
//
// Cors: A slice of CorsRules.
//
// HoursMetrics: A summary of request statistics grouped in hourly aggregatess for tables
//
// HoursMetrics: A summary of request statistics grouped in minute aggregates for tables
//
// Logging: Azure Analytics logging settings. If the service returns a non-successful HTTP
// status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *ServiceClient) SetProperties(ctx context.Context, properties ServiceProperties, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "ServiceClient.SetProperties", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	if options == nil {
		options = &SetPropertiesOptions{}
	}
	resp, err := t.service.SetProperties(ctx, *properties.toGenerated(), options.toGenerated())
	if err != nil {
		return SetPropertiesResponse{}, err
	}
	return setPropertiesResponseFromGenerated(&resp), err
}

// GetAccountSASURL is a convenience method for generating a SAS token for the currently pointed at account. This methods returns the full service URL and an error
// if there was an error during creation. This method can only be used by clients created by NewServiceClientWithSharedKey().
func (t ServiceClient) GetAccountSASURL(resources AccountSASResourceTypes, permissions AccountSASPermissions, start time.Time, expiry time.Time) (string, error) {
	if t.cred == nil {
		return "", errors.New("SAS can only be signed with a SharedKeyCredential")
	}
	qps, err := AccountSASSignatureValues{
		Version:       SASVersion,
		Protocol:      SASProtocolHTTPS,
		Permissions:   permissions.String(),
		Services:      "t",
		ResourceTypes: resources.String(),
		StartTime:     start.UTC(),
		ExpiryTime:    expiry.UTC(),
	}.Sign(t.cred)
	if err != nil {
		return "", err
	}
	endpoint := t.service.Endpoint()
	if !strings.HasSuffix(endpoint, "/") {
		endpoint += "/"
	}
	endpoint += "?" + qps.Encode()
	return endpoint, nil
}

func isCosmosEndpoint(url string) bool {
	isCosmosEmulator := strings.Contains(url, "localhost") && strings.Contains(url, "8902")
	return isCosmosEmulator || strings.Contains(url, cosmosTableDomain) || strings.Contains(url, legacyCosmosTableDomain)
}
