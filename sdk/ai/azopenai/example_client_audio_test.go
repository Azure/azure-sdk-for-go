//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
)

func ExampleClient_GetAudioTranscription() {
	azureOpenAIKey := os.Getenv("AOAI_WHISPER_API_KEY")

	// Ex: "https://<your-azure-openai-host>.openai.azure.com"
	azureOpenAIEndpoint := os.Getenv("AOAI_WHISPER_ENDPOINT")

	modelDeploymentID := os.Getenv("AOAI_WHISPER_MODEL")

	if azureOpenAIKey == "" || azureOpenAIEndpoint == "" || modelDeploymentID == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(azureOpenAIKey)

	client, err := azopenai.NewClientWithKeyCredential(azureOpenAIEndpoint, keyCredential, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	mp3Bytes, err := os.ReadFile("testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.mp3")

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	resp, err := client.GetAudioTranscription(context.TODO(), azopenai.AudioTranscriptionOptions{
		File: mp3Bytes,

		// this will return _just_ the translated text. Other formats are available, which return
		// different or additional metadata. See [azopenai.AudioTranscriptionFormat] for more examples.
		ResponseFormat: to.Ptr(azopenai.AudioTranscriptionFormatText),

		DeploymentName: &modelDeploymentID,
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	fmt.Fprintf(os.Stderr, "Transcribed text: %s\n", *resp.Text)

	// Output:
}

func ExampleClient_GenerateSpeechFromText() {
	openAIKey := os.Getenv("OPENAI_API_KEY")

	// Ex: "https://api.openai.com/v1"
	openAIEndpoint := os.Getenv("OPENAI_ENDPOINT")

	modelDeploymentID := "tts-1"

	if openAIKey == "" || openAIEndpoint == "" || modelDeploymentID == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	keyCredential := azcore.NewKeyCredential(openAIKey)

	client, err := azopenai.NewClientForOpenAI(openAIEndpoint, keyCredential, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	audioResp, err := client.GenerateSpeechFromText(context.Background(), azopenai.SpeechGenerationOptions{
		Input:          to.Ptr("i am a computer"),
		Voice:          to.Ptr(azopenai.SpeechVoiceAlloy),
		ResponseFormat: to.Ptr(azopenai.SpeechGenerationResponseFormatFlac),
		DeploymentName: to.Ptr("tts-1"),
	}, nil)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	defer audioResp.Body.Close()

	audioBytes, err := io.ReadAll(audioResp.Body)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		log.Printf("ERROR: %s", err)
		return
	}

	fmt.Fprintf(os.Stderr, "Got %d bytes of FLAC audio\n", len(audioBytes))

	// Output:
}
