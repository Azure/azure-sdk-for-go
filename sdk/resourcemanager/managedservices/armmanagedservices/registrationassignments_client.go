//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armmanagedservices

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

// RegistrationAssignmentsClient contains the methods for the RegistrationAssignments group.
// Don't use this type directly, use NewRegistrationAssignmentsClient() instead.
type RegistrationAssignmentsClient struct {
	internal *arm.Client
}

// NewRegistrationAssignmentsClient creates a new instance of RegistrationAssignmentsClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewRegistrationAssignmentsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*RegistrationAssignmentsClient, error) {
	cl, err := arm.NewClient(moduleName+".RegistrationAssignmentsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &RegistrationAssignmentsClient{
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Creates or updates a registration assignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
//   - scope - The scope of the resource.
//   - registrationAssignmentID - The GUID of the registration assignment.
//   - requestBody - The parameters required to create new registration assignment.
//   - options - RegistrationAssignmentsClientBeginCreateOrUpdateOptions contains the optional parameters for the RegistrationAssignmentsClient.BeginCreateOrUpdate
//     method.
func (client *RegistrationAssignmentsClient) BeginCreateOrUpdate(ctx context.Context, scope string, registrationAssignmentID string, requestBody RegistrationAssignment, options *RegistrationAssignmentsClientBeginCreateOrUpdateOptions) (*runtime.Poller[RegistrationAssignmentsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, scope, registrationAssignmentID, requestBody, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[RegistrationAssignmentsClientCreateOrUpdateResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[RegistrationAssignmentsClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Creates or updates a registration assignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
func (client *RegistrationAssignmentsClient) createOrUpdate(ctx context.Context, scope string, registrationAssignmentID string, requestBody RegistrationAssignment, options *RegistrationAssignmentsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, scope, registrationAssignmentID, requestBody, options)
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

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *RegistrationAssignmentsClient) createOrUpdateCreateRequest(ctx context.Context, scope string, registrationAssignmentID string, requestBody RegistrationAssignment, options *RegistrationAssignmentsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.ManagedServices/registrationAssignments/{registrationAssignmentId}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if registrationAssignmentID == "" {
		return nil, errors.New("parameter registrationAssignmentID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registrationAssignmentId}", url.PathEscape(registrationAssignmentID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, requestBody); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Deletes the specified registration assignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
//   - scope - The scope of the resource.
//   - registrationAssignmentID - The GUID of the registration assignment.
//   - options - RegistrationAssignmentsClientBeginDeleteOptions contains the optional parameters for the RegistrationAssignmentsClient.BeginDelete
//     method.
func (client *RegistrationAssignmentsClient) BeginDelete(ctx context.Context, scope string, registrationAssignmentID string, options *RegistrationAssignmentsClientBeginDeleteOptions) (*runtime.Poller[RegistrationAssignmentsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, scope, registrationAssignmentID, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[RegistrationAssignmentsClientDeleteResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[RegistrationAssignmentsClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Deletes the specified registration assignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
func (client *RegistrationAssignmentsClient) deleteOperation(ctx context.Context, scope string, registrationAssignmentID string, options *RegistrationAssignmentsClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, scope, registrationAssignmentID, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *RegistrationAssignmentsClient) deleteCreateRequest(ctx context.Context, scope string, registrationAssignmentID string, options *RegistrationAssignmentsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.ManagedServices/registrationAssignments/{registrationAssignmentId}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if registrationAssignmentID == "" {
		return nil, errors.New("parameter registrationAssignmentID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registrationAssignmentId}", url.PathEscape(registrationAssignmentID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets the details of the specified registration assignment.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-01-preview
//   - scope - The scope of the resource.
//   - registrationAssignmentID - The GUID of the registration assignment.
//   - options - RegistrationAssignmentsClientGetOptions contains the optional parameters for the RegistrationAssignmentsClient.Get
//     method.
func (client *RegistrationAssignmentsClient) Get(ctx context.Context, scope string, registrationAssignmentID string, options *RegistrationAssignmentsClientGetOptions) (RegistrationAssignmentsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, scope, registrationAssignmentID, options)
	if err != nil {
		return RegistrationAssignmentsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RegistrationAssignmentsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RegistrationAssignmentsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *RegistrationAssignmentsClient) getCreateRequest(ctx context.Context, scope string, registrationAssignmentID string, options *RegistrationAssignmentsClientGetOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.ManagedServices/registrationAssignments/{registrationAssignmentId}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if registrationAssignmentID == "" {
		return nil, errors.New("parameter registrationAssignmentID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{registrationAssignmentId}", url.PathEscape(registrationAssignmentID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.ExpandRegistrationDefinition != nil {
		reqQP.Set("$expandRegistrationDefinition", strconv.FormatBool(*options.ExpandRegistrationDefinition))
	}
	reqQP.Set("api-version", "2022-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *RegistrationAssignmentsClient) getHandleResponse(resp *http.Response) (RegistrationAssignmentsClientGetResponse, error) {
	result := RegistrationAssignmentsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RegistrationAssignment); err != nil {
		return RegistrationAssignmentsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets a list of the registration assignments.
//
// Generated from API version 2022-01-01-preview
//   - scope - The scope of the resource.
//   - options - RegistrationAssignmentsClientListOptions contains the optional parameters for the RegistrationAssignmentsClient.NewListPager
//     method.
func (client *RegistrationAssignmentsClient) NewListPager(scope string, options *RegistrationAssignmentsClientListOptions) (*runtime.Pager[RegistrationAssignmentsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[RegistrationAssignmentsClientListResponse]{
		More: func(page RegistrationAssignmentsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *RegistrationAssignmentsClientListResponse) (RegistrationAssignmentsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, scope, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return RegistrationAssignmentsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return RegistrationAssignmentsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return RegistrationAssignmentsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *RegistrationAssignmentsClient) listCreateRequest(ctx context.Context, scope string, options *RegistrationAssignmentsClientListOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.ManagedServices/registrationAssignments"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.ExpandRegistrationDefinition != nil {
		reqQP.Set("$expandRegistrationDefinition", strconv.FormatBool(*options.ExpandRegistrationDefinition))
	}
	reqQP.Set("api-version", "2022-01-01-preview")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *RegistrationAssignmentsClient) listHandleResponse(resp *http.Response) (RegistrationAssignmentsClientListResponse, error) {
	result := RegistrationAssignmentsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RegistrationAssignmentList); err != nil {
		return RegistrationAssignmentsClientListResponse{}, err
	}
	return result, nil
}

