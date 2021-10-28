//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armnetwork

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

// NetworkManagementClient contains the methods for the NetworkManagementClient group.
// Don't use this type directly, use NewNetworkManagementClient() instead.
type NetworkManagementClient struct {
	ep             string
	pl             runtime.Pipeline
	subscriptionID string
}

// NewNetworkManagementClient creates a new instance of NetworkManagementClient with the specified values.
func NewNetworkManagementClient(subscriptionID string, credential azcore.TokenCredential, options *arm.ClientOptions) *NetworkManagementClient {
	cp := arm.ClientOptions{}
	if options != nil {
		cp = *options
	}
	if len(cp.Host) == 0 {
		cp.Host = arm.AzurePublicCloud
	}
	return &NetworkManagementClient{subscriptionID: subscriptionID, ep: string(cp.Host), pl: armruntime.NewPipeline(module, version, credential, &cp)}
}

// CheckDNSNameAvailability - Checks whether a domain name in the cloudapp.azure.com zone is available for use.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) CheckDNSNameAvailability(ctx context.Context, location string, domainNameLabel string, options *NetworkManagementClientCheckDNSNameAvailabilityOptions) (NetworkManagementClientCheckDNSNameAvailabilityResponse, error) {
	req, err := client.checkDNSNameAvailabilityCreateRequest(ctx, location, domainNameLabel, options)
	if err != nil {
		return NetworkManagementClientCheckDNSNameAvailabilityResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return NetworkManagementClientCheckDNSNameAvailabilityResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return NetworkManagementClientCheckDNSNameAvailabilityResponse{}, client.checkDNSNameAvailabilityHandleError(resp)
	}
	return client.checkDNSNameAvailabilityHandleResponse(resp)
}

// checkDNSNameAvailabilityCreateRequest creates the CheckDNSNameAvailability request.
func (client *NetworkManagementClient) checkDNSNameAvailabilityCreateRequest(ctx context.Context, location string, domainNameLabel string, options *NetworkManagementClientCheckDNSNameAvailabilityOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/providers/Microsoft.Network/locations/{location}/CheckDnsNameAvailability"
	if location == "" {
		return nil, errors.New("parameter location cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{location}", url.PathEscape(location))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("domainNameLabel", domainNameLabel)
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// checkDNSNameAvailabilityHandleResponse handles the CheckDNSNameAvailability response.
func (client *NetworkManagementClient) checkDNSNameAvailabilityHandleResponse(resp *http.Response) (NetworkManagementClientCheckDNSNameAvailabilityResponse, error) {
	result := NetworkManagementClientCheckDNSNameAvailabilityResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.DNSNameAvailabilityResult); err != nil {
		return NetworkManagementClientCheckDNSNameAvailabilityResponse{}, err
	}
	return result, nil
}

// checkDNSNameAvailabilityHandleError handles the CheckDNSNameAvailability error response.
func (client *NetworkManagementClient) checkDNSNameAvailabilityHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginDeleteBastionShareableLink - Deletes the Bastion Shareable Links for all the VMs specified in the request.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) BeginDeleteBastionShareableLink(ctx context.Context, resourceGroupName string, bastionHostName string, bslRequest BastionShareableLinkListRequest, options *NetworkManagementClientBeginDeleteBastionShareableLinkOptions) (NetworkManagementClientDeleteBastionShareableLinkPollerResponse, error) {
	resp, err := client.deleteBastionShareableLink(ctx, resourceGroupName, bastionHostName, bslRequest, options)
	if err != nil {
		return NetworkManagementClientDeleteBastionShareableLinkPollerResponse{}, err
	}
	result := NetworkManagementClientDeleteBastionShareableLinkPollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("NetworkManagementClient.DeleteBastionShareableLink", "location", resp, client.pl, client.deleteBastionShareableLinkHandleError)
	if err != nil {
		return NetworkManagementClientDeleteBastionShareableLinkPollerResponse{}, err
	}
	result.Poller = &NetworkManagementClientDeleteBastionShareableLinkPoller{
		pt: pt,
	}
	return result, nil
}

