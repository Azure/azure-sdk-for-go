// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"errors"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	generated "github.com/Azure/azure-sdk-for-go/sdk/data/aztables/internal"
)

// A ServiceClient represents a client to the table service. It can be used to query the available tables, create/delete tables, and various other service level operations.
type ServiceClient struct {
	client  *generated.TableClient
	service *generated.ServiceClient
	cred    *SharedKeyCredential
	con     *generated.Connection
}

// NewServiceClient creates a ServiceClient struct using the specified serviceURL, credential, and options.
// Pass in nil for options to construct the client with the default ClientOptions.
func NewServiceClient(serviceURL string, cred azcore.TokenCredential, options *ClientOptions) (*ServiceClient, error) {
	conOptions := getConnectionOptions(serviceURL, options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, runtime.NewBearerTokenPolicy(cred, []string{"https://storage.azure.com/.default"}, nil))
	con := generated.NewConnection(serviceURL, conOptions)
	return &ServiceClient{
		client:  generated.NewTableClient(serviceURL, generated.Enum0TwoThousandNineteen0202, conOptions),
		service: generated.NewServiceClient(serviceURL, generated.Enum0TwoThousandNineteen0202, conOptions),
		con:     con,
	}, nil
}

// NewServiceClientWithNoCredential creates a ServiceClient struct using the specified serviceURL and options.
// Call this method when serviceURL contains a SAS token.
// Pass in nil for options to construct the client with the default ClientOptions.
func NewServiceClientWithNoCredential(serviceURL string, options *ClientOptions) (*ServiceClient, error) {
	conOptions := getConnectionOptions(serviceURL, options)
	con := generated.NewConnection(serviceURL, conOptions)
	return &ServiceClient{
		client:  generated.NewTableClient(serviceURL, generated.Enum0TwoThousandNineteen0202, conOptions),
		service: generated.NewServiceClient(serviceURL, generated.Enum0TwoThousandNineteen0202, conOptions),
		con:     con,
	}, nil
}

// NewServiceClientWithSharedKey creates a ServiceClient struct using the specified serviceURL, credential, and options.
// Pass in nil for options to construct the client with the default ClientOptions.
func NewServiceClientWithSharedKey(serviceURL string, cred *SharedKeyCredential, options *ClientOptions) (*ServiceClient, error) {
	conOptions := getConnectionOptions(serviceURL, options)
	conOptions.PerRetryPolicies = append(conOptions.PerRetryPolicies, newSharedKeyCredPolicy(cred))

	con := generated.NewConnection(serviceURL, conOptions)
	return &ServiceClient{
		client:  generated.NewTableClient(serviceURL, generated.Enum0TwoThousandNineteen0202, conOptions),
		service: generated.NewServiceClient(serviceURL, generated.Enum0TwoThousandNineteen0202, conOptions),
		cred:    cred,
		con:     con,
	}, nil
}

func getConnectionOptions(serviceURL string, options *ClientOptions) *policy.ClientOptions {
	if options == nil {
		options = &ClientOptions{}
	}
	conOptions := options.toPolicyOptions()
	if isCosmosEndpoint(serviceURL) {
		conOptions.PerCallPolicies = append(conOptions.PerCallPolicies, cosmosPatchTransformPolicy{})
	}
	return conOptions
}

// NewClient returns a pointer to a Client affinitized to the specified table name and initialized with the same serviceURL and credentials as this ServiceClient
func (t *ServiceClient) NewClient(tableName string) *Client {
	return &Client{
		client:  t.client,
		name:    tableName,
		service: t,
		con:     t.con,
		cred:    t.cred,
	}
}

// Options for Client.Create and ServiceClient.CreateTable method
type CreateTableOptions struct {
	// placeholder for future optional parameters
}

func (c *CreateTableOptions) toGenerated() *generated.TableClientCreateOptions {
	return &generated.TableClientCreateOptions{}
}

// Create creates a table with the specified name. If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *ServiceClient) CreateTable(ctx context.Context, name string, options *CreateTableOptions) (*Client, error) {
	if options == nil {
		options = &CreateTableOptions{}
	}
	_, err := t.client.Create(ctx, generated.Enum1Three0, generated.TableProperties{TableName: &name}, options.toGenerated(), &generated.QueryOptions{})
	return t.NewClient(name), err
}

// Options for Client.Delete and ServiceClient.DeleteTable methods
type DeleteTableOptions struct {
	// placeholder for future optional parameters
}

func (c *DeleteTableOptions) toGenerated() *generated.TableClientDeleteOptions {
	return &generated.TableClientDeleteOptions{}
}

// Response object from a ServiceClient.DeleteTable or Client.Delete operation
type DeleteTableResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func deleteTableResponseFromGen(g *generated.TableClientDeleteResponse) DeleteTableResponse {
	if g == nil {
		return DeleteTableResponse{}
	}
	return DeleteTableResponse{
		RawResponse: g.RawResponse,
	}
}

