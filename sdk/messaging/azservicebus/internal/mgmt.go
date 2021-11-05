// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	common "github.com/Azure/azure-amqp-common-go/v3"
	"github.com/Azure/azure-amqp-common-go/v3/rpc"
	"github.com/Azure/azure-amqp-common-go/v3/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/go-amqp"
	"github.com/devigned/tab"
)

type Disposition struct {
	Status                DispositionStatus
	LockTokens            []*uuid.UUID
	DeadLetterReason      *string
	DeadLetterDescription *string
}

type DispositionStatus string

const (
	CompletedDisposition DispositionStatus = "completed"
	AbandonedDisposition DispositionStatus = "abandoned"
	SuspendedDisposition DispositionStatus = "suspended"
	DeferredDisposition  DispositionStatus = "defered"
)

type (
	mgmtClient struct {
		ns    NamespaceForMgmtClient
		links AMQPLinks

		clientMu sync.RWMutex
		rpcLink  RPCLink

		sessionID          *string
		isSessionFilterSet bool
	}
)

type MgmtClient interface {
	Close(ctx context.Context) error
	SendDisposition(ctx context.Context, lockToken *amqp.UUID, state Disposition) error
	ReceiveDeferred(ctx context.Context, mode ReceiveMode, sequenceNumbers []int64) ([]*amqp.Message, error)
	PeekMessages(ctx context.Context, fromSequenceNumber int64, messageCount int32) ([]*amqp.Message, error)

	ScheduleMessages(ctx context.Context, enqueueTime time.Time, messages ...*amqp.Message) ([]int64, error)
	CancelScheduled(ctx context.Context, seq ...int64) error

	RenewLocks(ctx context.Context, linkName string, lockTokens []amqp.UUID) ([]time.Time, error)
	RenewSessionLock(ctx context.Context, sessionID string) (time.Time, error)

	GetSessionState(ctx context.Context, sessionID string) ([]byte, error)
	SetSessionState(ctx context.Context, sessionID string, state []byte) error
}

func newMgmtClient(ctx context.Context, links AMQPLinks, ns NamespaceForMgmtClient) (MgmtClient, error) {
	r := &mgmtClient{
		ns:    ns,
		links: links,
	}

	return r, nil
}

// Recover will attempt to close the current session and link, then rebuild them
func (mc *mgmtClient) recover(ctx context.Context) error {
	mc.clientMu.Lock()
	defer mc.clientMu.Unlock()

	ctx, span := mc.startSpanFromContext(ctx, string(tracing.SpanNameRecover))
	defer span.End()

	if mc.rpcLink != nil {
		if err := mc.rpcLink.Close(ctx); err != nil {
			tab.For(ctx).Debug(fmt.Sprintf("Error while closing old link in recovery: %s", err.Error()))
		}
		mc.rpcLink = nil
	}

	if _, err := mc.getLinkWithoutLock(ctx); err != nil {
		return err
	}

	return nil
}

// getLinkWithoutLock returns the currently cached link (or creates a new one)
func (mc *mgmtClient) getLinkWithoutLock(ctx context.Context) (RPCLink, error) {
	if mc.rpcLink != nil {
		return mc.rpcLink, nil
	}

	var err error
	mc.rpcLink, err = mc.ns.NewRPCLink(ctx, mc.links.ManagementPath())

	if err != nil {
		return nil, err
	}

	return mc.rpcLink, nil
}

// Close will close the AMQP connection
func (mc *mgmtClient) Close(ctx context.Context) error {
	mc.clientMu.Lock()
	defer mc.clientMu.Unlock()

	if mc.rpcLink == nil {
		return nil
	}

	err := mc.rpcLink.Close(ctx)
	mc.rpcLink = nil
	return err
}

