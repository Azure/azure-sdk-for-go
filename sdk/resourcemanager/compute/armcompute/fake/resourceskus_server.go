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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v6"
	"net/http"
	"net/url"
	"regexp"
)

// ResourceSKUsServer is a fake server for instances of the armcompute.ResourceSKUsClient type.
type ResourceSKUsServer struct {
	// NewListPager is the fake for method ResourceSKUsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armcompute.ResourceSKUsClientListOptions) (resp azfake.PagerResponder[armcompute.ResourceSKUsClientListResponse])
}

// NewResourceSKUsServerTransport creates a new instance of ResourceSKUsServerTransport with the provided implementation.
// The returned ResourceSKUsServerTransport instance is connected to an instance of armcompute.ResourceSKUsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewResourceSKUsServerTransport(srv *ResourceSKUsServer) *ResourceSKUsServerTransport {
	return &ResourceSKUsServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armcompute.ResourceSKUsClientListResponse]](),
	}
}

// ResourceSKUsServerTransport connects instances of armcompute.ResourceSKUsClient to instances of ResourceSKUsServer.
// Don't use this type directly, use NewResourceSKUsServerTransport instead.
type ResourceSKUsServerTransport struct {
	srv          *ResourceSKUsServer
	newListPager *tracker[azfake.PagerResponder[armcompute.ResourceSKUsClientListResponse]]
}

// Do implements the policy.Transporter interface for ResourceSKUsServerTransport.
func (r *ResourceSKUsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return r.dispatchToMethodFake(req, method)
}

func (r *ResourceSKUsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if resourceSkUsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = resourceSkUsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "ResourceSKUsClient.NewListPager":
				res.resp, res.err = r.dispatchNewListPager(req)
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

func (r *ResourceSKUsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := r.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/skus`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		includeExtendedLocationsUnescaped, err := url.QueryUnescape(qp.Get("includeExtendedLocations"))
		if err != nil {
			return nil, err
		}
		includeExtendedLocationsParam := getOptional(includeExtendedLocationsUnescaped)
		var options *armcompute.ResourceSKUsClientListOptions
		if filterParam != nil || includeExtendedLocationsParam != nil {
			options = &armcompute.ResourceSKUsClientListOptions{
				Filter:                   filterParam,
				IncludeExtendedLocations: includeExtendedLocationsParam,
			}
		}
		resp := r.srv.NewListPager(options)
		newListPager = &resp
		r.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armcompute.ResourceSKUsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		r.newListPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to ResourceSKUsServerTransport
var resourceSkUsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
