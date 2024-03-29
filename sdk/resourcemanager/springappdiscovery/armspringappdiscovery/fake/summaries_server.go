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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/springappdiscovery/armspringappdiscovery"
	"net/http"
	"net/url"
	"regexp"
)

// SummariesServer is a fake server for instances of the armspringappdiscovery.SummariesClient type.
type SummariesServer struct {
	// Get is the fake for method SummariesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, siteName string, summaryName string, options *armspringappdiscovery.SummariesClientGetOptions) (resp azfake.Responder[armspringappdiscovery.SummariesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListBySitePager is the fake for method SummariesClient.NewListBySitePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySitePager func(resourceGroupName string, siteName string, options *armspringappdiscovery.SummariesClientListBySiteOptions) (resp azfake.PagerResponder[armspringappdiscovery.SummariesClientListBySiteResponse])
}

// NewSummariesServerTransport creates a new instance of SummariesServerTransport with the provided implementation.
// The returned SummariesServerTransport instance is connected to an instance of armspringappdiscovery.SummariesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSummariesServerTransport(srv *SummariesServer) *SummariesServerTransport {
	return &SummariesServerTransport{
		srv:                srv,
		newListBySitePager: newTracker[azfake.PagerResponder[armspringappdiscovery.SummariesClientListBySiteResponse]](),
	}
}

// SummariesServerTransport connects instances of armspringappdiscovery.SummariesClient to instances of SummariesServer.
// Don't use this type directly, use NewSummariesServerTransport instead.
type SummariesServerTransport struct {
	srv                *SummariesServer
	newListBySitePager *tracker[azfake.PagerResponder[armspringappdiscovery.SummariesClientListBySiteResponse]]
}

// Do implements the policy.Transporter interface for SummariesServerTransport.
func (s *SummariesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "SummariesClient.Get":
		resp, err = s.dispatchGet(req)
	case "SummariesClient.NewListBySitePager":
		resp, err = s.dispatchNewListBySitePager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *SummariesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.OffAzureSpringBoot/springbootsites/(?P<siteName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/summaries/(?P<summaryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	siteNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("siteName")])
	if err != nil {
		return nil, err
	}
	summaryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("summaryName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), resourceGroupNameParam, siteNameParam, summaryNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Summary, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SummariesServerTransport) dispatchNewListBySitePager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListBySitePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySitePager not implemented")}
	}
	newListBySitePager := s.newListBySitePager.get(req)
	if newListBySitePager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.OffAzureSpringBoot/springbootsites/(?P<siteName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/summaries`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		siteNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("siteName")])
		if err != nil {
			return nil, err
		}
		resp := s.srv.NewListBySitePager(resourceGroupNameParam, siteNameParam, nil)
		newListBySitePager = &resp
		s.newListBySitePager.add(req, newListBySitePager)
		server.PagerResponderInjectNextLinks(newListBySitePager, req, func(page *armspringappdiscovery.SummariesClientListBySiteResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySitePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListBySitePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySitePager) {
		s.newListBySitePager.remove(req)
	}
	return resp, nil
}
