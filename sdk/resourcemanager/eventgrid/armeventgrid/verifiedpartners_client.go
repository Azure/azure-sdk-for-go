//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armeventgrid

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// VerifiedPartnersClient contains the methods for the VerifiedPartners group.
// Don't use this type directly, use NewVerifiedPartnersClient() instead.
type VerifiedPartnersClient struct {
	internal *arm.Client
}

// NewVerifiedPartnersClient creates a new instance of VerifiedPartnersClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewVerifiedPartnersClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*VerifiedPartnersClient, error) {
	cl, err := arm.NewClient(moduleName+".VerifiedPartnersClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &VerifiedPartnersClient{
		internal: cl,
	}
	return client, nil
}

// Get - Get properties of a verified partner.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-06-15
//   - verifiedPartnerName - Name of the verified partner.
//   - options - VerifiedPartnersClientGetOptions contains the optional parameters for the VerifiedPartnersClient.Get method.
func (client *VerifiedPartnersClient) Get(ctx context.Context, verifiedPartnerName string, options *VerifiedPartnersClientGetOptions) (VerifiedPartnersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, verifiedPartnerName, options)
	if err != nil {
		return VerifiedPartnersClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return VerifiedPartnersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return VerifiedPartnersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VerifiedPartnersClient) getCreateRequest(ctx context.Context, verifiedPartnerName string, options *VerifiedPartnersClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.EventGrid/verifiedPartners/{verifiedPartnerName}"
	if verifiedPartnerName == "" {
		return nil, errors.New("parameter verifiedPartnerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{verifiedPartnerName}", url.PathEscape(verifiedPartnerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-06-15")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VerifiedPartnersClient) getHandleResponse(resp *http.Response) (VerifiedPartnersClientGetResponse, error) {
	result := VerifiedPartnersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VerifiedPartner); err != nil {
		return VerifiedPartnersClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Get a list of all verified partners.
//
// Generated from API version 2022-06-15
//   - options - VerifiedPartnersClientListOptions contains the optional parameters for the VerifiedPartnersClient.NewListPager
//     method.
func (client *VerifiedPartnersClient) NewListPager(options *VerifiedPartnersClientListOptions) *runtime.Pager[VerifiedPartnersClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[VerifiedPartnersClientListResponse]{
		More: func(page VerifiedPartnersClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *VerifiedPartnersClientListResponse) (VerifiedPartnersClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return VerifiedPartnersClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return VerifiedPartnersClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return VerifiedPartnersClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *VerifiedPartnersClient) listCreateRequest(ctx context.Context, options *VerifiedPartnersClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.EventGrid/verifiedPartners"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-06-15")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *VerifiedPartnersClient) listHandleResponse(resp *http.Response) (VerifiedPartnersClientListResponse, error) {
	result := VerifiedPartnersClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.VerifiedPartnersListResult); err != nil {
		return VerifiedPartnersClientListResponse{}, err
	}
	return result, nil
}
