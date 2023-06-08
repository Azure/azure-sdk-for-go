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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v3"
	"net/http"
	"regexp"
)

// VipSwapServer is a fake server for instances of the armnetwork.VipSwapClient type.
type VipSwapServer struct {
	// BeginCreate is the fake for method VipSwapClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreate func(ctx context.Context, groupName string, resourceName string, parameters armnetwork.SwapResource, options *armnetwork.VipSwapClientBeginCreateOptions) (resp azfake.PollerResponder[armnetwork.VipSwapClientCreateResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method VipSwapClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, groupName string, resourceName string, options *armnetwork.VipSwapClientGetOptions) (resp azfake.Responder[armnetwork.VipSwapClientGetResponse], errResp azfake.ErrorResponder)

	// List is the fake for method VipSwapClient.List
	// HTTP status codes to indicate success: http.StatusOK
	List func(ctx context.Context, groupName string, resourceName string, options *armnetwork.VipSwapClientListOptions) (resp azfake.Responder[armnetwork.VipSwapClientListResponse], errResp azfake.ErrorResponder)
}

// NewVipSwapServerTransport creates a new instance of VipSwapServerTransport with the provided implementation.
// The returned VipSwapServerTransport instance is connected to an instance of armnetwork.VipSwapClient by way of the
// undefined.Transporter field.
func NewVipSwapServerTransport(srv *VipSwapServer) *VipSwapServerTransport {
	return &VipSwapServerTransport{srv: srv}
}

// VipSwapServerTransport connects instances of armnetwork.VipSwapClient to instances of VipSwapServer.
// Don't use this type directly, use NewVipSwapServerTransport instead.
type VipSwapServerTransport struct {
	srv         *VipSwapServer
	beginCreate *azfake.PollerResponder[armnetwork.VipSwapClientCreateResponse]
}

// Do implements the policy.Transporter interface for VipSwapServerTransport.
func (v *VipSwapServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "VipSwapClient.BeginCreate":
		resp, err = v.dispatchBeginCreate(req)
	case "VipSwapClient.Get":
		resp, err = v.dispatchGet(req)
	case "VipSwapClient.List":
		resp, err = v.dispatchList(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *VipSwapServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if v.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("method BeginCreate not implemented")}
	}
	if v.beginCreate == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<groupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Compute/cloudServices/(?P<resourceName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/cloudServiceSlots/(?P<singletonResource>[a-zA-Z0-9-_]+)"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.SwapResource](req)
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginCreate(req.Context(), matches[regex.SubexpIndex("groupName")], matches[regex.SubexpIndex("resourceName")], body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		v.beginCreate = &respr
	}

	resp, err := server.PollerResponderNext(v.beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(v.beginCreate) {
		v.beginCreate = nil
	}

	return resp, nil
}

func (v *VipSwapServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if v.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("method Get not implemented")}
	}
	const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<groupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Compute/cloudServices/(?P<resourceName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/cloudServiceSlots/(?P<singletonResource>[a-zA-Z0-9-_]+)"
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.Path)
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	respr, errRespr := v.srv.Get(req.Context(), matches[regex.SubexpIndex("groupName")], matches[regex.SubexpIndex("resourceName")], nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SwapResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (v *VipSwapServerTransport) dispatchList(req *http.Request) (*http.Response, error) {
	if v.srv.List == nil {
		return nil, &nonRetriableError{errors.New("method List not implemented")}
	}
	const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<groupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Compute/cloudServices/(?P<resourceName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/cloudServiceSlots"
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.Path)
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	respr, errRespr := v.srv.List(req.Context(), matches[regex.SubexpIndex("groupName")], matches[regex.SubexpIndex("resourceName")], nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SwapResourceListResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
