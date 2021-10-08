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

func (c *CreateKeyOptions) toKeyCreateParameters(keyType JSONWebKeyType) internal.KeyCreateParameters {
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
			Key:        jsonWebKeyToGenerated(g.Key),
			Tags:       g.Tags,
			Managed:    g.Managed,
		},
	}
}

func (c *Client) CreateKey(ctx context.Context, name string, keyType JSONWebKeyType, options *CreateKeyOptions) (CreateKeyResponse, error) {
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

func (c *CreateECKeyOptions) toKeyCreateParameters(keyType JSONWebKeyType) internal.KeyCreateParameters {
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
			Key:        jsonWebKeyToGenerated(g.Key),
			Tags:       g.Tags,
			Managed:    g.Managed,
		},
	}
}

func (c *Client) CreateECKey(ctx context.Context, name string, options *CreateECKeyOptions) (CreateECKeyResponse, error) {
	keyType := JSONWebKeyTypeEC

	if options != nil && options.HardwareProtected {
		keyType = JSONWebKeyTypeECHSM
	} else if options == nil {
		options = &CreateECKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &internal.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return CreateECKeyResponse{}, nil
	}

	return createECKeyResponseFromGenerated(resp), nil
}