// Delete deletes a table by name. If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *ServiceClient) DeleteTable(ctx context.Context, name string, options *DeleteTableOptions) (DeleteTableResponse, error) {
	if options == nil {
		options = &DeleteTableOptions{}
	}
	resp, err := t.client.Delete(ctx, name, options.toGenerated())
	return deleteTableResponseFromGen(&resp), err
}

// ListEntitiesOptions contains a group of parameters for the ServiceClient.QueryTables method.
type ListTablesOptions struct {
	// OData filter expression.
	Filter *string

	// Select expression using OData notation. Limits the columns on each record to just those requested, e.g. "$select=PolicyAssignmentId, ResourceId".
	Select *string

	// Maximum number of records to return.
	Top *int32
}

func (l *ListTablesOptions) toQueryOptions() *generated.QueryOptions {
	if l == nil {
		return &generated.QueryOptions{}
	}

	return &generated.QueryOptions{
		Filter: l.Filter,
		Format: generated.ODataMetadataFormatApplicationJSONODataMinimalmetadata.ToPtr(),
		Select: l.Select,
		Top:    l.Top,
	}
}

// ListTablesPageResponse contains the properties of a single page response from a ListTables operation
type ListTablesPageResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// ContinuationNextTableName contains the information returned from the x-ms-continuation-NextTableName header response.
	ContinuationNextTableName *string

	// The metadata response of the table.
	ODataMetadata *string `json:"odata.metadata,omitempty"`

	// List of tables.
	Tables []*TableProperties `json:"value,omitempty"`
}

func fromGeneratedTableQueryResponseEnvelope(g generated.TableClientQueryResponse) ListTablesPageResponse {
	var value []*TableProperties

	for _, v := range g.Value {
		value = append(value, fromGeneratedTableResponseProperties(v))
	}

	return ListTablesPageResponse{
		RawResponse:               g.RawResponse,
		ContinuationNextTableName: g.XMSContinuationNextTableName,
		ODataMetadata:             g.ODataMetadata,
		Tables:                    value,
	}
}

// TableProperties contains the properties for a single Table
type TableProperties struct {
	// The edit link of the table.
	ODataEditLink *string `json:"odata.editLink,omitempty"`

	// The ID of the table.
	ODataID *string `json:"odata.id,omitempty"`

	// The odata type of the table.
	ODataType *string `json:"odata.type,omitempty"`

	// The name of the table.
	Name *string `json:"TableName,omitempty"`
}

// Convets a generated TableResponseProperties to a ResponseProperties
func fromGeneratedTableResponseProperties(g *generated.TableResponseProperties) *TableProperties {
	if g == nil {
		return nil
	}

	return &TableProperties{
		Name:          g.TableName,
		ODataEditLink: g.ODataEditLink,
		ODataID:       g.ODataID,
		ODataType:     g.ODataType,
	}
}

// ListTablesPager is a Pager for Table List operations
//
// Call NextPage first to fetch the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on the Pager.
// If the result is false, the value of Err() will indicate if an error occurred.
//
// PageResponse returns the results from the page most recently fetched from the service.
type ListTablesPager struct {
	client            *generated.TableClient
	current           generated.TableClientQueryResponse
	tableQueryOptions *generated.TableClientQueryOptions
	listOptions       *ListTablesOptions
	nextTableName     *string
}

// NextPage fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaulated by calling PageResponse on this Pager.
func (p *ListTablesPager) NextPage(ctx context.Context) (ListTablesPageResponse, error) {
	req, err := p.client.QueryCreateRequest(ctx, generated.Enum1Three0, &generated.TableClientQueryOptions{
		NextTableName: p.nextTableName,
	}, p.listOptions.toQueryOptions())
	if err != nil {
		return ListTablesPageResponse{}, err
	}
	resp, err := p.client.Pl.Do(req)
	if err != nil {
		return ListTablesPageResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ListTablesPageResponse{}, runtime.NewResponseError(resp)
	}

	result, err := p.client.QueryHandleResponse(resp)
	if err != nil {
		return ListTablesPageResponse{}, err
	}
	p.current = result
	p.nextTableName = p.current.XMSContinuationNextTableName
	return fromGeneratedTableQueryResponseEnvelope(p.current), nil
}

// More returns true if there are more pages to retrieve
func (p *ListTablesPager) More() bool {
	if !reflect.ValueOf(p.current).IsZero() {
		if p.current.XMSContinuationNextTableName == nil || len(*p.current.XMSContinuationNextTableName) == 0 {
			return false
		}
	}
	return true
}

