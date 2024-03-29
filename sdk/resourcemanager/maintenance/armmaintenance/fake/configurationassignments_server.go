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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/maintenance/armmaintenance"
	"net/http"
	"net/url"
	"regexp"
)

// ConfigurationAssignmentsServer is a fake server for instances of the armmaintenance.ConfigurationAssignmentsClient type.
type ConfigurationAssignmentsServer struct {
	// CreateOrUpdate is the fake for method ConfigurationAssignmentsClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, providerName string, resourceType string, resourceName string, configurationAssignmentName string, configurationAssignment armmaintenance.ConfigurationAssignment, options *armmaintenance.ConfigurationAssignmentsClientCreateOrUpdateOptions) (resp azfake.Responder[armmaintenance.ConfigurationAssignmentsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// CreateOrUpdateParent is the fake for method ConfigurationAssignmentsClient.CreateOrUpdateParent
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdateParent func(ctx context.Context, resourceGroupName string, providerName string, resourceParentType string, resourceParentName string, resourceType string, resourceName string, configurationAssignmentName string, configurationAssignment armmaintenance.ConfigurationAssignment, options *armmaintenance.ConfigurationAssignmentsClientCreateOrUpdateParentOptions) (resp azfake.Responder[armmaintenance.ConfigurationAssignmentsClientCreateOrUpdateParentResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method ConfigurationAssignmentsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, providerName string, resourceType string, resourceName string, configurationAssignmentName string, options *armmaintenance.ConfigurationAssignmentsClientDeleteOptions) (resp azfake.Responder[armmaintenance.ConfigurationAssignmentsClientDeleteResponse], errResp azfake.ErrorResponder)

	// DeleteParent is the fake for method ConfigurationAssignmentsClient.DeleteParent
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	DeleteParent func(ctx context.Context, resourceGroupName string, providerName string, resourceParentType string, resourceParentName string, resourceType string, resourceName string, configurationAssignmentName string, options *armmaintenance.ConfigurationAssignmentsClientDeleteParentOptions) (resp azfake.Responder[armmaintenance.ConfigurationAssignmentsClientDeleteParentResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method ConfigurationAssignmentsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, providerName string, resourceType string, resourceName string, configurationAssignmentName string, options *armmaintenance.ConfigurationAssignmentsClientGetOptions) (resp azfake.Responder[armmaintenance.ConfigurationAssignmentsClientGetResponse], errResp azfake.ErrorResponder)

	// GetParent is the fake for method ConfigurationAssignmentsClient.GetParent
	// HTTP status codes to indicate success: http.StatusOK
	GetParent func(ctx context.Context, resourceGroupName string, providerName string, resourceParentType string, resourceParentName string, resourceType string, resourceName string, configurationAssignmentName string, options *armmaintenance.ConfigurationAssignmentsClientGetParentOptions) (resp azfake.Responder[armmaintenance.ConfigurationAssignmentsClientGetParentResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method ConfigurationAssignmentsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, providerName string, resourceType string, resourceName string, options *armmaintenance.ConfigurationAssignmentsClientListOptions) (resp azfake.PagerResponder[armmaintenance.ConfigurationAssignmentsClientListResponse])

	// NewListParentPager is the fake for method ConfigurationAssignmentsClient.NewListParentPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListParentPager func(resourceGroupName string, providerName string, resourceParentType string, resourceParentName string, resourceType string, resourceName string, options *armmaintenance.ConfigurationAssignmentsClientListParentOptions) (resp azfake.PagerResponder[armmaintenance.ConfigurationAssignmentsClientListParentResponse])
}

// NewConfigurationAssignmentsServerTransport creates a new instance of ConfigurationAssignmentsServerTransport with the provided implementation.
// The returned ConfigurationAssignmentsServerTransport instance is connected to an instance of armmaintenance.ConfigurationAssignmentsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewConfigurationAssignmentsServerTransport(srv *ConfigurationAssignmentsServer) *ConfigurationAssignmentsServerTransport {
	return &ConfigurationAssignmentsServerTransport{
		srv:                srv,
		newListPager:       newTracker[azfake.PagerResponder[armmaintenance.ConfigurationAssignmentsClientListResponse]](),
		newListParentPager: newTracker[azfake.PagerResponder[armmaintenance.ConfigurationAssignmentsClientListParentResponse]](),
	}
}

// ConfigurationAssignmentsServerTransport connects instances of armmaintenance.ConfigurationAssignmentsClient to instances of ConfigurationAssignmentsServer.
// Don't use this type directly, use NewConfigurationAssignmentsServerTransport instead.
type ConfigurationAssignmentsServerTransport struct {
	srv                *ConfigurationAssignmentsServer
	newListPager       *tracker[azfake.PagerResponder[armmaintenance.ConfigurationAssignmentsClientListResponse]]
	newListParentPager *tracker[azfake.PagerResponder[armmaintenance.ConfigurationAssignmentsClientListParentResponse]]
}

// Do implements the policy.Transporter interface for ConfigurationAssignmentsServerTransport.
func (c *ConfigurationAssignmentsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ConfigurationAssignmentsClient.CreateOrUpdate":
		resp, err = c.dispatchCreateOrUpdate(req)
	case "ConfigurationAssignmentsClient.CreateOrUpdateParent":
		resp, err = c.dispatchCreateOrUpdateParent(req)
	case "ConfigurationAssignmentsClient.Delete":
		resp, err = c.dispatchDelete(req)
	case "ConfigurationAssignmentsClient.DeleteParent":
		resp, err = c.dispatchDeleteParent(req)
	case "ConfigurationAssignmentsClient.Get":
		resp, err = c.dispatchGet(req)
	case "ConfigurationAssignmentsClient.GetParent":
		resp, err = c.dispatchGetParent(req)
	case "ConfigurationAssignmentsClient.NewListPager":
		resp, err = c.dispatchNewListPager(req)
	case "ConfigurationAssignmentsClient.NewListParentPager":
		resp, err = c.dispatchNewListParentPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *ConfigurationAssignmentsServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if c.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Maintenance/configurationAssignments/(?P<configurationAssignmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 6 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armmaintenance.ConfigurationAssignment](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
	if err != nil {
		return nil, err
	}
	resourceTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceType")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	configurationAssignmentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationAssignmentName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, providerNameParam, resourceTypeParam, resourceNameParam, configurationAssignmentNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ConfigurationAssignment, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ConfigurationAssignmentsServerTransport) dispatchCreateOrUpdateParent(req *http.Request) (*http.Response, error) {
	if c.srv.CreateOrUpdateParent == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdateParent not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceParentType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceParentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Maintenance/configurationAssignments/(?P<configurationAssignmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 8 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armmaintenance.ConfigurationAssignment](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
	if err != nil {
		return nil, err
	}
	resourceParentTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceParentType")])
	if err != nil {
		return nil, err
	}
	resourceParentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceParentName")])
	if err != nil {
		return nil, err
	}
	resourceTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceType")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	configurationAssignmentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationAssignmentName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.CreateOrUpdateParent(req.Context(), resourceGroupNameParam, providerNameParam, resourceParentTypeParam, resourceParentNameParam, resourceTypeParam, resourceNameParam, configurationAssignmentNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ConfigurationAssignment, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ConfigurationAssignmentsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if c.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Maintenance/configurationAssignments/(?P<configurationAssignmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 6 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
	if err != nil {
		return nil, err
	}
	resourceTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceType")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	configurationAssignmentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationAssignmentName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Delete(req.Context(), resourceGroupNameParam, providerNameParam, resourceTypeParam, resourceNameParam, configurationAssignmentNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusNoContent}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNoContent", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ConfigurationAssignment, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ConfigurationAssignmentsServerTransport) dispatchDeleteParent(req *http.Request) (*http.Response, error) {
	if c.srv.DeleteParent == nil {
		return nil, &nonRetriableError{errors.New("fake for method DeleteParent not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceParentType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceParentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Maintenance/configurationAssignments/(?P<configurationAssignmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 8 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
	if err != nil {
		return nil, err
	}
	resourceParentTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceParentType")])
	if err != nil {
		return nil, err
	}
	resourceParentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceParentName")])
	if err != nil {
		return nil, err
	}
	resourceTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceType")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	configurationAssignmentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationAssignmentName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.DeleteParent(req.Context(), resourceGroupNameParam, providerNameParam, resourceParentTypeParam, resourceParentNameParam, resourceTypeParam, resourceNameParam, configurationAssignmentNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusNoContent}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNoContent", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ConfigurationAssignment, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ConfigurationAssignmentsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if c.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Maintenance/configurationAssignments/(?P<configurationAssignmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 6 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
	if err != nil {
		return nil, err
	}
	resourceTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceType")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	configurationAssignmentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationAssignmentName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Get(req.Context(), resourceGroupNameParam, providerNameParam, resourceTypeParam, resourceNameParam, configurationAssignmentNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ConfigurationAssignment, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ConfigurationAssignmentsServerTransport) dispatchGetParent(req *http.Request) (*http.Response, error) {
	if c.srv.GetParent == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetParent not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceParentType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceParentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Maintenance/configurationAssignments/(?P<configurationAssignmentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 8 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
	if err != nil {
		return nil, err
	}
	resourceParentTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceParentType")])
	if err != nil {
		return nil, err
	}
	resourceParentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceParentName")])
	if err != nil {
		return nil, err
	}
	resourceTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceType")])
	if err != nil {
		return nil, err
	}
	resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
	if err != nil {
		return nil, err
	}
	configurationAssignmentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationAssignmentName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.GetParent(req.Context(), resourceGroupNameParam, providerNameParam, resourceParentTypeParam, resourceParentNameParam, resourceTypeParam, resourceNameParam, configurationAssignmentNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ConfigurationAssignment, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ConfigurationAssignmentsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := c.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Maintenance/configurationAssignments`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
		if err != nil {
			return nil, err
		}
		resourceTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceType")])
		if err != nil {
			return nil, err
		}
		resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListPager(resourceGroupNameParam, providerNameParam, resourceTypeParam, resourceNameParam, nil)
		newListPager = &resp
		c.newListPager.add(req, newListPager)
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		c.newListPager.remove(req)
	}
	return resp, nil
}

func (c *ConfigurationAssignmentsServerTransport) dispatchNewListParentPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListParentPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListParentPager not implemented")}
	}
	newListParentPager := c.newListParentPager.get(req)
	if newListParentPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceParentType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceParentName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceType>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/(?P<resourceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Maintenance/configurationAssignments`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 7 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
		if err != nil {
			return nil, err
		}
		resourceParentTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceParentType")])
		if err != nil {
			return nil, err
		}
		resourceParentNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceParentName")])
		if err != nil {
			return nil, err
		}
		resourceTypeParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceType")])
		if err != nil {
			return nil, err
		}
		resourceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceName")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListParentPager(resourceGroupNameParam, providerNameParam, resourceParentTypeParam, resourceParentNameParam, resourceTypeParam, resourceNameParam, nil)
		newListParentPager = &resp
		c.newListParentPager.add(req, newListParentPager)
	}
	resp, err := server.PagerResponderNext(newListParentPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListParentPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListParentPager) {
		c.newListParentPager.remove(req)
	}
	return resp, nil
}
