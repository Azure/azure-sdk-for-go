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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/support/armsupport/v2"
	"net/http"
)

// LookUpResourceIDServer is a fake server for instances of the armsupport.LookUpResourceIDClient type.
type LookUpResourceIDServer struct {
	// Post is the fake for method LookUpResourceIDClient.Post
	// HTTP status codes to indicate success: http.StatusOK
	Post func(ctx context.Context, lookUpResourceIDRequest armsupport.LookUpResourceIDRequest, options *armsupport.LookUpResourceIDClientPostOptions) (resp azfake.Responder[armsupport.LookUpResourceIDClientPostResponse], errResp azfake.ErrorResponder)
}

// NewLookUpResourceIDServerTransport creates a new instance of LookUpResourceIDServerTransport with the provided implementation.
// The returned LookUpResourceIDServerTransport instance is connected to an instance of armsupport.LookUpResourceIDClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewLookUpResourceIDServerTransport(srv *LookUpResourceIDServer) *LookUpResourceIDServerTransport {
	return &LookUpResourceIDServerTransport{srv: srv}
}

// LookUpResourceIDServerTransport connects instances of armsupport.LookUpResourceIDClient to instances of LookUpResourceIDServer.
// Don't use this type directly, use NewLookUpResourceIDServerTransport instead.
type LookUpResourceIDServerTransport struct {
	srv *LookUpResourceIDServer
}

// Do implements the policy.Transporter interface for LookUpResourceIDServerTransport.
func (l *LookUpResourceIDServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "LookUpResourceIDClient.Post":
		resp, err = l.dispatchPost(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (l *LookUpResourceIDServerTransport) dispatchPost(req *http.Request) (*http.Response, error) {
	if l.srv.Post == nil {
		return nil, &nonRetriableError{errors.New("fake for method Post not implemented")}
	}
	body, err := server.UnmarshalRequestAsJSON[armsupport.LookUpResourceIDRequest](req)
	if err != nil {
		return nil, err
	}
	respr, errRespr := l.srv.Post(req.Context(), body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).LookUpResourceIDResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
