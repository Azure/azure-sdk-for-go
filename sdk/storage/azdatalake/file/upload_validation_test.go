//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package file

import (
	"context"
	"os"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake/datalakeerror"
	"github.com/stretchr/testify/require"
)

// TestRejectPrecomputedValidation verifies the helper only rejects user-supplied precomputed
// checksums, while allowing computed CRC64 and structured message CRC64 (both validated per chunk).
func TestRejectPrecomputedValidation(t *testing.T) {
	// A precomputed CRC64 cannot validate the individual chunks of a multipart upload.
	require.ErrorIs(t, rejectPrecomputedValidation(TransferValidationTypeCRC64(12345)), datalakeerror.UnsupportedChecksum)

	// Computed CRC64 hashes each chunk on the fly, so it is allowed.
	require.NoError(t, rejectPrecomputedValidation(TransferValidationTypeComputeCRC64()))

	// Structured message CRC64 embeds per-segment checksums, so it is allowed.
	require.NoError(t, rejectPrecomputedValidation(TransferValidationTypeComputeStructuredMessageCRC64(0)))

	// No validation is allowed.
	require.NoError(t, rejectPrecomputedValidation(nil))
}

// TestUploadRejectsPrecomputedValidation verifies that each multipart upload entry point rejects a
// precomputed checksum before issuing any request. A no-credential client with an unreachable
// endpoint is sufficient because the guard returns before any network I/O.
func TestUploadRejectsPrecomputedValidation(t *testing.T) {
	client, err := NewClientWithNoCredential("https://fake.dfs.core.windows.net/fs/file", nil)
	require.NoError(t, err)

	content := []byte("some content to upload in multiple chunks")

	uploadErr := client.UploadBuffer(context.Background(), content, &UploadBufferOptions{
		TransactionalValidation: TransferValidationTypeCRC64(98765),
	})
	require.ErrorIs(t, uploadErr, datalakeerror.UnsupportedChecksum)

	uploadErr = client.UploadStream(context.Background(), strings.NewReader(string(content)), &UploadStreamOptions{
		TransactionalValidation: TransferValidationTypeCRC64(98765),
	})
	require.ErrorIs(t, uploadErr, datalakeerror.UnsupportedChecksum)

	tmp, err := os.CreateTemp(t.TempDir(), "sm-upload-*")
	require.NoError(t, err)
	_, err = tmp.Write(content)
	require.NoError(t, err)
	require.NoError(t, tmp.Close())

	f, err := os.Open(tmp.Name())
	require.NoError(t, err)
	defer f.Close()

	uploadErr = client.UploadFile(context.Background(), f, &UploadFileOptions{
		TransactionalValidation: TransferValidationTypeCRC64(98765),
	})
	require.ErrorIs(t, uploadErr, datalakeerror.UnsupportedChecksum)
}
