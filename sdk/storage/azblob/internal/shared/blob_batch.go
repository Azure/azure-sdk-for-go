//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package shared

import (
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"strconv"
	"strings"
)

const BatchIdPrefix = "batch_"

func CreateBatchID() (string, error) {
	batchID, err := uuid.New()
	if err != nil {
		return "", err
	}

	return BatchIdPrefix + batchID.String(), nil
}

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

func CreateSubReqHeader(batchID string, contentID int) string {
	var subReqHeader []string
	subReqHeader = append(subReqHeader, GetBatchRequestDelimiter(batchID, true, false))
	subReqHeader = append(subReqHeader, "Content-Type: application/http")
	subReqHeader = append(subReqHeader, "Content-Transfer-Encoding: binary")
	subReqHeader = append(subReqHeader, "Content-ID: "+strconv.Itoa(contentID))
	subReqHeader = append(subReqHeader, "")

	return strings.Join(subReqHeader, "\n")
}

//func SerializeBatchBodyForDelete(ctx context.Context, blobs []BatchDeleteOptions, batchid string) (string, error) {
//	/* Batch Body
//	--<batchid>
//	Content-Type: application/http
//	Content-Transfer-Encoding: binary
//	Content-ID: 0
//
//	DELETE /container0/blob0 HTTP/1.1
//	x-ms-date: <date>
//	Authorization: <auth header>
//	Content-Length: 0
//
//	--<batchid>
//	<subrequest>
//	.
//	.
//	--<batchid>--
//	*/
//
//	batchDelimiter := GetBatchRequestDelimiter(batchid, true, false)
//
//	var batchBody []string
//	for contentID, blob := range blobs {
//		// Put batch delimeter first
//		batchBody = append(batchBody, batchDelimiter)
//
//		// Below are fixed headers
//		batchBody = append(batchBody, "Content-Type: application/http")
//		batchBody = append(batchBody, "Content-Transfer-Encoding: binary")
//		batchBody = append(batchBody, "Content-ID: "+strconv.Itoa(contentID))
//		batchBody = append(batchBody, "")
//
//		// Sub request goes here
//		req, err := runtime.NewRequest(ctx, http.MethodDelete, blob)
//		if err != nil {
//			return "", err
//		}
//		req.Raw().Header.Set(HeaderXmsDate, time.Now().UTC().Format(http.TimeFormat))
//		req.Raw().Header.Set(HeaderContentLength, "0")
//		//req.Raw().Header.Set(HeaderAuthorization, "")
//
//		subRequest, err := ioutil.ReadAll(req.Body())
//		if err != nil {
//			return "", err
//		}
//
//		batchBody = append(batchBody, string(subRequest))
//		req.Close()
//
//		// Append empty line for batch item seperator
//		batchBody = append(batchBody, "")
//	}
//
//	batchTerminator := GetBatchRequestDelimiter(batchid, true, true)
//	batchBody = append(batchBody, batchTerminator)
//
//	return strings.Join(batchBody, "\n"), nil
//}
