// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapimanagement

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

// AllPoliciesClient contains the methods for the AllPolicies group.
// Don't use this type directly, use NewAllPoliciesClient() instead.
type AllPoliciesClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewAllPoliciesClient creates a new instance of AllPoliciesClient with the specified values.
//   - subscriptionID - The ID of the target subscription. The value must be an UUID.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAllPoliciesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*AllPoliciesClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AllPoliciesClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// NewListByServicePager - Status of all policies of API Management services.
//
// Generated from API version 2024-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - serviceName - The name of the API Management service.
//   - options - AllPoliciesClientListByServiceOptions contains the optional parameters for the AllPoliciesClient.NewListByServicePager
//     method.
func (client *AllPoliciesClient) NewListByServicePager(resourceGroupName string, serviceName string, options *AllPoliciesClientListByServiceOptions) *runtime.Pager[AllPoliciesClientListByServiceResponse] {
	return runtime.NewPager(runtime.PagingHandler[AllPoliciesClientListByServiceResponse]{
		More: func(page AllPoliciesClientListByServiceResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AllPoliciesClientListByServiceResponse) (AllPoliciesClientListByServiceResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "AllPoliciesClient.NewListByServicePager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByServiceCreateRequest(ctx, resourceGroupName, serviceName, options)
			}, nil)
			if err != nil {
				return AllPoliciesClientListByServiceResponse{}, err
			}
			return client.listByServiceHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByServiceCreateRequest creates the ListByService request.
func (client *AllPoliciesClient) listByServiceCreateRequest(ctx context.Context, resourceGroupName string, serviceName string, _ *AllPoliciesClientListByServiceOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ApiManagement/service/{serviceName}/allPolicies"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serviceName == "" {
		return nil, errors.New("parameter serviceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serviceName}", url.PathEscape(serviceName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByServiceHandleResponse handles the ListByService response.
func (client *AllPoliciesClient) listByServiceHandleResponse(resp *http.Response) (AllPoliciesClientListByServiceResponse, error) {
	result := AllPoliciesClientListByServiceResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AllPoliciesCollection); err != nil {
		return AllPoliciesClientListByServiceResponse{}, err
	}
	return result, nil
}
