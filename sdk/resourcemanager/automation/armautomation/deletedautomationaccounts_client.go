//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armautomation

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

// DeletedAutomationAccountsClient contains the methods for the DeletedAutomationAccounts group.
// Don't use this type directly, use NewDeletedAutomationAccountsClient() instead.
type DeletedAutomationAccountsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewDeletedAutomationAccountsClient creates a new instance of DeletedAutomationAccountsClient with the specified values.
//   - subscriptionID - Gets subscription credentials which uniquely identify Microsoft Azure subscription. The subscription ID
//     forms part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewDeletedAutomationAccountsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*DeletedAutomationAccountsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &DeletedAutomationAccountsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// ListBySubscription - Retrieve deleted automation account.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-01-31
//   - options - DeletedAutomationAccountsClientListBySubscriptionOptions contains the optional parameters for the DeletedAutomationAccountsClient.ListBySubscription
//     method.
func (client *DeletedAutomationAccountsClient) ListBySubscription(ctx context.Context, options *DeletedAutomationAccountsClientListBySubscriptionOptions) (DeletedAutomationAccountsClientListBySubscriptionResponse, error) {
	var err error
	const operationName = "DeletedAutomationAccountsClient.ListBySubscription"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.listBySubscriptionCreateRequest(ctx, options)
	if err != nil {
		return DeletedAutomationAccountsClientListBySubscriptionResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return DeletedAutomationAccountsClientListBySubscriptionResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return DeletedAutomationAccountsClientListBySubscriptionResponse{}, err
	}
	resp, err := client.listBySubscriptionHandleResponse(httpResp)
	return resp, err
}

// listBySubscriptionCreateRequest creates the ListBySubscription request.
func (client *DeletedAutomationAccountsClient) listBySubscriptionCreateRequest(ctx context.Context, options *DeletedAutomationAccountsClientListBySubscriptionOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Automation/deletedAutomationAccounts"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-01-31")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listBySubscriptionHandleResponse handles the ListBySubscription response.
func (client *DeletedAutomationAccountsClient) listBySubscriptionHandleResponse(resp *http.Response) (DeletedAutomationAccountsClientListBySubscriptionResponse, error) {
	result := DeletedAutomationAccountsClientListBySubscriptionResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DeletedAutomationAccountListResult); err != nil {
		return DeletedAutomationAccountsClientListBySubscriptionResponse{}, err
	}
	return result, nil
}
