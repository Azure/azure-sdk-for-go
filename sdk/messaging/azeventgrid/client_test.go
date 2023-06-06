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

	recvResp, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &azeventgrid.ReceiveCloudEventsOptions{
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

	events, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &azeventgrid.ReceiveCloudEventsOptions{
		MaxEvents: to.Ptr[int32](2),
	})
	require.NoError(t, err)

	// we'll ack one now so we can force a failure to happen.
	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
		LockTokens: []*string{events.Value[0].BrokerProperties.LockToken},
	}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)

	// this will result in a partial failure.
	ackResp, err = c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
		LockTokens: []*string{
			events.Value[0].BrokerProperties.LockToken,
			events.Value[1].BrokerProperties.LockToken,
		},
	}, nil)
	require.NoError(t, err)
	require.Equal(t, []*azeventgrid.FailedLockToken{
		{
			LockToken:        events.Value[0].BrokerProperties.LockToken,
			ErrorCode:        to.Ptr("TokenLost"),
			ErrorDescription: to.Ptr("Token has expired."),
		},
	}, ackResp.FailedLockTokens)
	require.Equal(t, []*string{events.Value[1].BrokerProperties.LockToken}, ackResp.SucceededLockTokens)
}

// func TestPartialAbandon(t *testing.T) {
// 	c := newClientForTest()
// 	defer c.cleanup()

// 	_, err := c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []*azeventgrid.CloudEvent{
// 		{
// 			Data:   []byte("event one"),
// 			Source: to.Ptr("hello-source"),
// 			Type:   to.Ptr("world"),
// 		},
// 		{
// 			Data:   []byte("abandon"),
// 			Source: to.Ptr("hello-source"),
// 			Type:   to.Ptr("world"),
// 		},
// 		{
// 			Data:   []byte("release"),
// 			Source: to.Ptr("hello-source"),
// 			Type:   to.Ptr("world"),
// 		},
// 	}, nil)
// 	require.NoError(t, err)

// 	events, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &azeventgrid.ReceiveCloudEventsOptions{
// 		MaxEvents: to.Ptr[int32](2),
// 	})
// 	require.NoError(t, err)
// }

func TestReject(t *testing.T) {
	c := newClientForTest()
	defer c.cleanup()

	_, err := c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []*azeventgrid.CloudEvent{
		{
			Data:   "event one",
			Source: to.Ptr("TestAbandon"),
			Type:   to.Ptr("world"),
		},
	}, nil)
	require.NoError(t, err)

	events, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, &azeventgrid.CloudEvent{
		Data:   "event one",
		Source: to.Ptr("TestAbandon"),
		Type:   to.Ptr("world"),
	}, events.Value[0].Event)

	require.Equal(t, int32(1), *events.Value[0].BrokerProperties.DeliveryCount, "DeliveryCount starts at 1")

	rejectResp, err := c.RejectCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.RejectOptions{
		LockTokens: []*string{events.Value[0].BrokerProperties.LockToken},
	}, nil)
	require.NoError(t, err)
	require.Empty(t, rejectResp.FailedLockTokens)

	events, err = c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &azeventgrid.ReceiveCloudEventsOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10),
	})
	require.NoError(t, err)
	require.Empty(t, events.Value)
}

func TestRelease(t *testing.T) {
	c := newClientForTest()
	defer c.cleanup()

	_, err := c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []*azeventgrid.CloudEvent{
		{
			Data:   "event one",
			Source: to.Ptr("TestAbandon"),
			Type:   to.Ptr("world"),
		},
	}, nil)
	require.NoError(t, err)

	events, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, &azeventgrid.CloudEvent{
		Data:   "event one",
		Source: to.Ptr("TestAbandon"),
		Type:   to.Ptr("world"),
	}, events.Value[0].Event)

	require.Equal(t, int32(1), *events.Value[0].BrokerProperties.DeliveryCount, "DeliveryCount starts at 1")

	rejectResp, err := c.ReleaseCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.ReleaseOptions{
		LockTokens: []*string{events.Value[0].BrokerProperties.LockToken},
	}, nil)
	require.NoError(t, err)
	require.Empty(t, rejectResp.FailedLockTokens)

	events, err = c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	require.Equal(t, int32(2), *events.Value[0].BrokerProperties.DeliveryCount, "DeliveryCount is incremented")
	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
		LockTokens: []*string{events.Value[0].BrokerProperties.LockToken},
	}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
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

// https://github.com/cloudevents/spec/blob/v1.0/json-format.md#31-handling-of-data
