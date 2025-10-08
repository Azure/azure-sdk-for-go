// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/azcore/log"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

const (
	apiVersion = "2020-11-05"
)

// Client is used to interact with the Azure Cosmos DB database service.
type Client struct {
	endpoint    string
	internal    *azcore.Client
	gem         *globalEndpointManager
	endpointUrl *url.URL
}

// Endpoint used to create the client.
func (c *Client) Endpoint() string {
	return c.endpoint
}

// NewClientWithKey creates a new instance of Cosmos client with shared key authentication. It uses the default pipeline configuration.
// endpoint - The cosmos service endpoint to use.
// cred - The credential used to authenticate with the cosmos service.
// options - Optional Cosmos client options.  Pass nil to accept default values.
func NewClientWithKey(endpoint string, cred KeyCredential, o *ClientOptions) (*Client, error) {
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}
	preferredRegions := []string{}
	enableCrossRegionRetries := true
	if o != nil {
		preferredRegions = o.PreferredRegions
	}

	gem, err := newGlobalEndpointManager(endpoint, newInternalPipeline(newSharedKeyCredPolicy(cred), o), preferredRegions, 0, enableCrossRegionRetries)
	if err != nil {
		return nil, err
	}

	internalClient, err := newClient(newSharedKeyCredPolicy(cred), gem, o)
	if err != nil {
		return nil, err
	}
	return &Client{endpoint: endpoint, endpointUrl: endpointUrl, internal: internalClient, gem: gem}, nil
}

// NewClient creates a new instance of Cosmos client with Azure AD access token authentication. It uses the default pipeline configuration.
// endpoint - The cosmos service endpoint to use.
// cred - The credential used to authenticate with the cosmos service.
// options - Optional Cosmos client options.  Pass nil to accept default values.
func NewClient(endpoint string, cred azcore.TokenCredential, o *ClientOptions) (*Client, error) {
	endpointUrl, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	var scope []string

	if o != nil && o.ClientOptions.Cloud.Services != nil {
		if svcCfg, ok := o.ClientOptions.Cloud.Services[ServiceName]; ok && svcCfg.Audience != "" {
			audience := svcCfg.Audience
			if !strings.HasSuffix(audience, "/.default") {
				audience = strings.TrimSuffix(audience, "/") + "/.default"
			}
			scope = []string{audience}
			log.Write(azlog.EventRequest, fmt.Sprintf("Using audience from client options: %s", audience))
		}
	}

	if scope == nil {
		// Fallback to account-scope
		scope, err = createScopeFromEndpoint(endpointUrl)
		if err != nil {
			return nil, err
		}
		log.Write(azlog.EventRequest, fmt.Sprintf("Using account scope from endpoint: %s", scope[0]))
	}

	preferredRegions := []string{}
	enableCrossRegionRetries := true
	if o != nil {
		preferredRegions = o.PreferredRegions
	}
	gem, err := newGlobalEndpointManager(endpoint, newInternalPipeline(newCosmosBearerTokenPolicy(cred, scope, nil), o), preferredRegions, 0, enableCrossRegionRetries)
	if err != nil {
		return nil, err
	}

	internalClient, err := newClient(newCosmosBearerTokenPolicy(cred, scope, nil), gem, o)
	if err != nil {
		return nil, err
	}
	return &Client{endpoint: endpoint, endpointUrl: endpointUrl, internal: internalClient, gem: gem}, nil
}

// NewClientFromConnectionString creates a new instance of Cosmos client from connection string. It uses the default pipeline configuration.
// connectionString - The cosmos service connection string.
// options - Optional Cosmos client options.  Pass nil to accept default values.
func NewClientFromConnectionString(connectionString string, o *ClientOptions) (*Client, error) {
	const (
		accountEndpoint = "AccountEndpoint"
		accountKey      = "AccountKey"
	)

	splits := strings.SplitN(connectionString, ";", 2)
	if len(splits) < 2 {
		return nil, errors.New("failed parsing connection string due to it not consist of two parts separated by ';'")
	}

	var endpoint string
	var cred KeyCredential
	for _, split := range splits {
		keyVal := strings.SplitN(split, "=", 2)
		if len(keyVal) < 2 {
			return nil, fmt.Errorf("failed parsing connection string due to unmatched key value separated by '='")
		}
		switch {
		case strings.EqualFold(accountEndpoint, keyVal[0]):
			endpoint = keyVal[1]
		case strings.EqualFold(accountKey, keyVal[0]):
			c, err := NewKeyCredential(strings.TrimSuffix(keyVal[1], ";"))
			if err != nil {
				return nil, err
			}
			cred = c
		}
	}

	return NewClientWithKey(endpoint, cred, o)
}

