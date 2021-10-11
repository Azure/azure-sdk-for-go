//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal"
)

// Client is the struct for interacting with a KeyVault Keys instance
type Client struct {
	kvClient *internal.KeyVaultClient
	vaultUrl string
}

// ClientOptions are the configurable options on a Client.
type ClientOptions struct {
	// Transport sets the transport for making HTTP requests.
	Transport policy.Transporter

	// Retry configures the built-in retry policy behavior.
	Retry policy.RetryOptions

	// Telemetry configures the built-in telemetry policy behavior.
	Telemetry policy.TelemetryOptions

	// Logging configures the built-in logging policy behavior.
	Logging policy.LogOptions

	// PerCallPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request.
	PerCallPolicies []policy.Policy

	// PerRetryPolicies contains custom policies to inject into the pipeline.
	// Each policy is executed once per request, and for each retry request.
	PerTryPolicies []policy.Policy
}

func (c *ClientOptions) toConnectionOptions() *internal.ConnectionOptions {
	if c == nil {
		return nil
	}

	return &internal.ConnectionOptions{
		HTTPClient:       c.Transport,
		Retry:            c.Retry,
		Telemetry:        c.Telemetry,
		Logging:          c.Logging,
		PerCallPolicies:  c.PerCallPolicies,
		PerRetryPolicies: c.PerTryPolicies,
	}
}

// NewClient returns a pointer to a Client object affinitized to a vaultUrl.
func NewClient(vaultUrl string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	conn := internal.NewConnection(credential, options.toConnectionOptions())

	return &Client{
		kvClient: &internal.KeyVaultClient{
			Con: conn,
		},
		vaultUrl: vaultUrl,
	}, nil
}

// CreateKeyOptions contains the optional parameters for the KeyVaultClient.CreateKey method.
type CreateKeyOptions struct {
	// Elliptic curve name. For valid values, see JsonWebKeyCurveName.
	Curve *internal.JSONWebKeyCurveName `json:"crv,omitempty"`

	// The attributes of a key managed by the key vault service.
	KeyAttributes *internal.KeyAttributes         `json:"attributes,omitempty"`
	KeyOps        []*internal.JSONWebKeyOperation `json:"key_ops,omitempty"`

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`

	// The public exponent for a RSA key.
	PublicExponent *int32 `json:"public_exponent,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

func (c *CreateKeyOptions) toGenerated() *internal.KeyVaultClientCreateKeyOptions {
	return &internal.KeyVaultClientCreateKeyOptions{}
}

func (c *CreateKeyOptions) toKeyCreateParameters(keyType KeyType) internal.KeyCreateParameters {
	return internal.KeyCreateParameters{
		Kty:            keyType.toGenerated(),
		Curve:          c.Curve,
		KeyAttributes:  c.KeyAttributes,
		KeyOps:         c.KeyOps,
		KeySize:        c.KeySize,
		PublicExponent: c.PublicExponent,
		Tags:           c.Tags,
	}
}

// KeyVaultClientCreateKeyResponse contains the response from method KeyVaultClient.CreateKey.
type CreateKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func createKeyResponseFromGenerated(g internal.KeyVaultClientCreateKeyResponse) CreateKeyResponse {
	return CreateKeyResponse{
		RawResponse: g.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(g.Attributes),
			Key:        jsonWebKeyFromGenerated(g.Key),
			Tags:       g.Tags,
			Managed:    g.Managed,
		},
	}
}

func (c *Client) CreateKey(ctx context.Context, name string, keyType KeyType, options *CreateKeyOptions) (CreateKeyResponse, error) {
	if options == nil {
		options = &CreateKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), options.toGenerated())
	if err != nil {
		return CreateKeyResponse{}, err
	}

	return createKeyResponseFromGenerated(resp), nil
}

type CreateECKeyOptions struct {
	// Elliptic curve name. For valid values, see JsonWebKeyCurveName.
	Curve *internal.JSONWebKeyCurveName `json:"crv,omitempty"`

	// The attributes of a key managed by the key vault service.
	KeyAttributes *internal.KeyAttributes         `json:"attributes,omitempty"`
	KeyOps        []*internal.JSONWebKeyOperation `json:"key_ops,omitempty"`

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`

	// The public exponent for a RSA key.
	PublicExponent *int32 `json:"public_exponent,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`

	// Whether to create an EC key with HSM protection
	HardwareProtected bool
}

func (c *CreateECKeyOptions) toKeyCreateParameters(keyType KeyType) internal.KeyCreateParameters {
	return internal.KeyCreateParameters{
		Kty:            keyType.toGenerated(),
		Curve:          c.Curve,
		KeyAttributes:  c.KeyAttributes,
		KeyOps:         c.KeyOps,
		KeySize:        c.KeySize,
		PublicExponent: c.PublicExponent,
		Tags:           c.Tags,
	}
}

type CreateECKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func createECKeyResponseFromGenerated(g internal.KeyVaultClientCreateKeyResponse) CreateECKeyResponse {
	return CreateECKeyResponse{
		RawResponse: g.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(g.Attributes),
			Key:        jsonWebKeyFromGenerated(g.Key),
			Tags:       g.Tags,
			Managed:    g.Managed,
		},
	}
}

func (c *Client) CreateECKey(ctx context.Context, name string, options *CreateECKeyOptions) (CreateECKeyResponse, error) {
	keyType := EC

	if options != nil && options.HardwareProtected {
		keyType = ECHSM
	} else if options == nil {
		options = &CreateECKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &internal.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return CreateECKeyResponse{}, nil
	}

	return createECKeyResponseFromGenerated(resp), nil
}

