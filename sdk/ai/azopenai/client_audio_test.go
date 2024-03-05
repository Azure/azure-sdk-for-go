//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestClient_GetAudioTranscription_AzureOpenAI(t *testing.T) {
	client := newTestClient(t, azureOpenAI.Whisper.Endpoint, withForgivingRetryOption())
	runTranscriptionTests(t, client, azureOpenAI.Whisper.Model, true)
}

func TestClient_GetAudioTranscription_OpenAI(t *testing.T) {
	client := newOpenAIClientForTest(t)
	runTranscriptionTests(t, client, openAI.Whisper.Model, false)
}

func TestClient_GetAudioTranslation_AzureOpenAI(t *testing.T) {
	client := newTestClient(t, azureOpenAI.Whisper.Endpoint, withForgivingRetryOption())
	runTranslationTests(t, client, azureOpenAI.Whisper.Model, true)
}

func TestClient_GetAudioTranslation_OpenAI(t *testing.T) {
	client := newOpenAIClientForTest(t)
	runTranslationTests(t, client, openAI.Whisper.Model, false)
}

func TestClient_GetAudioSpeech(t *testing.T) {
	client := newOpenAIClientForTest(t)

	audioResp, err := client.GenerateSpeechFromText(context.Background(), azopenai.SpeechGenerationOptions{
		Input:          to.Ptr("i am a computer"),
		Voice:          to.Ptr(azopenai.SpeechVoiceAlloy),
		ResponseFormat: to.Ptr(azopenai.SpeechGenerationResponseFormatFlac),
		DeploymentName: to.Ptr("tts-1"),
	}, nil)
	require.NoError(t, err)

	audioBytes, err := io.ReadAll(audioResp.Body)
	require.NoError(t, err)

	require.NotEmpty(t, audioBytes)
	require.Equal(t, "fLaC", string(audioBytes[0:4]))

	// now send _it_ back through the transcription API and see if we can get something useful.
	transcriptionResp, err := client.GetAudioTranscription(context.Background(), azopenai.AudioTranscriptionOptions{
		Filename:       to.Ptr("test.flac"),
		File:           audioBytes,
		ResponseFormat: to.Ptr(azopenai.AudioTranscriptionFormatVerboseJSON),
		DeploymentName: &openAI.Whisper.Model,
		Temperature:    to.Ptr[float32](0.0),
	}, nil)
	require.NoError(t, err)

	require.NotZero(t, *transcriptionResp.Duration)

	// it occasionally comes back with different punctuation or makes a complete sentence but
	// the major words always come through.
	require.Contains(t, *transcriptionResp.Text, "computer")
}

func runTranscriptionTests(t *testing.T, client *azopenai.Client, model string, isAzure bool) {
	// We're experiencing load issues on some of our shared test resources so we'll just spot check.
	// The bulk of the logic will test against OpenAI anyways.
	if isAzure {
		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatText, "m4a"), func(t *testing.T) {
			args := newTranscriptionOptions(azopenai.AudioTranscriptionFormatText, model, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.m4a")
			transcriptResp, err := client.GetAudioTranscription(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			requireEmptyAudioTranscription(t, transcriptResp.AudioTranscription)
		})

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatVerboseJSON, "mp3"), func(t *testing.T) {
			args := newTranscriptionOptions(azopenai.AudioTranscriptionFormatVerboseJSON, model, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.mp3")
			transcriptResp, err := client.GetAudioTranscription(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			require.Greater(t, *transcriptResp.Duration, float32(0.0))
			require.NotEmpty(t, *transcriptResp.Language)
			require.NotEmpty(t, transcriptResp.Segments)
			require.NotEmpty(t, transcriptResp.Segments[0])
			require.NotEmpty(t, transcriptResp.Task)
		})
		return
	}

	testFiles := []string{
		`testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.m4a`,
		`testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.mp3`,
	}

	for _, audioFile := range testFiles {
		ext := filepath.Ext(audioFile)

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatText, ext), func(t *testing.T) {
			args := newTranscriptionOptions(azopenai.AudioTranscriptionFormatText, model, audioFile)
			transcriptResp, err := client.GetAudioTranscription(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			requireEmptyAudioTranscription(t, transcriptResp.AudioTranscription)
		})

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatSrt, ext), func(t *testing.T) {
			args := newTranscriptionOptions(azopenai.AudioTranscriptionFormatSrt, model, audioFile)
			transcriptResp, err := client.GetAudioTranscription(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			requireEmptyAudioTranscription(t, transcriptResp.AudioTranscription)
		})

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatVtt, ext), func(t *testing.T) {
			args := newTranscriptionOptions(azopenai.AudioTranscriptionFormatVtt, model, audioFile)
			transcriptResp, err := client.GetAudioTranscription(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			requireEmptyAudioTranscription(t, transcriptResp.AudioTranscription)
		})

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatVerboseJSON, ext), func(t *testing.T) {
			args := newTranscriptionOptions(azopenai.AudioTranscriptionFormatVerboseJSON, model, audioFile)
			transcriptResp, err := client.GetAudioTranscription(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			require.Greater(t, *transcriptResp.Duration, float32(0.0))
			require.NotEmpty(t, *transcriptResp.Language)
			require.NotEmpty(t, transcriptResp.Segments)
			require.NotEmpty(t, transcriptResp.Segments[0])
			require.NotEmpty(t, transcriptResp.Task)
		})

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatJSON, ext), func(t *testing.T) {
			args := newTranscriptionOptions(azopenai.AudioTranscriptionFormatJSON, model, audioFile)
			transcriptResp, err := client.GetAudioTranscription(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			requireEmptyAudioTranscription(t, transcriptResp.AudioTranscription)
		})
	}
}

