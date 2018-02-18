package servicebus

import (
	"sync/atomic"

	"github.com/satori/go.uuid"
	"pack.ag/amqp"
)

type (
	// session is a wrapper for the AMQP session with some added information to help with Service Bus messaging
	session struct {
		*amqp.Session
		SessionID uuid.UUID
		counter   uint32
	}
)

// newSession is a constructor for a Service Bus session which will pre-populate the SessionID with a new UUID
func newSession(amqpSession *amqp.Session) *session {
	return &session{
		Session:   amqpSession,
		SessionID: uuid.NewV4(),
		counter:   0,
	}
}

// getNext gets and increments the next group sequence number for the session
func (s *session) getNext() uint32 {
	return atomic.AddUint32(&s.counter, 1)
}

func (s *session) String() string {
	return s.SessionID.String()
}
