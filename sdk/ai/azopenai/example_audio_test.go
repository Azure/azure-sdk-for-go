// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/azure"
)

func Example_audioTranscription() {
	endpoint := os.Getenv("AOAI_ENDPOINT")
	model := os.Getenv("AOAI_WHISPER_MODEL")

	if endpoint == "" || model == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)

	if err != nil {
		//  TODO: Update the following line with your application specific error handling logic
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-08-01-preview"),
		azure.WithTokenCredential(tokenCredential),
	)

	audio_file, err := os.Open("testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.mp3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}
	defer audio_file.Close()

	resp, err := client.Audio.Transcriptions.New(context.TODO(), openai.AudioTranscriptionNewParams{
		Model:          openai.AudioModel(model),
		File:           audio_file,
		ResponseFormat: openai.AudioResponseFormatJSON,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "Transcribed text: %s\n", resp.Text)
}

func Example_generateSpeechFromText() {
	endpoint := os.Getenv("AOAI_ENDPOINT")
	model := os.Getenv("AOAI_TTS_MODEL")

	if endpoint == "" || model == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-08-01-preview"),
		azure.WithTokenCredential(tokenCredential),
	)

	audioResp, err := client.Audio.Speech.New(context.Background(), openai.AudioSpeechNewParams{
		Model:          openai.SpeechModel(model),
		Input:          "i am a computer",
		Voice:          openai.AudioSpeechNewParamsVoiceAlloy,
		ResponseFormat: openai.AudioSpeechNewParamsResponseFormatFLAC,
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	defer audioResp.Body.Close()

	audioBytes, err := io.ReadAll(audioResp.Body)

	if err != nil {
		// TODO: Update the following line with your application specific error handling logic
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "Got %d bytes of FLAC audio\n", len(audioBytes))
}

func Example_audioTranslation() {
	endpoint := os.Getenv("AOAI_ENDPOINT")
	model := os.Getenv("AOAI_WHISPER_MODEL")

	if endpoint == "" || model == "" {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	tokenCredential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	client := openai.NewClient(
		azure.WithEndpoint(endpoint, "2024-08-01-preview"),
		azure.WithTokenCredential(tokenCredential),
	)

	audio_file, err := os.Open("testdata/sampleaudio_hindi_myVoiceIsMyPassportVerifyMe.mp3")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}
	defer audio_file.Close()

	resp, err := client.Audio.Translations.New(context.TODO(), openai.AudioTranslationNewParams{
		Model:  openai.AudioModel(model),
		File:   audio_file,
		Prompt: openai.String("Translate the following Hindi audio to English"),
	})

	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

	fmt.Fprintf(os.Stderr, "Translated text: %s\n", resp.Text)
}