// creates a new link and sends the RPC request, recovering and retrying on certain AMQP errors
func (mc *mgmtClient) doRPCWithRetry(ctx context.Context, msg *amqp.Message, times int, delay time.Duration, opts ...rpc.LinkOption) (*rpc.Response, error) {
	// track the number of times we attempt to perform the RPC call.
	// this is to avoid a potential infinite loop if the returned error
	// is always transient and Recover() doesn't fail.
	sendCount := 0

	for {
		mc.clientMu.RLock()
		rpcLink, err := mc.getLinkWithoutLock(ctx)
		mc.clientMu.RUnlock()

		var rsp *rpc.Response

		if err == nil {
			rsp, err = rpcLink.RetryableRPC(ctx, times, delay, msg)

			if err == nil {
				return rsp, err
			}
		}

		if sendCount >= amqpRetryDefaultTimes || !isAMQPTransientError(ctx, err) {
			return nil, err
		}
		sendCount++
		// if we get here, recover and try again
		tab.For(ctx).Debug("recovering RPC connection")

		_, retryErr := common.Retry(amqpRetryDefaultTimes, amqpRetryDefaultDelay, func() (interface{}, error) {
			ctx, sp := mc.startProducerSpanFromContext(ctx, string(tracing.SpanTryRecover))
			defer sp.End()

			if err := mc.recover(ctx); err == nil {
				tab.For(ctx).Debug("recovered RPC connection")
				return nil, nil
			}
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				return nil, common.Retryable(err.Error())
			}
		})

		if retryErr != nil {
			tab.For(ctx).Debug("RPC recovering retried, but error was unrecoverable")
			return nil, retryErr
		}
	}
}

// returns true if the AMQP error is considered transient
func isAMQPTransientError(ctx context.Context, err error) bool {
	// always retry on a detach error
	var amqpDetach *amqp.DetachError
	if errors.As(err, &amqpDetach) {
		return true
	}
	// for an AMQP error, only retry depending on the condition
	var amqpErr *amqp.Error
	if errors.As(err, &amqpErr) {
		switch amqpErr.Condition {
		case errorServerBusy, errorTimeout, errorOperationCancelled, errorContainerClose:
			return true
		default:
			tab.For(ctx).Debug(fmt.Sprintf("isAMQPTransientError: condition %s is not transient", amqpErr.Condition))
			return false
		}
	}
	tab.For(ctx).Debug(fmt.Sprintf("isAMQPTransientError: %T is not transient", err))
	return false
}

func (mc *mgmtClient) ReceiveDeferred(ctx context.Context, mode ReceiveMode, sequenceNumbers []int64) ([]*amqp.Message, error) {
	ctx, span := tracing.StartConsumerSpanFromContext(ctx, tracing.SpanReceiveDeferred, Version)
	defer span.End()

	const messagesField, messageField = "messages", "message"

	backwardsMode := uint32(0)
	if mode == PeekLock {
		backwardsMode = 1
	}

	values := map[string]interface{}{
		"sequence-numbers":     sequenceNumbers,
		"receiver-settle-mode": uint32(backwardsMode), // pick up messages with peek lock
	}

	var opts []rpc.LinkOption
	if mc.isSessionFilterSet {
		opts = append(opts, rpc.LinkWithSessionFilter(mc.sessionID))
		values["session-id"] = mc.sessionID
	}

	msg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:receive-by-sequence-number",
		},
		Value: values,
	}

	rsp, err := mc.doRPCWithRetry(ctx, msg, 5, 5*time.Second, opts...)
	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	if rsp.Code == 204 {
		return nil, ErrNoMessages{}
	}

	// Deferred messages come back in a relatively convoluted manner:
	// a map (always with one key: "messages")
	// 	of arrays
	// 		of maps (always with one key: "message")
	// 			of an array with raw encoded Service Bus messages
	val, ok := rsp.Message.Value.(map[string]interface{})
	if !ok {
		return nil, NewErrIncorrectType(messageField, map[string]interface{}{}, rsp.Message.Value)
	}

	rawMessages, ok := val[messagesField]
	if !ok {
		return nil, ErrMissingField(messagesField)
	}

	messages, ok := rawMessages.([]interface{})
	if !ok {
		return nil, NewErrIncorrectType(messagesField, []interface{}{}, rawMessages)
	}

	transformedMessages := make([]*amqp.Message, len(messages))
	for i := range messages {
		rawEntry, ok := messages[i].(map[string]interface{})
		if !ok {
			return nil, NewErrIncorrectType(messageField, map[string]interface{}{}, messages[i])
		}

		rawMessage, ok := rawEntry[messageField]
		if !ok {
			return nil, ErrMissingField(messageField)
		}

		marshaled, ok := rawMessage.([]byte)
		if !ok {
			return nil, new(ErrMalformedMessage)
		}

		var rehydrated amqp.Message
		err = rehydrated.UnmarshalBinary(marshaled)
		if err != nil {
			return nil, err
		}

		transformedMessages[i] = &rehydrated
	}

	return transformedMessages, nil
}

