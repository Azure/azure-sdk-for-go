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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/cosmos/armcosmos/v3"
	"net/http"
	"net/url"
	"regexp"
)

// CassandraClustersServer is a fake server for instances of the armcosmos.CassandraClustersClient type.
type CassandraClustersServer struct {
	// BeginCreateUpdate is the fake for method CassandraClustersClient.BeginCreateUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateUpdate func(ctx context.Context, resourceGroupName string, clusterName string, body armcosmos.ClusterResource, options *armcosmos.CassandraClustersClientBeginCreateUpdateOptions) (resp azfake.PollerResponder[armcosmos.CassandraClustersClientCreateUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDeallocate is the fake for method CassandraClustersClient.BeginDeallocate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDeallocate func(ctx context.Context, resourceGroupName string, clusterName string, options *armcosmos.CassandraClustersClientBeginDeallocateOptions) (resp azfake.PollerResponder[armcosmos.CassandraClustersClientDeallocateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method CassandraClustersClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, clusterName string, options *armcosmos.CassandraClustersClientBeginDeleteOptions) (resp azfake.PollerResponder[armcosmos.CassandraClustersClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method CassandraClustersClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, clusterName string, options *armcosmos.CassandraClustersClientGetOptions) (resp azfake.Responder[armcosmos.CassandraClustersClientGetResponse], errResp azfake.ErrorResponder)

	// BeginInvokeCommand is the fake for method CassandraClustersClient.BeginInvokeCommand
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginInvokeCommand func(ctx context.Context, resourceGroupName string, clusterName string, body armcosmos.CommandPostBody, options *armcosmos.CassandraClustersClientBeginInvokeCommandOptions) (resp azfake.PollerResponder[armcosmos.CassandraClustersClientInvokeCommandResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method CassandraClustersClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armcosmos.CassandraClustersClientListByResourceGroupOptions) (resp azfake.PagerResponder[armcosmos.CassandraClustersClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method CassandraClustersClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armcosmos.CassandraClustersClientListBySubscriptionOptions) (resp azfake.PagerResponder[armcosmos.CassandraClustersClientListBySubscriptionResponse])

	// BeginStart is the fake for method CassandraClustersClient.BeginStart
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginStart func(ctx context.Context, resourceGroupName string, clusterName string, options *armcosmos.CassandraClustersClientBeginStartOptions) (resp azfake.PollerResponder[armcosmos.CassandraClustersClientStartResponse], errResp azfake.ErrorResponder)

	// Status is the fake for method CassandraClustersClient.Status
	// HTTP status codes to indicate success: http.StatusOK
	Status func(ctx context.Context, resourceGroupName string, clusterName string, options *armcosmos.CassandraClustersClientStatusOptions) (resp azfake.Responder[armcosmos.CassandraClustersClientStatusResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method CassandraClustersClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, clusterName string, body armcosmos.ClusterResource, options *armcosmos.CassandraClustersClientBeginUpdateOptions) (resp azfake.PollerResponder[armcosmos.CassandraClustersClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewCassandraClustersServerTransport creates a new instance of CassandraClustersServerTransport with the provided implementation.
// The returned CassandraClustersServerTransport instance is connected to an instance of armcosmos.CassandraClustersClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewCassandraClustersServerTransport(srv *CassandraClustersServer) *CassandraClustersServerTransport {
	return &CassandraClustersServerTransport{
		srv:                         srv,
		beginCreateUpdate:           newTracker[azfake.PollerResponder[armcosmos.CassandraClustersClientCreateUpdateResponse]](),
		beginDeallocate:             newTracker[azfake.PollerResponder[armcosmos.CassandraClustersClientDeallocateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armcosmos.CassandraClustersClientDeleteResponse]](),
		beginInvokeCommand:          newTracker[azfake.PollerResponder[armcosmos.CassandraClustersClientInvokeCommandResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armcosmos.CassandraClustersClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armcosmos.CassandraClustersClientListBySubscriptionResponse]](),
		beginStart:                  newTracker[azfake.PollerResponder[armcosmos.CassandraClustersClientStartResponse]](),
		beginUpdate:                 newTracker[azfake.PollerResponder[armcosmos.CassandraClustersClientUpdateResponse]](),
	}
}

// CassandraClustersServerTransport connects instances of armcosmos.CassandraClustersClient to instances of CassandraClustersServer.
// Don't use this type directly, use NewCassandraClustersServerTransport instead.
type CassandraClustersServerTransport struct {
	srv                         *CassandraClustersServer
	beginCreateUpdate           *tracker[azfake.PollerResponder[armcosmos.CassandraClustersClientCreateUpdateResponse]]
	beginDeallocate             *tracker[azfake.PollerResponder[armcosmos.CassandraClustersClientDeallocateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armcosmos.CassandraClustersClientDeleteResponse]]
	beginInvokeCommand          *tracker[azfake.PollerResponder[armcosmos.CassandraClustersClientInvokeCommandResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armcosmos.CassandraClustersClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armcosmos.CassandraClustersClientListBySubscriptionResponse]]
	beginStart                  *tracker[azfake.PollerResponder[armcosmos.CassandraClustersClientStartResponse]]
	beginUpdate                 *tracker[azfake.PollerResponder[armcosmos.CassandraClustersClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for CassandraClustersServerTransport.
func (c *CassandraClustersServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return c.dispatchToMethodFake(req, method)
}

func (c *CassandraClustersServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if cassandraClustersServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = cassandraClustersServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "CassandraClustersClient.BeginCreateUpdate":
				res.resp, res.err = c.dispatchBeginCreateUpdate(req)
			case "CassandraClustersClient.BeginDeallocate":
				res.resp, res.err = c.dispatchBeginDeallocate(req)
			case "CassandraClustersClient.BeginDelete":
				res.resp, res.err = c.dispatchBeginDelete(req)
			case "CassandraClustersClient.Get":
				res.resp, res.err = c.dispatchGet(req)
			case "CassandraClustersClient.BeginInvokeCommand":
				res.resp, res.err = c.dispatchBeginInvokeCommand(req)
			case "CassandraClustersClient.NewListByResourceGroupPager":
				res.resp, res.err = c.dispatchNewListByResourceGroupPager(req)
			case "CassandraClustersClient.NewListBySubscriptionPager":
				res.resp, res.err = c.dispatchNewListBySubscriptionPager(req)
			case "CassandraClustersClient.BeginStart":
				res.resp, res.err = c.dispatchBeginStart(req)
			case "CassandraClustersClient.Status":
				res.resp, res.err = c.dispatchStatus(req)
			case "CassandraClustersClient.BeginUpdate":
				res.resp, res.err = c.dispatchBeginUpdate(req)
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

func (c *CassandraClustersServerTransport) dispatchBeginCreateUpdate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginCreateUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateUpdate not implemented")}
	}
	beginCreateUpdate := c.beginCreateUpdate.get(req)
	if beginCreateUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/cassandraClusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcosmos.ClusterResource](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginCreateUpdate(req.Context(), resourceGroupNameParam, clusterNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateUpdate = &respr
		c.beginCreateUpdate.add(req, beginCreateUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		c.beginCreateUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateUpdate) {
		c.beginCreateUpdate.remove(req)
	}

	return resp, nil
}

func (c *CassandraClustersServerTransport) dispatchBeginDeallocate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginDeallocate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDeallocate not implemented")}
	}
	beginDeallocate := c.beginDeallocate.get(req)
	if beginDeallocate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/cassandraClusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/deallocate`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginDeallocate(req.Context(), resourceGroupNameParam, clusterNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDeallocate = &respr
		c.beginDeallocate.add(req, beginDeallocate)
	}

	resp, err := server.PollerResponderNext(beginDeallocate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		c.beginDeallocate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDeallocate) {
		c.beginDeallocate.remove(req)
	}

	return resp, nil
}

func (c *CassandraClustersServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if c.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := c.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/cassandraClusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginDelete(req.Context(), resourceGroupNameParam, clusterNameParam, nil)
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

func (c *CassandraClustersServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if c.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/cassandraClusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Get(req.Context(), resourceGroupNameParam, clusterNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ClusterResource, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CassandraClustersServerTransport) dispatchBeginInvokeCommand(req *http.Request) (*http.Response, error) {
	if c.srv.BeginInvokeCommand == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginInvokeCommand not implemented")}
	}
	beginInvokeCommand := c.beginInvokeCommand.get(req)
	if beginInvokeCommand == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/cassandraClusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/invokeCommand`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcosmos.CommandPostBody](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginInvokeCommand(req.Context(), resourceGroupNameParam, clusterNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginInvokeCommand = &respr
		c.beginInvokeCommand.add(req, beginInvokeCommand)
	}

	resp, err := server.PollerResponderNext(beginInvokeCommand, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		c.beginInvokeCommand.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginInvokeCommand) {
		c.beginInvokeCommand.remove(req)
	}

	return resp, nil
}

func (c *CassandraClustersServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := c.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/cassandraClusters`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := c.srv.NewListByResourceGroupPager(resourceGroupNameParam, nil)
		newListByResourceGroupPager = &resp
		c.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
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

func (c *CassandraClustersServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := c.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/cassandraClusters`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := c.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		c.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
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

func (c *CassandraClustersServerTransport) dispatchBeginStart(req *http.Request) (*http.Response, error) {
	if c.srv.BeginStart == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginStart not implemented")}
	}
	beginStart := c.beginStart.get(req)
	if beginStart == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/cassandraClusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/start`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginStart(req.Context(), resourceGroupNameParam, clusterNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginStart = &respr
		c.beginStart.add(req, beginStart)
	}

	resp, err := server.PollerResponderNext(beginStart, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		c.beginStart.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginStart) {
		c.beginStart.remove(req)
	}

	return resp, nil
}

func (c *CassandraClustersServerTransport) dispatchStatus(req *http.Request) (*http.Response, error) {
	if c.srv.Status == nil {
		return nil, &nonRetriableError{errors.New("fake for method Status not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/cassandraClusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/status`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.Status(req.Context(), resourceGroupNameParam, clusterNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CassandraClusterPublicStatus, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *CassandraClustersServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := c.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.DocumentDB/cassandraClusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armcosmos.ClusterResource](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		clusterNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("clusterName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginUpdate(req.Context(), resourceGroupNameParam, clusterNameParam, body, nil)
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

// set this to conditionally intercept incoming requests to CassandraClustersServerTransport
var cassandraClustersServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
