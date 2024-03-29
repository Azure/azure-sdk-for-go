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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/azurestackhci/armazurestackhci/v2"
	"net/http"
	"net/url"
	"regexp"
)

// MarketplaceGalleryImagesServer is a fake server for instances of the armazurestackhci.MarketplaceGalleryImagesClient type.
type MarketplaceGalleryImagesServer struct {
	// BeginCreateOrUpdate is the fake for method MarketplaceGalleryImagesClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, marketplaceGalleryImageName string, marketplaceGalleryImages armazurestackhci.MarketplaceGalleryImages, options *armazurestackhci.MarketplaceGalleryImagesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armazurestackhci.MarketplaceGalleryImagesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method MarketplaceGalleryImagesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, marketplaceGalleryImageName string, options *armazurestackhci.MarketplaceGalleryImagesClientBeginDeleteOptions) (resp azfake.PollerResponder[armazurestackhci.MarketplaceGalleryImagesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method MarketplaceGalleryImagesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, marketplaceGalleryImageName string, options *armazurestackhci.MarketplaceGalleryImagesClientGetOptions) (resp azfake.Responder[armazurestackhci.MarketplaceGalleryImagesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method MarketplaceGalleryImagesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, options *armazurestackhci.MarketplaceGalleryImagesClientListOptions) (resp azfake.PagerResponder[armazurestackhci.MarketplaceGalleryImagesClientListResponse])

	// NewListAllPager is the fake for method MarketplaceGalleryImagesClient.NewListAllPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListAllPager func(options *armazurestackhci.MarketplaceGalleryImagesClientListAllOptions) (resp azfake.PagerResponder[armazurestackhci.MarketplaceGalleryImagesClientListAllResponse])

	// BeginUpdate is the fake for method MarketplaceGalleryImagesClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, marketplaceGalleryImageName string, marketplaceGalleryImages armazurestackhci.MarketplaceGalleryImagesUpdateRequest, options *armazurestackhci.MarketplaceGalleryImagesClientBeginUpdateOptions) (resp azfake.PollerResponder[armazurestackhci.MarketplaceGalleryImagesClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewMarketplaceGalleryImagesServerTransport creates a new instance of MarketplaceGalleryImagesServerTransport with the provided implementation.
// The returned MarketplaceGalleryImagesServerTransport instance is connected to an instance of armazurestackhci.MarketplaceGalleryImagesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewMarketplaceGalleryImagesServerTransport(srv *MarketplaceGalleryImagesServer) *MarketplaceGalleryImagesServerTransport {
	return &MarketplaceGalleryImagesServerTransport{
		srv:                 srv,
		beginCreateOrUpdate: newTracker[azfake.PollerResponder[armazurestackhci.MarketplaceGalleryImagesClientCreateOrUpdateResponse]](),
		beginDelete:         newTracker[azfake.PollerResponder[armazurestackhci.MarketplaceGalleryImagesClientDeleteResponse]](),
		newListPager:        newTracker[azfake.PagerResponder[armazurestackhci.MarketplaceGalleryImagesClientListResponse]](),
		newListAllPager:     newTracker[azfake.PagerResponder[armazurestackhci.MarketplaceGalleryImagesClientListAllResponse]](),
		beginUpdate:         newTracker[azfake.PollerResponder[armazurestackhci.MarketplaceGalleryImagesClientUpdateResponse]](),
	}
}

