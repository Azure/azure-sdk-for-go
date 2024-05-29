//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package aznamespaces_test

import (
	"context"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/messaging"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/recording"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/eventgrid/aznamespaces"
	"github.com/stretchr/testify/require"
)

func TestClients_UsingSASKey(t *testing.T) {
	sender, receiver := newClients(t, true)

	ce, err := messaging.NewCloudEvent("source", "eventType", "hello world", nil)
	require.NoError(t, err)

	sendResp, err := sender.Send(context.Background(), &ce, nil)
	require.NoError(t, err)
	require.Empty(t, sendResp)

	recvResp, err := receiver.Receive(context.Background(), nil)
	require.NoError(t, err)

	require.Equal(t, 1, len(recvResp.Details))
	lockToken := recvResp.Details[0].BrokerProperties.LockToken
	require.NotEmpty(t, *lockToken)

	// strings are json serialized (if you want to send the string as bytes you can just []byte("your string")
	// when creating the CloudEvent)
	require.Equal(t, "\"hello world\"", string(recvResp.Details[0].Event.Data.([]byte)))

	ackResp, err := receiver.Acknowledge(context.Background(), []string{*lockToken}, nil)
	require.NoError(t, err)
	require.Equal(t, []string{*lockToken}, ackResp.SucceededLockTokens)
}

