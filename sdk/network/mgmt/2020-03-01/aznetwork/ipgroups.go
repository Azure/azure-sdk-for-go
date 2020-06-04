// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package aznetwork

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// IPGroupsOperations contains the methods for the IPGroups group.
type IPGroupsOperations interface {
	// BeginCreateOrUpdate - Creates or updates an ipGroups in a specified resource group.
	BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, ipGroupsName string, parameters IPGroup) (*IPGroupResponse, error)
	// ResumeCreateOrUpdate - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeCreateOrUpdate(token string) (IPGroupPoller, error)
	// BeginDelete - Deletes the specified ipGroups.
	BeginDelete(ctx context.Context, resourceGroupName string, ipGroupsName string) (*HTTPResponse, error)
	// ResumeDelete - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeDelete(token string) (HTTPPoller, error)
	// Get - Gets the specified ipGroups.
	Get(ctx context.Context, resourceGroupName string, ipGroupsName string, ipGroupsGetOptions *IPGroupsGetOptions) (*IPGroupResponse, error)
	// List - Gets all IpGroups in a subscription.
	List() (IPGroupListResultPager, error)
	// ListByResourceGroup - Gets all IpGroups in a resource group.
	ListByResourceGroup(resourceGroupName string) (IPGroupListResultPager, error)
	// UpdateGroups - Updates tags of an IpGroups resource.
	UpdateGroups(ctx context.Context, resourceGroupName string, ipGroupsName string, parameters TagsObject) (*IPGroupResponse, error)
}

// ipGroupsOperations implements the IPGroupsOperations interface.
type ipGroupsOperations struct {
	*Client
	subscriptionID string
}

// CreateOrUpdate - Creates or updates an ipGroups in a specified resource group.
func (client *ipGroupsOperations) BeginCreateOrUpdate(ctx context.Context, resourceGroupName string, ipGroupsName string, parameters IPGroup) (*IPGroupResponse, error) {
	req, err := client.createOrUpdateCreateRequest(resourceGroupName, ipGroupsName, parameters)
	if err != nil {
		return nil, err
	}
	// send the first request to initialize the poller
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.createOrUpdateHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	pt, err := createPollingTracker("ipGroupsOperations.CreateOrUpdate", "azure-async-operation", resp, client.createOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	poller := &ipGroupPoller{
		pt:       pt,
		pipeline: client.p,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*IPGroupResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *ipGroupsOperations) ResumeCreateOrUpdate(token string) (IPGroupPoller, error) {
	pt, err := resumePollingTracker("ipGroupsOperations.CreateOrUpdate", token, client.createOrUpdateHandleError)
	if err != nil {
		return nil, err
	}
	return &ipGroupPoller{
		pipeline: client.p,
		pt:       pt,
	}, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *ipGroupsOperations) createOrUpdateCreateRequest(resourceGroupName string, ipGroupsName string, parameters IPGroup) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/ipGroups/{ipGroupsName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{ipGroupsName}", url.PathEscape(ipGroupsName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodPut, *u)
	return req, req.MarshalAsJSON(parameters)
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *ipGroupsOperations) createOrUpdateHandleResponse(resp *azcore.Response) (*IPGroupResponse, error) {
	if !resp.HasStatusCode(http.StatusOK, http.StatusCreated, http.StatusNoContent) {
		return nil, client.createOrUpdateHandleError(resp)
	}
	result := IPGroupResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.IPGroup)
}

// createOrUpdateHandleError handles the CreateOrUpdate error response.
func (client *ipGroupsOperations) createOrUpdateHandleError(resp *azcore.Response) error {
	var err Error
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// Delete - Deletes the specified ipGroups.
func (client *ipGroupsOperations) BeginDelete(ctx context.Context, resourceGroupName string, ipGroupsName string) (*HTTPResponse, error) {
	req, err := client.deleteCreateRequest(resourceGroupName, ipGroupsName)
	if err != nil {
		return nil, err
	}
	// send the first request to initialize the poller
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.deleteHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	pt, err := createPollingTracker("ipGroupsOperations.Delete", "location", resp, client.deleteHandleError)
	if err != nil {
		return nil, err
	}
	poller := &httpPoller{
		pt:       pt,
		pipeline: client.p,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*http.Response, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *ipGroupsOperations) ResumeDelete(token string) (HTTPPoller, error) {
	pt, err := resumePollingTracker("ipGroupsOperations.Delete", token, client.deleteHandleError)
	if err != nil {
		return nil, err
	}
	return &httpPoller{
		pipeline: client.p,
		pt:       pt,
	}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *ipGroupsOperations) deleteCreateRequest(resourceGroupName string, ipGroupsName string) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/ipGroups/{ipGroupsName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{ipGroupsName}", url.PathEscape(ipGroupsName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodDelete, *u)
	return req, nil
}

// deleteHandleResponse handles the Delete response.
func (client *ipGroupsOperations) deleteHandleResponse(resp *azcore.Response) (*HTTPResponse, error) {
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, client.deleteHandleError(resp)
	}
	return &HTTPResponse{RawResponse: resp.Response}, nil
}

// deleteHandleError handles the Delete error response.
func (client *ipGroupsOperations) deleteHandleError(resp *azcore.Response) error {
	var err Error
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// Get - Gets the specified ipGroups.
func (client *ipGroupsOperations) Get(ctx context.Context, resourceGroupName string, ipGroupsName string, ipGroupsGetOptions *IPGroupsGetOptions) (*IPGroupResponse, error) {
	req, err := client.getCreateRequest(resourceGroupName, ipGroupsName, ipGroupsGetOptions)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.getHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// getCreateRequest creates the Get request.
func (client *ipGroupsOperations) getCreateRequest(resourceGroupName string, ipGroupsName string, ipGroupsGetOptions *IPGroupsGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/ipGroups/{ipGroupsName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{ipGroupsName}", url.PathEscape(ipGroupsName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	if ipGroupsGetOptions != nil && ipGroupsGetOptions.Expand != nil {
		query.Set("$expand", *ipGroupsGetOptions.Expand)
	}
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodGet, *u)
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *ipGroupsOperations) getHandleResponse(resp *azcore.Response) (*IPGroupResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.getHandleError(resp)
	}
	result := IPGroupResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.IPGroup)
}

// getHandleError handles the Get error response.
func (client *ipGroupsOperations) getHandleError(resp *azcore.Response) error {
	var err Error
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// List - Gets all IpGroups in a subscription.
func (client *ipGroupsOperations) List() (IPGroupListResultPager, error) {
	req, err := client.listCreateRequest()
	if err != nil {
		return nil, err
	}
	return &ipGroupListResultPager{
		pipeline:  client.p,
		request:   req,
		responder: client.listHandleResponse,
		advancer: func(resp *IPGroupListResultResponse) (*azcore.Request, error) {
			u, err := url.Parse(*resp.IPGroupListResult.NextLink)
			if err != nil {
				return nil, fmt.Errorf("invalid NextLink: %w", err)
			}
			if u.Scheme == "" {
				return nil, fmt.Errorf("no scheme detected in NextLink %s", *resp.IPGroupListResult.NextLink)
			}
			return azcore.NewRequest(http.MethodGet, *u), nil
		},
	}, nil
}

// listCreateRequest creates the List request.
func (client *ipGroupsOperations) listCreateRequest() (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/ipGroups"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodGet, *u)
	return req, nil
}

// listHandleResponse handles the List response.
func (client *ipGroupsOperations) listHandleResponse(resp *azcore.Response) (*IPGroupListResultResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.listHandleError(resp)
	}
	result := IPGroupListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.IPGroupListResult)
}

