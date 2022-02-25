// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/sbauth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/auth"
	"github.com/Azure/go-amqp"
	"github.com/stretchr/testify/require"
)

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
	expires := time.Now().Add(24 * time.Hour)

	ns := &Namespace{
		retryOptions:  retryOptionsOnlyOnce,
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
}

func TestNamespaceNegotiateClaimRenewal(t *testing.T) {
	expires := time.Now().Add(24 * time.Hour)

	ns := &Namespace{
		retryOptions:  retryOptionsOnlyOnce,
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

	require.EqualValues(t, 2, nextRefreshDurationChecks)

	require.EqualValues(t, 2, cbsNegotiateClaimCalled)
	require.Empty(t, errorsLogged)

	cancel()
}

func TestNamespaceNegotiateClaimFailsToGetClient(t *testing.T) {
	ns := &Namespace{
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
