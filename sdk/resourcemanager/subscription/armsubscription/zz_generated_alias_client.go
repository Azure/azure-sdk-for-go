//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armsubscription

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// AliasClient contains the methods for the Alias group.
// Don't use this type directly, use NewAliasClient() instead.
type AliasClient struct {
	host string
	pl   runtime.Pipeline
}

// NewAliasClient creates a new instance of AliasClient with the specified values.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewAliasClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*AliasClient, error) {
	if options == nil {
		options = &arm.ClientOptions{}
	}
	ep := cloud.AzurePublic.Services[cloud.ResourceManager].Endpoint
	if c, ok := options.Cloud.Services[cloud.ResourceManager]; ok {
		ep = c.Endpoint
	}
	pl, err := armruntime.NewPipeline(moduleName, moduleVersion, credential, runtime.PipelineOptions{}, options)
	if err != nil {
		return nil, err
	}
	client := &AliasClient{
		host: ep,
		pl:   pl,
	}
	return client, nil
}

// BeginCreate - Create Alias Subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-10-01
// aliasName - AliasName is the name for the subscription creation request. Note that this is not the same as subscription
// name and this doesn’t have any other lifecycle need beyond the request for subscription
// creation.
// options - AliasClientBeginCreateOptions contains the optional parameters for the AliasClient.BeginCreate method.
func (client *AliasClient) BeginCreate(ctx context.Context, aliasName string, body PutAliasRequest, options *AliasClientBeginCreateOptions) (*runtime.Poller[AliasClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, aliasName, body, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[AliasClientCreateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[AliasClientCreateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Create - Create Alias Subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-10-01
func (client *AliasClient) create(ctx context.Context, aliasName string, body PutAliasRequest, options *AliasClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, aliasName, body, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *AliasClient) createCreateRequest(ctx context.Context, aliasName string, body PutAliasRequest, options *AliasClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Subscription/aliases/{aliasName}"
	if aliasName == "" {
		return nil, errors.New("parameter aliasName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{aliasName}", url.PathEscape(aliasName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, body)
}

// Delete - Delete Alias.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-10-01
// aliasName - AliasName is the name for the subscription creation request. Note that this is not the same as subscription
// name and this doesn’t have any other lifecycle need beyond the request for subscription
// creation.
// options - AliasClientDeleteOptions contains the optional parameters for the AliasClient.Delete method.
func (client *AliasClient) Delete(ctx context.Context, aliasName string, options *AliasClientDeleteOptions) (AliasClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, aliasName, options)
	if err != nil {
		return AliasClientDeleteResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AliasClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return AliasClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return AliasClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *AliasClient) deleteCreateRequest(ctx context.Context, aliasName string, options *AliasClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Subscription/aliases/{aliasName}"
	if aliasName == "" {
		return nil, errors.New("parameter aliasName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{aliasName}", url.PathEscape(aliasName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get Alias Subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-10-01
// aliasName - AliasName is the name for the subscription creation request. Note that this is not the same as subscription
// name and this doesn’t have any other lifecycle need beyond the request for subscription
// creation.
// options - AliasClientGetOptions contains the optional parameters for the AliasClient.Get method.
func (client *AliasClient) Get(ctx context.Context, aliasName string, options *AliasClientGetOptions) (AliasClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, aliasName, options)
	if err != nil {
		return AliasClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AliasClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AliasClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *AliasClient) getCreateRequest(ctx context.Context, aliasName string, options *AliasClientGetOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Subscription/aliases/{aliasName}"
	if aliasName == "" {
		return nil, errors.New("parameter aliasName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{aliasName}", url.PathEscape(aliasName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *AliasClient) getHandleResponse(resp *http.Response) (AliasClientGetResponse, error) {
	result := AliasClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AliasResponse); err != nil {
		return AliasClientGetResponse{}, err
	}
	return result, nil
}

// List - List Alias Subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2021-10-01
// options - AliasClientListOptions contains the optional parameters for the AliasClient.List method.
func (client *AliasClient) List(ctx context.Context, options *AliasClientListOptions) (AliasClientListResponse, error) {
	req, err := client.listCreateRequest(ctx, options)
	if err != nil {
		return AliasClientListResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return AliasClientListResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return AliasClientListResponse{}, runtime.NewResponseError(resp)
	}
	return client.listHandleResponse(resp)
}

// listCreateRequest creates the List request.
func (client *AliasClient) listCreateRequest(ctx context.Context, options *AliasClientListOptions) (*policy.Request, error) {
	urlPath := "/providers/Microsoft.Subscription/aliases"
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-10-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *AliasClient) listHandleResponse(resp *http.Response) (AliasClientListResponse, error) {
	result := AliasClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AliasListResult); err != nil {
		return AliasClientListResponse{}, err
	}
	return result, nil
}
