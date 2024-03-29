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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mariadb/armmariadb"
	"net/http"
	"net/url"
	"regexp"
)

// ManagementServer is a fake server for instances of the armmariadb.ManagementClient type.
type ManagementServer struct {
	// BeginCreateRecommendedActionSession is the fake for method ManagementClient.BeginCreateRecommendedActionSession
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreateRecommendedActionSession func(ctx context.Context, resourceGroupName string, serverName string, advisorName string, databaseName string, options *armmariadb.ManagementClientBeginCreateRecommendedActionSessionOptions) (resp azfake.PollerResponder[armmariadb.ManagementClientCreateRecommendedActionSessionResponse], errResp azfake.ErrorResponder)

	// ResetQueryPerformanceInsightData is the fake for method ManagementClient.ResetQueryPerformanceInsightData
	// HTTP status codes to indicate success: http.StatusOK
	ResetQueryPerformanceInsightData func(ctx context.Context, resourceGroupName string, serverName string, options *armmariadb.ManagementClientResetQueryPerformanceInsightDataOptions) (resp azfake.Responder[armmariadb.ManagementClientResetQueryPerformanceInsightDataResponse], errResp azfake.ErrorResponder)
}

// NewManagementServerTransport creates a new instance of ManagementServerTransport with the provided implementation.
// The returned ManagementServerTransport instance is connected to an instance of armmariadb.ManagementClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewManagementServerTransport(srv *ManagementServer) *ManagementServerTransport {
	return &ManagementServerTransport{
		srv:                                 srv,
		beginCreateRecommendedActionSession: newTracker[azfake.PollerResponder[armmariadb.ManagementClientCreateRecommendedActionSessionResponse]](),
	}
}

// ManagementServerTransport connects instances of armmariadb.ManagementClient to instances of ManagementServer.
// Don't use this type directly, use NewManagementServerTransport instead.
type ManagementServerTransport struct {
	srv                                 *ManagementServer
	beginCreateRecommendedActionSession *tracker[azfake.PollerResponder[armmariadb.ManagementClientCreateRecommendedActionSessionResponse]]
}

// Do implements the policy.Transporter interface for ManagementServerTransport.
func (m *ManagementServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ManagementClient.BeginCreateRecommendedActionSession":
		resp, err = m.dispatchBeginCreateRecommendedActionSession(req)
	case "ManagementClient.ResetQueryPerformanceInsightData":
		resp, err = m.dispatchResetQueryPerformanceInsightData(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *ManagementServerTransport) dispatchBeginCreateRecommendedActionSession(req *http.Request) (*http.Response, error) {
	if m.srv.BeginCreateRecommendedActionSession == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateRecommendedActionSession not implemented")}
	}
	beginCreateRecommendedActionSession := m.beginCreateRecommendedActionSession.get(req)
	if beginCreateRecommendedActionSession == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforMariaDB/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/advisors/(?P<advisorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/createRecommendedActionSession`
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
		serverNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("serverName")])
		if err != nil {
			return nil, err
		}
		advisorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("advisorName")])
		if err != nil {
			return nil, err
		}
		databaseNameParam, err := url.QueryUnescape(qp.Get("databaseName"))
		if err != nil {
			return nil, err
		}
		respr, errRespr := m.srv.BeginCreateRecommendedActionSession(req.Context(), resourceGroupNameParam, serverNameParam, advisorNameParam, databaseNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateRecommendedActionSession = &respr
		m.beginCreateRecommendedActionSession.add(req, beginCreateRecommendedActionSession)
	}

	resp, err := server.PollerResponderNext(beginCreateRecommendedActionSession, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		m.beginCreateRecommendedActionSession.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateRecommendedActionSession) {
		m.beginCreateRecommendedActionSession.remove(req)
	}

	return resp, nil
}

func (m *ManagementServerTransport) dispatchResetQueryPerformanceInsightData(req *http.Request) (*http.Response, error) {
	if m.srv.ResetQueryPerformanceInsightData == nil {
		return nil, &nonRetriableError{errors.New("fake for method ResetQueryPerformanceInsightData not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DBforMariaDB/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resetQueryPerformanceInsightData`
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
	respr, errRespr := m.srv.ResetQueryPerformanceInsightData(req.Context(), resourceGroupNameParam, serverNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).QueryPerformanceInsightResetDataResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
