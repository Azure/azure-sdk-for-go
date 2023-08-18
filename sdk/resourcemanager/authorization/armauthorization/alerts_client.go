//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armauthorization

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"strings"
)

// AlertsClient contains the methods for the Alerts group.
// Don't use this type directly, use NewAlertsClient() instead.
type AlertsClient struct {
	internal *arm.Client
}

// NewAlertsClient creates a new instance of AlertsClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAlertsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*AlertsClient, error) {
	cl, err := arm.NewClient(moduleName+".AlertsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AlertsClient{
		internal: cl,
	}
	return client, nil
}

// Get - Get the specified alert.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-01-preview
//   - scope - The scope of the alert. The scope can be any REST resource instance. For example, use '/providers/Microsoft.Subscription/subscriptions/{subscription-id}/'
//     for a subscription,
//     '/providers/Microsoft.Subscription/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}' for a resource
//     group, and
//     '/providers/Microsoft.Subscription/subscriptions/{subscription-id}/resourceGroups/{resource-group-name}/providers/{resource-provider}/{resource-type}/{resource-name}'
//     for a resource.
//   - alertID - The name of the alert to get.
//   - options - AlertsClientGetOptions contains the optional parameters for the AlertsClient.Get method.
func (client *AlertsClient) Get(ctx context.Context, scope string, alertID string, options *AlertsClientGetOptions) (AlertsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, scope, alertID, options)
	if err != nil {
		return AlertsClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AlertsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AlertsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *AlertsClient) getCreateRequest(ctx context.Context, scope string, alertID string, options *AlertsClientGetOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/roleManagementAlerts/{alertId}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	urlPath = strings.ReplaceAll(urlPath, "{alertId}", alertID)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AlertsClient) getHandleResponse(resp *http.Response) (AlertsClientGetResponse, error) {
	result := AlertsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Alert); err != nil {
		return AlertsClientGetResponse{}, err
	}
	return result, nil
}

// NewListForScopePager - Gets alerts for a resource scope.
//
// Generated from API version 2022-08-01-preview
//   - scope - The scope of the alert.
//   - options - AlertsClientListForScopeOptions contains the optional parameters for the AlertsClient.NewListForScopePager method.
func (client *AlertsClient) NewListForScopePager(scope string, options *AlertsClientListForScopeOptions) *runtime.Pager[AlertsClientListForScopeResponse] {
	return runtime.NewPager(runtime.PagingHandler[AlertsClientListForScopeResponse]{
		More: func(page AlertsClientListForScopeResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AlertsClientListForScopeResponse) (AlertsClientListForScopeResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listForScopeCreateRequest(ctx, scope, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return AlertsClientListForScopeResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return AlertsClientListForScopeResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return AlertsClientListForScopeResponse{}, runtime.NewResponseError(resp)
			}
			return client.listForScopeHandleResponse(resp)
		},
	})
}

// listForScopeCreateRequest creates the ListForScope request.
func (client *AlertsClient) listForScopeCreateRequest(ctx context.Context, scope string, options *AlertsClientListForScopeOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/roleManagementAlerts"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listForScopeHandleResponse handles the ListForScope response.
func (client *AlertsClient) listForScopeHandleResponse(resp *http.Response) (AlertsClientListForScopeResponse, error) {
	result := AlertsClientListForScopeResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AlertListResult); err != nil {
		return AlertsClientListForScopeResponse{}, err
	}
	return result, nil
}

// BeginRefresh - Refresh an alert.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-01-preview
//   - scope - The scope of the alert.
//   - alertID - The name of the alert to refresh.
//   - options - AlertsClientBeginRefreshOptions contains the optional parameters for the AlertsClient.BeginRefresh method.
func (client *AlertsClient) BeginRefresh(ctx context.Context, scope string, alertID string, options *AlertsClientBeginRefreshOptions) (*runtime.Poller[AlertsClientRefreshResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.refresh(ctx, scope, alertID, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AlertsClientRefreshResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[AlertsClientRefreshResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Refresh - Refresh an alert.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-01-preview
func (client *AlertsClient) refresh(ctx context.Context, scope string, alertID string, options *AlertsClientBeginRefreshOptions) (*http.Response, error) {
	req, err := client.refreshCreateRequest(ctx, scope, alertID, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// refreshCreateRequest creates the Refresh request.
func (client *AlertsClient) refreshCreateRequest(ctx context.Context, scope string, alertID string, options *AlertsClientBeginRefreshOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/roleManagementAlerts/{alertId}/refresh"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	urlPath = strings.ReplaceAll(urlPath, "{alertId}", alertID)
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginRefreshAll - Refresh all alerts for a resource scope.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-01-preview
//   - scope - The scope of the alert.
//   - options - AlertsClientBeginRefreshAllOptions contains the optional parameters for the AlertsClient.BeginRefreshAll method.
func (client *AlertsClient) BeginRefreshAll(ctx context.Context, scope string, options *AlertsClientBeginRefreshAllOptions) (*runtime.Poller[AlertsClientRefreshAllResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.refreshAll(ctx, scope, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[AlertsClientRefreshAllResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[AlertsClientRefreshAllResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// RefreshAll - Refresh all alerts for a resource scope.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-01-preview
func (client *AlertsClient) refreshAll(ctx context.Context, scope string, options *AlertsClientBeginRefreshAllOptions) (*http.Response, error) {
	req, err := client.refreshAllCreateRequest(ctx, scope, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// refreshAllCreateRequest creates the RefreshAll request.
func (client *AlertsClient) refreshAllCreateRequest(ctx context.Context, scope string, options *AlertsClientBeginRefreshAllOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/roleManagementAlerts/refresh"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Update - Update an alert.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-08-01-preview
//   - scope - The scope of the alert.
//   - alertID - The name of the alert to dismiss.
//   - parameters - Parameters for the alert.
//   - options - AlertsClientUpdateOptions contains the optional parameters for the AlertsClient.Update method.
func (client *AlertsClient) Update(ctx context.Context, scope string, alertID string, parameters Alert, options *AlertsClientUpdateOptions) (AlertsClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, scope, alertID, parameters, options)
	if err != nil {
		return AlertsClientUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AlertsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusNoContent) {
		return AlertsClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return AlertsClientUpdateResponse{}, nil
}

// updateCreateRequest creates the Update request.
func (client *AlertsClient) updateCreateRequest(ctx context.Context, scope string, alertID string, parameters Alert, options *AlertsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/{scope}/providers/Microsoft.Authorization/roleManagementAlerts/{alertId}"
	urlPath = strings.ReplaceAll(urlPath, "{scope}", scope)
	urlPath = strings.ReplaceAll(urlPath, "{alertId}", alertID)
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-08-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}
