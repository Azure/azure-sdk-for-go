//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"strings"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestListFiles(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	testFn := func(t *testing.T, azure bool) {
		client := mustGetClient(t, newClientArgs{Azure: azure})

		// upload some files to guarantee that we have at least some data
		expected := map[string]bool{}

		for i := 0; i < 5; i++ {
			fileName := getFileName(t, "txt")
			resp, err := client.UploadFile(context.Background(), strings.NewReader("hello world"), azopenaiassistants.FilePurposeAssistants, &azopenaiassistants.UploadFileOptions{
				Filename: fileName,
			})
			require.NoError(t, err)
			expected[*fileName] = true

			defer mustDeleteFile(t, client, *resp.ID)
		}

		listFilesResp, err := client.ListFiles(context.Background(), &azopenaiassistants.ListFilesOptions{
			Purpose: to.Ptr(azopenaiassistants.FilePurposeAssistants),
		})
		require.NoError(t, err)
		// files aren't a scoped resource - ie, all the assistant files for any assistants are listed in this call.
		// We expect at least 5 (that we uploaded), but there could be many more since we have concurrent
		// activity from other people.
		require.LessOrEqual(t, 5, len(listFilesResp.Data), "We have at least 5 files")

		for _, file := range listFilesResp.Data {
			if expected[*file.Filename] {
				delete(expected, *file.Filename)
			}
		}

		require.Empty(t, expected)
	}

	t.Run("OpenAI", func(t *testing.T) {
		testFn(t, false)
	})

	t.Run("AzureOpenAI", func(t *testing.T) {
		testFn(t, true)
	})
}
