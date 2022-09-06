//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armconsumption

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

// LotsClient contains the methods for the Lots group.
// Don't use this type directly, use NewLotsClient() instead.
type LotsClient struct {
	host string
	pl   runtime.Pipeline
}

// NewLotsClient creates a new instance of LotsClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewLotsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*LotsClient, error) {
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
	client := &LotsClient{
		host: ep,
		pl:   pl,
	}
	return client, nil
}

// NewListByBillingAccountPager - Lists all Microsoft Azure consumption commitments for a billing account. The API is only
// supported for Microsoft Customer Agreements (MCA) and Direct Enterprise Agreement (EA) billing accounts.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-10-01
// billingAccountID - BillingAccount ID
// options - LotsClientListByBillingAccountOptions contains the optional parameters for the LotsClient.ListByBillingAccount
// method.
func (client *LotsClient) NewListByBillingAccountPager(billingAccountID string, options *LotsClientListByBillingAccountOptions) *runtime.Pager[LotsClientListByBillingAccountResponse] {
	return runtime.NewPager(runtime.PagingHandler[LotsClientListByBillingAccountResponse]{
		More: func(page LotsClientListByBillingAccountResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *LotsClientListByBillingAccountResponse) (LotsClientListByBillingAccountResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByBillingAccountCreateRequest(ctx, billingAccountID, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return LotsClientListByBillingAccountResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return LotsClientListByBillingAccountResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return LotsClientListByBillingAccountResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByBillingAccountHandleResponse(resp)
		},
	})
}

// listByBillingAccountCreateRequest creates the ListByBillingAccount request.
func (client *LotsClient) listByBillingAccountCreateRequest(ctx context.Context, billingAccountID string, options *LotsClientListByBillingAccountOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/providers/Microsoft.Consumption/lots"
	if billingAccountID == "" {
		return nil, errors.New("parameter billingAccountID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountId}", url.PathEscape(billingAccountID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByBillingAccountHandleResponse handles the ListByBillingAccount response.
func (client *LotsClient) listByBillingAccountHandleResponse(resp *http.Response) (LotsClientListByBillingAccountResponse, error) {
	result := LotsClientListByBillingAccountResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Lots); err != nil {
		return LotsClientListByBillingAccountResponse{}, err
	}
	return result, nil
}

// NewListByBillingProfilePager - Lists all Azure credits for a billing account or a billing profile. The API is only supported
// for Microsoft Customer Agreements (MCA) billing accounts.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-10-01
// billingAccountID - BillingAccount ID
// billingProfileID - Azure Billing Profile ID.
// options - LotsClientListByBillingProfileOptions contains the optional parameters for the LotsClient.ListByBillingProfile
// method.
func (client *LotsClient) NewListByBillingProfilePager(billingAccountID string, billingProfileID string, options *LotsClientListByBillingProfileOptions) *runtime.Pager[LotsClientListByBillingProfileResponse] {
	return runtime.NewPager(runtime.PagingHandler[LotsClientListByBillingProfileResponse]{
		More: func(page LotsClientListByBillingProfileResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *LotsClientListByBillingProfileResponse) (LotsClientListByBillingProfileResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByBillingProfileCreateRequest(ctx, billingAccountID, billingProfileID, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return LotsClientListByBillingProfileResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return LotsClientListByBillingProfileResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return LotsClientListByBillingProfileResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByBillingProfileHandleResponse(resp)
		},
	})
}

// listByBillingProfileCreateRequest creates the ListByBillingProfile request.
func (client *LotsClient) listByBillingProfileCreateRequest(ctx context.Context, billingAccountID string, billingProfileID string, options *LotsClientListByBillingProfileOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/billingProfiles/{billingProfileId}/providers/Microsoft.Consumption/lots"
	if billingAccountID == "" {
		return nil, errors.New("parameter billingAccountID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountId}", url.PathEscape(billingAccountID))
	if billingProfileID == "" {
		return nil, errors.New("parameter billingProfileID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingProfileId}", url.PathEscape(billingProfileID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByBillingProfileHandleResponse handles the ListByBillingProfile response.
func (client *LotsClient) listByBillingProfileHandleResponse(resp *http.Response) (LotsClientListByBillingProfileResponse, error) {
	result := LotsClientListByBillingProfileResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Lots); err != nil {
		return LotsClientListByBillingProfileResponse{}, err
	}
	return result, nil
}

// NewListByCustomerPager - Lists all Azure credits for a customer. The API is only supported for Microsoft Partner Agreements
// (MPA) billing accounts.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-10-01
// billingAccountID - BillingAccount ID
// customerID - Customer ID
// options - LotsClientListByCustomerOptions contains the optional parameters for the LotsClient.ListByCustomer method.
func (client *LotsClient) NewListByCustomerPager(billingAccountID string, customerID string, options *LotsClientListByCustomerOptions) *runtime.Pager[LotsClientListByCustomerResponse] {
	return runtime.NewPager(runtime.PagingHandler[LotsClientListByCustomerResponse]{
		More: func(page LotsClientListByCustomerResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *LotsClientListByCustomerResponse) (LotsClientListByCustomerResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByCustomerCreateRequest(ctx, billingAccountID, customerID, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return LotsClientListByCustomerResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return LotsClientListByCustomerResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return LotsClientListByCustomerResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByCustomerHandleResponse(resp)
		},
	})
}

// listByCustomerCreateRequest creates the ListByCustomer request.
func (client *LotsClient) listByCustomerCreateRequest(ctx context.Context, billingAccountID string, customerID string, options *LotsClientListByCustomerOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountId}/customers/{customerId}/providers/Microsoft.Consumption/lots"
	if billingAccountID == "" {
		return nil, errors.New("parameter billingAccountID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountId}", url.PathEscape(billingAccountID))
	if customerID == "" {
		return nil, errors.New("parameter customerID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{customerId}", url.PathEscape(customerID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByCustomerHandleResponse handles the ListByCustomer response.
func (client *LotsClient) listByCustomerHandleResponse(resp *http.Response) (LotsClientListByCustomerResponse, error) {
	result := LotsClientListByCustomerResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Lots); err != nil {
		return LotsClientListByCustomerResponse{}, err
	}
	return result, nil
}
