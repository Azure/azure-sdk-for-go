// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azcosmos

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/cloud"
)

// ClientOptions configures a [Client].
//
// Unlike v1 this does not embed [azcore.ClientOptions]. v2 executes operations through the Cosmos
// driver rather than an azcore HTTP pipeline, so the transport, retry and per-call policy knobs on
// that type have no effect here. Advertising options the client would silently ignore is worse
// than not offering them, so the fields below are only the ones the driver honors.
//
// A nil *ClientOptions selects the defaults for every field.
type ClientOptions struct {
	// Cloud specifies which Azure cloud the account belongs to, which determines the audience
	// used for Microsoft Entra ID authentication. Defaults to [cloud.AzurePublic].
	Cloud cloud.Configuration

	// ApplicationID is an application-specific identifier appended to the user agent sent with
	// every request. Keep it short and free of personally identifiable information.
	ApplicationID string

	// PreferredRegions orders the regions the client routes to, most preferred first. When empty,
	// the client uses the account's own region order.
	PreferredRegions []string

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
	endpoint  string
	options   ClientOptions
	closeOnce sync.Once
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
// scoped down. options may be nil to accept the defaults.
func NewClientWithKey(endpoint string, cred KeyCredential, options *ClientOptions) (*Client, error) {
	if cred.accountKey == "" {
		return nil, errors.New("azcosmos: credential must be created with NewKeyCredential")
	}
	return newClient(endpoint, options)
}

// NewClientFromConnectionString creates a client from a Cosmos DB connection string, which carries
// both the account endpoint and an account key.
//
// Prefer [NewClient] where possible; account keys grant full access to the account and cannot be
// scoped down. options may be nil to accept the defaults.
func NewClientFromConnectionString(connectionString string, options *ClientOptions) (*Client, error) {
	endpoint, cred, err := parseConnectionString(connectionString)
	if err != nil {
		return nil, err
	}
	return NewClientWithKey(endpoint, cred, options)
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

	client := &Client{endpoint: endpoint}
	if options != nil {
		client.options = *options
		client.options.PreferredRegions = append([]string(nil), options.PreferredRegions...)
	}
	return client, nil
}

// parseConnectionString splits a connection string into its endpoint and account key. The format
// is a semicolon-separated list of key=value pairs; only AccountEndpoint and AccountKey are
// meaningful here, and any other pairs are ignored so that a string copied from the portal works
// as-is.
func parseConnectionString(connectionString string) (string, KeyCredential, error) {
	const (
		endpointKey = "AccountEndpoint"
		accountKey  = "AccountKey"
	)

	var endpoint, key string
	for _, pair := range strings.Split(connectionString, ";") {
		if pair == "" {
			continue
		}
		name, value, found := strings.Cut(pair, "=")
		if !found {
			return "", KeyCredential{}, errors.New("azcosmos: connection string must be a ';'-separated list of 'key=value' pairs")
		}
		switch {
		case strings.EqualFold(name, endpointKey):
			endpoint = value
		case strings.EqualFold(name, accountKey):
			key = value
		}
	}

	if endpoint == "" {
		return "", KeyCredential{}, fmt.Errorf("azcosmos: connection string is missing %q", endpointKey)
	}
	cred, err := NewKeyCredential(key)
	if err != nil {
		return "", KeyCredential{}, fmt.Errorf("azcosmos: connection string is missing %q", accountKey)
	}
	return endpoint, cred, nil
}

// Endpoint returns the Cosmos DB account endpoint the client was created with.
func (c *Client) Endpoint() string {
	return c.endpoint
}

// Close releases the driver resources the client owns and waits for its in-flight operations to
// finish. The client must not be used afterwards.
//
// Close is idempotent. It returns an error only when the client could not be torn down cleanly;
// the resources are released either way, so there is nothing to retry.
func (c *Client) Close() error {
	var err error
	c.closeOnce.Do(func() {
		// Nothing to release until the driver binding lands. The error return is part of the
		// contract from the start because teardown has to drain in-flight operations, and adding
		// it later would break every caller.
		err = nil
	})
	return err
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