// listHandleError handles the List error response.
func (client *ipGroupsOperations) listHandleError(resp *azcore.Response) error {
	var err Error
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// ListByResourceGroup - Gets all IpGroups in a resource group.
func (client *ipGroupsOperations) ListByResourceGroup(resourceGroupName string) (IPGroupListResultPager, error) {
	req, err := client.listByResourceGroupCreateRequest(resourceGroupName)
	if err != nil {
		return nil, err
	}
	return &ipGroupListResultPager{
		pipeline:  client.p,
		request:   req,
		responder: client.listByResourceGroupHandleResponse,
		advancer: func(resp *IPGroupListResultResponse) (*azcore.Request, error) {
			u, err := url.Parse(*resp.IPGroupListResult.NextLink)
			if err != nil {
				return nil, fmt.Errorf("invalid NextLink: %w", err)
			}
			if u.Scheme == "" {
				return nil, fmt.Errorf("no scheme detected in NextLink %s", *resp.IPGroupListResult.NextLink)
			}
			return azcore.NewRequest(http.MethodGet, *u), nil
		},
	}, nil
}

// listByResourceGroupCreateRequest creates the ListByResourceGroup request.
func (client *ipGroupsOperations) listByResourceGroupCreateRequest(resourceGroupName string) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/ipGroups"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodGet, *u)
	return req, nil
}

// listByResourceGroupHandleResponse handles the ListByResourceGroup response.
func (client *ipGroupsOperations) listByResourceGroupHandleResponse(resp *azcore.Response) (*IPGroupListResultResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.listByResourceGroupHandleError(resp)
	}
	result := IPGroupListResultResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.IPGroupListResult)
}

// listByResourceGroupHandleError handles the ListByResourceGroup error response.
func (client *ipGroupsOperations) listByResourceGroupHandleError(resp *azcore.Response) error {
	var err Error
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}

// UpdateGroups - Updates tags of an IpGroups resource.
func (client *ipGroupsOperations) UpdateGroups(ctx context.Context, resourceGroupName string, ipGroupsName string, parameters TagsObject) (*IPGroupResponse, error) {
	req, err := client.updateGroupsCreateRequest(resourceGroupName, ipGroupsName, parameters)
	if err != nil {
		return nil, err
	}
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.updateGroupsHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// updateGroupsCreateRequest creates the UpdateGroups request.
func (client *ipGroupsOperations) updateGroupsCreateRequest(resourceGroupName string, ipGroupsName string, parameters TagsObject) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/ipGroups/{ipGroupsName}"
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{ipGroupsName}", url.PathEscape(ipGroupsName))
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodPatch, *u)
	return req, req.MarshalAsJSON(parameters)
}

// updateGroupsHandleResponse handles the UpdateGroups response.
func (client *ipGroupsOperations) updateGroupsHandleResponse(resp *azcore.Response) (*IPGroupResponse, error) {
	if !resp.HasStatusCode(http.StatusOK) {
		return nil, client.updateGroupsHandleError(resp)
	}
	result := IPGroupResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.IPGroup)
}

// updateGroupsHandleError handles the UpdateGroups error response.
func (client *ipGroupsOperations) updateGroupsHandleError(resp *azcore.Response) error {
	var err Error
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}
