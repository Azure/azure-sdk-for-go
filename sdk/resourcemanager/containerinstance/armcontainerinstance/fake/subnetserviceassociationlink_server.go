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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/containerinstance/armcontainerinstance/v2"
	"net/http"
	"net/url"
	"regexp"
)

// SubnetServiceAssociationLinkServer is a fake server for instances of the armcontainerinstance.SubnetServiceAssociationLinkClient type.
type SubnetServiceAssociationLinkServer struct {
	// BeginDelete is the fake for method SubnetServiceAssociationLinkClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, virtualNetworkName string, subnetName string, options *armcontainerinstance.SubnetServiceAssociationLinkClientBeginDeleteOptions) (resp azfake.PollerResponder[armcontainerinstance.SubnetServiceAssociationLinkClientDeleteResponse], errResp azfake.ErrorResponder)
}

// NewSubnetServiceAssociationLinkServerTransport creates a new instance of SubnetServiceAssociationLinkServerTransport with the provided implementation.
// The returned SubnetServiceAssociationLinkServerTransport instance is connected to an instance of armcontainerinstance.SubnetServiceAssociationLinkClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSubnetServiceAssociationLinkServerTransport(srv *SubnetServiceAssociationLinkServer) *SubnetServiceAssociationLinkServerTransport {
	return &SubnetServiceAssociationLinkServerTransport{
		srv:         srv,
		beginDelete: newTracker[azfake.PollerResponder[armcontainerinstance.SubnetServiceAssociationLinkClientDeleteResponse]](),
	}
}

// SubnetServiceAssociationLinkServerTransport connects instances of armcontainerinstance.SubnetServiceAssociationLinkClient to instances of SubnetServiceAssociationLinkServer.
// Don't use this type directly, use NewSubnetServiceAssociationLinkServerTransport instead.
type SubnetServiceAssociationLinkServerTransport struct {
	srv         *SubnetServiceAssociationLinkServer
	beginDelete *tracker[azfake.PollerResponder[armcontainerinstance.SubnetServiceAssociationLinkClientDeleteResponse]]
}

// Do implements the policy.Transporter interface for SubnetServiceAssociationLinkServerTransport.
func (s *SubnetServiceAssociationLinkServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "SubnetServiceAssociationLinkClient.BeginDelete":
		resp, err = s.dispatchBeginDelete(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *SubnetServiceAssociationLinkServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if s.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := s.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourcegroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/virtualNetworks/(?P<virtualNetworkName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/subnets/(?P<subnetName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ContainerInstance/serviceAssociationLinks/default`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualNetworkNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualNetworkName")])
		if err != nil {
			return nil, err
		}
		subnetNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("subnetName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginDelete(req.Context(), resourceGroupNameParam, virtualNetworkNameParam, subnetNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		s.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		s.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		s.beginDelete.remove(req)
	}

	return resp, nil
}
