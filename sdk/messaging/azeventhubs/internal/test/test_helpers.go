// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package test

import (
	"fmt"
	"log"
	"sync"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
)

// CaptureLogsForTest adds a logging listener which captures messages to an
// internal channel.
// Returns a function that ends log capturing and returns any captured messages.
// It's safe to call endCapture() multiple times, so a simple call pattern is:
//
//   endCapture := CaptureLogsForTest()
//   defer endCapture()				// ensure cleanup in case of test assert failures
//
//   /* some test code */
//
//   messages := endCapture()
//   /* do inspection of log messages */
//
func CaptureLogsForTest() func() []string {
	messagesCh := make(chan string, 10000)
	return CaptureLogsForTestWithChannel(messagesCh)
}

func CaptureLogsForTestWithChannel(messagesCh chan string) func() []string {
	setAzLogListener(func(e azlog.Event, s string) {
		messagesCh <- fmt.Sprintf("[%s] %s", e, s)
	})

	return func() []string {
		if messagesCh == nil {
			// already been closed, probably manually.
			return nil
		}

		setAzLogListener(nil)
		close(messagesCh)

		var messages []string

		for msg := range messagesCh {
			messages = append(messages, msg)
		}

		messagesCh = nil
		return messages
	}
}

// EnableStdoutLogging turns on logging to stdout for diagnostics.
func EnableStdoutLogging() {
	setAzLogListener(func(e azlog.Event, s string) {
		log.Printf("%s %s", e, s)
	})
}

var logMu sync.Mutex

func setAzLogListener(listener func(e azlog.Event, s string)) {
	logMu.Lock()
	defer logMu.Unlock()
	azlog.SetListener(listener)
}
