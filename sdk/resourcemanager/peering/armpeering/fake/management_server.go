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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/peering/armpeering"
	"net/http"
	"regexp"
)

// ManagementServer is a fake server for instances of the armpeering.ManagementClient type.
type ManagementServer struct {
	// CheckServiceProviderAvailability is the fake for method ManagementClient.CheckServiceProviderAvailability
	// HTTP status codes to indicate success: http.StatusOK
	CheckServiceProviderAvailability func(ctx context.Context, checkServiceProviderAvailabilityInput armpeering.CheckServiceProviderAvailabilityInput, options *armpeering.ManagementClientCheckServiceProviderAvailabilityOptions) (resp azfake.Responder[armpeering.ManagementClientCheckServiceProviderAvailabilityResponse], errResp azfake.ErrorResponder)
}

// NewManagementServerTransport creates a new instance of ManagementServerTransport with the provided implementation.
// The returned ManagementServerTransport instance is connected to an instance of armpeering.ManagementClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewManagementServerTransport(srv *ManagementServer) *ManagementServerTransport {
	return &ManagementServerTransport{srv: srv}
}

// ManagementServerTransport connects instances of armpeering.ManagementClient to instances of ManagementServer.
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
	case "ManagementClient.CheckServiceProviderAvailability":
		resp, err = m.dispatchCheckServiceProviderAvailability(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *ManagementServerTransport) dispatchCheckServiceProviderAvailability(req *http.Request) (*http.Response, error) {
	if m.srv.CheckServiceProviderAvailability == nil {
		return nil, &nonRetriableError{errors.New("fake for method CheckServiceProviderAvailability not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Peering/checkServiceProviderAvailability`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armpeering.CheckServiceProviderAvailabilityInput](req)
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.CheckServiceProviderAvailability(req.Context(), body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Value, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