func newClient(authPolicy policy.Policy, gem *globalEndpointManager, options *ClientOptions) (*azcore.Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	return azcore.NewClient(moduleName, serviceLibVersion,
		azruntime.PipelineOptions{
			AllowedHeaders: getAllowedHeaders(),
			PerCall: []policy.Policy{
				&headerPolicies{
					enableContentResponseOnWrite: options.EnableContentResponseOnWrite,
				},
				&globalEndpointManagerPolicy{gem: gem},
			},
			PerRetry: []policy.Policy{
				authPolicy,
				&clientRetryPolicy{gem: gem},
			},
			Tracing: azruntime.TracingOptions{
				Namespace: "Microsoft.DocumentDB",
			},
		},
		&options.ClientOptions)
}

func newInternalPipeline(authPolicy policy.Policy, options *ClientOptions) azruntime.Pipeline {
	if options == nil {
		options = &ClientOptions{}
	}
	return azruntime.NewPipeline(moduleName, serviceLibVersion,
		azruntime.PipelineOptions{
			AllowedHeaders: getAllowedHeaders(),
			PerRetry: []policy.Policy{
				authPolicy,
			},
		},
		&options.ClientOptions)
}

func createScopeFromEndpoint(endpoint *url.URL) ([]string, error) {
	return []string{fmt.Sprintf("%s://%s/.default", endpoint.Scheme, endpoint.Hostname())}, nil
}

// NewDatabase returns a struct that represents a database and allows database level operations.
// id - The id of the database.
func (c *Client) NewDatabase(id string) (*DatabaseClient, error) {
	if id == "" {
		return nil, errors.New("id is required")
	}

	return newDatabase(id, c)
}

// NewContainer returns a struct that represents a container and allows container level operations.
// databaseId - The id of the database.
// containerId - The id of the container.
func (c *Client) NewContainer(databaseId string, containerId string) (*ContainerClient, error) {
	if databaseId == "" {
		return nil, errors.New("databaseId is required")
	}

	if containerId == "" {
		return nil, errors.New("containerId is required")
	}

	db, err := newDatabase(databaseId, c)
	if err != nil {
		return nil, err
	}

	return db.NewContainer(containerId)
}

// CreateDatabase creates a new database.
// ctx - The context for the request.
// databaseProperties - The definition of the database
// o - Options for the create database operation.
func (c *Client) CreateDatabase(
	ctx context.Context,
	databaseProperties DatabaseProperties,
	o *CreateDatabaseOptions) (DatabaseResponse, error) {
	var err error
	spanName, err := getSpanNameForDatabases(c.accountEndpointUrl(), operationTypeCreate, resourceTypeDatabase, databaseProperties.ID)
	if err != nil {
		return DatabaseResponse{}, err
	}
	ctx, endSpan := azruntime.StartSpan(ctx, spanName.name, c.internal.Tracer(), &spanName.options)
	defer func() { endSpan(err) }()

	if o == nil {
		o = &CreateDatabaseOptions{}
	}
	returnResponse := true
	h := &headerOptionsOverride{
		enableContentResponseOnWrite: &returnResponse,
	}

	operationContext := pipelineRequestOptions{
		resourceType:          resourceTypeDatabase,
		resourceAddress:       "",
		isWriteOperation:      true,
		headerOptionsOverride: h,
	}

	path, err := generatePathForNameBased(resourceTypeDatabase, "", true)
	if err != nil {
		return DatabaseResponse{}, err
	}

	azResponse, err := c.sendPostRequest(
		path,
		ctx,
		databaseProperties,
		operationContext,
		nil,
		o.ThroughputProperties.addHeadersToRequest)
	if err != nil {
		return DatabaseResponse{}, err
	}

	response, err := newDatabaseResponse(azResponse)
	return response, err
}

