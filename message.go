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

	"pack.ag/amqp"
)

type (
	// Message is an Service Bus message to be sent or received
	Message struct {
		Data          []byte
		Properties    map[string]interface{}
		ID            string
		GroupID       *string
		GroupSequence *uint32
		message       *amqp.Message
	}

	// DispositionAction represents the action to notify Azure Service Bus of the Message's disposition
	DispositionAction func(ctx context.Context)

	// MessageErrorCondition represents a well-known collection of AMQP errors
	MessageErrorCondition string
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

// Accept will notify Azure Service Bus that the message was successfully handled and should be deleted from the queue
func (m *Message) Accept() DispositionAction {
	return func(ctx context.Context) {
		span, _ := m.startSpanFromContext(ctx, "sb.Message.Accept")
		defer span.Finish()

		m.message.Accept()
	}
}

// FailButRetry will notify Azure Service Bus the message failed but should be re-queued for delivery.
func (m *Message) FailButRetry() DispositionAction {
	return func(ctx context.Context) {
		span, _ := m.startSpanFromContext(ctx, "sb.Message.FailButRetry")
		defer span.Finish()

		m.message.Modify(true, false, nil)
	}
}

// FailButRetryElsewhere will notify Azure Service Bus the message failed but should be re-queued for deliver to any
// other link but this one.
func (m *Message) FailButRetryElsewhere() DispositionAction {
	return func(ctx context.Context) {
		span, _ := m.startSpanFromContext(ctx, "sb.Message.FailButRetryElsewhere")
		defer span.Finish()

		m.message.Modify(true, true, nil)
	}
}

// Release will notify Azure Service Bus the message should be re-queued without failure.
func (m *Message) Release() DispositionAction {
	return func(ctx context.Context) {
		span, _ := m.startSpanFromContext(ctx, "sb.Message.Release")
		defer span.Finish()

		m.message.Release()
	}
}

// Reject will notify Azure Service Bus the message failed and should not re-queued
func (m *Message) Reject(err error) DispositionAction {
	return func(ctx context.Context) {
		span, _ := m.startSpanFromContext(ctx, "sb.Message.Reject")
		defer span.Finish()

		amqpErr := amqp.Error{
			Condition:   amqp.ErrorCondition(ErrorInternalError),
			Description: err.Error(),
		}
		m.message.Reject(&amqpErr)
	}
}

// RejectWithInfo will notify Azure Service Bus the message failed and should not be re-queued with additional context
func (m *Message) RejectWithInfo(err error, condition MessageErrorCondition, additionalData map[string]string) DispositionAction {
	var info map[string]interface{}
	if additionalData != nil {
		info = make(map[string]interface{}, len(additionalData))
		for key, val := range additionalData {
			info[key] = val
		}
	}

	return func(ctx context.Context) {
		span, _ := m.startSpanFromContext(ctx, "sb.Message.RejectWithInfo")
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
	if m.Properties == nil {
		m.Properties = make(map[string]interface{})
	}
	m.Properties[key] = value
}

// ForeachKey implements the opentracing.TextMapReader and gets properties on the event to be propagated from the message broker
func (m *Message) ForeachKey(handler func(key, val string) error) error {
	for key, value := range m.Properties {
		err := handler(key, value.(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Message) toMsg() *amqp.Message {
	msg := m.message
	if msg == nil {
		msg = amqp.NewMessage(m.Data)
	}

	msg.Properties = &amqp.MessageProperties{
		MessageID: m.ID,
	}

	if m.GroupID != nil && m.GroupSequence != nil {
		msg.Properties.GroupID = *m.GroupID
		msg.Properties.GroupSequence = *m.GroupSequence
	}

	if len(m.Properties) > 0 {
		msg.ApplicationProperties = make(map[string]interface{})
		for key, value := range m.Properties {
			msg.ApplicationProperties[key] = value
		}
	}

	return msg
}

func messageFromAMQPMessage(msg *amqp.Message) *Message {
	return newMessage(msg.Data[0], msg)
}

func newMessage(data []byte, msg *amqp.Message) *Message {
	message := &Message{
		Data:    data,
		message: msg,
	}

	if msg.Properties != nil {
		if id, ok := msg.Properties.MessageID.(string); ok {
			message.ID = id
		}
		message.GroupID = &msg.Properties.GroupID
		message.GroupSequence = &msg.Properties.GroupSequence
	}

	if msg != nil {
		message.Properties = make(map[string]interface{})
		for key, value := range msg.ApplicationProperties {
			message.Properties[key] = value
		}
	}
	return message
}
