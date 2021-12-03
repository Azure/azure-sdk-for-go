// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"

	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"golang.org/x/crypto/pkcs12"
)

// ClientCertificateCredentialOptions contains optional parameters for ClientCertificateCredential.
type ClientCertificateCredentialOptions struct {
	azcore.ClientOptions

	// SendCertificateChain controls whether the credential sends the public certificate chain in the x5c
	// header of each token request's JWT. This is required for Subject Name/Issuer (SNI) authentication.
	// Defaults to False.
	SendCertificateChain bool
	// AuthorityHost is the base URL of an Azure Active Directory authority. Defaults
	// to the value of environment variable AZURE_AUTHORITY_HOST, if set, or AzurePublicCloud.
	AuthorityHost AuthorityHost
}

// ClientCertificateCredential authenticates a service principal with a certificate.
type ClientCertificateCredential struct {
	client               *aadIdentityClient
	tenantID             string
	clientID             string
	cert                 *certContents
	sendCertificateChain bool
}

// NewClientCertificateCredential constructs a ClientCertificateCredential.
// tenantID: The application's Azure Active Directory tenant or directory ID.
// clientID: The application's client ID.
// certs: one or more certificates, for example as returned by ParseCertificates()
// key: the signing certificate's private key, for example as returned by ParseCertificates()
// options: Optional configuration.
func NewClientCertificateCredential(tenantID string, clientID string, certs []*x509.Certificate, key crypto.PrivateKey, options *ClientCertificateCredentialOptions) (*ClientCertificateCredential, error) {
	if len(certs) == 0 {
		return nil, errors.New("at least one certificate is required")
	}
	pk, ok := key.(*rsa.PrivateKey)
	if !ok {
		return nil, errors.New("'key' must be an *rsa.PrivateKey")
	}
	if !validTenantID(tenantID) {
		return nil, errors.New(tenantIDValidationErr)
	}
	cp := ClientCertificateCredentialOptions{}
	if options != nil {
		cp = *options
	}
	authorityHost, err := setAuthorityHost(cp.AuthorityHost)
	if err != nil {
		logCredentialError("Client Certificate Credential", err)
		return nil, err
	}
	cert, err := newCertContents(certs, pk, cp.SendCertificateChain)
	if err != nil {
		return nil, err
	}
	c, err := newAADIdentityClient(authorityHost, &cp.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &ClientCertificateCredential{tenantID: tenantID, clientID: clientID, cert: cert, sendCertificateChain: cp.SendCertificateChain, client: c}, nil
}

// GetToken obtains a token from Azure Active Directory. This method is called automatically by Azure SDK clients.
// ctx: Context controlling the request lifetime.
// opts: Options for the token request, in particular the desired scope of the access token.
func (c *ClientCertificateCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticateCertificate(ctx, c.tenantID, c.clientID, c.cert, c.sendCertificateChain, opts.Scopes)
	if err != nil {
		addGetTokenFailureLogs("Client Certificate Credential", err, true)
		return nil, err
	}
	logGetTokenSuccess(c, opts)
	return tk, nil
}

// ParseCertificates loads certificates and a private key for use with NewClientCertificateCredential.
// certData: certificate data encoded in PEM or PKCS12 format, including the certificate's private key.
// password: the password required to decrypt the private key. Pass nil if the key is not encrypted. This function can't decrypt keys in PEM format.
func ParseCertificates(certData []byte, password []byte) ([]*x509.Certificate, crypto.PrivateKey, error) {
	var blocks []*pem.Block
	var err error
	if len(password) == 0 {
		blocks, err = loadPEMCert(certData)
	}
	if len(blocks) == 0 || err != nil {
		blocks, err = loadPKCS12Cert(certData, string(password))
	}
	if err != nil {
		return nil, nil, err
	}
	var certs []*x509.Certificate
	var pk crypto.PrivateKey
	for _, block := range blocks {
		switch block.Type {
		case "CERTIFICATE":
			c, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, nil, err
			}
			certs = append(certs, c)
		case "PRIVATE KEY":
			if pk != nil {
				return nil, nil, errors.New("certData contains multiple private keys")
			}
			pk, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				pk, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			}
			if err != nil {
				return nil, nil, err
			}
		case "RSA PRIVATE KEY":
			if pk != nil {
				return nil, nil, errors.New("certData contains multiple private keys")
			}
			pk, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return nil, nil, err
			}
		}
	}
	if len(certs) == 0 {
		return nil, nil, errors.New("found no certificate")
	}
	if pk == nil {
		return nil, nil, errors.New("found no private key")
	}
	return certs, pk, nil
}

type certContents struct {
	fp  []byte          // the signing cert's fingerprint, a SHA-1 digest
	pk  *rsa.PrivateKey // the signing key
	x5c []string        // concatenation of every provided cert, base64 encoded
}

func newCertContents(certs []*x509.Certificate, key *rsa.PrivateKey, sendCertificateChain bool) (*certContents, error) {
	cc := certContents{pk: key}
	// need the the signing cert's fingerprint: identify that cert by matching its public key to the private key
	for _, cert := range certs {
		certKey, ok := cert.PublicKey.(*rsa.PublicKey)
		if ok && key.E == certKey.E && key.N.Cmp(certKey.N) == 0 {
			fp := sha1.Sum(cert.Raw)
			cc.fp = fp[:]
			if sendCertificateChain {
				// signing cert must be first in x5c
				cc.x5c = append([]string{base64.StdEncoding.EncodeToString(cert.Raw)}, cc.x5c...)
			}
		} else if sendCertificateChain {
			cc.x5c = append(cc.x5c, base64.StdEncoding.EncodeToString(cert.Raw))
		}
	}
	if len(cc.fp) == 0 {
		return nil, errors.New("found no certificate matching 'key'")
	}
	return &cc, nil
}

func loadPEMCert(certData []byte) ([]*pem.Block, error) {
	blocks := []*pem.Block{}
	for {
		var block *pem.Block
		block, certData = pem.Decode(certData)
		if block == nil {
			break
		}
		blocks = append(blocks, block)
	}
	if len(blocks) == 0 {
		return nil, errors.New("didn't find any PEM blocks")
	}
	return blocks, nil
}

func loadPKCS12Cert(certData []byte, password string) ([]*pem.Block, error) {
	blocks, err := pkcs12.ToPEM(certData, password)
	if err != nil {
		return nil, err
	}
	if len(blocks) == 0 {
		// not mentioning PKCS12 in this message because we end up here when certData is garbage
		return nil, errors.New("didn't find any certificate content")
	}
	return blocks, err
}

var _ azcore.TokenCredential = (*ClientCertificateCredential)(nil)
