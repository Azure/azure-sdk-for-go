//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azeventgrid_test

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventgrid"
	"github.com/stretchr/testify/require"
)

func TestFailedAck(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce, err := messaging.NewCloudEvent("hello-source", "world", []byte("ack this one"), nil)
	require.NoError(t, err)

	_, err = c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []messaging.CloudEvent{ce}, nil)
	require.NoError(t, err)

	recvResp, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &azeventgrid.ReceiveCloudEventsOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10),
	})
	require.NoError(t, err)

	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
		LockTokens: []string{*recvResp.Value[0].BrokerProperties.LockToken},
	}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
	require.Equal(t, []string{*recvResp.Value[0].BrokerProperties.LockToken}, ackResp.SucceededLockTokens)

	// now let's try to do stuff with an "out of date" token
	t.Run("AcknowledgeCloudEvents", func(t *testing.T) {
		resp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
			LockTokens: []string{*recvResp.Value[0].BrokerProperties.LockToken},
		}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		// TODO: these two fields are not symmetrical - FailedLockTokens carries a reason.
		require.Equal(t, []azeventgrid.FailedLockToken{
			{
				LockToken:        recvResp.Value[0].BrokerProperties.LockToken,
				ErrorCode:        to.Ptr("TokenLost"),
				ErrorDescription: to.Ptr("Token has expired."),
			},
		}, resp.FailedLockTokens)
	})

	t.Run("RejectCloudEvents", func(t *testing.T) {
		resp, err := c.RejectCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.RejectOptions{
			LockTokens: []string{*recvResp.Value[0].BrokerProperties.LockToken},
		}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		// TODO: these two fields are not symmetrical - FailedLockTokens carries a reason.
		require.Equal(t, []azeventgrid.FailedLockToken{
			{
				LockToken:        recvResp.Value[0].BrokerProperties.LockToken,
				ErrorCode:        to.Ptr("TokenLost"),
				ErrorDescription: to.Ptr("Token has expired."),
			},
		}, resp.FailedLockTokens)
	})

	t.Run("ReleaseCloudEvents", func(t *testing.T) {
		t.Skipf("Skipping, server-bug preventing release from working properly. https://github.com/Azure/azure-sdk-for-go/issues/21530")

		resp, err := c.ReleaseCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.ReleaseOptions{
			LockTokens: []string{*recvResp.Value[0].BrokerProperties.LockToken},
		}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		// TODO: these two fields are not symmetrical - FailedLockTokens carries a reason.
		require.Equal(t, []azeventgrid.FailedLockToken{
			{
				LockToken:        recvResp.Value[0].BrokerProperties.LockToken,
				ErrorCode:        to.Ptr("TokenLost"),
				ErrorDescription: to.Ptr("Token has expired."),
			},
		}, resp.FailedLockTokens)
	})
}

func TestPartialAckFailure(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce, err := messaging.NewCloudEvent("hello-source", "world", []byte("event one"), nil)
	require.NoError(t, err)

	ce2, err := messaging.NewCloudEvent("hello-source", "world", []byte("event two"), nil)
	require.NoError(t, err)

	_, err = c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []messaging.CloudEvent{ce, ce2}, nil)
	require.NoError(t, err)

	events, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &azeventgrid.ReceiveCloudEventsOptions{
		MaxEvents: to.Ptr[int32](2),
	})
	require.NoError(t, err)

	// we'll ack one now so we can force a failure to happen.
	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
		LockTokens: []string{*events.Value[0].BrokerProperties.LockToken},
	}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)

	// this will result in a partial failure.
	ackResp, err = c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
		LockTokens: []string{
			*events.Value[0].BrokerProperties.LockToken,
			*events.Value[1].BrokerProperties.LockToken,
		},
	}, nil)
	require.NoError(t, err)
	require.Equal(t, []azeventgrid.FailedLockToken{
		{
			LockToken:        events.Value[0].BrokerProperties.LockToken,
			ErrorCode:        to.Ptr("TokenLost"),
			ErrorDescription: to.Ptr("Token has expired."),
		},
	}, ackResp.FailedLockTokens)
	require.Equal(t, []string{*events.Value[1].BrokerProperties.LockToken}, ackResp.SucceededLockTokens)
}

