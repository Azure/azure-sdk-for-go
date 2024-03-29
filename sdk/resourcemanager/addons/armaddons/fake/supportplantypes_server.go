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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/addons/armaddons"
	"net/http"
	"net/url"
	"regexp"
)

// SupportPlanTypesServer is a fake server for instances of the armaddons.SupportPlanTypesClient type.
type SupportPlanTypesServer struct {
	// BeginCreateOrUpdate is the fake for method SupportPlanTypesClient.BeginCreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated, http.StatusNotFound
	BeginCreateOrUpdate func(ctx context.Context, providerName string, planTypeName armaddons.PlanTypeName, options *armaddons.SupportPlanTypesClientBeginCreateOrUpdateOptions) (resp azfake.PollerResponder[armaddons.SupportPlanTypesClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method SupportPlanTypesClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, providerName string, planTypeName armaddons.PlanTypeName, options *armaddons.SupportPlanTypesClientBeginDeleteOptions) (resp azfake.PollerResponder[armaddons.SupportPlanTypesClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method SupportPlanTypesClient.Get
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNotFound
	Get func(ctx context.Context, providerName string, planTypeName armaddons.PlanTypeName, options *armaddons.SupportPlanTypesClientGetOptions) (resp azfake.Responder[armaddons.SupportPlanTypesClientGetResponse], errResp azfake.ErrorResponder)

	// ListInfo is the fake for method SupportPlanTypesClient.ListInfo
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNotFound
	ListInfo func(ctx context.Context, options *armaddons.SupportPlanTypesClientListInfoOptions) (resp azfake.Responder[armaddons.SupportPlanTypesClientListInfoResponse], errResp azfake.ErrorResponder)
}

// NewSupportPlanTypesServerTransport creates a new instance of SupportPlanTypesServerTransport with the provided implementation.
// The returned SupportPlanTypesServerTransport instance is connected to an instance of armaddons.SupportPlanTypesClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewSupportPlanTypesServerTransport(srv *SupportPlanTypesServer) *SupportPlanTypesServerTransport {
	return &SupportPlanTypesServerTransport{
		srv:                 srv,
		beginCreateOrUpdate: newTracker[azfake.PollerResponder[armaddons.SupportPlanTypesClientCreateOrUpdateResponse]](),
		beginDelete:         newTracker[azfake.PollerResponder[armaddons.SupportPlanTypesClientDeleteResponse]](),
	}
}

// SupportPlanTypesServerTransport connects instances of armaddons.SupportPlanTypesClient to instances of SupportPlanTypesServer.
// Don't use this type directly, use NewSupportPlanTypesServerTransport instead.
type SupportPlanTypesServerTransport struct {
	srv                 *SupportPlanTypesServer
	beginCreateOrUpdate *tracker[azfake.PollerResponder[armaddons.SupportPlanTypesClientCreateOrUpdateResponse]]
	beginDelete         *tracker[azfake.PollerResponder[armaddons.SupportPlanTypesClientDeleteResponse]]
}

// Do implements the policy.Transporter interface for SupportPlanTypesServerTransport.
func (s *SupportPlanTypesServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "SupportPlanTypesClient.BeginCreateOrUpdate":
		resp, err = s.dispatchBeginCreateOrUpdate(req)
	case "SupportPlanTypesClient.BeginDelete":
		resp, err = s.dispatchBeginDelete(req)
	case "SupportPlanTypesClient.Get":
		resp, err = s.dispatchGet(req)
	case "SupportPlanTypesClient.ListInfo":
		resp, err = s.dispatchListInfo(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (s *SupportPlanTypesServerTransport) dispatchBeginCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if s.srv.BeginCreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateOrUpdate not implemented")}
	}
	beginCreateOrUpdate := s.beginCreateOrUpdate.get(req)
	if beginCreateOrUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Addons/supportProviders/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/supportPlanTypes/(?P<planTypeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
		if err != nil {
			return nil, err
		}
		planTypeNameParam, err := parseWithCast(matches[regex.SubexpIndex("planTypeName")], func(v string) (armaddons.PlanTypeName, error) {
			p, unescapeErr := url.PathUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armaddons.PlanTypeName(p), nil
		})
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginCreateOrUpdate(req.Context(), providerNameParam, planTypeNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateOrUpdate = &respr
		s.beginCreateOrUpdate.add(req, beginCreateOrUpdate)
	}

	resp, err := server.PollerResponderNext(beginCreateOrUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated, http.StatusNotFound}, resp.StatusCode) {
		s.beginCreateOrUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated, http.StatusNotFound", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateOrUpdate) {
		s.beginCreateOrUpdate.remove(req)
	}

	return resp, nil
}

func (s *SupportPlanTypesServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if s.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := s.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Addons/supportProviders/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/supportPlanTypes/(?P<planTypeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
		if err != nil {
			return nil, err
		}
		planTypeNameParam, err := parseWithCast(matches[regex.SubexpIndex("planTypeName")], func(v string) (armaddons.PlanTypeName, error) {
			p, unescapeErr := url.PathUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armaddons.PlanTypeName(p), nil
		})
		if err != nil {
			return nil, err
		}
		respr, errRespr := s.srv.BeginDelete(req.Context(), providerNameParam, planTypeNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		s.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		s.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		s.beginDelete.remove(req)
	}

	return resp, nil
}

func (s *SupportPlanTypesServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if s.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Addons/supportProviders/(?P<providerName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/supportPlanTypes/(?P<planTypeName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	providerNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("providerName")])
	if err != nil {
		return nil, err
	}
	planTypeNameParam, err := parseWithCast(matches[regex.SubexpIndex("planTypeName")], func(v string) (armaddons.PlanTypeName, error) {
		p, unescapeErr := url.PathUnescape(v)
		if unescapeErr != nil {
			return "", unescapeErr
		}
		return armaddons.PlanTypeName(p), nil
	})
	if err != nil {
		return nil, err
	}
	respr, errRespr := s.srv.Get(req.Context(), providerNameParam, planTypeNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusNotFound}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNotFound", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CanonicalSupportPlanResponseEnvelope, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *SupportPlanTypesServerTransport) dispatchListInfo(req *http.Request) (*http.Response, error) {
	if s.srv.ListInfo == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListInfo not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Addons/supportProviders/canonical/listSupportPlanInfo`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 1 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	respr, errRespr := s.srv.ListInfo(req.Context(), nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusNotFound}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNotFound", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CanonicalSupportPlanInfoDefinitionArray, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
