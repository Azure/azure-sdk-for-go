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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/logic/armlogic"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// IntegrationAccountSchemasServer is a fake server for instances of the armlogic.IntegrationAccountSchemasClient type.
type IntegrationAccountSchemasServer struct {
	// CreateOrUpdate is the fake for method IntegrationAccountSchemasClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, integrationAccountName string, schemaName string, schema armlogic.IntegrationAccountSchema, options *armlogic.IntegrationAccountSchemasClientCreateOrUpdateOptions) (resp azfake.Responder[armlogic.IntegrationAccountSchemasClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method IntegrationAccountSchemasClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, integrationAccountName string, schemaName string, options *armlogic.IntegrationAccountSchemasClientDeleteOptions) (resp azfake.Responder[armlogic.IntegrationAccountSchemasClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method IntegrationAccountSchemasClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, integrationAccountName string, schemaName string, options *armlogic.IntegrationAccountSchemasClientGetOptions) (resp azfake.Responder[armlogic.IntegrationAccountSchemasClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method IntegrationAccountSchemasClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, integrationAccountName string, options *armlogic.IntegrationAccountSchemasClientListOptions) (resp azfake.PagerResponder[armlogic.IntegrationAccountSchemasClientListResponse])

	// ListContentCallbackURL is the fake for method IntegrationAccountSchemasClient.ListContentCallbackURL
	// HTTP status codes to indicate success: http.StatusOK
	ListContentCallbackURL func(ctx context.Context, resourceGroupName string, integrationAccountName string, schemaName string, listContentCallbackURL armlogic.GetCallbackURLParameters, options *armlogic.IntegrationAccountSchemasClientListContentCallbackURLOptions) (resp azfake.Responder[armlogic.IntegrationAccountSchemasClientListContentCallbackURLResponse], errResp azfake.ErrorResponder)
}

// NewIntegrationAccountSchemasServerTransport creates a new instance of IntegrationAccountSchemasServerTransport with the provided implementation.
// The returned IntegrationAccountSchemasServerTransport instance is connected to an instance of armlogic.IntegrationAccountSchemasClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewIntegrationAccountSchemasServerTransport(srv *IntegrationAccountSchemasServer) *IntegrationAccountSchemasServerTransport {
	return &IntegrationAccountSchemasServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armlogic.IntegrationAccountSchemasClientListResponse]](),
	}
}

// IntegrationAccountSchemasServerTransport connects instances of armlogic.IntegrationAccountSchemasClient to instances of IntegrationAccountSchemasServer.
// Don't use this type directly, use NewIntegrationAccountSchemasServerTransport instead.
type IntegrationAccountSchemasServerTransport struct {
	srv          *IntegrationAccountSchemasServer
	newListPager *tracker[azfake.PagerResponder[armlogic.IntegrationAccountSchemasClientListResponse]]
}

// Do implements the policy.Transporter interface for IntegrationAccountSchemasServerTransport.
func (i *IntegrationAccountSchemasServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "IntegrationAccountSchemasClient.CreateOrUpdate":
		resp, err = i.dispatchCreateOrUpdate(req)
	case "IntegrationAccountSchemasClient.Delete":
		resp, err = i.dispatchDelete(req)
	case "IntegrationAccountSchemasClient.Get":
		resp, err = i.dispatchGet(req)
	case "IntegrationAccountSchemasClient.NewListPager":
		resp, err = i.dispatchNewListPager(req)
	case "IntegrationAccountSchemasClient.ListContentCallbackURL":
		resp, err = i.dispatchListContentCallbackURL(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (i *IntegrationAccountSchemasServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if i.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Logic/integrationAccounts/(?P<integrationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/schemas/(?P<schemaName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armlogic.IntegrationAccountSchema](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	integrationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("integrationAccountName")])
	if err != nil {
		return nil, err
	}
	schemaNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("schemaName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, integrationAccountNameParam, schemaNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).IntegrationAccountSchema, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *IntegrationAccountSchemasServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if i.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Logic/integrationAccounts/(?P<integrationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/schemas/(?P<schemaName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	integrationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("integrationAccountName")])
	if err != nil {
		return nil, err
	}
	schemaNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("schemaName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Delete(req.Context(), resourceGroupNameParam, integrationAccountNameParam, schemaNameParam, nil)
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

func (i *IntegrationAccountSchemasServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if i.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Logic/integrationAccounts/(?P<integrationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/schemas/(?P<schemaName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	integrationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("integrationAccountName")])
	if err != nil {
		return nil, err
	}
	schemaNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("schemaName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Get(req.Context(), resourceGroupNameParam, integrationAccountNameParam, schemaNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).IntegrationAccountSchema, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *IntegrationAccountSchemasServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if i.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := i.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Logic/integrationAccounts/(?P<integrationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/schemas`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		integrationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("integrationAccountName")])
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
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		var options *armlogic.IntegrationAccountSchemasClientListOptions
		if topParam != nil || filterParam != nil {
			options = &armlogic.IntegrationAccountSchemasClientListOptions{
				Top:    topParam,
				Filter: filterParam,
			}
		}
		resp := i.srv.NewListPager(resourceGroupNameParam, integrationAccountNameParam, options)
		newListPager = &resp
		i.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armlogic.IntegrationAccountSchemasClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		i.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		i.newListPager.remove(req)
	}
	return resp, nil
}

func (i *IntegrationAccountSchemasServerTransport) dispatchListContentCallbackURL(req *http.Request) (*http.Response, error) {
	if i.srv.ListContentCallbackURL == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListContentCallbackURL not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Logic/integrationAccounts/(?P<integrationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/schemas/(?P<schemaName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listContentCallbackUrl`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armlogic.GetCallbackURLParameters](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	integrationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("integrationAccountName")])
	if err != nil {
		return nil, err
	}
	schemaNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("schemaName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.ListContentCallbackURL(req.Context(), resourceGroupNameParam, integrationAccountNameParam, schemaNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).WorkflowTriggerCallbackURL, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
