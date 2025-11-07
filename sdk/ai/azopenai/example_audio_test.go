// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/openai/openai-go/v3"
)

// Example_audioTranscription demonstrates how to transcribe speech to text using Azure OpenAI's Whisper model.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Read an audio file and send it to the API
// - Convert spoken language to written text using the Whisper model
// - Process the transcription response
//
// The example uses environment variables for configuration:
// - AOAI_WHISPER_ENDPOINT: Your Azure OpenAI endpoint URL
// - AOAI_WHISPER_MODEL: The deployment name of your Whisper model
//
// Audio transcription is useful for accessibility features, creating searchable archives of audio content,
// generating captions or subtitles, and enabling voice commands in applications.
func Example_audioTranscription() {
	if !CheckRequiredEnvVars("AOAI_WHISPER_ENDPOINT", "AOAI_WHISPER_MODEL") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	endpoint := os.Getenv("AOAI_WHISPER_ENDPOINT")
	model := os.Getenv("AOAI_WHISPER_MODEL")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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

// Example_generateSpeechFromText demonstrates how to convert text to speech using Azure OpenAI's text-to-speech service.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Send text to be converted to speech
// - Specify voice and audio format parameters
// - Handle the audio response stream
//
// The example uses environment variables for configuration:
// - AOAI_TTS_ENDPOINT: Your Azure OpenAI endpoint URL
// - AOAI_TTS_MODEL: The deployment name of your text-to-speech model
//
// Text-to-speech conversion is valuable for creating audiobooks, virtual assistants,
// accessibility tools, and adding voice interfaces to applications.
func Example_generateSpeechFromText() {
	if !CheckRequiredEnvVars("AOAI_TTS_ENDPOINT", "AOAI_TTS_MODEL") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	endpoint := os.Getenv("AOAI_TTS_ENDPOINT")
	model := os.Getenv("AOAI_TTS_MODEL")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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

// Example_audioTranslation demonstrates how to translate speech from one language to English text.
// This example shows how to:
// - Create an Azure OpenAI client with token credentials
// - Read a non-English audio file
// - Translate the spoken content to English text
// - Process the translation response
//
// The example uses environment variables for configuration:
// - AOAI_WHISPER_ENDPOINT: Your Azure OpenAI endpoint URL
// - AOAI_WHISPER_MODEL: The deployment name of your Whisper model
//
// Speech translation is essential for cross-language communication, creating multilingual content,
// and building applications that break down language barriers.
func Example_audioTranslation() {
	if !CheckRequiredEnvVars("AOAI_WHISPER_ENDPOINT", "AOAI_WHISPER_MODEL") {
		fmt.Fprintf(os.Stderr, "Skipping example, environment variables missing\n")
		return
	}

	endpoint := os.Getenv("AOAI_WHISPER_ENDPOINT")
	model := os.Getenv("AOAI_WHISPER_MODEL")

	client, err := CreateOpenAIClientWithToken(endpoint, "")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		return
	}

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
