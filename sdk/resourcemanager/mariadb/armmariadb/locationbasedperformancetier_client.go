//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmariadb

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

// LocationBasedPerformanceTierClient contains the methods for the LocationBasedPerformanceTier group.
// Don't use this type directly, use NewLocationBasedPerformanceTierClient() instead.
type LocationBasedPerformanceTierClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewLocationBasedPerformanceTierClient creates a new instance of LocationBasedPerformanceTierClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewLocationBasedPerformanceTierClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*LocationBasedPerformanceTierClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &LocationBasedPerformanceTierClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// NewListPager - List all the performance tiers at specified location in a given subscription.
//
// Generated from API version 2018-06-01
//   - locationName - The name of the location.
//   - options - LocationBasedPerformanceTierClientListOptions contains the optional parameters for the LocationBasedPerformanceTierClient.NewListPager
//     method.
func (client *LocationBasedPerformanceTierClient) NewListPager(locationName string, options *LocationBasedPerformanceTierClientListOptions) *runtime.Pager[LocationBasedPerformanceTierClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[LocationBasedPerformanceTierClientListResponse]{
		More: func(page LocationBasedPerformanceTierClientListResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *LocationBasedPerformanceTierClientListResponse) (LocationBasedPerformanceTierClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "LocationBasedPerformanceTierClient.NewListPager")
			req, err := client.listCreateRequest(ctx, locationName, options)
			if err != nil {
				return LocationBasedPerformanceTierClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return LocationBasedPerformanceTierClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return LocationBasedPerformanceTierClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-06-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
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
