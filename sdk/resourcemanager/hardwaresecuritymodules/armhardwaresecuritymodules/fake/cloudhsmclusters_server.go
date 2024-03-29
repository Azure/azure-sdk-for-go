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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hardwaresecuritymodules/armhardwaresecuritymodules/v2"
	"net/http"
	"net/url"
	"regexp"
)

// CloudHsmClustersServer is a fake server for instances of the armhardwaresecuritymodules.CloudHsmClustersClient type.
type CloudHsmClustersServer struct {
	// BeginCreateOrUpdate is the fake for method CloudHsmClustersClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, body armhardwaresecuritymodules.CloudHsmCluster, options *armhardwaresecuritymodules.CloudHsmClustersClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armhardwaresecuritymodules.CloudHsmClustersClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method CloudHsmClustersClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, options *armhardwaresecuritymodules.CloudHsmClustersClientBeginDeleteOptions) (resp azfake.PollerResponder[armhardwaresecuritymodules.CloudHsmClustersClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method CloudHsmClustersClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, options *armhardwaresecuritymodules.CloudHsmClustersClientGetOptions) (resp azfake.Responder[armhardwaresecuritymodules.CloudHsmClustersClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method CloudHsmClustersClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armhardwaresecuritymodules.CloudHsmClustersClientListByResourceGroupOptions) (resp azfake.PagerResponder[armhardwaresecuritymodules.CloudHsmClustersClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method CloudHsmClustersClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armhardwaresecuritymodules.CloudHsmClustersClientListBySubscriptionOptions) (resp azfake.PagerResponder[armhardwaresecuritymodules.CloudHsmClustersClientListBySubscriptionResponse])

	// BeginUpdate is the fake for method CloudHsmClustersClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, cloudHsmClusterName string, body armhardwaresecuritymodules.CloudHsmClusterPatchParameters, options *armhardwaresecuritymodules.CloudHsmClustersClientBeginUpdateOptions) (resp azfake.PollerResponder[armhardwaresecuritymodules.CloudHsmClustersClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewCloudHsmClustersServerTransport creates a new instance of CloudHsmClustersServerTransport with the provided implementation.
// The returned CloudHsmClustersServerTransport instance is connected to an instance of armhardwaresecuritymodules.CloudHsmClustersClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewCloudHsmClustersServerTransport(srv *CloudHsmClustersServer) *CloudHsmClustersServerTransport {
	return &CloudHsmClustersServerTransport{
		srv:                         srv,
		beginCreateOrUpdate:         newTracker[azfake.PollerResponder[armhardwaresecuritymodules.CloudHsmClustersClientCreateOrUpdateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armhardwaresecuritymodules.CloudHsmClustersClientDeleteResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armhardwaresecuritymodules.CloudHsmClustersClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armhardwaresecuritymodules.CloudHsmClustersClientListBySubscriptionResponse]](),
		beginUpdate:                 newTracker[azfake.PollerResponder[armhardwaresecuritymodules.CloudHsmClustersClientUpdateResponse]](),
	}
}

// CloudHsmClustersServerTransport connects instances of armhardwaresecuritymodules.CloudHsmClustersClient to instances of CloudHsmClustersServer.
// Don't use this type directly, use NewCloudHsmClustersServerTransport instead.
type CloudHsmClustersServerTransport struct {
	srv                         *CloudHsmClustersServer
	beginCreateOrUpdate         *tracker[azfake.PollerResponder[armhardwaresecuritymodules.CloudHsmClustersClientCreateOrUpdateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armhardwaresecuritymodules.CloudHsmClustersClientDeleteResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armhardwaresecuritymodules.CloudHsmClustersClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armhardwaresecuritymodules.CloudHsmClustersClientListBySubscriptionResponse]]
	beginUpdate                 *tracker[azfake.PollerResponder[armhardwaresecuritymodules.CloudHsmClustersClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for CloudHsmClustersServerTransport.
func (c *CloudHsmClustersServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "CloudHsmClustersClient.BeginCreateOrUpdate":
		resp, err = c.dispatchBeginCreateOrUpdate(req)
	case "CloudHsmClustersClient.BeginDelete":
		resp, err = c.dispatchBeginDelete(req)
	case "CloudHsmClustersClient.Get":
		resp, err = c.dispatchGet(req)
	case "CloudHsmClustersClient.NewListByResourceGroupPager":
		resp, err = c.dispatchNewListByResourceGroupPager(req)
	case "CloudHsmClustersClient.NewListBySubscriptionPager":
		resp, err = c.dispatchNewListBySubscriptionPager(req)
	case "CloudHsmClustersClient.BeginUpdate":
		resp, err = c.dispatchBeginUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *CloudHsmClustersServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := c.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/cloudHsmClusters/(?P<cloudHsmClusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhardwaresecuritymodules.CloudHsmCluster](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		cloudHsmClusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("cloudHsmClusterName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, cloudHsmClusterNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		c.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		c.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		c.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (c *CloudHsmClustersServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if c.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := c.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/cloudHsmClusters/(?P<cloudHsmClusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		cloudHsmClusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("cloudHsmClusterName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginDelete(req.Context(), resourceGroupNameParam, cloudHsmClusterNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		c.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		c.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		c.beginDelete.remove(req)
	}

	return resp, nil
}

func (c *CloudHsmClustersServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if c.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/cloudHsmClusters/(?P<cloudHsmClusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	cloudHsmClusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("cloudHsmClusterName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Get(req.Context(), resourceGroupNameParam, cloudHsmClusterNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CloudHsmCluster, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CloudHsmClustersServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := c.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/cloudHsmClusters`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		skiptokenUnescaped, err := url.QueryUnescape(qp.Get("$skiptoken"))
		if err != nil {
			return nil, err
		}
		skiptokenParam := getOptional(skiptokenUnescaped)
		var options *armhardwaresecuritymodules.CloudHsmClustersClientListByResourceGroupOptions
		if skiptokenParam != nil {
			options = &armhardwaresecuritymodules.CloudHsmClustersClientListByResourceGroupOptions{
				Skiptoken: skiptokenParam,
			}
		}
		resp := c.srv.NewListByResourceGroupPager(resourceGroupNameParam, options)
		newListByResourceGroupPager = &resp
		c.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armhardwaresecuritymodules.CloudHsmClustersClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		c.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (c *CloudHsmClustersServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := c.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/cloudHsmClusters`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		skiptokenUnescaped, err := url.QueryUnescape(qp.Get("$skiptoken"))
		if err != nil {
			return nil, err
		}
		skiptokenParam := getOptional(skiptokenUnescaped)
		var options *armhardwaresecuritymodules.CloudHsmClustersClientListBySubscriptionOptions
		if skiptokenParam != nil {
			options = &armhardwaresecuritymodules.CloudHsmClustersClientListBySubscriptionOptions{
				Skiptoken: skiptokenParam,
			}
		}
		resp := c.srv.NewListBySubscriptionPager(options)
		newListBySubscriptionPager = &resp
		c.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armhardwaresecuritymodules.CloudHsmClustersClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		c.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (c *CloudHsmClustersServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := c.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HardwareSecurityModules/cloudHsmClusters/(?P<cloudHsmClusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhardwaresecuritymodules.CloudHsmClusterPatchParameters](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		cloudHsmClusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("cloudHsmClusterName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginUpdate(req.Context(), resourceGroupNameParam, cloudHsmClusterNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		c.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		c.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		c.beginUpdate.remove(req)
	}

	return resp, nil
}
