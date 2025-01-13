//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azopenaiassistants

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

// EventReader streams events dynamically from an OpenAI endpoint.
type EventReader[T any] struct {
	reader           io.ReadCloser // Required for Closing
	bufioReader      *bufio.Reader
	zeroT            T
	currentEventType string

	//nolint:unused // the parameters are there for documentation
	unmarshal func(eventType string, jsonData []byte) (T, error)
}

func newEventReader[T any](r io.ReadCloser, unmarshalFn func(eventType string, jsonData []byte) (T, error)) *EventReader[T] {
	return &EventReader[T]{
		reader:      r,
		bufioReader: bufio.NewReader(r),
		unmarshal:   unmarshalFn,
	}
}

// Read reads the next event from the stream.
// Returns io.EOF when there are no further events.
func (er *EventReader[T]) Read() (T, error) {
	// https://html.spec.whatwg.org/multipage/server-sent-events.html

	for {
		line, err := er.bufioReader.ReadString('\n')

		if err != nil {
			if errors.Is(err, io.EOF) {
				return er.zeroT, errors.New("incomplete stream")
			}

			return er.zeroT, err
		}

		if line == "" || line[0] == ':' {
			continue
		}

		if strings.Contains(line, ":") { // If the line contains a U+003A COLON character (:), process the field
			tokens := strings.SplitN(line, ":", 2)
			tokens[0], tokens[1] = strings.TrimSpace(tokens[0]), strings.TrimSpace(tokens[1])

			switch tokens[0] {
			case "event":
				er.currentEventType = tokens[1]
			case "data": // return the deserialized JSON object
				if tokens[1] == "[DONE]" { // If data is [DONE], end of stream was reached
					return er.zeroT, io.EOF
				}
				return er.unmarshal(er.currentEventType, []byte(tokens[1]))
			default: // Any other event type is an unexpected
				return er.zeroT, errors.New("unexpected event type: " + tokens[0])
			}
			// Unreachable
		}
	}
}

// Close closes the EventReader and any applicable inner stream state.
func (er *EventReader[T]) Close() error {
	return er.reader.Close()
}
