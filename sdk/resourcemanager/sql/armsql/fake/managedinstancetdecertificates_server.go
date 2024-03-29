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

// ManagedInstanceTdeCertificatesServer is a fake server for instances of the armsql.ManagedInstanceTdeCertificatesClient type.
type ManagedInstanceTdeCertificatesServer struct {
	// BeginCreate is the fake for method ManagedInstanceTdeCertificatesClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreate func(ctx context.Context, resourceGroupName string, managedInstanceName string, parameters armsql.TdeCertificate, options *armsql.ManagedInstanceTdeCertificatesClientBeginCreateOptions) (resp azfake.PollerResponder[armsql.ManagedInstanceTdeCertificatesClientCreateResponse], errResp azfake.ErrorResponder)
}

// NewManagedInstanceTdeCertificatesServerTransport creates a new instance of ManagedInstanceTdeCertificatesServerTransport with the provided implementation.
// The returned ManagedInstanceTdeCertificatesServerTransport instance is connected to an instance of armsql.ManagedInstanceTdeCertificatesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewManagedInstanceTdeCertificatesServerTransport(srv *ManagedInstanceTdeCertificatesServer) *ManagedInstanceTdeCertificatesServerTransport {
	return &ManagedInstanceTdeCertificatesServerTransport{
		srv:         srv,
		beginCreate: newTracker[azfake.PollerResponder[armsql.ManagedInstanceTdeCertificatesClientCreateResponse]](),
	}
}

// ManagedInstanceTdeCertificatesServerTransport connects instances of armsql.ManagedInstanceTdeCertificatesClient to instances of ManagedInstanceTdeCertificatesServer.
// Don't use this type directly, use NewManagedInstanceTdeCertificatesServerTransport instead.
type ManagedInstanceTdeCertificatesServerTransport struct {
	srv         *ManagedInstanceTdeCertificatesServer
	beginCreate *tracker[azfake.PollerResponder[armsql.ManagedInstanceTdeCertificatesClientCreateResponse]]
}

// Do implements the policy.Transporter interface for ManagedInstanceTdeCertificatesServerTransport.
func (m *ManagedInstanceTdeCertificatesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ManagedInstanceTdeCertificatesClient.BeginCreate":
		resp, err = m.dispatchBeginCreate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *ManagedInstanceTdeCertificatesServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if m.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := m.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sql/managedInstances/(?P<managedInstanceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/tdeCertificates`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armsql.TdeCertificate](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		managedInstanceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("managedInstanceName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := m.srv.BeginCreate(req.Context(), resourceGroupNameParam, managedInstanceNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		m.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		m.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		m.beginCreate.remove(req)
	}

	return resp, nil
}
