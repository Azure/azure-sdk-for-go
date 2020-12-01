// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcompute

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

// AvailabilitySetsClient contains the methods for the AvailabilitySets group.
// Don't use this type directly, use NewAvailabilitySetsClient() instead.
type AvailabilitySetsClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewAvailabilitySetsClient creates a new instance of AvailabilitySetsClient with the specified values.
func NewAvailabilitySetsClient(con *armcore.Connection, subscriptionID string) AvailabilitySetsClient {
	return AvailabilitySetsClient{con: con, subscriptionID: subscriptionID}
}

// Pipeline returns the pipeline associated with this client.
func (client AvailabilitySetsClient) Pipeline() azcore.Pipeline {
	return client.con.Pipeline()
}

// CreateOrUpdate - Create or update an availability set.
func (client AvailabilitySetsClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, availabilitySetName string, parameters AvailabilitySet, options *AvailabilitySetsCreateOrUpdateOptions) (AvailabilitySetResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, availabilitySetName, parameters, options)
	if err != nil {
		return AvailabilitySetResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return AvailabilitySetResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return AvailabilitySetResponse{}, client.createOrUpdateHandleError(resp)
	}
	result, err := client.createOrUpdateHandleResponse(resp)
	if err != nil {
		return AvailabilitySetResponse{}, err
	}
	return result, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client AvailabilitySetsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, availabilitySetName string, parameters AvailabilitySet, options *AvailabilitySetsCreateOrUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/availabilitySets/{availabilitySetName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{availabilitySetName}", url.PathEscape(availabilitySetName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPut, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client AvailabilitySetsClient) createOrUpdateHandleResponse(resp *azcore.Response) (AvailabilitySetResponse, error) {
	result := AvailabilitySetResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.AvailabilitySet)
	return result, err
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client AvailabilitySetsClient) createOrUpdateHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// Delete - Delete an availability set.
func (client AvailabilitySetsClient) Delete(ctx context.Context, resourceGroupName string, availabilitySetName string, options *AvailabilitySetsDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, availabilitySetName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	return resp.Response, nil
}

// deleteCreateRequest creates the Delete request.
func (client AvailabilitySetsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, availabilitySetName string, options *AvailabilitySetsDeleteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/availabilitySets/{availabilitySetName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{availabilitySetName}", url.PathEscape(availabilitySetName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodDelete, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client AvailabilitySetsClient) deleteHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// Get - Retrieves information about an availability set.
func (client AvailabilitySetsClient) Get(ctx context.Context, resourceGroupName string, availabilitySetName string, options *AvailabilitySetsGetOptions) (AvailabilitySetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, availabilitySetName, options)
	if err != nil {
		return AvailabilitySetResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return AvailabilitySetResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return AvailabilitySetResponse{}, client.getHandleError(resp)
	}
	result, err := client.getHandleResponse(resp)
	if err != nil {
		return AvailabilitySetResponse{}, err
	}
	return result, nil
}

// getCreateRequest creates the Get request.
func (client AvailabilitySetsClient) getCreateRequest(ctx context.Context, resourceGroupName string, availabilitySetName string, options *AvailabilitySetsGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/availabilitySets/{availabilitySetName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{availabilitySetName}", url.PathEscape(availabilitySetName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client AvailabilitySetsClient) getHandleResponse(resp *azcore.Response) (AvailabilitySetResponse, error) {
	result := AvailabilitySetResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.AvailabilitySet)
	return result, err
}

// getHandleError handles the Get error response.
func (client AvailabilitySetsClient) getHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// List - Lists all availability sets in a resource group.
func (client AvailabilitySetsClient) List(resourceGroupName string, options *AvailabilitySetsListOptions) AvailabilitySetListResultPager {
	return &availabilitySetListResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, options)
		},
		responder: client.listHandleResponse,
		errorer:   client.listHandleError,
		advancer: func(ctx context.Context, resp AvailabilitySetListResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.AvailabilitySetListResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listCreateRequest creates the List request.
func (client AvailabilitySetsClient) listCreateRequest(ctx context.Context, resourceGroupName string, options *AvailabilitySetsListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/availabilitySets"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client AvailabilitySetsClient) listHandleResponse(resp *azcore.Response) (AvailabilitySetListResultResponse, error) {
	result := AvailabilitySetListResultResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.AvailabilitySetListResult)
	return result, err
}

