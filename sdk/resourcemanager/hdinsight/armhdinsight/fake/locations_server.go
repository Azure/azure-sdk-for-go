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
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/hdinsight/armhdinsight"
	"net/http"
	"net/url"
	"regexp"
)

// LocationsServer is a fake server for instances of the armhdinsight.LocationsClient type.
type LocationsServer struct {
	// CheckNameAvailability is the fake for method LocationsClient.CheckNameAvailability
	// HTTP status codes to indicate success: http.StatusOK
	CheckNameAvailability func(ctx context.Context, location string, parameters armhdinsight.NameAvailabilityCheckRequestParameters, options *armhdinsight.LocationsClientCheckNameAvailabilityOptions) (resp azfake.Responder[armhdinsight.LocationsClientCheckNameAvailabilityResponse], errResp azfake.ErrorResponder)

	// GetAzureAsyncOperationStatus is the fake for method LocationsClient.GetAzureAsyncOperationStatus
	// HTTP status codes to indicate success: http.StatusOK
	GetAzureAsyncOperationStatus func(ctx context.Context, location string, operationID string, options *armhdinsight.LocationsClientGetAzureAsyncOperationStatusOptions) (resp azfake.Responder[armhdinsight.LocationsClientGetAzureAsyncOperationStatusResponse], errResp azfake.ErrorResponder)

	// GetCapabilities is the fake for method LocationsClient.GetCapabilities
	// HTTP status codes to indicate success: http.StatusOK
	GetCapabilities func(ctx context.Context, location string, options *armhdinsight.LocationsClientGetCapabilitiesOptions) (resp azfake.Responder[armhdinsight.LocationsClientGetCapabilitiesResponse], errResp azfake.ErrorResponder)

	// ListBillingSpecs is the fake for method LocationsClient.ListBillingSpecs
	// HTTP status codes to indicate success: http.StatusOK
	ListBillingSpecs func(ctx context.Context, location string, options *armhdinsight.LocationsClientListBillingSpecsOptions) (resp azfake.Responder[armhdinsight.LocationsClientListBillingSpecsResponse], errResp azfake.ErrorResponder)

	// ListUsages is the fake for method LocationsClient.ListUsages
	// HTTP status codes to indicate success: http.StatusOK
	ListUsages func(ctx context.Context, location string, options *armhdinsight.LocationsClientListUsagesOptions) (resp azfake.Responder[armhdinsight.LocationsClientListUsagesResponse], errResp azfake.ErrorResponder)

	// ValidateClusterCreateRequest is the fake for method LocationsClient.ValidateClusterCreateRequest
	// HTTP status codes to indicate success: http.StatusOK
	ValidateClusterCreateRequest func(ctx context.Context, location string, parameters armhdinsight.ClusterCreateRequestValidationParameters, options *armhdinsight.LocationsClientValidateClusterCreateRequestOptions) (resp azfake.Responder[armhdinsight.LocationsClientValidateClusterCreateRequestResponse], errResp azfake.ErrorResponder)
}

// NewLocationsServerTransport creates a new instance of LocationsServerTransport with the provided implementation.
// The returned LocationsServerTransport instance is connected to an instance of armhdinsight.LocationsClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewLocationsServerTransport(srv *LocationsServer) *LocationsServerTransport {
	return &LocationsServerTransport{srv: srv}
}

// LocationsServerTransport connects instances of armhdinsight.LocationsClient to instances of LocationsServer.
// Don't use this type directly, use NewLocationsServerTransport instead.
type LocationsServerTransport struct {
	srv *LocationsServer
}

// Do implements the policy.Transporter interface for LocationsServerTransport.
func (l *LocationsServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "LocationsClient.CheckNameAvailability":
		resp, err = l.dispatchCheckNameAvailability(req)
	case "LocationsClient.GetAzureAsyncOperationStatus":
		resp, err = l.dispatchGetAzureAsyncOperationStatus(req)
	case "LocationsClient.GetCapabilities":
		resp, err = l.dispatchGetCapabilities(req)
	case "LocationsClient.ListBillingSpecs":
		resp, err = l.dispatchListBillingSpecs(req)
	case "LocationsClient.ListUsages":
		resp, err = l.dispatchListUsages(req)
	case "LocationsClient.ValidateClusterCreateRequest":
		resp, err = l.dispatchValidateClusterCreateRequest(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (l *LocationsServerTransport) dispatchCheckNameAvailability(req *http.Request) (*http.Response, error) {
	if l.srv.CheckNameAvailability == nil {
		return nil, &nonRetriableError{errors.New("fake for method CheckNameAvailability not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/checkNameAvailability`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armhdinsight.NameAvailabilityCheckRequestParameters](req)
	if err != nil {
		return nil, err
	}
	locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := l.srv.CheckNameAvailability(req.Context(), locationParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).NameAvailabilityCheckResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *LocationsServerTransport) dispatchGetAzureAsyncOperationStatus(req *http.Request) (*http.Response, error) {
	if l.srv.GetAzureAsyncOperationStatus == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetAzureAsyncOperationStatus not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/azureasyncoperations/(?P<operationId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
	if err != nil {
		return nil, err
	}
	operationIDParam, err := url.PathUnescape(matches[regex.SubexpIndex("operationId")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := l.srv.GetAzureAsyncOperationStatus(req.Context(), locationParam, operationIDParam, nil)
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

func (l *LocationsServerTransport) dispatchGetCapabilities(req *http.Request) (*http.Response, error) {
	if l.srv.GetCapabilities == nil {
		return nil, &nonRetriableError{errors.New("fake for method GetCapabilities not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/capabilities`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := l.srv.GetCapabilities(req.Context(), locationParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CapabilitiesResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *LocationsServerTransport) dispatchListBillingSpecs(req *http.Request) (*http.Response, error) {
	if l.srv.ListBillingSpecs == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListBillingSpecs not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/billingSpecs`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := l.srv.ListBillingSpecs(req.Context(), locationParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).BillingResponseListResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *LocationsServerTransport) dispatchListUsages(req *http.Request) (*http.Response, error) {
	if l.srv.ListUsages == nil {
		return nil, &nonRetriableError{errors.New("fake for method ListUsages not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/usages`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := l.srv.ListUsages(req.Context(), locationParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).UsagesListResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (l *LocationsServerTransport) dispatchValidateClusterCreateRequest(req *http.Request) (*http.Response, error) {
	if l.srv.ValidateClusterCreateRequest == nil {
		return nil, &nonRetriableError{errors.New("fake for method ValidateClusterCreateRequest not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.HDInsight/locations/(?P<location>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/validateCreateRequest`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 2 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armhdinsight.ClusterCreateRequestValidationParameters](req)
	if err != nil {
		return nil, err
	}
	locationParam, err := url.PathUnescape(matches[regex.SubexpIndex("location")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := l.srv.ValidateClusterCreateRequest(req.Context(), locationParam, body, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).ClusterCreateValidationResult, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
