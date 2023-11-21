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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/deploymentmanager/armdeploymentmanager"
	"net/http"
	"net/url"
	"regexp"
)

// ServiceUnitsServer is a fake server for instances of the armdeploymentmanager.ServiceUnitsClient type.
type ServiceUnitsServer struct {
	// BeginCreateOrUpdate is the fake for method ServiceUnitsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, serviceTopologyName string, serviceName string, serviceUnitName string, serviceUnitInfo armdeploymentmanager.ServiceUnitResource, options *armdeploymentmanager.ServiceUnitsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armdeploymentmanager.ServiceUnitsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method ServiceUnitsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, serviceTopologyName string, serviceName string, serviceUnitName string, options *armdeploymentmanager.ServiceUnitsClientDeleteOptions) (resp azfake.Responder[armdeploymentmanager.ServiceUnitsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method ServiceUnitsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serviceTopologyName string, serviceName string, serviceUnitName string, options *armdeploymentmanager.ServiceUnitsClientGetOptions) (resp azfake.Responder[armdeploymentmanager.ServiceUnitsClientGetResponse], errResp azfake.ErrorResponder)

	// List is the fake for method ServiceUnitsClient.List
	// HTTP status codes to indicate success: http.StatusOK
	List func(ctx context.Context, resourceGroupName string, serviceTopologyName string, serviceName string, options *armdeploymentmanager.ServiceUnitsClientListOptions) (resp azfake.Responder[armdeploymentmanager.ServiceUnitsClientListResponse], errResp azfake.ErrorResponder)
}

// NewServiceUnitsServerTransport creates a new instance of ServiceUnitsServerTransport with the provided implementation.
// The returned ServiceUnitsServerTransport instance is connected to an instance of armdeploymentmanager.ServiceUnitsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewServiceUnitsServerTransport(srv *ServiceUnitsServer) *ServiceUnitsServerTransport {
	return &ServiceUnitsServerTransport{
		srv:                 srv,
		beginCreateOrUpdate: newTracker[azfake.PollerResponder[armdeploymentmanager.ServiceUnitsClientCreateOrUpdateResponse]](),
	}
}

// ServiceUnitsServerTransport connects instances of armdeploymentmanager.ServiceUnitsClient to instances of ServiceUnitsServer.
// Don't use this type directly, use NewServiceUnitsServerTransport instead.
type ServiceUnitsServerTransport struct {
	srv                 *ServiceUnitsServer
	beginCreateOrUpdate *tracker[azfake.PollerResponder[armdeploymentmanager.ServiceUnitsClientCreateOrUpdateResponse]]
}

// Do implements the policy.Transporter interface for ServiceUnitsServerTransport.
func (s *ServiceUnitsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ServiceUnitsClient.BeginCreateOrUpdate":
		resp, err = s.dispatchBeginCreateOrUpdate(req)
	case "ServiceUnitsClient.Delete":
		resp, err = s.dispatchDelete(req)
	case "ServiceUnitsClient.Get":
		resp, err = s.dispatchGet(req)
	case "ServiceUnitsClient.List":
		resp, err = s.dispatchList(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ServiceUnitsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := s.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DeploymentManager/serviceTopologies/(?P<serviceTopologyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/services/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/serviceUnits/(?P<serviceUnitName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armdeploymentmanager.ServiceUnitResource](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serviceTopologyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceTopologyName")])
		if err != nil {
			return nil, err
		}
		serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
		if err != nil {
			return nil, err
		}
		serviceUnitNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceUnitName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, serviceTopologyNameParam, serviceNameParam, serviceUnitNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		s.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusCreated}, resp.StatusCode) {
		s.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		s.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (s *ServiceUnitsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if s.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DeploymentManager/serviceTopologies/(?P<serviceTopologyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/services/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/serviceUnits/(?P<serviceUnitName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceTopologyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceTopologyName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	serviceUnitNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceUnitName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Delete(req.Context(), resourceGroupNameParam, serviceTopologyNameParam, serviceNameParam, serviceUnitNameParam, nil)
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

func (s *ServiceUnitsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DeploymentManager/serviceTopologies/(?P<serviceTopologyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/services/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/serviceUnits/(?P<serviceUnitName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceTopologyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceTopologyName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	serviceUnitNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceUnitName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), resourceGroupNameParam, serviceTopologyNameParam, serviceNameParam, serviceUnitNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ServiceUnitResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ServiceUnitsServerTransport) dispatchList(req *http.Request) (*http.Response, error) {
	if s.srv.List == nil {
		return nil, &nonRetriableError{errors.New("fake for method List not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DeploymentManager/serviceTopologies/(?P<serviceTopologyName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/services/(?P<serviceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/serviceUnits`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serviceTopologyNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceTopologyName")])
	if err != nil {
		return nil, err
	}
	serviceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serviceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.List(req.Context(), resourceGroupNameParam, serviceTopologyNameParam, serviceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ServiceUnitResourceArray, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
