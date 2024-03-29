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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/streamanalytics/armstreamanalytics/v2"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

// FunctionsServer is a fake server for instances of the armstreamanalytics.FunctionsClient type.
type FunctionsServer struct {
	// CreateOrReplace is the fake for method FunctionsClient.CreateOrReplace
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrReplace func(ctx context.Context, resourceGroupName string, jobName string, functionName string, function armstreamanalytics.Function, options *armstreamanalytics.FunctionsClientCreateOrReplaceOptions) (resp azfake.Responder[armstreamanalytics.FunctionsClientCreateOrReplaceResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method FunctionsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *armstreamanalytics.FunctionsClientDeleteOptions) (resp azfake.Responder[armstreamanalytics.FunctionsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method FunctionsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *armstreamanalytics.FunctionsClientGetOptions) (resp azfake.Responder[armstreamanalytics.FunctionsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByStreamingJobPager is the fake for method FunctionsClient.NewListByStreamingJobPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByStreamingJobPager func(resourceGroupName string, jobName string, options *armstreamanalytics.FunctionsClientListByStreamingJobOptions) (resp azfake.PagerResponder[armstreamanalytics.FunctionsClientListByStreamingJobResponse])

	// RetrieveDefaultDefinition is the fake for method FunctionsClient.RetrieveDefaultDefinition
	// HTTP status codes to indicate success: http.StatusOK
	RetrieveDefaultDefinition func(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *armstreamanalytics.FunctionsClientRetrieveDefaultDefinitionOptions) (resp azfake.Responder[armstreamanalytics.FunctionsClientRetrieveDefaultDefinitionResponse], errResp azfake.ErrorResponder)

	// BeginTest is the fake for method FunctionsClient.BeginTest
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginTest func(ctx context.Context, resourceGroupName string, jobName string, functionName string, options *armstreamanalytics.FunctionsClientBeginTestOptions) (resp azfake.PollerResponder[armstreamanalytics.FunctionsClientTestResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method FunctionsClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, jobName string, functionName string, function armstreamanalytics.Function, options *armstreamanalytics.FunctionsClientUpdateOptions) (resp azfake.Responder[armstreamanalytics.FunctionsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewFunctionsServerTransport creates a new instance of FunctionsServerTransport with the provided implementation.
// The returned FunctionsServerTransport instance is connected to an instance of armstreamanalytics.FunctionsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewFunctionsServerTransport(srv *FunctionsServer) *FunctionsServerTransport {
	return &FunctionsServerTransport{
		srv:                        srv,
		newListByStreamingJobPager: newTracker[azfake.PagerResponder[armstreamanalytics.FunctionsClientListByStreamingJobResponse]](),
		beginTest:                  newTracker[azfake.PollerResponder[armstreamanalytics.FunctionsClientTestResponse]](),
	}
}

// FunctionsServerTransport connects instances of armstreamanalytics.FunctionsClient to instances of FunctionsServer.
// Don't use this type directly, use NewFunctionsServerTransport instead.
type FunctionsServerTransport struct {
	srv                        *FunctionsServer
	newListByStreamingJobPager *tracker[azfake.PagerResponder[armstreamanalytics.FunctionsClientListByStreamingJobResponse]]
	beginTest                  *tracker[azfake.PollerResponder[armstreamanalytics.FunctionsClientTestResponse]]
}

// Do implements the policy.Transporter interface for FunctionsServerTransport.
func (f *FunctionsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "FunctionsClient.CreateOrReplace":
		resp, err = f.dispatchCreateOrReplace(req)
	case "FunctionsClient.Delete":
		resp, err = f.dispatchDelete(req)
	case "FunctionsClient.Get":
		resp, err = f.dispatchGet(req)
	case "FunctionsClient.NewListByStreamingJobPager":
		resp, err = f.dispatchNewListByStreamingJobPager(req)
	case "FunctionsClient.RetrieveDefaultDefinition":
		resp, err = f.dispatchRetrieveDefaultDefinition(req)
	case "FunctionsClient.BeginTest":
		resp, err = f.dispatchBeginTest(req)
	case "FunctionsClient.Update":
		resp, err = f.dispatchUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (f *FunctionsServerTransport) dispatchCreateOrReplace(req *http.Request) (*http.Response, error) {
	if f.srv.CreateOrReplace == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrReplace not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.StreamAnalytics/streamingjobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/functions/(?P<functionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armstreamanalytics.Function](req)
	if err != nil {
		return nil, err
	}
	ifMatchParam := getOptional(getHeaderValue(req.Header, "If-Match"))
	ifNoneMatchParam := getOptional(getHeaderValue(req.Header, "If-None-Match"))
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
	if err != nil {
		return nil, err
	}
	functionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("functionName")])
	if err != nil {
		return nil, err
	}
	var options *armstreamanalytics.FunctionsClientCreateOrReplaceOptions
	if ifMatchParam != nil || ifNoneMatchParam != nil {
		options = &armstreamanalytics.FunctionsClientCreateOrReplaceOptions{
			IfMatch:     ifMatchParam,
			IfNoneMatch: ifNoneMatchParam,
		}
	}
	respr, errRespr := f.srv.CreateOrReplace(req.Context(), resourceGroupNameParam, jobNameParam, functionNameParam, body, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Function, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (f *FunctionsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if f.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.StreamAnalytics/streamingjobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/functions/(?P<functionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
	if err != nil {
		return nil, err
	}
	functionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("functionName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := f.srv.Delete(req.Context(), resourceGroupNameParam, jobNameParam, functionNameParam, nil)
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

func (f *FunctionsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if f.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.StreamAnalytics/streamingjobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/functions/(?P<functionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
	if err != nil {
		return nil, err
	}
	functionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("functionName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := f.srv.Get(req.Context(), resourceGroupNameParam, jobNameParam, functionNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Function, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (f *FunctionsServerTransport) dispatchNewListByStreamingJobPager(req *http.Request) (*http.Response, error) {
	if f.srv.NewListByStreamingJobPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByStreamingJobPager not implemented")}
	}
	newListByStreamingJobPager := f.newListByStreamingJobPager.get(req)
	if newListByStreamingJobPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.StreamAnalytics/streamingjobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/functions`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		selectUnescaped, err := url.QueryUnescape(qp.Get("$select"))
		if err != nil {
			return nil, err
		}
		selectParam := getOptional(selectUnescaped)
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
		if err != nil {
			return nil, err
		}
		var options *armstreamanalytics.FunctionsClientListByStreamingJobOptions
		if selectParam != nil {
			options = &armstreamanalytics.FunctionsClientListByStreamingJobOptions{
				Select: selectParam,
			}
		}
		resp := f.srv.NewListByStreamingJobPager(resourceGroupNameParam, jobNameParam, options)
		newListByStreamingJobPager = &resp
		f.newListByStreamingJobPager.add(req, newListByStreamingJobPager)
		server.PagerResponderInjectNextLinks(newListByStreamingJobPager, req, func(page *armstreamanalytics.FunctionsClientListByStreamingJobResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByStreamingJobPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		f.newListByStreamingJobPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByStreamingJobPager) {
		f.newListByStreamingJobPager.remove(req)
	}
	return resp, nil
}

func (f *FunctionsServerTransport) dispatchRetrieveDefaultDefinition(req *http.Request) (*http.Response, error) {
	if f.srv.RetrieveDefaultDefinition == nil {
		return nil, &nonRetriableError{errors.New("fake for method RetrieveDefaultDefinition not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.StreamAnalytics/streamingjobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/functions/(?P<functionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/retrieveDefaultDefinition`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	raw, err := readRequestBody(req)
	if err != nil {
		return nil, err
	}
	body, err := unmarshalFunctionRetrieveDefaultDefinitionParametersClassification(raw)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
	if err != nil {
		return nil, err
	}
	functionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("functionName")])
	if err != nil {
		return nil, err
	}
	var options *armstreamanalytics.FunctionsClientRetrieveDefaultDefinitionOptions
	if !reflect.ValueOf(body).IsZero() {
		options = &armstreamanalytics.FunctionsClientRetrieveDefaultDefinitionOptions{
			FunctionRetrieveDefaultDefinitionParameters: body,
		}
	}
	respr, errRespr := f.srv.RetrieveDefaultDefinition(req.Context(), resourceGroupNameParam, jobNameParam, functionNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Function, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (f *FunctionsServerTransport) dispatchBeginTest(req *http.Request) (*http.Response, error) {
	if f.srv.BeginTest == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginTest not implemented")}
	}
	beginTest := f.beginTest.get(req)
	if beginTest == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.StreamAnalytics/streamingjobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/functions/(?P<functionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/test`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armstreamanalytics.Function](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
		if err != nil {
			return nil, err
		}
		functionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("functionName")])
		if err != nil {
			return nil, err
		}
		var options *armstreamanalytics.FunctionsClientBeginTestOptions
		if !reflect.ValueOf(body).IsZero() {
			options = &armstreamanalytics.FunctionsClientBeginTestOptions{
				Function: &body,
			}
		}
		respr, errRespr := f.srv.BeginTest(req.Context(), resourceGroupNameParam, jobNameParam, functionNameParam, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginTest = &respr
		f.beginTest.add(req, beginTest)
	}

	resp, err := server.PollerResponderNext(beginTest, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		f.beginTest.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginTest) {
		f.beginTest.remove(req)
	}

	return resp, nil
}

func (f *FunctionsServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if f.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.StreamAnalytics/streamingjobs/(?P<jobName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/functions/(?P<functionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armstreamanalytics.Function](req)
	if err != nil {
		return nil, err
	}
	ifMatchParam := getOptional(getHeaderValue(req.Header, "If-Match"))
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	jobNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("jobName")])
	if err != nil {
		return nil, err
	}
	functionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("functionName")])
	if err != nil {
		return nil, err
	}
	var options *armstreamanalytics.FunctionsClientUpdateOptions
	if ifMatchParam != nil {
		options = &armstreamanalytics.FunctionsClientUpdateOptions{
			IfMatch: ifMatchParam,
		}
	}
	respr, errRespr := f.srv.Update(req.Context(), resourceGroupNameParam, jobNameParam, functionNameParam, body, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Function, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}
