//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
)

const (
	BatchIdPrefix = "batch_"
	HttpVersion   = "HTTP/1.1"
	HttpNewline   = "\r\n"
)

// CreateBatchID is used for creating a new batch id which is used as batch boundary in the request body
func CreateBatchID() (string, error) {
	batchID, err := uuid.New()
	if err != nil {
		return "", err
	}

	return BatchIdPrefix + batchID.String(), nil
}

// GetBatchRequestDelimiter is used for creating the batch boundary
// e.g. --batch_357de4f7-6d0b-4e02-8cd2-6361411a9525
// last line of the request body: --batch_357de4f7-6d0b-4e02-8cd2-6361411a9525--
func GetBatchRequestDelimiter(batchID string, prefixDash bool, postfixDash bool) string {
	outString := ""

	if prefixDash {
		outString = "--"
	}

	outString += batchID

	if postfixDash {
		outString += "--"
	}

	return outString
}

// CreateSubReqHeader is used to create the sub-request header. Example:
// --batch_357de4f7-6d0b-4e02-8cd2-6361411a9525
// Content-Type: application/http
// Content-Transfer-Encoding: binary
// Content-ID: 0
func CreateSubReqHeader(batchID string, contentID int) string {
	var subReqHeader strings.Builder
	subReqHeader.WriteString(GetBatchRequestDelimiter(batchID, true, false) + HttpNewline)
	subReqHeader.WriteString("Content-Type: application/http" + HttpNewline)
	subReqHeader.WriteString("Content-Transfer-Encoding: binary" + HttpNewline)
	subReqHeader.WriteString("Content-ID: " + strconv.Itoa(contentID) + HttpNewline)
	subReqHeader.WriteString(HttpNewline)

	return subReqHeader.String()
}

// UpdateSubRequestHeaders updates the sub-request headers.
// Adds x-ms-date and removes x-ms-version header
func UpdateSubRequestHeaders(req *policy.Request) {
	// setting x-ms-date header
	dt := time.Now().UTC().Format(http.TimeFormat)
	req.Raw().Header[HeaderXmsDate] = []string{dt}

	// remove x-ms-version header from the request header
	req.Raw().Header.Del(HeaderXmsVersion)
	for k, _ := range req.Raw().Header {
		if strings.ToLower(k) == strings.ToLower(HeaderXmsVersion) {
			delete(req.Raw().Header, k)
		}
	}
}

// BuildSubRequest is used for building the sub-request. Example:
// DELETE /container0/blob0 HTTP/1.1
// x-ms-date: Thu, 14 Jun 2018 16:46:54 GMT
// Authorization: SharedKey account:G4jjBXA7LI/RnWKIOQ8i9xH4p76pAQ+4Fs4R1VxasaE=
// Content-Length: 0
func BuildSubRequest(req *policy.Request) string {
	var batchSubRequest strings.Builder
	blobPath := req.Raw().URL.Path
	if len(req.Raw().URL.RawQuery) > 0 {
		blobPath += "?" + req.Raw().URL.RawQuery
	}

	batchSubRequest.WriteString(fmt.Sprintf("%s %s %s%s", req.Raw().Method, blobPath, HttpVersion, HttpNewline))

	for k, v := range req.Raw().Header {
		if strings.ToLower(k) == HeaderXmsVersion {
			continue
		}
		if len(v) > 0 {
			batchSubRequest.WriteString(fmt.Sprintf("%v: %v%v", k, v[0], HttpNewline))
		}
	}

	batchSubRequest.WriteString(HttpNewline)
	return batchSubRequest.String()
}
