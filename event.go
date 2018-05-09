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
	// Event is an Event Hubs message to be sent or received
	Event struct {
		Data          []byte
		Properties    map[string]interface{}
		ID            string
		GroupID       *string
		GroupSequence *uint32
		message       *amqp.Message
	}

	// EventBatch is a batch of Event Hubs messages to be sent
	EventBatch struct {
		Events     []*Event
		Properties map[string]interface{}
		ID         string
	}
)

// NewEventFromString builds an Event from a string message
func NewEventFromString(message string) *Event {
	return NewEvent([]byte(message))
}

// NewEvent builds an Event from a slice of data
func NewEvent(data []byte) *Event {
	return &Event{
		Data: data,
	}
}

// NewEventBatch builds an EventBatch from an array of Events
func NewEventBatch(events []*Event) *EventBatch {
	return &EventBatch{
		Events: events,
	}
}

// Set implements opentracing.TextMapWriter and sets properties on the event to be propagated to the message broker
func (e *Event) Set(key, value string) {
	if e.Properties == nil {
		e.Properties = make(map[string]interface{})
	}
	e.Properties[key] = value
}

// ForeachKey implements the opentracing.TextMapReader and gets properties on the event to be propagated from the message broker
func (e *Event) ForeachKey(handler func(key string, val interface{}) error) error {
	for key, value := range e.Properties {
		err := handler(key, value)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Event) toMsg() *amqp.Message {
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

func eventFromMsg(msg *amqp.Message) *Event {
	return newEvent(msg.Data[0], msg)
}

func newEvent(data []byte, msg *amqp.Message) *Event {
	event := &Event{
		Data:    data,
		message: msg,
	}

	if msg.Properties != nil {
		if id, ok := msg.Properties.MessageID.(string); ok {
			event.ID = id
		}
		event.GroupID = &msg.Properties.GroupID
		event.GroupSequence = &msg.Properties.GroupSequence
	}

	if msg != nil {
		event.Properties = msg.ApplicationProperties
	}
	return event
}
