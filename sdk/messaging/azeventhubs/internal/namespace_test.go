// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/telemetry"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/sbauth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/test"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

type fakeTokenCredential struct {
	azcore.TokenCredential
	expires time.Time
}

func (ftc *fakeTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{
		ExpiresOn: ftc.expires,
	}, nil
}

var retryOptionsOnlyOnce = exported.RetryOptions{
	MaxRetries: 0,
}

func TestNamespaceUserAgent(t *testing.T) {
	ns := &Namespace{}

	// Examples:
	// User agent, no application ID  : 'azsdk-go-azeventhubs/v0.3.1 (go1.19.3; linux)'
	// User agent, with application ID: 'userApplicationID azsdk-go-azeventhubs/v0.3.1 (go1.19.3; linux)'

	baseUserAgent := telemetry.Format("azeventhubs", Version)
	require.NotEmpty(t, baseUserAgent)

	t.Logf("User agent, no application ID  : '%s'", ns.getUserAgent())
	require.Equal(t, baseUserAgent, ns.getUserAgent())

	opt := NamespaceWithUserAgent("userApplicationID")
	require.NoError(t, opt(ns))

	t.Logf("User agent, with application ID: '%s'", ns.getUserAgent())
	require.Equal(t, fmt.Sprintf("userApplicationID %s", baseUserAgent), ns.getUserAgent())
}

func TestNamespaceNegotiateClaim(t *testing.T) {
	expires := time.Now().Add(24 * time.Hour)

	ns := &Namespace{
		RetryOptions:  retryOptionsOnlyOnce,
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{expires: expires}),
	}

	cbsNegotiateClaimCalled := 0

	cbsNegotiateClaim := func(ctx context.Context, audience string, conn amqpwrap.AMQPClient, provider auth.TokenProvider) error {
		cbsNegotiateClaimCalled++
		return nil
	}

	newAMQPClientCalled := 0

	ns.newClientFn = func(ctx context.Context, connID uint64) (amqpwrap.AMQPClient, error) {
		newAMQPClientCalled++
		return &amqpwrap.AMQPClientWrapper{}, nil
	}

	// fire off a basic negotiate claim. The renewal duration is so long that it won't run - that's a separate test.
	cancel, _, err := ns.startNegotiateClaimRenewer(
		context.Background(),
		"my entity path",
		cbsNegotiateClaim,
		func(expirationTimeParam, currentTime time.Time) time.Duration {
			require.EqualValues(t, expires, expirationTimeParam)
			// wiggle room, but just want to check that they're passing me the time.Now() value (silly)
			require.GreaterOrEqual(t, time.Minute, time.Since(currentTime))

			// we're going to cancel out pretty much immediately
			return 24 * time.Hour
		})
	defer cancel()

	require.NoError(t, err)
	cancel()

	require.EqualValues(t, newAMQPClientCalled, 1)
	require.EqualValues(t, 1, cbsNegotiateClaimCalled)
}

func TestNamespaceNegotiateClaimRenewal(t *testing.T) {
	expires := time.Now().Add(24 * time.Hour)

	ns := &Namespace{
		RetryOptions:  retryOptionsOnlyOnce,
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{expires: expires}),
	}

	cbsNegotiateClaimCalled := 0

	cbsNegotiateClaim := func(ctx context.Context, audience string, conn amqpwrap.AMQPClient, provider auth.TokenProvider) error {
		cbsNegotiateClaimCalled++
		return nil
	}

	var errorsLogged []error
	nextRefreshDurationChecks := 0

	ns.newClientFn = func(ctx context.Context, connID uint64) (amqpwrap.AMQPClient, error) {
		return &amqpwrap.AMQPClientWrapper{Inner: &amqp.Conn{}}, nil
	}

	cancel, _, err := ns.startNegotiateClaimRenewer(
		context.Background(),
		"my entity path",
		cbsNegotiateClaim, func(expirationTimeParam, currentTime time.Time) time.Duration {
			require.EqualValues(t, expires, expirationTimeParam)
			nextRefreshDurationChecks++

			if nextRefreshDurationChecks == 1 {
				return 0
			}

			return 24 * time.Hour // ie, we don't need to do it again.
		})
	defer cancel()

	require.NoError(t, err)
	time.Sleep(3 * time.Second) // make sure, even with variability, we get at least one renewal

	require.EqualValues(t, 2, nextRefreshDurationChecks)
	require.EqualValues(t, 2, cbsNegotiateClaimCalled)
	require.Empty(t, errorsLogged)

	cancel()
}

