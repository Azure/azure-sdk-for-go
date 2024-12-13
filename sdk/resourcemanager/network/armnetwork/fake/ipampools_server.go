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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork/v6"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
)

// IpamPoolsServer is a fake server for instances of the armnetwork.IpamPoolsClient type.
type IpamPoolsServer struct {
	// BeginCreate is the fake for method IpamPoolsClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, resourceGroupName string, networkManagerName string, poolName string, body armnetwork.IpamPool, options *armnetwork.IpamPoolsClientBeginCreateOptions) (resp azfake.PollerResponder[armnetwork.IpamPoolsClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method IpamPoolsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, networkManagerName string, poolName string, options *armnetwork.IpamPoolsClientBeginDeleteOptions) (resp azfake.PollerResponder[armnetwork.IpamPoolsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method IpamPoolsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, networkManagerName string, poolName string, options *armnetwork.IpamPoolsClientGetOptions) (resp azfake.Responder[armnetwork.IpamPoolsClientGetResponse], errResp azfake.ErrorResponder)

	// GetPoolUsage is the fake for method IpamPoolsClient.GetPoolUsage
	// HTTP status codes to indicate success: http.StatusOK
	GetPoolUsage func(ctx context.Context, resourceGroupName string, networkManagerName string, poolName string, options *armnetwork.IpamPoolsClientGetPoolUsageOptions) (resp azfake.Responder[armnetwork.IpamPoolsClientGetPoolUsageResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method IpamPoolsClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, networkManagerName string, options *armnetwork.IpamPoolsClientListOptions) (resp azfake.PagerResponder[armnetwork.IpamPoolsClientListResponse])

	// NewListAssociatedResourcesPager is the fake for method IpamPoolsClient.NewListAssociatedResourcesPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListAssociatedResourcesPager func(resourceGroupName string, networkManagerName string, poolName string, options *armnetwork.IpamPoolsClientListAssociatedResourcesOptions) (resp azfake.PagerResponder[armnetwork.IpamPoolsClientListAssociatedResourcesResponse])

	// Update is the fake for method IpamPoolsClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, networkManagerName string, poolName string, options *armnetwork.IpamPoolsClientUpdateOptions) (resp azfake.Responder[armnetwork.IpamPoolsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewIpamPoolsServerTransport creates a new instance of IpamPoolsServerTransport with the provided implementation.
// The returned IpamPoolsServerTransport instance is connected to an instance of armnetwork.IpamPoolsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewIpamPoolsServerTransport(srv *IpamPoolsServer) *IpamPoolsServerTransport {
	return &IpamPoolsServerTransport{
		srv:                             srv,
		beginCreate:                     newTracker[azfake.PollerResponder[armnetwork.IpamPoolsClientCreateResponse]](),
		beginDelete:                     newTracker[azfake.PollerResponder[armnetwork.IpamPoolsClientDeleteResponse]](),
		newListPager:                    newTracker[azfake.PagerResponder[armnetwork.IpamPoolsClientListResponse]](),
		newListAssociatedResourcesPager: newTracker[azfake.PagerResponder[armnetwork.IpamPoolsClientListAssociatedResourcesResponse]](),
	}
}

// IpamPoolsServerTransport connects instances of armnetwork.IpamPoolsClient to instances of IpamPoolsServer.
// Don't use this type directly, use NewIpamPoolsServerTransport instead.
type IpamPoolsServerTransport struct {
	srv                             *IpamPoolsServer
	beginCreate                     *tracker[azfake.PollerResponder[armnetwork.IpamPoolsClientCreateResponse]]
	beginDelete                     *tracker[azfake.PollerResponder[armnetwork.IpamPoolsClientDeleteResponse]]
	newListPager                    *tracker[azfake.PagerResponder[armnetwork.IpamPoolsClientListResponse]]
	newListAssociatedResourcesPager *tracker[azfake.PagerResponder[armnetwork.IpamPoolsClientListAssociatedResourcesResponse]]
}

// Do implements the policy.Transporter interface for IpamPoolsServerTransport.
func (i *IpamPoolsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "IpamPoolsClient.BeginCreate":
		resp, err = i.dispatchBeginCreate(req)
	case "IpamPoolsClient.BeginDelete":
		resp, err = i.dispatchBeginDelete(req)
	case "IpamPoolsClient.Get":
		resp, err = i.dispatchGet(req)
	case "IpamPoolsClient.GetPoolUsage":
		resp, err = i.dispatchGetPoolUsage(req)
	case "IpamPoolsClient.NewListPager":
		resp, err = i.dispatchNewListPager(req)
	case "IpamPoolsClient.NewListAssociatedResourcesPager":
		resp, err = i.dispatchNewListAssociatedResourcesPager(req)
	case "IpamPoolsClient.Update":
		resp, err = i.dispatchUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (i *IpamPoolsServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if i.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := i.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkManagers/(?P<networkManagerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/ipamPools/(?P<poolName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armnetwork.IpamPool](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		networkManagerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkManagerName")])
		if err != nil {
			return nil, err
		}
		poolNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("poolName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginCreate(req.Context(), resourceGroupNameParam, networkManagerNameParam, poolNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		i.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		i.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		i.beginCreate.remove(req)
	}

	return resp, nil
}

func (i *IpamPoolsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if i.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := i.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkManagers/(?P<networkManagerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/ipamPools/(?P<poolName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		networkManagerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkManagerName")])
		if err != nil {
			return nil, err
		}
		poolNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("poolName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := i.srv.BeginDelete(req.Context(), resourceGroupNameParam, networkManagerNameParam, poolNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		i.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		i.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		i.beginDelete.remove(req)
	}

	return resp, nil
}

func (i *IpamPoolsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if i.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkManagers/(?P<networkManagerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/ipamPools/(?P<poolName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	networkManagerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkManagerName")])
	if err != nil {
		return nil, err
	}
	poolNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("poolName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Get(req.Context(), resourceGroupNameParam, networkManagerNameParam, poolNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).IpamPool, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *IpamPoolsServerTransport) dispatchGetPoolUsage(req *http.Request) (*http.Response, error) {
	if i.srv.GetPoolUsage == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetPoolUsage not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkManagers/(?P<networkManagerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/ipamPools/(?P<poolName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/getPoolUsage`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	networkManagerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkManagerName")])
	if err != nil {
		return nil, err
	}
	poolNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("poolName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.GetPoolUsage(req.Context(), resourceGroupNameParam, networkManagerNameParam, poolNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).PoolUsage, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *IpamPoolsServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if i.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := i.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkManagers/(?P<networkManagerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/ipamPools`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		networkManagerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkManagerName")])
		if err != nil {
			return nil, err
		}
		skipTokenUnescaped, err := url.QueryUnescape(qp.Get("skipToken"))
		if err != nil {
			return nil, err
		}
		skipTokenParam := getOptional(skipTokenUnescaped)
		skipUnescaped, err := url.QueryUnescape(qp.Get("skip"))
		if err != nil {
			return nil, err
		}
		skipParam, err := parseOptional(skipUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		topUnescaped, err := url.QueryUnescape(qp.Get("top"))
		if err != nil {
			return nil, err
		}
		topParam, err := parseOptional(topUnescaped, func(v string) (int32, error) {
			p, parseErr := strconv.ParseInt(v, 10, 32)
			if parseErr != nil {
				return 0, parseErr
			}
			return int32(p), nil
		})
		if err != nil {
			return nil, err
		}
		sortKeyUnescaped, err := url.QueryUnescape(qp.Get("sortKey"))
		if err != nil {
			return nil, err
		}
		sortKeyParam := getOptional(sortKeyUnescaped)
		sortValueUnescaped, err := url.QueryUnescape(qp.Get("sortValue"))
		if err != nil {
			return nil, err
		}
		sortValueParam := getOptional(sortValueUnescaped)
		var options *armnetwork.IpamPoolsClientListOptions
		if skipTokenParam != nil || skipParam != nil || topParam != nil || sortKeyParam != nil || sortValueParam != nil {
			options = &armnetwork.IpamPoolsClientListOptions{
				SkipToken: skipTokenParam,
				Skip:      skipParam,
				Top:       topParam,
				SortKey:   sortKeyParam,
				SortValue: sortValueParam,
			}
		}
		resp := i.srv.NewListPager(resourceGroupNameParam, networkManagerNameParam, options)
		newListPager = &resp
		i.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armnetwork.IpamPoolsClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		i.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		i.newListPager.remove(req)
	}
	return resp, nil
}

func (i *IpamPoolsServerTransport) dispatchNewListAssociatedResourcesPager(req *http.Request) (*http.Response, error) {
	if i.srv.NewListAssociatedResourcesPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListAssociatedResourcesPager not implemented")}
	}
	newListAssociatedResourcesPager := i.newListAssociatedResourcesPager.get(req)
	if newListAssociatedResourcesPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkManagers/(?P<networkManagerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/ipamPools/(?P<poolName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/listAssociatedResources`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		networkManagerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkManagerName")])
		if err != nil {
			return nil, err
		}
		poolNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("poolName")])
		if err != nil {
			return nil, err
		}
		resp := i.srv.NewListAssociatedResourcesPager(resourceGroupNameParam, networkManagerNameParam, poolNameParam, nil)
		newListAssociatedResourcesPager = &resp
		i.newListAssociatedResourcesPager.add(req, newListAssociatedResourcesPager)
		server.PagerResponderInjectNextLinks(newListAssociatedResourcesPager, req, func(page *armnetwork.IpamPoolsClientListAssociatedResourcesResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListAssociatedResourcesPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		i.newListAssociatedResourcesPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListAssociatedResourcesPager) {
		i.newListAssociatedResourcesPager.remove(req)
	}
	return resp, nil
}

func (i *IpamPoolsServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if i.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Network/networkManagers/(?P<networkManagerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/ipamPools/(?P<poolName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armnetwork.IpamPoolUpdate](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	networkManagerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("networkManagerName")])
	if err != nil {
		return nil, err
	}
	poolNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("poolName")])
	if err != nil {
		return nil, err
	}
	var options *armnetwork.IpamPoolsClientUpdateOptions
	if !reflect.ValueOf(body).IsZero() {
		options = &armnetwork.IpamPoolsClientUpdateOptions{
			Body: &body,
		}
	}
	respr, errRespr := i.srv.Update(req.Context(), resourceGroupNameParam, networkManagerNameParam, poolNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).IpamPool, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}