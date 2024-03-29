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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/customerinsights/armcustomerinsights"
	"net/http"
	"net/url"
	"regexp"
)

// RolesServer is a fake server for instances of the armcustomerinsights.RolesClient type.
type RolesServer struct {
	// NewListByHubPager is the fake for method RolesClient.NewListByHubPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByHubPager func(resourceGroupName string, hubName string, options *armcustomerinsights.RolesClientListByHubOptions) (resp azfake.PagerResponder[armcustomerinsights.RolesClientListByHubResponse])
}

// NewRolesServerTransport creates a new instance of RolesServerTransport with the provided implementation.
// The returned RolesServerTransport instance is connected to an instance of armcustomerinsights.RolesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewRolesServerTransport(srv *RolesServer) *RolesServerTransport {
	return &RolesServerTransport{
		srv:               srv,
		newListByHubPager: newTracker[azfake.PagerResponder[armcustomerinsights.RolesClientListByHubResponse]](),
	}
}

// RolesServerTransport connects instances of armcustomerinsights.RolesClient to instances of RolesServer.
// Don't use this type directly, use NewRolesServerTransport instead.
type RolesServerTransport struct {
	srv               *RolesServer
	newListByHubPager *tracker[azfake.PagerResponder[armcustomerinsights.RolesClientListByHubResponse]]
}

// Do implements the policy.Transporter interface for RolesServerTransport.
func (r *RolesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "RolesClient.NewListByHubPager":
		resp, err = r.dispatchNewListByHubPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *RolesServerTransport) dispatchNewListByHubPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListByHubPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByHubPager not implemented")}
	}
	newListByHubPager := r.newListByHubPager.get(req)
	if newListByHubPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.CustomerInsights/hubs/(?P<hubName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/roles`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		hubNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("hubName")])
		if err != nil {
			return nil, err
		}
		resp := r.srv.NewListByHubPager(resourceGroupNameParam, hubNameParam, nil)
		newListByHubPager = &resp
		r.newListByHubPager.add(req, newListByHubPager)
		server.PagerResponderInjectNextLinks(newListByHubPager, req, func(page *armcustomerinsights.RolesClientListByHubResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByHubPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.newListByHubPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByHubPager) {
		r.newListByHubPager.remove(req)
	}
	return resp, nil
}
