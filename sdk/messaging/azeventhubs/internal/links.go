// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"fmt"
	"sync"

	azlog "github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/amqpwrap"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azeventhubs/v2/internal/exported"
)

type AMQPLink interface {
	Close(ctx context.Context) error
	LinkName() string
}

// LinksForPartitionClient are the functions that the PartitionClient uses within Links[T]
// (for unit testing only)
type LinksForPartitionClient[LinkT AMQPLink] interface {
	// Retry is [Links.Retry]
	Retry(ctx context.Context, eventName azlog.Event, operation string, partitionID string, retryOptions exported.RetryOptions, fn func(ctx context.Context, lwid LinkWithID[LinkT]) error) error

	// Close is [Links.Close]
	Close(ctx context.Context) error
}

type Links[LinkT AMQPLink] struct {
	ns NamespaceForAMQPLinks

	linksMu            sync.RWMutex
	partitionCreation  creationGate
	managementCreation creationGate
	links              map[string]*linkState[LinkT]

	managementLink *linkState[amqpwrap.RPCLink]

	diagnostics    *Diagnostics
	component      string
	managementPath string
	newLinkFn      NewLinksFn[LinkT]
	entityPathFn   func(partitionID string) string

	lr LinkRetrier[LinkT]
	mr LinkRetrier[amqpwrap.RPCLink]
}

type NewLinksFn[LinkT AMQPLink] func(ctx context.Context, session amqpwrap.AMQPSession, entityPath string, partitionID string) (LinkT, error)

func NewLinks[LinkT AMQPLink](ns NamespaceForAMQPLinks, managementPath string, entityPathFn func(partitionID string) string, newLinkFn NewLinksFn[LinkT]) *Links[LinkT] {
	l := &Links[LinkT]{
		ns:             ns,
		links:          map[string]*linkState[LinkT]{},
		managementPath: managementPath,

		newLinkFn:    newLinkFn,
		entityPathFn: entityPathFn,
	}

	l.lr = LinkRetrier[LinkT]{
		GetLink:   l.GetLink,
		CloseLink: l.closePartitionLinkIfMatch,
		NSRecover: l.ns.Recover,
	}

	l.mr = LinkRetrier[amqpwrap.RPCLink]{
		GetLink: func(ctx context.Context, partitionID string) (LinkWithID[amqpwrap.RPCLink], error) {
			return l.GetManagementLink(ctx)
		},
		CloseLink: func(ctx context.Context, _, linkName string) error {
			return l.closeManagementLinkIfMatch(ctx, linkName)
		},
		NSRecover: l.ns.Recover,
	}

	return l
}

// WithDiagnostics configures AMQP link telemetry before the client is published.
func (l *Links[LinkT]) WithDiagnostics(d *Diagnostics, component string) *Links[LinkT] {
	l.diagnostics, l.component = d, component
	return l
}

func (l *Links[LinkT]) RetryManagement(ctx context.Context, eventName azlog.Event, operation string, retryOptions exported.RetryOptions, fn func(ctx context.Context, lwid LinkWithID[amqpwrap.RPCLink]) error) error {
	return l.mr.Retry(ctx, eventName, operation, "", retryOptions, fn)
}

func (l *Links[LinkT]) Retry(ctx context.Context, eventName azlog.Event, operation string, partitionID string, retryOptions exported.RetryOptions, fn func(ctx context.Context, lwid LinkWithID[LinkT]) error) error {
	return l.lr.Retry(ctx, eventName, operation, partitionID, retryOptions, fn)
}

func (l *Links[LinkT]) GetLink(ctx context.Context, partitionID string) (LinkWithID[LinkT], error) {
	return getOrCreateLink(ctx, &l.linksMu, &l.partitionCreation,
		func() (*linkState[LinkT], bool) { return l.links[partitionID], l.links == nil },
		func() (*linkState[LinkT], error) { return l.newLinkState(ctx, partitionID) },
		func(link *linkState[LinkT]) { l.links[partitionID] = link })
}