func (mc *mgmtClient) PeekMessages(ctx context.Context, fromSequenceNumber int64, messageCount int32) ([]*amqp.Message, error) {
	ctx, span := tracing.StartConsumerSpanFromContext(ctx, tracing.SpanPeekFromSequenceNumber, Version)
	defer span.End()

	const messagesField, messageField = "messages", "message"

	msg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:peek-message",
		},
		Value: map[string]interface{}{
			"from-sequence-number": fromSequenceNumber,
			"message-count":        messageCount,
		},
	}

	if deadline, ok := ctx.Deadline(); ok {
		msg.ApplicationProperties["server-timeout"] = uint(time.Until(deadline) / time.Millisecond)
	}

	rsp, err := mc.doRPCWithRetry(ctx, msg, 5, 5*time.Second)
	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	if rsp.Code == 204 {
		// no messages available
		return nil, nil
	}

	// Peeked messages come back in a relatively convoluted manner:
	// a map (always with one key: "messages")
	// 	of arrays
	// 		of maps (always with one key: "message")
	// 			of an array with raw encoded Service Bus messages
	val, ok := rsp.Message.Value.(map[string]interface{})
	if !ok {
		err = NewErrIncorrectType(messageField, map[string]interface{}{}, rsp.Message.Value)
		tab.For(ctx).Error(err)
		return nil, err
	}

	rawMessages, ok := val[messagesField]
	if !ok {
		err = ErrMissingField(messagesField)
		tab.For(ctx).Error(err)
		return nil, err
	}

	messages, ok := rawMessages.([]interface{})
	if !ok {
		err = NewErrIncorrectType(messagesField, []interface{}{}, rawMessages)
		tab.For(ctx).Error(err)
		return nil, err
	}

	transformedMessages := make([]*amqp.Message, len(messages))
	for i := range messages {
		rawEntry, ok := messages[i].(map[string]interface{})
		if !ok {
			err = NewErrIncorrectType(messageField, map[string]interface{}{}, messages[i])
			tab.For(ctx).Error(err)
			return nil, err
		}

		rawMessage, ok := rawEntry[messageField]
		if !ok {
			err = ErrMissingField(messageField)
			tab.For(ctx).Error(err)
			return nil, err
		}

		marshaled, ok := rawMessage.([]byte)
		if !ok {
			err = new(ErrMalformedMessage)
			tab.For(ctx).Error(err)
			return nil, err
		}

		var rehydrated amqp.Message
		err = rehydrated.UnmarshalBinary(marshaled)
		if err != nil {
			tab.For(ctx).Error(err)
			return nil, err
		}

		transformedMessages[i] = &rehydrated

		// transformedMessages[i], err = MessageFromAMQPMessage(&rehydrated)
		// if err != nil {
		// 	tab.For(ctx).Error(err)
		// 	return nil, err
		// }

		// transformedMessages[i].useSession = r.isSessionFilterSet
		// transformedMessages[i].sessionID = r.sessionID
	}

	// This sort is done to ensure that folks wanting to peek messages in sequence order may do so.
	// sort.Slice(transformedMessages, func(i, j int) bool {
	// 	iSeq := *transformedMessages[i].SystemProperties.SequenceNumber
	// 	jSeq := *transformedMessages[j].SystemProperties.SequenceNumber
	// 	return iSeq < jSeq
	// })

	return transformedMessages, nil
}

// RenewLocks renews the locks in a single 'com.microsoft:renew-lock' operation.
// NOTE: this function assumes all the messages received on the same link.
func (mc *mgmtClient) RenewLocks(ctx context.Context, linkName string, lockTokens []amqp.UUID) ([]time.Time, error) {
	ctx, span := tracing.StartConsumerSpanFromContext(ctx, tracing.SpanRenewLock, Version)
	defer span.End()

	renewRequestMsg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:renew-lock",
		},
		Value: map[string]interface{}{
			"lock-tokens": lockTokens,
		},
	}

	if linkName != "" {
		renewRequestMsg.ApplicationProperties["associated-link-name"] = linkName
	}

	response, err := mc.doRPCWithRetry(ctx, renewRequestMsg, 3, 1*time.Second)

	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	if response.Code != 200 {
		err := fmt.Errorf("error renewing locks: %v", response.Description)
		tab.For(ctx).Error(err)
		return nil, err
	}

	// extract the new lock renewal times from the response
	// response.Message.

	val, ok := response.Message.Value.(map[string]interface{})
	if !ok {
		return nil, NewErrIncorrectType("Message.Value", map[string]interface{}{}, response.Message.Value)
	}

	expirations, ok := val["expirations"]

	if !ok {
		return nil, NewErrIncorrectType("Message.Value[\"expirations\"]", map[string]interface{}{}, response.Message.Value)
	}

	asTimes, ok := expirations.([]time.Time)

	if !ok {
		return nil, NewErrIncorrectType("Message.Value[\"expirations\"] as times", map[string]interface{}{}, response.Message.Value)
	}

	return asTimes, nil
}

