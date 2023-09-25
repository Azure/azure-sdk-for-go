//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armportal

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
)

// ListTenantConfigurationViolationsClient contains the methods for the ListTenantConfigurationViolations group.
// Don't use this type directly, use NewListTenantConfigurationViolationsClient() instead.
type ListTenantConfigurationViolationsClient struct {
	internal *arm.Client
}

// NewListTenantConfigurationViolationsClient creates a new instance of ListTenantConfigurationViolationsClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewListTenantConfigurationViolationsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ListTenantConfigurationViolationsClient, error) {
	cl, err := arm.NewClient(moduleName+".ListTenantConfigurationViolationsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ListTenantConfigurationViolationsClient{
	internal: cl,
	}
	return client, nil
}

// NewListPager - Gets list of items that violate tenant's configuration.
//
// Generated from API version 2020-09-01-preview
//   - options - ListTenantConfigurationViolationsClientListOptions contains the optional parameters for the ListTenantConfigurationViolationsClient.NewListPager
//     method.
func (client *ListTenantConfigurationViolationsClient) NewListPager(options *ListTenantConfigurationViolationsClientListOptions) (*runtime.Pager[ListTenantConfigurationViolationsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ListTenantConfigurationViolationsClientListResponse]{
		More: func(page ListTenantConfigurationViolationsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ListTenantConfigurationViolationsClientListResponse) (ListTenantConfigurationViolationsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ListTenantConfigurationViolationsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ListTenantConfigurationViolationsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListTenantConfigurationViolationsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ListTenantConfigurationViolationsClient) listCreateRequest(ctx context.Context, options *ListTenantConfigurationViolationsClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Portal/listTenantConfigurationViolations"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ListTenantConfigurationViolationsClient) listHandleResponse(resp *http.Response) (ListTenantConfigurationViolationsClientListResponse, error) {
	result := ListTenantConfigurationViolationsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ViolationsList); err != nil {
		return ListTenantConfigurationViolationsClientListResponse{}, err
	}
	return result, nil
}

