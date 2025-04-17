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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"net/http"
	"net/url"
	"regexp"
)

// CloudServiceOperatingSystemsServer is a fake server for instances of the armcompute.CloudServiceOperatingSystemsClient type.
type CloudServiceOperatingSystemsServer struct {
	// GetOSFamily is the fake for method CloudServiceOperatingSystemsClient.GetOSFamily
	// HTTP status codes to indicate success: http.StatusOK
	GetOSFamily func(ctx context.Context, location string, osFamilyName string, options *armcompute.CloudServiceOperatingSystemsClientGetOSFamilyOptions) (resp azfake.Responder[armcompute.CloudServiceOperatingSystemsClientGetOSFamilyResponse], errResp azfake.ErrorResponder)

	// GetOSVersion is the fake for method CloudServiceOperatingSystemsClient.GetOSVersion
	// HTTP status codes to indicate success: http.StatusOK
	GetOSVersion func(ctx context.Context, location string, osVersionName string, options *armcompute.CloudServiceOperatingSystemsClientGetOSVersionOptions) (resp azfake.Responder[armcompute.CloudServiceOperatingSystemsClientGetOSVersionResponse], errResp azfake.ErrorResponder)

	// NewListOSFamiliesPager is the fake for method CloudServiceOperatingSystemsClient.NewListOSFamiliesPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListOSFamiliesPager func(location string, options *armcompute.CloudServiceOperatingSystemsClientListOSFamiliesOptions) (resp azfake.PagerResponder[armcompute.CloudServiceOperatingSystemsClientListOSFamiliesResponse])

	// NewListOSVersionsPager is the fake for method CloudServiceOperatingSystemsClient.NewListOSVersionsPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListOSVersionsPager func(location string, options *armcompute.CloudServiceOperatingSystemsClientListOSVersionsOptions) (resp azfake.PagerResponder[armcompute.CloudServiceOperatingSystemsClientListOSVersionsResponse])
}

// NewCloudServiceOperatingSystemsServerTransport creates a new instance of CloudServiceOperatingSystemsServerTransport with the provided implementation.
// The returned CloudServiceOperatingSystemsServerTransport instance is connected to an instance of armcompute.CloudServiceOperatingSystemsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewCloudServiceOperatingSystemsServerTransport(srv *CloudServiceOperatingSystemsServer) *CloudServiceOperatingSystemsServerTransport {
	return &CloudServiceOperatingSystemsServerTransport{
		srv:                    srv,
		newListOSFamiliesPager: newTracker[azfake.PagerResponder[armcompute.CloudServiceOperatingSystemsClientListOSFamiliesResponse]](),
		newListOSVersionsPager: newTracker[azfake.PagerResponder[armcompute.CloudServiceOperatingSystemsClientListOSVersionsResponse]](),
	}
}

// CloudServiceOperatingSystemsServerTransport connects instances of armcompute.CloudServiceOperatingSystemsClient to instances of CloudServiceOperatingSystemsServer.
// Don't use this type directly, use NewCloudServiceOperatingSystemsServerTransport instead.
type CloudServiceOperatingSystemsServerTransport struct {
	srv                    *CloudServiceOperatingSystemsServer
	newListOSFamiliesPager *tracker[azfake.PagerResponder[armcompute.CloudServiceOperatingSystemsClientListOSFamiliesResponse]]
	newListOSVersionsPager *tracker[azfake.PagerResponder[armcompute.CloudServiceOperatingSystemsClientListOSVersionsResponse]]
}

// Do implements the policy.Transporter interface for CloudServiceOperatingSystemsServerTransport.
func (c *CloudServiceOperatingSystemsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return c.dispatchToMethodFake(req, method)
}

func (c *CloudServiceOperatingSystemsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if cloudServiceOperatingSystemsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = cloudServiceOperatingSystemsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "CloudServiceOperatingSystemsClient.GetOSFamily":
				res.resp, res.err = c.dispatchGetOSFamily(req)
			case "CloudServiceOperatingSystemsClient.GetOSVersion":
				res.resp, res.err = c.dispatchGetOSVersion(req)
			case "CloudServiceOperatingSystemsClient.NewListOSFamiliesPager":
				res.resp, res.err = c.dispatchNewListOSFamiliesPager(req)
			case "CloudServiceOperatingSystemsClient.NewListOSVersionsPager":
				res.resp, res.err = c.dispatchNewListOSVersionsPager(req)
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

func (c *CloudServiceOperatingSystemsServerTransport) dispatchGetOSFamily(req *http.Request) (*http.Response, error) {
	if c.srv.GetOSFamily == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetOSFamily not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/cloudServiceOsFamilies/(?P<osFamilyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
	if err != nil {
		return nil, err
	}
	osFamilyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("osFamilyName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.GetOSFamily(req.Context(), locationParam, osFamilyNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).OSFamily, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CloudServiceOperatingSystemsServerTransport) dispatchGetOSVersion(req *http.Request) (*http.Response, error) {
	if c.srv.GetOSVersion == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetOSVersion not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/cloudServiceOsVersions/(?P<osVersionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
	if err != nil {
		return nil, err
	}
	osVersionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("osVersionName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.GetOSVersion(req.Context(), locationParam, osVersionNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).OSVersion, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CloudServiceOperatingSystemsServerTransport) dispatchNewListOSFamiliesPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListOSFamiliesPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListOSFamiliesPager not implemented")}
	}
	newListOSFamiliesPager := c.newListOSFamiliesPager.get(req)
	if newListOSFamiliesPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/cloudServiceOsFamilies`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListOSFamiliesPager(locationParam, nil)
		newListOSFamiliesPager = &resp
		c.newListOSFamiliesPager.add(req, newListOSFamiliesPager)
		server.PagerResponderInjectNextLinks(newListOSFamiliesPager, req, func(page *armcompute.CloudServiceOperatingSystemsClientListOSFamiliesResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListOSFamiliesPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListOSFamiliesPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListOSFamiliesPager) {
		c.newListOSFamiliesPager.remove(req)
	}
	return resp, nil
}

func (c *CloudServiceOperatingSystemsServerTransport) dispatchNewListOSVersionsPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListOSVersionsPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListOSVersionsPager not implemented")}
	}
	newListOSVersionsPager := c.newListOSVersionsPager.get(req)
	if newListOSVersionsPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/cloudServiceOsVersions`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListOSVersionsPager(locationParam, nil)
		newListOSVersionsPager = &resp
		c.newListOSVersionsPager.add(req, newListOSVersionsPager)
		server.PagerResponderInjectNextLinks(newListOSVersionsPager, req, func(page *armcompute.CloudServiceOperatingSystemsClientListOSVersionsResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListOSVersionsPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListOSVersionsPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListOSVersionsPager) {
		c.newListOSVersionsPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to CloudServiceOperatingSystemsServerTransport
var cloudServiceOperatingSystemsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
