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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v7"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
)

// VirtualHubsServer is a fake server for instances of the armnetwork.VirtualHubsClient type.
type VirtualHubsServer struct {
	// BeginCreateOrUpdate is the fake for method VirtualHubsClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, virtualHubName string, virtualHubParameters armnetwork.VirtualHub, options *armnetwork.VirtualHubsClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armnetwork.VirtualHubsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method VirtualHubsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, virtualHubName string, options *armnetwork.VirtualHubsClientBeginDeleteOptions) (resp azfake.PollerResponder[armnetwork.VirtualHubsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method VirtualHubsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, virtualHubName string, options *armnetwork.VirtualHubsClientGetOptions) (resp azfake.Responder[armnetwork.VirtualHubsClientGetResponse], errResp azfake.ErrorResponder)

	// BeginGetEffectiveVirtualHubRoutes is the fake for method VirtualHubsClient.BeginGetEffectiveVirtualHubRoutes
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginGetEffectiveVirtualHubRoutes func(ctx context.Context, resourceGroupName string, virtualHubName string, options *armnetwork.VirtualHubsClientBeginGetEffectiveVirtualHubRoutesOptions) (resp azfake.PollerResponder[armnetwork.VirtualHubsClientGetEffectiveVirtualHubRoutesResponse], errResp azfake.ErrorResponder)

	// BeginGetInboundRoutes is the fake for method VirtualHubsClient.BeginGetInboundRoutes
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginGetInboundRoutes func(ctx context.Context, resourceGroupName string, virtualHubName string, getInboundRoutesParameters armnetwork.GetInboundRoutesParameters, options *armnetwork.VirtualHubsClientBeginGetInboundRoutesOptions) (resp azfake.PollerResponder[armnetwork.VirtualHubsClientGetInboundRoutesResponse], errResp azfake.ErrorResponder)

	// BeginGetOutboundRoutes is the fake for method VirtualHubsClient.BeginGetOutboundRoutes
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginGetOutboundRoutes func(ctx context.Context, resourceGroupName string, virtualHubName string, getOutboundRoutesParameters armnetwork.GetOutboundRoutesParameters, options *armnetwork.VirtualHubsClientBeginGetOutboundRoutesOptions) (resp azfake.PollerResponder[armnetwork.VirtualHubsClientGetOutboundRoutesResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method VirtualHubsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armnetwork.VirtualHubsClientListOptions) (resp azfake.PagerResponder[armnetwork.VirtualHubsClientListResponse])

	// NewListByResourceGroupPager is the fake for method VirtualHubsClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armnetwork.VirtualHubsClientListByResourceGroupOptions) (resp azfake.PagerResponder[armnetwork.VirtualHubsClientListByResourceGroupResponse])

	// UpdateTags is the fake for method VirtualHubsClient.UpdateTags
	// HTTP status codes to indicate success: http.StatusOK
	UpdateTags func(ctx context.Context, resourceGroupName string, virtualHubName string, virtualHubParameters armnetwork.TagsObject, options *armnetwork.VirtualHubsClientUpdateTagsOptions) (resp azfake.Responder[armnetwork.VirtualHubsClientUpdateTagsResponse], errResp azfake.ErrorResponder)
}

// NewVirtualHubsServerTransport creates a new instance of VirtualHubsServerTransport with the provided implementation.
// The returned VirtualHubsServerTransport instance is connected to an instance of armnetwork.VirtualHubsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewVirtualHubsServerTransport(srv *VirtualHubsServer) *VirtualHubsServerTransport {
	return &VirtualHubsServerTransport{
		srv:                               srv,
		beginCreateOrUpdate:               newTracker[azfake.PollerResponder[armnetwork.VirtualHubsClientCreateOrUpdateResponse]](),
		beginDelete:                       newTracker[azfake.PollerResponder[armnetwork.VirtualHubsClientDeleteResponse]](),
		beginGetEffectiveVirtualHubRoutes: newTracker[azfake.PollerResponder[armnetwork.VirtualHubsClientGetEffectiveVirtualHubRoutesResponse]](),
		beginGetInboundRoutes:             newTracker[azfake.PollerResponder[armnetwork.VirtualHubsClientGetInboundRoutesResponse]](),
		beginGetOutboundRoutes:            newTracker[azfake.PollerResponder[armnetwork.VirtualHubsClientGetOutboundRoutesResponse]](),
		newListPager:                      newTracker[azfake.PagerResponder[armnetwork.VirtualHubsClientListResponse]](),
		newListByResourceGroupPager:       newTracker[azfake.PagerResponder[armnetwork.VirtualHubsClientListByResourceGroupResponse]](),
	}
}

