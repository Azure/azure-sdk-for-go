// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package aznetwork

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// VpnServerConfigurationsAssociatedWithVirtualWanOperations contains the methods for the VpnServerConfigurationsAssociatedWithVirtualWan group.
type VpnServerConfigurationsAssociatedWithVirtualWanOperations interface {
	// BeginList - Gives the list of VpnServerConfigurations associated with Virtual Wan in a resource group.
	BeginList(ctx context.Context, resourceGroupName string, virtualWanName string) (*VpnServerConfigurationsResponseResponse, error)
	// ResumeList - Used to create a new instance of this poller from the resume token of a previous instance of this poller type.
	ResumeList(token string) (VpnServerConfigurationsResponsePoller, error)
}

// vpnServerConfigurationsAssociatedWithVirtualWanOperations implements the VpnServerConfigurationsAssociatedWithVirtualWanOperations interface.
type vpnServerConfigurationsAssociatedWithVirtualWanOperations struct {
	*Client
	subscriptionID string
}

// List - Gives the list of VpnServerConfigurations associated with Virtual Wan in a resource group.
func (client *vpnServerConfigurationsAssociatedWithVirtualWanOperations) BeginList(ctx context.Context, resourceGroupName string, virtualWanName string) (*VpnServerConfigurationsResponseResponse, error) {
	req, err := client.listCreateRequest(resourceGroupName, virtualWanName)
	if err != nil {
		return nil, err
	}
	// send the first request to initialize the poller
	resp, err := client.p.Do(ctx, req)
	if err != nil {
		return nil, err
	}
	result, err := client.listHandleResponse(resp)
	if err != nil {
		return nil, err
	}
	pt, err := createPollingTracker("vpnServerConfigurationsAssociatedWithVirtualWanOperations.List", "location", resp, client.listHandleError)
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

func (client *vpnServerConfigurationsAssociatedWithVirtualWanOperations) ResumeList(token string) (VpnServerConfigurationsResponsePoller, error) {
	pt, err := resumePollingTracker("vpnServerConfigurationsAssociatedWithVirtualWanOperations.List", token, client.listHandleError)
	if err != nil {
		return nil, err
	}
	return &vpnServerConfigurationsResponsePoller{
		pipeline: client.p,
		pt:       pt,
	}, nil
}

// listCreateRequest creates the List request.
func (client *vpnServerConfigurationsAssociatedWithVirtualWanOperations) listCreateRequest(resourceGroupName string, virtualWanName string) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans/{virtualWANName}/vpnServerConfigurations"
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	urlPath = strings.ReplaceAll(urlPath, "{virtualWANName}", url.PathEscape(virtualWanName))
	u, err := client.u.Parse(urlPath)
	if err != nil {
		return nil, err
	}
	query := u.Query()
	query.Set("api-version", "2020-03-01")
	u.RawQuery = query.Encode()
	req := azcore.NewRequest(http.MethodPost, *u)
	return req, nil
}

// listHandleResponse handles the List response.
func (client *vpnServerConfigurationsAssociatedWithVirtualWanOperations) listHandleResponse(resp *azcore.Response) (*VpnServerConfigurationsResponseResponse, error) {
	if !resp.HasStatusCode(http.StatusOK, http.StatusAccepted, http.StatusNoContent) {
		return nil, client.listHandleError(resp)
	}
	result := VpnServerConfigurationsResponseResponse{RawResponse: resp.Response}
	return &result, resp.UnmarshalAsJSON(&result.VpnServerConfigurationsResponse)
}

// listHandleError handles the List error response.
func (client *vpnServerConfigurationsAssociatedWithVirtualWanOperations) listHandleError(resp *azcore.Response) error {
	var err CloudError
	if err := resp.UnmarshalAsJSON(&err); err != nil {
		return err
	}
	return err
}
