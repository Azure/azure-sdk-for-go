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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cdn/armcdn/v2"
	"net/http"
	"net/url"
	"regexp"
)

// ManagementServer is a fake server for instances of the armcdn.ManagementClient type.
type ManagementServer struct {
	// CheckEndpointNameAvailability is the fake for method ManagementClient.CheckEndpointNameAvailability
	// HTTP status codes to indicate success: http.StatusOK
	CheckEndpointNameAvailability func(ctx context.Context, resourceGroupName string, checkEndpointNameAvailabilityInput armcdn.CheckEndpointNameAvailabilityInput, options *armcdn.ManagementClientCheckEndpointNameAvailabilityOptions) (resp azfake.Responder[armcdn.ManagementClientCheckEndpointNameAvailabilityResponse], errResp azfake.ErrorResponder)

	// CheckNameAvailability is the fake for method ManagementClient.CheckNameAvailability
	// HTTP status codes to indicate success: http.StatusOK
	CheckNameAvailability func(ctx context.Context, checkNameAvailabilityInput armcdn.CheckNameAvailabilityInput, options *armcdn.ManagementClientCheckNameAvailabilityOptions) (resp azfake.Responder[armcdn.ManagementClientCheckNameAvailabilityResponse], errResp azfake.ErrorResponder)

	// CheckNameAvailabilityWithSubscription is the fake for method ManagementClient.CheckNameAvailabilityWithSubscription
	// HTTP status codes to indicate success: http.StatusOK
	CheckNameAvailabilityWithSubscription func(ctx context.Context, checkNameAvailabilityInput armcdn.CheckNameAvailabilityInput, options *armcdn.ManagementClientCheckNameAvailabilityWithSubscriptionOptions) (resp azfake.Responder[armcdn.ManagementClientCheckNameAvailabilityWithSubscriptionResponse], errResp azfake.ErrorResponder)

	// ValidateProbe is the fake for method ManagementClient.ValidateProbe
	// HTTP status codes to indicate success: http.StatusOK
	ValidateProbe func(ctx context.Context, validateProbeInput armcdn.ValidateProbeInput, options *armcdn.ManagementClientValidateProbeOptions) (resp azfake.Responder[armcdn.ManagementClientValidateProbeResponse], errResp azfake.ErrorResponder)
}

// NewManagementServerTransport creates a new instance of ManagementServerTransport with the provided implementation.
// The returned ManagementServerTransport instance is connected to an instance of armcdn.ManagementClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewManagementServerTransport(srv *ManagementServer) *ManagementServerTransport {
	return &ManagementServerTransport{srv: srv}
}

// ManagementServerTransport connects instances of armcdn.ManagementClient to instances of ManagementServer.
// Don't use this type directly, use NewManagementServerTransport instead.
type ManagementServerTransport struct {
	srv *ManagementServer
}

// Do implements the policy.Transporter interface for ManagementServerTransport.
func (m *ManagementServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ManagementClient.CheckEndpointNameAvailability":
		resp, err = m.dispatchCheckEndpointNameAvailability(req)
	case "ManagementClient.CheckNameAvailability":
		resp, err = m.dispatchCheckNameAvailability(req)
	case "ManagementClient.CheckNameAvailabilityWithSubscription":
		resp, err = m.dispatchCheckNameAvailabilityWithSubscription(req)
	case "ManagementClient.ValidateProbe":
		resp, err = m.dispatchValidateProbe(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *ManagementServerTransport) dispatchCheckEndpointNameAvailability(req *http.Request) (*http.Response, error) {
	if m.srv.CheckEndpointNameAvailability == nil {
		return nil, &nonRetriableError{errors.New("fake for method CheckEndpointNameAvailability not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Cdn/checkEndpointNameAvailability`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armcdn.CheckEndpointNameAvailabilityInput](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.CheckEndpointNameAvailability(req.Context(), resourceGroupNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CheckEndpointNameAvailabilityOutput, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *ManagementServerTransport) dispatchCheckNameAvailability(req *http.Request) (*http.Response, error) {
	if m.srv.CheckNameAvailability == nil {
		return nil, &nonRetriableError{errors.New("fake for method CheckNameAvailability not implemented")}
	}
	body, err := server.UnmarshalRequestAsJSON[armcdn.CheckNameAvailabilityInput](req)
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.CheckNameAvailability(req.Context(), body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CheckNameAvailabilityOutput, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *ManagementServerTransport) dispatchCheckNameAvailabilityWithSubscription(req *http.Request) (*http.Response, error) {
	if m.srv.CheckNameAvailabilityWithSubscription == nil {
		return nil, &nonRetriableError{errors.New("fake for method CheckNameAvailabilityWithSubscription not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Cdn/checkNameAvailability`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armcdn.CheckNameAvailabilityInput](req)
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.CheckNameAvailabilityWithSubscription(req.Context(), body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CheckNameAvailabilityOutput, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *ManagementServerTransport) dispatchValidateProbe(req *http.Request) (*http.Response, error) {
	if m.srv.ValidateProbe == nil {
		return nil, &nonRetriableError{errors.New("fake for method ValidateProbe not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Cdn/validateProbe`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armcdn.ValidateProbeInput](req)
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.ValidateProbe(req.Context(), body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ValidateProbeOutput, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
