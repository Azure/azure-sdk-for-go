//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package crypto

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	generated "github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// The Client performs cryptographic operations using Azure Key Vault Keys. This client
// will perform operations locally when it's initialized with the necessary key material or
// is able to get that material from Key Vault. When the required key material is unavailable,
// cryptographic operations are performed by the Key Vault service.
type Client struct {
	kvClient   *generated.KeyVaultClient
	vaultURL   string
	keyID      string
	keyVersion string
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

// NewClient creates a new azcrytpo.Client that will perform operations against the Key Vault service. The key should
// be an identifier of an Azure Key Vault key. Including a version is recommended but not required.
func NewClient(keyURL string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	// TODO: should this return by pointer or by reference, ask Joel
	if options == nil {
		options = &ClientOptions{}
	}
	genOptions := options.toConnectionOptions()

	genOptions.PerRetryPolicies = append(
		genOptions.PerRetryPolicies,
		shared.NewKeyVaultChallengePolicy(credential),
	)
	pl := runtime.NewPipeline(generated.ModuleName, generated.ModuleVersion, runtime.PipelineOptions{}, genOptions)

	vaultURL, err := parseVaultURL(keyURL)
	if err != nil {
		return nil, err
	}

	keyID, keyVersion, err := parseKeyIDAndVersion(keyURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		kvClient:   generated.NewKeyVaultClient(pl),
		vaultURL:   vaultURL,
		keyID:      keyID,
		keyVersion: keyVersion,
	}, nil
}

// Optional parameters for the crypto.Client.EncryptOptions method
type EncryptOptions struct {
	// Additional data to authenticate but not encrypt/decrypt when using authenticated crypto algorithms.
	AAD []byte `json:"aad,omitempty"`

	// Initialization vector for symmetric algorithms.
	IV []byte `json:"iv,omitempty"`

	// The tag to authenticate when performing decryption with an authenticated algorithm.
	Tag []byte `json:"tag,omitempty"`
}

func (e EncryptOptions) toGeneratedKeyOperationsParameters(alg EncryptionAlgorithm, value []byte) generated.KeyOperationsParameters {
	return generated.KeyOperationsParameters{
		Algorithm: (*generated.JSONWebKeyEncryptionAlgorithm)(&alg),
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

func encryptResponseFromGenerated(i generated.KeyVaultClientEncryptResponse) EncryptResponse {
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

	resp, err := c.kvClient.Encrypt(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		options.toGeneratedKeyOperationsParameters(alg, value),
		&generated.KeyVaultClientEncryptOptions{},
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

func (e DecryptOptions) toGeneratedKeyOperationsParameters(alg EncryptionAlgorithm, value []byte) generated.KeyOperationsParameters {
	return generated.KeyOperationsParameters{
		Algorithm: (*generated.JSONWebKeyEncryptionAlgorithm)(&alg),
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

func decryptResponseFromGenerated(i generated.KeyVaultClientDecryptResponse) DecryptResponse {
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

	resp, err := c.kvClient.Decrypt(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		options.toGeneratedKeyOperationsParameters(alg, ciphertext),
		&generated.KeyVaultClientDecryptOptions{},
	)

	if err != nil {
		return DecryptResponse{}, err
	}

	return decryptResponseFromGenerated(resp), nil
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

func (w WrapKeyOptions) toGeneratedKeyOperationsParameters(alg KeyWrapAlgorithm, value []byte) generated.KeyOperationsParameters {
	return generated.KeyOperationsParameters{
		Algorithm: (*generated.JSONWebKeyEncryptionAlgorithm)(&alg),
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

func wrapKeyResponseFromGenerated(i generated.KeyVaultClientWrapKeyResponse) WrapKeyResponse {
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

	resp, err := c.kvClient.WrapKey(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		options.toGeneratedKeyOperationsParameters(alg, key),
		&generated.KeyVaultClientWrapKeyOptions{},
	)

	if err != nil {
		return WrapKeyResponse{}, err
	}

	return wrapKeyResponseFromGenerated(resp), nil
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

func (w UnwrapKeyOptions) toGeneratedKeyOperationsParameters(alg KeyWrapAlgorithm, value []byte) generated.KeyOperationsParameters {
	return generated.KeyOperationsParameters{
		Algorithm: (*generated.JSONWebKeyEncryptionAlgorithm)(&alg),
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

func unwrapKeyResponseFromGenerated(i generated.KeyVaultClientUnwrapKeyResponse) UnwrapKeyResponse {
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

	resp, err := c.kvClient.UnwrapKey(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		options.toGeneratedKeyOperationsParameters(alg, encryptedKey),
		&generated.KeyVaultClientUnwrapKeyOptions{},
	)
	if err != nil {
		return UnwrapKeyResponse{}, err
	}

	return unwrapKeyResponseFromGenerated(resp), nil
}

// SignOptions contains the optional parameters for the Client.Sign method.
type SignOptions struct{}

func (s SignOptions) toGenerated() *generated.KeyVaultClientSignOptions {
	return &generated.KeyVaultClientSignOptions{}
}

// SignResponse contains the response for the Client.Sign method.
type SignResponse struct {
	KeyOperationResult
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func signResponseFromGenerated(i generated.KeyVaultClientSignResponse) SignResponse {
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

	resp, err := c.kvClient.Sign(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
		generated.KeySignParameters{
			Algorithm: (*generated.JSONWebKeySignatureAlgorithm)(&algorithm),
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

func (v VerifyOptions) toGenerated() *generated.KeyVaultClientVerifyOptions {
	return &generated.KeyVaultClientVerifyOptions{}
}

// VerifyResponse contains the response for the Client.Verify method
type VerifyResponse struct {
	// READ-ONLY; True if the signature is verified, otherwise false.
	IsValid *bool `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func verifyResponseFromGenerated(i generated.KeyVaultClientVerifyResponse) VerifyResponse {
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

	resp, err := c.kvClient.Verify(
		ctx,
		c.vaultURL,
		c.keyID,
		c.keyVersion,
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

	return verifyResponseFromGenerated(resp), nil
}
