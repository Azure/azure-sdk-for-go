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

// Set implements opentracing.TextMapWriter and sets properties on the event to be propagated to the message broker
func (e *Message) Set(key, value string) {
	if e.Properties == nil {
		e.Properties = make(map[string]interface{})
	}
	e.Properties[key] = value
}

// ForeachKey implements the opentracing.TextMapReader and gets properties on the event to be propagated from the message broker
func (e *Message) ForeachKey(handler func(key, val string) error) error {
	for key, value := range e.Properties {
		err := handler(key, value.(string))
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Message) toMsg() *amqp.Message {
	msg := e.message
	if msg == nil {
		msg = amqp.NewMessage(e.Data)
	}

	msg.Properties = &amqp.MessageProperties{
		MessageID: e.ID,
	}

	if e.GroupID != nil && e.GroupSequence != nil {
		msg.Properties.GroupID = *e.GroupID
		msg.Properties.GroupSequence = *e.GroupSequence
	}

	if len(e.Properties) > 0 {
		msg.ApplicationProperties = make(map[string]interface{})
		for key, value := range e.Properties {
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
