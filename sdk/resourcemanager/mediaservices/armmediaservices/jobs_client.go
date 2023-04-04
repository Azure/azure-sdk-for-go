//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.
// DO NOT EDIT.

package armmediaservices

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

// JobsClient contains the methods for the Jobs group.
// Don't use this type directly, use NewJobsClient() instead.
type JobsClient struct {
	internal       *arm.Client
	subscriptionID string
}

// NewJobsClient creates a new instance of JobsClient with the specified values.
//   - subscriptionID - The unique identifier for a Microsoft Azure subscription.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewJobsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*JobsClient, error) {
	cl, err := arm.NewClient(moduleName+".JobsClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &JobsClient{
		subscriptionID: subscriptionID,
		internal:       cl,
	}
	return client, nil
}

// CancelJob - Cancel a Job.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-07-01
//   - resourceGroupName - The name of the resource group within the Azure subscription.
//   - accountName - The Media Services account name.
//   - transformName - The Transform name.
//   - jobName - The Job name.
//   - options - JobsClientCancelJobOptions contains the optional parameters for the JobsClient.CancelJob method.
func (client *JobsClient) CancelJob(ctx context.Context, resourceGroupName string, accountName string, transformName string, jobName string, options *JobsClientCancelJobOptions) (JobsClientCancelJobResponse, error) {
	req, err := client.cancelJobCreateRequest(ctx, resourceGroupName, accountName, transformName, jobName, options)
	if err != nil {
		return JobsClientCancelJobResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return JobsClientCancelJobResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return JobsClientCancelJobResponse{}, runtime.NewResponseError(resp)
	}
	return JobsClientCancelJobResponse{}, nil
}

