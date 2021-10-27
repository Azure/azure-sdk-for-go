// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"golang.org/x/crypto/pkcs12"
)

// ClientCertificateCredentialOptions contain optional parameters that can be used when configuring a ClientCertificateCredential.
// All zero-value fields will be initialized with their default values.
type ClientCertificateCredentialOptions struct {
	azcore.ClientOptions

	// The password required to decrypt the private key.  Leave empty if there is no password.
	Password string
	// Set to true to include x5c header in client claims when acquiring a token to enable
	// SubjectName and Issuer based authentication for ClientCertificateCredential.
	SendCertificateChain bool
	// The host of the Azure Active Directory authority. The default is AzurePublicCloud.
	// Leave empty to allow overriding the value from the AZURE_AUTHORITY_HOST environment variable.
	AuthorityHost AuthorityHost
}

// ClientCertificateCredential enables authentication of a service principal to Azure Active Directory using a certificate that is assigned to its App Registration. More information
// on how to configure certificate authentication can be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-certificate-credentials#register-your-certificate-with-azure-ad
type ClientCertificateCredential struct {
	client               *aadIdentityClient
	tenantID             string        // The Azure Active Directory tenant (directory) ID of the service principal
	clientID             string        // The client (application) ID of the service principal
	cert                 *certContents // The contents of the certificate file
	sendCertificateChain bool          // Determines whether to include the certificate chain in the claims to retreive a token
}

// NewClientCertificateCredential creates an instance of ClientCertificateCredential with the details needed to authenticate against Azure Active Directory with the specified certificate.
// tenantID: The Azure Active Directory tenant (directory) ID of the service principal.
// clientID: The client (application) ID of the service principal.
// certData: The bytes of a certificate in PEM or PKCS12 format, including the private key.
// options: ClientCertificateCredentialOptions that can be used to provide additional configurations for the credential, such as the certificate password.
func NewClientCertificateCredential(tenantID string, clientID string, certData []byte, options *ClientCertificateCredentialOptions) (*ClientCertificateCredential, error) {
	if !validTenantID(tenantID) {
		return nil, errors.New(tenantIDValidationErr)
	}
	cp := ClientCertificateCredentialOptions{}
	if options != nil {
		cp = *options
	}
	cert, err := loadPEMCert(certData, cp.Password, cp.SendCertificateChain)
	if err != nil {
		cert, err = loadPKCS12Cert(certData, cp.Password, cp.SendCertificateChain)
	}
	if err != nil {
		logCredentialError("Client Certificate Credential", err)
		return nil, err
	}
	authorityHost, err := setAuthorityHost(cp.AuthorityHost)
	if err != nil {
		return nil, err
	}
	c, err := newAADIdentityClient(authorityHost, &cp.ClientOptions)
	if err != nil {
		return nil, err
	}
	return &ClientCertificateCredential{tenantID: tenantID, clientID: clientID, cert: cert, sendCertificateChain: cp.SendCertificateChain, client: c}, nil
}

// contains decoded cert contents we care about
type certContents struct {
	fp                 fingerprint
	pk                 *rsa.PrivateKey
	publicCertificates []string
}

func newCertContents(blocks []*pem.Block, fromPEM bool, sendCertificateChain bool) (*certContents, error) {
	cc := certContents{}
	// first extract the private key
	for _, block := range blocks {
		if block.Type == "PRIVATE KEY" {
			var key interface{}
			var err error
			if fromPEM {
				key, err = x509.ParsePKCS8PrivateKey(block.Bytes)
			} else {
				key, err = x509.ParsePKCS1PrivateKey(block.Bytes)
			}
			if err != nil {
				return nil, err
			}
			rsaKey, ok := key.(*rsa.PrivateKey)
			if !ok {
				return nil, errors.New("unexpected private key type")
			}
			cc.pk = rsaKey
			break
		}
	}
	if cc.pk == nil {
		return nil, errors.New("missing private key")
	}
	// now find the certificate with the matching public key of our private key
	for _, block := range blocks {
		if block.Type == "CERTIFICATE" {
			cert, err := x509.ParseCertificate(block.Bytes)
			if err != nil {
				return nil, err
			}
			certKey, ok := cert.PublicKey.(*rsa.PublicKey)
			if !ok {
				// keep looking
				continue
			}
			if cc.pk.E != certKey.E || cc.pk.N.Cmp(certKey.N) != 0 {
				// keep looking
				continue
			}
			// found a match
			fp, err := newFingerprint(block)
			if err != nil {
				return nil, err
			}
			cc.fp = fp
			break
		}
	}
	if cc.fp == nil {
		return nil, errors.New("missing certificate")
	}
	// now find all the public certificates to send in the x5c header
	if sendCertificateChain {
		for _, block := range blocks {
			if block.Type == "CERTIFICATE" {
				cc.publicCertificates = append(cc.publicCertificates, base64.StdEncoding.EncodeToString(block.Bytes))
			}
		}
	}
	return &cc, nil
}

func loadPEMCert(certData []byte, password string, sendCertificateChain bool) (*certContents, error) {
	// TODO: wire up support for password
	blocks := []*pem.Block{}
	// read all of the PEM blocks
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
	return newCertContents(blocks, true, sendCertificateChain)
}

func loadPKCS12Cert(certData []byte, password string, sendCertificateChain bool) (*certContents, error) {
	// convert data to PEM blocks
	blocks, err := pkcs12.ToPEM(certData, password)
	if err != nil {
		return nil, err
	}
	if len(blocks) == 0 {
		// not mentioning PKCS12 in this message because we end up here when certData is garbage
		return nil, errors.New("didn't find any certificate content")
	}
	return newCertContents(blocks, false, sendCertificateChain)
}

// GetToken obtains a token from Azure Active Directory, using the certificate in the file path.
// scopes: The list of scopes for which the token will have access.
// ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *ClientCertificateCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticateCertificate(ctx, c.tenantID, c.clientID, c.cert, c.sendCertificateChain, opts.Scopes)
	if err != nil {
		addGetTokenFailureLogs("Client Certificate Credential", err, true)
		return nil, err
	}
	logGetTokenSuccess(c, opts)
	return tk, nil
}

var _ azcore.TokenCredential = (*ClientCertificateCredential)(nil)
