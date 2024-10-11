//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"bytes"
	"context"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestFilesOperations(t *testing.T) {
	if recording.GetRecordMode() != recording.LiveMode {
		t.Skip("Skipping test in playback and record mode")
	}

	client := newTestClient(t, azureOpenAI.Files.Endpoint)

	t.Run("UploadParts", func(t *testing.T) {
		createUploadResp, err := client.CreateUpload(context.Background(), azopenai.CreateUploadRequest{
			Bytes:    to.Ptr(int32(10)),
			Filename: to.Ptr("test.txt"),
			MimeType: to.Ptr("text/plain"),
			Purpose:  to.Ptr(azopenai.CreateUploadRequestPurposeAssistants),
		}, nil)
		require.NoError(t, err)

		part1 := streaming.NopCloser(strings.NewReader("hello"))
		part2 := streaming.NopCloser(strings.NewReader("world"))

		// We can upload in any order, and in parallel if we wanted to. The CompleteUpload() call
		// specifies the ordering when the file is assembled.
		part2Resp, err := client.AddUploadPart(context.Background(), *createUploadResp.ID, part2, nil)
		require.NoError(t, err)

		part1Resp, err := client.AddUploadPart(context.Background(), *createUploadResp.ID, part1, nil)
		require.NoError(t, err)

		uploadResp, err := client.CompleteUpload(context.Background(), *createUploadResp.ID, azopenai.CompleteUploadRequest{
			PartIDs: []string{*part2Resp.ID, *part1Resp.ID},
		}, nil)
		require.NoError(t, err)

		// the total size of parts 1 and 2
		require.Equal(t, int64(10), *uploadResp.Bytes)

		// delete the uploaded file.
		_, err = client.DeleteFile(context.Background(), *uploadResp.File.ID, nil)
		require.NoError(t, err)
	})

	t.Run("UploadFile", func(t *testing.T) {
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
	})
}

func TestFileDownload(t *testing.T) {
	t.Skip("Need to find a file type we can download")
}

func TestBatchOperations(t *testing.T) {
	if recording.GetRecordMode() != recording.LiveMode {
		t.Skip("Skipping test in playback and record mode")
	}
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
