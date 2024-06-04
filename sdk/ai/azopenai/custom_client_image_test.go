//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/stretchr/testify/require"
)

func TestImageGeneration_AzureOpenAI(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skipf("Ignoring poller-based test")
	}

	client := newTestClient(t, azureOpenAI.DallE.Endpoint)
	testImageGeneration(t, client, azureOpenAI.DallE.Model, azopenai.ImageGenerationResponseFormatURL, true)
}

func TestImageGeneration_OpenAI(t *testing.T) {
	client := newTestClient(t, openAI.DallE.Endpoint)
	testImageGeneration(t, client, openAI.DallE.Model, azopenai.ImageGenerationResponseFormatURL, false)
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

	client := newTestClient(t, openAI.DallE.Endpoint)
	testImageGeneration(t, client, openAI.DallE.Model, azopenai.ImageGenerationResponseFormatBase64, false)
}

func testImageGeneration(t *testing.T, client *azopenai.Client, model string, responseFormat azopenai.ImageGenerationResponseFormat, azure bool) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	resp, err := client.GetImageGenerations(ctx, azopenai.ImageGenerationOptions{
		// saw this prompt in a thread about trying to _prevent_ Dall-E3 from rewriting your
		// propmt. When this is revised you'll see the text in the
		Prompt:         to.Ptr("acrylic painting of a sunflower with bees"),
		Size:           to.Ptr(azopenai.ImageSizeSize1024X1792),
		ResponseFormat: &responseFormat,
		DeploymentName: &model,
	}, nil)
	customRequireNoError(t, err, azure)

	if recording.GetRecordMode() == recording.LiveMode {
		switch responseFormat {
		case azopenai.ImageGenerationResponseFormatURL:
			headResp, err := http.DefaultClient.Head(*resp.Data[0].URL)
			require.NoError(t, err)

			headResp.Body.Close()
			require.Equal(t, http.StatusOK, headResp.StatusCode)
			require.NotEmpty(t, resp.Data[0].RevisedPrompt)
		case azopenai.ImageGenerationResponseFormatBase64:
			imgBytes, err := base64.StdEncoding.DecodeString(*resp.Data[0].Base64Data)
			require.NoError(t, err)
			require.NotEmpty(t, imgBytes)
			require.NotEmpty(t, resp.Data[0].RevisedPrompt)
		}
	}
}

func testImageGenerationFailure(t *testing.T, bogusClient *azopenai.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	resp, err := bogusClient.GetImageGenerations(ctx, azopenai.ImageGenerationOptions{
		Prompt:         to.Ptr("a cat"),
		Size:           to.Ptr(azopenai.ImageSizeSize256X256),
		ResponseFormat: to.Ptr(azopenai.ImageGenerationResponseFormatURL),
		DeploymentName: to.Ptr("ignored"),
	}, nil)
	require.Empty(t, resp)

	assertResponseIsError(t, err)
}
