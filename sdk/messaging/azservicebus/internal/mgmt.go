// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
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

type mgmtError struct {
	Resp    *RPCResponse
	Message string
}

func (me mgmtError) Error() string {
	return me.Message
}

func (me mgmtError) RPCCode() int {
	return me.Resp.Code
}

// creates a new link and sends the RPC request, recovering and retrying on certain AMQP errors
func doRPC(ctx context.Context, name string, rpcLink RPCLink, msg *amqp.Message) (*RPCResponse, error) {
	res, err := rpcLink.RPC(ctx, msg)

	if err != nil {
		return nil, err
	}

	if res.Code >= 200 && res.Code < 300 {
		tab.For(ctx).Debug(fmt.Sprintf("rpc: success, status code %d and description: %s", res.Code, res.Description))
		return res, nil
	}

	return nil, mgmtError{
		Message: fmt.Sprintf("rpc: failed, status code %d and description: %s", res.Code, res.Description),
		Resp:    res,
	}
}

func ReceiveDeferred(ctx context.Context, rpcLink RPCLink, mode ReceiveMode, sequenceNumbers []int64) ([]*amqp.Message, error) {
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

	msg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:receive-by-sequence-number",
		},
		Value: values,
	}

	rsp, err := doRPC(ctx, "receiveDeferred", rpcLink, msg)
	if err != nil {
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

func PeekMessages(ctx context.Context, rpcLink RPCLink, fromSequenceNumber int64, messageCount int32) ([]*amqp.Message, error) {
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

	rsp, err := doRPC(ctx, "peek", rpcLink, msg)
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
func RenewLocks(ctx context.Context, rpcLink RPCLink, linkName string, lockTokens []amqp.UUID) ([]time.Time, error) {
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

	response, err := doRPC(ctx, "renewlocks", rpcLink, renewRequestMsg)

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
func RenewSessionLock(ctx context.Context, rpcLink RPCLink, sessionID string) (time.Time, error) {
	body := map[string]interface{}{
		"session-id": sessionID,
	}

	msg := &amqp.Message{
		Value: body,
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:renew-session-lock",
		},
	}

	resp, err := doRPC(ctx, "renewsessionlock", rpcLink, msg)

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
func GetSessionState(ctx context.Context, rpcLink RPCLink, sessionID string) ([]byte, error) {
	amqpMsg := &amqp.Message{
		Value: map[string]interface{}{
			"session-id": sessionID,
		},
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:get-session-state",
		},
	}

	resp, err := doRPC(ctx, "getsessionstate", rpcLink, amqpMsg)

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
func SetSessionState(ctx context.Context, rpcLink RPCLink, sessionID string, state []byte) error {
	uuid, err := uuid.New()

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

	resp, err := doRPC(ctx, "setsessionstate", rpcLink, amqpMsg)

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
func SendDisposition(ctx context.Context, rpcLink RPCLink, lockToken *amqp.UUID, state Disposition, propertiesToModify map[string]interface{}) error {
	ctx, span := tracing.StartConsumerSpanFromContext(ctx, tracing.SpanSendDisposition, Version)
	defer span.End()

	if lockToken == nil {
		err := errors.New("lock token on the message is not set, thus cannot send disposition")
		tab.For(ctx).Error(err)
		return err
	}

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

	if propertiesToModify != nil {
		value["properties-to-modify"] = propertiesToModify
	}

	msg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			"operation": "com.microsoft:update-disposition",
		},
		Value: value,
	}

	// no error, then it was successful
	_, err := doRPC(ctx, "senddisposition", rpcLink, msg)
	if err != nil {
		tab.For(ctx).Error(err)
		return err
	}

	return nil
}

// ScheduleMessages will send a batch of messages to a Queue, schedule them to be enqueued, and return the sequence numbers
// that can be used to cancel each message.
func ScheduleMessages(ctx context.Context, rpcLink RPCLink, enqueueTime time.Time, messages []*amqp.Message) ([]int64, error) {
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
			id, err := uuid.New()
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

		if messages[i].Properties.GroupID != nil {
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

	resp, err := doRPC(ctx, "schedule", rpcLink, msg)
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

// CancelScheduledMessages allows for removal of messages that have been handed to the Service Bus broker for later delivery,
// but have not yet ben enqueued.
func CancelScheduledMessages(ctx context.Context, rpcLink RPCLink, seq []int64) error {
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

	resp, err := doRPC(ctx, "cancelscheduled", rpcLink, msg)
	if err != nil {
		tab.For(ctx).Error(err)
		return err
	}

	if resp.Code != 200 {
		return ErrAMQP(*resp)
	}

	return nil
}
