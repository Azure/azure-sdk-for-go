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

// APITagDescriptionServer is a fake server for instances of the armapimanagement.APITagDescriptionClient type.
type APITagDescriptionServer struct {
	// CreateOrUpdate is the fake for method APITagDescriptionClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, serviceName string, apiID string, tagDescriptionID string, parameters armapimanagement.TagDescriptionCreateParameters, options *armapimanagement.APITagDescriptionClientCreateOrUpdateOptions) (resp azfake.Responder[armapimanagement.APITagDescriptionClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method APITagDescriptionClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, serviceName string, apiID string, tagDescriptionID string, ifMatch string, options *armapimanagement.APITagDescriptionClientDeleteOptions) (resp azfake.Responder[armapimanagement.APITagDescriptionClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method APITagDescriptionClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serviceName string, apiID string, tagDescriptionID string, options *armapimanagement.APITagDescriptionClientGetOptions) (resp azfake.Responder[armapimanagement.APITagDescriptionClientGetResponse], errResp azfake.ErrorResponder)

	// GetEntityTag is the fake for method APITagDescriptionClient.GetEntityTag
	// HTTP status codes to indicate success: http.StatusOK
	GetEntityTag func(ctx context.Context, resourceGroupName string, serviceName string, apiID string, tagDescriptionID string, options *armapimanagement.APITagDescriptionClientGetEntityTagOptions) (resp azfake.Responder[armapimanagement.APITagDescriptionClientGetEntityTagResponse], errResp azfake.ErrorResponder)

	// NewListByServicePager is the fake for method APITagDescriptionClient.NewListByServicePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByServicePager func(resourceGroupName string, serviceName string, apiID string, options *armapimanagement.APITagDescriptionClientListByServiceOptions) (resp azfake.PagerResponder[armapimanagement.APITagDescriptionClientListByServiceResponse])
}

// NewAPITagDescriptionServerTransport creates a new instance of APITagDescriptionServerTransport with the provided implementation.
// The returned APITagDescriptionServerTransport instance is connected to an instance of armapimanagement.APITagDescriptionClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewAPITagDescriptionServerTransport(srv *APITagDescriptionServer) *APITagDescriptionServerTransport {
	return &APITagDescriptionServerTransport{
		srv:                   srv,
		newListByServicePager: newTracker[azfake.PagerResponder[armapimanagement.APITagDescriptionClientListByServiceResponse]](),
	}
}

// APITagDescriptionServerTransport connects instances of armapimanagement.APITagDescriptionClient to instances of APITagDescriptionServer.
// Don't use this type directly, use NewAPITagDescriptionServerTransport instead.
type APITagDescriptionServerTransport struct {
	srv                   *APITagDescriptionServer
	newListByServicePager *tracker[azfake.PagerResponder[armapimanagement.APITagDescriptionClientListByServiceResponse]]
}

// Do implements the policy.Transporter interface for APITagDescriptionServerTransport.
func (a *APITagDescriptionServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return a.dispatchToMethodFake(req, method)
}

func (a *APITagDescriptionServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if apiTagDescriptionServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = apiTagDescriptionServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "APITagDescriptionClient.CreateOrUpdate":
				res.resp, res.err = a.dispatchCreateOrUpdate(req)
			case "APITagDescriptionClient.Delete":
				res.resp, res.err = a.dispatchDelete(req)
			case "APITagDescriptionClient.Get":
				res.resp, res.err = a.dispatchGet(req)
			case "APITagDescriptionClient.GetEntityTag":
				res.resp, res.err = a.dispatchGetEntityTag(req)
			case "APITagDescriptionClient.NewListByServicePager":
				res.resp, res.err = a.dispatchNewListByServicePager(req)
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

func (a *APITagDescriptionServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if a.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/apis/(?P<apiId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/tagDescriptions/(?P<tagDescriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armapimanagement.TagDescriptionCreateParameters](req)
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
	apiIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("apiId")])
	if err != nil {
		return nil, err
	}
	tagDescriptionIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("tagDescriptionId")])
	if err != nil {
		return nil, err
	}
	ifMatchParam := getOptional(getHeaderValue(req.Header, "If-Match"))
	var options *armapimanagement.APITagDescriptionClientCreateOrUpdateOptions
	if ifMatchParam != nil {
		options = &armapimanagement.APITagDescriptionClientCreateOrUpdateOptions{
			IfMatch: ifMatchParam,
		}
	}
	respr, errRespr := a.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, serviceNameParam, apiIDParam, tagDescriptionIDParam, body, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).TagDescriptionContract, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (a *APITagDescriptionServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if a.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/apis/(?P<apiId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/tagDescriptions/(?P<tagDescriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	apiIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("apiId")])
	if err != nil {
		return nil, err
	}
	tagDescriptionIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("tagDescriptionId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Delete(req.Context(), resourceGroupNameParam, serviceNameParam, apiIDParam, tagDescriptionIDParam, getHeaderValue(req.Header, "If-Match"), nil)
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

func (a *APITagDescriptionServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if a.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/apis/(?P<apiId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/tagDescriptions/(?P<tagDescriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	apiIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("apiId")])
	if err != nil {
		return nil, err
	}
	tagDescriptionIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("tagDescriptionId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.Get(req.Context(), resourceGroupNameParam, serviceNameParam, apiIDParam, tagDescriptionIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).TagDescriptionContract, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (a *APITagDescriptionServerTransport) dispatchGetEntityTag(req *http.Request) (*http.Response, error) {
	if a.srv.GetEntityTag == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetEntityTag not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/apis/(?P<apiId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/tagDescriptions/(?P<tagDescriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	apiIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("apiId")])
	if err != nil {
		return nil, err
	}
	tagDescriptionIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("tagDescriptionId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := a.srv.GetEntityTag(req.Context(), resourceGroupNameParam, serviceNameParam, apiIDParam, tagDescriptionIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (a *APITagDescriptionServerTransport) dispatchNewListByServicePager(req *http.Request) (*http.Response, error) {
	if a.srv.NewListByServicePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByServicePager not implemented")}
	}
	newListByServicePager := a.newListByServicePager.get(req)
	if newListByServicePager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/apis/(?P<apiId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/tagDescriptions`
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
		apiIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("apiId")])
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
		var options *armapimanagement.APITagDescriptionClientListByServiceOptions
		if filterParam != nil || topParam != nil || skipParam != nil {
			options = &armapimanagement.APITagDescriptionClientListByServiceOptions{
				Filter: filterParam,
				Top:    topParam,
				Skip:   skipParam,
			}
		}
		resp := a.srv.NewListByServicePager(resourceGroupNameParam, serviceNameParam, apiIDParam, options)
		newListByServicePager = &resp
		a.newListByServicePager.add(req, newListByServicePager)
		server.PagerResponderInjectNextLinks(newListByServicePager, req, func(page *armapimanagement.APITagDescriptionClientListByServiceResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByServicePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		a.newListByServicePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByServicePager) {
		a.newListByServicePager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to APITagDescriptionServerTransport
var apiTagDescriptionServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
