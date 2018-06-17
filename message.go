package servicebus

//	MIT License
//
//	Copyright (c) Microsoft Corporation. All rights reserved.
//
//	Permission is hereby granted, free of charge, to any person obtaining a copy
//	of this software and associated documentation files (the "Software"), to deal
//	in the Software without restriction, including without limitation the rights
//	to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//	copies of the Software, and to permit persons to whom the Software is
//	furnished to do so, subject to the following conditions:
//
//	The above copyright notice and this permission notice shall be included in all
//	copies or substantial portions of the Software.
//
//	THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//	IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//	FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//	AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//	LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//	OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//	SOFTWARE

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/go-autorest/autorest/to"
	"github.com/mitchellh/mapstructure"
	"pack.ag/amqp"
)

type (
	// Message is an Service Bus message to be sent or received
	Message struct {
		ContentType      string
		CorrelationID    string
		Data             []byte
		DeliveryCount    uint32
		GroupID          *string
		GroupSequence    *uint32
		ID               string
		Label            string
		PartitionKey     string
		ReplyTo          string
		ReplyToGroupID   string
		To               string
		TTL              time.Duration
		LockToken        *string
		SystemProperties *SystemProperties
		UserProperties   map[string]interface{}
		message          *amqp.Message
	}

	// DispositionAction represents the action to notify Azure Service Bus of the Message's disposition
	DispositionAction func(ctx context.Context)

	// MessageErrorCondition represents a well-known collection of AMQP errors
	MessageErrorCondition string

	deliveryAnnotations struct {
		LockToken *amqp.UUID `mapstructure:"x-opt-lock-token"`
	}

	// SystemProperties are used to store properties that are set by the system.
	SystemProperties struct {
		LockedUntil            *time.Time `mapstructure:"x-opt-locked-until"`
		SequenceNumber         *int64     `mapstructure:"x-opt-sequence-number"`
		PartitionID            *int16     `mapstructure:"x-opt-partition-id"`
		PartitionKey           *string    `mapstructure:"x-opt-partition-key"`
		EnqueuedTime           *time.Time `mapstructure:"x-opt-enqueued-time"`
		DeadLetterSource       *string    `mapstructure:"x-opt-deadletter-source"`
		ScheduledEnqueuedTime  *time.Time `mapstructure:"x-opt-scheduled-enqueued-time"`
		EnqueuedSequenceNumber *int64     `mapstructure:"x-opt-enqueue-sequence-number"`
		ViaPartitionKey        *string    `mapstructure:"x-opt-via-partition-key"`
	}

	// MessageWithContext is a Service Bus message with its context which propagates the distributed trace information
	MessageWithContext struct {
		*Message
		Ctx context.Context
	}
)

// Error Conditions
const (
	ErrorInternalError         MessageErrorCondition = "amqp:internal-error"
	ErrorNotFound              MessageErrorCondition = "amqp:not-found"
	ErrorUnauthorizedAccess    MessageErrorCondition = "amqp:unauthorized-access"
	ErrorDecodeError           MessageErrorCondition = "amqp:decode-error"
	ErrorResourceLimitExceeded MessageErrorCondition = "amqp:resource-limit-exceeded"
	ErrorNotAllowed            MessageErrorCondition = "amqp:not-allowed"
	ErrorInvalidField          MessageErrorCondition = "amqp:invalid-field"
	ErrorNotImplemented        MessageErrorCondition = "amqp:not-implemented"
	ErrorResourceLocked        MessageErrorCondition = "amqp:resource-locked"
	ErrorPreconditionFailed    MessageErrorCondition = "amqp:precondition-failed"
	ErrorResourceDeleted       MessageErrorCondition = "amqp:resource-deleted"
	ErrorIllegalState          MessageErrorCondition = "amqp:illegal-state"
)

const (
	lockTokenName = "x-opt-lock-token"
)

// Complete will notify Azure Service Bus that the message was successfully handled and should be deleted from the queue
func (m *MessageWithContext) Complete() {
	m.Message.Complete()(m.Ctx)
}

// Abandon will notify Azure Service Bus the message failed but should be re-queued for delivery.
func (m *MessageWithContext) Abandon() {
	m.Message.Abandon()(m.Ctx)
}

// DeadLetter will notify Azure Service Bus the message failed and should not re-queued
func (m *MessageWithContext) DeadLetter(err error) {
	m.Message.DeadLetter(err)(m.Ctx)
}

// DeadLetterWithInfo will notify Azure Service Bus the message failed and should not be re-queued with additional context
func (m *MessageWithContext) DeadLetterWithInfo(err error, condition MessageErrorCondition, additionalData map[string]string) {
	m.Message.DeadLetterWithInfo(err, condition, additionalData)(m.Ctx)
}

// NewMessageFromString builds an Message from a string message
func NewMessageFromString(message string) *Message {
	return NewMessage([]byte(message))
}

// NewMessage builds an Message from a slice of data
func NewMessage(data []byte) *Message {
	return &Message{
		Data: data,
	}
}

// Complete will notify Azure Service Bus that the message was successfully handled and should be deleted from the queue
func (m *Message) Complete() DispositionAction {
	return func(ctx context.Context) {
		span, _ := m.startSpanFromContext(ctx, "sb.Message.Complete")
		defer span.Finish()

		m.message.Accept()
	}
}

// Abandon will notify Azure Service Bus the message failed but should be re-queued for delivery.
func (m *Message) Abandon() DispositionAction {
	return func(ctx context.Context) {
		span, _ := m.startSpanFromContext(ctx, "sb.Message.Abandon")
		defer span.Finish()

		m.message.Modify(false, false, nil)
	}
}

