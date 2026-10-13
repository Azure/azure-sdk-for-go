// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"encoding/xml"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestShareGetPropertiesHandleResponseCreationTime(t *testing.T) {
	client := &ShareClient{}
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{}}
	resp.Header.Set("x-ms-share-creation-time", "Tue, 13 Oct 2026 10:20:30 GMT")

	result, err := client.getPropertiesHandleResponse(resp, http.StatusOK)
	require.NoError(t, err)
	require.NotNil(t, result.CreationTime)
	require.True(t, time.Date(2026, time.October, 13, 10, 20, 30, 0, time.UTC).Equal(*result.CreationTime))
}

func TestShareGetPropertiesHandleResponseCreationTimeMissing(t *testing.T) {
	client := &ShareClient{}
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{}}

	result, err := client.getPropertiesHandleResponse(resp, http.StatusOK)
	require.NoError(t, err)
	require.Nil(t, result.CreationTime)
}

func TestShareGetPropertiesHandleResponseCreationTimeInvalid(t *testing.T) {
	client := &ShareClient{}
	resp := &http.Response{StatusCode: http.StatusOK, Header: http.Header{}}
	resp.Header.Set("x-ms-share-creation-time", "not a date")

	_, err := client.getPropertiesHandleResponse(resp, http.StatusOK)
	require.Error(t, err)
}

func TestSharePropertiesUnmarshalXMLCreationTime(t *testing.T) {
	const body = `<Properties>
  <Last-Modified>Tue, 13 Oct 2026 10:20:31 GMT</Last-Modified>
  <Etag>"0x8DE0000000000000"</Etag>
  <Quota>100</Quota>
  <Creation-Time>Tue, 13 Oct 2026 10:20:30 GMT</Creation-Time>
</Properties>`

	var props ShareProperties
	require.NoError(t, xml.Unmarshal([]byte(body), &props))
	require.NotNil(t, props.ETag)
	require.NotNil(t, props.CreationTime)
	require.True(t, time.Date(2026, time.October, 13, 10, 20, 30, 0, time.UTC).Equal(*props.CreationTime))
}

func TestSharePropertiesUnmarshalXMLCreationTimeMissing(t *testing.T) {
	const body = `<Properties>
  <Last-Modified>Tue, 13 Oct 2026 10:20:31 GMT</Last-Modified>
  <Etag>"0x8DE0000000000000"</Etag>
  <Quota>100</Quota>
</Properties>`

	var props ShareProperties
	require.NoError(t, xml.Unmarshal([]byte(body), &props))
	require.NotNil(t, props.ETag)
	require.Nil(t, props.CreationTime)
}

func TestListSharesResponseUnmarshalXMLCreationTime(t *testing.T) {
	const body = `<?xml version="1.0" encoding="utf-8"?>
<EnumerationResults ServiceEndpoint="https://account.file.core.windows.net/">
  <Shares>
    <Share>
      <Name>withcreationtime</Name>
      <Properties>
        <Last-Modified>Tue, 13 Oct 2026 10:20:31 GMT</Last-Modified>
        <Etag>"0x8DE0000000000001"</Etag>
        <Quota>100</Quota>
        <Creation-Time>Tue, 13 Oct 2026 10:20:30 GMT</Creation-Time>
      </Properties>
    </Share>
    <Share>
      <Name>withoutcreationtime</Name>
      <Properties>
        <Last-Modified>Tue, 13 Oct 2026 10:20:31 GMT</Last-Modified>
        <Etag>"0x8DE0000000000002"</Etag>
        <Quota>100</Quota>
      </Properties>
    </Share>
  </Shares>
  <NextMarker />
</EnumerationResults>`

	var resp ListSharesResponse
	require.NoError(t, xml.Unmarshal([]byte(body), &resp))
	require.Len(t, resp.Shares, 2)

	require.Equal(t, "withcreationtime", *resp.Shares[0].Name)
	require.NotNil(t, resp.Shares[0].Properties.CreationTime)
	require.True(t, time.Date(2026, time.October, 13, 10, 20, 30, 0, time.UTC).Equal(*resp.Shares[0].Properties.CreationTime))

	require.Equal(t, "withoutcreationtime", *resp.Shares[1].Name)
	require.Nil(t, resp.Shares[1].Properties.CreationTime)
}
