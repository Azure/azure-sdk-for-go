//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armcostmanagement

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// AlertsClient contains the methods for the Alerts group.
// Don't use this type directly, use NewAlertsClient() instead.
type AlertsClient struct {
	ep string
	pl runtime.Pipeline
}

// NewAlertsClient creates a new instance of AlertsClient with the specified values.
func NewAlertsClient(credential azcore.TokenCredential, options *arm.ClientOptions) *AlertsClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &AlertsClient{ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// Dismiss - Dismisses the specified alert
// If the operation fails it returns the *ErrorResponse error type.
func (client *AlertsClient) Dismiss(ctx context.Context, scope string, alertID string, parameters DismissAlertPayload, options *AlertsDismissOptions) (AlertsDismissResponse, error) {
	req, err := client.dismissCreateRequest(ctx, scope, alertID, parameters, options)
	if err != nil {
		return AlertsDismissResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AlertsDismissResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AlertsDismissResponse{}, client.dismissHandleError(resp)
	}
	return client.dismissHandleResponse(resp)
}

// dismissCreateRequest creates the Dismiss request.
func (client *AlertsClient) dismissCreateRequest(ctx context.Context, scope string, alertID string, parameters DismissAlertPayload, options *AlertsDismissOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.CostManagement/alerts/{alertId}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	urlPath = strings.ReplaceAll(urlPath, "{alertId}", alertID)
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// dismissHandleResponse handles the Dismiss response.
func (client *AlertsClient) dismissHandleResponse(resp *http.Response) (AlertsDismissResponse, error) {
	result := AlertsDismissResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.Alert); err != nil {
		return AlertsDismissResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// dismissHandleError handles the Dismiss error response.
func (client *AlertsClient) dismissHandleError(resp *http.Response) error {
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

// Get - Gets the alert for the scope by alert ID.
// If the operation fails it returns the *ErrorResponse error type.
func (client *AlertsClient) Get(ctx context.Context, scope string, alertID string, options *AlertsGetOptions) (AlertsGetResponse, error) {
	req, err := client.getCreateRequest(ctx, scope, alertID, options)
	if err != nil {
		return AlertsGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AlertsGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AlertsGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *AlertsClient) getCreateRequest(ctx context.Context, scope string, alertID string, options *AlertsGetOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.CostManagement/alerts/{alertId}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	urlPath = strings.ReplaceAll(urlPath, "{alertId}", alertID)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AlertsClient) getHandleResponse(resp *http.Response) (AlertsGetResponse, error) {
	result := AlertsGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.Alert); err != nil {
		return AlertsGetResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *AlertsClient) getHandleError(resp *http.Response) error {
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

// List - Lists the alerts for scope defined.
// If the operation fails it returns the *ErrorResponse error type.
func (client *AlertsClient) List(ctx context.Context, scope string, options *AlertsListOptions) (AlertsListResponse, error) {
	req, err := client.listCreateRequest(ctx, scope, options)
	if err != nil {
		return AlertsListResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AlertsListResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AlertsListResponse{}, client.listHandleError(resp)
	}
	return client.listHandleResponse(resp)
}

// listCreateRequest creates the List request.
func (client *AlertsClient) listCreateRequest(ctx context.Context, scope string, options *AlertsListOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.CostManagement/alerts"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *AlertsClient) listHandleResponse(resp *http.Response) (AlertsListResponse, error) {
	result := AlertsListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.AlertsResult); err != nil {
		return AlertsListResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *AlertsClient) listHandleError(resp *http.Response) error {
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

// ListExternal - Lists the Alerts for external cloud provider type defined.
// If the operation fails it returns the *ErrorResponse error type.
func (client *AlertsClient) ListExternal(ctx context.Context, externalCloudProviderType ExternalCloudProviderType, externalCloudProviderID string, options *AlertsListExternalOptions) (AlertsListExternalResponse, error) {
	req, err := client.listExternalCreateRequest(ctx, externalCloudProviderType, externalCloudProviderID, options)
	if err != nil {
		return AlertsListExternalResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AlertsListExternalResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AlertsListExternalResponse{}, client.listExternalHandleError(resp)
	}
	return client.listExternalHandleResponse(resp)
}

// listExternalCreateRequest creates the ListExternal request.
func (client *AlertsClient) listExternalCreateRequest(ctx context.Context, externalCloudProviderType ExternalCloudProviderType, externalCloudProviderID string, options *AlertsListExternalOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.CostManagement/{externalCloudProviderType}/{externalCloudProviderId}/alerts"
	if externalCloudProviderType == "" {
		return nil, errors.New("parameter externalCloudProviderType cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{externalCloudProviderType}", url.PathEscape(string(externalCloudProviderType)))
	if externalCloudProviderID == "" {
		return nil, errors.New("parameter externalCloudProviderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{externalCloudProviderId}", url.PathEscape(externalCloudProviderID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listExternalHandleResponse handles the ListExternal response.
func (client *AlertsClient) listExternalHandleResponse(resp *http.Response) (AlertsListExternalResponse, error) {
	result := AlertsListExternalResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.AlertsResult); err != nil {
		return AlertsListExternalResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// listExternalHandleError handles the ListExternal error response.
func (client *AlertsClient) listExternalHandleError(resp *http.Response) error {
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
