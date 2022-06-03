//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/base"
	generated "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Client performs cryptographic operations using Azure Key Vault Keys. It
// will perform operations locally when it's initialized with the necessary key material or
// is able to get that material from Key Vault. When the required key material is unavailable,
// cryptographic operations are performed by the Key Vault service.
type Client struct {
	base.CryptoClient
}

// ClientOptions are the configurable options on a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// converts ClientOptions to generated *generated.ConnectionOptions
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

	if !strings.HasPrefix(parsed.Path, "/keys/") {
		return "", "", fmt.Errorf("URL is not for a specific key, expect path to start with '/keys/', received %s", id)
	}

	path := strings.Split(strings.TrimPrefix(parsed.Path, "/keys/"), "/")

	if len(path) < 1 {
		return "", "", fmt.Errorf("could not parse Key ID from %s", id)
	}

	if len(path) == 1 {
		return path[0], "", nil
	}

	return path[0], path[1], nil
}

// Parse vault URL from the key identifier
func parseVaultURL(base string) (string, error) {
	parsed, err := url.Parse(base)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s://%s/", parsed.Scheme, parsed.Host), nil
}

// NewClient constructs a Client that performs cryptographic operations with a Key Vault's keys. The keyURL should
// be an identifier of an Azure Key Vault key. Including a version is recommended but not required.
func NewClient(keyURL string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}
	genOptions := options.toConnectionOptions()

	genOptions.PerRetryPolicies = append(
		genOptions.PerRetryPolicies,
		shared.NewKeyVaultChallengePolicy(credential),
	)
	pl := runtime.NewPipeline(internal.ModuleName, internal.ModuleVersion, runtime.PipelineOptions{}, genOptions)

	vaultURL, err := parseVaultURL(keyURL)
	if err != nil {
		return nil, err
	}

	keyID, keyVersion, err := parseKeyIDAndVersion(keyURL)
	if err != nil {
		return nil, err
	}

	return &Client{base.NewCryptoClient(vaultURL, keyID, keyVersion, pl)}, nil
}

// EncryptOptions contains optional parameters for Client.EncryptOptions
type EncryptOptions struct {
	// AuthData is additional data that is authenticated but not encrypted. For use in AES-GCM encryption.
	AuthData []byte `json:"aad,omitempty"`

	// AuthTag is a tag to authenticate when performing decryption with an authenticated algorithm.
	AuthTag []byte `json:"tag,omitempty"`

	// IV is the initialization vector for symmetric algorithms.
	IV []byte `json:"iv,omitempty"`
}

func (e EncryptOptions) toGeneratedKeyOperationsParameters(alg EncryptionAlg, value []byte) generated.KeyOperationsParameters {
	return generated.KeyOperationsParameters{
		Algorithm: (*generated.JSONWebKeyEncryptionAlgorithm)(&alg),
		Value:     value,
		AAD:       e.AuthData,
		Iv:        e.IV,
		Tag:       e.AuthTag,
	}
}

// EncryptResponse is returned by EncryptResponse.
type EncryptResponse struct {
	// Algorithm is the algorithm used to encrypt.
	Algorithm *EncryptionAlg

	// AuthData is additional data that is authenticated but not encrypted. For use in AES-GCM encryption.
	AuthData []byte `json:"aad,omitempty"`

	// AuthTag is a tag to authenticate when performing decryption with an authenticated algorithm.
	AuthTag []byte `json:"tag,omitempty"`

	// Ciphertext is the encrypted data.
	Ciphertext []byte

	// IV is the initialization vector for symmetric algorithms.
	IV []byte `json:"iv,omitempty"`

	// KeyID is the ID of the encrypting key.
	KeyID *string
}

func encryptResponseFromGenerated(i generated.KeyVaultClientEncryptResponse, alg EncryptionAlg) EncryptResponse {
	return EncryptResponse{
		AuthData:   i.AdditionalAuthenticatedData,
		AuthTag:    i.AuthenticationTag,
		IV:         i.Iv,
		KeyID:      i.Kid,
		Ciphertext: i.Result,
		Algorithm:  to.Ptr(alg),
	}
}

func (c *Client) client() *generated.KeyVaultClient {
	return base.Client(c.CryptoClient)
}

func (c *Client) vaultURL() string {
	return base.VaultURL(c.CryptoClient)
}

func (c *Client) keyID() string {
	return base.KeyName(c.CryptoClient)
}

func (c *Client) keyVersion() string {
	return base.KeyVersion(c.CryptoClient)
}

