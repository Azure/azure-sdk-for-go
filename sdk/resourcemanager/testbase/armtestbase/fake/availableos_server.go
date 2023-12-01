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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/testbase/armtestbase"
	"net/http"
	"net/url"
	"regexp"
)

// AvailableOSServer is a fake server for instances of the armtestbase.AvailableOSClient type.
type AvailableOSServer struct {
	// Get is the fake for method AvailableOSClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, testBaseAccountName string, availableOSResourceName string, options *armtestbase.AvailableOSClientGetOptions) (resp azfake.Responder[armtestbase.AvailableOSClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method AvailableOSClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, testBaseAccountName string, osUpdateType armtestbase.OsUpdateType, options *armtestbase.AvailableOSClientListOptions) (resp azfake.PagerResponder[armtestbase.AvailableOSClientListResponse])
}

// NewAvailableOSServerTransport creates a new instance of AvailableOSServerTransport with the provided implementation.
// The returned AvailableOSServerTransport instance is connected to an instance of armtestbase.AvailableOSClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAvailableOSServerTransport(srv *AvailableOSServer) *AvailableOSServerTransport {
	return &AvailableOSServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armtestbase.AvailableOSClientListResponse]](),
	}
}

// AvailableOSServerTransport connects instances of armtestbase.AvailableOSClient to instances of AvailableOSServer.
// Don't use this type directly, use NewAvailableOSServerTransport instead.
type AvailableOSServerTransport struct {
	srv          *AvailableOSServer
	newListPager *tracker[azfake.PagerResponder[armtestbase.AvailableOSClientListResponse]]
}

// Do implements the policy.Transporter interface for AvailableOSServerTransport.
func (a *AvailableOSServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "AvailableOSClient.Get":
		resp, err = a.dispatchGet(req)
	case "AvailableOSClient.NewListPager":
		resp, err = a.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AvailableOSServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if a.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.TestBase/testBaseAccounts/(?P<testBaseAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/availableOSs/(?P<availableOSResourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	testBaseAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("testBaseAccountName")])
	if err != nil {
		return nil, err
	}
	availableOSResourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("availableOSResourceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Get(req.Context(), resourceGroupNameParam, testBaseAccountNameParam, availableOSResourceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AvailableOSResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AvailableOSServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := a.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.TestBase/testBaseAccounts/(?P<testBaseAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/availableOSs`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		testBaseAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("testBaseAccountName")])
		if err != nil {
			return nil, err
		}
		osUpdateTypeParam, err := parseWithCast(qp.Get("osUpdateType"), func(v string) (armtestbase.OsUpdateType, error) {
			p, unescapeErr := url.QueryUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armtestbase.OsUpdateType(p), nil
		})
		if err != nil {
			return nil, err
		}
		resp := a.srv.NewListPager(resourceGroupNameParam, testBaseAccountNameParam, osUpdateTypeParam, nil)
		newListPager = &resp
		a.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armtestbase.AvailableOSClientListResponse, createLink func() string) {
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