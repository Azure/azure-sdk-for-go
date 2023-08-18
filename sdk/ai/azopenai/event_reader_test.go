//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenai

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEventReader_InvalidType(t *testing.T) {
	data := []string{
		"invaliddata: {\"name\":\"chatcmpl-7Z4kUpXX6HN85cWY28IXM4EwemLU3\",\"object\":\"chat.completion.chunk\",\"created\":1688594090,\"model\":\"gpt-4-0613\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"\"},\"finish_reason\":null}]}\n\n",
	}

	text := strings.NewReader(strings.Join(data, "\n"))
	eventReader := newEventReader[ChatCompletions](io.NopCloser(text))

	firstEvent, err := eventReader.Read()
	require.Empty(t, firstEvent)
	require.EqualError(t, err, "unexpected event type: invaliddata")
}

type badReader struct{}

func (br badReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

func TestEventReader_BadReader(t *testing.T) {
	eventReader := newEventReader[ChatCompletions](io.NopCloser(badReader{}))
	defer eventReader.Close()

	firstEvent, err := eventReader.Read()
	require.Empty(t, firstEvent)
	require.ErrorIs(t, io.ErrClosedPipe, err)
}

func TestEventReader_StreamIsClosedBeforeDone(t *testing.T) {
	buff := strings.NewReader("data: {}")

	eventReader := newEventReader[ChatCompletions](io.NopCloser(buff))

	evt, err := eventReader.Read()
	require.Empty(t, evt)
	require.NoError(t, err)

	evt, err = eventReader.Read()
	require.Empty(t, evt)
	require.EqualError(t, err, "incomplete stream")
}

func TestEventReader_SpacesAroundAreas(t *testing.T) {
	buff := strings.NewReader(
		// spaces between data
		"data: {\"name\":\"chatcmpl-7Z4kUpXX6HN85cWY28IXM4EwemLU3\",\"object\":\"chat.completion.chunk\",\"created\":1688594090,\"model\":\"gpt-4-0613\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"with-spaces\"},\"finish_reason\":null}]}\n" +
			// no spaces
			"data:{\"name\":\"chatcmpl-7Z4kUpXX6HN85cWY28IXM4EwemLU3\",\"object\":\"chat.completion.chunk\",\"created\":1688594090,\"model\":\"gpt-4-0613\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"without-spaces\"},\"finish_reason\":null}]}\n",
	)

	eventReader := newEventReader[ChatCompletions](io.NopCloser(buff))

	evt, err := eventReader.Read()
	require.NoError(t, err)
	require.Equal(t, "with-spaces", *evt.Choices[0].Delta.Content)

	evt, err = eventReader.Read()
	require.NoError(t, err)
	require.NotEmpty(t, evt)
	require.Equal(t, "without-spaces", *evt.Choices[0].Delta.Content)
}
