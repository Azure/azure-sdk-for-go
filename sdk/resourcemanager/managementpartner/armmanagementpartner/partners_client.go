//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmanagementpartner

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
)

// PartnersClient contains the methods for the Partners group.
// Don't use this type directly, use NewPartnersClient() instead.
type PartnersClient struct {
	internal *arm.Client
}

// NewPartnersClient creates a new instance of PartnersClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPartnersClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*PartnersClient, error) {
	cl, err := arm.NewClient(moduleName+".PartnersClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PartnersClient{
	internal: cl,
	}
	return client, nil
}

// Get - Get the management partner using the objectId and tenantId.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2018-02-01
//   - options - PartnersClientGetOptions contains the optional parameters for the PartnersClient.Get method.
func (client *PartnersClient) Get(ctx context.Context, options *PartnersClientGetOptions) (PartnersClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, options)
	if err != nil {
		return PartnersClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PartnersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PartnersClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PartnersClient) getCreateRequest(ctx context.Context, options *PartnersClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.ManagementPartner/partners"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2018-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PartnersClient) getHandleResponse(resp *http.Response) (PartnersClientGetResponse, error) {
	result := PartnersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PartnerResponse); err != nil {
		return PartnersClientGetResponse{}, err
	}
	return result, nil
}

