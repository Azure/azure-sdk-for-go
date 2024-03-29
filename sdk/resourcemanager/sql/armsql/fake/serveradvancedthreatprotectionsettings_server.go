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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql/v2"
	"net/http"
	"net/url"
	"regexp"
)

// ServerAdvancedThreatProtectionSettingsServer is a fake server for instances of the armsql.ServerAdvancedThreatProtectionSettingsClient type.
type ServerAdvancedThreatProtectionSettingsServer struct {
	// BeginCreateOrUpdate is the fake for method ServerAdvancedThreatProtectionSettingsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, serverName string, advancedThreatProtectionName armsql.AdvancedThreatProtectionName, parameters armsql.ServerAdvancedThreatProtection, options *armsql.ServerAdvancedThreatProtectionSettingsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armsql.ServerAdvancedThreatProtectionSettingsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method ServerAdvancedThreatProtectionSettingsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serverName string, advancedThreatProtectionName armsql.AdvancedThreatProtectionName, options *armsql.ServerAdvancedThreatProtectionSettingsClientGetOptions) (resp azfake.Responder[armsql.ServerAdvancedThreatProtectionSettingsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByServerPager is the fake for method ServerAdvancedThreatProtectionSettingsClient.NewListByServerPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByServerPager func(resourceGroupName string, serverName string, options *armsql.ServerAdvancedThreatProtectionSettingsClientListByServerOptions) (resp azfake.PagerResponder[armsql.ServerAdvancedThreatProtectionSettingsClientListByServerResponse])
}

// NewServerAdvancedThreatProtectionSettingsServerTransport creates a new instance of ServerAdvancedThreatProtectionSettingsServerTransport with the provided implementation.
// The returned ServerAdvancedThreatProtectionSettingsServerTransport instance is connected to an instance of armsql.ServerAdvancedThreatProtectionSettingsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewServerAdvancedThreatProtectionSettingsServerTransport(srv *ServerAdvancedThreatProtectionSettingsServer) *ServerAdvancedThreatProtectionSettingsServerTransport {
	return &ServerAdvancedThreatProtectionSettingsServerTransport{
		srv:                  srv,
		beginCreateOrUpdate:  newTracker[azfake.PollerResponder[armsql.ServerAdvancedThreatProtectionSettingsClientCreateOrUpdateResponse]](),
		newListByServerPager: newTracker[azfake.PagerResponder[armsql.ServerAdvancedThreatProtectionSettingsClientListByServerResponse]](),
	}
}

// ServerAdvancedThreatProtectionSettingsServerTransport connects instances of armsql.ServerAdvancedThreatProtectionSettingsClient to instances of ServerAdvancedThreatProtectionSettingsServer.
// Don't use this type directly, use NewServerAdvancedThreatProtectionSettingsServerTransport instead.
type ServerAdvancedThreatProtectionSettingsServerTransport struct {
	srv                  *ServerAdvancedThreatProtectionSettingsServer
	beginCreateOrUpdate  *tracker[azfake.PollerResponder[armsql.ServerAdvancedThreatProtectionSettingsClientCreateOrUpdateResponse]]
	newListByServerPager *tracker[azfake.PagerResponder[armsql.ServerAdvancedThreatProtectionSettingsClientListByServerResponse]]
}

// Do implements the policy.Transporter interface for ServerAdvancedThreatProtectionSettingsServerTransport.
func (s *ServerAdvancedThreatProtectionSettingsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ServerAdvancedThreatProtectionSettingsClient.BeginCreateOrUpdate":
		resp, err = s.dispatchBeginCreateOrUpdate(req)
	case "ServerAdvancedThreatProtectionSettingsClient.Get":
		resp, err = s.dispatchGet(req)
	case "ServerAdvancedThreatProtectionSettingsClient.NewListByServerPager":
		resp, err = s.dispatchNewListByServerPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *ServerAdvancedThreatProtectionSettingsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := s.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/advancedThreatProtectionSettings/(?P<advancedThreatProtectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armsql.ServerAdvancedThreatProtection](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
		if err != nil {
			return nil, err
		}
		advancedThreatProtectionNameParam, err := parseWithCast(matches[regex.SubexpIndex("advancedThreatProtectionName")], func(v string) (armsql.AdvancedThreatProtectionName, error) {
			p, unescapeErr := url.PathUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armsql.AdvancedThreatProtectionName(p), nil
		})
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, serverNameParam, advancedThreatProtectionNameParam, body, nil)
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

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		s.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		s.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (s *ServerAdvancedThreatProtectionSettingsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/advancedThreatProtectionSettings/(?P<advancedThreatProtectionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
	if err != nil {
		return nil, err
	}
	advancedThreatProtectionNameParam, err := parseWithCast(matches[regex.SubexpIndex("advancedThreatProtectionName")], func(v string) (armsql.AdvancedThreatProtectionName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armsql.AdvancedThreatProtectionName(p), nil
	})
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), resourceGroupNameParam, serverNameParam, advancedThreatProtectionNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ServerAdvancedThreatProtection, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *ServerAdvancedThreatProtectionSettingsServerTransport) dispatchNewListByServerPager(req *http.Request) (*http.Response, error) {
	if s.srv.NewListByServerPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByServerPager not implemented")}
	}
	newListByServerPager := s.newListByServerPager.get(req)
	if newListByServerPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/advancedThreatProtectionSettings`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
		if err != nil {
			return nil, err
		}
		resp := s.srv.NewListByServerPager(resourceGroupNameParam, serverNameParam, nil)
		newListByServerPager = &resp
		s.newListByServerPager.add(req, newListByServerPager)
		server.PagerResponderInjectNextLinks(newListByServerPager, req, func(page *armsql.ServerAdvancedThreatProtectionSettingsClientListByServerResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByServerPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		s.newListByServerPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByServerPager) {
		s.newListByServerPager.remove(req)
	}
	return resp, nil
}
