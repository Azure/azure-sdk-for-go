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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/applicationinsights/armapplicationinsights/v2"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
)

// WorkbooksServer is a fake server for instances of the armapplicationinsights.WorkbooksClient type.
type WorkbooksServer struct {
	// CreateOrUpdate is the fake for method WorkbooksClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, resourceName string, workbookProperties armapplicationinsights.Workbook, options *armapplicationinsights.WorkbooksClientCreateOrUpdateOptions) (resp azfake.Responder[armapplicationinsights.WorkbooksClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method WorkbooksClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, resourceName string, options *armapplicationinsights.WorkbooksClientDeleteOptions) (resp azfake.Responder[armapplicationinsights.WorkbooksClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method WorkbooksClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, resourceName string, options *armapplicationinsights.WorkbooksClientGetOptions) (resp azfake.Responder[armapplicationinsights.WorkbooksClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method WorkbooksClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, category armapplicationinsights.CategoryType, options *armapplicationinsights.WorkbooksClientListByResourceGroupOptions) (resp azfake.PagerResponder[armapplicationinsights.WorkbooksClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method WorkbooksClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(category armapplicationinsights.CategoryType, options *armapplicationinsights.WorkbooksClientListBySubscriptionOptions) (resp azfake.PagerResponder[armapplicationinsights.WorkbooksClientListBySubscriptionResponse])

	// RevisionGet is the fake for method WorkbooksClient.RevisionGet
	// HTTP status codes to indicate success: http.StatusOK
	RevisionGet func(ctx context.Context, resourceGroupName string, resourceName string, revisionID string, options *armapplicationinsights.WorkbooksClientRevisionGetOptions) (resp azfake.Responder[armapplicationinsights.WorkbooksClientRevisionGetResponse], errResp azfake.ErrorResponder)

	// NewRevisionsListPager is the fake for method WorkbooksClient.NewRevisionsListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewRevisionsListPager func(resourceGroupName string, resourceName string, options *armapplicationinsights.WorkbooksClientRevisionsListOptions) (resp azfake.PagerResponder[armapplicationinsights.WorkbooksClientRevisionsListResponse])

	// Update is the fake for method WorkbooksClient.Update
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	Update func(ctx context.Context, resourceGroupName string, resourceName string, options *armapplicationinsights.WorkbooksClientUpdateOptions) (resp azfake.Responder[armapplicationinsights.WorkbooksClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewWorkbooksServerTransport creates a new instance of WorkbooksServerTransport with the provided implementation.
// The returned WorkbooksServerTransport instance is connected to an instance of armapplicationinsights.WorkbooksClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewWorkbooksServerTransport(srv *WorkbooksServer) *WorkbooksServerTransport {
	return &WorkbooksServerTransport{
		srv:                         srv,
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armapplicationinsights.WorkbooksClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armapplicationinsights.WorkbooksClientListBySubscriptionResponse]](),
		newRevisionsListPager:       newTracker[azfake.PagerResponder[armapplicationinsights.WorkbooksClientRevisionsListResponse]](),
	}
}

// WorkbooksServerTransport connects instances of armapplicationinsights.WorkbooksClient to instances of WorkbooksServer.
// Don't use this type directly, use NewWorkbooksServerTransport instead.
type WorkbooksServerTransport struct {
	srv                         *WorkbooksServer
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armapplicationinsights.WorkbooksClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armapplicationinsights.WorkbooksClientListBySubscriptionResponse]]
	newRevisionsListPager       *tracker[azfake.PagerResponder[armapplicationinsights.WorkbooksClientRevisionsListResponse]]
}

// Do implements the policy.Transporter interface for WorkbooksServerTransport.
func (w *WorkbooksServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return w.dispatchToMethodFake(req, method)
}

