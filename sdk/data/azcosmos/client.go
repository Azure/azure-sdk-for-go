// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

// ClientOptions configures a [Client].
//
// This does not embed [azcore.ClientOptions]. v2 executes operations through the Cosmos driver
// rather than an azcore HTTP pipeline, so the transport, retry and per-call policy knobs on that
// type have no effect here, and advertising options the client would silently ignore is worse than
// not offering them.
//
// A nil *ClientOptions selects the defaults for every field.
type ClientOptions struct {
	// Runtime shares execution resources and same-account caches, including credentials.
	// Nil uses the process-wide runtime. Supply separate runtimes for identity isolation.
	Runtime *Runtime
	// Operation supplies account-level defaults for item operations and query page fetches.
	Operation OperationOptions
	// Routing decides the order in which the client considers the account's regions. The zero
	// value leaves the order to the account; prefer setting it with [PreferredRegions].
	//
	// [ProximityTo] expands a known application region to the SDK's estimated proximity order.
	Routing RoutingStrategy

	// DiagnosticsHandler observes completed item calls after their lifetime guard is released.
	// It must be safe for concurrent calls. Query-specific diagnostics are not configured here.
	DiagnosticsHandler func(context.Context, OperationDiagnostic)
	// DiagnosticsVerbosity controls how much detail every completion's [Diagnostics.JSON] renders.
	// The zero value, DiagnosticsVerbosityDefault, resolves to DiagnosticsVerbositySummary. It
	// applies to every operation on the client; there is no per-operation override.
	DiagnosticsVerbosity DiagnosticsVerbosity
}

// Client is a client for an Azure Cosmos DB account. It is the entry point to the databases and
// containers in that account.
//
// A Client is safe for concurrent use and is intended to be long lived: it owns the driver
// resources backing it, along with the routing and metadata caches that make requests cheap, so
// creating one per operation is expensive and defeats them. Call [Client.Close] when done.
type Client struct {
	endpoint string
	options  ClientOptions
	runtime  *Runtime
	// mu guards the client's lifetime rather than its fields. Operations hold it for read while
	// they run, so Close taking it for write is exactly "wait for in-flight operations to
	// finish". That matters more here than it would in a pure-Go client: closing releases handles
	// owned by the driver, and an operation still running would be using freed memory.
	mu        sync.RWMutex
	closed    bool
	closeOnce sync.Once
	closeErr  error

	// beforeItemAcquire is a test hook for observing the boundary before an item operation enters
	// the client's lifetime guard.
	beforeItemAcquire func()

	// driver holds the resources the client owns in the driver. It is nil in builds that are not
	// bound to the driver, where operations report that the driver is unavailable.
	driver *nativeDriver

	// usesKeyAuth records whether the client was constructed with NewClientWithKey. Azure Cosmos
	// DB only accepts master-key authentication for database and container management operations
	// issued through the data-plane SQL API, so management methods reject a token-credential
	// client rather than letting it fail server-side with a less specific error.
	usesKeyAuth bool
}

// NewClient creates a client that authenticates with Microsoft Entra ID.
//
// endpoint is the Cosmos DB account endpoint. cred is any [azcore.TokenCredential], such as the
// implementations in the azidentity module. options may be nil to accept the defaults. Construction
// initializes account metadata and routing before returning success.
// Failure returns a nil client. Context cancellation bounds the wait; native cleanup may finish later.
//
// Clients sharing a runtime and account also share native authentication state; cached resources
// may use a credential supplied by another client, even after that client closes. This reuse
// reduces duplicate token requests and authentication-service load. A nil Runtime uses the
// process-wide runtime; supply separate runtimes when credential identities must be isolated.
func NewClient(ctx context.Context, endpoint string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if cred == nil {
		return nil, errors.New("azcosmos: credential must not be nil")
	}
	return newInitializedClient(ctx, endpoint, "", cred, options)
}

