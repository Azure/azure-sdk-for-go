//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armpostgresqlhsc

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

// ServerGroupsClient contains the methods for the ServerGroups group.
// Don't use this type directly, use NewServerGroupsClient() instead.
type ServerGroupsClient struct {
	host           string
	subscriptionID string
	pl             runtime.Pipeline
}

// NewServerGroupsClient creates a new instance of ServerGroupsClient with the specified values.
// subscriptionID - The ID of the target subscription.
// credential - used to authorize requests. Usually a credential from azidentity.
// options - pass nil to accept the default values.
func NewServerGroupsClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) (*ServerGroupsClient, error) {
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
	client := &ServerGroupsClient{
		subscriptionID: subscriptionID,
		host:           ep,
		pl:             pl,
	}
	return client, nil
}

// CheckNameAvailability - Check the availability of name for resource
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
// nameAvailabilityRequest - The required parameters for checking if resource name is available.
// options - ServerGroupsClientCheckNameAvailabilityOptions contains the optional parameters for the ServerGroupsClient.CheckNameAvailability
// method.
func (client *ServerGroupsClient) CheckNameAvailability(ctx context.Context, nameAvailabilityRequest NameAvailabilityRequest, options *ServerGroupsClientCheckNameAvailabilityOptions) (ServerGroupsClientCheckNameAvailabilityResponse, error) {
	req, err := client.checkNameAvailabilityCreateRequest(ctx, nameAvailabilityRequest, options)
	if err != nil {
		return ServerGroupsClientCheckNameAvailabilityResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ServerGroupsClientCheckNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServerGroupsClientCheckNameAvailabilityResponse{}, runtime.NewResponseError(resp)
	}
	return client.checkNameAvailabilityHandleResponse(resp)
}

// checkNameAvailabilityCreateRequest creates the CheckNameAvailability request.
func (client *ServerGroupsClient) checkNameAvailabilityCreateRequest(ctx context.Context, nameAvailabilityRequest NameAvailabilityRequest, options *ServerGroupsClientCheckNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.DBForPostgreSql/checkNameAvailability"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-05-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, nameAvailabilityRequest)
}

// checkNameAvailabilityHandleResponse handles the CheckNameAvailability response.
func (client *ServerGroupsClient) checkNameAvailabilityHandleResponse(resp *http.Response) (ServerGroupsClientCheckNameAvailabilityResponse, error) {
	result := ServerGroupsClientCheckNameAvailabilityResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.NameAvailability); err != nil {
		return ServerGroupsClientCheckNameAvailabilityResponse{}, err
	}
	return result, nil
}

