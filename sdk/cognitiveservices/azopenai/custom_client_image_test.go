//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"image/png"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/cognitiveservices/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestImageGeneration_AzureOpenAI(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skipf("Ignoring poller-based test")
	}

	cred, err := azopenai.NewKeyCredential(azureOpenAI.APIKey)
	require.NoError(t, err)

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAI.Endpoint, cred, newClientOptionsForTest(t))
	require.NoError(t, err)

	testImageGeneration(t, client, azopenai.ImageGenerationResponseFormatURL)
}

func TestImageGeneration_OpenAI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping OpenAI tests when attempting to do quick tests")
	}

	client := newOpenAIClientForTest(t)
	testImageGeneration(t, client, azopenai.ImageGenerationResponseFormatURL)
}

func TestImageGeneration_AzureOpenAI_WithError(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip()
	}

	client := newBogusAzureOpenAIClient(t)
	testImageGenerationFailure(t, client)
}

func TestImageGeneration_OpenAI_WithError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping OpenAI tests when attempting to do quick tests")
	}

	client := newBogusOpenAIClient(t)
	testImageGenerationFailure(t, client)
}

func TestImageGeneration_OpenAI_Base64(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping OpenAI tests when attempting to do quick tests")
	}

	client := newOpenAIClientForTest(t)
	testImageGeneration(t, client, azopenai.ImageGenerationResponseFormatB64JSON)
}

func testImageGeneration(t *testing.T, client *azopenai.Client, responseFormat azopenai.ImageGenerationResponseFormat) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	resp, err := client.CreateImage(ctx, azopenai.ImageGenerationOptions{
		Prompt:         to.Ptr("a cat"),
		Size:           to.Ptr(azopenai.ImageSize256x256),
		ResponseFormat: &responseFormat,
	}, nil)
	require.NoError(t, err)

	if recording.GetRecordMode() == recording.LiveMode {
		switch responseFormat {
		case azopenai.ImageGenerationResponseFormatURL:
			headResp, err := http.DefaultClient.Head(*resp.Data[0].URL)
			require.NoError(t, err)
			require.Equal(t, http.StatusOK, headResp.StatusCode)
		case azopenai.ImageGenerationResponseFormatB64JSON:
			pngBytes, err := base64.StdEncoding.DecodeString(*resp.Data[0].Base64Data)
			require.NoError(t, err)
			require.NotEmpty(t, pngBytes)

			// the bytes here should just be a valid PNG
			buff := bytes.NewBuffer(pngBytes)

			// just check that it's a valid PNG
			_, err = png.Decode(buff)
			require.NoError(t, err)
		}
	}
}

func testImageGenerationFailure(t *testing.T, bogusClient *azopenai.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	resp, err := bogusClient.CreateImage(ctx, azopenai.ImageGenerationOptions{
		Prompt:         to.Ptr("a cat"),
		Size:           to.Ptr(azopenai.ImageSize256x256),
		ResponseFormat: to.Ptr(azopenai.ImageGenerationResponseFormatURL),
	}, nil)
	require.Empty(t, resp)

	assertResponseIsError(t, err)
}
