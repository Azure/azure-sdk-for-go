//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants_test

import (
	"context"
	"io"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenaiassistants"
)

func TestDownloadFileContent(t *testing.T) {
	skipRecordingsCantMatchRoutesTestHack(t)

	// I haven't come up with a way to guarantee I always get a file download here, which makes this test fail on occasion.
	// If you're attempting (and failing) to get a recording for this you might just have to give it a few tries.
	//
	// So something like this, so you don't have to rerun the entire suite each time:
	// AZURE_RECORD_MODE=record go test -run TestDownloadFileContent
	if recording.GetRecordMode() == recording.LiveMode {
		t.Skip("Skipping non-deterministic test for live tests. Only runs in record/playback.")
	}

	args := runThreadArgs{
		Assistant: azopenaiassistants.CreateAssistantBody{
			DeploymentName: &assistantsModel,
			Instructions:   to.Ptr("You are a helpful assistant that always draws images and provides files for download."),
			Tools: []azopenaiassistants.ToolDefinitionClassification{
				&azopenaiassistants.CodeInterpreterToolDefinition{},
			},
		},
		Thread: azopenaiassistants.CreateAndRunThreadBody{
			Thread: &azopenaiassistants.CreateThreadBody{
				Messages: []azopenaiassistants.CreateMessageBody{
					{
						Role:    to.Ptr(azopenaiassistants.MessageRoleUser),
						Content: to.Ptr("Draw an image of two squares, connected by a line, as a PNG file and make it available for download"),
					},
				},
			},
		},
	}

	client, messages := mustRunThread(context.Background(), t, args)
	fileFound := false

	for _, m := range messages {
		// the assistants reply should contain a file ID for an image to download
		for _, c := range m.Content {
			switch v := c.(type) {
			case *azopenaiassistants.MessageImageFileContent:
				resp, err := client.GetFileContent(context.Background(), *v.ImageFile.FileID, nil)
				require.NoError(t, err)

				defer func() {
					err := resp.Content.Close()
					require.NoError(t, err)
				}()

				fileBytes, err := io.ReadAll(resp.Content)
				require.NoError(t, err)
				require.NotEmpty(t, fileBytes)
				fileFound = true

				t.Logf("[%s] image file ID: %s, file is %d bytes", *m.Role, *v.ImageFile.FileID, len(fileBytes))
			case *azopenaiassistants.MessageTextContent:
				t.Logf("[%s] %s", *m.Role, *v.Text.Value)
			}
		}
	}

	require.True(t, fileFound)
}