func (l *Links[LinkT]) GetManagementLink(ctx context.Context) (LinkWithID[amqpwrap.RPCLink], error) {
	return getOrCreateLink(ctx, &l.linksMu, &l.managementCreation,
		func() (*linkState[amqpwrap.RPCLink], bool) { return l.managementLink, l.links == nil },
		func() (*linkState[amqpwrap.RPCLink], error) { return l.newManagementLinkState(ctx) },
		func(link *linkState[amqpwrap.RPCLink]) { l.managementLink = link })
}

// lookup and store run under the cache lock. Creation and cleanup can call
// application code and run outside that lock.
func getOrCreateLink[T AMQPLink](ctx context.Context, mu *sync.RWMutex, gate *creationGate,
	lookup func() (*linkState[T], bool), create func() (*linkState[T], error), store func(*linkState[T])) (LinkWithID[T], error) {
	mu.RLock()
	link, closed := lookup()
	mu.RUnlock()
	if closed {
		return nil, ErrClientClosed
	}
	if link != nil {
		return link, nil
	}
	if err := gate.enter(ctx); err != nil {
		return nil, err
	}
	defer gate.leave()
	mu.RLock()
	link, closed = lookup()
	mu.RUnlock()
	if closed {
		return nil, ErrClientClosed
	}
	if link != nil {
		return link, nil
	}
	link, err := create()
	if err != nil {
		return nil, err
	}
	mu.Lock()
	_, closed = lookup()
	if !closed {
		store(link)
	}
	mu.Unlock()
	if closed {
		_ = link.Close(ctx)
		return nil, ErrClientClosed
	}
	return link, nil
}

func (l *Links[LinkT]) newLinkState(ctx context.Context, partitionID string) (*linkState[LinkT], error) {
	azlog.Writef(exported.EventConn, "Creating link for partition ID '%s'", partitionID)

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

	attachCtx, end := l.diagnostics.StartLifecycle(ctx, l.component, "attach")
	tmpLink, err := l.newLinkFn(attachCtx, session, l.entityPathFn(partitionID), partitionID)
	end(err)

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
	l.linksMu.Lock()
	management, partitions := l.managementLink, l.links
	l.managementLink, l.links = nil, nil
	l.linksMu.Unlock()
	cancelled := false
	if management != nil {
		if err := management.Close(ctx); err != nil {
			azlog.Writef(exported.EventConn, "Error while cleaning up management link: %s", err.Error())
			cancelled = IsCancelError(err)
		}
	}
	for partitionID, link := range partitions {
		if err := link.Close(ctx); err != nil {
			azlog.Writef(exported.EventConn, "Error while cleaning up link for partition ID '%s': %s", partitionID, err.Error())
			cancelled = cancelled || IsCancelError(err)
		}
	}
	if cancelled {
		return ctx.Err()
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
	l.linksMu.Lock()
	current, exists := l.links[partitionID]
	if !exists || current.Link().LinkName() != linkName {
		l.linksMu.Unlock()
		return nil
	}
	delete(l.links, partitionID)
	l.linksMu.Unlock()
	err := current.Close(ctx)
	l.diagnostics.recordLifecycle(ctx, l.component, "invalidated")
	return err
}

func (l *Links[LinkT]) closeManagementLinkIfMatch(ctx context.Context, linkName string) error {
	l.linksMu.Lock()
	link := l.managementLink
	if link == nil || link.Link().LinkName() != linkName {
		l.linksMu.Unlock()
		return nil
	}
	l.managementLink = nil
	l.linksMu.Unlock()
	return link.Close(ctx)
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

	var linkCloseErr error

	if ls.link != nil {
		// we're more interested in a link failing to close than we are in
		// the session.
		linkCloseErr = ls.Link().Close(ctx)
	}

	if ls.session != nil {
		_ = ls.session.Close(ctx)
	}

	return linkCloseErr
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
