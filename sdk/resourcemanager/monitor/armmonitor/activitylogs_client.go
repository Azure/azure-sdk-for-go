//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmonitor

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

// ActivityLogsClient contains the methods for the ActivityLogs group.
// Don't use this type directly, use NewActivityLogsClient() instead.
type ActivityLogsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewActivityLogsClient creates a new instance of ActivityLogsClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewActivityLogsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ActivityLogsClient, error) {
	cl, err := arm.NewClient(moduleName+".ActivityLogsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ActivityLogsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// NewListPager - Provides the list of records from the activity logs.
//
// Generated from API version 2015-04-01
//   - filter - Reduces the set of data collected.
//     This argument is required and it also requires at least the start date/time.
//     The $filter argument is very restricted and allows only the following patterns.
//     - List events for a resource group: $filter=eventTimestamp ge '2014-07-16T04:36:37.6407898Z' and eventTimestamp le '2014-07-20T04:36:37.6407898Z'
//     and resourceGroupName eq 'resourceGroupName'.
//     - List events for resource: $filter=eventTimestamp ge '2014-07-16T04:36:37.6407898Z' and eventTimestamp le '2014-07-20T04:36:37.6407898Z'
//     and resourceUri eq 'resourceURI'.
//     - List events for a subscription in a time range: $filter=eventTimestamp ge '2014-07-16T04:36:37.6407898Z' and eventTimestamp
//     le '2014-07-20T04:36:37.6407898Z'.
//     - List events for a resource provider: $filter=eventTimestamp ge '2014-07-16T04:36:37.6407898Z' and eventTimestamp le '2014-07-20T04:36:37.6407898Z'
//     and resourceProvider eq 'resourceProviderName'.
//     - List events for a correlation Id: $filter=eventTimestamp ge '2014-07-16T04:36:37.6407898Z' and eventTimestamp le '2014-07-20T04:36:37.6407898Z'
//     and correlationId eq 'correlationID'.
//     NOTE: No other syntax is allowed.
//   - options - ActivityLogsClientListOptions contains the optional parameters for the ActivityLogsClient.NewListPager method.
func (client *ActivityLogsClient) NewListPager(filter string, options *ActivityLogsClientListOptions) (*runtime.Pager[ActivityLogsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ActivityLogsClientListResponse]{
		More: func(page ActivityLogsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ActivityLogsClientListResponse) (ActivityLogsClientListResponse, error) {
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ActivityLogsClient.NewListPager")
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, filter, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ActivityLogsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ActivityLogsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ActivityLogsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *ActivityLogsClient) listCreateRequest(ctx context.Context, filter string, options *ActivityLogsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Insights/eventtypes/management/values"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2015-04-01")
	reqQP.Set("$filter", filter)
	if options != nil && options.Select != nil {
		reqQP.Set("$select", *options.Select)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ActivityLogsClient) listHandleResponse(resp *http.Response) (ActivityLogsClientListResponse, error) {
	result := ActivityLogsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EventDataCollection); err != nil {
		return ActivityLogsClientListResponse{}, err
	}
	return result, nil
}

