package servicebus

import (
	"context"
	"errors"
	"github.com/Azure/azure-amqp-common-go/rpc"
	"pack.ag/amqp"
	"sort"
	"time"
)

type (
	MessageIterator interface {
		Done() bool
		Next(context.Context) (*Message, error)
	}

	MessageSliceIterator struct {
		Target []*Message
		Cursor int
	}

	peekIterator struct {
		entity             *entity
		connection         *amqp.Client
		buffer             chan *Message
		lastSequenceNumber int64
	}

	PeekOption func(*peekIterator) error
)

const (
	defaultPeekPageSize = 10
)

// AsMessageSliceIterator wraps a slice of Message pointers to allow it to be made into a MessageIterator.
func AsMessageSliceIterator(target []*Message) *MessageSliceIterator {
	return &MessageSliceIterator{
		Target: target,
	}
}

func (ms MessageSliceIterator) Done() bool {
	return ms.Cursor >= len(ms.Target)
}

func (ms *MessageSliceIterator) Next(_ context.Context) (*Message, error) {
	if ms.Done() {
		return nil, ErrNoMessages{}
	}

	retval := ms.Target[ms.Cursor]
	ms.Cursor++
	return retval, nil
}

func newPeekIterator(entity *entity, connection *amqp.Client, options ...PeekOption) (*peekIterator, error) {
	retval := &peekIterator{
		entity:     entity,
		connection: connection,
	}

	foundPageSize := false
	for i := range options {
		options[i](retval)

		if retval.buffer != nil {
			foundPageSize = true
		}
	}

	if !foundPageSize {
		err := PeekWithPageSize(defaultPeekPageSize)(retval)
		if err != nil {
			return nil, err
		}
	}

	return retval, nil
}

func PeekWithPageSize(pageSize int) PeekOption {
	return func(pi *peekIterator) error {
		if pageSize < 0 {
			return errors.New("page size must not be less than zero")
		}

		if pi.buffer != nil {
			return errors.New("cannot modify an existing peekIterator's buffer")
		}

		pi.buffer = make(chan *Message, pageSize)
		return nil
	}
}

func PeekFromSequenceNumber(seq int64) PeekOption {
	return func(pi *peekIterator) error {
		pi.lastSequenceNumber = seq + 1
		return nil
	}
}

func (pi peekIterator) Done() bool {
	return false
}

func (pi *peekIterator) Next(ctx context.Context) (*Message, error) {
	if len(pi.buffer) == 0 {
		if err := pi.getNextPage(ctx); err != nil {
			return nil, err
		}
	}

	select {
	case next := <-pi.buffer:
		return next, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func (pi *peekIterator) getNextPage(ctx context.Context) error {
	const messagesField, messageField = "messages", "message"

	msg := &amqp.Message{
		ApplicationProperties: map[string]interface{}{
			operationFieldName: peekMessageOperationID,
		},
		Value: map[string]interface{}{
			"from-sequence-number": pi.lastSequenceNumber,
			"message-count":        int32(cap(pi.buffer)),
		},
	}

	if deadline, ok := ctx.Deadline(); ok {
		msg.ApplicationProperties["server-timeout"] = uint(time.Until(deadline) / time.Millisecond)
	}

	link, err := rpc.NewLink(pi.connection, pi.entity.ManagementPath())
	if err != nil {
		return err
	}

	rsp, err := link.RetryableRPC(ctx, 5, 5*time.Second, msg)
	if err != nil {
		return err
	}

	if rsp.Code == 204 {
		return ErrNoMessages{}
	}

	// Peeked messages come back in a relatively convoluted manner:
	// a map (always with one key: "messages")
	// 	of arrays
	// 		of maps (always with one key: "message")
	// 			of an array with raw encoded Service Bus messages
	if val, ok := rsp.Message.Value.(map[string]interface{}); ok {
		if rawMessages, ok := val[messagesField]; ok {
			if messages, ok := rawMessages.([]interface{}); ok {
				transformedMessages := make([]*Message, len(messages))

				for i := range messages {
					if rawEntry, ok := messages[i].(map[string]interface{}); ok {
						if rawMessage, ok := rawEntry[messageField]; ok {
							if marshaled, ok := rawMessage.([]byte); ok {
								var rehydrated amqp.Message
								err = rehydrated.UnmarshalBinary(marshaled)
								if err != nil {
									return err
								}

								transformedMessages[i], err = messageFromAMQPMessage(&rehydrated)
								if err != nil {
									return err
								}
								continue
							}
						}
						return ErrMissingField(messageField)
					}
					return newErrIncorrectType(messageField, map[string]interface{}{}, messages[i])
				}

				// This sort is done to ensure that folks wanting to peek messages in sequence order may do so.
				sort.Slice(transformedMessages, func(i, j int) bool {
					iSeq := *transformedMessages[i].SystemProperties.SequenceNumber
					jSeq := *transformedMessages[j].SystemProperties.SequenceNumber
					return iSeq < jSeq
				})

				for i := range transformedMessages {
					select {
					case pi.buffer <- transformedMessages[i]:
						// Intentionally Left Blank
					case <-ctx.Done():
						return ctx.Err()
					}
				}

				// Update last seen sequence number so that the next read starts from where this ended.
				pi.lastSequenceNumber = *transformedMessages[len(transformedMessages)-1].SystemProperties.SequenceNumber + 1
				return nil
			}
			return newErrIncorrectType(messagesField, []interface{}{}, rawMessages)
		}
		return ErrMissingField(messagesField)
	}
	return newErrIncorrectType(messageField, map[string]interface{}{}, rsp.Message.Value)
}