// TODO: Defer - will move to the "defer" queue and user will need to track the sequence number
// FailButRetryElsewhere will notify Azure Service Bus the message failed but should be re-queued for deliver to any
// other link but this one.
//func (m *Message) FailButRetryElsewhere() DispositionAction {
//	return func(ctx context.Context) {
//		span, _ := m.startSpanFromContext(ctx, "sb.Message.FailButRetryElsewhere")
//		defer span.Finish()
//
//		m.message.Modify(true, true, nil)
//	}
//}

// Release will notify Azure Service Bus the message should be re-queued without failure.
//func (m *Message) Release() DispositionAction {
//	return func(ctx context.Context) {
//		span, _ := m.startSpanFromContext(ctx, "sb.Message.Release")
//		defer span.Finish()
//
//		m.message.Release()
//	}
//}

// DeadLetter will notify Azure Service Bus the message failed and should not re-queued
func (m *Message) DeadLetter(err error) DispositionAction {
	return func(ctx context.Context) {
		span, _ := m.startSpanFromContext(ctx, "sb.Message.DeadLetter")
		defer span.Finish()

		amqpErr := amqp.Error{
			Condition:   amqp.ErrorCondition(ErrorInternalError),
			Description: err.Error(),
		}
		m.message.Reject(&amqpErr)
	}
}

// DeadLetterWithInfo will notify Azure Service Bus the message failed and should not be re-queued with additional
// context
func (m *Message) DeadLetterWithInfo(err error, condition MessageErrorCondition, additionalData map[string]string) DispositionAction {
	var info map[string]interface{}
	if additionalData != nil {
		info = make(map[string]interface{}, len(additionalData))
		for key, val := range additionalData {
			info[key] = val
		}
	}

	return func(ctx context.Context) {
		span, _ := m.startSpanFromContext(ctx, "sb.Message.DeadLetterWithInfo")
		defer span.Finish()

		amqpErr := amqp.Error{
			Condition:   amqp.ErrorCondition(condition),
			Description: err.Error(),
			Info:        info,
		}
		m.message.Reject(&amqpErr)
	}
}

// Set implements opentracing.TextMapWriter and sets properties on the event to be propagated to the message broker
func (m *Message) Set(key, value string) {
	if m.UserProperties == nil {
		m.UserProperties = make(map[string]interface{})
	}
	m.UserProperties[key] = value
}

// ForeachKey implements the opentracing.TextMapReader and gets properties on the event to be propagated from the message broker
func (m *Message) ForeachKey(handler func(key, val string) error) error {
	for key, value := range m.UserProperties {
		err := handler(key, value.(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Message) toMsg() *amqp.Message {
	amqpMsg := m.message
	if amqpMsg == nil {
		amqpMsg = amqp.NewMessage(m.Data)
	}

	amqpMsg.Properties = &amqp.MessageProperties{
		MessageID: m.ID,
	}

	if m.GroupID != nil && m.GroupSequence != nil {
		amqpMsg.Properties.GroupID = *m.GroupID
		amqpMsg.Properties.GroupSequence = *m.GroupSequence
	}

	amqpMsg.Properties.CorrelationID = m.CorrelationID
	amqpMsg.Properties.Subject = m.Label
	amqpMsg.Properties.To = m.To
	amqpMsg.Properties.ReplyTo = m.ReplyTo
	amqpMsg.Properties.ReplyToGroupID = m.ReplyToGroupID

	if len(m.UserProperties) > 0 {
		amqpMsg.ApplicationProperties = make(map[string]interface{})
		for key, value := range m.UserProperties {
			amqpMsg.ApplicationProperties[key] = value
		}
	}

	if m.LockToken != nil {
		if amqpMsg.DeliveryAnnotations == nil {
			amqpMsg.DeliveryAnnotations = make(amqp.Annotations)
		}
		amqpMsg.DeliveryAnnotations[lockTokenName] = m.LockToken
	}

	return amqpMsg
}

func messageFromAMQPMessage(msg *amqp.Message) (*Message, error) {
	return newMessage(msg.Data[0], msg)
}

func newMessage(data []byte, amqpMsg *amqp.Message) (*Message, error) {
	msg := &Message{
		Data:    data,
		message: amqpMsg,
	}

	if amqpMsg == nil {
		return msg, nil
	}

	if amqpMsg.Properties != nil {
		if id, ok := amqpMsg.Properties.MessageID.(string); ok {
			msg.ID = id
		}
		msg.GroupID = &amqpMsg.Properties.GroupID
		msg.GroupSequence = &amqpMsg.Properties.GroupSequence
		if id, ok := amqpMsg.Properties.CorrelationID.(string); ok {
			msg.CorrelationID = id
		}
		msg.ContentType = amqpMsg.Properties.ContentType
		msg.Label = amqpMsg.Properties.Subject
		msg.To = amqpMsg.Properties.To
		msg.ReplyTo = amqpMsg.Properties.ReplyTo
		msg.ReplyToGroupID = amqpMsg.Properties.ReplyToGroupID
		msg.DeliveryCount = amqpMsg.Header.DeliveryCount + 1
		msg.TTL = amqpMsg.Header.TTL
	}

	if amqpMsg.Annotations != nil {
		if err := mapstructure.Decode(amqpMsg.Annotations, &msg.SystemProperties); err != nil {
			return msg, err
		}
	}

	if amqpMsg.DeliveryAnnotations != nil {
		var da deliveryAnnotations
		if err := mapstructure.Decode(amqpMsg.DeliveryAnnotations, &da); err != nil {
			fmt.Println("ERROR!!", err.Error())
			return msg, err
		}
		if da.LockToken != nil {
			msg.LockToken = to.StringPtr(da.LockToken.String())
		}
	}

	return msg, nil
}
