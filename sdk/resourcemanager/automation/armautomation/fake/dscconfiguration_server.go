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
	"reflect"
	"regexp"
	"strconv"
)

// DscConfigurationServer is a fake server for instances of the armautomation.DscConfigurationClient type.
type DscConfigurationServer struct {
	// CreateOrUpdateWithJSON is the fake for method DscConfigurationClient.CreateOrUpdateWithJSON
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdateWithJSON func(ctx context.Context, resourceGroupName string, automationAccountName string, configurationName string, parameters armautomation.DscConfigurationCreateOrUpdateParameters, options *armautomation.DscConfigurationClientCreateOrUpdateWithJSONOptions) (resp azfake.Responder[armautomation.DscConfigurationClientCreateOrUpdateWithJSONResponse], errResp azfake.ErrorResponder)

	// CreateOrUpdateWithText is the fake for method DscConfigurationClient.CreateOrUpdateWithText
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdateWithText func(ctx context.Context, resourceGroupName string, automationAccountName string, configurationName string, parameters string, options *armautomation.DscConfigurationClientCreateOrUpdateWithTextOptions) (resp azfake.Responder[armautomation.DscConfigurationClientCreateOrUpdateWithTextResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method DscConfigurationClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, automationAccountName string, configurationName string, options *armautomation.DscConfigurationClientDeleteOptions) (resp azfake.Responder[armautomation.DscConfigurationClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method DscConfigurationClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, automationAccountName string, configurationName string, options *armautomation.DscConfigurationClientGetOptions) (resp azfake.Responder[armautomation.DscConfigurationClientGetResponse], errResp azfake.ErrorResponder)

	// GetContent is the fake for method DscConfigurationClient.GetContent
	// HTTP status codes to indicate success: http.StatusOK
	GetContent func(ctx context.Context, resourceGroupName string, automationAccountName string, configurationName string, options *armautomation.DscConfigurationClientGetContentOptions) (resp azfake.Responder[armautomation.DscConfigurationClientGetContentResponse], errResp azfake.ErrorResponder)

	// NewListByAutomationAccountPager is the fake for method DscConfigurationClient.NewListByAutomationAccountPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByAutomationAccountPager func(resourceGroupName string, automationAccountName string, options *armautomation.DscConfigurationClientListByAutomationAccountOptions) (resp azfake.PagerResponder[armautomation.DscConfigurationClientListByAutomationAccountResponse])

	// UpdateWithJSON is the fake for method DscConfigurationClient.UpdateWithJSON
	// HTTP status codes to indicate success: http.StatusOK
	UpdateWithJSON func(ctx context.Context, resourceGroupName string, automationAccountName string, configurationName string, options *armautomation.DscConfigurationClientUpdateWithJSONOptions) (resp azfake.Responder[armautomation.DscConfigurationClientUpdateWithJSONResponse], errResp azfake.ErrorResponder)

	// UpdateWithText is the fake for method DscConfigurationClient.UpdateWithText
	// HTTP status codes to indicate success: http.StatusOK
	UpdateWithText func(ctx context.Context, resourceGroupName string, automationAccountName string, configurationName string, options *armautomation.DscConfigurationClientUpdateWithTextOptions) (resp azfake.Responder[armautomation.DscConfigurationClientUpdateWithTextResponse], errResp azfake.ErrorResponder)
}

// NewDscConfigurationServerTransport creates a new instance of DscConfigurationServerTransport with the provided implementation.
// The returned DscConfigurationServerTransport instance is connected to an instance of armautomation.DscConfigurationClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDscConfigurationServerTransport(srv *DscConfigurationServer) *DscConfigurationServerTransport {
	return &DscConfigurationServerTransport{
		srv:                             srv,
		newListByAutomationAccountPager: newTracker[azfake.PagerResponder[armautomation.DscConfigurationClientListByAutomationAccountResponse]](),
	}
}

// DscConfigurationServerTransport connects instances of armautomation.DscConfigurationClient to instances of DscConfigurationServer.
// Don't use this type directly, use NewDscConfigurationServerTransport instead.
type DscConfigurationServerTransport struct {
	srv                             *DscConfigurationServer
	newListByAutomationAccountPager *tracker[azfake.PagerResponder[armautomation.DscConfigurationClientListByAutomationAccountResponse]]
}

// Do implements the policy.Transporter interface for DscConfigurationServerTransport.
func (d *DscConfigurationServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "DscConfigurationClient.CreateOrUpdateWithJSON":
		resp, err = d.dispatchCreateOrUpdateWithJSON(req)
	case "DscConfigurationClient.CreateOrUpdateWithText":
		resp, err = d.dispatchCreateOrUpdateWithText(req)
	case "DscConfigurationClient.Delete":
		resp, err = d.dispatchDelete(req)
	case "DscConfigurationClient.Get":
		resp, err = d.dispatchGet(req)
	case "DscConfigurationClient.GetContent":
		resp, err = d.dispatchGetContent(req)
	case "DscConfigurationClient.NewListByAutomationAccountPager":
		resp, err = d.dispatchNewListByAutomationAccountPager(req)
	case "DscConfigurationClient.UpdateWithJSON":
		resp, err = d.dispatchUpdateWithJSON(req)
	case "DscConfigurationClient.UpdateWithText":
		resp, err = d.dispatchUpdateWithText(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DscConfigurationServerTransport) dispatchCreateOrUpdateWithJSON(req *http.Request) (*http.Response, error) {
	if d.srv.CreateOrUpdateWithJSON == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdateWithJSON not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automation/automationAccounts/(?P<automationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armautomation.DscConfigurationCreateOrUpdateParameters](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	automationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("automationAccountName")])
	if err != nil {
		return nil, err
	}
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.CreateOrUpdateWithJSON(req.Context(), resourceGroupNameParam, automationAccountNameParam, configurationNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DscConfiguration, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DscConfigurationServerTransport) dispatchCreateOrUpdateWithText(req *http.Request) (*http.Response, error) {
	if d.srv.CreateOrUpdateWithText == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdateWithText not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automation/automationAccounts/(?P<automationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsText(req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	automationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("automationAccountName")])
	if err != nil {
		return nil, err
	}
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.CreateOrUpdateWithText(req.Context(), resourceGroupNameParam, automationAccountNameParam, configurationNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DscConfiguration, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DscConfigurationServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if d.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automation/automationAccounts/(?P<automationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.Delete(req.Context(), resourceGroupNameParam, automationAccountNameParam, configurationNameParam, nil)
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

func (d *DscConfigurationServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if d.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automation/automationAccounts/(?P<automationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.Get(req.Context(), resourceGroupNameParam, automationAccountNameParam, configurationNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DscConfiguration, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DscConfigurationServerTransport) dispatchGetContent(req *http.Request) (*http.Response, error) {
	if d.srv.GetContent == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetContent not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automation/automationAccounts/(?P<automationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/content`
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
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.GetContent(req.Context(), resourceGroupNameParam, automationAccountNameParam, configurationNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsText(respContent, server.GetResponse(respr).Value, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DscConfigurationServerTransport) dispatchNewListByAutomationAccountPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListByAutomationAccountPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByAutomationAccountPager not implemented")}
	}
	newListByAutomationAccountPager := d.newListByAutomationAccountPager.get(req)
	if newListByAutomationAccountPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automation/automationAccounts/(?P<automationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations`
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
		automationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("automationAccountName")])
		if err != nil {
			return nil, err
		}
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		skipUnescaped, err := url.QueryUnescape(qp.Get("$skip"))
		if err != nil {
			return nil, err
		}
		skipParam, err := parseOptional(skipUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
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
		inlinecountUnescaped, err := url.QueryUnescape(qp.Get("$inlinecount"))
		if err != nil {
			return nil, err
		}
		inlinecountParam := getOptional(inlinecountUnescaped)
		var options *armautomation.DscConfigurationClientListByAutomationAccountOptions
		if filterParam != nil || skipParam != nil || topParam != nil || inlinecountParam != nil {
			options = &armautomation.DscConfigurationClientListByAutomationAccountOptions{
				Filter:      filterParam,
				Skip:        skipParam,
				Top:         topParam,
				Inlinecount: inlinecountParam,
			}
		}
		resp := d.srv.NewListByAutomationAccountPager(resourceGroupNameParam, automationAccountNameParam, options)
		newListByAutomationAccountPager = &resp
		d.newListByAutomationAccountPager.add(req, newListByAutomationAccountPager)
		server.PagerResponderInjectNextLinks(newListByAutomationAccountPager, req, func(page *armautomation.DscConfigurationClientListByAutomationAccountResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByAutomationAccountPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListByAutomationAccountPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByAutomationAccountPager) {
		d.newListByAutomationAccountPager.remove(req)
	}
	return resp, nil
}

func (d *DscConfigurationServerTransport) dispatchUpdateWithJSON(req *http.Request) (*http.Response, error) {
	if d.srv.UpdateWithJSON == nil {
		return nil, &nonRetriableError{errors.New("fake for method UpdateWithJSON not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automation/automationAccounts/(?P<automationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armautomation.DscConfigurationUpdateParameters](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	automationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("automationAccountName")])
	if err != nil {
		return nil, err
	}
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	var options *armautomation.DscConfigurationClientUpdateWithJSONOptions
	if !reflect.ValueOf(body).IsZero() {
		options = &armautomation.DscConfigurationClientUpdateWithJSONOptions{
			Parameters: &body,
		}
	}
	respr, errRespr := d.srv.UpdateWithJSON(req.Context(), resourceGroupNameParam, automationAccountNameParam, configurationNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DscConfiguration, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DscConfigurationServerTransport) dispatchUpdateWithText(req *http.Request) (*http.Response, error) {
	if d.srv.UpdateWithText == nil {
		return nil, &nonRetriableError{errors.New("fake for method UpdateWithText not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automation/automationAccounts/(?P<automationAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsText(req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	automationAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("automationAccountName")])
	if err != nil {
		return nil, err
	}
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	var options *armautomation.DscConfigurationClientUpdateWithTextOptions
	if !reflect.ValueOf(body).IsZero() {
		options = &armautomation.DscConfigurationClientUpdateWithTextOptions{
			Parameters: &body,
		}
	}
	respr, errRespr := d.srv.UpdateWithText(req.Context(), resourceGroupNameParam, automationAccountNameParam, configurationNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DscConfiguration, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
