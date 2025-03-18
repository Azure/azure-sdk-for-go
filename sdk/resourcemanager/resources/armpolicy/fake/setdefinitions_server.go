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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armpolicy"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// SetDefinitionsServer is a fake server for instances of the armpolicy.SetDefinitionsClient type.
type SetDefinitionsServer struct {
	// CreateOrUpdate is the fake for method SetDefinitionsClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, policySetDefinitionName string, parameters armpolicy.SetDefinition, options *armpolicy.SetDefinitionsClientCreateOrUpdateOptions) (resp azfake.Responder[armpolicy.SetDefinitionsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// CreateOrUpdateAtManagementGroup is the fake for method SetDefinitionsClient.CreateOrUpdateAtManagementGroup
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdateAtManagementGroup func(ctx context.Context, managementGroupID string, policySetDefinitionName string, parameters armpolicy.SetDefinition, options *armpolicy.SetDefinitionsClientCreateOrUpdateAtManagementGroupOptions) (resp azfake.Responder[armpolicy.SetDefinitionsClientCreateOrUpdateAtManagementGroupResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method SetDefinitionsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, policySetDefinitionName string, options *armpolicy.SetDefinitionsClientDeleteOptions) (resp azfake.Responder[armpolicy.SetDefinitionsClientDeleteResponse], errResp azfake.ErrorResponder)

	// DeleteAtManagementGroup is the fake for method SetDefinitionsClient.DeleteAtManagementGroup
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	DeleteAtManagementGroup func(ctx context.Context, managementGroupID string, policySetDefinitionName string, options *armpolicy.SetDefinitionsClientDeleteAtManagementGroupOptions) (resp azfake.Responder[armpolicy.SetDefinitionsClientDeleteAtManagementGroupResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method SetDefinitionsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, policySetDefinitionName string, options *armpolicy.SetDefinitionsClientGetOptions) (resp azfake.Responder[armpolicy.SetDefinitionsClientGetResponse], errResp azfake.ErrorResponder)

	// GetAtManagementGroup is the fake for method SetDefinitionsClient.GetAtManagementGroup
	// HTTP status codes to indicate success: http.StatusOK
	GetAtManagementGroup func(ctx context.Context, managementGroupID string, policySetDefinitionName string, options *armpolicy.SetDefinitionsClientGetAtManagementGroupOptions) (resp azfake.Responder[armpolicy.SetDefinitionsClientGetAtManagementGroupResponse], errResp azfake.ErrorResponder)

	// GetBuiltIn is the fake for method SetDefinitionsClient.GetBuiltIn
	// HTTP status codes to indicate success: http.StatusOK
	GetBuiltIn func(ctx context.Context, policySetDefinitionName string, options *armpolicy.SetDefinitionsClientGetBuiltInOptions) (resp azfake.Responder[armpolicy.SetDefinitionsClientGetBuiltInResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method SetDefinitionsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armpolicy.SetDefinitionsClientListOptions) (resp azfake.PagerResponder[armpolicy.SetDefinitionsClientListResponse])

	// NewListBuiltInPager is the fake for method SetDefinitionsClient.NewListBuiltInPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBuiltInPager func(options *armpolicy.SetDefinitionsClientListBuiltInOptions) (resp azfake.PagerResponder[armpolicy.SetDefinitionsClientListBuiltInResponse])

	// NewListByManagementGroupPager is the fake for method SetDefinitionsClient.NewListByManagementGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByManagementGroupPager func(managementGroupID string, options *armpolicy.SetDefinitionsClientListByManagementGroupOptions) (resp azfake.PagerResponder[armpolicy.SetDefinitionsClientListByManagementGroupResponse])
}

// NewSetDefinitionsServerTransport creates a new instance of SetDefinitionsServerTransport with the provided implementation.
// The returned SetDefinitionsServerTransport instance is connected to an instance of armpolicy.SetDefinitionsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSetDefinitionsServerTransport(srv *SetDefinitionsServer) *SetDefinitionsServerTransport {
	return &SetDefinitionsServerTransport{
		srv:                           srv,
		newListPager:                  newTracker[azfake.PagerResponder[armpolicy.SetDefinitionsClientListResponse]](),
		newListBuiltInPager:           newTracker[azfake.PagerResponder[armpolicy.SetDefinitionsClientListBuiltInResponse]](),
		newListByManagementGroupPager: newTracker[azfake.PagerResponder[armpolicy.SetDefinitionsClientListByManagementGroupResponse]](),
	}
}

// SetDefinitionsServerTransport connects instances of armpolicy.SetDefinitionsClient to instances of SetDefinitionsServer.
// Don't use this type directly, use NewSetDefinitionsServerTransport instead.
type SetDefinitionsServerTransport struct {
	srv                           *SetDefinitionsServer
	newListPager                  *tracker[azfake.PagerResponder[armpolicy.SetDefinitionsClientListResponse]]
	newListBuiltInPager           *tracker[azfake.PagerResponder[armpolicy.SetDefinitionsClientListBuiltInResponse]]
	newListByManagementGroupPager *tracker[azfake.PagerResponder[armpolicy.SetDefinitionsClientListByManagementGroupResponse]]
}

// Do implements the policy.Transporter interface for SetDefinitionsServerTransport.
func (s *SetDefinitionsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return s.dispatchToMethodFake(req, method)
}

func (s *SetDefinitionsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if setDefinitionsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = setDefinitionsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "SetDefinitionsClient.CreateOrUpdate":
				res.resp, res.err = s.dispatchCreateOrUpdate(req)
			case "SetDefinitionsClient.CreateOrUpdateAtManagementGroup":
				res.resp, res.err = s.dispatchCreateOrUpdateAtManagementGroup(req)
			case "SetDefinitionsClient.Delete":
				res.resp, res.err = s.dispatchDelete(req)
			case "SetDefinitionsClient.DeleteAtManagementGroup":
				res.resp, res.err = s.dispatchDeleteAtManagementGroup(req)
			case "SetDefinitionsClient.Get":
				res.resp, res.err = s.dispatchGet(req)
			case "SetDefinitionsClient.GetAtManagementGroup":
				res.resp, res.err = s.dispatchGetAtManagementGroup(req)
			case "SetDefinitionsClient.GetBuiltIn":
				res.resp, res.err = s.dispatchGetBuiltIn(req)
			case "SetDefinitionsClient.NewListPager":
				res.resp, res.err = s.dispatchNewListPager(req)
			case "SetDefinitionsClient.NewListBuiltInPager":
				res.resp, res.err = s.dispatchNewListBuiltInPager(req)
			case "SetDefinitionsClient.NewListByManagementGroupPager":
				res.resp, res.err = s.dispatchNewListByManagementGroupPager(req)
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

func (s *SetDefinitionsServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Authorization/policySetDefinitions/(?P<policySetDefinitionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armpolicy.SetDefinition](req)
	if err != nil {
		return nil, err
	}
	policySetDefinitionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("policySetDefinitionName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.CreateOrUpdate(req.Context(), policySetDefinitionNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SetDefinition, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SetDefinitionsServerTransport) dispatchCreateOrUpdateAtManagementGroup(req *http.Request) (*http.Response, error) {
	if s.srv.CreateOrUpdateAtManagementGroup == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdateAtManagementGroup not implemented")}
	}
	const regexStr = `/providers/Microsoft\.Management/managementGroups/(?P<managementGroupId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Authorization/policySetDefinitions/(?P<policySetDefinitionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armpolicy.SetDefinition](req)
	if err != nil {
		return nil, err
	}
	managementGroupIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("managementGroupId")])
	if err != nil {
		return nil, err
	}
	policySetDefinitionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("policySetDefinitionName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.CreateOrUpdateAtManagementGroup(req.Context(), managementGroupIDParam, policySetDefinitionNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SetDefinition, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SetDefinitionsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if s.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Authorization/policySetDefinitions/(?P<policySetDefinitionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	policySetDefinitionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("policySetDefinitionName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Delete(req.Context(), policySetDefinitionNameParam, nil)
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

func (s *SetDefinitionsServerTransport) dispatchDeleteAtManagementGroup(req *http.Request) (*http.Response, error) {
	if s.srv.DeleteAtManagementGroup == nil {
		return nil, &nonRetriableError{errors.New("fake for method DeleteAtManagementGroup not implemented")}
	}
	const regexStr = `/providers/Microsoft\.Management/managementGroups/(?P<managementGroupId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Authorization/policySetDefinitions/(?P<policySetDefinitionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	managementGroupIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("managementGroupId")])
	if err != nil {
		return nil, err
	}
	policySetDefinitionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("policySetDefinitionName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.DeleteAtManagementGroup(req.Context(), managementGroupIDParam, policySetDefinitionNameParam, nil)
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

func (s *SetDefinitionsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Authorization/policySetDefinitions/(?P<policySetDefinitionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	policySetDefinitionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("policySetDefinitionName")])
	if err != nil {
		return nil, err
	}
	expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
	if err != nil {
		return nil, err
	}
	expandParam := getOptional(expandUnescaped)
	var options *armpolicy.SetDefinitionsClientGetOptions
	if expandParam != nil {
		options = &armpolicy.SetDefinitionsClientGetOptions{
			Expand: expandParam,
		}
	}
	respr, errRespr := s.srv.Get(req.Context(), policySetDefinitionNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SetDefinition, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SetDefinitionsServerTransport) dispatchGetAtManagementGroup(req *http.Request) (*http.Response, error) {
	if s.srv.GetAtManagementGroup == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetAtManagementGroup not implemented")}
	}
	const regexStr = `/providers/Microsoft\.Management/managementGroups/(?P<managementGroupId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Authorization/policySetDefinitions/(?P<policySetDefinitionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	managementGroupIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("managementGroupId")])
	if err != nil {
		return nil, err
	}
	policySetDefinitionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("policySetDefinitionName")])
	if err != nil {
		return nil, err
	}
	expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
	if err != nil {
		return nil, err
	}
	expandParam := getOptional(expandUnescaped)
	var options *armpolicy.SetDefinitionsClientGetAtManagementGroupOptions
	if expandParam != nil {
		options = &armpolicy.SetDefinitionsClientGetAtManagementGroupOptions{
			Expand: expandParam,
		}
	}
	respr, errRespr := s.srv.GetAtManagementGroup(req.Context(), managementGroupIDParam, policySetDefinitionNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SetDefinition, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SetDefinitionsServerTransport) dispatchGetBuiltIn(req *http.Request) (*http.Response, error) {
	if s.srv.GetBuiltIn == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetBuiltIn not implemented")}
	}
	const regexStr = `/providers/Microsoft\.Authorization/policySetDefinitions/(?P<policySetDefinitionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	policySetDefinitionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("policySetDefinitionName")])
	if err != nil {
		return nil, err
	}
	expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
	if err != nil {
		return nil, err
	}
	expandParam := getOptional(expandUnescaped)
	var options *armpolicy.SetDefinitionsClientGetBuiltInOptions
	if expandParam != nil {
		options = &armpolicy.SetDefinitionsClientGetBuiltInOptions{
			Expand: expandParam,
		}
	}
	respr, errRespr := s.srv.GetBuiltIn(req.Context(), policySetDefinitionNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).SetDefinition, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SetDefinitionsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := s.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Authorization/policySetDefinitions`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
		if err != nil {
			return nil, err
		}
		expandParam := getOptional(expandUnescaped)
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
		var options *armpolicy.SetDefinitionsClientListOptions
		if filterParam != nil || expandParam != nil || topParam != nil {
			options = &armpolicy.SetDefinitionsClientListOptions{
				Filter: filterParam,
				Expand: expandParam,
				Top:    topParam,
			}
		}
		resp := s.srv.NewListPager(options)
		newListPager = &resp
		s.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armpolicy.SetDefinitionsClientListResponse, createLink func() string) {
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

func (s *SetDefinitionsServerTransport) dispatchNewListBuiltInPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListBuiltInPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBuiltInPager not implemented")}
	}
	newListBuiltInPager := s.newListBuiltInPager.get(req)
	if newListBuiltInPager == nil {
		qp := req.URL.Query()
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
		if err != nil {
			return nil, err
		}
		expandParam := getOptional(expandUnescaped)
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
		var options *armpolicy.SetDefinitionsClientListBuiltInOptions
		if filterParam != nil || expandParam != nil || topParam != nil {
			options = &armpolicy.SetDefinitionsClientListBuiltInOptions{
				Filter: filterParam,
				Expand: expandParam,
				Top:    topParam,
			}
		}
		resp := s.srv.NewListBuiltInPager(options)
		newListBuiltInPager = &resp
		s.newListBuiltInPager.add(req, newListBuiltInPager)
		server.PagerResponderInjectNextLinks(newListBuiltInPager, req, func(page *armpolicy.SetDefinitionsClientListBuiltInResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBuiltInPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListBuiltInPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBuiltInPager) {
		s.newListBuiltInPager.remove(req)
	}
	return resp, nil
}

func (s *SetDefinitionsServerTransport) dispatchNewListByManagementGroupPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByManagementGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByManagementGroupPager not implemented")}
	}
	newListByManagementGroupPager := s.newListByManagementGroupPager.get(req)
	if newListByManagementGroupPager == nil {
		const regexStr = `/providers/Microsoft\.Management/managementGroups/(?P<managementGroupId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Authorization/policySetDefinitions`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		managementGroupIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("managementGroupId")])
		if err != nil {
			return nil, err
		}
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		expandUnescaped, err := url.QueryUnescape(qp.Get("$expand"))
		if err != nil {
			return nil, err
		}
		expandParam := getOptional(expandUnescaped)
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
		var options *armpolicy.SetDefinitionsClientListByManagementGroupOptions
		if filterParam != nil || expandParam != nil || topParam != nil {
			options = &armpolicy.SetDefinitionsClientListByManagementGroupOptions{
				Filter: filterParam,
				Expand: expandParam,
				Top:    topParam,
			}
		}
		resp := s.srv.NewListByManagementGroupPager(managementGroupIDParam, options)
		newListByManagementGroupPager = &resp
		s.newListByManagementGroupPager.add(req, newListByManagementGroupPager)
		server.PagerResponderInjectNextLinks(newListByManagementGroupPager, req, func(page *armpolicy.SetDefinitionsClientListByManagementGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByManagementGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListByManagementGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByManagementGroupPager) {
		s.newListByManagementGroupPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to SetDefinitionsServerTransport
var setDefinitionsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
