// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/devigned/tab"
)

type errClosedPermanently struct {
}

func (e errClosedPermanently) Error() string { return "Link has been closed permanently" }
func (e errClosedPermanently) NonRetriable() {}

type AMQPLinks interface {
	EntityPath() string
	Audience() string

	// Get will initialize a session and call its link.linkCreator function.
	// If this link has been closed via Close() it will return an non retriable error.
	Get(ctx context.Context) (AMQPSender, AMQPReceiver, MgmtClient, uint64, error)

	// Recover will recycle all associated links (mgmt, receiver, sender and session)
	Recover(ctx context.Context) error

	// Close will close the the link.
	// If permanent is true the link will not be auto-recreated if Get/Recover
	// are called. All functions will return `ErrLinksClosed`
	Close(ctx context.Context, permanent bool) error
}

// amqpLinks manages the set of AMQP links (and detritus) typically needed to work
//  within Service Bus:
// - An *goamqp.Sender or *goamqp.Receiver AMQP link (could also be 'both' if needed)
// - A `$management` link
// - an *goamqp.Session
//
// State management can be done through Recover (close and reopen), Close (close permanently, return failures)
// and Get() (retrieve latest version of all amqpLinks, or create if needed).
type amqpLinks struct {
	entityPath     string
	managementPath string
	audience       string
	createLink     CreateLinkFunc

	mu sync.RWMutex

	// mgmt lets you interact with the $management link for your entity.
	mgmt MgmtClient

	// the AMQP session for either the 'sender' or 'receiver' link
	session AMQPSessionCloser

	// these are populated by your `createLinkFunc` when you construct
	// the amqpLinks
	sender   AMQPSenderCloser
	receiver AMQPReceiverCloser

	// the current 'revision' of our set of links.
	// starts at 1, increments each time you call Recover().
	revision uint64

	// whether this links set has been closed permanently (via Close)
	// Recover() does not affect this value.
	closedPermanently bool

	cancelAuthRefreshLink     func() <-chan struct{}
	cancelAuthRefreshMgmtLink func() <-chan struct{}

	ns NamespaceForAMQPLinks
}

// CreateLinkFunc creates the links, using the given session. Typically you'll only create either an
// *amqp.Sender or a *amqp.Receiver. AMQPLinks handles it either way.
type CreateLinkFunc func(ctx context.Context, session AMQPSession) (AMQPSenderCloser, AMQPReceiverCloser, error)

// NewAMQPLinks creates a session, starts the claim refresher and creates an associated
// management link for a specific entity path.
func newAMQPLinks(ns NamespaceForAMQPLinks, entityPath string, createLink CreateLinkFunc) AMQPLinks {
	l := &amqpLinks{
		entityPath:        entityPath,
		managementPath:    fmt.Sprintf("%s/$management", entityPath),
		audience:          ns.GetEntityAudience(entityPath),
		createLink:        createLink,
		closedPermanently: false,
		revision:          1,
		ns:                ns,
	}

	return l
}

// Recover will recycle all associated links (mgmt, receiver, sender and session)
// and recreate them using the link.linkCreator function.
func (links *amqpLinks) Recover(ctx context.Context) error {
	links.mu.RLock()
	closedPermanently := links.closedPermanently
	links.mu.RUnlock()

	if closedPermanently {
		return errClosedPermanently{}
	}

	links.mu.Lock()
	defer links.mu.Unlock()

	links.revision++

	if err := links.closeWithoutLocking(ctx, false); err != nil {
		tab.For(ctx)
	}

	return links.initWithoutLocking(ctx)
}

