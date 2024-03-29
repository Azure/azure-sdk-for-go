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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/paloaltonetworksngfw/armpanngfw"
	"net/http"
	"net/url"
	"regexp"
)

// FqdnListLocalRulestackServer is a fake server for instances of the armpanngfw.FqdnListLocalRulestackClient type.
type FqdnListLocalRulestackServer struct {
	// BeginCreateOrUpdate is the fake for method FqdnListLocalRulestackClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, localRulestackName string, name string, resource armpanngfw.FqdnListLocalRulestackResource, options *armpanngfw.FqdnListLocalRulestackClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armpanngfw.FqdnListLocalRulestackClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method FqdnListLocalRulestackClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, localRulestackName string, name string, options *armpanngfw.FqdnListLocalRulestackClientBeginDeleteOptions) (resp azfake.PollerResponder[armpanngfw.FqdnListLocalRulestackClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method FqdnListLocalRulestackClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, localRulestackName string, name string, options *armpanngfw.FqdnListLocalRulestackClientGetOptions) (resp azfake.Responder[armpanngfw.FqdnListLocalRulestackClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByLocalRulestacksPager is the fake for method FqdnListLocalRulestackClient.NewListByLocalRulestacksPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByLocalRulestacksPager func(resourceGroupName string, localRulestackName string, options *armpanngfw.FqdnListLocalRulestackClientListByLocalRulestacksOptions) (resp azfake.PagerResponder[armpanngfw.FqdnListLocalRulestackClientListByLocalRulestacksResponse])
}

// NewFqdnListLocalRulestackServerTransport creates a new instance of FqdnListLocalRulestackServerTransport with the provided implementation.
// The returned FqdnListLocalRulestackServerTransport instance is connected to an instance of armpanngfw.FqdnListLocalRulestackClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewFqdnListLocalRulestackServerTransport(srv *FqdnListLocalRulestackServer) *FqdnListLocalRulestackServerTransport {
	return &FqdnListLocalRulestackServerTransport{
		srv:                           srv,
		beginCreateOrUpdate:           newTracker[azfake.PollerResponder[armpanngfw.FqdnListLocalRulestackClientCreateOrUpdateResponse]](),
		beginDelete:                   newTracker[azfake.PollerResponder[armpanngfw.FqdnListLocalRulestackClientDeleteResponse]](),
		newListByLocalRulestacksPager: newTracker[azfake.PagerResponder[armpanngfw.FqdnListLocalRulestackClientListByLocalRulestacksResponse]](),
	}
}

// FqdnListLocalRulestackServerTransport connects instances of armpanngfw.FqdnListLocalRulestackClient to instances of FqdnListLocalRulestackServer.
// Don't use this type directly, use NewFqdnListLocalRulestackServerTransport instead.
type FqdnListLocalRulestackServerTransport struct {
	srv                           *FqdnListLocalRulestackServer
	beginCreateOrUpdate           *tracker[azfake.PollerResponder[armpanngfw.FqdnListLocalRulestackClientCreateOrUpdateResponse]]
	beginDelete                   *tracker[azfake.PollerResponder[armpanngfw.FqdnListLocalRulestackClientDeleteResponse]]
	newListByLocalRulestacksPager *tracker[azfake.PagerResponder[armpanngfw.FqdnListLocalRulestackClientListByLocalRulestacksResponse]]
}

// Do implements the policy.Transporter interface for FqdnListLocalRulestackServerTransport.
func (f *FqdnListLocalRulestackServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "FqdnListLocalRulestackClient.BeginCreateOrUpdate":
		resp, err = f.dispatchBeginCreateOrUpdate(req)
	case "FqdnListLocalRulestackClient.BeginDelete":
		resp, err = f.dispatchBeginDelete(req)
	case "FqdnListLocalRulestackClient.Get":
		resp, err = f.dispatchGet(req)
	case "FqdnListLocalRulestackClient.NewListByLocalRulestacksPager":
		resp, err = f.dispatchNewListByLocalRulestacksPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (f *FqdnListLocalRulestackServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if f.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := f.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/PaloAltoNetworks\.Cloudngfw/localRulestacks/(?P<localRulestackName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/fqdnlists/(?P<name>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armpanngfw.FqdnListLocalRulestackResource](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		localRulestackNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("localRulestackName")])
		if err != nil {
			return nil, err
		}
		nameParam, err := url.PathUnescape(matches[regex.SubexpIndex("name")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := f.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, localRulestackNameParam, nameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		f.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		f.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		f.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (f *FqdnListLocalRulestackServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if f.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := f.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/PaloAltoNetworks\.Cloudngfw/localRulestacks/(?P<localRulestackName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/fqdnlists/(?P<name>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		localRulestackNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("localRulestackName")])
		if err != nil {
			return nil, err
		}
		nameParam, err := url.PathUnescape(matches[regex.SubexpIndex("name")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := f.srv.BeginDelete(req.Context(), resourceGroupNameParam, localRulestackNameParam, nameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		f.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		f.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		f.beginDelete.remove(req)
	}

	return resp, nil
}

func (f *FqdnListLocalRulestackServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if f.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/PaloAltoNetworks\.Cloudngfw/localRulestacks/(?P<localRulestackName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/fqdnlists/(?P<name>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	localRulestackNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("localRulestackName")])
	if err != nil {
		return nil, err
	}
	nameParam, err := url.PathUnescape(matches[regex.SubexpIndex("name")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := f.srv.Get(req.Context(), resourceGroupNameParam, localRulestackNameParam, nameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).FqdnListLocalRulestackResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (f *FqdnListLocalRulestackServerTransport) dispatchNewListByLocalRulestacksPager(req *http.Request) (*http.Response, error) {
	if f.srv.NewListByLocalRulestacksPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByLocalRulestacksPager not implemented")}
	}
	newListByLocalRulestacksPager := f.newListByLocalRulestacksPager.get(req)
	if newListByLocalRulestacksPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/PaloAltoNetworks\.Cloudngfw/localRulestacks/(?P<localRulestackName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/fqdnlists`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		localRulestackNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("localRulestackName")])
		if err != nil {
			return nil, err
		}
		resp := f.srv.NewListByLocalRulestacksPager(resourceGroupNameParam, localRulestackNameParam, nil)
		newListByLocalRulestacksPager = &resp
		f.newListByLocalRulestacksPager.add(req, newListByLocalRulestacksPager)
		server.PagerResponderInjectNextLinks(newListByLocalRulestacksPager, req, func(page *armpanngfw.FqdnListLocalRulestackClientListByLocalRulestacksResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByLocalRulestacksPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		f.newListByLocalRulestacksPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByLocalRulestacksPager) {
		f.newListByLocalRulestacksPager.remove(req)
	}
	return resp, nil
}
