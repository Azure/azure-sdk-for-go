/*
Package easykey provides access to Azure Keyvault storage in a more Go like interface
than the normal REST to Go conversion. In addition, it provides methos for providing
needed conversions from the internal Keyvault formats for keys to common formats like
PEM and TLS certs needed for Go. To allow for hermetic testing, a Keyvault mock is
provided. See keyvault_test.go for examples of using the mocker.

Create a keyvault without env variables:

	kv, err := New(ctx, "vaultName", Config{ClientID: "client id", TenantID: "tenant id", SubscriptionID: "sub id"})
	if err != nil {
		// Do something
	}

Create a keyvault with env variables (preferred):

	kv, err := New(ctx, "vaultName", Config{})
	if err != nil {
		// Do something
	}

Get a secret example:

	secret, err := kv.Secret(ctx, "yourKey", LatestVersion)
	if err != nil {
		// Do something
	}
	fmt.Println("secret was: %s", secret.Value)

Get a certificate example:

	// This retrieves the public certificate part of the .pfx file uploaded to Keyvault or generated.
	cert, err :=  kv.Certificate(ctx, "certName", LatestVersion)
	if err != nil {
		// Do something
	}

	// This translates the cert given back to a native Go X509 that you can use.
	xCert, err := cert.X509()
	if err != nil {
		// Do something
	}

Get a private key:

	key, err := kv.PrivateKey(ctx, "certName", LatestVersion)
	if err != nil {
		// Do something
	}

Get a TLS cert (from certificate and key):

	tlsConfig, err := kv.TLSCert(ctx, "certName", LatestVersion)
	if err != nil {
		// Do something
	}

+build go1.9
*/

package easykey

//go:generate stringer -type=DeletionRecoveryLevel

// Copyright (c) Microsoft and contributors.  All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//
// See the License for the specific language governing permissions and
// limitations under the License.

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Azure/azure-sdk-for-go/services/keyvault/auth"
	"github.com/Azure/azure-sdk-for-go/services/keyvault/v7.0/keyvault"
)

type vault interface {
	// GetSecret grabs a secret from the keyvault at url for key at version.
	GetSecret(ctx context.Context, url string, key string, version string) (keyvault.SecretBundle, error)

	// GetCertificate grabs a certificate from the keyvault at url for certificate name at version.
	GetCertificate(ctx context.Context, url string, name string, version string) (keyvault.CertificateBundle, error)
}

// ID is a generic keyvault identifier.
type ID struct {
	// Name is the identifier name.
	Name string
	// Version is the identifier version.
	Version string

	url string
}

func (i ID) String() string {
	return i.url
}

func urlToID(u *url.URL) (ID, error) {
	sp := strings.Split(u.Path, "/")[1:]
	if len(sp) != 3 {
		return ID{}, fmt.Errorf("vault has URL that we don't understand: %s", u)
	}
	return ID{Name: sp[1], Version: sp[2], url: u.String()}, nil
}

// DeletionRecoveryLevel indicates what level of recovery is associated with a particular secret.
// Details at: https://docs.microsoft.com/en-us/rest/api/keyvault/getsecretversions/getsecretversions#deletionrecoverylevel
type DeletionRecoveryLevel = keyvault.DeletionRecoveryLevel

const (
	// Purgeable indicates soft-delete is not enabled for this vault. A DELETE operation results in immediate and
	// irreversible data loss.
	Purgeable DeletionRecoveryLevel = "Purgeable"
	// Recoverable indicates soft-delete is enabled for this vault and purge has been disabled. A deleted entity
	// will remain in this state until recovered, or the end of the retention interval.
	Recoverable DeletionRecoveryLevel = "Recoverable"
	// RecoverableProtectedSubscription indicates soft-delete is enabled for this vault, and the subscription is
	// protected against immediate deletion.
	RecoverableProtectedSubscription DeletionRecoveryLevel = "Recoverable+ProtectedSubscription"
	// RecoverablePurgeable indicates soft-delete is enabled for this vault; A privileged user may trigger an
	// immediate, irreversible deletion(purge) of a deleted entity.
	RecoverablePurgeable DeletionRecoveryLevel = "Recoverable+Purgeable"
)

// LatestVersion is used to return the latest version of a key's data when calling Secret().
const LatestVersion = ""