// NewQueryDatabasesPager executes query for databases.
// query - The SQL query to execute.
// o - Options for the operation.
func (c *Client) NewQueryDatabasesPager(query string, o *QueryDatabasesOptions) *azruntime.Pager[QueryDatabasesResponse] {
	queryOptions := &QueryDatabasesOptions{}
	if o != nil {
		originalOptions := *o
		queryOptions = &originalOptions
	}

	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: "",
	}

	path, _ := generatePathForNameBased(resourceTypeDatabase, operationContext.resourceAddress, true)

	return azruntime.NewPager(azruntime.PagingHandler[QueryDatabasesResponse]{
		More: func(page QueryDatabasesResponse) bool {
			return page.ContinuationToken != nil
		},
		Fetcher: func(ctx context.Context, page *QueryDatabasesResponse) (QueryDatabasesResponse, error) {
			var err error
			spanName, err := getSpanNameForClient(c.accountEndpointUrl(), operationTypeQuery, resourceTypeDatabase, c.accountEndpointUrl().Hostname())
			if err != nil {
				return QueryDatabasesResponse{}, err
			}
			ctx, endSpan := azruntime.StartSpan(ctx, spanName.name, c.internal.Tracer(), &spanName.options)
			defer func() { endSpan(err) }()
			if page != nil {
				if page.ContinuationToken != nil {
					// Use the previous page continuation if available
					queryOptions.ContinuationToken = page.ContinuationToken
				}
			}

			azResponse, err := c.sendQueryRequest(
				path,
				ctx,
				query,
				queryOptions.QueryParameters,
				operationContext,
				queryOptions,
				nil)

			if err != nil {
				return QueryDatabasesResponse{}, err
			}

			return newDatabasesQueryResponse(azResponse)
		},
	})
}

func (c *Client) sendPostRequest(
	path string,
	ctx context.Context,
	content interface{},
	operationContext pipelineRequestOptions,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodPost, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	err = c.attachContent(content, req)
	if err != nil {
		return nil, err
	}

	return c.executeAndEnsureSuccessResponse(ctx, req)
}

func (c *Client) sendQueryRequest(
	path string,
	ctx context.Context,
	query string,
	parameters []QueryParameter,
	operationContext pipelineRequestOptions,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodPost, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	err = azruntime.MarshalAsJSON(req, queryBody{
		Query:      query,
		Parameters: parameters,
	})

	if err != nil {
		return nil, err
	}

	req.Raw().Header.Add(cosmosHeaderQuery, "True")
	// Override content type for query
	req.Raw().Header.Set(headerContentType, cosmosHeaderValuesQuery)

	return c.executeAndEnsureSuccessResponse(ctx, req)
}

func (c *Client) sendPutRequest(
	path string,
	ctx context.Context,
	content interface{},
	operationContext pipelineRequestOptions,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodPut, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	err = c.attachContent(content, req)
	if err != nil {
		return nil, err
	}

	return c.executeAndEnsureSuccessResponse(ctx, req)
}

func (c *Client) sendGetRequest(
	path string,
	ctx context.Context,
	operationContext pipelineRequestOptions,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodGet, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	return c.executeAndEnsureSuccessResponse(ctx, req)
}

func (c *Client) sendDeleteRequest(
	path string,
	ctx context.Context,
	operationContext pipelineRequestOptions,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodDelete, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	return c.executeAndEnsureSuccessResponse(ctx, req)
}

func (c *Client) sendBatchRequest(
	ctx context.Context,
	path string,
	batch []batchOperation,
	operationContext pipelineRequestOptions,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodPost, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	err = c.attachContent(batch, req)
	if err != nil {
		return nil, err
	}

	return c.executeAndEnsureSuccessResponse(ctx, req)
}

func (c *Client) sendPatchRequest(
	path string,
	ctx context.Context,
	content interface{},
	operationContext pipelineRequestOptions,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodPatch, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	err = c.attachContent(content, req)
	if err != nil {
		return nil, err
	}

	return c.executeAndEnsureSuccessResponse(ctx, req)
}

func (c *Client) createRequest(
	path string,
	ctx context.Context,
	method string,
	operationContext pipelineRequestOptions,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*policy.Request, error) {

	// todo: endpoint will be set originally by globalendpointmanager
	finalURL := c.endpoint

	if path != "" {
		finalURL = azruntime.JoinPaths(c.endpoint, path)
	}

	req, err := azruntime.NewRequest(ctx, method, finalURL)
	if err != nil {
		return nil, err
	}

	if requestOptions != nil {
		headers := requestOptions.toHeaders()
		if headers != nil {
			for k, v := range *headers {
				req.Raw().Header.Set(k, v)
			}
		}
	}

	addDefaultHeaders(req)

	req.SetOperationValue(operationContext)

	if requestEnricher != nil {
		requestEnricher(req)
	}

	return req, nil
}

