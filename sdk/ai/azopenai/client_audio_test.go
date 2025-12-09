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

	"github.com/openai/openai-go/v3"
	"github.com/stretchr/testify/require"
)

func TestClient_GetAudioTranscription(t *testing.T) {
	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.Whisper.Endpoint)
	model := azureOpenAI.Whisper.Model

	// We're experiencing load issues on some of our shared test resources so we'll just spot check.
	t.Run(fmt.Sprintf("%s (%s)", openai.AudioResponseFormatText, "m4a"), func(t *testing.T) {
		transcriptResp, err := client.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
			Model:          openai.AudioModel(model),
			File:           getFile(t, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.m4a"),
			ResponseFormat: openai.AudioResponseFormatText,
			Language:       openai.String("en"),
			Temperature:    openai.Float(0.0),
		})
		require.Empty(t, transcriptResp)
		require.EqualError(t, err, "expected destination type of 'string' or '[]byte' for responses with content-type 'text/plain; charset=utf-8' that is not 'application/json'")
	})

	t.Run(fmt.Sprintf("%s (%s)", openai.AudioResponseFormatJSON, "mp3"), func(t *testing.T) {
		transcriptResp, err := client.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
			Model:          openai.AudioModel(model),
			File:           getFile(t, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.mp3"),
			ResponseFormat: openai.AudioResponseFormatJSON,
			Language:       openai.String("en"),
			Temperature:    openai.Float(0.0),
		})
		customRequireNoError(t, err)
		t.Logf("Transcription: %s", transcriptResp.Text)
		require.NotEmpty(t, transcriptResp)
	})
}

func TestClient_GetAudioTranslation(t *testing.T) {
	client := newStainlessTestClientWithAzureURL(t, azureOpenAI.Whisper.Endpoint)
	model := azureOpenAI.Whisper.Model

	resp, err := client.Audio.Translations.New(context.Background(), openai.AudioTranslationNewParams{
		Model:          openai.AudioModel(model),
		File:           getFile(t, "testdata/sampledata_audiofiles_myVoiceIsMyPassportVerifyMe01.m4a"),
		ResponseFormat: openai.AudioTranslationNewParamsResponseFormatVerboseJSON,
		Temperature:    openai.Float(0.0),
	})
	customRequireNoError(t, err)

	t.Logf("Translation: %s", resp.Text)
	require.NotEmpty(t, resp.Text)
}

// fakeFlacFile works around a problem with the Stainless client's use of .Name() on the
// passed in file and how it causes our test recordings to not match if the filename or
// path is randomized.
type fakeFlacFile struct {
	inner io.Reader
}

func (f *fakeFlacFile) Read(p []byte) (n int, err error) {
	return f.inner.Read(p)
}

func (f *fakeFlacFile) Name() string {
	return "audio.flac"
}

func TestClient_GetAudioSpeech(t *testing.T) {
	var tempFile *os.File

	// Generate some speech from text.
	{
		speechClient := newStainlessTestClientWithAzureURL(t, azureOpenAI.Speech.Endpoint)

		audioResp, err := speechClient.Audio.Speech.New(context.Background(), openai.AudioSpeechNewParams{
			Input:          "i am a computer",
			Voice:          openai.AudioSpeechNewParamsVoiceAlloy,
			ResponseFormat: openai.AudioSpeechNewParamsResponseFormatFLAC,
			Model:          openai.AudioModel(azureOpenAI.Speech.Model),
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

		// For test recordings, make sure we write the FLAC to a temp file with a consistent base name - the
		// Stainless API uses the filename of the file when it sends the request
		flacPath := filepath.Join(t.TempDir(), "audio.flac")
		require.NoError(t, err)

		writer, err := os.Create(flacPath)
		require.NoError(t, err)

		tempFile = writer

		_, err = tempFile.Write(audioBytes)
		require.NoError(t, err)

		_, err = tempFile.Seek(0, io.SeekStart)
		require.NoError(t, err)

		t.Cleanup(func() {
			_ = tempFile.Close()
		})
	}

	// as a simple check we'll now transcribe the audio file we just generated...
	transcriptClient := newStainlessTestClientWithAzureURL(t, azureOpenAI.Whisper.Endpoint)

	// now send _it_ back through the transcription API and see if we can get something useful.
	transcriptResp, err := transcriptClient.Audio.Transcriptions.New(context.Background(), openai.AudioTranscriptionNewParams{
		Model:          openai.AudioModel(azureOpenAI.Whisper.Model),
		File:           &fakeFlacFile{tempFile},
		ResponseFormat: openai.AudioResponseFormatVerboseJSON,
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
