//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azcrypto

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal"
)

// The Client performs cryptographic operations using Azure Key Vault Keys. This client
// will perform operations locally when it's initialized with the necessary key material or
// is able to get that material from Key Vault. When the required key material is unavailable,
// cryptographic operations are performed by the Key Vault service.
type Client struct {
	kvClient   *internal.KeyVaultClient
	vaultURL   string
	keyID      string
	keyVersion string
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

// converts ClientOptions to generated *internal.ConnectionOptions
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

func parseKeyIDAndVersion(id string) (string, string, error) {
	parsed, err := url.Parse(id)
	if err != nil {
		return "", "", err
	}

	path := strings.Split(parsed.Path, "/")

	if len(path) != 4 {
		return "", "", fmt.Errorf("could not parse Key ID from %s", id)
	}

	return path[2], path[3], nil
}

func parseVaultURL(base string) (string, error) {
	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s/", parsed.Scheme, parsed.Host), nil
}

func NewClient(key string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	vaultURL, err := parseVaultURL(key)
	if err != nil {
		return &Client{}, err
	}

	keyID, keyVersion, err := parseKeyIDAndVersion(key)
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		kvClient: &internal.KeyVaultClient{
			Con: internal.NewConnection(cred, options.toConnectionOptions()),
		},
		vaultURL:   vaultURL,
		keyID:      keyID,
		keyVersion: keyVersion,
	}, nil
}

// Creates a new Client that can only perform cryptographic operations locally
func NewClientFromJSONWebKey(key azkeys.JSONWebKey, options *ClientOptions) (*Client, error) {
	return &Client{}, nil
}

// Optional parameters for the azcrypto.Client.EncryptOptions method
type EncryptOptions struct {
	// Additional data to authenticate but not encrypt/decrypt when using authenticated crypto algorithms.
	AAD []byte `json:"aad,omitempty"`

	// Initialization vector for symmetric algorithms.
	IV []byte `json:"iv,omitempty"`

	// The tag to authenticate when performing decryption with an authenticated algorithm.
	Tag []byte `json:"tag,omitempty"`
}

func (e EncryptOptions) toGeneratedKeyOperationsParameters(alg EncryptionAlgorithm, value []byte) internal.KeyOperationsParameters {
	return internal.KeyOperationsParameters{
		Algorithm: (*internal.JSONWebKeyEncryptionAlgorithm)(&alg),
		Value:     value,
		AAD:       e.AAD,
		Iv:        e.IV,
		Tag:       e.Tag,
	}
}

