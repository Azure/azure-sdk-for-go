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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/costmanagement/armcostmanagement/v2"
	"net/http"
	"net/url"
	"regexp"
)

// GenerateDetailedCostReportServer is a fake server for instances of the armcostmanagement.GenerateDetailedCostReportClient type.
type GenerateDetailedCostReportServer struct {
	// BeginCreateOperation is the fake for method GenerateDetailedCostReportClient.BeginCreateOperation
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreateOperation func(ctx context.Context, scope string, parameters armcostmanagement.GenerateDetailedCostReportDefinition, options *armcostmanagement.GenerateDetailedCostReportClientBeginCreateOperationOptions) (resp azfake.PollerResponder[armcostmanagement.GenerateDetailedCostReportClientCreateOperationResponse], errResp azfake.ErrorResponder)
}

// NewGenerateDetailedCostReportServerTransport creates a new instance of GenerateDetailedCostReportServerTransport with the provided implementation.
// The returned GenerateDetailedCostReportServerTransport instance is connected to an instance of armcostmanagement.GenerateDetailedCostReportClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewGenerateDetailedCostReportServerTransport(srv *GenerateDetailedCostReportServer) *GenerateDetailedCostReportServerTransport {
	return &GenerateDetailedCostReportServerTransport{
		srv:                  srv,
		beginCreateOperation: newTracker[azfake.PollerResponder[armcostmanagement.GenerateDetailedCostReportClientCreateOperationResponse]](),
	}
}

// GenerateDetailedCostReportServerTransport connects instances of armcostmanagement.GenerateDetailedCostReportClient to instances of GenerateDetailedCostReportServer.
// Don't use this type directly, use NewGenerateDetailedCostReportServerTransport instead.
type GenerateDetailedCostReportServerTransport struct {
	srv                  *GenerateDetailedCostReportServer
	beginCreateOperation *tracker[azfake.PollerResponder[armcostmanagement.GenerateDetailedCostReportClientCreateOperationResponse]]
}

// Do implements the policy.Transporter interface for GenerateDetailedCostReportServerTransport.
func (g *GenerateDetailedCostReportServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "GenerateDetailedCostReportClient.BeginCreateOperation":
		resp, err = g.dispatchBeginCreateOperation(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *GenerateDetailedCostReportServerTransport) dispatchBeginCreateOperation(req *http.Request) (*http.Response, error) {
	if g.srv.BeginCreateOperation == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOperation not implemented")}
	}
	beginCreateOperation := g.beginCreateOperation.get(req)
	if beginCreateOperation == nil {
		const regexStr = `/(?P<scope>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.CostManagement/generateDetailedCostReport`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcostmanagement.GenerateDetailedCostReportDefinition](req)
		if err != nil {
			return nil, err
		}
		scopeParam, err := url.PathUnescape(matches[regex.SubexpIndex("scope")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := g.srv.BeginCreateOperation(req.Context(), scopeParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOperation = &respr
		g.beginCreateOperation.add(req, beginCreateOperation)
	}

	resp, err := server.PollerResponderNext(beginCreateOperation, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		g.beginCreateOperation.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOperation) {
		g.beginCreateOperation.remove(req)
	}

	return resp, nil
}
