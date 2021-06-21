// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armauthorization

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
)

// RoleManagementPoliciesClient contains the methods for the RoleManagementPolicies group.
// Don't use this type directly, use NewRoleManagementPoliciesClient() instead.
type RoleManagementPoliciesClient struct {
	con *armcore.Connection
}

// NewRoleManagementPoliciesClient creates a new instance of RoleManagementPoliciesClient with the specified values.
func NewRoleManagementPoliciesClient(con *armcore.Connection) *RoleManagementPoliciesClient {
	return &RoleManagementPoliciesClient{con: con}
}

// Delete - Delete a role management policy
// If the operation fails it returns the *CloudError error type.
func (client *RoleManagementPoliciesClient) Delete(ctx context.Context, scope string, roleManagementPolicyName string, options *RoleManagementPoliciesDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, scope, roleManagementPolicyName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	return resp.Response, nil
}

// deleteCreateRequest creates the Delete request.
func (client *RoleManagementPoliciesClient) deleteCreateRequest(ctx context.Context, scope string, roleManagementPolicyName string, options *RoleManagementPoliciesDeleteOptions) (*azcore.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/roleManagementPolicies/{roleManagementPolicyName}"
	if scope == "" {
		return nil, errors.New("parameter scope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if roleManagementPolicyName == "" {
		return nil, errors.New("parameter roleManagementPolicyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{roleManagementPolicyName}", url.PathEscape(roleManagementPolicyName))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2020-10-01-preview")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *RoleManagementPoliciesClient) deleteHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
		errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}

// Get - Get the specified role management policy for a resource scope
// If the operation fails it returns the *CloudError error type.
func (client *RoleManagementPoliciesClient) Get(ctx context.Context, scope string, roleManagementPolicyName string, options *RoleManagementPoliciesGetOptions) (RoleManagementPolicyResponse, error) {
	req, err := client.getCreateRequest(ctx, scope, roleManagementPolicyName, options)
	if err != nil {
		return RoleManagementPolicyResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return RoleManagementPolicyResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return RoleManagementPolicyResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *RoleManagementPoliciesClient) getCreateRequest(ctx context.Context, scope string, roleManagementPolicyName string, options *RoleManagementPoliciesGetOptions) (*azcore.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/roleManagementPolicies/{roleManagementPolicyName}"
	if scope == "" {
		return nil, errors.New("parameter scope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if roleManagementPolicyName == "" {
		return nil, errors.New("parameter roleManagementPolicyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{roleManagementPolicyName}", url.PathEscape(roleManagementPolicyName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2020-10-01-preview")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *RoleManagementPoliciesClient) getHandleResponse(resp *azcore.Response) (RoleManagementPolicyResponse, error) {
	var val *RoleManagementPolicy
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return RoleManagementPolicyResponse{}, err
	}
return RoleManagementPolicyResponse{RawResponse: resp.Response, RoleManagementPolicy: val}, nil
}

// getHandleError handles the Get error response.
func (client *RoleManagementPoliciesClient) getHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
		errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}

// ListForScope - Gets role management policies for a resource scope.
// If the operation fails it returns the *CloudError error type.
func (client *RoleManagementPoliciesClient) ListForScope(scope string, options *RoleManagementPoliciesListForScopeOptions) (RoleManagementPolicyListResultPager) {
	return &roleManagementPolicyListResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listForScopeCreateRequest(ctx, scope, options)
		},
		responder: client.listForScopeHandleResponse,
		errorer:   client.listForScopeHandleError,
		advancer: func(ctx context.Context, resp RoleManagementPolicyListResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.RoleManagementPolicyListResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listForScopeCreateRequest creates the ListForScope request.
func (client *RoleManagementPoliciesClient) listForScopeCreateRequest(ctx context.Context, scope string, options *RoleManagementPoliciesListForScopeOptions) (*azcore.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/roleManagementPolicies"
	if scope == "" {
		return nil, errors.New("parameter scope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2020-10-01-preview")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listForScopeHandleResponse handles the ListForScope response.
func (client *RoleManagementPoliciesClient) listForScopeHandleResponse(resp *azcore.Response) (RoleManagementPolicyListResultResponse, error) {
	var val *RoleManagementPolicyListResult
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return RoleManagementPolicyListResultResponse{}, err
	}
return RoleManagementPolicyListResultResponse{RawResponse: resp.Response, RoleManagementPolicyListResult: val}, nil
}

// listForScopeHandleError handles the ListForScope error response.
func (client *RoleManagementPoliciesClient) listForScopeHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
		errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}

// Update - Update a role management policy
// If the operation fails it returns the *CloudError error type.
func (client *RoleManagementPoliciesClient) Update(ctx context.Context, scope string, roleManagementPolicyName string, parameters RoleManagementPolicy, options *RoleManagementPoliciesUpdateOptions) (RoleManagementPolicyResponse, error) {
	req, err := client.updateCreateRequest(ctx, scope, roleManagementPolicyName, parameters, options)
	if err != nil {
		return RoleManagementPolicyResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return RoleManagementPolicyResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return RoleManagementPolicyResponse{}, client.updateHandleError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *RoleManagementPoliciesClient) updateCreateRequest(ctx context.Context, scope string, roleManagementPolicyName string, parameters RoleManagementPolicy, options *RoleManagementPoliciesUpdateOptions) (*azcore.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/roleManagementPolicies/{roleManagementPolicyName}"
	if scope == "" {
		return nil, errors.New("parameter scope cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	if roleManagementPolicyName == "" {
		return nil, errors.New("parameter roleManagementPolicyName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{roleManagementPolicyName}", url.PathEscape(roleManagementPolicyName))
	req, err := azcore.NewRequest(ctx, http.MethodPatch, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2020-10-01-preview")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// updateHandleResponse handles the Update response.
func (client *RoleManagementPoliciesClient) updateHandleResponse(resp *azcore.Response) (RoleManagementPolicyResponse, error) {
	var val *RoleManagementPolicy
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return RoleManagementPolicyResponse{}, err
	}
return RoleManagementPolicyResponse{RawResponse: resp.Response, RoleManagementPolicy: val}, nil
}

// updateHandleError handles the Update error response.
func (client *RoleManagementPoliciesClient) updateHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
		errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}

