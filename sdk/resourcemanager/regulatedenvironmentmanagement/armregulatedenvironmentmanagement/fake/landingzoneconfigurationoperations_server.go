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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/regulatedenvironmentmanagement/armregulatedenvironmentmanagement"
	"net/http"
	"net/url"
	"regexp"
)

// LandingZoneConfigurationOperationsServer is a fake server for instances of the armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClient type.
type LandingZoneConfigurationOperationsServer struct {
	// BeginCreate is the fake for method LandingZoneConfigurationOperationsClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusCreated
	BeginCreate func(ctx context.Context, resourceGroupName string, landingZoneAccountName string, landingZoneConfigurationName string, resource armregulatedenvironmentmanagement.LZConfiguration, options *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientBeginCreateOptions) (resp azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginCreateCopy is the fake for method LandingZoneConfigurationOperationsClient.BeginCreateCopy
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginCreateCopy func(ctx context.Context, resourceGroupName string, landingZoneAccountName string, landingZoneConfigurationName string, body armregulatedenvironmentmanagement.CreateLZConfigurationCopyRequest, options *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientBeginCreateCopyOptions) (resp azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientCreateCopyResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method LandingZoneConfigurationOperationsClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, landingZoneAccountName string, landingZoneConfigurationName string, options *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientBeginDeleteOptions) (resp azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientDeleteResponse], errResp azfake.ErrorResponder)

	// BeginGenerateLandingZone is the fake for method LandingZoneConfigurationOperationsClient.BeginGenerateLandingZone
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginGenerateLandingZone func(ctx context.Context, resourceGroupName string, landingZoneAccountName string, landingZoneConfigurationName string, body armregulatedenvironmentmanagement.GenerateLandingZoneRequest, options *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientBeginGenerateLandingZoneOptions) (resp azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientGenerateLandingZoneResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method LandingZoneConfigurationOperationsClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, landingZoneAccountName string, landingZoneConfigurationName string, options *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientGetOptions) (resp azfake.Responder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientGetResponse], errResp azfake.ErrorResponder)

	// NewListByResourceGroupPager is the fake for method LandingZoneConfigurationOperationsClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, landingZoneAccountName string, options *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientListByResourceGroupOptions) (resp azfake.PagerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientListByResourceGroupResponse])

	// NewListBySubscriptionPager is the fake for method LandingZoneConfigurationOperationsClient.NewListBySubscriptionPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListBySubscriptionPager func(landingZoneAccountName string, options *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientListBySubscriptionOptions) (resp azfake.PagerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientListBySubscriptionResponse])

	// BeginUpdate is the fake for method LandingZoneConfigurationOperationsClient.BeginUpdate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdate func(ctx context.Context, resourceGroupName string, landingZoneAccountName string, landingZoneConfigurationName string, properties armregulatedenvironmentmanagement.LZConfiguration, options *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientBeginUpdateOptions) (resp azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientUpdateResponse], errResp azfake.ErrorResponder)

	// BeginUpdateAuthoringStatus is the fake for method LandingZoneConfigurationOperationsClient.BeginUpdateAuthoringStatus
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdateAuthoringStatus func(ctx context.Context, resourceGroupName string, landingZoneAccountName string, landingZoneConfigurationName string, body armregulatedenvironmentmanagement.UpdateAuthoringStatusRequest, options *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientBeginUpdateAuthoringStatusOptions) (resp azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientUpdateAuthoringStatusResponse], errResp azfake.ErrorResponder)
}

// NewLandingZoneConfigurationOperationsServerTransport creates a new instance of LandingZoneConfigurationOperationsServerTransport with the provided implementation.
// The returned LandingZoneConfigurationOperationsServerTransport instance is connected to an instance of armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewLandingZoneConfigurationOperationsServerTransport(srv *LandingZoneConfigurationOperationsServer) *LandingZoneConfigurationOperationsServerTransport {
	return &LandingZoneConfigurationOperationsServerTransport{
		srv:                         srv,
		beginCreate:                 newTracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientCreateResponse]](),
		beginCreateCopy:             newTracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientCreateCopyResponse]](),
		beginDelete:                 newTracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientDeleteResponse]](),
		beginGenerateLandingZone:    newTracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientGenerateLandingZoneResponse]](),
		newListByResourceGroupPager: newTracker[azfake.PagerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientListByResourceGroupResponse]](),
		newListBySubscriptionPager:  newTracker[azfake.PagerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientListBySubscriptionResponse]](),
		beginUpdate:                 newTracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientUpdateResponse]](),
		beginUpdateAuthoringStatus:  newTracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientUpdateAuthoringStatusResponse]](),
	}
}

