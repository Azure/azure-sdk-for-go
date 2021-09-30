// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	azruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// cosmosClientConnection maintains a Pipeline for the client.
// The Pipeline is build based on the CosmosClientOptions.
type cosmosClientConnection struct {
	endpoint string
	Pipeline azruntime.Pipeline
}

// newConnection creates an instance of the connection type with the specified endpoint.
// Pass nil to accept the default options; this is the same as passing a zero-value options.
func newCosmosClientConnection(endpoint string, cred azcore.Credential, options *CosmosClientOptions) *cosmosClientConnection {
	policies := []policy.Policy{
		azruntime.NewTelemetryPolicy("azcosmos", serviceLibVersion, &options.Telemetry),
	}
	policies = append(policies, options.PerCallPolicies...)
	policies = append(policies, azruntime.NewRetryPolicy(&options.Retry))
	policies = append(policies, options.PerRetryPolicies...)
	policies = append(policies, options.getSDKInternalPolicies()...)
	policies = append(policies, cred.NewAuthenticationPolicy(azruntime.AuthenticationOptions{}))
	policies = append(policies, azruntime.NewLogPolicy(&options.Logging))
	return &cosmosClientConnection{endpoint: endpoint, Pipeline: azruntime.NewPipeline(options.HTTPClient, policies...)}
}

func (c *cosmosClientConnection) sendPostRequest(
	path string,
	ctx context.Context,
	content interface{},
	operationContext cosmosOperationContext,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodPost, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	err = azruntime.MarshalAsJSON(req, content)
	if err != nil {
		return nil, err
	}

	return c.executeAndEnsureSuccessResponse(req)
}

func (c *cosmosClientConnection) sendQueryRequest(
	path string,
	ctx context.Context,
	query string,
	operationContext cosmosOperationContext,
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

func (c *cosmosClientConnection) sendPutRequest(
	path string,
	ctx context.Context,
	content interface{},
	operationContext cosmosOperationContext,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodPut, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	err = azruntime.MarshalAsJSON(req, content)
	if err != nil {
		return nil, err
	}

	return c.executeAndEnsureSuccessResponse(req)
}

func (c *cosmosClientConnection) sendGetRequest(
	path string,
	ctx context.Context,
	operationContext cosmosOperationContext,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodGet, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	return c.executeAndEnsureSuccessResponse(req)
}

func (c *cosmosClientConnection) sendDeleteRequest(
	path string,
	ctx context.Context,
	operationContext cosmosOperationContext,
	requestOptions cosmosRequestOptions,
	requestEnricher func(*policy.Request)) (*http.Response, error) {
	req, err := c.createRequest(path, ctx, http.MethodDelete, operationContext, requestOptions, requestEnricher)
	if err != nil {
		return nil, err
	}

	return c.executeAndEnsureSuccessResponse(req)
}

func (c *cosmosClientConnection) createRequest(
	path string,
	ctx context.Context,
	method string,
	operationContext cosmosOperationContext,
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

	headers := requestOptions.toHeaders()
	if headers != nil {
		for k, v := range *headers {
			req.Raw().Header.Set(k, v)
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

func (c *cosmosClientConnection) executeAndEnsureSuccessResponse(request *policy.Request) (*http.Response, error) {
	response, err := c.Pipeline.Do(request)
	if err != nil {
		return nil, err
	}

	successResponse := (response.StatusCode >= 200 && response.StatusCode < 300) || response.StatusCode == 304
	if successResponse {
		return response, nil
	}

	return nil, newCosmosError(response)
}