func TestReject(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce, err := messaging.NewCloudEvent("TestAbandon", "world", []byte("event one"), nil)
	require.NoError(t, err)

	_, err = c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []messaging.CloudEvent{ce}, nil)
	require.NoError(t, err)

	events, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		SpecVersion: "1.0",
		Data:        []byte("event one"),
		Source:      "TestAbandon",
		Type:        "world",
	}, events.Value[0].Event)

	require.Equal(t, int32(1), *events.Value[0].BrokerProperties.DeliveryCount, "DeliveryCount starts at 1")

	rejectResp, err := c.RejectCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.RejectOptions{
		LockTokens: []string{*events.Value[0].BrokerProperties.LockToken},
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
	t.Skipf("Skipping, server-bug preventing release from working properly. https://github.com/Azure/azure-sdk-for-go/issues/21530")

	c := newClientWrapper(t, nil)

	ce, err := messaging.NewCloudEvent("TestAbandon", "world", []byte("event one"), nil)
	require.NoError(t, err)

	_, err = c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []messaging.CloudEvent{ce}, nil)
	require.NoError(t, err)

	events, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, ce, events.Value[0].Event)

	require.Equal(t, int32(1), *events.Value[0].BrokerProperties.DeliveryCount, "DeliveryCount starts at 1")

	rejectResp, err := c.ReleaseCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.ReleaseOptions{
		LockTokens: []string{*events.Value[0].BrokerProperties.LockToken},
	}, nil)
	require.NoError(t, err)

	if len(rejectResp.FailedLockTokens) > 0 {
		for _, flt := range rejectResp.FailedLockTokens {
			t.Logf("FailedLockToken:\n  ec: %s\n  desc: %s\n  locktoken:%s", *flt.ErrorCode, *flt.ErrorDescription, *flt.LockToken)
		}
		require.Fail(t, "Failed to release events")
	}

	require.Empty(t, rejectResp.FailedLockTokens)

	events, err = c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	require.Equal(t, int32(2), *events.Value[0].BrokerProperties.DeliveryCount, "DeliveryCount is incremented")
	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, azeventgrid.AcknowledgeOptions{
		LockTokens: []string{*events.Value[0].BrokerProperties.LockToken},
	}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
}

func TestPublishingAndReceivingCloudEvents(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce1, err := messaging.NewCloudEvent("hello-source", "eventType", "hello world 1", nil)
	require.NoError(t, err)

	ce2, err := messaging.NewCloudEvent("hello-source", "eventType", "hello world 2", &messaging.CloudEventOptions{
		DataContentType: to.Ptr("data content type"),
		DataSchema:      to.Ptr("https://dataschema"),
		Extensions: map[string]any{
			"extension1": "extension1value",
		},
		Subject: to.Ptr("subject"),
	})
	require.NoError(t, err)

	type simpleType struct {
		Name string
	}

	ce3, err := messaging.NewCloudEvent("hello-source", "eventType", simpleType{Name: "simple type name"}, nil)
	require.NoError(t, err)

	_, err = c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []messaging.CloudEvent{ce1, ce2, ce3}, nil)
	require.NoError(t, err)

	resp, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &azeventgrid.ReceiveCloudEventsOptions{
		MaxEvents: to.Ptr[int32](3),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Value)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		SpecVersion: "1.0",
		Source:      "hello-source",
		Type:        "eventType",
		Data:        []byte("\"hello world 1\""),
	}, resp.Value[0].Event)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		SpecVersion:     "1.0",
		Source:          "hello-source",
		Type:            "eventType",
		DataSchema:      to.Ptr("https://dataschema"),
		Data:            []byte("\"hello world 2\""),
		DataContentType: to.Ptr("data content type"),
		Subject:         to.Ptr("subject"),
		Extensions: map[string]any{
			"extension1": "extension1value",
		},
	}, resp.Value[1].Event)

	bytes, err := json.Marshal(simpleType{Name: "simple type name"})
	require.NoError(t, err)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		SpecVersion: "1.0",
		Source:      "hello-source",
		Type:        "eventType",
		Data:        []byte(bytes),
	}, resp.Value[2].Event)

	ackArgs := azeventgrid.AcknowledgeOptions{}

	for _, e := range resp.Value {
		require.NotNil(t, e.BrokerProperties.LockToken)
		ackArgs.LockTokens = append(ackArgs.LockTokens, *e.BrokerProperties.LockToken)
	}

	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, ackArgs, nil)
	require.NoError(t, err)

	require.Empty(t, ackResp.FailedLockTokens)
	require.NotEmpty(t, ackResp.SucceededLockTokens)
}

func TestSimpleErrors(t *testing.T) {
	c := newClientWrapper(t, nil)

	_, err := c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []messaging.CloudEvent{
		{},
	}, nil)
	var respErr *azcore.ResponseError

	require.ErrorAs(t, err, &respErr)
	require.Equal(t, http.StatusBadRequest, respErr.StatusCode)
	require.Contains(t, respErr.Error(), "'data' attribute is required")
}
