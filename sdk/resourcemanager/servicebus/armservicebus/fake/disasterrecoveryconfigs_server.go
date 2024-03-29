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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicebus/armservicebus"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

// DisasterRecoveryConfigsServer is a fake server for instances of the armservicebus.DisasterRecoveryConfigsClient type.
type DisasterRecoveryConfigsServer struct {
	// BreakPairing is the fake for method DisasterRecoveryConfigsClient.BreakPairing
	// HTTP status codes to indicate success: http.StatusOK
	BreakPairing func(ctx context.Context, resourceGroupName string, namespaceName string, alias string, options *armservicebus.DisasterRecoveryConfigsClientBreakPairingOptions) (resp azfake.Responder[armservicebus.DisasterRecoveryConfigsClientBreakPairingResponse], errResp azfake.ErrorResponder)

	// CheckNameAvailability is the fake for method DisasterRecoveryConfigsClient.CheckNameAvailability
	// HTTP status codes to indicate success: http.StatusOK
	CheckNameAvailability func(ctx context.Context, resourceGroupName string, namespaceName string, parameters armservicebus.CheckNameAvailability, options *armservicebus.DisasterRecoveryConfigsClientCheckNameAvailabilityOptions) (resp azfake.Responder[armservicebus.DisasterRecoveryConfigsClientCheckNameAvailabilityResponse], errResp azfake.ErrorResponder)

	// CreateOrUpdate is the fake for method DisasterRecoveryConfigsClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, namespaceName string, alias string, parameters armservicebus.ArmDisasterRecovery, options *armservicebus.DisasterRecoveryConfigsClientCreateOrUpdateOptions) (resp azfake.Responder[armservicebus.DisasterRecoveryConfigsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method DisasterRecoveryConfigsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, namespaceName string, alias string, options *armservicebus.DisasterRecoveryConfigsClientDeleteOptions) (resp azfake.Responder[armservicebus.DisasterRecoveryConfigsClientDeleteResponse], errResp azfake.ErrorResponder)

	// FailOver is the fake for method DisasterRecoveryConfigsClient.FailOver
	// HTTP status codes to indicate success: http.StatusOK
	FailOver func(ctx context.Context, resourceGroupName string, namespaceName string, alias string, options *armservicebus.DisasterRecoveryConfigsClientFailOverOptions) (resp azfake.Responder[armservicebus.DisasterRecoveryConfigsClientFailOverResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method DisasterRecoveryConfigsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, namespaceName string, alias string, options *armservicebus.DisasterRecoveryConfigsClientGetOptions) (resp azfake.Responder[armservicebus.DisasterRecoveryConfigsClientGetResponse], errResp azfake.ErrorResponder)

	// GetAuthorizationRule is the fake for method DisasterRecoveryConfigsClient.GetAuthorizationRule
	// HTTP status codes to indicate success: http.StatusOK
	GetAuthorizationRule func(ctx context.Context, resourceGroupName string, namespaceName string, alias string, authorizationRuleName string, options *armservicebus.DisasterRecoveryConfigsClientGetAuthorizationRuleOptions) (resp azfake.Responder[armservicebus.DisasterRecoveryConfigsClientGetAuthorizationRuleResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method DisasterRecoveryConfigsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, namespaceName string, options *armservicebus.DisasterRecoveryConfigsClientListOptions) (resp azfake.PagerResponder[armservicebus.DisasterRecoveryConfigsClientListResponse])

	// NewListAuthorizationRulesPager is the fake for method DisasterRecoveryConfigsClient.NewListAuthorizationRulesPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListAuthorizationRulesPager func(resourceGroupName string, namespaceName string, alias string, options *armservicebus.DisasterRecoveryConfigsClientListAuthorizationRulesOptions) (resp azfake.PagerResponder[armservicebus.DisasterRecoveryConfigsClientListAuthorizationRulesResponse])

	// ListKeys is the fake for method DisasterRecoveryConfigsClient.ListKeys
	// HTTP status codes to indicate success: http.StatusOK
	ListKeys func(ctx context.Context, resourceGroupName string, namespaceName string, alias string, authorizationRuleName string, options *armservicebus.DisasterRecoveryConfigsClientListKeysOptions) (resp azfake.Responder[armservicebus.DisasterRecoveryConfigsClientListKeysResponse], errResp azfake.ErrorResponder)
}

// NewDisasterRecoveryConfigsServerTransport creates a new instance of DisasterRecoveryConfigsServerTransport with the provided implementation.
// The returned DisasterRecoveryConfigsServerTransport instance is connected to an instance of armservicebus.DisasterRecoveryConfigsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDisasterRecoveryConfigsServerTransport(srv *DisasterRecoveryConfigsServer) *DisasterRecoveryConfigsServerTransport {
	return &DisasterRecoveryConfigsServerTransport{
		srv:                            srv,
		newListPager:                   newTracker[azfake.PagerResponder[armservicebus.DisasterRecoveryConfigsClientListResponse]](),
		newListAuthorizationRulesPager: newTracker[azfake.PagerResponder[armservicebus.DisasterRecoveryConfigsClientListAuthorizationRulesResponse]](),
	}
}

// DisasterRecoveryConfigsServerTransport connects instances of armservicebus.DisasterRecoveryConfigsClient to instances of DisasterRecoveryConfigsServer.
// Don't use this type directly, use NewDisasterRecoveryConfigsServerTransport instead.
type DisasterRecoveryConfigsServerTransport struct {
	srv                            *DisasterRecoveryConfigsServer
	newListPager                   *tracker[azfake.PagerResponder[armservicebus.DisasterRecoveryConfigsClientListResponse]]
	newListAuthorizationRulesPager *tracker[azfake.PagerResponder[armservicebus.DisasterRecoveryConfigsClientListAuthorizationRulesResponse]]
}

// Do implements the policy.Transporter interface for DisasterRecoveryConfigsServerTransport.
func (d *DisasterRecoveryConfigsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "DisasterRecoveryConfigsClient.BreakPairing":
		resp, err = d.dispatchBreakPairing(req)
	case "DisasterRecoveryConfigsClient.CheckNameAvailability":
		resp, err = d.dispatchCheckNameAvailability(req)
	case "DisasterRecoveryConfigsClient.CreateOrUpdate":
		resp, err = d.dispatchCreateOrUpdate(req)
	case "DisasterRecoveryConfigsClient.Delete":
		resp, err = d.dispatchDelete(req)
	case "DisasterRecoveryConfigsClient.FailOver":
		resp, err = d.dispatchFailOver(req)
	case "DisasterRecoveryConfigsClient.Get":
		resp, err = d.dispatchGet(req)
	case "DisasterRecoveryConfigsClient.GetAuthorizationRule":
		resp, err = d.dispatchGetAuthorizationRule(req)
	case "DisasterRecoveryConfigsClient.NewListPager":
		resp, err = d.dispatchNewListPager(req)
	case "DisasterRecoveryConfigsClient.NewListAuthorizationRulesPager":
		resp, err = d.dispatchNewListAuthorizationRulesPager(req)
	case "DisasterRecoveryConfigsClient.ListKeys":
		resp, err = d.dispatchListKeys(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DisasterRecoveryConfigsServerTransport) dispatchBreakPairing(req *http.Request) (*http.Response, error) {
	if d.srv.BreakPairing == nil {
		return nil, &nonRetriableError{errors.New("fake for method BreakPairing not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceBus/namespaces/(?P<namespaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/disasterRecoveryConfigs/(?P<alias>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/breakPairing`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	namespaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("namespaceName")])
	if err != nil {
		return nil, err
	}
	aliasParam, err := url.PathUnescape(matches[regex.SubexpIndex("alias")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.BreakPairing(req.Context(), resourceGroupNameParam, namespaceNameParam, aliasParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DisasterRecoveryConfigsServerTransport) dispatchCheckNameAvailability(req *http.Request) (*http.Response, error) {
	if d.srv.CheckNameAvailability == nil {
		return nil, &nonRetriableError{errors.New("fake for method CheckNameAvailability not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceBus/namespaces/(?P<namespaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/disasterRecoveryConfigs/CheckNameAvailability`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armservicebus.CheckNameAvailability](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	namespaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("namespaceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.CheckNameAvailability(req.Context(), resourceGroupNameParam, namespaceNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CheckNameAvailabilityResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DisasterRecoveryConfigsServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if d.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceBus/namespaces/(?P<namespaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/disasterRecoveryConfigs/(?P<alias>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armservicebus.ArmDisasterRecovery](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	namespaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("namespaceName")])
	if err != nil {
		return nil, err
	}
	aliasParam, err := url.PathUnescape(matches[regex.SubexpIndex("alias")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, namespaceNameParam, aliasParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ArmDisasterRecovery, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DisasterRecoveryConfigsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if d.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceBus/namespaces/(?P<namespaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/disasterRecoveryConfigs/(?P<alias>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	namespaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("namespaceName")])
	if err != nil {
		return nil, err
	}
	aliasParam, err := url.PathUnescape(matches[regex.SubexpIndex("alias")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.Delete(req.Context(), resourceGroupNameParam, namespaceNameParam, aliasParam, nil)
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

func (d *DisasterRecoveryConfigsServerTransport) dispatchFailOver(req *http.Request) (*http.Response, error) {
	if d.srv.FailOver == nil {
		return nil, &nonRetriableError{errors.New("fake for method FailOver not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceBus/namespaces/(?P<namespaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/disasterRecoveryConfigs/(?P<alias>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/failover`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armservicebus.FailoverProperties](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	namespaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("namespaceName")])
	if err != nil {
		return nil, err
	}
	aliasParam, err := url.PathUnescape(matches[regex.SubexpIndex("alias")])
	if err != nil {
		return nil, err
	}
	var options *armservicebus.DisasterRecoveryConfigsClientFailOverOptions
	if !reflect.ValueOf(body).IsZero() {
		options = &armservicebus.DisasterRecoveryConfigsClientFailOverOptions{
			Parameters: &body,
		}
	}
	respr, errRespr := d.srv.FailOver(req.Context(), resourceGroupNameParam, namespaceNameParam, aliasParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DisasterRecoveryConfigsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if d.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceBus/namespaces/(?P<namespaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/disasterRecoveryConfigs/(?P<alias>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	namespaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("namespaceName")])
	if err != nil {
		return nil, err
	}
	aliasParam, err := url.PathUnescape(matches[regex.SubexpIndex("alias")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.Get(req.Context(), resourceGroupNameParam, namespaceNameParam, aliasParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ArmDisasterRecovery, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DisasterRecoveryConfigsServerTransport) dispatchGetAuthorizationRule(req *http.Request) (*http.Response, error) {
	if d.srv.GetAuthorizationRule == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetAuthorizationRule not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceBus/namespaces/(?P<namespaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/disasterRecoveryConfigs/(?P<alias>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizationRules/(?P<authorizationRuleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	namespaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("namespaceName")])
	if err != nil {
		return nil, err
	}
	aliasParam, err := url.PathUnescape(matches[regex.SubexpIndex("alias")])
	if err != nil {
		return nil, err
	}
	authorizationRuleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationRuleName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.GetAuthorizationRule(req.Context(), resourceGroupNameParam, namespaceNameParam, aliasParam, authorizationRuleNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SBAuthorizationRule, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DisasterRecoveryConfigsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := d.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceBus/namespaces/(?P<namespaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/disasterRecoveryConfigs`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		namespaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("namespaceName")])
		if err != nil {
			return nil, err
		}
		resp := d.srv.NewListPager(resourceGroupNameParam, namespaceNameParam, nil)
		newListPager = &resp
		d.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armservicebus.DisasterRecoveryConfigsClientListResponse, createLink func() string) {
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

func (d *DisasterRecoveryConfigsServerTransport) dispatchNewListAuthorizationRulesPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListAuthorizationRulesPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListAuthorizationRulesPager not implemented")}
	}
	newListAuthorizationRulesPager := d.newListAuthorizationRulesPager.get(req)
	if newListAuthorizationRulesPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceBus/namespaces/(?P<namespaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/disasterRecoveryConfigs/(?P<alias>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizationRules`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		namespaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("namespaceName")])
		if err != nil {
			return nil, err
		}
		aliasParam, err := url.PathUnescape(matches[regex.SubexpIndex("alias")])
		if err != nil {
			return nil, err
		}
		resp := d.srv.NewListAuthorizationRulesPager(resourceGroupNameParam, namespaceNameParam, aliasParam, nil)
		newListAuthorizationRulesPager = &resp
		d.newListAuthorizationRulesPager.add(req, newListAuthorizationRulesPager)
		server.PagerResponderInjectNextLinks(newListAuthorizationRulesPager, req, func(page *armservicebus.DisasterRecoveryConfigsClientListAuthorizationRulesResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListAuthorizationRulesPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListAuthorizationRulesPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListAuthorizationRulesPager) {
		d.newListAuthorizationRulesPager.remove(req)
	}
	return resp, nil
}

func (d *DisasterRecoveryConfigsServerTransport) dispatchListKeys(req *http.Request) (*http.Response, error) {
	if d.srv.ListKeys == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListKeys not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceBus/namespaces/(?P<namespaceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/disasterRecoveryConfigs/(?P<alias>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/authorizationRules/(?P<authorizationRuleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listKeys`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	namespaceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("namespaceName")])
	if err != nil {
		return nil, err
	}
	aliasParam, err := url.PathUnescape(matches[regex.SubexpIndex("alias")])
	if err != nil {
		return nil, err
	}
	authorizationRuleNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("authorizationRuleName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.ListKeys(req.Context(), resourceGroupNameParam, namespaceNameParam, aliasParam, authorizationRuleNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AccessKeys, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
