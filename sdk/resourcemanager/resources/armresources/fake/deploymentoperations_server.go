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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armresources"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// DeploymentOperationsServer is a fake server for instances of the armresources.DeploymentOperationsClient type.
type DeploymentOperationsServer struct {
	// Get is the fake for method DeploymentOperationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, deploymentName string, operationID string, options *armresources.DeploymentOperationsClientGetOptions) (resp azfake.Responder[armresources.DeploymentOperationsClientGetResponse], errResp azfake.ErrorResponder)

	// GetAtManagementGroupScope is the fake for method DeploymentOperationsClient.GetAtManagementGroupScope
	// HTTP status codes to indicate success: http.StatusOK
	GetAtManagementGroupScope func(ctx context.Context, groupID string, deploymentName string, operationID string, options *armresources.DeploymentOperationsClientGetAtManagementGroupScopeOptions) (resp azfake.Responder[armresources.DeploymentOperationsClientGetAtManagementGroupScopeResponse], errResp azfake.ErrorResponder)

	// GetAtScope is the fake for method DeploymentOperationsClient.GetAtScope
	// HTTP status codes to indicate success: http.StatusOK
	GetAtScope func(ctx context.Context, scope string, deploymentName string, operationID string, options *armresources.DeploymentOperationsClientGetAtScopeOptions) (resp azfake.Responder[armresources.DeploymentOperationsClientGetAtScopeResponse], errResp azfake.ErrorResponder)

	// GetAtSubscriptionScope is the fake for method DeploymentOperationsClient.GetAtSubscriptionScope
	// HTTP status codes to indicate success: http.StatusOK
	GetAtSubscriptionScope func(ctx context.Context, deploymentName string, operationID string, options *armresources.DeploymentOperationsClientGetAtSubscriptionScopeOptions) (resp azfake.Responder[armresources.DeploymentOperationsClientGetAtSubscriptionScopeResponse], errResp azfake.ErrorResponder)

	// GetAtTenantScope is the fake for method DeploymentOperationsClient.GetAtTenantScope
	// HTTP status codes to indicate success: http.StatusOK
	GetAtTenantScope func(ctx context.Context, deploymentName string, operationID string, options *armresources.DeploymentOperationsClientGetAtTenantScopeOptions) (resp azfake.Responder[armresources.DeploymentOperationsClientGetAtTenantScopeResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method DeploymentOperationsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, deploymentName string, options *armresources.DeploymentOperationsClientListOptions) (resp azfake.PagerResponder[armresources.DeploymentOperationsClientListResponse])

	// NewListAtManagementGroupScopePager is the fake for method DeploymentOperationsClient.NewListAtManagementGroupScopePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListAtManagementGroupScopePager func(groupID string, deploymentName string, options *armresources.DeploymentOperationsClientListAtManagementGroupScopeOptions) (resp azfake.PagerResponder[armresources.DeploymentOperationsClientListAtManagementGroupScopeResponse])

	// NewListAtScopePager is the fake for method DeploymentOperationsClient.NewListAtScopePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListAtScopePager func(scope string, deploymentName string, options *armresources.DeploymentOperationsClientListAtScopeOptions) (resp azfake.PagerResponder[armresources.DeploymentOperationsClientListAtScopeResponse])

	// NewListAtSubscriptionScopePager is the fake for method DeploymentOperationsClient.NewListAtSubscriptionScopePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListAtSubscriptionScopePager func(deploymentName string, options *armresources.DeploymentOperationsClientListAtSubscriptionScopeOptions) (resp azfake.PagerResponder[armresources.DeploymentOperationsClientListAtSubscriptionScopeResponse])

	// NewListAtTenantScopePager is the fake for method DeploymentOperationsClient.NewListAtTenantScopePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListAtTenantScopePager func(deploymentName string, options *armresources.DeploymentOperationsClientListAtTenantScopeOptions) (resp azfake.PagerResponder[armresources.DeploymentOperationsClientListAtTenantScopeResponse])
}

// NewDeploymentOperationsServerTransport creates a new instance of DeploymentOperationsServerTransport with the provided implementation.
// The returned DeploymentOperationsServerTransport instance is connected to an instance of armresources.DeploymentOperationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDeploymentOperationsServerTransport(srv *DeploymentOperationsServer) *DeploymentOperationsServerTransport {
	return &DeploymentOperationsServerTransport{
		srv:                                srv,
		newListPager:                       newTracker[azfake.PagerResponder[armresources.DeploymentOperationsClientListResponse]](),
		newListAtManagementGroupScopePager: newTracker[azfake.PagerResponder[armresources.DeploymentOperationsClientListAtManagementGroupScopeResponse]](),
		newListAtScopePager:                newTracker[azfake.PagerResponder[armresources.DeploymentOperationsClientListAtScopeResponse]](),
		newListAtSubscriptionScopePager:    newTracker[azfake.PagerResponder[armresources.DeploymentOperationsClientListAtSubscriptionScopeResponse]](),
		newListAtTenantScopePager:          newTracker[azfake.PagerResponder[armresources.DeploymentOperationsClientListAtTenantScopeResponse]](),
	}
}