// NewClientWithKey creates a client that authenticates with an account key.
//
// Prefer [NewClient] where possible; account keys grant full access to the account and cannot be
// scoped down. options may be nil to accept the defaults. Construction initializes account
// metadata and routing before returning success.
// Failure returns a nil client. Context cancellation bounds the wait; native cleanup may finish later.
func NewClientWithKey(ctx context.Context, endpoint string, cred KeyCredential, options *ClientOptions) (*Client, error) {
	if cred.accountKey == "" {
		return nil, errors.New("azcosmos: credential must be created with NewKeyCredential")
	}
	return newInitializedClient(ctx, endpoint, cred.accountKey, nil, options)
}

func newInitializedClient(ctx context.Context, endpoint, accountKey string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if ctx == nil {
		return nil, errors.New("azcosmos: context must not be nil")
	}
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	client, err := newClient(endpoint, accountKey, cred, options)
	if err != nil {
		return nil, err
	}
	if err := client.initialize(ctx); err != nil {
		// Cancellation can leave native submissions in flight. Teardown drains their late
		// completions before releasing handles and detaching the runtime, without delaying ctx.
		go func() {
			_ = client.close() // Native and stub teardown currently return nil.
		}()
		return nil, err
	}
	return client, nil
}

// newClient validates the inputs shared by every constructor and opens the client's driver
// resources.
//
// accountKey is empty when the caller supplied a token credential. In builds that are not bound to
// the driver no resources are acquired, and operations report that the driver is unavailable.
func newClient(
	endpoint string,
	accountKey string,
	tokenCredential azcore.TokenCredential,
	options *ClientOptions,
) (*Client, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("azcosmos: parsing endpoint: %w", err)
	}
	// Go's url.Parse is permissive where the Rust SDK's URL parser is not: it accepts a relative
	// path and an endpoint whose authority is only a port, so both are rejected here to reach the
	// same effective validation. Hostname() rather than Host, because Host includes the port.
	if !parsed.IsAbs() || parsed.Hostname() == "" {
		return nil, fmt.Errorf("azcosmos: endpoint %q must be an absolute URL, for example https://myaccount.documents.azure.com", endpoint)
	}
	// The driver rejects http:// for anything that is not an emulator host, when the driver is
	// created and so before any credential is sent, and documents itself as the single source of
	// truth for that rule. Emulator detection is deliberately not duplicated here; this only
	// catches an obviously wrong scheme early.
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return nil, fmt.Errorf("azcosmos: endpoint %q must use the http or https scheme", endpoint)
	}

	client := &Client{endpoint: endpoint, usesKeyAuth: accountKey != ""}
	if options != nil {
		client.options = *options
		client.options.Routing = options.Routing.clone()
		client.options.Operation = options.Operation.clone()
	}
	if err := client.options.validate(); err != nil {
		return nil, err
	}
	client.runtime = client.options.Runtime
	if client.runtime == nil && driverAvailable {
		client.runtime, err = defaultRuntime()
		if err != nil {
			return nil, err
		}
	}
	// Registration and local handle construction are one transaction against runtime shutdown.
	if client.runtime != nil {
		client.runtime.mu.Lock()
		if client.runtime.closed || client.runtime.native == nil {
			client.runtime.mu.Unlock()
			return nil, &Error{Code: CodeClientClosed, Message: "the runtime is closed or uninitialized"}
		}
	}

	driver, err := openDriver(driverConfig{
		endpoint:        endpoint,
		accountKey:      accountKey,
		tokenCredential: tokenCredential,
		options:         client.options,
		runtime:         client.runtime,
	})
	client.driver = driver
	if client.runtime != nil {
		if err == nil {
			client.runtime.clients++
		}
		client.runtime.mu.Unlock()
	}
	if err != nil {
		return nil, err
	}
	return client, nil
}

