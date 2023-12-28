//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/maps/armmaps"
	"net/http"
	"regexp"
)

// Server is a fake server for instances of the armmaps.Client type.
type Server struct {
	// NewListOperationsPager is the fake for method Client.NewListOperationsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListOperationsPager func(options *armmaps.ClientListOperationsOptions) (resp azfake.PagerResponder[armmaps.ClientListOperationsResponse])

	// NewListSubscriptionOperationsPager is the fake for method Client.NewListSubscriptionOperationsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListSubscriptionOperationsPager func(options *armmaps.ClientListSubscriptionOperationsOptions) (resp azfake.PagerResponder[armmaps.ClientListSubscriptionOperationsResponse])
}

// NewServerTransport creates a new instance of ServerTransport with the provided implementation.
// The returned ServerTransport instance is connected to an instance of armmaps.Client via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewServerTransport(srv *Server) *ServerTransport {
	return &ServerTransport{
		srv:                                srv,
		newListOperationsPager:             newTracker[azfake.PagerResponder[armmaps.ClientListOperationsResponse]](),
		newListSubscriptionOperationsPager: newTracker[azfake.PagerResponder[armmaps.ClientListSubscriptionOperationsResponse]](),
	}
}

// ServerTransport connects instances of armmaps.Client to instances of Server.
// Don't use this type directly, use NewServerTransport instead.
type ServerTransport struct {
	srv                                *Server
	newListOperationsPager             *tracker[azfake.PagerResponder[armmaps.ClientListOperationsResponse]]
	newListSubscriptionOperationsPager *tracker[azfake.PagerResponder[armmaps.ClientListSubscriptionOperationsResponse]]
}

// Do implements the policy.Transporter interface for ServerTransport.
func (s *ServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "Client.NewListOperationsPager":
		resp, err = s.dispatchNewListOperationsPager(req)
	case "Client.NewListSubscriptionOperationsPager":
		resp, err = s.dispatchNewListSubscriptionOperationsPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ServerTransport) dispatchNewListOperationsPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListOperationsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListOperationsPager not implemented")}
	}
	newListOperationsPager := s.newListOperationsPager.get(req)
	if newListOperationsPager == nil {
		resp := s.srv.NewListOperationsPager(nil)
		newListOperationsPager = &resp
		s.newListOperationsPager.add(req, newListOperationsPager)
		server.PagerResponderInjectNextLinks(newListOperationsPager, req, func(page *armmaps.ClientListOperationsResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListOperationsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListOperationsPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListOperationsPager) {
		s.newListOperationsPager.remove(req)
	}
	return resp, nil
}

func (s *ServerTransport) dispatchNewListSubscriptionOperationsPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListSubscriptionOperationsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListSubscriptionOperationsPager not implemented")}
	}
	newListSubscriptionOperationsPager := s.newListSubscriptionOperationsPager.get(req)
	if newListSubscriptionOperationsPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Maps/operations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := s.srv.NewListSubscriptionOperationsPager(nil)
		newListSubscriptionOperationsPager = &resp
		s.newListSubscriptionOperationsPager.add(req, newListSubscriptionOperationsPager)
		server.PagerResponderInjectNextLinks(newListSubscriptionOperationsPager, req, func(page *armmaps.ClientListSubscriptionOperationsResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListSubscriptionOperationsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListSubscriptionOperationsPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListSubscriptionOperationsPager) {
		s.newListSubscriptionOperationsPager.remove(req)
	}
	return resp, nil
}