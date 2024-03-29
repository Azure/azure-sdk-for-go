//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package fake

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	azfake "github.com/Azure/azure-sdk-for-go/sdk/azcore/fake"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/fake/server"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/deviceprovisioningservices/armdeviceprovisioningservices"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

// DpsCertificateServer is a fake server for instances of the armdeviceprovisioningservices.DpsCertificateClient type.
type DpsCertificateServer struct {
	// CreateOrUpdate is the fake for method DpsCertificateClient.CreateOrUpdate
	// HTTP status codes to indicate success: http.StatusOK
	CreateOrUpdate func(ctx context.Context, resourceGroupName string, provisioningServiceName string, certificateName string, certificateDescription armdeviceprovisioningservices.CertificateResponse, options *armdeviceprovisioningservices.DpsCertificateClientCreateOrUpdateOptions) (resp azfake.Responder[armdeviceprovisioningservices.DpsCertificateClientCreateOrUpdateResponse], errResp azfake.ErrorResponder)

	// Delete is the fake for method DpsCertificateClient.Delete
	// HTTP status codes to indicate success: http.StatusOK, http.StatusNoContent
	Delete func(ctx context.Context, resourceGroupName string, ifMatch string, provisioningServiceName string, certificateName string, options *armdeviceprovisioningservices.DpsCertificateClientDeleteOptions) (resp azfake.Responder[armdeviceprovisioningservices.DpsCertificateClientDeleteResponse], errResp azfake.ErrorResponder)

	// GenerateVerificationCode is the fake for method DpsCertificateClient.GenerateVerificationCode
	// HTTP status codes to indicate success: http.StatusOK
	GenerateVerificationCode func(ctx context.Context, certificateName string, ifMatch string, resourceGroupName string, provisioningServiceName string, options *armdeviceprovisioningservices.DpsCertificateClientGenerateVerificationCodeOptions) (resp azfake.Responder[armdeviceprovisioningservices.DpsCertificateClientGenerateVerificationCodeResponse], errResp azfake.ErrorResponder)

	// Get is the fake for method DpsCertificateClient.Get
	// HTTP status codes to indicate success: http.StatusOK
	Get func(ctx context.Context, certificateName string, resourceGroupName string, provisioningServiceName string, options *armdeviceprovisioningservices.DpsCertificateClientGetOptions) (resp azfake.Responder[armdeviceprovisioningservices.DpsCertificateClientGetResponse], errResp azfake.ErrorResponder)

	// List is the fake for method DpsCertificateClient.List
	// HTTP status codes to indicate success: http.StatusOK
	List func(ctx context.Context, resourceGroupName string, provisioningServiceName string, options *armdeviceprovisioningservices.DpsCertificateClientListOptions) (resp azfake.Responder[armdeviceprovisioningservices.DpsCertificateClientListResponse], errResp azfake.ErrorResponder)

	// VerifyCertificate is the fake for method DpsCertificateClient.VerifyCertificate
	// HTTP status codes to indicate success: http.StatusOK
	VerifyCertificate func(ctx context.Context, certificateName string, ifMatch string, resourceGroupName string, provisioningServiceName string, request armdeviceprovisioningservices.VerificationCodeRequest, options *armdeviceprovisioningservices.DpsCertificateClientVerifyCertificateOptions) (resp azfake.Responder[armdeviceprovisioningservices.DpsCertificateClientVerifyCertificateResponse], errResp azfake.ErrorResponder)
}

// NewDpsCertificateServerTransport creates a new instance of DpsCertificateServerTransport with the provided implementation.
// The returned DpsCertificateServerTransport instance is connected to an instance of armdeviceprovisioningservices.DpsCertificateClient via the
// azcore.ClientOptions.Transporter field in the client's constructor parameters.
func NewDpsCertificateServerTransport(srv *DpsCertificateServer) *DpsCertificateServerTransport {
	return &DpsCertificateServerTransport{srv: srv}
}

// DpsCertificateServerTransport connects instances of armdeviceprovisioningservices.DpsCertificateClient to instances of DpsCertificateServer.
// Don't use this type directly, use NewDpsCertificateServerTransport instead.
type DpsCertificateServerTransport struct {
	srv *DpsCertificateServer
}

