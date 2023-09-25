//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armlogic

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// IntegrationAccountPartnersClient contains the methods for the IntegrationAccountPartners group.
// Don't use this type directly, use NewIntegrationAccountPartnersClient() instead.
type IntegrationAccountPartnersClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewIntegrationAccountPartnersClient creates a new instance of IntegrationAccountPartnersClient with the specified values.
//   - subscriptionID - The subscription id.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewIntegrationAccountPartnersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*IntegrationAccountPartnersClient, error) {
	cl, err := arm.NewClient(moduleName+".IntegrationAccountPartnersClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &IntegrationAccountPartnersClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Creates or updates an integration account partner.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - integrationAccountName - The integration account name.
//   - partnerName - The integration account partner name.
//   - partner - The integration account partner.
//   - options - IntegrationAccountPartnersClientCreateOrUpdateOptions contains the optional parameters for the IntegrationAccountPartnersClient.CreateOrUpdate
//     method.
func (client *IntegrationAccountPartnersClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, partner IntegrationAccountPartner, options *IntegrationAccountPartnersClientCreateOrUpdateOptions) (IntegrationAccountPartnersClientCreateOrUpdateResponse, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, integrationAccountName, partnerName, partner, options)
	if err != nil {
		return IntegrationAccountPartnersClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return IntegrationAccountPartnersClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return IntegrationAccountPartnersClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *IntegrationAccountPartnersClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, partner IntegrationAccountPartner, options *IntegrationAccountPartnersClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{integrationAccountName}/partners/{partnerName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if integrationAccountName == "" {
		return nil, errors.New("parameter integrationAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{integrationAccountName}", url.PathEscape(integrationAccountName))
	if partnerName == "" {
		return nil, errors.New("parameter partnerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partnerName}", url.PathEscape(partnerName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, partner); err != nil {
	return nil, err
}
	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *IntegrationAccountPartnersClient) createOrUpdateHandleResponse(resp *http.Response) (IntegrationAccountPartnersClientCreateOrUpdateResponse, error) {
	result := IntegrationAccountPartnersClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.IntegrationAccountPartner); err != nil {
		return IntegrationAccountPartnersClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes an integration account partner.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - integrationAccountName - The integration account name.
//   - partnerName - The integration account partner name.
//   - options - IntegrationAccountPartnersClientDeleteOptions contains the optional parameters for the IntegrationAccountPartnersClient.Delete
//     method.
func (client *IntegrationAccountPartnersClient) Delete(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, options *IntegrationAccountPartnersClientDeleteOptions) (IntegrationAccountPartnersClientDeleteResponse, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, integrationAccountName, partnerName, options)
	if err != nil {
		return IntegrationAccountPartnersClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return IntegrationAccountPartnersClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return IntegrationAccountPartnersClientDeleteResponse{}, err
	}
	return IntegrationAccountPartnersClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *IntegrationAccountPartnersClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, options *IntegrationAccountPartnersClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{integrationAccountName}/partners/{partnerName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if integrationAccountName == "" {
		return nil, errors.New("parameter integrationAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{integrationAccountName}", url.PathEscape(integrationAccountName))
	if partnerName == "" {
		return nil, errors.New("parameter partnerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partnerName}", url.PathEscape(partnerName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets an integration account partner.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - integrationAccountName - The integration account name.
//   - partnerName - The integration account partner name.
//   - options - IntegrationAccountPartnersClientGetOptions contains the optional parameters for the IntegrationAccountPartnersClient.Get
//     method.
func (client *IntegrationAccountPartnersClient) Get(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, options *IntegrationAccountPartnersClientGetOptions) (IntegrationAccountPartnersClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, resourceGroupName, integrationAccountName, partnerName, options)
	if err != nil {
		return IntegrationAccountPartnersClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return IntegrationAccountPartnersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return IntegrationAccountPartnersClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *IntegrationAccountPartnersClient) getCreateRequest(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, options *IntegrationAccountPartnersClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{integrationAccountName}/partners/{partnerName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if integrationAccountName == "" {
		return nil, errors.New("parameter integrationAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{integrationAccountName}", url.PathEscape(integrationAccountName))
	if partnerName == "" {
		return nil, errors.New("parameter partnerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partnerName}", url.PathEscape(partnerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *IntegrationAccountPartnersClient) getHandleResponse(resp *http.Response) (IntegrationAccountPartnersClientGetResponse, error) {
	result := IntegrationAccountPartnersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.IntegrationAccountPartner); err != nil {
		return IntegrationAccountPartnersClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Gets a list of integration account partners.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - integrationAccountName - The integration account name.
//   - options - IntegrationAccountPartnersClientListOptions contains the optional parameters for the IntegrationAccountPartnersClient.NewListPager
//     method.
func (client *IntegrationAccountPartnersClient) NewListPager(resourceGroupName string, integrationAccountName string, options *IntegrationAccountPartnersClientListOptions) (*runtime.Pager[IntegrationAccountPartnersClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[IntegrationAccountPartnersClientListResponse]{
		More: func(page IntegrationAccountPartnersClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *IntegrationAccountPartnersClientListResponse) (IntegrationAccountPartnersClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, integrationAccountName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return IntegrationAccountPartnersClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return IntegrationAccountPartnersClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return IntegrationAccountPartnersClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *IntegrationAccountPartnersClient) listCreateRequest(ctx context.Context, resourceGroupName string, integrationAccountName string, options *IntegrationAccountPartnersClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{integrationAccountName}/partners"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if integrationAccountName == "" {
		return nil, errors.New("parameter integrationAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{integrationAccountName}", url.PathEscape(integrationAccountName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	if options != nil && options.Top != nil {
		reqQP.Set("$top", strconv.FormatInt(int64(*options.Top), 10))
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *IntegrationAccountPartnersClient) listHandleResponse(resp *http.Response) (IntegrationAccountPartnersClientListResponse, error) {
	result := IntegrationAccountPartnersClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.IntegrationAccountPartnerListResult); err != nil {
		return IntegrationAccountPartnersClientListResponse{}, err
	}
	return result, nil
}

// ListContentCallbackURL - Get the content callback url.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2019-05-01
//   - resourceGroupName - The resource group name.
//   - integrationAccountName - The integration account name.
//   - partnerName - The integration account partner name.
//   - options - IntegrationAccountPartnersClientListContentCallbackURLOptions contains the optional parameters for the IntegrationAccountPartnersClient.ListContentCallbackURL
//     method.
func (client *IntegrationAccountPartnersClient) ListContentCallbackURL(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, listContentCallbackURL GetCallbackURLParameters, options *IntegrationAccountPartnersClientListContentCallbackURLOptions) (IntegrationAccountPartnersClientListContentCallbackURLResponse, error) {
	var err error
	req, err := client.listContentCallbackURLCreateRequest(ctx, resourceGroupName, integrationAccountName, partnerName, listContentCallbackURL, options)
	if err != nil {
		return IntegrationAccountPartnersClientListContentCallbackURLResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return IntegrationAccountPartnersClientListContentCallbackURLResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return IntegrationAccountPartnersClientListContentCallbackURLResponse{}, err
	}
	resp, err := client.listContentCallbackURLHandleResponse(httpResp)
	return resp, err
}

// listContentCallbackURLCreateRequest creates the ListContentCallbackURL request.
func (client *IntegrationAccountPartnersClient) listContentCallbackURLCreateRequest(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, listContentCallbackURL GetCallbackURLParameters, options *IntegrationAccountPartnersClientListContentCallbackURLOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Logic/integrationAccounts/{integrationAccountName}/partners/{partnerName}/listContentCallbackUrl"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if integrationAccountName == "" {
		return nil, errors.New("parameter integrationAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{integrationAccountName}", url.PathEscape(integrationAccountName))
	if partnerName == "" {
		return nil, errors.New("parameter partnerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{partnerName}", url.PathEscape(partnerName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, listContentCallbackURL); err != nil {
	return nil, err
}
	return req, nil
}

// listContentCallbackURLHandleResponse handles the ListContentCallbackURL response.
func (client *IntegrationAccountPartnersClient) listContentCallbackURLHandleResponse(resp *http.Response) (IntegrationAccountPartnersClientListContentCallbackURLResponse, error) {
	result := IntegrationAccountPartnersClientListContentCallbackURLResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkflowTriggerCallbackURL); err != nil {
		return IntegrationAccountPartnersClientListContentCallbackURLResponse{}, err
	}
	return result, nil
}

