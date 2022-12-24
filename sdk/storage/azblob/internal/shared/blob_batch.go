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

const (
	BatchIdPrefix = "batch_"
	HttpVersion   = "HTTP/1.1"
	HttpNewline   = "\r\n"
)

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
	var subReqHeader strings.Builder
	subReqHeader.WriteString(GetBatchRequestDelimiter(batchID, true, false) + HttpNewline)
	subReqHeader.WriteString("Content-Type: application/http" + HttpNewline)
	subReqHeader.WriteString("Content-Transfer-Encoding: binary" + HttpNewline)
	subReqHeader.WriteString("Content-ID: " + strconv.Itoa(contentID) + HttpNewline)
	subReqHeader.WriteString(HttpNewline)

	return subReqHeader.String()
}