// DeleteBastionShareableLink - Deletes the Bastion Shareable Links for all the VMs specified in the request.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) deleteBastionShareableLink(ctx context.Context, resourceGroupName string, bastionHostName string, bslRequest BastionShareableLinkListRequest, options *NetworkManagementClientBeginDeleteBastionShareableLinkOptions) (*http.Response, error) {
	req, err := client.deleteBastionShareableLinkCreateRequest(ctx, resourceGroupName, bastionHostName, bslRequest, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.deleteBastionShareableLinkHandleError(resp)
	}
	return resp, nil
}

// deleteBastionShareableLinkCreateRequest creates the DeleteBastionShareableLink request.
func (client *NetworkManagementClient) deleteBastionShareableLinkCreateRequest(ctx context.Context, resourceGroupName string, bastionHostName string, bslRequest BastionShareableLinkListRequest, options *NetworkManagementClientBeginDeleteBastionShareableLinkOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/bastionHosts/{bastionHostName}/deleteShareableLinks"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if bastionHostName == "" {
		return nil, errors.New("parameter bastionHostName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bastionHostName}", url.PathEscape(bastionHostName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, bslRequest)
}

// deleteBastionShareableLinkHandleError handles the DeleteBastionShareableLink error response.
func (client *NetworkManagementClient) deleteBastionShareableLinkHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// DisconnectActiveSessions - Returns the list of currently active sessions on the Bastion.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) DisconnectActiveSessions(resourceGroupName string, bastionHostName string, sessionIDs SessionIDs, options *NetworkManagementClientDisconnectActiveSessionsOptions) *NetworkManagementClientDisconnectActiveSessionsPager {
	return &NetworkManagementClientDisconnectActiveSessionsPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.disconnectActiveSessionsCreateRequest(ctx, resourceGroupName, bastionHostName, sessionIDs, options)
		},
		advancer: func(ctx context.Context, resp NetworkManagementClientDisconnectActiveSessionsResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.BastionSessionDeleteResult.NextLink)
		},
	}
}

// disconnectActiveSessionsCreateRequest creates the DisconnectActiveSessions request.
func (client *NetworkManagementClient) disconnectActiveSessionsCreateRequest(ctx context.Context, resourceGroupName string, bastionHostName string, sessionIDs SessionIDs, options *NetworkManagementClientDisconnectActiveSessionsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/bastionHosts/{bastionHostName}/disconnectActiveSessions"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if bastionHostName == "" {
		return nil, errors.New("parameter bastionHostName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bastionHostName}", url.PathEscape(bastionHostName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, sessionIDs)
}

// disconnectActiveSessionsHandleResponse handles the DisconnectActiveSessions response.
func (client *NetworkManagementClient) disconnectActiveSessionsHandleResponse(resp *http.Response) (NetworkManagementClientDisconnectActiveSessionsResponse, error) {
	result := NetworkManagementClientDisconnectActiveSessionsResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.BastionSessionDeleteResult); err != nil {
		return NetworkManagementClientDisconnectActiveSessionsResponse{}, err
	}
	return result, nil
}

// disconnectActiveSessionsHandleError handles the DisconnectActiveSessions error response.
func (client *NetworkManagementClient) disconnectActiveSessionsHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginGeneratevirtualwanvpnserverconfigurationvpnprofile - Generates a unique VPN profile for P2S clients for VirtualWan and associated VpnServerConfiguration
// combination in the specified resource group.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) BeginGeneratevirtualwanvpnserverconfigurationvpnprofile(ctx context.Context, resourceGroupName string, virtualWANName string, vpnClientParams VirtualWanVPNProfileParameters, options *NetworkManagementClientBeginGeneratevirtualwanvpnserverconfigurationvpnprofileOptions) (NetworkManagementClientGeneratevirtualwanvpnserverconfigurationvpnprofilePollerResponse, error) {
	resp, err := client.generatevirtualwanvpnserverconfigurationvpnprofile(ctx, resourceGroupName, virtualWANName, vpnClientParams, options)
	if err != nil {
		return NetworkManagementClientGeneratevirtualwanvpnserverconfigurationvpnprofilePollerResponse{}, err
	}
	result := NetworkManagementClientGeneratevirtualwanvpnserverconfigurationvpnprofilePollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("NetworkManagementClient.Generatevirtualwanvpnserverconfigurationvpnprofile", "location", resp, client.pl, client.generatevirtualwanvpnserverconfigurationvpnprofileHandleError)
	if err != nil {
		return NetworkManagementClientGeneratevirtualwanvpnserverconfigurationvpnprofilePollerResponse{}, err
	}
	result.Poller = &NetworkManagementClientGeneratevirtualwanvpnserverconfigurationvpnprofilePoller{
		pt: pt,
	}
	return result, nil
}

