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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v4"
	"net/http"
	"net/url"
	"regexp"
)

// ChaosFaultServer is a fake server for instances of the armcosmos.ChaosFaultClient type.
type ChaosFaultServer struct {
	// BeginEnableDisable is the fake for method ChaosFaultClient.BeginEnableDisable
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginEnableDisable func(ctx context.Context, resourceGroupName string, accountName string, chaosFault string, chaosFaultRequest armcosmos.ChaosFaultResource, options *armcosmos.ChaosFaultClientBeginEnableDisableOptions) (resp azfake.PollerResponder[armcosmos.ChaosFaultClientEnableDisableResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method ChaosFaultClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, accountName string, chaosFault string, options *armcosmos.ChaosFaultClientGetOptions) (resp azfake.Responder[armcosmos.ChaosFaultClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method ChaosFaultClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, accountName string, options *armcosmos.ChaosFaultClientListOptions) (resp azfake.PagerResponder[armcosmos.ChaosFaultClientListResponse])
}

// NewChaosFaultServerTransport creates a new instance of ChaosFaultServerTransport with the provided implementation.
// The returned ChaosFaultServerTransport instance is connected to an instance of armcosmos.ChaosFaultClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewChaosFaultServerTransport(srv *ChaosFaultServer) *ChaosFaultServerTransport {
	return &ChaosFaultServerTransport{
		srv:                srv,
		beginEnableDisable: newTracker[azfake.PollerResponder[armcosmos.ChaosFaultClientEnableDisableResponse]](),
		newListPager:       newTracker[azfake.PagerResponder[armcosmos.ChaosFaultClientListResponse]](),
	}
}

// ChaosFaultServerTransport connects instances of armcosmos.ChaosFaultClient to instances of ChaosFaultServer.
// Don't use this type directly, use NewChaosFaultServerTransport instead.
type ChaosFaultServerTransport struct {
	srv                *ChaosFaultServer
	beginEnableDisable *tracker[azfake.PollerResponder[armcosmos.ChaosFaultClientEnableDisableResponse]]
	newListPager       *tracker[azfake.PagerResponder[armcosmos.ChaosFaultClientListResponse]]
}

// Do implements the policy.Transporter interface for ChaosFaultServerTransport.
func (c *ChaosFaultServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ChaosFaultClient.BeginEnableDisable":
		resp, err = c.dispatchBeginEnableDisable(req)
	case "ChaosFaultClient.Get":
		resp, err = c.dispatchGet(req)
	case "ChaosFaultClient.NewListPager":
		resp, err = c.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *ChaosFaultServerTransport) dispatchBeginEnableDisable(req *http.Request) (*http.Response, error) {
	if c.srv.BeginEnableDisable == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginEnableDisable not implemented")}
	}
	beginEnableDisable := c.beginEnableDisable.get(req)
	if beginEnableDisable == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/databaseAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/chaosFaults/(?P<chaosFault>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcosmos.ChaosFaultResource](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
		if err != nil {
			return nil, err
		}
		chaosFaultParam, err := url.PathUnescape(matches[regex.SubexpIndex("chaosFault")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginEnableDisable(req.Context(), resourceGroupNameParam, accountNameParam, chaosFaultParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginEnableDisable = &respr
		c.beginEnableDisable.add(req, beginEnableDisable)
	}

	resp, err := server.PollerResponderNext(beginEnableDisable, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		c.beginEnableDisable.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginEnableDisable) {
		c.beginEnableDisable.remove(req)
	}

	return resp, nil
}

func (c *ChaosFaultServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if c.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/databaseAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/chaosFaults/(?P<chaosFault>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
	if err != nil {
		return nil, err
	}
	chaosFaultParam, err := url.PathUnescape(matches[regex.SubexpIndex("chaosFault")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Get(req.Context(), resourceGroupNameParam, accountNameParam, chaosFaultParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ChaosFaultResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ChaosFaultServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := c.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/databaseAccounts/(?P<accountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/chaosFaults`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		accountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("accountName")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListPager(resourceGroupNameParam, accountNameParam, nil)
		newListPager = &resp
		c.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armcosmos.ChaosFaultClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		c.newListPager.remove(req)
	}
	return resp, nil
}