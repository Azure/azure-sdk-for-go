//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

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

// AutomationsClient contains the methods for the Automations group.
// Don't use this type directly, use NewAutomationsClient() instead.
type AutomationsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewAutomationsClient creates a new instance of AutomationsClient with the specified values.
//   - subscriptionID - Azure subscription ID
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAutomationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AutomationsClient, error) {
	cl, err := arm.NewClient(moduleName+".AutomationsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AutomationsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates a security automation. If a security automation is already created and a subsequent
// request is issued for the same automation id, then it will be updated.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-01-01-preview
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - automationName - The security automation name.
//   - automation - The security automation resource
//   - options - AutomationsClientCreateOrUpdateOptions contains the optional parameters for the AutomationsClient.CreateOrUpdate
//     method.
func (client *AutomationsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, automationName string, automation Automation, options *AutomationsClientCreateOrUpdateOptions) (AutomationsClientCreateOrUpdateResponse, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, automationName, automation, options)
	if err != nil {
		return AutomationsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AutomationsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return AutomationsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *AutomationsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, automationName string, automation Automation, options *AutomationsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/automations/{automationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if automationName == "" {
		return nil, errors.New("parameter automationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{automationName}", url.PathEscape(automationName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, automation); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *AutomationsClient) createOrUpdateHandleResponse(resp *http.Response) (AutomationsClientCreateOrUpdateResponse, error) {
	result := AutomationsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Automation); err != nil {
		return AutomationsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a security automation.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-01-01-preview
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - automationName - The security automation name.
//   - options - AutomationsClientDeleteOptions contains the optional parameters for the AutomationsClient.Delete method.
func (client *AutomationsClient) Delete(ctx context.Context, resourceGroupName string, automationName string, options *AutomationsClientDeleteOptions) (AutomationsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, automationName, options)
	if err != nil {
		return AutomationsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AutomationsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return AutomationsClientDeleteResponse{}, err
	}
	return AutomationsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *AutomationsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, automationName string, options *AutomationsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/automations/{automationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if automationName == "" {
		return nil, errors.New("parameter automationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{automationName}", url.PathEscape(automationName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Retrieves information about the model of a security automation.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-01-01-preview
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - automationName - The security automation name.
//   - options - AutomationsClientGetOptions contains the optional parameters for the AutomationsClient.Get method.
func (client *AutomationsClient) Get(ctx context.Context, resourceGroupName string, automationName string, options *AutomationsClientGetOptions) (AutomationsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, automationName, options)
	if err != nil {
		return AutomationsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AutomationsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AutomationsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *AutomationsClient) getCreateRequest(ctx context.Context, resourceGroupName string, automationName string, options *AutomationsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/automations/{automationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if automationName == "" {
		return nil, errors.New("parameter automationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{automationName}", url.PathEscape(automationName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AutomationsClient) getHandleResponse(resp *http.Response) (AutomationsClientGetResponse, error) {
	result := AutomationsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Automation); err != nil {
		return AutomationsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all the security automations in the specified subscription. Use the 'nextLink' property in the response
// to get the next page of security automations for the specified subscription.
//
// Generated from API version 2019-01-01-preview
//   - options - AutomationsClientListOptions contains the optional parameters for the AutomationsClient.NewListPager method.
func (client *AutomationsClient) NewListPager(options *AutomationsClientListOptions) (*runtime.Pager[AutomationsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[AutomationsClientListResponse]{
		More: func(page AutomationsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AutomationsClientListResponse) (AutomationsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AutomationsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return AutomationsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AutomationsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *AutomationsClient) listCreateRequest(ctx context.Context, options *AutomationsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/automations"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *AutomationsClient) listHandleResponse(resp *http.Response) (AutomationsClientListResponse, error) {
	result := AutomationsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AutomationList); err != nil {
		return AutomationsClientListResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Lists all the security automations in the specified resource group. Use the 'nextLink' property
// in the response to get the next page of security automations for the specified resource group.
//
// Generated from API version 2019-01-01-preview
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - options - AutomationsClientListByResourceGroupOptions contains the optional parameters for the AutomationsClient.NewListByResourceGroupPager
//     method.
func (client *AutomationsClient) NewListByResourceGroupPager(resourceGroupName string, options *AutomationsClientListByResourceGroupOptions) (*runtime.Pager[AutomationsClientListByResourceGroupResponse]) {
	return runtime.NewPager(runtime.PagingHandler[AutomationsClientListByResourceGroupResponse]{
		More: func(page AutomationsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AutomationsClientListByResourceGroupResponse) (AutomationsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AutomationsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return AutomationsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AutomationsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *AutomationsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *AutomationsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/automations"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *AutomationsClient) listByResourceGroupHandleResponse(resp *http.Response) (AutomationsClientListByResourceGroupResponse, error) {
	result := AutomationsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AutomationList); err != nil {
		return AutomationsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// Validate - Validates the security automation model before create or update. Any validation errors are returned to the client.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-01-01-preview
//   - resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
//   - automationName - The security automation name.
//   - automation - The security automation resource
//   - options - AutomationsClientValidateOptions contains the optional parameters for the AutomationsClient.Validate method.
func (client *AutomationsClient) Validate(ctx context.Context, resourceGroupName string, automationName string, automation Automation, options *AutomationsClientValidateOptions) (AutomationsClientValidateResponse, error) {
	var err error
	req, err := client.validateCreateRequest(ctx, resourceGroupName, automationName, automation, options)
	if err != nil {
		return AutomationsClientValidateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AutomationsClientValidateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AutomationsClientValidateResponse{}, err
	}
	resp, err := client.validateHandleResponse(httpResp)
	return resp, err
}

// validateCreateRequest creates the Validate request.
func (client *AutomationsClient) validateCreateRequest(ctx context.Context, resourceGroupName string, automationName string, automation Automation, options *AutomationsClientValidateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Security/automations/{automationName}/validate"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if automationName == "" {
		return nil, errors.New("parameter automationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{automationName}", url.PathEscape(automationName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, automation); err != nil {
	return nil, err
}
	return req, nil
}

// validateHandleResponse handles the Validate response.
func (client *AutomationsClient) validateHandleResponse(resp *http.Response) (AutomationsClientValidateResponse, error) {
	result := AutomationsClientValidateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AutomationValidationStatus); err != nil {
		return AutomationsClientValidateResponse{}, err
	}
	return result, nil
}

