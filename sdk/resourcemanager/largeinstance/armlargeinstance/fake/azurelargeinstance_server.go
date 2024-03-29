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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/largeinstance/armlargeinstance"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

// AzureLargeInstanceServer is a fake server for instances of the armlargeinstance.AzureLargeInstanceClient type.
type AzureLargeInstanceServer struct {
	// Get is the fake for method AzureLargeInstanceClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, azureLargeInstanceName string, options *armlargeinstance.AzureLargeInstanceClientGetOptions) (resp azfake.Responder[armlargeinstance.AzureLargeInstanceClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method AzureLargeInstanceClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armlargeinstance.AzureLargeInstanceClientListByResourceGroupOptions) (resp azfake.PagerResponder[armlargeinstance.AzureLargeInstanceClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method AzureLargeInstanceClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armlargeinstance.AzureLargeInstanceClientListBySubscriptionOptions) (resp azfake.PagerResponder[armlargeinstance.AzureLargeInstanceClientListBySubscriptionResponse])

	// BeginRestart is the fake for method AzureLargeInstanceClient.BeginRestart
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginRestart func(ctx context.Context, resourceGroupName string, azureLargeInstanceName string, options *armlargeinstance.AzureLargeInstanceClientBeginRestartOptions) (resp azfake.PollerResponder[armlargeinstance.AzureLargeInstanceClientRestartResponse], errResp azfake.ErrorResponder)

	// BeginShutdown is the fake for method AzureLargeInstanceClient.BeginShutdown
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginShutdown func(ctx context.Context, resourceGroupName string, azureLargeInstanceName string, options *armlargeinstance.AzureLargeInstanceClientBeginShutdownOptions) (resp azfake.PollerResponder[armlargeinstance.AzureLargeInstanceClientShutdownResponse], errResp azfake.ErrorResponder)

	// BeginStart is the fake for method AzureLargeInstanceClient.BeginStart
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginStart func(ctx context.Context, resourceGroupName string, azureLargeInstanceName string, options *armlargeinstance.AzureLargeInstanceClientBeginStartOptions) (resp azfake.PollerResponder[armlargeinstance.AzureLargeInstanceClientStartResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method AzureLargeInstanceClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, azureLargeInstanceName string, properties armlargeinstance.AzureLargeInstanceTagsUpdate, options *armlargeinstance.AzureLargeInstanceClientUpdateOptions) (resp azfake.Responder[armlargeinstance.AzureLargeInstanceClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewAzureLargeInstanceServerTransport creates a new instance of AzureLargeInstanceServerTransport with the provided implementation.
// The returned AzureLargeInstanceServerTransport instance is connected to an instance of armlargeinstance.AzureLargeInstanceClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAzureLargeInstanceServerTransport(srv *AzureLargeInstanceServer) *AzureLargeInstanceServerTransport {
	return &AzureLargeInstanceServerTransport{
		srv:                         srv,
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armlargeinstance.AzureLargeInstanceClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armlargeinstance.AzureLargeInstanceClientListBySubscriptionResponse]](),
		beginRestart:                newTracker[azfake.PollerResponder[armlargeinstance.AzureLargeInstanceClientRestartResponse]](),
		beginShutdown:               newTracker[azfake.PollerResponder[armlargeinstance.AzureLargeInstanceClientShutdownResponse]](),
		beginStart:                  newTracker[azfake.PollerResponder[armlargeinstance.AzureLargeInstanceClientStartResponse]](),
	}
}

// AzureLargeInstanceServerTransport connects instances of armlargeinstance.AzureLargeInstanceClient to instances of AzureLargeInstanceServer.
// Don't use this type directly, use NewAzureLargeInstanceServerTransport instead.
type AzureLargeInstanceServerTransport struct {
	srv                         *AzureLargeInstanceServer
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armlargeinstance.AzureLargeInstanceClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armlargeinstance.AzureLargeInstanceClientListBySubscriptionResponse]]
	beginRestart                *tracker[azfake.PollerResponder[armlargeinstance.AzureLargeInstanceClientRestartResponse]]
	beginShutdown               *tracker[azfake.PollerResponder[armlargeinstance.AzureLargeInstanceClientShutdownResponse]]
	beginStart                  *tracker[azfake.PollerResponder[armlargeinstance.AzureLargeInstanceClientStartResponse]]
}

