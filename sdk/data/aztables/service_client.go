// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	generated "github.com/Azure/azure-sdk-for-go/sdk/data/aztables/internal"
)

// A ServiceClient represents a client to the table service. It can be used to query the available tables, create/delete tables, and various other service level operations.
type ServiceClient struct {
	client  *generated.TableClient
	service *generated.ServiceClient
	cred    azcore.Credential
}

// NewServiceClient creates a ServiceClient struct using the specified serviceURL, credential, and options.
func NewServiceClient(serviceURL string, cred azcore.Credential, options *ClientOptions) (*ServiceClient, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	conOptions := options.getConnectionOptions()
	if isCosmosEndpoint(serviceURL) {
		conOptions.PerCallPolicies = append(conOptions.PerCallPolicies, cosmosPatchTransformPolicy{})
	}
	conOptions.PerCallPolicies = append(conOptions.PerCallPolicies, options.PerCallOptions...)
	con := generated.NewConnection(serviceURL, cred, conOptions)
	return &ServiceClient{
		client:  generated.NewTableClient(con),
		service: generated.NewServiceClient(con),
		cred:    cred,
	}, nil
}

// NewClient returns a pointer to a Client affinitized to the specified table name and initialized with the same serviceURL and credentials as this ServiceClient
func (t *ServiceClient) NewClient(tableName string) *Client {
	return &Client{
		client:  t.client,
		cred:    t.cred,
		name:    tableName,
		service: t,
	}
}

// Options for Client.Create and ServiceClient.CreateTable method
type CreateTableOptions struct {
}

func (c *CreateTableOptions) toGenerated() *generated.TableCreateOptions {
	return &generated.TableCreateOptions{}
}

// Create creates a table with the specified name.
func (t *ServiceClient) CreateTable(ctx context.Context, name string, options *CreateTableOptions) (*Client, error) {
	if options == nil {
		options = &CreateTableOptions{}
	}
	_, err := t.client.Create(ctx, generated.TableProperties{TableName: &name}, options.toGenerated(), &generated.QueryOptions{})
	return t.NewClient(name), err
}

// Options for Client.Delete and ServiceClient.DeleteTable methods
type DeleteTableOptions struct {
}

func (c *DeleteTableOptions) toGenerated() *generated.TableDeleteOptions {
	return &generated.TableDeleteOptions{}
}

// Response object from a ServiceClient.DeleteTable or Client.Delete operation
type DeleteTableResponse struct {
	RawResponse *http.Response
}

func deleteTableResponseFromGen(g *generated.TableDeleteResponse) DeleteTableResponse {
	if g == nil {
		return DeleteTableResponse{}
	}
	return DeleteTableResponse{
		RawResponse: g.RawResponse,
	}
}

// Delete deletes a table by name.
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

// ListTablesPager is a Pager for Table List operations
//
// NextPage should be called first. It fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaluated by calling PageResponse on this Pager.
// If the result is false, the value of Err() will indicate if an error occurred.
//
// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.Printf("The page contains %i results.\n", len(resp.Tables))
// }
// err := pager.Err()
type ListTablesPager interface {
	// PageResponse returns the current TableQueryResponseResponse.
	PageResponse() ListTablesPage
	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
	// Err returns an error if there was an error on the last request
	Err() error
}

// ListTablesPage contains the properties of a single page response from a ListTables operation
type ListTablesPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// ContinuationNextTableName contains the information returned from the x-ms-continuation-NextTableName header response.
	ContinuationNextTableName *string

	// The metadata response of the table.
	ODataMetadata *string `json:"odata.metadata,omitempty"`

	// List of tables.
	Tables []*ResponseProperties `json:"value,omitempty"`
}

func fromGeneratedTableQueryResponseEnvelope(g *generated.TableQueryResponseEnvelope) *ListTablesPage {
	if g == nil {
		return nil
	}

	var value []*ResponseProperties

	for _, v := range g.Value {
		value = append(value, fromGeneratedTableResponseProperties(v))
	}

	return &ListTablesPage{
		RawResponse:               g.RawResponse,
		ContinuationNextTableName: g.XMSContinuationNextTableName,
		ODataMetadata:             g.ODataMetadata,
		Tables:                    value,
	}
}