// Encrypt encrypts plaintext using the client's key. This method encrypts only a single block of data, whose
// size dependens on the key and algorithm.
func (c *Client) Encrypt(ctx context.Context, alg EncryptionAlg, plaintext []byte, options *EncryptOptions) (EncryptResponse, error) {
	if options == nil {
		options = &EncryptOptions{}
	}

	resp, err := c.client().Encrypt(
		ctx,
		c.vaultURL(),
		c.keyID(),
		c.keyVersion(),
		options.toGeneratedKeyOperationsParameters(alg, plaintext),
		&generated.KeyVaultClientEncryptOptions{},
	)
	if err != nil {
		return EncryptResponse{}, err
	}

	return encryptResponseFromGenerated(resp, alg), nil
}

// DecryptOptions contains optional parameters for Decrypt.
type DecryptOptions struct {
	// AuthData is additional data that is authenticated but not encrypted. For use in AES-GCM encryption.
	AuthData []byte `json:"aad,omitempty"`

	// AuthTag is a tag to authenticate when performing decryption with an authenticated algorithm.
	AuthTag []byte `json:"tag,omitempty"`

	// IV is the initialization vector for symmetric algorithms.
	IV []byte `json:"iv,omitempty"`
}

func (e DecryptOptions) toGeneratedKeyOperationsParameters(alg EncryptionAlg, value []byte) generated.KeyOperationsParameters {
	return generated.KeyOperationsParameters{
		Algorithm: (*generated.JSONWebKeyEncryptionAlgorithm)(&alg),
		Value:     value,
		AAD:       e.AuthData,
		Iv:        e.IV,
		Tag:       e.AuthTag,
	}
}

// DecryptResponse is returned by Decrypt.
type DecryptResponse struct {
	// Algorithm is the decryption algorithm.
	Algorithm *EncryptionAlg

	// KeyID is the ID of the decrypting key.
	KeyID *string

	// Plaintext is the decrypted bytes.
	Plaintext []byte
}

func decryptResponseFromGenerated(i generated.KeyVaultClientDecryptResponse, alg EncryptionAlg) DecryptResponse {
	return DecryptResponse{
		Algorithm: to.Ptr(alg),
		KeyID:     i.Kid,
		Plaintext: i.Result,
	}
}

// Decrypt decrypts the specified ciphertext. This method decrypts only a single block of data, whose
// size dependens on the key and algorithm.
func (c *Client) Decrypt(ctx context.Context, alg EncryptionAlg, ciphertext []byte, options *DecryptOptions) (DecryptResponse, error) {
	if options == nil {
		options = &DecryptOptions{}
	}

	resp, err := c.client().Decrypt(
		ctx,
		c.vaultURL(),
		c.keyID(),
		c.keyVersion(),
		options.toGeneratedKeyOperationsParameters(alg, ciphertext),
		&generated.KeyVaultClientDecryptOptions{},
	)

	if err != nil {
		return DecryptResponse{}, err
	}

	return decryptResponseFromGenerated(resp, alg), nil
}

// WrapKeyOptions contains optional parameters for WrapKey.
type WrapKeyOptions struct {
	// placeholder for future optional parameters
}

func (w WrapKeyOptions) toGeneratedKeyOperationsParameters(alg WrapAlg, value []byte) generated.KeyOperationsParameters {
	return generated.KeyOperationsParameters{
		Algorithm: (*generated.JSONWebKeyEncryptionAlgorithm)(&alg),
		Value:     value,
	}
}

// WrapKeyResponse is returned by WrapKey.
type WrapKeyResponse struct {
	// Algorithm is the key wrapping algorithm.
	Algorithm *WrapAlg

	// EncryptedKey is the wrapped key.
	EncryptedKey []byte

	// KeyID is the ID of the wrapping key.
	KeyID *string
}

func wrapKeyResponseFromGenerated(i generated.KeyVaultClientWrapKeyResponse, alg WrapAlg) WrapKeyResponse {
	return WrapKeyResponse{
		Algorithm:    to.Ptr(alg),
		KeyID:        i.Kid,
		EncryptedKey: i.Result,
	}
}

// WrapKey encrypts the specified key.
func (c *Client) WrapKey(ctx context.Context, alg WrapAlg, key []byte, options *WrapKeyOptions) (WrapKeyResponse, error) {
	if options == nil {
		options = &WrapKeyOptions{}
	}

	resp, err := c.client().WrapKey(
		ctx,
		c.vaultURL(),
		c.keyID(),
		c.keyVersion(),
		options.toGeneratedKeyOperationsParameters(alg, key),
		&generated.KeyVaultClientWrapKeyOptions{},
	)

	if err != nil {
		return WrapKeyResponse{}, err
	}

	return wrapKeyResponseFromGenerated(resp, alg), nil
}

// UnwrapKeyOptions contains optional parameters for UnwrapKey.
type UnwrapKeyOptions struct {
	// placeholder for future optional parameters
}

func (w UnwrapKeyOptions) toGeneratedKeyOperationsParameters(alg WrapAlg, value []byte) generated.KeyOperationsParameters {
	return generated.KeyOperationsParameters{
		Algorithm: (*generated.JSONWebKeyEncryptionAlgorithm)(&alg),
		Value:     value,
	}
}

