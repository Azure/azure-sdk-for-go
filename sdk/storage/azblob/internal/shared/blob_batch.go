//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"io"
	"net/http"
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
func CreateBatchRequest(bb *BlobBatchBuilder) (string, string, error) {
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
				// TODO: handle error: log this error
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

// BlobBatchResponse contains the response for the individual sub-requests.
type BlobBatchResponse struct {
	ContentID     *int
	ContainerName *string
	BlobName      *string
	Response      *http.Response
	Error         error
}

func getResponseBoundary(contentType *string) (string, error) {
	if contentType == nil {
		return "", fmt.Errorf("Content-Type returned in SubmitBatch response is nil")
	}
	boundaryIdx := strings.Index(*contentType, "batchresponse_")
	if boundaryIdx == -1 {
		// TODO: add log
		return "", fmt.Errorf("batch boundary not present in Content-Type header of the SubmitBatch response.\nContent-Type: %v\n", *contentType)
	}
	return "--" + (*contentType)[boundaryIdx:], nil
}

func getContentID(part string) *int {
	contentIDIdx := strings.Index(part, "Content-ID")
	if contentIDIdx == -1 {
		return nil
	}
	contentID := -1
	for i := contentIDIdx + len("Content-ID"); i < len(part); i++ {
		if part[i] >= '0' && part[i] <= '9' {
			if contentID == -1 {
				contentID = 0
			}
			contentID = contentID*10 + int(part[i]-'0')
		} else if contentID != -1 || part[i] == '\n' {
			break
		}
	}
	if contentID == -1 {
		return nil
	}
	return &contentID
}

func ParseBlobBatchResponse(respBody io.ReadCloser, contentType *string, subRequests []*policy.Request) ([]*BlobBatchResponse, error) {
	boundary, err := getResponseBoundary(contentType)
	if err != nil {
		return nil, err
	}

	bytesBody, err := runtime.Payload(&http.Response{
		Body: respBody,
	})
	if err != nil {
		return nil, err
	}
	body := string(bytesBody)

	parts := strings.Split(body, boundary)
	var responses []*BlobBatchResponse
	batchPartialError := false

	for i, part := range parts {
		if i == 0 || i == len(parts)-1 {
			continue
		}

		batchResponse := &BlobBatchResponse{}
		batchResponse.ContentID = getContentID(part)

		if batchResponse.ContentID != nil {
			path := strings.Trim(subRequests[*batchResponse.ContentID].Raw().URL.Path, "/")
			p := strings.Split(path, "/")
			batchResponse.ContainerName = to.Ptr(p[0])
			batchResponse.BlobName = to.Ptr(strings.Join(p[1:], "/"))
		}

		respStringIdx := strings.Index(part, "HTTP/1.1")
		if respStringIdx == -1 {
			// TODO: add log
			continue
		}

		buf := bytes.NewBufferString(part[respStringIdx:])
		resp, err := http.ReadResponse(bufio.NewReader(buf), nil)
		batchResponse.Response = resp
		batchResponse.Error = err
		if err != nil {
			// TODO: add log
			batchPartialError = true
		}
		if resp != nil && (resp.StatusCode < 200 || resp.StatusCode >= 300) {
			// TODO: add log
			batchPartialError = true
			if i == 1 && batchResponse.ContentID == nil {
				// this case can happen when the parent request fails.
				// For example, batch request having more than 256 sub-requests.
				return nil, fmt.Errorf("%v", part[respStringIdx:])
			}
		}

		responses = append(responses, batchResponse)
	}

	if len(responses) != len(subRequests) {
		return responses, fmt.Errorf("expected %v responses, got %v for the batch request", len(subRequests), len(responses))
	}

	if batchPartialError {
		return responses, fmt.Errorf("there is a partial failure in the batch operation")
	} else {
		return responses, nil
	}
}
