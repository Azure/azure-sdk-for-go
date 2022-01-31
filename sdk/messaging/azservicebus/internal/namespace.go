// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"runtime"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/internal/log"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/sbauth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/utils"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/cbs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/conn"
	"github.com/Azure/go-amqp"
	"github.com/devigned/tab"
)

const (
	rootUserAgent = "/azsdk-go-azservicebus/" + Version
)

type (
	// Namespace is an abstraction over an amqp.Client, allowing us to hold onto a single
	// instance of a connection per ServiceBusClient.
	Namespace struct {
		// NOTE: values need to be 64-bit aligned. Simplest way to make sure this happens
		// is just to make it the first value in the struct
		// See:
		//   Godoc: https://pkg.go.dev/sync/atomic#pkg-note-BUG
		//   PR: https://github.com/Azure/azure-sdk-for-go/pull/16847
		connID uint64

		FQDN          string
		TokenProvider *sbauth.TokenProvider
		tlsConfig     *tls.Config
		userAgent     string

		newWebSocketConn func(ctx context.Context, args NewWebSocketConnArgs) (net.Conn, error)

		retryOptions utils.RetryOptions

		clientMu         sync.RWMutex
		client           *amqp.Client
		negotiateClaimMu sync.Mutex
	}

	// NamespaceOption provides structure for configuring a new Service Bus namespace
	NamespaceOption func(h *Namespace) error
)

// NamespaceWithNewAMQPLinks is the Namespace surface for consumers of AMQPLinks.
type NamespaceWithNewAMQPLinks interface {
	NewAMQPLinks(entityPath string, createLinkFunc CreateLinkFunc) AMQPLinks
}

// NamespaceForAMQPLinks is the Namespace surface needed for the internals of AMQPLinks.
type NamespaceForAMQPLinks interface {
	NegotiateClaim(ctx context.Context, entityPath string) (func() <-chan struct{}, error)
	NewAMQPSession(ctx context.Context) (AMQPSessionCloser, uint64, error)
	NewRPCLink(ctx context.Context, managementPath string) (RPCLink, error)
	GetEntityAudience(entityPath string) string
	Recover(ctx context.Context, clientRevision uint64) (bool, error)
}

// NamespaceWithConnectionString configures a namespace with the information provided in a Service Bus connection string
func NamespaceWithConnectionString(connStr string) NamespaceOption {
	return func(ns *Namespace) error {
		parsed, err := conn.ParsedConnectionFromStr(connStr)
		if err != nil {
			return err
		}

		if parsed.Namespace != "" {
			ns.FQDN = parsed.Namespace
		}

		provider, err := sbauth.NewTokenProviderWithConnectionString(parsed.KeyName, parsed.Key)
		if err != nil {
			return err
		}

		ns.TokenProvider = provider
		return nil
	}
}

// NamespaceWithTLSConfig appends to the TLS config.
func NamespaceWithTLSConfig(tlsConfig *tls.Config) NamespaceOption {
	return func(ns *Namespace) error {
		ns.tlsConfig = tlsConfig
		return nil
	}
}

// NamespaceWithUserAgent appends to the root user-agent value.
func NamespaceWithUserAgent(userAgent string) NamespaceOption {
	return func(ns *Namespace) error {
		ns.userAgent = userAgent
		return nil
	}
}

// NewWebSocketConnArgs are the arguments to the NewWebSocketConn function you pass if you want
// to enable websockets.
type NewWebSocketConnArgs struct {
	// NOTE: this struct is exported via client.go:NewWebSocketConnArgs

	// Host is the the `wss://<host>` to connect to
	Host string
}

// NamespaceWithWebSocket configures the namespace and all entities to use wss:// rather than amqps://
func NamespaceWithWebSocket(newWebSocketConn func(ctx context.Context, args NewWebSocketConnArgs) (net.Conn, error)) NamespaceOption {
	return func(ns *Namespace) error {
		ns.newWebSocketConn = newWebSocketConn
		return nil
	}
}