// DeploymentOperationsServerTransport connects instances of armresources.DeploymentOperationsClient to instances of DeploymentOperationsServer.
// Don't use this type directly, use NewDeploymentOperationsServerTransport instead.
type DeploymentOperationsServerTransport struct {
	srv                                *DeploymentOperationsServer
	newListPager                       *tracker[azfake.PagerResponder[armresources.DeploymentOperationsClientListResponse]]
	newListAtManagementGroupScopePager *tracker[azfake.PagerResponder[armresources.DeploymentOperationsClientListAtManagementGroupScopeResponse]]
	newListAtScopePager                *tracker[azfake.PagerResponder[armresources.DeploymentOperationsClientListAtScopeResponse]]
	newListAtSubscriptionScopePager    *tracker[azfake.PagerResponder[armresources.DeploymentOperationsClientListAtSubscriptionScopeResponse]]
	newListAtTenantScopePager          *tracker[azfake.PagerResponder[armresources.DeploymentOperationsClientListAtTenantScopeResponse]]
}

// Do implements the policy.Transporter interface for DeploymentOperationsServerTransport.
func (d *DeploymentOperationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "DeploymentOperationsClient.Get":
		resp, err = d.dispatchGet(req)
	case "DeploymentOperationsClient.GetAtManagementGroupScope":
		resp, err = d.dispatchGetAtManagementGroupScope(req)
	case "DeploymentOperationsClient.GetAtScope":
		resp, err = d.dispatchGetAtScope(req)
	case "DeploymentOperationsClient.GetAtSubscriptionScope":
		resp, err = d.dispatchGetAtSubscriptionScope(req)
	case "DeploymentOperationsClient.GetAtTenantScope":
		resp, err = d.dispatchGetAtTenantScope(req)
	case "DeploymentOperationsClient.NewListPager":
		resp, err = d.dispatchNewListPager(req)
	case "DeploymentOperationsClient.NewListAtManagementGroupScopePager":
		resp, err = d.dispatchNewListAtManagementGroupScopePager(req)
	case "DeploymentOperationsClient.NewListAtScopePager":
		resp, err = d.dispatchNewListAtScopePager(req)
	case "DeploymentOperationsClient.NewListAtSubscriptionScopePager":
		resp, err = d.dispatchNewListAtSubscriptionScopePager(req)
	case "DeploymentOperationsClient.NewListAtTenantScopePager":
		resp, err = d.dispatchNewListAtTenantScopePager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DeploymentOperationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if d.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/deployments/(?P<deploymentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operations/(?P<operationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	deploymentNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("deploymentName")])
	if err != nil {
		return nil, err
	}
	operationIDUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("operationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.Get(req.Context(), resourceGroupNameUnescaped, deploymentNameUnescaped, operationIDUnescaped, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DeploymentOperation, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DeploymentOperationsServerTransport) dispatchGetAtManagementGroupScope(req *http.Request) (*http.Response, error) {
	if d.srv.GetAtManagementGroupScope == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetAtManagementGroupScope not implemented")}
	}
	const regexStr = `/providers/Microsoft.Management/managementGroups/(?P<groupId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Resources/deployments/(?P<deploymentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operations/(?P<operationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	groupIDUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("groupId")])
	if err != nil {
		return nil, err
	}
	deploymentNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("deploymentName")])
	if err != nil {
		return nil, err
	}
	operationIDUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("operationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.GetAtManagementGroupScope(req.Context(), groupIDUnescaped, deploymentNameUnescaped, operationIDUnescaped, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DeploymentOperation, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DeploymentOperationsServerTransport) dispatchGetAtScope(req *http.Request) (*http.Response, error) {
	if d.srv.GetAtScope == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetAtScope not implemented")}
	}
	const regexStr = `/(?P<scope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Resources/deployments/(?P<deploymentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operations/(?P<operationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	scopeUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("scope")])
	if err != nil {
		return nil, err
	}
	deploymentNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("deploymentName")])
	if err != nil {
		return nil, err
	}
	operationIDUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("operationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.GetAtScope(req.Context(), scopeUnescaped, deploymentNameUnescaped, operationIDUnescaped, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DeploymentOperation, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DeploymentOperationsServerTransport) dispatchGetAtSubscriptionScope(req *http.Request) (*http.Response, error) {
	if d.srv.GetAtSubscriptionScope == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetAtSubscriptionScope not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Resources/deployments/(?P<deploymentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operations/(?P<operationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	deploymentNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("deploymentName")])
	if err != nil {
		return nil, err
	}
	operationIDUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("operationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.GetAtSubscriptionScope(req.Context(), deploymentNameUnescaped, operationIDUnescaped, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DeploymentOperation, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DeploymentOperationsServerTransport) dispatchGetAtTenantScope(req *http.Request) (*http.Response, error) {
	if d.srv.GetAtTenantScope == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetAtTenantScope not implemented")}
	}
	const regexStr = `/providers/Microsoft.Resources/deployments/(?P<deploymentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operations/(?P<operationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	deploymentNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("deploymentName")])
	if err != nil {
		return nil, err
	}
	operationIDUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("operationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.GetAtTenantScope(req.Context(), deploymentNameUnescaped, operationIDUnescaped, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DeploymentOperation, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DeploymentOperationsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := d.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/deployments/(?P<deploymentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		deploymentNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("deploymentName")])
		if err != nil {
			return nil, err
		}
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armresources.DeploymentOperationsClientListOptions
		if topParam != nil {
			options = &armresources.DeploymentOperationsClientListOptions{
				Top: topParam,
			}
		}
		resp := d.srv.NewListPager(resourceGroupNameUnescaped, deploymentNameUnescaped, options)
		newListPager = &resp
		d.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armresources.DeploymentOperationsClientListResponse, createLink func() string) {
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

func (d *DeploymentOperationsServerTransport) dispatchNewListAtManagementGroupScopePager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListAtManagementGroupScopePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListAtManagementGroupScopePager not implemented")}
	}
	newListAtManagementGroupScopePager := d.newListAtManagementGroupScopePager.get(req)
	if newListAtManagementGroupScopePager == nil {
		const regexStr = `/providers/Microsoft.Management/managementGroups/(?P<groupId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Resources/deployments/(?P<deploymentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		groupIDUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("groupId")])
		if err != nil {
			return nil, err
		}
		deploymentNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("deploymentName")])
		if err != nil {
			return nil, err
		}
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armresources.DeploymentOperationsClientListAtManagementGroupScopeOptions
		if topParam != nil {
			options = &armresources.DeploymentOperationsClientListAtManagementGroupScopeOptions{
				Top: topParam,
			}
		}
		resp := d.srv.NewListAtManagementGroupScopePager(groupIDUnescaped, deploymentNameUnescaped, options)
		newListAtManagementGroupScopePager = &resp
		d.newListAtManagementGroupScopePager.add(req, newListAtManagementGroupScopePager)
		server.PagerResponderInjectNextLinks(newListAtManagementGroupScopePager, req, func(page *armresources.DeploymentOperationsClientListAtManagementGroupScopeResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListAtManagementGroupScopePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListAtManagementGroupScopePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListAtManagementGroupScopePager) {
		d.newListAtManagementGroupScopePager.remove(req)
	}
	return resp, nil
}