// ResponseProperties contains the properties for a single Table
type ResponseProperties struct {
	// The edit link of the table.
	ODataEditLink *string `json:"odata.editLink,omitempty"`

	// The id of the table.
	ODataID *string `json:"odata.id,omitempty"`

	// The odata type of the table.
	ODataType *string `json:"odata.type,omitempty"`

	// The name of the table.
	TableName *string `json:"TableName,omitempty"`
}

// Convets a generated TableResponseProperties to a ResponseProperties
func fromGeneratedTableResponseProperties(g *generated.TableResponseProperties) *ResponseProperties {
	if g == nil {
		return nil
	}

	return &ResponseProperties{
		TableName:     g.TableName,
		ODataEditLink: g.ODataEditLink,
		ODataID:       g.ODataID,
		ODataType:     g.ODataType,
	}
}

type tableQueryResponsePager struct {
	client            *generated.TableClient
	current           *generated.TableQueryResponseEnvelope
	tableQueryOptions *generated.TableQueryOptions
	listOptions       *ListTablesOptions
	err               error
}

// NextPage fetches the next available page of results from the service.
// If the fetched page contains results, the return value is true, else false.
// Results fetched from the service can be evaulated by calling PageResponse on this Pager.
func (p *tableQueryResponsePager) NextPage(ctx context.Context) bool {
	if p.err != nil || (p.current != nil && p.current.XMSContinuationNextTableName == nil) {
		return false
	}
	var resp generated.TableQueryResponseEnvelope
	resp, p.err = p.client.Query(ctx, p.tableQueryOptions, p.listOptions.toQueryOptions())
	p.current = &resp
	p.tableQueryOptions.NextTableName = resp.XMSContinuationNextTableName
	return p.err == nil && resp.TableQueryResponse.Value != nil && len(resp.TableQueryResponse.Value) > 0
}

