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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
)

// Client is used to interact with the Azure Cosmos DB database service.
type Client struct {
	endpoint string
	pipeline azruntime.Pipeline
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
	return &Client{endpoint: endpoint, pipeline: newPipeline(newSharedKeyCredPolicy(cred), o)}, nil
}

// NewClient creates a new instance of Cosmos client with Azure AD access token authentication. It uses the default pipeline configuration.
// endpoint - The cosmos service endpoint to use.
// cred - The credential used to authenticate with the cosmos service.
// options - Optional Cosmos client options.  Pass nil to accept default values.
func NewClient(endpoint string, cred azcore.TokenCredential, o *ClientOptions) (*Client, error) {
	scope, err := createScopeFromEndpoint(endpoint)
	if err != nil {
		return nil, err
	}
	return &Client{endpoint: endpoint, pipeline: newPipeline(newCosmosBearerTokenPolicy(cred, scope, nil), o)}, nil
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

func newPipeline(authPolicy policy.Policy, options *ClientOptions) azruntime.Pipeline {
	if options == nil {
		options = &ClientOptions{}
	}

	return azruntime.NewPipeline("azcosmos", serviceLibVersion,
		azruntime.PipelineOptions{
			PerCall: []policy.Policy{
				&headerPolicies{
					enableContentResponseOnWrite: options.EnableContentResponseOnWrite,
				},
			},
			PerRetry: []policy.Policy{
				authPolicy,
			},
		},
		&options.ClientOptions)
}

func createScopeFromEndpoint(endpoint string) ([]string, error) {
	u, err := url.Parse(endpoint)
	if err != nil {
		return nil, err
	}

	return []string{fmt.Sprintf("%s://%s/.default", u.Scheme, u.Hostname())}, nil
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
	if o == nil {
		o = &CreateDatabaseOptions{}
	}

	operationContext := pipelineRequestOptions{
		resourceType:    resourceTypeDatabase,
		resourceAddress: ""}

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

	return newDatabaseResponse(azResponse)
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

	return c.executeAndEnsureSuccessResponse(req)
}

func (c *Client) sendQueryRequest(
	path string,
	ctx context.Context,
	query string,
	operationContext pipelineRequestOptions,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodPost, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	type queryBody struct {
		Query string `json:"query"`
	}

	err = azruntime.MarshalAsJSON(req, queryBody{
		Query: query,
	})
	if err != nil {
		return nil, err
	}

	req.Raw().Header.Add(cosmosHeaderQuery, "True")
	// Override content type for query
	req.Raw().Header.Set(headerContentType, cosmosHeaderValuesQuery)

	return c.executeAndEnsureSuccessResponse(req)
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

	return c.executeAndEnsureSuccessResponse(req)
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

	return c.executeAndEnsureSuccessResponse(req)
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

	return c.executeAndEnsureSuccessResponse(req)
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

	return c.executeAndEnsureSuccessResponse(req)
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

	req.Raw().Header.Set(headerXmsDate, time.Now().UTC().Format(http.TimeFormat))
	req.Raw().Header.Set(headerXmsVersion, "2020-11-05")

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

func (c *Client) executeAndEnsureSuccessResponse(request *policy.Request) (*http.Response, error) {
	response, err := c.pipeline.Do(request)
	if err != nil {
		return nil, err
	}

	successResponse := (response.StatusCode >= 200 && response.StatusCode < 300) || response.StatusCode == 304
	if successResponse {
		return response, nil
	}

	return nil, newCosmosError(response)
}

type pipelineRequestOptions struct {
	headerOptionsOverride *headerOptionsOverride
	resourceType          resourceType
	resourceAddress       string
	isRidBased            bool
	isWriteOperation      bool
}
