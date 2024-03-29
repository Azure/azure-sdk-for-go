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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/newrelic/armnewrelicobservability"
	"net/http"
	"net/url"
	"regexp"
)

// AccountsServer is a fake server for instances of the armnewrelicobservability.AccountsClient type.
type AccountsServer struct {
	// NewListPager is the fake for method AccountsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(userEmail string, location string, options *armnewrelicobservability.AccountsClientListOptions) (resp azfake.PagerResponder[armnewrelicobservability.AccountsClientListResponse])
}

// NewAccountsServerTransport creates a new instance of AccountsServerTransport with the provided implementation.
// The returned AccountsServerTransport instance is connected to an instance of armnewrelicobservability.AccountsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAccountsServerTransport(srv *AccountsServer) *AccountsServerTransport {
	return &AccountsServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armnewrelicobservability.AccountsClientListResponse]](),
	}
}

// AccountsServerTransport connects instances of armnewrelicobservability.AccountsClient to instances of AccountsServer.
// Don't use this type directly, use NewAccountsServerTransport instead.
type AccountsServerTransport struct {
	srv          *AccountsServer
	newListPager *tracker[azfake.PagerResponder[armnewrelicobservability.AccountsClientListResponse]]
}

// Do implements the policy.Transporter interface for AccountsServerTransport.
func (a *AccountsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "AccountsClient.NewListPager":
		resp, err = a.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AccountsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := a.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/NewRelic\.Observability/accounts`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		userEmailParam, err := url.QueryUnescape(qp.Get("userEmail"))
		if err != nil {
			return nil, err
		}
		locationParam, err := url.QueryUnescape(qp.Get("location"))
		if err != nil {
			return nil, err
		}
		resp := a.srv.NewListPager(userEmailParam, locationParam, nil)
		newListPager = &resp
		a.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armnewrelicobservability.AccountsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		a.newListPager.remove(req)
	}
	return resp, nil
}
