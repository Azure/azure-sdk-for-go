// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) Go Code Generator. DO NOT EDIT.

package fake

import (
	"context"
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/impactreporting/armimpactreporting"
	"net/http"
	"net/url"
	"regexp"
)

// InsightsServer is a fake server for instances of the armimpactreporting.InsightsClient type.
type InsightsServer struct {
	// Create is the fake for method InsightsClient.Create
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	Create func(ctx context.Context, workloadImpactName string, insightName string, resource armimpactreporting.Insight, options *armimpactreporting.InsightsClientCreateOptions) (resp azfake.Responder[armimpactreporting.InsightsClientCreateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method InsightsClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, workloadImpactName string, insightName string, options *armimpactreporting.InsightsClientDeleteOptions) (resp azfake.Responder[armimpactreporting.InsightsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method InsightsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, workloadImpactName string, insightName string, options *armimpactreporting.InsightsClientGetOptions) (resp azfake.Responder[armimpactreporting.InsightsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListBySubscriptionPager is the fake for method InsightsClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(workloadImpactName string, options *armimpactreporting.InsightsClientListBySubscriptionOptions) (resp azfake.PagerResponder[armimpactreporting.InsightsClientListBySubscriptionResponse])
}

// NewInsightsServerTransport creates a new instance of InsightsServerTransport with the provided implementation.
// The returned InsightsServerTransport instance is connected to an instance of armimpactreporting.InsightsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewInsightsServerTransport(srv *InsightsServer) *InsightsServerTransport {
	return &InsightsServerTransport{
		srv:                        srv,
		newListBySubscriptionPager: newTracker[azfake.PagerResponder[armimpactreporting.InsightsClientListBySubscriptionResponse]](),
	}
}

// InsightsServerTransport connects instances of armimpactreporting.InsightsClient to instances of InsightsServer.
// Don't use this type directly, use NewInsightsServerTransport instead.
type InsightsServerTransport struct {
	srv                        *InsightsServer
	newListBySubscriptionPager *tracker[azfake.PagerResponder[armimpactreporting.InsightsClientListBySubscriptionResponse]]
}

// Do implements the policy.Transporter interface for InsightsServerTransport.
func (i *InsightsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return i.dispatchToMethodFake(req, method)
}

func (i *InsightsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if insightsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = insightsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "InsightsClient.Create":
				res.resp, res.err = i.dispatchCreate(req)
			case "InsightsClient.Delete":
				res.resp, res.err = i.dispatchDelete(req)
			case "InsightsClient.Get":
				res.resp, res.err = i.dispatchGet(req)
			case "InsightsClient.NewListBySubscriptionPager":
				res.resp, res.err = i.dispatchNewListBySubscriptionPager(req)
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

func (i *InsightsServerTransport) dispatchCreate(req *http.Request) (*http.Response, error) {
	if i.srv.Create == nil {
		return nil, &nonRetriableError{errors.New("fake for method Create not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Impact/workloadImpacts/(?P<workloadImpactName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/insights/(?P<insightName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armimpactreporting.Insight](req)
	if err != nil {
		return nil, err
	}
	workloadImpactNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workloadImpactName")])
	if err != nil {
		return nil, err
	}
	insightNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("insightName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Create(req.Context(), workloadImpactNameParam, insightNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Insight, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *InsightsServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if i.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Impact/workloadImpacts/(?P<workloadImpactName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/insights/(?P<insightName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	workloadImpactNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workloadImpactName")])
	if err != nil {
		return nil, err
	}
	insightNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("insightName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Delete(req.Context(), workloadImpactNameParam, insightNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusNoContent}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusNoContent", respContent.HTTPStatus)}
	}
	resp, err := server.NewResponse(respContent, req, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *InsightsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if i.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Impact/workloadImpacts/(?P<workloadImpactName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/insights/(?P<insightName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	workloadImpactNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workloadImpactName")])
	if err != nil {
		return nil, err
	}
	insightNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("insightName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := i.srv.Get(req.Context(), workloadImpactNameParam, insightNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Insight, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (i *InsightsServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if i.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := i.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Impact/workloadImpacts/(?P<workloadImpactName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/insights`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		workloadImpactNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("workloadImpactName")])
		if err != nil {
			return nil, err
		}
		resp := i.srv.NewListBySubscriptionPager(workloadImpactNameParam, nil)
		newListBySubscriptionPager = &resp
		i.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armimpactreporting.InsightsClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		i.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		i.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to InsightsServerTransport
var insightsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
