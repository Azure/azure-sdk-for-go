//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdomainservices

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
)

// OuContainerOperationsClient contains the methods for the OuContainerOperations group.
// Don't use this type directly, use NewOuContainerOperationsClient() instead.
type OuContainerOperationsClient struct {
	host string
	pl   runtime.Pipeline
}

// NewOuContainerOperationsClient creates a new instance of OuContainerOperationsClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewOuContainerOperationsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*OuContainerOperationsClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublicCloud.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &OuContainerOperationsClient{
		host: ep,
		pl:   pl,
	}
	return client, nil
}

// NewListPager - Lists all the available OuContainer operations.
// If the operation fails it returns an *azcore.ResponseError type.
// options - OuContainerOperationsClientListOptions contains the optional parameters for the OuContainerOperationsClient.List
// method.
func (client *OuContainerOperationsClient) NewListPager(options *OuContainerOperationsClientListOptions) *runtime.Pager[OuContainerOperationsClientListResponse] {
	return runtime.NewPager(runtime.PageProcessor[OuContainerOperationsClientListResponse]{
		More: func(page OuContainerOperationsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *OuContainerOperationsClientListResponse) (OuContainerOperationsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return OuContainerOperationsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return OuContainerOperationsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return OuContainerOperationsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *OuContainerOperationsClient) listCreateRequest(ctx context.Context, options *OuContainerOperationsClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Aad/operations"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *OuContainerOperationsClient) listHandleResponse(resp *http.Response) (OuContainerOperationsClientListResponse, error) {
	result := OuContainerOperationsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.OperationEntityListResult); err != nil {
		return OuContainerOperationsClientListResponse{}, err
	}
	return result, nil
}
