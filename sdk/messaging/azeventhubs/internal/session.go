// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"github.com/Azure/azure-amqp-common-go/v3/uuid"
	"github.com/Azure/go-amqp"
)

type (
	// session is a wrapper for the AMQP session with some added information to help with Service Bus messaging
	session struct {
		*amqp.Session
		SessionID string
	}
)

// newSession is a constructor for a Service Bus session which will pre-populate the SessionID with a new UUID
func newSession(amqpSession *amqp.Session) (*session, error) {
	sessionID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	return &session{
		Session:   amqpSession,
		SessionID: sessionID.String(),
	}, nil
}

func (s *session) String() string {
	return s.SessionID
}