func (w *WorkbooksServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if workbooksServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = workbooksServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "WorkbooksClient.CreateOrUpdate":
				res.resp, res.err = w.dispatchCreateOrUpdate(req)
			case "WorkbooksClient.Delete":
				res.resp, res.err = w.dispatchDelete(req)
			case "WorkbooksClient.Get":
				res.resp, res.err = w.dispatchGet(req)
			case "WorkbooksClient.NewListByResourceGroupPager":
				res.resp, res.err = w.dispatchNewListByResourceGroupPager(req)
			case "WorkbooksClient.NewListBySubscriptionPager":
				res.resp, res.err = w.dispatchNewListBySubscriptionPager(req)
			case "WorkbooksClient.RevisionGet":
				res.resp, res.err = w.dispatchRevisionGet(req)
			case "WorkbooksClient.NewRevisionsListPager":
				res.resp, res.err = w.dispatchNewRevisionsListPager(req)
			case "WorkbooksClient.Update":
				res.resp, res.err = w.dispatchUpdate(req)
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

func (w *WorkbooksServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if w.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Insights/workbooks/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	body, err := server.UnmarshalRequestAsJSON[armapplicationinsights.Workbook](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	sourceIDUnescaped, err := url.QueryUnescape(qp.Get("sourceId"))
	if err != nil {
		return nil, err
	}
	sourceIDParam := getOptional(sourceIDUnescaped)
	var options *armapplicationinsights.WorkbooksClientCreateOrUpdateOptions
	if sourceIDParam != nil {
		options = &armapplicationinsights.WorkbooksClientCreateOrUpdateOptions{
			SourceID: sourceIDParam,
		}
	}
	respr, errRespr := w.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, resourceNameParam, body, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Workbook, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (w *WorkbooksServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if w.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Insights/workbooks/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := w.srv.Delete(req.Context(), resourceGroupNameParam, resourceNameParam, nil)
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

func (w *WorkbooksServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if w.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Insights/workbooks/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	canFetchContentUnescaped, err := url.QueryUnescape(qp.Get("canFetchContent"))
	if err != nil {
		return nil, err
	}
	canFetchContentParam, err := parseOptional(canFetchContentUnescaped, strconv.ParseBool)
	if err != nil {
		return nil, err
	}
	var options *armapplicationinsights.WorkbooksClientGetOptions
	if canFetchContentParam != nil {
		options = &armapplicationinsights.WorkbooksClientGetOptions{
			CanFetchContent: canFetchContentParam,
		}
	}
	respr, errRespr := w.srv.Get(req.Context(), resourceGroupNameParam, resourceNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Workbook, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (w *WorkbooksServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if w.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := w.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Insights/workbooks`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		categoryParam, err := parseWithCast(qp.Get("category"), func(v string) (armapplicationinsights.CategoryType, error) {
			p, unescapeErr := url.QueryUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armapplicationinsights.CategoryType(p), nil
		})
		if err != nil {
			return nil, err
		}
		tagsUnescaped, err := url.QueryUnescape(qp.Get("tags"))
		if err != nil {
			return nil, err
		}
		tagsParam := splitHelper(tagsUnescaped, ",")
		sourceIDUnescaped, err := url.QueryUnescape(qp.Get("sourceId"))
		if err != nil {
			return nil, err
		}
		sourceIDParam := getOptional(sourceIDUnescaped)
		canFetchContentUnescaped, err := url.QueryUnescape(qp.Get("canFetchContent"))
		if err != nil {
			return nil, err
		}
		canFetchContentParam, err := parseOptional(canFetchContentUnescaped, strconv.ParseBool)
		if err != nil {
			return nil, err
		}
		var options *armapplicationinsights.WorkbooksClientListByResourceGroupOptions
		if len(tagsParam) > 0 || sourceIDParam != nil || canFetchContentParam != nil {
			options = &armapplicationinsights.WorkbooksClientListByResourceGroupOptions{
				Tags:            tagsParam,
				SourceID:        sourceIDParam,
				CanFetchContent: canFetchContentParam,
			}
		}
		resp := w.srv.NewListByResourceGroupPager(resourceGroupNameParam, categoryParam, options)
		newListByResourceGroupPager = &resp
		w.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armapplicationinsights.WorkbooksClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		w.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		w.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (w *WorkbooksServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if w.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := w.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Insights/workbooks`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		categoryParam, err := parseWithCast(qp.Get("category"), func(v string) (armapplicationinsights.CategoryType, error) {
			p, unescapeErr := url.QueryUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armapplicationinsights.CategoryType(p), nil
		})
		if err != nil {
			return nil, err
		}
		tagsUnescaped, err := url.QueryUnescape(qp.Get("tags"))
		if err != nil {
			return nil, err
		}
		tagsParam := splitHelper(tagsUnescaped, ",")
		canFetchContentUnescaped, err := url.QueryUnescape(qp.Get("canFetchContent"))
		if err != nil {
			return nil, err
		}
		canFetchContentParam, err := parseOptional(canFetchContentUnescaped, strconv.ParseBool)
		if err != nil {
			return nil, err
		}
		var options *armapplicationinsights.WorkbooksClientListBySubscriptionOptions
		if len(tagsParam) > 0 || canFetchContentParam != nil {
			options = &armapplicationinsights.WorkbooksClientListBySubscriptionOptions{
				Tags:            tagsParam,
				CanFetchContent: canFetchContentParam,
			}
		}
		resp := w.srv.NewListBySubscriptionPager(categoryParam, options)
		newListBySubscriptionPager = &resp
		w.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armapplicationinsights.WorkbooksClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		w.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		w.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (w *WorkbooksServerTransport) dispatchRevisionGet(req *http.Request) (*http.Response, error) {
	if w.srv.RevisionGet == nil {
		return nil, &nonRetriableError{errors.New("fake for method RevisionGet not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Insights/workbooks/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/revisions/(?P<revisionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	revisionIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("revisionId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := w.srv.RevisionGet(req.Context(), resourceGroupNameParam, resourceNameParam, revisionIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Workbook, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (w *WorkbooksServerTransport) dispatchNewRevisionsListPager(req *http.Request) (*http.Response, error) {
	if w.srv.NewRevisionsListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewRevisionsListPager not implemented")}
	}
	newRevisionsListPager := w.newRevisionsListPager.get(req)
	if newRevisionsListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Insights/workbooks/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/revisions`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
		if err != nil {
			return nil, err
		}
		resp := w.srv.NewRevisionsListPager(resourceGroupNameParam, resourceNameParam, nil)
		newRevisionsListPager = &resp
		w.newRevisionsListPager.add(req, newRevisionsListPager)
		server.PagerResponderInjectNextLinks(newRevisionsListPager, req, func(page *armapplicationinsights.WorkbooksClientRevisionsListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newRevisionsListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		w.newRevisionsListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newRevisionsListPager) {
		w.newRevisionsListPager.remove(req)
	}
	return resp, nil
}

func (w *WorkbooksServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if w.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Insights/workbooks/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	body, err := server.UnmarshalRequestAsJSON[armapplicationinsights.WorkbookUpdateParameters](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	sourceIDUnescaped, err := url.QueryUnescape(qp.Get("sourceId"))
	if err != nil {
		return nil, err
	}
	sourceIDParam := getOptional(sourceIDUnescaped)
	var options *armapplicationinsights.WorkbooksClientUpdateOptions
	if sourceIDParam != nil || !reflect.ValueOf(body).IsZero() {
		options = &armapplicationinsights.WorkbooksClientUpdateOptions{
			SourceID:                 sourceIDParam,
			WorkbookUpdateParameters: &body,
		}
	}
	respr, errRespr := w.srv.Update(req.Context(), resourceGroupNameParam, resourceNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Workbook, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to WorkbooksServerTransport
var workbooksServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
