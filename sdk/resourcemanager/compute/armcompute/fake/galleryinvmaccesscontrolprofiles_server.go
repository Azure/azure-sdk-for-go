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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/compute/armcompute/v7"
	"net/http"
	"net/url"
	"regexp"
)

// GalleryInVMAccessControlProfilesServer is a fake server for instances of the armcompute.GalleryInVMAccessControlProfilesClient type.
type GalleryInVMAccessControlProfilesServer struct {
	// BeginCreateOrUpdate is the fake for method GalleryInVMAccessControlProfilesClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, galleryInVMAccessControlProfile armcompute.GalleryInVMAccessControlProfile, options *armcompute.GalleryInVMAccessControlProfilesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armcompute.GalleryInVMAccessControlProfilesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method GalleryInVMAccessControlProfilesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, options *armcompute.GalleryInVMAccessControlProfilesClientBeginDeleteOptions) (resp azfake.PollerResponder[armcompute.GalleryInVMAccessControlProfilesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method GalleryInVMAccessControlProfilesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, options *armcompute.GalleryInVMAccessControlProfilesClientGetOptions) (resp azfake.Responder[armcompute.GalleryInVMAccessControlProfilesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByGalleryPager is the fake for method GalleryInVMAccessControlProfilesClient.NewListByGalleryPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByGalleryPager func(resourceGroupName string, galleryName string, options *armcompute.GalleryInVMAccessControlProfilesClientListByGalleryOptions) (resp azfake.PagerResponder[armcompute.GalleryInVMAccessControlProfilesClientListByGalleryResponse])

	// BeginUpdate is the fake for method GalleryInVMAccessControlProfilesClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK
	BeginUpdate func(ctx context.Context, resourceGroupName string, galleryName string, inVMAccessControlProfileName string, galleryInVMAccessControlProfile armcompute.GalleryInVMAccessControlProfileUpdate, options *armcompute.GalleryInVMAccessControlProfilesClientBeginUpdateOptions) (resp azfake.PollerResponder[armcompute.GalleryInVMAccessControlProfilesClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewGalleryInVMAccessControlProfilesServerTransport creates a new instance of GalleryInVMAccessControlProfilesServerTransport with the provided implementation.
// The returned GalleryInVMAccessControlProfilesServerTransport instance is connected to an instance of armcompute.GalleryInVMAccessControlProfilesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewGalleryInVMAccessControlProfilesServerTransport(srv *GalleryInVMAccessControlProfilesServer) *GalleryInVMAccessControlProfilesServerTransport {
	return &GalleryInVMAccessControlProfilesServerTransport{
		srv:                   srv,
		beginCreateOrUpdate:   newTracker[azfake.PollerResponder[armcompute.GalleryInVMAccessControlProfilesClientCreateOrUpdateResponse]](),
		beginDelete:           newTracker[azfake.PollerResponder[armcompute.GalleryInVMAccessControlProfilesClientDeleteResponse]](),
		newListByGalleryPager: newTracker[azfake.PagerResponder[armcompute.GalleryInVMAccessControlProfilesClientListByGalleryResponse]](),
		beginUpdate:           newTracker[azfake.PollerResponder[armcompute.GalleryInVMAccessControlProfilesClientUpdateResponse]](),
	}
}

// GalleryInVMAccessControlProfilesServerTransport connects instances of armcompute.GalleryInVMAccessControlProfilesClient to instances of GalleryInVMAccessControlProfilesServer.
// Don't use this type directly, use NewGalleryInVMAccessControlProfilesServerTransport instead.
type GalleryInVMAccessControlProfilesServerTransport struct {
	srv                   *GalleryInVMAccessControlProfilesServer
	beginCreateOrUpdate   *tracker[azfake.PollerResponder[armcompute.GalleryInVMAccessControlProfilesClientCreateOrUpdateResponse]]
	beginDelete           *tracker[azfake.PollerResponder[armcompute.GalleryInVMAccessControlProfilesClientDeleteResponse]]
	newListByGalleryPager *tracker[azfake.PagerResponder[armcompute.GalleryInVMAccessControlProfilesClientListByGalleryResponse]]
	beginUpdate           *tracker[azfake.PollerResponder[armcompute.GalleryInVMAccessControlProfilesClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for GalleryInVMAccessControlProfilesServerTransport.
func (g *GalleryInVMAccessControlProfilesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return g.dispatchToMethodFake(req, method)
}

func (g *GalleryInVMAccessControlProfilesServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if galleryInVMAccessControlProfilesServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = galleryInVMAccessControlProfilesServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "GalleryInVMAccessControlProfilesClient.BeginCreateOrUpdate":
				res.resp, res.err = g.dispatchBeginCreateOrUpdate(req)
			case "GalleryInVMAccessControlProfilesClient.BeginDelete":
				res.resp, res.err = g.dispatchBeginDelete(req)
			case "GalleryInVMAccessControlProfilesClient.Get":
				res.resp, res.err = g.dispatchGet(req)
			case "GalleryInVMAccessControlProfilesClient.NewListByGalleryPager":
				res.resp, res.err = g.dispatchNewListByGalleryPager(req)
			case "GalleryInVMAccessControlProfilesClient.BeginUpdate":
				res.resp, res.err = g.dispatchBeginUpdate(req)
			default:
				res.err = fmt.Errorf("unhandled API %s", method)
			}

		}
		select {
		case resultChan <- res:
		case <-req.Context().Done():
		}
	}()

	select {
	case <-req.Context().Done():
		return nil, req.Context().Err()
	case res := <-resultChan:
		return res.resp, res.err
	}
}

func (g *GalleryInVMAccessControlProfilesServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if g.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := g.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/galleries/(?P<galleryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/inVMAccessControlProfiles/(?P<inVMAccessControlProfileName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcompute.GalleryInVMAccessControlProfile](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		galleryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("galleryName")])
		if err != nil {
			return nil, err
		}
		inVMAccessControlProfileNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("inVMAccessControlProfileName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := g.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, galleryNameParam, inVMAccessControlProfileNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		g.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		g.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		g.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (g *GalleryInVMAccessControlProfilesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if g.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := g.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/galleries/(?P<galleryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/inVMAccessControlProfiles/(?P<inVMAccessControlProfileName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		galleryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("galleryName")])
		if err != nil {
			return nil, err
		}
		inVMAccessControlProfileNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("inVMAccessControlProfileName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := g.srv.BeginDelete(req.Context(), resourceGroupNameParam, galleryNameParam, inVMAccessControlProfileNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		g.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		g.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		g.beginDelete.remove(req)
	}

	return resp, nil
}

func (g *GalleryInVMAccessControlProfilesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if g.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/galleries/(?P<galleryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/inVMAccessControlProfiles/(?P<inVMAccessControlProfileName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if len(matches) < 5 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	galleryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("galleryName")])
	if err != nil {
		return nil, err
	}
	inVMAccessControlProfileNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("inVMAccessControlProfileName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := g.srv.Get(req.Context(), resourceGroupNameParam, galleryNameParam, inVMAccessControlProfileNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).GalleryInVMAccessControlProfile, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (g *GalleryInVMAccessControlProfilesServerTransport) dispatchNewListByGalleryPager(req *http.Request) (*http.Response, error) {
	if g.srv.NewListByGalleryPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByGalleryPager not implemented")}
	}
	newListByGalleryPager := g.newListByGalleryPager.get(req)
	if newListByGalleryPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/galleries/(?P<galleryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/inVMAccessControlProfiles`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		galleryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("galleryName")])
		if err != nil {
			return nil, err
		}
		resp := g.srv.NewListByGalleryPager(resourceGroupNameParam, galleryNameParam, nil)
		newListByGalleryPager = &resp
		g.newListByGalleryPager.add(req, newListByGalleryPager)
		server.PagerResponderInjectNextLinks(newListByGalleryPager, req, func(page *armcompute.GalleryInVMAccessControlProfilesClientListByGalleryResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByGalleryPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		g.newListByGalleryPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByGalleryPager) {
		g.newListByGalleryPager.remove(req)
	}
	return resp, nil
}

func (g *GalleryInVMAccessControlProfilesServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if g.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := g.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Compute/galleries/(?P<galleryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/inVMAccessControlProfiles/(?P<inVMAccessControlProfileName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if len(matches) < 5 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcompute.GalleryInVMAccessControlProfileUpdate](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		galleryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("galleryName")])
		if err != nil {
			return nil, err
		}
		inVMAccessControlProfileNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("inVMAccessControlProfileName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := g.srv.BeginUpdate(req.Context(), resourceGroupNameParam, galleryNameParam, inVMAccessControlProfileNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		g.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		g.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		g.beginUpdate.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to GalleryInVMAccessControlProfilesServerTransport
var galleryInVMAccessControlProfilesServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