// NamespaceWithTokenCredential sets the token provider on the namespace
// fullyQualifiedNamespace is the Service Bus namespace name (ex: myservicebus.servicebus.windows.net)
func NamespaceWithTokenCredential(fullyQualifiedNamespace string, tokenCredential azcore.TokenCredential) NamespaceOption {
	return func(ns *Namespace) error {
		ns.TokenProvider = sbauth.NewTokenProvider(tokenCredential)
		ns.FQDN = fullyQualifiedNamespace
		return nil
	}
}

func NamespaceWithRetryOptions(retryOptions utils.RetryOptions) NamespaceOption {
	return func(ns *Namespace) error {
		ns.retryOptions = retryOptions
		return nil
	}
}

// NewNamespace creates a new namespace configured through NamespaceOption(s)
func NewNamespace(opts ...NamespaceOption) (*Namespace, error) {
	ns := &Namespace{}

	for _, opt := range opts {
		err := opt(ns)
		if err != nil {
			return nil, err
		}
	}

	return ns, nil
}

func (ns *Namespace) newClient(ctx context.Context) (*amqp.Client, error) {
	ctx, span := ns.startSpanFromContext(ctx, "sb.namespace.newClient")
	defer span.End()
	defaultConnOptions := []amqp.ConnOption{
		amqp.ConnSASLAnonymous(),
		amqp.ConnMaxSessions(65535),
		amqp.ConnProperty("product", "MSGolangClient"),
		amqp.ConnProperty("version", Version),
		amqp.ConnProperty("platform", runtime.GOOS),
		amqp.ConnProperty("framework", runtime.Version()),
		amqp.ConnProperty("user-agent", ns.getUserAgent()),
	}

	if ns.tlsConfig != nil {
		defaultConnOptions = append(
			defaultConnOptions,
			amqp.ConnTLS(true),
			amqp.ConnTLSConfig(ns.tlsConfig),
		)
	}

	if ns.newWebSocketConn != nil {
		nConn, err := ns.newWebSocketConn(ctx, NewWebSocketConnArgs{
			Host: ns.getWSSHostURI() + "$servicebus/websocket",
		})

		if err != nil {
			return nil, err
		}

		return amqp.New(nConn, append(defaultConnOptions, amqp.ConnServerHostname(ns.FQDN))...)
	}

	return amqp.Dial(ns.getAMQPHostURI(), defaultConnOptions...)
}

// NewAMQPSession creates a new AMQP session with the internally cached *amqp.Client.
// Returns a closeable AMQP session and the current client revision.
func (ns *Namespace) NewAMQPSession(ctx context.Context) (AMQPSessionCloser, uint64, error) {
	client, clientRevision, err := ns.GetAMQPClientImpl(ctx)

	if err != nil {
		return nil, 0, err
	}

	session, err := client.NewSession()

	if err != nil {
		return nil, 0, err
	}

	return session, clientRevision, err
}

// NewRPCLink creates a new amqp-common *rpc.Link with the internally cached *amqp.Client.
func (ns *Namespace) NewRPCLink(ctx context.Context, managementPath string) (RPCLink, error) {
	client, _, err := ns.GetAMQPClientImpl(ctx)

	if err != nil {
		return nil, err
	}

	return NewRPCLink(client, managementPath)
}

// NewAMQPLinks creates an AMQPLinks struct, which groups together the commonly needed links for
// working with Service Bus.
func (ns *Namespace) NewAMQPLinks(entityPath string, createLinkFunc CreateLinkFunc) AMQPLinks {
	return NewAMQPLinks(ns, entityPath, createLinkFunc)
}

// Close closes the current cached client.
func (ns *Namespace) Close(ctx context.Context) error {
	ns.clientMu.Lock()
	defer ns.clientMu.Unlock()

	if ns.client != nil {
		return ns.client.Close()
	}

	return nil
}