// Get will initialize a session and call its link.linkCreator function.
// If this link has been closed via Close() it will return an non retriable error.
func (l *amqpLinks) Get(ctx context.Context) (AMQPSender, AMQPReceiver, MgmtClient, uint64, error) {
	l.mu.RLock()
	sender, receiver, mgmt, closedPermanently := l.sender, l.receiver, l.mgmt, l.closedPermanently
	l.mu.RUnlock()

	if closedPermanently {
		return nil, nil, nil, 0, errClosedPermanently{}
	}

	if sender != nil || receiver != nil {
		return sender, receiver, mgmt, 0, nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.initWithoutLocking(ctx); err != nil {
		return nil, nil, nil, 0, err
	}

	return l.sender, l.receiver, l.mgmt, 0, nil
}

// EntityPath is the full entity path for the queue/topic/subscription.
func (l *amqpLinks) EntityPath() string {
	return l.entityPath
}

// EntityPath is the audience for the queue/topic/subscription.
func (l *amqpLinks) Audience() string {
	return l.audience
}

// Close will close the the link permanently.
// Any further calls to Get()/Recover() to return ErrLinksClosed.
func (l *amqpLinks) Close(ctx context.Context, permanent bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.closeWithoutLocking(ctx, permanent)
}

// initWithoutLocking will create a new link, unconditionally.
func (l *amqpLinks) initWithoutLocking(ctx context.Context) error {
	var err error
	l.cancelAuthRefreshLink, err = l.ns.NegotiateClaim(ctx, l.entityPath)

	if err != nil {
		if err := l.closeWithoutLocking(ctx, false); err != nil {
			tab.For(ctx).Debug(fmt.Sprintf("Failure during link cleanup after negotiateClaim: %s", err.Error()))
		}
		return err
	}

	l.cancelAuthRefreshMgmtLink, err = l.ns.NegotiateClaim(ctx, l.managementPath)

	if err != nil {
		if err := l.closeWithoutLocking(ctx, false); err != nil {
			tab.For(ctx).Debug(fmt.Sprintf("Failure during link cleanup after negotiate claim for mgmt link: %s", err.Error()))
		}
		return err
	}

	l.session, err = l.ns.NewAMQPSession(ctx)

	if err != nil {
		if err := l.closeWithoutLocking(ctx, false); err != nil {
			tab.For(ctx).Debug(fmt.Sprintf("Failure during link cleanup after creating AMQP session: %s", err.Error()))
		}
		return err
	}

	l.sender, l.receiver, err = l.createLink(ctx, l.session)

	if err != nil {
		if err := l.closeWithoutLocking(ctx, false); err != nil {
			tab.For(ctx).Debug(fmt.Sprintf("Failure during link cleanup after creating link: %s", err.Error()))
		}
		return err
	}

	l.mgmt, err = l.ns.NewMgmtClient(ctx, l.managementPath)

	if err != nil {
		if err := l.closeWithoutLocking(ctx, false); err != nil {
			tab.For(ctx).Debug(fmt.Sprintf("Failure during link cleanup after creating mgmt client: %s", err.Error()))
		}
		return err
	}

	return nil
}

// close closes the link.
// NOTE: No locking is done in this function, call `Close` if you require locking.
func (l *amqpLinks) closeWithoutLocking(ctx context.Context, permanent bool) error {
	if l.closedPermanently {
		return nil
	}

	defer func() {
		if permanent {
			l.closedPermanently = true
		}
	}()

	var messages []string

	if l.cancelAuthRefreshLink != nil {
		l.cancelAuthRefreshLink()
	}

	if l.cancelAuthRefreshMgmtLink != nil {
		l.cancelAuthRefreshMgmtLink()
	}

	if l.sender != nil {
		if err := l.sender.Close(ctx); err != nil {
			messages = append(messages, fmt.Sprintf("amqp sender close error: %s", err.Error()))
		}
		l.sender = nil
	}

	if l.receiver != nil {
		if err := l.receiver.Close(ctx); err != nil {
			messages = append(messages, fmt.Sprintf("amqp receiver close error: %s", err.Error()))
		}
		l.receiver = nil
	}

	if l.session != nil {
		if err := l.session.Close(ctx); err != nil {
			messages = append(messages, fmt.Sprintf("amqp session close error: %s", err.Error()))
		}
		l.session = nil
	}

	if l.mgmt != nil {
		if err := l.mgmt.Close(ctx); err != nil {
			messages = append(messages, fmt.Sprintf("$management link close error: %s", err.Error()))
		}
		l.mgmt = nil
	}

	if len(messages) > 0 {
		return errors.New(strings.Join(messages, "\n"))
	}

	return nil
}
