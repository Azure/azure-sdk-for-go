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
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/sbauth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/azservicebus/internal/tracing"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/auth"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/cbs"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/conn"
	"github.com/Azure/azure-sdk-for-go/sdk/messaging/internal/rpc"
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
		FQDN          string
		TokenProvider *sbauth.TokenProvider
		tlsConfig     *tls.Config
		userAgent     string

		newWebSocketConn func(ctx context.Context, args NewWebSocketConnArgs) (net.Conn, error)

		baseRetrier Retrier

		clientMu       sync.Mutex
		clientRevision uint64
		client         *amqp.Client

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
	NewMgmtClient(ctx context.Context, links AMQPLinks) (MgmtClient, error)
	GetEntityAudience(entityPath string) string
	Recover(ctx context.Context, clientRevision uint64) error
}

// NamespaceForAMQPLinks is the Namespace surface needed for the *MgmtClient.
type NamespaceForMgmtClient interface {
	NewRPCLink(ctx context.Context, managementPath string) (*rpc.Link, error)
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

// NamespacesWithTokenCredential sets the token provider on the namespace
// fullyQualifiedNamespace is the Service Bus namespace name (ex: myservicebus.servicebus.windows.net)
func NamespacesWithTokenCredential(fullyQualifiedNamespace string, tokenCredential azcore.TokenCredential) NamespaceOption {
	return func(ns *Namespace) error {
		ns.TokenProvider = sbauth.NewTokenProvider(tokenCredential)
		ns.FQDN = fullyQualifiedNamespace
		return nil
	}
}

// NewNamespace creates a new namespace configured through NamespaceOption(s)
func NewNamespace(opts ...NamespaceOption) (*Namespace, error) {
	ns := &Namespace{
		baseRetrier: NewBackoffRetrier(struct {
			MaxRetries int
			Factor     float64
			Jitter     bool
			Min        time.Duration
			Max        time.Duration
		}{
			Factor:     2,
			Min:        time.Second,
			Max:        time.Minute,
			MaxRetries: 10,
		}),
	}

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
func (ns *Namespace) NewAMQPSession(ctx context.Context) (AMQPSessionCloser, uint64, error) {
	client, clientRevision, err := ns.getAMQPClientImpl(ctx)

	if err != nil {
		return nil, 0, err
	}

	session, err := client.NewSession()

	if err != nil {
		return nil, 0, err
	}

	return session, clientRevision, err
}

// NewMgmtClient creates a new management client with the internally cached *amqp.Client.
func (ns *Namespace) NewMgmtClient(ctx context.Context, l AMQPLinks) (MgmtClient, error) {
	return newMgmtClient(ctx, l, ns)
}

// NewRPCLink creates a new amqp-common *rpc.Link with the internally cached *amqp.Client.
func (ns *Namespace) NewRPCLink(ctx context.Context, managementPath string) (*rpc.Link, error) {
	client, _, err := ns.getAMQPClientImpl(ctx)

	if err != nil {
		return nil, err
	}

	return rpc.NewLink(client, managementPath)
}

// NewAMQPLinks creates an AMQPLinks struct, which groups together the commonly needed links for
// working with Service Bus.
func (ns *Namespace) NewAMQPLinks(entityPath string, createLinkFunc CreateLinkFunc) AMQPLinks {
	return newAMQPLinks(ns, entityPath, ns.baseRetrier, createLinkFunc)
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

// Recover destroys the currently held client and recreates it.
// clientRevision being nil will recover without a revision check.
func (ns *Namespace) Recover(ctx context.Context, clientRevision uint64) error {
	ns.clientMu.Lock()
	defer ns.clientMu.Unlock()

	_, span := tab.StartSpan(ctx, tracing.SpanRecoverClient)
	defer span.End()

	span.AddAttributes(
		tab.Int64Attribute("revision", int64(ns.clientRevision)),
		tab.Int64Attribute("requested", int64(clientRevision)))

	if ns.clientRevision > clientRevision {
		span.Logger().Info(fmt.Sprintf("Skipping recovery, already recovered: %d vs %d", ns.clientRevision, clientRevision))
		// we've already recovered since the client last tried.
		return nil
	}

	if ns.client != nil {
		oldClient := ns.client
		ns.client = nil

		// the error on close isn't critical
		go func() {
			span.Logger().Info(fmt.Sprintf("Closing old client (client:%d,passed in:%d)", ns.clientRevision, clientRevision))
			err := oldClient.Close()
			tab.For(ctx).Error(err)
		}()
	}

	var err error
	span.Logger().Info(fmt.Sprintf("Creating a new client (client:%d,passed in:%d)", ns.clientRevision, clientRevision))
	ns.client, err = ns.newClient(ctx)

	if err == nil {
		span.AddAttributes(tab.Int64Attribute("newcr", int64(ns.clientRevision)))
		ns.clientRevision++
	}

	return err
}

// negotiateClaim performs initial authentication and starts periodic refresh of credentials.
// the returned func is to cancel() the refresh goroutine.
func (ns *Namespace) NegotiateClaim(ctx context.Context, entityPath string) (func() <-chan struct{}, error) {
	return ns.startNegotiateClaimRenewer(ctx,
		entityPath,
		cbs.NegotiateClaim,
		ns.getAMQPClientImpl,
		nextClaimRefreshDuration)
}

func (ns *Namespace) startNegotiateClaimRenewer(ctx context.Context,
	entityPath string,
	cbsNegotiateClaim func(ctx context.Context, audience string, conn *amqp.Client, provider auth.TokenProvider) error,
	nsGetAMQPClientImpl func(ctx context.Context) (*amqp.Client, uint64, error),
	nextClaimRefreshDurationFn func(expirationTime time.Time, currentTime time.Time) time.Duration) (func() <-chan struct{}, error) {
	audience := ns.GetEntityAudience(entityPath)

	refreshClaim := func() (time.Time, error) {
		retrier := ns.baseRetrier.Copy()

		var lastErr error
		var expiration time.Time

		for retrier.Try(ctx) {
			expiration, lastErr = func() (time.Time, error) {
				ctx, span := ns.startSpanFromContext(ctx, tracing.SpanNegotiateClaim)
				defer span.End()

				amqpClient, clientRevision, err := nsGetAMQPClientImpl(ctx)

				if err != nil {
					span.Logger().Error(err)
					return time.Time{}, err
				}

				token, expiration, err := ns.TokenProvider.GetTokenAsTokenProvider(audience)

				if err != nil {
					span.Logger().Error(err)
					return time.Time{}, err
				}

				// You're not allowed to have multiple $cbs links open in a single connection.
				// The current cbs.NegotiateClaim implementation automatically creates and shuts
				// down it's own link so we have to guard against that here.
				ns.negotiateClaimMu.Lock()
				err = cbsNegotiateClaim(ctx, audience, amqpClient, token)
				ns.negotiateClaimMu.Unlock()

				if err != nil {
					if shouldRecreateConnection(ctx, err) {
						if err := ns.Recover(ctx, clientRevision); err != nil {
							span.Logger().Error(fmt.Errorf("connection recovery failed: %w", err))
						}
					}

					span.Logger().Error(err)
					return time.Time{}, err
				}

				return expiration, nil
			}()

			if lastErr == nil {
				break
			}
		}

		return expiration, lastErr
	}

	expiresOn, err := refreshClaim()

	if err != nil {
		return nil, err
	}

	// start the periodic refresh of credentials
	refreshCtx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-refreshCtx.Done():
				return
			case <-time.After(nextClaimRefreshDurationFn(expiresOn, time.Now())):
				tmpExpiresOn, err := refreshClaim() // logging will report the error for now

				if err == nil {
					expiresOn = tmpExpiresOn
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

func (ns *Namespace) getAMQPClientImpl(ctx context.Context) (*amqp.Client, uint64, error) {
	ns.clientMu.Lock()
	defer ns.clientMu.Unlock()

	if ns.client != nil {
		return ns.client, ns.clientRevision, nil
	}

	var err error
	retrier := ns.baseRetrier.Copy()

	for retrier.Try(ctx) {
		ns.client, err = ns.newClient(ctx)

		if err == nil {
			break
		}
	}

	return ns.client, ns.clientRevision, err
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
