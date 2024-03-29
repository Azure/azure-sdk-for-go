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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hdinsight/armhdinsight"
	"net/http"
	"net/url"
	"regexp"
)

// ClustersServer is a fake server for instances of the armhdinsight.ClustersClient type.
type ClustersServer struct {
	// BeginCreate is the fake for method ClustersClient.BeginCreate
	// HTTP status codes to indicate success: http.StatusOK
	BeginCreate func(ctx context.Context, resourceGroupName string, clusterName string, parameters armhdinsight.ClusterCreateParametersExtended, options *armhdinsight.ClustersClientBeginCreateOptions) (resp azfake.PollerResponder[armhdinsight.ClustersClientCreateResponse], errResp azfake.ErrorResponder)

	// BeginDelete is the fake for method ClustersClient.BeginDelete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted, http.StatusNoContent
	BeginDelete func(ctx context.Context, resourceGroupName string, clusterName string, options *armhdinsight.ClustersClientBeginDeleteOptions) (resp azfake.PollerResponder[armhdinsight.ClustersClientDeleteResponse], errResp azfake.ErrorResponder)

	// BeginExecuteScriptActions is the fake for method ClustersClient.BeginExecuteScriptActions
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginExecuteScriptActions func(ctx context.Context, resourceGroupName string, clusterName string, parameters armhdinsight.ExecuteScriptActionParameters, options *armhdinsight.ClustersClientBeginExecuteScriptActionsOptions) (resp azfake.PollerResponder[armhdinsight.ClustersClientExecuteScriptActionsResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method ClustersClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, resourceGroupName string, clusterName string, options *armhdinsight.ClustersClientGetOptions) (resp azfake.Responder[armhdinsight.ClustersClientGetResponse], errResp azfake.ErrorResponder)

	// GetAzureAsyncOperationStatus is the fake for method ClustersClient.GetAzureAsyncOperationStatus
	// HTTP status codes to indicate success: http.StatusOK
	GetAzureAsyncOperationStatus func(ctx context.Context, resourceGroupName string, clusterName string, operationID string, options *armhdinsight.ClustersClientGetAzureAsyncOperationStatusOptions) (resp azfake.Responder[armhdinsight.ClustersClientGetAzureAsyncOperationStatusResponse], errResp azfake.ErrorResponder)

	// GetGatewaySettings is the fake for method ClustersClient.GetGatewaySettings
	// HTTP status codes to indicate success: http.StatusOK
	GetGatewaySettings func(ctx context.Context, resourceGroupName string, clusterName string, options *armhdinsight.ClustersClientGetGatewaySettingsOptions) (resp azfake.Responder[armhdinsight.ClustersClientGetGatewaySettingsResponse], errResp azfake.ErrorResponder)

	// NewListPager is the fake for method ClustersClient.NewListPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListPager func(options *armhdinsight.ClustersClientListOptions) (resp azfake.PagerResponder[armhdinsight.ClustersClientListResponse])

	// NewListByResourceGroupPager is the fake for method ClustersClient.NewListByResourceGroupPager
	// HTTP status codes to indicate success: http.StatusOK
	NewListByResourceGroupPager func(resourceGroupName string, options *armhdinsight.ClustersClientListByResourceGroupOptions) (resp azfake.PagerResponder[armhdinsight.ClustersClientListByResourceGroupResponse])

	// BeginResize is the fake for method ClustersClient.BeginResize
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginResize func(ctx context.Context, resourceGroupName string, clusterName string, roleName armhdinsight.RoleName, parameters armhdinsight.ClusterResizeParameters, options *armhdinsight.ClustersClientBeginResizeOptions) (resp azfake.PollerResponder[armhdinsight.ClustersClientResizeResponse], errResp azfake.ErrorResponder)

	// BeginRotateDiskEncryptionKey is the fake for method ClustersClient.BeginRotateDiskEncryptionKey
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginRotateDiskEncryptionKey func(ctx context.Context, resourceGroupName string, clusterName string, parameters armhdinsight.ClusterDiskEncryptionParameters, options *armhdinsight.ClustersClientBeginRotateDiskEncryptionKeyOptions) (resp azfake.PollerResponder[armhdinsight.ClustersClientRotateDiskEncryptionKeyResponse], errResp azfake.ErrorResponder)

	// Update is the fake for method ClustersClient.Update
	// HTTP status codes to indicate success: http.StatusOK
	Update func(ctx context.Context, resourceGroupName string, clusterName string, parameters armhdinsight.ClusterPatchParameters, options *armhdinsight.ClustersClientUpdateOptions) (resp azfake.Responder[armhdinsight.ClustersClientUpdateResponse], errResp azfake.ErrorResponder)

	// BeginUpdateAutoScaleConfiguration is the fake for method ClustersClient.BeginUpdateAutoScaleConfiguration
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdateAutoScaleConfiguration func(ctx context.Context, resourceGroupName string, clusterName string, roleName armhdinsight.RoleName, parameters armhdinsight.AutoscaleConfigurationUpdateParameter, options *armhdinsight.ClustersClientBeginUpdateAutoScaleConfigurationOptions) (resp azfake.PollerResponder[armhdinsight.ClustersClientUpdateAutoScaleConfigurationResponse], errResp azfake.ErrorResponder)

	// BeginUpdateGatewaySettings is the fake for method ClustersClient.BeginUpdateGatewaySettings
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdateGatewaySettings func(ctx context.Context, resourceGroupName string, clusterName string, parameters armhdinsight.UpdateGatewaySettingsParameters, options *armhdinsight.ClustersClientBeginUpdateGatewaySettingsOptions) (resp azfake.PollerResponder[armhdinsight.ClustersClientUpdateGatewaySettingsResponse], errResp azfake.ErrorResponder)

	// BeginUpdateIdentityCertificate is the fake for method ClustersClient.BeginUpdateIdentityCertificate
	// HTTP status codes to indicate success: http.StatusOK, http.StatusAccepted
	BeginUpdateIdentityCertificate func(ctx context.Context, resourceGroupName string, clusterName string, parameters armhdinsight.UpdateClusterIdentityCertificateParameters, options *armhdinsight.ClustersClientBeginUpdateIdentityCertificateOptions) (resp azfake.PollerResponder[armhdinsight.ClustersClientUpdateIdentityCertificateResponse], errResp azfake.ErrorResponder)
}

// NewClustersServerTransport creates a new instance of ClustersServerTransport with the provided implementation.
// The returned ClustersServerTransport instance is connected to an instance of armhdinsight.ClustersClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewClustersServerTransport(srv *ClustersServer) *ClustersServerTransport {
	return &ClustersServerTransport{
		srv:                               srv,
		beginCreate:                       newTracker[azfake.PollerResponder[armhdinsight.ClustersClientCreateResponse]](),
		beginDelete:                       newTracker[azfake.PollerResponder[armhdinsight.ClustersClientDeleteResponse]](),
		beginExecuteScriptActions:         newTracker[azfake.PollerResponder[armhdinsight.ClustersClientExecuteScriptActionsResponse]](),
		newListPager:                      newTracker[azfake.PagerResponder[armhdinsight.ClustersClientListResponse]](),
		newListByResourceGroupPager:       newTracker[azfake.PagerResponder[armhdinsight.ClustersClientListByResourceGroupResponse]](),
		beginResize:                       newTracker[azfake.PollerResponder[armhdinsight.ClustersClientResizeResponse]](),
		beginRotateDiskEncryptionKey:      newTracker[azfake.PollerResponder[armhdinsight.ClustersClientRotateDiskEncryptionKeyResponse]](),
		beginUpdateAutoScaleConfiguration: newTracker[azfake.PollerResponder[armhdinsight.ClustersClientUpdateAutoScaleConfigurationResponse]](),
		beginUpdateGatewaySettings:        newTracker[azfake.PollerResponder[armhdinsight.ClustersClientUpdateGatewaySettingsResponse]](),
		beginUpdateIdentityCertificate:    newTracker[azfake.PollerResponder[armhdinsight.ClustersClientUpdateIdentityCertificateResponse]](),
	}
}

