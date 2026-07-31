//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package blockblob_test

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blockblob"
	"github.com/stretchr/testify/require"
)

const maxBlobHashBytes = int64(4000 * 1024 * 1024)

type getBlobHashTransportResponse struct {
	statusCode int
	headers    http.Header
	body       string
	err        error
}

type getBlobHashTransport struct {
	mu        sync.Mutex
	requests  []*http.Request
	responses []getBlobHashTransportResponse
}

func (t *getBlobHashTransport) Do(req *http.Request) (*http.Response, error) {
	if err := req.Context().Err(); err != nil {
		return nil, err
	}

	t.mu.Lock()
	index := len(t.requests)
	cloned := req.Clone(req.Context())
	cloned.Header = req.Header.Clone()
	t.requests = append(t.requests, cloned)
	response := t.responses[len(t.responses)-1]
	if index < len(t.responses) {
		response = t.responses[index]
	}
	t.mu.Unlock()

	if response.err != nil {
		return nil, response.err
	}
	statusCode := response.statusCode
	if statusCode == 0 {
		statusCode = http.StatusOK
	}
	return &http.Response{
		Request:    req,
		Status:     fmt.Sprintf("%d %s", statusCode, http.StatusText(statusCode)),
		StatusCode: statusCode,
		Header:     response.headers.Clone(),
		Body:       io.NopCloser(strings.NewReader(response.body)),
	}, nil
}

func (t *getBlobHashTransport) requestCount() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return len(t.requests)
}

func (t *getBlobHashTransport) request(index int) *http.Request {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.requests[index]
}

type getBlobHashRequestIDPolicy struct{}

func (getBlobHashRequestIDPolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Set("x-ms-client-request-id", "azblob-test-request-id")
	return req.Next()
}

