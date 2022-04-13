//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// CustomAssessmentAutomationsClient contains the methods for the CustomAssessmentAutomations group.
// Don't use this type directly, use NewCustomAssessmentAutomationsClient() instead.
type CustomAssessmentAutomationsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewCustomAssessmentAutomationsClient creates a new instance of CustomAssessmentAutomationsClient with the specified values.
// subscriptionID - Azure subscription ID
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewCustomAssessmentAutomationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CustomAssessmentAutomationsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublicCloud.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &CustomAssessmentAutomationsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Create - Creates or updates a custom assessment automation for the provided subscription. Please note that providing an
// existing custom assessment automation will replace the existing record.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
// customAssessmentAutomationName - Name of the Custom Assessment Automation.
// customAssessmentAutomationBody - Custom Assessment Automation body
// options - CustomAssessmentAutomationsClientCreateOptions contains the optional parameters for the CustomAssessmentAutomationsClient.Create
// method.
func (client *CustomAssessmentAutomationsClient) Create(ctx context.Context, resourceGroupName string, customAssessmentAutomationName string, customAssessmentAutomationBody CustomAssessmentAutomationRequest, options *CustomAssessmentAutomationsClientCreateOptions) (CustomAssessmentAutomationsClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, customAssessmentAutomationName, customAssessmentAutomationBody, options)
	if err != nil {
		return CustomAssessmentAutomationsClientCreateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return CustomAssessmentAutomationsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return CustomAssessmentAutomationsClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *CustomAssessmentAutomationsClient) createCreateRequest(ctx context.Context, resourceGroupName string, customAssessmentAutomationName string, customAssessmentAutomationBody CustomAssessmentAutomationRequest, options *CustomAssessmentAutomationsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Security/customAssessmentAutomations/{customAssessmentAutomationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if customAssessmentAutomationName == "" {
		return nil, errors.New("parameter customAssessmentAutomationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{customAssessmentAutomationName}", url.PathEscape(customAssessmentAutomationName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-07-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, customAssessmentAutomationBody)
}

// createHandleResponse handles the Create response.
func (client *CustomAssessmentAutomationsClient) createHandleResponse(resp *http.Response) (CustomAssessmentAutomationsClientCreateResponse, error) {
	result := CustomAssessmentAutomationsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CustomAssessmentAutomation); err != nil {
		return CustomAssessmentAutomationsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a custom assessment automation by name for a provided subscription
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
// customAssessmentAutomationName - Name of the Custom Assessment Automation.
// options - CustomAssessmentAutomationsClientDeleteOptions contains the optional parameters for the CustomAssessmentAutomationsClient.Delete
// method.
func (client *CustomAssessmentAutomationsClient) Delete(ctx context.Context, resourceGroupName string, customAssessmentAutomationName string, options *CustomAssessmentAutomationsClientDeleteOptions) (CustomAssessmentAutomationsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, customAssessmentAutomationName, options)
	if err != nil {
		return CustomAssessmentAutomationsClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return CustomAssessmentAutomationsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return CustomAssessmentAutomationsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return CustomAssessmentAutomationsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *CustomAssessmentAutomationsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, customAssessmentAutomationName string, options *CustomAssessmentAutomationsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Security/customAssessmentAutomations/{customAssessmentAutomationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if customAssessmentAutomationName == "" {
		return nil, errors.New("parameter customAssessmentAutomationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{customAssessmentAutomationName}", url.PathEscape(customAssessmentAutomationName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-07-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Gets a single custom assessment automation by name for the provided subscription and resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
// customAssessmentAutomationName - Name of the Custom Assessment Automation.
// options - CustomAssessmentAutomationsClientGetOptions contains the optional parameters for the CustomAssessmentAutomationsClient.Get
// method.
func (client *CustomAssessmentAutomationsClient) Get(ctx context.Context, resourceGroupName string, customAssessmentAutomationName string, options *CustomAssessmentAutomationsClientGetOptions) (CustomAssessmentAutomationsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, customAssessmentAutomationName, options)
	if err != nil {
		return CustomAssessmentAutomationsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return CustomAssessmentAutomationsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return CustomAssessmentAutomationsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *CustomAssessmentAutomationsClient) getCreateRequest(ctx context.Context, resourceGroupName string, customAssessmentAutomationName string, options *CustomAssessmentAutomationsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Security/customAssessmentAutomations/{customAssessmentAutomationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if customAssessmentAutomationName == "" {
		return nil, errors.New("parameter customAssessmentAutomationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{customAssessmentAutomationName}", url.PathEscape(customAssessmentAutomationName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-07-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *CustomAssessmentAutomationsClient) getHandleResponse(resp *http.Response) (CustomAssessmentAutomationsClientGetResponse, error) {
	result := CustomAssessmentAutomationsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CustomAssessmentAutomation); err != nil {
		return CustomAssessmentAutomationsClientGetResponse{}, err
	}
	return result, nil
}

// ListByResourceGroup - List custom assessment automations by provided subscription and resource group
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group within the user's subscription. The name is case insensitive.
// options - CustomAssessmentAutomationsClientListByResourceGroupOptions contains the optional parameters for the CustomAssessmentAutomationsClient.ListByResourceGroup
// method.
func (client *CustomAssessmentAutomationsClient) ListByResourceGroup(resourceGroupName string, options *CustomAssessmentAutomationsClientListByResourceGroupOptions) *runtime.Pager[CustomAssessmentAutomationsClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PageProcessor[CustomAssessmentAutomationsClientListByResourceGroupResponse]{
		More: func(page CustomAssessmentAutomationsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CustomAssessmentAutomationsClientListByResourceGroupResponse) (CustomAssessmentAutomationsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CustomAssessmentAutomationsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return CustomAssessmentAutomationsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CustomAssessmentAutomationsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *CustomAssessmentAutomationsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *CustomAssessmentAutomationsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourcegroups/{resourceGroupName}/providers/Microsoft.Security/customAssessmentAutomations"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-07-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *CustomAssessmentAutomationsClient) listByResourceGroupHandleResponse(resp *http.Response) (CustomAssessmentAutomationsClientListByResourceGroupResponse, error) {
	result := CustomAssessmentAutomationsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CustomAssessmentAutomationsListResult); err != nil {
		return CustomAssessmentAutomationsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// ListBySubscription - List custom assessment automations by provided subscription
// If the operation fails it returns an *azcore.ResponseError type.
// options - CustomAssessmentAutomationsClientListBySubscriptionOptions contains the optional parameters for the CustomAssessmentAutomationsClient.ListBySubscription
// method.
func (client *CustomAssessmentAutomationsClient) ListBySubscription(options *CustomAssessmentAutomationsClientListBySubscriptionOptions) *runtime.Pager[CustomAssessmentAutomationsClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PageProcessor[CustomAssessmentAutomationsClientListBySubscriptionResponse]{
		More: func(page CustomAssessmentAutomationsClientListBySubscriptionResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *CustomAssessmentAutomationsClientListBySubscriptionResponse) (CustomAssessmentAutomationsClientListBySubscriptionResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listBySubscriptionCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return CustomAssessmentAutomationsClientListBySubscriptionResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return CustomAssessmentAutomationsClientListBySubscriptionResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return CustomAssessmentAutomationsClientListBySubscriptionResponse{}, runtime.NewResponseError(resp)
			}
			return client.listBySubscriptionHandleResponse(resp)
		},
	})
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *CustomAssessmentAutomationsClient) listBySubscriptionCreateRequest(ctx context.Context, options *CustomAssessmentAutomationsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/customAssessmentAutomations"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-07-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *CustomAssessmentAutomationsClient) listBySubscriptionHandleResponse(resp *http.Response) (CustomAssessmentAutomationsClientListBySubscriptionResponse, error) {
	result := CustomAssessmentAutomationsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CustomAssessmentAutomationsListResult); err != nil {
		return CustomAssessmentAutomationsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}
