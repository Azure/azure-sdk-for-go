//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"bytes"
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/stretchr/testify/require"
)

func TestFilesOperations(t *testing.T) {
	client := newTestClient(t, azureOpenAI.Files.Endpoint)

	uploadResp, err := client.UploadFile(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("hello world"))), azopenai.FilePurposeAssistants, nil)
	require.NoError(t, err)

	t.Cleanup(func() {
		_, err := client.DeleteFile(context.Background(), *uploadResp.ID, nil)
		require.NoError(t, err)
	})

	getFileResp, err := client.GetFile(context.Background(), *uploadResp.ID, nil)
	require.NoError(t, err)

	require.Equal(t, azopenai.FilePurposeAssistants, getFileResp.Purpose)

	filesResp, err := client.ListFiles(context.Background(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, filesResp.Data)
}