// MarketplaceGalleryImagesServerTransport connects instances of armazurestackhci.MarketplaceGalleryImagesClient to instances of MarketplaceGalleryImagesServer.
// Don't use this type directly, use NewMarketplaceGalleryImagesServerTransport instead.
type MarketplaceGalleryImagesServerTransport struct {
	srv                 *MarketplaceGalleryImagesServer
	beginCreateOrUpdate *tracker[azfake.PollerResponder[armazurestackhci.MarketplaceGalleryImagesClientCreateOrUpdateResponse]]
	beginDelete         *tracker[azfake.PollerResponder[armazurestackhci.MarketplaceGalleryImagesClientDeleteResponse]]
	newListPager        *tracker[azfake.PagerResponder[armazurestackhci.MarketplaceGalleryImagesClientListResponse]]
	newListAllPager     *tracker[azfake.PagerResponder[armazurestackhci.MarketplaceGalleryImagesClientListAllResponse]]
	beginUpdate         *tracker[azfake.PollerResponder[armazurestackhci.MarketplaceGalleryImagesClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for MarketplaceGalleryImagesServerTransport.
func (m *MarketplaceGalleryImagesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "MarketplaceGalleryImagesClient.BeginCreateOrUpdate":
		resp, err = m.dispatchBeginCreateOrUpdate(req)
	case "MarketplaceGalleryImagesClient.BeginDelete":
		resp, err = m.dispatchBeginDelete(req)
	case "MarketplaceGalleryImagesClient.Get":
		resp, err = m.dispatchGet(req)
	case "MarketplaceGalleryImagesClient.NewListPager":
		resp, err = m.dispatchNewListPager(req)
	case "MarketplaceGalleryImagesClient.NewListAllPager":
		resp, err = m.dispatchNewListAllPager(req)
	case "MarketplaceGalleryImagesClient.BeginUpdate":
		resp, err = m.dispatchBeginUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (m *MarketplaceGalleryImagesServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if m.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := m.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureStackHCI/marketplaceGalleryImages/(?P<marketplaceGalleryImageName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armazurestackhci.MarketplaceGalleryImages](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		marketplaceGalleryImageNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("marketplaceGalleryImageName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := m.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, marketplaceGalleryImageNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		m.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		m.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		m.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (m *MarketplaceGalleryImagesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if m.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := m.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureStackHCI/marketplaceGalleryImages/(?P<marketplaceGalleryImageName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		marketplaceGalleryImageNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("marketplaceGalleryImageName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := m.srv.BeginDelete(req.Context(), resourceGroupNameParam, marketplaceGalleryImageNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		m.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		m.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		m.beginDelete.remove(req)
	}

	return resp, nil
}

func (m *MarketplaceGalleryImagesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if m.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureStackHCI/marketplaceGalleryImages/(?P<marketplaceGalleryImageName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	marketplaceGalleryImageNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("marketplaceGalleryImageName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := m.srv.Get(req.Context(), resourceGroupNameParam, marketplaceGalleryImageNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).MarketplaceGalleryImages, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *MarketplaceGalleryImagesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if m.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := m.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureStackHCI/marketplaceGalleryImages`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := m.srv.NewListPager(resourceGroupNameParam, nil)
		newListPager = &resp
		m.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armazurestackhci.MarketplaceGalleryImagesClientListResponse, createLink func() string) {
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

func (m *MarketplaceGalleryImagesServerTransport) dispatchNewListAllPager(req *http.Request) (*http.Response, error) {
	if m.srv.NewListAllPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListAllPager not implemented")}
	}
	newListAllPager := m.newListAllPager.get(req)
	if newListAllPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureStackHCI/marketplaceGalleryImages`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := m.srv.NewListAllPager(nil)
		newListAllPager = &resp
		m.newListAllPager.add(req, newListAllPager)
		server.PagerResponderInjectNextLinks(newListAllPager, req, func(page *armazurestackhci.MarketplaceGalleryImagesClientListAllResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListAllPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		m.newListAllPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListAllPager) {
		m.newListAllPager.remove(req)
	}
	return resp, nil
}

func (m *MarketplaceGalleryImagesServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if m.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := m.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.AzureStackHCI/marketplaceGalleryImages/(?P<marketplaceGalleryImageName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armazurestackhci.MarketplaceGalleryImagesUpdateRequest](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		marketplaceGalleryImageNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("marketplaceGalleryImageName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := m.srv.BeginUpdate(req.Context(), resourceGroupNameParam, marketplaceGalleryImageNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		m.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		m.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		m.beginUpdate.remove(req)
	}

	return resp, nil
}
