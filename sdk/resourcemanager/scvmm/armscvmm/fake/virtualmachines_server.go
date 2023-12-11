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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/scvmm/armscvmm"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
)

// VirtualMachinesServer is a fake server for instances of the armscvmm.VirtualMachinesClient type.
type VirtualMachinesServer struct {
	// BeginCreateCheckpoint is the fake for method VirtualMachinesClient.BeginCreateCheckpoint
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreateCheckpoint func(ctx context.Context, resourceGroupName string, virtualMachineName string, options *armscvmm.VirtualMachinesClientBeginCreateCheckpointOptions) (resp azfake.PollerResponder[armscvmm.VirtualMachinesClientCreateCheckpointResponse], errResp azfake.ErrorResponder)

	// BeginCreateOrUpdate is the fake for method VirtualMachinesClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, virtualMachineName string, body armscvmm.VirtualMachine, options *armscvmm.VirtualMachinesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armscvmm.VirtualMachinesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method VirtualMachinesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, virtualMachineName string, options *armscvmm.VirtualMachinesClientBeginDeleteOptions) (resp azfake.PollerResponder[armscvmm.VirtualMachinesClientDeleteResponse], errResp azfake.ErrorResponder)

	// BeginDeleteCheckpoint is the fake for method VirtualMachinesClient.BeginDeleteCheckpoint
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginDeleteCheckpoint func(ctx context.Context, resourceGroupName string, virtualMachineName string, options *armscvmm.VirtualMachinesClientBeginDeleteCheckpointOptions) (resp azfake.PollerResponder[armscvmm.VirtualMachinesClientDeleteCheckpointResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method VirtualMachinesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, virtualMachineName string, options *armscvmm.VirtualMachinesClientGetOptions) (resp azfake.Responder[armscvmm.VirtualMachinesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method VirtualMachinesClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armscvmm.VirtualMachinesClientListByResourceGroupOptions) (resp azfake.PagerResponder[armscvmm.VirtualMachinesClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method VirtualMachinesClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armscvmm.VirtualMachinesClientListBySubscriptionOptions) (resp azfake.PagerResponder[armscvmm.VirtualMachinesClientListBySubscriptionResponse])

	// BeginRestart is the fake for method VirtualMachinesClient.BeginRestart
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginRestart func(ctx context.Context, resourceGroupName string, virtualMachineName string, options *armscvmm.VirtualMachinesClientBeginRestartOptions) (resp azfake.PollerResponder[armscvmm.VirtualMachinesClientRestartResponse], errResp azfake.ErrorResponder)

	// BeginRestoreCheckpoint is the fake for method VirtualMachinesClient.BeginRestoreCheckpoint
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginRestoreCheckpoint func(ctx context.Context, resourceGroupName string, virtualMachineName string, options *armscvmm.VirtualMachinesClientBeginRestoreCheckpointOptions) (resp azfake.PollerResponder[armscvmm.VirtualMachinesClientRestoreCheckpointResponse], errResp azfake.ErrorResponder)

	// BeginStart is the fake for method VirtualMachinesClient.BeginStart
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginStart func(ctx context.Context, resourceGroupName string, virtualMachineName string, options *armscvmm.VirtualMachinesClientBeginStartOptions) (resp azfake.PollerResponder[armscvmm.VirtualMachinesClientStartResponse], errResp azfake.ErrorResponder)

	// BeginStop is the fake for method VirtualMachinesClient.BeginStop
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginStop func(ctx context.Context, resourceGroupName string, virtualMachineName string, options *armscvmm.VirtualMachinesClientBeginStopOptions) (resp azfake.PollerResponder[armscvmm.VirtualMachinesClientStopResponse], errResp azfake.ErrorResponder)

	// BeginUpdate is the fake for method VirtualMachinesClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, virtualMachineName string, body armscvmm.VirtualMachineUpdate, options *armscvmm.VirtualMachinesClientBeginUpdateOptions) (resp azfake.PollerResponder[armscvmm.VirtualMachinesClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewVirtualMachinesServerTransport creates a new instance of VirtualMachinesServerTransport with the provided implementation.
// The returned VirtualMachinesServerTransport instance is connected to an instance of armscvmm.VirtualMachinesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewVirtualMachinesServerTransport(srv *VirtualMachinesServer) *VirtualMachinesServerTransport {
	return &VirtualMachinesServerTransport{
		srv:                         srv,
		beginCreateCheckpoint:       newTracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientCreateCheckpointResponse]](),
		beginCreateOrUpdate:         newTracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientCreateOrUpdateResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientDeleteResponse]](),
		beginDeleteCheckpoint:       newTracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientDeleteCheckpointResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armscvmm.VirtualMachinesClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armscvmm.VirtualMachinesClientListBySubscriptionResponse]](),
		beginRestart:                newTracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientRestartResponse]](),
		beginRestoreCheckpoint:      newTracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientRestoreCheckpointResponse]](),
		beginStart:                  newTracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientStartResponse]](),
		beginStop:                   newTracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientStopResponse]](),
		beginUpdate:                 newTracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientUpdateResponse]](),
	}
}