// Generatevirtualwanvpnserverconfigurationvpnprofile - Generates a unique VPN profile for P2S clients for VirtualWan and associated VpnServerConfiguration
// combination in the specified resource group.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) generatevirtualwanvpnserverconfigurationvpnprofile(ctx context.Context, resourceGroupName string, virtualWANName string, vpnClientParams VirtualWanVPNProfileParameters, options *NetworkManagementClientBeginGeneratevirtualwanvpnserverconfigurationvpnprofileOptions) (*http.Response, error) {
	req, err := client.generatevirtualwanvpnserverconfigurationvpnprofileCreateRequest(ctx, resourceGroupName, virtualWANName, vpnClientParams, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.generatevirtualwanvpnserverconfigurationvpnprofileHandleError(resp)
	}
	return resp, nil
}

// generatevirtualwanvpnserverconfigurationvpnprofileCreateRequest creates the Generatevirtualwanvpnserverconfigurationvpnprofile request.
func (client *NetworkManagementClient) generatevirtualwanvpnserverconfigurationvpnprofileCreateRequest(ctx context.Context, resourceGroupName string, virtualWANName string, vpnClientParams VirtualWanVPNProfileParameters, options *NetworkManagementClientBeginGeneratevirtualwanvpnserverconfigurationvpnprofileOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans/{virtualWANName}/GenerateVpnProfile"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if virtualWANName == "" {
		return nil, errors.New("parameter virtualWANName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{virtualWANName}", url.PathEscape(virtualWANName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, vpnClientParams)
}

// generatevirtualwanvpnserverconfigurationvpnprofileHandleError handles the Generatevirtualwanvpnserverconfigurationvpnprofile error response.
func (client *NetworkManagementClient) generatevirtualwanvpnserverconfigurationvpnprofileHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginGetActiveSessions - Returns the list of currently active sessions on the Bastion.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) BeginGetActiveSessions(ctx context.Context, resourceGroupName string, bastionHostName string, options *NetworkManagementClientBeginGetActiveSessionsOptions) (NetworkManagementClientGetActiveSessionsPollerResponse, error) {
	resp, err := client.getActiveSessions(ctx, resourceGroupName, bastionHostName, options)
	if err != nil {
		return NetworkManagementClientGetActiveSessionsPollerResponse{}, err
	}
	result := NetworkManagementClientGetActiveSessionsPollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("NetworkManagementClient.GetActiveSessions", "location", resp, client.pl, client.getActiveSessionsHandleError)
	if err != nil {
		return NetworkManagementClientGetActiveSessionsPollerResponse{}, err
	}
	result.Poller = &NetworkManagementClientGetActiveSessionsPoller{
		pt:     pt,
		client: client,
	}
	return result, nil
}

// GetActiveSessions - Returns the list of currently active sessions on the Bastion.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) getActiveSessions(ctx context.Context, resourceGroupName string, bastionHostName string, options *NetworkManagementClientBeginGetActiveSessionsOptions) (*http.Response, error) {
	req, err := client.getActiveSessionsCreateRequest(ctx, resourceGroupName, bastionHostName, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.getActiveSessionsHandleError(resp)
	}
	return resp, nil
}

