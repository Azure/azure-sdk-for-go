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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/securityinsights/armsecurityinsights/v2"
	"net/http"
	"net/url"
	"regexp"
)

// EntityQueryTemplatesServer is a fake server for instances of the armsecurityinsights.EntityQueryTemplatesClient type.
type EntityQueryTemplatesServer struct {
	// Get is the fake for method EntityQueryTemplatesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, workspaceName string, entityQueryTemplateID string, options *armsecurityinsights.EntityQueryTemplatesClientGetOptions) (resp azfake.Responder[armsecurityinsights.EntityQueryTemplatesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method EntityQueryTemplatesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, workspaceName string, options *armsecurityinsights.EntityQueryTemplatesClientListOptions) (resp azfake.PagerResponder[armsecurityinsights.EntityQueryTemplatesClientListResponse])
}

// NewEntityQueryTemplatesServerTransport creates a new instance of EntityQueryTemplatesServerTransport with the provided implementation.
// The returned EntityQueryTemplatesServerTransport instance is connected to an instance of armsecurityinsights.EntityQueryTemplatesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewEntityQueryTemplatesServerTransport(srv *EntityQueryTemplatesServer) *EntityQueryTemplatesServerTransport {
	return &EntityQueryTemplatesServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armsecurityinsights.EntityQueryTemplatesClientListResponse]](),
	}
}

// EntityQueryTemplatesServerTransport connects instances of armsecurityinsights.EntityQueryTemplatesClient to instances of EntityQueryTemplatesServer.
// Don't use this type directly, use NewEntityQueryTemplatesServerTransport instead.
type EntityQueryTemplatesServerTransport struct {
	srv          *EntityQueryTemplatesServer
	newListPager *tracker[azfake.PagerResponder[armsecurityinsights.EntityQueryTemplatesClientListResponse]]
}

// Do implements the policy.Transporter interface for EntityQueryTemplatesServerTransport.
func (e *EntityQueryTemplatesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "EntityQueryTemplatesClient.Get":
		resp, err = e.dispatchGet(req)
	case "EntityQueryTemplatesClient.NewListPager":
		resp, err = e.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (e *EntityQueryTemplatesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if e.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.OperationalInsights/workspaces/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.SecurityInsights/entityQueryTemplates/(?P<entityQueryTemplateId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	entityQueryTemplateIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("entityQueryTemplateId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := e.srv.Get(req.Context(), resourceGroupNameParam, workspaceNameParam, entityQueryTemplateIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).EntityQueryTemplateClassification, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *EntityQueryTemplatesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if e.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := e.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.OperationalInsights/workspaces/(?P<workspaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.SecurityInsights/entityQueryTemplates`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		kindUnescaped, err := url.QueryUnescape(qp.Get("kind"))
		if err != nil {
			return nil, err
		}
		kindParam := getOptional(armsecurityinsights.Enum15(kindUnescaped))
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		workspaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceName")])
		if err != nil {
			return nil, err
		}
		var options *armsecurityinsights.EntityQueryTemplatesClientListOptions
		if kindParam != nil {
			options = &armsecurityinsights.EntityQueryTemplatesClientListOptions{
				Kind: kindParam,
			}
		}
		resp := e.srv.NewListPager(resourceGroupNameParam, workspaceNameParam, options)
		newListPager = &resp
		e.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armsecurityinsights.EntityQueryTemplatesClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
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