// UnwrapKeyResponse is returned by UnwrapKey.
type UnwrapKeyResponse struct {
	// Algorithm is the wrapping algorithm.
	Algorithm *WrapAlg

	// Key is the unwrapped key.
	Key []byte

	// KeyID is the ID of the wrapping key.
	KeyID *string
}

func unwrapKeyResponseFromGenerated(i generated.KeyVaultClientUnwrapKeyResponse, alg WrapAlg) UnwrapKeyResponse {
	return UnwrapKeyResponse{
		KeyID:     i.Kid,
		Key:       i.Result,
		Algorithm: to.Ptr(alg),
	}
}

// UnwrapKey decrypts an encrypted key.
func (c *Client) UnwrapKey(ctx context.Context, alg WrapAlg, encryptedKey []byte, options *UnwrapKeyOptions) (UnwrapKeyResponse, error) {
	if options == nil {
		options = &UnwrapKeyOptions{}
	}

	resp, err := c.client().UnwrapKey(
		ctx,
		c.vaultURL(),
		c.keyID(),
		c.keyVersion(),
		options.toGeneratedKeyOperationsParameters(alg, encryptedKey),
		&generated.KeyVaultClientUnwrapKeyOptions{},
	)
	if err != nil {
		return UnwrapKeyResponse{}, err
	}

	return unwrapKeyResponseFromGenerated(resp, alg), nil
}

// SignOptions contains optional parameters for Sign.
type SignOptions struct {
	// placeholder for future optional parameters
}

func (s SignOptions) toGenerated() *generated.KeyVaultClientSignOptions {
	return &generated.KeyVaultClientSignOptions{}
}

// SignResponse is returned by Sign.
type SignResponse struct {
	// Algorithm is the signing algorithm.
	Algorithm *SignatureAlg

	// KeyID is the ID of the signing key.
	KeyID *string

	// Signature is the signed data.
	Signature []byte
}

func signResponseFromGenerated(i generated.KeyVaultClientSignResponse, alg SignatureAlg) SignResponse {
	return SignResponse{
		Algorithm: to.Ptr(alg),
		KeyID:     i.Kid,
		Signature: i.Result,
	}
}

// Sign signs the specified digest. The hash algorithm used to compute the digest must be compatible with the specified algorithm.
func (c *Client) Sign(ctx context.Context, algorithm SignatureAlg, digest []byte, options *SignOptions) (SignResponse, error) {
	if options == nil {
		options = &SignOptions{}
	}

	resp, err := c.client().Sign(
		ctx,
		c.vaultURL(),
		c.keyID(),
		c.keyVersion(),
		generated.KeySignParameters{
			Algorithm: (*generated.JSONWebKeySignatureAlgorithm)(&algorithm),
			Value:     digest,
		},
		options.toGenerated(),
	)
	if err != nil {
		return SignResponse{}, err
	}

	return signResponseFromGenerated(resp, algorithm), nil
}

// VerifyOptions contains optional parameters for Verify.
type VerifyOptions struct {
	// placeholder for future optional parameters
}

func (v VerifyOptions) toGenerated() *generated.KeyVaultClientVerifyOptions {
	return &generated.KeyVaultClientVerifyOptions{}
}

// VerifyResponse is returned by Verify.
type VerifyResponse struct {
	// Algorithm is the verification algorithm.
	Algorithm *SignatureAlg

	// IsValid is true when the signature is verified.
	IsValid *bool `json:"value,omitempty" azure:"ro"`

	// KeyID is the ID of the verifying key.
	KeyID *string
}

func verifyResponseFromGenerated(i generated.KeyVaultClientVerifyResponse, id string, alg SignatureAlg) VerifyResponse {
	return VerifyResponse{
		IsValid:   i.Value,
		KeyID:     &id,
		Algorithm: to.Ptr(alg),
	}
}

// Verify verifies the specified signature. The algorithm must be the same algorithm used to sign the digest, and
// compatible with the hash algorithm used to compute the digest.
func (c *Client) Verify(ctx context.Context, algorithm SignatureAlg, digest []byte, signature []byte, options *VerifyOptions) (VerifyResponse, error) {
	if options == nil {
		options = &VerifyOptions{}
	}

	resp, err := c.client().Verify(
		ctx,
		c.vaultURL(),
		c.keyID(),
		c.keyVersion(),
		generated.KeyVerifyParameters{
			Algorithm: (*generated.JSONWebKeySignatureAlgorithm)(&algorithm),
			Digest:    digest,
			Signature: signature,
		},
		options.toGenerated(),
	)

	if err != nil {
		return VerifyResponse{}, err
	}

	return verifyResponseFromGenerated(resp, c.keyID(), algorithm), nil
}