// initialize is also used by offline tests to exercise native initialization ownership.
func (c *Client) initialize(ctx context.Context) error {
	if ctx == nil {
		return errors.New("azcosmos: context must not be nil")
	}
	release, err := c.acquire()
	if err != nil {
		return err
	}
	defer release()
	if err := ctx.Err(); err != nil {
		return err
	}
	return c.driver.initialize(ctx)
}

// Endpoint returns the Cosmos DB account endpoint the client was created with.
func (c *Client) Endpoint() string {
	return c.endpoint
}

// Close releases this client's driver resources after all submitted native work finishes.
// Cached credentials remain alive with the runtime and are not cancelled by closing one client.
// Each token acquisition has a 60-second timeout; custom credentials must honor context cancellation.
// Afterwards every operation on the client
// fails with [CodeClientClosed] rather than reaching the driver.
//
// Close is idempotent and safe to call concurrently; every caller observes the same result. It
// returns an error only when the client could not be torn down cleanly, in which case the
// resources are released anyway, so there is nothing to retry.
func (c *Client) Close() error {
	return c.close()
}

func (c *Client) close() error {
	c.closeOnce.Do(func() {
		// Taking the write lock blocks until every operation holding it for read has finished,
		// and keeps later operations out once closed is set.
		c.mu.Lock()
		defer c.mu.Unlock()
		c.closed = true
		// The error is stored on the client rather than in a local so that the second and
		// subsequent callers see it too.
		c.closeErr = c.driver.close()
		c.driver = nil
		if c.runtime != nil {
			c.runtime.detach()
		}
	})
	return c.closeErr
}

// acquire registers an operation as in flight, or reports that the client has been closed. The
// returned function must be called when the operation finishes.
func (c *Client) acquire() (release func(), err error) {
	c.mu.RLock()
	if c.closed {
		c.mu.RUnlock()
		return nil, &Error{Code: CodeClientClosed, Message: "the client has been closed"}
	}
	return c.mu.RUnlock, nil
}

// NewDatabase returns a client for a database in the account. It does not contact the service, so
// it succeeds whether or not the database exists.
func (c *Client) NewDatabase(id string) (*DatabaseClient, error) {
	if id == "" {
		return nil, errors.New("azcosmos: database id must not be empty")
	}
	return &DatabaseClient{id: id, client: c}, nil
}

// CreateDatabase creates a new database.
//
// properties.ID is required; every other field is server-assigned and ignored. options may be
// nil. The operation always requests response content so the returned properties contain the
// server-assigned fields; CreateDatabaseOptions.Operation.EnableContentResponseOnWrite is ignored.
//
// Azure Cosmos DB only accepts master-key authentication for this operation: a client created
// with [NewClient] (Microsoft Entra ID) returns an error rather than reaching the service. Use
// [NewClientWithKey].
func (c *Client) CreateDatabase(ctx context.Context, properties DatabaseProperties, options *CreateDatabaseOptions) (DatabaseResponse, error) {
	return c.createDatabase(ctx, properties, options)
}

// NewContainer resolves existing container metadata through the native cache before returning.
// Cache misses contact the service; this method does not create a container.
func (c *Client) NewContainer(ctx context.Context, databaseID string, containerID string) (*ContainerClient, error) {
	database, err := c.NewDatabase(databaseID)
	if err != nil {
		return nil, err
	}
	return database.NewContainer(ctx, containerID)
}

// validate reports option values that cannot be passed through the C ABI. Values the driver
// understands are passed through and validated there, so this package does not duplicate its rules.
func (o ClientOptions) validate() error {
	switch o.DiagnosticsVerbosity {
	case DiagnosticsVerbosityDefault, DiagnosticsVerbositySummary, DiagnosticsVerbosityDetailed:
	default:
		return fmt.Errorf("azcosmos: ClientOptions.DiagnosticsVerbosity: unrecognized %s", o.DiagnosticsVerbosity)
	}
	return o.Operation.validate()
}
