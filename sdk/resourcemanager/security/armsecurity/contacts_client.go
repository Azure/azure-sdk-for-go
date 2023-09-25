//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsecurity

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

// ContactsClient contains the methods for the SecurityContacts group.
// Don't use this type directly, use NewContactsClient() instead.
type ContactsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewContactsClient creates a new instance of ContactsClient with the specified values.
//   - subscriptionID - Azure subscription ID
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewContactsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ContactsClient, error) {
	cl, err := arm.NewClient(moduleName+".ContactsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ContactsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// Create - Create security contact configurations for the subscription
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-01-01-preview
//   - securityContactName - Name of the security contact object
//   - securityContact - Security contact object
//   - options - ContactsClientCreateOptions contains the optional parameters for the ContactsClient.Create method.
func (client *ContactsClient) Create(ctx context.Context, securityContactName string, securityContact Contact, options *ContactsClientCreateOptions) (ContactsClientCreateResponse, error) {
	var err error
	req, err := client.createCreateRequest(ctx, securityContactName, securityContact, options)
	if err != nil {
		return ContactsClientCreateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ContactsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return ContactsClientCreateResponse{}, err
	}
	resp, err := client.createHandleResponse(httpResp)
	return resp, err
}

// createCreateRequest creates the Create request.
func (client *ContactsClient) createCreateRequest(ctx context.Context, securityContactName string, securityContact Contact, options *ContactsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/securityContacts/{securityContactName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if securityContactName == "" {
		return nil, errors.New("parameter securityContactName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{securityContactName}", url.PathEscape(securityContactName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, securityContact); err != nil {
	return nil, err
}
	return req, nil
}

// createHandleResponse handles the Create response.
func (client *ContactsClient) createHandleResponse(resp *http.Response) (ContactsClientCreateResponse, error) {
	result := ContactsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Contact); err != nil {
		return ContactsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Delete security contact configurations for the subscription
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-01-01-preview
//   - securityContactName - Name of the security contact object
//   - options - ContactsClientDeleteOptions contains the optional parameters for the ContactsClient.Delete method.
func (client *ContactsClient) Delete(ctx context.Context, securityContactName string, options *ContactsClientDeleteOptions) (ContactsClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, securityContactName, options)
	if err != nil {
		return ContactsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ContactsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return ContactsClientDeleteResponse{}, err
	}
	return ContactsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ContactsClient) deleteCreateRequest(ctx context.Context, securityContactName string, options *ContactsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/securityContacts/{securityContactName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if securityContactName == "" {
		return nil, errors.New("parameter securityContactName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{securityContactName}", url.PathEscape(securityContactName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get Default Security contact configurations for the subscription
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-01-01-preview
//   - securityContactName - Name of the security contact object
//   - options - ContactsClientGetOptions contains the optional parameters for the ContactsClient.Get method.
func (client *ContactsClient) Get(ctx context.Context, securityContactName string, options *ContactsClientGetOptions) (ContactsClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, securityContactName, options)
	if err != nil {
		return ContactsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ContactsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ContactsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ContactsClient) getCreateRequest(ctx context.Context, securityContactName string, options *ContactsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/securityContacts/{securityContactName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if securityContactName == "" {
		return nil, errors.New("parameter securityContactName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{securityContactName}", url.PathEscape(securityContactName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ContactsClient) getHandleResponse(resp *http.Response) (ContactsClientGetResponse, error) {
	result := ContactsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Contact); err != nil {
		return ContactsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List all security contact configurations for the subscription
//
// Generated from API version 2020-01-01-preview
//   - options - ContactsClientListOptions contains the optional parameters for the ContactsClient.NewListPager method.
func (client *ContactsClient) NewListPager(options *ContactsClientListOptions) (*runtime.Pager[ContactsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[ContactsClientListResponse]{
		More: func(page ContactsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ContactsClientListResponse) (ContactsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ContactsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return ContactsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ContactsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ContactsClient) listCreateRequest(ctx context.Context, options *ContactsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Security/securityContacts"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-01-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ContactsClient) listHandleResponse(resp *http.Response) (ContactsClientListResponse, error) {
	result := ContactsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ContactList); err != nil {
		return ContactsClientListResponse{}, err
	}
	return result, nil
}

