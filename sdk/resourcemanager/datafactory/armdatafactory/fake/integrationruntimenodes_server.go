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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v5"
	"net/http"
	"net/url"
	"regexp"
)

// IntegrationRuntimeNodesServer is a fake server for instances of the armdatafactory.IntegrationRuntimeNodesClient type.
type IntegrationRuntimeNodesServer struct {
	// Delete is the fake for method IntegrationRuntimeNodesClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, nodeName string, options *armdatafactory.IntegrationRuntimeNodesClientDeleteOptions) (resp azfake.Responder[armdatafactory.IntegrationRuntimeNodesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method IntegrationRuntimeNodesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, nodeName string, options *armdatafactory.IntegrationRuntimeNodesClientGetOptions) (resp azfake.Responder[armdatafactory.IntegrationRuntimeNodesClientGetResponse], errResp azfake.ErrorResponder)

	// GetIPAddress is the fake for method IntegrationRuntimeNodesClient.GetIPAddress
	// HTTP status codes to indicate success: http.StatusOK
	GetIPAddress func(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, nodeName string, options *armdatafactory.IntegrationRuntimeNodesClientGetIPAddressOptions) (resp azfake.Responder[armdatafactory.IntegrationRuntimeNodesClientGetIPAddressResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method IntegrationRuntimeNodesClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, factoryName string, integrationRuntimeName string, nodeName string, updateIntegrationRuntimeNodeRequest armdatafactory.UpdateIntegrationRuntimeNodeRequest, options *armdatafactory.IntegrationRuntimeNodesClientUpdateOptions) (resp azfake.Responder[armdatafactory.IntegrationRuntimeNodesClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewIntegrationRuntimeNodesServerTransport creates a new instance of IntegrationRuntimeNodesServerTransport with the provided implementation.
// The returned IntegrationRuntimeNodesServerTransport instance is connected to an instance of armdatafactory.IntegrationRuntimeNodesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewIntegrationRuntimeNodesServerTransport(srv *IntegrationRuntimeNodesServer) *IntegrationRuntimeNodesServerTransport {
	return &IntegrationRuntimeNodesServerTransport{srv: srv}
}

// IntegrationRuntimeNodesServerTransport connects instances of armdatafactory.IntegrationRuntimeNodesClient to instances of IntegrationRuntimeNodesServer.
// Don't use this type directly, use NewIntegrationRuntimeNodesServerTransport instead.
type IntegrationRuntimeNodesServerTransport struct {
	srv *IntegrationRuntimeNodesServer
}

// Do implements the policy.Transporter interface for IntegrationRuntimeNodesServerTransport.
func (i *IntegrationRuntimeNodesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "IntegrationRuntimeNodesClient.Delete":
		resp, err = i.dispatchDelete(req)
	case "IntegrationRuntimeNodesClient.Get":
		resp, err = i.dispatchGet(req)
	case "IntegrationRuntimeNodesClient.GetIPAddress":
		resp, err = i.dispatchGetIPAddress(req)
	case "IntegrationRuntimeNodesClient.Update":
		resp, err = i.dispatchUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (i *IntegrationRuntimeNodesServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if i.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/integrationRuntimes/(?P<integrationRuntimeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/nodes/(?P<nodeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	factoryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("factoryName")])
	if err != nil {
		return nil, err
	}
	integrationRuntimeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("integrationRuntimeName")])
	if err != nil {
		return nil, err
	}
	nodeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("nodeName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Delete(req.Context(), resourceGroupNameParam, factoryNameParam, integrationRuntimeNameParam, nodeNameParam, nil)
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

func (i *IntegrationRuntimeNodesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if i.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/integrationRuntimes/(?P<integrationRuntimeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/nodes/(?P<nodeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	factoryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("factoryName")])
	if err != nil {
		return nil, err
	}
	integrationRuntimeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("integrationRuntimeName")])
	if err != nil {
		return nil, err
	}
	nodeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("nodeName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Get(req.Context(), resourceGroupNameParam, factoryNameParam, integrationRuntimeNameParam, nodeNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SelfHostedIntegrationRuntimeNode, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *IntegrationRuntimeNodesServerTransport) dispatchGetIPAddress(req *http.Request) (*http.Response, error) {
	if i.srv.GetIPAddress == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetIPAddress not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/integrationRuntimes/(?P<integrationRuntimeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/nodes/(?P<nodeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/ipAddress`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	factoryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("factoryName")])
	if err != nil {
		return nil, err
	}
	integrationRuntimeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("integrationRuntimeName")])
	if err != nil {
		return nil, err
	}
	nodeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("nodeName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.GetIPAddress(req.Context(), resourceGroupNameParam, factoryNameParam, integrationRuntimeNameParam, nodeNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).IntegrationRuntimeNodeIPAddress, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *IntegrationRuntimeNodesServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if i.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/integrationRuntimes/(?P<integrationRuntimeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/nodes/(?P<nodeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdatafactory.UpdateIntegrationRuntimeNodeRequest](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	factoryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("factoryName")])
	if err != nil {
		return nil, err
	}
	integrationRuntimeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("integrationRuntimeName")])
	if err != nil {
		return nil, err
	}
	nodeNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("nodeName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Update(req.Context(), resourceGroupNameParam, factoryNameParam, integrationRuntimeNameParam, nodeNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SelfHostedIntegrationRuntimeNode, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
