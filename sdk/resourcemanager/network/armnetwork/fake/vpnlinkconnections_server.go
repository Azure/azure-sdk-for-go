//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"context"
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v3"
	"net/http"
	"regexp"
)

// VPNLinkConnectionsServer is a fake server for instances of the armnetwork.VPNLinkConnectionsClient type.
type VPNLinkConnectionsServer struct {
	// BeginGetIkeSas is the fake for method VPNLinkConnectionsClient.BeginGetIkeSas
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginGetIkeSas func(ctx context.Context, resourceGroupName string, gatewayName string, connectionName string, linkConnectionName string, options *armnetwork.VPNLinkConnectionsClientBeginGetIkeSasOptions) (resp azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientGetIkeSasResponse], errResp azfake.ErrorResponder)

	// NewListByVPNConnectionPager is the fake for method VPNLinkConnectionsClient.NewListByVPNConnectionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByVPNConnectionPager func(resourceGroupName string, gatewayName string, connectionName string, options *armnetwork.VPNLinkConnectionsClientListByVPNConnectionOptions) (resp azfake.PagerResponder[armnetwork.VPNLinkConnectionsClientListByVPNConnectionResponse])

	// BeginResetConnection is the fake for method VPNLinkConnectionsClient.BeginResetConnection
	// HTTP status codes to indicate success: http.StatusAccepted
	BeginResetConnection func(ctx context.Context, resourceGroupName string, gatewayName string, connectionName string, linkConnectionName string, options *armnetwork.VPNLinkConnectionsClientBeginResetConnectionOptions) (resp azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientResetConnectionResponse], errResp azfake.ErrorResponder)
}

// NewVPNLinkConnectionsServerTransport creates a new instance of VPNLinkConnectionsServerTransport with the provided implementation.
// The returned VPNLinkConnectionsServerTransport instance is connected to an instance of armnetwork.VPNLinkConnectionsClient by way of the
// undefined.Transporter field.
func NewVPNLinkConnectionsServerTransport(srv *VPNLinkConnectionsServer) *VPNLinkConnectionsServerTransport {
	return &VPNLinkConnectionsServerTransport{srv: srv}
}

// VPNLinkConnectionsServerTransport connects instances of armnetwork.VPNLinkConnectionsClient to instances of VPNLinkConnectionsServer.
// Don't use this type directly, use NewVPNLinkConnectionsServerTransport instead.
type VPNLinkConnectionsServerTransport struct {
	srv                         *VPNLinkConnectionsServer
	beginGetIkeSas              *azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientGetIkeSasResponse]
	newListByVPNConnectionPager *azfake.PagerResponder[armnetwork.VPNLinkConnectionsClientListByVPNConnectionResponse]
	beginResetConnection        *azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientResetConnectionResponse]
}

// Do implements the policy.Transporter interface for VPNLinkConnectionsServerTransport.
func (v *VPNLinkConnectionsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "VPNLinkConnectionsClient.BeginGetIkeSas":
		resp, err = v.dispatchBeginGetIkeSas(req)
	case "VPNLinkConnectionsClient.NewListByVPNConnectionPager":
		resp, err = v.dispatchNewListByVPNConnectionPager(req)
	case "VPNLinkConnectionsClient.BeginResetConnection":
		resp, err = v.dispatchBeginResetConnection(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *VPNLinkConnectionsServerTransport) dispatchBeginGetIkeSas(req *http.Request) (*http.Response, error) {
	if v.srv.BeginGetIkeSas == nil {
		return nil, &nonRetriableError{errors.New("method BeginGetIkeSas not implemented")}
	}
	if v.beginGetIkeSas == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/vpnGateways/(?P<gatewayName>[a-zA-Z0-9-_]+)/vpnConnections/(?P<connectionName>[a-zA-Z0-9-_]+)/vpnLinkConnections/(?P<linkConnectionName>[a-zA-Z0-9-_]+)/getikesas"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		respr, errRespr := v.srv.BeginGetIkeSas(req.Context(), matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("gatewayName")], matches[regex.SubexpIndex("connectionName")], matches[regex.SubexpIndex("linkConnectionName")], nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		v.beginGetIkeSas = &respr
	}

	resp, err := server.PollerResponderNext(v.beginGetIkeSas, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(v.beginGetIkeSas) {
		v.beginGetIkeSas = nil
	}

	return resp, nil
}

func (v *VPNLinkConnectionsServerTransport) dispatchNewListByVPNConnectionPager(req *http.Request) (*http.Response, error) {
	if v.srv.NewListByVPNConnectionPager == nil {
		return nil, &nonRetriableError{errors.New("method NewListByVPNConnectionPager not implemented")}
	}
	if v.newListByVPNConnectionPager == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/vpnGateways/(?P<gatewayName>[a-zA-Z0-9-_]+)/vpnConnections/(?P<connectionName>[a-zA-Z0-9-_]+)/vpnLinkConnections"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := v.srv.NewListByVPNConnectionPager(matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("gatewayName")], matches[regex.SubexpIndex("connectionName")], nil)
		v.newListByVPNConnectionPager = &resp
		server.PagerResponderInjectNextLinks(v.newListByVPNConnectionPager, req, func(page *armnetwork.VPNLinkConnectionsClientListByVPNConnectionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(v.newListByVPNConnectionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(v.newListByVPNConnectionPager) {
		v.newListByVPNConnectionPager = nil
	}
	return resp, nil
}

func (v *VPNLinkConnectionsServerTransport) dispatchBeginResetConnection(req *http.Request) (*http.Response, error) {
	if v.srv.BeginResetConnection == nil {
		return nil, &nonRetriableError{errors.New("method BeginResetConnection not implemented")}
	}
	if v.beginResetConnection == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/vpnGateways/(?P<gatewayName>[a-zA-Z0-9-_]+)/vpnConnections/(?P<connectionName>[a-zA-Z0-9-_]+)/vpnLinkConnections/(?P<linkConnectionName>[a-zA-Z0-9-_]+)/resetconnection"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		respr, errRespr := v.srv.BeginResetConnection(req.Context(), matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("gatewayName")], matches[regex.SubexpIndex("connectionName")], matches[regex.SubexpIndex("linkConnectionName")], nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		v.beginResetConnection = &respr
	}

	resp, err := server.PollerResponderNext(v.beginResetConnection, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusAccepted}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(v.beginResetConnection) {
		v.beginResetConnection = nil
	}

	return resp, nil
}
