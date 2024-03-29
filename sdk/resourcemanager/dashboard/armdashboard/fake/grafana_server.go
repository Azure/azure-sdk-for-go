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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/dashboard/armdashboard"
	"net/http"
	"net/url"
	"regexp"
)

// GrafanaServer is a fake server for instances of the armdashboard.GrafanaClient type.
type GrafanaServer struct {
	// CheckEnterpriseDetails is the fake for method GrafanaClient.CheckEnterpriseDetails
	// HTTP status codes to indicate success: http.StatusOK
	CheckEnterpriseDetails func(ctx context.Context, resourceGroupName string, workspaceName string, options *armdashboard.GrafanaClientCheckEnterpriseDetailsOptions) (resp azfake.Responder[armdashboard.GrafanaClientCheckEnterpriseDetailsResponse], errResp azfake.ErrorResponder)

	// BeginCreate is the fake for method GrafanaClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, resourceGroupName string, workspaceName string, requestBodyParameters armdashboard.ManagedGrafana, options *armdashboard.GrafanaClientBeginCreateOptions) (resp azfake.PollerResponder[armdashboard.GrafanaClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method GrafanaClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, workspaceName string, options *armdashboard.GrafanaClientBeginDeleteOptions) (resp azfake.PollerResponder[armdashboard.GrafanaClientDeleteResponse], errResp azfake.ErrorResponder)

	// FetchAvailablePlugins is the fake for method GrafanaClient.FetchAvailablePlugins
	// HTTP status codes to indicate success: http.StatusOK
	FetchAvailablePlugins func(ctx context.Context, resourceGroupName string, workspaceName string, options *armdashboard.GrafanaClientFetchAvailablePluginsOptions) (resp azfake.Responder[armdashboard.GrafanaClientFetchAvailablePluginsResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method GrafanaClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, workspaceName string, options *armdashboard.GrafanaClientGetOptions) (resp azfake.Responder[armdashboard.GrafanaClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method GrafanaClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armdashboard.GrafanaClientListOptions) (resp azfake.PagerResponder[armdashboard.GrafanaClientListResponse])

	// NewListByResourceGroupPager is the fake for method GrafanaClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armdashboard.GrafanaClientListByResourceGroupOptions) (resp azfake.PagerResponder[armdashboard.GrafanaClientListByResourceGroupResponse])

	// Update is the fake for method GrafanaClient.Update
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	Update func(ctx context.Context, resourceGroupName string, workspaceName string, requestBodyParameters armdashboard.ManagedGrafanaUpdateParameters, options *armdashboard.GrafanaClientUpdateOptions) (resp azfake.Responder[armdashboard.GrafanaClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewGrafanaServerTransport creates a new instance of GrafanaServerTransport with the provided implementation.
// The returned GrafanaServerTransport instance is connected to an instance of armdashboard.GrafanaClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewGrafanaServerTransport(srv *GrafanaServer) *GrafanaServerTransport {
	return &GrafanaServerTransport{
		srv:                         srv,
		beginCreate:                 newTracker[azfake.PollerResponder[armdashboard.GrafanaClientCreateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armdashboard.GrafanaClientDeleteResponse]](),
		newListPager:                newTracker[azfake.PagerResponder[armdashboard.GrafanaClientListResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armdashboard.GrafanaClientListByResourceGroupResponse]](),
	}
}

// GrafanaServerTransport connects instances of armdashboard.GrafanaClient to instances of GrafanaServer.
// Don't use this type directly, use NewGrafanaServerTransport instead.
type GrafanaServerTransport struct {
	srv                         *GrafanaServer
	beginCreate                 *tracker[azfake.PollerResponder[armdashboard.GrafanaClientCreateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armdashboard.GrafanaClientDeleteResponse]]
	newListPager                *tracker[azfake.PagerResponder[armdashboard.GrafanaClientListResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armdashboard.GrafanaClientListByResourceGroupResponse]]
}

// Do implements the policy.Transporter interface for GrafanaServerTransport.
func (g *GrafanaServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "GrafanaClient.CheckEnterpriseDetails":
		resp, err = g.dispatchCheckEnterpriseDetails(req)
	case "GrafanaClient.BeginCreate":
		resp, err = g.dispatchBeginCreate(req)
	case "GrafanaClient.BeginDelete":
		resp, err = g.dispatchBeginDelete(req)
	case "GrafanaClient.FetchAvailablePlugins":
		resp, err = g.dispatchFetchAvailablePlugins(req)
	case "GrafanaClient.Get":
		resp, err = g.dispatchGet(req)
	case "GrafanaClient.NewListPager":
		resp, err = g.dispatchNewListPager(req)
	case "GrafanaClient.NewListByResourceGroupPager":
		resp, err = g.dispatchNewListByResourceGroupPager(req)
	case "GrafanaClient.Update":
		resp, err = g.dispatchUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *GrafanaServerTransport) dispatchCheckEnterpriseDetails(req *http.Request) (*http.Response, error) {
	if g.srv.CheckEnterpriseDetails == nil {
		return nil, &nonRetriableError{errors.New("fake for method CheckEnterpriseDetails not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Dashboard/grafana/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/checkEnterpriseDetails`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.CheckEnterpriseDetails(req.Context(), resourceGroupNameParam, workspaceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).EnterpriseDetails, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GrafanaServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if g.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := g.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Dashboard/grafana/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armdashboard.ManagedGrafana](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := g.srv.BeginCreate(req.Context(), resourceGroupNameParam, workspaceNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		g.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		g.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		g.beginCreate.remove(req)
	}

	return resp, nil
}

func (g *GrafanaServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if g.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := g.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Dashboard/grafana/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := g.srv.BeginDelete(req.Context(), resourceGroupNameParam, workspaceNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		g.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		g.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		g.beginDelete.remove(req)
	}

	return resp, nil
}

func (g *GrafanaServerTransport) dispatchFetchAvailablePlugins(req *http.Request) (*http.Response, error) {
	if g.srv.FetchAvailablePlugins == nil {
		return nil, &nonRetriableError{errors.New("fake for method FetchAvailablePlugins not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Dashboard/grafana/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/fetchAvailablePlugins`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.FetchAvailablePlugins(req.Context(), resourceGroupNameParam, workspaceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).GrafanaAvailablePluginListResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GrafanaServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if g.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Dashboard/grafana/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.Get(req.Context(), resourceGroupNameParam, workspaceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ManagedGrafana, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GrafanaServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if g.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := g.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Dashboard/grafana`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := g.srv.NewListPager(nil)
		newListPager = &resp
		g.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armdashboard.GrafanaClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		g.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		g.newListPager.remove(req)
	}
	return resp, nil
}

func (g *GrafanaServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if g.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := g.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Dashboard/grafana`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := g.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		g.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armdashboard.GrafanaClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		g.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		g.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (g *GrafanaServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if g.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Dashboard/grafana/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdashboard.ManagedGrafanaUpdateParameters](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.Update(req.Context(), resourceGroupNameParam, workspaceNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusAccepted}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ManagedGrafana, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).AzureAsyncOperation; val != nil {
		resp.Header.Set("Azure-AsyncOperation", *val)
	}
	return resp, nil
}
