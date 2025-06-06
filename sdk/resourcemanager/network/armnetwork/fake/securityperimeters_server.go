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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// SecurityPerimetersServer is a fake server for instances of the armnetwork.SecurityPerimetersClient type.
type SecurityPerimetersServer struct {
	// CreateOrUpdate is the fake for method SecurityPerimetersClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, networkSecurityPerimeterName string, parameters armnetwork.SecurityPerimeter, options *armnetwork.SecurityPerimetersClientCreateOrUpdateOptions) (resp azfake.Responder[armnetwork.SecurityPerimetersClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method SecurityPerimetersClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, networkSecurityPerimeterName string, options *armnetwork.SecurityPerimetersClientBeginDeleteOptions) (resp azfake.PollerResponder[armnetwork.SecurityPerimetersClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method SecurityPerimetersClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, networkSecurityPerimeterName string, options *armnetwork.SecurityPerimetersClientGetOptions) (resp azfake.Responder[armnetwork.SecurityPerimetersClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method SecurityPerimetersClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, options *armnetwork.SecurityPerimetersClientListOptions) (resp azfake.PagerResponder[armnetwork.SecurityPerimetersClientListResponse])

	// NewListBySubscriptionPager is the fake for method SecurityPerimetersClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armnetwork.SecurityPerimetersClientListBySubscriptionOptions) (resp azfake.PagerResponder[armnetwork.SecurityPerimetersClientListBySubscriptionResponse])

	// Patch is the fake for method SecurityPerimetersClient.Patch
	// HTTP status codes to indicate success: http.StatusOK
	Patch func(ctx context.Context, resourceGroupName string, networkSecurityPerimeterName string, parameters armnetwork.UpdateTagsRequest, options *armnetwork.SecurityPerimetersClientPatchOptions) (resp azfake.Responder[armnetwork.SecurityPerimetersClientPatchResponse], errResp azfake.ErrorResponder)
}

// NewSecurityPerimetersServerTransport creates a new instance of SecurityPerimetersServerTransport with the provided implementation.
// The returned SecurityPerimetersServerTransport instance is connected to an instance of armnetwork.SecurityPerimetersClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSecurityPerimetersServerTransport(srv *SecurityPerimetersServer) *SecurityPerimetersServerTransport {
	return &SecurityPerimetersServerTransport{
		srv:                        srv,
		beginDelete:                newTracker[azfake.PollerResponder[armnetwork.SecurityPerimetersClientDeleteResponse]](),
		newListPager:               newTracker[azfake.PagerResponder[armnetwork.SecurityPerimetersClientListResponse]](),
		newListBySubscriptionPager: newTracker[azfake.PagerResponder[armnetwork.SecurityPerimetersClientListBySubscriptionResponse]](),
	}
}

// SecurityPerimetersServerTransport connects instances of armnetwork.SecurityPerimetersClient to instances of SecurityPerimetersServer.
// Don't use this type directly, use NewSecurityPerimetersServerTransport instead.
type SecurityPerimetersServerTransport struct {
	srv                        *SecurityPerimetersServer
	beginDelete                *tracker[azfake.PollerResponder[armnetwork.SecurityPerimetersClientDeleteResponse]]
	newListPager               *tracker[azfake.PagerResponder[armnetwork.SecurityPerimetersClientListResponse]]
	newListBySubscriptionPager *tracker[azfake.PagerResponder[armnetwork.SecurityPerimetersClientListBySubscriptionResponse]]
}

// Do implements the policy.Transporter interface for SecurityPerimetersServerTransport.
func (s *SecurityPerimetersServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return s.dispatchToMethodFake(req, method)
}

