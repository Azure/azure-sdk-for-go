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
	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/atom"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

var (
	letterRunes = []rune("abcdefghijklmnopqrstuvwxyz123456789")
)

// these are created by the test-resources.bicep template - they're useful for tests where we don't need to guaranteee any state, just existence, like our connectivity/recovery tests.
const (
	BuiltInTestQueue             = "testQueue"
	BuildInTestQueueWithSessions = "testQueueWithSessions"
)

func init() {
	addSwappableLogger()
}

// RandomString generates a random string with prefix
func RandomString(prefix string, length int) string {
	rand := rand.New(rand.NewSource(time.Now().Unix()))

	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return prefix + string(b)
}

func GetConnectionString(t *testing.T) string {
	return getEnvOrSkipTest(t, "SERVICEBUS_CONNECTION_STRING")
}

func GetConnectionStringForPremiumSB(t *testing.T) string {
	return getEnvOrSkipTest(t, "SERVICEBUS_CONNECTION_STRING_PREMIUM")
}

func GetConnectionStringWithoutManagePerms(t *testing.T) string {
	return getEnvOrSkipTest(t, "SERVICEBUS_CONNECTION_STRING_NO_MANAGE")
}

func GetConnectionStringSendOnly(t *testing.T) string {
	return getEnvOrSkipTest(t, "SERVICEBUS_CONNECTION_STRING_SEND_ONLY")
}

func GetConnectionStringListenOnly(t *testing.T) string {
	return getEnvOrSkipTest(t, "SERVICEBUS_CONNECTION_STRING_LISTEN_ONLY")
}

func GetIdentityVars(t *testing.T) *struct {
	TenantID string
	ClientID string
	Secret   string
	Endpoint string
} {
	runningLiveTest := GetConnectionString(t) != ""

	if !runningLiveTest {
		return nil
	}

	return &struct {
		TenantID string
		ClientID string
		Secret   string
		Endpoint string
	}{
		TenantID: getEnvOrSkipTest(t, "AZURE_TENANT_ID"),
		ClientID: getEnvOrSkipTest(t, "AZURE_CLIENT_ID"),
		Endpoint: getEnvOrSkipTest(t, "SERVICEBUS_ENDPOINT"),
		Secret:   getEnvOrSkipTest(t, "AZURE_CLIENT_SECRET"),
	}
}

func getEnvOrSkipTest(t *testing.T, name string) string {
	cs := os.Getenv(name)

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
func CaptureLogsForTest(echo bool) func() []string {
	return CaptureLogsForTestWithChannel(nil, echo)
}

func CaptureLogsForTestWithChannel(messagesCh chan string, echo bool) func() []string {
	if messagesCh == nil {
		messagesCh = make(chan string, 10000)
	}

	LoggingChannelValue.Store(&messagesCh)

	return func() []string {
		if messagesCh == nil {
			// already been closed, probably manually.
			return nil
		}

		var messages []string

	Loop:
		for {
			select {
			case msg := <-messagesCh:
				if echo {
					log.Printf("%s", msg)
				}
				messages = append(messages, msg)
			default:
				break Loop
			}
		}

		return messages
	}
}

// EnableStdoutLogging turns on logging to stdout for diagnostics.
func EnableStdoutLogging(t *testing.T) {
	ch := make(chan string, 10000)
	cleanupLogs := CaptureLogsForTestWithChannel(ch, true)
	ctx, cancel := context.WithCancel(context.Background())

	t.Cleanup(func() {
		cancel()
	})

	go func() {
		<-ctx.Done()
		_ = cleanupLogs()
	}()
}

func RequireClose(t *testing.T, closeable interface {
	Close(ctx context.Context) error
}) {
	err := closeable.Close(context.Background())
	require.NoError(t, err)
}

func RequireLinksClose(t *testing.T, closeable interface {
	Close(ctx context.Context, permanent bool) error
}) {
	err := closeable.Close(context.Background(), true)
	require.NoError(t, err)
}

func RequireNSClose(t *testing.T, closeable interface {
	Close(permanent bool) error
}) {
	err := closeable.Close(true)
	require.NoError(t, err)
}

func MustAMQPUUID() amqp.UUID {
	id, err := uuid.New()

	if err != nil {
		panic(err)
	}

	return amqp.UUID(id)
}
