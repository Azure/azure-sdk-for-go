// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package arrow

import (
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/exported"
)

const (
	// ArrowContentType is the Content-Type for Apache Arrow IPC stream responses.
	ArrowContentType = "application/vnd.apache.arrow.stream"

	// ArrowAcceptHeader is the Accept header value to request Arrow format.
	ArrowAcceptHeader = ArrowContentType
)

// ErrNotSupported is returned when the Arrow format is requested explicitly
// and the module was built with the azblob_noarrow tag.
var ErrNotSupported = errors.New("the Arrow response format is not available: built with the azblob_noarrow tag")

// UseArrow reports whether a list operation requests the Arrow format.
// Auto resolves to XML when Arrow is not compiled in; an explicit Arrow request fails.
func UseArrow(f exported.StorageResponseFormat) (bool, error) {
	if f == exported.StorageResponseFormatAuto && !Supported {
		return false, nil
	}
	if exported.ResolveAutoFormat(f) != exported.StorageResponseFormatArrow {
		return false, nil
	}
	if !Supported {
		return false, ErrNotSupported
	}
	return true, nil
}
