// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
)

type Operation struct {
	Op       string
	Entity   string
	LinkName string
	Role     string

	Credits uint32
	M       *amqp.Message
}

func (op *Operation) String() string {
	data := ""

	if op.M != nil && op.M.ApplicationProperties != nil {
		opAny, exists := op.M.ApplicationProperties["operation"]

		if exists {
			data = fmt.Sprintf("m.ap.op:%s", opAny.(string))
		}
	}

	return fmt.Sprintf("e:%s l:%s c:%d %s", op.Entity, op.LinkName, op.Credits, data)
}

type Queue struct {
	name      string
	creditsCh chan int
	src       chan *amqp.Message
	dest      chan *amqp.Message
	pumpFn    sync.Once
	events    *Events
}

func NewQueue(name string, events *Events) *Queue {
	return &Queue{
		creditsCh: make(chan int, 1000),
		dest:      make(chan *amqp.Message, 1000),
		name:      name,
		src:       make(chan *amqp.Message, 1000),
		events:    events,
	}
}

func (q *Queue) Send(ctx context.Context, msg *amqp.Message, evt LinkEvent, status *Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-status.Done():
		return status.Err()
	default:
		log.Printf("[%s] send...", q.name)

		q.src <- msg
		q.events.Send(SendEvent{
			LinkEvent: evt,
			Message:   msg,
		})
	}

	return nil
}

func (q *Queue) IssueCredit(credit uint32, evt LinkEvent, status *Status) error {
	select {
	case <-status.Done():
		return status.Err()
	default:
		log.Printf("[%s] Issuing credits %d", q.name, credit)

		q.creditsCh <- int(credit)
		q.events.IssueCredit(CreditEvent{
			Credit:    credit,
			LinkEvent: evt,
		})
		return nil
	}
}

func (q *Queue) Receive(ctx context.Context, evt LinkEvent, status *Status) (*amqp.Message, error) {
	q.pumpFn.Do(q.pumpMessages)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-status.Done():
		return nil, status.Err()
	default:
		log.Printf("[%s] Receiving on queue of length src: %d, dest: %d", q.name, len(q.src), len(q.dest))

		// only attempt to receive if we've guaranteed that we weren't closed at the start.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-status.Done():
			return nil, status.Err()
		case msg := <-q.dest:
			q.events.Receive(ReceiveEvent{
				LinkEvent: evt,
				Message:   msg,
			})
			return msg, nil
		}
	}
}

func (q *Queue) AcceptMessage(ctx context.Context, msg *amqp.Message, evt LinkEvent, status *Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-status.Done():
		return status.Err()
	default:
		q.events.Disposition(DispositionEvent{
			evt,
			DispositionTypeAccept,
			msg,
		})
		return nil
	}
}

func (q *Queue) RejectMessage(ctx context.Context, msg *amqp.Message, e *amqp.Error, evt LinkEvent, status *Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-status.Done():
		return status.Err()
	default:
		q.events.Disposition(DispositionEvent{
			evt,
			DispositionTypeReject,
			msg,
		})

		msg.Header.DeliveryCount++
		q.dest <- msg
		return nil
	}

}

func (q *Queue) ReleaseMessage(ctx context.Context, msg *amqp.Message, evt LinkEvent) error {
	q.events.Disposition(DispositionEvent{
		evt,
		DispositionTypeRelease,
		msg,
	})

	q.dest <- msg
	return nil
}

func (q *Queue) ModifyMessage(ctx context.Context, msg *amqp.Message, options *amqp.ModifyMessageOptions, evt LinkEvent) error {
	panic("ModifyMessage not implemented")
}

func (q *Queue) pumpMessages() {
	defer log.Printf("[%s] pumpMessages starting...", q.name)

	go func() {
		defer log.Printf("[%s] pumpMessages stopping...", q.name)

		for {
			credit := <-q.creditsCh

			log.Printf("[%s] pumpMessages issued %d credits...", q.name, credit)

			if credit == 0 {
				break
			}

			for i := 0; i < credit; i++ {
				log.Printf("[%s] waiting for message...", q.name)
				msg := <-q.src

				if msg == nil {
					break
				}

				log.Printf("[%s] pumped single message...", q.name)
				q.dest <- msg
			}
		}
	}()
}

func (q *Queue) Close() {
	close(q.creditsCh)
	close(q.src)
	close(q.dest)
}
