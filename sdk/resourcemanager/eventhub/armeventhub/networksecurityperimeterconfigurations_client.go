//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armeventhub

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

// NetworkSecurityPerimeterConfigurationsClient contains the methods for the NetworkSecurityPerimeterConfigurations group.
// Don't use this type directly, use NewNetworkSecurityPerimeterConfigurationsClient() instead.
type NetworkSecurityPerimeterConfigurationsClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewNetworkSecurityPerimeterConfigurationsClient creates a new instance of NetworkSecurityPerimeterConfigurationsClient with the specified values.
//   - subscriptionID - Subscription credentials that uniquely identify a Microsoft Azure subscription. The subscription ID forms
//     part of the URI for every service call.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewNetworkSecurityPerimeterConfigurationsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*NetworkSecurityPerimeterConfigurationsClient, error) {
	cl, err := arm.NewClient(moduleName+".NetworkSecurityPerimeterConfigurationsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &NetworkSecurityPerimeterConfigurationsClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Refreshes any information about the association.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-01-preview
//   - resourceGroupName - Name of the resource group within the azure subscription.
//   - namespaceName - The Namespace name
//   - resourceAssociationName - The ResourceAssociation Name
//   - options - NetworkSecurityPerimeterConfigurationsClientBeginCreateOrUpdateOptions contains the optional parameters for the
//     NetworkSecurityPerimeterConfigurationsClient.BeginCreateOrUpdate method.
func (client *NetworkSecurityPerimeterConfigurationsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, resourceAssociationName string, options *NetworkSecurityPerimeterConfigurationsClientBeginCreateOrUpdateOptions) (*runtime.Poller[NetworkSecurityPerimeterConfigurationsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, namespaceName, resourceAssociationName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller[NetworkSecurityPerimeterConfigurationsClientCreateOrUpdateResponse](resp, client.internal.Pipeline(), nil)
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[NetworkSecurityPerimeterConfigurationsClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Refreshes any information about the association.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-10-01-preview
func (client *NetworkSecurityPerimeterConfigurationsClient) createOrUpdate(ctx context.Context, resourceGroupName string, namespaceName string, resourceAssociationName string, options *NetworkSecurityPerimeterConfigurationsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, namespaceName, resourceAssociationName, options)
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

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *NetworkSecurityPerimeterConfigurationsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, namespaceName string, resourceAssociationName string, options *NetworkSecurityPerimeterConfigurationsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.EventHub/namespaces/{namespaceName}/networkSecurityPerimeterConfigurations/{resourceAssociationName}/reconcile"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if namespaceName == "" {
		return nil, errors.New("parameter namespaceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{namespaceName}", url.PathEscape(namespaceName))
	if resourceAssociationName == "" {
		return nil, errors.New("parameter resourceAssociationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceAssociationName}", url.PathEscape(resourceAssociationName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