// listHandleError handles the List error response.
func (client AvailabilitySetsClient) listHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// ListAvailableSizes - Lists all available virtual machine sizes that can be used to create a new virtual machine in an existing availability set.
func (client AvailabilitySetsClient) ListAvailableSizes(ctx context.Context, resourceGroupName string, availabilitySetName string, options *AvailabilitySetsListAvailableSizesOptions) (VirtualMachineSizeListResultResponse, error) {
	req, err := client.listAvailableSizesCreateRequest(ctx, resourceGroupName, availabilitySetName, options)
	if err != nil {
		return VirtualMachineSizeListResultResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return VirtualMachineSizeListResultResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return VirtualMachineSizeListResultResponse{}, client.listAvailableSizesHandleError(resp)
	}
	result, err := client.listAvailableSizesHandleResponse(resp)
	if err != nil {
		return VirtualMachineSizeListResultResponse{}, err
	}
	return result, nil
}

// listAvailableSizesCreateRequest creates the ListAvailableSizes request.
func (client AvailabilitySetsClient) listAvailableSizesCreateRequest(ctx context.Context, resourceGroupName string, availabilitySetName string, options *AvailabilitySetsListAvailableSizesOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/availabilitySets/{availabilitySetName}/vmSizes"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{availabilitySetName}", url.PathEscape(availabilitySetName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listAvailableSizesHandleResponse handles the ListAvailableSizes response.
func (client AvailabilitySetsClient) listAvailableSizesHandleResponse(resp *azcore.Response) (VirtualMachineSizeListResultResponse, error) {
	result := VirtualMachineSizeListResultResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.VirtualMachineSizeListResult)
	return result, err
}

// listAvailableSizesHandleError handles the ListAvailableSizes error response.
func (client AvailabilitySetsClient) listAvailableSizesHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// ListBySubscription - Lists all availability sets in a subscription.
func (client AvailabilitySetsClient) ListBySubscription(options *AvailabilitySetsListBySubscriptionOptions) AvailabilitySetListResultPager {
	return &availabilitySetListResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listBySubscriptionCreateRequest(ctx, options)
		},
		responder: client.listBySubscriptionHandleResponse,
		errorer:   client.listBySubscriptionHandleError,
		advancer: func(ctx context.Context, resp AvailabilitySetListResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.AvailabilitySetListResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client AvailabilitySetsClient) listBySubscriptionCreateRequest(ctx context.Context, options *AvailabilitySetsListBySubscriptionOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Compute/availabilitySets"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	if options != nil && options.Expand != nil {
		query.Set("$expand", *options.Expand)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client AvailabilitySetsClient) listBySubscriptionHandleResponse(resp *azcore.Response) (AvailabilitySetListResultResponse, error) {
	result := AvailabilitySetListResultResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.AvailabilitySetListResult)
	return result, err
}

// listBySubscriptionHandleError handles the ListBySubscription error response.
func (client AvailabilitySetsClient) listBySubscriptionHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}

// Update - Update an availability set.
func (client AvailabilitySetsClient) Update(ctx context.Context, resourceGroupName string, availabilitySetName string, parameters AvailabilitySetUpdate, options *AvailabilitySetsUpdateOptions) (AvailabilitySetResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, availabilitySetName, parameters, options)
	if err != nil {
		return AvailabilitySetResponse{}, err
	}
	resp, err := client.Pipeline().Do(req)
	if err != nil {
		return AvailabilitySetResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return AvailabilitySetResponse{}, client.updateHandleError(resp)
	}
	result, err := client.updateHandleResponse(resp)
	if err != nil {
		return AvailabilitySetResponse{}, err
	}
	return result, nil
}

// updateCreateRequest creates the Update request.
func (client AvailabilitySetsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, availabilitySetName string, parameters AvailabilitySetUpdate, options *AvailabilitySetsUpdateOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Compute/availabilitySets/{availabilitySetName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{availabilitySetName}", url.PathEscape(availabilitySetName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := azcore.NewRequest(ctx, http.MethodPatch, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	query := req.URL.Query()
	query.Set("api-version", "2020-06-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, req.MarshalAsJSON(parameters)
}

// updateHandleResponse handles the Update response.
func (client AvailabilitySetsClient) updateHandleResponse(resp *azcore.Response) (AvailabilitySetResponse, error) {
	result := AvailabilitySetResponse{RawResponse: resp.Response}
	err := resp.UnmarshalAsJSON(&result.AvailabilitySet)
	return result, err
}

// updateHandleError handles the Update error response.
func (client AvailabilitySetsClient) updateHandleError(resp *azcore.Response) error {
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("%s; failed to read response body: %w", resp.Status, err)
	}
	if len(body) == 0 {
		return azcore.NewResponseError(errors.New(resp.Status), resp.Response)
	}
	return azcore.NewResponseError(errors.New(string(body)), resp.Response)
}