// VirtualHubsServerTransport connects instances of armnetwork.VirtualHubsClient to instances of VirtualHubsServer.
// Don't use this type directly, use NewVirtualHubsServerTransport instead.
type VirtualHubsServerTransport struct {
	srv                               *VirtualHubsServer
	beginCreateOrUpdate               *tracker[azfake.PollerResponder[armnetwork.VirtualHubsClientCreateOrUpdateResponse]]
	beginDelete                       *tracker[azfake.PollerResponder[armnetwork.VirtualHubsClientDeleteResponse]]
	beginGetEffectiveVirtualHubRoutes *tracker[azfake.PollerResponder[armnetwork.VirtualHubsClientGetEffectiveVirtualHubRoutesResponse]]
	beginGetInboundRoutes             *tracker[azfake.PollerResponder[armnetwork.VirtualHubsClientGetInboundRoutesResponse]]
	beginGetOutboundRoutes            *tracker[azfake.PollerResponder[armnetwork.VirtualHubsClientGetOutboundRoutesResponse]]
	newListPager                      *tracker[azfake.PagerResponder[armnetwork.VirtualHubsClientListResponse]]
	newListByResourceGroupPager       *tracker[azfake.PagerResponder[armnetwork.VirtualHubsClientListByResourceGroupResponse]]
}

// Do implements the policy.Transporter interface for VirtualHubsServerTransport.
func (v *VirtualHubsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return v.dispatchToMethodFake(req, method)
}

