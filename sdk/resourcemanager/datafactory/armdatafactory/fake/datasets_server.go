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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/datafactory/armdatafactory/v4"
	"net/http"
	"net/url"
	"regexp"
)

// DatasetsServer is a fake server for instances of the armdatafactory.DatasetsClient type.
type DatasetsServer struct {
	// CreateOrUpdate is the fake for method DatasetsClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, factoryName string, datasetName string, dataset armdatafactory.DatasetResource, options *armdatafactory.DatasetsClientCreateOrUpdateOptions) (resp azfake.Responder[armdatafactory.DatasetsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method DatasetsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, factoryName string, datasetName string, options *armdatafactory.DatasetsClientDeleteOptions) (resp azfake.Responder[armdatafactory.DatasetsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method DatasetsClient.Get
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNotModified
	Get func(ctx context.Context, resourceGroupName string, factoryName string, datasetName string, options *armdatafactory.DatasetsClientGetOptions) (resp azfake.Responder[armdatafactory.DatasetsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByFactoryPager is the fake for method DatasetsClient.NewListByFactoryPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByFactoryPager func(resourceGroupName string, factoryName string, options *armdatafactory.DatasetsClientListByFactoryOptions) (resp azfake.PagerResponder[armdatafactory.DatasetsClientListByFactoryResponse])
}

// NewDatasetsServerTransport creates a new instance of DatasetsServerTransport with the provided implementation.
// The returned DatasetsServerTransport instance is connected to an instance of armdatafactory.DatasetsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDatasetsServerTransport(srv *DatasetsServer) *DatasetsServerTransport {
	return &DatasetsServerTransport{
		srv:                   srv,
		newListByFactoryPager: newTracker[azfake.PagerResponder[armdatafactory.DatasetsClientListByFactoryResponse]](),
	}
}

// DatasetsServerTransport connects instances of armdatafactory.DatasetsClient to instances of DatasetsServer.
// Don't use this type directly, use NewDatasetsServerTransport instead.
type DatasetsServerTransport struct {
	srv                   *DatasetsServer
	newListByFactoryPager *tracker[azfake.PagerResponder[armdatafactory.DatasetsClientListByFactoryResponse]]
}

// Do implements the policy.Transporter interface for DatasetsServerTransport.
func (d *DatasetsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "DatasetsClient.CreateOrUpdate":
		resp, err = d.dispatchCreateOrUpdate(req)
	case "DatasetsClient.Delete":
		resp, err = d.dispatchDelete(req)
	case "DatasetsClient.Get":
		resp, err = d.dispatchGet(req)
	case "DatasetsClient.NewListByFactoryPager":
		resp, err = d.dispatchNewListByFactoryPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DatasetsServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if d.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/datasets/(?P<datasetName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdatafactory.DatasetResource](req)
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
	datasetNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("datasetName")])
	if err != nil {
		return nil, err
	}
	ifMatchParam := getOptional(getHeaderValue(req.Header, "If-Match"))
	var options *armdatafactory.DatasetsClientCreateOrUpdateOptions
	if ifMatchParam != nil {
		options = &armdatafactory.DatasetsClientCreateOrUpdateOptions{
			IfMatch: ifMatchParam,
		}
	}
	respr, errRespr := d.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, factoryNameParam, datasetNameParam, body, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DatasetResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DatasetsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if d.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/datasets/(?P<datasetName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	datasetNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("datasetName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.Delete(req.Context(), resourceGroupNameParam, factoryNameParam, datasetNameParam, nil)
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

func (d *DatasetsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if d.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/datasets/(?P<datasetName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	datasetNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("datasetName")])
	if err != nil {
		return nil, err
	}
	ifNoneMatchParam := getOptional(getHeaderValue(req.Header, "If-None-Match"))
	var options *armdatafactory.DatasetsClientGetOptions
	if ifNoneMatchParam != nil {
		options = &armdatafactory.DatasetsClientGetOptions{
			IfNoneMatch: ifNoneMatchParam,
		}
	}
	respr, errRespr := d.srv.Get(req.Context(), resourceGroupNameParam, factoryNameParam, datasetNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusNotModified}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNotModified", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).DatasetResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DatasetsServerTransport) dispatchNewListByFactoryPager(req *http.Request) (*http.Response, error) {
	if d.srv.NewListByFactoryPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByFactoryPager not implemented")}
	}
	newListByFactoryPager := d.newListByFactoryPager.get(req)
	if newListByFactoryPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DataFactory/factories/(?P<factoryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/datasets`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
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
		resp := d.srv.NewListByFactoryPager(resourceGroupNameParam, factoryNameParam, nil)
		newListByFactoryPager = &resp
		d.newListByFactoryPager.add(req, newListByFactoryPager)
		server.PagerResponderInjectNextLinks(newListByFactoryPager, req, func(page *armdatafactory.DatasetsClientListByFactoryResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByFactoryPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		d.newListByFactoryPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByFactoryPager) {
		d.newListByFactoryPager.remove(req)
	}
	return resp, nil
}