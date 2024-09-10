//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiextensions_test

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
	"github.com/stretchr/testify/require"
)

func TestClient_GetAudioTranscription(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}

	client := newStainlessTestClient(t, azureOpenAI.Whisper.Endpoint)
	model := azureOpenAI.Whisper.Model

	// We're experiencing load issues on some of our shared test resources so we'll just spot check.
	t.Run(fmt.Sprintf("%s (%s)", openai.AudioTranscriptionNewParamsResponseFormatText, "m4a"), func(t *testing.T) {
		// TODO: BUG: I think. I'm not quite sure how to request any format other than JSON because the bare formats
		// cause a deserialization error in the Stainless client.
		//
		// transcriptResp, err := client.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
		// 	Model:          openai.F(openai.AudioTranscriptionNewParamsModel(model)),
		// 	File:           openai.F(getFile(t, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.m4a")),
		// 	ResponseFormat: openai.F(openai.AudioTranscriptionNewParamsResponseFormatText),
		// 	Language:       openai.String("en"),
		// 	Temperature:    openai.Float(0.0),
		// })
		// require.Empty(t, transcriptResp)
		// require.EqualError(t, err, "expected destination type of 'string' or '[]byte' for responses with content-type that is not 'application/json'")

		var text *string

		transcriptResp, err := client.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
			Model:          openai.F(openai.AudioModel(model)),
			File:           openai.F(getFile(t, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.m4a")),
			ResponseFormat: openai.F(openai.AudioTranscriptionNewParamsResponseFormatText),
			Language:       openai.String("en"),
			Temperature:    openai.Float(0.0),
		}, option.WithResponseBodyInto(&text))
		require.Empty(t, transcriptResp)
		require.NoError(t, err)
		require.NotEmpty(t, *text)
	})

	t.Run(fmt.Sprintf("%s (%s)", openai.AudioTranscriptionNewParamsResponseFormatJSON, "mp3"), func(t *testing.T) {
		if recording.GetRecordMode() == recording.PlaybackMode {
			t.Skip("TODO: sanitization break: needs to be looked at")
		}

		transcriptResp, err := client.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
			Model:          openai.F(openai.AudioModel(model)),
			File:           openai.F(getFile(t, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.mp3")),
			ResponseFormat: openai.F(openai.AudioTranscriptionNewParamsResponseFormatVerboseJSON),
			Language:       openai.String("en"),
			Temperature:    openai.Float(0.0),
		})
		customRequireNoError(t, err, true)
		t.Logf("Transcription: %s", transcriptResp.Text)
		require.NotEmpty(t, transcriptResp)
	})
}

func TestClient_GetAudioTranslation(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}

	client := newStainlessTestClient(t, azureOpenAI.Whisper.Endpoint)
	model := azureOpenAI.Whisper.Model

	resp, err := client.Audio.Translations.New(context.Background(), openai.AudioTranslationNewParams{
		Model: openai.F(openai.AudioModel(model)),
		File:  openai.F(getFile(t, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.m4a")),
		// TODO: no specific enumeration for Translations format?
		ResponseFormat: openai.F(string(openai.AudioTranscriptionNewParamsResponseFormatVerboseJSON)),
		Temperature:    openai.Float(0.0),
	})
	customRequireNoError(t, err, true)

	t.Logf("Translation: %s", resp.Text)
	require.NotEmpty(t, resp.Text)
}

func TestClient_GetAudioSpeech(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}

	var tempFile *os.File

	// Generate some speech from text.
	{
		speechClient := newStainlessTestClient(t, azureOpenAI.Speech.Endpoint)

		audioResp, err := speechClient.Audio.Speech.New(context.Background(), openai.AudioSpeechNewParams{
			Input:          openai.String("i am a computer"),
			Voice:          openai.F(openai.AudioSpeechNewParamsVoiceAlloy),
			ResponseFormat: openai.F(openai.AudioSpeechNewParamsResponseFormatFLAC),
			Model:          openai.F(openai.AudioModel(azureOpenAI.Speech.Model)),
		})
		require.NoError(t, err)

		defer func() {
			err := audioResp.Body.Close()
			require.NoError(t, err)
		}()

		audioBytes, err := io.ReadAll(audioResp.Body)
		require.NoError(t, err)

		require.NotEmpty(t, audioBytes)
		require.Equal(t, "fLaC", string(audioBytes[0:4]))

		// write the FLAC to a temp file - the Stainless API uses the filename of the file
		// when it sends the request.
		tempFile, err = os.CreateTemp("", "audio*.flac")
		require.NoError(t, err)
		defer tempFile.Close()

		_, err = tempFile.Write(audioBytes)
		require.NoError(t, err)

		_, err = tempFile.Seek(0, io.SeekStart)
		require.NoError(t, err)
	}

	// as a simple check we'll now transcribe the audio file we just generated...
	transcriptClient := newStainlessTestClient(t, azureOpenAI.Whisper.Endpoint)

	// now send _it_ back through the transcription API and see if we can get something useful.
	transcriptResp, err := transcriptClient.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
		Model:          openai.F(openai.AudioModel(azureOpenAI.Whisper.Model)),
		File:           openai.F[io.Reader](tempFile),
		ResponseFormat: openai.F(openai.AudioTranscriptionNewParamsResponseFormatVerboseJSON),
		Language:       openai.String("en"),
		Temperature:    openai.Float(0.0),
	})
	require.NoError(t, err)

	// it occasionally comes back with different punctuation or makes a complete sentence but
	// the major words always come through.
	require.Contains(t, transcriptResp.Text, "computer")
}

func getFile(t *testing.T, path string) io.Reader {
	file, err := os.Open(path)
	require.NoError(t, err)

	t.Cleanup(func() {
		err := file.Close()
		require.NoError(t, err)
	})

	return file
}
