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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"net/http"
	"net/url"
	"regexp"
)

// VPNLinkConnectionsServer is a fake server for instances of the armnetwork.VPNLinkConnectionsClient type.
type VPNLinkConnectionsServer struct {
	// NewGetAllSharedKeysPager is the fake for method VPNLinkConnectionsClient.NewGetAllSharedKeysPager
	// HTTP status codes to indicate success: http.StatusOK
	NewGetAllSharedKeysPager func(resourceGroupName string, gatewayName string, connectionName string, linkConnectionName string, options *armnetwork.VPNLinkConnectionsClientGetAllSharedKeysOptions) (resp azfake.PagerResponder[armnetwork.VPNLinkConnectionsClientGetAllSharedKeysResponse])

	// GetDefaultSharedKey is the fake for method VPNLinkConnectionsClient.GetDefaultSharedKey
	// HTTP status codes to indicate success: http.StatusOK
	GetDefaultSharedKey func(ctx context.Context, resourceGroupName string, gatewayName string, connectionName string, linkConnectionName string, options *armnetwork.VPNLinkConnectionsClientGetDefaultSharedKeyOptions) (resp azfake.Responder[armnetwork.VPNLinkConnectionsClientGetDefaultSharedKeyResponse], errResp azfake.ErrorResponder)

	// BeginGetIkeSas is the fake for method VPNLinkConnectionsClient.BeginGetIkeSas
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginGetIkeSas func(ctx context.Context, resourceGroupName string, gatewayName string, connectionName string, linkConnectionName string, options *armnetwork.VPNLinkConnectionsClientBeginGetIkeSasOptions) (resp azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientGetIkeSasResponse], errResp azfake.ErrorResponder)

	// NewListByVPNConnectionPager is the fake for method VPNLinkConnectionsClient.NewListByVPNConnectionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByVPNConnectionPager func(resourceGroupName string, gatewayName string, connectionName string, options *armnetwork.VPNLinkConnectionsClientListByVPNConnectionOptions) (resp azfake.PagerResponder[armnetwork.VPNLinkConnectionsClientListByVPNConnectionResponse])

	// ListDefaultSharedKey is the fake for method VPNLinkConnectionsClient.ListDefaultSharedKey
	// HTTP status codes to indicate success: http.StatusOK
	ListDefaultSharedKey func(ctx context.Context, resourceGroupName string, gatewayName string, connectionName string, linkConnectionName string, options *armnetwork.VPNLinkConnectionsClientListDefaultSharedKeyOptions) (resp azfake.Responder[armnetwork.VPNLinkConnectionsClientListDefaultSharedKeyResponse], errResp azfake.ErrorResponder)

	// BeginResetConnection is the fake for method VPNLinkConnectionsClient.BeginResetConnection
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginResetConnection func(ctx context.Context, resourceGroupName string, gatewayName string, connectionName string, linkConnectionName string, options *armnetwork.VPNLinkConnectionsClientBeginResetConnectionOptions) (resp azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientResetConnectionResponse], errResp azfake.ErrorResponder)

	// BeginSetOrInitDefaultSharedKey is the fake for method VPNLinkConnectionsClient.BeginSetOrInitDefaultSharedKey
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginSetOrInitDefaultSharedKey func(ctx context.Context, resourceGroupName string, gatewayName string, connectionName string, linkConnectionName string, connectionSharedKeyParameters armnetwork.ConnectionSharedKeyResult, options *armnetwork.VPNLinkConnectionsClientBeginSetOrInitDefaultSharedKeyOptions) (resp azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientSetOrInitDefaultSharedKeyResponse], errResp azfake.ErrorResponder)
}

// NewVPNLinkConnectionsServerTransport creates a new instance of VPNLinkConnectionsServerTransport with the provided implementation.
// The returned VPNLinkConnectionsServerTransport instance is connected to an instance of armnetwork.VPNLinkConnectionsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewVPNLinkConnectionsServerTransport(srv *VPNLinkConnectionsServer) *VPNLinkConnectionsServerTransport {
	return &VPNLinkConnectionsServerTransport{
		srv:                            srv,
		newGetAllSharedKeysPager:       newTracker[azfake.PagerResponder[armnetwork.VPNLinkConnectionsClientGetAllSharedKeysResponse]](),
		beginGetIkeSas:                 newTracker[azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientGetIkeSasResponse]](),
		newListByVPNConnectionPager:    newTracker[azfake.PagerResponder[armnetwork.VPNLinkConnectionsClientListByVPNConnectionResponse]](),
		beginResetConnection:           newTracker[azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientResetConnectionResponse]](),
		beginSetOrInitDefaultSharedKey: newTracker[azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientSetOrInitDefaultSharedKeyResponse]](),
	}
}

// VPNLinkConnectionsServerTransport connects instances of armnetwork.VPNLinkConnectionsClient to instances of VPNLinkConnectionsServer.
// Don't use this type directly, use NewVPNLinkConnectionsServerTransport instead.
type VPNLinkConnectionsServerTransport struct {
	srv                            *VPNLinkConnectionsServer
	newGetAllSharedKeysPager       *tracker[azfake.PagerResponder[armnetwork.VPNLinkConnectionsClientGetAllSharedKeysResponse]]
	beginGetIkeSas                 *tracker[azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientGetIkeSasResponse]]
	newListByVPNConnectionPager    *tracker[azfake.PagerResponder[armnetwork.VPNLinkConnectionsClientListByVPNConnectionResponse]]
	beginResetConnection           *tracker[azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientResetConnectionResponse]]
	beginSetOrInitDefaultSharedKey *tracker[azfake.PollerResponder[armnetwork.VPNLinkConnectionsClientSetOrInitDefaultSharedKeyResponse]]
}

