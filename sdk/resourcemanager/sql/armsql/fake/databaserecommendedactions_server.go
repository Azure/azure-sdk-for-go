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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/sql/armsql/v2"
	"net/http"
	"net/url"
	"regexp"
)

// DatabaseRecommendedActionsServer is a fake server for instances of the armsql.DatabaseRecommendedActionsClient type.
type DatabaseRecommendedActionsServer struct {
	// Get is the fake for method DatabaseRecommendedActionsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, serverName string, databaseName string, advisorName string, recommendedActionName string, options *armsql.DatabaseRecommendedActionsClientGetOptions) (resp azfake.Responder[armsql.DatabaseRecommendedActionsClientGetResponse], errResp azfake.ErrorResponder)

	// ListByDatabaseAdvisor is the fake for method DatabaseRecommendedActionsClient.ListByDatabaseAdvisor
	// HTTP status codes to indicate success: http.StatusOK
	ListByDatabaseAdvisor func(ctx context.Context, resourceGroupName string, serverName string, databaseName string, advisorName string, options *armsql.DatabaseRecommendedActionsClientListByDatabaseAdvisorOptions) (resp azfake.Responder[armsql.DatabaseRecommendedActionsClientListByDatabaseAdvisorResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method DatabaseRecommendedActionsClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, serverName string, databaseName string, advisorName string, recommendedActionName string, parameters armsql.RecommendedAction, options *armsql.DatabaseRecommendedActionsClientUpdateOptions) (resp azfake.Responder[armsql.DatabaseRecommendedActionsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewDatabaseRecommendedActionsServerTransport creates a new instance of DatabaseRecommendedActionsServerTransport with the provided implementation.
// The returned DatabaseRecommendedActionsServerTransport instance is connected to an instance of armsql.DatabaseRecommendedActionsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDatabaseRecommendedActionsServerTransport(srv *DatabaseRecommendedActionsServer) *DatabaseRecommendedActionsServerTransport {
	return &DatabaseRecommendedActionsServerTransport{srv: srv}
}

// DatabaseRecommendedActionsServerTransport connects instances of armsql.DatabaseRecommendedActionsClient to instances of DatabaseRecommendedActionsServer.
// Don't use this type directly, use NewDatabaseRecommendedActionsServerTransport instead.
type DatabaseRecommendedActionsServerTransport struct {
	srv *DatabaseRecommendedActionsServer
}

// Do implements the policy.Transporter interface for DatabaseRecommendedActionsServerTransport.
func (d *DatabaseRecommendedActionsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "DatabaseRecommendedActionsClient.Get":
		resp, err = d.dispatchGet(req)
	case "DatabaseRecommendedActionsClient.ListByDatabaseAdvisor":
		resp, err = d.dispatchListByDatabaseAdvisor(req)
	case "DatabaseRecommendedActionsClient.Update":
		resp, err = d.dispatchUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DatabaseRecommendedActionsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if d.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/advisors/(?P<advisorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/recommendedActions/(?P<recommendedActionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 6 {
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
	databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
	if err != nil {
		return nil, err
	}
	advisorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("advisorName")])
	if err != nil {
		return nil, err
	}
	recommendedActionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("recommendedActionName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.Get(req.Context(), resourceGroupNameParam, serverNameParam, databaseNameParam, advisorNameParam, recommendedActionNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).RecommendedAction, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DatabaseRecommendedActionsServerTransport) dispatchListByDatabaseAdvisor(req *http.Request) (*http.Response, error) {
	if d.srv.ListByDatabaseAdvisor == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListByDatabaseAdvisor not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/advisors/(?P<advisorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/recommendedActions`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 5 {
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
	databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
	if err != nil {
		return nil, err
	}
	advisorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("advisorName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.ListByDatabaseAdvisor(req.Context(), resourceGroupNameParam, serverNameParam, databaseNameParam, advisorNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).RecommendedActionArray, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DatabaseRecommendedActionsServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if d.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/servers/(?P<serverName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/databases/(?P<databaseName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/advisors/(?P<advisorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/recommendedActions/(?P<recommendedActionName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 6 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armsql.RecommendedAction](req)
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
	databaseNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("databaseName")])
	if err != nil {
		return nil, err
	}
	advisorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("advisorName")])
	if err != nil {
		return nil, err
	}
	recommendedActionNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("recommendedActionName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.Update(req.Context(), resourceGroupNameParam, serverNameParam, databaseNameParam, advisorNameParam, recommendedActionNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).RecommendedAction, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