func (d *DeploymentOperationsServerTransport) dispatchNewListAtScopePager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListAtScopePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListAtScopePager not implemented")}
	}
	newListAtScopePager := d.newListAtScopePager.get(req)
	if newListAtScopePager == nil {
		const regexStr = `/(?P<scope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Resources/deployments/(?P<deploymentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		scopeUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("scope")])
		if err != nil {
			return nil, err
		}
		deploymentNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("deploymentName")])
		if err != nil {
			return nil, err
		}
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armresources.DeploymentOperationsClientListAtScopeOptions
		if topParam != nil {
			options = &armresources.DeploymentOperationsClientListAtScopeOptions{
				Top: topParam,
			}
		}
		resp := d.srv.NewListAtScopePager(scopeUnescaped, deploymentNameUnescaped, options)
		newListAtScopePager = &resp
		d.newListAtScopePager.add(req, newListAtScopePager)
		server.PagerResponderInjectNextLinks(newListAtScopePager, req, func(page *armresources.DeploymentOperationsClientListAtScopeResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListAtScopePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListAtScopePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListAtScopePager) {
		d.newListAtScopePager.remove(req)
	}
	return resp, nil
}

func (d *DeploymentOperationsServerTransport) dispatchNewListAtSubscriptionScopePager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListAtSubscriptionScopePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListAtSubscriptionScopePager not implemented")}
	}
	newListAtSubscriptionScopePager := d.newListAtSubscriptionScopePager.get(req)
	if newListAtSubscriptionScopePager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft.Resources/deployments/(?P<deploymentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		deploymentNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("deploymentName")])
		if err != nil {
			return nil, err
		}
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armresources.DeploymentOperationsClientListAtSubscriptionScopeOptions
		if topParam != nil {
			options = &armresources.DeploymentOperationsClientListAtSubscriptionScopeOptions{
				Top: topParam,
			}
		}
		resp := d.srv.NewListAtSubscriptionScopePager(deploymentNameUnescaped, options)
		newListAtSubscriptionScopePager = &resp
		d.newListAtSubscriptionScopePager.add(req, newListAtSubscriptionScopePager)
		server.PagerResponderInjectNextLinks(newListAtSubscriptionScopePager, req, func(page *armresources.DeploymentOperationsClientListAtSubscriptionScopeResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListAtSubscriptionScopePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListAtSubscriptionScopePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListAtSubscriptionScopePager) {
		d.newListAtSubscriptionScopePager.remove(req)
	}
	return resp, nil
}

func (d *DeploymentOperationsServerTransport) dispatchNewListAtTenantScopePager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListAtTenantScopePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListAtTenantScopePager not implemented")}
	}
	newListAtTenantScopePager := d.newListAtTenantScopePager.get(req)
	if newListAtTenantScopePager == nil {
		const regexStr = `/providers/Microsoft.Resources/deployments/(?P<deploymentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/operations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		deploymentNameUnescaped, err := url.PathUnescape(matches[regex.SubexpIndex("deploymentName")])
		if err != nil {
			return nil, err
		}
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armresources.DeploymentOperationsClientListAtTenantScopeOptions
		if topParam != nil {
			options = &armresources.DeploymentOperationsClientListAtTenantScopeOptions{
				Top: topParam,
			}
		}
		resp := d.srv.NewListAtTenantScopePager(deploymentNameUnescaped, options)
		newListAtTenantScopePager = &resp
		d.newListAtTenantScopePager.add(req, newListAtTenantScopePager)
		server.PagerResponderInjectNextLinks(newListAtTenantScopePager, req, func(page *armresources.DeploymentOperationsClientListAtTenantScopeResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListAtTenantScopePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListAtTenantScopePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListAtTenantScopePager) {
		d.newListAtTenantScopePager.remove(req)
	}
	return resp, nil
}
