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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/apimanagement/armapimanagement/v3"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// WorkspaceProductGroupLinkServer is a fake server for instances of the armapimanagement.WorkspaceProductGroupLinkClient type.
type WorkspaceProductGroupLinkServer struct {
	// CreateOrUpdate is the fake for method WorkspaceProductGroupLinkClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, serviceName string, workspaceID string, productID string, groupLinkID string, parameters armapimanagement.ProductGroupLinkContract, options *armapimanagement.WorkspaceProductGroupLinkClientCreateOrUpdateOptions) (resp azfake.Responder[armapimanagement.WorkspaceProductGroupLinkClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method WorkspaceProductGroupLinkClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, serviceName string, workspaceID string, productID string, groupLinkID string, options *armapimanagement.WorkspaceProductGroupLinkClientDeleteOptions) (resp azfake.Responder[armapimanagement.WorkspaceProductGroupLinkClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method WorkspaceProductGroupLinkClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serviceName string, workspaceID string, productID string, groupLinkID string, options *armapimanagement.WorkspaceProductGroupLinkClientGetOptions) (resp azfake.Responder[armapimanagement.WorkspaceProductGroupLinkClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByProductPager is the fake for method WorkspaceProductGroupLinkClient.NewListByProductPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByProductPager func(resourceGroupName string, serviceName string, workspaceID string, productID string, options *armapimanagement.WorkspaceProductGroupLinkClientListByProductOptions) (resp azfake.PagerResponder[armapimanagement.WorkspaceProductGroupLinkClientListByProductResponse])
}

// NewWorkspaceProductGroupLinkServerTransport creates a new instance of WorkspaceProductGroupLinkServerTransport with the provided implementation.
// The returned WorkspaceProductGroupLinkServerTransport instance is connected to an instance of armapimanagement.WorkspaceProductGroupLinkClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewWorkspaceProductGroupLinkServerTransport(srv *WorkspaceProductGroupLinkServer) *WorkspaceProductGroupLinkServerTransport {
	return &WorkspaceProductGroupLinkServerTransport{
		srv:                   srv,
		newListByProductPager: newTracker[azfake.PagerResponder[armapimanagement.WorkspaceProductGroupLinkClientListByProductResponse]](),
	}
}

// WorkspaceProductGroupLinkServerTransport connects instances of armapimanagement.WorkspaceProductGroupLinkClient to instances of WorkspaceProductGroupLinkServer.
// Don't use this type directly, use NewWorkspaceProductGroupLinkServerTransport instead.
type WorkspaceProductGroupLinkServerTransport struct {
	srv                   *WorkspaceProductGroupLinkServer
	newListByProductPager *tracker[azfake.PagerResponder[armapimanagement.WorkspaceProductGroupLinkClientListByProductResponse]]
}

// Do implements the policy.Transporter interface for WorkspaceProductGroupLinkServerTransport.
func (w *WorkspaceProductGroupLinkServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return w.dispatchToMethodFake(req, method)
}

func (w *WorkspaceProductGroupLinkServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if workspaceProductGroupLinkServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = workspaceProductGroupLinkServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "WorkspaceProductGroupLinkClient.CreateOrUpdate":
				res.resp, res.err = w.dispatchCreateOrUpdate(req)
			case "WorkspaceProductGroupLinkClient.Delete":
				res.resp, res.err = w.dispatchDelete(req)
			case "WorkspaceProductGroupLinkClient.Get":
				res.resp, res.err = w.dispatchGet(req)
			case "WorkspaceProductGroupLinkClient.NewListByProductPager":
				res.resp, res.err = w.dispatchNewListByProductPager(req)
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

func (w *WorkspaceProductGroupLinkServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if w.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/workspaces/(?P<workspaceId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/products/(?P<productId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/groupLinks/(?P<groupLinkId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 6 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armapimanagement.ProductGroupLinkContract](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	workspaceIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceId")])
	if err != nil {
		return nil, err
	}
	productIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("productId")])
	if err != nil {
		return nil, err
	}
	groupLinkIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("groupLinkId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := w.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, serviceNameParam, workspaceIDParam, productIDParam, groupLinkIDParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ProductGroupLinkContract, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (w *WorkspaceProductGroupLinkServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if w.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/workspaces/(?P<workspaceId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/products/(?P<productId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/groupLinks/(?P<groupLinkId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 6 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	workspaceIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceId")])
	if err != nil {
		return nil, err
	}
	productIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("productId")])
	if err != nil {
		return nil, err
	}
	groupLinkIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("groupLinkId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := w.srv.Delete(req.Context(), resourceGroupNameParam, serviceNameParam, workspaceIDParam, productIDParam, groupLinkIDParam, nil)
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

func (w *WorkspaceProductGroupLinkServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if w.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/workspaces/(?P<workspaceId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/products/(?P<productId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/groupLinks/(?P<groupLinkId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 6 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	workspaceIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceId")])
	if err != nil {
		return nil, err
	}
	productIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("productId")])
	if err != nil {
		return nil, err
	}
	groupLinkIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("groupLinkId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := w.srv.Get(req.Context(), resourceGroupNameParam, serviceNameParam, workspaceIDParam, productIDParam, groupLinkIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ProductGroupLinkContract, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (w *WorkspaceProductGroupLinkServerTransport) dispatchNewListByProductPager(req *http.Request) (*http.Response, error) {
	if w.srv.NewListByProductPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByProductPager not implemented")}
	}
	newListByProductPager := w.newListByProductPager.get(req)
	if newListByProductPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/workspaces/(?P<workspaceId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/products/(?P<productId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/groupLinks`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
		if err != nil {
			return nil, err
		}
		workspaceIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("workspaceId")])
		if err != nil {
			return nil, err
		}
		productIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("productId")])
		if err != nil {
			return nil, err
		}
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		skipUnescaped, err := url.QueryUnescape(qp.Get("$skip"))
		if err != nil {
			return nil, err
		}
		skipParam, err := parseOptional(skipUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		var options *armapimanagement.WorkspaceProductGroupLinkClientListByProductOptions
		if filterParam != nil || topParam != nil || skipParam != nil {
			options = &armapimanagement.WorkspaceProductGroupLinkClientListByProductOptions{
				Filter: filterParam,
				Top:    topParam,
				Skip:   skipParam,
			}
		}
		resp := w.srv.NewListByProductPager(resourceGroupNameParam, serviceNameParam, workspaceIDParam, productIDParam, options)
		newListByProductPager = &resp
		w.newListByProductPager.add(req, newListByProductPager)
		server.PagerResponderInjectNextLinks(newListByProductPager, req, func(page *armapimanagement.WorkspaceProductGroupLinkClientListByProductResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByProductPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		w.newListByProductPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByProductPager) {
		w.newListByProductPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to WorkspaceProductGroupLinkServerTransport
var workspaceProductGroupLinkServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
