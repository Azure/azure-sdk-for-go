// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package test

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/stretchr/testify/require"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")
)

func init() {
	addSwappableLogger()
	rand.Seed(time.Now().Unix())
}

// RandomString generates a random string with prefix
func RandomString(prefix string, length int) string {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return prefix + string(b)
}

func GetConnectionString(t *testing.T) string {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING")

	if cs == "" {
		t.Skip()
	}

	return cs
}

func GetConnectionStringForPremiumSB(t *testing.T) string {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING_PREMIUM")

	if cs == "" {
		t.Skip()
	}

	return cs
}

func GetConnectionStringWithoutManagePerms(t *testing.T) string {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING_NO_MANAGE")

	if cs == "" {
		t.Skip()
	}

	return cs
}

func CreateExpiringQueue(t *testing.T, qd *atom.QueueDescription) (string, func()) {
	cs := GetConnectionString(t)
	em, err := atom.NewEntityManagerWithConnectionString(cs, "", nil)
	require.NoError(t, err)

	queueName := RandomString("queue", 5)

	if qd == nil {
		qd = &atom.QueueDescription{}
	}

	qd.AutoDeleteOnIdle = to.Ptr("PT5M")

	env := atom.WrapWithQueueEnvelope(qd, em.TokenProvider())

	var qe *atom.QueueEnvelope
	resp, err := em.Put(context.Background(), queueName, env, &qe, nil)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusCreated, resp.StatusCode)

	return queueName, func() {
		_, err := em.Delete(context.Background(), queueName)
		require.NoError(t, err)
	}
}

var LoggingChannelValue atomic.Value

func addSwappableLogger() {
	azlog.SetListener(func(e azlog.Event, s string) {
		ch, ok := LoggingChannelValue.Load().(*chan string)

		if !ok || ch == nil {
			return
		}

		select {
		case *ch <- fmt.Sprintf("[%s] %s", e, s):
		default:
		}
	})
}

// CaptureLogsForTest adds a logging listener which captures messages to an
// internal channel.
// Returns a function that ends log capturing and returns any captured messages.
// It's safe to call endCapture() multiple times, so a simple call pattern is:
//
//	endCapture := CaptureLogsForTest()
//	defer endCapture()				// ensure cleanup in case of test assert failures
//
//	/* some test code */
//
//	messages := endCapture()
//	/* do inspection of log messages */
func CaptureLogsForTest() func() []string {
	return CaptureLogsForTestWithChannel(nil)
}

func CaptureLogsForTestWithChannel(messagesCh chan string) func() []string {
	if messagesCh == nil {
		messagesCh = make(chan string, 10000)
	}

	LoggingChannelValue.Store(&messagesCh)

	return func() []string {
		if messagesCh == nil {
			// already been closed, probably manually.
			return nil
		}

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
func EnableStdoutLogging() func() {
	ch := make(chan string, 10000)
	cleanupLogs := CaptureLogsForTestWithChannel(ch)

	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		for msg := range ch {
			log.Printf("%s", msg)
		}
	}()

	return func() {
		_ = cleanupLogs()
		<-doneCh
	}
}