func (s *SecurityPerimetersServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if securityPerimetersServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = securityPerimetersServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "SecurityPerimetersClient.CreateOrUpdate":
				res.resp, res.err = s.dispatchCreateOrUpdate(req)
			case "SecurityPerimetersClient.BeginDelete":
				res.resp, res.err = s.dispatchBeginDelete(req)
			case "SecurityPerimetersClient.Get":
				res.resp, res.err = s.dispatchGet(req)
			case "SecurityPerimetersClient.NewListPager":
				res.resp, res.err = s.dispatchNewListPager(req)
			case "SecurityPerimetersClient.NewListBySubscriptionPager":
				res.resp, res.err = s.dispatchNewListBySubscriptionPager(req)
			case "SecurityPerimetersClient.Patch":
				res.resp, res.err = s.dispatchPatch(req)
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

func (s *SecurityPerimetersServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkSecurityPerimeters/(?P<networkSecurityPerimeterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armnetwork.SecurityPerimeter](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	networkSecurityPerimeterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkSecurityPerimeterName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, networkSecurityPerimeterNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SecurityPerimeter, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SecurityPerimetersServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if s.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := s.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkSecurityPerimeters/(?P<networkSecurityPerimeterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
		networkSecurityPerimeterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkSecurityPerimeterName")])
		if err != nil {
			return nil, err
		}
		forceDeletionUnescaped, err := url.QueryUnescape(qp.Get("forceDeletion"))
		if err != nil {
			return nil, err
		}
		forceDeletionParam, err := parseOptional(forceDeletionUnescaped, strconv.ParseBool)
		if err != nil {
			return nil, err
		}
		var options *armnetwork.SecurityPerimetersClientBeginDeleteOptions
		if forceDeletionParam != nil {
			options = &armnetwork.SecurityPerimetersClientBeginDeleteOptions{
				ForceDeletion: forceDeletionParam,
			}
		}
		respr, errRespr := s.srv.BeginDelete(req.Context(), resourceGroupNameParam, networkSecurityPerimeterNameParam, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		s.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		s.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		s.beginDelete.remove(req)
	}

	return resp, nil
}

func (s *SecurityPerimetersServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkSecurityPerimeters/(?P<networkSecurityPerimeterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	networkSecurityPerimeterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkSecurityPerimeterName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), resourceGroupNameParam, networkSecurityPerimeterNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SecurityPerimeter, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SecurityPerimetersServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := s.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkSecurityPerimeters`
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
		skipTokenUnescaped, err := url.QueryUnescape(qp.Get("$skipToken"))
		if err != nil {
			return nil, err
		}
		skipTokenParam := getOptional(skipTokenUnescaped)
		var options *armnetwork.SecurityPerimetersClientListOptions
		if topParam != nil || skipTokenParam != nil {
			options = &armnetwork.SecurityPerimetersClientListOptions{
				Top:       topParam,
				SkipToken: skipTokenParam,
			}
		}
		resp := s.srv.NewListPager(resourceGroupNameParam, options)
		newListPager = &resp
		s.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armnetwork.SecurityPerimetersClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		s.newListPager.remove(req)
	}
	return resp, nil
}

func (s *SecurityPerimetersServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := s.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkSecurityPerimeters`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
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
		skipTokenUnescaped, err := url.QueryUnescape(qp.Get("$skipToken"))
		if err != nil {
			return nil, err
		}
		skipTokenParam := getOptional(skipTokenUnescaped)
		var options *armnetwork.SecurityPerimetersClientListBySubscriptionOptions
		if topParam != nil || skipTokenParam != nil {
			options = &armnetwork.SecurityPerimetersClientListBySubscriptionOptions{
				Top:       topParam,
				SkipToken: skipTokenParam,
			}
		}
		resp := s.srv.NewListBySubscriptionPager(options)
		newListBySubscriptionPager = &resp
		s.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armnetwork.SecurityPerimetersClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		s.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (s *SecurityPerimetersServerTransport) dispatchPatch(req *http.Request) (*http.Response, error) {
	if s.srv.Patch == nil {
		return nil, &nonRetriableError{errors.New("fake for method Patch not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkSecurityPerimeters/(?P<networkSecurityPerimeterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armnetwork.UpdateTagsRequest](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	networkSecurityPerimeterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkSecurityPerimeterName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Patch(req.Context(), resourceGroupNameParam, networkSecurityPerimeterNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SecurityPerimeter, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to SecurityPerimetersServerTransport
var securityPerimetersServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