// PageResponse returns the results from the page most recently fetched from the service.
// Example usage of this in combination with NextPage would look like the following:
//
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.Printf("The page contains %i results.\n", len(resp.Tables))
// }
func (p *tableQueryResponsePager) PageResponse() ListTablesPage {
	return *fromGeneratedTableQueryResponseEnvelope(p.current)
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (p *tableQueryResponsePager) Err() error {
	return p.err
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
// List returns a Pager, which allows iteration through each page of results. Example:
//
// options := &ListTablesOptions{Filter: to.StringPtr("PartitionKey eq 'pk001'"), Top: to.Int32Ptr(25)}
// pager := client.List(options) // Pass in 'nil' if you want to return all Tables for an account.
// for pager.NextPage(ctx) {
//     resp = pager.PageResponse()
//     fmt.Printf("The page contains %i results.\n", len(resp.Tables))
// }
// err := pager.Err()
func (t *ServiceClient) ListTables(listOptions *ListTablesOptions) ListTablesPager {
	return &tableQueryResponsePager{
		client:            t.client,
		listOptions:       listOptions,
		tableQueryOptions: new(generated.TableQueryOptions),
	}
}

// GetStatisticsOptions are the options for a ServiceClient.GetStatistics call
type GetStatisticsOptions struct {
}

type GetStatisticsResponse struct {
	RawResponse    *http.Response
	GeoReplication *GeoReplication `xml:"GeoReplication"`
}

func getStatisticsResponseFromGenerated(g *generated.ServiceGetStatisticsResponse) GetStatisticsResponse {
	return GetStatisticsResponse{
		RawResponse:    g.RawResponse,
		GeoReplication: fromGeneratedGeoReplication(g.GeoReplication),
	}
}

func (g *GetStatisticsOptions) toGenerated() *generated.ServiceGetStatisticsOptions {
	return &generated.ServiceGetStatisticsOptions{}
}

// GetStatistics retrieves all the statistics for an account with Geo-redundancy established.
//
// response, err := client.GetStatistics(context.Background, nil)
// handle(err)
// fmt.Println("Status: ", response.GeoReplication.Status)
// fmt.Println(Last Sync Time: ", response.GeoReplication.LastSyncTime)
func (t *ServiceClient) GetStatistics(ctx context.Context, options *GetStatisticsOptions) (GetStatisticsResponse, error) {
	if options == nil {
		options = &GetStatisticsOptions{}
	}
	resp, err := t.service.GetStatistics(ctx, options.toGenerated())
	return getStatisticsResponseFromGenerated(&resp), err
}

type GetPropertiesOptions struct {
}

func (g *GetPropertiesOptions) toGenerated() *generated.ServiceGetPropertiesOptions {
	return &generated.ServiceGetPropertiesOptions{}
}

type GetPropertiesResponse struct {
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

func getPropertiesResponseFromGenerated(g *generated.ServiceGetPropertiesResponse) GetPropertiesResponse {
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
//
// response, err := client.GetProperties(context.Background, nil)
// handle(err)
// fmt.Println(response.Cors)
// fmt.Println(response.HourMetrics)
// fmt.Println(response.Logging)
// fmt.Println(response.MinuteMetrics)
func (t *ServiceClient) GetProperties(ctx context.Context, options *GetPropertiesOptions) (GetPropertiesResponse, error) {
	if options == nil {
		options = &GetPropertiesOptions{}
	}
	resp, err := t.service.GetProperties(ctx, options.toGenerated())
	return getPropertiesResponseFromGenerated(&resp), err
}

type SetPropertiesOptions struct{}

func (s *SetPropertiesOptions) toGenerated() *generated.ServiceSetPropertiesOptions {
	return &generated.ServiceSetPropertiesOptions{}
}

type SetPropertiesResponse struct {
	RawResponse *http.Response
}

func setPropertiesResponseFromGenerated(g *generated.ServiceSetPropertiesResponse) SetPropertiesResponse {
	return SetPropertiesResponse{
		RawResponse: g.RawResponse,
	}
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
// logging := aztables.Logging{
// 		Read:    to.BoolPtr(true),
// 		Write:   to.BoolPtr(true),
// 		Delete:  to.BoolPtr(true),
// 		Version: to.StringPtr("1.0"),
// 		RetentionPolicy: &aztables.RetentionPolicy{
// 			Enabled: to.BoolPtr(true),
// 		Days:    to.Int32Ptr(5),
// 		},
// }
// props := aztables.ServiceProperties{Logging: &logging}
// resp, err := client.SetProperties(ctx, props, nil)
// handle(err)
func (t *ServiceClient) SetProperties(ctx context.Context, properties ServiceProperties, options *SetPropertiesOptions) (SetPropertiesResponse, error) {
	if options == nil {
		options = &SetPropertiesOptions{}
	}
	resp, err := t.service.SetProperties(ctx, *properties.toGenerated(), options.toGenerated())
	return setPropertiesResponseFromGenerated(&resp), err
}

// GetAccountSASToken is a convenience method for generating a SAS token for the currently pointed at account. This methods returns the full service URL and an error
// if there was an error during creation. This method can only be used if the supplied azcore.Credential during creation was a SharedKeyCredential.
func (t ServiceClient) GetAccountSASToken(resources AccountSASResourceTypes, permissions AccountSASPermissions, start time.Time, expiry time.Time) (string, error) {
	cred, ok := t.cred.(*SharedKeyCredential)
	if !ok {
		return "", errors.New("credential is not a SharedKeyCredential. SAS can only be signed with a SharedKeyCredential")
	}
	qps, err := AccountSASSignatureValues{
		Version:       SASVersion,
		Protocol:      SASProtocolHTTPS,
		Permissions:   permissions.String(),
		Services:      "t",
		ResourceTypes: resources.String(),
		StartTime:     start.UTC(),
		ExpiryTime:    expiry.UTC(),
	}.Sign(cred)
	if err != nil {
		return "", err
	}
	endpoint := t.client.Con.Endpoint()
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