func getBlobHashHeaderValue(header http.Header, name string) string {
	for headerName, values := range header {
		if strings.EqualFold(headerName, name) && len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

func newGetBlobHashClient(t *testing.T, url string, transport *getBlobHashTransport, configure func(*blockblob.ClientOptions)) *blockblob.Client {
	t.Helper()
	options := &blockblob.ClientOptions{ClientOptions: policy.ClientOptions{
		Transport:       transport,
		Telemetry:       policy.TelemetryOptions{ApplicationID: "testApp/1.0.0-preview.2"},
		PerCallPolicies: []policy.Policy{getBlobHashRequestIDPolicy{}},
	}}
	if configure != nil {
		configure(options)
	}
	client, err := blockblob.NewClientWithNoCredential(url, options)
	require.NoError(t, err)
	return client
}

func blobHashOptions() *blockblob.GetBlobHashOptions {
	etag := azcore.ETag(`"opaque-epoch"`)
	return &blockblob.GetBlobHashOptions{
		AccessConditions: &blob.AccessConditions{
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &etag},
		},
	}
}

func blobHashForRange(rnge blockblob.BlobHashRange) []byte {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%d:%d", rnge.Offset, rnge.Count)))
	return hash[:]
}

func rangeHashXML(ranges []blockblob.BlobHashRange, unknownFields bool) string {
	var body strings.Builder
	body.WriteString(`<?xml version="1.0" encoding="utf-8"?><RangeHashList>`)
	if unknownFields {
		body.WriteString(`<FutureRootField>ignored</FutureRootField>`)
	}
	for _, rnge := range ranges {
		fmt.Fprintf(&body, `<RangeHash><Offset>%d</Offset><Length>%d</Length><Sha256>%s</Sha256>`,
			rnge.Offset, rnge.Count, base64.StdEncoding.EncodeToString(blobHashForRange(rnge)))
		if unknownFields {
			body.WriteString(`<FutureRangeField>ignored</FutureRangeField>`)
		}
		body.WriteString(`</RangeHash>`)
	}
	body.WriteString(`</RangeHashList>`)
	return body.String()
}

func blobHashResponse(ranges []blockblob.BlobHashRange) getBlobHashTransportResponse {
	headers := http.Header{}
	headers.Set("Content-Type", "application/xml")
	headers.Set("ETag", `"opaque-epoch"`)
	headers.Set("Last-Modified", "Mon, 02 Jan 2006 15:04:05 GMT")
	headers.Set("x-ms-blob-content-length", "8192")
	headers.Set("x-ms-hash-algorithm", "Sha256")
	headers.Set("x-ms-request-id", "service-request")
	headers.Set("x-ms-client-request-id", "client-request")
	headers.Set("x-ms-test-dedupe-sha256-cpu-time-us", "321")
	headers.Set("x-ms-version", "2025-11-05")
	return getBlobHashTransportResponse{headers: headers, body: rangeHashXML(ranges, false)}
}

func TestGetBlobHash(t *testing.T) {
	requested := []blockblob.BlobHashRange{{Offset: 0, Count: 1024}, {Offset: 4096, Count: 2048}}
	responseOrder := []blockblob.BlobHashRange{requested[1], requested[0]}
	transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{blobHashResponse(responseOrder)}}
	client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob?sig=secret", transport, nil)

	response, err := client.GetBlobHash(context.Background(), requested, blobHashOptions())
	require.NoError(t, err)
	require.Len(t, response.RangeHashes, 2)
	require.Equal(t, responseOrder[0].Offset, response.RangeHashes[0].Offset)
	require.Equal(t, responseOrder[0].Count, response.RangeHashes[0].Count)
	require.Equal(t, blobHashForRange(responseOrder[0]), response.RangeHashes[0].SHA256)
	require.Equal(t, azcore.ETag(`"opaque-epoch"`), *response.ETag)
	require.Equal(t, int64(8192), *response.BlobContentLength)
	require.Equal(t, "Sha256", *response.HashAlgorithm)
	require.Equal(t, "service-request", *response.RequestID)
	require.Equal(t, "client-request", *response.ClientRequestID)
	require.Equal(t, int64(321), *response.SHA256CPUTimeUS)
	require.Equal(t, "2025-11-05", *response.Version)

	require.Equal(t, 1, transport.requestCount())
	req := transport.request(0)
	require.Equal(t, http.MethodGet, req.Method)
	require.Equal(t, "hash", req.URL.Query().Get("comp"))
	require.Equal(t, "secret", req.URL.Query().Get("sig"))
	require.Equal(t, "bytes=0-1023,4096-6143", getBlobHashHeaderValue(req.Header, "x-ms-multi-range"))
	require.Equal(t, "Sha256", getBlobHashHeaderValue(req.Header, "x-ms-hash-algorithm"))
	require.Equal(t, `"opaque-epoch"`, getBlobHashHeaderValue(req.Header, "If-Match"))
	require.Equal(t, "2025-11-05", getBlobHashHeaderValue(req.Header, "x-ms-version"))
	require.Equal(t, "application/xml", getBlobHashHeaderValue(req.Header, "Accept"))
	require.True(t, req.Body == nil || req.Body == http.NoBody)
	require.Zero(t, req.ContentLength)
}

