//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armdatabox

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

// JobsClient contains the methods for the Jobs group.
// Don't use this type directly, use NewJobsClient() instead.
type JobsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewJobsClient creates a new instance of JobsClient with the specified values.
// subscriptionID - The Subscription Id
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewJobsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*JobsClient, error) {
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
	client := &JobsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// BookShipmentPickUp - Book shipment pick up.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// resourceGroupName - The Resource Group Name
// jobName - The name of the job Resource within the specified resource group. job names must be between 3 and 24 characters
// in length and use any alphanumeric and underscore only
// shipmentPickUpRequest - Details of shipment pick up request.
// options - JobsClientBookShipmentPickUpOptions contains the optional parameters for the JobsClient.BookShipmentPickUp method.
func (client *JobsClient) BookShipmentPickUp(ctx context.Context, resourceGroupName string, jobName string, shipmentPickUpRequest ShipmentPickUpRequest, options *JobsClientBookShipmentPickUpOptions) (JobsClientBookShipmentPickUpResponse, error) {
	req, err := client.bookShipmentPickUpCreateRequest(ctx, resourceGroupName, jobName, shipmentPickUpRequest, options)
	if err != nil {
		return JobsClientBookShipmentPickUpResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return JobsClientBookShipmentPickUpResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return JobsClientBookShipmentPickUpResponse{}, runtime.NewResponseError(resp)
	}
	return client.bookShipmentPickUpHandleResponse(resp)
}

// bookShipmentPickUpCreateRequest creates the BookShipmentPickUp request.
func (client *JobsClient) bookShipmentPickUpCreateRequest(ctx context.Context, resourceGroupName string, jobName string, shipmentPickUpRequest ShipmentPickUpRequest, options *JobsClientBookShipmentPickUpOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataBox/jobs/{jobName}/bookShipmentPickUp"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, shipmentPickUpRequest)
}

// bookShipmentPickUpHandleResponse handles the BookShipmentPickUp response.
func (client *JobsClient) bookShipmentPickUpHandleResponse(resp *http.Response) (JobsClientBookShipmentPickUpResponse, error) {
	result := JobsClientBookShipmentPickUpResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ShipmentPickUpResponse); err != nil {
		return JobsClientBookShipmentPickUpResponse{}, err
	}
	return result, nil
}

// Cancel - CancelJob.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// resourceGroupName - The Resource Group Name
// jobName - The name of the job Resource within the specified resource group. job names must be between 3 and 24 characters
// in length and use any alphanumeric and underscore only
// cancellationReason - Reason for cancellation.
// options - JobsClientCancelOptions contains the optional parameters for the JobsClient.Cancel method.
func (client *JobsClient) Cancel(ctx context.Context, resourceGroupName string, jobName string, cancellationReason CancellationReason, options *JobsClientCancelOptions) (JobsClientCancelResponse, error) {
	req, err := client.cancelCreateRequest(ctx, resourceGroupName, jobName, cancellationReason, options)
	if err != nil {
		return JobsClientCancelResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return JobsClientCancelResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusNoContent) {
		return JobsClientCancelResponse{}, runtime.NewResponseError(resp)
	}
	return JobsClientCancelResponse{}, nil
}

// cancelCreateRequest creates the Cancel request.
func (client *JobsClient) cancelCreateRequest(ctx context.Context, resourceGroupName string, jobName string, cancellationReason CancellationReason, options *JobsClientCancelOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataBox/jobs/{jobName}/cancel"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, cancellationReason)
}

