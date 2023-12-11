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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hybridconnectivity/armhybridconnectivity"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
)

// EndpointsServer is a fake server for instances of the armhybridconnectivity.EndpointsClient type.
type EndpointsServer struct {
	// CreateOrUpdate is the fake for method EndpointsClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK
	CreateOrUpdate func(ctx context.Context, resourceURI string, endpointName string, endpointResource armhybridconnectivity.EndpointResource, options *armhybridconnectivity.EndpointsClientCreateOrUpdateOptions) (resp azfake.Responder[armhybridconnectivity.EndpointsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method EndpointsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceURI string, endpointName string, options *armhybridconnectivity.EndpointsClientDeleteOptions) (resp azfake.Responder[armhybridconnectivity.EndpointsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method EndpointsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceURI string, endpointName string, options *armhybridconnectivity.EndpointsClientGetOptions) (resp azfake.Responder[armhybridconnectivity.EndpointsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method EndpointsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceURI string, options *armhybridconnectivity.EndpointsClientListOptions) (resp azfake.PagerResponder[armhybridconnectivity.EndpointsClientListResponse])

	// ListCredentials is the fake for method EndpointsClient.ListCredentials
	// HTTP status codes to indicate success: http.StatusOK
	ListCredentials func(ctx context.Context, resourceURI string, endpointName string, options *armhybridconnectivity.EndpointsClientListCredentialsOptions) (resp azfake.Responder[armhybridconnectivity.EndpointsClientListCredentialsResponse], errResp azfake.ErrorResponder)

	// ListIngressGatewayCredentials is the fake for method EndpointsClient.ListIngressGatewayCredentials
	// HTTP status codes to indicate success: http.StatusOK
	ListIngressGatewayCredentials func(ctx context.Context, resourceURI string, endpointName string, options *armhybridconnectivity.EndpointsClientListIngressGatewayCredentialsOptions) (resp azfake.Responder[armhybridconnectivity.EndpointsClientListIngressGatewayCredentialsResponse], errResp azfake.ErrorResponder)

	// ListManagedProxyDetails is the fake for method EndpointsClient.ListManagedProxyDetails
	// HTTP status codes to indicate success: http.StatusOK
	ListManagedProxyDetails func(ctx context.Context, resourceURI string, endpointName string, managedProxyRequest armhybridconnectivity.ManagedProxyRequest, options *armhybridconnectivity.EndpointsClientListManagedProxyDetailsOptions) (resp azfake.Responder[armhybridconnectivity.EndpointsClientListManagedProxyDetailsResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method EndpointsClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceURI string, endpointName string, endpointResource armhybridconnectivity.EndpointResource, options *armhybridconnectivity.EndpointsClientUpdateOptions) (resp azfake.Responder[armhybridconnectivity.EndpointsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewEndpointsServerTransport creates a new instance of EndpointsServerTransport with the provided implementation.
// The returned EndpointsServerTransport instance is connected to an instance of armhybridconnectivity.EndpointsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewEndpointsServerTransport(srv *EndpointsServer) *EndpointsServerTransport {
	return &EndpointsServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armhybridconnectivity.EndpointsClientListResponse]](),
	}
}

// EndpointsServerTransport connects instances of armhybridconnectivity.EndpointsClient to instances of EndpointsServer.
// Don't use this type directly, use NewEndpointsServerTransport instead.
type EndpointsServerTransport struct {
	srv          *EndpointsServer
	newListPager *tracker[azfake.PagerResponder[armhybridconnectivity.EndpointsClientListResponse]]
}

// Do implements the policy.Transporter interface for EndpointsServerTransport.
func (e *EndpointsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "EndpointsClient.CreateOrUpdate":
		resp, err = e.dispatchCreateOrUpdate(req)
	case "EndpointsClient.Delete":
		resp, err = e.dispatchDelete(req)
	case "EndpointsClient.Get":
		resp, err = e.dispatchGet(req)
	case "EndpointsClient.NewListPager":
		resp, err = e.dispatchNewListPager(req)
	case "EndpointsClient.ListCredentials":
		resp, err = e.dispatchListCredentials(req)
	case "EndpointsClient.ListIngressGatewayCredentials":
		resp, err = e.dispatchListIngressGatewayCredentials(req)
	case "EndpointsClient.ListManagedProxyDetails":
		resp, err = e.dispatchListManagedProxyDetails(req)
	case "EndpointsClient.Update":
		resp, err = e.dispatchUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (e *EndpointsServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if e.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/(?P<resourceUri>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/endpoints/(?P<endpointName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armhybridconnectivity.EndpointResource](req)
	if err != nil {
		return nil, err
	}
	resourceURIParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceUri")])
	if err != nil {
		return nil, err
	}
	endpointNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("endpointName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := e.srv.CreateOrUpdate(req.Context(), resourceURIParam, endpointNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).EndpointResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *EndpointsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if e.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/(?P<resourceUri>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/endpoints/(?P<endpointName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceURIParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceUri")])
	if err != nil {
		return nil, err
	}
	endpointNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("endpointName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := e.srv.Delete(req.Context(), resourceURIParam, endpointNameParam, nil)
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

func (e *EndpointsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if e.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/(?P<resourceUri>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/endpoints/(?P<endpointName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceURIParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceUri")])
	if err != nil {
		return nil, err
	}
	endpointNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("endpointName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := e.srv.Get(req.Context(), resourceURIParam, endpointNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).EndpointResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *EndpointsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if e.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := e.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/(?P<resourceUri>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/endpoints`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceURIParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceUri")])
		if err != nil {
			return nil, err
		}
		resp := e.srv.NewListPager(resourceURIParam, nil)
		newListPager = &resp
		e.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armhybridconnectivity.EndpointsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		e.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		e.newListPager.remove(req)
	}
	return resp, nil
}

func (e *EndpointsServerTransport) dispatchListCredentials(req *http.Request) (*http.Response, error) {
	if e.srv.ListCredentials == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListCredentials not implemented")}
	}
	const regexStr = `/(?P<resourceUri>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/endpoints/(?P<endpointName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listCredentials`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	body, err := server.UnmarshalRequestAsJSON[armhybridconnectivity.ListCredentialsRequest](req)
	if err != nil {
		return nil, err
	}
	resourceURIParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceUri")])
	if err != nil {
		return nil, err
	}
	endpointNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("endpointName")])
	if err != nil {
		return nil, err
	}
	expiresinUnescaped, err := url.QueryUnescape(qp.Get("expiresin"))
	if err != nil {
		return nil, err
	}
	expiresinParam, err := parseOptional(expiresinUnescaped, func(v string) (int64, error) {
		p, parseErr := strconv.ParseInt(v, 10, 64)
		if parseErr != nil {
			return 0, parseErr
		}
		return p, nil
	})
	if err != nil {
		return nil, err
	}
	var options *armhybridconnectivity.EndpointsClientListCredentialsOptions
	if expiresinParam != nil || !reflect.ValueOf(body).IsZero() {
		options = &armhybridconnectivity.EndpointsClientListCredentialsOptions{
			Expiresin:              expiresinParam,
			ListCredentialsRequest: &body,
		}
	}
	respr, errRespr := e.srv.ListCredentials(req.Context(), resourceURIParam, endpointNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).EndpointAccessResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *EndpointsServerTransport) dispatchListIngressGatewayCredentials(req *http.Request) (*http.Response, error) {
	if e.srv.ListIngressGatewayCredentials == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListIngressGatewayCredentials not implemented")}
	}
	const regexStr = `/(?P<resourceUri>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/endpoints/(?P<endpointName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listIngressGatewayCredentials`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	body, err := server.UnmarshalRequestAsJSON[armhybridconnectivity.ListIngressGatewayCredentialsRequest](req)
	if err != nil {
		return nil, err
	}
	resourceURIParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceUri")])
	if err != nil {
		return nil, err
	}
	endpointNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("endpointName")])
	if err != nil {
		return nil, err
	}
	expiresinUnescaped, err := url.QueryUnescape(qp.Get("expiresin"))
	if err != nil {
		return nil, err
	}
	expiresinParam, err := parseOptional(expiresinUnescaped, func(v string) (int64, error) {
		p, parseErr := strconv.ParseInt(v, 10, 64)
		if parseErr != nil {
			return 0, parseErr
		}
		return p, nil
	})
	if err != nil {
		return nil, err
	}
	var options *armhybridconnectivity.EndpointsClientListIngressGatewayCredentialsOptions
	if expiresinParam != nil || !reflect.ValueOf(body).IsZero() {
		options = &armhybridconnectivity.EndpointsClientListIngressGatewayCredentialsOptions{
			Expiresin:                            expiresinParam,
			ListIngressGatewayCredentialsRequest: &body,
		}
	}
	respr, errRespr := e.srv.ListIngressGatewayCredentials(req.Context(), resourceURIParam, endpointNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).IngressGatewayResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *EndpointsServerTransport) dispatchListManagedProxyDetails(req *http.Request) (*http.Response, error) {
	if e.srv.ListManagedProxyDetails == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListManagedProxyDetails not implemented")}
	}
	const regexStr = `/(?P<resourceUri>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/endpoints/(?P<endpointName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listManagedProxyDetails`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armhybridconnectivity.ManagedProxyRequest](req)
	if err != nil {
		return nil, err
	}
	resourceURIParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceUri")])
	if err != nil {
		return nil, err
	}
	endpointNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("endpointName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := e.srv.ListManagedProxyDetails(req.Context(), resourceURIParam, endpointNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ManagedProxyResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (e *EndpointsServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if e.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/(?P<resourceUri>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HybridConnectivity/endpoints/(?P<endpointName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armhybridconnectivity.EndpointResource](req)
	if err != nil {
		return nil, err
	}
	resourceURIParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceUri")])
	if err != nil {
		return nil, err
	}
	endpointNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("endpointName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := e.srv.Update(req.Context(), resourceURIParam, endpointNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).EndpointResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}