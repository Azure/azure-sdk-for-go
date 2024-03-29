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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/security/armsecurity"
	"net/http"
	"net/url"
	"regexp"
)

// GitLabSubgroupsServer is a fake server for instances of the armsecurity.GitLabSubgroupsClient type.
type GitLabSubgroupsServer struct {
	// List is the fake for method GitLabSubgroupsClient.List
	// HTTP status codes to indicate success: http.StatusOK
	List func(ctx context.Context, resourceGroupName string, securityConnectorName string, groupFQName string, options *armsecurity.GitLabSubgroupsClientListOptions) (resp azfake.Responder[armsecurity.GitLabSubgroupsClientListResponse], errResp azfake.ErrorResponder)
}

// NewGitLabSubgroupsServerTransport creates a new instance of GitLabSubgroupsServerTransport with the provided implementation.
// The returned GitLabSubgroupsServerTransport instance is connected to an instance of armsecurity.GitLabSubgroupsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewGitLabSubgroupsServerTransport(srv *GitLabSubgroupsServer) *GitLabSubgroupsServerTransport {
	return &GitLabSubgroupsServerTransport{srv: srv}
}

// GitLabSubgroupsServerTransport connects instances of armsecurity.GitLabSubgroupsClient to instances of GitLabSubgroupsServer.
// Don't use this type directly, use NewGitLabSubgroupsServerTransport instead.
type GitLabSubgroupsServerTransport struct {
	srv *GitLabSubgroupsServer
}

// Do implements the policy.Transporter interface for GitLabSubgroupsServerTransport.
func (g *GitLabSubgroupsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "GitLabSubgroupsClient.List":
		resp, err = g.dispatchList(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *GitLabSubgroupsServerTransport) dispatchList(req *http.Request) (*http.Response, error) {
	if g.srv.List == nil {
		return nil, &nonRetriableError{errors.New("fake for method List not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Security/securityConnectors/(?P<securityConnectorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/devops/default/gitLabGroups/(?P<groupFQName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listSubgroups`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	securityConnectorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("securityConnectorName")])
	if err != nil {
		return nil, err
	}
	groupFQNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("groupFQName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.List(req.Context(), resourceGroupNameParam, securityConnectorNameParam, groupFQNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).GitLabGroupListResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
