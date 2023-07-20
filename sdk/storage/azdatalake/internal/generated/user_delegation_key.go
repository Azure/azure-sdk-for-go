//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"strconv"
	"time"
)

// KeyInfo - Key information
type KeyInfo struct {
	// REQUIRED; The date-time the key expires in ISO 8601 UTC time
	Expiry *string `xml:"Expiry"`

	// REQUIRED; The date-time the key is active in ISO 8601 UTC time
	Start *string `xml:"Start"`
}

// ServiceClientGetUserDelegationKeyOptions contains the optional parameters for the ServiceClient.GetUserDelegationKey method.
type ServiceClientGetUserDelegationKeyOptions struct {
	// Provides a client-generated, opaque value with a 1 KB character limit that is recorded in the analytics logs when storage
	// analytics logging is enabled.
	RequestID *string
	// The timeout parameter is expressed in seconds. For more information, see Setting Timeouts for Blob Service Operations.
	// [https://docs.microsoft.com/en-us/rest/api/storageservices/fileservices/setting-timeouts-for-blob-service-operations]
	Timeout *int32
}

// ServiceClientGetUserDelegationKeyResponse contains the response from method ServiceClient.GetUserDelegationKey.
type ServiceClientGetUserDelegationKeyResponse struct {
	UserDelegationKey
	// ClientRequestID contains the information returned from the x-ms-client-request-id header response.
	ClientRequestID *string `xml:"ClientRequestID"`

	// Date contains the information returned from the Date header response.
	Date *time.Time `xml:"Date"`

	// RequestID contains the information returned from the x-ms-request-id header response.
	RequestID *string `xml:"RequestID"`

	// Version contains the information returned from the x-ms-version header response.
	Version *string `xml:"Version"`
}

// UserDelegationKey - A user delegation key
type UserDelegationKey struct {
	// REQUIRED; The date-time the key expires
	SignedExpiry *time.Time `xml:"SignedExpiry"`

	// REQUIRED; The Azure Active Directory object ID in GUID format.
	SignedOID *string `xml:"SignedOid"`

	// REQUIRED; Abbreviation of the Azure Storage service that accepts the key
	SignedService *string `xml:"SignedService"`

	// REQUIRED; The date-time the key is active
	SignedStart *time.Time `xml:"SignedStart"`

	// REQUIRED; The Azure Active Directory tenant ID in GUID format
	SignedTID *string `xml:"SignedTid"`

	// REQUIRED; The service version that created the key
	SignedVersion *string `xml:"SignedVersion"`

	// REQUIRED; The key as a base64 string
	Value *string `xml:"Value"`
}

// GetUserDelegationKey - Retrieves a user delegation key for the Blob service. This is only a valid operation when using
// bearer token authentication.
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2020-10-02
//   - keyInfo - Key information
//   - options - ServiceClientGetUserDelegationKeyOptions contains the optional parameters for the ServiceClient.GetUserDelegationKey
//     method.
func (client *ServiceClient) GetUserDelegationKey(ctx context.Context, keyInfo KeyInfo, options *ServiceClientGetUserDelegationKeyOptions) (ServiceClientGetUserDelegationKeyResponse, error) {
	req, err := client.getUserDelegationKeyCreateRequest(ctx, keyInfo, options)
	if err != nil {
		return ServiceClientGetUserDelegationKeyResponse{}, err
	}
	resp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return ServiceClientGetUserDelegationKeyResponse{}, err
	}
	if !runtime.HasStatusCode(resp, http.StatusOK) {
		return ServiceClientGetUserDelegationKeyResponse{}, runtime.NewResponseError(resp)
	}
	return client.getUserDelegationKeyHandleResponse(resp)
}

// getUserDelegationKeyCreateRequest creates the GetUserDelegationKey request.
func (client *ServiceClient) getUserDelegationKeyCreateRequest(ctx context.Context, keyInfo KeyInfo, options *ServiceClientGetUserDelegationKeyOptions) (*policy.Request, error) {
	req, err := runtime.NewRequest(ctx, http.MethodPost, client.endpoint)
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("restype", "service")
	reqQP.Set("comp", "userdelegationkey")
	if options != nil && options.Timeout != nil {
		reqQP.Set("timeout", strconv.FormatInt(int64(*options.Timeout), 10))
	}
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["x-ms-version"] = []string{"2020-10-02"}
	if options != nil && options.RequestID != nil {
		req.Raw().Header["x-ms-client-request-id"] = []string{*options.RequestID}
	}
	req.Raw().Header["Accept"] = []string{"application/xml"}
	return req, runtime.MarshalAsXML(req, keyInfo)
}

// getUserDelegationKeyHandleResponse handles the GetUserDelegationKey response.
func (client *ServiceClient) getUserDelegationKeyHandleResponse(resp *http.Response) (ServiceClientGetUserDelegationKeyResponse, error) {
	result := ServiceClientGetUserDelegationKeyResponse{}
	if val := resp.Header.Get("x-ms-client-request-id"); val != "" {
		result.ClientRequestID = &val
	}
	if val := resp.Header.Get("x-ms-request-id"); val != "" {
		result.RequestID = &val
	}
	if val := resp.Header.Get("x-ms-version"); val != "" {
		result.Version = &val
	}
	if val := resp.Header.Get("Date"); val != "" {
		date, err := time.Parse(time.RFC1123, val)
		if err != nil {
			return ServiceClientGetUserDelegationKeyResponse{}, err
		}
		result.Date = &date
	}
	if err := runtime.UnmarshalAsXML(resp, &result.UserDelegationKey); err != nil {
		return ServiceClientGetUserDelegationKeyResponse{}, err
	}
	return result, nil
}
