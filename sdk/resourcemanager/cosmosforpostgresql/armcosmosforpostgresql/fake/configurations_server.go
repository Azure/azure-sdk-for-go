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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmosforpostgresql/armcosmosforpostgresql"
	"net/http"
	"net/url"
	"regexp"
)

// ConfigurationsServer is a fake server for instances of the armcosmosforpostgresql.ConfigurationsClient type.
type ConfigurationsServer struct {
	// Get is the fake for method ConfigurationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, clusterName string, configurationName string, options *armcosmosforpostgresql.ConfigurationsClientGetOptions) (resp azfake.Responder[armcosmosforpostgresql.ConfigurationsClientGetResponse], errResp azfake.ErrorResponder)

	// GetCoordinator is the fake for method ConfigurationsClient.GetCoordinator
	// HTTP status codes to indicate success: http.StatusOK
	GetCoordinator func(ctx context.Context, resourceGroupName string, clusterName string, configurationName string, options *armcosmosforpostgresql.ConfigurationsClientGetCoordinatorOptions) (resp azfake.Responder[armcosmosforpostgresql.ConfigurationsClientGetCoordinatorResponse], errResp azfake.ErrorResponder)

	// GetNode is the fake for method ConfigurationsClient.GetNode
	// HTTP status codes to indicate success: http.StatusOK
	GetNode func(ctx context.Context, resourceGroupName string, clusterName string, configurationName string, options *armcosmosforpostgresql.ConfigurationsClientGetNodeOptions) (resp azfake.Responder[armcosmosforpostgresql.ConfigurationsClientGetNodeResponse], errResp azfake.ErrorResponder)

	// NewListByClusterPager is the fake for method ConfigurationsClient.NewListByClusterPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByClusterPager func(resourceGroupName string, clusterName string, options *armcosmosforpostgresql.ConfigurationsClientListByClusterOptions) (resp azfake.PagerResponder[armcosmosforpostgresql.ConfigurationsClientListByClusterResponse])

	// NewListByServerPager is the fake for method ConfigurationsClient.NewListByServerPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByServerPager func(resourceGroupName string, clusterName string, serverName string, options *armcosmosforpostgresql.ConfigurationsClientListByServerOptions) (resp azfake.PagerResponder[armcosmosforpostgresql.ConfigurationsClientListByServerResponse])

	// BeginUpdateOnCoordinator is the fake for method ConfigurationsClient.BeginUpdateOnCoordinator
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginUpdateOnCoordinator func(ctx context.Context, resourceGroupName string, clusterName string, configurationName string, parameters armcosmosforpostgresql.ServerConfiguration, options *armcosmosforpostgresql.ConfigurationsClientBeginUpdateOnCoordinatorOptions) (resp azfake.PollerResponder[armcosmosforpostgresql.ConfigurationsClientUpdateOnCoordinatorResponse], errResp azfake.ErrorResponder)

	// BeginUpdateOnNode is the fake for method ConfigurationsClient.BeginUpdateOnNode
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginUpdateOnNode func(ctx context.Context, resourceGroupName string, clusterName string, configurationName string, parameters armcosmosforpostgresql.ServerConfiguration, options *armcosmosforpostgresql.ConfigurationsClientBeginUpdateOnNodeOptions) (resp azfake.PollerResponder[armcosmosforpostgresql.ConfigurationsClientUpdateOnNodeResponse], errResp azfake.ErrorResponder)
}

// NewConfigurationsServerTransport creates a new instance of ConfigurationsServerTransport with the provided implementation.
// The returned ConfigurationsServerTransport instance is connected to an instance of armcosmosforpostgresql.ConfigurationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewConfigurationsServerTransport(srv *ConfigurationsServer) *ConfigurationsServerTransport {
	return &ConfigurationsServerTransport{
		srv:                      srv,
		newListByClusterPager:    newTracker[azfake.PagerResponder[armcosmosforpostgresql.ConfigurationsClientListByClusterResponse]](),
		newListByServerPager:     newTracker[azfake.PagerResponder[armcosmosforpostgresql.ConfigurationsClientListByServerResponse]](),
		beginUpdateOnCoordinator: newTracker[azfake.PollerResponder[armcosmosforpostgresql.ConfigurationsClientUpdateOnCoordinatorResponse]](),
		beginUpdateOnNode:        newTracker[azfake.PollerResponder[armcosmosforpostgresql.ConfigurationsClientUpdateOnNodeResponse]](),
	}
}

// ConfigurationsServerTransport connects instances of armcosmosforpostgresql.ConfigurationsClient to instances of ConfigurationsServer.
// Don't use this type directly, use NewConfigurationsServerTransport instead.
type ConfigurationsServerTransport struct {
	srv                      *ConfigurationsServer
	newListByClusterPager    *tracker[azfake.PagerResponder[armcosmosforpostgresql.ConfigurationsClientListByClusterResponse]]
	newListByServerPager     *tracker[azfake.PagerResponder[armcosmosforpostgresql.ConfigurationsClientListByServerResponse]]
	beginUpdateOnCoordinator *tracker[azfake.PollerResponder[armcosmosforpostgresql.ConfigurationsClientUpdateOnCoordinatorResponse]]
	beginUpdateOnNode        *tracker[azfake.PollerResponder[armcosmosforpostgresql.ConfigurationsClientUpdateOnNodeResponse]]
}