func TestGetBlobHashConditions(t *testing.T) {
	ranges := []blockblob.BlobHashRange{{Offset: 0, Count: 1}}
	transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{blobHashResponse(ranges)}}
	client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob", transport, nil)

	leaseID := "lease-id"
	ifTags := `"tag" = 'value'`
	etag := azcore.ETag(`"etag"`)
	algorithm := blob.EncryptionAlgorithmTypeAES256
	encryptionKey := "base64-key"
	encryptionKeySHA256 := "base64-key-sha256"
	_, err := client.GetBlobHash(context.Background(), ranges, &blockblob.GetBlobHashOptions{
		AccessConditions: &blob.AccessConditions{
			LeaseAccessConditions: &blob.LeaseAccessConditions{LeaseID: &leaseID},
			ModifiedAccessConditions: &blob.ModifiedAccessConditions{
				IfMatch: &etag,
				IfTags:  &ifTags,
			},
		},
		CPKInfo: &blob.CPKInfo{
			EncryptionAlgorithm: &algorithm,
			EncryptionKey:       &encryptionKey,
			EncryptionKeySHA256: &encryptionKeySHA256,
		},
	})
	require.NoError(t, err)

	req := transport.request(0)
	require.Equal(t, leaseID, getBlobHashHeaderValue(req.Header, "x-ms-lease-id"))
	require.Equal(t, ifTags, getBlobHashHeaderValue(req.Header, "x-ms-if-tags"))
	require.Equal(t, string(etag), getBlobHashHeaderValue(req.Header, "If-Match"))
	require.Equal(t, string(algorithm), getBlobHashHeaderValue(req.Header, "x-ms-encryption-algorithm"))
	require.Equal(t, encryptionKey, getBlobHashHeaderValue(req.Header, "x-ms-encryption-key"))
	require.Equal(t, encryptionKeySHA256, getBlobHashHeaderValue(req.Header, "x-ms-encryption-key-sha256"))
}

func exactHeaderLimitRanges() []blockblob.BlobHashRange {
	ranges := []blockblob.BlobHashRange{
		{Offset: 99999999999995, Count: 1},
		{Offset: 99999999999997, Count: 1},
		{Offset: 99999999999999, Count: 2},
	}
	for i := 0; i < 253; i++ {
		ranges = append(ranges, blockblob.BlobHashRange{Offset: 100000000000002 + int64(i*2), Count: 1})
	}
	return ranges
}

func TestGetBlobHashBoundaryLimits(t *testing.T) {
	tests := []struct {
		name            string
		ranges          []blockblob.BlobHashRange
		wantHeaderBytes int
	}{
		{name: "256 ranges and 8192-byte header", ranges: exactHeaderLimitRanges(), wantHeaderBytes: 8192},
		{name: "individual 4000 MiB", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: maxBlobHashBytes}}},
		{name: "aggregate 4000 MiB", ranges: []blockblob.BlobHashRange{
			{Offset: 0, Count: maxBlobHashBytes / 2},
			{Offset: maxBlobHashBytes / 2, Count: maxBlobHashBytes / 2},
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{blobHashResponse(test.ranges)}}
			client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob", transport, nil)

			response, err := client.GetBlobHash(context.Background(), test.ranges, blobHashOptions())
			require.NoError(t, err)
			require.Len(t, response.RangeHashes, len(test.ranges))
			if test.wantHeaderBytes != 0 {
				require.Len(t, getBlobHashHeaderValue(transport.request(0).Header, "x-ms-multi-range"), test.wantHeaderBytes)
			}
		})
	}
}