func TestNamespaceNegotiateClaimFailsToGetClient(t *testing.T) {
	ns := &Namespace{
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{expires: time.Now()}),
	}

	ns.newClientFn = func(ctx context.Context, connID uint64) (amqpwrap.AMQPClient, error) {
		return nil, errors.New("Getting *amqp.Client failed")
	}

	cancel, _, err := ns.startNegotiateClaimRenewer(
		context.Background(),
		"entity path",
		func(ctx context.Context, audience string, conn amqpwrap.AMQPClient, provider auth.TokenProvider) error {
			return errors.New("NegotiateClaim amqp.Client failed")
		}, func(expirationTime, currentTime time.Time) time.Duration {
			// refresh immediately since we're in a unit test.
			return 0
		})

	require.EqualError(t, err, "Getting *amqp.Client failed")
	require.Nil(t, cancel)
}

func TestNamespaceNegotiateClaimNonRenewableToken(t *testing.T) {
	ns := &Namespace{
		RetryOptions: retryOptionsOnlyOnce,
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{
			// credentials that don't renew return a zero-initialized time.
			expires: time.Time{},
		}),
	}

	cbsNegotiateClaimCalled := 0

	cbsNegotiateClaim := func(ctx context.Context, audience string, conn amqpwrap.AMQPClient, provider auth.TokenProvider) error {
		cbsNegotiateClaimCalled++
		return nil
	}

	ns.newClientFn = func(ctx context.Context, connID uint64) (amqpwrap.AMQPClient, error) {
		return &amqpwrap.AMQPClientWrapper{Inner: &amqp.Conn{}}, nil
	}

	// since the token is non-renewable we will just do the single cbsNegotiateClaim call and never renew.
	_, done, err := ns.startNegotiateClaimRenewer(
		context.Background(),
		"my entity path",
		cbsNegotiateClaim,
		func(expirationTimeParam, currentTime time.Time) time.Duration {
			panic("Won't be called, no refreshing of claims will be done")
		})

	require.NoError(t, err)
	require.Equal(t, 1, cbsNegotiateClaimCalled)

	select {
	case <-done:
	default:
		require.Fail(t, "cancel() returns a channel that is already Done()")
	}
}

func TestNamespaceNegotiateClaimFails(t *testing.T) {
	ns := &Namespace{
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{expires: time.Now()}),
	}

	ns.newClientFn = func(ctx context.Context, connID uint64) (amqpwrap.AMQPClient, error) {
		return &fakeAMQPClient{}, nil
	}

	cancel, _, err := ns.startNegotiateClaimRenewer(
		context.Background(),
		"entity path",
		func(ctx context.Context, audience string, conn amqpwrap.AMQPClient, provider auth.TokenProvider) error {
			return errors.New("NegotiateClaim amqp.Client failed")
		}, func(expirationTime, currentTime time.Time) time.Duration {
			// not even used.
			return 0
		})

	require.EqualError(t, err, "NegotiateClaim amqp.Client failed")
	require.Nil(t, cancel)
}

func TestNamespaceNegotiateClaimFatalErrors(t *testing.T) {
	ns := &Namespace{
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{expires: time.Now()}),
	}

	cbsNegotiateClaimCalled := 0

	cbsNegotiateClaim := func(ctx context.Context, audience string, conn amqpwrap.AMQPClient, provider auth.TokenProvider) error {
		cbsNegotiateClaimCalled++

		// work the first time, fail on renewals.
		if cbsNegotiateClaimCalled > 1 {
			return errNonRetriable{Message: "non retriable error message"}
		}

		return nil
	}

	endCapture := test.CaptureLogsForTest()
	defer endCapture()

	ns.newClientFn = func(ctx context.Context, connID uint64) (amqpwrap.AMQPClient, error) {
		return &amqpwrap.AMQPClientWrapper{Inner: &amqp.Conn{}}, nil
	}

	_, done, err := ns.startNegotiateClaimRenewer(
		context.Background(),
		"entity path",
		cbsNegotiateClaim, func(expirationTime, currentTime time.Time) time.Duration {
			// instant renewals.
			return 0
		})

	require.NoError(t, err)

	select {
	case <-done:
		logs := endCapture()
		// check the log messages - we should have one telling us why we stopped the claims loop
		require.Contains(t, logs, "[azeh.Auth] [entity path] fatal error, stopping token refresh loop: non retriable error message")
	case <-time.After(3 * time.Second):
		// was locked! Should have been closed.
		require.Fail(t, "claim renewal was automatically cancelled because of a non-retriable error")
	}
}

