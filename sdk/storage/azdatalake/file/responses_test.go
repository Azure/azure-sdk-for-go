// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/stretchr/testify/require"
)

// FormatDownloadStreamResponse must not panic when the raw HTTP response was not
// captured (rawResponse is nil). This mirrors the guard tested for
// FormatGetPropertiesResponse. See issue #25490.
func TestFormatDownloadStreamResponseNilRawResponse(t *testing.T) {
	require.NotPanics(t, func() {
		resp := FormatDownloadStreamResponse(&blob.DownloadStreamResponse{}, nil)
		// The datalake-specific fields come from the raw response, so they stay nil.
		require.Nil(t, resp.AccessControlList)
		require.Nil(t, resp.EncryptionContext)
	})
}

// When the raw HTTP response is present, the datalake-specific headers are read.
func TestFormatDownloadStreamResponseWithRawResponse(t *testing.T) {
	rawResponse := &http.Response{Header: http.Header{}}
	rawResponse.Header.Set("x-ms-acl", "user::rwx,group::r-x,other::---")
	rawResponse.Header.Set("x-ms-encryption-context", "context1")

	resp := FormatDownloadStreamResponse(&blob.DownloadStreamResponse{}, rawResponse)
	require.NotNil(t, resp.AccessControlList)
	require.Equal(t, "user::rwx,group::r-x,other::---", *resp.AccessControlList)
	require.NotNil(t, resp.EncryptionContext)
	require.Equal(t, "context1", *resp.EncryptionContext)
}