func TestGetBlobHashValidation(t *testing.T) {
	tooManyRanges := make([]blockblob.BlobHashRange, 257)
	for i := range tooManyRanges {
		tooManyRanges[i] = blockblob.BlobHashRange{Offset: int64(i * 2), Count: 1}
	}
	oversizedHeader := exactHeaderLimitRanges()
	oversizedHeader[len(oversizedHeader)-1].Offset = 1000000000000000

	tests := []struct {
		name    string
		ranges  []blockblob.BlobHashRange
		options func() *blockblob.GetBlobHashOptions
	}{
		{name: "empty ranges", ranges: nil, options: blobHashOptions},
		{name: "negative offset", ranges: []blockblob.BlobHashRange{{Offset: -1, Count: 1}}, options: blobHashOptions},
		{name: "zero count", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 0}}, options: blobHashOptions},
		{name: "negative count", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: -1}}, options: blobHashOptions},
		{name: "range overflow", ranges: []blockblob.BlobHashRange{{Offset: math.MaxInt64, Count: 2}}, options: blobHashOptions},
		{name: "unsorted ranges", ranges: []blockblob.BlobHashRange{{Offset: 10, Count: 1}, {Offset: 0, Count: 1}}, options: blobHashOptions},
		{name: "overlapping ranges", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 10}, {Offset: 9, Count: 2}}, options: blobHashOptions},
		{name: "257 ranges", ranges: tooManyRanges, options: blobHashOptions},
		{name: "individual data limit", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: maxBlobHashBytes + 1}}, options: blobHashOptions},
		{name: "aggregate data limit", ranges: []blockblob.BlobHashRange{
			{Offset: 0, Count: maxBlobHashBytes / 2},
			{Offset: maxBlobHashBytes / 2, Count: maxBlobHashBytes/2 + 1},
		}, options: blobHashOptions},
		{name: "header size limit", ranges: oversizedHeader, options: blobHashOptions},
		{name: "nil options", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 1}}},
		{name: "nil access conditions", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, options: func() *blockblob.GetBlobHashOptions {
			return &blockblob.GetBlobHashOptions{}
		}},
		{name: "nil modified conditions", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, options: func() *blockblob.GetBlobHashOptions {
			return &blockblob.GetBlobHashOptions{AccessConditions: &blob.AccessConditions{}}
		}},
		{name: "missing If-Match", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, options: func() *blockblob.GetBlobHashOptions {
			return &blockblob.GetBlobHashOptions{AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{},
			}}
		}},
		{name: "empty If-Match", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, options: func() *blockblob.GetBlobHashOptions {
			etag := azcore.ETag("")
			return &blockblob.GetBlobHashOptions{AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &etag},
			}}
		}},
		{name: "wildcard If-Match", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, options: func() *blockblob.GetBlobHashOptions {
			etag := azcore.ETagAny
			return &blockblob.GetBlobHashOptions{AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &etag},
			}}
		}},
		{name: "whitespace wildcard If-Match", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, options: func() *blockblob.GetBlobHashOptions {
			etag := azcore.ETag(" * ")
			return &blockblob.GetBlobHashOptions{AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &etag},
			}}
		}},
		{name: "weak If-Match", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, options: func() *blockblob.GetBlobHashOptions {
			etag := azcore.ETag(`W/"etag"`)
			return &blockblob.GetBlobHashOptions{AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &etag},
			}}
		}},
		{name: "If-Match list", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, options: func() *blockblob.GetBlobHashOptions {
			etag := azcore.ETag(`"one","two"`)
			return &blockblob.GetBlobHashOptions{AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &etag},
			}}
		}},
		{name: "unquoted If-Match", ranges: []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, options: func() *blockblob.GetBlobHashOptions {
			etag := azcore.ETag("etag")
			return &blockblob.GetBlobHashOptions{AccessConditions: &blob.AccessConditions{
				ModifiedAccessConditions: &blob.ModifiedAccessConditions{IfMatch: &etag},
			}}
		}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{{}}}
			client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob", transport, nil)
			var options *blockblob.GetBlobHashOptions
			if test.options != nil {
				options = test.options()
			}

			_, err := client.GetBlobHash(context.Background(), test.ranges, options)
			require.Error(t, err)
			require.Zero(t, transport.requestCount())
		})
	}
}

