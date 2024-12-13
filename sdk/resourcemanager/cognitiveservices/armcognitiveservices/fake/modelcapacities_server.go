//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cognitiveservices/armcognitiveservices"
	"net/http"
	"net/url"
	"regexp"
)

// ModelCapacitiesServer is a fake server for instances of the armcognitiveservices.ModelCapacitiesClient type.
type ModelCapacitiesServer struct {
	// NewListPager is the fake for method ModelCapacitiesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(modelFormat string, modelName string, modelVersion string, options *armcognitiveservices.ModelCapacitiesClientListOptions) (resp azfake.PagerResponder[armcognitiveservices.ModelCapacitiesClientListResponse])
}

// NewModelCapacitiesServerTransport creates a new instance of ModelCapacitiesServerTransport with the provided implementation.
// The returned ModelCapacitiesServerTransport instance is connected to an instance of armcognitiveservices.ModelCapacitiesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewModelCapacitiesServerTransport(srv *ModelCapacitiesServer) *ModelCapacitiesServerTransport {
	return &ModelCapacitiesServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armcognitiveservices.ModelCapacitiesClientListResponse]](),
	}
}

// ModelCapacitiesServerTransport connects instances of armcognitiveservices.ModelCapacitiesClient to instances of ModelCapacitiesServer.
// Don't use this type directly, use NewModelCapacitiesServerTransport instead.
type ModelCapacitiesServerTransport struct {
	srv          *ModelCapacitiesServer
	newListPager *tracker[azfake.PagerResponder[armcognitiveservices.ModelCapacitiesClientListResponse]]
}

// Do implements the policy.Transporter interface for ModelCapacitiesServerTransport.
func (m *ModelCapacitiesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ModelCapacitiesClient.NewListPager":
		resp, err = m.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *ModelCapacitiesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if m.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := m.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.CognitiveServices/modelCapacities`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		modelFormatParam, err := url.QueryUnescape(qp.Get("modelFormat"))
		if err != nil {
			return nil, err
		}
		modelNameParam, err := url.QueryUnescape(qp.Get("modelName"))
		if err != nil {
			return nil, err
		}
		modelVersionParam, err := url.QueryUnescape(qp.Get("modelVersion"))
		if err != nil {
			return nil, err
		}
		resp := m.srv.NewListPager(modelFormatParam, modelNameParam, modelVersionParam, nil)
		newListPager = &resp
		m.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armcognitiveservices.ModelCapacitiesClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		m.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		m.newListPager.remove(req)
	}
	return resp, nil
}