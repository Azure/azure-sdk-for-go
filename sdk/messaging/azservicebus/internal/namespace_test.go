// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/sbauth"
	"github.com/Azure/go-amqp"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/stretchr/testify/require"
)

func TestNewNamespaceWithAzureEnvironment(t *testing.T) {
	ns, err := NewNamespace(NamespaceWithAzureEnvironment("namespaceName", "AzureGermanCloud"))
	if err != nil {
		t.Fatalf("unexpected error creating namespace: %s", err)
	}
	if ns.Environment != azure.GermanCloud {
		t.Fatalf("expected namespace environment to be %q but was %q", azure.GermanCloud, ns.Environment)
	}
	if !strings.EqualFold(ns.Suffix, azure.GermanCloud.ServiceBusEndpointSuffix) {
		t.Fatalf("expected suffix to be %q but was %q", azure.GermanCloud.ServiceBusEndpointSuffix, ns.Suffix)
	}
	if ns.Name != "namespaceName" {
		t.Fatalf("expected namespace name to be %q but was %q", "namespaceName", ns.Name)
	}
}

// implements `Retrier` interface.
type fakeRetrier struct {
	tryCalled  int
	copyCalled int
}

func (r *fakeRetrier) Copy() Retrier {
	r.copyCalled++

	// NOTE: purposefully not making a copy so I can keep track of the
	// try/copy counts.
	return r
}

func (r *fakeRetrier) Exhausted() bool {
	return false
}

func (r *fakeRetrier) Try(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return false
	default:
	}

	r.tryCalled++
	return true
}

func (r *fakeRetrier) CurrentTry() int {
	return r.tryCalled
}

type fakeTokenCredential struct {
	azcore.TokenCredential
	expires time.Time
}

func (ftc *fakeTokenCredential) GetToken(ctx context.Context, options policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	return &azcore.AccessToken{
		ExpiresOn: ftc.expires,
	}, nil
}

func TestNamespaceNegotiateClaim(t *testing.T) {
	retrier := &fakeRetrier{}

	expires := time.Now().Add(24 * time.Hour)

	ns := &Namespace{
		baseRetrier:   retrier,
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{expires: expires}),
	}

	cbsNegotiateClaimCalled := 0

	cbsNegotiateClaim := func(ctx context.Context, audience string, conn *amqp.Client, provider auth.TokenProvider) error {
		cbsNegotiateClaimCalled++
		return nil
	}

	getAMQPClientCalled := 0

	getAMQPClient := func(ctx context.Context) (*amqp.Client, uint64, error) {
		getAMQPClientCalled++
		return &amqp.Client{}, 0, nil
	}

	// fire off a basic negotiate claim. The renewal duration is so long that it won't run - that's a separate test.
	cancel, err := ns.startNegotiateClaimRenewer(
		context.Background(),
		"my entity path",
		cbsNegotiateClaim,
		getAMQPClient,
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

	require.EqualValues(t, getAMQPClientCalled, 1)
	require.EqualValues(t, 1, cbsNegotiateClaimCalled)
	require.EqualValues(t, 1, retrier.copyCalled)
	require.EqualValues(t, 1, retrier.tryCalled)
}

func TestNamespaceNegotiateClaimRenewal(t *testing.T) {
	retrier := &fakeRetrier{}

	expires := time.Now().Add(24 * time.Hour)

	ns := &Namespace{
		baseRetrier:   retrier,
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{expires: expires}),
	}

	cbsNegotiateClaimCalled := 0

	cbsNegotiateClaim := func(ctx context.Context, audience string, conn *amqp.Client, provider auth.TokenProvider) error {
		cbsNegotiateClaimCalled++
		return nil
	}

	getAMQPClientCalled := 0

	notify := make(chan struct{})

	getAMQPClient := func(ctx context.Context) (*amqp.Client, uint64, error) {
		getAMQPClientCalled++

		if getAMQPClientCalled == 3 {
			close(notify)
			<-ctx.Done()
		}

		return &amqp.Client{}, 0, nil
	}

	var errorsLogged []error
	nextRefreshDurationChecks := 0

	// fire off a basic negotiate claim. The renewal duration is so long that it won't run - that's a separate test.
	cancel, err := ns.startNegotiateClaimRenewer(
		context.Background(),
		"my entity path",
		cbsNegotiateClaim, getAMQPClient, func(expirationTimeParam, currentTime time.Time) time.Duration {
			require.EqualValues(t, expires, expirationTimeParam)
			nextRefreshDurationChecks++
			return 0
		})
	defer cancel()

	require.NoError(t, err)
	time.Sleep(3 * time.Second) // make sure, even with variability, we get at least one renewal

	<-notify

	require.GreaterOrEqual(t, getAMQPClientCalled, 2+1) // that last +1 is when we blocked to prevent us renewing too much for our test!

	require.EqualValues(t, 3, retrier.copyCalled)
	require.EqualValues(t, 3, retrier.tryCalled)
	require.EqualValues(t, 2, nextRefreshDurationChecks)

	require.EqualValues(t, 2, cbsNegotiateClaimCalled)
	require.Empty(t, errorsLogged)

	cancel()
}

func TestNamespaceNegotiateClaimFailsToGetClient(t *testing.T) {
	ns := &Namespace{
		baseRetrier:   noRetryRetrier.Copy(),
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{expires: time.Now()}),
	}

	cancel, err := ns.startNegotiateClaimRenewer(
		context.Background(),
		"entity path",
		func(ctx context.Context, audience string, conn *amqp.Client, provider auth.TokenProvider) error {
			return errors.New("NegotiateClaim amqp.Client failed")
		}, func(ctx context.Context) (*amqp.Client, uint64, error) {
			return nil, 0, errors.New("Getting *amqp.Client failed")
		}, func(expirationTime, currentTime time.Time) time.Duration {
			// refresh immediately since we're in a unit test.
			return 0
		})

	require.EqualError(t, err, "Getting *amqp.Client failed")
	require.Nil(t, cancel)
}

func TestNamespaceNegotiateClaimFails(t *testing.T) {
	ns := &Namespace{
		baseRetrier:   noRetryRetrier.Copy(),
		TokenProvider: sbauth.NewTokenProvider(&fakeTokenCredential{expires: time.Now()}),
	}

	cancel, err := ns.startNegotiateClaimRenewer(
		context.Background(),
		"entity path",
		func(ctx context.Context, audience string, conn *amqp.Client, provider auth.TokenProvider) error {
			return errors.New("NegotiateClaim amqp.Client failed")
		}, func(ctx context.Context) (*amqp.Client, uint64, error) {
			return &amqp.Client{}, 0, nil
		}, func(expirationTime, currentTime time.Time) time.Duration {
			// not even used.
			return 0
		})

	require.EqualError(t, err, "NegotiateClaim amqp.Client failed")
	require.Nil(t, cancel)
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

var noRetryRetrier = NewBackoffRetrier(struct {
	MaxRetries int
	Factor     float64
	Jitter     bool
	Min        time.Duration
	Max        time.Duration
}{
	MaxRetries: 0,
})
