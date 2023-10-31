//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armazurestackhci

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"strings"
)

// GuestAgentClient contains the methods for the GuestAgent group.
// Don't use this type directly, use NewGuestAgentClient() instead.
type GuestAgentClient struct {
	internal *arm.Client
}

// NewGuestAgentClient creates a new instance of GuestAgentClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewGuestAgentClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*GuestAgentClient, error) {
	cl, err := arm.NewClient(moduleName+".GuestAgentClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &GuestAgentClient{
		internal: cl,
	}
	return client, nil
}

// BeginCreate - Create Or Update GuestAgent.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the Hybrid Compute machine resource to be extended.
//   - options - GuestAgentClientBeginCreateOptions contains the optional parameters for the GuestAgentClient.BeginCreate method.
func (client *GuestAgentClient) BeginCreate(ctx context.Context, resourceURI string, options *GuestAgentClientBeginCreateOptions) (*runtime.Poller[GuestAgentClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceURI, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[GuestAgentClientCreateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[GuestAgentClientCreateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Create - Create Or Update GuestAgent.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01-preview
func (client *GuestAgentClient) create(ctx context.Context, resourceURI string, options *GuestAgentClientBeginCreateOptions) (*http.Response, error) {
	var err error
	req, err := client.createCreateRequest(ctx, resourceURI, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createCreateRequest creates the Create request.
func (client *GuestAgentClient) createCreateRequest(ctx context.Context, resourceURI string, options *GuestAgentClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.AzureStackHCI/virtualMachineInstances/default/guestAgents/default"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if options != nil && options.Body != nil {
		if err := runtime.MarshalAsJSON(req, *options.Body); err != nil {
			return nil, err
		}
		return req, nil
	}
	return req, nil
}

// BeginDelete - Implements GuestAgent DELETE method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the Hybrid Compute machine resource to be extended.
//   - options - GuestAgentClientBeginDeleteOptions contains the optional parameters for the GuestAgentClient.BeginDelete method.
func (client *GuestAgentClient) BeginDelete(ctx context.Context, resourceURI string, options *GuestAgentClientBeginDeleteOptions) (*runtime.Poller[GuestAgentClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceURI, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[GuestAgentClientDeleteResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[GuestAgentClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Implements GuestAgent DELETE method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01-preview
func (client *GuestAgentClient) deleteOperation(ctx context.Context, resourceURI string, options *GuestAgentClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceURI, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *GuestAgentClient) deleteCreateRequest(ctx context.Context, resourceURI string, options *GuestAgentClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.AzureStackHCI/virtualMachineInstances/default/guestAgents/default"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Implements GuestAgent GET method.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-09-01-preview
//   - resourceURI - The fully qualified Azure Resource manager identifier of the Hybrid Compute machine resource to be extended.
//   - options - GuestAgentClientGetOptions contains the optional parameters for the GuestAgentClient.Get method.
func (client *GuestAgentClient) Get(ctx context.Context, resourceURI string, options *GuestAgentClientGetOptions) (GuestAgentClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceURI, options)
	if err != nil {
		return GuestAgentClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return GuestAgentClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return GuestAgentClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *GuestAgentClient) getCreateRequest(ctx context.Context, resourceURI string, options *GuestAgentClientGetOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.AzureStackHCI/virtualMachineInstances/default/guestAgents/default"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-09-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *GuestAgentClient) getHandleResponse(resp *http.Response) (GuestAgentClientGetResponse, error) {
	result := GuestAgentClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.GuestAgent); err != nil {
		return GuestAgentClientGetResponse{}, err
	}
	return result, nil
}