//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armresourcehealth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ChildAvailabilityStatusesClient contains the methods for the ChildAvailabilityStatuses group.
// Don't use this type directly, use NewChildAvailabilityStatusesClient() instead.
type ChildAvailabilityStatusesClient struct {
	ep string
	pl runtime.Pipeline
}

// NewChildAvailabilityStatusesClient creates a new instance of ChildAvailabilityStatusesClient with the specified values.
func NewChildAvailabilityStatusesClient(con *arm.Connection) *ChildAvailabilityStatusesClient {
	return &ChildAvailabilityStatusesClient{ep: con.Endpoint(), pl: con.NewPipeline(module, version)}
}

// GetByResource - Gets current availability status for a single resource
// If the operation fails it returns the *ErrorResponse error type.
func (client *ChildAvailabilityStatusesClient) GetByResource(ctx context.Context, resourceURI string, options *ChildAvailabilityStatusesGetByResourceOptions) (ChildAvailabilityStatusesGetByResourceResponse, error) {
	req, err := client.getByResourceCreateRequest(ctx, resourceURI, options)
	if err != nil {
		return ChildAvailabilityStatusesGetByResourceResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ChildAvailabilityStatusesGetByResourceResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ChildAvailabilityStatusesGetByResourceResponse{}, client.getByResourceHandleError(resp)
	}
	return client.getByResourceHandleResponse(resp)
}

// getByResourceCreateRequest creates the GetByResource request.
func (client *ChildAvailabilityStatusesClient) getByResourceCreateRequest(ctx context.Context, resourceURI string, options *ChildAvailabilityStatusesGetByResourceOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.ResourceHealth/childAvailabilityStatuses/current"
	if resourceURI == "" {
		return nil, errors.New("parameter resourceURI cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2017-07-01")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getByResourceHandleResponse handles the GetByResource response.
func (client *ChildAvailabilityStatusesClient) getByResourceHandleResponse(resp *http.Response) (ChildAvailabilityStatusesGetByResourceResponse, error) {
	result := ChildAvailabilityStatusesGetByResourceResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.AvailabilityStatus); err != nil {
		return ChildAvailabilityStatusesGetByResourceResponse{}, err
	}
	return result, nil
}

// getByResourceHandleError handles the GetByResource error response.
func (client *ChildAvailabilityStatusesClient) getByResourceHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// List - Lists the historical availability statuses for a single child resource. Use the nextLink property in the response to get the next page of availability
// status
// If the operation fails it returns the *ErrorResponse error type.
func (client *ChildAvailabilityStatusesClient) List(resourceURI string, options *ChildAvailabilityStatusesListOptions) *ChildAvailabilityStatusesListPager {
	return &ChildAvailabilityStatusesListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceURI, options)
		},
		advancer: func(ctx context.Context, resp ChildAvailabilityStatusesListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.AvailabilityStatusListResult.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *ChildAvailabilityStatusesClient) listCreateRequest(ctx context.Context, resourceURI string, options *ChildAvailabilityStatusesListOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.ResourceHealth/childAvailabilityStatuses"
	if resourceURI == "" {
		return nil, errors.New("parameter resourceURI cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2017-07-01")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ChildAvailabilityStatusesClient) listHandleResponse(resp *http.Response) (ChildAvailabilityStatusesListResponse, error) {
	result := ChildAvailabilityStatusesListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.AvailabilityStatusListResult); err != nil {
		return ChildAvailabilityStatusesListResponse{}, err
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *ChildAvailabilityStatusesClient) listHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
