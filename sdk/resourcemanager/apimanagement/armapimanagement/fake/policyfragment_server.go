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

// PolicyFragmentServer is a fake server for instances of the armapimanagement.PolicyFragmentClient type.
type PolicyFragmentServer struct {
	// BeginCreateOrUpdate is the fake for method PolicyFragmentClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, serviceName string, id string, parameters armapimanagement.PolicyFragmentContract, options *armapimanagement.PolicyFragmentClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armapimanagement.PolicyFragmentClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method PolicyFragmentClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, serviceName string, id string, ifMatch string, options *armapimanagement.PolicyFragmentClientDeleteOptions) (resp azfake.Responder[armapimanagement.PolicyFragmentClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method PolicyFragmentClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serviceName string, id string, options *armapimanagement.PolicyFragmentClientGetOptions) (resp azfake.Responder[armapimanagement.PolicyFragmentClientGetResponse], errResp azfake.ErrorResponder)

	// GetEntityTag is the fake for method PolicyFragmentClient.GetEntityTag
	// HTTP status codes to indicate success: http.StatusOK
	GetEntityTag func(ctx context.Context, resourceGroupName string, serviceName string, id string, options *armapimanagement.PolicyFragmentClientGetEntityTagOptions) (resp azfake.Responder[armapimanagement.PolicyFragmentClientGetEntityTagResponse], errResp azfake.ErrorResponder)

	// NewListByServicePager is the fake for method PolicyFragmentClient.NewListByServicePager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByServicePager func(resourceGroupName string, serviceName string, options *armapimanagement.PolicyFragmentClientListByServiceOptions) (resp azfake.PagerResponder[armapimanagement.PolicyFragmentClientListByServiceResponse])

	// ListReferences is the fake for method PolicyFragmentClient.ListReferences
	// HTTP status codes to indicate success: http.StatusOK
	ListReferences func(ctx context.Context, resourceGroupName string, serviceName string, id string, options *armapimanagement.PolicyFragmentClientListReferencesOptions) (resp azfake.Responder[armapimanagement.PolicyFragmentClientListReferencesResponse], errResp azfake.ErrorResponder)
}

// NewPolicyFragmentServerTransport creates a new instance of PolicyFragmentServerTransport with the provided implementation.
// The returned PolicyFragmentServerTransport instance is connected to an instance of armapimanagement.PolicyFragmentClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewPolicyFragmentServerTransport(srv *PolicyFragmentServer) *PolicyFragmentServerTransport {
	return &PolicyFragmentServerTransport{
		srv:                   srv,
		beginCreateOrUpdate:   newTracker[azfake.PollerResponder[armapimanagement.PolicyFragmentClientCreateOrUpdateResponse]](),
		newListByServicePager: newTracker[azfake.PagerResponder[armapimanagement.PolicyFragmentClientListByServiceResponse]](),
	}
}

// PolicyFragmentServerTransport connects instances of armapimanagement.PolicyFragmentClient to instances of PolicyFragmentServer.
// Don't use this type directly, use NewPolicyFragmentServerTransport instead.
type PolicyFragmentServerTransport struct {
	srv                   *PolicyFragmentServer
	beginCreateOrUpdate   *tracker[azfake.PollerResponder[armapimanagement.PolicyFragmentClientCreateOrUpdateResponse]]
	newListByServicePager *tracker[azfake.PagerResponder[armapimanagement.PolicyFragmentClientListByServiceResponse]]
}

// Do implements the policy.Transporter interface for PolicyFragmentServerTransport.
func (p *PolicyFragmentServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return p.dispatchToMethodFake(req, method)
}

