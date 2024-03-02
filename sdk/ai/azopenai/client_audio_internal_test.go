//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"bytes"
	"context"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/streaming"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/stretchr/testify/require"
)

func TestSetMultipartFormData(t *testing.T) {
	t.Run("getAudioTranscriptionInternalOptions", func(t *testing.T) {
		req, err := runtime.NewRequest(context.Background(), "POST", "http://localhost")
		require.NoError(t, err)

		err = setMultipartFormData(req, streaming.NopCloser(bytes.NewReader([]byte{1, 2, 3})), getAudioTranscriptionInternalOptions{
			Language:       to.Ptr("en"),
			DeploymentName: to.Ptr("hello"),
			Prompt:         to.Ptr("my prompt"),
			ResponseFormat: to.Ptr(AudioTranscriptionFormatJSON),
			Temperature:    to.Ptr[float32](1.0),
		})
		require.NoError(t, err)

		_, params, err := mime.ParseMediaType(req.Raw().Header["Content-Type"][0])
		require.NoError(t, err)

		parts := getParts(t, req.Body(), params["boundary"])

		require.Equal(t, []kv{
			{Key: "file", Value: "\x01\x02\x03"},
			{Key: "language", Value: "en"},
			{Key: "model", Value: "hello"},
			{Key: "prompt", Value: "my prompt"},
			{Key: "response_format", Value: string(AudioTranscriptionFormatJSON)},
			{Key: "temperature", Value: "1.000000"},
		}, parts)
	})

	t.Run("getAudioTranslationInternalOptions", func(t *testing.T) {
		req, err := runtime.NewRequest(context.Background(), "POST", "http://localhost")
		require.NoError(t, err)

		err = setMultipartFormData(req, streaming.NopCloser(bytes.NewReader([]byte{1, 2, 3})), getAudioTranslationInternalOptions{
			DeploymentName: to.Ptr("hello"),
			Prompt:         to.Ptr("my prompt"),
			ResponseFormat: to.Ptr(AudioTranslationFormatJSON),
			Temperature:    to.Ptr[float32](1.0),
		})
		require.NoError(t, err)

		_, params, err := mime.ParseMediaType(req.Raw().Header["Content-Type"][0])
		require.NoError(t, err)

		parts := getParts(t, req.Body(), params["boundary"])

		require.Equal(t, []kv{
			{Key: "file", Value: "\x01\x02\x03"},
			{Key: "model", Value: "hello"},
			{Key: "prompt", Value: "my prompt"},
			{Key: "response_format", Value: string(AudioTranscriptionFormatJSON)},
			{Key: "temperature", Value: "1.000000"},
		}, parts)
	})
}

func TestWriteField(t *testing.T) {
	buff := &bytes.Buffer{}
	writer := multipart.NewWriter(buff)

	// value is nil - we skip the field
	err := writeField(writer, "mynil", (*string)(nil))
	require.NoError(t, err)

	err = writeField(writer, "mystring", to.Ptr("hello world"))
	require.NoError(t, err)

	err = writeField(writer, "myfloat", to.Ptr[float32](1.0))
	require.NoError(t, err)

	err = writeField(writer, "myformatsrt", to.Ptr(AudioTranscriptionFormatSrt))
	require.NoError(t, err)

	err = writer.Close()
	require.NoError(t, err)

	parts := getParts(t, buff, writer.Boundary())
	require.Equal(t, []kv{
		// ...nil field was skipped...
		{Key: "mystring", Value: "hello world"},
		{Key: "myfloat", Value: "1.000000"},
		{Key: "myformatsrt", Value: string(AudioTranscriptionFormatSrt)},
	}, parts)
}

type kv struct {
	Key   string
	Value string
}

func getParts(t *testing.T, buff io.Reader, boundary string) []kv {
	reader := multipart.NewReader(buff, boundary)

	var parts []kv

	for {
		part, err := reader.NextPart()

		if errors.Is(err, io.EOF) {
			break
		}
		require.NoError(t, err)

		text, err := io.ReadAll(part)
		require.NoError(t, err)

		parts = append(parts, kv{Key: part.FormName(), Value: string(text)})
	}

	return parts
}