// Do implements the policy.Transporter interface for DpsCertificateServerTransport.
func (d *DpsCertificateServerTransport) Do(req *http.Request) (*http.Response, error) {
	rawMethod := req.Context().Value(runtime.CtxAPINameKey{})
	method, ok := rawMethod.(string)
	if !ok {
		return nil, nonRetriableError{errors.New("unable to dispatch request, missing value for CtxAPINameKey")}
	}

	var resp *http.Response
	var err error

	switch method {
	case "DpsCertificateClient.CreateOrUpdate":
		resp, err = d.dispatchCreateOrUpdate(req)
	case "DpsCertificateClient.Delete":
		resp, err = d.dispatchDelete(req)
	case "DpsCertificateClient.GenerateVerificationCode":
		resp, err = d.dispatchGenerateVerificationCode(req)
	case "DpsCertificateClient.Get":
		resp, err = d.dispatchGet(req)
	case "DpsCertificateClient.List":
		resp, err = d.dispatchList(req)
	case "DpsCertificateClient.VerifyCertificate":
		resp, err = d.dispatchVerifyCertificate(req)
	default:
		err = fmt.Errorf("unhandled API %s", method)
	}

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (d *DpsCertificateServerTransport) dispatchCreateOrUpdate(req *http.Request) (*http.Response, error) {
	if d.srv.CreateOrUpdate == nil {
		return nil, &nonRetriableError{errors.New("fake for method CreateOrUpdate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Devices/provisioningServices/(?P<provisioningServiceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/certificates/(?P<certificateName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	body, err := server.UnmarshalRequestAsJSON[armdeviceprovisioningservices.CertificateResponse](req)
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	provisioningServiceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("provisioningServiceName")])
	if err != nil {
		return nil, err
	}
	certificateNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("certificateName")])
	if err != nil {
		return nil, err
	}
	ifMatchParam := getOptional(getHeaderValue(req.Header, "If-Match"))
	var options *armdeviceprovisioningservices.DpsCertificateClientCreateOrUpdateOptions
	if ifMatchParam != nil {
		options = &armdeviceprovisioningservices.DpsCertificateClientCreateOrUpdateOptions{
			IfMatch: ifMatchParam,
		}
	}
	respr, errRespr := d.srv.CreateOrUpdate(req.Context(), resourceGroupNameParam, provisioningServiceNameParam, certificateNameParam, body, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CertificateResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DpsCertificateServerTransport) dispatchDelete(req *http.Request) (*http.Response, error) {
	if d.srv.Delete == nil {
		return nil, &nonRetriableError{errors.New("fake for method Delete not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Devices/provisioningServices/(?P<provisioningServiceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/certificates/(?P<certificateName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	provisioningServiceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("provisioningServiceName")])
	if err != nil {
		return nil, err
	}
	certificateNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("certificateName")])
	if err != nil {
		return nil, err
	}
	certificateName1Unescaped, err := url.QueryUnescape(qp.Get("certificate.name"))
	if err != nil {
		return nil, err
	}
	certificateName1Param := getOptional(certificateName1Unescaped)
	certificateRawBytesUnescaped, err := url.QueryUnescape(qp.Get("certificate.rawBytes"))
	if err != nil {
		return nil, err
	}
	certificateRawBytesParam, err := base64.StdEncoding.DecodeString(certificateRawBytesUnescaped)
	if err != nil {
		return nil, err
	}
	certificateIsVerifiedUnescaped, err := url.QueryUnescape(qp.Get("certificate.isVerified"))
	if err != nil {
		return nil, err
	}
	certificateIsVerifiedParam, err := parseOptional(certificateIsVerifiedUnescaped, strconv.ParseBool)
	if err != nil {
		return nil, err
	}
	certificatePurposeUnescaped, err := url.QueryUnescape(qp.Get("certificate.purpose"))
	if err != nil {
		return nil, err
	}
	certificatePurposeParam := getOptional(armdeviceprovisioningservices.CertificatePurpose(certificatePurposeUnescaped))
	certificateCreatedUnescaped, err := url.QueryUnescape(qp.Get("certificate.created"))
	if err != nil {
		return nil, err
	}
	certificateCreatedParam, err := parseOptional(certificateCreatedUnescaped, func(v string) (time.Time, error) { return time.Parse(time.RFC3339Nano, v) })
	if err != nil {
		return nil, err
	}
	certificateLastUpdatedUnescaped, err := url.QueryUnescape(qp.Get("certificate.lastUpdated"))
	if err != nil {
		return nil, err
	}
	certificateLastUpdatedParam, err := parseOptional(certificateLastUpdatedUnescaped, func(v string) (time.Time, error) { return time.Parse(time.RFC3339Nano, v) })
	if err != nil {
		return nil, err
	}
	certificateHasPrivateKeyUnescaped, err := url.QueryUnescape(qp.Get("certificate.hasPrivateKey"))
	if err != nil {
		return nil, err
	}
	certificateHasPrivateKeyParam, err := parseOptional(certificateHasPrivateKeyUnescaped, strconv.ParseBool)
	if err != nil {
		return nil, err
	}
	certificateNonceUnescaped, err := url.QueryUnescape(qp.Get("certificate.nonce"))
	if err != nil {
		return nil, err
	}
	certificateNonceParam := getOptional(certificateNonceUnescaped)
	var options *armdeviceprovisioningservices.DpsCertificateClientDeleteOptions
	if certificateName1Param != nil || certificateRawBytesParam != nil || certificateIsVerifiedParam != nil || certificatePurposeParam != nil || certificateCreatedParam != nil || certificateLastUpdatedParam != nil || certificateHasPrivateKeyParam != nil || certificateNonceParam != nil {
		options = &armdeviceprovisioningservices.DpsCertificateClientDeleteOptions{
			CertificateName1:         certificateName1Param,
			CertificateRawBytes:      certificateRawBytesParam,
			CertificateIsVerified:    certificateIsVerifiedParam,
			CertificatePurpose:       certificatePurposeParam,
			CertificateCreated:       certificateCreatedParam,
			CertificateLastUpdated:   certificateLastUpdatedParam,
			CertificateHasPrivateKey: certificateHasPrivateKeyParam,
			CertificateNonce:         certificateNonceParam,
		}
	}
	respr, errRespr := d.srv.Delete(req.Context(), resourceGroupNameParam, getHeaderValue(req.Header, "If-Match"), provisioningServiceNameParam, certificateNameParam, options)
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

func (d *DpsCertificateServerTransport) dispatchGenerateVerificationCode(req *http.Request) (*http.Response, error) {
	if d.srv.GenerateVerificationCode == nil {
		return nil, &nonRetriableError{errors.New("fake for method GenerateVerificationCode not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Devices/provisioningServices/(?P<provisioningServiceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/certificates/(?P<certificateName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/generateVerificationCode`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	certificateNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("certificateName")])
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	provisioningServiceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("provisioningServiceName")])
	if err != nil {
		return nil, err
	}
	certificateName1Unescaped, err := url.QueryUnescape(qp.Get("certificate.name"))
	if err != nil {
		return nil, err
	}
	certificateName1Param := getOptional(certificateName1Unescaped)
	certificateRawBytesUnescaped, err := url.QueryUnescape(qp.Get("certificate.rawBytes"))
	if err != nil {
		return nil, err
	}
	certificateRawBytesParam, err := base64.StdEncoding.DecodeString(certificateRawBytesUnescaped)
	if err != nil {
		return nil, err
	}
	certificateIsVerifiedUnescaped, err := url.QueryUnescape(qp.Get("certificate.isVerified"))
	if err != nil {
		return nil, err
	}
	certificateIsVerifiedParam, err := parseOptional(certificateIsVerifiedUnescaped, strconv.ParseBool)
	if err != nil {
		return nil, err
	}
	certificatePurposeUnescaped, err := url.QueryUnescape(qp.Get("certificate.purpose"))
	if err != nil {
		return nil, err
	}
	certificatePurposeParam := getOptional(armdeviceprovisioningservices.CertificatePurpose(certificatePurposeUnescaped))
	certificateCreatedUnescaped, err := url.QueryUnescape(qp.Get("certificate.created"))
	if err != nil {
		return nil, err
	}
	certificateCreatedParam, err := parseOptional(certificateCreatedUnescaped, func(v string) (time.Time, error) { return time.Parse(time.RFC3339Nano, v) })
	if err != nil {
		return nil, err
	}
	certificateLastUpdatedUnescaped, err := url.QueryUnescape(qp.Get("certificate.lastUpdated"))
	if err != nil {
		return nil, err
	}
	certificateLastUpdatedParam, err := parseOptional(certificateLastUpdatedUnescaped, func(v string) (time.Time, error) { return time.Parse(time.RFC3339Nano, v) })
	if err != nil {
		return nil, err
	}
	certificateHasPrivateKeyUnescaped, err := url.QueryUnescape(qp.Get("certificate.hasPrivateKey"))
	if err != nil {
		return nil, err
	}
	certificateHasPrivateKeyParam, err := parseOptional(certificateHasPrivateKeyUnescaped, strconv.ParseBool)
	if err != nil {
		return nil, err
	}
	certificateNonceUnescaped, err := url.QueryUnescape(qp.Get("certificate.nonce"))
	if err != nil {
		return nil, err
	}
	certificateNonceParam := getOptional(certificateNonceUnescaped)
	var options *armdeviceprovisioningservices.DpsCertificateClientGenerateVerificationCodeOptions
	if certificateName1Param != nil || certificateRawBytesParam != nil || certificateIsVerifiedParam != nil || certificatePurposeParam != nil || certificateCreatedParam != nil || certificateLastUpdatedParam != nil || certificateHasPrivateKeyParam != nil || certificateNonceParam != nil {
		options = &armdeviceprovisioningservices.DpsCertificateClientGenerateVerificationCodeOptions{
			CertificateName1:         certificateName1Param,
			CertificateRawBytes:      certificateRawBytesParam,
			CertificateIsVerified:    certificateIsVerifiedParam,
			CertificatePurpose:       certificatePurposeParam,
			CertificateCreated:       certificateCreatedParam,
			CertificateLastUpdated:   certificateLastUpdatedParam,
			CertificateHasPrivateKey: certificateHasPrivateKeyParam,
			CertificateNonce:         certificateNonceParam,
		}
	}
	respr, errRespr := d.srv.GenerateVerificationCode(req.Context(), certificateNameParam, getHeaderValue(req.Header, "If-Match"), resourceGroupNameParam, provisioningServiceNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).VerificationCodeResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DpsCertificateServerTransport) dispatchGet(req *http.Request) (*http.Response, error) {
	if d.srv.Get == nil {
		return nil, &nonRetriableError{errors.New("fake for method Get not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Devices/provisioningServices/(?P<provisioningServiceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/certificates/(?P<certificateName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	certificateNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("certificateName")])
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	provisioningServiceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("provisioningServiceName")])
	if err != nil {
		return nil, err
	}
	ifMatchParam := getOptional(getHeaderValue(req.Header, "If-Match"))
	var options *armdeviceprovisioningservices.DpsCertificateClientGetOptions
	if ifMatchParam != nil {
		options = &armdeviceprovisioningservices.DpsCertificateClientGetOptions{
			IfMatch: ifMatchParam,
		}
	}
	respr, errRespr := d.srv.Get(req.Context(), certificateNameParam, resourceGroupNameParam, provisioningServiceNameParam, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CertificateResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DpsCertificateServerTransport) dispatchList(req *http.Request) (*http.Response, error) {
	if d.srv.List == nil {
		return nil, &nonRetriableError{errors.New("fake for method List not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Devices/provisioningServices/(?P<provisioningServiceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/certificates`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 3 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	provisioningServiceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("provisioningServiceName")])
	if err != nil {
		return nil, err
	}
	respr, errRespr := d.srv.List(req.Context(), resourceGroupNameParam, provisioningServiceNameParam, nil)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CertificateListDescription, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (d *DpsCertificateServerTransport) dispatchVerifyCertificate(req *http.Request) (*http.Response, error) {
	if d.srv.VerifyCertificate == nil {
		return nil, &nonRetriableError{errors.New("fake for method VerifyCertificate not implemented")}
	}
	const regexStr = `/subscriptions/(?P<subscriptionId>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/resourceGroups/(?P<resourceGroupName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/providers/Microsoft\.Devices/provisioningServices/(?P<provisioningServiceName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/certificates/(?P<certificateName>[!#&$-;=?-\[\]_a-zA-Z0-9~%@]+)/verify`
	regex := regexp.MustCompile(regexStr)
	matches := regex.FindStringSubmatch(req.URL.EscapedPath())
	if matches == nil || len(matches) < 4 {
		return nil, fmt.Errorf("failed to parse path %s", req.URL.Path)
	}
	qp := req.URL.Query()
	body, err := server.UnmarshalRequestAsJSON[armdeviceprovisioningservices.VerificationCodeRequest](req)
	if err != nil {
		return nil, err
	}
	certificateNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("certificateName")])
	if err != nil {
		return nil, err
	}
	resourceGroupNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("resourceGroupName")])
	if err != nil {
		return nil, err
	}
	provisioningServiceNameParam, err := url.PathUnescape(matches[regex.SubexpIndex("provisioningServiceName")])
	if err != nil {
		return nil, err
	}
	certificateName1Unescaped, err := url.QueryUnescape(qp.Get("certificate.name"))
	if err != nil {
		return nil, err
	}
	certificateName1Param := getOptional(certificateName1Unescaped)
	certificateRawBytesUnescaped, err := url.QueryUnescape(qp.Get("certificate.rawBytes"))
	if err != nil {
		return nil, err
	}
	certificateRawBytesParam, err := base64.StdEncoding.DecodeString(certificateRawBytesUnescaped)
	if err != nil {
		return nil, err
	}
	certificateIsVerifiedUnescaped, err := url.QueryUnescape(qp.Get("certificate.isVerified"))
	if err != nil {
		return nil, err
	}
	certificateIsVerifiedParam, err := parseOptional(certificateIsVerifiedUnescaped, strconv.ParseBool)
	if err != nil {
		return nil, err
	}
	certificatePurposeUnescaped, err := url.QueryUnescape(qp.Get("certificate.purpose"))
	if err != nil {
		return nil, err
	}
	certificatePurposeParam := getOptional(armdeviceprovisioningservices.CertificatePurpose(certificatePurposeUnescaped))
	certificateCreatedUnescaped, err := url.QueryUnescape(qp.Get("certificate.created"))
	if err != nil {
		return nil, err
	}
	certificateCreatedParam, err := parseOptional(certificateCreatedUnescaped, func(v string) (time.Time, error) { return time.Parse(time.RFC3339Nano, v) })
	if err != nil {
		return nil, err
	}
	certificateLastUpdatedUnescaped, err := url.QueryUnescape(qp.Get("certificate.lastUpdated"))
	if err != nil {
		return nil, err
	}
	certificateLastUpdatedParam, err := parseOptional(certificateLastUpdatedUnescaped, func(v string) (time.Time, error) { return time.Parse(time.RFC3339Nano, v) })
	if err != nil {
		return nil, err
	}
	certificateHasPrivateKeyUnescaped, err := url.QueryUnescape(qp.Get("certificate.hasPrivateKey"))
	if err != nil {
		return nil, err
	}
	certificateHasPrivateKeyParam, err := parseOptional(certificateHasPrivateKeyUnescaped, strconv.ParseBool)
	if err != nil {
		return nil, err
	}
	certificateNonceUnescaped, err := url.QueryUnescape(qp.Get("certificate.nonce"))
	if err != nil {
		return nil, err
	}
	certificateNonceParam := getOptional(certificateNonceUnescaped)
	var options *armdeviceprovisioningservices.DpsCertificateClientVerifyCertificateOptions
	if certificateName1Param != nil || certificateRawBytesParam != nil || certificateIsVerifiedParam != nil || certificatePurposeParam != nil || certificateCreatedParam != nil || certificateLastUpdatedParam != nil || certificateHasPrivateKeyParam != nil || certificateNonceParam != nil {
		options = &armdeviceprovisioningservices.DpsCertificateClientVerifyCertificateOptions{
			CertificateName1:         certificateName1Param,
			CertificateRawBytes:      certificateRawBytesParam,
			CertificateIsVerified:    certificateIsVerifiedParam,
			CertificatePurpose:       certificatePurposeParam,
			CertificateCreated:       certificateCreatedParam,
			CertificateLastUpdated:   certificateLastUpdatedParam,
			CertificateHasPrivateKey: certificateHasPrivateKeyParam,
			CertificateNonce:         certificateNonceParam,
		}
	}
	respr, errRespr := d.srv.VerifyCertificate(req.Context(), certificateNameParam, getHeaderValue(req.Header, "If-Match"), resourceGroupNameParam, provisioningServiceNameParam, body, options)
	if respErr := server.GetError(errRespr, req); respErr != nil {
		return nil, respErr
	}
	respContent := server.GetResponseContent(respr)
	if !contains([]int{http.StatusOK}, respContent.HTTPStatus) {
		return nil, &nonRetriableError{fmt.Errorf("unexpected status code %d. acceptable values are http.StatusOK", respContent.HTTPStatus)}
	}
	resp, err := server.MarshalResponseAsJSON(respContent, server.GetResponse(respr).CertificateResponse, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
