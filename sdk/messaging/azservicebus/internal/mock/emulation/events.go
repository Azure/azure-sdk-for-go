// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"fmt"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
)

type EventType string

const (
	EventTypeConnOpen  EventType = "conn-open"
	EventTypeConnClose EventType = "conn-close"

	EventTypeLinkOpen   EventType = "link-attach"
	EventTypeLinkCredit EventType = "link-issue-credit"
	EventTypeLinkClose  EventType = "link-detach"

	EventTypeLinkDisposition EventType = "link-disposition"

	EventTypeReceive EventType = "link-receive"
	EventTypeSend    EventType = "link-send"
)

type Event struct {
	Type EventType
	Data any
}

type LinkRole string

const (
	LinkRoleReceiver LinkRole = "receiver"
	LinkRoleSender   LinkRole = "sender"
)

type LinkEvent struct {
	ConnID        string
	SessID        string
	Entity        string
	Name          string
	Role          LinkRole
	TargetAddress string
}

type CreditEvent struct {
	LinkEvent
	Credit uint32
}

type SendEvent struct {
	LinkEvent
	*amqp.Message
}

type ReceiveEvent struct {
	LinkEvent
	*amqp.Message
}

type DispositionType string

const (
	DispositionTypeAccept  DispositionType = "accept"
	DispositionTypeReject  DispositionType = "reject"
	DispositionTypeRelease DispositionType = "release"
)

type DispositionEvent struct {
	LinkEvent
	DispType DispositionType
	*amqp.Message
}

type Events struct {
	mu  sync.Mutex
	all []Event
	ch  chan Event
}

func NewEvents() *Events {
	return &Events{
		ch: make(chan Event, 10000),
	}
}

func (e *Events) Chan() <-chan Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.ch
}

func (e *Events) All() []Event {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.all
}

func (e *Events) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.all = nil
}

func (e *Events) addEvent(evt Event) {
	e.mu.Lock()
	e.all = append(e.all, evt)
	e.mu.Unlock()

	e.ch <- evt
}

func (e *Events) Custom(name string, data any) {
	e.addEvent(Event{
		Type: EventType(name),
		Data: data,
	})
}

func (e *Events) OpenLink(args LinkEvent) {
	e.addEvent(Event{
		Type: EventTypeLinkOpen,
		Data: args,
	})
}

func (e *Events) CloseLink(args LinkEvent) {
	e.addEvent(Event{
		Type: EventTypeLinkClose,
		Data: args,
	})
}

func (e *Events) IssueCredit(args CreditEvent) {
	e.addEvent(Event{
		Type: EventTypeLinkCredit,
		Data: args,
	})
}

func (e *Events) Send(args SendEvent) {
	e.addEvent(Event{
		Type: EventTypeSend,
		Data: args,
	})
}

func (e *Events) Receive(args ReceiveEvent) {
	e.addEvent(Event{
		Type: EventTypeReceive,
		Data: args,
	})
}

func (e *Events) Disposition(args DispositionEvent) {
	e.addEvent(Event{
		Type: EventTypeLinkDisposition,
		Data: args,
	})
}

func (e *Events) OpenConnection(name string) {
	e.addEvent(Event{
		Type: EventTypeConnOpen,
		Data: name,
	})
}

func (e *Events) CloseConnection(name string) {
	e.addEvent(Event{
		Type: EventTypeConnClose,
		Data: name,
	})
}

func (e *Events) GetOpenLinks() []string {
	opened := map[string]bool{}

	for _, evt := range e.all {
		switch evt.Type {
		case EventTypeLinkOpen:
			opened[evt.Data.(LinkEvent).Name] = true
		case EventTypeLinkClose:
			n := evt.Data.(LinkEvent).Name

			if !opened[n] {
				panic(fmt.Sprintf("We closed a link that was never open: %s", n))
			}

			opened[n] = false
		}
	}

	var leaks []string

	for k := range opened {
		if opened[k] {
			leaks = append(leaks, k)
		}
	}

	return leaks
}

func (e *Events) GetOpenConns() []string {
	opened := map[string]bool{}

	for _, evt := range e.all {
		switch evt.Type {
		case EventTypeConnOpen:
			opened[evt.Data.(string)] = true
		case EventTypeConnClose:
			n := evt.Data.(string)

			if !opened[n] {
				panic(fmt.Sprintf("We closed a connection that was never open: %s", n))
			}

			opened[n] = false
		}
	}

	var leaks []string

	for k := range opened {
		if opened[k] {
			leaks = append(leaks, k)
		}
	}

	return leaks
}
