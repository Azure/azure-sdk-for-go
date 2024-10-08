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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v4"
	"net/http"
	"net/url"
	"regexp"
)

// ThroughputPoolAccountsServer is a fake server for instances of the armcosmos.ThroughputPoolAccountsClient type.
type ThroughputPoolAccountsServer struct {
	// NewListPager is the fake for method ThroughputPoolAccountsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, throughputPoolName string, options *armcosmos.ThroughputPoolAccountsClientListOptions) (resp azfake.PagerResponder[armcosmos.ThroughputPoolAccountsClientListResponse])
}

// NewThroughputPoolAccountsServerTransport creates a new instance of ThroughputPoolAccountsServerTransport with the provided implementation.
// The returned ThroughputPoolAccountsServerTransport instance is connected to an instance of armcosmos.ThroughputPoolAccountsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewThroughputPoolAccountsServerTransport(srv *ThroughputPoolAccountsServer) *ThroughputPoolAccountsServerTransport {
	return &ThroughputPoolAccountsServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armcosmos.ThroughputPoolAccountsClientListResponse]](),
	}
}

// ThroughputPoolAccountsServerTransport connects instances of armcosmos.ThroughputPoolAccountsClient to instances of ThroughputPoolAccountsServer.
// Don't use this type directly, use NewThroughputPoolAccountsServerTransport instead.
type ThroughputPoolAccountsServerTransport struct {
	srv          *ThroughputPoolAccountsServer
	newListPager *tracker[azfake.PagerResponder[armcosmos.ThroughputPoolAccountsClientListResponse]]
}

// Do implements the policy.Transporter interface for ThroughputPoolAccountsServerTransport.
func (t *ThroughputPoolAccountsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ThroughputPoolAccountsClient.NewListPager":
		resp, err = t.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (t *ThroughputPoolAccountsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if t.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := t.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/throughputPools/(?P<throughputPoolName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/throughputPoolAccounts`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		throughputPoolNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("throughputPoolName")])
		if err != nil {
			return nil, err
		}
		resp := t.srv.NewListPager(resourceGroupNameParam, throughputPoolNameParam, nil)
		newListPager = &resp
		t.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armcosmos.ThroughputPoolAccountsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		t.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		t.newListPager.remove(req)
	}
	return resp, nil
}
