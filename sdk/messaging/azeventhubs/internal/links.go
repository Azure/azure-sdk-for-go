// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"errors"
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

// LinksForPartitionClient are the functions that the PartitionClient uses within Links[T]
// (for unit testing only)
type LinksForPartitionClient[LinkT AMQPLink] interface {
	RecoverIfNeeded(ctx context.Context, err error) error
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
	newLinkFn      NewLinksFn[LinkT]
	entityPathFn   func(partitionID string) string
}

type NewLinksFn[LinkT AMQPLink] func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string, partitionID string) (LinkT, error)

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

func (l *Links[LinkT]) RecoverIfNeeded(ctx context.Context, err error) error {
	rk := GetRecoveryKind(err)

	switch rk {
	case RecoveryKindNone:
		return nil
	case RecoveryKindLink:
		var awErr amqpwrap.Error

		if !errors.As(err, &awErr) {
			log.Writef(exported.EventConn, "RecoveryKindLink, but not an amqpwrap.Error: %T,%v", err, err)
			return nil
		}

		if err := l.closePartitionLinkIfMatch(ctx, awErr.PartitionID, awErr.LinkName); err != nil {
			azlog.Writef(exported.EventConn, "(%s) Error when cleaning up old link for link recovery: %s", formatLogPrefix(awErr.ConnID, awErr.LinkName, awErr.PartitionID), err)
			return err
		}

		return nil
	case RecoveryKindConn:
		var awErr amqpwrap.Error

		if !errors.As(err, &awErr) {
			log.Writef(exported.EventConn, "RecoveryKindConn, but not an amqpwrap.Error: %T,%v", err, err)
			return nil
		}

		// We only close _this_ partition's link. Other partitions will also get an error, and will recover.
		// We used to close _all_ the links, but no longer do that since it's possible (when we do receiver
		// redirect) to have more than one active connection at a time which means not all links would be
		// affected when a single connection goes down.
		if err := l.closePartitionLinkIfMatch(ctx, awErr.PartitionID, awErr.LinkName); err != nil {
			azlog.Writef(exported.EventConn, "(%s) Error when cleaning up old link: %s", formatLogPrefix(awErr.ConnID, awErr.LinkName, awErr.PartitionID), err)

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
		err := l.ns.Recover(ctx, awErr.ConnID)

		if err != nil {
			azlog.Writef(exported.EventConn, "(%s) Failure recovering connection for link: %s", formatLogPrefix(awErr.ConnID, awErr.LinkName, awErr.PartitionID), err)
			return err
		}

		return nil
	default:
		return err
	}
}

func (l *Links[LinkT]) Retry(ctx context.Context, eventName log.Event, operation string, partitionID string, retryOptions exported.RetryOptions, fn func(ctx context.Context, lwid LinkWithID[LinkT]) error) error {
	didQuickRetry := false

	isFatalErrorFunc := func(err error) bool {
		return GetRecoveryKind(err) == RecoveryKindFatal
	}

	currentPrefix := ""

	prefix := func() string {
		return currentPrefix
	}

	return utils.Retry(ctx, eventName, prefix, retryOptions, func(ctx context.Context, args *utils.RetryFnArgs) error {
		if err := l.RecoverIfNeeded(ctx, args.LastErr); err != nil {
			return err
		}

		linkWithID, err := l.GetLink(ctx, partitionID)

		if err != nil {
			return err
		}

		currentPrefix = linkWithID.String()

		if err := fn(ctx, linkWithID); err != nil {
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

func (l *Links[LinkT]) GetLink(ctx context.Context, partitionID string) (LinkWithID[LinkT], error) {
	if err := l.checkOpen(); err != nil {
		return nil, err
	}

	l.linksMu.RLock()
	current := l.links[partitionID]
	l.linksMu.RUnlock()

	if current != nil {
		return current, nil
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
		current = ls
	}

	return current, nil
}

func (l *Links[LinkT]) GetManagementLink(ctx context.Context) (LinkWithID[amqpwrap.RPCLink], error) {
	if err := l.checkOpen(); err != nil {
		return nil, err
	}

	l.managementLinkMu.Lock()
	defer l.managementLinkMu.Unlock()

	if l.managementLink == nil {
		ls, err := l.newManagementLinkState(ctx)

		if err != nil {
			return nil, err
		}

		l.managementLink = ls
	}

	return l.managementLink, nil
}

func (l *Links[LinkT]) newLinkState(ctx context.Context, partitionID string) (*linkState[LinkT], error) {
	azlog.Writef(exported.EventConn, "Creating link for partition ID '%s'", partitionID)

	// check again now that we have the write lock
	ls := &linkState[LinkT]{
		partitionID: partitionID,
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
	ls.connID = connID

	tmpLink, err := l.newLinkFn(ctx, session, l.entityPathFn(partitionID), partitionID)

	if err != nil {
		azlog.Writef(exported.EventConn, "(%s): Failed to create link for partition ID '%s': %s", ls.String(), partitionID, err)
		_ = ls.Close(ctx)
		return nil, err
	}

	ls.link = &tmpLink

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

	ls.connID = connID
	ls.link = &tmpRPCLink

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
		current.Link().LinkName() != linkName { // we've already created a new link, their link was stale.
		return nil
	}

	l.linksMu.Lock()
	defer l.linksMu.Unlock()

	current, exists = l.links[partitionID]

	if !exists ||
		current.Link().LinkName() != linkName { // we've already created a new link, their link was stale.
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
	// connID is an arbitrary (but unique) integer that represents the
	// current connection. This comes back from the Namespace, anytime
	// it hands back a connection.
	connID uint64

	// link will be either an [amqpwrap.AMQPSenderCloser], [amqpwrap.AMQPReceiverCloser] or [amqpwrap.RPCLink]
	link *LinkT

	// partitionID, if available.
	partitionID string

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

	if ls.link != nil {
		linkName = ls.Link().LinkName()
	}

	return formatLogPrefix(ls.connID, linkName, ls.partitionID)
}

// Close cancels the background authentication loop for this link and
// then closes the AMQP links.
// NOTE: this avoids any issues where closing fails on the broker-side or
// locally and we leak a goroutine.
func (ls *linkState[LinkT]) Close(ctx context.Context) error {
	if ls.cancelAuth != nil {
		ls.cancelAuth()
	}

	if ls.link != nil {
		return ls.Link().Close(ctx)
	}

	return nil
}

func (ls *linkState[LinkT]) PartitionID() string {
	return ls.partitionID
}

func (ls *linkState[LinkT]) ConnID() uint64 {
	return ls.connID
}

func (ls *linkState[LinkT]) Link() LinkT {
	return *ls.link
}

// LinkWithID is a readonly interface over the top of a linkState.
type LinkWithID[LinkT AMQPLink] interface {
	ConnID() uint64
	Link() LinkT
	PartitionID() string
	Close(ctx context.Context) error
	String() string
}

func formatLogPrefix(connID uint64, linkName, partitionID string) string {
	return fmt.Sprintf("c:%d,l:%.5s,p:%s", connID, linkName, partitionID)
}
