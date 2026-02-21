// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package test

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/exported"
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
		setAzLogListener(func(azlog.Event, string) {})

		var messages []string

		for {
			select {
			case msg := <-messagesCh:
				messages = append(messages, msg)
			default:
				return messages
			}
		}
	}
}

// EnableStdoutLogging turns on logging to stdout for diagnostics.
func EnableStdoutLogging() {
	azlog.SetEvents(exported.EventAuth, exported.EventConn, exported.EventConsumer, exported.EventProducer)
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
	ClientID              string
	Cred                  azcore.TokenCredential
	EventHubName          string
	EventHubLinksOnlyName string
	EventHubNamespace     string
	GeoDRNamespace        string // optional: resource requires special setup to create, so it's manual
	GeoDRHubName          string // optional: resource requires special setup to create, so it's manual
	GeoDRStorageEndpoint  string // optional: resource requires special setup to create, so it's manual
	StorageEndpoint       string
	ResourceGroup         string
	SubscriptionID        string
	TenantID              string
}

func (c ConnectionParamsForTest) CS(t *testing.T) struct{ Primary, ListenOnly, SendOnly, Storage string } {
	if val, exists := os.LookupEnv("EVENTHUB_CONNECTION_STRING_LISTEN_ONLY"); exists && val == "" {
		// This happens if we're not in the TME subscription - the variable will just be set to an empty string
		// rather than not existing, altogether.
		t.Skip("Not in TME, skipping connection string tests")
	}

	envVars := mustGetEnvironmentVars(t, []string{
		"EVENTHUB_CONNECTION_STRING_LISTEN_ONLY",
		"EVENTHUB_CONNECTION_STRING_SEND_ONLY",
		"EVENTHUB_CONNECTION_STRING",
	})

	return struct{ Primary, ListenOnly, SendOnly, Storage string }{
		envVars["EVENTHUB_CONNECTION_STRING"],
		envVars["EVENTHUB_CONNECTION_STRING_LISTEN_ONLY"],
		envVars["EVENTHUB_CONNECTION_STRING_SEND_ONLY"],
		envVars["EVENTHUB_CONNECTION_STRING_SEND_ONLY"],
	}
}

func GetConnectionParamsForTest(t *testing.T) ConnectionParamsForTest {
	if _, err := os.Stat("../.env"); err == nil {
		_ = godotenv.Load("../.env")
	} else {
		_ = godotenv.Load()
	}

	envVars := mustGetEnvironmentVars(t, []string{
		"AZURE_SUBSCRIPTION_ID",
		"CHECKPOINTSTORE_STORAGE_ENDPOINT",
		"EVENTHUB_NAME",
		"EVENTHUB_NAMESPACE",
		"EVENTHUB_LINKSONLY_NAME",
		"RESOURCE_GROUP",
	})

	cred, err := credential.New(nil)
	require.NoError(t, err)

	return ConnectionParamsForTest{
		Cred:                  cred,
		EventHubName:          envVars["EVENTHUB_NAME"],
		EventHubLinksOnlyName: envVars["EVENTHUB_LINKSONLY_NAME"],
		EventHubNamespace:     envVars["EVENTHUB_NAMESPACE"],
		GeoDRNamespace:        os.Getenv("EVENTHUBS_GEODR_NAMESPACE"),
		GeoDRHubName:          os.Getenv("EVENTHUBS_GEODR_HUBNAME"),
		GeoDRStorageEndpoint:  os.Getenv("EVENTHUBS_GEODR_CHECKPOINTSTORE_STORAGE_ENDPOINT"),
		ResourceGroup:         envVars["RESOURCE_GROUP"],
		StorageEndpoint:       envVars["CHECKPOINTSTORE_STORAGE_ENDPOINT"],
		SubscriptionID:        envVars["AZURE_SUBSCRIPTION_ID"],
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

func RequireClose[T interface {
	Close(ctx context.Context) error
}](t *testing.T, closeable T) {
	require.NoErrorf(t, closeable.Close(context.Background()), "%T closes cleanly", closeable)
}

func RequireNSClose(t *testing.T, closeable interface {
	Close(ctx context.Context, permanent bool) error
}) {
	require.NoError(t, closeable.Close(context.Background(), true))
}

// RequireContextHasDefaultTimeout checks that the context has a deadline set, and that it's
// using the right timeout.
// NOTE: There's some wiggle room since some time will expire before this is called.
func RequireContextHasDefaultTimeout(t *testing.T, ctx context.Context, timeout time.Duration) {
	tm, hasDeadline := ctx.Deadline()

	require.True(t, hasDeadline, "deadline must exist, we always set an operation timeout")
	duration := time.Until(tm)

	require.Greater(t, duration, time.Duration(0))
	require.LessOrEqual(t, duration, timeout)
}

func URLJoinPaths(base string, subPath string) string {
	slash := "/"

	if strings.HasSuffix(base, "/") {
		slash = ""
	}

	return base + slash + subPath
}
