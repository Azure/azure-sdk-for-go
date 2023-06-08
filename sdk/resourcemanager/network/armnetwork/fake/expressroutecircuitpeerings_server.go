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

// ExpressRouteCircuitPeeringsServer is a fake server for instances of the armnetwork.ExpressRouteCircuitPeeringsClient type.
type ExpressRouteCircuitPeeringsServer struct {
	// BeginCreateOrUpdate is the fake for method ExpressRouteCircuitPeeringsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, circuitName string, peeringName string, peeringParameters armnetwork.ExpressRouteCircuitPeering, options *armnetwork.ExpressRouteCircuitPeeringsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armnetwork.ExpressRouteCircuitPeeringsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method ExpressRouteCircuitPeeringsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, circuitName string, peeringName string, options *armnetwork.ExpressRouteCircuitPeeringsClientBeginDeleteOptions) (resp azfake.PollerResponder[armnetwork.ExpressRouteCircuitPeeringsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method ExpressRouteCircuitPeeringsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, circuitName string, peeringName string, options *armnetwork.ExpressRouteCircuitPeeringsClientGetOptions) (resp azfake.Responder[armnetwork.ExpressRouteCircuitPeeringsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method ExpressRouteCircuitPeeringsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, circuitName string, options *armnetwork.ExpressRouteCircuitPeeringsClientListOptions) (resp azfake.PagerResponder[armnetwork.ExpressRouteCircuitPeeringsClientListResponse])
}

// NewExpressRouteCircuitPeeringsServerTransport creates a new instance of ExpressRouteCircuitPeeringsServerTransport with the provided implementation.
// The returned ExpressRouteCircuitPeeringsServerTransport instance is connected to an instance of armnetwork.ExpressRouteCircuitPeeringsClient by way of the
// undefined.Transporter field.
func NewExpressRouteCircuitPeeringsServerTransport(srv *ExpressRouteCircuitPeeringsServer) *ExpressRouteCircuitPeeringsServerTransport {
	return &ExpressRouteCircuitPeeringsServerTransport{srv: srv}
}

// ExpressRouteCircuitPeeringsServerTransport connects instances of armnetwork.ExpressRouteCircuitPeeringsClient to instances of ExpressRouteCircuitPeeringsServer.
// Don't use this type directly, use NewExpressRouteCircuitPeeringsServerTransport instead.
type ExpressRouteCircuitPeeringsServerTransport struct {
	srv                 *ExpressRouteCircuitPeeringsServer
	beginCreateOrUpdate *azfake.PollerResponder[armnetwork.ExpressRouteCircuitPeeringsClientCreateOrUpdateResponse]
	beginDelete         *azfake.PollerResponder[armnetwork.ExpressRouteCircuitPeeringsClientDeleteResponse]
	newListPager        *azfake.PagerResponder[armnetwork.ExpressRouteCircuitPeeringsClientListResponse]
}

// Do implements the policy.Transporter interface for ExpressRouteCircuitPeeringsServerTransport.
func (e *ExpressRouteCircuitPeeringsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ExpressRouteCircuitPeeringsClient.BeginCreateOrUpdate":
		resp, err = e.dispatchBeginCreateOrUpdate(req)
	case "ExpressRouteCircuitPeeringsClient.BeginDelete":
		resp, err = e.dispatchBeginDelete(req)
	case "ExpressRouteCircuitPeeringsClient.Get":
		resp, err = e.dispatchGet(req)
	case "ExpressRouteCircuitPeeringsClient.NewListPager":
		resp, err = e.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (e *ExpressRouteCircuitPeeringsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if e.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("method BeginCreateOrUpdate not implemented")}
	}
	if e.beginCreateOrUpdate == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/expressRouteCircuits/(?P<circuitName>[a-zA-Z0-9-_]+)/peerings/(?P<peeringName>[a-zA-Z0-9-_]+)"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.ExpressRouteCircuitPeering](req)
		if err != nil {
			return nil, err
		}
		respr, errRespr := e.srv.BeginCreateOrUpdate(req.Context(), matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("circuitName")], matches[regex.SubexpIndex("peeringName")], body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		e.beginCreateOrUpdate = &respr
	}

	resp, err := server.PollerResponderNext(e.beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(e.beginCreateOrUpdate) {
		e.beginCreateOrUpdate = nil
	}

	return resp, nil
}

func (e *ExpressRouteCircuitPeeringsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if e.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("method BeginDelete not implemented")}
	}
	if e.beginDelete == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/expressRouteCircuits/(?P<circuitName>[a-zA-Z0-9-_]+)/peerings/(?P<peeringName>[a-zA-Z0-9-_]+)"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		respr, errRespr := e.srv.BeginDelete(req.Context(), matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("circuitName")], matches[regex.SubexpIndex("peeringName")], nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		e.beginDelete = &respr
	}

	resp, err := server.PollerResponderNext(e.beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(e.beginDelete) {
		e.beginDelete = nil
	}

	return resp, nil
}

func (e *ExpressRouteCircuitPeeringsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if e.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("method Get not implemented")}
	}
	const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/expressRouteCircuits/(?P<circuitName>[a-zA-Z0-9-_]+)/peerings/(?P<peeringName>[a-zA-Z0-9-_]+)"
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.Path)
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	respr, errRespr := e.srv.Get(req.Context(), matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("circuitName")], matches[regex.SubexpIndex("peeringName")], nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ExpressRouteCircuitPeering, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *ExpressRouteCircuitPeeringsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if e.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("method NewListPager not implemented")}
	}
	if e.newListPager == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/expressRouteCircuits/(?P<circuitName>[a-zA-Z0-9-_]+)/peerings"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := e.srv.NewListPager(matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("circuitName")], nil)
		e.newListPager = &resp
		server.PagerResponderInjectNextLinks(e.newListPager, req, func(page *armnetwork.ExpressRouteCircuitPeeringsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(e.newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(e.newListPager) {
		e.newListPager = nil
	}
	return resp, nil
}
