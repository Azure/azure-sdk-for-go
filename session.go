package servicebus

import (
	"sync/atomic"

	"github.com/Azure/azure-amqp-common-go/uuid"
	"pack.ag/amqp"
)

type (
	// session is a wrapper for the AMQP session with some added information to help with Service Bus messaging
	session struct {
		*amqp.Session
		SessionID string
		counter   uint32
	}
)

// newSession is a constructor for a Service Bus session which will pre-populate the SessionID with a new UUID
func newSession(amqpSession *amqp.Session) (*session, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &session{
		Session:   amqpSession,
		SessionID: id.String(),
		counter:   0,
	}, nil
}

// getNext gets and increments the next group sequence number for the session
func (s *session) getNext() uint32 {
	return atomic.AddUint32(&s.counter, 1)
}

func (s *session) String() string {
	return s.SessionID
}
