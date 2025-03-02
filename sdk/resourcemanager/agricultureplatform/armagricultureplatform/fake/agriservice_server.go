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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/agricultureplatform/armagricultureplatform"
	"net/http"
	"net/url"
	"regexp"
)

// AgriServiceServer is a fake server for instances of the armagricultureplatform.AgriServiceClient type.
type AgriServiceServer struct {
	// BeginCreateOrUpdate is the fake for method AgriServiceClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, agriServiceResourceName string, resource armagricultureplatform.AgriServiceResource, options *armagricultureplatform.AgriServiceClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armagricultureplatform.AgriServiceClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method AgriServiceClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, agriServiceResourceName string, options *armagricultureplatform.AgriServiceClientBeginDeleteOptions) (resp azfake.PollerResponder[armagricultureplatform.AgriServiceClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method AgriServiceClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, agriServiceResourceName string, options *armagricultureplatform.AgriServiceClientGetOptions) (resp azfake.Responder[armagricultureplatform.AgriServiceClientGetResponse], errResp azfake.ErrorResponder)

	// ListAvailableSolutions is the fake for method AgriServiceClient.ListAvailableSolutions
	// HTTP status codes to indicate success: http.StatusOK
	ListAvailableSolutions func(ctx context.Context, resourceGroupName string, agriServiceResourceName string, options *armagricultureplatform.AgriServiceClientListAvailableSolutionsOptions) (resp azfake.Responder[armagricultureplatform.AgriServiceClientListAvailableSolutionsResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method AgriServiceClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armagricultureplatform.AgriServiceClientListByResourceGroupOptions) (resp azfake.PagerResponder[armagricultureplatform.AgriServiceClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method AgriServiceClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armagricultureplatform.AgriServiceClientListBySubscriptionOptions) (resp azfake.PagerResponder[armagricultureplatform.AgriServiceClientListBySubscriptionResponse])

	// BeginUpdate is the fake for method AgriServiceClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, agriServiceResourceName string, properties armagricultureplatform.AgriServiceResourceUpdate, options *armagricultureplatform.AgriServiceClientBeginUpdateOptions) (resp azfake.PollerResponder[armagricultureplatform.AgriServiceClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewAgriServiceServerTransport creates a new instance of AgriServiceServerTransport with the provided implementation.
// The returned AgriServiceServerTransport instance is connected to an instance of armagricultureplatform.AgriServiceClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAgriServiceServerTransport(srv *AgriServiceServer) *AgriServiceServerTransport {
	return &AgriServiceServerTransport{
		srv:                         srv,
		beginCreateOrUpdate:         newTracker[azfake.PollerResponder[armagricultureplatform.AgriServiceClientCreateOrUpdateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armagricultureplatform.AgriServiceClientDeleteResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armagricultureplatform.AgriServiceClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armagricultureplatform.AgriServiceClientListBySubscriptionResponse]](),
		beginUpdate:                 newTracker[azfake.PollerResponder[armagricultureplatform.AgriServiceClientUpdateResponse]](),
	}
}

// AgriServiceServerTransport connects instances of armagricultureplatform.AgriServiceClient to instances of AgriServiceServer.
// Don't use this type directly, use NewAgriServiceServerTransport instead.
type AgriServiceServerTransport struct {
	srv                         *AgriServiceServer
	beginCreateOrUpdate         *tracker[azfake.PollerResponder[armagricultureplatform.AgriServiceClientCreateOrUpdateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armagricultureplatform.AgriServiceClientDeleteResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armagricultureplatform.AgriServiceClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armagricultureplatform.AgriServiceClientListBySubscriptionResponse]]
	beginUpdate                 *tracker[azfake.PollerResponder[armagricultureplatform.AgriServiceClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for AgriServiceServerTransport.
func (a *AgriServiceServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return a.dispatchToMethodFake(req, method)
}

func (a *AgriServiceServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if agriServiceServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = agriServiceServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "AgriServiceClient.BeginCreateOrUpdate":
				res.resp, res.err = a.dispatchBeginCreateOrUpdate(req)
			case "AgriServiceClient.BeginDelete":
				res.resp, res.err = a.dispatchBeginDelete(req)
			case "AgriServiceClient.Get":
				res.resp, res.err = a.dispatchGet(req)
			case "AgriServiceClient.ListAvailableSolutions":
				res.resp, res.err = a.dispatchListAvailableSolutions(req)
			case "AgriServiceClient.NewListByResourceGroupPager":
				res.resp, res.err = a.dispatchNewListByResourceGroupPager(req)
			case "AgriServiceClient.NewListBySubscriptionPager":
				res.resp, res.err = a.dispatchNewListBySubscriptionPager(req)
			case "AgriServiceClient.BeginUpdate":
				res.resp, res.err = a.dispatchBeginUpdate(req)
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

func (a *AgriServiceServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if a.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := a.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AgriculturePlatform/agriServices/(?P<agriServiceResourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armagricultureplatform.AgriServiceResource](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		agriServiceResourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("agriServiceResourceName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := a.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, agriServiceResourceNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		a.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		a.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		a.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (a *AgriServiceServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if a.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := a.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AgriculturePlatform/agriServices/(?P<agriServiceResourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		agriServiceResourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("agriServiceResourceName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := a.srv.BeginDelete(req.Context(), resourceGroupNameParam, agriServiceResourceNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		a.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		a.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		a.beginDelete.remove(req)
	}

	return resp, nil
}

func (a *AgriServiceServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if a.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AgriculturePlatform/agriServices/(?P<agriServiceResourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	agriServiceResourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("agriServiceResourceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Get(req.Context(), resourceGroupNameParam, agriServiceResourceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AgriServiceResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AgriServiceServerTransport) dispatchListAvailableSolutions(req *http.Request) (*http.Response, error) {
	if a.srv.ListAvailableSolutions == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListAvailableSolutions not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AgriculturePlatform/agriServices/(?P<agriServiceResourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listAvailableSolutions`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	agriServiceResourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("agriServiceResourceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.ListAvailableSolutions(req.Context(), resourceGroupNameParam, agriServiceResourceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AvailableAgriSolutionListResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (a *AgriServiceServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := a.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AgriculturePlatform/agriServices`
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
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armagricultureplatform.AgriServiceClientListByResourceGroupResponse, createLink func() string) {
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

func (a *AgriServiceServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := a.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AgriculturePlatform/agriServices`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := a.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		a.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armagricultureplatform.AgriServiceClientListBySubscriptionResponse, createLink func() string) {
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

func (a *AgriServiceServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if a.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := a.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AgriculturePlatform/agriServices/(?P<agriServiceResourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armagricultureplatform.AgriServiceResourceUpdate](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		agriServiceResourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("agriServiceResourceName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := a.srv.BeginUpdate(req.Context(), resourceGroupNameParam, agriServiceResourceNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		a.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		a.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		a.beginUpdate.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to AgriServiceServerTransport
var agriServiceServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