type EncryptResponse struct {
	KeyOperationResult

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func encryptResponseFromGenerated(i internal.KeyVaultClientEncryptResponse) EncryptResponse {
	return EncryptResponse{
		RawResponse: i.RawResponse,
		KeyOperationResult: KeyOperationResult{
			AdditionalAuthenticatedData: i.AdditionalAuthenticatedData,
			AuthenticationTag:           i.AuthenticationTag,
			IV:                          i.Iv,
			KeyID:                       i.Kid,
			Result:                      i.Result,
		},
	}
}

func (c *Client) Encrypt(ctx context.Context, alg EncryptionAlgorithm, value []byte, options *EncryptOptions) (EncryptResponse, error) {
	if options == nil {
		options = &EncryptOptions{}
	}

	resp, err := c.kvClient.Encrypt(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		options.toGeneratedKeyOperationsParameters(alg, value),
		&internal.KeyVaultClientEncryptOptions{},
	)
	if err != nil {
		return EncryptResponse{}, err
	}

	return encryptResponseFromGenerated(resp), nil
}

// DecryptOptions contains the optional parameters for the Client.Decrypt function
type DecryptOptions struct {
	// Additional data to authenticate but not encrypt/decrypt when using authenticated crypto algorithms.
	AAD []byte `json:"aad,omitempty"`

	// Initialization vector for symmetric algorithms.
	IV []byte `json:"iv,omitempty"`

	// The tag to authenticate when performing decryption with an authenticated algorithm.
	Tag []byte `json:"tag,omitempty"`
}

func (e DecryptOptions) toGeneratedKeyOperationsParameters(alg EncryptionAlgorithm, value []byte) internal.KeyOperationsParameters {
	return internal.KeyOperationsParameters{
		Algorithm: (*internal.JSONWebKeyEncryptionAlgorithm)(&alg),
		Value:     value,
		AAD:       e.AAD,
		Iv:        e.IV,
		Tag:       e.Tag,
	}
}

type DecryptResponse struct {
	KeyOperationResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func decryptResponseFromGenerated(i internal.KeyVaultClientDecryptResponse) DecryptResponse {
	return DecryptResponse{
		RawResponse: i.RawResponse,
		KeyOperationResult: KeyOperationResult{
			AdditionalAuthenticatedData: i.AdditionalAuthenticatedData,
			AuthenticationTag:           i.AuthenticationTag,
			IV:                          i.Iv,
			KeyID:                       i.Kid,
			Result:                      i.Result,
		},
	}
}

func (c *Client) Decrypt(ctx context.Context, alg EncryptionAlgorithm, ciphertext []byte, options *DecryptOptions) (DecryptResponse, error) {
	if options == nil {
		options = &DecryptOptions{}
	}

	resp, err := c.kvClient.Decrypt(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		options.toGeneratedKeyOperationsParameters(alg, ciphertext),
		&internal.KeyVaultClientDecryptOptions{},
	)

	if err != nil {
		return DecryptResponse{}, err
	}

	return decryptResponseFromGenerated(resp), nil
}

type WrapKeyOptions struct {
	// Additional data to authenticate but not encrypt/decrypt when using authenticated crypto algorithms.
	AAD []byte `json:"aad,omitempty"`

	// Initialization vector for symmetric algorithms.
	IV []byte `json:"iv,omitempty"`

	// The tag to authenticate when performing decryption with an authenticated algorithm.
	Tag []byte `json:"tag,omitempty"`
}

func (w WrapKeyOptions) toGeneratedKeyOperationsParameters(alg KeyWrapAlgorithm, value []byte) internal.KeyOperationsParameters {
	return internal.KeyOperationsParameters{
		Algorithm: (*internal.JSONWebKeyEncryptionAlgorithm)(&alg),
		Value:     value,
		AAD:       w.AAD,
		Iv:        w.IV,
		Tag:       w.Tag,
	}
}

// WrapKeyResponse contains the response for the Client.WrapKey method
type WrapKeyResponse struct {
	KeyOperationResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func wrapKeyResponseFromGenerated(i internal.KeyVaultClientWrapKeyResponse) WrapKeyResponse {
	return WrapKeyResponse{
		RawResponse: i.RawResponse,
		KeyOperationResult: KeyOperationResult{
			AdditionalAuthenticatedData: i.AdditionalAuthenticatedData,
			AuthenticationTag:           i.AuthenticationTag,
			IV:                          i.Iv,
			KeyID:                       i.Kid,
			Result:                      i.Result,
		},
	}
}

func (c *Client) WrapKey(ctx context.Context, alg KeyWrapAlgorithm, key []byte, options *WrapKeyOptions) (WrapKeyResponse, error) {
	if options == nil {
		options = &WrapKeyOptions{}
	}

	resp, err := c.kvClient.WrapKey(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		options.toGeneratedKeyOperationsParameters(alg, key),
		&internal.KeyVaultClientWrapKeyOptions{},
	)

	if err != nil {
		return WrapKeyResponse{}, err
	}

	return wrapKeyResponseFromGenerated(resp), nil
}

type UnwrapKeyOptions struct {
	// Additional data to authenticate but not encrypt/decrypt when using authenticated crypto algorithms.
	AAD []byte `json:"aad,omitempty"`

	// Initialization vector for symmetric algorithms.
	IV []byte `json:"iv,omitempty"`

	// The tag to authenticate when performing decryption with an authenticated algorithm.
	Tag []byte `json:"tag,omitempty"`
}

func (w UnwrapKeyOptions) toGeneratedKeyOperationsParameters(alg KeyWrapAlgorithm, value []byte) internal.KeyOperationsParameters {
	return internal.KeyOperationsParameters{
		Algorithm: (*internal.JSONWebKeyEncryptionAlgorithm)(&alg),
		Value:     value,
		AAD:       w.AAD,
		Iv:        w.IV,
		Tag:       w.Tag,
	}
}

type UnwrapKeyResponse struct {
	KeyOperationResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func unwrapKeyResponseFromGenerated(i internal.KeyVaultClientUnwrapKeyResponse) UnwrapKeyResponse {
	return UnwrapKeyResponse{
		RawResponse: i.RawResponse,
		KeyOperationResult: KeyOperationResult{
			AdditionalAuthenticatedData: i.AdditionalAuthenticatedData,
			AuthenticationTag:           i.AuthenticationTag,
			IV:                          i.Iv,
			KeyID:                       i.Kid,
			Result:                      i.Result,
		},
	}
}

func (c *Client) UnwrapKey(ctx context.Context, alg KeyWrapAlgorithm, encryptedKey []byte, options *UnwrapKeyOptions) (UnwrapKeyResponse, error) {
	if options == nil {
		options = &UnwrapKeyOptions{}
	}

	resp, err := c.kvClient.UnwrapKey(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		options.toGeneratedKeyOperationsParameters(alg, encryptedKey),
		&internal.KeyVaultClientUnwrapKeyOptions{},
	)
	if err != nil {
		return UnwrapKeyResponse{}, err
	}

	return unwrapKeyResponseFromGenerated(resp), nil
}

// SignOptions contains the optional parameters for the Client.Sign method.
type SignOptions struct{}

func (s SignOptions) toGenerated() *internal.KeyVaultClientSignOptions {
	return &internal.KeyVaultClientSignOptions{}
}

// SignResponse contains the response for the Client.Sign method.
type SignResponse struct {
	KeyOperationResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func signResponseFromGenerated(i internal.KeyVaultClientSignResponse) SignResponse {
	return SignResponse{
		RawResponse: i.RawResponse,
		KeyOperationResult: KeyOperationResult{
			AdditionalAuthenticatedData: i.AdditionalAuthenticatedData,
			AuthenticationTag:           i.AuthenticationTag,
			IV:                          i.Iv,
			KeyID:                       i.Kid,
			Result:                      i.Result,
		},
	}
}

func (c *Client) Sign(ctx context.Context, algorithm SignatureAlgorithm, digest []byte, options *SignOptions) (SignResponse, error) {
	if options == nil {
		options = &SignOptions{}
	}

	resp, err := c.kvClient.Sign(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		internal.KeySignParameters{
			Algorithm: (*internal.JSONWebKeySignatureAlgorithm)(&algorithm),
			Value:     digest,
		},
		options.toGenerated(),
	)
	if err != nil {
		return SignResponse{}, err
	}

	return signResponseFromGenerated(resp), nil
}

// VerifyOptions contains the optional parameters for the Client.Verify method
type VerifyOptions struct{}

func (v VerifyOptions) toGenerated() *internal.KeyVaultClientVerifyOptions {
	return &internal.KeyVaultClientVerifyOptions{}
}

// VerifyResponse contains the response for the Client.Verify method
type VerifyResponse struct {
	// READ-ONLY; True if the signature is verified, otherwise false.
	IsValid *bool `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func verifyResponseFromGenerated(i internal.KeyVaultClientVerifyResponse) VerifyResponse {
	return VerifyResponse{
		RawResponse: i.RawResponse,
		IsValid:     i.Value,
	}
}

func (c *Client) Verify(ctx context.Context, algorithm SignatureAlgorithm, digest []byte, signature []byte, options *VerifyOptions) (VerifyResponse, error) {
	if options == nil {
		options = &VerifyOptions{}
	}

	resp, err := c.kvClient.Verify(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		internal.KeyVerifyParameters{
			Algorithm: (*internal.JSONWebKeySignatureAlgorithm)(&algorithm),
			Digest:    digest,
			Signature: signature,
		},
		options.toGenerated(),
	)

	if err != nil {
		return VerifyResponse{}, err
	}

	return verifyResponseFromGenerated(resp), nil
}
