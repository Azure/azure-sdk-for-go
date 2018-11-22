package easykey

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"github.com/Azure/go-autorest/autorest/date"
	"golang.org/x/crypto/pkcs12"
)

// Certificate describes certificate information stored in KeyVault.
type Certificate struct {
	// ID identifying the certificate within keyvault.
	ID ID
	// Thumbprint is the X509 thumbprint of the certificate a (URL-encoded base64 string).
	Thumbprint string
	// Content is the CER contents of x509 certificate.
	Content []byte
	// Attr is the attributes of the certificate.
	Attr CertAttr
}

// version implements versioner.
func (c Certificate) version() string {
	return c.ID.Version
}

func (c *Certificate) fromBundle(bundle keyvault.CertificateBundle) error {
	u, _ := url.Parse(*bundle.ID)

	var err error
	c.ID, err = urlToID(u)
	if err != nil {
		return err
	}
	if bundle.X509Thumbprint != nil {
		c.Thumbprint = *bundle.X509Thumbprint
	}
	if bundle.Cer != nil {
		c.Content = *bundle.Cer
	}
	if bundle.Attributes != nil {
		return c.Attr.fromBundle(*bundle.Attributes)
	}
	return nil
}

func (c Certificate) toBundle() keyvault.CertificateBundle {
	id := c.ID.String()
	content := make([]byte, len(c.Content))
	copy(content, c.Content)
	return keyvault.CertificateBundle{
		ID:             &id,
		X509Thumbprint: &c.Thumbprint,
		Cer:            &content,
		Attributes:     c.Attr.toBundle(),
	}
}

// X509 translates .Content into a x509.Certificate.
func (c Certificate) X509() (*x509.Certificate, error) {
	return x509.ParseCertificate(c.Content)
}

// CertAttr is attributes related to the storage of this certificate.
type CertAttr struct {
	// RecoveryLevel the level of recovery for this certificate when deleted.  See the description of
	// DeletionRecoveryLevel above.
	RecoveryLevel DeletionRecoveryLevel
	// Enabled indicates if the secret is currently enabled.
	Enabled bool
	// Expires indicates the time the certificate will expire. If set to the zero value, it indicates
	// this was not set.
	Expires time.Time
	// Created indicates the time the secret was created in UTC. If set to the zero value, it indicates
	// this was not set.
	Created time.Time
	// NotBefore indicate that the key isn't valid before this time in UTC. If set to the zero value, it indicates
	// this was not set.
	NotBefore time.Time
	// Updated indicates the last time the secret was updated in UTC. If set to the zero value, it indicates
	// this was not set.
	Updated time.Time
}

func (c *CertAttr) fromBundle(bundle keyvault.CertificateAttributes) error {
	c.RecoveryLevel = bundle.RecoveryLevel
	c.Enabled = *bundle.Enabled
	if bundle.Expires != nil {
		c.Expires = time.Time(*bundle.Expires)
	}
	if bundle.NotBefore != nil {
		c.NotBefore = time.Time(*bundle.NotBefore)
	}
	if bundle.Created != nil {
		c.Created = time.Time(*bundle.Created)
	}
	if bundle.Updated != nil {
		c.Updated = time.Time(*bundle.Updated)
	}
	return nil
}

func (c CertAttr) toBundle() *keyvault.CertificateAttributes {
	a := &keyvault.CertificateAttributes{
		RecoveryLevel: c.RecoveryLevel,
		Enabled:       &c.Enabled,
	}

	if c.NotBefore.IsZero() {
		z := date.NewUnixTimeFromSeconds(0.0)
		a.NotBefore = &z
	} else {
		d := date.UnixTime(c.NotBefore)
		a.NotBefore = &d
	}
	if c.Expires.IsZero() {
		z := date.NewUnixTimeFromSeconds(0.0)
		a.Expires = &z
	} else {
		d := date.UnixTime(c.Expires)
		a.Expires = &d
	}
	if c.Created.IsZero() {
		z := date.NewUnixTimeFromSeconds(0.0)
		a.Created = &z
	} else {
		d := date.UnixTime(c.Created)
		a.Created = &d
	}
	if c.Updated.IsZero() {
		z := date.NewUnixTimeFromSeconds(0.0)
		a.Updated = &z
	} else {
		d := date.UnixTime(c.Updated)
		a.Updated = &d
	}
	return a
}

