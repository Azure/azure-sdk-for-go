// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/devigned/tab"
)

type LinksWithID struct {
	Sender   AMQPSender
	Receiver AMQPReceiver
	RPC      RPCLink
	ID       LinkID
}

type RetryWithLinksFn func(ctx context.Context, lwid *LinksWithID, args *utils.RetryFnArgs) error

type AMQPLinks interface {
	EntityPath() string
	ManagementPath() string

	Audience() string

	// Get will initialize a session and call its link.linkCreator function.
	// If this link has been closed via Close() it will return an non retriable error.
	Get(ctx context.Context) (*LinksWithID, error)

	// Retry will run your callback, recovering links when necessary.
	Retry(ctx context.Context, name string, fn RetryWithLinksFn, o utils.RetryOptions) error

	// RecoverIfNeeded will check if an error requires recovery, and will recover
	// the link or, possibly, the connection.
	RecoverIfNeeded(ctx context.Context, linkID LinkID, err error) error

	// Close will close the the link.
	// If permanent is true the link will not be auto-recreated if Get/Recover
	// are called. All functions will return `ErrLinksClosed`
	Close(ctx context.Context, permanent bool) error

	// ClosedPermanently is true if AMQPLinks.Close(ctx, true) has been called.
	ClosedPermanently() bool
}

