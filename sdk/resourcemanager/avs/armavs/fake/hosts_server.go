// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package fake

import (
	"context"
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/avs/armavs/v2"
	"net/http"
	"net/url"
	"regexp"
)

// HostsServer is a fake server for instances of the armavs.HostsClient type.
type HostsServer struct {
	// Get is the fake for method HostsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, privateCloudName string, clusterName string, hostID string, options *armavs.HostsClientGetOptions) (resp azfake.Responder[armavs.HostsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method HostsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, privateCloudName string, clusterName string, options *armavs.HostsClientListOptions) (resp azfake.PagerResponder[armavs.HostsClientListResponse])
}

// NewHostsServerTransport creates a new instance of HostsServerTransport with the provided implementation.
// The returned HostsServerTransport instance is connected to an instance of armavs.HostsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewHostsServerTransport(srv *HostsServer) *HostsServerTransport {
	return &HostsServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armavs.HostsClientListResponse]](),
	}
}

// HostsServerTransport connects instances of armavs.HostsClient to instances of HostsServer.
// Don't use this type directly, use NewHostsServerTransport instead.
type HostsServerTransport struct {
	srv          *HostsServer
	newListPager *tracker[azfake.PagerResponder[armavs.HostsClientListResponse]]
}

// Do implements the policy.Transporter interface for HostsServerTransport.
func (h *HostsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return h.dispatchToMethodFake(req, method)
}

func (h *HostsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if hostsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = hostsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "HostsClient.Get":
				res.resp, res.err = h.dispatchGet(req)
			case "HostsClient.NewListPager":
				res.resp, res.err = h.dispatchNewListPager(req)
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

func (h *HostsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if h.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AVS/privateClouds/(?P<privateCloudName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/hosts/(?P<hostId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if len(matches) < 6 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	privateCloudNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("privateCloudName")])
	if err != nil {
		return nil, err
	}
	clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
	if err != nil {
		return nil, err
	}
	hostIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("hostId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := h.srv.Get(req.Context(), resourceGroupNameParam, privateCloudNameParam, clusterNameParam, hostIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Host, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (h *HostsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if h.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := h.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AVS/privateClouds/(?P<privateCloudName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/hosts`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		privateCloudNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("privateCloudName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		resp := h.srv.NewListPager(resourceGroupNameParam, privateCloudNameParam, clusterNameParam, nil)
		newListPager = &resp
		h.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armavs.HostsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		h.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		h.newListPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to HostsServerTransport
var hostsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
