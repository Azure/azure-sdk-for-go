package servicebus

import (
	"github.com/satori/go.uuid"
	"pack.ag/amqp"
)

// Session is a wrapper for the AMQP session with some added information to help with Service Bus messaging
type Session struct {
	*amqp.Session
	SessionID string
}

// NewSession is a constructor for a Service Bus Session which will pre-populate the SessionID with a new UUID
func NewSession(amqpSession *amqp.Session) *Session {
	return &Session{
		Session:   amqpSession,
		SessionID: uuid.NewV4().String(),
	}
}
