// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azblob

import (
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/internal/shared"
)

const (
	// EventUpload is used for logging events related to upload operation.
	EventUpload = shared.EventUpload

	// EventSubmitBatch is used for logging events related to submit blob batch operation.
	EventSubmitBatch = shared.EventSubmitBatch
)
