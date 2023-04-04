//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armhybridconnectivity

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"strconv"
	"strings"
)

// EndpointsClient contains the methods for the Endpoints group.
// Don't use this type directly, use NewEndpointsClient() instead.
type EndpointsClient struct {
	internal *arm.Client
}

// NewEndpointsClient creates a new instance of EndpointsClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewEndpointsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*EndpointsClient, error) {
	cl, err := arm.NewClient(moduleName+".EndpointsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &EndpointsClient{
		internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update the endpoint to the target resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-05-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the resource to be connected.
//   - endpointName - The endpoint name.
//   - endpointResource - Endpoint details
//   - options - EndpointsClientCreateOrUpdateOptions contains the optional parameters for the EndpointsClient.CreateOrUpdate
//     method.
func (client *EndpointsClient) CreateOrUpdate(ctx context.Context, resourceURI string, endpointName string, endpointResource EndpointResource, options *EndpointsClientCreateOrUpdateOptions) (EndpointsClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceURI, endpointName, endpointResource, options)
	if err != nil {
		return EndpointsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EndpointsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return EndpointsClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *EndpointsClient) createOrUpdateCreateRequest(ctx context.Context, resourceURI string, endpointName string, endpointResource EndpointResource, options *EndpointsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.HybridConnectivity/endpoints/{endpointName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", endpointName)
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, endpointResource)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *EndpointsClient) createOrUpdateHandleResponse(resp *http.Response) (EndpointsClientCreateOrUpdateResponse, error) {
	result := EndpointsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EndpointResource); err != nil {
		return EndpointsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes the endpoint access to the target resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-05-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the resource to be connected.
//   - endpointName - The endpoint name.
//   - options - EndpointsClientDeleteOptions contains the optional parameters for the EndpointsClient.Delete method.
func (client *EndpointsClient) Delete(ctx context.Context, resourceURI string, endpointName string, options *EndpointsClientDeleteOptions) (EndpointsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceURI, endpointName, options)
	if err != nil {
		return EndpointsClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EndpointsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return EndpointsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return EndpointsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *EndpointsClient) deleteCreateRequest(ctx context.Context, resourceURI string, endpointName string, options *EndpointsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.HybridConnectivity/endpoints/{endpointName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", endpointName)
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the endpoint to the resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-05-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the resource to be connected.
//   - endpointName - The endpoint name.
//   - options - EndpointsClientGetOptions contains the optional parameters for the EndpointsClient.Get method.
func (client *EndpointsClient) Get(ctx context.Context, resourceURI string, endpointName string, options *EndpointsClientGetOptions) (EndpointsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceURI, endpointName, options)
	if err != nil {
		return EndpointsClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EndpointsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return EndpointsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *EndpointsClient) getCreateRequest(ctx context.Context, resourceURI string, endpointName string, options *EndpointsClientGetOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.HybridConnectivity/endpoints/{endpointName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", endpointName)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *EndpointsClient) getHandleResponse(resp *http.Response) (EndpointsClientGetResponse, error) {
	result := EndpointsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EndpointResource); err != nil {
		return EndpointsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List of endpoints to the target resource.
//
// Generated from API version 2022-05-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the resource to be connected.
//   - options - EndpointsClientListOptions contains the optional parameters for the EndpointsClient.NewListPager method.
func (client *EndpointsClient) NewListPager(resourceURI string, options *EndpointsClientListOptions) *runtime.Pager[EndpointsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[EndpointsClientListResponse]{
		More: func(page EndpointsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *EndpointsClientListResponse) (EndpointsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceURI, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return EndpointsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return EndpointsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return EndpointsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *EndpointsClient) listCreateRequest(ctx context.Context, resourceURI string, options *EndpointsClientListOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.HybridConnectivity/endpoints"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *EndpointsClient) listHandleResponse(resp *http.Response) (EndpointsClientListResponse, error) {
	result := EndpointsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EndpointsList); err != nil {
		return EndpointsClientListResponse{}, err
	}
	return result, nil
}

// ListCredentials - Gets the endpoint access credentials to the resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-05-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the resource to be connected.
//   - endpointName - The endpoint name.
//   - options - EndpointsClientListCredentialsOptions contains the optional parameters for the EndpointsClient.ListCredentials
//     method.
func (client *EndpointsClient) ListCredentials(ctx context.Context, resourceURI string, endpointName string, options *EndpointsClientListCredentialsOptions) (EndpointsClientListCredentialsResponse, error) {
	req, err := client.listCredentialsCreateRequest(ctx, resourceURI, endpointName, options)
	if err != nil {
		return EndpointsClientListCredentialsResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EndpointsClientListCredentialsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return EndpointsClientListCredentialsResponse{}, runtime.NewResponseError(resp)
	}
	return client.listCredentialsHandleResponse(resp)
}

// listCredentialsCreateRequest creates the ListCredentials request.
func (client *EndpointsClient) listCredentialsCreateRequest(ctx context.Context, resourceURI string, endpointName string, options *EndpointsClientListCredentialsOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.HybridConnectivity/endpoints/{endpointName}/listCredentials"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", endpointName)
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01-preview")
	if options != nil && options.Expiresin != nil {
		reqQP.Set("expiresin", strconv.FormatInt(*options.Expiresin, 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listCredentialsHandleResponse handles the ListCredentials response.
func (client *EndpointsClient) listCredentialsHandleResponse(resp *http.Response) (EndpointsClientListCredentialsResponse, error) {
	result := EndpointsClientListCredentialsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EndpointAccessResource); err != nil {
		return EndpointsClientListCredentialsResponse{}, err
	}
	return result, nil
}

// ListManagedProxyDetails - Fetches the managed proxy details
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-05-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the resource to be connected.
//   - endpointName - The endpoint name.
//   - managedProxyRequest - Object of type ManagedProxyRequest
//   - options - EndpointsClientListManagedProxyDetailsOptions contains the optional parameters for the EndpointsClient.ListManagedProxyDetails
//     method.
func (client *EndpointsClient) ListManagedProxyDetails(ctx context.Context, resourceURI string, endpointName string, managedProxyRequest ManagedProxyRequest, options *EndpointsClientListManagedProxyDetailsOptions) (EndpointsClientListManagedProxyDetailsResponse, error) {
	req, err := client.listManagedProxyDetailsCreateRequest(ctx, resourceURI, endpointName, managedProxyRequest, options)
	if err != nil {
		return EndpointsClientListManagedProxyDetailsResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EndpointsClientListManagedProxyDetailsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return EndpointsClientListManagedProxyDetailsResponse{}, runtime.NewResponseError(resp)
	}
	return client.listManagedProxyDetailsHandleResponse(resp)
}

// listManagedProxyDetailsCreateRequest creates the ListManagedProxyDetails request.
func (client *EndpointsClient) listManagedProxyDetailsCreateRequest(ctx context.Context, resourceURI string, endpointName string, managedProxyRequest ManagedProxyRequest, options *EndpointsClientListManagedProxyDetailsOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.HybridConnectivity/endpoints/{endpointName}/listManagedProxyDetails"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", endpointName)
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, managedProxyRequest)
}

// listManagedProxyDetailsHandleResponse handles the ListManagedProxyDetails response.
func (client *EndpointsClient) listManagedProxyDetailsHandleResponse(resp *http.Response) (EndpointsClientListManagedProxyDetailsResponse, error) {
	result := EndpointsClientListManagedProxyDetailsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ManagedProxyResource); err != nil {
		return EndpointsClientListManagedProxyDetailsResponse{}, err
	}
	return result, nil
}

// Update - Update the endpoint to the target resource.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-05-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the resource to be connected.
//   - endpointName - The endpoint name.
//   - endpointResource - Endpoint details
//   - options - EndpointsClientUpdateOptions contains the optional parameters for the EndpointsClient.Update method.
func (client *EndpointsClient) Update(ctx context.Context, resourceURI string, endpointName string, endpointResource EndpointResource, options *EndpointsClientUpdateOptions) (EndpointsClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceURI, endpointName, endpointResource, options)
	if err != nil {
		return EndpointsClientUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return EndpointsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return EndpointsClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *EndpointsClient) updateCreateRequest(ctx context.Context, resourceURI string, endpointName string, endpointResource EndpointResource, options *EndpointsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.HybridConnectivity/endpoints/{endpointName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	urlPath = strings.ReplaceAll(urlPath, "{endpointName}", endpointName)
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-05-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, endpointResource)
}

// updateHandleResponse handles the Update response.
func (client *EndpointsClient) updateHandleResponse(resp *http.Response) (EndpointsClientUpdateResponse, error) {
	result := EndpointsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.EndpointResource); err != nil {
		return EndpointsClientUpdateResponse{}, err
	}
	return result, nil
}
