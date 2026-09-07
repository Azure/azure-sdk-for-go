// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package path

import (
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/stretchr/testify/require"
)

// FormatGetPropertiesResponse must not panic when the raw HTTP response was not
// captured (rawResponse is nil). This happens for clients whose pipeline does not
// include the IncludeBlobResponsePolicy. See issue #25490.
func TestFormatGetPropertiesResponseNilRawResponse(t *testing.T) {
	require.NotPanics(t, func() {
		resp := FormatGetPropertiesResponse(&blob.GetPropertiesResponse{}, nil)
		// The datalake-specific fields come from the raw response, so they stay nil.
		require.Nil(t, resp.Owner)
		require.Nil(t, resp.Group)
		require.Nil(t, resp.Permissions)
		require.Nil(t, resp.AccessControlList)
		require.Nil(t, resp.ResourceType)
	})
}

// When the raw HTTP response is present, the datalake-specific headers are read.
func TestFormatGetPropertiesResponseWithRawResponse(t *testing.T) {
	rawResponse := &http.Response{Header: http.Header{}}
	rawResponse.Header.Set("x-ms-owner", "owner1")
	rawResponse.Header.Set("x-ms-group", "group1")
	rawResponse.Header.Set("x-ms-permissions", "rwxr-x---")
	rawResponse.Header.Set("x-ms-resource-type", "file")

	resp := FormatGetPropertiesResponse(&blob.GetPropertiesResponse{}, rawResponse)
	require.Equal(t, "owner1", *resp.Owner)
	require.Equal(t, "group1", *resp.Group)
	require.Equal(t, "rwxr-x---", *resp.Permissions)
	require.Equal(t, "file", *resp.ResourceType)
}
