//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapimanagement

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// PortalRevisionClient contains the methods for the PortalRevision group.
// Don't use this type directly, use NewPortalRevisionClient() instead.
type PortalRevisionClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewPortalRevisionClient creates a new instance of PortalRevisionClient with the specified values.
func NewPortalRevisionClient(con *arm.Connection, subscriptionID string) *PortalRevisionClient {
	return &PortalRevisionClient{ep: con.Endpoint(), pl: con.NewPipeline(module, version), subscriptionID: subscriptionID}
}

// BeginCreateOrUpdate - Creates a new developer portal's revision by running the portal's publishing. The isCurrent property indicates if the revision
// is publicly accessible.
// If the operation fails it returns the *ErrorResponse error type.
func (client *PortalRevisionClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, portalRevisionID string, parameters PortalRevisionContract, options *PortalRevisionBeginCreateOrUpdateOptions) (PortalRevisionCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceGroupName, serviceName, portalRevisionID, parameters, options)
	if err != nil {
		return PortalRevisionCreateOrUpdatePollerResponse{}, err
	}
	result := PortalRevisionCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("PortalRevisionClient.CreateOrUpdate", "location", resp, client.pl, client.createOrUpdateHandleError)
	if err != nil {
		return PortalRevisionCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &PortalRevisionCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Creates a new developer portal's revision by running the portal's publishing. The isCurrent property indicates if the revision is publicly
// accessible.
// If the operation fails it returns the *ErrorResponse error type.
func (client *PortalRevisionClient) createOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, portalRevisionID string, parameters PortalRevisionContract, options *PortalRevisionBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serviceName, portalRevisionID, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusCreated, http.StatusAccepted) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *PortalRevisionClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, portalRevisionID string, parameters PortalRevisionContract, options *PortalRevisionBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/portalRevisions/{portalRevisionId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if portalRevisionID == "" {
		return nil, errors.New("parameter portalRevisionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{portalRevisionId}", url.PathEscape(portalRevisionID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *PortalRevisionClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType.InnerError); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Gets the developer portal's revision specified by its identifier.
// If the operation fails it returns the *ErrorResponse error type.
func (client *PortalRevisionClient) Get(ctx context.Context, resourceGroupName string, serviceName string, portalRevisionID string, options *PortalRevisionGetOptions) (PortalRevisionGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, serviceName, portalRevisionID, options)
	if err != nil {
		return PortalRevisionGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PortalRevisionGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return PortalRevisionGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *PortalRevisionClient) getCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, portalRevisionID string, options *PortalRevisionGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/portalRevisions/{portalRevisionId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if portalRevisionID == "" {
		return nil, errors.New("parameter portalRevisionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{portalRevisionId}", url.PathEscape(portalRevisionID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PortalRevisionClient) getHandleResponse(resp *http.Response) (PortalRevisionGetResponse, error) {
	result := PortalRevisionGetResponse{RawResponse: resp}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.PortalRevisionContract); err != nil {
		return PortalRevisionGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *PortalRevisionClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType.InnerError); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// GetEntityTag - Gets the developer portal revision specified by its identifier.
// If the operation fails it returns the *ErrorResponse error type.
func (client *PortalRevisionClient) GetEntityTag(ctx context.Context, resourceGroupName string, serviceName string, portalRevisionID string, options *PortalRevisionGetEntityTagOptions) (PortalRevisionGetEntityTagResponse, error) {
	req, err := client.getEntityTagCreateRequest(ctx, resourceGroupName, serviceName, portalRevisionID, options)
	if err != nil {
		return PortalRevisionGetEntityTagResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return PortalRevisionGetEntityTagResponse{}, err
	}
	return client.getEntityTagHandleResponse(resp)
}

// getEntityTagCreateRequest creates the GetEntityTag request.
func (client *PortalRevisionClient) getEntityTagCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, portalRevisionID string, options *PortalRevisionGetEntityTagOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/portalRevisions/{portalRevisionId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if portalRevisionID == "" {
		return nil, errors.New("parameter portalRevisionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{portalRevisionId}", url.PathEscape(portalRevisionID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodHead, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getEntityTagHandleResponse handles the GetEntityTag response.
func (client *PortalRevisionClient) getEntityTagHandleResponse(resp *http.Response) (PortalRevisionGetEntityTagResponse, error) {
	result := PortalRevisionGetEntityTagResponse{RawResponse: resp}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		result.Success = true
	}
	return result, nil
}

// ListByService - Lists developer portal's revisions.
// If the operation fails it returns the *ErrorResponse error type.
func (client *PortalRevisionClient) ListByService(resourceGroupName string, serviceName string, options *PortalRevisionListByServiceOptions) *PortalRevisionListByServicePager {
	return &PortalRevisionListByServicePager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listByServiceCreateRequest(ctx, resourceGroupName, serviceName, options)
		},
		advancer: func(ctx context.Context, resp PortalRevisionListByServiceResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.PortalRevisionCollection.NextLink)
		},
	}
}

// listByServiceCreateRequest creates the ListByService request.
func (client *PortalRevisionClient) listByServiceCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, options *PortalRevisionListByServiceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/portalRevisions"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Skip != nil {
		reqQP.Set("$skip", strconv.FormatInt(int64(*options.Skip), 10))
	}
	reqQP.Set("api-version", "2021-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByServiceHandleResponse handles the ListByService response.
func (client *PortalRevisionClient) listByServiceHandleResponse(resp *http.Response) (PortalRevisionListByServiceResponse, error) {
	result := PortalRevisionListByServiceResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.PortalRevisionCollection); err != nil {
		return PortalRevisionListByServiceResponse{}, err
	}
	return result, nil
}

// listByServiceHandleError handles the ListByService error response.
func (client *PortalRevisionClient) listByServiceHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType.InnerError); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginUpdate - Updates the description of specified portal revision or makes it current.
// If the operation fails it returns the *ErrorResponse error type.
func (client *PortalRevisionClient) BeginUpdate(ctx context.Context, resourceGroupName string, serviceName string, portalRevisionID string, ifMatch string, parameters PortalRevisionContract, options *PortalRevisionBeginUpdateOptions) (PortalRevisionUpdatePollerResponse, error) {
	resp, err := client.update(ctx, resourceGroupName, serviceName, portalRevisionID, ifMatch, parameters, options)
	if err != nil {
		return PortalRevisionUpdatePollerResponse{}, err
	}
	result := PortalRevisionUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("PortalRevisionClient.Update", "location", resp, client.pl, client.updateHandleError)
	if err != nil {
		return PortalRevisionUpdatePollerResponse{}, err
	}
	result.Poller = &PortalRevisionUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// Update - Updates the description of specified portal revision or makes it current.
// If the operation fails it returns the *ErrorResponse error type.
func (client *PortalRevisionClient) update(ctx context.Context, resourceGroupName string, serviceName string, portalRevisionID string, ifMatch string, parameters PortalRevisionContract, options *PortalRevisionBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, serviceName, portalRevisionID, ifMatch, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.updateHandleError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *PortalRevisionClient) updateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, portalRevisionID string, ifMatch string, parameters PortalRevisionContract, options *PortalRevisionBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/portalRevisions/{portalRevisionId}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if portalRevisionID == "" {
		return nil, errors.New("parameter portalRevisionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{portalRevisionId}", url.PathEscape(portalRevisionID))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("If-Match", ifMatch)
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateHandleError handles the Update error response.
func (client *PortalRevisionClient) updateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType.InnerError); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
