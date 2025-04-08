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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/eventgrid/armeventgrid/v2"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// PartnerTopicsServer is a fake server for instances of the armeventgrid.PartnerTopicsClient type.
type PartnerTopicsServer struct {
	// Activate is the fake for method PartnerTopicsClient.Activate
	// HTTP status codes to indicate success: http.StatusOK
	Activate func(ctx context.Context, resourceGroupName string, partnerTopicName string, options *armeventgrid.PartnerTopicsClientActivateOptions) (resp azfake.Responder[armeventgrid.PartnerTopicsClientActivateResponse], errResp azfake.ErrorResponder)

	// CreateOrUpdate is the fake for method PartnerTopicsClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, partnerTopicName string, partnerTopicInfo armeventgrid.PartnerTopic, options *armeventgrid.PartnerTopicsClientCreateOrUpdateOptions) (resp azfake.Responder[armeventgrid.PartnerTopicsClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Deactivate is the fake for method PartnerTopicsClient.Deactivate
	// HTTP status codes to indicate success: http.StatusOK
	Deactivate func(ctx context.Context, resourceGroupName string, partnerTopicName string, options *armeventgrid.PartnerTopicsClientDeactivateOptions) (resp azfake.Responder[armeventgrid.PartnerTopicsClientDeactivateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method PartnerTopicsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, partnerTopicName string, options *armeventgrid.PartnerTopicsClientBeginDeleteOptions) (resp azfake.PollerResponder[armeventgrid.PartnerTopicsClientDeleteResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method PartnerTopicsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, partnerTopicName string, options *armeventgrid.PartnerTopicsClientGetOptions) (resp azfake.Responder[armeventgrid.PartnerTopicsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method PartnerTopicsClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armeventgrid.PartnerTopicsClientListByResourceGroupOptions) (resp azfake.PagerResponder[armeventgrid.PartnerTopicsClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method PartnerTopicsClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(options *armeventgrid.PartnerTopicsClientListBySubscriptionOptions) (resp azfake.PagerResponder[armeventgrid.PartnerTopicsClientListBySubscriptionResponse])

	// Update is the fake for method PartnerTopicsClient.Update
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	Update func(ctx context.Context, resourceGroupName string, partnerTopicName string, partnerTopicUpdateParameters armeventgrid.PartnerTopicUpdateParameters, options *armeventgrid.PartnerTopicsClientUpdateOptions) (resp azfake.Responder[armeventgrid.PartnerTopicsClientUpdateResponse], errResp azfake.ErrorResponder)
}

// NewPartnerTopicsServerTransport creates a new instance of PartnerTopicsServerTransport with the provided implementation.
// The returned PartnerTopicsServerTransport instance is connected to an instance of armeventgrid.PartnerTopicsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewPartnerTopicsServerTransport(srv *PartnerTopicsServer) *PartnerTopicsServerTransport {
	return &PartnerTopicsServerTransport{
		srv:                         srv,
		beginDelete:                 newTracker[azfake.PollerResponder[armeventgrid.PartnerTopicsClientDeleteResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armeventgrid.PartnerTopicsClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armeventgrid.PartnerTopicsClientListBySubscriptionResponse]](),
	}
}

// PartnerTopicsServerTransport connects instances of armeventgrid.PartnerTopicsClient to instances of PartnerTopicsServer.
// Don't use this type directly, use NewPartnerTopicsServerTransport instead.
type PartnerTopicsServerTransport struct {
	srv                         *PartnerTopicsServer
	beginDelete                 *tracker[azfake.PollerResponder[armeventgrid.PartnerTopicsClientDeleteResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armeventgrid.PartnerTopicsClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armeventgrid.PartnerTopicsClientListBySubscriptionResponse]]
}

// Do implements the policy.Transporter interface for PartnerTopicsServerTransport.
func (p *PartnerTopicsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return p.dispatchToMethodFake(req, method)
}

func (p *PartnerTopicsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if partnerTopicsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = partnerTopicsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "PartnerTopicsClient.Activate":
				res.resp, res.err = p.dispatchActivate(req)
			case "PartnerTopicsClient.CreateOrUpdate":
				res.resp, res.err = p.dispatchCreateOrUpdate(req)
			case "PartnerTopicsClient.Deactivate":
				res.resp, res.err = p.dispatchDeactivate(req)
			case "PartnerTopicsClient.BeginDelete":
				res.resp, res.err = p.dispatchBeginDelete(req)
			case "PartnerTopicsClient.Get":
				res.resp, res.err = p.dispatchGet(req)
			case "PartnerTopicsClient.NewListByResourceGroupPager":
				res.resp, res.err = p.dispatchNewListByResourceGroupPager(req)
			case "PartnerTopicsClient.NewListBySubscriptionPager":
				res.resp, res.err = p.dispatchNewListBySubscriptionPager(req)
			case "PartnerTopicsClient.Update":
				res.resp, res.err = p.dispatchUpdate(req)
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

func (p *PartnerTopicsServerTransport) dispatchActivate(req *http.Request) (*http.Response, error) {
	if p.srv.Activate == nil {
		return nil, &nonRetriableError{errors.New("fake for method Activate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.EventGrid/partnerTopics/(?P<partnerTopicName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/activate`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	partnerTopicNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("partnerTopicName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.Activate(req.Context(), resourceGroupNameParam, partnerTopicNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).PartnerTopic, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *PartnerTopicsServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if p.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.EventGrid/partnerTopics/(?P<partnerTopicName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armeventgrid.PartnerTopic](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	partnerTopicNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("partnerTopicName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, partnerTopicNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).PartnerTopic, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *PartnerTopicsServerTransport) dispatchDeactivate(req *http.Request) (*http.Response, error) {
	if p.srv.Deactivate == nil {
		return nil, &nonRetriableError{errors.New("fake for method Deactivate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.EventGrid/partnerTopics/(?P<partnerTopicName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/deactivate`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	partnerTopicNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("partnerTopicName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.Deactivate(req.Context(), resourceGroupNameParam, partnerTopicNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).PartnerTopic, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *PartnerTopicsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if p.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := p.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.EventGrid/partnerTopics/(?P<partnerTopicName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		partnerTopicNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("partnerTopicName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := p.srv.BeginDelete(req.Context(), resourceGroupNameParam, partnerTopicNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		p.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		p.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		p.beginDelete.remove(req)
	}

	return resp, nil
}

func (p *PartnerTopicsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if p.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.EventGrid/partnerTopics/(?P<partnerTopicName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	partnerTopicNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("partnerTopicName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.Get(req.Context(), resourceGroupNameParam, partnerTopicNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).PartnerTopic, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (p *PartnerTopicsServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if p.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := p.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.EventGrid/partnerTopics`
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
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
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
		var options *armeventgrid.PartnerTopicsClientListByResourceGroupOptions
		if filterParam != nil || topParam != nil {
			options = &armeventgrid.PartnerTopicsClientListByResourceGroupOptions{
				Filter: filterParam,
				Top:    topParam,
			}
		}
		resp := p.srv.NewListByResourceGroupPager(resourceGroupNameParam, options)
		newListByResourceGroupPager = &resp
		p.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armeventgrid.PartnerTopicsClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		p.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		p.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (p *PartnerTopicsServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if p.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := p.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.EventGrid/partnerTopics`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		qp := req.URL.Query()
		filterUnescaped, err := url.QueryUnescape(qp.Get("$filter"))
		if err != nil {
			return nil, err
		}
		filterParam := getOptional(filterUnescaped)
		topUnescaped, err := url.QueryUnescape(qp.Get("$top"))
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
		var options *armeventgrid.PartnerTopicsClientListBySubscriptionOptions
		if filterParam != nil || topParam != nil {
			options = &armeventgrid.PartnerTopicsClientListBySubscriptionOptions{
				Filter: filterParam,
				Top:    topParam,
			}
		}
		resp := p.srv.NewListBySubscriptionPager(options)
		newListBySubscriptionPager = &resp
		p.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armeventgrid.PartnerTopicsClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		p.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		p.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (p *PartnerTopicsServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if p.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.EventGrid/partnerTopics/(?P<partnerTopicName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armeventgrid.PartnerTopicUpdateParameters](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	partnerTopicNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("partnerTopicName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := p.srv.Update(req.Context(), resourceGroupNameParam, partnerTopicNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK, http.StatusCreated}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).PartnerTopic, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// set this to conditionally intercept incoming requests to PartnerTopicsServerTransport
var partnerTopicsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