func TestFailedAck(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}

	ce, err := messaging.NewCloudEvent("TestFailedAck", "world", []byte("ack this one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	sender, receiver := newClients(t, false)

	_, err = sender.SendEvents(context.Background(), []*messaging.CloudEvent{&ce}, nil)
	require.NoError(t, err)

	recvResp, err := receiver.Receive(context.Background(), &aznamespaces.ReceiveOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10),
	})
	require.NoError(t, err)

	ackResp, err := receiver.Acknowledge(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
	require.Equal(t, []string{*recvResp.Details[0].BrokerProperties.LockToken}, ackResp.SucceededLockTokens)

	// now let's try to do stuff with an "out of date" token
	t.Run("Acknowledge", func(t *testing.T) {
		resp, err := receiver.Acknowledge(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		requireFailedLockTokens(t, []string{*recvResp.Details[0].BrokerProperties.LockToken}, resp.FailedLockTokens)
	})

	t.Run("RejectCloudEvents", func(t *testing.T) {
		resp, err := receiver.Reject(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		requireFailedLockTokens(t, []string{*recvResp.Details[0].BrokerProperties.LockToken}, resp.FailedLockTokens)
	})

	t.Run("ReleaseCloudEvents", func(t *testing.T) {
		resp, err := receiver.Release(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)
		require.Empty(t, resp.SucceededLockTokens)
		requireFailedLockTokens(t, []string{*recvResp.Details[0].BrokerProperties.LockToken}, resp.FailedLockTokens)
	})
}

func TestPartialAckFailure(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}

	ce, err := messaging.NewCloudEvent("TestPartialAckFailure", "world", []byte("event one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	ce2, err := messaging.NewCloudEvent("TestPartialAckFailure", "world", []byte("event two"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	sender, receiver := newClients(t, false)

	_, err = sender.SendEvents(context.Background(), []*messaging.CloudEvent{&ce, &ce2}, nil)
	require.NoError(t, err)

	const numExpectedEvents = 2

	var events []aznamespaces.ReceiveDetails

	receiveCtx, cancelReceive := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancelReceive()

	for len(events) < numExpectedEvents {
		eventsResp, err := receiver.Receive(receiveCtx, &aznamespaces.ReceiveOptions{
			MaxEvents: to.Ptr(int32(numExpectedEvents - len(events))),
		})
		require.NoError(t, err)
		events = append(events, eventsResp.Details...)
	}

	// we'll ack one now so we can force a failure to happen.
	ackResp, err := receiver.Acknowledge(context.Background(), []string{*events[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)

	// this will result in a partial failure.
	ackResp, err = receiver.Acknowledge(context.Background(), []string{
		*events[0].BrokerProperties.LockToken,
		*events[1].BrokerProperties.LockToken,
	}, nil)
	require.NoError(t, err)

	requireFailedLockTokens(t, []string{*events[0].BrokerProperties.LockToken}, ackResp.FailedLockTokens)
	require.Equal(t, []string{*events[1].BrokerProperties.LockToken}, ackResp.SucceededLockTokens)
}

func TestReject(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}

	ce, err := messaging.NewCloudEvent("TestAbandon", "world", []byte("event one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	sender, receiver := newClients(t, false)

	t.Logf("Publishing cloud events")
	_, err = sender.Send(context.Background(), &ce, nil)
	require.NoError(t, err)

	t.Logf("Receiving cloud events")
	events, err := receiver.Receive(context.Background(), nil)
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
	rejectResp, err := receiver.Reject(context.Background(), []string{*events.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, rejectResp.FailedLockTokens)
	t.Logf("Done rejecting cloud events")

	events, err = receiver.Receive(context.Background(), &aznamespaces.ReceiveOptions{
		MaxEvents:   to.Ptr[int32](1),
		MaxWaitTime: to.Ptr[int32](10),
	})
	require.NoError(t, err)
	require.Empty(t, events.Details)
}

func TestRelease(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}

	ce, err := messaging.NewCloudEvent("TestRelease", "world", []byte("event one"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	sender, receiver := newClients(t, false)

	_, err = sender.Send(context.Background(), &ce, nil)
	require.NoError(t, err)

	events, err := receiver.Receive(context.Background(), nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, ce, events.Details[0].Event)

	require.Equal(t, int32(1), *events.Details[0].BrokerProperties.DeliveryCount, "DeliveryCount starts at 1")

	rejectResp, err := receiver.Release(context.Background(), []string{*events.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)

	if len(rejectResp.FailedLockTokens) > 0 {
		for _, flt := range rejectResp.FailedLockTokens {
			t.Logf("FailedLockToken:\n  ec: %s\n  desc: %s\n  locktoken:%s", *flt.Error.Code, flt.Error.Error(), *flt.LockToken)
		}
		require.Fail(t, "Failed to release events")
	}

	require.Empty(t, rejectResp.FailedLockTokens)

	events, err = receiver.Receive(context.Background(), nil)
	require.NoError(t, err)

	require.Equal(t, int32(2), *events.Details[0].BrokerProperties.DeliveryCount, "DeliveryCount is incremented")
	ackResp, err := receiver.Acknowledge(context.Background(), []string{*events.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
}

func TestPublishBytes(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}

	ce, err := messaging.NewCloudEvent("TestPublishBytes", "eventType", []byte("TestPublishBytes"), &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	sender, receiver := newClients(t, false)

	_, err = sender.Send(context.Background(), &ce, nil)
	require.NoError(t, err)

	recvResp, err := receiver.Receive(context.Background(), nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		Source:          "TestPublishBytes",
		SpecVersion:     "1.0",
		Type:            "eventType",
		Data:            []byte("TestPublishBytes"),
		DataContentType: to.Ptr("application/octet-stream"),
	}, recvResp.Details[0].Event)
}

func TestPublishString(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	sender, receiver := newClients(t, false)

	ce, err := messaging.NewCloudEvent("TestPublishString", "eventType", "TestPublishString", &messaging.CloudEventOptions{
		DataContentType: to.Ptr("application/json"),
	})
	require.NoError(t, err)

	_, err = sender.Send(context.Background(), &ce, nil)
	require.NoError(t, err)

	recvResp, err := receiver.Receive(context.Background(), nil)
	require.NoError(t, err)

	requireEqualCloudEvent(t, messaging.CloudEvent{
		Source:          "TestPublishString",
		SpecVersion:     "1.0",
		Type:            "eventType",
		Data:            []byte("\"TestPublishString\""), // non []byte returns as the JSON bytes
		DataContentType: to.Ptr("application/json"),
	}, recvResp.Details[0].Event)
}

func TestPublishingAndReceivingMultipleCloudEvents(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
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

	_, err := sender.SendEvents(context.Background(), sentEvents, nil)
	require.NoError(t, err)

	resp, err := receiver.Receive(context.Background(), &aznamespaces.ReceiveOptions{
		MaxEvents:   to.Ptr(int32(len(sentEvents))),
		MaxWaitTime: to.Ptr(int32(60)),
	})
	require.NoError(t, err)
	require.NotEmpty(t, resp.Details)
	require.Equal(t, len(sentEvents), len(resp.Details))

	for i := 0; i < len(sentEvents); i++ {
		_, err := receiver.Acknowledge(context.Background(), []string{*resp.Details[i].BrokerProperties.LockToken}, nil)
		require.NoError(t, err)

		requireEqualCloudEvent(t, testData[i].Expected, resp.Details[i].Event)
	}
}

func TestSimpleErrors(t *testing.T) {
	sender, _ := newClients(t, false)

	_, err := sender.SendEvents(context.Background(), []*messaging.CloudEvent{
		{},
	}, nil)

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

func TestRenewCloudEventLocks(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	sender, receiver := newClients(t, false)

	ce := mustCreateEvent(t, "TestRenewCloudEventLocks", "eventType", "hello world", nil)
	_, err := sender.Send(context.Background(), &ce, nil)
	require.NoError(t, err)

	recvResp, err := receiver.Receive(context.Background(), nil)
	require.NoError(t, err)

	_, err = receiver.RenewLocks(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)

	ackResp, err := receiver.Acknowledge(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
}

func TestReleaseWithDelay(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	sender, receiver := newClients(t, false)

	ce := mustCreateEvent(t, "TestReleaseWithDelay", "eventType", "hello world", nil)
	_, err := sender.Send(context.Background(), &ce, nil)
	require.NoError(t, err)

	recvResp, err := receiver.Receive(context.Background(), nil)
	require.NoError(t, err)

	releaseResp, err := receiver.Release(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, &aznamespaces.ReleaseOptions{
		ReleaseDelayInSeconds: to.Ptr(aznamespaces.ReleaseDelayTenSeconds),
	})
	require.NoError(t, err)
	require.Empty(t, releaseResp.FailedLockTokens)

	now := time.Now()

	// message will be available, but not immediately.
	recvResp, err = receiver.Receive(context.Background(), nil)
	require.NoError(t, err)
	require.NotEmpty(t, recvResp.Details)
	require.Equal(t, int32(2), *recvResp.Details[0].BrokerProperties.DeliveryCount)

	if recording.GetRecordMode() == recording.LiveMode {
		// doesn't work when recording but it's somewhat unimportant there.
		elapsed := time.Since(now)
		require.GreaterOrEqual(t, int(elapsed/time.Second), 8) // give a little wiggle room for potential delays between requests, etc...
	}

	ackResp, err := receiver.Acknowledge(context.Background(), []string{*recvResp.Details[0].BrokerProperties.LockToken}, nil)
	require.NoError(t, err)
	require.Empty(t, ackResp.FailedLockTokens)
}

func TestSend_binaryMode(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}

	sender, receiver := newClients(t, false)

	// you can send a variety of different payloads, all of which can be encoded by messaging.CloudEvent
	binaryPayload := []byte{1, 2, 3}
	customContentType := "application/customcontenttype"
	eventToSend, err := messaging.NewCloudEvent("TestSend_binaryMode", "eventType", binaryPayload, &messaging.CloudEventOptions{
		// this will be used as the "Content-Type" for the request.
		DataContentType: &customContentType,
	})
	require.NoError(t, err)
	fixCloudEvent(t, &eventToSend)

	// want to validate that we're actually doing the binary mode encoding
	var actualResp *http.Response
	captureCtx := policy.WithCaptureResponse(context.Background(), &actualResp)

	// binary mode publish (CloudEvent attributes encoded as headers, CloudEvent.Data used as the request body)
	_, err = sender.Send(captureCtx, &eventToSend, &aznamespaces.SendOptions{
		BinaryMode: true,
	})
	require.NoError(t, err)

	require.NotNil(t, actualResp.Request)
	// binary mode propagates the content type of the event itself.
	// require.Equal(t, customContentType, actualResp.Request.Header["Content-Type"][0])

	// (the body bytes stream is exhausted so we can't compare the contents but we can check that we only sent the body bytes.
	require.Equal(t, int64(len(eventToSend.Data.([]byte))), actualResp.Request.ContentLength)

	event := receiveAll(t, receiver, 1)[0]

	// zero out fields that always change
	eventToSend.Time, event.Event.Time = nil, nil

	require.Equal(t, eventToSend, event.Event)
}

// send an event using the two methods (binary and non-binary mode) and make sure they both
// produce the same content, when received.
func TestSend_worksSameInBinaryVsNonBinaryMode(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	sender, receiver := newClients(t, false)

	// you can send a variety of different payloads, all of which can be encoded by messaging.CloudEvent
	binaryPayload := []byte{1, 2, 3}
	customContentType := "application/customcontenttype"
	eventToSend, err := messaging.NewCloudEvent("source", "eventType", binaryPayload, &messaging.CloudEventOptions{
		// this will be used as the "Content-Type" for the request.
		DataContentType: &customContentType,
	})
	require.NoError(t, err)
	fixCloudEvent(t, &eventToSend)

	// want to validate that we're actually doing the binary mode encoding
	var actualResp *http.Response
	captureCtx := policy.WithCaptureResponse(context.Background(), &actualResp)

	// binary mode publish (CloudEvent attributes encoded as headers, CloudEvent.Data used as the request body)
	{
		_, err = sender.Send(captureCtx, &eventToSend, &aznamespaces.SendOptions{
			BinaryMode: true,
		})
		require.NoError(t, err)

		require.NotNil(t, actualResp.Request)
		// binary mode propagates the content type of the event itself.
		require.Equal(t, customContentType, actualResp.Request.Header["Content-Type"][0])

		// (the body bytes stream is exhausted so we can't compare the contents but we can check that we only sent the body bytes.
		require.Equal(t, int64(len(eventToSend.Data.([]byte))), actualResp.Request.ContentLength)
	}

	// non-binary mode publish (ie, entire CloudEvent encoded as JSON as the body)
	{
		_, err = sender.Send(captureCtx, &eventToSend, &aznamespaces.SendOptions{
			// Default is binary mode OFF.
			//BinaryMode: true,
		})
		require.NoError(t, err)

		require.NotNil(t, actualResp.Request)
		require.Equal(t, "application/cloudevents+json; charset=utf-8", actualResp.Request.Header["Content-Type"][0])
		require.Less(t, int64(len(eventToSend.Data.([]byte))), actualResp.Request.ContentLength, "the body is a JSON blob of the entire event, it will be larger than the payload I attempted to send")
	}

	// receive the two events
	events := receiveAll(t, receiver, 2)

	// events sent via binary content mode and non-binary-content-mode should be the same.
	// it only affects transport, Event Grid should persist them the same.

	// scrub out fields that change no matter what.
	for i := 0; i < len(events); i++ {
		require.NotEmpty(t, events[i].BrokerProperties.LockToken)
		events[i].BrokerProperties.LockToken = nil

		require.NotZero(t, events[i].Event.Time)
		events[i].Event.Time = &time.Time{}
	}

	require.Equal(t, events[0], events[1])
}

func TestSend_binaryModeNonByteSlicePayloadFails(t *testing.T) {
	sender, _ := newClients(t, false)

	eventToSend, err := messaging.NewCloudEvent("source", "eventType", "you can't send strings (or any non []byte) using binary mode!", &messaging.CloudEventOptions{
		// this will be used as the "Content-Type" for the request.
		DataContentType: to.Ptr("application/octet-stream"),
	})
	require.NoError(t, err)

	_, err = sender.Send(context.Background(), &eventToSend, &aznamespaces.SendOptions{
		BinaryMode: true,
	})
	require.EqualError(t, err, "CloudEvent.Data must be of type []byte, was type string")
}

func TestSend_binaryModeNoContentTypeFails(t *testing.T) {
	sender, _ := newClients(t, false)

	eventToSend, err := messaging.NewCloudEvent("source", "eventType", []byte{1, 2, 3}, nil)
	require.NoError(t, err)
	fixCloudEvent(t, &eventToSend)

	_, err = sender.Send(context.Background(), &eventToSend, &aznamespaces.SendOptions{
		BinaryMode: true,
	})
	var respErr *azcore.ResponseError
	require.ErrorAs(t, err, &respErr)
	require.Contains(t, err.Error(), "CONTENT-TYPE is not provided")
	require.Equal(t, "BadRequest", respErr.ErrorCode)
}

func TestSend_binaryModeUseOptionalValues(t *testing.T) {
	if recording.GetRecordMode() == recording.PlaybackMode {
		t.Skip("https://github.com/Azure/azure-sdk-for-go/issues/22869")
	}
	sender, receiver := newClients(t, false)

	tm, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	require.NoError(t, err)

	arbitraryURL, err := url.Parse("https://microsoft.com/someschema")
	require.NoError(t, err)

	binaryPayload := []byte{1, 2, 3}
	customContentType := "application/customcontenttype"

	eventToSend, err := messaging.NewCloudEvent("TestSend_binaryModeUseOptionalValues", "eventType", binaryPayload, &messaging.CloudEventOptions{
		DataContentType: &customContentType,
		Extensions: map[string]any{
			"extensiondatastring":  "hello",
			"extensiondatastring2": stringableType{},
			"extensiondataint":     101,
			"extensiondataurl":     arbitraryURL,
			"extensiondatauint":    uint(202),
			"extensiondatatime":    tm,
			"extensiondatabytes":   []byte{4, 5, 6, 7},
		},
		DataSchema: to.Ptr("https://microsoft.com"),
		Subject:    to.Ptr("my subject"),
	})
	require.NoError(t, err)
	fixCloudEvent(t, &eventToSend)

	// want to validate that we're actually doing the binary mode encoding
	var actualResp *http.Response
	captureCtx := policy.WithCaptureResponse(context.Background(), &actualResp)

	// binary mode publish (CloudEvent attributes encoded as headers, CloudEvent.Data used as the request body)
	_, err = sender.Send(captureCtx, &eventToSend, &aznamespaces.SendOptions{
		BinaryMode: true,
	})
	require.NoError(t, err)

	req := actualResp.Request
	require.NotNil(t, req)

	message := receiveAll(t, receiver, 1)[0]

	require.Equal(t, customContentType, *message.Event.DataContentType)
	require.Equal(t, []byte{1, 2, 3}, message.Event.Data)
	require.Equal(t, "https://microsoft.com", *message.Event.DataSchema)
	require.Equal(t, "TestSend_binaryModeUseOptionalValues", message.Event.Source)
	require.Equal(t, "1.0", message.Event.SpecVersion)
	require.Equal(t, "my subject", *message.Event.Subject)

	require.Equal(t, map[string]any{
		"extensiondatastring":  "hello",
		"extensiondatastring2": stringableType{}.String(),
		"extensiondataint":     "101",
		"extensiondataurl":     arbitraryURL.String(),
		"extensiondatauint":    "202",
		"extensiondatatime":    tm.Format(time.RFC3339),
		"extensiondatabytes":   "BAUGBw==", // byte data comes back as a base64 string
	}, message.Event.Extensions)
}

func TestSend_binaryModeNoContentType(t *testing.T) {
	sender, _ := newClients(t, false)

	binaryPayload := []byte{1, 2, 3}
	eventToSend, err := messaging.NewCloudEvent("source", "eventType", binaryPayload, nil)
	require.NoError(t, err)

	fixCloudEvent(t, &eventToSend)

	// want to validate that we're actually doing the binary mode encoding
	var actualResp *http.Response
	captureCtx := policy.WithCaptureResponse(context.Background(), &actualResp)

	_, err = sender.Send(captureCtx, &eventToSend, &aznamespaces.SendOptions{
		BinaryMode: true,
	})
	var resp *azcore.ResponseError
	require.ErrorAs(t, err, &resp)
	require.Equal(t, "BadRequest", resp.ErrorCode)
	require.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestSend_binaryModeUnstringableExtension(t *testing.T) {
	sender, _ := newClients(t, false)

	binaryPayload := []byte{1, 2, 3}
	eventToSend, err := messaging.NewCloudEvent("source", "eventType", binaryPayload, &messaging.CloudEventOptions{
		Extensions: map[string]any{
			"wontwork": aznamespaces.ReceiverClient{},
		},
	})
	require.NoError(t, err)

	_, err = sender.Send(context.Background(), &eventToSend, &aznamespaces.SendOptions{
		BinaryMode: true,
	})
	require.EqualError(t, err, "type aznamespaces.ReceiverClient cannot be converted to a string")
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

type stringableType struct{}

func (st stringableType) String() string {
	return "hello"
}

func fixCloudEvent(t *testing.T, ce *messaging.CloudEvent) {
	// (this matches what our sanitizers do)
	ce.ID = "00000000-0000-0000-0000-000000000000"
	tm, err := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	require.NoError(t, err)
	ce.Time = &tm
}
