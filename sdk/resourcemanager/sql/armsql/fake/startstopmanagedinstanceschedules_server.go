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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql/v2"
	"net/http"
	"net/url"
	"regexp"
)

// StartStopManagedInstanceSchedulesServer is a fake server for instances of the armsql.StartStopManagedInstanceSchedulesClient type.
type StartStopManagedInstanceSchedulesServer struct {
	// CreateOrUpdate is the fake for method StartStopManagedInstanceSchedulesClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, managedInstanceName string, startStopScheduleName armsql.StartStopScheduleName, parameters armsql.StartStopManagedInstanceSchedule, options *armsql.StartStopManagedInstanceSchedulesClientCreateOrUpdateOptions) (resp azfake.Responder[armsql.StartStopManagedInstanceSchedulesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method StartStopManagedInstanceSchedulesClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, managedInstanceName string, startStopScheduleName armsql.StartStopScheduleName, options *armsql.StartStopManagedInstanceSchedulesClientDeleteOptions) (resp azfake.Responder[armsql.StartStopManagedInstanceSchedulesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method StartStopManagedInstanceSchedulesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, managedInstanceName string, startStopScheduleName armsql.StartStopScheduleName, options *armsql.StartStopManagedInstanceSchedulesClientGetOptions) (resp azfake.Responder[armsql.StartStopManagedInstanceSchedulesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByInstancePager is the fake for method StartStopManagedInstanceSchedulesClient.NewListByInstancePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByInstancePager func(resourceGroupName string, managedInstanceName string, options *armsql.StartStopManagedInstanceSchedulesClientListByInstanceOptions) (resp azfake.PagerResponder[armsql.StartStopManagedInstanceSchedulesClientListByInstanceResponse])
}

// NewStartStopManagedInstanceSchedulesServerTransport creates a new instance of StartStopManagedInstanceSchedulesServerTransport with the provided implementation.
// The returned StartStopManagedInstanceSchedulesServerTransport instance is connected to an instance of armsql.StartStopManagedInstanceSchedulesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewStartStopManagedInstanceSchedulesServerTransport(srv *StartStopManagedInstanceSchedulesServer) *StartStopManagedInstanceSchedulesServerTransport {
	return &StartStopManagedInstanceSchedulesServerTransport{
		srv:                    srv,
		newListByInstancePager: newTracker[azfake.PagerResponder[armsql.StartStopManagedInstanceSchedulesClientListByInstanceResponse]](),
	}
}

// StartStopManagedInstanceSchedulesServerTransport connects instances of armsql.StartStopManagedInstanceSchedulesClient to instances of StartStopManagedInstanceSchedulesServer.
// Don't use this type directly, use NewStartStopManagedInstanceSchedulesServerTransport instead.
type StartStopManagedInstanceSchedulesServerTransport struct {
	srv                    *StartStopManagedInstanceSchedulesServer
	newListByInstancePager *tracker[azfake.PagerResponder[armsql.StartStopManagedInstanceSchedulesClientListByInstanceResponse]]
}

// Do implements the policy.Transporter interface for StartStopManagedInstanceSchedulesServerTransport.
func (s *StartStopManagedInstanceSchedulesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "StartStopManagedInstanceSchedulesClient.CreateOrUpdate":
		resp, err = s.dispatchCreateOrUpdate(req)
	case "StartStopManagedInstanceSchedulesClient.Delete":
		resp, err = s.dispatchDelete(req)
	case "StartStopManagedInstanceSchedulesClient.Get":
		resp, err = s.dispatchGet(req)
	case "StartStopManagedInstanceSchedulesClient.NewListByInstancePager":
		resp, err = s.dispatchNewListByInstancePager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *StartStopManagedInstanceSchedulesServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/managedInstances/(?P<managedInstanceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/startStopSchedules/(?P<startStopScheduleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armsql.StartStopManagedInstanceSchedule](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	managedInstanceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("managedInstanceName")])
	if err != nil {
		return nil, err
	}
	startStopScheduleNameParam, err := parseWithCast(matches[regex.SubexpIndex("startStopScheduleName")], func(v string) (armsql.StartStopScheduleName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armsql.StartStopScheduleName(p), nil
	})
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, managedInstanceNameParam, startStopScheduleNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).StartStopManagedInstanceSchedule, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *StartStopManagedInstanceSchedulesServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if s.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/managedInstances/(?P<managedInstanceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/startStopSchedules/(?P<startStopScheduleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	managedInstanceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("managedInstanceName")])
	if err != nil {
		return nil, err
	}
	startStopScheduleNameParam, err := parseWithCast(matches[regex.SubexpIndex("startStopScheduleName")], func(v string) (armsql.StartStopScheduleName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armsql.StartStopScheduleName(p), nil
	})
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Delete(req.Context(), resourceGroupNameParam, managedInstanceNameParam, startStopScheduleNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusNoContent}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNoContent", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *StartStopManagedInstanceSchedulesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/managedInstances/(?P<managedInstanceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/startStopSchedules/(?P<startStopScheduleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	managedInstanceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("managedInstanceName")])
	if err != nil {
		return nil, err
	}
	startStopScheduleNameParam, err := parseWithCast(matches[regex.SubexpIndex("startStopScheduleName")], func(v string) (armsql.StartStopScheduleName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armsql.StartStopScheduleName(p), nil
	})
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), resourceGroupNameParam, managedInstanceNameParam, startStopScheduleNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).StartStopManagedInstanceSchedule, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *StartStopManagedInstanceSchedulesServerTransport) dispatchNewListByInstancePager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByInstancePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByInstancePager not implemented")}
	}
	newListByInstancePager := s.newListByInstancePager.get(req)
	if newListByInstancePager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/managedInstances/(?P<managedInstanceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/startStopSchedules`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		managedInstanceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("managedInstanceName")])
		if err != nil {
			return nil, err
		}
		resp := s.srv.NewListByInstancePager(resourceGroupNameParam, managedInstanceNameParam, nil)
		newListByInstancePager = &resp
		s.newListByInstancePager.add(req, newListByInstancePager)
		server.PagerResponderInjectNextLinks(newListByInstancePager, req, func(page *armsql.StartStopManagedInstanceSchedulesClientListByInstanceResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByInstancePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListByInstancePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByInstancePager) {
		s.newListByInstancePager.remove(req)
	}
	return resp, nil
}