// Do implements the policy.Transporter interface for ConfigurationsServerTransport.
func (c *ConfigurationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ConfigurationsClient.Get":
		resp, err = c.dispatchGet(req)
	case "ConfigurationsClient.GetCoordinator":
		resp, err = c.dispatchGetCoordinator(req)
	case "ConfigurationsClient.GetNode":
		resp, err = c.dispatchGetNode(req)
	case "ConfigurationsClient.NewListByClusterPager":
		resp, err = c.dispatchNewListByClusterPager(req)
	case "ConfigurationsClient.NewListByServerPager":
		resp, err = c.dispatchNewListByServerPager(req)
	case "ConfigurationsClient.BeginUpdateOnCoordinator":
		resp, err = c.dispatchBeginUpdateOnCoordinator(req)
	case "ConfigurationsClient.BeginUpdateOnNode":
		resp, err = c.dispatchBeginUpdateOnNode(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *ConfigurationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if c.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/serverGroupsv2/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
	if err != nil {
		return nil, err
	}
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Get(req.Context(), resourceGroupNameParam, clusterNameParam, configurationNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Configuration, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ConfigurationsServerTransport) dispatchGetCoordinator(req *http.Request) (*http.Response, error) {
	if c.srv.GetCoordinator == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetCoordinator not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/serverGroupsv2/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/coordinatorConfigurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
	if err != nil {
		return nil, err
	}
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.GetCoordinator(req.Context(), resourceGroupNameParam, clusterNameParam, configurationNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ServerConfiguration, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ConfigurationsServerTransport) dispatchGetNode(req *http.Request) (*http.Response, error) {
	if c.srv.GetNode == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetNode not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/serverGroupsv2/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/nodeConfigurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
	if err != nil {
		return nil, err
	}
	configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.GetNode(req.Context(), resourceGroupNameParam, clusterNameParam, configurationNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ServerConfiguration, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ConfigurationsServerTransport) dispatchNewListByClusterPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListByClusterPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByClusterPager not implemented")}
	}
	newListByClusterPager := c.newListByClusterPager.get(req)
	if newListByClusterPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/serverGroupsv2/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListByClusterPager(resourceGroupNameParam, clusterNameParam, nil)
		newListByClusterPager = &resp
		c.newListByClusterPager.add(req, newListByClusterPager)
		server.PagerResponderInjectNextLinks(newListByClusterPager, req, func(page *armcosmosforpostgresql.ConfigurationsClientListByClusterResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByClusterPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListByClusterPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByClusterPager) {
		c.newListByClusterPager.remove(req)
	}
	return resp, nil
}

func (c *ConfigurationsServerTransport) dispatchNewListByServerPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListByServerPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByServerPager not implemented")}
	}
	newListByServerPager := c.newListByServerPager.get(req)
	if newListByServerPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/serverGroupsv2/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/configurations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListByServerPager(resourceGroupNameParam, clusterNameParam, serverNameParam, nil)
		newListByServerPager = &resp
		c.newListByServerPager.add(req, newListByServerPager)
		server.PagerResponderInjectNextLinks(newListByServerPager, req, func(page *armcosmosforpostgresql.ConfigurationsClientListByServerResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByServerPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListByServerPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByServerPager) {
		c.newListByServerPager.remove(req)
	}
	return resp, nil
}

func (c *ConfigurationsServerTransport) dispatchBeginUpdateOnCoordinator(req *http.Request) (*http.Response, error) {
	if c.srv.BeginUpdateOnCoordinator == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdateOnCoordinator not implemented")}
	}
	beginUpdateOnCoordinator := c.beginUpdateOnCoordinator.get(req)
	if beginUpdateOnCoordinator == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/serverGroupsv2/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/coordinatorConfigurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcosmosforpostgresql.ServerConfiguration](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginUpdateOnCoordinator(req.Context(), resourceGroupNameParam, clusterNameParam, configurationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdateOnCoordinator = &respr
		c.beginUpdateOnCoordinator.add(req, beginUpdateOnCoordinator)
	}

	resp, err := server.PollerResponderNext(beginUpdateOnCoordinator, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		c.beginUpdateOnCoordinator.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdateOnCoordinator) {
		c.beginUpdateOnCoordinator.remove(req)
	}

	return resp, nil
}

func (c *ConfigurationsServerTransport) dispatchBeginUpdateOnNode(req *http.Request) (*http.Response, error) {
	if c.srv.BeginUpdateOnNode == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdateOnNode not implemented")}
	}
	beginUpdateOnNode := c.beginUpdateOnNode.get(req)
	if beginUpdateOnNode == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforPostgreSQL/serverGroupsv2/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/nodeConfigurations/(?P<configurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcosmosforpostgresql.ServerConfiguration](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		configurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("configurationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginUpdateOnNode(req.Context(), resourceGroupNameParam, clusterNameParam, configurationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdateOnNode = &respr
		c.beginUpdateOnNode.add(req, beginUpdateOnNode)
	}

	resp, err := server.PollerResponderNext(beginUpdateOnNode, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		c.beginUpdateOnNode.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdateOnNode) {
		c.beginUpdateOnNode.remove(req)
	}

	return resp, nil
}
