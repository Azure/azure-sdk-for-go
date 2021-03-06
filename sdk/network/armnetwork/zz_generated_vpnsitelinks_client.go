// +build go1.13

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/armcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"net/http"
	"net/url"
	"strings"
)

// VPNSiteLinksClient contains the methods for the VPNSiteLinks group.
// Don't use this type directly, use NewVPNSiteLinksClient() instead.
type VPNSiteLinksClient struct {
	con            *armcore.Connection
	subscriptionID string
}

// NewVPNSiteLinksClient creates a new instance of VPNSiteLinksClient with the specified values.
func NewVPNSiteLinksClient(con *armcore.Connection, subscriptionID string) *VPNSiteLinksClient {
	return &VPNSiteLinksClient{con: con, subscriptionID: subscriptionID}
}

// Get - Retrieves the details of a VPN site link.
// If the operation fails it returns the *CloudError error type.
func (client *VPNSiteLinksClient) Get(ctx context.Context, resourceGroupName string, vpnSiteName string, vpnSiteLinkName string, options *VPNSiteLinksGetOptions) (VPNSiteLinkResponse, error) {
	req, err := client.getCreateRequest(ctx, resourceGroupName, vpnSiteName, vpnSiteLinkName, options)
	if err != nil {
		return VPNSiteLinkResponse{}, err
	}
	resp, err := client.con.Pipeline().Do(req)
	if err != nil {
		return VPNSiteLinkResponse{}, err
	}
	if !resp.HasStatusCode(http.StatusOK) {
		return VPNSiteLinkResponse{}, client.getHandleError(resp)
	}
	return client.getHandleResponse(resp)
}

// getCreateRequest creates the Get request.
func (client *VPNSiteLinksClient) getCreateRequest(ctx context.Context, resourceGroupName string, vpnSiteName string, vpnSiteLinkName string, options *VPNSiteLinksGetOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnSites/{vpnSiteName}/vpnSiteLinks/{vpnSiteLinkName}"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vpnSiteName == "" {
		return nil, errors.New("parameter vpnSiteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vpnSiteName}", url.PathEscape(vpnSiteName))
	if vpnSiteLinkName == "" {
		return nil, errors.New("parameter vpnSiteLinkName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vpnSiteLinkName}", url.PathEscape(vpnSiteLinkName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-02-01")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *VPNSiteLinksClient) getHandleResponse(resp *azcore.Response) (VPNSiteLinkResponse, error) {
	var val *VPNSiteLink
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return VPNSiteLinkResponse{}, err
	}
	return VPNSiteLinkResponse{RawResponse: resp.Response, VPNSiteLink: val}, nil
}

// getHandleError handles the Get error response.
func (client *VPNSiteLinksClient) getHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}

// ListByVPNSite - Lists all the vpnSiteLinks in a resource group for a vpn site.
// If the operation fails it returns the *CloudError error type.
func (client *VPNSiteLinksClient) ListByVPNSite(resourceGroupName string, vpnSiteName string, options *VPNSiteLinksListByVPNSiteOptions) ListVPNSiteLinksResultPager {
	return &listVPNSiteLinksResultPager{
		pipeline: client.con.Pipeline(),
		requester: func(ctx context.Context) (*azcore.Request, error) {
			return client.listByVPNSiteCreateRequest(ctx, resourceGroupName, vpnSiteName, options)
		},
		responder: client.listByVPNSiteHandleResponse,
		errorer:   client.listByVPNSiteHandleError,
		advancer: func(ctx context.Context, resp ListVPNSiteLinksResultResponse) (*azcore.Request, error) {
			return azcore.NewRequest(ctx, http.MethodGet, *resp.ListVPNSiteLinksResult.NextLink)
		},
		statusCodes: []int{http.StatusOK},
	}
}

// listByVPNSiteCreateRequest creates the ListByVPNSite request.
func (client *VPNSiteLinksClient) listByVPNSiteCreateRequest(ctx context.Context, resourceGroupName string, vpnSiteName string, options *VPNSiteLinksListByVPNSiteOptions) (*azcore.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/vpnSites/{vpnSiteName}/vpnSiteLinks"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if vpnSiteName == "" {
		return nil, errors.New("parameter vpnSiteName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{vpnSiteName}", url.PathEscape(vpnSiteName))
	req, err := azcore.NewRequest(ctx, http.MethodGet, azcore.JoinPaths(client.con.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	req.Telemetry(telemetryInfo)
	reqQP := req.URL.Query()
	reqQP.Set("api-version", "2021-02-01")
	req.URL.RawQuery = reqQP.Encode()
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// listByVPNSiteHandleResponse handles the ListByVPNSite response.
func (client *VPNSiteLinksClient) listByVPNSiteHandleResponse(resp *azcore.Response) (ListVPNSiteLinksResultResponse, error) {
	var val *ListVPNSiteLinksResult
	if err := resp.UnmarshalAsJSON(&val); err != nil {
		return ListVPNSiteLinksResultResponse{}, err
	}
	return ListVPNSiteLinksResultResponse{RawResponse: resp.Response, ListVPNSiteLinksResult: val}, nil
}

// listByVPNSiteHandleError handles the ListByVPNSite error response.
func (client *VPNSiteLinksClient) listByVPNSiteHandleError(resp *azcore.Response) error {
	body, err := resp.Payload()
	if err != nil {
		return azcore.NewResponseError(err, resp.Response)
	}
	errType := CloudError{raw: string(body)}
	if err := resp.UnmarshalAsJSON(&errType); err != nil {
		return azcore.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp.Response)
	}
	return azcore.NewResponseError(&errType, resp.Response)
}
