//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/stretchr/testify/require"
)

func TestFormatDownloadStreamResponseWithNilRawResponse(t *testing.T) {
	// Test case that reproduces the null pointer exception from the issue
	var resp *blob.DownloadStreamResponse = nil
	var rawResponse *http.Response = nil

	// This should not panic even when rawResponse is nil
	result := FormatDownloadStreamResponse(resp, rawResponse)

	// Verify the result is properly initialized
	require.NotNil(t, result)
	require.Nil(t, result.EncryptionContext)
	require.Nil(t, result.AccessControlList)
}

func TestFormatDownloadStreamResponseWithNilBlobResponse(t *testing.T) {
	// Test case where blob response is nil but rawResponse is not
	var resp *blob.DownloadStreamResponse = nil
	rawResponse := &http.Response{
		Header: make(http.Header),
	}
	rawResponse.Header.Set("x-ms-encryption-context", "test-context")
	rawResponse.Header.Set("x-ms-acl", "test-acl")

	// This should not panic and should handle headers from rawResponse
	result := FormatDownloadStreamResponse(resp, rawResponse)

	// Verify the result includes data from rawResponse headers
	require.NotNil(t, result)
	require.NotNil(t, result.EncryptionContext)
	require.Equal(t, "test-context", *result.EncryptionContext)
	require.NotNil(t, result.AccessControlList)
	require.Equal(t, "test-acl", *result.AccessControlList)
}

func TestFormatDownloadStreamResponseWithValidInputs(t *testing.T) {
	// Test case with valid inputs to ensure normal functionality still works
	resp := &blob.DownloadStreamResponse{}
	rawResponse := &http.Response{
		Header: make(http.Header),
	}
	rawResponse.Header.Set("x-ms-encryption-context", "test-context")
	rawResponse.Header.Set("x-ms-acl", "test-acl")

	result := FormatDownloadStreamResponse(resp, rawResponse)

	// Verify the result includes data from both responses
	require.NotNil(t, result)
	require.NotNil(t, result.EncryptionContext)
	require.Equal(t, "test-context", *result.EncryptionContext)
	require.NotNil(t, result.AccessControlList)
	require.Equal(t, "test-acl", *result.AccessControlList)
}
