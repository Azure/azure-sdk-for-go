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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/changeanalysis/armchangeanalysis"
	"net/http"
	"net/url"
	"regexp"
	"time"
)

// ResourceChangesServer is a fake server for instances of the armchangeanalysis.ResourceChangesClient type.
type ResourceChangesServer struct {
	// NewListPager is the fake for method ResourceChangesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceID string, startTime time.Time, endTime time.Time, options *armchangeanalysis.ResourceChangesClientListOptions) (resp azfake.PagerResponder[armchangeanalysis.ResourceChangesClientListResponse])
}

// NewResourceChangesServerTransport creates a new instance of ResourceChangesServerTransport with the provided implementation.
// The returned ResourceChangesServerTransport instance is connected to an instance of armchangeanalysis.ResourceChangesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewResourceChangesServerTransport(srv *ResourceChangesServer) *ResourceChangesServerTransport {
	return &ResourceChangesServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armchangeanalysis.ResourceChangesClientListResponse]](),
	}
}

// ResourceChangesServerTransport connects instances of armchangeanalysis.ResourceChangesClient to instances of ResourceChangesServer.
// Don't use this type directly, use NewResourceChangesServerTransport instead.
type ResourceChangesServerTransport struct {
	srv          *ResourceChangesServer
	newListPager *tracker[azfake.PagerResponder[armchangeanalysis.ResourceChangesClientListResponse]]
}

// Do implements the policy.Transporter interface for ResourceChangesServerTransport.
func (r *ResourceChangesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ResourceChangesClient.NewListPager":
		resp, err = r.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *ResourceChangesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := r.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/(?P<resourceId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ChangeAnalysis/resourceChanges`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceId")])
		if err != nil {
			return nil, err
		}
		startTimeUnescaped, err := url.QueryUnescape(qp.Get("$startTime"))
		if err != nil {
			return nil, err
		}
		startTimeParam, err := time.Parse(time.RFC3339Nano, startTimeUnescaped)
		if err != nil {
			return nil, err
		}
		endTimeUnescaped, err := url.QueryUnescape(qp.Get("$endTime"))
		if err != nil {
			return nil, err
		}
		endTimeParam, err := time.Parse(time.RFC3339Nano, endTimeUnescaped)
		if err != nil {
			return nil, err
		}
		skipTokenUnescaped, err := url.QueryUnescape(qp.Get("$skipToken"))
		if err != nil {
			return nil, err
		}
		skipTokenParam := getOptional(skipTokenUnescaped)
		var options *armchangeanalysis.ResourceChangesClientListOptions
		if skipTokenParam != nil {
			options = &armchangeanalysis.ResourceChangesClientListOptions{
				SkipToken: skipTokenParam,
			}
		}
		resp := r.srv.NewListPager(resourceIDParam, startTimeParam, endTimeParam, options)
		newListPager = &resp
		r.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armchangeanalysis.ResourceChangesClientListResponse, createLink func() string) {
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