func runTranslationTests(t *testing.T, client *azopenai.Client, model string, isAzure bool) {
	// We're experiencing load issues on some of our shared test resources so we'll just spot check.
	// The bulk of the logic will test against OpenAI anyways.
	if isAzure {
		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatText, "m4a"), func(t *testing.T) {
			args := newTranslationOptions(azopenai.AudioTranslationFormatText, model, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.m4a")
			transcriptResp, err := client.GetAudioTranslation(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			requireEmptyAudioTranslation(t, transcriptResp.AudioTranslation)
		})

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatVerboseJSON, "mp3"), func(t *testing.T) {
			args := newTranslationOptions(azopenai.AudioTranslationFormatVerboseJSON, model, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.mp3")
			transcriptResp, err := client.GetAudioTranslation(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			require.Greater(t, *transcriptResp.Duration, float32(0.0))
			require.NotEmpty(t, *transcriptResp.Language)
			require.NotEmpty(t, transcriptResp.Segments)
			require.NotEmpty(t, transcriptResp.Segments[0])
			require.NotEmpty(t, transcriptResp.Task)
		})

		return
	}

	testFiles := []string{
		`testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.m4a`,
		`testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.mp3`,
	}

	for _, audioFile := range testFiles {
		ext := filepath.Ext(audioFile)

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatText, ext), func(t *testing.T) {
			args := newTranslationOptions(azopenai.AudioTranslationFormatText, model, audioFile)
			transcriptResp, err := client.GetAudioTranslation(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			requireEmptyAudioTranslation(t, transcriptResp.AudioTranslation)
		})

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatSrt, ext), func(t *testing.T) {
			args := newTranslationOptions(azopenai.AudioTranslationFormatSrt, model, audioFile)
			transcriptResp, err := client.GetAudioTranslation(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			requireEmptyAudioTranslation(t, transcriptResp.AudioTranslation)
		})

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatVtt, ext), func(t *testing.T) {
			args := newTranslationOptions(azopenai.AudioTranslationFormatVtt, model, audioFile)
			transcriptResp, err := client.GetAudioTranslation(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			requireEmptyAudioTranslation(t, transcriptResp.AudioTranslation)
		})

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatVerboseJSON, ext), func(t *testing.T) {
			args := newTranslationOptions(azopenai.AudioTranslationFormatVerboseJSON, model, audioFile)
			transcriptResp, err := client.GetAudioTranslation(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			require.Greater(t, *transcriptResp.Duration, float32(0.0))
			require.NotEmpty(t, *transcriptResp.Language)
			require.NotEmpty(t, transcriptResp.Segments)
			require.NotEmpty(t, transcriptResp.Segments[0])
			require.NotEmpty(t, transcriptResp.Task)
		})

		t.Run(fmt.Sprintf("%s (%s)", azopenai.AudioTranscriptionFormatJSON, ext), func(t *testing.T) {
			args := newTranslationOptions(azopenai.AudioTranslationFormatJSON, model, audioFile)
			transcriptResp, err := client.GetAudioTranslation(context.Background(), args, nil)
			require.NoError(t, err)
			require.NotEmpty(t, transcriptResp)

			require.NotEmpty(t, *transcriptResp.Text)
			requireEmptyAudioTranslation(t, transcriptResp.AudioTranslation)
		})
	}
}

func newTranscriptionOptions(format azopenai.AudioTranscriptionFormat, model string, path string) azopenai.AudioTranscriptionOptions {
	audioBytes, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	return azopenai.AudioTranscriptionOptions{
		DeploymentName: to.Ptr(model),
		File:           audioBytes,
		Filename:       &path,
		ResponseFormat: &format,
		Language:       to.Ptr("en"),
		Temperature:    to.Ptr[float32](0.0),
	}
}

func newTranslationOptions(format azopenai.AudioTranslationFormat, model string, path string) azopenai.AudioTranslationOptions {
	audioBytes, err := os.ReadFile(path)

	if err != nil {
		panic(err)
	}

	var filename *string

	if filepath.Ext(path) != ".mp3" {
		filename = &path
	}

	return azopenai.AudioTranslationOptions{
		DeploymentName: to.Ptr(model),
		File:           audioBytes,
		Filename:       filename,
		ResponseFormat: &format,
		Temperature:    to.Ptr[float32](0.0),
	}
}

// requireEmptyAudioTranscription checks that all the attributes are empty (aside
// from Text)
func requireEmptyAudioTranscription(t *testing.T, at azopenai.AudioTranscription) {
	// Text is always filled out for

	require.Empty(t, at.Duration)
	require.Empty(t, at.Language)
	require.Empty(t, at.Segments)
	require.Empty(t, at.Task)
}

func requireEmptyAudioTranslation(t *testing.T, at azopenai.AudioTranslation) {
	// Text is always filled out for

	require.Empty(t, at.Duration)
	require.Empty(t, at.Language)
	require.Empty(t, at.Segments)
	require.Empty(t, at.Task)
}
