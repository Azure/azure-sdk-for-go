//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsql

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

// CapabilitiesClient contains the methods for the Capabilities group.
// Don't use this type directly, use NewCapabilitiesClient() instead.
type CapabilitiesClient struct {
	internal *arm.Client
	subscriptionID string
}

// NewCapabilitiesClient creates a new instance of CapabilitiesClient with the specified values.
//   - subscriptionID - The subscription ID that identifies an Azure subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewCapabilitiesClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*CapabilitiesClient, error) {
	cl, err := arm.NewClient(moduleName+".CapabilitiesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &CapabilitiesClient{
		subscriptionID: subscriptionID,
	internal: cl,
	}
	return client, nil
}

// ListByLocation - Gets the subscription capabilities available for the specified location.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-11-01-preview
//   - locationName - The location name whose capabilities are retrieved.
//   - options - CapabilitiesClientListByLocationOptions contains the optional parameters for the CapabilitiesClient.ListByLocation
//     method.
func (client *CapabilitiesClient) ListByLocation(ctx context.Context, locationName string, options *CapabilitiesClientListByLocationOptions) (CapabilitiesClientListByLocationResponse, error) {
	var err error
	req, err := client.listByLocationCreateRequest(ctx, locationName, options)
	if err != nil {
		return CapabilitiesClientListByLocationResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return CapabilitiesClientListByLocationResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return CapabilitiesClientListByLocationResponse{}, err
	}
	resp, err := client.listByLocationHandleResponse(httpResp)
	return resp, err
}

// listByLocationCreateRequest creates the ListByLocation request.
func (client *CapabilitiesClient) listByLocationCreateRequest(ctx context.Context, locationName string, options *CapabilitiesClientListByLocationOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Sql/locations/{locationName}/capabilities"
	if locationName == "" {
		return nil, errors.New("parameter locationName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{locationName}", url.PathEscape(locationName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Include != nil {
		reqQP.Set("include", string(*options.Include))
	}
	reqQP.Set("api-version", "2020-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByLocationHandleResponse handles the ListByLocation response.
func (client *CapabilitiesClient) listByLocationHandleResponse(resp *http.Response) (CapabilitiesClientListByLocationResponse, error) {
	result := CapabilitiesClientListByLocationResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.LocationCapabilities); err != nil {
		return CapabilitiesClientListByLocationResponse{}, err
	}
	return result, nil
}