// Recover destroys the currently held AMQP connection and recreates it, if needed.
// If a new is actually created (rather than just cached) then the returned bool
// will be true. Any links that were created from the original connection will need to
// be recreated.
func (ns *Namespace) Recover(ctx context.Context, theirConnID uint64) (bool, error) {
	ns.clientMu.Lock()
	defer ns.clientMu.Unlock()

	_, span := tab.StartSpan(ctx, tracing.SpanRecoverClient)
	defer span.End()

	span.AddAttributes(
		tab.Int64Attribute("connID", int64(ns.connID)),
		tab.Int64Attribute("theirConnID", int64(theirConnID)))

	if ns.connID != theirConnID {
		log.Writef(EventConn, "Skipping connection recovery, already recovered: %d vs %d", ns.connID, theirConnID)
		// we've already recovered since the client last tried.
		return false, nil
	}

	if ns.client != nil {
		oldClient := ns.client
		ns.client = nil

		// the error on close isn't critical
		_ = oldClient.Close()
	}

	var err error
	log.Writef(EventConn, "Creating a new client (rev:%d)", ns.connID)
	ns.client, err = ns.newClient(ctx)

	if err != nil {
		return false, err
	}

	ns.connID++
	log.Writef(EventConn, "New client created, (rev: %d)", ns.connID)
	return true, nil
}

// negotiateClaim performs initial authentication and starts periodic refresh of credentials.
// the returned func is to cancel() the refresh goroutine.
func (ns *Namespace) NegotiateClaim(ctx context.Context, entityPath string) (func() <-chan struct{}, error) {
	return ns.startNegotiateClaimRenewer(ctx,
		entityPath,
		cbs.NegotiateClaim,
		ns.GetAMQPClientImpl,
		nextClaimRefreshDuration)
}

func (ns *Namespace) startNegotiateClaimRenewer(ctx context.Context,
	entityPath string,
	cbsNegotiateClaim func(ctx context.Context, audience string, conn *amqp.Client, provider auth.TokenProvider) error,
	nsGetAMQPClientImpl func(ctx context.Context) (*amqp.Client, uint64, error),
	nextClaimRefreshDurationFn func(expirationTime time.Time, currentTime time.Time) time.Duration) (func() <-chan struct{}, error) {
	audience := ns.GetEntityAudience(entityPath)

	refreshClaim := func(ctx context.Context) (time.Time, error) {
		log.Writef(EventAuth, "(%s) refreshing claim", entityPath)
		ctx, span := ns.startSpanFromContext(ctx, tracing.SpanNegotiateClaim)
		defer span.End()

		amqpClient, clientRevision, err := nsGetAMQPClientImpl(ctx)

		if err != nil {
			return time.Time{}, err
		}

		token, expiration, err := ns.TokenProvider.GetTokenAsTokenProvider(audience)

		if err != nil {
			log.Writef(EventAuth, "(%s) negotiate claim, failed getting token: %s", entityPath, err.Error())
			return time.Time{}, err
		}

		log.Writef(EventAuth, "(%s) negotiate claim, token expires on %s", entityPath, expiration.Format(time.RFC3339))

		// You're not allowed to have multiple $cbs links open in a single connection.
		// The current cbs.NegotiateClaim implementation automatically creates and shuts
		// down it's own link so we have to guard against that here.
		ns.negotiateClaimMu.Lock()
		err = cbsNegotiateClaim(ctx, audience, amqpClient, token)
		ns.negotiateClaimMu.Unlock()

		sbe := GetSBErrInfo(err)

		if sbe != nil {
			// Note we only handle connection recovery here since (currently)
			// the negotiateClaim code creates it's own link each time.
			if sbe.RecoveryKind == RecoveryKindConn {
				if _, err := ns.Recover(ctx, clientRevision); err != nil {
					log.Writef(EventAuth, "(%s) negotiate claim, failed in connection recovery: %s", entityPath, err)
				}
			}

			log.Writef(EventAuth, "(%s) negotiate claim, failed: %s", entityPath, err.Error())
			return time.Time{}, err
		}

		return expiration, nil
	}

	expiresOn, err := refreshClaim(ctx)

	if err != nil {
		return nil, err
	}

	// start the periodic refresh of credentials
	refreshCtx, cancel := context.WithCancel(context.Background())

	go func() {
	TokenRefreshLoop:
		for {
			nextClaimAt := nextClaimRefreshDurationFn(expiresOn, time.Now())

			log.Writef(EventAuth, "(%s) next refresh in %s", entityPath, nextClaimAt)

			select {
			case <-refreshCtx.Done():
				return
			case <-time.After(nextClaimAt):
				for {
					err := utils.Retry(refreshCtx, "claimrefresh", func(ctx context.Context, args *utils.RetryFnArgs) error {
						tmpExpiresOn, err := refreshClaim(ctx)

						if err != nil {
							return err
						}

						expiresOn = tmpExpiresOn
						return nil
					}, IsFatalSBError, ns.retryOptions)

					if err == nil {
						break
					}

					// if we fail our retries _and_ we've exceeded the window where our token would have
					// been good we can just stop.
					if time.Since(expiresOn) <= 0 {
						log.Writef(EventAuth, "[%s] token has expired, stopping refresh loop", entityPath)
						break TokenRefreshLoop
					}
				}
			}
		}
	}()

	cancelRefresh := func() <-chan struct{} {
		cancel()
		return refreshCtx.Done()
	}

	return cancelRefresh, nil
}

