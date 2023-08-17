//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package exported

import (
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/blob"
	"strconv"
	"strings"
)

const SnapshotTimeFormat = "2006-01-02T15:04:05.0000000Z07:00"

// HTTPRange defines a range of bytes within an HTTP resource, starting at offset and
// ending at offset+count. A zero-value HTTPRange indicates the entire resource. An HTTPRange
// which has an offset but no zero value count indicates from the offset to the resource's end.
type HTTPRange = blob.HTTPRange

// FormatHTTPRange converts an HTTPRange to its string format.
func FormatHTTPRange(r HTTPRange) *string {
	if r.Offset == 0 && r.Count == 0 {
		return nil // No specified range
	}
	endOffset := "" // if count == CountToEnd (0)
	if r.Count > 0 {
		endOffset = strconv.FormatInt((r.Offset+r.Count)-1, 10)
	}
	dataRange := fmt.Sprintf("bytes=%v-%s", r.Offset, endOffset)
	return &dataRange
}

func ConvertToDFSError(err error) error {
	if err == nil {
		return nil
	}
	var responseErr *azcore.ResponseError
	isRespErr := errors.As(err, &responseErr)
	if isRespErr {
		responseErr.ErrorCode = strings.Replace(responseErr.ErrorCode, "blob", "path", -1)
		responseErr.ErrorCode = strings.Replace(responseErr.ErrorCode, "Blob", "Path", -1)
		responseErr.ErrorCode = strings.Replace(responseErr.ErrorCode, "container", "filesystem", -1)
		responseErr.ErrorCode = strings.Replace(responseErr.ErrorCode, "Container", "FileSystem", -1)
		return responseErr
	}
	return err
}