// BeginCreateOrUpdate - Creates a new server group with servers.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// serverGroupName - The name of the server group.
// parameters - The required parameters for creating or updating a server group.
// options - ServerGroupsClientBeginCreateOrUpdateOptions contains the optional parameters for the ServerGroupsClient.BeginCreateOrUpdate
// method.
func (client *ServerGroupsClient) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, serverGroupName string, parameters ServerGroup, options *ServerGroupsClientBeginCreateOrUpdateOptions) (*runtime.Poller[ServerGroupsClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, resourceGroupName, serverGroupName, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[ServerGroupsClientCreateOrUpdateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[ServerGroupsClientCreateOrUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// CreateOrUpdate - Creates a new server group with servers.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
func (client *ServerGroupsClient) createOrUpdate(ctx context.Context, resourceGroupName string, serverGroupName string, parameters ServerGroup, options *ServerGroupsClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	req, err := client.createOrUpdateCreateRequest(ctx, resourceGroupName, serverGroupName, parameters, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusCreated, http.StatusAccepted) {
		return nil, runtime.NewResponseError(resp)
	}
	return resp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ServerGroupsClient) createOrUpdateCreateRequest(ctx context.Context, resourceGroupName string, serverGroupName string, parameters ServerGroup, options *ServerGroupsClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBForPostgreSql/serverGroupsv2/{serverGroupName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverGroupName == "" {
		return nil, errors.New("parameter serverGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverGroupName}", url.PathEscape(serverGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-05-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}

// BeginDelete - Deletes a server group together with servers in it.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// serverGroupName - The name of the server group.
// options - ServerGroupsClientBeginDeleteOptions contains the optional parameters for the ServerGroupsClient.BeginDelete
// method.
func (client *ServerGroupsClient) BeginDelete(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginDeleteOptions) (*runtime.Poller[ServerGroupsClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, resourceGroupName, serverGroupName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[ServerGroupsClientDeleteResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[ServerGroupsClientDeleteResponse](options.ResumeToken, client.pl, nil)
	}
}

// Delete - Deletes a server group together with servers in it.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
func (client *ServerGroupsClient) deleteOperation(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginDeleteOptions) (*http.Response, error) {
	req, err := client.deleteCreateRequest(ctx, resourceGroupName, serverGroupName, options)
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
func (client *ServerGroupsClient) deleteCreateRequest(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBForPostgreSql/serverGroupsv2/{serverGroupName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverGroupName == "" {
		return nil, errors.New("parameter serverGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverGroupName}", url.PathEscape(serverGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-05-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Gets information about a server group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// serverGroupName - The name of the server group.
// options - ServerGroupsClientGetOptions contains the optional parameters for the ServerGroupsClient.Get method.
func (client *ServerGroupsClient) Get(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientGetOptions) (ServerGroupsClientGetResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, serverGroupName, options)
	if err != nil {
		return ServerGroupsClientGetResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return ServerGroupsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServerGroupsClientGetResponse{}, runtime.NewResponseError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *ServerGroupsClient) getCreateRequest(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientGetOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBForPostgreSql/serverGroupsv2/{serverGroupName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverGroupName == "" {
		return nil, errors.New("parameter serverGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverGroupName}", url.PathEscape(serverGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-05-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ServerGroupsClient) getHandleResponse(resp *http.Response) (ServerGroupsClientGetResponse, error) {
	result := ServerGroupsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ServerGroup); err != nil {
		return ServerGroupsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List all the server groups in a given subscription.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
// options - ServerGroupsClientListOptions contains the optional parameters for the ServerGroupsClient.List method.
func (client *ServerGroupsClient) NewListPager(options *ServerGroupsClientListOptions) *runtime.Pager[ServerGroupsClientListResponse] {
	return runtime.NewPager(runtime.PagingHandler[ServerGroupsClientListResponse]{
		More: func(page ServerGroupsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ServerGroupsClientListResponse) (ServerGroupsClientListResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ServerGroupsClientListResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return ServerGroupsClientListResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ServerGroupsClientListResponse{}, runtime.NewResponseError(resp)
			}
			return client.listHandleResponse(resp)
		},
	})
}

// listCreateRequest creates the List request.
func (client *ServerGroupsClient) listCreateRequest(ctx context.Context, options *ServerGroupsClientListOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.DBForPostgreSql/serverGroupsv2"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-05-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ServerGroupsClient) listHandleResponse(resp *http.Response) (ServerGroupsClientListResponse, error) {
	result := ServerGroupsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ServerGroupListResult); err != nil {
		return ServerGroupsClientListResponse{}, err
	}
	return result, nil
}

// NewListByResourceGroupPager - List all the server groups in a given resource group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// options - ServerGroupsClientListByResourceGroupOptions contains the optional parameters for the ServerGroupsClient.ListByResourceGroup
// method.
func (client *ServerGroupsClient) NewListByResourceGroupPager(resourceGroupName string, options *ServerGroupsClientListByResourceGroupOptions) *runtime.Pager[ServerGroupsClientListByResourceGroupResponse] {
	return runtime.NewPager(runtime.PagingHandler[ServerGroupsClientListByResourceGroupResponse]{
		More: func(page ServerGroupsClientListByResourceGroupResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ServerGroupsClientListByResourceGroupResponse) (ServerGroupsClientListByResourceGroupResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByResourceGroupCreateRequest(ctx, resourceGroupName, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ServerGroupsClientListByResourceGroupResponse{}, err
			}
			resp, err := client.pl.Do(req)
			if err != nil {
				return ServerGroupsClientListByResourceGroupResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ServerGroupsClientListByResourceGroupResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByResourceGroupHandleResponse(resp)
		},
	})
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *ServerGroupsClient) listByResourceGroupCreateRequest(ctx context.Context, resourceGroupName string, options *ServerGroupsClientListByResourceGroupOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBForPostgreSql/serverGroupsv2"
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
	reqQP.Set("api-version", "2020-10-05-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *ServerGroupsClient) listByResourceGroupHandleResponse(resp *http.Response) (ServerGroupsClientListByResourceGroupResponse, error) {
	result := ServerGroupsClientListByResourceGroupResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.ServerGroupListResult); err != nil {
		return ServerGroupsClientListByResourceGroupResponse{}, err
	}
	return result, nil
}

// BeginRestart - Restarts the server group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// serverGroupName - The name of the server group.
// options - ServerGroupsClientBeginRestartOptions contains the optional parameters for the ServerGroupsClient.BeginRestart
// method.
func (client *ServerGroupsClient) BeginRestart(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginRestartOptions) (*runtime.Poller[ServerGroupsClientRestartResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.restart(ctx, resourceGroupName, serverGroupName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[ServerGroupsClientRestartResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[ServerGroupsClientRestartResponse](options.ResumeToken, client.pl, nil)
	}
}

// Restart - Restarts the server group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
func (client *ServerGroupsClient) restart(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginRestartOptions) (*http.Response, error) {
	req, err := client.restartCreateRequest(ctx, resourceGroupName, serverGroupName, options)
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

// restartCreateRequest creates the Restart request.
func (client *ServerGroupsClient) restartCreateRequest(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginRestartOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBForPostgreSql/serverGroupsv2/{serverGroupName}/restart"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverGroupName == "" {
		return nil, errors.New("parameter serverGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverGroupName}", url.PathEscape(serverGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-05-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginStart - Starts the server group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// serverGroupName - The name of the server group.
// options - ServerGroupsClientBeginStartOptions contains the optional parameters for the ServerGroupsClient.BeginStart method.
func (client *ServerGroupsClient) BeginStart(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginStartOptions) (*runtime.Poller[ServerGroupsClientStartResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.start(ctx, resourceGroupName, serverGroupName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[ServerGroupsClientStartResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[ServerGroupsClientStartResponse](options.ResumeToken, client.pl, nil)
	}
}

// Start - Starts the server group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
func (client *ServerGroupsClient) start(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginStartOptions) (*http.Response, error) {
	req, err := client.startCreateRequest(ctx, resourceGroupName, serverGroupName, options)
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

// startCreateRequest creates the Start request.
func (client *ServerGroupsClient) startCreateRequest(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginStartOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBForPostgreSql/serverGroupsv2/{serverGroupName}/start"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverGroupName == "" {
		return nil, errors.New("parameter serverGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverGroupName}", url.PathEscape(serverGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-05-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginStop - Stops the server group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// serverGroupName - The name of the server group.
// options - ServerGroupsClientBeginStopOptions contains the optional parameters for the ServerGroupsClient.BeginStop method.
func (client *ServerGroupsClient) BeginStop(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginStopOptions) (*runtime.Poller[ServerGroupsClientStopResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.stop(ctx, resourceGroupName, serverGroupName, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[ServerGroupsClientStopResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[ServerGroupsClientStopResponse](options.ResumeToken, client.pl, nil)
	}
}

// Stop - Stops the server group.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
func (client *ServerGroupsClient) stop(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginStopOptions) (*http.Response, error) {
	req, err := client.stopCreateRequest(ctx, resourceGroupName, serverGroupName, options)
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

// stopCreateRequest creates the Stop request.
func (client *ServerGroupsClient) stopCreateRequest(ctx context.Context, resourceGroupName string, serverGroupName string, options *ServerGroupsClientBeginStopOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBForPostgreSql/serverGroupsv2/{serverGroupName}/stop"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverGroupName == "" {
		return nil, errors.New("parameter serverGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverGroupName}", url.PathEscape(serverGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-05-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// BeginUpdate - Updates an existing server group. The request body can contain one to many of the properties present in the
// normal server group definition.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
// resourceGroupName - The name of the resource group. The name is case insensitive.
// serverGroupName - The name of the server group.
// parameters - The parameters for updating a server group.
// options - ServerGroupsClientBeginUpdateOptions contains the optional parameters for the ServerGroupsClient.BeginUpdate
// method.
func (client *ServerGroupsClient) BeginUpdate(ctx context.Context, resourceGroupName string, serverGroupName string, parameters ServerGroupForUpdate, options *ServerGroupsClientBeginUpdateOptions) (*runtime.Poller[ServerGroupsClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, resourceGroupName, serverGroupName, parameters, options)
		if err != nil {
			return nil, err
		}
		return runtime.NewPoller[ServerGroupsClientUpdateResponse](resp, client.pl, nil)
	} else {
		return runtime.NewPollerFromResumeToken[ServerGroupsClientUpdateResponse](options.ResumeToken, client.pl, nil)
	}
}

// Update - Updates an existing server group. The request body can contain one to many of the properties present in the normal
// server group definition.
// If the operation fails it returns an *azcore.ResponseError type.
// Generated from API version 2020-10-05-privatepreview
func (client *ServerGroupsClient) update(ctx context.Context, resourceGroupName string, serverGroupName string, parameters ServerGroupForUpdate, options *ServerGroupsClientBeginUpdateOptions) (*http.Response, error) {
	req, err := client.updateCreateRequest(ctx, resourceGroupName, serverGroupName, parameters, options)
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
func (client *ServerGroupsClient) updateCreateRequest(ctx context.Context, resourceGroupName string, serverGroupName string, parameters ServerGroupForUpdate, options *ServerGroupsClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.DBForPostgreSql/serverGroupsv2/{serverGroupName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if serverGroupName == "" {
		return nil, errors.New("parameter serverGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{serverGroupName}", url.PathEscape(serverGroupName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.host, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2020-10-05-privatepreview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, runtime.MarshalAsJSON(req, parameters)
}
