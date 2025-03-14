// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armapplicationinsights

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

// ComponentQuotaStatusClient contains the methods for the ComponentQuotaStatus group.
// Don't use this type directly, use NewComponentQuotaStatusClient() instead.
type ComponentQuotaStatusClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewComponentQuotaStatusClient creates a new instance of ComponentQuotaStatusClient with the specified values.
//   - subscriptionID - The ID of the target subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewComponentQuotaStatusClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ComponentQuotaStatusClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &ComponentQuotaStatusClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// Get - Returns daily data volume cap (quota) status for an Application Insights component.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2015-05-01
//   - resourceGroupName - The name of the resource group. The name is case insensitive.
//   - resourceName - The name of the Application Insights component resource.
//   - options - ComponentQuotaStatusClientGetOptions contains the optional parameters for the ComponentQuotaStatusClient.Get
//     method.
func (client *ComponentQuotaStatusClient) Get(ctx context.Context, resourceGroupName string, resourceName string, options *ComponentQuotaStatusClientGetOptions) (ComponentQuotaStatusClientGetResponse, error) {
	var err error
	const operationName = "ComponentQuotaStatusClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, resourceGroupName, resourceName, options)
	if err != nil {
		return ComponentQuotaStatusClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ComponentQuotaStatusClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return ComponentQuotaStatusClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *ComponentQuotaStatusClient) getCreateRequest(ctx context.Context, resourceGroupName string, resourceName string, _ *ComponentQuotaStatusClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Insights/components/{resourceName}/quotastatus"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceName == "" {
		return nil, errors.New("parameter resourceName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceName}", url.PathEscape(resourceName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2015-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ComponentQuotaStatusClient) getHandleResponse(resp *http.Response) (ComponentQuotaStatusClientGetResponse, error) {
	result := ComponentQuotaStatusClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ComponentQuotaStatus); err != nil {
		return ComponentQuotaStatusClientGetResponse{}, err
	}
	return result, nil
}
