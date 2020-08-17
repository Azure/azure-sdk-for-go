// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License.

package azidentity

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"golang.org/x/crypto/pkcs12"
)

// ClientCertificateCredential enables authentication of a service principal to Azure Active Directory using a certificate that is assigned to its App Registration. More information
// on how to configure certificate authentication can be found here:
// https://docs.microsoft.com/en-us/azure/active-directory/develop/active-directory-certificate-credentials#register-your-certificate-with-azure-ad
type ClientCertificateCredential struct {
	client   *aadIdentityClient
	tenantID string        // The Azure Active Directory tenant (directory) ID of the service principal
	clientID string        // The client (application) ID of the service principal
	cert     *certContents // The contents of the certificate file
}

// NewClientCertificateCredential creates an instance of ClientCertificateCredential with the details needed to authenticate against Azure Active Directory with the specified certificate.
// tenantID: The Azure Active Directory tenant (directory) ID of the service principal.
// clientID: The client (application) ID of the service principal.
// clientCertificate: The path to the client certificate used to authenticate the client.  Supported formats are PEM and PFX.
// password: The password required to decrypt the private key.  Pass nil if there is no password.
// options: configure the management of the requests sent to Azure Active Directory.
func NewClientCertificateCredential(tenantID string, clientID string, clientCertificate string, password *string, options *TokenCredentialOptions) (*ClientCertificateCredential, error) {
	_, err := os.Stat(clientCertificate)
	if err != nil {
		credErr := &CredentialUnavailableError{CredentialType: "Client Certificate Credential", Message: "Certificate file not found in path: " + clientCertificate}
		azcore.Log().Write(azcore.LogError, logCredentialError(credErr.CredentialType, credErr))
		return nil, credErr
	}
	certData, err := ioutil.ReadFile(clientCertificate)
	if err != nil {
		credErr := &CredentialUnavailableError{CredentialType: "Client Certificate Credential", Message: err.Error()}
		azcore.Log().Write(azcore.LogError, logCredentialError(credErr.CredentialType, credErr))
		return nil, credErr
	}
	var cert *certContents
	clientCertificate = strings.ToUpper(clientCertificate)
	if strings.HasSuffix(clientCertificate, ".PEM") {
		cert, err = extractFromPEMFile(certData, password)
	} else if strings.HasSuffix(clientCertificate, ".PFX") {
		cert, err = extractFromPFXFile(certData, password)
	} else {
		err = errors.New("only PEM and PFX files are supported")
	}
	if err != nil {
		credErr := &CredentialUnavailableError{CredentialType: "Client Certificate Credential", Message: err.Error()}
		azcore.Log().Write(azcore.LogError, logCredentialError(credErr.CredentialType, credErr))
		return nil, credErr
	}
	c, err := newAADIdentityClient(options)
	if err != nil {
		return nil, err
	}
	return &ClientCertificateCredential{tenantID: tenantID, clientID: clientID, cert: cert, client: c}, nil
}

// contains decoded cert contents we care about
type certContents struct {
	fp fingerprint
	pk *rsa.PrivateKey
}

func newCertContents(blocks []*pem.Block, fromPEM bool) (*certContents, error) {
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
	return &cc, nil
}

func extractFromPEMFile(certData []byte, password *string) (*certContents, error) {
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
		return nil, errors.New("didn't find any blocks in PEM file")
	}
	return newCertContents(blocks, true)
}

func extractFromPFXFile(certData []byte, password *string) (*certContents, error) {
	if password == nil {
		empty := ""
		password = &empty
	}
	// convert PFX binary data to PEM blocks
	blocks, err := pkcs12.ToPEM(certData, *password)
	if err != nil {
		return nil, err
	}
	if len(blocks) == 0 {
		return nil, errors.New("didn't find any blocks in PFX file")
	}
	return newCertContents(blocks, false)
}

// GetToken obtains a token from Azure Active Directory, using the certificate in the file path.
// scopes: The list of scopes for which the token will have access.
// ctx: controlling the request lifetime.
// Returns an AccessToken which can be used to authenticate service client calls.
func (c *ClientCertificateCredential) GetToken(ctx context.Context, opts azcore.TokenRequestOptions) (*azcore.AccessToken, error) {
	tk, err := c.client.authenticateCertificate(ctx, c.tenantID, c.clientID, c.cert, opts.Scopes)
	if err != nil {
		addGetTokenFailureLogs("Client Certificate Credential", err)
		return nil, err
	}
	azcore.Log().Write(LogCredential, logGetTokenSuccess(c, opts))
	return tk, nil
}

// AuthenticationPolicy implements the azcore.Credential interface on ClientSecretCredential.
func (c *ClientCertificateCredential) AuthenticationPolicy(options azcore.AuthenticationPolicyOptions) azcore.Policy {
	return newBearerTokenPolicy(c, options)
}