// cancelJobCreateRequest creates the CancelJob request.
func (client *JobsClient) cancelJobCreateRequest(ctx context.Context, resourceGroupName string, accountName string, transformName string, jobName string, options *JobsClientCancelJobOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/transforms/{transformName}/jobs/{jobName}/cancelJob"
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
	if transformName == "" {
		return nil, errors.New("parameter transformName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{transformName}", url.PathEscape(transformName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Create - Creates a Job.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-07-01
//   - resourceGroupName - The name of the resource group within the Azure subscription.
//   - accountName - The Media Services account name.
//   - transformName - The Transform name.
//   - jobName - The Job name.
//   - parameters - The request parameters
//   - options - JobsClientCreateOptions contains the optional parameters for the JobsClient.Create method.
func (client *JobsClient) Create(ctx context.Context, resourceGroupName string, accountName string, transformName string, jobName string, parameters Job, options *JobsClientCreateOptions) (JobsClientCreateResponse, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, accountName, transformName, jobName, parameters, options)
	if err != nil {
		return JobsClientCreateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return JobsClientCreateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusCreated) {
		return JobsClientCreateResponse{}, runtime.NewResponseError(resp)
	}
	return client.createHandleResponse(resp)
}

// createCreateRequest creates the Create request.
func (client *JobsClient) createCreateRequest(ctx context.Context, resourceGroupName string, accountName string, transformName string, jobName string, parameters Job, options *JobsClientCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/transforms/{transformName}/jobs/{jobName}"
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
	if transformName == "" {
		return nil, errors.New("parameter transformName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{transformName}", url.PathEscape(transformName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// createHandleResponse handles the Create response.
func (client *JobsClient) createHandleResponse(resp *http.Response) (JobsClientCreateResponse, error) {
	result := JobsClientCreateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Job); err != nil {
		return JobsClientCreateResponse{}, err
	}
	return result, nil
}

// Delete - Deletes a Job.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-07-01
//   - resourceGroupName - The name of the resource group within the Azure subscription.
//   - accountName - The Media Services account name.
//   - transformName - The Transform name.
//   - jobName - The Job name.
//   - options - JobsClientDeleteOptions contains the optional parameters for the JobsClient.Delete method.
func (client *JobsClient) Delete(ctx context.Context, resourceGroupName string, accountName string, transformName string, jobName string, options *JobsClientDeleteOptions) (JobsClientDeleteResponse, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, accountName, transformName, jobName, options)
	if err != nil {
		return JobsClientDeleteResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return JobsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusNoContent) {
		return JobsClientDeleteResponse{}, runtime.NewResponseError(resp)
	}
	return JobsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *JobsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, accountName string, transformName string, jobName string, options *JobsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/transforms/{transformName}/jobs/{jobName}"
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
	if transformName == "" {
		return nil, errors.New("parameter transformName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{transformName}", url.PathEscape(transformName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets a Job.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-07-01
//   - resourceGroupName - The name of the resource group within the Azure subscription.
//   - accountName - The Media Services account name.
//   - transformName - The Transform name.
//   - jobName - The Job name.
//   - options - JobsClientGetOptions contains the optional parameters for the JobsClient.Get method.
func (client *JobsClient) Get(ctx context.Context, resourceGroupName string, accountName string, transformName string, jobName string, options *JobsClientGetOptions) (JobsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, accountName, transformName, jobName, options)
	if err != nil {
		return JobsClientGetResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return JobsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return JobsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *JobsClient) getCreateRequest(ctx context.Context, resourceGroupName string, accountName string, transformName string, jobName string, options *JobsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/transforms/{transformName}/jobs/{jobName}"
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
	if transformName == "" {
		return nil, errors.New("parameter transformName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{transformName}", url.PathEscape(transformName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *JobsClient) getHandleResponse(resp *http.Response) (JobsClientGetResponse, error) {
	result := JobsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Job); err != nil {
		return JobsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all of the Jobs for the Transform.
//
// Generated from API version 2022-07-01
//   - resourceGroupName - The name of the resource group within the Azure subscription.
//   - accountName - The Media Services account name.
//   - transformName - The Transform name.
//   - options - JobsClientListOptions contains the optional parameters for the JobsClient.NewListPager method.
func (client *JobsClient) NewListPager(resourceGroupName string, accountName string, transformName string, options *JobsClientListOptions) *runtime.Pager[JobsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[JobsClientListResponse]{
		More: func(page JobsClientListResponse) bool {
			return page.ODataNextLink != nil && len(*page.ODataNextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *JobsClientListResponse) (JobsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, resourceGroupName, accountName, transformName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.ODataNextLink)
			}
			if err != nil {
				return JobsClientListResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return JobsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return JobsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *JobsClient) listCreateRequest(ctx context.Context, resourceGroupName string, accountName string, transformName string, options *JobsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/transforms/{transformName}/jobs"
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
	if transformName == "" {
		return nil, errors.New("parameter transformName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{transformName}", url.PathEscape(transformName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-07-01")
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

// listHandleResponse handles the List response.
func (client *JobsClient) listHandleResponse(resp *http.Response) (JobsClientListResponse, error) {
	result := JobsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.JobCollection); err != nil {
		return JobsClientListResponse{}, err
	}
	return result, nil
}

// Update - Update is only supported for description and priority. Updating Priority will take effect when the Job state is
// Queued or Scheduled and depending on the timing the priority update may be ignored.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2022-07-01
//   - resourceGroupName - The name of the resource group within the Azure subscription.
//   - accountName - The Media Services account name.
//   - transformName - The Transform name.
//   - jobName - The Job name.
//   - parameters - The request parameters
//   - options - JobsClientUpdateOptions contains the optional parameters for the JobsClient.Update method.
func (client *JobsClient) Update(ctx context.Context, resourceGroupName string, accountName string, transformName string, jobName string, parameters Job, options *JobsClientUpdateOptions) (JobsClientUpdateResponse, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, accountName, transformName, jobName, parameters, options)
	if err != nil {
		return JobsClientUpdateResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return JobsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return JobsClientUpdateResponse{}, runtime.NewResponseError(resp)
	}
	return client.updateHandleResponse(resp)
}

// updateCreateRequest creates the Update request.
func (client *JobsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, accountName string, transformName string, jobName string, parameters Job, options *JobsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Media/mediaServices/{accountName}/transforms/{transformName}/jobs/{jobName}"
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
	if transformName == "" {
		return nil, errors.New("parameter transformName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{transformName}", url.PathEscape(transformName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-07-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// updateHandleResponse handles the Update response.
func (client *JobsClient) updateHandleResponse(resp *http.Response) (JobsClientUpdateResponse, error) {
	result := JobsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.Job); err != nil {
		return JobsClientUpdateResponse{}, err
	}
	return result, nil
}
