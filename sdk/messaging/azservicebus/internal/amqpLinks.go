// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/devigned/tab"
)

type errClosedPermanently struct{}

func (e errClosedPermanently) Error() string { return "Link has been closed permanently" }
func (e errClosedPermanently) NonRetriable() {}

func ShouldRecover(ctx context.Context, err error) bool {
	return shouldRecreateConnection(ctx, err) || shouldRecreateLink(err)
}

type AMQPLinks interface {
	EntityPath() string
	ManagementPath() string

	Audience() string

	// Get will initialize a session and call its link.linkCreator function.
	// If this link has been closed via Close() it will return an non retriable error.
	Get(ctx context.Context) (AMQPSender, AMQPReceiver, MgmtClient, uint64, error)

	// RecoverIfNeeded will check if an error requires recovery, and will recover
	// the link or, possibly, the connection.
	RecoverIfNeeded(ctx context.Context, linksRevision uint64, err error) error

	// Close will close the the link.
	// If permanent is true the link will not be auto-recreated if Get/Recover
	// are called. All functions will return `ErrLinksClosed`
	Close(ctx context.Context, permanent bool) error

	// ClosedPermanently is true if AMQPLinks.Close(ctx, true) has been called.
	ClosedPermanently() bool
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
	baseRetrier    Retrier

	mu sync.RWMutex

	// mgmt lets you interact with the $management link for your entity.
	mgmt MgmtClient

	// the AMQP session for either the 'sender' or 'receiver' link
	session AMQPSessionCloser

	// these are populated by your `createLinkFunc` when you construct
	// the amqpLinks
	sender   AMQPSenderCloser
	receiver AMQPReceiverCloser

	// last connection revision seen by this links instance.
	clientRevision uint64

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
func newAMQPLinks(ns NamespaceForAMQPLinks, entityPath string, baseRetrier Retrier, createLink CreateLinkFunc) AMQPLinks {
	l := &amqpLinks{
		entityPath:        entityPath,
		managementPath:    fmt.Sprintf("%s/$management", entityPath),
		audience:          ns.GetEntityAudience(entityPath),
		createLink:        createLink,
		baseRetrier:       baseRetrier,
		closedPermanently: false,
		revision:          1,
		ns:                ns,
	}

	return l
}

// ManagementPath is the management path for the associated entity.
func (links *amqpLinks) ManagementPath() string {
	return links.managementPath
}

// recoverLink will recycle all associated links (mgmt, receiver, sender and session)
// and recreate them using the link.linkCreator function.
func (links *amqpLinks) recoverLink(ctx context.Context, theirLinkRevision *uint64) error {
	ctx, span := tab.StartSpan(ctx, tracing.SpanRecoverLink)
	defer span.End()

	links.mu.RLock()
	closedPermanently := links.closedPermanently
	ourLinkRevision := links.revision
	links.mu.RUnlock()

	if closedPermanently {
		span.AddAttributes(tab.StringAttribute("outcome", "was_closed_permanently"))
		return errClosedPermanently{}
	}

	if theirLinkRevision != nil && ourLinkRevision > *theirLinkRevision {
		// we've already recovered past their failure.
		span.AddAttributes(
			tab.StringAttribute("outcome", "already_recovered"),
			tab.StringAttribute("lock", "readlock"),
			tab.StringAttribute("revisions", fmt.Sprintf("ours(%d), theirs(%d)", ourLinkRevision, *theirLinkRevision)),
		)
		return nil
	}

	links.mu.Lock()
	defer links.mu.Unlock()

	if theirLinkRevision != nil && ourLinkRevision > *theirLinkRevision {
		// we've already recovered past their failure.
		span.AddAttributes(
			tab.StringAttribute("outcome", "already_recovered"),
			tab.StringAttribute("lock", "writelock"),
			tab.StringAttribute("revisions", fmt.Sprintf("ours(%d), theirs(%d)", ourLinkRevision, *theirLinkRevision)),
		)
		return nil
	}

	if err := links.closeWithoutLocking(ctx, false); err != nil {
		span.Logger().Error(err)
	}

	err := links.initWithoutLocking(ctx)

	if err != nil {
		span.AddAttributes(tab.StringAttribute("init_error", err.Error()))
		return err
	}

	links.revision++

	span.AddAttributes(
		tab.StringAttribute("outcome", "recovered"),
		tab.StringAttribute("revision_new", fmt.Sprintf("%d", links.revision)),
	)
	return nil
}

// Recover will recover the links or the connection, depending
// on the severity of the error. This function uses the `baseRetrier`
// defined in the links struct.
func (links *amqpLinks) RecoverIfNeeded(ctx context.Context, linksRevision uint64, origErr error) error {
	ctx, span := tab.StartSpan(ctx, tracing.SpanRecover)
	defer span.End()

	var err error = origErr

	retrier := links.baseRetrier.Copy()

	for retrier.Try(ctx) {
		span.AddAttributes(tab.StringAttribute("recover_attempt", fmt.Sprintf("%d", retrier.CurrentTry())))

		err = links.recoverImpl(ctx, retrier.CurrentTry(), linksRevision, err)

		if err == nil {
			return nil
		}
	}

	return err
}

func (links *amqpLinks) recoverImpl(ctx context.Context, try int, linksRevision uint64, origErr error) error {
	_, span := tab.StartSpan(ctx, tracing.SpanRecoverLink)
	defer span.End()

	if origErr == nil || IsCancelError(origErr) {
		return nil
	}

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	span.AddAttributes(tab.Int64Attribute("attempt", int64(try)))

	if shouldRecreateLink(origErr) {
		span.AddAttributes(
			tab.StringAttribute("recovery_kind", "link"),
			tab.StringAttribute("error", origErr.Error()),
			tab.StringAttribute("error_type", fmt.Sprintf("%T", origErr)))

		if err := links.recoverLink(ctx, &linksRevision); err != nil {
			span.AddAttributes(tab.StringAttribute("recoveryFailure", err.Error()))
			return err
		}

		return nil
	} else if shouldRecreateConnection(ctx, origErr) {
		span.AddAttributes(
			tab.StringAttribute("recovery_kind", "connection"),
			tab.StringAttribute("error", origErr.Error()),
			tab.StringAttribute("error_type", fmt.Sprintf("%T", origErr)))

		if err := links.recoverConnection(ctx); err != nil {
			span.Logger().Error(fmt.Errorf("failed to recreate connection: %w", err))
			return err
		}

		// unconditionally recover the link if the connection died.
		if err := links.recoverLink(ctx, nil); err != nil {
			span.Logger().Error(fmt.Errorf("failed to recover links after connection restarted: %w", err))
			return err
		}

		return nil
	}

	span.AddAttributes(
		tab.StringAttribute("recovery", "none"),
		tab.StringAttribute("error", origErr.Error()),
		tab.StringAttribute("errorType", fmt.Sprintf("%T", origErr)))

	return nil
}

func (links *amqpLinks) recoverConnection(ctx context.Context) error {
	tab.For(ctx).Info("Connection is dead, recovering")

	links.mu.RLock()
	clientRevision := links.clientRevision
	links.mu.RUnlock()

	err := links.ns.Recover(ctx, clientRevision)

	if err != nil {
		tab.For(ctx).Error(fmt.Errorf("Recover connection failure: %w", err))
		return err
	}

	return nil
}

// Get will initialize a session and call its link.linkCreator function.
// If this link has been closed via Close() it will return an non retriable error.
func (l *amqpLinks) Get(ctx context.Context) (AMQPSender, AMQPReceiver, MgmtClient, uint64, error) {
	l.mu.RLock()
	sender, receiver, mgmt, revision, closedPermanently := l.sender, l.receiver, l.mgmt, l.revision, l.closedPermanently
	l.mu.RUnlock()

	if closedPermanently {
		return nil, nil, nil, 0, errClosedPermanently{}
	}

	if sender != nil || receiver != nil {
		return sender, receiver, mgmt, revision, nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.initWithoutLocking(ctx); err != nil {
		return nil, nil, nil, 0, err
	}

	return l.sender, l.receiver, l.mgmt, l.revision, nil
}

// EntityPath is the full entity path for the queue/topic/subscription.
func (l *amqpLinks) EntityPath() string {
	return l.entityPath
}

// EntityPath is the audience for the queue/topic/subscription.
func (l *amqpLinks) Audience() string {
	return l.audience
}

// ClosedPermanently is true if AMQPLinks.Close(ctx, true) has been called.
func (l *amqpLinks) ClosedPermanently() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.closedPermanently
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

	l.session, l.clientRevision, err = l.ns.NewAMQPSession(ctx)

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

	l.mgmt, err = l.ns.NewMgmtClient(ctx, l)

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
