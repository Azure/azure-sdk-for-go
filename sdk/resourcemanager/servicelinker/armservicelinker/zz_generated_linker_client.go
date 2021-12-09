//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armservicelinker

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	armruntime "github.com/Azure/azure-sdk-for-go/sdk/azcore/arm/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// LinkerClient contains the methods for the Linker group.
// Don't use this type directly, use NewLinkerClient() instead.
type LinkerClient struct {
	ep string
	pl runtime.Pipeline
}

// NewLinkerClient creates a new instance of LinkerClient with the specified values.
func NewLinkerClient(credential azcore.TokenCredential, options *arm.ClientOptions) *LinkerClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &LinkerClient{ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// BeginCreateOrUpdate - Create or update linker resource.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) BeginCreateOrUpdate(ctx context.Context, resourceURI string, linkerName string, parameters LinkerResource, options *LinkerBeginCreateOrUpdateOptions) (LinkerCreateOrUpdatePollerResponse, error) {
	resp, err := client.createOrUpdate(ctx, resourceURI, linkerName, parameters, options)
	if err != nil {
		return LinkerCreateOrUpdatePollerResponse{}, err
	}
	result := LinkerCreateOrUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("LinkerClient.CreateOrUpdate", "azure-async-operation", resp, client.pl, client.createOrUpdateHandleError)
	if err != nil {
		return LinkerCreateOrUpdatePollerResponse{}, err
	}
	result.Poller = &LinkerCreateOrUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// CreateOrUpdate - Create or update linker resource.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) createOrUpdate(ctx context.Context, resourceURI string, linkerName string, parameters LinkerResource, options *LinkerBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceURI, linkerName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *LinkerClient) createOrUpdateCreateRequest(ctx context.Context, resourceURI string, linkerName string, parameters LinkerResource, options *LinkerBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.ServiceLinker/linkers/{linkerName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	if linkerName == "" {
		return nil, errors.New("parameter linkerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{linkerName}", url.PathEscape(linkerName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *LinkerClient) createOrUpdateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginDelete - Delete a link.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) BeginDelete(ctx context.Context, resourceURI string, linkerName string, options *LinkerBeginDeleteOptions) (LinkerDeletePollerResponse, error) {
	resp, err := client.deleteOperation(ctx, resourceURI, linkerName, options)
	if err != nil {
		return LinkerDeletePollerResponse{}, err
	}
	result := LinkerDeletePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("LinkerClient.Delete", "azure-async-operation", resp, client.pl, client.deleteHandleError)
	if err != nil {
		return LinkerDeletePollerResponse{}, err
	}
	result.Poller = &LinkerDeletePoller{
		pt: pt,
	}
	return result, nil
}

// Delete - Delete a link.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) deleteOperation(ctx context.Context, resourceURI string, linkerName string, options *LinkerBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceURI, linkerName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *LinkerClient) deleteCreateRequest(ctx context.Context, resourceURI string, linkerName string, options *LinkerBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.ServiceLinker/linkers/{linkerName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	if linkerName == "" {
		return nil, errors.New("parameter linkerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{linkerName}", url.PathEscape(linkerName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// deleteHandleError handles the Delete error response.
func (client *LinkerClient) deleteHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// Get - Returns Linker resource for a given name.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) Get(ctx context.Context, resourceURI string, linkerName string, options *LinkerGetOptions) (LinkerGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceURI, linkerName, options)
	if err != nil {
		return LinkerGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return LinkerGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return LinkerGetResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *LinkerClient) getCreateRequest(ctx context.Context, resourceURI string, linkerName string, options *LinkerGetOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.ServiceLinker/linkers/{linkerName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	if linkerName == "" {
		return nil, errors.New("parameter linkerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{linkerName}", url.PathEscape(linkerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *LinkerClient) getHandleResponse(resp *http.Response) (LinkerGetResponse, error) {
	result := LinkerGetResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.LinkerResource); err != nil {
		return LinkerGetResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// getHandleError handles the Get error response.
func (client *LinkerClient) getHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// List - Returns list of Linkers which connects to the resource.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) List(resourceURI string, options *LinkerListOptions) *LinkerListPager {
	return &LinkerListPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.listCreateRequest(ctx, resourceURI, options)
		},
		advancer: func(ctx context.Context, resp LinkerListResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.LinkerList.NextLink)
		},
	}
}

// listCreateRequest creates the List request.
func (client *LinkerClient) listCreateRequest(ctx context.Context, resourceURI string, options *LinkerListOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.ServiceLinker/linkers"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listHandleResponse handles the List response.
func (client *LinkerClient) listHandleResponse(resp *http.Response) (LinkerListResponse, error) {
	result := LinkerListResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.LinkerList); err != nil {
		return LinkerListResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// listHandleError handles the List error response.
func (client *LinkerClient) listHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// ListConfigurations - list source configurations for a linker.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) ListConfigurations(ctx context.Context, resourceURI string, linkerName string, options *LinkerListConfigurationsOptions) (LinkerListConfigurationsResponse, error) {
	req, err := client.listConfigurationsCreateRequest(ctx, resourceURI, linkerName, options)
	if err != nil {
		return LinkerListConfigurationsResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return LinkerListConfigurationsResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return LinkerListConfigurationsResponse{}, client.listConfigurationsHandleError(resp)
	}
	return client.listConfigurationsHandleResponse(resp)
}

// listConfigurationsCreateRequest creates the ListConfigurations request.
func (client *LinkerClient) listConfigurationsCreateRequest(ctx context.Context, resourceURI string, linkerName string, options *LinkerListConfigurationsOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.ServiceLinker/linkers/{linkerName}/listConfigurations"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	if linkerName == "" {
		return nil, errors.New("parameter linkerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{linkerName}", url.PathEscape(linkerName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// listConfigurationsHandleResponse handles the ListConfigurations response.
func (client *LinkerClient) listConfigurationsHandleResponse(resp *http.Response) (LinkerListConfigurationsResponse, error) {
	result := LinkerListConfigurationsResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.SourceConfigurationResult); err != nil {
		return LinkerListConfigurationsResponse{}, runtime.NewResponseError(err, resp)
	}
	return result, nil
}

// listConfigurationsHandleError handles the ListConfigurations error response.
func (client *LinkerClient) listConfigurationsHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginUpdate - Operation to update an existing link.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) BeginUpdate(ctx context.Context, resourceURI string, linkerName string, parameters LinkerPatch, options *LinkerBeginUpdateOptions) (LinkerUpdatePollerResponse, error) {
	resp, err := client.update(ctx, resourceURI, linkerName, parameters, options)
	if err != nil {
		return LinkerUpdatePollerResponse{}, err
	}
	result := LinkerUpdatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("LinkerClient.Update", "azure-async-operation", resp, client.pl, client.updateHandleError)
	if err != nil {
		return LinkerUpdatePollerResponse{}, err
	}
	result.Poller = &LinkerUpdatePoller{
		pt: pt,
	}
	return result, nil
}

// Update - Operation to update an existing link.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) update(ctx context.Context, resourceURI string, linkerName string, parameters LinkerPatch, options *LinkerBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceURI, linkerName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated) {
		return nil, client.updateHandleError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *LinkerClient) updateCreateRequest(ctx context.Context, resourceURI string, linkerName string, parameters LinkerPatch, options *LinkerBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.ServiceLinker/linkers/{linkerName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	if linkerName == "" {
		return nil, errors.New("parameter linkerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{linkerName}", url.PathEscape(linkerName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateHandleError handles the Update error response.
func (client *LinkerClient) updateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginValidate - Validate a link.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) BeginValidate(ctx context.Context, resourceURI string, linkerName string, options *LinkerBeginValidateOptions) (LinkerValidatePollerResponse, error) {
	resp, err := client.validate(ctx, resourceURI, linkerName, options)
	if err != nil {
		return LinkerValidatePollerResponse{}, err
	}
	result := LinkerValidatePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("LinkerClient.Validate", "azure-async-operation", resp, client.pl, client.validateHandleError)
	if err != nil {
		return LinkerValidatePollerResponse{}, err
	}
	result.Poller = &LinkerValidatePoller{
		pt: pt,
	}
	return result, nil
}

// Validate - Validate a link.
// If the operation fails it returns the *ErrorResponse error type.
func (client *LinkerClient) validate(ctx context.Context, resourceURI string, linkerName string, options *LinkerBeginValidateOptions) (*http.Response, error) {
	req, err := client.validateCreateRequest(ctx, resourceURI, linkerName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.validateHandleError(resp)
	}
	return resp, nil
}

// validateCreateRequest creates the Validate request.
func (client *LinkerClient) validateCreateRequest(ctx context.Context, resourceURI string, linkerName string, options *LinkerBeginValidateOptions) (*policy.Request, error) {
	urlPath := "/{resourceUri}/providers/Microsoft.ServiceLinker/linkers/{linkerName}/validateLinker"
	urlPath = strings.ReplaceAll(urlPath, "{resourceUri}", resourceURI)
	if linkerName == "" {
		return nil, errors.New("parameter linkerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{linkerName}", url.PathEscape(linkerName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-11-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// validateHandleError handles the Validate error response.
func (client *LinkerClient) validateHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := ErrorResponse{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
