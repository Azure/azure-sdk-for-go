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

	require.Equal(t, azopenai.FilePurposeAssistants, *getFileResp.Purpose)

	// fileContentsResp, err := client.GetFileContent(context.Background(), *getFileResp.ID, nil)
	// require.NoError(t, err)
	// require.NotEmpty(t, fileContentsResp.Value)

	filesResp, err := client.ListFiles(context.Background(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, filesResp.Data)
}

func TestFileDownload(t *testing.T) {
	t.Skip("Need to find a file type we can download")
}

func TestBatchOperations(t *testing.T) {
	client := newTestClient(t, azureOpenAI.Files.Endpoint)

	// TODO: this is a little tricky because the files aren't instantly uploaded, so we can't just proceed
	// with the rest of the test.

	// uploadResp, err := client.UploadFile(context.Background(), streaming.NopCloser(bytes.NewReader([]byte("{}"))), azopenai.FilePurposeBatch, &azopenai.UploadFileOptions{
	// 	Filename: to.Ptr("file.jsonl"),
	// })
	// require.NoError(t, err)

	// t.Cleanup(func() {
	// 	_, err := client.DeleteFile(context.Background(), *uploadResp.ID, nil)
	// 	require.NoError(t, err)
	// })

	// createResp, err := client.CreateBatch(context.Background(), azopenai.BatchCreateRequest{
	// 	InputFileID: uploadResp.ID,
	// }, nil)
	// require.NoError(t, err)
	// require.NotEmpty(t, createResp)

	batchPager := client.NewListBatchesPager(nil)

	for batchPager.More() {
		resp, err := batchPager.NextPage(context.Background())
		require.NoError(t, err)
		require.NotEmpty(t, resp)
	}
}