// Certificate returns a certificate from the vault that is defined by name at version.
// version can be set to constant LatestVersion to retreive the current version of the cert.
func (k *KeyVault) Certificate(ctx context.Context, name, version string) (Certificate, error) {
	bundle, err := k.vault.GetCertificate(ctx, k.vaultURL, name, version)
	if err != nil {
		if strings.Contains(err.Error(), "StatusCode=403") {
			return Certificate{}, fmt.Errorf("access denied to keyvault(%s) for cert name(%s), "+
				"often due to not having an access policy set for the user in AZURE_CLIENT_ID or missing "+
				"environmental variable AZURE_CLIENT_SECRET: %s", k.vaultName, name, err)
		}
		return Certificate{}, fmt.Errorf("could not connect to vault %s and retrieve key %q: %s", k.vaultURL, name, err)
	}
	cert := Certificate{}
	if err := cert.fromBundle(bundle); err != nil {
		return Certificate{}, nil
	}

	if len(cert.Content) == 0 {
		return Certificate{}, fmt.Errorf("certificate had no content")
	}
	return cert, nil
}

// PrivateKey returns the private key after it has been bases64 decoded.  This will be in pkcs12 format
// with an empty password.  Use golang.org/x/crypto/pkcs12 to decode it further.
// It trying to use with TLS for a net.HTTP server, TLSCert() is probably what you want.
func (k *KeyVault) PrivateKey(ctx context.Context, name, version string) ([]byte, error) {
	secret, err := k.Secret(ctx, name, version)
	if err != nil {
		return nil, err
	}

	secretData, err := base64.StdEncoding.DecodeString(secret.Value)
	if err != nil {
		return nil, fmt.Errorf("problem base64 decoding our private key: %s", err)
	}

	return secretData, nil
}

/*
TLSCert returns a tls.Certificate that can be used in TLS connections in the net/http package.
This assumes that the cert held in Keyvault has the Private key in the first decoded PEM block and the
public certificate in the second PEM block. Because of the many different ways key pairs can be generated
and labelled we don't know which PEM block will have your private key and public certificate. Normally keyBlockNum == 0
and certBlockNum == 1. If that doesn't work, use PEM() and loop through each PEM block printing out the .Type
to locate the indexes you need.

NOTE: This requires that you are careful about how your certificates are generated in the future. If not and you
change the order, this could cause issues with your service.

Here is a quick way to use the cert in your service.

	cert, err := kv.TLSCert(ctx, "certname", LatestVersion)
	if err != nil {
		panic(err)
	}

	cfg := &tls.Config{Certificates: []tls.Certificate{cert}}
	srv := &http.Server{
		TLSConfig:    cfg,
		ReadTimeout:  time.Minute,
		WriteTimeout: time.Minute,
	}
	log.Fatal(srv.ListenAndServeTLS("", ""))
*/
func (k *KeyVault) TLSCert(ctx context.Context, name, version string, certBlockIndex, keyBlockIndex int) (tls.Certificate, error) {
	key, err := k.PrivateKey(ctx, name, version)
	if err != nil {
		return tls.Certificate{}, err
	}

	blocks, err := pkcs12.ToPEM(key, "")
	if err != nil {
		return tls.Certificate{}, err
	}

	if len(blocks) < 2 {
		return tls.Certificate{}, fmt.Errorf("decoded less than two blocks in PKCS12 key, so could not create a X509KeyPair")
	}

	return tls.X509KeyPair(pem.EncodeToMemory(blocks[certBlockIndex]), pem.EncodeToMemory(blocks[keyBlockIndex]))
}

// PEM decodes our PKCS12 private key and returns it as PEM encoded blocks for analysis.
// Generally useful if you wish to look at your key/certificate to determine where the public/private
// ones are located for use in TLSCert().
func (k *KeyVault) PEM(ctx context.Context, name, version string) ([]*pem.Block, error) {
	key, err := k.PrivateKey(ctx, name, version)
	if err != nil {
		return nil, err
	}

	return pkcs12.ToPEM(key, "")
}