func (ns *Namespace) GetAMQPClientImpl(ctx context.Context) (*amqp.Client, uint64, error) {
	ns.clientMu.Lock()
	defer ns.clientMu.Unlock()

	if ns.client != nil {
		return ns.client, ns.connID, nil
	}

	var err error
	ns.client, err = ns.newClient(ctx)

	if err == nil {
		ns.connID++
	}

	return ns.client, ns.connID, err
}

func (ns *Namespace) getWSSHostURI() string {
	return fmt.Sprintf("wss://%s/", ns.FQDN)
}

func (ns *Namespace) getAMQPHostURI() string {
	return fmt.Sprintf("amqps://%s/", ns.FQDN)
}

func (ns *Namespace) GetHTTPSHostURI() string {
	return fmt.Sprintf("https://%s/", ns.FQDN)
}

func (ns *Namespace) GetEntityAudience(entityPath string) string {
	return ns.getAMQPHostURI() + entityPath
}

func (ns *Namespace) getUserAgent() string {
	userAgent := rootUserAgent
	if ns.userAgent != "" {
		userAgent = fmt.Sprintf("%s/%s", userAgent, ns.userAgent)
	}
	return userAgent
}

func (ns *Namespace) startSpanFromContext(ctx context.Context, operationName string) (context.Context, tab.Spanner) {
	ctx, span := tab.StartSpan(ctx, operationName)
	tracing.ApplyComponentInfo(span, Version)
	return ctx, span
}

// nextClaimRefreshDuration figures out the proper interval for the next authorization
// refresh.
//
// It applies a few real world adjustments:
// - We assume the expiration time is 10 minutes ahead of when it actually is, to adjust for clock drift.
// - We don't let the refresh interval fall below 2 minutes
// - We don't let the refresh interval go above 49 days
//
// This logic is from here:
// https://github.com/Azure/azure-sdk-for-net/blob/bfd3109d0f9afa763131731d78a31e39c81101b3/sdk/servicebus/Azure.Messaging.ServiceBus/src/Amqp/AmqpConnectionScope.cs#L998
func nextClaimRefreshDuration(expirationTime time.Time, currentTime time.Time) time.Duration {
	const min = 2 * time.Minute
	const max = 49 * 24 * time.Hour
	const clockDrift = 10 * time.Minute

	var refreshDuration = expirationTime.Sub(currentTime) - clockDrift

	if refreshDuration < min {
		return min
	} else if refreshDuration > max {
		return max
	}

	return refreshDuration
}