// getActiveSessionsCreateRequest creates the GetActiveSessions request.
func (client *NetworkManagementClient) getActiveSessionsCreateRequest(ctx context.Context, resourceGroupName string, bastionHostName string, options *NetworkManagementClientBeginGetActiveSessionsOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/bastionHosts/{bastionHostName}/getActiveSessions"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if bastionHostName == "" {
		return nil, errors.New("parameter bastionHostName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bastionHostName}", url.PathEscape(bastionHostName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// getActiveSessionsHandleResponse handles the GetActiveSessions response.
func (client *NetworkManagementClient) getActiveSessionsHandleResponse(resp *http.Response) (NetworkManagementClientGetActiveSessionsResponse, error) {
	result := NetworkManagementClientGetActiveSessionsResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.BastionActiveSessionListResult); err != nil {
		return NetworkManagementClientGetActiveSessionsResponse{}, err
	}
	return result, nil
}

// getActiveSessionsHandleError handles the GetActiveSessions error response.
func (client *NetworkManagementClient) getActiveSessionsHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// GetBastionShareableLink - Return the Bastion Shareable Links for all the VMs specified in the request.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) GetBastionShareableLink(resourceGroupName string, bastionHostName string, bslRequest BastionShareableLinkListRequest, options *NetworkManagementClientGetBastionShareableLinkOptions) *NetworkManagementClientGetBastionShareableLinkPager {
	return &NetworkManagementClientGetBastionShareableLinkPager{
		client: client,
		requester: func(ctx context.Context) (*policy.Request, error) {
			return client.getBastionShareableLinkCreateRequest(ctx, resourceGroupName, bastionHostName, bslRequest, options)
		},
		advancer: func(ctx context.Context, resp NetworkManagementClientGetBastionShareableLinkResponse) (*policy.Request, error) {
			return runtime.NewRequest(ctx, http.MethodGet, *resp.BastionShareableLinkListResult.NextLink)
		},
	}
}

// getBastionShareableLinkCreateRequest creates the GetBastionShareableLink request.
func (client *NetworkManagementClient) getBastionShareableLinkCreateRequest(ctx context.Context, resourceGroupName string, bastionHostName string, bslRequest BastionShareableLinkListRequest, options *NetworkManagementClientGetBastionShareableLinkOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/bastionHosts/{bastionHostName}/getShareableLinks"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if bastionHostName == "" {
		return nil, errors.New("parameter bastionHostName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bastionHostName}", url.PathEscape(bastionHostName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, bslRequest)
}

// getBastionShareableLinkHandleResponse handles the GetBastionShareableLink response.
func (client *NetworkManagementClient) getBastionShareableLinkHandleResponse(resp *http.Response) (NetworkManagementClientGetBastionShareableLinkResponse, error) {
	result := NetworkManagementClientGetBastionShareableLinkResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.BastionShareableLinkListResult); err != nil {
		return NetworkManagementClientGetBastionShareableLinkResponse{}, err
	}
	return result, nil
}

// getBastionShareableLinkHandleError handles the GetBastionShareableLink error response.
func (client *NetworkManagementClient) getBastionShareableLinkHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// BeginPutBastionShareableLink - Creates a Bastion Shareable Links for all the VMs specified in the request.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) BeginPutBastionShareableLink(ctx context.Context, resourceGroupName string, bastionHostName string, bslRequest BastionShareableLinkListRequest, options *NetworkManagementClientBeginPutBastionShareableLinkOptions) (NetworkManagementClientPutBastionShareableLinkPollerResponse, error) {
	resp, err := client.putBastionShareableLink(ctx, resourceGroupName, bastionHostName, bslRequest, options)
	if err != nil {
		return NetworkManagementClientPutBastionShareableLinkPollerResponse{}, err
	}
	result := NetworkManagementClientPutBastionShareableLinkPollerResponse{
		RawResponse: resp,
	}
	pt, err := armruntime.NewPoller("NetworkManagementClient.PutBastionShareableLink", "location", resp, client.pl, client.putBastionShareableLinkHandleError)
	if err != nil {
		return NetworkManagementClientPutBastionShareableLinkPollerResponse{}, err
	}
	result.Poller = &NetworkManagementClientPutBastionShareableLinkPoller{
		pt:     pt,
		client: client,
	}
	return result, nil
}

// PutBastionShareableLink - Creates a Bastion Shareable Links for all the VMs specified in the request.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) putBastionShareableLink(ctx context.Context, resourceGroupName string, bastionHostName string, bslRequest BastionShareableLinkListRequest, options *NetworkManagementClientBeginPutBastionShareableLinkOptions) (*http.Response, error) {
	req, err := client.putBastionShareableLinkCreateRequest(ctx, resourceGroupName, bastionHostName, bslRequest, options)
	if err != nil {
		return nil, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK, http.StatusAccepted) {
		return nil, client.putBastionShareableLinkHandleError(resp)
	}
	return resp, nil
}

