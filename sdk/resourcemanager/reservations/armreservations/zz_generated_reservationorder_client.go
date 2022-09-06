//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armreservations

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// ReservationOrderClient contains the methods for the ReservationOrder group.
// Don't use this type directly, use NewReservationOrderClient() instead.
type ReservationOrderClient struct {
	host string
	pl   runtime.Pipeline
}

// NewReservationOrderClient creates a new instance of ReservationOrderClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewReservationOrderClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*ReservationOrderClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &ReservationOrderClient{
		host: ep,
		pl:   pl,
	}
	return client, nil
}

// Calculate - Calculate price for placing a ReservationOrder.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// body - Information needed for calculate or purchase reservation
// options - ReservationOrderClientCalculateOptions contains the optional parameters for the ReservationOrderClient.Calculate
// method.
func (client *ReservationOrderClient) Calculate(ctx context.Context, body PurchaseRequest, options *ReservationOrderClientCalculateOptions) (ReservationOrderClientCalculateResponse, error) {
	req, err := client.calculateCreateRequest(ctx, body, options)
	if err != nil {
		return ReservationOrderClientCalculateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ReservationOrderClientCalculateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ReservationOrderClientCalculateResponse{}, runtime.NewResponseError(resp)
	}
	return client.calculateHandleResponse(resp)
}

// calculateCreateRequest creates the Calculate request.
func (client *ReservationOrderClient) calculateCreateRequest(ctx context.Context, body PurchaseRequest, options *ReservationOrderClientCalculateOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Capacity/calculatePrice"
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}

// calculateHandleResponse handles the Calculate response.
func (client *ReservationOrderClient) calculateHandleResponse(resp *http.Response) (ReservationOrderClientCalculateResponse, error) {
	result := ReservationOrderClientCalculateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.CalculatePriceResponse); err != nil {
		return ReservationOrderClientCalculateResponse{}, err
	}
	return result, nil
}

// ChangeDirectory - Change directory (tenant) of ReservationOrder and all Reservation under it to specified tenant id
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// reservationOrderID - Order Id of the reservation
// body - Information needed to change directory of reservation order
// options - ReservationOrderClientChangeDirectoryOptions contains the optional parameters for the ReservationOrderClient.ChangeDirectory
// method.
func (client *ReservationOrderClient) ChangeDirectory(ctx context.Context, reservationOrderID string, body ChangeDirectoryRequest, options *ReservationOrderClientChangeDirectoryOptions) (ReservationOrderClientChangeDirectoryResponse, error) {
	req, err := client.changeDirectoryCreateRequest(ctx, reservationOrderID, body, options)
	if err != nil {
		return ReservationOrderClientChangeDirectoryResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ReservationOrderClientChangeDirectoryResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ReservationOrderClientChangeDirectoryResponse{}, runtime.NewResponseError(resp)
	}
	return client.changeDirectoryHandleResponse(resp)
}

// changeDirectoryCreateRequest creates the ChangeDirectory request.
func (client *ReservationOrderClient) changeDirectoryCreateRequest(ctx context.Context, reservationOrderID string, body ChangeDirectoryRequest, options *ReservationOrderClientChangeDirectoryOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Capacity/reservationOrders/{reservationOrderId}/changeDirectory"
	if reservationOrderID == "" {
		return nil, errors.New("parameter reservationOrderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reservationOrderId}", url.PathEscape(reservationOrderID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}

// changeDirectoryHandleResponse handles the ChangeDirectory response.
func (client *ReservationOrderClient) changeDirectoryHandleResponse(resp *http.Response) (ReservationOrderClientChangeDirectoryResponse, error) {
	result := ReservationOrderClientChangeDirectoryResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ChangeDirectoryResponse); err != nil {
		return ReservationOrderClientChangeDirectoryResponse{}, err
	}
	return result, nil
}

// Get - Get the details of the ReservationOrder.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// reservationOrderID - Order Id of the reservation
// options - ReservationOrderClientGetOptions contains the optional parameters for the ReservationOrderClient.Get method.
func (client *ReservationOrderClient) Get(ctx context.Context, reservationOrderID string, options *ReservationOrderClientGetOptions) (ReservationOrderClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, reservationOrderID, options)
	if err != nil {
		return ReservationOrderClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ReservationOrderClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ReservationOrderClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ReservationOrderClient) getCreateRequest(ctx context.Context, reservationOrderID string, options *ReservationOrderClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Capacity/reservationOrders/{reservationOrderId}"
	if reservationOrderID == "" {
		return nil, errors.New("parameter reservationOrderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reservationOrderId}", url.PathEscape(reservationOrderID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ReservationOrderClient) getHandleResponse(resp *http.Response) (ReservationOrderClientGetResponse, error) {
	result := ReservationOrderClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReservationOrderResponse); err != nil {
		return ReservationOrderClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List of all the ReservationOrders that the user has access to in the current tenant.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// options - ReservationOrderClientListOptions contains the optional parameters for the ReservationOrderClient.List method.
func (client *ReservationOrderClient) NewListPager(options *ReservationOrderClientListOptions) *runtime.Pager[ReservationOrderClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[ReservationOrderClientListResponse]{
		More: func(page ReservationOrderClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ReservationOrderClientListResponse) (ReservationOrderClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ReservationOrderClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return ReservationOrderClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ReservationOrderClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ReservationOrderClient) listCreateRequest(ctx context.Context, options *ReservationOrderClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Capacity/reservationOrders"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ReservationOrderClient) listHandleResponse(resp *http.Response) (ReservationOrderClientListResponse, error) {
	result := ReservationOrderClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ReservationOrderList); err != nil {
		return ReservationOrderClientListResponse{}, err
	}
	return result, nil
}

// BeginPurchase - Purchase ReservationOrder and create resource under the specified URI.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
// reservationOrderID - Order Id of the reservation
// body - Information needed for calculate or purchase reservation
// options - ReservationOrderClientBeginPurchaseOptions contains the optional parameters for the ReservationOrderClient.BeginPurchase
// method.
func (client *ReservationOrderClient) BeginPurchase(ctx context.Context, reservationOrderID string, body PurchaseRequest, options *ReservationOrderClientBeginPurchaseOptions) (*runtime.Poller[ReservationOrderClientPurchaseResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.purchase(ctx, reservationOrderID, body, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller(resp, client.pl, &runtime.NewPollerOptions[ReservationOrderClientPurchaseResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
	} else {
		return runtime.NewPollerFromResumeToken[ReservationOrderClientPurchaseResponse](options.ResumeToken, client.pl, nil)
	}
}

// Purchase - Purchase ReservationOrder and create resource under the specified URI.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-03-01
func (client *ReservationOrderClient) purchase(ctx context.Context, reservationOrderID string, body PurchaseRequest, options *ReservationOrderClientBeginPurchaseOptions) (*http.Response, error) {
	req, err := client.purchaseCreateRequest(ctx, reservationOrderID, body, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// purchaseCreateRequest creates the Purchase request.
func (client *ReservationOrderClient) purchaseCreateRequest(ctx context.Context, reservationOrderID string, body PurchaseRequest, options *ReservationOrderClientBeginPurchaseOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Capacity/reservationOrders/{reservationOrderId}"
	if reservationOrderID == "" {
		return nil, errors.New("parameter reservationOrderID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{reservationOrderId}", url.PathEscape(reservationOrderID))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-03-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}
