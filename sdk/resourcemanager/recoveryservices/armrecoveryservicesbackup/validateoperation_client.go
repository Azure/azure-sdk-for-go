//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armrecoveryservicesbackup

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

// ValidateOperationClient contains the methods for the ValidateOperation group.
// Don't use this type directly, use NewValidateOperationClient() instead.
type ValidateOperationClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewValidateOperationClient creates a new instance of ValidateOperationClient with the specified values.
//   - subscriptionID - The subscription Id.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewValidateOperationClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ValidateOperationClient, error) {
	cl, err := arm.NewClient(moduleName+".ValidateOperationClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ValidateOperationClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// BeginTrigger - Validate operation for specified backed up item in the form of an asynchronous operation. Returns tracking
// headers which can be tracked using GetValidateOperationResult API.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
//   - vaultName - The name of the recovery services vault.
//   - resourceGroupName - The name of the resource group where the recovery services vault is present.
//   - parameters - resource validate operation request
//   - options - ValidateOperationClientBeginTriggerOptions contains the optional parameters for the ValidateOperationClient.BeginTrigger
//     method.
func (client *ValidateOperationClient) BeginTrigger(ctx context.Context, vaultName string, resourceGroupName string, parameters ValidateOperationRequestClassification, options *ValidateOperationClientBeginTriggerOptions) (*runtime.Poller[ValidateOperationClientTriggerResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.trigger(ctx, vaultName, resourceGroupName, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[ValidateOperationClientTriggerResponse](resp, client.internal.Pipeline(), nil)
	} else {
		return runtime.NewPollerFromResumeToken[ValidateOperationClientTriggerResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Trigger - Validate operation for specified backed up item in the form of an asynchronous operation. Returns tracking headers
// which can be tracked using GetValidateOperationResult API.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-02-01
func (client *ValidateOperationClient) trigger(ctx context.Context, vaultName string, resourceGroupName string, parameters ValidateOperationRequestClassification, options *ValidateOperationClientBeginTriggerOptions) (*http.Response, error) {
	req, err := client.triggerCreateRequest(ctx, vaultName, resourceGroupName, parameters, options)
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

// triggerCreateRequest creates the Trigger request.
func (client *ValidateOperationClient) triggerCreateRequest(ctx context.Context, vaultName string, resourceGroupName string, parameters ValidateOperationRequestClassification, options *ValidateOperationClientBeginTriggerOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.RecoveryServices/vaults/{vaultName}/backupTriggerValidateOperation"
	if vaultName == "" {
		return nil, errors.New("parameter vaultName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vaultName}", url.PathEscape(vaultName))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}