func TestGetBlobHashResponseValidation(t *testing.T) {
	ranges := []blockblob.BlobHashRange{{Offset: 0, Count: 10}, {Offset: 20, Count: 10}}
	validSHA := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x11}, sha256.Size))
	shortSHA := base64.StdEncoding.EncodeToString(bytes.Repeat([]byte{0x11}, sha256.Size-1))

	tests := []struct {
		name string
		body string
	}{
		{name: "invalid base64", body: `<RangeHashList><RangeHash><Offset>0</Offset><Length>10</Length><Sha256>not-base64</Sha256></RangeHash><RangeHash><Offset>20</Offset><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash></RangeHashList>`},
		{name: "wrong SHA256 length", body: `<RangeHashList><RangeHash><Offset>0</Offset><Length>10</Length><Sha256>` + shortSHA + `</Sha256></RangeHash><RangeHash><Offset>20</Offset><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash></RangeHashList>`},
		{name: "missing SHA256", body: `<RangeHashList><RangeHash><Offset>0</Offset><Length>10</Length></RangeHash><RangeHash><Offset>20</Offset><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash></RangeHashList>`},
		{name: "missing offset", body: `<RangeHashList><RangeHash><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash><RangeHash><Offset>20</Offset><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash></RangeHashList>`},
		{name: "missing length", body: `<RangeHashList><RangeHash><Offset>0</Offset><Sha256>` + validSHA + `</Sha256></RangeHash><RangeHash><Offset>20</Offset><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash></RangeHashList>`},
		{name: "response count mismatch", body: `<RangeHashList><RangeHash><Offset>0</Offset><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash></RangeHashList>`},
		{name: "response range mismatch", body: `<RangeHashList><RangeHash><Offset>0</Offset><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash><RangeHash><Offset>21</Offset><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash></RangeHashList>`},
		{name: "duplicate response range", body: `<RangeHashList><RangeHash><Offset>0</Offset><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash><RangeHash><Offset>0</Offset><Length>10</Length><Sha256>` + validSHA + `</Sha256></RangeHash></RangeHashList>`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{{headers: http.Header{"Content-Type": []string{"application/xml"}}, body: test.body}}}
			client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob", transport, nil)

			_, err := client.GetBlobHash(context.Background(), ranges, blobHashOptions())
			require.Error(t, err)
		})
	}
}

func TestGetBlobHashIgnoresUnknownXMLFields(t *testing.T) {
	ranges := []blockblob.BlobHashRange{{Offset: 0, Count: 10}}
	response := blobHashResponse(ranges)
	response.body = rangeHashXML(ranges, true)
	transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{response}}
	client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob", transport, nil)

	result, err := client.GetBlobHash(context.Background(), ranges, blobHashOptions())
	require.NoError(t, err)
	require.Len(t, result.RangeHashes, 1)
	require.Equal(t, blobHashForRange(ranges[0]), result.RangeHashes[0].SHA256)
}

func TestGetBlobHashSnapshotAndVersionQueries(t *testing.T) {
	ranges := []blockblob.BlobHashRange{{Offset: 0, Count: 1}}
	transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{blobHashResponse(ranges), blobHashResponse(ranges)}}
	client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob?sig=secret", transport, nil)

	snapshotClient, err := client.WithSnapshot("snapshot-value")
	require.NoError(t, err)
	_, err = snapshotClient.GetBlobHash(context.Background(), ranges, blobHashOptions())
	require.NoError(t, err)

	versionClient, err := client.WithVersionID("version-value")
	require.NoError(t, err)
	_, err = versionClient.GetBlobHash(context.Background(), ranges, blobHashOptions())
	require.NoError(t, err)

	require.Equal(t, "snapshot-value", transport.request(0).URL.Query().Get("snapshot"))
	require.Equal(t, "secret", transport.request(0).URL.Query().Get("sig"))
	require.Equal(t, "version-value", transport.request(1).URL.Query().Get("versionid"))
	require.Equal(t, "secret", transport.request(1).URL.Query().Get("sig"))
}

func TestGetBlobHashErrors(t *testing.T) {
	tests := []struct {
		statusCode int
		errorCode  string
	}{
		{statusCode: http.StatusBadRequest, errorCode: "InvalidHeaderValue"},
		{statusCode: http.StatusNotFound, errorCode: "BlobNotFound"},
		{statusCode: http.StatusConflict, errorCode: "FeatureVersionMismatch"},
		{statusCode: http.StatusPreconditionFailed, errorCode: "ConditionNotMet"},
		{statusCode: http.StatusRequestedRangeNotSatisfiable, errorCode: "InvalidRange"},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d", test.statusCode), func(t *testing.T) {
			headers := http.Header{}
			headers.Set("x-ms-error-code", test.errorCode)
			transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{{
				statusCode: test.statusCode,
				headers:    headers,
				body:       `<Error><Code>` + test.errorCode + `</Code><Message>request failed</Message></Error>`,
			}}}
			client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob", transport, nil)

			_, err := client.GetBlobHash(context.Background(), []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, blobHashOptions())
			require.Error(t, err)
			var responseError *azcore.ResponseError
			require.ErrorAs(t, err, &responseError)
			require.Equal(t, test.statusCode, responseError.StatusCode)
			require.Equal(t, test.errorCode, responseError.ErrorCode)
		})
	}
}

