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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v3"
	"net/http"
	"regexp"
)

// ManagerCommitsServer is a fake server for instances of the armnetwork.ManagerCommitsClient type.
type ManagerCommitsServer struct {
	// BeginPost is the fake for method ManagerCommitsClient.BeginPost
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginPost func(ctx context.Context, resourceGroupName string, networkManagerName string, parameters armnetwork.ManagerCommit, options *armnetwork.ManagerCommitsClientBeginPostOptions) (resp azfake.PollerResponder[armnetwork.ManagerCommitsClientPostResponse], errResp azfake.ErrorResponder)
}

// NewManagerCommitsServerTransport creates a new instance of ManagerCommitsServerTransport with the provided implementation.
// The returned ManagerCommitsServerTransport instance is connected to an instance of armnetwork.ManagerCommitsClient by way of the
// undefined.Transporter field.
func NewManagerCommitsServerTransport(srv *ManagerCommitsServer) *ManagerCommitsServerTransport {
	return &ManagerCommitsServerTransport{srv: srv}
}

// ManagerCommitsServerTransport connects instances of armnetwork.ManagerCommitsClient to instances of ManagerCommitsServer.
// Don't use this type directly, use NewManagerCommitsServerTransport instead.
type ManagerCommitsServerTransport struct {
	srv       *ManagerCommitsServer
	beginPost *azfake.PollerResponder[armnetwork.ManagerCommitsClientPostResponse]
}

// Do implements the policy.Transporter interface for ManagerCommitsServerTransport.
func (m *ManagerCommitsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ManagerCommitsClient.BeginPost":
		resp, err = m.dispatchBeginPost(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *ManagerCommitsServerTransport) dispatchBeginPost(req *http.Request) (*http.Response, error) {
	if m.srv.BeginPost == nil {
		return nil, &nonRetriableError{errors.New("method BeginPost not implemented")}
	}
	if m.beginPost == nil {
		const regexStr = "/subscriptions/(?P<subscriptionId>[a-zA-Z0-9-_]+)/resourceGroups/(?P<resourceGroupName>[a-zA-Z0-9-_]+)/providers/Microsoft.Network/networkManagers/(?P<networkManagerName>[a-zA-Z0-9-_]+)/commit"
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.Path)
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.ManagerCommit](req)
		if err != nil {
			return nil, err
		}
		respr, errRespr := m.srv.BeginPost(req.Context(), matches[regex.SubexpIndex("resourceGroupName")], matches[regex.SubexpIndex("networkManagerName")], body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		m.beginPost = &respr
	}

	resp, err := server.PollerResponderNext(m.beginPost, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(m.beginPost) {
		m.beginPost = nil
	}

	return resp, nil
}
