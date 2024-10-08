//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armbilling

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

// ReservationsClient contains the methods for the Reservations group.
// Don't use this type directly, use NewReservationsClient() instead.
type ReservationsClient struct {
	internal *arm.Client
}

// NewReservationsClient creates a new instance of ReservationsClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewReservationsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ReservationsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ReservationsClient{
		internal: cl,
	}
	return client, nil
}

// GetByReservationOrder - Get specific Reservation details in the billing account.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-04-01
//   - billingAccountName - The ID that uniquely identifies a billing account.
//   - reservationOrderID - Order Id of the reservation
//   - reservationID - Id of the reservation item
//   - options - ReservationsClientGetByReservationOrderOptions contains the optional parameters for the ReservationsClient.GetByReservationOrder
//     method.
func (client *ReservationsClient) GetByReservationOrder(ctx context.Context, billingAccountName string, reservationOrderID string, reservationID string, options *ReservationsClientGetByReservationOrderOptions) (ReservationsClientGetByReservationOrderResponse, error) {
	var err error
	const operationName = "ReservationsClient.GetByReservationOrder"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getByReservationOrderCreateRequest(ctx, billingAccountName, reservationOrderID, reservationID, options)
	if err != nil {
		return ReservationsClientGetByReservationOrderResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ReservationsClientGetByReservationOrderResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ReservationsClientGetByReservationOrderResponse{}, err
	}
	resp, err := client.getByReservationOrderHandleResponse(httpResp)
	return resp, err
}

