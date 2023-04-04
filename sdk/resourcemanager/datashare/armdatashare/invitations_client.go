//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armdatashare

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

// InvitationsClient contains the methods for the Invitations group.
// Don't use this type directly, use NewInvitationsClient() instead.
type InvitationsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewInvitationsClient creates a new instance of InvitationsClient with the specified values.
//   - subscriptionID - The subscription identifier
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewInvitationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*InvitationsClient, error) {
	cl, err := arm.NewClient(moduleName+".InvitationsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &InvitationsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Create - Create an invitation
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareName - The name of the share to send the invitation for.
//   - invitationName - The name of the invitation.
//   - invitation - Invitation details.
//   - options - InvitationsClientCreateOptions contains the optional parameters for the InvitationsClient.Create method.
func (client *InvitationsClient) Create(ctx context.Context, resourceGroupName string, accountName string, shareName string, invitationName string, invitation Invitation, options *InvitationsClientCreateOptions) (InvitationsClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, accountName, shareName, invitationName, invitation, options)
	if err != nil {
		return InvitationsClientCreateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return InvitationsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return InvitationsClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *InvitationsClient) createCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, invitationName string, invitation Invitation, options *InvitationsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}/invitations/{invitationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	if invitationName == "" {
		return nil, errors.New("parameter invitationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{invitationName}", url.PathEscape(invitationName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, invitation)
}

// createHandleResponse handles the Create response.
func (client *InvitationsClient) createHandleResponse(resp *http.Response) (InvitationsClientCreateResponse, error) {
	result := InvitationsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Invitation); err != nil {
		return InvitationsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Delete an invitation in a share
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareName - The name of the share.
//   - invitationName - The name of the invitation.
//   - options - InvitationsClientDeleteOptions contains the optional parameters for the InvitationsClient.Delete method.
func (client *InvitationsClient) Delete(ctx context.Context, resourceGroupName string, accountName string, shareName string, invitationName string, options *InvitationsClientDeleteOptions) (InvitationsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, shareName, invitationName, options)
	if err != nil {
		return InvitationsClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return InvitationsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return InvitationsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return InvitationsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *InvitationsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, invitationName string, options *InvitationsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}/invitations/{invitationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	if invitationName == "" {
		return nil, errors.New("parameter invitationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{invitationName}", url.PathEscape(invitationName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get an invitation in a share
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareName - The name of the share.
//   - invitationName - The name of the invitation.
//   - options - InvitationsClientGetOptions contains the optional parameters for the InvitationsClient.Get method.
func (client *InvitationsClient) Get(ctx context.Context, resourceGroupName string, accountName string, shareName string, invitationName string, options *InvitationsClientGetOptions) (InvitationsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, shareName, invitationName, options)
	if err != nil {
		return InvitationsClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return InvitationsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return InvitationsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *InvitationsClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, invitationName string, options *InvitationsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}/invitations/{invitationName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	if invitationName == "" {
		return nil, errors.New("parameter invitationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{invitationName}", url.PathEscape(invitationName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *InvitationsClient) getHandleResponse(resp *http.Response) (InvitationsClientGetResponse, error) {
	result := InvitationsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Invitation); err != nil {
		return InvitationsClientGetResponse{}, err
	}
	return result, nil
}

// NewListBySharePager - List invitations in a share
//
// Generated from API version 2020-09-01
//   - resourceGroupName - The resource group name.
//   - accountName - The name of the share account.
//   - shareName - The name of the share.
//   - options - InvitationsClientListByShareOptions contains the optional parameters for the InvitationsClient.NewListBySharePager
//     method.
func (client *InvitationsClient) NewListBySharePager(resourceGroupName string, accountName string, shareName string, options *InvitationsClientListByShareOptions) *runtime.Pager[InvitationsClientListByShareResponse] {
	return runtime.NewPager(runtime.PagingHandler[InvitationsClientListByShareResponse]{
		More: func(page InvitationsClientListByShareResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *InvitationsClientListByShareResponse) (InvitationsClientListByShareResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByShareCreateRequest(ctx, resourceGroupName, accountName, shareName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return InvitationsClientListByShareResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return InvitationsClientListByShareResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return InvitationsClientListByShareResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByShareHandleResponse(resp)
		},
	})
}

// listByShareCreateRequest creates the ListByShare request.
func (client *InvitationsClient) listByShareCreateRequest(ctx context.Context, resourceGroupName string, accountName string, shareName string, options *InvitationsClientListByShareOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataShare/accounts/{accountName}/shares/{shareName}/invitations"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if accountName == "" {
		return nil, errors.New("parameter accountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{accountName}", url.PathEscape(accountName))
	if shareName == "" {
		return nil, errors.New("parameter shareName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{shareName}", url.PathEscape(shareName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-09-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	if options != nil && options.Orderby != nil {
		reqQP.Set("$orderby", *options.Orderby)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByShareHandleResponse handles the ListByShare response.
func (client *InvitationsClient) listByShareHandleResponse(resp *http.Response) (InvitationsClientListByShareResponse, error) {
	result := InvitationsClientListByShareResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.InvitationList); err != nil {
		return InvitationsClientListByShareResponse{}, err
	}
	return result, nil
}