func (v *VirtualHubsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if virtualHubsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = virtualHubsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "VirtualHubsClient.BeginCreateOrUpdate":
				res.resp, res.err = v.dispatchBeginCreateOrUpdate(req)
			case "VirtualHubsClient.BeginDelete":
				res.resp, res.err = v.dispatchBeginDelete(req)
			case "VirtualHubsClient.Get":
				res.resp, res.err = v.dispatchGet(req)
			case "VirtualHubsClient.BeginGetEffectiveVirtualHubRoutes":
				res.resp, res.err = v.dispatchBeginGetEffectiveVirtualHubRoutes(req)
			case "VirtualHubsClient.BeginGetInboundRoutes":
				res.resp, res.err = v.dispatchBeginGetInboundRoutes(req)
			case "VirtualHubsClient.BeginGetOutboundRoutes":
				res.resp, res.err = v.dispatchBeginGetOutboundRoutes(req)
			case "VirtualHubsClient.NewListPager":
				res.resp, res.err = v.dispatchNewListPager(req)
			case "VirtualHubsClient.NewListByResourceGroupPager":
				res.resp, res.err = v.dispatchNewListByResourceGroupPager(req)
			case "VirtualHubsClient.UpdateTags":
				res.resp, res.err = v.dispatchUpdateTags(req)
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

func (v *VirtualHubsServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if v.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := v.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/virtualHubs/(?P<virtualHubName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.VirtualHub](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualHubNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualHubName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, virtualHubNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		v.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		v.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		v.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (v *VirtualHubsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if v.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := v.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/virtualHubs/(?P<virtualHubName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualHubNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualHubName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginDelete(req.Context(), resourceGroupNameParam, virtualHubNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		v.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		v.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		v.beginDelete.remove(req)
	}

	return resp, nil
}

func (v *VirtualHubsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if v.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/virtualHubs/(?P<virtualHubName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	virtualHubNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualHubName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := v.srv.Get(req.Context(), resourceGroupNameParam, virtualHubNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).VirtualHub, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (v *VirtualHubsServerTransport) dispatchBeginGetEffectiveVirtualHubRoutes(req *http.Request) (*http.Response, error) {
	if v.srv.BeginGetEffectiveVirtualHubRoutes == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginGetEffectiveVirtualHubRoutes not implemented")}
	}
	beginGetEffectiveVirtualHubRoutes := v.beginGetEffectiveVirtualHubRoutes.get(req)
	if beginGetEffectiveVirtualHubRoutes == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/virtualHubs/(?P<virtualHubName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/effectiveRoutes`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.EffectiveRoutesParameters](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualHubNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualHubName")])
		if err != nil {
			return nil, err
		}
		var options *armnetwork.VirtualHubsClientBeginGetEffectiveVirtualHubRoutesOptions
		if !reflect.ValueOf(body).IsZero() {
			options = &armnetwork.VirtualHubsClientBeginGetEffectiveVirtualHubRoutesOptions{
				EffectiveRoutesParameters: &body,
			}
		}
		respr, errRespr := v.srv.BeginGetEffectiveVirtualHubRoutes(req.Context(), resourceGroupNameParam, virtualHubNameParam, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginGetEffectiveVirtualHubRoutes = &respr
		v.beginGetEffectiveVirtualHubRoutes.add(req, beginGetEffectiveVirtualHubRoutes)
	}

	resp, err := server.PollerResponderNext(beginGetEffectiveVirtualHubRoutes, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginGetEffectiveVirtualHubRoutes.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginGetEffectiveVirtualHubRoutes) {
		v.beginGetEffectiveVirtualHubRoutes.remove(req)
	}

	return resp, nil
}

func (v *VirtualHubsServerTransport) dispatchBeginGetInboundRoutes(req *http.Request) (*http.Response, error) {
	if v.srv.BeginGetInboundRoutes == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginGetInboundRoutes not implemented")}
	}
	beginGetInboundRoutes := v.beginGetInboundRoutes.get(req)
	if beginGetInboundRoutes == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/virtualHubs/(?P<virtualHubName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/inboundRoutes`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.GetInboundRoutesParameters](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualHubNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualHubName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginGetInboundRoutes(req.Context(), resourceGroupNameParam, virtualHubNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginGetInboundRoutes = &respr
		v.beginGetInboundRoutes.add(req, beginGetInboundRoutes)
	}

	resp, err := server.PollerResponderNext(beginGetInboundRoutes, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginGetInboundRoutes.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginGetInboundRoutes) {
		v.beginGetInboundRoutes.remove(req)
	}

	return resp, nil
}

func (v *VirtualHubsServerTransport) dispatchBeginGetOutboundRoutes(req *http.Request) (*http.Response, error) {
	if v.srv.BeginGetOutboundRoutes == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginGetOutboundRoutes not implemented")}
	}
	beginGetOutboundRoutes := v.beginGetOutboundRoutes.get(req)
	if beginGetOutboundRoutes == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/virtualHubs/(?P<virtualHubName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/outboundRoutes`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.GetOutboundRoutesParameters](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualHubNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualHubName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginGetOutboundRoutes(req.Context(), resourceGroupNameParam, virtualHubNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginGetOutboundRoutes = &respr
		v.beginGetOutboundRoutes.add(req, beginGetOutboundRoutes)
	}

	resp, err := server.PollerResponderNext(beginGetOutboundRoutes, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginGetOutboundRoutes.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginGetOutboundRoutes) {
		v.beginGetOutboundRoutes.remove(req)
	}

	return resp, nil
}

func (v *VirtualHubsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if v.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := v.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/virtualHubs`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := v.srv.NewListPager(nil)
		newListPager = &resp
		v.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armnetwork.VirtualHubsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		v.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		v.newListPager.remove(req)
	}
	return resp, nil
}

func (v *VirtualHubsServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if v.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := v.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/virtualHubs`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := v.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		v.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armnetwork.VirtualHubsClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		v.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		v.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (v *VirtualHubsServerTransport) dispatchUpdateTags(req *http.Request) (*http.Response, error) {
	if v.srv.UpdateTags == nil {
		return nil, &nonRetriableError{errors.New("fake for method UpdateTags not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/virtualHubs/(?P<virtualHubName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armnetwork.TagsObject](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	virtualHubNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualHubName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := v.srv.UpdateTags(req.Context(), resourceGroupNameParam, virtualHubNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).VirtualHub, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to VirtualHubsServerTransport
var virtualHubsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
