// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
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

	// PartitionID, if available.
	PartitionID string
}

func (lwid *LinkWithID[LinkT]) String() string {
	if lwid == nil {
		return "none"
	}

	return fmt.Sprintf("c:%d,l:%.5s,p:%s", lwid.ConnID, lwid.Link.LinkName(), lwid.PartitionID)
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
	managementLink   *linkState[amqpwrap.RPCLink]

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
		managementPath:   managementPath,

		newLinkFn:    newLinkFn,
		entityPathFn: entityPathFn,
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
		if err := l.closePartitionLinkIfMatch(ctx, partitionID, lwid.Link.LinkName()); err != nil {
			azlog.Writef(exported.EventConn, "(%s) Error when cleaning up old link for link recovery: %s", lwid.String(), err)
			return err
		}

		return nil
	case RecoveryKindConn:
		// We only close _this_ partition's link. Other partitions will also get an error, and will recover.
		// We used to close _all_ the links, but no longer do that since it's possible (when we do receiver
		// redirect) to have more than one active connection at a time which means not all links would be
		// affected when a single connection goes down.
		if err := l.closePartitionLinkIfMatch(ctx, partitionID, lwid.Link.LinkName()); err != nil {
			azlog.Writef(exported.EventConn, "(%s) Error when cleaning up old link: %s", lwid.String(), err)

			// NOTE: this is best effort - it's probable the connection is dead anyways so we'll log
			// but ignore the error for recovery purposes.
		}

		// There are two possibilities here:
		//
		// 1. (stale) The caller got this error but the `lwid` they're passing us is 'stale' - ie, '
		//    the connection the error happened on doesn't exist anymore (we recovered already) or
		//    the link itself is no longer active in our cache.
		//
		// 2. (current) The caller got this error and is the current link and/or connection, so we're going to
		//    need to recycle the connection (possibly) and links.
		//
		// For #1, we basically don't need to do anything. Recover(old-connection-id) will be a no-op
		// and the closePartitionLinkIfMatch() will no-op as well since the link they passed us will
		// not match the current link.
		//
		// For #2, we may recreate the connection. It's possible we won't if the connection itself
		// has already been recovered by another goroutine.
		err := l.ns.Recover(ctx, lwid.ConnID)

		if err != nil {
			azlog.Writef(exported.EventConn, "(%s) Failure recovering connection for link: %s", lwid.String(), err)
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

	prefix := func() string {
		return prevLinkWithID.String()
	}

	return utils.Retry(ctx, eventName, prefix, retryOptions, func(ctx context.Context, args *utils.RetryFnArgs) error {
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
				azlog.Writef(exported.EventConn, "(%s, %s) Link was previously detached. Attempting quick reconnect to recover from error: %s", linkWithID.String(), operation, err.Error())
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
			ConnID:      l.links[partitionID].ConnID,
			Link:        *l.links[partitionID].Link,
			PartitionID: partitionID,
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
		ConnID:      l.links[partitionID].ConnID,
		Link:        *l.links[partitionID].Link,
		PartitionID: partitionID,
	}, nil
}

func (l *Links[LinkT]) GetManagementLink(ctx context.Context) (LinkWithID[amqpwrap.RPCLink], error) {
	if err := l.checkOpen(); err != nil {
		return LinkWithID[amqpwrap.RPCLink]{}, err
	}

	l.managementLinkMu.Lock()
	defer l.managementLinkMu.Unlock()

	if l.managementLink == nil {
		ls, err := l.newManagementLinkState(ctx)

		if err != nil {
			return LinkWithID[amqpwrap.RPCLink]{}, err
		}

		l.managementLink = ls
	}

	return LinkWithID[amqpwrap.RPCLink]{
		ConnID: l.managementLink.ConnID,
		Link:   *l.managementLink.Link,
	}, nil
}

func (l *Links[LinkT]) newLinkState(ctx context.Context, partitionID string) (*linkState[LinkT], error) {
	azlog.Writef(exported.EventConn, "Creating link for partition ID '%s'", partitionID)

	// check again now that we have the write lock
	ls := &linkState[LinkT]{
		PartitionID: partitionID,
	}

	cancelAuth, _, err := l.ns.NegotiateClaim(ctx, l.entityPathFn(partitionID))

	if err != nil {
		azlog.Writef(exported.EventConn, "(%s): Failed to negotiate claim for partition ID '%s': %s", ls.String(), partitionID, err)
		return nil, err
	}

	ls.cancelAuth = cancelAuth

	session, connID, err := l.ns.NewAMQPSession(ctx)

	if err != nil {
		azlog.Writef(exported.EventConn, "(%s): Failed to create AMQP session for partition ID '%s': %s", ls.String(), partitionID, err)
		_ = ls.Close(ctx)
		return nil, err
	}

	ls.session = session
	ls.ConnID = connID

	tmpLink, err := l.newLinkFn(ctx, session, l.entityPathFn(partitionID))

	if err != nil {
		azlog.Writef(exported.EventConn, "(%s): Failed to create link for partition ID '%s': %s", ls.String(), partitionID, err)
		_ = ls.Close(ctx)
		return nil, err
	}

	ls.Link = &tmpLink
	azlog.Writef(exported.EventConn, "(%s): Succesfully created link for partition ID '%s'", ls.String(), partitionID)
	return ls, nil
}

func (l *Links[LinkT]) newManagementLinkState(ctx context.Context) (*linkState[amqpwrap.RPCLink], error) {
	ls := &linkState[amqpwrap.RPCLink]{}

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
	cancelled := false

	if err := l.closeManagementLink(ctx); err != nil {
		azlog.Writef(exported.EventConn, "Error while cleaning up management link while doing connection recovery: %s", err.Error())

		if IsCancelError(err) {
			cancelled = true
		}
	}

	l.linksMu.Lock()
	defer l.linksMu.Unlock()

	tmpLinks := l.links
	l.links = nil

	for partitionID, link := range tmpLinks {
		if err := link.Close(ctx); err != nil {
			azlog.Writef(exported.EventConn, "Error while cleaning up link for partition ID '%s' while doing connection recovery: %s", partitionID, err.Error())

			if IsCancelError(err) {
				cancelled = true
			}
		}
	}

	if !permanent {
		l.links = map[string]*linkState[LinkT]{}
	}

	if cancelled {
		// this is the only kind of error I'd consider usable from Close() - it'll indicate
		// that some of the links haven't been cleanly closed.
		return ctx.Err()
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

// closePartitionLinkIfMatch will close the link in the cache if it matches the passed in linkName.
// This is similar to how an etag works - we'll only close it if you are working with the latest link -
// if not, it's a no-op since somebody else has already 'saved' (recovered) before you.
//
// Note that the only error that can be returned here will come from go-amqp. Cleanup of _our_ internal state
// will always happen, if needed.
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

	// PartitionID, if available.
	PartitionID string

	// cancelAuth cancels the backround claim negotation for this link.
	cancelAuth func()

	// optional session, if we created one for this
	// link.
	session amqpwrap.AMQPSession
}

// String returns a string that can be used for logging, of the format:
// (c:<connid>,l:<5 characters of link id>)
//
// It can also handle nil and partial initialization.
func (ls *linkState[LinkT]) String() string {
	if ls == nil {
		return "none"
	}

	linkName := ""

	if ls.Link != nil {
		linkName = (*ls.Link).LinkName()
	}

	return fmt.Sprintf("c:%d,l:%.5s,p:%s", ls.ConnID, linkName, ls.PartitionID)
}

// Close cancels the background authentication loop for this link and
// then closes the AMQP links.
// NOTE: this avoids any issues where closing fails on the broker-side or
// locally and we leak a goroutine.
func (ls *linkState[LinkT]) Close(ctx context.Context) error {
	if ls.cancelAuth != nil {
		ls.cancelAuth()
	}

	if ls.Link != nil {
		return (*ls.Link).Close(ctx)
	}

	return nil
}
