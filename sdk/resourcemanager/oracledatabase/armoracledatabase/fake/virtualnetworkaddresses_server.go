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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/oracledatabase/armoracledatabase"
	"net/http"
	"net/url"
	"regexp"
)

// VirtualNetworkAddressesServer is a fake server for instances of the armoracledatabase.VirtualNetworkAddressesClient type.
type VirtualNetworkAddressesServer struct {
	// BeginCreateOrUpdate is the fake for method VirtualNetworkAddressesClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, cloudvmclustername string, virtualnetworkaddressname string, resource armoracledatabase.VirtualNetworkAddress, options *armoracledatabase.VirtualNetworkAddressesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armoracledatabase.VirtualNetworkAddressesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method VirtualNetworkAddressesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, cloudvmclustername string, virtualnetworkaddressname string, options *armoracledatabase.VirtualNetworkAddressesClientBeginDeleteOptions) (resp azfake.PollerResponder[armoracledatabase.VirtualNetworkAddressesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method VirtualNetworkAddressesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, cloudvmclustername string, virtualnetworkaddressname string, options *armoracledatabase.VirtualNetworkAddressesClientGetOptions) (resp azfake.Responder[armoracledatabase.VirtualNetworkAddressesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByCloudVMClusterPager is the fake for method VirtualNetworkAddressesClient.NewListByCloudVMClusterPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByCloudVMClusterPager func(resourceGroupName string, cloudvmclustername string, options *armoracledatabase.VirtualNetworkAddressesClientListByCloudVMClusterOptions) (resp azfake.PagerResponder[armoracledatabase.VirtualNetworkAddressesClientListByCloudVMClusterResponse])
}

// NewVirtualNetworkAddressesServerTransport creates a new instance of VirtualNetworkAddressesServerTransport with the provided implementation.
// The returned VirtualNetworkAddressesServerTransport instance is connected to an instance of armoracledatabase.VirtualNetworkAddressesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewVirtualNetworkAddressesServerTransport(srv *VirtualNetworkAddressesServer) *VirtualNetworkAddressesServerTransport {
	return &VirtualNetworkAddressesServerTransport{
		srv:                          srv,
		beginCreateOrUpdate:          newTracker[azfake.PollerResponder[armoracledatabase.VirtualNetworkAddressesClientCreateOrUpdateResponse]](),
		beginDelete:                  newTracker[azfake.PollerResponder[armoracledatabase.VirtualNetworkAddressesClientDeleteResponse]](),
		newListByCloudVMClusterPager: newTracker[azfake.PagerResponder[armoracledatabase.VirtualNetworkAddressesClientListByCloudVMClusterResponse]](),
	}
}

// VirtualNetworkAddressesServerTransport connects instances of armoracledatabase.VirtualNetworkAddressesClient to instances of VirtualNetworkAddressesServer.
// Don't use this type directly, use NewVirtualNetworkAddressesServerTransport instead.
type VirtualNetworkAddressesServerTransport struct {
	srv                          *VirtualNetworkAddressesServer
	beginCreateOrUpdate          *tracker[azfake.PollerResponder[armoracledatabase.VirtualNetworkAddressesClientCreateOrUpdateResponse]]
	beginDelete                  *tracker[azfake.PollerResponder[armoracledatabase.VirtualNetworkAddressesClientDeleteResponse]]
	newListByCloudVMClusterPager *tracker[azfake.PagerResponder[armoracledatabase.VirtualNetworkAddressesClientListByCloudVMClusterResponse]]
}

// Do implements the policy.Transporter interface for VirtualNetworkAddressesServerTransport.
func (v *VirtualNetworkAddressesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "VirtualNetworkAddressesClient.BeginCreateOrUpdate":
		resp, err = v.dispatchBeginCreateOrUpdate(req)
	case "VirtualNetworkAddressesClient.BeginDelete":
		resp, err = v.dispatchBeginDelete(req)
	case "VirtualNetworkAddressesClient.Get":
		resp, err = v.dispatchGet(req)
	case "VirtualNetworkAddressesClient.NewListByCloudVMClusterPager":
		resp, err = v.dispatchNewListByCloudVMClusterPager(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *VirtualNetworkAddressesServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if v.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := v.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Oracle\.Database/cloudVmClusters/(?P<cloudvmclustername>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/virtualNetworkAddresses/(?P<virtualnetworkaddressname>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armoracledatabase.VirtualNetworkAddress](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		cloudvmclusternameParam, err := url.PathUnescape(matches[regex.SubexpIndex("cloudvmclustername")])
		if err != nil {
			return nil, err
		}
		virtualnetworkaddressnameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualnetworkaddressname")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, cloudvmclusternameParam, virtualnetworkaddressnameParam, body, nil)
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

func (v *VirtualNetworkAddressesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if v.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := v.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Oracle\.Database/cloudVmClusters/(?P<cloudvmclustername>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/virtualNetworkAddresses/(?P<virtualnetworkaddressname>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		cloudvmclusternameParam, err := url.PathUnescape(matches[regex.SubexpIndex("cloudvmclustername")])
		if err != nil {
			return nil, err
		}
		virtualnetworkaddressnameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualnetworkaddressname")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginDelete(req.Context(), resourceGroupNameParam, cloudvmclusternameParam, virtualnetworkaddressnameParam, nil)
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

	if !contains([]int{http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		v.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		v.beginDelete.remove(req)
	}

	return resp, nil
}

func (v *VirtualNetworkAddressesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if v.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Oracle\.Database/cloudVmClusters/(?P<cloudvmclustername>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/virtualNetworkAddresses/(?P<virtualnetworkaddressname>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	cloudvmclusternameParam, err := url.PathUnescape(matches[regex.SubexpIndex("cloudvmclustername")])
	if err != nil {
		return nil, err
	}
	virtualnetworkaddressnameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualnetworkaddressname")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := v.srv.Get(req.Context(), resourceGroupNameParam, cloudvmclusternameParam, virtualnetworkaddressnameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).VirtualNetworkAddress, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (v *VirtualNetworkAddressesServerTransport) dispatchNewListByCloudVMClusterPager(req *http.Request) (*http.Response, error) {
	if v.srv.NewListByCloudVMClusterPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByCloudVMClusterPager not implemented")}
	}
	newListByCloudVMClusterPager := v.newListByCloudVMClusterPager.get(req)
	if newListByCloudVMClusterPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Oracle\.Database/cloudVmClusters/(?P<cloudvmclustername>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/virtualNetworkAddresses`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		cloudvmclusternameParam, err := url.PathUnescape(matches[regex.SubexpIndex("cloudvmclustername")])
		if err != nil {
			return nil, err
		}
		resp := v.srv.NewListByCloudVMClusterPager(resourceGroupNameParam, cloudvmclusternameParam, nil)
		newListByCloudVMClusterPager = &resp
		v.newListByCloudVMClusterPager.add(req, newListByCloudVMClusterPager)
		server.PagerResponderInjectNextLinks(newListByCloudVMClusterPager, req, func(page *armoracledatabase.VirtualNetworkAddressesClientListByCloudVMClusterResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByCloudVMClusterPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		v.newListByCloudVMClusterPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByCloudVMClusterPager) {
		v.newListByCloudVMClusterPager.remove(req)
	}
	return resp, nil
}