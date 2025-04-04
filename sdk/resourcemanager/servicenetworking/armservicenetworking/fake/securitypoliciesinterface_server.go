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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/servicenetworking/armservicenetworking"
	"net/http"
	"net/url"
	"regexp"
)

// SecurityPoliciesInterfaceServer is a fake server for instances of the armservicenetworking.SecurityPoliciesInterfaceClient type.
type SecurityPoliciesInterfaceServer struct {
	// BeginCreateOrUpdate is the fake for method SecurityPoliciesInterfaceClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, trafficControllerName string, securityPolicyName string, resource armservicenetworking.SecurityPolicy, options *armservicenetworking.SecurityPoliciesInterfaceClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armservicenetworking.SecurityPoliciesInterfaceClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method SecurityPoliciesInterfaceClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, trafficControllerName string, securityPolicyName string, options *armservicenetworking.SecurityPoliciesInterfaceClientBeginDeleteOptions) (resp azfake.PollerResponder[armservicenetworking.SecurityPoliciesInterfaceClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method SecurityPoliciesInterfaceClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, trafficControllerName string, securityPolicyName string, options *armservicenetworking.SecurityPoliciesInterfaceClientGetOptions) (resp azfake.Responder[armservicenetworking.SecurityPoliciesInterfaceClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByTrafficControllerPager is the fake for method SecurityPoliciesInterfaceClient.NewListByTrafficControllerPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByTrafficControllerPager func(resourceGroupName string, trafficControllerName string, options *armservicenetworking.SecurityPoliciesInterfaceClientListByTrafficControllerOptions) (resp azfake.PagerResponder[armservicenetworking.SecurityPoliciesInterfaceClientListByTrafficControllerResponse])

	// Update is the fake for method SecurityPoliciesInterfaceClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, trafficControllerName string, securityPolicyName string, properties armservicenetworking.SecurityPolicyUpdate, options *armservicenetworking.SecurityPoliciesInterfaceClientUpdateOptions) (resp azfake.Responder[armservicenetworking.SecurityPoliciesInterfaceClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewSecurityPoliciesInterfaceServerTransport creates a new instance of SecurityPoliciesInterfaceServerTransport with the provided implementation.
// The returned SecurityPoliciesInterfaceServerTransport instance is connected to an instance of armservicenetworking.SecurityPoliciesInterfaceClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSecurityPoliciesInterfaceServerTransport(srv *SecurityPoliciesInterfaceServer) *SecurityPoliciesInterfaceServerTransport {
	return &SecurityPoliciesInterfaceServerTransport{
		srv:                             srv,
		beginCreateOrUpdate:             newTracker[azfake.PollerResponder[armservicenetworking.SecurityPoliciesInterfaceClientCreateOrUpdateResponse]](),
		beginDelete:                     newTracker[azfake.PollerResponder[armservicenetworking.SecurityPoliciesInterfaceClientDeleteResponse]](),
		newListByTrafficControllerPager: newTracker[azfake.PagerResponder[armservicenetworking.SecurityPoliciesInterfaceClientListByTrafficControllerResponse]](),
	}
}

// SecurityPoliciesInterfaceServerTransport connects instances of armservicenetworking.SecurityPoliciesInterfaceClient to instances of SecurityPoliciesInterfaceServer.
// Don't use this type directly, use NewSecurityPoliciesInterfaceServerTransport instead.
type SecurityPoliciesInterfaceServerTransport struct {
	srv                             *SecurityPoliciesInterfaceServer
	beginCreateOrUpdate             *tracker[azfake.PollerResponder[armservicenetworking.SecurityPoliciesInterfaceClientCreateOrUpdateResponse]]
	beginDelete                     *tracker[azfake.PollerResponder[armservicenetworking.SecurityPoliciesInterfaceClientDeleteResponse]]
	newListByTrafficControllerPager *tracker[azfake.PagerResponder[armservicenetworking.SecurityPoliciesInterfaceClientListByTrafficControllerResponse]]
}

// Do implements the policy.Transporter interface for SecurityPoliciesInterfaceServerTransport.
func (s *SecurityPoliciesInterfaceServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return s.dispatchToMethodFake(req, method)
}

func (s *SecurityPoliciesInterfaceServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if securityPoliciesInterfaceServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = securityPoliciesInterfaceServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "SecurityPoliciesInterfaceClient.BeginCreateOrUpdate":
				res.resp, res.err = s.dispatchBeginCreateOrUpdate(req)
			case "SecurityPoliciesInterfaceClient.BeginDelete":
				res.resp, res.err = s.dispatchBeginDelete(req)
			case "SecurityPoliciesInterfaceClient.Get":
				res.resp, res.err = s.dispatchGet(req)
			case "SecurityPoliciesInterfaceClient.NewListByTrafficControllerPager":
				res.resp, res.err = s.dispatchNewListByTrafficControllerPager(req)
			case "SecurityPoliciesInterfaceClient.Update":
				res.resp, res.err = s.dispatchUpdate(req)
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

func (s *SecurityPoliciesInterfaceServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := s.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers/(?P<trafficControllerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/securityPolicies/(?P<securityPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armservicenetworking.SecurityPolicy](req)
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
		securityPolicyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("securityPolicyName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, trafficControllerNameParam, securityPolicyNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		s.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		s.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		s.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (s *SecurityPoliciesInterfaceServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if s.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := s.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers/(?P<trafficControllerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/securityPolicies/(?P<securityPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
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
		securityPolicyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("securityPolicyName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginDelete(req.Context(), resourceGroupNameParam, trafficControllerNameParam, securityPolicyNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		s.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		s.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		s.beginDelete.remove(req)
	}

	return resp, nil
}

func (s *SecurityPoliciesInterfaceServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers/(?P<trafficControllerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/securityPolicies/(?P<securityPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	securityPolicyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("securityPolicyName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), resourceGroupNameParam, trafficControllerNameParam, securityPolicyNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SecurityPolicy, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SecurityPoliciesInterfaceServerTransport) dispatchNewListByTrafficControllerPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByTrafficControllerPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByTrafficControllerPager not implemented")}
	}
	newListByTrafficControllerPager := s.newListByTrafficControllerPager.get(req)
	if newListByTrafficControllerPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers/(?P<trafficControllerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/securityPolicies`
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
		resp := s.srv.NewListByTrafficControllerPager(resourceGroupNameParam, trafficControllerNameParam, nil)
		newListByTrafficControllerPager = &resp
		s.newListByTrafficControllerPager.add(req, newListByTrafficControllerPager)
		server.PagerResponderInjectNextLinks(newListByTrafficControllerPager, req, func(page *armservicenetworking.SecurityPoliciesInterfaceClientListByTrafficControllerResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByTrafficControllerPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListByTrafficControllerPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByTrafficControllerPager) {
		s.newListByTrafficControllerPager.remove(req)
	}
	return resp, nil
}

func (s *SecurityPoliciesInterfaceServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ServiceNetworking/trafficControllers/(?P<trafficControllerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/securityPolicies/(?P<securityPolicyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armservicenetworking.SecurityPolicyUpdate](req)
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
	securityPolicyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("securityPolicyName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Update(req.Context(), resourceGroupNameParam, trafficControllerNameParam, securityPolicyNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SecurityPolicy, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to SecurityPoliciesInterfaceServerTransport
var securityPoliciesInterfaceServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
