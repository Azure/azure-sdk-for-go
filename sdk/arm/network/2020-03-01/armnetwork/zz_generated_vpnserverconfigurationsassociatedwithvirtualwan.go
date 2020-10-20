// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// VpnServerConfigurationsAssociatedWithVirtualWanOperations contains the methods for the VpnServerConfigurationsAssociatedWithVirtualWan group.
type VpnServerConfigurationsAssociatedWithVirtualWanOperations interface {
	// BeginList - Gives the list of VpnServerConfigurations associated with Virtual Wan in a resource group.
	BeginList(ctx context.Context, resourceGroupName string, virtualWanName string, options *VpnServerConfigurationsAssociatedWithVirtualWanListOptions) (*VpnServerConfigurationsResponsePollerResponse, error)
	// ResumeList - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeList(token string) (VpnServerConfigurationsResponsePoller, error)
}

// VpnServerConfigurationsAssociatedWithVirtualWanClient implements the VpnServerConfigurationsAssociatedWithVirtualWanOperations interface.
// Don't use this type directly, use NewVpnServerConfigurationsAssociatedWithVirtualWanClient() instead.
type VpnServerConfigurationsAssociatedWithVirtualWanClient struct {
	*Client
	subscriptionID string
}

// NewVpnServerConfigurationsAssociatedWithVirtualWanClient creates a new instance of VpnServerConfigurationsAssociatedWithVirtualWanClient with the specified values.
func NewVpnServerConfigurationsAssociatedWithVirtualWanClient(c *Client, subscriptionID string) VpnServerConfigurationsAssociatedWithVirtualWanOperations {
	return &VpnServerConfigurationsAssociatedWithVirtualWanClient{Client: c, subscriptionID: subscriptionID}
}

// Do invokes the Do() method on the pipeline associated with this client.
func (client *VpnServerConfigurationsAssociatedWithVirtualWanClient) Do(req *azcore.Request) (*azcore.Response, error) {
	return client.p.Do(req)
}

func (client *VpnServerConfigurationsAssociatedWithVirtualWanClient) BeginList(ctx context.Context, resourceGroupName string, virtualWanName string, options *VpnServerConfigurationsAssociatedWithVirtualWanListOptions) (*VpnServerConfigurationsResponsePollerResponse, error) {
	resp, err := client.List(ctx, resourceGroupName, virtualWanName, options)
	if err != nil {
		return nil, err
	}
	result := &VpnServerConfigurationsResponsePollerResponse{
		RawResponse: resp.Response,
	}
	pt, err := armcore.NewPoller("VpnServerConfigurationsAssociatedWithVirtualWanClient.List", "location", resp, client.ListHandleError)
	if err != nil {
		return nil, err
	}
	poller := &vpnServerConfigurationsResponsePoller{
		pt:       pt,
		pipeline: client.p,
	}
	result.Poller = poller
	result.PollUntilDone = func(ctx context.Context, frequency time.Duration) (*VpnServerConfigurationsResponseResponse, error) {
		return poller.pollUntilDone(ctx, frequency)
	}
	return result, nil
}

func (client *VpnServerConfigurationsAssociatedWithVirtualWanClient) ResumeList(token string) (VpnServerConfigurationsResponsePoller, error) {
	pt, err := armcore.NewPollerFromResumeToken("VpnServerConfigurationsAssociatedWithVirtualWanClient.List", token, client.ListHandleError)
	if err != nil {
		return nil, err
	}
	return &vpnServerConfigurationsResponsePoller{
		pipeline: client.p,
		pt:       pt,
	}, nil
}

// List - Gives the list of VpnServerConfigurations associated with Virtual Wan in a resource group.
func (client *VpnServerConfigurationsAssociatedWithVirtualWanClient) List(ctx context.Context, resourceGroupName string, virtualWanName string, options *VpnServerConfigurationsAssociatedWithVirtualWanListOptions) (*azcore.Response, error) {
	req, err := client.ListCreateRequest(ctx, resourceGroupName, virtualWanName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted) {
		return nil, client.ListHandleError(resp)
	}
	return resp, nil
}

// ListCreateRequest creates the List request.
func (client *VpnServerConfigurationsAssociatedWithVirtualWanClient) ListCreateRequest(ctx context.Context, resourceGroupName string, virtualWanName string, options *VpnServerConfigurationsAssociatedWithVirtualWanListOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans/{virtualWANName}/vpnServerConfigurations"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{virtualWANName}", url.PathEscape(virtualWanName))
	req, err := azcore.NewRequest(ctx, http.MethodPost, azcore.JoinPaths(client.u, urlPath))
	if err != nil {
		return nil, err
	}
	query := req.URL.Query()
	query.Set("api-version", "2020-03-01")
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// ListHandleResponse handles the List response.
func (client *VpnServerConfigurationsAssociatedWithVirtualWanClient) ListHandleResponse(resp *azcore.Response) (*VpnServerConfigurationsResponseResponse, error) {
	result := VpnServerConfigurationsResponseResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.VpnServerConfigurationsResponse)
}

// ListHandleError handles the List error response.
func (client *VpnServerConfigurationsAssociatedWithVirtualWanClient) ListHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return azcore.NewResponseError(&err, resp.Response)
}
