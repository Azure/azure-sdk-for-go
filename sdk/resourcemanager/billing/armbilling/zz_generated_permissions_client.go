//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armbilling

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

// PermissionsClient contains the methods for the BillingPermissions group.
// Don't use this type directly, use NewPermissionsClient() instead.
type PermissionsClient struct {
	host string
	pl   runtime.Pipeline
}

// NewPermissionsClient creates a new instance of PermissionsClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewPermissionsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*PermissionsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublicCloud.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &PermissionsClient{
		host: ep,
		pl:   pl,
	}
	return client, nil
}

// NewListByBillingAccountPager - Lists the billing permissions the caller has on a billing account.
// If the operation fails it returns an *azcore.ResponseError type.
// billingAccountName - The ID that uniquely identifies a billing account.
// options - PermissionsClientListByBillingAccountOptions contains the optional parameters for the PermissionsClient.ListByBillingAccount
// method.
func (client *PermissionsClient) NewListByBillingAccountPager(billingAccountName string, options *PermissionsClientListByBillingAccountOptions) *runtime.Pager[PermissionsClientListByBillingAccountResponse] {
	return runtime.NewPager(runtime.PageProcessor[PermissionsClientListByBillingAccountResponse]{
		More: func(page PermissionsClientListByBillingAccountResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PermissionsClientListByBillingAccountResponse) (PermissionsClientListByBillingAccountResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByBillingAccountCreateRequest(ctx, billingAccountName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PermissionsClientListByBillingAccountResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return PermissionsClientListByBillingAccountResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PermissionsClientListByBillingAccountResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByBillingAccountHandleResponse(resp)
		},
	})
}

// listByBillingAccountCreateRequest creates the ListByBillingAccount request.
func (client *PermissionsClient) listByBillingAccountCreateRequest(ctx context.Context, billingAccountName string, options *PermissionsClientListByBillingAccountOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingPermissions"
	if billingAccountName == "" {
		return nil, errors.New("parameter billingAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountName}", url.PathEscape(billingAccountName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByBillingAccountHandleResponse handles the ListByBillingAccount response.
func (client *PermissionsClient) listByBillingAccountHandleResponse(resp *http.Response) (PermissionsClientListByBillingAccountResponse, error) {
	result := PermissionsClientListByBillingAccountResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PermissionsListResult); err != nil {
		return PermissionsClientListByBillingAccountResponse{}, err
	}
	return result, nil
}

// NewListByBillingProfilePager - Lists the billing permissions the caller has on a billing profile.
// If the operation fails it returns an *azcore.ResponseError type.
// billingAccountName - The ID that uniquely identifies a billing account.
// billingProfileName - The ID that uniquely identifies a billing profile.
// options - PermissionsClientListByBillingProfileOptions contains the optional parameters for the PermissionsClient.ListByBillingProfile
// method.
func (client *PermissionsClient) NewListByBillingProfilePager(billingAccountName string, billingProfileName string, options *PermissionsClientListByBillingProfileOptions) *runtime.Pager[PermissionsClientListByBillingProfileResponse] {
	return runtime.NewPager(runtime.PageProcessor[PermissionsClientListByBillingProfileResponse]{
		More: func(page PermissionsClientListByBillingProfileResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PermissionsClientListByBillingProfileResponse) (PermissionsClientListByBillingProfileResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByBillingProfileCreateRequest(ctx, billingAccountName, billingProfileName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PermissionsClientListByBillingProfileResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return PermissionsClientListByBillingProfileResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PermissionsClientListByBillingProfileResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByBillingProfileHandleResponse(resp)
		},
	})
}

// listByBillingProfileCreateRequest creates the ListByBillingProfile request.
func (client *PermissionsClient) listByBillingProfileCreateRequest(ctx context.Context, billingAccountName string, billingProfileName string, options *PermissionsClientListByBillingProfileOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/billingPermissions"
	if billingAccountName == "" {
		return nil, errors.New("parameter billingAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountName}", url.PathEscape(billingAccountName))
	if billingProfileName == "" {
		return nil, errors.New("parameter billingProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingProfileName}", url.PathEscape(billingProfileName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByBillingProfileHandleResponse handles the ListByBillingProfile response.
func (client *PermissionsClient) listByBillingProfileHandleResponse(resp *http.Response) (PermissionsClientListByBillingProfileResponse, error) {
	result := PermissionsClientListByBillingProfileResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PermissionsListResult); err != nil {
		return PermissionsClientListByBillingProfileResponse{}, err
	}
	return result, nil
}

// NewListByCustomerPager - Lists the billing permissions the caller has for a customer.
// If the operation fails it returns an *azcore.ResponseError type.
// billingAccountName - The ID that uniquely identifies a billing account.
// customerName - The ID that uniquely identifies a customer.
// options - PermissionsClientListByCustomerOptions contains the optional parameters for the PermissionsClient.ListByCustomer
// method.
func (client *PermissionsClient) NewListByCustomerPager(billingAccountName string, customerName string, options *PermissionsClientListByCustomerOptions) *runtime.Pager[PermissionsClientListByCustomerResponse] {
	return runtime.NewPager(runtime.PageProcessor[PermissionsClientListByCustomerResponse]{
		More: func(page PermissionsClientListByCustomerResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PermissionsClientListByCustomerResponse) (PermissionsClientListByCustomerResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByCustomerCreateRequest(ctx, billingAccountName, customerName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PermissionsClientListByCustomerResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return PermissionsClientListByCustomerResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PermissionsClientListByCustomerResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByCustomerHandleResponse(resp)
		},
	})
}

// listByCustomerCreateRequest creates the ListByCustomer request.
func (client *PermissionsClient) listByCustomerCreateRequest(ctx context.Context, billingAccountName string, customerName string, options *PermissionsClientListByCustomerOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/customers/{customerName}/billingPermissions"
	if billingAccountName == "" {
		return nil, errors.New("parameter billingAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountName}", url.PathEscape(billingAccountName))
	if customerName == "" {
		return nil, errors.New("parameter customerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{customerName}", url.PathEscape(customerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByCustomerHandleResponse handles the ListByCustomer response.
func (client *PermissionsClient) listByCustomerHandleResponse(resp *http.Response) (PermissionsClientListByCustomerResponse, error) {
	result := PermissionsClientListByCustomerResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PermissionsListResult); err != nil {
		return PermissionsClientListByCustomerResponse{}, err
	}
	return result, nil
}

// NewListByInvoiceSectionsPager - Lists the billing permissions the caller has on an invoice section.
// If the operation fails it returns an *azcore.ResponseError type.
// billingAccountName - The ID that uniquely identifies a billing account.
// billingProfileName - The ID that uniquely identifies a billing profile.
// invoiceSectionName - The ID that uniquely identifies an invoice section.
// options - PermissionsClientListByInvoiceSectionsOptions contains the optional parameters for the PermissionsClient.ListByInvoiceSections
// method.
func (client *PermissionsClient) NewListByInvoiceSectionsPager(billingAccountName string, billingProfileName string, invoiceSectionName string, options *PermissionsClientListByInvoiceSectionsOptions) *runtime.Pager[PermissionsClientListByInvoiceSectionsResponse] {
	return runtime.NewPager(runtime.PageProcessor[PermissionsClientListByInvoiceSectionsResponse]{
		More: func(page PermissionsClientListByInvoiceSectionsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PermissionsClientListByInvoiceSectionsResponse) (PermissionsClientListByInvoiceSectionsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByInvoiceSectionsCreateRequest(ctx, billingAccountName, billingProfileName, invoiceSectionName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PermissionsClientListByInvoiceSectionsResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return PermissionsClientListByInvoiceSectionsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PermissionsClientListByInvoiceSectionsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByInvoiceSectionsHandleResponse(resp)
		},
	})
}

// listByInvoiceSectionsCreateRequest creates the ListByInvoiceSections request.
func (client *PermissionsClient) listByInvoiceSectionsCreateRequest(ctx context.Context, billingAccountName string, billingProfileName string, invoiceSectionName string, options *PermissionsClientListByInvoiceSectionsOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Billing/billingAccounts/{billingAccountName}/billingProfiles/{billingProfileName}/invoiceSections/{invoiceSectionName}/billingPermissions"
	if billingAccountName == "" {
		return nil, errors.New("parameter billingAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingAccountName}", url.PathEscape(billingAccountName))
	if billingProfileName == "" {
		return nil, errors.New("parameter billingProfileName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{billingProfileName}", url.PathEscape(billingProfileName))
	if invoiceSectionName == "" {
		return nil, errors.New("parameter invoiceSectionName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{invoiceSectionName}", url.PathEscape(invoiceSectionName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listByInvoiceSectionsHandleResponse handles the ListByInvoiceSections response.
func (client *PermissionsClient) listByInvoiceSectionsHandleResponse(resp *http.Response) (PermissionsClientListByInvoiceSectionsResponse, error) {
	result := PermissionsClientListByInvoiceSectionsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.PermissionsListResult); err != nil {
		return PermissionsClientListByInvoiceSectionsResponse{}, err
	}
	return result, nil
}
