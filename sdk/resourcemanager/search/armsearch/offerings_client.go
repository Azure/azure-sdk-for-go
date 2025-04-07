// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsearch

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
)

// OfferingsClient contains the methods for the Offerings group.
// Don't use this type directly, use NewOfferingsClient() instead.
type OfferingsClient struct {
	internal *arm.Client
}

// NewOfferingsClient creates a new instance of OfferingsClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewOfferingsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*OfferingsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &OfferingsClient{
		internal: cl,
	}
	return client, nil
}

// NewListPager - Lists all of the features and SKUs offered by the Azure AI Search service in each region.
//
// Generated from API version 2025-02-01-preview
//   - options - OfferingsClientListOptions contains the optional parameters for the OfferingsClient.NewListPager method.
func (client *OfferingsClient) NewListPager(options *OfferingsClientListOptions) *runtime.Pager[OfferingsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[OfferingsClientListResponse]{
		More: func(page OfferingsClientListResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *OfferingsClientListResponse) (OfferingsClientListResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "OfferingsClient.NewListPager")
			req, err := client.listCreateRequest(ctx, options)
			if err != nil {
				return OfferingsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return OfferingsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return OfferingsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *OfferingsClient) listCreateRequest(ctx context.Context, _ *OfferingsClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Search/offerings"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2025-02-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *OfferingsClient) listHandleResponse(resp *http.Response) (OfferingsClientListResponse, error) {
	result := OfferingsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.OfferingsListResult); err != nil {
		return OfferingsClientListResponse{}, err
	}
	return result, nil
}
