// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package emulation

import (
	"context"
	"fmt"
	"sync"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/go-amqp"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/test"
)

type Operation struct {
	Op       string
	Entity   string
	LinkName string
	Role     string

	Credits uint32
	M       *amqp.Message
}

const (
	EventEmulator azlog.Event = "emulator"
)

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
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewQueue(name string, events *Events) *Queue {
	ctx, cancel := context.WithCancel(context.Background())

	return &Queue{
		creditsCh: make(chan int, 1000),
		dest:      make(chan *amqp.Message, 1000),
		name:      name,
		src:       make(chan *amqp.Message, 1000),
		events:    events,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (q *Queue) Send(ctx context.Context, tempMsg *amqp.Message, evt LinkEvent, status *Status) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-status.Done():
		return status.Err()
	default:
		azlog.Writef(EventEmulator, "[%s] send...", q.name)
		msg := *tempMsg

		// this is a bit silly because it's not 1:1 but we just
		// need _some_ bytes that are unique for the lock token.
		id := test.MustAMQPUUID()
		msg.DeliveryTag = id[:16]

		q.src <- &msg
		q.events.Send(SendEvent{
			LinkEvent: evt,
			Message:   &msg,
		})
	}

	return nil
}

func (q *Queue) IssueCredit(credit uint32, evt LinkEvent, status *Status) error {
	azlog.Writef(EventEmulator, "[%s] Issuing credits %d", q.name, credit)

	q.creditsCh <- int(credit)
	q.events.IssueCredit(CreditEvent{
		Credit:    credit,
		LinkEvent: evt,
	})

	// TODO: the only time we ever get an error is if you chose manual
	// credits and attempted to add credits.
	return nil
}

func (q *Queue) Receive(ctx context.Context, evt LinkEvent, status *Status) (*amqp.Message, error) {
	q.pumpFn.Do(q.pumpMessages)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-status.Done():
		return nil, status.Err()
	default:
		azlog.Writef(EventEmulator, "[%s] Receiving on queue of length src: %d, dest: %d", q.name, len(q.src), len(q.dest))

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
	defer azlog.Writef(EventEmulator, "[%s] pumpMessages starting...", q.name)

	go func() {
		defer azlog.Writef(EventEmulator, "[%s] pumpMessages stopping...", q.name)

		defer func() {
			close(q.creditsCh)
			close(q.src)
			close(q.dest)
		}()

		for {
			select {
			case credit := <-q.creditsCh:
				azlog.Writef(EventEmulator, "[%s] pumpMessages issued %d credits...", q.name, credit)

				if credit == 0 {
					break
				}

				for i := 0; i < credit; i++ {
					azlog.Writef(EventEmulator, "[%s] waiting for message...", q.name)
					msg := <-q.src

					if msg == nil {
						break
					}

					azlog.Writef(EventEmulator, "[%s] pumped single message...", q.name)

					select {
					case <-q.ctx.Done():
						return
					default:
						q.dest <- msg
					}

				}
			case <-q.ctx.Done():
				return
			}
		}
	}()
}

func (q *Queue) Close() {
	q.cancel()
}
