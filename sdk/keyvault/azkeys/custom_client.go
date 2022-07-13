//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

// this file contains handwritten additions to the generated code

import (
	"context"
	"crypto"
	"crypto/ecdsa"
	"crypto/rsa"
	"errors"
	"io"
	"math/big"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
	"golang.org/x/crypto/cryptobyte"
	"golang.org/x/crypto/cryptobyte/asn1"
)

// NewClient creates a client that accesses a Key Vault's keys.
func NewClient(vaultURL string, credential azcore.TokenCredential, options *azcore.ClientOptions) *Client {
	authPolicy := internal.NewKeyVaultChallengePolicy(credential)
	pl := runtime.NewPipeline(moduleName, version, runtime.PipelineOptions{PerRetry: []policy.Policy{authPolicy}}, options)
	return &Client{endpoint: vaultURL, pl: pl}
}

// ID is a key's unique ID, containing its version, if any, and name.
type ID string

// Name of the key.
func (i *ID) Name() string {
	_, name, _ := internal.ParseID((*string)(i))
	return *name
}

// Version of the key. This returns an empty string when the ID contains no version.
func (i *ID) Version() string {
	_, _, version := internal.ParseID((*string)(i))
	if version == nil {
		return ""
	}
	return *version
}

// SignerForKey - Returns an object that implements the crypto.Signer interface, using asymmetric keys stored in Azure Key Vault.
// The context is maintained for all operations performed by the signer object.
// The operations performed by the Signer require the keys/sign and keys/list permissions.
// name - The name of the key.
// version - The version of the key.
func (client *Client) SignerForKey(ctx context.Context, name string, version string) crypto.Signer {
	return &signer{
		client:     client,
		keyName:    name,
		keyVersion: version,
		ctx:        ctx,
	}
}

// signer implements the crypto.Signer interface and allows performing signing operations using asymmetric
type signer struct {
	client     *Client
	keyName    string
	keyVersion string
	ctx        context.Context

	// Cached public key
	publicKey crypto.PublicKey
}

// Public returns the public key corresponding to the private key stored in the Key Vault..
func (signer *signer) Public() crypto.PublicKey {
	pk, _ := signer.getPublicKey()
	return pk
}

// Sign signs a value with the private key stored in the Key Vault.
// For an RSA key, the resulting signature should be either a
// PKCS #1 v1.5 or PSS signature (as indicated by opts). For an (EC)DSA
// key, it should be a DER-serialised, ASN.1 signature structure.
//
// Hash implements the SignerOpts interface and, in most cases, one can
// simply pass in the hash function used as opts. Sign may also attempt
// to type assert opts to other types in order to obtain algorithm
// specific values. See the documentation in each package for details.
func (signer *signer) Sign(_ io.Reader, value []byte, opts crypto.SignerOpts) ([]byte, error) {
	params := SignParameters{
		Value: value,
	}

	// We need the public key to be able to determine the algorithm to use
	pk, err := signer.getPublicKey()
	if err != nil {
		return nil, err
	}

	// Determine the algorithm
	signAlg := ""
	switch pk.(type) {
	case *rsa.PublicKey:
		// If we have an opts object and that is of type *rsa.PSSOptions, then use RSA-PSS
		// (however, we are ignoring the SaltLength property of the object as it's not supported by Key Vault).
		// Otherwise, use RSA in PKCS #1 v1.5 mode.
		if opts != nil {
			if _, ok := opts.(*rsa.PSSOptions); ok {
				signAlg = "PS"
			}
		}
		if signAlg == "" {
			signAlg = "RS"
		}
	case *ecdsa.PublicKey:
		signAlg = "ES"
	default:
		return nil, errors.New("unsupported key type")
	}

	// Determine the hashing function
	var hf crypto.Hash = crypto.SHA256
	if opts != nil {
		hf = opts.HashFunc()
	}
	switch hf {
	case crypto.SHA256:
		params.Algorithm = to.Ptr(JSONWebKeySignatureAlgorithm(signAlg + "256"))
	case crypto.SHA384:
		params.Algorithm = to.Ptr(JSONWebKeySignatureAlgorithm(signAlg + "384"))
	case crypto.SHA512:
		params.Algorithm = to.Ptr(JSONWebKeySignatureAlgorithm(signAlg + "512"))
	default:
		return nil, errors.New("unsupported hashing function")
	}

	// Make the request
	res, err := signer.client.Sign(signer.ctx, signer.keyName, signer.keyVersion, params, nil)
	if err != nil {
		return nil, err
	}
	l := len(res.Result)
	if l == 0 {
		return nil, errors.New("response did not contain a result")
	}

	// When using ECDSA, Azure Key Vault returns a signature formatted per ANSI X9.62 (that is: "r|s")
	// We need to convert that to ASN DER.1 which is the format crypto.Signer.Sign expects
	if signAlg == "ES" {
		r := (&big.Int{}).SetBytes(res.Result[:(l / 2)])
		s := (&big.Int{}).SetBytes(res.Result[(l / 2):])
		var b cryptobyte.Builder
		b.AddASN1(asn1.SEQUENCE, func(b *cryptobyte.Builder) {
			b.AddASN1BigInt(r)
			b.AddASN1BigInt(s)
		})
		return b.Bytes()
	}

	return res.Result, nil
}

func (signer *signer) getPublicKey() (crypto.PublicKey, error) {
	if signer.publicKey != nil {
		return signer.publicKey, nil
	}

	res, err := signer.client.GetKey(signer.ctx, signer.keyName, signer.keyVersion, nil)
	if err != nil || res.Key == nil {
		return nil, err
	}

	pk, err := res.Key.Public()
	if err != nil {
		return nil, err
	}

	signer.publicKey = pk

	if signer.keyVersion == "" && res.Key.KID != nil {
		signer.keyVersion = res.Key.KID.Version()
	}

	return pk, nil
}
