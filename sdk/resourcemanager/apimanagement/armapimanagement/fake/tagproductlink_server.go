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

// TagProductLinkServer is a fake server for instances of the armapimanagement.TagProductLinkClient type.
type TagProductLinkServer struct {
	// CreateOrUpdate is the fake for method TagProductLinkClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, serviceName string, tagID string, productLinkID string, parameters armapimanagement.TagProductLinkContract, options *armapimanagement.TagProductLinkClientCreateOrUpdateOptions) (resp azfake.Responder[armapimanagement.TagProductLinkClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method TagProductLinkClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, serviceName string, tagID string, productLinkID string, options *armapimanagement.TagProductLinkClientDeleteOptions) (resp azfake.Responder[armapimanagement.TagProductLinkClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method TagProductLinkClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serviceName string, tagID string, productLinkID string, options *armapimanagement.TagProductLinkClientGetOptions) (resp azfake.Responder[armapimanagement.TagProductLinkClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByProductPager is the fake for method TagProductLinkClient.NewListByProductPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByProductPager func(resourceGroupName string, serviceName string, tagID string, options *armapimanagement.TagProductLinkClientListByProductOptions) (resp azfake.PagerResponder[armapimanagement.TagProductLinkClientListByProductResponse])
}

// NewTagProductLinkServerTransport creates a new instance of TagProductLinkServerTransport with the provided implementation.
// The returned TagProductLinkServerTransport instance is connected to an instance of armapimanagement.TagProductLinkClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewTagProductLinkServerTransport(srv *TagProductLinkServer) *TagProductLinkServerTransport {
	return &TagProductLinkServerTransport{
		srv:                   srv,
		newListByProductPager: newTracker[azfake.PagerResponder[armapimanagement.TagProductLinkClientListByProductResponse]](),
	}
}

// TagProductLinkServerTransport connects instances of armapimanagement.TagProductLinkClient to instances of TagProductLinkServer.
// Don't use this type directly, use NewTagProductLinkServerTransport instead.
type TagProductLinkServerTransport struct {
	srv                   *TagProductLinkServer
	newListByProductPager *tracker[azfake.PagerResponder[armapimanagement.TagProductLinkClientListByProductResponse]]
}

// Do implements the policy.Transporter interface for TagProductLinkServerTransport.
func (t *TagProductLinkServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return t.dispatchToMethodFake(req, method)
}

func (t *TagProductLinkServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if tagProductLinkServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = tagProductLinkServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "TagProductLinkClient.CreateOrUpdate":
				res.resp, res.err = t.dispatchCreateOrUpdate(req)
			case "TagProductLinkClient.Delete":
				res.resp, res.err = t.dispatchDelete(req)
			case "TagProductLinkClient.Get":
				res.resp, res.err = t.dispatchGet(req)
			case "TagProductLinkClient.NewListByProductPager":
				res.resp, res.err = t.dispatchNewListByProductPager(req)
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

func (t *TagProductLinkServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if t.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/tags/(?P<tagId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/productLinks/(?P<productLinkId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armapimanagement.TagProductLinkContract](req)
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
	tagIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("tagId")])
	if err != nil {
		return nil, err
	}
	productLinkIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("productLinkId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := t.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, serviceNameParam, tagIDParam, productLinkIDParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).TagProductLinkContract, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (t *TagProductLinkServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if t.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/tags/(?P<tagId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/productLinks/(?P<productLinkId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
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
	tagIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("tagId")])
	if err != nil {
		return nil, err
	}
	productLinkIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("productLinkId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := t.srv.Delete(req.Context(), resourceGroupNameParam, serviceNameParam, tagIDParam, productLinkIDParam, nil)
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

func (t *TagProductLinkServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if t.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/tags/(?P<tagId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/productLinks/(?P<productLinkId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
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
	tagIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("tagId")])
	if err != nil {
		return nil, err
	}
	productLinkIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("productLinkId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := t.srv.Get(req.Context(), resourceGroupNameParam, serviceNameParam, tagIDParam, productLinkIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).TagProductLinkContract, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (t *TagProductLinkServerTransport) dispatchNewListByProductPager(req *http.Request) (*http.Response, error) {
	if t.srv.NewListByProductPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByProductPager not implemented")}
	}
	newListByProductPager := t.newListByProductPager.get(req)
	if newListByProductPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/tags/(?P<tagId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/productLinks`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
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
		tagIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("tagId")])
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
		var options *armapimanagement.TagProductLinkClientListByProductOptions
		if filterParam != nil || topParam != nil || skipParam != nil {
			options = &armapimanagement.TagProductLinkClientListByProductOptions{
				Filter: filterParam,
				Top:    topParam,
				Skip:   skipParam,
			}
		}
		resp := t.srv.NewListByProductPager(resourceGroupNameParam, serviceNameParam, tagIDParam, options)
		newListByProductPager = &resp
		t.newListByProductPager.add(req, newListByProductPager)
		server.PagerResponderInjectNextLinks(newListByProductPager, req, func(page *armapimanagement.TagProductLinkClientListByProductResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByProductPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		t.newListByProductPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByProductPager) {
		t.newListByProductPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to TagProductLinkServerTransport
var tagProductLinkServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
