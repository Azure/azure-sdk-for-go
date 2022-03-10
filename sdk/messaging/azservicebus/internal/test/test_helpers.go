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
	"testing"
	"time"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/stretchr/testify/require"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")
)

func init() {
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

func GetConnectionStringWithoutManagePerms(t *testing.T) string {
	cs := os.Getenv("SERVICEBUS_CONNECTION_STRING_NO_MANAGE")

	if cs == "" {
		t.Skip()
	}

	return cs
}

func CreateExpiringQueue(t *testing.T, qd *atom.QueueDescription) (string, func()) {
	cs := GetConnectionString(t)
	em, err := atom.NewEntityManagerWithConnectionString(cs, "")
	require.NoError(t, err)

	queueName := RandomString("queue", 5)

	if qd == nil {
		qd = &atom.QueueDescription{}
	}

	deleteAfter := 5 * time.Minute
	qd.AutoDeleteOnIdle = utils.DurationToStringPtr(&deleteAfter)

	env, _ := atom.WrapWithQueueEnvelope(qd, em.TokenProvider())

	var qe *atom.QueueEnvelope
	resp, err := em.Put(context.Background(), queueName, env, &qe)
	require.NoError(t, err)
	require.EqualValues(t, http.StatusCreated, resp.StatusCode)

	return queueName, func() {
		_, err := em.Delete(context.Background(), queueName)
		require.NoError(t, err)
	}
}

// CaptureLogsForTest adds a listener that appends to your `logMessages`
// slice.
// Returns a function that should be called to unregister our listener.
func CaptureLogsForTest(logMessages *[]string) func() {
	azlog.SetListener(func(e azlog.Event, s string) {
		*logMessages = append(*logMessages, fmt.Sprintf("[%s] %s", e, s))
	})
	return func() { azlog.SetListener(nil) }
}

// EnableStdoutLogging turns on logging to stdout for diagnostics.
func EnableStdoutLogging() {
	azlog.SetListener(func(e azlog.Event, s string) {
		log.Printf("%s %s", e, s)
	})
}