// getByReservationOrderCreateRequest creates the GetByReservationOrder request.
func (client *ReservationsClient) getByReservationOrderCreateRequest(ctx context.Context, billingAccountName string, reservationOrderID string, reservationID string, options *ReservationsClientGetByReservationOrderOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/reservationOrders/{reservationOrderId}/reservations/{reservationId}"
	if billingAccountName == "" {
		return nil, errors.New("parameter billingAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountName}", url.PathEscape(billingAccountName))
	if reservationOrderID == "" {
		return nil, errors.New("parameter reservationOrderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reservationOrderId}", url.PathEscape(reservationOrderID))
	if reservationID == "" {
		return nil, errors.New("parameter reservationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reservationId}", url.PathEscape(reservationID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-04-01")
	if options != nil && options.Expand != nil {
		reqQP.Set("expand", *options.Expand)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getByReservationOrderHandleResponse handles the GetByReservationOrder response.
func (client *ReservationsClient) getByReservationOrderHandleResponse(resp *http.Response) (ReservationsClientGetByReservationOrderResponse, error) {
	result := ReservationsClientGetByReservationOrderResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Reservation); err != nil {
		return ReservationsClientGetByReservationOrderResponse{}, err
	}
	return result, nil
}

// NewListByBillingAccountPager - Lists the reservations in the billing account and the roll up counts of reservations group
// by provisioning states.
//
// Generated from API version 2024-04-01
//   - billingAccountName - The ID that uniquely identifies a billing account.
//   - options - ReservationsClientListByBillingAccountOptions contains the optional parameters for the ReservationsClient.NewListByBillingAccountPager
//     method.
func (client *ReservationsClient) NewListByBillingAccountPager(billingAccountName string, options *ReservationsClientListByBillingAccountOptions) *runtime.Pager[ReservationsClientListByBillingAccountResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReservationsClientListByBillingAccountResponse]{
		More: func(page ReservationsClientListByBillingAccountResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReservationsClientListByBillingAccountResponse) (ReservationsClientListByBillingAccountResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ReservationsClient.NewListByBillingAccountPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByBillingAccountCreateRequest(ctx, billingAccountName, options)
			}, nil)
			if err != nil {
				return ReservationsClientListByBillingAccountResponse{}, err
			}
			return client.listByBillingAccountHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByBillingAccountCreateRequest creates the ListByBillingAccount request.
func (client *ReservationsClient) listByBillingAccountCreateRequest(ctx context.Context, billingAccountName string, options *ReservationsClientListByBillingAccountOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/reservations"
	if billingAccountName == "" {
		return nil, errors.New("parameter billingAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountName}", url.PathEscape(billingAccountName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-04-01")
	if options != nil && options.Filter != nil {
		reqQP.Set("filter", *options.Filter)
	}
	if options != nil && options.OrderBy != nil {
		reqQP.Set("orderBy", *options.OrderBy)
	}
	if options != nil && options.RefreshSummary != nil {
		reqQP.Set("refreshSummary", *options.RefreshSummary)
	}
	if options != nil && options.SelectedState != nil {
		reqQP.Set("selectedState", *options.SelectedState)
	}
	if options != nil && options.Skiptoken != nil {
		reqQP.Set("skiptoken", strconv.FormatFloat(float64(*options.Skiptoken), 'f', -1, 32))
	}
	if options != nil && options.Take != nil {
		reqQP.Set("take", strconv.FormatFloat(float64(*options.Take), 'f', -1, 32))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByBillingAccountHandleResponse handles the ListByBillingAccount response.
func (client *ReservationsClient) listByBillingAccountHandleResponse(resp *http.Response) (ReservationsClientListByBillingAccountResponse, error) {
	result := ReservationsClientListByBillingAccountResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReservationsListResult); err != nil {
		return ReservationsClientListByBillingAccountResponse{}, err
	}
	return result, nil
}

// NewListByBillingProfilePager - Lists the reservations for a billing profile and the roll up counts of reservations group
// by provisioning state.
//
// Generated from API version 2024-04-01
//   - billingAccountName - The ID that uniquely identifies a billing account.
//   - billingProfileName - The ID that uniquely identifies a billing profile.
//   - options - ReservationsClientListByBillingProfileOptions contains the optional parameters for the ReservationsClient.NewListByBillingProfilePager
//     method.
func (client *ReservationsClient) NewListByBillingProfilePager(billingAccountName string, billingProfileName string, options *ReservationsClientListByBillingProfileOptions) *runtime.Pager[ReservationsClientListByBillingProfileResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReservationsClientListByBillingProfileResponse]{
		More: func(page ReservationsClientListByBillingProfileResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReservationsClientListByBillingProfileResponse) (ReservationsClientListByBillingProfileResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ReservationsClient.NewListByBillingProfilePager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByBillingProfileCreateRequest(ctx, billingAccountName, billingProfileName, options)
			}, nil)
			if err != nil {
				return ReservationsClientListByBillingProfileResponse{}, err
			}
			return client.listByBillingProfileHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByBillingProfileCreateRequest creates the ListByBillingProfile request.
func (client *ReservationsClient) listByBillingProfileCreateRequest(ctx context.Context, billingAccountName string, billingProfileName string, options *ReservationsClientListByBillingProfileOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/reservations"
	if billingAccountName == "" {
		return nil, errors.New("parameter billingAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountName}", url.PathEscape(billingAccountName))
	if billingProfileName == "" {
		return nil, errors.New("parameter billingProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingProfileName}", url.PathEscape(billingProfileName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-04-01")
	if options != nil && options.Filter != nil {
		reqQP.Set("filter", *options.Filter)
	}
	if options != nil && options.OrderBy != nil {
		reqQP.Set("orderBy", *options.OrderBy)
	}
	if options != nil && options.RefreshSummary != nil {
		reqQP.Set("refreshSummary", *options.RefreshSummary)
	}
	if options != nil && options.SelectedState != nil {
		reqQP.Set("selectedState", *options.SelectedState)
	}
	if options != nil && options.Skiptoken != nil {
		reqQP.Set("skiptoken", strconv.FormatFloat(float64(*options.Skiptoken), 'f', -1, 32))
	}
	if options != nil && options.Take != nil {
		reqQP.Set("take", strconv.FormatFloat(float64(*options.Take), 'f', -1, 32))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByBillingProfileHandleResponse handles the ListByBillingProfile response.
func (client *ReservationsClient) listByBillingProfileHandleResponse(resp *http.Response) (ReservationsClientListByBillingProfileResponse, error) {
	result := ReservationsClientListByBillingProfileResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReservationsListResult); err != nil {
		return ReservationsClientListByBillingProfileResponse{}, err
	}
	return result, nil
}

// NewListByReservationOrderPager - List Reservations within a single ReservationOrder in the billing account.
//
// Generated from API version 2024-04-01
//   - billingAccountName - The ID that uniquely identifies a billing account.
//   - reservationOrderID - Order Id of the reservation
//   - options - ReservationsClientListByReservationOrderOptions contains the optional parameters for the ReservationsClient.NewListByReservationOrderPager
//     method.
func (client *ReservationsClient) NewListByReservationOrderPager(billingAccountName string, reservationOrderID string, options *ReservationsClientListByReservationOrderOptions) *runtime.Pager[ReservationsClientListByReservationOrderResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReservationsClientListByReservationOrderResponse]{
		More: func(page ReservationsClientListByReservationOrderResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReservationsClientListByReservationOrderResponse) (ReservationsClientListByReservationOrderResponse, error) {
			ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "ReservationsClient.NewListByReservationOrderPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listByReservationOrderCreateRequest(ctx, billingAccountName, reservationOrderID, options)
			}, nil)
			if err != nil {
				return ReservationsClientListByReservationOrderResponse{}, err
			}
			return client.listByReservationOrderHandleResponse(resp)
		},
		Tracer: client.internal.Tracer(),
	})
}

// listByReservationOrderCreateRequest creates the ListByReservationOrder request.
func (client *ReservationsClient) listByReservationOrderCreateRequest(ctx context.Context, billingAccountName string, reservationOrderID string, options *ReservationsClientListByReservationOrderOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/reservationOrders/{reservationOrderId}/reservations"
	if billingAccountName == "" {
		return nil, errors.New("parameter billingAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountName}", url.PathEscape(billingAccountName))
	if reservationOrderID == "" {
		return nil, errors.New("parameter reservationOrderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reservationOrderId}", url.PathEscape(reservationOrderID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByReservationOrderHandleResponse handles the ListByReservationOrder response.
func (client *ReservationsClient) listByReservationOrderHandleResponse(resp *http.Response) (ReservationsClientListByReservationOrderResponse, error) {
	result := ReservationsClientListByReservationOrderResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReservationList); err != nil {
		return ReservationsClientListByReservationOrderResponse{}, err
	}
	return result, nil
}

// BeginUpdateByBillingAccount - Update reservation by billing account.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-04-01
//   - billingAccountName - The ID that uniquely identifies a billing account.
//   - reservationOrderID - Order Id of the reservation
//   - reservationID - Id of the reservation item
//   - body - Request body for patching a reservation
//   - options - ReservationsClientBeginUpdateByBillingAccountOptions contains the optional parameters for the ReservationsClient.BeginUpdateByBillingAccount
//     method.
func (client *ReservationsClient) BeginUpdateByBillingAccount(ctx context.Context, billingAccountName string, reservationOrderID string, reservationID string, body Patch, options *ReservationsClientBeginUpdateByBillingAccountOptions) (*runtime.Poller[ReservationsClientUpdateByBillingAccountResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.updateByBillingAccount(ctx, billingAccountName, reservationOrderID, reservationID, body, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[ReservationsClientUpdateByBillingAccountResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
			Tracer:        client.internal.Tracer(),
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, client.internal.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[ReservationsClientUpdateByBillingAccountResponse]{
			Tracer: client.internal.Tracer(),
		})
	}
}

// UpdateByBillingAccount - Update reservation by billing account.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2024-04-01
func (client *ReservationsClient) updateByBillingAccount(ctx context.Context, billingAccountName string, reservationOrderID string, reservationID string, body Patch, options *ReservationsClientBeginUpdateByBillingAccountOptions) (*http.Response, error) {
	var err error
	const operationName = "ReservationsClient.BeginUpdateByBillingAccount"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateByBillingAccountCreateRequest(ctx, billingAccountName, reservationOrderID, reservationID, body, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateByBillingAccountCreateRequest creates the UpdateByBillingAccount request.
func (client *ReservationsClient) updateByBillingAccountCreateRequest(ctx context.Context, billingAccountName string, reservationOrderID string, reservationID string, body Patch, options *ReservationsClientBeginUpdateByBillingAccountOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/reservationOrders/{reservationOrderId}/reservations/{reservationId}"
	if billingAccountName == "" {
		return nil, errors.New("parameter billingAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountName}", url.PathEscape(billingAccountName))
	if reservationOrderID == "" {
		return nil, errors.New("parameter reservationOrderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reservationOrderId}", url.PathEscape(reservationOrderID))
	if reservationID == "" {
		return nil, errors.New("parameter reservationID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reservationId}", url.PathEscape(reservationID))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2024-04-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, body); err != nil {
		return nil, err
	}
	return req, nil
}