// Do implements the policy.Transporter interface for VPNLinkConnectionsServerTransport.
func (v *VPNLinkConnectionsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return v.dispatchToMethodFake(req, method)
}

func (v *VPNLinkConnectionsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if vpnLinkConnectionsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = vpnLinkConnectionsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "VPNLinkConnectionsClient.NewGetAllSharedKeysPager":
				res.resp, res.err = v.dispatchNewGetAllSharedKeysPager(req)
			case "VPNLinkConnectionsClient.GetDefaultSharedKey":
				res.resp, res.err = v.dispatchGetDefaultSharedKey(req)
			case "VPNLinkConnectionsClient.BeginGetIkeSas":
				res.resp, res.err = v.dispatchBeginGetIkeSas(req)
			case "VPNLinkConnectionsClient.NewListByVPNConnectionPager":
				res.resp, res.err = v.dispatchNewListByVPNConnectionPager(req)
			case "VPNLinkConnectionsClient.ListDefaultSharedKey":
				res.resp, res.err = v.dispatchListDefaultSharedKey(req)
			case "VPNLinkConnectionsClient.BeginResetConnection":
				res.resp, res.err = v.dispatchBeginResetConnection(req)
			case "VPNLinkConnectionsClient.BeginSetOrInitDefaultSharedKey":
				res.resp, res.err = v.dispatchBeginSetOrInitDefaultSharedKey(req)
			default:
				res.err = fmt.Errorf("unhandled API %s", method)
			}

		}
		select {
		case resultChan <- res:
		case <-req.Context().Done():
		}
	}()

	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

func (v *VPNLinkConnectionsServerTransport) dispatchNewGetAllSharedKeysPager(req *http.Request) (*http.Response, error) {
	if v.srv.NewGetAllSharedKeysPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewGetAllSharedKeysPager not implemented")}
	}
	newGetAllSharedKeysPager := v.newGetAllSharedKeysPager.get(req)
	if newGetAllSharedKeysPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/vpnGateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnConnections/(?P<connectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnLinkConnections/(?P<linkConnectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/sharedKeys`
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
		resp := v.srv.NewGetAllSharedKeysPager(resourceGroupNameParam, gatewayNameParam, connectionNameParam, linkConnectionNameParam, nil)
		newGetAllSharedKeysPager = &resp
		v.newGetAllSharedKeysPager.add(req, newGetAllSharedKeysPager)
		server.PagerResponderInjectNextLinks(newGetAllSharedKeysPager, req, func(page *armnetwork.VPNLinkConnectionsClientGetAllSharedKeysResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newGetAllSharedKeysPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		v.newGetAllSharedKeysPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newGetAllSharedKeysPager) {
		v.newGetAllSharedKeysPager.remove(req)
	}
	return resp, nil
}

func (v *VPNLinkConnectionsServerTransport) dispatchGetDefaultSharedKey(req *http.Request) (*http.Response, error) {
	if v.srv.GetDefaultSharedKey == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetDefaultSharedKey not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/vpnGateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnConnections/(?P<connectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnLinkConnections/(?P<linkConnectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/sharedKeys/default`
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
	respr, errRespr := v.srv.GetDefaultSharedKey(req.Context(), resourceGroupNameParam, gatewayNameParam, connectionNameParam, linkConnectionNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ConnectionSharedKeyResult, req)
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

func (v *VPNLinkConnectionsServerTransport) dispatchListDefaultSharedKey(req *http.Request) (*http.Response, error) {
	if v.srv.ListDefaultSharedKey == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListDefaultSharedKey not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/vpnGateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnConnections/(?P<connectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnLinkConnections/(?P<linkConnectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/sharedKeys/default/listSharedKey`
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
	respr, errRespr := v.srv.ListDefaultSharedKey(req.Context(), resourceGroupNameParam, gatewayNameParam, connectionNameParam, linkConnectionNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ConnectionSharedKeyResult, req)
	if err != nil {
		return nil, err
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

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		v.beginResetConnection.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginResetConnection) {
		v.beginResetConnection.remove(req)
	}

	return resp, nil
}

func (v *VPNLinkConnectionsServerTransport) dispatchBeginSetOrInitDefaultSharedKey(req *http.Request) (*http.Response, error) {
	if v.srv.BeginSetOrInitDefaultSharedKey == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginSetOrInitDefaultSharedKey not implemented")}
	}
	beginSetOrInitDefaultSharedKey := v.beginSetOrInitDefaultSharedKey.get(req)
	if beginSetOrInitDefaultSharedKey == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/vpnGateways/(?P<gatewayName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnConnections/(?P<connectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vpnLinkConnections/(?P<linkConnectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/sharedKeys/default`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.ConnectionSharedKeyResult](req)
		if err != nil {
			return nil, err
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
		respr, errRespr := v.srv.BeginSetOrInitDefaultSharedKey(req.Context(), resourceGroupNameParam, gatewayNameParam, connectionNameParam, linkConnectionNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginSetOrInitDefaultSharedKey = &respr
		v.beginSetOrInitDefaultSharedKey.add(req, beginSetOrInitDefaultSharedKey)
	}

	resp, err := server.PollerResponderNext(beginSetOrInitDefaultSharedKey, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		v.beginSetOrInitDefaultSharedKey.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginSetOrInitDefaultSharedKey) {
		v.beginSetOrInitDefaultSharedKey.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to VPNLinkConnectionsServerTransport
var vpnLinkConnectionsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
