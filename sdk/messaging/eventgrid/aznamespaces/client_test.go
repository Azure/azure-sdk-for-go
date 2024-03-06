//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces"
	"github.com/stretchr/testify/require"
)

func TestFailedAck(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce, err := messaging.NewCloudEvent("hello-source", "world", []byte("ack this one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	_, err = c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []messaging.CloudEvent{ce}, nil)
	require.NoError(t, err)

	recvResp, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &aznamespaces.ReceiveCloudEventsOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10),
	})
	require.NoError(t, err)

	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*recvResp.Value[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
	require.Equal(t, []string{*recvResp.Value[0].BrokerProperties.LockToken}, ackResp.SucceededLockTokens)

	// now let's try to do stuff with an "out of date" token
	t.Run("AcknowledgeCloudEvents", func(t *testing.T) {
		resp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*recvResp.Value[0].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		requireFailedLockTokens(t, []string{*recvResp.Value[0].BrokerProperties.LockToken}, resp.FailedLockTokens)
	})

	t.Run("RejectCloudEvents", func(t *testing.T) {
		resp, err := c.RejectCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*recvResp.Value[0].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		requireFailedLockTokens(t, []string{*recvResp.Value[0].BrokerProperties.LockToken}, resp.FailedLockTokens)
	})

	t.Run("ReleaseCloudEvents", func(t *testing.T) {
		resp, err := c.ReleaseCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*recvResp.Value[0].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		requireFailedLockTokens(t, []string{*recvResp.Value[0].BrokerProperties.LockToken}, resp.FailedLockTokens)
	})
}

func TestPartialAckFailure(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce, err := messaging.NewCloudEvent("hello-source", "world", []byte("event one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	ce2, err := messaging.NewCloudEvent("hello-source", "world", []byte("event two"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	_, err = c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []messaging.CloudEvent{ce, ce2}, nil)
	require.NoError(t, err)

	const numExpectedEvents = 2

	var events []aznamespaces.ReceiveDetails

	receiveCtx, cancelReceive := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelReceive()

	for len(events) < numExpectedEvents {
		eventsResp, err := c.ReceiveCloudEvents(receiveCtx, c.TestVars.Topic, c.TestVars.Subscription, &aznamespaces.ReceiveCloudEventsOptions{
			MaxEvents: to.Ptr(int32(numExpectedEvents - len(events))),
		})
		require.NoError(t, err)
		events = append(events, eventsResp.Value...)
	}

	// we'll ack one now so we can force a failure to happen.
	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*events[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)

	// this will result in a partial failure.
	ackResp, err = c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{
		*events[0].BrokerProperties.LockToken,
		*events[1].BrokerProperties.LockToken,
	}, nil)
	require.NoError(t, err)

	requireFailedLockTokens(t, []string{*events[0].BrokerProperties.LockToken}, ackResp.FailedLockTokens)
	require.Equal(t, []string{*events[1].BrokerProperties.LockToken}, ackResp.SucceededLockTokens)
}

func TestReject(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce, err := messaging.NewCloudEvent("TestAbandon", "world", []byte("event one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	t.Logf("Publishing cloud events")
	_, err = c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []messaging.CloudEvent{ce}, nil)
	require.NoError(t, err)

	t.Logf("Receiving cloud events")
	events, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		SpecVersion:     "1.0",
		DataContentType: to.Ptr("application/octet-stream"),
		Data:            []byte("event one"),
		Source:          "TestAbandon",
		Type:            "world",
	}, events.Value[0].Event)

	require.Equal(t, int32(1), *events.Value[0].BrokerProperties.DeliveryCount, "DeliveryCount starts at 1")

	t.Logf("Rejecting cloud events")
	rejectResp, err := c.RejectCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*events.Value[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, rejectResp.FailedLockTokens)
	t.Logf("Done rejecting cloud events")

	events, err = c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &aznamespaces.ReceiveCloudEventsOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10),
	})
	require.NoError(t, err)
	require.Empty(t, events.Value)
}

