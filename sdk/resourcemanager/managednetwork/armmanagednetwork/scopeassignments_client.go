//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmanagednetwork

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

// ScopeAssignmentsClient contains the methods for the ScopeAssignments group.
// Don't use this type directly, use NewScopeAssignmentsClient() instead.
type ScopeAssignmentsClient struct {
	internal *arm.Client
}

// NewScopeAssignmentsClient creates a new instance of ScopeAssignmentsClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewScopeAssignmentsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ScopeAssignmentsClient, error) {
	cl, err := arm.NewClient(moduleName+".ScopeAssignmentsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ScopeAssignmentsClient{
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates a scope assignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-06-01-preview
//   - scope - The base resource of the scope assignment to create. The scope can be any REST resource instance. For example,
//     use 'subscriptions/{subscription-id}' for a subscription,
//     'subscriptions/{subscription-id}/resourceGroups/{resource-group-name}' for a resource group, and
//     'subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/{resource-provider}/{resource-type}/{resource-name}'
//     for a resource.
//   - scopeAssignmentName - The name of the scope assignment to create.
//   - parameters - Parameters supplied to the specify which Managed Network this scope is being assigned
//   - options - ScopeAssignmentsClientCreateOrUpdateOptions contains the optional parameters for the ScopeAssignmentsClient.CreateOrUpdate
//     method.
func (client *ScopeAssignmentsClient) CreateOrUpdate(ctx context.Context, scope string, scopeAssignmentName string, parameters ScopeAssignment, options *ScopeAssignmentsClientCreateOrUpdateOptions) (ScopeAssignmentsClientCreateOrUpdateResponse, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, scope, scopeAssignmentName, parameters, options)
	if err != nil {
		return ScopeAssignmentsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ScopeAssignmentsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return ScopeAssignmentsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ScopeAssignmentsClient) createOrUpdateCreateRequest(ctx context.Context, scope string, scopeAssignmentName string, parameters ScopeAssignment, options *ScopeAssignmentsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.ManagedNetwork/scopeAssignments/{scopeAssignmentName}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if scopeAssignmentName == "" {
		return nil, errors.New("parameter scopeAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scopeAssignmentName}", url.PathEscape(scopeAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, parameters); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ScopeAssignmentsClient) createOrUpdateHandleResponse(resp *http.Response) (ScopeAssignmentsClientCreateOrUpdateResponse, error) {
	result := ScopeAssignmentsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ScopeAssignment); err != nil {
		return ScopeAssignmentsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a scope assignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-06-01-preview
//   - scope - The scope of the scope assignment to delete.
//   - scopeAssignmentName - The name of the scope assignment to delete.
//   - options - ScopeAssignmentsClientDeleteOptions contains the optional parameters for the ScopeAssignmentsClient.Delete method.
func (client *ScopeAssignmentsClient) Delete(ctx context.Context, scope string, scopeAssignmentName string, options *ScopeAssignmentsClientDeleteOptions) (ScopeAssignmentsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, scope, scopeAssignmentName, options)
	if err != nil {
		return ScopeAssignmentsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ScopeAssignmentsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ScopeAssignmentsClientDeleteResponse{}, err
	}
	return ScopeAssignmentsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ScopeAssignmentsClient) deleteCreateRequest(ctx context.Context, scope string, scopeAssignmentName string, options *ScopeAssignmentsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.ManagedNetwork/scopeAssignments/{scopeAssignmentName}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if scopeAssignmentName == "" {
		return nil, errors.New("parameter scopeAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scopeAssignmentName}", url.PathEscape(scopeAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	return req, nil
}

// Get - Get the specified scope assignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-06-01-preview
//   - scope - The base resource of the scope assignment.
//   - scopeAssignmentName - The name of the scope assignment to get.
//   - options - ScopeAssignmentsClientGetOptions contains the optional parameters for the ScopeAssignmentsClient.Get method.
func (client *ScopeAssignmentsClient) Get(ctx context.Context, scope string, scopeAssignmentName string, options *ScopeAssignmentsClientGetOptions) (ScopeAssignmentsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, scope, scopeAssignmentName, options)
	if err != nil {
		return ScopeAssignmentsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ScopeAssignmentsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ScopeAssignmentsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ScopeAssignmentsClient) getCreateRequest(ctx context.Context, scope string, scopeAssignmentName string, options *ScopeAssignmentsClientGetOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.ManagedNetwork/scopeAssignments/{scopeAssignmentName}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if scopeAssignmentName == "" {
		return nil, errors.New("parameter scopeAssignmentName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scopeAssignmentName}", url.PathEscape(scopeAssignmentName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ScopeAssignmentsClient) getHandleResponse(resp *http.Response) (ScopeAssignmentsClientGetResponse, error) {
	result := ScopeAssignmentsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ScopeAssignment); err != nil {
		return ScopeAssignmentsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Get the specified scope assignment.
//
// Generated from API version 2019-06-01-preview
//   - scope - The base resource of the scope assignment.
//   - options - ScopeAssignmentsClientListOptions contains the optional parameters for the ScopeAssignmentsClient.NewListPager
//     method.
func (client *ScopeAssignmentsClient) NewListPager(scope string, options *ScopeAssignmentsClientListOptions) (*runtime.Pager[ScopeAssignmentsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ScopeAssignmentsClientListResponse]{
		More: func(page ScopeAssignmentsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ScopeAssignmentsClientListResponse) (ScopeAssignmentsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, scope, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ScopeAssignmentsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ScopeAssignmentsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ScopeAssignmentsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ScopeAssignmentsClient) listCreateRequest(ctx context.Context, scope string, options *ScopeAssignmentsClientListOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.ManagedNetwork/scopeAssignments"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-06-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ScopeAssignmentsClient) listHandleResponse(resp *http.Response) (ScopeAssignmentsClientListResponse, error) {
	result := ScopeAssignmentsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ScopeAssignmentListResult); err != nil {
		return ScopeAssignmentsClientListResponse{}, err
	}
	return result, nil
}

