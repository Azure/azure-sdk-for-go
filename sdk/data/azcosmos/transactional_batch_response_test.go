// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/stretchr/testify/require"
)

func TestTransactionalBatchResponsePreservesMetadataAndOwnership(t *testing.T) {
	diagnostics := &Diagnostics{AttemptCount: 2, TotalRequestCharge: 6.25}
	metadata := ItemResponse{
		Response: Response{
			StatusCode: http.StatusOK, SubStatus: 100, RequestCharge: 6.25, ActivityID: "batch",
			Diagnostics: diagnostics, AttemptCount: 2,
		},
		SessionToken: "0:-1#42", ETag: `"batch-etag"`,
	}
	body := []byte(`[
		{"statusCode":201,"subStatusCode":17,"requestCharge":1.125,"eTag":"\"created\"",
		 "resourceBody":{ "n":9007199254740993, "s":"\u0041" },"retryAfterMilliseconds":25},
		{"statusCode":200,"resourceBody":{"id":"read"},"requestCharge":2.25},
		{"statusCode":204},
		{"statusCode":299}
	]`)
	response, err := decodeTransactionalBatchResponse(metadata, body, 125*time.Millisecond, true, 4)
	require.NoError(t, err)
	require.Equal(t, metadata.Response, response.Response)
	require.Same(t, diagnostics, response.Diagnostics)
	require.Equal(t, metadata.SessionToken, response.SessionToken)
	require.Equal(t, metadata.ETag, response.ETag)
	require.Equal(t, 125*time.Millisecond, response.RetryAfter)
	require.Equal(t, []TransactionalBatchResult{
		{StatusCode: 201, SubStatus: 17, RequestCharge: 1.125, ETag: `"created"`,
			ResourceBody: []byte(`{ "n":9007199254740993, "s":"\u0041" }`), RetryAfter: 25 * time.Millisecond},
		{StatusCode: 200, RequestCharge: 2.25, ResourceBody: []byte(`{"id":"read"}`)},
		{StatusCode: 204}, {StatusCode: 299},
	}, response.OperationResults)
	ownedBody := string(response.Body)
	for i := range body {
		body[i] = 'x'
	}
	require.Equal(t, ownedBody, string(response.Body))
	for i := range response.Body {
		response.Body[i] = 'x'
	}
	require.Equal(t, `{ "n":9007199254740993, "s":"\u0041" }`, string(response.OperationResults[0].ResourceBody))
}

func TestTransactionalBatchResponsePreservesFailureResults(t *testing.T) {
	for _, failure := range []struct {
		status int
		body   string
	}{
		{404, `[{"statusCode":424},{"statusCode":404,"resourceBody":{"code":"NotFound"}},{"statusCode":424}]`},
		{409, `[{"statusCode":424},{"statusCode":409,"eTag":"original"},{"statusCode":424}]`},
		{412, `[{"statusCode":424},{"statusCode":412,"substatusCode":1001,"requestCharge":1.5},{"statusCode":424}]`},
		{429, `[{"statusCode":424},{"statusCode":429,"retryAfterMilliseconds":500},{"statusCode":424}]`},
		{777, `[{"statusCode":424},{"statusCode":777},{"statusCode":424}]`},
	} {
		response, err := decodeTransactionalBatchResponse(ItemResponse{Response: Response{StatusCode: http.StatusMultiStatus}},
			[]byte(failure.body), 0, true, 3)
		require.NoError(t, err)
		require.Equal(t, http.StatusMultiStatus, response.StatusCode)
		require.Len(t, response.OperationResults, 3)
		require.Equal(t, failure.status, response.OperationResults[1].StatusCode)
	}
}

func TestTransactionalBatchResponsePreservesRawOutcomes(t *testing.T) {
	for _, test := range []struct {
		name     string
		status   int
		body     string
		statuses []int
	}{
		{"conditional read", 207, `[{"statusCode":304}]`, []int{304}},
		{"conditional read with write", 207, `[{"statusCode":201},{"statusCode":304}]`, []int{201, 304}},
		{"OK with conditional read", 200, `[{"statusCode":304}]`, []int{304}},
		{"OK with failed operation", 200, `[{"statusCode":404}]`, []int{404}},
		{"multi-status with created item", 207, `[{"statusCode":201}]`, []int{201}},
		{"only dependency failures", 207, `[{"statusCode":424}]`, []int{424}},
		{"multiple failures", 207, `[{"statusCode":404},{"statusCode":412}]`, []int{404, 412}},
		{"other native completion status", 201, `[{"statusCode":201}]`, []int{201}},
		{"unknown operation status", 207, `[{"statusCode":65535}]`, []int{65535}},
	} {
		t.Run(test.name, func(t *testing.T) {
			response, err := decodeTransactionalBatchResponse(ItemResponse{Response: Response{StatusCode: test.status}},
				[]byte(test.body), 0, true, len(test.statuses))
			require.NoError(t, err, "native statuses are preserved without deriving a transaction verdict")
			require.Equal(t, test.status, response.StatusCode)
			require.Equal(t, test.body, string(response.Body))
			require.Len(t, response.OperationResults, len(test.statuses))
			for i, status := range test.statuses {
				require.Equal(t, status, response.OperationResults[i].StatusCode)
			}
		})
	}
}