// List queries the existing tables using the specified ListTablesOptions.
// listOptions can specify the following properties to affect the query results returned:
//
// Filter: An OData filter expression that limits results to those tables that satisfy the filter expression.
// For example, the following expression would return only tables with a TableName of 'foo': "TableName eq 'foo'"
//
// Top: The maximum number of tables that will be returned per page of results.
// Note: This value does not limit the total number of results if NextPage is called on the returned Pager until it returns false.
//
// List returns a Pager, which allows iteration through each page of results. Specify nil for listOptions if you want to use the default options.
func (t *ServiceClient) ListTables(listOptions *ListTablesOptions) ListTablesPager {
	return ListTablesPager{
		client:            t.client,
		tableQueryOptions: &generated.TableClientQueryOptions{},
		listOptions:       listOptions,
	}
}

// GetStatisticsOptions are the options for a ServiceClient.GetStatistics call
type GetStatisticsOptions struct {
	// placeholder for future optional parameters
}

type GetStatisticsResponse struct {
	RawResponse    *http.Response
	GeoReplication *GeoReplication `xml:"GeoReplication"`
}

func getStatisticsResponseFromGenerated(g *generated.ServiceClientGetStatisticsResponse) GetStatisticsResponse {
	return GetStatisticsResponse{
		RawResponse:    g.RawResponse,
		GeoReplication: fromGeneratedGeoReplication(g.GeoReplication),
	}
}

func (g *GetStatisticsOptions) toGenerated() *generated.ServiceClientGetStatisticsOptions {
	return &generated.ServiceClientGetStatisticsOptions{}
}

// GetStatistics retrieves all the statistics for an account with Geo-redundancy established. If the service returns a non-successful
// HTTP status code, the function returns an *azcore.ResponseError type. Specify nil for options if you want to use the default options.
func (t *ServiceClient) GetStatistics(ctx context.Context, options *GetStatisticsOptions) (GetStatisticsResponse, error) {
	if options == nil {
		options = &GetStatisticsOptions{}
	}
	resp, err := t.service.GetStatistics(ctx, generated.Enum5Service, generated.Enum7Stats, options.toGenerated())
	return getStatisticsResponseFromGenerated(&resp), err
}

// GetPropertiesOptions contains the optional parameters for the Client.GetProperties function
type GetPropertiesOptions struct {
	// placeholder for future optional parameters
}

func (g *GetPropertiesOptions) toGenerated() *generated.ServiceClientGetPropertiesOptions {
	return &generated.ServiceClientGetPropertiesOptions{}
}

type GetPropertiesResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// The set of CORS rules.
	Cors []*CorsRule `xml:"Cors>CorsRule"`

	// A summary of request statistics grouped by API in hourly aggregates for tables.
	HourMetrics *Metrics `xml:"HourMetrics"`

	// Azure Analytics Logging settings.
	Logging *Logging `xml:"Logging"`

	// A summary of request statistics grouped by API in minute aggregates for tables.
	MinuteMetrics *Metrics `xml:"MinuteMetrics"`
}

func getPropertiesResponseFromGenerated(g *generated.ServiceClientGetPropertiesResponse) GetPropertiesResponse {
	var cors []*CorsRule
	for _, c := range g.Cors {
		cors = append(cors, fromGeneratedCors(c))
	}
	return GetPropertiesResponse{
		RawResponse:   g.RawResponse,
		Cors:          cors,
		HourMetrics:   fromGeneratedMetrics(g.HourMetrics),
		Logging:       fromGeneratedLogging(g.Logging),
		MinuteMetrics: fromGeneratedMetrics(g.MinuteMetrics),
	}
}

// GetProperties retrieves the properties for an account including the metrics, logging, and cors rules established.
// If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *ServiceClient) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	if options == nil {
		options = &GetPropertiesOptions{}
	}
	resp, err := t.service.GetProperties(ctx, generated.Enum5Service, generated.Enum6Properties, options.toGenerated())
	return getPropertiesResponseFromGenerated(&resp), err
}

// SetPropertiesOptions contains the optional parameters for the Client.SetProperties method.
type SetPropertiesOptions struct {
	// placeholder for future optional parameters
}

func (s *SetPropertiesOptions) toGenerated() *generated.ServiceClientSetPropertiesOptions {
	return &generated.ServiceClientSetPropertiesOptions{}
}

type SetPropertiesResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func setPropertiesResponseFromGenerated(g *generated.ServiceClientSetPropertiesResponse) SetPropertiesResponse {
	return SetPropertiesResponse{
		RawResponse: g.RawResponse,
	}
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
	if options == nil {
		options = &SetPropertiesOptions{}
	}
	resp, err := t.service.SetProperties(ctx, generated.Enum5Service, generated.Enum6Properties, *properties.toGenerated(), options.toGenerated())
	return setPropertiesResponseFromGenerated(&resp), err
}

// GetAccountSASToken is a convenience method for generating a SAS token for the currently pointed at account. This methods returns the full service URL and an error
// if there was an error during creation. This method can only be used by clients created by NewServiceClientWithSharedKey().
func (t ServiceClient) GetAccountSASToken(resources AccountSASResourceTypes, permissions AccountSASPermissions, start time.Time, expiry time.Time) (string, error) {
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
	endpoint := t.con.Endpoint()
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
