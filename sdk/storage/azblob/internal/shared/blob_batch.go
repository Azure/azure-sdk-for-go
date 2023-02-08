//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"strconv"
	"strings"
)

type BlobBatchOperationType string

const (
	BatchIdPrefix                                    = "batch_"
	HttpVersion                                      = "HTTP/1.1"
	HttpNewline                                      = "\r\n"
	BatchDeleteOperationType  BlobBatchOperationType = "delete"
	BatchSetTierOperationType BlobBatchOperationType = "set tier"
)

type BlobBatchBuilder struct {
	Endpoint    *string
	AuthPolicy  policy.Policy
	SubRequests []*policy.Request
}

// createBatchID is used for creating a new batch id which is used as batch boundary in the request body
func createBatchID() (string, error) {
	batchID, err := uuid.New()
	if err != nil {
		return "", err
	}

	return BatchIdPrefix + batchID.String(), nil
}

// getBatchRequestDelimiter is used for creating the batch boundary
// e.g. --batch_357de4f7-6d0b-4e02-8cd2-6361411a9525
// last line of the request body: --batch_357de4f7-6d0b-4e02-8cd2-6361411a9525--
func getBatchRequestDelimiter(batchID *string, prefixDash bool, postfixDash bool) string {
	outString := ""

	if prefixDash {
		outString = "--"
	}

	outString += *batchID

	if postfixDash {
		outString += "--"
	}

	return outString
}

// createSubReqHeader is used to create the sub-request header. Example:
// --batch_357de4f7-6d0b-4e02-8cd2-6361411a9525
// Content-Type: application/http
// Content-Transfer-Encoding: binary
// Content-ID: 0
func createSubReqHeader(batchID *string, contentID *int) string {
	var subReqHeader strings.Builder
	subReqHeader.WriteString(getBatchRequestDelimiter(batchID, true, false) + HttpNewline)
	subReqHeader.WriteString("Content-Type: application/http" + HttpNewline)
	subReqHeader.WriteString("Content-Transfer-Encoding: binary" + HttpNewline)
	subReqHeader.WriteString("Content-ID: " + strconv.Itoa(*contentID) + HttpNewline)
	subReqHeader.WriteString(HttpNewline)

	return subReqHeader.String()
}

// buildSubRequest is used for building the sub-request. Example:
// DELETE /container0/blob0 HTTP/1.1
// x-ms-date: Thu, 14 Jun 2018 16:46:54 GMT
// Authorization: SharedKey account:G4jjBXA7LI/RnWKIOQ8i9xH4p76pAQ+4Fs4R1VxasaE=
// Content-Length: 0
func buildSubRequest(req *policy.Request) string {
	var batchSubRequest strings.Builder
	blobPath := req.Raw().URL.Path
	if len(req.Raw().URL.RawQuery) > 0 {
		blobPath += "?" + req.Raw().URL.RawQuery
	}

	batchSubRequest.WriteString(fmt.Sprintf("%s %s %s%s", req.Raw().Method, blobPath, HttpVersion, HttpNewline))

	for k, v := range req.Raw().Header {
		if strings.EqualFold(k, HeaderXmsVersion) {
			continue
		}
		if len(v) > 0 {
			batchSubRequest.WriteString(fmt.Sprintf("%v: %v%v", k, v[0], HttpNewline))
		}
	}

	batchSubRequest.WriteString(HttpNewline)
	return batchSubRequest.String()
}

// CreateBatchRequest creates a new batch request using the sub-requests present in the BlobBatchBuilder.
//
// Example of a sub-request in the batch request body:
//
//	--batch_357de4f7-6d0b-4e02-8cd2-6361411a9525
//	Content-Type: application/http
//	Content-Transfer-Encoding: binary
//	Content-ID: 0
//
//	DELETE /container0/blob0 HTTP/1.1
//	x-ms-date: Thu, 14 Jun 2018 16:46:54 GMT
//	Authorization: SharedKey account:G4jjBXA7LI/RnWKIOQ8i9xH4p76pAQ+4Fs4R1VxasaE=
//	Content-Length: 0
func CreateBatchRequest(ctx context.Context, bb *BlobBatchBuilder) (string, string, error) {
	batchID, err := createBatchID()
	if err != nil {
		return "", "", err
	}

	contentID := 0
	var batchRequest strings.Builder

	for _, req := range bb.SubRequests {
		if bb.AuthPolicy != nil {
			resp, err := bb.AuthPolicy.Do(req)
			if err != nil && resp != nil {
				// TODO: handle error
				continue
			}
		}

		batchRequest.WriteString(createSubReqHeader(&batchID, &contentID))
		batchRequest.WriteString(buildSubRequest(req))
		contentID++
	}

	// add the last line of the request body. It looks like,
	// --batch_357de4f7-6d0b-4e02-8cd2-6361411a9525--
	batchRequest.WriteString(getBatchRequestDelimiter(&batchID, true, true) + HttpNewline)

	return batchRequest.String(), batchID, nil
}

// UpdateSubRequestHeaders updates the sub-request headers.
// Removes x-ms-version header.
func UpdateSubRequestHeaders(req *policy.Request) {
	// remove x-ms-version header from the request header
	for k := range req.Raw().Header {
		if strings.EqualFold(k, HeaderXmsVersion) {
			delete(req.Raw().Header, k)
		}
	}
}
