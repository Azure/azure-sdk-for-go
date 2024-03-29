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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/blockchain/armblockchain"
	"net/http"
	"regexp"
)

// SKUsServer is a fake server for instances of the armblockchain.SKUsClient type.
type SKUsServer struct {
	// List is the fake for method SKUsClient.List
	// HTTP status codes to indicate success: http.StatusOK
	List func(ctx context.Context, options *armblockchain.SKUsClientListOptions) (resp azfake.Responder[armblockchain.SKUsClientListResponse], errResp azfake.ErrorResponder)
}

// NewSKUsServerTransport creates a new instance of SKUsServerTransport with the provided implementation.
// The returned SKUsServerTransport instance is connected to an instance of armblockchain.SKUsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSKUsServerTransport(srv *SKUsServer) *SKUsServerTransport {
	return &SKUsServerTransport{srv: srv}
}

// SKUsServerTransport connects instances of armblockchain.SKUsClient to instances of SKUsServer.
// Don't use this type directly, use NewSKUsServerTransport instead.
type SKUsServerTransport struct {
	srv *SKUsServer
}

// Do implements the policy.Transporter interface for SKUsServerTransport.
func (s *SKUsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "SKUsClient.List":
		resp, err = s.dispatchList(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *SKUsServerTransport) dispatchList(req *http.Request) (*http.Response, error) {
	if s.srv.List == nil {
		return nil, &nonRetriableError{errors.New("fake for method List not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Blockchain/skus`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	respr, errRespr := s.srv.List(req.Context(), nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ResourceTypeSKUCollection, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
