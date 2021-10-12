// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azservicebus

import (
	"context"
	"fmt"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal"
	"github.com/Azure/go-amqp"
)

type SessionReceiver struct {
	*Receiver

	sessionID *string
}

const sessionFilterName = "com.microsoft:session-filter"

func newSessionReceiver(sessionID *string, ns internal.NamespaceWithNewAMQPLinks, entity *entity, cleanupOnClose func(), options *ReceiverOptions) (*SessionReceiver, error) {
	const code = uint64(0x00000137000000C)

	sessionReceiver := &SessionReceiver{
		sessionID: sessionID,
	}

	var err error

	sessionReceiver.Receiver, err = newReceiver(ns, entity, cleanupOnClose, options, func(ctx context.Context, session internal.AMQPSession) (internal.AMQPSenderCloser, internal.AMQPReceiverCloser, error) {
		linkOptions := createLinkOptions(sessionReceiver.Receiver.receiveMode, sessionReceiver.amqpLinks.EntityPath())

		if sessionID == nil {
			linkOptions = append(linkOptions, amqp.LinkSourceFilter(sessionFilterName, code, nil))
		} else {
			linkOptions = append(linkOptions, amqp.LinkSourceFilter(sessionFilterName, code, sessionID))
		}

		_, link, err := createReceiverLink(ctx, session, linkOptions)

		if err != nil {
			return nil, nil, err
		}

		// check the session ID that came back - if we asked for a named session ID and didn't get it then
		// we failed to get the lock.
		// if we specified nil then we can _set_ our internally held session ID now that we know the value.
		receivedSessionID := link.LinkSourceFilterValue(sessionFilterName)
		asStr, ok := receivedSessionID.(string)

		if !ok || (sessionID != nil && asStr != *sessionID) {
			return nil, nil, fmt.Errorf("invalid type/value for returned sessionID(type:%T, value:%v)", receivedSessionID, receivedSessionID)
		}

		sessionReceiver.sessionID = &asStr
		return nil, link, nil
	})

	if err != nil {
		return nil, err
	}

	return sessionReceiver, nil
}

func (sr *SessionReceiver) SessionID() string {
	// return the ultimately assigned session ID for this link (anonymous will get it from the
	// link filter options, non-anonymous is set in newSessionReceiver)
	return *sr.sessionID
}

// init ensures the link was created, guaranteeing that we get our expected session lock.
func (sr *SessionReceiver) init(ctx context.Context) error {
	// initialize the links
	_, _, _, _, err := sr.amqpLinks.Get(ctx)
	return err
}

// TODO: correct name? Also, correct return value?
func (sr *SessionReceiver) RenewSessionLock() time.Time {
	panic("RenewSessionLock not implemented")
}
