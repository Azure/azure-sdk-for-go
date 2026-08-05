// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
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
// type have no effect here. Advertising options the client would silently ignore is worse than not
// offering them, so the fields below are only the ones the driver honors.
//
// A nil *ClientOptions selects the defaults for every field.
type ClientOptions struct {
	// Routing decides the order in which the client considers the account's regions. The zero
	// value leaves the order to the account; prefer setting it with [ProximityTo] or
	// [PreferredRegions].
	Routing RoutingStrategy

	// ApplicationID is an application-specific identifier appended to the user agent sent with
	// every request. Keep it short and free of personally identifiable information.
	ApplicationID string

	// EnableContentResponseOnWrite requests that writes return the resulting item. Leaving it
	// false reduces network and CPU cost, because the service does not send the item back and the
	// client does not deserialize it. Can be overridden per operation.
	EnableContentResponseOnWrite bool
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

	// mu guards the client's lifetime rather than its fields. Operations hold it for read while
	// they run, so Close taking it for write is exactly "wait for in-flight operations to
	// finish". That matters more here than it would in a pure-Go client: closing releases handles
	// owned by the driver, and an operation still running would be using freed memory.
	mu     sync.RWMutex
	closed bool

	closeOnce sync.Once
	closeErr  error
}

// NewClient creates a client that authenticates with Microsoft Entra ID.
//
// endpoint is the Cosmos DB account endpoint. cred is any [azcore.TokenCredential], such as the
// implementations in the azidentity module. options may be nil to accept the defaults.
func NewClient(endpoint string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if cred == nil {
		return nil, errors.New("azcosmos: credential must not be nil")
	}
	return newClient(endpoint, options)
}

// NewClientWithKey creates a client that authenticates with an account key.
//
// Prefer [NewClient] where possible; account keys grant full access to the account and cannot be
// scoped down. Rotate the key in place with [azcore.KeyCredential.Update]. options may be nil to
// accept the defaults.
func NewClientWithKey(endpoint string, cred *azcore.KeyCredential, options *ClientOptions) (*Client, error) {
	if cred == nil {
		return nil, errors.New("azcosmos: credential must not be nil")
	}
	return newClient(endpoint, options)
}

// newClient validates the inputs shared by every constructor and returns a client handle.
//
// The handle is not yet backed by the driver: operations on it report that they are not
// implemented. Construction still succeeds so that the surface can be compiled and explored
// against during the preview.
func newClient(endpoint string, options *ClientOptions) (*Client, error) {
	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("azcosmos: parsing endpoint: %w", err)
	}
	if !parsed.IsAbs() || parsed.Host == "" {
		return nil, fmt.Errorf("azcosmos: endpoint %q must be an absolute URL, for example https://myaccount.documents.azure.com", endpoint)
	}
	if parsed.Scheme != "https" && parsed.Scheme != "http" {
		return nil, fmt.Errorf("azcosmos: endpoint %q must use the http or https scheme", endpoint)
	}

	client := &Client{endpoint: endpoint}
	if options != nil {
		client.options = *options
		client.options.Routing = options.Routing.clone()
	}
	return client, nil
}

// Endpoint returns the Cosmos DB account endpoint the client was created with.
func (c *Client) Endpoint() string {
	return c.endpoint
}

// Close releases the driver resources the client owns. It first waits for the client's in-flight
// operations to finish, and afterwards every operation on the client fails with [CodeClientClosed]
// rather than reaching the driver.
//
// Close is idempotent and safe to call concurrently; every caller observes the same result. It
// returns an error only when the client could not be torn down cleanly, in which case the
// resources are released anyway, so there is nothing to retry.
func (c *Client) Close() error {
	c.closeOnce.Do(func() {
		// Taking the write lock blocks until every operation holding it for read has finished,
		// and keeps later operations out once closed is set.
		c.mu.Lock()
		defer c.mu.Unlock()
		c.closed = true
		// Driver resources are released here once the binding lands, recording a teardown that
		// did not complete cleanly in c.closeErr. The error is stored on the client rather than
		// in a local so that the second and subsequent callers see it too.
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
