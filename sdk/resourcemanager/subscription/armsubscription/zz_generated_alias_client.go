//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsubscription

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

// AliasClient contains the methods for the Alias group.
// Don't use this type directly, use NewAliasClient() instead.
type AliasClient struct {
	ep string
	pl runtime.Pipeline
}

// NewAliasClient creates a new instance of AliasClient with the specified values.
func NewAliasClient(credential azcore.TokenCredential, options *arm.ClientOptions) *AliasClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &AliasClient{ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// BeginCreate - Create Alias Subscription.
// If the operation fails it returns the *ErrorResponseBody error type.
func (client *AliasClient) BeginCreate(ctx context.Context, aliasName string, body PutAliasRequest, options *AliasBeginCreateOptions) (AliasCreatePollerResponse, error) {
	resp, err := client.create(ctx, aliasName, body, options)
	if err != nil {
		return AliasCreatePollerResponse{}, err
	}
	result := AliasCreatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("AliasClient.Create", "", resp, client.pl, client.createHandleError)
	if err != nil {
		return AliasCreatePollerResponse{}, err
	}
	result.Poller = &AliasCreatePoller{
		pt: pt,
	}
	return result, nil
}

// Create - Create Alias Subscription.
// If the operation fails it returns the *ErrorResponseBody error type.
func (client *AliasClient) create(ctx context.Context, aliasName string, body PutAliasRequest, options *AliasBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, aliasName, body, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, client.createHandleError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *AliasClient) createCreateRequest(ctx context.Context, aliasName string, body PutAliasRequest, options *AliasBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Subscription/aliases/{aliasName}"
	if aliasName == "" {
		return nil, errors.New("parameter aliasName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{aliasName}", url.PathEscape(aliasName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, body)
}

// createHandleError handles the Create error response.
func (client *AliasClient) createHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponseBody{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Delete - Delete Alias.
// If the operation fails it returns the *ErrorResponseBody error type.
func (client *AliasClient) Delete(ctx context.Context, aliasName string, options *AliasDeleteOptions) (AliasDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, aliasName, options)
	if err != nil {
		return AliasDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AliasDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return AliasDeleteResponse{}, client.deleteHandleError(resp)
	}
	return AliasDeleteResponse{RawResponse: resp}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *AliasClient) deleteCreateRequest(ctx context.Context, aliasName string, options *AliasDeleteOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Subscription/aliases/{aliasName}"
	if aliasName == "" {
		return nil, errors.New("parameter aliasName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{aliasName}", url.PathEscape(aliasName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *AliasClient) deleteHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponseBody{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Get Alias Subscription.
// If the operation fails it returns the *ErrorResponseBody error type.
func (client *AliasClient) Get(ctx context.Context, aliasName string, options *AliasGetOptions) (AliasGetResponse, error) {
	req, err := client.getCreateRequest(ctx, aliasName, options)
	if err != nil {
		return AliasGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AliasGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AliasGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *AliasClient) getCreateRequest(ctx context.Context, aliasName string, options *AliasGetOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Subscription/aliases/{aliasName}"
	if aliasName == "" {
		return nil, errors.New("parameter aliasName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{aliasName}", url.PathEscape(aliasName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AliasClient) getHandleResponse(resp *http.Response) (AliasGetResponse, error) {
	result := AliasGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SubscriptionAliasResponse); err != nil {
		return AliasGetResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *AliasClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponseBody{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// List - List Alias Subscription.
// If the operation fails it returns the *ErrorResponseBody error type.
func (client *AliasClient) List(ctx context.Context, options *AliasListOptions) (AliasListResponse, error) {
	req, err := client.listCreateRequest(ctx, options)
	if err != nil {
		return AliasListResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AliasListResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AliasListResponse{}, client.listHandleError(resp)
	}
	return client.listHandleResponse(resp)
}

// listCreateRequest creates the List request.
func (client *AliasClient) listCreateRequest(ctx context.Context, options *AliasListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Subscription/aliases"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *AliasClient) listHandleResponse(resp *http.Response) (AliasListResponse, error) {
	result := AliasListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SubscriptionAliasListResult); err != nil {
		return AliasListResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *AliasClient) listHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponseBody{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