// ClustersServerTransport connects instances of armhdinsight.ClustersClient to instances of ClustersServer.
// Don't use this type directly, use NewClustersServerTransport instead.
type ClustersServerTransport struct {
	srv                               *ClustersServer
	beginCreate                       *tracker[azfake.PollerResponder[armhdinsight.ClustersClientCreateResponse]]
	beginDelete                       *tracker[azfake.PollerResponder[armhdinsight.ClustersClientDeleteResponse]]
	beginExecuteScriptActions         *tracker[azfake.PollerResponder[armhdinsight.ClustersClientExecuteScriptActionsResponse]]
	newListPager                      *tracker[azfake.PagerResponder[armhdinsight.ClustersClientListResponse]]
	newListByResourceGroupPager       *tracker[azfake.PagerResponder[armhdinsight.ClustersClientListByResourceGroupResponse]]
	beginResize                       *tracker[azfake.PollerResponder[armhdinsight.ClustersClientResizeResponse]]
	beginRotateDiskEncryptionKey      *tracker[azfake.PollerResponder[armhdinsight.ClustersClientRotateDiskEncryptionKeyResponse]]
	beginUpdateAutoScaleConfiguration *tracker[azfake.PollerResponder[armhdinsight.ClustersClientUpdateAutoScaleConfigurationResponse]]
	beginUpdateGatewaySettings        *tracker[azfake.PollerResponder[armhdinsight.ClustersClientUpdateGatewaySettingsResponse]]
	beginUpdateIdentityCertificate    *tracker[azfake.PollerResponder[armhdinsight.ClustersClientUpdateIdentityCertificateResponse]]
}

