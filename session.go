package servicebus

import (
	"github.com/satori/go.uuid"
	"pack.ag/amqp"
	"sync/atomic"
)

// Session is a wrapper for the AMQP session with some added information to help with Service Bus messaging
type Session struct {
	*amqp.Session
	SessionID string
	counter   uint32
}

// NewSession is a constructor for a Service Bus Session which will pre-populate the SessionID with a new UUID
func NewSession(amqpSession *amqp.Session) *Session {
	return &Session{
		Session:   amqpSession,
		SessionID: uuid.NewV4().String(),
		counter:   0,
	}
}

// GetNext gets and increments the next group sequence number for the session
func (s *Session) GetNext() uint32 {
	return atomic.AddUint32(&s.counter, 1)
}