func TestRelease(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce, err := messaging.NewCloudEvent("TestAbandon", "world", []byte("event one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	_, err = c.PublishCloudEvents(context.Background(), c.TestVars.Topic, []messaging.CloudEvent{ce}, nil)
	require.NoError(t, err)

	events, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, ce, events.Value[0].Event)

	require.Equal(t, int32(1), *events.Value[0].BrokerProperties.DeliveryCount, "DeliveryCount starts at 1")

	rejectResp, err := c.ReleaseCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*events.Value[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)

	if len(rejectResp.FailedLockTokens) > 0 {
		for _, flt := range rejectResp.FailedLockTokens {
			t.Logf("FailedLockToken:\n  ec: %s\n  desc: %s\n  locktoken:%s", *flt.Error.Code, flt.Error.Error(), *flt.LockToken)
		}
		require.Fail(t, "Failed to release events")
	}

	require.Empty(t, rejectResp.FailedLockTokens)

	events, err = c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	require.Equal(t, int32(2), *events.Value[0].BrokerProperties.DeliveryCount, "DeliveryCount is incremented")
	ackResp, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*events.Value[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
}

func TestPublishBytes(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce, err := messaging.NewCloudEvent("hello-source", "eventType", []byte("TestPublishBytes"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	_, err = c.PublishCloudEvent(context.Background(), c.TestVars.Topic, ce, nil)
	require.NoError(t, err)

	recvResp, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		Source:          "hello-source",
		SpecVersion:     "1.0",
		Type:            "eventType",
		Data:            []byte("TestPublishBytes"),
		DataContentType: to.Ptr("application/octet-stream"),
	}, recvResp.Value[0].Event)
}

func TestPublishString(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce, err := messaging.NewCloudEvent("hello-source", "eventType", "TestPublishString", &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/json"),
	})
	require.NoError(t, err)

	_, err = c.PublishCloudEvent(context.Background(), c.TestVars.Topic, ce, nil)
	require.NoError(t, err)

	recvResp, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		Source:          "hello-source",
		SpecVersion:     "1.0",
		Type:            "eventType",
		Data:            []byte("\"TestPublishString\""), // non []byte returns as the JSON bytes
		DataContentType: to.Ptr("application/json"),
	}, recvResp.Value[0].Event)
}

func TestPublishingAndReceivingMultipleCloudEvents(t *testing.T) {
	c := newClientWrapper(t, nil)

	testData := []struct {
		Send     messaging.CloudEvent
		Expected messaging.CloudEvent
	}{
		{
			Send: mustCreateEvent(t, "hello-source", "eventType", []byte("TestPublishingAndReceivingMultipleCloudEvents 1"), &messaging.CloudEventOptions{
				DataContentType: to.Ptr("application/octet-stream"),
			}),
			Expected: messaging.CloudEvent{
				SpecVersion:     "1.0",
				Source:          "hello-source",
				Type:            "eventType",
				DataContentType: to.Ptr("application/octet-stream"),
				Data:            []byte("TestPublishingAndReceivingMultipleCloudEvents 1"),
			},
		},
		{
			Send: mustCreateEvent(t, "hello-source", "eventType", "TestPublishingAndReceivingMultipleCloudEvents 2", &messaging.CloudEventOptions{
				DataContentType: to.Ptr("application/json"),
				DataSchema:      to.Ptr("https://dataschema"),
				Extensions: map[string]any{
					"extension1": "extension1value",
				},
				Subject: to.Ptr("subject"),
			}),
			Expected: messaging.CloudEvent{
				SpecVersion:     "1.0",
				Source:          "hello-source",
				Type:            "eventType",
				DataSchema:      to.Ptr("https://dataschema"),
				Data:            []byte("\"TestPublishingAndReceivingMultipleCloudEvents 2\""),
				DataContentType: to.Ptr("application/json"),
				Subject:         to.Ptr("subject"),
				Extensions: map[string]any{
					"extension1": "extension1value",
				},
			},
		},
	}

	var batch []messaging.CloudEvent

	for _, td := range testData {
		batch = append(batch, td.Send)
	}

	// type simpleType struct {
	// 	Name string
	// }

	// ce3, err := messaging.NewCloudEvent("hello-source", "eventType", simpleType{Name: "simple type name"}, &messaging.CloudEventOptions{
	// 	DataContentType: to.Ptr("application/octet-stream"),
	// })
	// require.NoError(t, err)
	// toSend = append(toSend, ce3)

	// _, err := c.PublishCloudEvents(context.Background(), c.TestVars.Topic, batch, nil)
	// require.NoError(t, err)

	t.Logf("\n\n\n=====> starting our test, publishing\n\n\n")

	// _, err := c.PublishCloudEvent(context.Background(), c.TestVars.Topic, batch[0], nil)
	// require.NoError(t, err)

	_, err := c.PublishCloudEvents(context.Background(), c.TestVars.Topic, batch, nil)
	require.NoError(t, err)

	t.Logf("\n\n\n=====> starting our test, receiving\n\n\n")

	resp, err := c.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, &aznamespaces.ReceiveCloudEventsOptions{
		MaxEvents:   to.Ptr(int32(len(batch))),
		MaxWaitTime: to.Ptr(int32(60)),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Value)

	for i := 0; i < len(batch); i++ {
		_, err := c.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*resp.Value[i].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)

		requireEqualCloudEvent(t, testData[i].Expected, resp.Value[i].Event)
	}

	// bytes, err := json.Marshal(simpleType{Name: "simple type name"})
	// require.NoError(t, err)

	// requireEqualCloudEvent(t, messaging.CloudEvent{
	// 	SpecVersion: "1.0",
	// 	Source:      "hello-source",
	// 	Type:        "eventType",
	// 	Data:        []byte(bytes),
	// }, resp.Value[2].Event)
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

