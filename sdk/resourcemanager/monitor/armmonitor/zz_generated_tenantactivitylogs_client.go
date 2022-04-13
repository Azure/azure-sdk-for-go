//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmonitor

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
)

// TenantActivityLogsClient contains the methods for the TenantActivityLogs group.
// Don't use this type directly, use NewTenantActivityLogsClient() instead.
type TenantActivityLogsClient struct {
	host string
	pl   runtime.Pipeline
}

// NewTenantActivityLogsClient creates a new instance of TenantActivityLogsClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewTenantActivityLogsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*TenantActivityLogsClient, error) {
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
	client := &TenantActivityLogsClient{
		host: ep,
		pl:   pl,
	}
	return client, nil
}

// List - Gets the Activity Logs for the Tenant. Everything that is applicable to the API to get the Activity Logs for the
// subscription is applicable to this API (the parameters, $filter, etc.). One thing to
// point out here is that this API does not retrieve the logs at the individual subscription of the tenant but only surfaces
// the logs that were generated at the tenant level.
// If the operation fails it returns an *azcore.ResponseError type.
// options - TenantActivityLogsClientListOptions contains the optional parameters for the TenantActivityLogsClient.List method.
func (client *TenantActivityLogsClient) List(options *TenantActivityLogsClientListOptions) *runtime.Pager[TenantActivityLogsClientListResponse] {
	return runtime.NewPager(runtime.PageProcessor[TenantActivityLogsClientListResponse]{
		More: func(page TenantActivityLogsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *TenantActivityLogsClientListResponse) (TenantActivityLogsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return TenantActivityLogsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return TenantActivityLogsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return TenantActivityLogsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *TenantActivityLogsClient) listCreateRequest(ctx context.Context, options *TenantActivityLogsClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Insights/eventtypes/management/values"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2015-04-01")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Select != nil {
		reqQP.Set("$select", *options.Select)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *TenantActivityLogsClient) listHandleResponse(resp *http.Response) (TenantActivityLogsClientListResponse, error) {
	result := TenantActivityLogsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EventDataCollection); err != nil {
		return TenantActivityLogsClientListResponse{}, err
	}
	return result, nil
}
