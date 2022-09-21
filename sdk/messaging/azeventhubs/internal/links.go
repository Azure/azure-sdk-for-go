// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.
package internal

import (
	"context"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/exported"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/internal/utils"
)

type AMQPLink interface {
	Close(ctx context.Context) error
	LinkName() string
}

type LinkWithID[LinkT AMQPLink] struct {
	// ConnID is an arbitrary (but unique) integer that represents the
	// current connection. This comes back from the Namespace, anytime
	// it hands back a connection.
	ConnID uint64

	// Link will be an amqp.Receiver or amqp.Sender link.
	Link LinkT
}

// LinksForPartitionClient are the functions that the PartitionClient uses within Links[T]
// (for unit testing only)
type LinksForPartitionClient[LinkT AMQPLink] interface {
	RecoverIfNeeded(ctx context.Context, partitionID string, lwid *LinkWithID[LinkT], err error) error
	Retry(ctx context.Context, eventName log.Event, operation string, partitionID string, retryOptions exported.RetryOptions, fn func(ctx context.Context, lwid LinkWithID[LinkT]) error) error
	Close(ctx context.Context) error
}

type Links[LinkT AMQPLink] struct {
	ns NamespaceForAMQPLinks

	linksMu *sync.RWMutex
	links   map[string]*linkState[LinkT]

	managementLinkMu *sync.RWMutex
	managementLink   *linkState[RPCLink]

	managementPath string
	newLinkFn      func(ctx context.Context, session amqpwrap.AMQPSession, partitionID string) (LinkT, error)
	entityPathFn   func(partitionID string) string
}

type NewLinksFn[LinkT AMQPLink] func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string) (LinkT, error)

func NewLinks[LinkT AMQPLink](ns NamespaceForAMQPLinks, managementPath string, entityPathFn func(partitionID string) string, newLinkFn NewLinksFn[LinkT]) *Links[LinkT] {
	return &Links[LinkT]{
		ns:               ns,
		linksMu:          &sync.RWMutex{},
		links:            map[string]*linkState[LinkT]{},
		managementLinkMu: &sync.RWMutex{},
		newLinkFn:        newLinkFn,
		entityPathFn:     entityPathFn,
		managementPath:   managementPath,
	}
}

func (l *Links[LinkT]) RecoverIfNeeded(ctx context.Context, partitionID string, lwid *LinkWithID[LinkT], err error) error {
	if lwid == nil {
		return nil
	}

	rk := GetRecoveryKind(err)

	switch rk {
	case RecoveryKindNone:
		return nil
	case RecoveryKindLink:
		// close and recreate
		return l.closePartitionLinkIfMatch(ctx, partitionID, lwid.Link.LinkName())
	case RecoveryKindConn:
		created, err := l.ns.Recover(ctx, lwid.ConnID)

		if err != nil {
			return err
		}

		if !created {
			return nil
		}

		if err := l.closeLinks(ctx, false); err != nil {
			return err
		}

		return nil
	default:
		return err
	}
}

func (l *Links[LinkT]) Retry(ctx context.Context, eventName log.Event, operation string, partitionID string, retryOptions exported.RetryOptions, fn func(ctx context.Context, lwid LinkWithID[LinkT]) error) error {
	var prevLinkWithID *LinkWithID[LinkT]

	didQuickRetry := false

	isFatalErrorFunc := func(err error) bool {
		return GetRecoveryKind(err) == RecoveryKindFatal
	}

	return utils.Retry(ctx, eventName, operation, retryOptions, func(ctx context.Context, args *utils.RetryFnArgs) error {
		if err := l.RecoverIfNeeded(ctx, partitionID, prevLinkWithID, args.LastErr); err != nil {
			return err
		}

		linkWithID, err := l.GetLink(ctx, partitionID)

		if err != nil {
			return err
		}

		prevLinkWithID = linkWithID

		if err := fn(ctx, *linkWithID); err != nil {
			if args.I == 0 && !didQuickRetry && IsQuickRecoveryError(err) {
				// go-amqp will asynchronously handle detaches. This means errors that you get
				// back from Send(), for instance, can actually be from much earlier in time
				// depending on the last time you called into Send().
				//
				// This means we'll sometimes do an unneeded sleep after a failed retry when
				// it would have just immediately worked. To counteract that we'll do a one-time
				// quick attempt to recreate link immediately if we see a detach error. This might
				// waste a bit of time attempting to do the creation, but since it's just link creation
				// it should be fairly fast.
				//
				// So when we've received a detach is:
				//   0th attempt
				//   extra immediate 0th attempt (if last error was detach)
				//   (actual retries)
				//
				// Whereas normally you'd do (for non-detach errors):
				//   0th attempt
				//   (actual retries)
				azlog.Writef(exported.EventConn, "(%s) Link was previously detached. Attempting quick reconnect to recover from error: %s", operation, err.Error())
				didQuickRetry = true
				args.ResetAttempts()
			}

			return err
		}

		return nil
	}, isFatalErrorFunc)
}

func (l *Links[LinkT]) CloseLink(ctx context.Context, partitionID string) error {
	l.linksMu.RLock()
	current := l.links[partitionID]
	l.linksMu.RUnlock()

	if current == nil {
		return nil
	}

	l.linksMu.Lock()
	defer l.linksMu.Unlock()

	current = l.links[partitionID]

	if current == nil {
		return nil
	}

	_ = current.Close(ctx)
	delete(l.links, partitionID)

	return nil
}

