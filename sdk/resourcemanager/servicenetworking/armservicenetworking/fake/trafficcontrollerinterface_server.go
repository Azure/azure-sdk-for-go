// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package fake

import (
	"context"
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicenetworking/armservicenetworking/v2"
	"net/http"
	"net/url"
	"regexp"
)

// TrafficControllerInterfaceServer is a fake server for instances of the armservicenetworking.TrafficControllerInterfaceClient type.
type TrafficControllerInterfaceServer struct {
	// BeginCreateOrUpdate is the fake for method TrafficControllerInterfaceClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, trafficControllerName string, resource armservicenetworking.TrafficController, options *armservicenetworking.TrafficControllerInterfaceClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armservicenetworking.TrafficControllerInterfaceClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method TrafficControllerInterfaceClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, trafficControllerName string, options *armservicenetworking.TrafficControllerInterfaceClientBeginDeleteOptions) (resp azfake.PollerResponder[armservicenetworking.TrafficControllerInterfaceClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method TrafficControllerInterfaceClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, trafficControllerName string, options *armservicenetworking.TrafficControllerInterfaceClientGetOptions) (resp azfake.Responder[armservicenetworking.TrafficControllerInterfaceClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method TrafficControllerInterfaceClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armservicenetworking.TrafficControllerInterfaceClientListByResourceGroupOptions) (resp azfake.PagerResponder[armservicenetworking.TrafficControllerInterfaceClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method TrafficControllerInterfaceClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armservicenetworking.TrafficControllerInterfaceClientListBySubscriptionOptions) (resp azfake.PagerResponder[armservicenetworking.TrafficControllerInterfaceClientListBySubscriptionResponse])

	// Update is the fake for method TrafficControllerInterfaceClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, trafficControllerName string, properties armservicenetworking.TrafficControllerUpdate, options *armservicenetworking.TrafficControllerInterfaceClientUpdateOptions) (resp azfake.Responder[armservicenetworking.TrafficControllerInterfaceClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewTrafficControllerInterfaceServerTransport creates a new instance of TrafficControllerInterfaceServerTransport with the provided implementation.
// The returned TrafficControllerInterfaceServerTransport instance is connected to an instance of armservicenetworking.TrafficControllerInterfaceClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewTrafficControllerInterfaceServerTransport(srv *TrafficControllerInterfaceServer) *TrafficControllerInterfaceServerTransport {
	return &TrafficControllerInterfaceServerTransport{
		srv:                         srv,
		beginCreateOrUpdate:         newTracker[azfake.PollerResponder[armservicenetworking.TrafficControllerInterfaceClientCreateOrUpdateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armservicenetworking.TrafficControllerInterfaceClientDeleteResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armservicenetworking.TrafficControllerInterfaceClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armservicenetworking.TrafficControllerInterfaceClientListBySubscriptionResponse]](),
	}
}

// TrafficControllerInterfaceServerTransport connects instances of armservicenetworking.TrafficControllerInterfaceClient to instances of TrafficControllerInterfaceServer.
// Don't use this type directly, use NewTrafficControllerInterfaceServerTransport instead.
type TrafficControllerInterfaceServerTransport struct {
	srv                         *TrafficControllerInterfaceServer
	beginCreateOrUpdate         *tracker[azfake.PollerResponder[armservicenetworking.TrafficControllerInterfaceClientCreateOrUpdateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armservicenetworking.TrafficControllerInterfaceClientDeleteResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armservicenetworking.TrafficControllerInterfaceClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armservicenetworking.TrafficControllerInterfaceClientListBySubscriptionResponse]]
}

// Do implements the policy.Transporter interface for TrafficControllerInterfaceServerTransport.
func (t *TrafficControllerInterfaceServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return t.dispatchToMethodFake(req, method)
}

func (t *TrafficControllerInterfaceServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if trafficControllerInterfaceServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = trafficControllerInterfaceServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "TrafficControllerInterfaceClient.BeginCreateOrUpdate":
				res.resp, res.err = t.dispatchBeginCreateOrUpdate(req)
			case "TrafficControllerInterfaceClient.BeginDelete":
				res.resp, res.err = t.dispatchBeginDelete(req)
			case "TrafficControllerInterfaceClient.Get":
				res.resp, res.err = t.dispatchGet(req)
			case "TrafficControllerInterfaceClient.NewListByResourceGroupPager":
				res.resp, res.err = t.dispatchNewListByResourceGroupPager(req)
			case "TrafficControllerInterfaceClient.NewListBySubscriptionPager":
				res.resp, res.err = t.dispatchNewListBySubscriptionPager(req)
			case "TrafficControllerInterfaceClient.Update":
				res.resp, res.err = t.dispatchUpdate(req)
			default:
				res.err = fmt.Errorf("unhandled API %s", method)
			}

		}
		select {
		case resultChan <- res:
		case <-req.Context().Done():
		}
	}()

	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

func (t *TrafficControllerInterfaceServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if t.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := t.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers/(?P<trafficControllerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armservicenetworking.TrafficController](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		trafficControllerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("trafficControllerName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := t.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, trafficControllerNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		t.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		t.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		t.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (t *TrafficControllerInterfaceServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if t.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := t.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers/(?P<trafficControllerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		trafficControllerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("trafficControllerName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := t.srv.BeginDelete(req.Context(), resourceGroupNameParam, trafficControllerNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		t.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		t.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		t.beginDelete.remove(req)
	}

	return resp, nil
}

func (t *TrafficControllerInterfaceServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if t.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers/(?P<trafficControllerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	trafficControllerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("trafficControllerName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := t.srv.Get(req.Context(), resourceGroupNameParam, trafficControllerNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).TrafficController, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *TrafficControllerInterfaceServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if t.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := t.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := t.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		t.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armservicenetworking.TrafficControllerInterfaceClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		t.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		t.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (t *TrafficControllerInterfaceServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if t.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := t.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := t.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		t.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armservicenetworking.TrafficControllerInterfaceClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		t.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		t.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (t *TrafficControllerInterfaceServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if t.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers/(?P<trafficControllerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armservicenetworking.TrafficControllerUpdate](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	trafficControllerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("trafficControllerName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := t.srv.Update(req.Context(), resourceGroupNameParam, trafficControllerNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).TrafficController, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to TrafficControllerInterfaceServerTransport
var trafficControllerInterfaceServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
