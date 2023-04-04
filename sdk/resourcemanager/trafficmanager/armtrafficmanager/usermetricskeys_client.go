//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armtrafficmanager

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

// UserMetricsKeysClient contains the methods for the TrafficManagerUserMetricsKeys group.
// Don't use this type directly, use NewUserMetricsKeysClient() instead.
type UserMetricsKeysClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewUserMetricsKeysClient creates a new instance of UserMetricsKeysClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewUserMetricsKeysClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*UserMetricsKeysClient, error) {
	cl, err := arm.NewClient(moduleName+".UserMetricsKeysClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &UserMetricsKeysClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update a subscription-level key used for Real User Metrics collection.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01-preview
//   - options - UserMetricsKeysClientCreateOrUpdateOptions contains the optional parameters for the UserMetricsKeysClient.CreateOrUpdate
//     method.
func (client *UserMetricsKeysClient) CreateOrUpdate(ctx context.Context, options *UserMetricsKeysClientCreateOrUpdateOptions) (UserMetricsKeysClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, options)
	if err != nil {
		return UserMetricsKeysClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return UserMetricsKeysClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusCreated) {
		return UserMetricsKeysClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *UserMetricsKeysClient) createOrUpdateCreateRequest(ctx context.Context, options *UserMetricsKeysClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/trafficManagerUserMetricsKeys/default"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *UserMetricsKeysClient) createOrUpdateHandleResponse(resp *http.Response) (UserMetricsKeysClientCreateOrUpdateResponse, error) {
	result := UserMetricsKeysClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UserMetricsModel); err != nil {
		return UserMetricsKeysClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete a subscription-level key used for Real User Metrics collection.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01-preview
//   - options - UserMetricsKeysClientDeleteOptions contains the optional parameters for the UserMetricsKeysClient.Delete method.
func (client *UserMetricsKeysClient) Delete(ctx context.Context, options *UserMetricsKeysClientDeleteOptions) (UserMetricsKeysClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, options)
	if err != nil {
		return UserMetricsKeysClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return UserMetricsKeysClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return UserMetricsKeysClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return client.deleteHandleResponse(resp)
}

// deleteCreateRequest creates the Delete request.
func (client *UserMetricsKeysClient) deleteCreateRequest(ctx context.Context, options *UserMetricsKeysClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/trafficManagerUserMetricsKeys/default"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *UserMetricsKeysClient) deleteHandleResponse(resp *http.Response) (UserMetricsKeysClientDeleteResponse, error) {
	result := UserMetricsKeysClientDeleteResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeleteOperationResult); err != nil {
		return UserMetricsKeysClientDeleteResponse{}, err
	}
	return result, nil
}

// Get - Get the subscription-level key used for Real User Metrics collection.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-04-01-preview
//   - options - UserMetricsKeysClientGetOptions contains the optional parameters for the UserMetricsKeysClient.Get method.
func (client *UserMetricsKeysClient) Get(ctx context.Context, options *UserMetricsKeysClientGetOptions) (UserMetricsKeysClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, options)
	if err != nil {
		return UserMetricsKeysClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return UserMetricsKeysClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return UserMetricsKeysClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *UserMetricsKeysClient) getCreateRequest(ctx context.Context, options *UserMetricsKeysClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/trafficManagerUserMetricsKeys/default"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-04-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *UserMetricsKeysClient) getHandleResponse(resp *http.Response) (UserMetricsKeysClientGetResponse, error) {
	result := UserMetricsKeysClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UserMetricsModel); err != nil {
		return UserMetricsKeysClientGetResponse{}, err
	}
	return result, nil
}