type CreateOCTKeyOptions struct {
	// Hardware Protected OCT Key
	HardwareProtected bool

	// Elliptic curve name. For valid values, see JsonWebKeyCurveName.
	Curve *internal.JSONWebKeyCurveName `json:"crv,omitempty"`

	// The attributes of a key managed by the key vault service.
	KeyAttributes *internal.KeyAttributes         `json:"attributes,omitempty"`
	KeyOps        []*internal.JSONWebKeyOperation `json:"key_ops,omitempty"`

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`

	// The public exponent for a RSA key.
	PublicExponent *int32 `json:"public_exponent,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]*string `json:"tags,omitempty"`
}

func (c *CreateOCTKeyOptions) toKeyCreateParameters(keyType KeyType) internal.KeyCreateParameters {
	return internal.KeyCreateParameters{
		Kty:            keyType.toGenerated(),
		Curve:          c.Curve,
		KeyAttributes:  c.KeyAttributes,
		KeyOps:         c.KeyOps,
		KeySize:        c.KeySize,
		PublicExponent: c.PublicExponent,
		Tags:           c.Tags,
	}
}

type CreateOCTKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func createOCTKeyResponseFromGenerated(i internal.KeyVaultClientCreateKeyResponse) CreateOCTKeyResponse {
	return CreateOCTKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       i.Tags,
			Managed:    i.Managed,
		},
	}
}

func (c *Client) CreateOCTKey(ctx context.Context, name string, options *CreateOCTKeyOptions) (CreateOCTKeyResponse, error) {
	keyType := Oct

	if options != nil && options.HardwareProtected {
		keyType = OctHSM
	} else if options == nil {
		options = &CreateOCTKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &internal.KeyVaultClientCreateKeyOptions{})

	return createOCTKeyResponseFromGenerated(resp), err
}

type ListKeysPager interface {
	// PageResponse returns the current ListSecretsPage
	PageResponse() ListKeysPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

type listKeysPager struct {
	genPager internal.KeyVaultClientGetKeysPager
	err      error
}

func (l *listKeysPager) PageResponse() ListKeysPage {
	return listSecretsPageFromGenerated(l.genPager.PageResponse())
}

func (l *listKeysPager) Err() error {
	return l.err
}

func (l *listKeysPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

type ListKeysOptions struct {
	MaxResults *int32
}

func (l ListKeysOptions) toGenerated() *internal.KeyVaultClientGetKeysOptions {
	return &internal.KeyVaultClientGetKeysOptions{Maxresults: l.MaxResults}
}

type ListKeysPage struct {
	// READ-ONLY; The URL to get the next set of keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of keys in the key vault along with a link to the next page of keys.
	Keys []*KeyItem `json:"value,omitempty" azure:"ro"`
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func listSecretsPageFromGenerated(i internal.KeyVaultClientGetKeysResponse) ListKeysPage {
	var keys []*KeyItem
	for _, k := range i.Value {
		keys = append(keys, keyItemFromGenerated(k))
	}
	return ListKeysPage{
		RawResponse: i.RawResponse,
		NextLink:    i.NextLink,
		Keys:        keys,
	}
}

func (c *Client) ListKeys(options *ListKeysOptions) ListKeysPager {
	if options == nil {
		options = &ListKeysOptions{}
	}
	p := c.kvClient.GetKeys(c.vaultUrl, options.toGenerated())

	return &listKeysPager{
		genPager: *p,
	}
}

type GetKeyOptions struct {
	Version string
}

type GetKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func getKeyResponseFromGenerated(i internal.KeyVaultClientGetKeyResponse) GetKeyResponse {
	return GetKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       i.Tags,
			Managed:    i.Managed,
		},
	}
}

func (c *Client) GetKey(ctx context.Context, keyName string, options *GetKeyOptions) (GetKeyResponse, error) {
	if options == nil {
		options = &GetKeyOptions{}
	}

	resp, err := c.kvClient.GetKey(ctx, c.vaultUrl, keyName, options.Version, &internal.KeyVaultClientGetKeyOptions{})
	if err != nil {
		return GetKeyResponse{}, err
	}

	return getKeyResponseFromGenerated(resp), err
}

/*
type DeleteKeyPoller interface {
	// Done returns true if the LRO has reached a terminal state
	Done() bool

	// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
	// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
	// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
	Poll(context.Context) (*http.Response, error)

	// FinalResponse returns the final response after the operations has finished
	FinalResponse(context.Context) (DeleteKeyResponse, error)
}

type DeleteKeyPollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (DeleteKeyResponse, error)

	// Poller contains an initialized WidgetPoller
	Poller DeleteKeyPoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

type DeleteKeyResponse struct{}

type BeginDeleteKeyOptions struct{}

func (b *BeginDeleteKeyOptions) toGenerated() *internal.KeyVaultClientDeleteKeyOptions {
	return &internal.KeyVaultClientDeleteKeyOptions{}
}

func (c *Client) BeginDeleteKey(ctx context.Context, keyName string, options *BeginDeleteKeyOptions) (DeleteKeyPollerResponse, error) {
	if options == nil {
		options = &BeginDeleteKeyOptions{}
	}

	resp, err := c.kvClient.DeleteKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return DeleteKeyPollerResponse{}, err
	}
	_ = resp
}
*/
