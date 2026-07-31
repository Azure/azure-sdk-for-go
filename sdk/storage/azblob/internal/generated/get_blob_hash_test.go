// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package generated

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

type trackingGetBlobHashBody struct {
	*strings.Reader
	closed bool
}

func (b *trackingGetBlobHashBody) Close() error {
	b.closed = true
	return nil
}

func getBlobHashHeaderValue(header http.Header, name string) string {
	for headerName, values := range header {
		if strings.EqualFold(headerName, name) && len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

func TestGetBlobHashCreateRequest(t *testing.T) {
	client := &BlockBlobClient{endpoint: "https://acct.blob.core.windows.net/container/blob?snapshot=snapshot-value&sig=secret"}
	etag := azcore.ETag(`"opaque-epoch"`)
	ifTags := `"tag" = 'value'`
	leaseID := "lease-id"
	requestID := "client-request"
	timeout := int32(30)
	algorithm := EncryptionAlgorithmTypeAES256
	encryptionKey := "base64-key"
	encryptionKeySHA256 := "base64-key-sha256"

	req, err := client.getBlobHashCreateRequest(
		context.Background(),
		"bytes=0-1023,4096-6143",
		&BlockBlobClientGetBlobHashOptions{RequestID: &requestID, Timeout: &timeout},
		&LeaseAccessConditions{LeaseID: &leaseID},
		&CPKInfo{
			EncryptionAlgorithm: &algorithm,
			EncryptionKey:       &encryptionKey,
			EncryptionKeySHA256: &encryptionKeySHA256,
		},
		&ModifiedAccessConditions{IfMatch: &etag, IfTags: &ifTags},
	)
	require.NoError(t, err)
	require.Equal(t, http.MethodGet, req.Raw().Method)
	require.Equal(t, "hash", req.Raw().URL.Query().Get("comp"))
	require.Equal(t, "snapshot-value", req.Raw().URL.Query().Get("snapshot"))
	require.Equal(t, "secret", req.Raw().URL.Query().Get("sig"))
	require.Equal(t, "30", req.Raw().URL.Query().Get("timeout"))
	require.Equal(t, "Sha256", getBlobHashHeaderValue(req.Raw().Header, "x-ms-hash-algorithm"))
	require.Equal(t, "bytes=0-1023,4096-6143", getBlobHashHeaderValue(req.Raw().Header, "x-ms-multi-range"))
	require.Equal(t, string(etag), getBlobHashHeaderValue(req.Raw().Header, "If-Match"))
	require.Equal(t, ifTags, getBlobHashHeaderValue(req.Raw().Header, "x-ms-if-tags"))
	require.Equal(t, leaseID, getBlobHashHeaderValue(req.Raw().Header, "x-ms-lease-id"))
	require.Equal(t, string(algorithm), getBlobHashHeaderValue(req.Raw().Header, "x-ms-encryption-algorithm"))
	require.Equal(t, encryptionKey, getBlobHashHeaderValue(req.Raw().Header, "x-ms-encryption-key"))
	require.Equal(t, encryptionKeySHA256, getBlobHashHeaderValue(req.Raw().Header, "x-ms-encryption-key-sha256"))
	require.Equal(t, requestID, getBlobHashHeaderValue(req.Raw().Header, "x-ms-client-request-id"))
	require.Equal(t, "2025-11-05", getBlobHashHeaderValue(req.Raw().Header, "x-ms-version"))
	require.True(t, req.Raw().Body == nil || req.Raw().Body == http.NoBody)
	require.Zero(t, req.Raw().ContentLength)
}

func TestGetBlobHashHandleResponse(t *testing.T) {
	firstHash := bytes.Repeat([]byte{0x11}, sha256.Size)
	secondHash := bytes.Repeat([]byte{0x22}, sha256.Size)
	body := `<RangeHashList>
		<FutureRootField>ignored</FutureRootField>
		<RangeHash><Offset>0</Offset><Length>1024</Length><Sha256>` + base64.StdEncoding.EncodeToString(firstHash) + `</Sha256><FutureRangeField>ignored</FutureRangeField></RangeHash>
		<RangeHash><Offset>4096</Offset><Length>2048</Length><Sha256>` + base64.StdEncoding.EncodeToString(secondHash) + `</Sha256></RangeHash>
	</RangeHashList>`
	headers := http.Header{}
	headers.Set("Content-Type", "application/xml")
	headers.Set("ETag", `"opaque-epoch"`)
	headers.Set("Last-Modified", "Mon, 02 Jan 2006 15:04:05 GMT")
	headers.Set("x-ms-blob-content-length", "8192")
	headers.Set("x-ms-hash-algorithm", "Sha256")
	headers.Set("x-ms-request-id", "service-request")
	headers.Set("x-ms-client-request-id", "client-request")
	headers.Set("x-ms-test-dedupe-sha256-cpu-time-us", "123")
	headers.Set("x-ms-version", "2025-11-05")

	response, err := (&BlockBlobClient{}).getBlobHashHandleResponse(&http.Response{
		StatusCode: http.StatusOK,
		Header:     headers,
		Body:       io.NopCloser(strings.NewReader(body)),
	})
	require.NoError(t, err)
	require.Len(t, response.RangeHashes, 2)
	require.Equal(t, int64(0), *response.RangeHashes[0].Offset)
	require.Equal(t, int64(1024), *response.RangeHashes[0].Length)
	require.Equal(t, firstHash, response.RangeHashes[0].Sha256)
	require.Equal(t, int64(4096), *response.RangeHashes[1].Offset)
	require.Equal(t, int64(2048), *response.RangeHashes[1].Length)
	require.Equal(t, secondHash, response.RangeHashes[1].Sha256)
	require.Equal(t, azcore.ETag(`"opaque-epoch"`), *response.ETag)
	require.Equal(t, int64(8192), *response.BlobContentLength)
	require.Equal(t, "Sha256", *response.HashAlgorithm)
	require.Equal(t, "service-request", *response.RequestID)
	require.Equal(t, "client-request", *response.ClientRequestID)
	require.Equal(t, int64(123), *response.SHA256CPUTimeUS)
	require.Equal(t, "2025-11-05", *response.Version)
	require.True(t, response.LastModified.Equal(time.Date(2006, time.January, 2, 15, 4, 5, 0, time.UTC)))
}

func TestGetBlobHashHandleResponseToleratesInvalidCPUTime(t *testing.T) {
	hash := bytes.Repeat([]byte{0x11}, sha256.Size)
	body := `<RangeHashList><RangeHash><Offset>0</Offset><Length>1</Length><Sha256>` +
		base64.StdEncoding.EncodeToString(hash) +
		`</Sha256></RangeHash></RangeHashList>`

	for _, value := range []string{"", "-1", "not-a-number", "9223372036854775808"} {
		t.Run(value, func(t *testing.T) {
			headers := http.Header{"Content-Type": []string{"application/xml"}}
			if value != "" {
				headers.Set("x-ms-test-dedupe-sha256-cpu-time-us", value)
			}
			response, err := (&BlockBlobClient{}).getBlobHashHandleResponse(&http.Response{
				StatusCode: http.StatusOK,
				Header:     headers,
				Body:       io.NopCloser(strings.NewReader(body)),
			})
			require.NoError(t, err)
			require.Nil(t, response.SHA256CPUTimeUS)
		})
	}

	headers := http.Header{"Content-Type": []string{"application/xml"}}
	headers.Set("x-ms-test-dedupe-sha256-cpu-time-us", "0")
	response, err := (&BlockBlobClient{}).getBlobHashHandleResponse(&http.Response{
		StatusCode: http.StatusOK,
		Header:     headers,
		Body:       io.NopCloser(strings.NewReader(body)),
	})
	require.NoError(t, err)
	require.Equal(t, int64(0), *response.SHA256CPUTimeUS)
}

func TestGetBlobHashHandleResponseRejectsInvalidBase64(t *testing.T) {
	for _, value := range []string{"not-base64", "&quot;"} {
		t.Run(value, func(t *testing.T) {
			body := `<RangeHashList><RangeHash><Offset>0</Offset><Length>1</Length><Sha256>` + value + `</Sha256></RangeHash></RangeHashList>`
			_, err := (&BlockBlobClient{}).getBlobHashHandleResponse(&http.Response{
				StatusCode: http.StatusOK,
				Header:     http.Header{"Content-Type": []string{"application/xml"}},
				Body:       io.NopCloser(strings.NewReader(body)),
			})
			require.Error(t, err)
		})
	}
}

func TestGetBlobHashHandleResponseDrainsBodyOnHeaderError(t *testing.T) {
	body := &trackingGetBlobHashBody{Reader: strings.NewReader("response payload")}
	headers := http.Header{}
	headers.Set("x-ms-blob-content-length", "not-an-integer")

	_, err := (&BlockBlobClient{}).getBlobHashHandleResponse(&http.Response{
		StatusCode: http.StatusOK,
		Header:     headers,
		Body:       body,
	})

	require.Error(t, err)
	require.True(t, body.closed)
	require.Zero(t, body.Len())
}
