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
// not offering them. Every field below is one the driver honors.
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

	// BinaryEncoding sets the SDK's client encoding default for create/read/replace/upsert.
	// Nil resolves AZURE_COSMOS_BINARY_ENCODING_ENABLED at construction, then defaults to enabled.
	// Explicit request encoding overrides this value. PATCH and delete use Operation inheritance.
	BinaryEncoding *BinaryEncodingOptions
	// FaultInjectionRules configures client-specific native fault injection for testing.
	FaultInjectionRules []FaultInjectionRule
	// DiagnosticsHandler observes completed item calls after their lifetime guard is released.
	// It must be safe for concurrent calls. Query-specific diagnostics are not configured here.
	DiagnosticsHandler func(context.Context, OperationDiagnostic)
}

// Client is a client for an Azure Cosmos DB account. It is the entry point to the databases and
// containers in that account.
//
// A Client is safe for concurrent use and is intended to be long lived: it owns the driver
// resources backing it, along with the routing and metadata caches that make requests cheap, so
// creating one per operation is expensive and defeats them. Call [Client.Close] when done.
type Client struct {
	endpoint       string
	options        ClientOptions
	runtime        *Runtime
	binaryEncoding BinaryEncodingOptions
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
}

// NewClient creates a client that authenticates with Microsoft Entra ID.
//
// endpoint is the Cosmos DB account endpoint. cred is any [azcore.TokenCredential], such as the
// implementations in the azidentity module. options may be nil to accept the defaults. Construction
// performs no network I/O; call [Client.Initialize] to fill account-level caches eagerly.
func NewClient(endpoint string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if cred == nil {
		return nil, errors.New("azcosmos: credential must not be nil")
	}
	return newClient(endpoint, "", cred, options)
}

// NewClientWithKey creates a client that authenticates with an account key.
//
// Prefer [NewClient] where possible; account keys grant full access to the account and cannot be
// scoped down. options may be nil to accept the defaults. Construction performs no network I/O;
// call [Client.Initialize] to fill account-level caches eagerly.
func NewClientWithKey(endpoint string, cred KeyCredential, options *ClientOptions) (*Client, error) {
	if cred.accountKey == "" {
		return nil, errors.New("azcosmos: credential must be created with NewKeyCredential")
	}
	return newClient(endpoint, cred.accountKey, nil, options)
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

	client := &Client{endpoint: endpoint}
	if options != nil {
		client.options = *options
		client.options.Routing = options.Routing.clone()
		client.options.Operation = options.Operation.clone()
		client.options.BinaryEncoding = options.BinaryEncoding.clone()
		client.options.FaultInjectionRules = cloneFaultInjectionRules(options.FaultInjectionRules)
	}
	if err := client.options.validate(); err != nil {
		return nil, err
	}
	client.binaryEncoding = resolveClientBinaryEncoding(client.options.BinaryEncoding)
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

// Initialize eagerly creates the driver and fills its account-properties and routing caches.
//
// Initialize is idempotent and safe for concurrent use. Operations also initialize lazily, so
// callers only need this when they want readiness and initialization failures before the first
// operation.
func (c *Client) Initialize(ctx context.Context) error {
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

// NewContainer returns a client for a container in the account. It does not contact the service,
// so it succeeds whether or not the container exists.
func (c *Client) NewContainer(databaseID string, containerID string) (*ContainerClient, error) {
	database, err := c.NewDatabase(databaseID)
	if err != nil {
		return nil, err
	}
	return database.NewContainer(containerID)
}

// validate reports option values that cannot be passed through the C ABI. Values the driver
// understands are passed through and validated there, so this package does not duplicate its rules.
func (o ClientOptions) validate() error {
	if err := validateFaultInjectionRules(o.FaultInjectionRules); err != nil {
		return err
	}
	return o.Operation.validate()
}