// Do implements the policy.Transporter interface for AzureLargeInstanceServerTransport.
func (a *AzureLargeInstanceServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "AzureLargeInstanceClient.Get":
		resp, err = a.dispatchGet(req)
	case "AzureLargeInstanceClient.NewListByResourceGroupPager":
		resp, err = a.dispatchNewListByResourceGroupPager(req)
	case "AzureLargeInstanceClient.NewListBySubscriptionPager":
		resp, err = a.dispatchNewListBySubscriptionPager(req)
	case "AzureLargeInstanceClient.BeginRestart":
		resp, err = a.dispatchBeginRestart(req)
	case "AzureLargeInstanceClient.BeginShutdown":
		resp, err = a.dispatchBeginShutdown(req)
	case "AzureLargeInstanceClient.BeginStart":
		resp, err = a.dispatchBeginStart(req)
	case "AzureLargeInstanceClient.Update":
		resp, err = a.dispatchUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AzureLargeInstanceServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if a.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureLargeInstance/azureLargeInstances/(?P<azureLargeInstanceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	azureLargeInstanceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("azureLargeInstanceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Get(req.Context(), resourceGroupNameParam, azureLargeInstanceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AzureLargeInstance, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AzureLargeInstanceServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := a.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureLargeInstance/azureLargeInstances`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := a.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		a.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armlargeinstance.AzureLargeInstanceClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		a.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (a *AzureLargeInstanceServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := a.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureLargeInstance/azureLargeInstances`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := a.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		a.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armlargeinstance.AzureLargeInstanceClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		a.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (a *AzureLargeInstanceServerTransport) dispatchBeginRestart(req *http.Request) (*http.Response, error) {
	if a.srv.BeginRestart == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginRestart not implemented")}
	}
	beginRestart := a.beginRestart.get(req)
	if beginRestart == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureLargeInstance/azureLargeInstances/(?P<azureLargeInstanceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/restart`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armlargeinstance.ForceState](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		azureLargeInstanceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("azureLargeInstanceName")])
		if err != nil {
			return nil, err
		}
		var options *armlargeinstance.AzureLargeInstanceClientBeginRestartOptions
		if !reflect.ValueOf(body).IsZero() {
			options = &armlargeinstance.AzureLargeInstanceClientBeginRestartOptions{
				ForceParameter: &body,
			}
		}
		respr, errRespr := a.srv.BeginRestart(req.Context(), resourceGroupNameParam, azureLargeInstanceNameParam, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginRestart = &respr
		a.beginRestart.add(req, beginRestart)
	}

	resp, err := server.PollerResponderNext(beginRestart, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		a.beginRestart.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginRestart) {
		a.beginRestart.remove(req)
	}

	return resp, nil
}

func (a *AzureLargeInstanceServerTransport) dispatchBeginShutdown(req *http.Request) (*http.Response, error) {
	if a.srv.BeginShutdown == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginShutdown not implemented")}
	}
	beginShutdown := a.beginShutdown.get(req)
	if beginShutdown == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureLargeInstance/azureLargeInstances/(?P<azureLargeInstanceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/shutdown`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		azureLargeInstanceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("azureLargeInstanceName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := a.srv.BeginShutdown(req.Context(), resourceGroupNameParam, azureLargeInstanceNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginShutdown = &respr
		a.beginShutdown.add(req, beginShutdown)
	}

	resp, err := server.PollerResponderNext(beginShutdown, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		a.beginShutdown.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginShutdown) {
		a.beginShutdown.remove(req)
	}

	return resp, nil
}

func (a *AzureLargeInstanceServerTransport) dispatchBeginStart(req *http.Request) (*http.Response, error) {
	if a.srv.BeginStart == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginStart not implemented")}
	}
	beginStart := a.beginStart.get(req)
	if beginStart == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureLargeInstance/azureLargeInstances/(?P<azureLargeInstanceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/start`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		azureLargeInstanceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("azureLargeInstanceName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := a.srv.BeginStart(req.Context(), resourceGroupNameParam, azureLargeInstanceNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginStart = &respr
		a.beginStart.add(req, beginStart)
	}

	resp, err := server.PollerResponderNext(beginStart, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		a.beginStart.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginStart) {
		a.beginStart.remove(req)
	}

	return resp, nil
}

func (a *AzureLargeInstanceServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if a.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureLargeInstance/azureLargeInstances/(?P<azureLargeInstanceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armlargeinstance.AzureLargeInstanceTagsUpdate](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	azureLargeInstanceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("azureLargeInstanceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Update(req.Context(), resourceGroupNameParam, azureLargeInstanceNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AzureLargeInstance, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
