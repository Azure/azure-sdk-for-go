//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleClient_GetAudioTranscription() {
	azureOpenAIKey := os.Getenv("AOAI_API_KEY_WHISPER")

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_ENDPOINT_WHISPER")

	modelDeploymentID := os.Getenv("AOAI_MODEL_WHISPER")

	if azureOpenAIKey == "" || azureOpenAIEndpoint == "" || modelDeploymentID == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	mp3Bytes, err := os.ReadFile("testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.mp3")

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	resp, err := client.GetAudioTranscription(context.TODO(), azopenai.AudioTranscriptionOptions{
		File: mp3Bytes,

		// this will return _just_ the translated text. Other formats are available, which return
		// different or additional metadata. See [azopenai.AudioTranscriptionFormat] for more examples.
		ResponseFormat: to.Ptr(azopenai.AudioTranscriptionFormatText),

		DeploymentName: &modelDeploymentID,
	}, nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		log.Fatalf("ERROR: %s", err)
	}

	fmt.Fprintf(os.Stderr, "Transcribed text: %s\n", *resp.Text)

	// Output:
}