// RenewSessionLocks renews a session lock.
func (mc *mgmtClient) RenewSessionLock(ctx context.Context, sessionID string) (time.Time, error) {
	body := map[string]interface{}{
		"session-id": sessionID,
	}

	msg := &amqp.Message{
		Value: body,
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:renew-session-lock",
		},
	}

	resp, err := mc.doRPCWithRetry(ctx, msg, 5, 5*time.Second)

	if err != nil {
		return time.Time{}, err
	}

	m, ok := resp.Message.Value.(map[string]interface{})

	if !ok {
		return time.Time{}, NewErrIncorrectType("Message.Value", map[string]interface{}{}, resp.Message.Value)
	}

	lockedUntil, ok := m["expiration"].(time.Time)

	if !ok {
		return time.Time{}, NewErrIncorrectType("Message.Value[\"expiration\"] as times", time.Time{}, resp.Message.Value)
	}

	return lockedUntil, nil
}

// GetSessionState retrieves state associated with the session.
func (mc *mgmtClient) GetSessionState(ctx context.Context, sessionID string) ([]byte, error) {
	amqpMsg := &amqp.Message{
		Value: map[string]interface{}{
			"session-id": sessionID,
		},
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:get-session-state",
		},
	}

	resp, err := mc.doRPCWithRetry(ctx, amqpMsg, 5, 5*time.Second)

	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, ErrAMQP(*resp)
	}

	asMap, ok := resp.Message.Value.(map[string]interface{})

	if !ok {
		return nil, NewErrIncorrectType("Value", map[string]interface{}{}, resp.Message.Value)
	}

	val := asMap["session-state"]

	if val == nil {
		// no session state set
		return nil, nil
	}

	asBytes, ok := val.([]byte)

	if !ok {
		return nil, NewErrIncorrectType("Value['session-state']", []byte{}, asMap["session-state"])
	}

	return asBytes, nil
}

// SetSessionState sets the state associated with the session.
func (mc *mgmtClient) SetSessionState(ctx context.Context, sessionID string, state []byte) error {
	uuid, err := uuid.NewV4()

	if err != nil {
		return err
	}

	amqpMsg := &amqp.Message{
		Value: map[string]interface{}{
			"session-id":    sessionID,
			"session-state": state,
		},
		ApplicationProperties: map[string]interface{}{
			"operation":                 "com.microsoft:set-session-state",
			"com.microsoft:tracking-id": uuid.String(),
		},
	}

	resp, err := mc.doRPCWithRetry(ctx, amqpMsg, 5, 5*time.Second)

	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return ErrAMQP(*resp)
	}

	return nil
}

// SendDisposition allows you settle a message using the management link, rather than via your
// *amqp.Receiver. Use this if the receiver has been closed/lost or if the message isn't associated
// with a link (ex: deferred messages).
func (mc *mgmtClient) SendDisposition(ctx context.Context, lockToken *amqp.UUID, state Disposition) error {
	ctx, span := tracing.StartConsumerSpanFromContext(ctx, tracing.SpanSendDisposition, Version)
	defer span.End()

	if lockToken == nil {
		err := errors.New("lock token on the message is not set, thus cannot send disposition")
		tab.For(ctx).Error(err)
		return err
	}

	var opts []rpc.LinkOption
	value := map[string]interface{}{
		"disposition-status": string(state.Status),
		"lock-tokens":        []amqp.UUID{*lockToken},
	}

	if state.DeadLetterReason != nil {
		value["deadletter-reason"] = state.DeadLetterReason
	}

	if state.DeadLetterDescription != nil {
		value["deadletter-description"] = state.DeadLetterDescription
	}

	msg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:update-disposition",
		},
		Value: value,
	}

	// no error, then it was successful
	_, err := mc.doRPCWithRetry(ctx, msg, 5, 5*time.Second, opts...)
	if err != nil {
		tab.For(ctx).Error(err)
		return err
	}

	return nil
}

