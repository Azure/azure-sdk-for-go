//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armtestbase

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
	"strings"
)

// UsageClient contains the methods for the Usage group.
// Don't use this type directly, use NewUsageClient() instead.
type UsageClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewUsageClient creates a new instance of UsageClient with the specified values.
// subscriptionID - The Azure subscription ID. This is a GUID-formatted string.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewUsageClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *UsageClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Endpoint) == 0 {
		cp.Endpoint = arm.AzurePublicCloud
	}
	client := &UsageClient{
		subscriptionID: subscriptionID,
		host:           string(cp.Endpoint),
		pl:             armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, &cp),
	}
	return client
}

// List - Lists the usage data of a Test Base Account.
// If the operation fails it returns an *azcore.ResponseError type.
// resourceGroupName - The name of the resource group that contains the resource.
// testBaseAccountName - The resource name of the Test Base Account.
// options - UsageClientListOptions contains the optional parameters for the UsageClient.List method.
func (client *UsageClient) List(resourceGroupName string, testBaseAccountName string, options *UsageClientListOptions) *UsageClientListPager {
	return &UsageClientListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceGroupName, testBaseAccountName, options)
		},
		advancer: func(ctx context.Context, resp UsageClientListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.AccountUsageDataList.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *UsageClient) listCreateRequest(ctx context.Context, resourceGroupName string, testBaseAccountName string, options *UsageClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.TestBase/testBaseAccounts/{testBaseAccountName}/usages"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if testBaseAccountName == "" {
		return nil, errors.New("parameter testBaseAccountName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{testBaseAccountName}", url.PathEscape(testBaseAccountName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	if options != nil && options.Filter != nil {
		reqQP.Set("$filter", *options.Filter)
	}
	reqQP.Set("api-version", "2020-12-16-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *UsageClient) listHandleResponse(resp *http.Response) (UsageClientListResponse, error) {
	result := UsageClientListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.AccountUsageDataList); err != nil {
		return UsageClientListResponse{}, err
	}
	return result, nil
}