func (l *Links[LinkT]) GetLink(ctx context.Context, partitionID string) (*LinkWithID[LinkT], error) {
	if err := l.checkOpen(); err != nil {
		return nil, err
	}

	l.linksMu.RLock()
	current := l.links[partitionID]
	l.linksMu.RUnlock()

	if current != nil {
		return &LinkWithID[LinkT]{
			ConnID: l.links[partitionID].ConnID,
			Link:   *l.links[partitionID].Link,
		}, nil
	}

	// no existing link, let's create a new one within the write lock.
	l.linksMu.Lock()
	defer l.linksMu.Unlock()

	// check again now that we have the write lock
	current = l.links[partitionID]

	if current == nil {
		ls, err := l.newLinkState(ctx, partitionID)

		if err != nil {
			return nil, err
		}

		l.links[partitionID] = ls
	}

	return &LinkWithID[LinkT]{
		ConnID: l.links[partitionID].ConnID,
		Link:   *l.links[partitionID].Link,
	}, nil
}

func (l *Links[LinkT]) GetManagementLink(ctx context.Context) (LinkWithID[RPCLink], error) {
	if err := l.checkOpen(); err != nil {
		return LinkWithID[RPCLink]{}, err
	}

	l.managementLinkMu.Lock()
	defer l.managementLinkMu.Unlock()

	if l.managementLink == nil {
		ls, err := l.newManagementLinkState(ctx)

		if err != nil {
			return LinkWithID[RPCLink]{}, err
		}

		l.managementLink = ls
	}

	return LinkWithID[RPCLink]{
		ConnID: l.managementLink.ConnID,
		Link:   *l.managementLink.Link,
	}, nil
}

func (l *Links[LinkT]) newLinkState(ctx context.Context, partitionID string) (*linkState[LinkT], error) {
	// check again now that we have the write lock
	ls := &linkState[LinkT]{}

	cancelAuth, _, err := l.ns.NegotiateClaim(ctx, l.entityPathFn(partitionID))

	if err != nil {
		return nil, err
	}

	ls.cancelAuth = cancelAuth

	session, connID, err := l.ns.NewAMQPSession(ctx)

	if err != nil {
		_ = ls.Close(ctx)
		return nil, err
	}

	ls.session = session
	ls.ConnID = connID

	tmpLink, err := l.newLinkFn(ctx, session, l.entityPathFn(partitionID))

	if err != nil {
		_ = ls.Close(ctx)
		return nil, err
	}

	ls.Link = &tmpLink
	return ls, nil
}

func (l *Links[LinkT]) newManagementLinkState(ctx context.Context) (*linkState[RPCLink], error) {
	ls := &linkState[RPCLink]{}

	cancelAuth, _, err := l.ns.NegotiateClaim(ctx, l.managementPath)

	if err != nil {
		return nil, err
	}

	ls.cancelAuth = cancelAuth

	tmpRPCLink, connID, err := l.ns.NewRPCLink(ctx, "$management")

	if err != nil {
		_ = ls.Close(ctx)
		return nil, err
	}

	ls.ConnID = connID
	ls.Link = &tmpRPCLink

	return ls, nil
}

func (l *Links[LinkT]) Close(ctx context.Context) error {
	return l.closeLinks(ctx, true)
}

func (l *Links[LinkT]) closeLinks(ctx context.Context, permanent bool) error {
	if err := l.closeManagementLink(ctx); err != nil {
		azlog.Writef(exported.EventConn, "Error while cleaning up management link while doing connection recovery: %s", err.Error())
	}

	l.linksMu.Lock()
	defer l.linksMu.Unlock()

	tmpLinks := l.links
	l.links = nil

	for partitionID, link := range tmpLinks {
		if err := link.Close(ctx); err != nil {
			azlog.Writef(exported.EventConn, "Error while cleaning up link for partition ID '%s' while doing connection recovery: %s", partitionID, err.Error())
		}
	}

	if !permanent {
		l.links = map[string]*linkState[LinkT]{}
	}
	return nil
}

func (l *Links[LinkT]) checkOpen() error {
	l.linksMu.RLock()
	defer l.linksMu.RUnlock()

	if l.links == nil {
		return NewErrNonRetriable("client has been closed by user")
	}

	return nil
}

func (l *Links[LinkT]) closePartitionLinkIfMatch(ctx context.Context, partitionID string, linkName string) error {
	l.linksMu.RLock()
	current, exists := l.links[partitionID]
	l.linksMu.RUnlock()

	if !exists ||
		(*current.Link).LinkName() != linkName { // we've already created a new link, their link was stale.
		return nil
	}

	l.linksMu.Lock()
	defer l.linksMu.Unlock()

	current, exists = l.links[partitionID]

	if !exists ||
		(*current.Link).LinkName() != linkName { // we've already created a new link, their link was stale.
		return nil
	}

	current.cancelAuth()
	delete(l.links, partitionID)

	return current.Close(ctx)
}

func (l *Links[LinkT]) closeManagementLink(ctx context.Context) error {
	l.managementLinkMu.Lock()
	defer l.managementLinkMu.Unlock()

	if l.managementLink != nil {
		err := l.managementLink.Close(ctx)
		l.managementLink = nil
		return err
	}

	return nil
}

type linkState[LinkT AMQPLink] struct {
	// ConnID is an arbitrary (but unique) integer that represents the
	// current connection. This comes back from the Namespace, anytime
	// it hands back a connection.
	ConnID uint64

	// Link will be an amqp.Receiver, an amqp.Sender link, or an RPCLink.
	Link *LinkT

	// cancelAuth cancels the backround claim negotation for this link.
	cancelAuth func()

	// optional session, if we created one for this
	// link.
	session amqpwrap.AMQPSession
}

func (ls *linkState[LinkT]) Close(ctx context.Context) error {
	if ls.cancelAuth != nil {
		ls.cancelAuth()
	}

	if ls.Link != nil {
		return (*ls.Link).Close(ctx)
	}

	return nil
}
