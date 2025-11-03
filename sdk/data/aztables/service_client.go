// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aztables

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
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
	cl := cloud.AzurePublic
	if options != nil && !reflect.ValueOf(options.Cloud).IsZero() {
		cl = options.Cloud
	}

	cfg, ok := cl.Services[ServiceName]
	if !ok || cfg.Audience == "" {
		return nil, errors.New("cloud configuration is missing for Azure Tables")
	}

	audience := cfg.Audience
	if isCosmosEndpoint(serviceURL) {
		audience = cfg.CosmosAudience
	}

	plOpts := runtime.PipelineOptions{
		PerRetry: []policy.Policy{runtime.NewBearerTokenPolicy(cred, []string{audience + "/.default"}, nil)},
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

// DeleteTable deletes a table by name. If the service returns a non-successful HTTP status code, the function returns an *azcore.ResponseError type.
// Specify nil for options if you want to use the default options.
func (t *ServiceClient) DeleteTable(ctx context.Context, name string, options *DeleteTableOptions) (DeleteTableResponse, error) {
	var err error
	ctx, endSpan := runtime.StartSpan(ctx, "ServiceClient.DeleteTable", t.client.Tracer(), nil)
	defer func() { endSpan(err) }()

	_, err = t.client.Delete(ctx, name, options.toGenerated())
	return DeleteTableResponse{}, err
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
				tableName = page.NextTableName
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

			tableProps := make([]*TableProperties, len(resp.Value))
			for i := range resp.Value {
				odataValues := map[string]any{}
				if resp.Value[i].ODataEditLink != nil {
					odataValues["odata.editLink"] = *resp.Value[i].ODataEditLink
				}
				if resp.Value[i].ODataID != nil {
					odataValues["odata.id"] = *resp.Value[i].ODataID
				}
				if resp.Value[i].ODataType != nil {
					odataValues["odata.type"] = *resp.Value[i].ODataType
				}
				var odataJSON []byte
				if len(odataValues) > 0 {
					odataJSON, err = json.Marshal(odataValues)
					if err != nil {
						return ListTablesResponse{}, err
					}
				}
				tableProps[i] = &TableProperties{
					Name:  resp.Value[i].TableName,
					Value: odataJSON,
				}
			}

			return ListTablesResponse{
				NextTableName: resp.XMSContinuationNextTableName,
				Tables:        tableProps,
			}, nil
		},
		Tracer: t.client.Tracer(),
	})
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
	return GetStatisticsResponse{
		GeoReplication: fromGeneratedGeoReplication(resp.GeoReplication),
	}, nil
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
	var cors []*CorsRule
	if len(resp.Cors) > 0 {
		cors = make([]*CorsRule, len(resp.Cors))
		for i := range resp.Cors {
			cors[i] = fromGeneratedCors(resp.Cors[i])
		}
	}
	return GetPropertiesResponse{
		ServiceProperties: ServiceProperties{
			Cors:          cors,
			HourMetrics:   fromGeneratedMetrics(resp.HourMetrics),
			Logging:       fromGeneratedLogging(resp.Logging),
			MinuteMetrics: fromGeneratedMetrics(resp.MinuteMetrics),
		},
	}, nil
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
	_, err = t.service.SetProperties(ctx, *properties.toGenerated(), options.toGenerated())
	return SetPropertiesResponse{}, err
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
