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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/machinelearning/armmachinelearning/v4"
	"net/http"
	"net/url"
	"regexp"
)

// RegistriesServer is a fake server for instances of the armmachinelearning.RegistriesClient type.
type RegistriesServer struct {
	// BeginCreateOrUpdate is the fake for method RegistriesClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreateOrUpdate func(ctx context.Context, resourceGroupName string, registryName string, body armmachinelearning.Registry, options *armmachinelearning.RegistriesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armmachinelearning.RegistriesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method RegistriesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, registryName string, options *armmachinelearning.RegistriesClientBeginDeleteOptions) (resp azfake.PollerResponder[armmachinelearning.RegistriesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method RegistriesClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, registryName string, options *armmachinelearning.RegistriesClientGetOptions) (resp azfake.Responder[armmachinelearning.RegistriesClientGetResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method RegistriesClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(resourceGroupName string, options *armmachinelearning.RegistriesClientListOptions) (resp azfake.PagerResponder[armmachinelearning.RegistriesClientListResponse])

	// NewListBySubscriptionPager is the fake for method RegistriesClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armmachinelearning.RegistriesClientListBySubscriptionOptions) (resp azfake.PagerResponder[armmachinelearning.RegistriesClientListBySubscriptionResponse])

	// BeginRemoveRegions is the fake for method RegistriesClient.BeginRemoveRegions
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginRemoveRegions func(ctx context.Context, resourceGroupName string, registryName string, body armmachinelearning.Registry, options *armmachinelearning.RegistriesClientBeginRemoveRegionsOptions) (resp azfake.PollerResponder[armmachinelearning.RegistriesClientRemoveRegionsResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method RegistriesClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, registryName string, body armmachinelearning.PartialRegistryPartialTrackedResource, options *armmachinelearning.RegistriesClientUpdateOptions) (resp azfake.Responder[armmachinelearning.RegistriesClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewRegistriesServerTransport creates a new instance of RegistriesServerTransport with the provided implementation.
// The returned RegistriesServerTransport instance is connected to an instance of armmachinelearning.RegistriesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewRegistriesServerTransport(srv *RegistriesServer) *RegistriesServerTransport {
	return &RegistriesServerTransport{
		srv:                        srv,
		beginCreateOrUpdate:        newTracker[azfake.PollerResponder[armmachinelearning.RegistriesClientCreateOrUpdateResponse]](),
		beginDelete:                newTracker[azfake.PollerResponder[armmachinelearning.RegistriesClientDeleteResponse]](),
		newListPager:               newTracker[azfake.PagerResponder[armmachinelearning.RegistriesClientListResponse]](),
		newListBySubscriptionPager: newTracker[azfake.PagerResponder[armmachinelearning.RegistriesClientListBySubscriptionResponse]](),
		beginRemoveRegions:         newTracker[azfake.PollerResponder[armmachinelearning.RegistriesClientRemoveRegionsResponse]](),
	}
}

// RegistriesServerTransport connects instances of armmachinelearning.RegistriesClient to instances of RegistriesServer.
// Don't use this type directly, use NewRegistriesServerTransport instead.
type RegistriesServerTransport struct {
	srv                        *RegistriesServer
	beginCreateOrUpdate        *tracker[azfake.PollerResponder[armmachinelearning.RegistriesClientCreateOrUpdateResponse]]
	beginDelete                *tracker[azfake.PollerResponder[armmachinelearning.RegistriesClientDeleteResponse]]
	newListPager               *tracker[azfake.PagerResponder[armmachinelearning.RegistriesClientListResponse]]
	newListBySubscriptionPager *tracker[azfake.PagerResponder[armmachinelearning.RegistriesClientListBySubscriptionResponse]]
	beginRemoveRegions         *tracker[azfake.PollerResponder[armmachinelearning.RegistriesClientRemoveRegionsResponse]]
}

// Do implements the policy.Transporter interface for RegistriesServerTransport.
func (r *RegistriesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "RegistriesClient.BeginCreateOrUpdate":
		resp, err = r.dispatchBeginCreateOrUpdate(req)
	case "RegistriesClient.BeginDelete":
		resp, err = r.dispatchBeginDelete(req)
	case "RegistriesClient.Get":
		resp, err = r.dispatchGet(req)
	case "RegistriesClient.NewListPager":
		resp, err = r.dispatchNewListPager(req)
	case "RegistriesClient.NewListBySubscriptionPager":
		resp, err = r.dispatchNewListBySubscriptionPager(req)
	case "RegistriesClient.BeginRemoveRegions":
		resp, err = r.dispatchBeginRemoveRegions(req)
	case "RegistriesClient.Update":
		resp, err = r.dispatchUpdate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (r *RegistriesServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if r.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := r.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmachinelearning.Registry](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := r.srv.BeginCreateOrUpdate(req.Context(), resourceGroupNameParam, registryNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		r.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		r.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		r.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (r *RegistriesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if r.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := r.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := r.srv.BeginDelete(req.Context(), resourceGroupNameParam, registryNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		r.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		r.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		r.beginDelete.remove(req)
	}

	return resp, nil
}

func (r *RegistriesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if r.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	registryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := r.srv.Get(req.Context(), resourceGroupNameParam, registryNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Registry, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (r *RegistriesServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := r.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/registries`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		resp := r.srv.NewListPager(resourceGroupNameParam, nil)
		newListPager = &resp
		r.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armmachinelearning.RegistriesClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		r.newListPager.remove(req)
	}
	return resp, nil
}

func (r *RegistriesServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if r.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := r.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/registries`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := r.srv.NewListBySubscriptionPager(nil)
		newListBySubscriptionPager = &resp
		r.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armmachinelearning.RegistriesClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		r.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		r.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (r *RegistriesServerTransport) dispatchBeginRemoveRegions(req *http.Request) (*http.Response, error) {
	if r.srv.BeginRemoveRegions == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginRemoveRegions not implemented")}
	}
	beginRemoveRegions := r.beginRemoveRegions.get(req)
	if beginRemoveRegions == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/removeRegions`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armmachinelearning.Registry](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		registryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := r.srv.BeginRemoveRegions(req.Context(), resourceGroupNameParam, registryNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginRemoveRegions = &respr
		r.beginRemoveRegions.add(req, beginRemoveRegions)
	}

	resp, err := server.PollerResponderNext(beginRemoveRegions, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		r.beginRemoveRegions.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginRemoveRegions) {
		r.beginRemoveRegions.remove(req)
	}

	return resp, nil
}

func (r *RegistriesServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if r.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.MachineLearningServices/registries/(?P<registryName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armmachinelearning.PartialRegistryPartialTrackedResource](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	registryNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("registryName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := r.srv.Update(req.Context(), resourceGroupNameParam, registryNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Registry, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