func TestNamespaceNextClaimRefreshDuration(t *testing.T) {
	now := time.Now()

	clockDrift := 10 * time.Minute
	lessThanMin := now.Add(119 * time.Second).Add(clockDrift)
	greaterThanMax := now.Add(49*24*time.Hour + time.Second).Add(clockDrift)

	require.EqualValues(t, 2*time.Minute, nextClaimRefreshDuration(lessThanMin, now),
		"Just under the min refresh time, so we get the min instead")

	require.EqualValues(t, 49*24*time.Hour, nextClaimRefreshDuration(greaterThanMax, now),
		"Just over the max refresh time, so we just get the max instead")

	require.EqualValues(t, 3*time.Minute, nextClaimRefreshDuration(now.Add(3*time.Minute+clockDrift), now))
}

func TestNamespaceStaleConnection(t *testing.T) {
	ns := &Namespace{
		RetryOptions: retryOptionsOnlyOnce,
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{
			// credentials that don't renew return a zero-initialized time.
			expires: time.Time{},
		}),
	}

	fakeClient := &fakeAMQPClient{}

	ns.client = fakeClient
	ns.connID = 101

	require.NoError(t, ns.Close(context.Background(), false))
	require.Equal(t, 1, fakeClient.closeCalled)
	require.Nil(t, ns.client)

	ns.newClientFn = func(ctx context.Context, connID uint64) (amqpwrap.AMQPClient, error) {
		return &fakeAMQPClient{}, nil
	}

	client, clientID, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.NotSame(t, fakeClient, client, "A new client should be created")
	require.Equal(t, uint64(101+1), clientID, "Client ID is incremented since we had to recreate it")
	require.NotNil(t, client)
}

func TestNamespaceUpdateClientWithoutLock(t *testing.T) {
	newClient := 0
	var clientToReturn amqpwrap.AMQPClient
	var err error

	ns := &Namespace{
		newClientFn: func(ctx context.Context, connID uint64) (amqpwrap.AMQPClient, error) {
			newClient++
			return clientToReturn, err
		},
		connID: 101,
	}

	err = errors.New("client error")

	client, clientID, err := ns.updateClientWithoutLock(context.Background())
	require.Error(t, err, "client error")
	require.Equal(t, uint64(0), clientID)
	require.Nil(t, client)

	// when they create a new client they'll get this one.
	clientToReturn = &fakeAMQPClient{}
	err = nil

	client, clientID, err = ns.updateClientWithoutLock(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(101+1), clientID)
	require.Same(t, clientToReturn, client)

	// change out the returned client (it won't get used because we return the cached one in ns.client)
	origClient := client
	clientToReturn = &fakeAMQPClient{}

	client, clientID, err = ns.updateClientWithoutLock(context.Background())
	require.NoError(t, err)
	require.Equal(t, uint64(101+1), clientID)
	require.Same(t, origClient, client)
}

