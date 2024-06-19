//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants

import (
	"encoding/json"
	"fmt"
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
	eventReader := newEventReader[any](io.NopCloser(text), func(eventType string, jsonData []byte) (any, error) {
		return string(jsonData), nil
	})

	firstEvent, err := eventReader.Read()
	require.Empty(t, firstEvent)
	require.EqualError(t, err, "unexpected event type: invaliddata")
}

type badReader struct{}

func (br badReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

func TestEventReader_BadReader(t *testing.T) {
	eventReader := newEventReader[any](io.NopCloser(badReader{}), func(eventType string, jsonData []byte) (any, error) {
		return string(jsonData), nil
	})
	defer eventReader.Close()

	firstEvent, err := eventReader.Read()
	require.Empty(t, firstEvent)
	require.ErrorIs(t, io.ErrClosedPipe, err)
}

func TestEventReader_StreamIsClosedBeforeDone(t *testing.T) {
	buff := strings.NewReader("data: {}\n")

	eventReader := newEventReader(io.NopCloser(buff), func(eventType string, jsonData []byte) (any, error) {
		return string(jsonData), nil
	})

	evt, err := eventReader.Read()
	require.Equal(t, "{}", evt)
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

	eventReader := newEventReader[eventReaderTestModel](io.NopCloser(buff), func(eventType string, jsonData []byte) (eventReaderTestModel, error) {
		var v *eventReaderTestModel
		err := json.Unmarshal(jsonData, &v)
		return *v, err
	})

	evt, err := eventReader.Read()
	require.NoError(t, err)
	require.Equal(t, "with-spaces", *evt.Choices[0].Delta.Content)

	evt, err = eventReader.Read()
	require.NoError(t, err)
	require.NotEmpty(t, evt)
	require.Equal(t, "without-spaces", *evt.Choices[0].Delta.Content)
}

func TestEventReader_CanReadHugeChunk(t *testing.T) {
	// Ran into this with a customer that gets _huge_ chunks of text in their stream:
	// https://github.com/Azure/azure-sdk-for-go/pull/22646
	// bufio.Scanner has a limitation of 64k (which is huge, but not big enough).

	bigBytes := make([]byte, 64*1024+1)

	for i := 0; i < len(bigBytes); i++ {
		bigBytes[i] = 'A'
	}

	buff := strings.NewReader(
		fmt.Sprintf("data: {\"name\":\"chatcmpl-7Z4kUpXX6HN85cWY28IXM4EwemLU3\",\"object\":\"chat.completion.chunk\",\"created\":1688594090,\"model\":\"gpt-4-0613\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"%s\"},\"finish_reason\":null}]}\n", string(bigBytes)) + fmt.Sprintf("data: {\"name\":\"chatcmpl-7Z4kUpXX6HN85cWY28IXM4EwemLU3\",\"object\":\"chat.completion.chunk\",\"created\":1688594090,\"model\":\"gpt-4-0613\",\"choices\":[{\"index\":0,\"delta\":{\"role\":\"assistant\",\"content\":\"%s\"},\"finish_reason\":null}]}\n", "small message"),
	)

	eventReader := newEventReader[eventReaderTestModel](io.NopCloser(buff), func(eventType string, jsonData []byte) (eventReaderTestModel, error) {
		var v *eventReaderTestModel
		err := json.Unmarshal(jsonData, &v)
		return *v, err
	})

	evt, err := eventReader.Read()
	require.NoError(t, err)
	require.Equal(t, string(bigBytes), *evt.Choices[0].Delta.Content)

	evt, err = eventReader.Read()
	require.NoError(t, err)
	require.Equal(t, "small message", *evt.Choices[0].Delta.Content)
}

type eventReaderTestModel struct {
	Choices []struct {
		Delta struct {
			Content *string
		}
	}
}