// LandingZoneConfigurationOperationsServerTransport connects instances of armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClient to instances of LandingZoneConfigurationOperationsServer.
// Don't use this type directly, use NewLandingZoneConfigurationOperationsServerTransport instead.
type LandingZoneConfigurationOperationsServerTransport struct {
	srv                         *LandingZoneConfigurationOperationsServer
	beginCreate                 *tracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientCreateResponse]]
	beginCreateCopy             *tracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientCreateCopyResponse]]
	beginDelete                 *tracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientDeleteResponse]]
	beginGenerateLandingZone    *tracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientGenerateLandingZoneResponse]]
	newListByResourceGroupPager *tracker[azfake.PagerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientListByResourceGroupResponse]]
	newListBySubscriptionPager  *tracker[azfake.PagerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientListBySubscriptionResponse]]
	beginUpdate                 *tracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientUpdateResponse]]
	beginUpdateAuthoringStatus  *tracker[azfake.PollerResponder[armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientUpdateAuthoringStatusResponse]]
}

// Do implements the policy.Transporter interface for LandingZoneConfigurationOperationsServerTransport.
func (l *LandingZoneConfigurationOperationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	return l.dispatchToMethodFake(req, method)
}

func (l *LandingZoneConfigurationOperationsServerTransport) dispatchToMethodFake(req *http.Request, method string) (*http.Response, error) {
	resultChan := make(chan result)
	defer close(resultChan)

	go func() {
		var intercepted bool
		var res result
		if landingZoneConfigurationOperationsServerTransportInterceptor != nil {
			res.resp, res.err, intercepted = landingZoneConfigurationOperationsServerTransportInterceptor.Do(req)
		}
		if !intercepted {
			switch method {
			case "LandingZoneConfigurationOperationsClient.BeginCreate":
				res.resp, res.err = l.dispatchBeginCreate(req)
			case "LandingZoneConfigurationOperationsClient.BeginCreateCopy":
				res.resp, res.err = l.dispatchBeginCreateCopy(req)
			case "LandingZoneConfigurationOperationsClient.BeginDelete":
				res.resp, res.err = l.dispatchBeginDelete(req)
			case "LandingZoneConfigurationOperationsClient.BeginGenerateLandingZone":
				res.resp, res.err = l.dispatchBeginGenerateLandingZone(req)
			case "LandingZoneConfigurationOperationsClient.Get":
				res.resp, res.err = l.dispatchGet(req)
			case "LandingZoneConfigurationOperationsClient.NewListByResourceGroupPager":
				res.resp, res.err = l.dispatchNewListByResourceGroupPager(req)
			case "LandingZoneConfigurationOperationsClient.NewListBySubscriptionPager":
				res.resp, res.err = l.dispatchNewListBySubscriptionPager(req)
			case "LandingZoneConfigurationOperationsClient.BeginUpdate":
				res.resp, res.err = l.dispatchBeginUpdate(req)
			case "LandingZoneConfigurationOperationsClient.BeginUpdateAuthoringStatus":
				res.resp, res.err = l.dispatchBeginUpdateAuthoringStatus(req)
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

func (l *LandingZoneConfigurationOperationsServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if l.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := l.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sovereign/landingZoneAccounts/(?P<landingZoneAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/landingZoneConfigurations/(?P<landingZoneConfigurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armregulatedenvironmentmanagement.LZConfiguration](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		landingZoneAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneAccountName")])
		if err != nil {
			return nil, err
		}
		landingZoneConfigurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneConfigurationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginCreate(req.Context(), resourceGroupNameParam, landingZoneAccountNameParam, landingZoneConfigurationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		l.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusCreated}, resp.StatusCode) {
		l.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusCreated", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		l.beginCreate.remove(req)
	}

	return resp, nil
}

func (l *LandingZoneConfigurationOperationsServerTransport) dispatchBeginCreateCopy(req *http.Request) (*http.Response, error) {
	if l.srv.BeginCreateCopy == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreateCopy not implemented")}
	}
	beginCreateCopy := l.beginCreateCopy.get(req)
	if beginCreateCopy == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sovereign/landingZoneAccounts/(?P<landingZoneAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/landingZoneConfigurations/(?P<landingZoneConfigurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/createCopy`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armregulatedenvironmentmanagement.CreateLZConfigurationCopyRequest](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		landingZoneAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneAccountName")])
		if err != nil {
			return nil, err
		}
		landingZoneConfigurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneConfigurationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginCreateCopy(req.Context(), resourceGroupNameParam, landingZoneAccountNameParam, landingZoneConfigurationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreateCopy = &respr
		l.beginCreateCopy.add(req, beginCreateCopy)
	}

	resp, err := server.PollerResponderNext(beginCreateCopy, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		l.beginCreateCopy.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreateCopy) {
		l.beginCreateCopy.remove(req)
	}

	return resp, nil
}

func (l *LandingZoneConfigurationOperationsServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if l.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := l.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sovereign/landingZoneAccounts/(?P<landingZoneAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/landingZoneConfigurations/(?P<landingZoneConfigurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		landingZoneAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneAccountName")])
		if err != nil {
			return nil, err
		}
		landingZoneConfigurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneConfigurationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginDelete(req.Context(), resourceGroupNameParam, landingZoneAccountNameParam, landingZoneConfigurationNameParam, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginDelete = &respr
		l.beginDelete.add(req, beginDelete)
	}

	resp, err := server.PollerResponderNext(beginDelete, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted, http.StatusNoContent}, resp.StatusCode) {
		l.beginDelete.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted, http.StatusNoContent", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginDelete) {
		l.beginDelete.remove(req)
	}

	return resp, nil
}

func (l *LandingZoneConfigurationOperationsServerTransport) dispatchBeginGenerateLandingZone(req *http.Request) (*http.Response, error) {
	if l.srv.BeginGenerateLandingZone == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginGenerateLandingZone not implemented")}
	}
	beginGenerateLandingZone := l.beginGenerateLandingZone.get(req)
	if beginGenerateLandingZone == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sovereign/landingZoneAccounts/(?P<landingZoneAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/landingZoneConfigurations/(?P<landingZoneConfigurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/generateLandingZone`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armregulatedenvironmentmanagement.GenerateLandingZoneRequest](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		landingZoneAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneAccountName")])
		if err != nil {
			return nil, err
		}
		landingZoneConfigurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneConfigurationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginGenerateLandingZone(req.Context(), resourceGroupNameParam, landingZoneAccountNameParam, landingZoneConfigurationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginGenerateLandingZone = &respr
		l.beginGenerateLandingZone.add(req, beginGenerateLandingZone)
	}

	resp, err := server.PollerResponderNext(beginGenerateLandingZone, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		l.beginGenerateLandingZone.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginGenerateLandingZone) {
		l.beginGenerateLandingZone.remove(req)
	}

	return resp, nil
}

func (l *LandingZoneConfigurationOperationsServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if l.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sovereign/landingZoneAccounts/(?P<landingZoneAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/landingZoneConfigurations/(?P<landingZoneConfigurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	landingZoneAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneAccountName")])
	if err != nil {
		return nil, err
	}
	landingZoneConfigurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneConfigurationName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := l.srv.Get(req.Context(), resourceGroupNameParam, landingZoneAccountNameParam, landingZoneConfigurationNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).LZConfiguration, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *LandingZoneConfigurationOperationsServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if l.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := l.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sovereign/landingZoneAccounts/(?P<landingZoneAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/landingZoneConfigurations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		landingZoneAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneAccountName")])
		if err != nil {
			return nil, err
		}
		resp := l.srv.NewListByResourceGroupPager(resourceGroupNameParam, landingZoneAccountNameParam, nil)
		newListByResourceGroupPager = &resp
		l.newListByResourceGroupPager.add(req, newListByResourceGroupPager)
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientListByResourceGroupResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListByResourceGroupPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		l.newListByResourceGroupPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListByResourceGroupPager) {
		l.newListByResourceGroupPager.remove(req)
	}
	return resp, nil
}

func (l *LandingZoneConfigurationOperationsServerTransport) dispatchNewListBySubscriptionPager(req *http.Request) (*http.Response, error) {
	if l.srv.NewListBySubscriptionPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListBySubscriptionPager not implemented")}
	}
	newListBySubscriptionPager := l.newListBySubscriptionPager.get(req)
	if newListBySubscriptionPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sovereign/landingZoneAccounts/(?P<landingZoneAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/landingZoneConfigurations`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 2 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		landingZoneAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneAccountName")])
		if err != nil {
			return nil, err
		}
		resp := l.srv.NewListBySubscriptionPager(landingZoneAccountNameParam, nil)
		newListBySubscriptionPager = &resp
		l.newListBySubscriptionPager.add(req, newListBySubscriptionPager)
		server.PagerResponderInjectNextLinks(newListBySubscriptionPager, req, func(page *armregulatedenvironmentmanagement.LandingZoneConfigurationOperationsClientListBySubscriptionResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListBySubscriptionPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		l.newListBySubscriptionPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListBySubscriptionPager) {
		l.newListBySubscriptionPager.remove(req)
	}
	return resp, nil
}

func (l *LandingZoneConfigurationOperationsServerTransport) dispatchBeginUpdate(req *http.Request) (*http.Response, error) {
	if l.srv.BeginUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdate not implemented")}
	}
	beginUpdate := l.beginUpdate.get(req)
	if beginUpdate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sovereign/landingZoneAccounts/(?P<landingZoneAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/landingZoneConfigurations/(?P<landingZoneConfigurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armregulatedenvironmentmanagement.LZConfiguration](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		landingZoneAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneAccountName")])
		if err != nil {
			return nil, err
		}
		landingZoneConfigurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneConfigurationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginUpdate(req.Context(), resourceGroupNameParam, landingZoneAccountNameParam, landingZoneConfigurationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdate = &respr
		l.beginUpdate.add(req, beginUpdate)
	}

	resp, err := server.PollerResponderNext(beginUpdate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		l.beginUpdate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdate) {
		l.beginUpdate.remove(req)
	}

	return resp, nil
}

func (l *LandingZoneConfigurationOperationsServerTransport) dispatchBeginUpdateAuthoringStatus(req *http.Request) (*http.Response, error) {
	if l.srv.BeginUpdateAuthoringStatus == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdateAuthoringStatus not implemented")}
	}
	beginUpdateAuthoringStatus := l.beginUpdateAuthoringStatus.get(req)
	if beginUpdateAuthoringStatus == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Sovereign/landingZoneAccounts/(?P<landingZoneAccountName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/landingZoneConfigurations/(?P<landingZoneConfigurationName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/updateAuthoringStatus`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armregulatedenvironmentmanagement.UpdateAuthoringStatusRequest](req)
		if err != nil {
			return nil, err
		}
		resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
		if err != nil {
			return nil, err
		}
		landingZoneAccountNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneAccountName")])
		if err != nil {
			return nil, err
		}
		landingZoneConfigurationNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("landingZoneConfigurationName")])
		if err != nil {
			return nil, err
		}
		respr, errRespr := l.srv.BeginUpdateAuthoringStatus(req.Context(), resourceGroupNameParam, landingZoneAccountNameParam, landingZoneConfigurationNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdateAuthoringStatus = &respr
		l.beginUpdateAuthoringStatus.add(req, beginUpdateAuthoringStatus)
	}

	resp, err := server.PollerResponderNext(beginUpdateAuthoringStatus, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		l.beginUpdateAuthoringStatus.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdateAuthoringStatus) {
		l.beginUpdateAuthoringStatus.remove(req)
	}

	return resp, nil
}

// set this to conditionally intercept incoming requests to LandingZoneConfigurationOperationsServerTransport
var landingZoneConfigurationOperationsServerTransportInterceptor interface {
	// Do returns true if the server transport should use the returned response/error
	Do(*http.Request) (*http.Response, error, bool)
}