// BeginCreate - Creates a new job with the specified parameters. Existing job cannot be updated with this API and should
// instead be updated with the Update job API.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// resourceGroupName - The Resource Group Name
// jobName - The name of the job Resource within the specified resource group. job names must be between 3 and 24 characters
// in length and use any alphanumeric and underscore only
// jobResource - Job details from request body.
// options - JobsClientBeginCreateOptions contains the optional parameters for the JobsClient.BeginCreate method.
func (client *JobsClient) BeginCreate(ctx context.Context, resourceGroupName string, jobName string, jobResource JobResource, options *JobsClientBeginCreateOptions) (*runtime.Poller[JobsClientCreateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.create(ctx, resourceGroupName, jobName, jobResource, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[JobsClientCreateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[JobsClientCreateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Create - Creates a new job with the specified parameters. Existing job cannot be updated with this API and should instead
// be updated with the Update job API.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
func (client *JobsClient) create(ctx context.Context, resourceGroupName string, jobName string, jobResource JobResource, options *JobsClientBeginCreateOptions) (*http.Response, error) {
	req, err := client.createCreateRequest(ctx, resourceGroupName, jobName, jobResource, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createCreateRequest creates the Create request.
func (client *JobsClient) createCreateRequest(ctx context.Context, resourceGroupName string, jobName string, jobResource JobResource, options *JobsClientBeginCreateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataBox/jobs/{jobName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, jobResource)
}

// BeginDelete - Deletes a job.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// resourceGroupName - The Resource Group Name
// jobName - The name of the job Resource within the specified resource group. job names must be between 3 and 24 characters
// in length and use any alphanumeric and underscore only
// options - JobsClientBeginDeleteOptions contains the optional parameters for the JobsClient.BeginDelete method.
func (client *JobsClient) BeginDelete(ctx context.Context, resourceGroupName string, jobName string, options *JobsClientBeginDeleteOptions) (*runtime.Poller[JobsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, jobName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[JobsClientDeleteResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[JobsClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Deletes a job.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
func (client *JobsClient) deleteOperation(ctx context.Context, resourceGroupName string, jobName string, options *JobsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, jobName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *JobsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, jobName string, options *JobsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataBox/jobs/{jobName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets information about the specified job.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// resourceGroupName - The Resource Group Name
// jobName - The name of the job Resource within the specified resource group. job names must be between 3 and 24 characters
// in length and use any alphanumeric and underscore only
// options - JobsClientGetOptions contains the optional parameters for the JobsClient.Get method.
func (client *JobsClient) Get(ctx context.Context, resourceGroupName string, jobName string, options *JobsClientGetOptions) (JobsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, jobName, options)
	if err != nil {
		return JobsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return JobsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return JobsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *JobsClient) getCreateRequest(ctx context.Context, resourceGroupName string, jobName string, options *JobsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataBox/jobs/{jobName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	if options != nil && options.Expand != nil {
		reqQP.Set("$expand", *options.Expand)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *JobsClient) getHandleResponse(resp *http.Response) (JobsClientGetResponse, error) {
	result := JobsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.JobResource); err != nil {
		return JobsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - Lists all the jobs available under the subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// options - JobsClientListOptions contains the optional parameters for the JobsClient.List method.
func (client *JobsClient) NewListPager(options *JobsClientListOptions) *runtime.Pager[JobsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[JobsClientListResponse]{
		More: func(page JobsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *JobsClientListResponse) (JobsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return JobsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
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
func (client *JobsClient) listCreateRequest(ctx context.Context, options *JobsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.DataBox/jobs"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *JobsClient) listHandleResponse(resp *http.Response) (JobsClientListResponse, error) {
	result := JobsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.JobResourceList); err != nil {
		return JobsClientListResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - Lists all the jobs available under the given resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// resourceGroupName - The Resource Group Name
// options - JobsClientListByResourceGroupOptions contains the optional parameters for the JobsClient.ListByResourceGroup
// method.
func (client *JobsClient) NewListByResourceGroupPager(resourceGroupName string, options *JobsClientListByResourceGroupOptions) *runtime.Pager[JobsClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[JobsClientListByResourceGroupResponse]{
		More: func(page JobsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *JobsClientListByResourceGroupResponse) (JobsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return JobsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return JobsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return JobsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *JobsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *JobsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataBox/jobs"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	if options != nil && options.SkipToken != nil {
		reqQP.Set("$skipToken", *options.SkipToken)
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *JobsClient) listByResourceGroupHandleResponse(resp *http.Response) (JobsClientListByResourceGroupResponse, error) {
	result := JobsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.JobResourceList); err != nil {
		return JobsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// NewListCredentialsPager - This method gets the unencrypted secrets related to the job.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// resourceGroupName - The Resource Group Name
// jobName - The name of the job Resource within the specified resource group. job names must be between 3 and 24 characters
// in length and use any alphanumeric and underscore only
// options - JobsClientListCredentialsOptions contains the optional parameters for the JobsClient.ListCredentials method.
func (client *JobsClient) NewListCredentialsPager(resourceGroupName string, jobName string, options *JobsClientListCredentialsOptions) *runtime.Pager[JobsClientListCredentialsResponse] {
	return runtime.NewPager(runtime.PagingHandler[JobsClientListCredentialsResponse]{
		More: func(page JobsClientListCredentialsResponse) bool {
			return false
		},
		Fetcher: func(ctx context.Context, page *JobsClientListCredentialsResponse) (JobsClientListCredentialsResponse, error) {
			req, err := client.listCredentialsCreateRequest(ctx, resourceGroupName, jobName, options)
			if err != nil {
				return JobsClientListCredentialsResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return JobsClientListCredentialsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return JobsClientListCredentialsResponse{}, runtime.NewResponseError(resp)
			}
			return client.listCredentialsHandleResponse(resp)
		},
	})
}

// listCredentialsCreateRequest creates the ListCredentials request.
func (client *JobsClient) listCredentialsCreateRequest(ctx context.Context, resourceGroupName string, jobName string, options *JobsClientListCredentialsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataBox/jobs/{jobName}/listCredentials"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listCredentialsHandleResponse handles the ListCredentials response.
func (client *JobsClient) listCredentialsHandleResponse(resp *http.Response) (JobsClientListCredentialsResponse, error) {
	result := JobsClientListCredentialsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.UnencryptedCredentialsList); err != nil {
		return JobsClientListCredentialsResponse{}, err
	}
	return result, nil
}

// MarkDevicesShipped - Request to mark devices for a given job as shipped
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// jobName - The name of the job Resource within the specified resource group. job names must be between 3 and 24 characters
// in length and use any alphanumeric and underscore only
// resourceGroupName - The Resource Group Name
// markDevicesShippedRequest - Mark Devices Shipped Request
// options - JobsClientMarkDevicesShippedOptions contains the optional parameters for the JobsClient.MarkDevicesShipped method.
func (client *JobsClient) MarkDevicesShipped(ctx context.Context, jobName string, resourceGroupName string, markDevicesShippedRequest MarkDevicesShippedRequest, options *JobsClientMarkDevicesShippedOptions) (JobsClientMarkDevicesShippedResponse, error) {
	req, err := client.markDevicesShippedCreateRequest(ctx, jobName, resourceGroupName, markDevicesShippedRequest, options)
	if err != nil {
		return JobsClientMarkDevicesShippedResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return JobsClientMarkDevicesShippedResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusNoContent) {
		return JobsClientMarkDevicesShippedResponse{}, runtime.NewResponseError(resp)
	}
	return JobsClientMarkDevicesShippedResponse{}, nil
}

// markDevicesShippedCreateRequest creates the MarkDevicesShipped request.
func (client *JobsClient) markDevicesShippedCreateRequest(ctx context.Context, jobName string, resourceGroupName string, markDevicesShippedRequest MarkDevicesShippedRequest, options *JobsClientMarkDevicesShippedOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataBox/jobs/{jobName}/markDevicesShipped"
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, markDevicesShippedRequest)
}

// BeginUpdate - Updates the properties of an existing job.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
// resourceGroupName - The Resource Group Name
// jobName - The name of the job Resource within the specified resource group. job names must be between 3 and 24 characters
// in length and use any alphanumeric and underscore only
// jobResourceUpdateParameter - Job update parameters from request body.
// options - JobsClientBeginUpdateOptions contains the optional parameters for the JobsClient.BeginUpdate method.
func (client *JobsClient) BeginUpdate(ctx context.Context, resourceGroupName string, jobName string, jobResourceUpdateParameter JobResourceUpdateParameter, options *JobsClientBeginUpdateOptions) (*runtime.Poller[JobsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, jobName, jobResourceUpdateParameter, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[JobsClientUpdateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[JobsClientUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Update - Updates the properties of an existing job.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2022-02-01
func (client *JobsClient) update(ctx context.Context, resourceGroupName string, jobName string, jobResourceUpdateParameter JobResourceUpdateParameter, options *JobsClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, jobName, jobResourceUpdateParameter, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// updateCreateRequest creates the Update request.
func (client *JobsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, jobName string, jobResourceUpdateParameter JobResourceUpdateParameter, options *JobsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DataBox/jobs/{jobName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if jobName == "" {
		return nil, errors.New("parameter jobName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{jobName}", url.PathEscape(jobName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2022-02-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	if options != nil && options.IfMatch != nil {
		req.Raw().Header["If-Match"] = []string{*options.IfMatch}
	}
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, jobResourceUpdateParameter)
}
