// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/uuid"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/exported"
	"github.com/Azure/go-amqp"
)

// defaultServerTimeout bounds a management RPC attempt when the caller's context
// has no deadline. Without it the service is given no bound at all, so a stalled
// broker holds the call until the AMQP link itself fails.
// See https://github.com/Azure/azure-sdk-for-go/issues/26421.
const defaultServerTimeout = 60 * time.Second

// serverTimeoutBuffer is taken off the caller's remaining deadline so the broker's
// timeout fires before the context does. Sending the full remaining time makes the
// two expire together, and the context wins by the round trip the reply still needs,
// so the caller sees "context deadline exceeded" rather than the broker's own
// timeout. The Java SDK subtracts the same second in
// MessageUtils.adjustServerTimeout ("Pass little less than client timeout to the
// server so client doesn't time out before server times out"), though it subtracts
// from a fixed per-attempt timeout that is never itself under a second, where this
// subtracts from whatever the caller has left. See serverTimeoutMillis for what
// that difference means below one second.
const serverTimeoutBuffer = time.Second

// serverTimeoutMillis returns the server-timeout, in milliseconds, for a management
// operation. With a deadline it is the remaining time less serverTimeoutBuffer,
// clamped at zero, so the broker answers first and the caller gets a service-side
// timeout it can act on. Without a deadline it is defaultServerTimeout.
//
// Under a second of remaining time there is no room to answer first, so the value
// clamps to zero and the buffer simply stops helping. That is not a regression: the
// caller's own context still ends the call within that same second, exactly as it
// does today.
//
// This is a broker-side bound, not a client-side one. rpcLink.RPC waits only for a
// response or ctx.Done(), so the value asks the broker to answer within the window
// but never caps the client's own wait. A call with no deadline therefore depends on
// the broker honoring it: MaxRetries+1 attempts of this length plus backoff when it
// does, and no bound at all if the broker goes silent with the link still open,
// since the loop counts retries after the initial attempt. A client-side per-attempt
// TryTimeout is tracked in https://github.com/Azure/azure-sdk-for-go/issues/27269.
func serverTimeoutMillis(ctx context.Context) uint {
	timeout := defaultServerTimeout

	if deadline, ok := ctx.Deadline(); ok {
		timeout = time.Until(deadline) - serverTimeoutBuffer
		if timeout < 0 {
			// Less than the buffer remains, so the caller is about to give up anyway.
			// Clamp to avoid an unsigned underflow.
			timeout = 0
		}
	}

	return uint(timeout / time.Millisecond)
}

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

