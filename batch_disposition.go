package servicebus

import (
	"context"
	"fmt"

	"github.com/Azure/azure-amqp-common-go/uuid"
)

type (
	// MessageStatus defines an acceptable Message disposition status.
	MessageStatus dispositionStatus
	// BatchDispositionIterator provides an iterator over LockTokenIDs
	BatchDispositionIterator struct {
		LockTokenIDs []*uuid.UUID
		Status       MessageStatus
		cursor       int
	}
	// BatchDispositionError is an error which returns a collection of DispositionError.
	BatchDispositionError struct {
		Errors []DispositionError
	}
	// DispositionError is an error associated with a LockTokenID.
	DispositionError struct {
		LockTokenID *uuid.UUID
		err         error
	}
)

const (
	// Complete exposes completedDisposition
	Complete MessageStatus = MessageStatus(completedDisposition)
	// Abort exposes abandonedDisposition
	Abort MessageStatus = MessageStatus(abandonedDisposition)
)

func (bde BatchDispositionError) Error() string {
	msg := ""
	if len(bde.Errors) != 0 {
		msg = fmt.Sprintf("Operation failed, %d error(s) reported.", len(bde.Errors))
	}
	return msg
}

func (de DispositionError) Error() string {
	return de.err.Error()
}

// UnWrap will return the private error.
func (de DispositionError) UnWrap() error {
	return de.err
}

// Done communicates whether there are more messages remaining to be iterated over.
func (bdi *BatchDispositionIterator) Done() bool {
	return len(bdi.LockTokenIDs) == bdi.cursor
}

// Next iterates to the next LockToken
func (bdi *BatchDispositionIterator) Next() (uuid *uuid.UUID) {
	if done := bdi.Done(); done == false {
		uuid = bdi.LockTokenIDs[bdi.cursor]
		bdi.cursor++
	}
	return uuid
}

func (bdi *BatchDispositionIterator) doUpdate(ctx context.Context, ec entityConnector) BatchDispositionError {
	batchError := BatchDispositionError{}
	for !bdi.Done() {
		if uuid := bdi.Next(); uuid != nil {
			m := &Message{
				LockToken: uuid,
			}
			m.ec = ec
			err := m.sendDisposition(ctx, bdi.Status)
			if err != nil {
				batchError.Errors = append(batchError.Errors, DispositionError{
					LockTokenID: uuid,
					err:         err,
				})
			}
		}
	}
	return batchError
}

// SendBatchDisposition updates the LockTokenIDs to the disposition status.
func (q *Queue) SendBatchDisposition(ctx context.Context, iterator BatchDispositionIterator) error {
	span, ctx := q.startSpanFromContext(ctx, "sb.Queue.SendBatchDisposition")
	defer span.Finish()
	return iterator.doUpdate(ctx, q)
}

// SendBatchDisposition updates the LockTokenIDs to the desired disposition status.
func (s *Subscription) SendBatchDisposition(ctx context.Context, iterator BatchDispositionIterator) error {
	span, ctx := s.startSpanFromContext(ctx, "sb.Subscription.SendBatchDisposition")
	defer span.Finish()
	return iterator.doUpdate(ctx, s)
}

func (m *Message) sendDisposition(ctx context.Context, dispositionStatus MessageStatus) (err error) {
	switch dispositionStatus {
	case Complete:
		err = m.Complete(ctx)
	case Abort:
		err = m.Abandon(ctx)
	default:
		err = fmt.Errorf("unsupported bulk disposition status %q", dispositionStatus)
	}
	return err
}
