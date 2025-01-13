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

func TestClients_UsingSASKey(t *testing.T) {
	sender, receiver := newClients(t, true)

	ce, err := messaging.NewCloudEvent("source", "eventType", "hello world", nil)
	require.NoError(t, err)

	sendResp, err := sender.SendEvent(context.Background(), &ce, nil)
	require.NoError(t, err)
	require.Empty(t, sendResp)

	recvResp, err := receiver.ReceiveEvents(context.Background(), nil)
	require.NoError(t, err)

	require.Equal(t, 1, len(recvResp.Details))
	lockToken := recvResp.Details[0].BrokerProperties.LockToken
	require.NotEmpty(t, *lockToken)

	// strings are json serialized (if you want to send the string as bytes you can just []byte("your string")
	// when creating the CloudEvent)
	require.Equal(t, "\"hello world\"", string(recvResp.Details[0].Event.Data.([]byte)))

	ackResp, err := receiver.AcknowledgeEvents(context.Background(), []string{*lockToken}, nil)
	require.NoError(t, err)
	require.Equal(t, []string{*lockToken}, ackResp.SucceededLockTokens)
}

func TestFailedAck(t *testing.T) {
	ce, err := messaging.NewCloudEvent("TestFailedAck", "world", []byte("ack this one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	sender, receiver := newClients(t, false)

	sendResp, err := sender.SendEvents(context.Background(), []*messaging.CloudEvent{&ce}, nil)
	require.Empty(t, sendResp)
	require.NoError(t, err)

	recvResp, err := receiver.ReceiveEvents(context.Background(), &aznamespaces.ReceiveEventsOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10),
	})
	require.NoError(t, err)

	ackResp, err := receiver.AcknowledgeEvents(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
	require.Equal(t, []string{*recvResp.Details[0].BrokerProperties.LockToken}, ackResp.SucceededLockTokens)

	// now let's try to do stuff with an "out of date" token
	t.Run("Acknowledge", func(t *testing.T) {
		resp, err := receiver.AcknowledgeEvents(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		requireFailedLockTokens(t, []string{*recvResp.Details[0].BrokerProperties.LockToken}, resp.FailedLockTokens)
	})

	t.Run("RejectCloudEvents", func(t *testing.T) {
		resp, err := receiver.RejectEvents(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		requireFailedLockTokens(t, []string{*recvResp.Details[0].BrokerProperties.LockToken}, resp.FailedLockTokens)
	})

	t.Run("ReleaseCloudEvents", func(t *testing.T) {
		resp, err := receiver.ReleaseEvents(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		requireFailedLockTokens(t, []string{*recvResp.Details[0].BrokerProperties.LockToken}, resp.FailedLockTokens)
	})
}

func TestPartialAckFailure(t *testing.T) {
	ce, err := messaging.NewCloudEvent("TestPartialAckFailure", "world", []byte("event one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	ce2, err := messaging.NewCloudEvent("TestPartialAckFailure", "world", []byte("event two"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	sender, receiver := newClients(t, false)

	sendResp, err := sender.SendEvents(context.Background(), []*messaging.CloudEvent{&ce, &ce2}, nil)
	require.NoError(t, err)
	require.Empty(t, sendResp)

	const numExpectedEvents = 2

	var events []aznamespaces.ReceiveDetails

	receiveCtx, cancelReceive := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelReceive()

	for len(events) < numExpectedEvents {
		eventsResp, err := receiver.ReceiveEvents(receiveCtx, &aznamespaces.ReceiveEventsOptions{
			MaxEvents: to.Ptr(int32(numExpectedEvents - len(events))),
		})
		require.NoError(t, err)
		events = append(events, eventsResp.Details...)
	}

	// we'll ack one now so we can force a failure to happen.
	ackResp, err := receiver.AcknowledgeEvents(context.Background(), []string{*events[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)

	// this will result in a partial failure.
	ackResp, err = receiver.AcknowledgeEvents(context.Background(), []string{
		*events[0].BrokerProperties.LockToken,
		*events[1].BrokerProperties.LockToken,
	}, nil)
	require.NoError(t, err)

	requireFailedLockTokens(t, []string{*events[0].BrokerProperties.LockToken}, ackResp.FailedLockTokens)
	require.Equal(t, []string{*events[1].BrokerProperties.LockToken}, ackResp.SucceededLockTokens)
}

func TestRejectEvents(t *testing.T) {

	ce, err := messaging.NewCloudEvent("TestAbandon", "world", []byte("event one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	sender, receiver := newClients(t, false)

	t.Logf("Publishing cloud events")
	_, err = sender.SendEvent(context.Background(), &ce, nil)
	require.NoError(t, err)

	t.Logf("Receiving cloud events")
	events, err := receiver.ReceiveEvents(context.Background(), nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		SpecVersion:     "1.0",
		DataContentType: to.Ptr("application/octet-stream"),
		Data:            []byte("event one"),
		Source:          "TestAbandon",
		Type:            "world",
	}, events.Details[0].Event)

	require.Equal(t, int32(1), *events.Details[0].BrokerProperties.DeliveryCount, "DeliveryCount starts at 1")

	t.Logf("Rejecting cloud events")
	rejectResp, err := receiver.RejectEvents(context.Background(), []string{*events.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, rejectResp.FailedLockTokens)
	t.Logf("Done rejecting cloud events")

	events, err = receiver.ReceiveEvents(context.Background(), &aznamespaces.ReceiveEventsOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10),
	})
	require.NoError(t, err)
	require.Empty(t, events.Details)
}

func TestReleaseEvents(t *testing.T) {
	ce, err := messaging.NewCloudEvent("TestRelease", "world", []byte("event one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	sender, receiver := newClients(t, false)

	_, err = sender.SendEvent(context.Background(), &ce, nil)
	require.NoError(t, err)

	events, err := receiver.ReceiveEvents(context.Background(), nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, ce, events.Details[0].Event)

	require.Equal(t, int32(1), *events.Details[0].BrokerProperties.DeliveryCount, "DeliveryCount starts at 1")

	rejectResp, err := receiver.ReleaseEvents(context.Background(), []string{*events.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)

	if len(rejectResp.FailedLockTokens) > 0 {
		for _, flt := range rejectResp.FailedLockTokens {
			t.Logf("FailedLockToken:\n  ec: %s\n  desc: %s\n  locktoken:%s", *flt.Error.Code, flt.Error.Error(), *flt.LockToken)
		}
		require.Fail(t, "Failed to release events")
	}

	require.Empty(t, rejectResp.FailedLockTokens)

	events, err = receiver.ReceiveEvents(context.Background(), nil)
	require.NoError(t, err)

	require.Equal(t, int32(2), *events.Details[0].BrokerProperties.DeliveryCount, "DeliveryCount is incremented")
	ackResp, err := receiver.AcknowledgeEvents(context.Background(), []string{*events.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
}

func TestPublishBytes(t *testing.T) {
	ce, err := messaging.NewCloudEvent("TestPublishBytes", "eventType", []byte("TestPublishBytes"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	sender, receiver := newClients(t, false)

	_, err = sender.SendEvent(context.Background(), &ce, nil)
	require.NoError(t, err)

	recvResp, err := receiver.ReceiveEvents(context.Background(), nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		Source:          "TestPublishBytes",
		SpecVersion:     "1.0",
		Type:            "eventType",
		Data:            []byte("TestPublishBytes"),
		DataContentType: to.Ptr("application/octet-stream"),
	}, recvResp.Details[0].Event)
}

func TestSendEventWithStringPayload(t *testing.T) {
	sender, receiver := newClients(t, false)

	ce, err := messaging.NewCloudEvent("TestPublishString", "eventType", "TestPublishString", &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/json"),
	})
	require.NoError(t, err)

	_, err = sender.SendEvent(context.Background(), &ce, nil)
	require.NoError(t, err)

	recvResp, err := receiver.ReceiveEvents(context.Background(), nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		Source:          "TestPublishString",
		SpecVersion:     "1.0",
		Type:            "eventType",
		Data:            []byte("\"TestPublishString\""), // non []byte returns as the JSON bytes
		DataContentType: to.Ptr("application/json"),
	}, recvResp.Details[0].Event)
}

func TestSendEventsAndReceiveEvents(t *testing.T) {
	sender, receiver := newClients(t, false)

	testData := []struct {
		Send     messaging.CloudEvent
		Expected messaging.CloudEvent
	}{
		{
			Send: mustCreateEvent(t, "TestPublishingAndReceivingMultipleCloudEvents", "eventType", []byte("TestPublishingAndReceivingMultipleCloudEvents 1"), &messaging.CloudEventOptions{
				DataContentType: to.Ptr("application/octet-stream"),
			}),
			Expected: messaging.CloudEvent{
				SpecVersion:     "1.0",
				Source:          "TestPublishingAndReceivingMultipleCloudEvents",
				Type:            "eventType",
				DataContentType: to.Ptr("application/octet-stream"),
				Data:            []byte("TestPublishingAndReceivingMultipleCloudEvents 1"),
			},
		},
		{
			Send: mustCreateEvent(t, "TestPublishingAndReceivingMultipleCloudEvents", "eventType", "TestPublishingAndReceivingMultipleCloudEvents 2", &messaging.CloudEventOptions{
				DataContentType: to.Ptr("application/json"),
				DataSchema:      to.Ptr("https://dataschema"),
				Extensions: map[string]any{
					"extension2": "extension2value",
				},
				Subject: to.Ptr("subject"),
			}),
			Expected: messaging.CloudEvent{
				SpecVersion:     "1.0",
				Source:          "TestPublishingAndReceivingMultipleCloudEvents",
				Type:            "eventType",
				DataSchema:      to.Ptr("https://dataschema"),
				Data:            []byte("\"TestPublishingAndReceivingMultipleCloudEvents 2\""),
				DataContentType: to.Ptr("application/json"),
				Subject:         to.Ptr("subject"),
				Extensions: map[string]any{
					"extension2": "extension2value",
				},
			},
		},
	}

	var sentEvents []*messaging.CloudEvent

	for _, td := range testData {
		e := td.Send
		sentEvents = append(sentEvents, &e)
	}

	sendResp, err := sender.SendEvents(context.Background(), sentEvents, nil)
	require.NoError(t, err)
	require.Empty(t, sendResp)

	resp, err := receiver.ReceiveEvents(context.Background(), &aznamespaces.ReceiveEventsOptions{
		MaxEvents:   to.Ptr(int32(len(sentEvents))),
		MaxWaitTime: to.Ptr(int32(60)),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Details)
	require.Equal(t, len(sentEvents), len(resp.Details))

	for i := 0; i < len(sentEvents); i++ {
		_, err := receiver.AcknowledgeEvents(context.Background(), []string{*resp.Details[i].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)

		requireEqualCloudEvent(t, testData[i].Expected, resp.Details[i].Event)
	}
}

func TestSimpleErrors(t *testing.T) {
	sender, _ := newClients(t, false)

	sendResp, err := sender.SendEvents(context.Background(), []*messaging.CloudEvent{
		{},
	}, nil)
	require.Empty(t, sendResp)

	// this'll lead to an error message like this:
	// "error": {
	// 	"code":"BadRequest",
	// 	"message":"`id` must be non-empty string",
	// 	"timestamp_utc":"2024-03-22T21:53:04.035019068+00:00",
	// 	"tracking_id":"<some tracking ID>"
	// }
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Equal(t, "BadRequest", respErr.ErrorCode)
	require.Equal(t, http.StatusBadRequest, respErr.StatusCode)
}

func TestRenewEventLocks(t *testing.T) {
	sender, receiver := newClients(t, false)

	ce := mustCreateEvent(t, "TestRenewCloudEventLocks", "eventType", "hello world", nil)
	_, err := sender.SendEvent(context.Background(), &ce, nil)
	require.NoError(t, err)

	recvResp, err := receiver.ReceiveEvents(context.Background(), nil)
	require.NoError(t, err)

	_, err = receiver.RenewEventLocks(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)

	ackResp, err := receiver.AcknowledgeEvents(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
}

func TestReleaseWithDelay(t *testing.T) {
	sender, receiver := newClients(t, false)

	ce := mustCreateEvent(t, "TestReleaseWithDelay", "eventType", "hello world", nil)
	_, err := sender.SendEvent(context.Background(), &ce, nil)
	require.NoError(t, err)

	recvResp, err := receiver.ReceiveEvents(context.Background(), nil)
	require.NoError(t, err)

	releaseResp, err := receiver.ReleaseEvents(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, &aznamespaces.ReleaseEventsOptions{
		ReleaseDelayInSeconds: to.Ptr(aznamespaces.ReleaseDelayTenSeconds),
	})
	require.NoError(t, err)
	require.Empty(t, releaseResp.FailedLockTokens)

	now := time.Now()

	// message will be available, but not immediately.
	recvResp, err = receiver.ReceiveEvents(context.Background(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, recvResp.Details)
	require.Equal(t, int32(2), *recvResp.Details[0].BrokerProperties.DeliveryCount)

	if recording.GetRecordMode() == recording.LiveMode {
		// doesn't work when recording but it's somewhat unimportant there.
		elapsed := time.Since(now)
		require.GreaterOrEqual(t, int(elapsed/time.Second), 8) // give a little wiggle room for potential delays between requests, etc...
	}

	ackResp, err := receiver.AcknowledgeEvents(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
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