// putBastionShareableLinkCreateRequest creates the PutBastionShareableLink request.
func (client *NetworkManagementClient) putBastionShareableLinkCreateRequest(ctx context.Context, resourceGroupName string, bastionHostName string, bslRequest BastionShareableLinkListRequest, options *NetworkManagementClientBeginPutBastionShareableLinkOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/bastionHosts/{bastionHostName}/createShareableLinks"
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if bastionHostName == "" {
		return nil, errors.New("parameter bastionHostName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{bastionHostName}", url.PathEscape(bastionHostName))
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, runtime.MarshalAsJSON(req, bslRequest)
}

// putBastionShareableLinkHandleResponse handles the PutBastionShareableLink response.
func (client *NetworkManagementClient) putBastionShareableLinkHandleResponse(resp *http.Response) (NetworkManagementClientPutBastionShareableLinkResponse, error) {
	result := NetworkManagementClientPutBastionShareableLinkResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.BastionShareableLinkListResult); err != nil {
		return NetworkManagementClientPutBastionShareableLinkResponse{}, err
	}
	return result, nil
}

// putBastionShareableLinkHandleError handles the PutBastionShareableLink error response.
func (client *NetworkManagementClient) putBastionShareableLinkHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}

// SupportedSecurityProviders - Gives the supported security providers for the virtual wan.
// If the operation fails it returns the *CloudError error type.
func (client *NetworkManagementClient) SupportedSecurityProviders(ctx context.Context, resourceGroupName string, virtualWANName string, options *NetworkManagementClientSupportedSecurityProvidersOptions) (NetworkManagementClientSupportedSecurityProvidersResponse, error) {
	req, err := client.supportedSecurityProvidersCreateRequest(ctx, resourceGroupName, virtualWANName, options)
	if err != nil {
		return NetworkManagementClientSupportedSecurityProvidersResponse{}, err
	}
	resp, err := client.pl.Do(req)
	if err != nil {
		return NetworkManagementClientSupportedSecurityProvidersResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return NetworkManagementClientSupportedSecurityProvidersResponse{}, client.supportedSecurityProvidersHandleError(resp)
	}
	return client.supportedSecurityProvidersHandleResponse(resp)
}

// supportedSecurityProvidersCreateRequest creates the SupportedSecurityProviders request.
func (client *NetworkManagementClient) supportedSecurityProvidersCreateRequest(ctx context.Context, resourceGroupName string, virtualWANName string, options *NetworkManagementClientSupportedSecurityProvidersOptions) (*policy.Request, error) {
	urlPath := "/subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.Network/virtualWans/{virtualWANName}/supportedSecurityProviders"
	if client.subscriptionID == "" {
		return nil, errors.New("parameter client.subscriptionID cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{subscriptionId}", url.PathEscape(client.subscriptionID))
	if resourceGroupName == "" {
		return nil, errors.New("parameter resourceGroupName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{resourceGroupName}", url.PathEscape(resourceGroupName))
	if virtualWANName == "" {
		return nil, errors.New("parameter virtualWANName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{virtualWANName}", url.PathEscape(virtualWANName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.ep, urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2021-05-01")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header.Set("Accept", "application/json")
	return req, nil
}

// supportedSecurityProvidersHandleResponse handles the SupportedSecurityProviders response.
func (client *NetworkManagementClient) supportedSecurityProvidersHandleResponse(resp *http.Response) (NetworkManagementClientSupportedSecurityProvidersResponse, error) {
	result := NetworkManagementClientSupportedSecurityProvidersResponse{RawResponse: resp}
	if err := runtime.UnmarshalAsJSON(resp, &result.VirtualWanSecurityProviders); err != nil {
		return NetworkManagementClientSupportedSecurityProvidersResponse{}, err
	}
	return result, nil
}

// supportedSecurityProvidersHandleError handles the SupportedSecurityProviders error response.
func (client *NetworkManagementClient) supportedSecurityProvidersHandleError(resp *http.Response) error {
	body, err := runtime.Payload(resp)
	if err != nil {
		return runtime.NewResponseError(err, resp)
	}
	errType := CloudError{raw: string(body)}
	if err := runtime.UnmarshalAsJSON(resp, &errType); err != nil {
		return runtime.NewResponseError(fmt.Errorf("%s\n%s", string(body), err), resp)
	}
	return runtime.NewResponseError(&errType, resp)
}
