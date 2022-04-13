//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmariadb

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

// LocationBasedPerformanceTierClient contains the methods for the LocationBasedPerformanceTier group.
// Don't use this type directly, use NewLocationBasedPerformanceTierClient() instead.
type LocationBasedPerformanceTierClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewLocationBasedPerformanceTierClient creates a new instance of LocationBasedPerformanceTierClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewLocationBasedPerformanceTierClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*LocationBasedPerformanceTierClient, error) {
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
	client := &LocationBasedPerformanceTierClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// List - List all the performance tiers at specified location in a given subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// locationName - The name of the location.
// options - LocationBasedPerformanceTierClientListOptions contains the optional parameters for the LocationBasedPerformanceTierClient.List
// method.
func (client *LocationBasedPerformanceTierClient) List(locationName string, options *LocationBasedPerformanceTierClientListOptions) *runtime.Pager[LocationBasedPerformanceTierClientListResponse] {
	return runtime.NewPager(runtime.PageProcessor[LocationBasedPerformanceTierClientListResponse]{
		More: func(page LocationBasedPerformanceTierClientListResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *LocationBasedPerformanceTierClientListResponse) (LocationBasedPerformanceTierClientListResponse, error) {
			req, err := client.listCreateRequest(ctx, locationName, options)
			if err != nil {
				return LocationBasedPerformanceTierClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return LocationBasedPerformanceTierClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return LocationBasedPerformanceTierClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *LocationBasedPerformanceTierClient) listCreateRequest(ctx context.Context, locationName string, options *LocationBasedPerformanceTierClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.DBforMariaDB/locations/{locationName}/performanceTiers"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *LocationBasedPerformanceTierClient) listHandleResponse(resp *http.Response) (LocationBasedPerformanceTierClientListResponse, error) {
	result := LocationBasedPerformanceTierClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PerformanceTierListResult); err != nil {
		return LocationBasedPerformanceTierClientListResponse{}, err
	}
	return result, nil
}