// Do implements the policy.Transporter interface for ClustersServerTransport.
func (c *ClustersServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "ClustersClient.BeginCreate":
		resp, err = c.dispatchBeginCreate(req)
	case "ClustersClient.BeginDelete":
		resp, err = c.dispatchBeginDelete(req)
	case "ClustersClient.BeginExecuteScriptActions":
		resp, err = c.dispatchBeginExecuteScriptActions(req)
	case "ClustersClient.Get":
		resp, err = c.dispatchGet(req)
	case "ClustersClient.GetAzureAsyncOperationStatus":
		resp, err = c.dispatchGetAzureAsyncOperationStatus(req)
	case "ClustersClient.GetGatewaySettings":
		resp, err = c.dispatchGetGatewaySettings(req)
	case "ClustersClient.NewListPager":
		resp, err = c.dispatchNewListPager(req)
	case "ClustersClient.NewListByResourceGroupPager":
		resp, err = c.dispatchNewListByResourceGroupPager(req)
	case "ClustersClient.BeginResize":
		resp, err = c.dispatchBeginResize(req)
	case "ClustersClient.BeginRotateDiskEncryptionKey":
		resp, err = c.dispatchBeginRotateDiskEncryptionKey(req)
	case "ClustersClient.Update":
		resp, err = c.dispatchUpdate(req)
	case "ClustersClient.BeginUpdateAutoScaleConfiguration":
		resp, err = c.dispatchBeginUpdateAutoScaleConfiguration(req)
	case "ClustersClient.BeginUpdateGatewaySettings":
		resp, err = c.dispatchBeginUpdateGatewaySettings(req)
	case "ClustersClient.BeginUpdateIdentityCertificate":
		resp, err = c.dispatchBeginUpdateIdentityCertificate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (c *ClustersServerTransport) dispatchBeginCreate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginCreate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginCreate not implemented")}
	}
	beginCreate := c.beginCreate.get(req)
	if beginCreate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhdinsight.ClusterCreateParametersExtended](req)
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
		respr, errRespr := c.srv.BeginCreate(req.Context(), resourceGroupNameParam, clusterNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginCreate = &respr
		c.beginCreate.add(req, beginCreate)
	}

	resp, err := server.PollerResponderNext(beginCreate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.beginCreate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginCreate) {
		c.beginCreate.remove(req)
	}

	return resp, nil
}

func (c *ClustersServerTransport) dispatchBeginDelete(req *http.Request) (*http.Response, error) {
	if c.srv.BeginDelete == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginDelete not implemented")}
	}
	beginDelete := c.beginDelete.get(req)
	if beginDelete == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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

