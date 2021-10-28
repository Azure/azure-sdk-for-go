//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

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

// LoadBalancerProbesClient contains the methods for the LoadBalancerProbes group.
// Don't use this type directly, use NewLoadBalancerProbesClient() instead.
type LoadBalancerProbesClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewLoadBalancerProbesClient creates a new instance of LoadBalancerProbesClient with the specified values.
func NewLoadBalancerProbesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *LoadBalancerProbesClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &LoadBalancerProbesClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// Get - Gets load balancer probe.
// If the operation fails it returns the *CloudError error type.
func (client *LoadBalancerProbesClient) Get(ctx context.Context, resourceGroupName string, loadBalancerName string, probeName string, options *LoadBalancerProbesGetOptions) (LoadBalancerProbesGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, loadBalancerName, probeName, options)
	if err != nil {
		return LoadBalancerProbesGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return LoadBalancerProbesGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return LoadBalancerProbesGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *LoadBalancerProbesClient) getCreateRequest(ctx context.Context, resourceGroupName string, loadBalancerName string, probeName string, options *LoadBalancerProbesGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/loadBalancers/{loadBalancerName}/probes/{probeName}"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if loadBalancerName == "" {
		return nil, errors.New("parameter loadBalancerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{loadBalancerName}", url.PathEscape(loadBalancerName))
	if probeName == "" {
		return nil, errors.New("parameter probeName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{probeName}", url.PathEscape(probeName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *LoadBalancerProbesClient) getHandleResponse(resp *http.Response) (LoadBalancerProbesGetResponse, error) {
	result := LoadBalancerProbesGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.Probe); err != nil {
		return LoadBalancerProbesGetResponse{}, err
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *LoadBalancerProbesClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// List - Gets all the load balancer probes.
// If the operation fails it returns the *CloudError error type.
func (client *LoadBalancerProbesClient) List(resourceGroupName string, loadBalancerName string, options *LoadBalancerProbesListOptions) *LoadBalancerProbesListPager {
	return &LoadBalancerProbesListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, loadBalancerName, options)
		},
		advancer: func(ctx context.Context, resp LoadBalancerProbesListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.LoadBalancerProbeListResult.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *LoadBalancerProbesClient) listCreateRequest(ctx context.Context, resourceGroupName string, loadBalancerName string, options *LoadBalancerProbesListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/loadBalancers/{loadBalancerName}/probes"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if loadBalancerName == "" {
		return nil, errors.New("parameter loadBalancerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{loadBalancerName}", url.PathEscape(loadBalancerName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *LoadBalancerProbesClient) listHandleResponse(resp *http.Response) (LoadBalancerProbesListResponse, error) {
	result := LoadBalancerProbesListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.LoadBalancerProbeListResult); err != nil {
		return LoadBalancerProbesListResponse{}, err
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *LoadBalancerProbesClient) listHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