func TestRenewCloudEventLocks(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce := mustCreateEvent(t, "source", "eventType", "hello world", nil)
	_, err := c.Client.PublishCloudEvent(context.Background(), c.TestVars.Topic, ce, nil)
	require.NoError(t, err)

	recvResp, err := c.Client.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	_, err = c.Client.RenewCloudEventLocks(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*recvResp.Value[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)

	ackResp, err := c.Client.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*recvResp.Value[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
}

func TestReleaseWithDelay(t *testing.T) {
	c := newClientWrapper(t, nil)

	ce := mustCreateEvent(t, "source", "eventType", "hello world", nil)
	_, err := c.Client.PublishCloudEvent(context.Background(), c.TestVars.Topic, ce, nil)
	require.NoError(t, err)

	recvResp, err := c.Client.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)

	releaseResp, err := c.Client.ReleaseCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*recvResp.Value[0].BrokerProperties.LockToken}, &aznamespaces.ReleaseCloudEventsOptions{
		ReleaseDelayInSeconds: to.Ptr(aznamespaces.ReleaseDelayBy10Seconds),
	})
	require.NoError(t, err)
	require.Empty(t, releaseResp.FailedLockTokens)

	now := time.Now()

	// message will be available, but not immediately.
	recvResp, err = c.Client.ReceiveCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, nil)
	require.NoError(t, err)
	require.NotEmpty(t, recvResp.Value)
	require.Equal(t, int32(2), *recvResp.Value[0].BrokerProperties.DeliveryCount)

	if recording.GetRecordMode() == recording.LiveMode {
		// doesn't work when recording but it's somewhat unimportant there.
		elapsed := time.Since(now)
		require.GreaterOrEqual(t, int(elapsed/time.Second), 8) // give a little wiggle room for potential delays between requests, etc...
	}

	ackResp, err := c.Client.AcknowledgeCloudEvents(context.Background(), c.TestVars.Topic, c.TestVars.Subscription, []string{*recvResp.Value[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
}

func mustCreateEvent(t *testing.T, source string, eventType string, data any, options *messaging.CloudEventOptions) messaging.CloudEvent {
	event, err := messaging.NewCloudEvent(source, eventType, data, options)
	require.NoError(t, err)

	return event
}

func requireFailedLockTokens(t *testing.T, lockTokens []string, flts []aznamespaces.FailedLockToken) {
	for i, flt := range flts {
		// make sure the lock tokens line up
		require.Equal(t, lockTokens[i], *flt.LockToken)
		require.Equal(t, flt.Error.Code, to.Ptr("TokenLost"))
		require.EqualError(t, flt.Error, "Token has expired.")
	}
}