func (c *ClustersServerTransport) dispatchBeginExecuteScriptActions(req *http.Request) (*http.Response, error) {
	if c.srv.BeginExecuteScriptActions == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginExecuteScriptActions not implemented")}
	}
	beginExecuteScriptActions := c.beginExecuteScriptActions.get(req)
	if beginExecuteScriptActions == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/executeScriptActions`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhdinsight.ExecuteScriptActionParameters](req)
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
		respr, errRespr := c.srv.BeginExecuteScriptActions(req.Context(), resourceGroupNameParam, clusterNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginExecuteScriptActions = &respr
		c.beginExecuteScriptActions.add(req, beginExecuteScriptActions)
	}

	resp, err := server.PollerResponderNext(beginExecuteScriptActions, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		c.beginExecuteScriptActions.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginExecuteScriptActions) {
		c.beginExecuteScriptActions.remove(req)
	}

	return resp, nil
}

func (c *ClustersServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if c.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
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
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Cluster, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ClustersServerTransport) dispatchGetAzureAsyncOperationStatus(req *http.Request) (*http.Response, error) {
	if c.srv.GetAzureAsyncOperationStatus == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetAzureAsyncOperationStatus not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/azureasyncoperations/(?P<operationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
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
	operationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("operationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := c.srv.GetAzureAsyncOperationStatus(req.Context(), resourceGroupNameParam, clusterNameParam, operationIDParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).AsyncOperationResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ClustersServerTransport) dispatchGetGatewaySettings(req *http.Request) (*http.Response, error) {
	if c.srv.GetGatewaySettings == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetGatewaySettings not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/getGatewaySettings`
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
	respr, errRespr := c.srv.GetGatewaySettings(req.Context(), resourceGroupNameParam, clusterNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).GatewaySettings, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ClustersServerTransport) dispatchNewListPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListPager not implemented")}
	}
	newListPager := c.newListPager.get(req)
	if newListPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 1 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		resp := c.srv.NewListPager(nil)
		newListPager = &resp
		c.newListPager.add(req, newListPager)
		server.PagerResponderInjectNextLinks(newListPager, req, func(page *armhdinsight.ClustersClientListResponse, createLink func() string) {
			page.NextLink = to.Ptr(createLink())
		})
	}
	resp, err := server.PagerResponderNext(newListPager, req)
	if err != nil {
		return nil, err
	}
	if !contains([]int{http.StatusOK}, resp.StatusCode) {
		c.newListPager.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", resp.StatusCode)}
	}
	if !server.PagerResponderMore(newListPager) {
		c.newListPager.remove(req)
	}
	return resp, nil
}

func (c *ClustersServerTransport) dispatchNewListByResourceGroupPager(req *http.Request) (*http.Response, error) {
	if c.srv.NewListByResourceGroupPager == nil {
		return nil, &nonRetriableError{errors.New("fake for method NewListByResourceGroupPager not implemented")}
	}
	newListByResourceGroupPager := c.newListByResourceGroupPager.get(req)
	if newListByResourceGroupPager == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters`
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
		server.PagerResponderInjectNextLinks(newListByResourceGroupPager, req, func(page *armhdinsight.ClustersClientListByResourceGroupResponse, createLink func() string) {
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

func (c *ClustersServerTransport) dispatchBeginResize(req *http.Request) (*http.Response, error) {
	if c.srv.BeginResize == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginResize not implemented")}
	}
	beginResize := c.beginResize.get(req)
	if beginResize == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/roles/(?P<roleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resize`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhdinsight.ClusterResizeParameters](req)
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
		roleNameParam, err := parseWithCast(matches[regex.SubexpIndex("roleName")], func(v string) (armhdinsight.RoleName, error) {
			p, unescapeErr := url.PathUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armhdinsight.RoleName(p), nil
		})
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginResize(req.Context(), resourceGroupNameParam, clusterNameParam, roleNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginResize = &respr
		c.beginResize.add(req, beginResize)
	}

	resp, err := server.PollerResponderNext(beginResize, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		c.beginResize.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginResize) {
		c.beginResize.remove(req)
	}

	return resp, nil
}

