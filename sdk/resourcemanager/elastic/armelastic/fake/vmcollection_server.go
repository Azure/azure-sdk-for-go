//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Code generated by Microsoft (R) AutoRest Code Generator.Changes may cause incorrect behavior and will be lost if the code
// is regenerated.
// Code generated by @autorest/go. DO NOT EDIT.

package fake

import (
	"context"
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/elastic/armelastic"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

// VMCollectionServer is a fake server for instances of the armelastic.VMCollectionClient type.
type VMCollectionServer struct {
	// Update is the fake for method VMCollectionClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, monitorName string, options *armelastic.VMCollectionClientUpdateOptions) (resp azfake.Responder[armelastic.VMCollectionClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewVMCollectionServerTransport creates a new instance of VMCollectionServerTransport with the provided implementation.
// The returned VMCollectionServerTransport instance is connected to an instance of armelastic.VMCollectionClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewVMCollectionServerTransport(srv *VMCollectionServer) *VMCollectionServerTransport {
	return &VMCollectionServerTransport{srv: srv}
}

// VMCollectionServerTransport connects instances of armelastic.VMCollectionClient to instances of VMCollectionServer.
// Don't use this type directly, use NewVMCollectionServerTransport instead.
type VMCollectionServerTransport struct {
	srv *VMCollectionServer
}

// Do implements the policy.Transporter interface for VMCollectionServerTransport.
func (v *VMCollectionServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "VMCollectionClient.Update":
		resp, err = v.dispatchUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *VMCollectionServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if v.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Elastic/monitors/(?P<monitorName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/vmCollectionUpdate`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armelastic.VMCollectionUpdate](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	monitorNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("monitorName")])
	if err != nil {
		return nil, err
	}
	var options *armelastic.VMCollectionClientUpdateOptions
	if !reflect.ValueOf(body).IsZero() {
		options = &armelastic.VMCollectionClientUpdateOptions{
			Body: &body,
		}
	}
	respr, errRespr := v.srv.Update(req.Context(), resourceGroupNameParam, monitorNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
