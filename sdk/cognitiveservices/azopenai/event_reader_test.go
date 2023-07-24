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
		"invaliddata: {\u0022name\u0022:\u0022chatcmpl-7Z4kUpXX6HN85cWY28IXM4EwemLU3\u0022,\u0022object\u0022:\u0022chat.completion.chunk\u0022,\u0022created\u0022:1688594090,\u0022model\u0022:\u0022gpt-4-0613\u0022,\u0022choices\u0022:[{\u0022index\u0022:0,\u0022delta\u0022:{\u0022role\u0022:\u0022assistant\u0022,\u0022content\u0022:\u0022\u0022},\u0022finish_reason\u0022:null}]}\n\n",
	}

	text := strings.NewReader(strings.Join(data, "\n"))
	eventReader := newEventReader[ChatCompletions](text)

	firstEvent, err := eventReader.Read()
	require.Empty(t, firstEvent)
	require.EqualError(t, err, "Unexpected event type: invaliddata")
}

type badReader struct{}

func (br badReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrClosedPipe
}

func TestEventReader_BadReader(t *testing.T) {
	eventReader := newEventReader[ChatCompletions](badReader{})

	firstEvent, err := eventReader.Read()
	require.Empty(t, firstEvent)
	require.ErrorIs(t, io.ErrClosedPipe, err)
}