func (c *ClustersServerTransport) dispatchBeginRotateDiskEncryptionKey(req *http.Request) (*http.Response, error) {
	if c.srv.BeginRotateDiskEncryptionKey == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginRotateDiskEncryptionKey not implemented")}
	}
	beginRotateDiskEncryptionKey := c.beginRotateDiskEncryptionKey.get(req)
	if beginRotateDiskEncryptionKey == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/rotatediskencryptionkey`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhdinsight.ClusterDiskEncryptionParameters](req)
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
		respr, errRespr := c.srv.BeginRotateDiskEncryptionKey(req.Context(), resourceGroupNameParam, clusterNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginRotateDiskEncryptionKey = &respr
		c.beginRotateDiskEncryptionKey.add(req, beginRotateDiskEncryptionKey)
	}

	resp, err := server.PollerResponderNext(beginRotateDiskEncryptionKey, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		c.beginRotateDiskEncryptionKey.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginRotateDiskEncryptionKey) {
		c.beginRotateDiskEncryptionKey.remove(req)
	}

	return resp, nil
}

func (c *ClustersServerTransport) dispatchUpdate(req *http.Request) (*http.Response, error) {
	if c.srv.Update == nil {
		return nil, &nonRetriableError{errors.New("fake for method Update not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armhdinsight.ClusterPatchParameters](req)
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
	respr, errRespr := c.srv.Update(req.Context(), resourceGroupNameParam, clusterNameParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).Cluster, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *ClustersServerTransport) dispatchBeginUpdateAutoScaleConfiguration(req *http.Request) (*http.Response, error) {
	if c.srv.BeginUpdateAutoScaleConfiguration == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdateAutoScaleConfiguration not implemented")}
	}
	beginUpdateAutoScaleConfiguration := c.beginUpdateAutoScaleConfiguration.get(req)
	if beginUpdateAutoScaleConfiguration == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/roles/(?P<roleName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/autoscale`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 4 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhdinsight.AutoscaleConfigurationUpdateParameter](req)
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
		roleNameParam, err := parseWithCast(matches[regex.SubexpIndex("roleName")], func(v string) (armhdinsight.RoleName, error) {
			p, unescapeErr := url.PathUnescape(v)
			if unescapeErr != nil {
				return "", unescapeErr
			}
			return armhdinsight.RoleName(p), nil
		})
		if err != nil {
			return nil, err
		}
		respr, errRespr := c.srv.BeginUpdateAutoScaleConfiguration(req.Context(), resourceGroupNameParam, clusterNameParam, roleNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdateAutoScaleConfiguration = &respr
		c.beginUpdateAutoScaleConfiguration.add(req, beginUpdateAutoScaleConfiguration)
	}

	resp, err := server.PollerResponderNext(beginUpdateAutoScaleConfiguration, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		c.beginUpdateAutoScaleConfiguration.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdateAutoScaleConfiguration) {
		c.beginUpdateAutoScaleConfiguration.remove(req)
	}

	return resp, nil
}

func (c *ClustersServerTransport) dispatchBeginUpdateGatewaySettings(req *http.Request) (*http.Response, error) {
	if c.srv.BeginUpdateGatewaySettings == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdateGatewaySettings not implemented")}
	}
	beginUpdateGatewaySettings := c.beginUpdateGatewaySettings.get(req)
	if beginUpdateGatewaySettings == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/updateGatewaySettings`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhdinsight.UpdateGatewaySettingsParameters](req)
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
		respr, errRespr := c.srv.BeginUpdateGatewaySettings(req.Context(), resourceGroupNameParam, clusterNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdateGatewaySettings = &respr
		c.beginUpdateGatewaySettings.add(req, beginUpdateGatewaySettings)
	}

	resp, err := server.PollerResponderNext(beginUpdateGatewaySettings, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		c.beginUpdateGatewaySettings.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdateGatewaySettings) {
		c.beginUpdateGatewaySettings.remove(req)
	}

	return resp, nil
}

func (c *ClustersServerTransport) dispatchBeginUpdateIdentityCertificate(req *http.Request) (*http.Response, error) {
	if c.srv.BeginUpdateIdentityCertificate == nil {
		return nil, &nonRetriableError{errors.New("fake for method BeginUpdateIdentityCertificate not implemented")}
	}
	beginUpdateIdentityCertificate := c.beginUpdateIdentityCertificate.get(req)
	if beginUpdateIdentityCertificate == nil {
		const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/clusters/(?P<clusterName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/updateClusterIdentityCertificate`
		regex := regexp.MustCompile(regexStr)
		matches := regex.FindStringSubmatch(req.URL.EscapedPath())
		if matches == nil || len(matches) < 3 {
			return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
		}
		body, err := server.UnmarshalRequestAsJSON[armhdinsight.UpdateClusterIdentityCertificateParameters](req)
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
		respr, errRespr := c.srv.BeginUpdateIdentityCertificate(req.Context(), resourceGroupNameParam, clusterNameParam, body, nil)
		if respErr := server.GetError(errRespr, req); respErr != nil {
			return nil, respErr
		}
		beginUpdateIdentityCertificate = &respr
		c.beginUpdateIdentityCertificate.add(req, beginUpdateIdentityCertificate)
	}

	resp, err := server.PollerResponderNext(beginUpdateIdentityCertificate, req)
	if err != nil {
		return nil, err
	}

	if !contains([]int{http.StatusOK, http.StatusAccepted}, resp.StatusCode) {
		c.beginUpdateIdentityCertificate.remove(req)
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK, http.StatusAccepted", resp.StatusCode)}
	}
	if !server.PollerResponderMore(beginUpdateIdentityCertificate) {
		c.beginUpdateIdentityCertificate.remove(req)
	}

	return resp, nil
}
