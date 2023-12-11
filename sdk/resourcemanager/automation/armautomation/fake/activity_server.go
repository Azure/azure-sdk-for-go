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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automation/armautomation"
	"net/http"
	"net/url"
	"regexp"
)

// ActivityServer is a fake server for instances of the armautomation.ActivityClient type.
type ActivityServer struct {
	// Get is the fake for method ActivityClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, automationAccountName string, moduleName string, activityName string, options *armautomation.ActivityClientGetOptions) (resp azfake.Responder[armautomation.ActivityClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByModulePager is the fake for method ActivityClient.NewListByModulePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByModulePager func(resourceGroupName string, automationAccountName string, moduleName string, options *armautomation.ActivityClientListByModuleOptions) (resp azfake.PagerResponder[armautomation.ActivityClientListByModuleResponse])
}

// NewActivityServerTransport creates a new instance of ActivityServerTransport with the provided implementation.
// The returned ActivityServerTransport instance is connected to an instance of armautomation.ActivityClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewActivityServerTransport(srv *ActivityServer) *ActivityServerTransport {
	return &ActivityServerTransport{
		srv:                  srv,
		newListByModulePager: newTracker[azfake.PagerResponder[armautomation.ActivityClientListByModuleResponse]](),
	}
}

// ActivityServerTransport connects instances of armautomation.ActivityClient to instances of ActivityServer.
// Don't use this type directly, use NewActivityServerTransport instead.
type ActivityServerTransport struct {
	srv                  *ActivityServer
	newListByModulePager *tracker[azfake.PagerResponder[armautomation.ActivityClientListByModuleResponse]]
}

// Do implements the policy.Transporter interface for ActivityServerTransport.
func (a *ActivityServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ActivityClient.Get":
		resp, err = a.dispatchGet(req)
	case "ActivityClient.NewListByModulePager":
		resp, err = a.dispatchNewListByModulePager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *ActivityServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if a.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automation/automationAccounts/(?P<automationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/modules/(?P<moduleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/activities/(?P<activityName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	automationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("automationAccountName")])
	if err != nil {
		return nil, err
	}
	moduleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("moduleName")])
	if err != nil {
		return nil, err
	}
	activityNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("activityName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Get(req.Context(), resourceGroupNameParam, automationAccountNameParam, moduleNameParam, activityNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Activity, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *ActivityServerTransport) dispatchNewListByModulePager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListByModulePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByModulePager not implemented")}
	}
	newListByModulePager := a.newListByModulePager.get(req)
	if newListByModulePager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automation/automationAccounts/(?P<automationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/modules/(?P<moduleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/activities`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		automationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("automationAccountName")])
		if err != nil {
			return nil, err
		}
		moduleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("moduleName")])
		if err != nil {
			return nil, err
		}
		resp := a.srv.NewListByModulePager(resourceGroupNameParam, automationAccountNameParam, moduleNameParam, nil)
		newListByModulePager = &resp
		a.newListByModulePager.add(req, newListByModulePager)
		server.PagerResponderInjectNextLinks(newListByModulePager, req, func(page *armautomation.ActivityClientListByModuleResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByModulePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListByModulePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByModulePager) {
		a.newListByModulePager.remove(req)
	}
	return resp, nil
}