// AMQPLinksImpl manages the set of AMQP links (and detritus) typically needed to work
//  within Service Bus:
// - An *goamqp.Sender or *goamqp.Receiver AMQP link (could also be 'both' if needed)
// - A `$management` link
// - an *goamqp.Session
//
// State management can be done through Recover (close and reopen), Close (close permanently, return failures)
// and Get() (retrieve latest version of all AMQPLinksImpl, or create if needed).
type AMQPLinksImpl struct {
	// NOTE: values need to be 64-bit aligned. Simplest way to make sure this happens
	// is just to make it the first value in the struct
	// See:
	//   Godoc: https://pkg.go.dev/sync/atomic#pkg-note-BUG
	//   PR: https://github.com/Azure/azure-sdk-for-go/pull/16847
	id LinkID

	entityPath     string
	managementPath string
	audience       string
	createLink     CreateLinkFunc

	mu sync.RWMutex

	// RPCLink lets you interact with the $management link for your entity.
	RPCLink RPCLink

	// the AMQP session for either the 'sender' or 'receiver' link
	session AMQPSessionCloser

	// these are populated by your `createLinkFunc` when you construct
	// the amqpLinks
	Sender   AMQPSenderCloser
	Receiver AMQPReceiverCloser

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
func NewAMQPLinks(ns NamespaceForAMQPLinks, entityPath string, createLink CreateLinkFunc) AMQPLinks {
	l := &AMQPLinksImpl{
		entityPath:        entityPath,
		managementPath:    fmt.Sprintf("%s/$management", entityPath),
		audience:          ns.GetEntityAudience(entityPath),
		createLink:        createLink,
		closedPermanently: false,
		ns:                ns,
	}

	return l
}

// ManagementPath is the management path for the associated entity.
func (links *AMQPLinksImpl) ManagementPath() string {
	return links.managementPath
}

// recoverLink will recycle all associated links (mgmt, receiver, sender and session)
// and recreate them using the link.linkCreator function.
func (links *AMQPLinksImpl) recoverLink(ctx context.Context, theirLinkRevision LinkID) error {
	ctx, span := tab.StartSpan(ctx, tracing.SpanRecoverLink)
	defer span.End()

	links.mu.RLock()
	closedPermanently := links.closedPermanently
	ourLinkRevision := links.id
	links.mu.RUnlock()

	if closedPermanently {
		span.AddAttributes(tab.StringAttribute("outcome", "was_closed_permanently"))
		return ErrNonRetriable{Message: "Link has been closed permanently"}
	}

	// cheap check before we do the lock
	if ourLinkRevision.Link != theirLinkRevision.Link {
		// we've already recovered past their failure.
		return nil
	}

	links.mu.Lock()
	defer links.mu.Unlock()

	// check once more, just in case someone else modified it before we took
	// the lock.
	if links.id.Link != theirLinkRevision.Link {
		// we've already recovered past their failure.
		return nil
	}

	err := links.initWithoutLocking(ctx)

	if err != nil {
		span.AddAttributes(tab.StringAttribute("init_error", err.Error()))
		return err
	}

	span.AddAttributes(
		tab.StringAttribute("outcome", "recovered"),
		tab.StringAttribute("revision_new", fmt.Sprintf("%d", links.id)),
	)
	return nil
}

// Recover will recover the links or the connection, depending
// on the severity of the error.
func (links *AMQPLinksImpl) RecoverIfNeeded(ctx context.Context, theirID LinkID, origErr error) error {
	ctx, span := tab.StartSpan(ctx, tracing.SpanRecover)
	defer span.End()

	if origErr == nil || IsCancelError(origErr) {
		return nil
	}

	log.Writef(EventConn, "Recovering link for error %s", origErr.Error())

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	sbe := GetSBErrInfo(origErr)

	if sbe.RecoveryKind == RecoveryKindLink {
		if err := links.recoverLink(ctx, theirID); err != nil {
			log.Writef(EventConn, "failed to recreate link: %s", err.Error())
			return err
		}

		log.Writef(EventConn, "Recovered links")
		return nil
	} else if sbe.RecoveryKind == RecoveryKindConn {
		if err := links.recoverConnection(ctx, theirID); err != nil {
			log.Writef(EventConn, "failed to recreate connection: %s", err.Error())
			return err
		}

		log.Writef(EventConn, "Recovered connection and links")
		return nil
	}

	log.Writef(EventConn, "Recovered, no action needed")
	return nil
}

func (links *AMQPLinksImpl) recoverConnection(ctx context.Context, theirID LinkID) error {
	tab.For(ctx).Info("Connection is dead, recovering")

	links.mu.Lock()
	defer links.mu.Unlock()

	created, err := links.ns.Recover(ctx, uint64(theirID.Conn))

	if err != nil {
		log.Writef(EventConn, "Recover connection failure: %s", err)
		return err
	}

	// We'll recreate the link if:
	// - `created` is true, meaning we recreated the AMQP connection (ie, all old links are invalid)
	// - the link they received an error on is our current link, so it needs to be recreated.
	//   (if it wasn't the same then we've already recovered and created a new link,
	//    so no recovery would be needed)
	if created || theirID.Link == links.id.Link {
		log.Writef(EventConn, "recreating link: c: %v, current:%v, old:%v", created, links.id, theirID)
		if err := links.initWithoutLocking(ctx); err != nil {
			return err
		}
	}

	return nil
}

// LinkID is ID that represent our current link and the client used to create it.
// These are used when trying to determine what parts need to be recreated when
// an error occurs, to prevent recovering a connection/link repeatedly.
// See amqpLinks.RecoverIfNeeded() for usage.
type LinkID struct {
	// Conn is the ID of the connection we used to create our links.
	Conn uint64

	// Link is the ID of our current link.
	Link uint64
}

// Get will initialize a session and call its link.linkCreator function.
// If this link has been closed via Close() it will return an non retriable error.
func (l *AMQPLinksImpl) Get(ctx context.Context) (*LinksWithID, error) {
	l.mu.RLock()
	sender, receiver, mgmtLink, revision, closedPermanently := l.Sender, l.Receiver, l.RPCLink, l.id, l.closedPermanently
	l.mu.RUnlock()

	if closedPermanently {
		return nil, ErrNonRetriable{}
	}

	if sender != nil || receiver != nil {
		return &LinksWithID{
			Sender:   sender,
			Receiver: receiver,
			RPC:      mgmtLink,
			ID:       revision,
		}, nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	if err := l.initWithoutLocking(ctx); err != nil {
		return nil, err
	}

	return &LinksWithID{
		Sender:   l.Sender,
		Receiver: l.Receiver,
		RPC:      l.RPCLink,
		ID:       l.id,
	}, nil
}

func (l *AMQPLinksImpl) Retry(ctx context.Context, name string, fn RetryWithLinksFn, o utils.RetryOptions) error {
	var lastID LinkID

	return utils.Retry(ctx, name, func(ctx context.Context, args *utils.RetryFnArgs) error {
		if err := l.RecoverIfNeeded(ctx, lastID, args.LastErr); err != nil {
			return err
		}

		linksWithVersion, err := l.Get(ctx)

		if err != nil {
			return err
		}

		lastID = linksWithVersion.ID
		return fn(ctx, linksWithVersion, args)
	}, IsFatalSBError, o)
}

// EntityPath is the full entity path for the queue/topic/subscription.
func (l *AMQPLinksImpl) EntityPath() string {
	return l.entityPath
}

// EntityPath is the audience for the queue/topic/subscription.
func (l *AMQPLinksImpl) Audience() string {
	return l.audience
}

// ClosedPermanently is true if AMQPLinks.Close(ctx, true) has been called.
func (l *AMQPLinksImpl) ClosedPermanently() bool {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.closedPermanently
}

// Close will close the the link permanently.
// Any further calls to Get()/Recover() to return ErrLinksClosed.
func (l *AMQPLinksImpl) Close(ctx context.Context, permanent bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.closeWithoutLocking(ctx, permanent)
}

// initWithoutLocking will create a new link, unconditionally.
func (l *AMQPLinksImpl) initWithoutLocking(ctx context.Context) error {
	// shut down any links we have
	_ = l.closeWithoutLocking(ctx, false)

	var err error
	l.cancelAuthRefreshLink, err = l.ns.NegotiateClaim(ctx, l.entityPath)

	if err != nil {
		if err := l.closeWithoutLocking(ctx, false); err != nil {
			log.Writef(EventConn, "Failure during link cleanup after negotiateClaim: %s", err.Error())
		}
		return err
	}

	l.cancelAuthRefreshMgmtLink, err = l.ns.NegotiateClaim(ctx, l.managementPath)

	if err != nil {
		if err := l.closeWithoutLocking(ctx, false); err != nil {
			log.Writef(EventConn, "Failure during link cleanup after negotiate claim for mgmt link: %s", err.Error())
		}
		return err
	}

	sess, cr, err := l.ns.NewAMQPSession(ctx)

	if err != nil {
		if err := l.closeWithoutLocking(ctx, false); err != nil {
			log.Writef(EventConn, "Failure during link cleanup after creating AMQP session: %s", err.Error())
		}
		return err
	}

	l.session = sess
	l.id.Conn = cr

	l.Sender, l.Receiver, err = l.createLink(ctx, l.session)

	if err != nil {
		if err := l.closeWithoutLocking(ctx, false); err != nil {
			log.Writef(EventConn, "Failure during link cleanup after creating link: %s", err.Error())
		}
		return err
	}

	rpcLink, err := l.ns.NewRPCLink(ctx, l.ManagementPath())

	if err != nil {
		if err := l.closeWithoutLocking(ctx, false); err != nil {
			log.Writef("Failure during link cleanup after creating mgmt client: %s", err.Error())
		}
		return err
	}

	l.RPCLink = rpcLink
	l.id.Link++
	return nil
}

// close closes the link.
// NOTE: No locking is done in this function, call `Close` if you require locking.
func (l *AMQPLinksImpl) closeWithoutLocking(ctx context.Context, permanent bool) error {
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

	if l.Sender != nil {
		if err := l.Sender.Close(ctx); err != nil {
			messages = append(messages, fmt.Sprintf("amqp sender close error: %s", err.Error()))
		}
		l.Sender = nil
	}

	if l.Receiver != nil {
		if err := l.Receiver.Close(ctx); err != nil {
			messages = append(messages, fmt.Sprintf("amqp receiver close error: %s", err.Error()))
		}
		l.Receiver = nil
	}

	if l.session != nil {
		if err := l.session.Close(ctx); err != nil {
			messages = append(messages, fmt.Sprintf("amqp session close error: %s", err.Error()))
		}
		l.session = nil
	}

	if l.RPCLink != nil {
		if err := l.RPCLink.Close(ctx); err != nil {
			messages = append(messages, fmt.Sprintf("$management link close error: %s", err.Error()))
		}
		l.RPCLink = nil
	}

	if len(messages) > 0 {
		return errors.New(strings.Join(messages, "\n"))
	}

	return nil
}