// Config provides for configuration information needed to identify your application to KeyVault.
// To locate these values, please see:
// https://docs.microsoft.com/en-us/azure/azure-resource-manager/resource-group-create-service-principal-portal
// If you pass an empty Config, we will try to read the values from environmental variables.
// The mappings from environmental variables names to Config attributs are listed below.
type Config struct {
	// ClientID (required) is application ID given to your application.  If getting from the environment,
	// the env variable is "AZURE_CLIENT_ID".
	ClientID string
	// ClientSecret (optional) provides the secret that authorizes this application to access the keyvault.
	// If getting from the enivronment, the env variable is "AZURE_CLIENT_SECRET".
	// While this is optional, it is only optional if the environment using keyvault can use Managed Service Identities.
	// These would not be available unless running in Azure.
	ClientSecret string
	// TenantID (required) provides the ID that is used to identify your Active Directory instance that is the authority
	// for your services.  If getting from the enivronment, the env variable is "AZURE_TENANT_ID".
	TenantID string
	// SubscriptionID (required) provides the ID of the subscription that is used for charging your Azure account.
	// If getting from the enivronment, the env variable is "AZURE_SUBSCRIPTION_ID".
	SubscriptionID string
}

// getEnv() is used to populate the Config from the environment in case the Config is a zero value.
func (c *Config) getEnv() {
	c.ClientID = os.Getenv("AZURE_CLIENT_ID")
	c.ClientSecret = os.Getenv("AZURE_CLIENT_SECRET")
	c.TenantID = os.Getenv("AZURE_TENANT_ID")
	c.SubscriptionID = os.Getenv("AZURE_SUBSCRIPTION_ID")
}

// validate validates we have set the necessary config values.  It does not validate if the values are incorrect.
func (c *Config) validate() error {
	switch "" {
	case c.ClientID, c.TenantID, c.SubscriptionID:
		return fmt.Errorf("must have ClientID, TenantID and SubscriptionID set in Config passed to keyvault.New()")
	case c.ClientSecret:
		// TODO(jdoak): This should end up being a warning if we cannot test for MSI compatibility.
		// I have been unable to make MSI work on MSI enabled VMs, but I'm sure this is my fault.
		return fmt.Errorf("ClientSecret was not set in keyvault.Config. This means Managed Service Identity must be usable from the environment")
	}
	return nil
}

// setEnv is used to set our environmental variables to the values in the Config.
func (c *Config) setEnv() {
	os.Setenv("AZURE_CLIENT_ID", c.ClientID)
	os.Setenv("AZURE_TENANT_ID", c.TenantID)
	os.Setenv("AZURE_SUBSCRIPTION_ID", c.SubscriptionID)

	if c.ClientSecret == "" {
		os.Setenv("AZURE_CLIENT_SECRET", c.ClientSecret)
	}
}

// zeroConfig is used to test if a passed Config is empty.
var zeroConfig = Config{}

// Option is an optional argument for New().
type Option func(k *KeyVault)

// MockService instructs the KeyVault client to use the Mock as the KeyVault service.
func MockService(m *Mock) Option {
	return func(k *KeyVault) {
		k.vault = m
	}
}

// KeyVault provides access to the Azure KeyVault service for storing secrets and certificates.
type KeyVault struct {
	vaultName string
	vaultURL  string
	vault     vault // *keyvault.BaseClient implements this.
}

// New is the constructor for a KeyVault.
func New(ctx context.Context, vault string, config Config, options ...Option) (*KeyVault, error) {
	if vault == "" {
		return nil, fmt.Errorf("keyvault.New() cannot have vault argument == empty string")
	}

	client := &KeyVault{
		vaultName: vault,
		vaultURL:  fmt.Sprintf("https://%s.vault.azure.net/", vault),
	}

	for _, o := range options {
		o(client)
	}

	// Indicates we didn't receive the MockService option.
	if client.vault == nil {
		if config == zeroConfig {
			config.getEnv()
			if err := config.validate(); err != nil {
				return nil, err
			}
		}

		// TODO(jdoak): Make sure if config.ClientSecret is "", that we check that MSI is going to work.

		config.setEnv()
		// This sets up some auth stuff we need from our Config into internal settings and
		// some other things keyvault needs.  Unlike our Config, it does not error when things
		// that are needed are not there.
		var err error
		kv := keyvault.New()
		kv.Authorizer, err = auth.NewAuthorizerFromEnvironment()
		if err != nil {
			return nil, fmt.Errorf("could not get Keyvault authorizer: %s", err)
		}
		client.vault = &kv
	}
	return client, nil
}