// VirtualMachinesServerTransport connects instances of armscvmm.VirtualMachinesClient to instances of VirtualMachinesServer.
// Don't use this type directly, use NewVirtualMachinesServerTransport instead.
type VirtualMachinesServerTransport struct {
	srv                         *VirtualMachinesServer
	beginCreateCheckpoint       *tracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientCreateCheckpointResponse]]
	beginCreateOrUpdate         *tracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientCreateOrUpdateResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientDeleteResponse]]
	beginDeleteCheckpoint       *tracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientDeleteCheckpointResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armscvmm.VirtualMachinesClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armscvmm.VirtualMachinesClientListBySubscriptionResponse]]
	beginRestart                *tracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientRestartResponse]]
	beginRestoreCheckpoint      *tracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientRestoreCheckpointResponse]]
	beginStart                  *tracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientStartResponse]]
	beginStop                   *tracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientStopResponse]]
	beginUpdate                 *tracker[azfake.PollerResponder[armscvmm.VirtualMachinesClientUpdateResponse]]
}

// Do implements the policy.Transporter interface for VirtualMachinesServerTransport.
func (v *VirtualMachinesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "VirtualMachinesClient.BeginCreateCheckpoint":
		resp, err = v.dispatchBeginCreateCheckpoint(req)
	case "VirtualMachinesClient.BeginCreateOrUpdate":
		resp, err = v.dispatchBeginCreateOrUpdate(req)
	case "VirtualMachinesClient.BeginDelete":
		resp, err = v.dispatchBeginDelete(req)
	case "VirtualMachinesClient.BeginDeleteCheckpoint":
		resp, err = v.dispatchBeginDeleteCheckpoint(req)
	case "VirtualMachinesClient.Get":
		resp, err = v.dispatchGet(req)
	case "VirtualMachinesClient.NewListByResourceGroupPager":
		resp, err = v.dispatchNewListByResourceGroupPager(req)
	case "VirtualMachinesClient.NewListBySubscriptionPager":
		resp, err = v.dispatchNewListBySubscriptionPager(req)
	case "VirtualMachinesClient.BeginRestart":
		resp, err = v.dispatchBeginRestart(req)
	case "VirtualMachinesClient.BeginRestoreCheckpoint":
		resp, err = v.dispatchBeginRestoreCheckpoint(req)
	case "VirtualMachinesClient.BeginStart":
		resp, err = v.dispatchBeginStart(req)
	case "VirtualMachinesClient.BeginStop":
		resp, err = v.dispatchBeginStop(req)
	case "VirtualMachinesClient.BeginUpdate":
		resp, err = v.dispatchBeginUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (v *VirtualMachinesServerTransport) dispatchBeginCreateCheckpoint(req *http.Request) (*http.Response, error) {
	if v.srv.BeginCreateCheckpoint == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateCheckpoint not implemented")}
	}
	beginCreateCheckpoint := v.beginCreateCheckpoint.get(req)
	if beginCreateCheckpoint == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines/(?P<virtualMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/createCheckpoint`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armscvmm.VirtualMachineCreateCheckpoint](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualMachineName")])
		if err != nil {
			return nil, err
		}
		var options *armscvmm.VirtualMachinesClientBeginCreateCheckpointOptions
		if !reflect.ValueOf(body).IsZero() {
			options = &armscvmm.VirtualMachinesClientBeginCreateCheckpointOptions{
				Body: &body,
			}
		}
		respr, errRespr := v.srv.BeginCreateCheckpoint(req.Context(), resourceGroupNameParam, virtualMachineNameParam, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateCheckpoint = &respr
		v.beginCreateCheckpoint.add(req, beginCreateCheckpoint)
	}

	resp, err := server.PollerResponderNext(beginCreateCheckpoint, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginCreateCheckpoint.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateCheckpoint) {
		v.beginCreateCheckpoint.remove(req)
	}

	return resp, nil
}

func (v *VirtualMachinesServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if v.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := v.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines/(?P<virtualMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armscvmm.VirtualMachine](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualMachineName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, virtualMachineNameParam, body, nil)
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

func (v *VirtualMachinesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if v.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := v.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines/(?P<virtualMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
		virtualMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualMachineName")])
		if err != nil {
			return nil, err
		}
		retainUnescaped, err := url.QueryUnescape(qp.Get("retain"))
		if err != nil {
			return nil, err
		}
		retainParam, err := parseOptional(retainUnescaped, strconv.ParseBool)
		if err != nil {
			return nil, err
		}
		forceUnescaped, err := url.QueryUnescape(qp.Get("force"))
		if err != nil {
			return nil, err
		}
		forceParam, err := parseOptional(forceUnescaped, strconv.ParseBool)
		if err != nil {
			return nil, err
		}
		var options *armscvmm.VirtualMachinesClientBeginDeleteOptions
		if retainParam != nil || forceParam != nil {
			options = &armscvmm.VirtualMachinesClientBeginDeleteOptions{
				Retain: retainParam,
				Force:  forceParam,
			}
		}
		respr, errRespr := v.srv.BeginDelete(req.Context(), resourceGroupNameParam, virtualMachineNameParam, options)
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

func (v *VirtualMachinesServerTransport) dispatchBeginDeleteCheckpoint(req *http.Request) (*http.Response, error) {
	if v.srv.BeginDeleteCheckpoint == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDeleteCheckpoint not implemented")}
	}
	beginDeleteCheckpoint := v.beginDeleteCheckpoint.get(req)
	if beginDeleteCheckpoint == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines/(?P<virtualMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/deleteCheckpoint`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armscvmm.VirtualMachineDeleteCheckpoint](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualMachineName")])
		if err != nil {
			return nil, err
		}
		var options *armscvmm.VirtualMachinesClientBeginDeleteCheckpointOptions
		if !reflect.ValueOf(body).IsZero() {
			options = &armscvmm.VirtualMachinesClientBeginDeleteCheckpointOptions{
				Body: &body,
			}
		}
		respr, errRespr := v.srv.BeginDeleteCheckpoint(req.Context(), resourceGroupNameParam, virtualMachineNameParam, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDeleteCheckpoint = &respr
		v.beginDeleteCheckpoint.add(req, beginDeleteCheckpoint)
	}

	resp, err := server.PollerResponderNext(beginDeleteCheckpoint, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginDeleteCheckpoint.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDeleteCheckpoint) {
		v.beginDeleteCheckpoint.remove(req)
	}

	return resp, nil
}

func (v *VirtualMachinesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if v.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines/(?P<virtualMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	virtualMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualMachineName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := v.srv.Get(req.Context(), resourceGroupNameParam, virtualMachineNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).VirtualMachine, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (v *VirtualMachinesServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if v.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := v.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines`
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
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armscvmm.VirtualMachinesClientListByResourceGroupResponse, createLink func() string) {
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

func (v *VirtualMachinesServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if v.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := v.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := v.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		v.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armscvmm.VirtualMachinesClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		v.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		v.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (v *VirtualMachinesServerTransport) dispatchBeginRestart(req *http.Request) (*http.Response, error) {
	if v.srv.BeginRestart == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginRestart not implemented")}
	}
	beginRestart := v.beginRestart.get(req)
	if beginRestart == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines/(?P<virtualMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/restart`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualMachineName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginRestart(req.Context(), resourceGroupNameParam, virtualMachineNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginRestart = &respr
		v.beginRestart.add(req, beginRestart)
	}

	resp, err := server.PollerResponderNext(beginRestart, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginRestart.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginRestart) {
		v.beginRestart.remove(req)
	}

	return resp, nil
}

func (v *VirtualMachinesServerTransport) dispatchBeginRestoreCheckpoint(req *http.Request) (*http.Response, error) {
	if v.srv.BeginRestoreCheckpoint == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginRestoreCheckpoint not implemented")}
	}
	beginRestoreCheckpoint := v.beginRestoreCheckpoint.get(req)
	if beginRestoreCheckpoint == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines/(?P<virtualMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/restoreCheckpoint`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armscvmm.VirtualMachineRestoreCheckpoint](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualMachineName")])
		if err != nil {
			return nil, err
		}
		var options *armscvmm.VirtualMachinesClientBeginRestoreCheckpointOptions
		if !reflect.ValueOf(body).IsZero() {
			options = &armscvmm.VirtualMachinesClientBeginRestoreCheckpointOptions{
				Body: &body,
			}
		}
		respr, errRespr := v.srv.BeginRestoreCheckpoint(req.Context(), resourceGroupNameParam, virtualMachineNameParam, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginRestoreCheckpoint = &respr
		v.beginRestoreCheckpoint.add(req, beginRestoreCheckpoint)
	}

	resp, err := server.PollerResponderNext(beginRestoreCheckpoint, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginRestoreCheckpoint.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginRestoreCheckpoint) {
		v.beginRestoreCheckpoint.remove(req)
	}

	return resp, nil
}

func (v *VirtualMachinesServerTransport) dispatchBeginStart(req *http.Request) (*http.Response, error) {
	if v.srv.BeginStart == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginStart not implemented")}
	}
	beginStart := v.beginStart.get(req)
	if beginStart == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines/(?P<virtualMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/start`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualMachineName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginStart(req.Context(), resourceGroupNameParam, virtualMachineNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginStart = &respr
		v.beginStart.add(req, beginStart)
	}

	resp, err := server.PollerResponderNext(beginStart, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginStart.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginStart) {
		v.beginStart.remove(req)
	}

	return resp, nil
}

func (v *VirtualMachinesServerTransport) dispatchBeginStop(req *http.Request) (*http.Response, error) {
	if v.srv.BeginStop == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginStop not implemented")}
	}
	beginStop := v.beginStop.get(req)
	if beginStop == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines/(?P<virtualMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/stop`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armscvmm.StopVirtualMachineOptions](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualMachineName")])
		if err != nil {
			return nil, err
		}
		var options *armscvmm.VirtualMachinesClientBeginStopOptions
		if !reflect.ValueOf(body).IsZero() {
			options = &armscvmm.VirtualMachinesClientBeginStopOptions{
				Body: &body,
			}
		}
		respr, errRespr := v.srv.BeginStop(req.Context(), resourceGroupNameParam, virtualMachineNameParam, options)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginStop = &respr
		v.beginStop.add(req, beginStop)
	}

	resp, err := server.PollerResponderNext(beginStop, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		v.beginStop.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginStop) {
		v.beginStop.remove(req)
	}

	return resp, nil
}

func (v *VirtualMachinesServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if v.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := v.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.ScVmm/virtualMachines/(?P<virtualMachineName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armscvmm.VirtualMachineUpdate](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		virtualMachineNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("virtualMachineName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := v.srv.BeginUpdate(req.Context(), resourceGroupNameParam, virtualMachineNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		v.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated, http.StatusAccepted}, resp.StatusCode) {
		v.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		v.beginUpdate.remove(req)
	}

	return resp, nil
}