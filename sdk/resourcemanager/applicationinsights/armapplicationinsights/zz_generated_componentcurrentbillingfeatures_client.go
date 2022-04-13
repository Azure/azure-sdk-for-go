//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapplicationinsights

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

// ComponentCurrentBillingFeaturesClient contains the methods for the ComponentCurrentBillingFeatures group.
// Don't use this type directly, use NewComponentCurrentBillingFeaturesClient() instead.
type ComponentCurrentBillingFeaturesClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewComponentCurrentBillingFeaturesClient creates a new instance of ComponentCurrentBillingFeaturesClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewComponentCurrentBillingFeaturesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ComponentCurrentBillingFeaturesClient, error) {
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
	client := &ComponentCurrentBillingFeaturesClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// Get - Returns current billing features for an Application Insights component.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// resourceName - The name of the Application Insights component resource.
// options - ComponentCurrentBillingFeaturesClientGetOptions contains the optional parameters for the ComponentCurrentBillingFeaturesClient.Get
// method.
func (client *ComponentCurrentBillingFeaturesClient) Get(ctx context.Context, resourceGroupName string, resourceName string, options *ComponentCurrentBillingFeaturesClientGetOptions) (ComponentCurrentBillingFeaturesClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return ComponentCurrentBillingFeaturesClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ComponentCurrentBillingFeaturesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ComponentCurrentBillingFeaturesClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ComponentCurrentBillingFeaturesClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, options *ComponentCurrentBillingFeaturesClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/currentbillingfeatures"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2015-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ComponentCurrentBillingFeaturesClient) getHandleResponse(resp *http.Response) (ComponentCurrentBillingFeaturesClientGetResponse, error) {
	result := ComponentCurrentBillingFeaturesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ComponentBillingFeatures); err != nil {
		return ComponentCurrentBillingFeaturesClientGetResponse{}, err
	}
	return result, nil
}

// Update - Update current billing features for an Application Insights component.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group. The name is case insensitive.
// resourceName - The name of the Application Insights component resource.
// billingFeaturesProperties - Properties that need to be specified to update billing features for an Application Insights
// component.
// options - ComponentCurrentBillingFeaturesClientUpdateOptions contains the optional parameters for the ComponentCurrentBillingFeaturesClient.Update
// method.
func (client *ComponentCurrentBillingFeaturesClient) Update(ctx context.Context, resourceGroupName string, resourceName string, billingFeaturesProperties ComponentBillingFeatures, options *ComponentCurrentBillingFeaturesClientUpdateOptions) (ComponentCurrentBillingFeaturesClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, resourceName, billingFeaturesProperties, options)
	if err != nil {
		return ComponentCurrentBillingFeaturesClientUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ComponentCurrentBillingFeaturesClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ComponentCurrentBillingFeaturesClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *ComponentCurrentBillingFeaturesClient) updateCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, billingFeaturesProperties ComponentBillingFeatures, options *ComponentCurrentBillingFeaturesClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/currentbillingfeatures"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2015-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, billingFeaturesProperties)
}

// updateHandleResponse handles the Update response.
func (client *ComponentCurrentBillingFeaturesClient) updateHandleResponse(resp *http.Response) (ComponentCurrentBillingFeaturesClientUpdateResponse, error) {
	result := ComponentCurrentBillingFeaturesClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ComponentBillingFeatures); err != nil {
		return ComponentCurrentBillingFeaturesClientUpdateResponse{}, err
	}
	return result, nil
}
