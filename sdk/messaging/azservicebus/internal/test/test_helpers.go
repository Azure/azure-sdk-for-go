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
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/test/credential"
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

type IdentityVars struct {
	Endpoint        string
	PremiumEndpoint string
	Cred            azcore.TokenCredential
}

func GetIdentityVars(t *testing.T) *IdentityVars {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skipf("Skipping live test when in recording.PlaybackMode")
		return nil
	}

	tokenCred, err := credential.New(nil)
	require.NoError(t, err)

	envVars := MustGetEnvVars([]EnvKey{EnvKeyEndpoint, EnvKeyEndpointPremium})

	return &IdentityVars{
		Endpoint:        envVars[EnvKeyEndpoint],
		PremiumEndpoint: envVars[EnvKeyEndpointPremium],
		Cred:            tokenCred,
	}
}

type EnvKey string

const (
	EnvKeyEndpoint                   EnvKey = "SERVICEBUS_ENDPOINT"
	EnvKeyEndpointPremium            EnvKey = "SERVICEBUS_ENDPOINT_PREMIUM"
	EnvKeyConnectionString           EnvKey = "SERVICEBUS_CONNECTION_STRING"
	EnvKeyConnectionStringPremium    EnvKey = "SERVICEBUS_CONNECTION_STRING_PREMIUM"
	EnvKeyConnectionStringNoManage   EnvKey = "SERVICEBUS_CONNECTION_STRING_NO_MANAGE"
	EnvKeyConnectionStringSendOnly   EnvKey = "SERVICEBUS_CONNECTION_STRING_SEND_ONLY"
	EnvKeyConnectionStringListenOnly EnvKey = "SERVICEBUS_CONNECTION_STRING_LISTEN_ONLY"
)

func MustGetEnvVars[KeyT ~string](keys []KeyT) map[KeyT]string {
	m := map[KeyT]string{}
	var missing []string

	for _, k := range keys {
		v, exists := os.LookupEnv(string(k))

		if !exists {
			missing = append(missing, string(k))
		}

		m[k] = v
	}

	if len(missing) != 0 {
		sort.Strings(missing)
		log.Fatalf("Required env variables are missing: %#v", missing)
	}

	return m
}

func GetConnectionString(t *testing.T, name EnvKey) string {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skipf("Skipping live test when in recording.PlaybackMode")
		return ""
	}

	val, exists := os.LookupEnv(string(name))

	if exists && val == "" {
		// This happens if we're not in the TME subscription - the variable will just be set to an empty string
		// rather than not existing, altogether.
		t.Skip("Not in TME, skipping connection string tests")
	}

	return val
}

func MustGetEnvVar(t *testing.T, name EnvKey) string {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skipf("Skipping live test when in recording.PlaybackMode")
		return ""
	}

	cs := os.Getenv(string(name))
	require.NotEmptyf(t, cs, "Env var %s should be set", name)

	return cs
}

type NewClientArgs[OptionsT any, ClientT any] struct {
	NewClientFromConnectionString func(connectionString string, options *OptionsT) (*ClientT, error)
	NewClient                     func(endpoint string, cred azcore.TokenCredential, options *OptionsT) (*ClientT, error)
}

type NewClientOptions[OptionsT any] struct {
	ClientOptions       *OptionsT
	UseConnectionString bool
	UsePremium          bool
}

func NewClient[OptionsT any, ClientT any](t *testing.T, args NewClientArgs[OptionsT, ClientT], options *NewClientOptions[OptionsT]) *ClientT {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skipf("Skipping live test when in recording.PlaybackMode")
		return nil
	}

	if options == nil {
		options = &NewClientOptions[OptionsT]{}
	}

	if options.UseConnectionString {
		cs := GetConnectionString(t, EnvKeyConnectionString)
		env := MustGetEnvVars([]EnvKey{EnvKeyConnectionString, EnvKeyConnectionStringPremium})

		if options.UsePremium {
			cs = env[EnvKeyConnectionStringPremium]
		}

		client, err := args.NewClientFromConnectionString(cs, options.ClientOptions)
		require.NoError(t, err)

		return client
	}

	iv := GetIdentityVars(t)

	endpoint := iv.Endpoint

	if options.UsePremium {
		endpoint = iv.PremiumEndpoint
	}

	client, err := args.NewClient(endpoint, iv.Cred, options.ClientOptions)
	require.NoError(t, err)

	return client
}

func CreateExpiringQueue(t *testing.T, qd *atom.QueueDescription) (string, func()) {
	iv := GetIdentityVars(t)

	em, err := atom.NewEntityManager(iv.Endpoint, iv.Cred, "", nil)
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
