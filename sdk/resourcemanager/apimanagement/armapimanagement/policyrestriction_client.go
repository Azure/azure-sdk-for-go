// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapimanagement

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// PolicyRestrictionClient contains the methods for the PolicyRestriction group.
// Don't use this type directly, use NewPolicyRestrictionClient() instead.
type PolicyRestrictionClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewPolicyRestrictionClient creates a new instance of PolicyRestrictionClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPolicyRestrictionClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PolicyRestrictionClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PolicyRestrictionClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates the policy restriction configuration of the Api Management service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - policyRestrictionID - Policy restrictions after an entity level
//   - parameters - The policy restriction to apply.
//   - options - PolicyRestrictionClientCreateOrUpdateOptions contains the optional parameters for the PolicyRestrictionClient.CreateOrUpdate
//     method.
func (client *PolicyRestrictionClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, serviceName string, policyRestrictionID string, parameters PolicyRestrictionContract, options *PolicyRestrictionClientCreateOrUpdateOptions) (PolicyRestrictionClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "PolicyRestrictionClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serviceName, policyRestrictionID, parameters, options)
	if err != nil {
		return PolicyRestrictionClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PolicyRestrictionClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return PolicyRestrictionClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *PolicyRestrictionClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, policyRestrictionID string, parameters PolicyRestrictionContract, options *PolicyRestrictionClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/policyRestrictions/{policyRestrictionId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if policyRestrictionID == "" {
		return nil, errors.New("parameter policyRestrictionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyRestrictionId}", url.PathEscape(policyRestrictionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.IfMatch != nil {
		req.Raw().Header["If-Match"] = []string{*options.IfMatch}
	}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *PolicyRestrictionClient) createOrUpdateHandleResponse(resp *http.Response) (PolicyRestrictionClientCreateOrUpdateResponse, error) {
	result := PolicyRestrictionClientCreateOrUpdateResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.PolicyRestrictionContract); err != nil {
		return PolicyRestrictionClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the policy restriction configuration of the Api Management Service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - policyRestrictionID - Policy restrictions after an entity level
//   - options - PolicyRestrictionClientDeleteOptions contains the optional parameters for the PolicyRestrictionClient.Delete
//     method.
func (client *PolicyRestrictionClient) Delete(ctx context.Context, resourceGroupName string, serviceName string, policyRestrictionID string, options *PolicyRestrictionClientDeleteOptions) (PolicyRestrictionClientDeleteResponse, error) {
	var err error
	const operationName = "PolicyRestrictionClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serviceName, policyRestrictionID, options)
	if err != nil {
		return PolicyRestrictionClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PolicyRestrictionClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return PolicyRestrictionClientDeleteResponse{}, err
	}
	return PolicyRestrictionClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *PolicyRestrictionClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, policyRestrictionID string, options *PolicyRestrictionClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/policyRestrictions/{policyRestrictionId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if policyRestrictionID == "" {
		return nil, errors.New("parameter policyRestrictionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyRestrictionId}", url.PathEscape(policyRestrictionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.IfMatch != nil {
		req.Raw().Header["If-Match"] = []string{*options.IfMatch}
	}
	return req, nil
}

// Get - Get the policy restriction of the Api Management service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - policyRestrictionID - Policy restrictions after an entity level
//   - options - PolicyRestrictionClientGetOptions contains the optional parameters for the PolicyRestrictionClient.Get method.
func (client *PolicyRestrictionClient) Get(ctx context.Context, resourceGroupName string, serviceName string, policyRestrictionID string, options *PolicyRestrictionClientGetOptions) (PolicyRestrictionClientGetResponse, error) {
	var err error
	const operationName = "PolicyRestrictionClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, serviceName, policyRestrictionID, options)
	if err != nil {
		return PolicyRestrictionClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PolicyRestrictionClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PolicyRestrictionClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PolicyRestrictionClient) getCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, policyRestrictionID string, _ *PolicyRestrictionClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/policyRestrictions/{policyRestrictionId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if policyRestrictionID == "" {
		return nil, errors.New("parameter policyRestrictionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyRestrictionId}", url.PathEscape(policyRestrictionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PolicyRestrictionClient) getHandleResponse(resp *http.Response) (PolicyRestrictionClientGetResponse, error) {
	result := PolicyRestrictionClientGetResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.PolicyRestrictionContract); err != nil {
		return PolicyRestrictionClientGetResponse{}, err
	}
	return result, nil
}

// GetEntityTag - Gets the entity state (Etag) version of the policy restriction in the Api Management service.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - policyRestrictionID - Policy restrictions after an entity level
//   - options - PolicyRestrictionClientGetEntityTagOptions contains the optional parameters for the PolicyRestrictionClient.GetEntityTag
//     method.
func (client *PolicyRestrictionClient) GetEntityTag(ctx context.Context, resourceGroupName string, serviceName string, policyRestrictionID string, options *PolicyRestrictionClientGetEntityTagOptions) (PolicyRestrictionClientGetEntityTagResponse, error) {
	var err error
	const operationName = "PolicyRestrictionClient.GetEntityTag"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getEntityTagCreateRequest(ctx, resourceGroupName, serviceName, policyRestrictionID, options)
	if err != nil {
		return PolicyRestrictionClientGetEntityTagResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PolicyRestrictionClientGetEntityTagResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PolicyRestrictionClientGetEntityTagResponse{}, err
	}
	resp, err := client.getEntityTagHandleResponse(httpResp)
	return resp, err
}

// getEntityTagCreateRequest creates the GetEntityTag request.
func (client *PolicyRestrictionClient) getEntityTagCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, policyRestrictionID string, _ *PolicyRestrictionClientGetEntityTagOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/policyRestrictions/{policyRestrictionId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if policyRestrictionID == "" {
		return nil, errors.New("parameter policyRestrictionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyRestrictionId}", url.PathEscape(policyRestrictionID))
	req, err := runtime.NewRequest(ctx, http.MethodHead, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getEntityTagHandleResponse handles the GetEntityTag response.
func (client *PolicyRestrictionClient) getEntityTagHandleResponse(resp *http.Response) (PolicyRestrictionClientGetEntityTagResponse, error) {
	result := PolicyRestrictionClientGetEntityTagResponse{Success: resp.StatusCode >= 200 && resp.StatusCode < 300}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	return result, nil
}

// NewListByServicePager - Gets all policy restrictions of API Management services.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - options - PolicyRestrictionClientListByServiceOptions contains the optional parameters for the PolicyRestrictionClient.NewListByServicePager
//     method.
func (client *PolicyRestrictionClient) NewListByServicePager(resourceGroupName string, serviceName string, options *PolicyRestrictionClientListByServiceOptions) *runtime.Pager[PolicyRestrictionClientListByServiceResponse] {
	return runtime.NewPager(runtime.PagingHandler[PolicyRestrictionClientListByServiceResponse]{
		More: func(page PolicyRestrictionClientListByServiceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PolicyRestrictionClientListByServiceResponse) (PolicyRestrictionClientListByServiceResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "PolicyRestrictionClient.NewListByServicePager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByServiceCreateRequest(ctx, resourceGroupName, serviceName, options)
			}, nil)
			if err != nil {
				return PolicyRestrictionClientListByServiceResponse{}, err
			}
			return client.listByServiceHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByServiceCreateRequest creates the ListByService request.
func (client *PolicyRestrictionClient) listByServiceCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, _ *PolicyRestrictionClientListByServiceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/policyRestrictions"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByServiceHandleResponse handles the ListByService response.
func (client *PolicyRestrictionClient) listByServiceHandleResponse(resp *http.Response) (PolicyRestrictionClientListByServiceResponse, error) {
	result := PolicyRestrictionClientListByServiceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PolicyRestrictionCollection); err != nil {
		return PolicyRestrictionClientListByServiceResponse{}, err
	}
	return result, nil
}

// Update - Updates the policy restriction configuration of the Api Management service.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - policyRestrictionID - Policy restrictions after an entity level
//   - ifMatch - ETag of the Entity. ETag should match the current entity state from the header response of the GET request or
//     it should be * for unconditional update.
//   - parameters - The policy restriction to apply.
//   - options - PolicyRestrictionClientUpdateOptions contains the optional parameters for the PolicyRestrictionClient.Update
//     method.
func (client *PolicyRestrictionClient) Update(ctx context.Context, resourceGroupName string, serviceName string, policyRestrictionID string, ifMatch string, parameters PolicyRestrictionUpdateContract, options *PolicyRestrictionClientUpdateOptions) (PolicyRestrictionClientUpdateResponse, error) {
	var err error
	const operationName = "PolicyRestrictionClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, resourceGroupName, serviceName, policyRestrictionID, ifMatch, parameters, options)
	if err != nil {
		return PolicyRestrictionClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PolicyRestrictionClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PolicyRestrictionClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *PolicyRestrictionClient) updateCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, policyRestrictionID string, ifMatch string, parameters PolicyRestrictionUpdateContract, _ *PolicyRestrictionClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/policyRestrictions/{policyRestrictionId}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if policyRestrictionID == "" {
		return nil, errors.New("parameter policyRestrictionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{policyRestrictionId}", url.PathEscape(policyRestrictionID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	req.Raw().Header["If-Match"] = []string{ifMatch}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
		return nil, err
	}
	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *PolicyRestrictionClient) updateHandleResponse(resp *http.Response) (PolicyRestrictionClientUpdateResponse, error) {
	result := PolicyRestrictionClientUpdateResponse{}
	if val := resp.Header.Get("ETag"); val != "" {
		result.ETag = &val
	}
	if err := runtime.UnmarshalAsJSON(resp, &result.PolicyRestrictionContract); err != nil {
		return PolicyRestrictionClientUpdateResponse{}, err
	}
	return result, nil
}
