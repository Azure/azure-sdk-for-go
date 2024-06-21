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
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/mobilenetwork/armmobilenetwork/v4"
	"net/http"
	"net/url"
	"regexp"
)

// RoutingInfoServer is a fake server for instances of the armmobilenetwork.RoutingInfoClient type.
type RoutingInfoServer struct {
	// Get is the fake for method RoutingInfoClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, packetCoreControlPlaneName string, options *armmobilenetwork.RoutingInfoClientGetOptions) (resp azfake.Responder[armmobilenetwork.RoutingInfoClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method RoutingInfoClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, packetCoreControlPlaneName string, options *armmobilenetwork.RoutingInfoClientListOptions) (resp azfake.PagerResponder[armmobilenetwork.RoutingInfoClientListResponse])
}

// NewRoutingInfoServerTransport creates a new instance of RoutingInfoServerTransport with the provided implementation.
// The returned RoutingInfoServerTransport instance is connected to an instance of armmobilenetwork.RoutingInfoClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewRoutingInfoServerTransport(srv *RoutingInfoServer) *RoutingInfoServerTransport {
	return &RoutingInfoServerTransport{
		srv:          srv,
		newListPager: newTracker[azfake.PagerResponder[armmobilenetwork.RoutingInfoClientListResponse]](),
	}
}

// RoutingInfoServerTransport connects instances of armmobilenetwork.RoutingInfoClient to instances of RoutingInfoServer.
// Don't use this type directly, use NewRoutingInfoServerTransport instead.
type RoutingInfoServerTransport struct {
	srv          *RoutingInfoServer
	newListPager *tracker[azfake.PagerResponder[armmobilenetwork.RoutingInfoClientListResponse]]
}

// Do implements the policy.Transporter interface for RoutingInfoServerTransport.
func (r *RoutingInfoServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "RoutingInfoClient.Get":
		resp, err = r.dispatchGet(req)
	case "RoutingInfoClient.NewListPager":
		resp, err = r.dispatchNewListPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *RoutingInfoServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if r.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MobileNetwork/packetCoreControlPlanes/(?P<packetCoreControlPlaneName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/routingInfo/default`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	packetCoreControlPlaneNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("packetCoreControlPlaneName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := r.srv.Get(req.Context(), resourceGroupNameParam, packetCoreControlPlaneNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).RoutingInfoModel, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *RoutingInfoServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := r.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MobileNetwork/packetCoreControlPlanes/(?P<packetCoreControlPlaneName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/routingInfo`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		packetCoreControlPlaneNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("packetCoreControlPlaneName")])
		if err != nil {
			return nil, err
		}
		resp := r.srv.NewListPager(resourceGroupNameParam, packetCoreControlPlaneNameParam, nil)
		newListPager = &resp
		r.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armmobilenetwork.RoutingInfoClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		r.newListPager.remove(req)
	}
	return resp, nil
}
