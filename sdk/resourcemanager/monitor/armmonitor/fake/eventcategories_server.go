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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/monitor/armmonitor"
	"net/http"
)

// EventCategoriesServer is a fake server for instances of the armmonitor.EventCategoriesClient type.
type EventCategoriesServer struct {
	// NewListPager is the fake for method EventCategoriesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armmonitor.EventCategoriesClientListOptions) (resp azfake.PagerResponder[armmonitor.EventCategoriesClientListResponse])
}

// NewEventCategoriesServerTransport creates a new instance of EventCategoriesServerTransport with the provided implementation.
// The returned EventCategoriesServerTransport instance is connected to an instance of armmonitor.EventCategoriesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewEventCategoriesServerTransport(srv *EventCategoriesServer) *EventCategoriesServerTransport {
	return &EventCategoriesServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armmonitor.EventCategoriesClientListResponse]](),
	}
}

// EventCategoriesServerTransport connects instances of armmonitor.EventCategoriesClient to instances of EventCategoriesServer.
// Don't use this type directly, use NewEventCategoriesServerTransport instead.
type EventCategoriesServerTransport struct {
	srv          *EventCategoriesServer
	newListPager *tracker[azfake.PagerResponder[armmonitor.EventCategoriesClientListResponse]]
}

// Do implements the policy.Transporter interface for EventCategoriesServerTransport.
func (e *EventCategoriesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "EventCategoriesClient.NewListPager":
		resp, err = e.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (e *EventCategoriesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if e.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := e.newListPager.get(req)
	if newListPager == nil {
		resp := e.srv.NewListPager(nil)
		newListPager = &resp
		e.newListPager.add(req, newListPager)
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		e.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		e.newListPager.remove(req)
	}
	return resp, nil
}
