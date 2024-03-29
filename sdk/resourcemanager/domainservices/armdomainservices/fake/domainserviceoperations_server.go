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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/domainservices/armdomainservices"
	"net/http"
)

// DomainServiceOperationsServer is a fake server for instances of the armdomainservices.DomainServiceOperationsClient type.
type DomainServiceOperationsServer struct {
	// NewListPager is the fake for method DomainServiceOperationsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armdomainservices.DomainServiceOperationsClientListOptions) (resp azfake.PagerResponder[armdomainservices.DomainServiceOperationsClientListResponse])
}

// NewDomainServiceOperationsServerTransport creates a new instance of DomainServiceOperationsServerTransport with the provided implementation.
// The returned DomainServiceOperationsServerTransport instance is connected to an instance of armdomainservices.DomainServiceOperationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDomainServiceOperationsServerTransport(srv *DomainServiceOperationsServer) *DomainServiceOperationsServerTransport {
	return &DomainServiceOperationsServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armdomainservices.DomainServiceOperationsClientListResponse]](),
	}
}

// DomainServiceOperationsServerTransport connects instances of armdomainservices.DomainServiceOperationsClient to instances of DomainServiceOperationsServer.
// Don't use this type directly, use NewDomainServiceOperationsServerTransport instead.
type DomainServiceOperationsServerTransport struct {
	srv          *DomainServiceOperationsServer
	newListPager *tracker[azfake.PagerResponder[armdomainservices.DomainServiceOperationsClientListResponse]]
}

// Do implements the policy.Transporter interface for DomainServiceOperationsServerTransport.
func (d *DomainServiceOperationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "DomainServiceOperationsClient.NewListPager":
		resp, err = d.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DomainServiceOperationsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := d.newListPager.get(req)
	if newListPager == nil {
		resp := d.srv.NewListPager(nil)
		newListPager = &resp
		d.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armdomainservices.DomainServiceOperationsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		d.newListPager.remove(req)
	}
	return resp, nil
}
