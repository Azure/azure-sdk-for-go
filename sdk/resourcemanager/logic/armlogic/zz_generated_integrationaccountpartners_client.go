//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armlogic

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
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
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewIntegrationAccountPartnersClient creates a new instance of IntegrationAccountPartnersClient with the specified values.
// subscriptionID - The subscription id.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewIntegrationAccountPartnersClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *IntegrationAccountPartnersClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Endpoint) == 0 {
		cp.Endpoint = arm.AzurePublicCloud
	}
	client := &IntegrationAccountPartnersClient{
		subscriptionID: subscriptionID,
		host:           string(cp.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, &cp),
	}
	return client
}

// CreateOrUpdate - Creates or updates an integration account partner.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// integrationAccountName - The integration account name.
// partnerName - The integration account partner name.
// partner - The integration account partner.
// options - IntegrationAccountPartnersClientCreateOrUpdateOptions contains the optional parameters for the IntegrationAccountPartnersClient.CreateOrUpdate
// method.
func (client *IntegrationAccountPartnersClient) CreateOrUpdate(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, partner IntegrationAccountPartner, options *IntegrationAccountPartnersClientCreateOrUpdateOptions) (IntegrationAccountPartnersClientCreateOrUpdateResponse, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, integrationAccountName, partnerName, partner, options)
	if err != nil {
		return IntegrationAccountPartnersClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return IntegrationAccountPartnersClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return IntegrationAccountPartnersClientCreateOrUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createOrUpdateHandleResponse(resp)
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
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, partner)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *IntegrationAccountPartnersClient) createOrUpdateHandleResponse(resp *http.Response) (IntegrationAccountPartnersClientCreateOrUpdateResponse, error) {
	result := IntegrationAccountPartnersClientCreateOrUpdateResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.IntegrationAccountPartner); err != nil {
		return IntegrationAccountPartnersClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes an integration account partner.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// integrationAccountName - The integration account name.
// partnerName - The integration account partner name.
// options - IntegrationAccountPartnersClientDeleteOptions contains the optional parameters for the IntegrationAccountPartnersClient.Delete
// method.
func (client *IntegrationAccountPartnersClient) Delete(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, options *IntegrationAccountPartnersClientDeleteOptions) (IntegrationAccountPartnersClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, integrationAccountName, partnerName, options)
	if err != nil {
		return IntegrationAccountPartnersClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return IntegrationAccountPartnersClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return IntegrationAccountPartnersClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return IntegrationAccountPartnersClientDeleteResponse{RawResponse: resp}, nil
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
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// Get - Gets an integration account partner.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// integrationAccountName - The integration account name.
// partnerName - The integration account partner name.
// options - IntegrationAccountPartnersClientGetOptions contains the optional parameters for the IntegrationAccountPartnersClient.Get
// method.
func (client *IntegrationAccountPartnersClient) Get(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, options *IntegrationAccountPartnersClientGetOptions) (IntegrationAccountPartnersClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, integrationAccountName, partnerName, options)
	if err != nil {
		return IntegrationAccountPartnersClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return IntegrationAccountPartnersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return IntegrationAccountPartnersClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *IntegrationAccountPartnersClient) getHandleResponse(resp *http.Response) (IntegrationAccountPartnersClientGetResponse, error) {
	result := IntegrationAccountPartnersClientGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.IntegrationAccountPartner); err != nil {
		return IntegrationAccountPartnersClientGetResponse{}, err
	}
	return result, nil
}

// List - Gets a list of integration account partners.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// integrationAccountName - The integration account name.
// options - IntegrationAccountPartnersClientListOptions contains the optional parameters for the IntegrationAccountPartnersClient.List
// method.
func (client *IntegrationAccountPartnersClient) List(resourceGroupName string, integrationAccountName string, options *IntegrationAccountPartnersClientListOptions) *IntegrationAccountPartnersClientListPager {
	return &IntegrationAccountPartnersClientListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, integrationAccountName, options)
		},
		advancer: func(ctx context.Context, resp IntegrationAccountPartnersClientListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.IntegrationAccountPartnerListResult.NextLink)
		},
	}
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
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
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
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *IntegrationAccountPartnersClient) listHandleResponse(resp *http.Response) (IntegrationAccountPartnersClientListResponse, error) {
	result := IntegrationAccountPartnersClientListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.IntegrationAccountPartnerListResult); err != nil {
		return IntegrationAccountPartnersClientListResponse{}, err
	}
	return result, nil
}

// ListContentCallbackURL - Get the content callback url.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The resource group name.
// integrationAccountName - The integration account name.
// partnerName - The integration account partner name.
// options - IntegrationAccountPartnersClientListContentCallbackURLOptions contains the optional parameters for the IntegrationAccountPartnersClient.ListContentCallbackURL
// method.
func (client *IntegrationAccountPartnersClient) ListContentCallbackURL(ctx context.Context, resourceGroupName string, integrationAccountName string, partnerName string, listContentCallbackURL GetCallbackURLParameters, options *IntegrationAccountPartnersClientListContentCallbackURLOptions) (IntegrationAccountPartnersClientListContentCallbackURLResponse, error) {
	req, err := client.listContentCallbackURLCreateRequest(ctx, resourceGroupName, integrationAccountName, partnerName, listContentCallbackURL, options)
	if err != nil {
		return IntegrationAccountPartnersClientListContentCallbackURLResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return IntegrationAccountPartnersClientListContentCallbackURLResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return IntegrationAccountPartnersClientListContentCallbackURLResponse{}, runtime.NewResponseError(resp)
	}
	return client.listContentCallbackURLHandleResponse(resp)
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
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2019-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, listContentCallbackURL)
}

// listContentCallbackURLHandleResponse handles the ListContentCallbackURL response.
func (client *IntegrationAccountPartnersClient) listContentCallbackURLHandleResponse(resp *http.Response) (IntegrationAccountPartnersClientListContentCallbackURLResponse, error) {
	result := IntegrationAccountPartnersClientListContentCallbackURLResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.WorkflowTriggerCallbackURL); err != nil {
		return IntegrationAccountPartnersClientListContentCallbackURLResponse{}, err
	}
	return result, nil
}
