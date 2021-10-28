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
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/auth"
	internal "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
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
	useService bool
	cred       azcore.TokenCredential
	transport  policy.Transporter
}

// ClientOptions are the configurable options on a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// converts ClientOptions to generated *internal.ConnectionOptions
func (c *ClientOptions) toConnectionOptions() *policy.ClientOptions {
	if c == nil {
		return nil
	}
	return &policy.ClientOptions{
		Transport:        c.Transport,
		Retry:            c.Retry,
		Telemetry:        c.Telemetry,
		Logging:          c.Logging,
		PerCallPolicies:  c.PerCallPolicies,
		PerRetryPolicies: c.PerRetryPolicies,
	}
}

// parse the KeyID and Version. If no version is present, return a blank string.
func parseKeyIDAndVersion(id string) (string, string, error) {
	parsed, err := url.Parse(id)
	if err != nil {
		return "", "", err
	}

	path := strings.Split(parsed.Path, "/")

	if len(path) < 3 {
		return "", "", fmt.Errorf("could not parse Key ID from %s", id)
	}

	if len(path) == 3 {
		return path[2], "", nil
	}

	return path[2], path[3], nil
}

// Parse vault URL from the key identifier
func parseVaultURL(base string) (string, error) {
	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s/", parsed.Scheme, parsed.Host), nil
}

// NewClient creates a new azcrytpo.Client that will perform operations against the Key Vault service. The key should
// be an identifier of an Azure Key Vault key. Including a version is recommended but not required.
func NewClient(key string, cred azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	// Have to have a transport for the challenge policy
	if options.Transport == nil {
		options.Transport = http.DefaultClient
	}

	options.PerRetryPolicies = append(
		options.PerRetryPolicies,
		&auth.KeyVaultChallengePolicy{
			Cred:      cred,
			Transport: options.Transport,
		},
	)
	conn := internal.NewConnection(options.toConnectionOptions())

	vaultURL, err := parseVaultURL(key)
	if err != nil {
		return &Client{}, err
	}

	keyID, keyVersion, err := parseKeyIDAndVersion(key)
	if err != nil {
		return &Client{}, err
	}

	return &Client{
		kvClient:   internal.NewKeyVaultClient(conn),
		vaultURL:   vaultURL,
		keyID:      keyID,
		keyVersion: keyVersion,
		useService: true,
		cred:       cred,
		transport:  options.Transport,
	}, nil
}

// Creates a new Client that can only perform cryptographic operations locally.
func NewClientFromJSONWebKey(key azkeys.JSONWebKey) (*Client, error) {
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

// The ENCRYPT operation encrypts an arbitrary sequence of bytes using an encryption key that is stored in
// Azure Key Vault. Note that the ENCRYPT operation only supports a single block of data, the size of which
// is dependent on the target key and the encryption algorithm to be used. The ENCRYPT operation is only
// strictly necessary for symmetric keys stored in Azure Key Vault since protection with an asymmetric key
// can be performed using public portion of the key. This operation is supported for asymmetric keys as a
// convenience for callers that have a key-reference but do not have access to the public key material.
// This operation requires the keys/encrypt permission.
func (c *Client) Encrypt(ctx context.Context, alg EncryptionAlgorithm, value []byte, options *EncryptOptions) (EncryptResponse, error) {
	if options == nil {
		options = &EncryptOptions{}
	}

	if c.useService {
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
	return EncryptResponse{}, nil
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

// The DECRYPT operation decrypts a well-formed block of ciphertext using the target encryption key and
// specified algorithm. This operation is the reverse of the ENCRYPT operation; only a single block of
// data may be decrypted, the size of this block is dependent on the target key and the algorithm to be
// used. The DECRYPT operation applies to asymmetric and symmetric keys stored in Azure Key Vault since
// it uses the private portion of the key. This operation requires the keys/decrypt permission.
func (c *Client) Decrypt(ctx context.Context, alg EncryptionAlgorithm, ciphertext []byte, options *DecryptOptions) (DecryptResponse, error) {
	if options == nil {
		options = &DecryptOptions{}
	}

	if c.useService {
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
	return DecryptResponse{}, nil
}

// WrapKeyOptions contains the optional parameters for the Client.WrapKey method.
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

// The WRAP operation supports encryption of a symmetric key using a key encryption key that has previously
// been stored in an Azure Key Vault. The WRAP operation is only strictly necessary for symmetric keys stored
//in Azure Key Vault since protection with an asymmetric key can be performed using the public portion of
// the key. This operation is supported for asymmetric keys as a convenience for callers that have a
// key-reference but do not have access to the public key material. This operation requires the keys/wrapKey permission.
func (c *Client) WrapKey(ctx context.Context, alg KeyWrapAlgorithm, key []byte, options *WrapKeyOptions) (WrapKeyResponse, error) {
	if options == nil {
		options = &WrapKeyOptions{}
	}

	if c.useService {
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
	return WrapKeyResponse{}, nil
}

// UnwrapKeyOptions contains the optional parameters for the Client.UnwrapKey method
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

// UnwrapKeyResponse contains the response for the Client.UnwrapKey method
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

// UnwrapKey - The UNWRAP operation supports decryption of a symmetric key using the target key encryption key.
// This operation is the reverse of the WRAP operation. The UNWRAP operation applies to asymmetric and symmetric
// keys stored in Azure Key Vault since it uses the private portion of the key. This operation requires the
// keys/unwrapKey permission.
func (c *Client) UnwrapKey(ctx context.Context, alg KeyWrapAlgorithm, encryptedKey []byte, options *UnwrapKeyOptions) (UnwrapKeyResponse, error) {
	if options == nil {
		options = &UnwrapKeyOptions{}
	}

	if c.useService {
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
	return UnwrapKeyResponse{}, nil
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

// The SIGN operation is applicable to asymmetric and symmetric keys stored in Azure Key Vault since
// this operation uses the private portion of the key. This operation requires the keys/sign permission.
func (c *Client) Sign(ctx context.Context, algorithm SignatureAlgorithm, digest []byte, options *SignOptions) (SignResponse, error) {
	if options == nil {
		options = &SignOptions{}
	}

	if c.useService {
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
	return SignResponse{}, nil
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

// The VERIFY operation is applicable to symmetric keys stored in Azure Key Vault. VERIFY is not strictly
// necessary for asymmetric keys stored in Azure Key Vault since signature verification can be performed
// using the public portion of the key but this operation is supported as a convenience for callers that
// only have a key-reference and not the public portion of the key. This operation requires the keys/verify permission.
func (c *Client) Verify(ctx context.Context, algorithm SignatureAlgorithm, digest []byte, signature []byte, options *VerifyOptions) (VerifyResponse, error) {
	if options == nil {
		options = &VerifyOptions{}
	}

	if c.useService {
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
	return VerifyResponse{}, nil
}
