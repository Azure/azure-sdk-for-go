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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/guestconfiguration/armguestconfiguration"
	"net/http"
	"net/url"
	"regexp"
)

// AssignmentReportsServer is a fake server for instances of the armguestconfiguration.AssignmentReportsClient type.
type AssignmentReportsServer struct {
	// Get is the fake for method AssignmentReportsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, guestConfigurationAssignmentName string, reportID string, vmName string, options *armguestconfiguration.AssignmentReportsClientGetOptions) (resp azfake.Responder[armguestconfiguration.AssignmentReportsClientGetResponse], errResp azfake.ErrorResponder)

	// List is the fake for method AssignmentReportsClient.List
	// HTTP status codes to indicate success: http.StatusOK
	List func(ctx context.Context, resourceGroupName string, guestConfigurationAssignmentName string, vmName string, options *armguestconfiguration.AssignmentReportsClientListOptions) (resp azfake.Responder[armguestconfiguration.AssignmentReportsClientListResponse], errResp azfake.ErrorResponder)
}

// NewAssignmentReportsServerTransport creates a new instance of AssignmentReportsServerTransport with the provided implementation.
// The returned AssignmentReportsServerTransport instance is connected to an instance of armguestconfiguration.AssignmentReportsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAssignmentReportsServerTransport(srv *AssignmentReportsServer) *AssignmentReportsServerTransport {
	return &AssignmentReportsServerTransport{srv: srv}
}

// AssignmentReportsServerTransport connects instances of armguestconfiguration.AssignmentReportsClient to instances of AssignmentReportsServer.
// Don't use this type directly, use NewAssignmentReportsServerTransport instead.
type AssignmentReportsServerTransport struct {
	srv *AssignmentReportsServer
}

// Do implements the policy.Transporter interface for AssignmentReportsServerTransport.
func (a *AssignmentReportsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "AssignmentReportsClient.Get":
		resp, err = a.dispatchGet(req)
	case "AssignmentReportsClient.List":
		resp, err = a.dispatchList(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AssignmentReportsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if a.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/virtualMachines/(?P<vmName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.GuestConfiguration/guestConfigurationAssignments/(?P<guestConfigurationAssignmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/reports/(?P<reportId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	guestConfigurationAssignmentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("guestConfigurationAssignmentName")])
	if err != nil {
		return nil, err
	}
	reportIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("reportId")])
	if err != nil {
		return nil, err
	}
	vmNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("vmName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Get(req.Context(), resourceGroupNameParam, guestConfigurationAssignmentNameParam, reportIDParam, vmNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AssignmentReport, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AssignmentReportsServerTransport) dispatchList(req *http.Request) (*http.Response, error) {
	if a.srv.List == nil {
		return nil, &nonRetriableError{errors.New("fake for method List not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/virtualMachines/(?P<vmName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.GuestConfiguration/guestConfigurationAssignments/(?P<guestConfigurationAssignmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/reports`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	guestConfigurationAssignmentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("guestConfigurationAssignmentName")])
	if err != nil {
		return nil, err
	}
	vmNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("vmName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.List(req.Context(), resourceGroupNameParam, guestConfigurationAssignmentNameParam, vmNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AssignmentReportList, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