func ReceiveDeferred(ctx context.Context, rpcLink amqpwrap.RPCLink, linkName string, mode exported.ReceiveMode, sequenceNumbers []int64) ([]*amqp.Message, error) {
	const messagesField, messageField = "messages", "message"

	backwardsMode := uint32(0)
	if mode == exported.PeekLock {
		backwardsMode = 1
	}

	values := map[string]any{
		"sequence-numbers":     sequenceNumbers,
		"receiver-settle-mode": uint32(backwardsMode), // pick up messages with peek lock
	}

	msg := &amqp.Message{
		ApplicationProperties: map[string]any{
			"operation": "com.microsoft:receive-by-sequence-number",
		},
		Value: values,
	}

	addAssociatedLinkName(linkName, msg)

	rsp, err := rpcLink.RPC(ctx, msg)
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
	val, ok := rsp.Message.Value.(map[string]any)
	if !ok {
		return nil, NewErrIncorrectType(messageField, map[string]any{}, rsp.Message.Value)
	}

	rawMessages, ok := val[messagesField]
	if !ok {
		return nil, ErrMissingField(messagesField)
	}

	messages, ok := rawMessages.([]any)
	if !ok {
		return nil, NewErrIncorrectType(messagesField, []any{}, rawMessages)
	}

	transformedMessages := make([]*amqp.Message, len(messages))
	for i := range messages {
		rawEntry, ok := messages[i].(map[string]any)
		if !ok {
			return nil, NewErrIncorrectType(messageField, map[string]any{}, messages[i])
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

func PeekMessages(ctx context.Context, rpcLink amqpwrap.RPCLink, linkName string, fromSequenceNumber int64, messageCount int32) ([]*amqp.Message, error) {
	const messagesField, messageField = "messages", "message"

	msg := &amqp.Message{
		ApplicationProperties: map[string]any{
			"operation": "com.microsoft:peek-message",
		},
		Value: map[string]any{
			"from-sequence-number": fromSequenceNumber,
			"message-count":        messageCount,
		},
	}

	addAssociatedLinkName(linkName, msg)

	// server-timeout is left to rpcLink.RPC, which sets the bare key for any
	// com.microsoft: operation that has not chosen one. Computing it there keeps
	// the value closer to the send, so it reflects the time actually left.

	rsp, err := rpcLink.RPC(ctx, msg)
	if err != nil {
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
	val, ok := rsp.Message.Value.(map[string]any)
	if !ok {
		err = NewErrIncorrectType(messageField, map[string]any{}, rsp.Message.Value)
		return nil, err
	}

	rawMessages, ok := val[messagesField]
	if !ok {
		err = ErrMissingField(messagesField)
		return nil, err
	}

	messages, ok := rawMessages.([]any)
	if !ok {
		err = NewErrIncorrectType(messagesField, []any{}, rawMessages)
		return nil, err
	}

	transformedMessages := make([]*amqp.Message, len(messages))
	for i := range messages {
		rawEntry, ok := messages[i].(map[string]any)
		if !ok {
			err = NewErrIncorrectType(messageField, map[string]any{}, messages[i])
			return nil, err
		}

		rawMessage, ok := rawEntry[messageField]
		if !ok {
			err = ErrMissingField(messageField)
			return nil, err
		}

		marshaled, ok := rawMessage.([]byte)
		if !ok {
			err = new(ErrMalformedMessage)
			return nil, err
		}

		var rehydrated amqp.Message
		err = rehydrated.UnmarshalBinary(marshaled)
		if err != nil {
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
func RenewLocks(ctx context.Context, rpcLink amqpwrap.RPCLink, linkName string, lockTokens []amqp.UUID) ([]time.Time, error) {
	renewRequestMsg := &amqp.Message{
		ApplicationProperties: map[string]any{
			"operation": "com.microsoft:renew-lock",
		},
		Value: map[string]any{
			"lock-tokens": lockTokens,
		},
	}

	addAssociatedLinkName(linkName, renewRequestMsg)

	response, err := rpcLink.RPC(ctx, renewRequestMsg)

	if err != nil {
		return nil, err
	}

	if response.Code != 200 {
		err := fmt.Errorf("error renewing locks: %v", response.Description)
		return nil, err
	}

	// extract the new lock renewal times from the response
	// response.Message.

	val, ok := response.Message.Value.(map[string]any)
	if !ok {
		return nil, NewErrIncorrectType("Message.Value", map[string]any{}, response.Message.Value)
	}

	expirations, ok := val["expirations"]

	if !ok {
		return nil, NewErrIncorrectType("Message.Value[\"expirations\"]", map[string]any{}, response.Message.Value)
	}

	asTimes, ok := expirations.([]time.Time)

	if !ok {
		return nil, NewErrIncorrectType("Message.Value[\"expirations\"] as times", map[string]any{}, response.Message.Value)
	}

	return asTimes, nil
}

// RenewSessionLocks renews a session lock.
func RenewSessionLock(ctx context.Context, rpcLink amqpwrap.RPCLink, linkName string, sessionID string) (time.Time, error) {
	body := map[string]any{
		"session-id": sessionID,
	}

	msg := &amqp.Message{
		Value: body,
		ApplicationProperties: map[string]any{
			"operation": "com.microsoft:renew-session-lock",
		},
	}

	addAssociatedLinkName(linkName, msg)

	resp, err := rpcLink.RPC(ctx, msg)

	if err != nil {
		return time.Time{}, err
	}

	m, ok := resp.Message.Value.(map[string]any)

	if !ok {
		return time.Time{}, NewErrIncorrectType("Message.Value", map[string]any{}, resp.Message.Value)
	}

	lockedUntil, ok := m["expiration"].(time.Time)

	if !ok {
		return time.Time{}, NewErrIncorrectType("Message.Value[\"expiration\"] as times", time.Time{}, resp.Message.Value)
	}

	return lockedUntil, nil
}

// GetSessionState retrieves state associated with the session.
func GetSessionState(ctx context.Context, rpcLink amqpwrap.RPCLink, linkName string, sessionID string) ([]byte, error) {
	amqpMsg := &amqp.Message{
		Value: map[string]any{
			"session-id": sessionID,
		},
		ApplicationProperties: map[string]any{
			"operation": "com.microsoft:get-session-state",
		},
	}

	addAssociatedLinkName(linkName, amqpMsg)

	resp, err := rpcLink.RPC(ctx, amqpMsg)

	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, ErrAMQP(*resp)
	}

	asMap, ok := resp.Message.Value.(map[string]any)

	if !ok {
		return nil, NewErrIncorrectType("Value", map[string]any{}, resp.Message.Value)
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
func SetSessionState(ctx context.Context, rpcLink amqpwrap.RPCLink, linkName string, sessionID string, state []byte) error {
	uuid, err := uuid.New()

	if err != nil {
		return err
	}

	amqpMsg := &amqp.Message{
		Value: map[string]any{
			"session-id":    sessionID,
			"session-state": state,
		},
		ApplicationProperties: map[string]any{
			"operation":                 "com.microsoft:set-session-state",
			"com.microsoft:tracking-id": uuid.String(),
		},
	}

	addAssociatedLinkName(linkName, amqpMsg)

	resp, err := rpcLink.RPC(ctx, amqpMsg)

	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return ErrAMQP(*resp)
	}

	return nil
}

// SettleOnMgmtLink allows you settle a message using the management link, rather than via your
// *amqp.Receiver. Use this if the receiver has been closed/lost or if the message isn't associated
// with a link (ex: deferred messages).
func SettleOnMgmtLink(ctx context.Context, rpcLink amqpwrap.RPCLink, linkName string, lockToken *amqp.UUID, state Disposition, propertiesToModify map[string]any) error {
	if lockToken == nil {
		err := errors.New("lock token on the message is not set, thus cannot send disposition")
		return err
	}

	value := map[string]any{
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
		ApplicationProperties: map[string]any{
			"operation": "com.microsoft:update-disposition",
		},
		Value: value,
	}

	addAssociatedLinkName(linkName, msg)

	// no error, then it was successful
	_, err := rpcLink.RPC(ctx, msg)
	if err != nil {
		return err
	}

	return nil
}

// ScheduleMessages will send a batch of messages to a Queue, schedule them to be enqueued, and return the sequence numbers
// that can be used to cancel each message.
func ScheduleMessages(ctx context.Context, rpcLink amqpwrap.RPCLink, linkName string, enqueueTime time.Time, messages []*amqp.Message) ([]int64, error) {
	if len(messages) == 0 {
		return nil, errors.New("expected one or more messages")
	}

	transformed := make([]any, 0, len(messages))
	enqueueTimeAsUTC := enqueueTime.UTC()

	for i := range messages {
		if messages[i].Annotations == nil {
			return nil, errors.New("message does not have an initialized Annotations property")
		}

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

		individualMessage := map[string]any{
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
		ApplicationProperties: map[string]any{
			"operation": "com.microsoft:schedule-message",
		},
		Value: map[string]any{
			"messages": transformed,
		},
	}

	addAssociatedLinkName(linkName, msg)

	msg.ApplicationProperties["com.microsoft:server-timeout"] = serverTimeoutMillis(ctx)

	resp, err := rpcLink.RPC(ctx, msg)
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, ErrAMQP(*resp)
	}

	retval := make([]int64, 0, len(messages))
	if rawVal, ok := resp.Message.Value.(map[string]any); ok {
		const sequenceFieldName = "sequence-numbers"
		if rawArr, ok := rawVal[sequenceFieldName]; ok {
			if arr, ok := rawArr.([]int64); ok {
				retval = append(retval, arr...)
				return retval, nil
			}
			return nil, NewErrIncorrectType(sequenceFieldName, []int64{}, rawArr)
		}
		return nil, ErrMissingField(sequenceFieldName)
	}
	return nil, NewErrIncorrectType("value", map[string]any{}, resp.Message.Value)
}

// CancelScheduledMessages allows for removal of messages that have been handed to the Service Bus broker for later delivery,
// but have not yet ben enqueued.
func CancelScheduledMessages(ctx context.Context, rpcLink amqpwrap.RPCLink, linkName string, seq []int64) error {
	msg := &amqp.Message{
		ApplicationProperties: map[string]any{
			"operation": "com.microsoft:cancel-scheduled-message",
		},
		Value: map[string]any{
			"sequence-numbers": seq,
		},
	}

	addAssociatedLinkName(linkName, msg)

	msg.ApplicationProperties["com.microsoft:server-timeout"] = serverTimeoutMillis(ctx)

	resp, err := rpcLink.RPC(ctx, msg)
	if err != nil {
		return err
	}

	if resp.Code != 200 {
		return ErrAMQP(*resp)
	}

	return nil
}

// addAssociatedLinkName adds the 'associated-link-name' application
// property to the AMQP message. Setting this property associates
// management link activity with a sender or receiver link, which can
// prevent it from idling out.
func addAssociatedLinkName(linkName string, msg *amqp.Message) {
	if linkName != "" {
		msg.ApplicationProperties["associated-link-name"] = linkName
	}
}

// GetMessageSessions lists session IDs from an entity. The behavior depends on lastUpdatedTime:
//   - Pass time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC) to list sessions with active
//     messages, as well as sessions that have session state set but no active messages.
//     This encodes to 253402300800000 ms on the AMQP wire, which the service's
//     .NET AMQP decoder clamps to DateTime.MaxValue, triggering default listing mode.
//   - Pass a real timestamp to list sessions whose session state was updated after that time.
//
// No associated link name is needed (entity-level operation).
func GetMessageSessions(ctx context.Context, rpcLink amqpwrap.RPCLink, lastUpdatedTime time.Time, skip int32, top int32) ([]string, error) {
	msg := &amqp.Message{
		Value: map[string]any{
			"last-updated-time": lastUpdatedTime,
			"skip":              skip,
			"top":               top,
		},
		ApplicationProperties: map[string]any{
			"operation": "com.microsoft:get-message-sessions",
		},
	}
	// No associated link name for entity-level operations

	resp, err := rpcLink.RPC(ctx, msg)
	if err != nil {
		// rpcLink.RPC converts every non-2xx status into an RPCError with the response
		// attached, so a 404 arrives here as an error rather than as resp.Code == 404. A
		// 404 with a "session not found" description means the entity exists but has no
		// sessions - treat it as an empty result. Other 404s (e.g. entity not found) and
		// all other errors surface so callers can distinguish "no sessions" from "the
		// entity does not exist".
		var rpcErr RPCError
		if errors.As(err, &rpcErr) && rpcErr.Resp != nil &&
			rpcErr.Resp.Code == 404 && isSessionNotFound(rpcErr.Resp.Description) {
			return nil, nil
		}
		return nil, err
	}

	if resp.Code == 204 {
		return nil, nil
	}

	if resp.Code != 200 {
		return nil, ErrAMQP(*resp)
	}

	asMap, ok := resp.Message.Value.(map[string]any)
	if !ok {
		return nil, NewErrIncorrectType("Value", map[string]any{}, resp.Message.Value)
	}

	sessionsObj := asMap["sessions-ids"]
	if sessionsObj == nil {
		return nil, nil
	}

	// The sessions-ids field is a string array from the service
	switch v := sessionsObj.(type) {
	case []string:
		return v, nil
	case []any:
		result := make([]string, len(v))
		for i, s := range v {
			sessionID, ok := s.(string)
			if !ok {
				return nil, NewErrIncorrectType("Value['sessions-ids']", []string{}, sessionsObj)
			}
			result[i] = sessionID
		}
		return result, nil
	default:
		return nil, NewErrIncorrectType("Value['sessions-ids']", []string{}, sessionsObj)
	}
}

// isSessionNotFound reports whether an AMQP error description indicates the
// "no sessions for this entity" condition for the get-message-sessions
// operation. The service emits the description as human-readable prose
// (e.g. "Session not found") rather than a structured condition code, so the
// match is case-insensitive and tolerates either the camel-case token form
// ("SessionNotFound") or the spaced prose form ("session not found").
func isSessionNotFound(description string) bool {
	lower := strings.ToLower(description)
	return strings.Contains(lower, "sessionnotfound") ||
		strings.Contains(lower, "session not found")
}
