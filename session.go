package servicebus

import (
	"github.com/satori/go.uuid"
	"pack.ag/amqp"
	"sync/atomic"
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
func newSession(amqpSession *amqp.Session) *session {
	return &session{
		Session:   amqpSession,
		SessionID: uuid.NewV4().String(),
		counter:   0,
	}
}

// getNext gets and increments the next group sequence number for the session
func (s *session) getNext() uint32 {
	return atomic.AddUint32(&s.counter, 1)
}