func (c *Client) attachContent(content interface{}, req *policy.Request) error {
	var err error
	switch v := content.(type) {
	case []byte:
		// If its a raw byte array, we can just set the body
		err = req.SetBody(streaming.NopCloser(bytes.NewReader(v)), "application/json")
	default:
		// Otherwise, we need to marshal it
		err = azruntime.MarshalAsJSON(req, content)
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) executeAndEnsureSuccessResponse(ctx context.Context, request *policy.Request) (*http.Response, error) {
	log.Write(azlog.EventResponse, fmt.Sprintf("\n===== Client preferred regions:\n%v\n=====\n", c.gem.preferredLocations))
	response, err := c.internal.Pipeline().Do(request)
	if err != nil {
		return nil, err
	}

	c.addResponseValuesToSpan(ctx, response)

	successResponse := (response.StatusCode >= 200 && response.StatusCode < 300) || response.StatusCode == 304
	if successResponse {
		return response, nil
	}

	return nil, azruntime.NewResponseErrorWithErrorCode(response, response.Status)
}

func (c *Client) accountEndpointUrl() *url.URL {
	return c.endpointUrl
}

func (c *Client) addResponseValuesToSpan(ctx context.Context, resp *http.Response) {
	span := c.internal.Tracer().SpanFromContext(ctx)
	span.SetAttributes(
		tracing.Attribute{Key: "db.cosmosdb.request_charge", Value: newResponse(resp).RequestCharge},
		tracing.Attribute{Key: "db.cosmosdb.status_code", Value: resp.StatusCode},
	)
}

type pipelineRequestOptions struct {
	headerOptionsOverride *headerOptionsOverride
	resourceType          resourceType
	resourceAddress       string
	isRidBased            bool
	isWriteOperation      bool
}

func addDefaultHeaders(req *policy.Request) {
	req.Raw().Header.Set(headerXmsDate, time.Now().UTC().Format(http.TimeFormat))
	req.Raw().Header.Set(headerXmsVersion, apiVersion)
	req.Raw().Header.Set(cosmosHeaderSDKSupportedCapabilities, supportedCapabilitiesHeaderValue)
}

func getAllowedHeaders() []string {
	return []string{
		cosmosHeaderRequestCharge,
		cosmosHeaderActivityId,
		cosmosHeaderEtag,
		cosmosHeaderSubstatus,
		cosmosHeaderPopulateQuotaInfo,
		cosmosHeaderPreTriggerInclude,
		cosmosHeaderPostTriggerInclude,
		cosmosHeaderIndexingDirective,
		cosmosHeaderSessionToken,
		cosmosHeaderConsistencyLevel,
		cosmosHeaderPrefer,
		cosmosHeaderIsUpsert,
		cosmosHeaderOfferThroughput,
		cosmosHeaderOfferAutoscale,
		cosmosHeaderQuery,
		cosmosHeaderOfferReplacePending,
		cosmosHeaderOfferMinimumThroughput,
		cosmosHeaderResponseContinuationTokenLimitInKb,
		cosmosHeaderEnableScanInQuery,
		cosmosHeaderMaxItemCount,
		cosmosHeaderContinuationToken,
		cosmosHeaderPopulateIndexMetrics,
		cosmosHeaderPopulateQueryMetrics,
		cosmosHeaderQueryMetrics,
		cosmosHeaderIndexUtilization,
		cosmosHeaderCorrelatedActivityId,
		cosmosHeaderIsBatchRequest,
		cosmosHeaderIsBatchAtomic,
		cosmosHeaderIsBatchOrdered,
		cosmosHeaderSDKSupportedCapabilities,
		headerXmsDate,
		headerContentType,
		headerIfMatch,
		headerIfNoneMatch,
		headerXmsVersion,
		headerContentLocation,
		headerXmsGatewayVersion,
		headerLsn,
		headerXmsCosmosLlsn,
		headerXmsCosmosItemLlsn,
		headerXmsItemLsn,
		headerXmsCosmosQuorumAckedLlsn,
		headerXmsCurrentReplicaSetSize,
		headerXmsCurrentWriteQuorum,
		headerXmsGlobalCommittedLsn,
		headerXmsLastStateChangeUtc,
		headerXmsNumberOfReadRegions,
		headerXmsQuorumAckedLsn,
		headerXmsRequestDurationMs,
		headerXmsResourceQuota,
		headerXmsResourceUsage,
		headerXmsSchemaVersion,
		headerXmsServiceVersion,
		headerXmsTransportRequestId,
		headerXmsXpRole,
		headerCollectionPartitionIndex,
		headerCollectionServiceIndex,
		headerXmsDocumentDbPartitionKeyRangeId,
		cosmosHeaderPhysicalPartitionId,
		headerStrictTransportSecurity,
		headerXmsDatabaseAccountConsumedMb,
		headerXmsDatabaseAccountProvisionedMb,
		headerXmsDatabaseAccountReservedMb,
		headerXmsMaxMediaStorageUsageMb,
		headerXmsMediaStorageUsageMb,
		headerXmsContentPath,
		headerXmsAltContentPath,
		cosmosHeaderMaxContentLength,
		cosmosHeaderIsPartitionKeyDeletePending,
		cosmosHeaderQueryExecutionInfo,
		headerXmsItemCount,
	}
}