// ScheduleMessages will send a batch of messages to a Queue, schedule them to be enqueued, and return the sequence numbers
// that can be used to cancel each message.
func (mc *mgmtClient) ScheduleMessages(ctx context.Context, enqueueTime time.Time, messages ...*amqp.Message) ([]int64, error) {
	ctx, span := tracing.StartConsumerSpanFromContext(ctx, tracing.SpanScheduleMessage, Version)
	defer span.End()

	if len(messages) <= 0 {
		return nil, errors.New("expected one or more messages")
	}

	transformed := make([]interface{}, 0, len(messages))
	enqueueTimeAsUTC := enqueueTime.UTC()

	for i := range messages {
		// TODO: don't like that we're modifying the underlying message here
		messages[i].Annotations["x-opt-scheduled-enqueue-time"] = &enqueueTimeAsUTC

		if messages[i].Properties == nil {
			messages[i].Properties = &amqp.MessageProperties{}
		}

		// TODO: this is in two spots now (in Message, and here). Since this one
		// could potentially take the raw AMQP message we need to check it, and we assume
		// that 'nil' is the only zero value that matters.
		if messages[i].Properties.MessageID == nil {
			id, err := uuid.NewV4()
			if err != nil {
				return nil, err
			}
			messages[i].Properties.MessageID = id.String()
		}

		encoded, err := messages[i].MarshalBinary()
		if err != nil {
			return nil, err
		}

		individualMessage := map[string]interface{}{
			"message-id": messages[i].Properties.MessageID,
			"message":    encoded,
		}

		// TODO: I believe empty string should be allowed here. There isn't a way for the
		// user to opt out of session related information.
		if messages[i].Properties.GroupID != "" {
			individualMessage["session-id"] = messages[i].Properties.GroupID
		}

		if value, ok := messages[i].Annotations["x-opt-partition-key"]; ok {
			individualMessage["partition-key"] = value.(string)
		}

		if value, ok := messages[i].Annotations["x-opt-via-partition-key"]; ok {
			individualMessage["via-partition-key"] = value.(string)
		}

		transformed = append(transformed, individualMessage)
	}

	msg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:schedule-message",
		},
		Value: map[string]interface{}{
			"messages": transformed,
		},
	}

	if deadline, ok := ctx.Deadline(); ok {
		msg.ApplicationProperties["com.microsoft:server-timeout"] = uint(time.Until(deadline) / time.Millisecond)
	}

	resp, err := mc.doRPCWithRetry(ctx, msg, 5, 5*time.Second)
	if err != nil {
		tab.For(ctx).Error(err)
		return nil, err
	}

	if resp.Code != 200 {
		return nil, ErrAMQP(*resp)
	}

	retval := make([]int64, 0, len(messages))
	if rawVal, ok := resp.Message.Value.(map[string]interface{}); ok {
		const sequenceFieldName = "sequence-numbers"
		if rawArr, ok := rawVal[sequenceFieldName]; ok {
			if arr, ok := rawArr.([]int64); ok {
				for i := range arr {
					retval = append(retval, arr[i])
				}
				return retval, nil
			}
			return nil, NewErrIncorrectType(sequenceFieldName, []int64{}, rawArr)
		}
		return nil, ErrMissingField(sequenceFieldName)
	}
	return nil, NewErrIncorrectType("value", map[string]interface{}{}, resp.Message.Value)
}

// CancelScheduled allows for removal of messages that have been handed to the Service Bus broker for later delivery,
// but have not yet ben enqueued.
func (mc *mgmtClient) CancelScheduled(ctx context.Context, seq ...int64) error {
	ctx, span := tracing.StartConsumerSpanFromContext(ctx, tracing.SpanCancelScheduledMessage, Version)
	defer span.End()

	msg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:cancel-scheduled-message",
		},
		Value: map[string]interface{}{
			"sequence-numbers": seq,
		},
	}

	if deadline, ok := ctx.Deadline(); ok {
		msg.ApplicationProperties["com.microsoft:server-timeout"] = uint(time.Until(deadline) / time.Millisecond)
	}

	resp, err := mc.doRPCWithRetry(ctx, msg, 5, 5*time.Second)
	if err != nil {
		tab.For(ctx).Error(err)
		return err
	}

	if resp.Code != 200 {
		return ErrAMQP(*resp)
	}

	return nil
}

func (mc *mgmtClient) startSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tracing.StartConsumerSpanFromContext(ctx, operationName, Version)
	span.AddAttributes(tab.StringAttribute("message_bus.destination", mc.links.ManagementPath()))
	return ctx, span
}

func (mc *mgmtClient) startProducerSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	tracing.ApplyComponentInfo(span, Version)
	span.AddAttributes(
		tab.StringAttribute("span.kind", "producer"),
		tab.StringAttribute("message_bus.destination", mc.links.ManagementPath()),
	)
	return ctx, span
}
