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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/quota/armquota"
	"net/http"
	"net/url"
	"regexp"
)

// GroupQuotaLimitsServer is a fake server for instances of the armquota.GroupQuotaLimitsClient type.
type GroupQuotaLimitsServer struct {
	// List is the fake for method GroupQuotaLimitsClient.List
	// HTTP status codes to indicate success: http.StatusOK
	List func(ctx context.Context, managementGroupID string, groupQuotaName string, resourceProviderName string, location string, options *armquota.GroupQuotaLimitsClientListOptions) (resp azfake.Responder[armquota.GroupQuotaLimitsClientListResponse], errResp azfake.ErrorResponder)
}

// NewGroupQuotaLimitsServerTransport creates a new instance of GroupQuotaLimitsServerTransport with the provided implementation.
// The returned GroupQuotaLimitsServerTransport instance is connected to an instance of armquota.GroupQuotaLimitsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewGroupQuotaLimitsServerTransport(srv *GroupQuotaLimitsServer) *GroupQuotaLimitsServerTransport {
	return &GroupQuotaLimitsServerTransport{srv: srv}
}

// GroupQuotaLimitsServerTransport connects instances of armquota.GroupQuotaLimitsClient to instances of GroupQuotaLimitsServer.
// Don't use this type directly, use NewGroupQuotaLimitsServerTransport instead.
type GroupQuotaLimitsServerTransport struct {
	srv *GroupQuotaLimitsServer
}

// Do implements the policy.Transporter interface for GroupQuotaLimitsServerTransport.
func (g *GroupQuotaLimitsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "GroupQuotaLimitsClient.List":
		resp, err = g.dispatchList(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (g *GroupQuotaLimitsServerTransport) dispatchList(req *http.Request) (*http.Response, error) {
	if g.srv.List == nil {
		return nil, &nonRetriableError{errors.New("fake for method List not implemented")}
	}
	const regexStr = `/providers/Microsoft\.Management/managementGroups/(?P<managementGroupId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Quota/groupQuotas/(?P<groupQuotaName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceProviders/(?P<resourceProviderName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/groupQuotaLimits/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	managementGroupIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("managementGroupId")])
	if err != nil {
		return nil, err
	}
	groupQuotaNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("groupQuotaName")])
	if err != nil {
		return nil, err
	}
	resourceProviderNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceProviderName")])
	if err != nil {
		return nil, err
	}
	locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.List(req.Context(), managementGroupIDParam, groupQuotaNameParam, resourceProviderNameParam, locationParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).GroupQuotaLimitList, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