func (p *PolicyFragmentServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if policyFragmentServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = policyFragmentServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "PolicyFragmentClient.BeginCreateOrUpdate":
				res.resp, res.err = p.dispatchBeginCreateOrUpdate(req)
			case "PolicyFragmentClient.Delete":
				res.resp, res.err = p.dispatchDelete(req)
			case "PolicyFragmentClient.Get":
				res.resp, res.err = p.dispatchGet(req)
			case "PolicyFragmentClient.GetEntityTag":
				res.resp, res.err = p.dispatchGetEntityTag(req)
			case "PolicyFragmentClient.NewListByServicePager":
				res.resp, res.err = p.dispatchNewListByServicePager(req)
			case "PolicyFragmentClient.ListReferences":
				res.resp, res.err = p.dispatchListReferences(req)
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

func (p *PolicyFragmentServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if p.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := p.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/policyFragments/(?P<id>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armapimanagement.PolicyFragmentContract](req)
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
		idParam, err := url.PathUnescape(matches[regex.SubexpIndex("id")])
		if err != nil {
			return nil, err
		}
		ifMatchParam := getOptional(getHeaderValue(req.Header, "If-Match"))
		var options *armapimanagement.PolicyFragmentClientBeginCreateOrUpdateOptions
		if ifMatchParam != nil {
			options = &armapimanagement.PolicyFragmentClientBeginCreateOrUpdateOptions{
				IfMatch: ifMatchParam,
			}
		}
		respr, errRespr := p.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, serviceNameParam, idParam, body, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		p.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		p.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		p.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (p *PolicyFragmentServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if p.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/policyFragments/(?P<id>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	idParam, err := url.PathUnescape(matches[regex.SubexpIndex("id")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.Delete(req.Context(), resourceGroupNameParam, serviceNameParam, idParam, getHeaderValue(req.Header, "If-Match"), nil)
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

func (p *PolicyFragmentServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if p.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/policyFragments/(?P<id>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	idParam, err := url.PathUnescape(matches[regex.SubexpIndex("id")])
	if err != nil {
		return nil, err
	}
	formatUnescaped, err := url.QueryUnescape(qp.Get("format"))
	if err != nil {
		return nil, err
	}
	formatParam := getOptional(armapimanagement.PolicyFragmentContentFormat(formatUnescaped))
	var options *armapimanagement.PolicyFragmentClientGetOptions
	if formatParam != nil {
		options = &armapimanagement.PolicyFragmentClientGetOptions{
			Format: formatParam,
		}
	}
	respr, errRespr := p.srv.Get(req.Context(), resourceGroupNameParam, serviceNameParam, idParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).PolicyFragmentContract, req)
	if err != nil {
		return nil, err
	}
	if val := server.GetResponse(respr).ETag; val != nil {
		resp.Header.Set("ETag", *val)
	}
	return resp, nil
}

func (p *PolicyFragmentServerTransport) dispatchGetEntityTag(req *http.Request) (*http.Response, error) {
	if p.srv.GetEntityTag == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetEntityTag not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/policyFragments/(?P<id>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	idParam, err := url.PathUnescape(matches[regex.SubexpIndex("id")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.GetEntityTag(req.Context(), resourceGroupNameParam, serviceNameParam, idParam, nil)
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

func (p *PolicyFragmentServerTransport) dispatchNewListByServicePager(req *http.Request) (*http.Response, error) {
	if p.srv.NewListByServicePager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByServicePager not implemented")}
	}
	newListByServicePager := p.newListByServicePager.get(req)
	if newListByServicePager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/policyFragments`
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
		serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
		if err != nil {
			return nil, err
		}
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		orderbyUnescaped, err := url.QueryUnescape(qp.Get("$orderby"))
		if err != nil {
			return nil, err
		}
		orderbyParam := getOptional(orderbyUnescaped)
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
		var options *armapimanagement.PolicyFragmentClientListByServiceOptions
		if filterParam != nil || orderbyParam != nil || topParam != nil || skipParam != nil {
			options = &armapimanagement.PolicyFragmentClientListByServiceOptions{
				Filter:  filterParam,
				Orderby: orderbyParam,
				Top:     topParam,
				Skip:    skipParam,
			}
		}
		resp := p.srv.NewListByServicePager(resourceGroupNameParam, serviceNameParam, options)
		newListByServicePager = &resp
		p.newListByServicePager.add(req, newListByServicePager)
		server.PagerResponderInjectNextLinks(newListByServicePager, req, func(page *armapimanagement.PolicyFragmentClientListByServiceResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByServicePager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		p.newListByServicePager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByServicePager) {
		p.newListByServicePager.remove(req)
	}
	return resp, nil
}

func (p *PolicyFragmentServerTransport) dispatchListReferences(req *http.Request) (*http.Response, error) {
	if p.srv.ListReferences == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListReferences not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ApiManagement/service/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/policyFragments/(?P<id>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listReferences`
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
	idParam, err := url.PathUnescape(matches[regex.SubexpIndex("id")])
	if err != nil {
		return nil, err
	}
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
	var options *armapimanagement.PolicyFragmentClientListReferencesOptions
	if topParam != nil || skipParam != nil {
		options = &armapimanagement.PolicyFragmentClientListReferencesOptions{
			Top:  topParam,
			Skip: skipParam,
		}
	}
	respr, errRespr := p.srv.ListReferences(req.Context(), resourceGroupNameParam, serviceNameParam, idParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ResourceCollection, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to PolicyFragmentServerTransport
var policyFragmentServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
