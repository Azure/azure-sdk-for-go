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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/automanage/armautomanage"
	"net/http"
	"regexp"
)

// ServicePrincipalsServer is a fake server for instances of the armautomanage.ServicePrincipalsClient type.
type ServicePrincipalsServer struct {
	// Get is the fake for method ServicePrincipalsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, options *armautomanage.ServicePrincipalsClientGetOptions) (resp azfake.Responder[armautomanage.ServicePrincipalsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListBySubscriptionPager is the fake for method ServicePrincipalsClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armautomanage.ServicePrincipalsClientListBySubscriptionOptions) (resp azfake.PagerResponder[armautomanage.ServicePrincipalsClientListBySubscriptionResponse])
}

// NewServicePrincipalsServerTransport creates a new instance of ServicePrincipalsServerTransport with the provided implementation.
// The returned ServicePrincipalsServerTransport instance is connected to an instance of armautomanage.ServicePrincipalsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewServicePrincipalsServerTransport(srv *ServicePrincipalsServer) *ServicePrincipalsServerTransport {
	return &ServicePrincipalsServerTransport{
		srv:                        srv,
		newListBySubscriptionPager: newTracker[azfake.PagerResponder[armautomanage.ServicePrincipalsClientListBySubscriptionResponse]](),
	}
}

// ServicePrincipalsServerTransport connects instances of armautomanage.ServicePrincipalsClient to instances of ServicePrincipalsServer.
// Don't use this type directly, use NewServicePrincipalsServerTransport instead.
type ServicePrincipalsServerTransport struct {
	srv                        *ServicePrincipalsServer
	newListBySubscriptionPager *tracker[azfake.PagerResponder[armautomanage.ServicePrincipalsClientListBySubscriptionResponse]]
}

// Do implements the policy.Transporter interface for ServicePrincipalsServerTransport.
func (s *ServicePrincipalsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ServicePrincipalsClient.Get":
		resp, err = s.dispatchGet(req)
	case "ServicePrincipalsClient.NewListBySubscriptionPager":
		resp, err = s.dispatchNewListBySubscriptionPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ServicePrincipalsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automanage/servicePrincipals/default`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	respr, errRespr := s.srv.Get(req.Context(), nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ServicePrincipal, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ServicePrincipalsServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := s.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Automanage/servicePrincipals`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := s.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		s.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		s.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}
