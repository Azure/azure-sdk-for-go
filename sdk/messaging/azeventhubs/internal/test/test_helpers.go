// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package test

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/conn"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/require"
)

// CaptureLogsForTest adds a logging listener which captures messages to an
// internal channel.
// Returns a function that ends log capturing and returns any captured messages.
// It's safe to call endCapture() multiple times, so a simple call pattern is:
//
//	endCapture := CaptureLogsForTest()
//	defer endCapture()  // ensure cleanup in case of test assert failures
//
//	/* some test code */
//
//	messages := endCapture()
//	/* do inspection of log messages */
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

// RandomString generates a random string with prefix
func RandomString(prefix string, length int) string {
	b := make([]byte, length)
	_, err := rand.Read(b)

	if err != nil {
		panic(err)
	}

	return prefix + hex.EncodeToString(b)
}

type ConnectionParamsForTest struct {
	ConnectionString  string
	EventHubName      string
	EventHubNamespace string

	StorageConnectionString string
}

func GetConnectionParamsForTest(t *testing.T) ConnectionParamsForTest {
	_ = godotenv.Load()

	envVars := mustGetEnvironmentVars(t, []string{
		"EVENTHUB_CONNECTION_STRING",
		"EVENTHUB_NAME",
		"CHECKPOINTSTORE_STORAGE_CONNECTION_STRING",
	})

	parsedConn, err := conn.ParsedConnectionFromStr(envVars["EVENTHUB_CONNECTION_STRING"])
	require.NoError(t, err)

	return ConnectionParamsForTest{
		ConnectionString:        envVars["EVENTHUB_CONNECTION_STRING"],
		EventHubName:            envVars["EVENTHUB_NAME"],
		EventHubNamespace:       parsedConn.Namespace,
		StorageConnectionString: envVars["CHECKPOINTSTORE_STORAGE_CONNECTION_STRING"],
	}
}

func mustGetEnvironmentVars(t *testing.T, names []string) map[string]string {
	var missingVars []string
	envVars := map[string]string{}

	for _, name := range names {
		val := os.Getenv(name)

		if val == "" {
			missingVars = append(missingVars, name)
			continue
		}

		envVars[name] = val
	}

	if len(missingVars) > 0 {
		t.Skipf("Missing env vars for live test: %s. Skipping...", strings.Join(missingVars, ","))
		return nil
	}

	return envVars
}