func TestNamespaceConnectionRecovery(t *testing.T) {
	type testData struct {
		NS              *Namespace
		NewClientCount  int
		FakeClientError error
		FakeClient      *fakeAMQPClient
	}

	init := func() *testData {
		td := &testData{}
		td.NS = &Namespace{
			connID: 2,
			newClientFn: func(ctx context.Context, connID uint64) (amqpwrap.AMQPClient, error) {
				td.NewClientCount++
				return td.FakeClient, td.FakeClientError
			},
		}
		return td
	}

	t.Run("stale connection ID", func(t *testing.T) {
		testData := init()

		// ie, my connection is stale (it doesn't actually matter if the connID is >, although that's impossible
		// since it means their connection came from the future)
		origConnID := testData.NS.connID

		err := testData.NS.Recover(context.Background(), testData.NS.connID-1)
		require.Zero(t, testData.NewClientCount, "existing client is re-used")
		require.Equal(t, origConnID, testData.NS.connID, "no new client created, connID is unchanged")
		require.NoError(t, err)
	})

	t.Run("connection matches", func(t *testing.T) {
		testData := init()
		testData.FakeClient = &fakeAMQPClient{} // new client that was "created" for our recovery

		// this time the connection must be having errors AND it matches our current ID
		origConnID := testData.NS.connID
		origClient := &fakeAMQPClient{}
		testData.NS.client = origClient

		err := testData.NS.Recover(context.Background(), testData.NS.connID)
		require.Equal(t, 1, testData.NewClientCount, "new client is created (assumption is if it matches then our current connection is returning errors)")
		require.Equal(t, origConnID+1, testData.NS.connID, "new client created, connID increments")
		require.NoError(t, err)
		require.Equal(t, 1, origClient.closeCalled, "old client is closed")
		require.NotSame(t, origClient, testData.NS.client)
	})

	t.Run("recover but failed", func(t *testing.T) {
		testData := init()

		// and the last outcome - we did try to recover, but failed. We will end up in a state
		// where the client will be nil, so the next attempt to get the client will create
		// a new one.
		testData.FakeClientError = errors.New("we failed to create the connection!")
		origConnID := testData.NS.connID

		err := testData.NS.Recover(context.Background(), origConnID)
		require.Equal(t, testData.FakeClientError, err)
		require.Equal(t, 1, testData.NewClientCount, "we did attempt to create a new client, it just failed.")
		require.Equal(t, origConnID, testData.NS.connID, "new client failed to be created so the conn ID is unchanged")

		// if the namespace is closed then this function fails.
		_ = testData.NS.Close(context.Background(), true)
		err = testData.NS.Recover(context.Background(), origConnID)
		require.ErrorIs(t, err, ErrClientClosed)
	})
}

func TestNamespaceCantStopRecoverFromClosingConn(t *testing.T) {
	numCancels := 0
	numClients := 0

	ns := &Namespace{
		newClientFn: func(ctx context.Context, connID uint64) (amqpwrap.AMQPClient, error) {
			select {
			case <-ctx.Done():
				numCancels++
				return nil, ctx.Err()
			default:
				numClients++
				client := &fakeAMQPClient{}
				return client, nil
			}
		},
	}

	conn, id, err := ns.GetAMQPClientImpl(context.Background())
	require.NoError(t, err)
	require.NotNil(t, conn)
	require.Equal(t, uint64(1), id)

	require.Equal(t, 1, numClients)
	require.Equal(t, 0, numCancels)

	canceledCtx, cancel := context.WithCancel(context.Background())
	cancel()

	err = ns.Recover(canceledCtx, id)

	// two key things:
	// 1. the old client gets closed, even when the 'ctx' is cancelled.
	// 2. since the context is cancelled we don't create a new one.
	require.ErrorIs(t, err, context.Canceled)
	require.Equal(t, 1, numClients, "we did NOT create a new client")
	require.Equal(t, 1, numCancels, "we cancelled a client creation")
	require.False(t, ns.closedPermanently)
}

func TestNamespaceDisablingAMQPS(t *testing.T) {
	t.Run("UseDevelopmentEmulator", func(t *testing.T) {
		cs := "Endpoint=sb://localhost:6765;SharedAccessKeyName=" + "MyKey" + ";SharedAccessKey=" + "MySecret" + ";UseDevelopmentEmulator=true"
		ns, err := NewNamespace(NamespaceWithConnectionString(cs))
		require.NoError(t, err)

		audience := ns.GetEntityAudience("hub1")
		require.Equal(t, "amqp://localhost:6765/hub1", audience)
	})

	t.Run("Normal", func(t *testing.T) {
		cs := "Endpoint=sb://localhost:6765;SharedAccessKeyName=" + "MyKey" + ";SharedAccessKey=" + "MySecret"
		ns, err := NewNamespace(NamespaceWithConnectionString(cs))
		require.NoError(t, err)

		audience := ns.GetEntityAudience("hub1")
		require.Equal(t, "amqps://localhost:6765/hub1", audience)
	})

	t.Run("TokenCredential", func(t *testing.T) {
		ns, err := NewNamespace(NamespaceWithTokenCredential("localhost:6765", &fakeTokenCredential{}))
		require.NoError(t, err)

		audience := ns.GetEntityAudience("hub1")
		require.Equal(t, "amqps://localhost:6765/hub1", audience)
	})
}
