//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armreservations

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// QuotaClient contains the methods for the Quota group.
// Don't use this type directly, use NewQuotaClient() instead.
type QuotaClient struct {
	ep string
	pl runtime.Pipeline
}

// NewQuotaClient creates a new instance of QuotaClient with the specified values.
func NewQuotaClient(credential azcore.TokenCredential, options *arm.ClientOptions) *QuotaClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &QuotaClient{ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// BeginCreateOrUpdate - Create or update the quota (service limits) of a resource to the requested value. Steps:
// 1. Make the Get request to get the quota information for specific resource.
//
//
// 2. To increase the quota, update the limit field in the response from Get request to new value.
//
//
// 3. Submit the JSON to the quota request API to update the quota. The Create quota request may be constructed as follows. The PUT operation can be used
// to update the quota.
// If the operation fails it returns the *ExceptionResponse error type.
func (client *QuotaClient) BeginCreateOrUpdate(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaBeginCreateOrUpdateOptions) (QuotaCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, subscriptionID, providerID, location, resourceName, createQuotaRequest, options)
	if err != nil {
		return QuotaCreateOrUpdatePollerResponse{}, err
	}
	result := QuotaCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("QuotaClient.CreateOrUpdate", "location", resp, client.pl, client.createOrUpdateHandleError)
	if err != nil {
		return QuotaCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &QuotaCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Create or update the quota (service limits) of a resource to the requested value. Steps:
// 1. Make the Get request to get the quota information for specific resource.
//
//
// 2. To increase the quota, update the limit field in the response from Get request to new value.
//
//
// 3. Submit the JSON to the quota request API to update the quota. The Create quota request may be constructed as follows. The PUT operation can be used
// to update the quota.
// If the operation fails it returns the *ExceptionResponse error type.
func (client *QuotaClient) createOrUpdate(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, subscriptionID, providerID, location, resourceName, createQuotaRequest, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *QuotaClient) createOrUpdateCreateRequest(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Capacity/resourceProviders/{providerId}/locations/{location}/serviceLimits/{resourceName}"
	if subscriptionID == "" {
		return nil, errors.New("parameter subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	if providerID == "" {
		return nil, errors.New("parameter providerID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{providerId}", url.PathEscape(providerID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-25")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, createQuotaRequest)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *QuotaClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ExceptionResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Get the current quota (service limit) and usage of a resource. You can use the response from the GET operation to submit quota update request.
// If the operation fails it returns the *ExceptionResponse error type.
func (client *QuotaClient) Get(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, options *QuotaGetOptions) (QuotaGetResponse, error) {
	req, err := client.getCreateRequest(ctx, subscriptionID, providerID, location, resourceName, options)
	if err != nil {
		return QuotaGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return QuotaGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return QuotaGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *QuotaClient) getCreateRequest(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, options *QuotaGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Capacity/resourceProviders/{providerId}/locations/{location}/serviceLimits/{resourceName}"
	if subscriptionID == "" {
		return nil, errors.New("parameter subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	if providerID == "" {
		return nil, errors.New("parameter providerID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{providerId}", url.PathEscape(providerID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-25")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *QuotaClient) getHandleResponse(resp *http.Response) (QuotaGetResponse, error) {
	result := QuotaGetResponse{RawResponse: resp}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.CurrentQuotaLimitBase); err != nil {
		return QuotaGetResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *QuotaClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ExceptionResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// List - Gets a list of current quotas (service limits) and usage for all resources. The response from the list quota operation can be leveraged to request
// quota updates.
// If the operation fails it returns the *ExceptionResponse error type.
func (client *QuotaClient) List(subscriptionID string, providerID string, location string, options *QuotaListOptions) *QuotaListPager {
	return &QuotaListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, subscriptionID, providerID, location, options)
		},
		advancer: func(ctx context.Context, resp QuotaListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.QuotaLimits.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *QuotaClient) listCreateRequest(ctx context.Context, subscriptionID string, providerID string, location string, options *QuotaListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Capacity/resourceProviders/{providerId}/locations/{location}/serviceLimits"
	if subscriptionID == "" {
		return nil, errors.New("parameter subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	if providerID == "" {
		return nil, errors.New("parameter providerID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{providerId}", url.PathEscape(providerID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-25")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *QuotaClient) listHandleResponse(resp *http.Response) (QuotaListResponse, error) {
	result := QuotaListResponse{RawResponse: resp}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.QuotaLimits); err != nil {
		return QuotaListResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *QuotaClient) listHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ExceptionResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginUpdate - Update the quota (service limits) of this resource to the requested value.
// • To get the quota information for specific resource, send a GET request.
// • To increase the quota, update the limit field from the GET response to a new value.
// • To update the quota value, submit the JSON response to the quota request API to update the quota. • To update the quota. use the PATCH operation.
// If the operation fails it returns the *ExceptionResponse error type.
func (client *QuotaClient) BeginUpdate(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaBeginUpdateOptions) (QuotaUpdatePollerResponse, error) {
	resp, err := client.update(ctx, subscriptionID, providerID, location, resourceName, createQuotaRequest, options)
	if err != nil {
		return QuotaUpdatePollerResponse{}, err
	}
	result := QuotaUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("QuotaClient.Update", "location", resp, client.pl, client.updateHandleError)
	if err != nil {
		return QuotaUpdatePollerResponse{}, err
	}
	result.Poller = &QuotaUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// Update - Update the quota (service limits) of this resource to the requested value.
// • To get the quota information for specific resource, send a GET request.
// • To increase the quota, update the limit field from the GET response to a new value.
// • To update the quota value, submit the JSON response to the quota request API to update the quota. • To update the quota. use the PATCH operation.
// If the operation fails it returns the *ExceptionResponse error type.
func (client *QuotaClient) update(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, subscriptionID, providerID, location, resourceName, createQuotaRequest, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, client.updateHandleError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *QuotaClient) updateCreateRequest(ctx context.Context, subscriptionID string, providerID string, location string, resourceName string, createQuotaRequest CurrentQuotaLimitBase, options *QuotaBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Capacity/resourceProviders/{providerId}/locations/{location}/serviceLimits/{resourceName}"
	if subscriptionID == "" {
		return nil, errors.New("parameter subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(subscriptionID))
	if providerID == "" {
		return nil, errors.New("parameter providerID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{providerId}", url.PathEscape(providerID))
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-25")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, createQuotaRequest)
}

// updateHandleError handles the Update error response.
func (client *QuotaClient) updateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ExceptionResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
