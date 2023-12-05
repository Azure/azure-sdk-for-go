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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v5"
	"net/http"
	"net/url"
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
// The returned VPNLinkConnectionsServerTransport instance is connected to an instance of armnetwork.VPNLinkConnectionsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewVPNLinkConnectionsServerTransport(srv *VPNLinkConnectionsServer) *VPNLinkConnectionsServerTransport {
	return &VPNLinkConnectionsServerTransport{
		srv:                         srv,
		beginGetIkeSas:              newTracker[azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientGetIkeSasResponse]](),
		newListByVPNConnectionPager: newTracker[azfake.PagerResponder[armnetwork.VPNLinkConnectionsClientListByVPNConnectionResponse]](),
		beginResetConnection:        newTracker[azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientResetConnectionResponse]](),
	}
}

// VPNLinkConnectionsServerTransport connects instances of armnetwork.VPNLinkConnectionsClient to instances of VPNLinkConnectionsServer.
// Don't use this type directly, use NewVPNLinkConnectionsServerTransport instead.
type VPNLinkConnectionsServerTransport struct {
	srv                         *VPNLinkConnectionsServer
	beginGetIkeSas              *tracker[azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientGetIkeSasResponse]]
	newListByVPNConnectionPager *tracker[azfake.PagerResponder[armnetwork.VPNLinkConnectionsClientListByVPNConnectionResponse]]
	beginResetConnection        *tracker[azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientResetConnectionResponse]]
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
		return nil, &nonRetriableError{errors.New("fake for method BeginGetIkeSas not implemented")}
	}
	beginGetIkeSas := v.beginGetIkeSas.get(req)
	if beginGetIkeSas == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/vpnGateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnConnections/(?P<connectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnLinkConnections/(?P<linkConnectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/getikesas`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		gatewayNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("gatewayName")])
		if err != nil {
			return nil, err
		}
		connectionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("connectionName")])
		if err != nil {
			return nil, err
		}
		linkConnectionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("linkConnectionName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginGetIkeSas(req.Context(), resourceGroupNameParam, gatewayNameParam, connectionNameParam, linkConnectionNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginGetIkeSas = &respr
		v.beginGetIkeSas.add(req, beginGetIkeSas)
	}

	resp, err := server.PollerResponderNext(beginGetIkeSas, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginGetIkeSas.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginGetIkeSas) {
		v.beginGetIkeSas.remove(req)
	}

	return resp, nil
}

func (v *VPNLinkConnectionsServerTransport) dispatchNewListByVPNConnectionPager(req *http.Request) (*http.Response, error) {
	if v.srv.NewListByVPNConnectionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByVPNConnectionPager not implemented")}
	}
	newListByVPNConnectionPager := v.newListByVPNConnectionPager.get(req)
	if newListByVPNConnectionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/vpnGateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnConnections/(?P<connectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnLinkConnections`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		gatewayNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("gatewayName")])
		if err != nil {
			return nil, err
		}
		connectionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("connectionName")])
		if err != nil {
			return nil, err
		}
		resp := v.srv.NewListByVPNConnectionPager(resourceGroupNameParam, gatewayNameParam, connectionNameParam, nil)
		newListByVPNConnectionPager = &resp
		v.newListByVPNConnectionPager.add(req, newListByVPNConnectionPager)
		server.PagerResponderInjectNextLinks(newListByVPNConnectionPager, req, func(page *armnetwork.VPNLinkConnectionsClientListByVPNConnectionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByVPNConnectionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		v.newListByVPNConnectionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByVPNConnectionPager) {
		v.newListByVPNConnectionPager.remove(req)
	}
	return resp, nil
}

func (v *VPNLinkConnectionsServerTransport) dispatchBeginResetConnection(req *http.Request) (*http.Response, error) {
	if v.srv.BeginResetConnection == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginResetConnection not implemented")}
	}
	beginResetConnection := v.beginResetConnection.get(req)
	if beginResetConnection == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/vpnGateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnConnections/(?P<connectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnLinkConnections/(?P<linkConnectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resetconnection`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		gatewayNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("gatewayName")])
		if err != nil {
			return nil, err
		}
		connectionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("connectionName")])
		if err != nil {
			return nil, err
		}
		linkConnectionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("linkConnectionName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginResetConnection(req.Context(), resourceGroupNameParam, gatewayNameParam, connectionNameParam, linkConnectionNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginResetConnection = &respr
		v.beginResetConnection.add(req, beginResetConnection)
	}

	resp, err := server.PollerResponderNext(beginResetConnection, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusAccepted}, resp.StatusCode) {
		v.beginResetConnection.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginResetConnection) {
		v.beginResetConnection.remove(req)
	}

	return resp, nil
}
