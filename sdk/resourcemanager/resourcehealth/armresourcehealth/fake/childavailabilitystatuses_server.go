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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resourcehealth/armresourcehealth"
	"net/http"
	"net/url"
	"regexp"
)

// ChildAvailabilityStatusesServer is a fake server for instances of the armresourcehealth.ChildAvailabilityStatusesClient type.
type ChildAvailabilityStatusesServer struct {
	// GetByResource is the fake for method ChildAvailabilityStatusesClient.GetByResource
	// HTTP status codes to indicate success: http.StatusOK
	GetByResource func(ctx context.Context, resourceURI string, options *armresourcehealth.ChildAvailabilityStatusesClientGetByResourceOptions) (resp azfake.Responder[armresourcehealth.ChildAvailabilityStatusesClientGetByResourceResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method ChildAvailabilityStatusesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceURI string, options *armresourcehealth.ChildAvailabilityStatusesClientListOptions) (resp azfake.PagerResponder[armresourcehealth.ChildAvailabilityStatusesClientListResponse])
}

// NewChildAvailabilityStatusesServerTransport creates a new instance of ChildAvailabilityStatusesServerTransport with the provided implementation.
// The returned ChildAvailabilityStatusesServerTransport instance is connected to an instance of armresourcehealth.ChildAvailabilityStatusesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewChildAvailabilityStatusesServerTransport(srv *ChildAvailabilityStatusesServer) *ChildAvailabilityStatusesServerTransport {
	return &ChildAvailabilityStatusesServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armresourcehealth.ChildAvailabilityStatusesClientListResponse]](),
	}
}

// ChildAvailabilityStatusesServerTransport connects instances of armresourcehealth.ChildAvailabilityStatusesClient to instances of ChildAvailabilityStatusesServer.
// Don't use this type directly, use NewChildAvailabilityStatusesServerTransport instead.
type ChildAvailabilityStatusesServerTransport struct {
	srv          *ChildAvailabilityStatusesServer
	newListPager *tracker[azfake.PagerResponder[armresourcehealth.ChildAvailabilityStatusesClientListResponse]]
}

// Do implements the policy.Transporter interface for ChildAvailabilityStatusesServerTransport.
func (c *ChildAvailabilityStatusesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ChildAvailabilityStatusesClient.GetByResource":
		resp, err = c.dispatchGetByResource(req)
	case "ChildAvailabilityStatusesClient.NewListPager":
		resp, err = c.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *ChildAvailabilityStatusesServerTransport) dispatchGetByResource(req *http.Request) (*http.Response, error) {
	if c.srv.GetByResource == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetByResource not implemented")}
	}
	const regexStr = `/(?P<resourceUri>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ResourceHealth/childAvailabilityStatuses/current`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	resourceURIParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceUri")])
	if err != nil {
		return nil, err
	}
	filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
	if err != nil {
		return nil, err
	}
	filterParam := getOptional(filterUnescaped)
	expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
	if err != nil {
		return nil, err
	}
	expandParam := getOptional(expandUnescaped)
	var options *armresourcehealth.ChildAvailabilityStatusesClientGetByResourceOptions
	if filterParam != nil || expandParam != nil {
		options = &armresourcehealth.ChildAvailabilityStatusesClientGetByResourceOptions{
			Filter: filterParam,
			Expand: expandParam,
		}
	}
	respr, errRespr := c.srv.GetByResource(req.Context(), resourceURIParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AvailabilityStatus, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ChildAvailabilityStatusesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := c.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/(?P<resourceUri>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ResourceHealth/childAvailabilityStatuses`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceURIParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceUri")])
		if err != nil {
			return nil, err
		}
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
		if err != nil {
			return nil, err
		}
		expandParam := getOptional(expandUnescaped)
		var options *armresourcehealth.ChildAvailabilityStatusesClientListOptions
		if filterParam != nil || expandParam != nil {
			options = &armresourcehealth.ChildAvailabilityStatusesClientListOptions{
				Filter: filterParam,
				Expand: expandParam,
			}
		}
		resp := c.srv.NewListPager(resourceURIParam, options)
		newListPager = &resp
		c.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armresourcehealth.ChildAvailabilityStatusesClientListResponse, createLink func() string) {
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
