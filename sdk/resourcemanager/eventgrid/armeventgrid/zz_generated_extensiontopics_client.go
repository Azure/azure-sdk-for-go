//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armeventgrid

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

// ExtensionTopicsClient contains the methods for the ExtensionTopics group.
// Don't use this type directly, use NewExtensionTopicsClient() instead.
type ExtensionTopicsClient struct {
	host string
	pl   runtime.Pipeline
}

// NewExtensionTopicsClient creates a new instance of ExtensionTopicsClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewExtensionTopicsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ExtensionTopicsClient, error) {
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
	client := &ExtensionTopicsClient{
		host: ep,
		pl:   pl,
	}
	return client, nil
}

// Get - Get the properties of an extension topic.
// If the operation fails it returns an *azcore.ResponseError type.
// scope - The identifier of the resource to which extension topic is queried. The scope can be a subscription, or a resource
// group, or a top level resource belonging to a resource provider namespace. For
// example, use '/subscriptions/{subscriptionId}/' for a subscription, '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}'
// for a resource group, and
// '/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}'
// for Azure resource.
// options - ExtensionTopicsClientGetOptions contains the optional parameters for the ExtensionTopicsClient.Get method.
func (client *ExtensionTopicsClient) Get(ctx context.Context, scope string, options *ExtensionTopicsClientGetOptions) (ExtensionTopicsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, scope, options)
	if err != nil {
		return ExtensionTopicsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ExtensionTopicsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ExtensionTopicsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ExtensionTopicsClient) getCreateRequest(ctx context.Context, scope string, options *ExtensionTopicsClientGetOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.EventGrid/extensionTopics/default"
	if scope == "" {
		return nil, errors.New("parameter scope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scope}", url.PathEscape(scope))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-15-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ExtensionTopicsClient) getHandleResponse(resp *http.Response) (ExtensionTopicsClientGetResponse, error) {
	result := ExtensionTopicsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ExtensionTopic); err != nil {
		return ExtensionTopicsClientGetResponse{}, err
	}
	return result, nil
}
