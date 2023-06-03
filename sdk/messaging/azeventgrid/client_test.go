//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid"
	"github.com/stretchr/testify/require"
)

func TestFailedAck(t *testing.T) {
	c := newClientForTest()
	defer c.cleanup()

	pubResp, err := c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []*azeventgrid.CloudEvent{
		{
			Data:   []byte("ack this one"),
			Source: to.Ptr("hello-source"),
			Type:   to.Ptr("world"),
		},
	}, nil)
	require.NoError(t, err)

	// just documenting this, I don't think the return value is useful.
	require.Equal(t, map[string]interface{}{}, pubResp.Interface)

	recvResp, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &azeventgrid.ClientReceiveCloudEventsOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10),
	})
	require.NoError(t, err)

	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
		LockTokens: []*string{recvResp.Value[0].BrokerProperties.LockToken},
	}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
	require.Equal(t, []*string{recvResp.Value[0].BrokerProperties.LockToken}, ackResp.SucceededLockTokens)

	// now let's try to do stuff with an "out of date" token
	t.Run("AcknowledgeCloudEvents", func(t *testing.T) {
		resp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
			LockTokens: []*string{recvResp.Value[0].BrokerProperties.LockToken},
		}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		// TODO: these two fields are not symmetrical - FailedLockTokens carries a reason.
		require.Equal(t, []*azeventgrid.FailedLockToken{
			{
				LockToken:        recvResp.Value[0].BrokerProperties.LockToken,
				ErrorCode:        to.Ptr("TokenLost"),
				ErrorDescription: to.Ptr("Token has expired."),
			},
		}, resp.FailedLockTokens)
	})

	t.Run("RejectCloudEvents", func(t *testing.T) {
		resp, err := c.RejectCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.RejectOptions{
			LockTokens: []*string{recvResp.Value[0].BrokerProperties.LockToken},
		}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		// TODO: these two fields are not symmetrical - FailedLockTokens carries a reason.
		require.Equal(t, []*azeventgrid.FailedLockToken{
			{
				LockToken:        recvResp.Value[0].BrokerProperties.LockToken,
				ErrorCode:        to.Ptr("TokenLost"),
				ErrorDescription: to.Ptr("Token has expired."),
			},
		}, resp.FailedLockTokens)
	})

	t.Run("AcknowledgeCloudEvents", func(t *testing.T) {
		resp, err := c.ReleaseCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.ReleaseOptions{
			LockTokens: []*string{recvResp.Value[0].BrokerProperties.LockToken},
		}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		// TODO: these two fields are not symmetrical - FailedLockTokens carries a reason.
		require.Equal(t, []*azeventgrid.FailedLockToken{
			{
				LockToken:        recvResp.Value[0].BrokerProperties.LockToken,
				ErrorCode:        to.Ptr("TokenLost"),
				ErrorDescription: to.Ptr("Token has expired."),
			},
		}, resp.FailedLockTokens)
	})
}

func TestPartialAckFailure(t *testing.T) {
	// this API seems to have some sharp edges:
	//
	// 1. The return result lists the lock tokens that failed and succeeded in separate lists which means
	//    (AFAICT) that I'd need to iterate through both lists and the original list of lock tokens multiple
	//    times in order to whittle them down.
	//
	// 2. The doc comment makes it sound like we might end up returning both a response and an error, but we
	//    don't want it to be that weird to use. Either we return a failure or we return a response but we
	//    don't want to return both because nobody will check both values.
	// ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{}, nil)
	// require.NoError(t, err)

	c := newClientForTest()
	defer c.cleanup()

	_, err := c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []*azeventgrid.CloudEvent{
		{
			Data:   []byte("event one"),
			Source: to.Ptr("hello-source"),
			Type:   to.Ptr("world"),
		},
		{
			Data:   []byte("event two"),
			Source: to.Ptr("hello-source"),
			Type:   to.Ptr("world"),
		},
	}, nil)
	require.NoError(t, err)
}

func TestPublishingAndReceivingCloudEvents(t *testing.T) {
	c := newClientForTest()
	defer c.cleanup()

	_, err := c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []*azeventgrid.CloudEvent{
		{
			Data:   "hello world",
			Source: to.Ptr("hello-source"),
			Type:   to.Ptr("world"),
		},
	}, nil)
	require.NoError(t, err)

	resp, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)
	require.NotEmpty(t, resp.Value)

	// this doesn't work - it comes back as a base64 encoded string.
	require.Equal(t, "hello world", resp.Value[0].Event.Data)
	require.Equal(t, "hello-source", *resp.Value[0].Event.Source)
	require.Equal(t, "world", *resp.Value[0].Event.Type)

	ackArgs := azeventgrid.AcknowledgeOptions{}

	for _, e := range resp.Value {
		require.NotNil(t, e.BrokerProperties.LockToken)
		ackArgs.LockTokens = append(ackArgs.LockTokens, e.BrokerProperties.LockToken)
	}

	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, ackArgs, nil)
	require.NoError(t, err)

	require.Empty(t, ackResp.FailedLockTokens)
	require.NotEmpty(t, ackResp.SucceededLockTokens)
}
