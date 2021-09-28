// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package internal

import (
	"context"
	"crypto/tls"
	"fmt"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/Azure/azure-amqp-common-go/v3/auth"
	"github.com/Azure/azure-amqp-common-go/v3/cbs"
	"github.com/Azure/azure-amqp-common-go/v3/conn"
	"github.com/Azure/azure-amqp-common-go/v3/rpc"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/go-amqp"
	"github.com/Azure/go-autorest/autorest/azure"
	"github.com/devigned/tab"
	"nhooyr.io/websocket"
)

const (
	rootUserAgent = "/golang-service-bus"
)

type (
	// Namespace is an abstraction over an amqp.Client, allowing us to hold onto a single
	// instance of a connection per ServiceBusClient.
	Namespace struct {
		Name          string
		Suffix        string
		TokenProvider *tokenProvider
		Environment   azure.Environment
		tlsConfig     *tls.Config
		userAgent     string
		useWebSocket  bool

		baseRetrier Retrier

		clientMu sync.Mutex
		client   *amqp.Client

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
	NewAMQPSession(ctx context.Context) (AMQPSessionCloser, error)
	NewMgmtClient(ctx context.Context, managementPath string) (MgmtClient, error)
	GetEntityAudience(entityPath string) string
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
			ns.Name = parsed.Namespace
		}

		if parsed.Suffix != "" {
			ns.Suffix = parsed.Suffix
		}

		provider, err := newTokenProviderWithConnectionString(parsed.KeyName, parsed.Key)
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

// NamespaceWithWebSocket configures the namespace and all entities to use wss:// rather than amqps://
func NamespaceWithWebSocket() NamespaceOption {
	return func(ns *Namespace) error {
		ns.useWebSocket = true
		return nil
	}
}

// // NamespaceWithEnvironmentBinding configures a namespace using the environment details. It uses one of the following methods:
// //
// // 1. Client Credentials: attempt to authenticate with a Service Principal via "AZURE_TENANT_ID", "AZURE_CLIENT_ID" and
// //    "AZURE_CLIENT_SECRET"
// //
// // 2. Client Certificate: attempt to authenticate with a Service Principal via "AZURE_TENANT_ID", "AZURE_CLIENT_ID",
// //    "AZURE_CERTIFICATE_PATH" and "AZURE_CERTIFICATE_PASSWORD"
// //
// // 3. Managed Identity (MI): attempt to authenticate via the MI assigned to the Azure resource
// //
// //
// // The Azure Environment used can be specified using the name of the Azure Environment set in "AZURE_ENVIRONMENT" var.
// func NamespaceWithEnvironmentBinding(name string) NamespaceOption {
// 	return func(ns *Namespace) error {
// 		provider, err := aad.NewJWTProvider(
// 			aad.JWTProviderWithEnvironmentVars(),
// 			// TODO: fix bug upstream to use environment resourceURI
// 			aad.JWTProviderWithResourceURI(ns.getResourceURI()),
// 		)
// 		if err != nil {
// 			return err
// 		}

// 		ns.TokenProvider = provider
// 		ns.Name = name
// 		return nil
// 	}
// }

// NamespaceWithAzureEnvironment sets the namespace's Environment, Suffix and ResourceURI parameters according
// to the Azure Environment defined in "github.com/Azure/go-autorest/autorest/azure" package.
// This allows to configure the library to be used in the different Azure clouds.
// environmentName is the name of the cloud as defined in autorest : https://github.com/Azure/go-autorest/blob/b076c1437d051bf4c328db428b70f4fe22ad38b0/autorest/azure/environments.go#L34-L39
func NamespaceWithAzureEnvironment(namespaceName, environmentName string) NamespaceOption {
	return func(ns *Namespace) error {
		azureEnv, err := azure.EnvironmentFromName(environmentName)
		if err != nil {
			return err
		}
		ns.Environment = azureEnv
		ns.Suffix = azureEnv.ServiceBusEndpointSuffix
		ns.Name = namespaceName
		return nil
	}
}

// NamespacesWithTokenCredential sets the token provider on the namespace
func NamespacesWithTokenCredential(namespaceName string, tokenCredential azcore.TokenCredential) NamespaceOption {
	return func(ns *Namespace) error {
		ns.TokenProvider = newTokenProviderWithTokenCredential(tokenCredential)

		parts := strings.SplitN(namespaceName, ".", 2)

		if len(parts) != 2 {
			ns.Name = parts[0]
		} else {
			ns.Name, ns.Suffix = parts[0], parts[1]
		}

		return nil
	}
}

// NewNamespace creates a new namespace configured through NamespaceOption(s)
func NewNamespace(opts ...NamespaceOption) (*Namespace, error) {
	ns := &Namespace{
		Environment: azure.PublicCloud,
		baseRetrier: NewBackoffRetrier(struct {
			MaxRetries int
			Factor     float64
			Jitter     bool
			Min        time.Duration
			Max        time.Duration
		}{
			MaxRetries: 5,
			Factor:     1,
			Min:        5 * time.Second,
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

	if ns.useWebSocket {
		wssHost := ns.getWSSHostURI() + "$servicebus/websocket"
		opts := &websocket.DialOptions{Subprotocols: []string{"amqp"}}
		wssConn, _, err := websocket.Dial(ctx, wssHost, opts)

		if err != nil {
			return nil, err
		}
		nConn := websocket.NetConn(context.Background(), wssConn, websocket.MessageBinary)

		return amqp.New(nConn, append(defaultConnOptions, amqp.ConnServerHostname(ns.GetHostname()))...)
	}

	return amqp.Dial(ns.getAMQPHostURI(), defaultConnOptions...)
}

// NewAMQPSession creates a new AMQP session with the internally cached *amqp.Client.
func (ns *Namespace) NewAMQPSession(ctx context.Context) (AMQPSessionCloser, error) {
	client, err := ns.getAMQPClientImpl(ctx)

	if err != nil {
		return nil, err
	}

	return client.NewSession()
}

// NewMgmtClient creates a new management client with the internally cached *amqp.Client.
func (ns *Namespace) NewMgmtClient(ctx context.Context, managementPath string) (MgmtClient, error) {
	return newMgmtClient(ctx, managementPath, ns)
}

// NewRPCLink creates a new amqp-common *rpc.Link with the internally cached *amqp.Client.
func (ns *Namespace) NewRPCLink(ctx context.Context, managementPath string) (*rpc.Link, error) {
	client, err := ns.getAMQPClientImpl(ctx)

	if err != nil {
		return nil, err
	}

	return rpc.NewLink(client, managementPath)
}

// NewAMQPLinks creates an AMQPLinks struct, which groups together the commonly needed links for
// working with Service Bus.
func (ns *Namespace) NewAMQPLinks(entityPath string, createLinkFunc CreateLinkFunc) AMQPLinks {
	return newAMQPLinks(ns, entityPath, createLinkFunc)
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
func (ns *Namespace) Recover(ctx context.Context) (*amqp.Client, error) {
	ns.clientMu.Lock()
	defer ns.clientMu.Unlock()

	if ns.client != nil {
		// the error on close isn't critical
		err := ns.client.Close()
		tab.For(ctx).Error(err)
	}

	var err error
	ns.client, err = ns.newClient(ctx)

	return ns.client, err
}

// negotiateClaim performs initial authentication and starts periodic refresh of credentials.
// the returned func is to cancel() the refresh goroutine.
func (ns *Namespace) NegotiateClaim(ctx context.Context, entityPath string) (func() <-chan struct{}, error) {
	return ns.startNegotiateClaimRenewer(ctx,
		entityPath,
		cbs.NegotiateClaim,
		ns.getAMQPClientImpl)
}

func (ns *Namespace) startNegotiateClaimRenewer(ctx context.Context,
	entityPath string,
	cbsNegotiateClaim func(ctx context.Context, audience string, conn *amqp.Client, provider auth.TokenProvider) error,
	nsGetAMQPClientImpl func(ctx context.Context) (*amqp.Client, error)) (func() <-chan struct{}, error) {
	audience := ns.GetEntityAudience(entityPath)

	refreshClaim := func() (time.Time, error) {
		retrier := ns.baseRetrier.Copy()

		var lastErr error
		var expiration time.Time

		for retrier.Try(ctx) {
			expiration, lastErr = func() (time.Time, error) {
				ctx, span := ns.startSpanFromContext(ctx, "sb.namespace.negotiateClaim")
				defer span.End()

				amqpClient, err := nsGetAMQPClientImpl(ctx)

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

	renewalDuration := time.Until(expiresOn)

	// start the periodic refresh of credentials
	refreshCtx, cancel := context.WithCancel(context.Background())

	go func() {
		for {
			select {
			case <-refreshCtx.Done():
				return
			case <-time.After(renewalDuration):
				expiresOn, err = refreshClaim() // logging will report the error for now

				if err != nil {
					expiresOn = time.Now().Add(time.Minute) // check again in a minute
				}

				renewalDuration = time.Until(expiresOn)
			}
		}
	}()

	cancelRefresh := func() <-chan struct{} {
		cancel()
		return refreshCtx.Done()
	}

	return cancelRefresh, nil
}

func (ns *Namespace) getAMQPClientImpl(ctx context.Context) (*amqp.Client, error) {
	ns.clientMu.Lock()
	defer ns.clientMu.Unlock()

	if ns.client != nil {
		return ns.client, nil
	}

	var err error
	retrier := ns.baseRetrier.Copy()

	for retrier.Try(ctx) {
		ns.client, err = ns.newClient(ctx)

		if err == nil {
			break
		}
	}

	return ns.client, err
}

func (ns *Namespace) getWSSHostURI() string {
	suffix := ns.resolveSuffix()
	if strings.HasSuffix(suffix, "onebox.windows-int.net") {
		return fmt.Sprintf("wss://%s:4446/", ns.GetHostname())
	}
	return fmt.Sprintf("wss://%s/", ns.GetHostname())
}

func (ns *Namespace) getAMQPHostURI() string {
	return fmt.Sprintf("amqps://%s/", ns.GetHostname())
}

func (ns *Namespace) GetHTTPSHostURI() string {
	suffix := ns.resolveSuffix()
	if strings.HasSuffix(suffix, "onebox.windows-int.net") {
		return fmt.Sprintf("https://%s:4446/", ns.GetHostname())
	}
	return fmt.Sprintf("https://%s/", ns.GetHostname())
}

func (ns *Namespace) GetHostname() string {
	return strings.Join([]string{ns.Name, ns.resolveSuffix()}, ".")
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

func (ns *Namespace) resolveSuffix() string {
	if ns.Suffix != "" {
		return ns.Suffix
	}
	return azure.PublicCloud.ServiceBusEndpointSuffix
}