func TestTransactionalBatchResponsePreservesConditionalReadMetadata(t *testing.T) {
	body := []byte(`[{"statusCode":304,"substatusCode":17,"requestCharge":1.25,"eTag":"\"unchanged\"","retryAfterMilliseconds":25}]`)
	response, err := decodeTransactionalBatchResponse(ItemResponse{Response: Response{StatusCode: 207, RequestCharge: 2.5},
		SessionToken: "0:1", ETag: `"batch"`}, body, time.Second, true, 1)
	require.NoError(t, err)
	require.Equal(t, 207, response.StatusCode)
	require.Equal(t, 2.5, response.RequestCharge)
	require.Equal(t, SessionToken("0:1"), response.SessionToken)
	require.Equal(t, azcore.ETag(`"batch"`), response.ETag)
	require.Equal(t, time.Second, response.RetryAfter)
	require.Equal(t, []TransactionalBatchResult{{
		StatusCode: 304, SubStatus: 17, RequestCharge: 1.25, ETag: `"unchanged"`, RetryAfter: 25 * time.Millisecond,
	}}, response.OperationResults)
	original := string(body)
	for i := range body {
		body[i] = 'x'
	}
	require.Equal(t, original, string(response.Body))
}

func TestTransactionalBatchResponseRejectsMalformedResults(t *testing.T) {
	for _, test := range []struct {
		name   string
		status int
		count  int
		body   string
	}{
		{"invalid JSON", 200, 1, `[`},
		{"invalid UTF-8", 200, 1, "[{\"statusCode\":200,\"eTag\":\"\xff\"}]"},
		{"object envelope", 200, 1, `{"results":[{"statusCode":200}]}`},
		{"null envelope", 200, 1, `null`},
		{"empty envelope", 200, 1, `[]`},
		{"missing result", 200, 2, `[{"statusCode":200}]`},
		{"extra result", 200, 1, `[{"statusCode":200},{"statusCode":200}]`},
		{"null result", 200, 1, `[null]`},
		{"missing status", 200, 1, `[{}]`},
		{"null status", 200, 1, `[{"statusCode":null}]`},
		{"zero status", 200, 1, `[{"statusCode":0}]`},
		{"status overflow", 200, 1, `[{"statusCode":65536}]`},
		{"negative substatus", 200, 1, `[{"statusCode":200,"substatusCode":-1}]`},
		{"substatus overflow", 200, 1, `[{"statusCode":200,"substatusCode":4294967296}]`},
		{"negative charge", 200, 1, `[{"statusCode":200,"requestCharge":-1}]`},
		{"negative retry", 200, 1, `[{"statusCode":200,"retryAfterMilliseconds":-1}]`},
		{"retry overflow", 200, 1, `[{"statusCode":200,"retryAfterMilliseconds":18446744073709551615}]`},
		{"no submitted operations", 200, 0, `[]`},
	} {
		t.Run(test.name, func(t *testing.T) {
			body := []byte(test.body)
			diagnostics := &Diagnostics{AttemptCount: 3}
			metadata := ItemResponse{Response: Response{
				StatusCode: test.status, SubStatus: 17, RequestCharge: 5.25, ActivityID: "activity",
				Diagnostics: diagnostics, AttemptCount: 3,
			}, SessionToken: "0:1", ETag: azcore.ETag(`"etag"`)}
			response, err := decodeTransactionalBatchResponse(metadata, body, time.Second, false, test.count)
			require.Zero(t, response)
			var cosmosErr *Error
			require.ErrorAs(t, err, &cosmosErr)
			require.Equal(t, CodeSerializationFailed, cosmosErr.Code)
			require.Equal(t, test.status, cosmosErr.StatusCode)
			require.Equal(t, 17, cosmosErr.SubStatus)
			require.Equal(t, 5.25, cosmosErr.RequestCharge)
			require.Equal(t, "activity", cosmosErr.ActivityID)
			require.Same(t, diagnostics, cosmosErr.Diagnostics)
			require.Equal(t, uint32(3), cosmosErr.AttemptCount)
			require.Equal(t, SessionToken("0:1"), cosmosErr.SessionToken)
			require.Equal(t, metadata.ETag, cosmosErr.ETag)
			require.Equal(t, time.Second, cosmosErr.RetryAfter)
			require.False(t, cosmosErr.FromWire, "do not infer wire origin from status")
			require.NotNil(t, cosmosErr.Unwrap())
			for i := range body {
				body[i] = 'x'
			}
			require.Equal(t, test.body, string(cosmosErr.Body))
		})
	}
}

func TestTransactionalBatchResponseNullAndAbsentBodies(t *testing.T) {
	body := []byte(`[{"statusCode":200,"resourceBody":null},{"statusCode":204}]`)
	response, err := decodeTransactionalBatchResponse(ItemResponse{Response: Response{StatusCode: 200}}, body, 0, true, 2)
	require.NoError(t, err)
	require.Equal(t, "null", string(response.OperationResults[0].ResourceBody))
	require.Nil(t, response.OperationResults[1].ResourceBody)
	require.False(t, strings.Contains(string(response.Body), `"retryAfterMilliseconds"`))
}