func TestGetBlobHashETagMutation(t *testing.T) {
	headers := http.Header{}
	headers.Set("x-ms-error-code", "ConditionNotMet")
	transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{{
		statusCode: http.StatusPreconditionFailed,
		headers:    headers,
		body:       `<Error><Code>ConditionNotMet</Code><Message>The blob changed.</Message></Error>`,
	}}}
	client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob", transport, nil)

	_, err := client.GetBlobHash(context.Background(), []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, blobHashOptions())
	var responseError *azcore.ResponseError
	require.ErrorAs(t, err, &responseError)
	require.Equal(t, http.StatusPreconditionFailed, responseError.StatusCode)
	require.Equal(t, "ConditionNotMet", responseError.ErrorCode)
}

func TestGetBlobHashCancellation(t *testing.T) {
	transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{{}}}
	client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob", transport, nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	_, err := client.GetBlobHash(ctx, []blockblob.BlobHashRange{{Offset: 0, Count: 1}}, blobHashOptions())
	require.Error(t, err)
	require.True(t, errors.Is(err, context.Canceled))
	require.LessOrEqual(t, transport.requestCount(), 1)
}

func TestGetBlobHashRetry(t *testing.T) {
	ranges := []blockblob.BlobHashRange{{Offset: 0, Count: 1}}
	transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{
		{statusCode: http.StatusInternalServerError, body: `<Error><Code>InternalError</Code></Error>`},
		blobHashResponse(ranges),
	}}
	client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob", transport, func(options *blockblob.ClientOptions) {
		options.Retry = policy.RetryOptions{MaxRetries: 1, RetryDelay: time.Millisecond, MaxRetryDelay: time.Millisecond}
	})

	response, err := client.GetBlobHash(context.Background(), ranges, blobHashOptions())
	require.NoError(t, err)
	require.Len(t, response.RangeHashes, 1)
	require.Equal(t, 2, transport.requestCount())
	require.Equal(t,
		getBlobHashHeaderValue(transport.request(0).Header, "x-ms-multi-range"),
		getBlobHashHeaderValue(transport.request(1).Header, "x-ms-multi-range"))
}

func TestGetBlockListDoesNotUseGetBlobHashWireContract(t *testing.T) {
	const body = `<?xml version="1.0" encoding="utf-8"?><BlockList><CommittedBlocks/></BlockList>`
	transport := &getBlobHashTransport{responses: []getBlobHashTransportResponse{{body: body}}}
	client := newGetBlobHashClient(t, "https://fake.blob.core.windows.net/container/blob", transport, nil)

	_, err := client.GetBlockList(context.Background(), blockblob.BlockListTypeCommitted, &blockblob.GetBlockListOptions{
		Include: []blockblob.BlockListIncludeItem{blockblob.BlockListIncludeItemSha256},
	})
	require.NoError(t, err)
	req := transport.request(0)
	require.Equal(t, "blocklist", req.URL.Query().Get("comp"))
	require.Equal(t, "sha256", req.URL.Query().Get("include"))
	require.Empty(t, getBlobHashHeaderValue(req.Header, "x-ms-hash-algorithm"))
	require.Empty(t, getBlobHashHeaderValue(req.Header, "x-ms-multi-range"))
	require.Empty(t, getBlobHashHeaderValue(req.Header, "x-ms-sha256-block-indices"))
}
