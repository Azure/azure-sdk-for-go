package easykey

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/url"
	"sync"
	"time"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
	"golang.org/x/crypto/pkcs12"
)

type versioner interface {
	version() string
}

type keyLookup map[string][]versioner

func (k keyLookup) value(key string, version string) (interface{}, error) {
	v := k[key]
	if version == LatestVersion {
		if len(v) > 0 {
			return v[len(v)-1], nil
		}
		return nil, fmt.Errorf("%s not stored", key)
	}
	for _, s := range v {
		if s.version() == version {
			return s, nil
		}
	}
	return nil, fmt.Errorf("key %s at version %s was not found", key, version)
}

// Mock provides a mock implementation of the internal keyvault structure.
// This allows for hermetic testing without using the actual keyvault service.
type Mock struct {
	vaultName string
	secrets   keyLookup
	certs     keyLookup

	mu                sync.Mutex
	getSecretFailures map[string]map[int]bool
	counts            map[string]int
	getCertFailures   map[string]map[int]bool
	certCounts        map[string]int
}

// NewMock creates a new KeyVault service mocker that represents vaultName.  vaultName is the just the name,
// not the fully qualified URL.
func NewMock(vaultName string) (*Mock, error) {
	return &Mock{
		vaultName:         vaultName,
		secrets:           keyLookup{},
		certs:             keyLookup{},
		getSecretFailures: map[string]map[int]bool{},
		getCertFailures:   map[string]map[int]bool{},
		counts:            map[string]int{},
		certCounts:        map[string]int{},
	}, nil
}

// AddSecret adds a secret to the mock. The last secret added for a name will be the
// secret retrieved if using LatestVersion.  Secret.ID can be left blank as it will be
// auto populated. Do not use "/" in name or version to avoid unintentional errors.
func (m *Mock) AddSecret(name, version string, s Secret) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, _ := url.Parse(fmt.Sprintf("https://%s.vault.azure.net/secrets/%s/%s", m.vaultName, name, version))

	var err error
	s.ID, err = urlToID(u)
	if err != nil {
		return err
	}

	if v := m.secrets[name]; v == nil {
		m.secrets[name] = []versioner{s}
		return nil
	}
	m.secrets[name] = append(m.secrets[name], s)
	m.counts[name] = 0
	return nil
}

/*
AddPKCS12 adds a PKCS12 encoded private key/public certificate. This should be the content of
a .pfx file.  If the pfx file content is not password protected, pass the empty string.
You will also need to provide the index in the decoded PEM
blocks to where the public cert is located. If you need help with this, use
PKCS12ToPEM() to decode and look at the order.
Using this will cause a certificate to be created that contains the public cert and a secret accessible via Secret() that contains the private cert (with the same name).

You can convert existing test key pairs with OpenSSL to use in the mock (should be test keys, not real keys):

	openssl pkcs12 -export -out certificate.pfx -inkey privateKey.key -in certificate.crt -certfile more.crt

	Breaking down the command:

	-export -out certificate.pfx – export and save the PFX file as certificate.pfx
	-inkey privateKey.key – use the private key file privateKey.key as the private key to combine with the certificate.
	-in certificate.crt – use certificate.crt as the certificate the private key will be combined with.
	-certfile more.crt – This is optional, this is if you have any additional certificates you would like to include in the PFX file.
*/
func (m *Mock) AddPKCS12(name, version string, pfx []byte, password string, certBlockIndex int) error {
	base64.StdEncoding.EncodeToString(pfx)
	s := Secret{
		Value: base64.StdEncoding.EncodeToString(pfx),
		Attr: SecretAttr{
			Enabled: true,
			Created: time.Now(),
			Updated: time.Now(),
		},
	}

	blocks, err := PKCS12ToPEM(pfx, password)
	if err != nil {
		return err
	}

	c := Certificate{
		Content: blocks[certBlockIndex].Bytes,
		Attr: CertAttr{
			Enabled: true,
			Created: time.Now(),
			Updated: time.Now(),
		},
	}

	if err := m.AddCertificate(name, version, c); err != nil {
		return err
	}

	if err := m.AddSecret(name, version, s); err != nil {
		return err
	}
	return nil
}

// AddCertificate adds a public certificate for a key pair to Keyvault. However, this can only add the public certificate.
// If you want to add a public/private key, use AddPKCS12().
func (m *Mock) AddCertificate(name, version string, c Certificate) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	u, _ := url.Parse(fmt.Sprintf("https://%s.vault.azure.net/secrets/%s/%s", m.vaultName, name, version))

	var err error
	c.ID, err = urlToID(u)
	if err != nil {
		return err
	}

	if v := m.certs[name]; v == nil {
		m.certs[name] = []versioner{c}
		return nil
	}
	m.certs[name] = append(m.certs[name], c)
	m.certCounts[name] = 0
	return nil
}

// GetSecretFailures allows you to set failures that will occur when a numbered attempt to retrieve a named secret occurs.
// This allows users to simulate failures in their tests.
// Setting SetFailures("key", 3, 8, 12) will causes the third, eighth and tweleth
// fetches of "key" to fail with an error.  SetFailures replaces a previous call to SetFailures if called again and
// resets the call count.
func (m *Mock) GetSecretFailures(name string, attempt ...int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	attempts := make(map[int]bool, len(attempt))
	for _, a := range attempt {
		attempts[a] = true
	}

	m.getSecretFailures[name] = attempts
	m.counts[name] = 0
	return
}

// GetSecrets implements keyvault.BaseClient.GetSecret.
func (m *Mock) GetSecret(ctx context.Context, url string, key string, version string) (keyvault.SecretBundle, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	defer func() { m.counts[key]++ }()

	count := m.counts[key]

	if m.getSecretFailures[key][count] {
		return keyvault.SecretBundle{}, fmt.Errorf("error getting secret: testing error handling on attempt %d for key %s", count, key)
	}

	i, err := m.secrets.value(key, version)
	if err != nil {
		return keyvault.SecretBundle{}, err
	}
	s := i.(Secret)

	return s.toBundle(), nil
}

// GetCertificate implements keyvault.BaseClient.GetCertificate.
func (m *Mock) GetCertificate(ctx context.Context, url string, name string, version string) (keyvault.CertificateBundle, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	defer func() { m.certCounts[name]++ }()

	count := m.certCounts[name]

	if m.getCertFailures[name][count] {
		return keyvault.CertificateBundle{}, fmt.Errorf("error getting certificate: testing error handling on attempt %d for name %s", count, name)
	}

	i, err := m.certs.value(name, version)
	if err != nil {
		return keyvault.CertificateBundle{}, err
	}
	c := i.(Certificate)

	return c.toBundle(), nil
}

// PKCS12ToPEM converts the contexts of a pfx file to PEM blocks. You can use the PEM
// blocks to look for the location of your various certificates. This will handle
// base64 standard encoding or non-base64 standard encoded formats. If the pfx data
// has no password, pass an empty string.
func PKCS12ToPEM(pfx []byte, password string) ([]*pem.Block, error) {
	blocks, err := pkcs12.ToPEM(pfx, password)
	if err != nil {
		b, err := base64.StdEncoding.DecodeString(string(pfx))
		if err != nil {
			return nil, fmt.Errorf("pfx(PKCS12) does not look to be decodable. Tried to decode raw []byte and then tried to base64 decode when that failed.  Neither worked.")
		}
		blocks, err = pkcs12.ToPEM(b, password)
		if err != nil {
			return nil, fmt.Errorf("pfx(PKCS12) is not decodable. Would not decode as passed and would not decode after base64 decoding: %s", err)
		}
	}
	return blocks, nil
}
