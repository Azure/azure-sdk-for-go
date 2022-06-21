//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/crypto"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/base"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Client interacts with Key Vault keys.
type Client struct {
	kvClient *generated.KeyVaultClient
	vaultURL string
}

// ClientOptions are the configurable options for a Client.
type ClientOptions struct {
	azcore.ClientOptions
}

// converts ClientOptions to generated *generated.ConnectionOptions
func (c *ClientOptions) toConnectionOptions() *policy.ClientOptions {
	if c == nil {
		return nil
	}

	return &policy.ClientOptions{
		Logging:          c.Logging,
		Retry:            c.Retry,
		Telemetry:        c.Telemetry,
		Transport:        c.Transport,
		PerCallPolicies:  c.PerCallPolicies,
		PerRetryPolicies: c.PerRetryPolicies,
	}
}

// NewClient constructs a Client that accesses a Key Vault's keys.
func NewClient(vaultURL string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	genOptions := options.toConnectionOptions()

	genOptions.PerRetryPolicies = append(
		genOptions.PerRetryPolicies,
		shared.NewKeyVaultChallengePolicy(credential),
	)

	pl := runtime.NewPipeline(internal.ModuleName, internal.ModuleVersion, runtime.PipelineOptions{}, genOptions)
	return &Client{
		kvClient: generated.NewKeyVaultClient(pl),
		vaultURL: vaultURL,
	}, nil
}

// VaultURL returns the URL for the client's Key Vault.
func (c *Client) VaultURL() string {
	return c.vaultURL
}

// NewCryptoClient creates a new *crypto.Client for the specified key and optional version.
// The created client uses the same vault URL and options as this Client.
func (c *Client) NewCryptoClient(keyName string, keyVersion *string) *crypto.Client {
	keyVer := ""
	if keyVersion != nil {
		keyVer = *keyVersion
	}
	return &crypto.Client{CryptoClient: base.NewCryptoClient(c.vaultURL, keyName, keyVer, c.kvClient.Pipeline())}
}

// CreateKeyOptions contains optional parameters for CreateKey.
type CreateKeyOptions struct {
	// Curve is the elliptic curve name. For valid values, see PossibleCurveNameValues.
	Curve *CurveName

	// Properties is the key's management properties.
	Properties *Properties

	// Operations are the operations Key Vault will allow for the key.
	Operations []*Operation

	// ReleasePolicy specifies conditions under which the key can be exported.
	ReleasePolicy *ReleasePolicy

	// Size is the key size in bits. For example: 2048, 3072, or 4096 for RSA.
	Size *int32

	// PublicExponent is the public exponent of an RSA key.
	PublicExponent *int32

	// Tags is application specific metadata in the form of key-value pairs.
	Tags map[string]*string
}

// convert CreateKeyOptions to *generated.KeyVaultClientCreateKeyOptions
func (c *CreateKeyOptions) toGenerated() *generated.KeyVaultClientCreateKeyOptions {
	return &generated.KeyVaultClientCreateKeyOptions{}
}

// convert CreateKeyOptions to generated.KeyCreateParameters
func (c *CreateKeyOptions) toKeyCreateParameters(keyType KeyType) generated.KeyCreateParameters {
	var attribs *generated.KeyAttributes
	if c.Properties != nil {
		attribs = c.Properties.toGenerated()
	}

	var ops []*generated.JSONWebKeyOperation
	if c.Operations != nil {
		ops = make([]*generated.JSONWebKeyOperation, len(c.Operations))
		for i, o := range c.Operations {
			ops[i] = (*generated.JSONWebKeyOperation)(o)
		}
	}

	return generated.KeyCreateParameters{
		Kty:            keyType.toGenerated(),
		Curve:          (*generated.JSONWebKeyCurveName)(c.Curve),
		KeyAttributes:  attribs,
		KeyOps:         ops,
		KeySize:        c.Size,
		PublicExponent: c.PublicExponent,
		Tags:           c.Tags,
		ReleasePolicy:  c.ReleasePolicy.toGenerated(),
	}
}

// CreateKeyResponse is returned by CreateKey.
type CreateKeyResponse struct {
	Key
}

// creates CreateKeyResponse from generated.KeyVaultClient.CreateKeyResponse
func createKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) CreateKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return CreateKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags, g.ReleasePolicy),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// CreateKey creates a key. If the named key already exists, it creates a new version of the key. This method can create
// a key of any type. CreateRSAKey, CreateECKey, and CreateOctKey are more convenient for creating specific key types.
//  Pass nil for options to accept default values.
func (c *Client) CreateKey(ctx context.Context, name string, keyType KeyType, options *CreateKeyOptions) (CreateKeyResponse, error) {
	if options == nil {
		options = &CreateKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultURL, name, options.toKeyCreateParameters(keyType), options.toGenerated())
	if err != nil {
		return CreateKeyResponse{}, err
	}

	return createKeyResponseFromGenerated(resp), nil
}

// CreateECKeyOptions contains optional parameters for CreateECKey
type CreateECKeyOptions struct {
	// Curve is the elliptic curve name. For valid values, see PossibleCurveNameValues.
	Curve *CurveName

	// Tags is application specific metadata in the form of key-value pairs.
	Tags map[string]*string

	// HardwareProtected determines whether the key is is created in a hardware security module (HSM).
	HardwareProtected *bool

	// Properties is the key's management properties.
	Properties *Properties

	// Operations are the operations Key Vault will allow for the key.
	Operations []*Operation

	// ReleasePolicy specifies conditions under which the key can be exported
	ReleasePolicy *ReleasePolicy
}

// convert CreateECKeyOptions to generated.KeyCreateParameters
func (c *CreateECKeyOptions) toKeyCreateParameters(keyType KeyType) generated.KeyCreateParameters {
	var keyOps []*generated.JSONWebKeyOperation
	if c.Operations != nil {
		keyOps = make([]*generated.JSONWebKeyOperation, len(c.Operations))
		for i, k := range c.Operations {
			keyOps[i] = (*generated.JSONWebKeyOperation)(k)
		}
	}
	return generated.KeyCreateParameters{
		Kty:           keyType.toGenerated(),
		Curve:         (*generated.JSONWebKeyCurveName)(c.Curve),
		Tags:          c.Tags,
		KeyOps:        keyOps,
		ReleasePolicy: c.ReleasePolicy.toGenerated(),
		KeyAttributes: c.Properties.toGenerated(),
	}
}

// CreateECKeyResponse is returned by CreateECKey.
type CreateECKeyResponse struct {
	Key
}

// convert the generated.KeyVaultClientCreateKeyResponse to CreateECKeyResponse
func createECKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) CreateECKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return CreateECKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags, g.ReleasePolicy),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// CreateECKey creates a new elliptic curve key. If the named key already exists, this creates a new version of the key.
// Pass nil for options to accept default values.
func (c *Client) CreateECKey(ctx context.Context, name string, options *CreateECKeyOptions) (CreateECKeyResponse, error) {
	keyType := KeyTypeEC

	if options == nil {
		options = &CreateECKeyOptions{}
	}
	if options.HardwareProtected != nil && *options.HardwareProtected {
		keyType = KeyTypeECHSM
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultURL, name, options.toKeyCreateParameters(keyType), &generated.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return CreateECKeyResponse{}, err
	}

	return createECKeyResponseFromGenerated(resp), nil
}

// CreateOctKeyOptions contains optional parameters for CreateOctKey
type CreateOctKeyOptions struct {
	// HardwareProtected determines whether the key is is created in a hardware security module (HSM).
	HardwareProtected *bool

	// Size is the key size in bits. For example: 128, 192 or 256.
	Size *int32

	// Properties is the key's management properties.
	Properties *Properties

	// Operations are the operations Key Vault will allow for the key.
	Operations []*Operation

	// ReleasePolicy specifies conditions under which the key can be exported
	ReleasePolicy *ReleasePolicy

	// Tags is application specific metadata in the form of key-value pairs.
	Tags map[string]*string
}

// conver the CreateOctKeyOptions to generated.KeyCreateParameters
func (c *CreateOctKeyOptions) toKeyCreateParameters(keyType KeyType) generated.KeyCreateParameters {
	var keyOps []*generated.JSONWebKeyOperation
	if c.Operations != nil {
		keyOps = make([]*generated.JSONWebKeyOperation, len(c.Operations))
		for i, k := range c.Operations {
			keyOps[i] = (*generated.JSONWebKeyOperation)(k)
		}
	}
	return generated.KeyCreateParameters{
		Kty:           keyType.toGenerated(),
		KeySize:       c.Size,
		Tags:          c.Tags,
		ReleasePolicy: c.ReleasePolicy.toGenerated(),
		KeyAttributes: c.Properties.toGenerated(),
		KeyOps:        keyOps,
	}
}

// CreateOctKeyResponse is returned by CreateOctKey.
type CreateOctKeyResponse struct {
	Key
}

// convert generated response to CreateOctKeyResponse
func createOctKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) CreateOctKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return CreateOctKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags, g.ReleasePolicy),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// CreateOctKey creates a new AES key. If the named key already exists, this creates a new version of the key. Pass nil for options to accept default values.
func (c *Client) CreateOctKey(ctx context.Context, name string, options *CreateOctKeyOptions) (CreateOctKeyResponse, error) {
	keyType := KeyTypeOctHSM

	if options != nil && options.HardwareProtected != nil && !*options.HardwareProtected {
		keyType = KeyTypeOct
	} else if options == nil {
		options = &CreateOctKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultURL, name, options.toKeyCreateParameters(keyType), &generated.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return CreateOctKeyResponse{}, err
	}

	return createOctKeyResponseFromGenerated(resp), nil
}

// CreateRSAKeyOptions contains optional parameters for CreateRSAKey.
type CreateRSAKeyOptions struct {
	// HardwareProtected determines whether the key is is created in a hardware security module (HSM).
	HardwareProtected *bool

	// Size is the key size in bits. For example: 2048, 3072, or 4096.
	Size *int32

	// PublicExponent is the key's public exponent.
	PublicExponent *int32

	// Tags is application specific metadata in the form of key-value pairs.
	Tags map[string]*string

	// Properties is the key's management properties.
	Properties *Properties

	// Operations are the operations Key Vault will allow for the key.
	Operations []*Operation

	// ReleasePolicy specifies conditions under which the key can be exported
	ReleasePolicy *ReleasePolicy
}

// convert CreateRSAKeyOptions to generated.KeyCreateParameters
func (c CreateRSAKeyOptions) toKeyCreateParameters(k KeyType) generated.KeyCreateParameters {
	var keyOps []*generated.JSONWebKeyOperation
	if c.Operations != nil {
		keyOps = make([]*generated.JSONWebKeyOperation, len(c.Operations))
		for i, k := range c.Operations {
			keyOps[i] = (*generated.JSONWebKeyOperation)(k)
		}
	}
	return generated.KeyCreateParameters{
		Kty:            k.toGenerated(),
		KeySize:        c.Size,
		PublicExponent: c.PublicExponent,
		Tags:           c.Tags,
		KeyAttributes:  c.Properties.toGenerated(),
		KeyOps:         keyOps,
		ReleasePolicy:  c.ReleasePolicy.toGenerated(),
	}
}

// CreateRSAKeyResponse is returned by CreateRSAKey.
type CreateRSAKeyResponse struct {
	Key
}

// convert internal response to CreateRSAKeyResponse
func createRSAKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) CreateRSAKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return CreateRSAKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags, g.ReleasePolicy),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// CreateRSAKey can be used to create a new RSA key in Azure Key Vault. If the named key already
// exists, Azure Key Vault creates a new version of the key. RSA keys can be created in Standard or
// Premium SKU vaults, RSAHSM can be created in Premium SKU vaults or Managed HSMs.
// It requires the keys/create permission. Pass nil for options to accept default values.
func (c *Client) CreateRSAKey(ctx context.Context, name string, options *CreateRSAKeyOptions) (CreateRSAKeyResponse, error) {
	keyType := KeyTypeRSA

	if options == nil {
		options = &CreateRSAKeyOptions{}
	}
	if options.HardwareProtected != nil && *options.HardwareProtected {
		keyType = KeyTypeRSAHSM
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultURL, name, options.toKeyCreateParameters(keyType), &generated.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return CreateRSAKeyResponse{}, err
	}

	return createRSAKeyResponseFromGenerated(resp), nil
}

// ListPropertiesOfKeysOptions contains optional parameters for ListKeys
type ListPropertiesOfKeysOptions struct {
	// placeholder for future optional parameters
}

// ListPropertiesOfKeysResponse contains a page of key properties.
type ListPropertiesOfKeysResponse struct {
	// NextLink is the URL to get the next page.
	NextLink *string

	// Keys is the page's content.
	Keys []*KeyItem
}

// convert internal Response to ListKeysPage
func listKeysPageFromGenerated(i generated.KeyVaultClientGetKeysResponse) ListPropertiesOfKeysResponse {
	var keys []*KeyItem
	for _, k := range i.Value {
		keys = append(keys, keyItemFromGenerated(k))
	}
	return ListPropertiesOfKeysResponse{
		NextLink: i.NextLink,
		Keys:     keys,
	}
}

// NewListPropertiesOfKeysPager retrieves a list of the keys in the Key Vault as JSON Web Key structures that contain the
// public part of a stored key. The LIST operation is applicable to all key types, however only the
// base key identifier, attributes, and tags are provided in the response. Individual versions of a
// key are not listed in the response. This operation requires the keys/list permission.
func (c *Client) NewListPropertiesOfKeysPager(options *ListPropertiesOfKeysOptions) *runtime.Pager[ListPropertiesOfKeysResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListPropertiesOfKeysResponse]{
		More: func(page ListPropertiesOfKeysResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ListPropertiesOfKeysResponse) (ListPropertiesOfKeysResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = c.kvClient.GetKeysCreateRequest(ctx, c.vaultURL, &generated.KeyVaultClientGetKeysOptions{})
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ListPropertiesOfKeysResponse{}, err
			}
			resp, err := c.kvClient.Pipeline().Do(req)
			if err != nil {
				return ListPropertiesOfKeysResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListPropertiesOfKeysResponse{}, runtime.NewResponseError(resp)
			}
			genResp, err := c.kvClient.GetKeysHandleResponse(resp)
			if err != nil {
				return ListPropertiesOfKeysResponse{}, err
			}
			return listKeysPageFromGenerated(genResp), nil
		},
	})
}

// GetKeyOptions contains the options for the Client.GetKey method
type GetKeyOptions struct {
	Version string
}

// GetKeyResponse is returned by GetResponse.
type GetKeyResponse struct {
	Key
}

// convert internal response to GetKeyResponse
func getKeyResponseFromGenerated(g generated.KeyVaultClientGetKeyResponse) GetKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return GetKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags, g.ReleasePolicy),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// GetKey is used to retrieve the content for any single Key. If the requested key is symmetric, then
// no key material is released in the response. This operation requires the keys/get permission.
// Pass nil for options to accept default values.
func (c *Client) GetKey(ctx context.Context, name string, options *GetKeyOptions) (GetKeyResponse, error) {
	if options == nil {
		options = &GetKeyOptions{}
	}

	resp, err := c.kvClient.GetKey(ctx, c.vaultURL, name, options.Version, &generated.KeyVaultClientGetKeyOptions{})
	if err != nil {
		return GetKeyResponse{}, err
	}

	return getKeyResponseFromGenerated(resp), err
}

// GetDeletedKeyOptions contains optional parameters for GetDeletedKey
type GetDeletedKeyOptions struct {
	// placeholder for future optional parameters
}

// convert the GetDeletedKeyOptions to the internal representation
func (g GetDeletedKeyOptions) toGenerated() *generated.KeyVaultClientGetDeletedKeyOptions {
	return &generated.KeyVaultClientGetDeletedKeyOptions{}
}

// GetDeletedKeyResponse contains the response from a Client.GetDeletedKey
type GetDeletedKeyResponse struct {
	DeletedKey
}

// convert generated response to GetDeletedKeyResponse
func getDeletedKeyResponseFromGenerated(g generated.KeyVaultClientGetDeletedKeyResponse) GetDeletedKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return GetDeletedKeyResponse{
		DeletedKey: DeletedKey{
			Properties:         keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags, g.ReleasePolicy),
			Key:                jsonWebKeyFromGenerated(g.Key),
			RecoveryID:         g.RecoveryID,
			DeletedOn:          g.DeletedDate,
			ScheduledPurgeDate: g.ScheduledPurgeDate,
		},
	}
}

// GetDeletedKey is used to retrieve information about a deleted key. This operation is only
// applicable for soft-delete enabled vaults. While the operation can be invoked on any vault,
// it will return an error if invoked on a non soft-delete enabled vault. This operation requires
// the keys/get permission. Pass nil for options to accept default values.
func (c *Client) GetDeletedKey(ctx context.Context, name string, options *GetDeletedKeyOptions) (GetDeletedKeyResponse, error) {
	if options == nil {
		options = &GetDeletedKeyOptions{}
	}

	resp, err := c.kvClient.GetDeletedKey(ctx, c.vaultURL, name, options.toGenerated())
	if err != nil {
		return GetDeletedKeyResponse{}, err
	}

	return getDeletedKeyResponseFromGenerated(resp), nil
}

// PurgeDeletedKeyOptions contains optional parameters for PurgeDeletedKey.
type PurgeDeletedKeyOptions struct {
	// placeholder for future optional parameters
}

// convert options to internal options
func (p *PurgeDeletedKeyOptions) toGenerated() *generated.KeyVaultClientPurgeDeletedKeyOptions {
	return &generated.KeyVaultClientPurgeDeletedKeyOptions{}
}

// PurgeDeletedKeyResponse is returned by PurgeDeletedKey.
type PurgeDeletedKeyResponse struct {
	// placeholder for future response values
}

// Converts the generated response to the publicly exposed version.
func purgeDeletedKeyResponseFromGenerated(i generated.KeyVaultClientPurgeDeletedKeyResponse) PurgeDeletedKeyResponse {
	return PurgeDeletedKeyResponse{}
}

// PurgeDeletedKey permanently deletes a deleted key. Key Vault may require several seconds to finish purging the key after this
// method returns. Pass nil for options to accept default values.
func (c *Client) PurgeDeletedKey(ctx context.Context, name string, options *PurgeDeletedKeyOptions) (PurgeDeletedKeyResponse, error) {
	if options == nil {
		options = &PurgeDeletedKeyOptions{}
	}
	resp, err := c.kvClient.PurgeDeletedKey(ctx, c.vaultURL, name, options.toGenerated())
	return purgeDeletedKeyResponseFromGenerated(resp), err
}

// DeleteKeyResponse contains the response for a Client.BeginDeleteKey operation.
type DeleteKeyResponse struct {
	DeletedKey
}

// BeginDeleteKeyOptions contains optional parameters for BeginDeleteKey.
type BeginDeleteKeyOptions struct {
	// ResumeToken is a string to rehydrate a poller for an operation that has already begun.
	ResumeToken string
}

// BeginDeleteKey deletes all versions of a key. It returns a Poller that enables waiting for the deletion to finish.
// Pass nil for options to accept default values.
func (c *Client) BeginDeleteKey(ctx context.Context, name string, options *BeginDeleteKeyOptions) (*runtime.Poller[DeleteKeyResponse], error) {
	if options == nil {
		options = &BeginDeleteKeyOptions{}
	}

	handler := beginDeleteKeyPoller{
		poll: func(ctx context.Context) (*http.Response, error) {
			req, err := c.kvClient.GetDeletedKeyCreateRequest(ctx, c.vaultURL, name, nil)
			if err != nil {
				return nil, err
			}
			return c.kvClient.Pipeline().Do(req)
		},
	}

	if options.ResumeToken != "" {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, c.kvClient.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[DeleteKeyResponse]{
			Handler: &handler,
		})
	}

	var rawResp *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &rawResp)
	if _, err := c.kvClient.DeleteKey(ctx, c.vaultURL, name, nil); err != nil {
		return nil, err
	}

	return runtime.NewPoller(rawResp, c.kvClient.Pipeline(), &runtime.NewPollerOptions[DeleteKeyResponse]{
		Handler: &handler,
	})
}

// BackupKeyOptions contains optional parameters for BackupKey
type BackupKeyOptions struct {
	// placeholder for future optional parameters
}

// convert Options to generated version
func (b BackupKeyOptions) toGenerated() *generated.KeyVaultClientBackupKeyOptions {
	return &generated.KeyVaultClientBackupKeyOptions{}
}

// BackupKeyResponse contains the response from the Client.BackupKey method
type BackupKeyResponse struct {
	// READ-ONLY; The backup blob containing the backed up key.
	Value []byte
}

// convert internal reponse to BackupKeyResponse
func backupKeyResponseFromGenerated(i generated.KeyVaultClientBackupKeyResponse) BackupKeyResponse {
	return BackupKeyResponse{
		Value: i.Value,
	}
}

// BackupKey exports all versions of a key from Azure Key Vault in a protected form.
//
// Note that this operation does NOT return key material in a form that can be used outside the
// Azure Key Vault system. The returned key material is either protected to an Azure Key Vault
// HSM or to Azure Key Vault itself. The intent of this operation is to allow a client to GENERATE
// a key in one Azure Key Vault instance, BACKUP the key, and then RESTORE it into another Azure
// Key Vault instance. BACKUP / RESTORE can be performed within geographical boundaries only; a
//  BACKUP from one geographical area cannot be restored to another geographical area. For example,
//  a backup from the US geographical area cannot be restored in an EU geographical area.
// Pass nil for options to accept default values.
func (c *Client) BackupKey(ctx context.Context, name string, options *BackupKeyOptions) (BackupKeyResponse, error) {
	if options == nil {
		options = &BackupKeyOptions{}
	}

	resp, err := c.kvClient.BackupKey(ctx, c.vaultURL, name, options.toGenerated())
	if err != nil {
		return BackupKeyResponse{}, err
	}

	return backupKeyResponseFromGenerated(resp), nil
}

// BeginRecoverDeletedKeyOptions contains the optional parameters for the Client.BeginRecoverDeletedKey operation
type BeginRecoverDeletedKeyOptions struct {
	// ResumeToken returns a string for creating a new poller to begin polling again
	ResumeToken string
}

// RecoverDeletedKeyResponse is the response object for the Client.RecoverDeletedKey operation.
type RecoverDeletedKeyResponse struct {
	Key
}

// change recover deleted key reponse to the generated version.
func recoverDeletedKeyResponseFromGenerated(g generated.KeyVaultClientRecoverDeletedKeyResponse) RecoverDeletedKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return RecoverDeletedKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags, g.ReleasePolicy),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// BeginRecoverDeletedKey recovers a deleted key to its latest version. Pass nil for options to accept default values.
func (c *Client) BeginRecoverDeletedKey(ctx context.Context, name string, options *BeginRecoverDeletedKeyOptions) (*runtime.Poller[RecoverDeletedKeyResponse], error) {
	if options == nil {
		options = &BeginRecoverDeletedKeyOptions{}
	}

	handler := beginRecoverDeletedKeyPoller{
		poll: func(ctx context.Context) (*http.Response, error) {
			req, err := c.kvClient.GetKeyCreateRequest(ctx, c.vaultURL, name, "", nil)
			if err != nil {
				return nil, err
			}
			return c.kvClient.Pipeline().Do(req)
		},
	}

	if options.ResumeToken != "" {
		return runtime.NewPollerFromResumeToken(options.ResumeToken, c.kvClient.Pipeline(), &runtime.NewPollerFromResumeTokenOptions[RecoverDeletedKeyResponse]{
			Handler: &handler,
		})
	}

	var rawResp *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &rawResp)
	if _, err := c.kvClient.RecoverDeletedKey(ctx, c.vaultURL, name, nil); err != nil {
		return nil, err
	}

	return runtime.NewPoller(rawResp, c.kvClient.Pipeline(), &runtime.NewPollerOptions[RecoverDeletedKeyResponse]{
		Handler: &handler,
	})
}

// UpdateKeyPropertiesOptions contains optional parameters for UpdateKeyProperties
type UpdateKeyPropertiesOptions struct {
	// Operations are the operations Key Vault will allow for the key.
	Operations []*Operation
}

// UpdateKeyPropertiesResponse is returned by UpdateKeyProperties.
type UpdateKeyPropertiesResponse struct {
	Key
}

// convert the internal response to UpdateKeyPropertiesResponse
func updateKeyPropertiesFromGenerated(g generated.KeyVaultClientUpdateKeyResponse) UpdateKeyPropertiesResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return UpdateKeyPropertiesResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags, g.ReleasePolicy),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// UpdateKeyProperties updates the management properties of a key, but not its cryptographic material.
// Pass nil for options to accept default values.
func (c *Client) UpdateKeyProperties(ctx context.Context, properties Properties, options *UpdateKeyPropertiesOptions) (UpdateKeyPropertiesResponse, error) {
	if options == nil {
		options = &UpdateKeyPropertiesOptions{}
	}
	name, version := "", ""
	if properties.Name != nil {
		name = *properties.Name
	}
	if properties.Version != nil {
		version = *properties.Version
	}
	params := generated.KeyUpdateParameters{
		KeyAttributes: properties.toGenerated(),
		ReleasePolicy: properties.ReleasePolicy.toGenerated(),
		Tags:          properties.Tags,
	}
	if options.Operations != nil {
		params.KeyOps = make([]*generated.JSONWebKeyOperation, len(options.Operations))
		for i, op := range options.Operations {
			params.KeyOps[i] = (*generated.JSONWebKeyOperation)(op)
		}
	}
	resp, err := c.kvClient.UpdateKey(ctx, c.vaultURL, name, version, params, nil)
	if err != nil {
		return UpdateKeyPropertiesResponse{}, err
	}

	return updateKeyPropertiesFromGenerated(resp), nil
}

// ListDeletedKeysResponse holds the data for a single page.
type ListDeletedKeysResponse struct {
	// NextLink is the URL to get the next page.
	NextLink *string

	// DeletedKeys is the page's content.
	DeletedKeys []*DeletedKeyItem
}

// ListDeletedKeysOptions contains optional parameters for NewListDeletedKeysPager.
type ListDeletedKeysOptions struct {
	// placeholder for future optional parameters
}

// Convert publicly exposed options to the generated version.a
func (l *ListDeletedKeysOptions) toGenerated() *generated.KeyVaultClientGetDeletedKeysOptions {
	return &generated.KeyVaultClientGetDeletedKeysOptions{}
}

// NewListDeletedKeysPager creates a pager that lists deleted keys. Pass nil for options to accept default values.
func (c *Client) NewListDeletedKeysPager(options *ListDeletedKeysOptions) *runtime.Pager[ListDeletedKeysResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListDeletedKeysResponse]{
		More: func(page ListDeletedKeysResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ListDeletedKeysResponse) (ListDeletedKeysResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = c.kvClient.GetDeletedKeysCreateRequest(ctx, c.vaultURL, options.toGenerated())
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ListDeletedKeysResponse{}, err
			}
			resp, err := c.kvClient.Pipeline().Do(req)
			if err != nil {
				return ListDeletedKeysResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListDeletedKeysResponse{}, runtime.NewResponseError(resp)
			}
			genResp, err := c.kvClient.GetDeletedKeysHandleResponse(resp)
			if err != nil {
				return ListDeletedKeysResponse{}, runtime.NewResponseError(resp)
			}

			var values []*DeletedKeyItem
			for _, d := range genResp.Value {
				values = append(values, deletedKeyItemFromGenerated(d))
			}

			return ListDeletedKeysResponse{
				NextLink:    genResp.NextLink,
				DeletedKeys: values,
			}, nil
		},
	})
}

// ListPropertiesOfKeyVersionsOptions contains optional parameters for NewListPropertiesOfKeyVersionsPager.
type ListPropertiesOfKeyVersionsOptions struct {
	// placeholder for future optional parameters
}

// convert the public ListKeyVersionsOptions to the generated version
func (l *ListPropertiesOfKeyVersionsOptions) toGenerated() *generated.KeyVaultClientGetKeyVersionsOptions {
	if l == nil {
		return &generated.KeyVaultClientGetKeyVersionsOptions{}
	}
	return &generated.KeyVaultClientGetKeyVersionsOptions{}
}

// ListPropertiesOfKeyVersionsResponse contains a page of key versions.
type ListPropertiesOfKeyVersionsResponse struct {
	// NextLink is the URL to get the next page.
	NextLink *string

	// Keys is the page's content.
	Keys []*KeyItem
}

// create ListKeysPage from generated pager
func listKeyVersionsPageFromGenerated(i generated.KeyVaultClientGetKeyVersionsResponse) ListPropertiesOfKeyVersionsResponse {
	var keys []*KeyItem
	for _, s := range i.Value {
		if s != nil {
			keys = append(keys, keyItemFromGenerated(s))
		}
	}
	return ListPropertiesOfKeyVersionsResponse{
		NextLink: i.NextLink,
		Keys:     keys,
	}
}

// NewListPropertiesOfKeyVersionsPager creates a pager that lists properties of a key's versions, not including key material.
func (c *Client) NewListPropertiesOfKeyVersionsPager(keyName string, options *ListPropertiesOfKeyVersionsOptions) *runtime.Pager[ListPropertiesOfKeyVersionsResponse] {
	return runtime.NewPager(runtime.PagingHandler[ListPropertiesOfKeyVersionsResponse]{
		More: func(page ListPropertiesOfKeyVersionsResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *ListPropertiesOfKeyVersionsResponse) (ListPropertiesOfKeyVersionsResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = c.kvClient.GetKeyVersionsCreateRequest(ctx, c.vaultURL, keyName, options.toGenerated())
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return ListPropertiesOfKeyVersionsResponse{}, err
			}
			resp, err := c.kvClient.Pipeline().Do(req)
			if err != nil {
				return ListPropertiesOfKeyVersionsResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return ListPropertiesOfKeyVersionsResponse{}, runtime.NewResponseError(resp)
			}
			genResp, err := c.kvClient.GetKeyVersionsHandleResponse(resp)
			if err != nil {
				return ListPropertiesOfKeyVersionsResponse{}, runtime.NewResponseError(resp)
			}
			return listKeyVersionsPageFromGenerated(genResp), nil
		},
	})
}

// RestoreKeyBackupOptions contains optional parameters for RestoreKey.
type RestoreKeyBackupOptions struct {
	// placeholder for future optional parameters
}

func (r RestoreKeyBackupOptions) toGenerated() *generated.KeyVaultClientRestoreKeyOptions {
	return &generated.KeyVaultClientRestoreKeyOptions{}
}

// RestoreKeyBackupResponse is returned by RestoreKeyBackup.
type RestoreKeyBackupResponse struct {
	Key
}

// converts the generated response to the publicly exposed version.
func restoreKeyBackupResponseFromGenerated(g generated.KeyVaultClientRestoreKeyResponse) RestoreKeyBackupResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return RestoreKeyBackupResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags, g.ReleasePolicy),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// RestoreKeyBackup restores all versions of a backed up key to the vault. The keyBackup parameter is the bytes of a key backup as returned by BackupKey.
// Pass nil for options to accept default values.
func (c *Client) RestoreKeyBackup(ctx context.Context, keyBackup []byte, options *RestoreKeyBackupOptions) (RestoreKeyBackupResponse, error) {
	if options == nil {
		options = &RestoreKeyBackupOptions{}
	}

	resp, err := c.kvClient.RestoreKey(ctx, c.vaultURL, generated.KeyRestoreParameters{KeyBundleBackup: keyBackup}, options.toGenerated())
	if err != nil {
		return RestoreKeyBackupResponse{}, err
	}

	return restoreKeyBackupResponseFromGenerated(resp), nil
}

// ImportKeyOptions contains optional parameters for ImportKeyOptions.
type ImportKeyOptions struct {
	// HardwareProtected determines whether Key Vault protects the imported key with an HSM.
	HardwareProtected *bool

	// Properties is the properties of the key.
	Properties *Properties
}

func (i ImportKeyOptions) toImportKeyParameters(key JSONWebKey) generated.KeyImportParameters {
	var attribs *generated.KeyAttributes
	if i.Properties != nil {
		attribs = i.Properties.toGenerated()
	}
	return generated.KeyImportParameters{
		Key:           key.toGenerated(),
		Hsm:           i.HardwareProtected,
		KeyAttributes: attribs,
	}
}

// ImportKeyResponse is returned by ImportKey.
type ImportKeyResponse struct {
	Key
}

// convert the generated response to the ImportKeyResponse
func importKeyResponseFromGenerated(g generated.KeyVaultClientImportKeyResponse) ImportKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return ImportKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags, g.ReleasePolicy),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// ImportKey imports a key into the vault. If the named key already exists, this creates a new version of the key. Pass nil for options to accept default values.
func (c *Client) ImportKey(ctx context.Context, name string, key JSONWebKey, options *ImportKeyOptions) (ImportKeyResponse, error) {
	if options == nil {
		options = &ImportKeyOptions{}
	}

	resp, err := c.kvClient.ImportKey(ctx, c.vaultURL, name, options.toImportKeyParameters(key), &generated.KeyVaultClientImportKeyOptions{})
	if err != nil {
		return ImportKeyResponse{}, err
	}

	return importKeyResponseFromGenerated(resp), nil
}

// GetRandomBytesOptions contains optional parameters for GetRandomBytes.
type GetRandomBytesOptions struct {
	// placeholder for future optional parameters
}

func (g GetRandomBytesOptions) toGenerated() *generated.KeyVaultClientGetRandomBytesOptions {
	return &generated.KeyVaultClientGetRandomBytesOptions{}
}

// GetRandomBytesResponse is returned by GetRandomBytes.
type GetRandomBytesResponse struct {
	// Value is the random bytes.
	Value []byte
}

// GetRandomBytes gets the requested number of random bytes from Azure Managed HSM. Pass nil for options to accept default values.
func (c *Client) GetRandomBytes(ctx context.Context, count *int32, options *GetRandomBytesOptions) (GetRandomBytesResponse, error) {
	if options == nil {
		options = &GetRandomBytesOptions{}
	}

	resp, err := c.kvClient.GetRandomBytes(
		ctx,
		c.vaultURL,
		generated.GetRandomBytesRequest{Count: count},
		options.toGenerated(),
	)

	if err != nil {
		return GetRandomBytesResponse{}, err
	}

	return GetRandomBytesResponse{
		Value: resp.Value,
	}, nil
}

// RotateKeyOptions contains optional parameters for RotateKey.
type RotateKeyOptions struct {
	// placeholder for future optional parameters
}

func (r RotateKeyOptions) toGenerated() *generated.KeyVaultClientRotateKeyOptions {
	return &generated.KeyVaultClientRotateKeyOptions{}
}

// RotateKeyResponse is returned by RotateKey.
type RotateKeyResponse struct {
	Key
}

// RotateKey rotates a key based on the key's rotation policy, creating a new version in the specified key. Pass nil for options to accept default values.
func (c *Client) RotateKey(ctx context.Context, keyName string, options *RotateKeyOptions) (RotateKeyResponse, error) {
	if options == nil {
		options = &RotateKeyOptions{}
	}

	resp, err := c.kvClient.RotateKey(
		ctx,
		c.vaultURL,
		keyName,
		options.toGenerated(),
	)
	if err != nil {
		return RotateKeyResponse{}, err
	}

	vaultURL, name, version := shared.ParseID(resp.Key.Kid)
	return RotateKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(resp.Attributes, resp.Key.Kid, name, version, resp.Managed, vaultURL, resp.Tags, resp.ReleasePolicy),
			JSONWebKey: jsonWebKeyFromGenerated(resp.Key),
			ID:         resp.Key.Kid,
			Name:       name,
		},
	}, nil
}

// GetKeyRotationPolicyOptions contains optional parameters for GetKeyRotationPolicy.
type GetKeyRotationPolicyOptions struct {
	// placeholder for future optional parameters
}

func (g GetKeyRotationPolicyOptions) toGenerated() *generated.KeyVaultClientGetKeyRotationPolicyOptions {
	return &generated.KeyVaultClientGetKeyRotationPolicyOptions{}
}

// GetKeyRotationPolicyResponse is returned by GetKeyRotationPolicy.
type GetKeyRotationPolicyResponse struct {
	RotationPolicy
}

func getKeyRotationPolicyResponseFromGenerated(i generated.KeyVaultClientGetKeyRotationPolicyResponse) GetKeyRotationPolicyResponse {
	var acts []*LifetimeActions
	for _, a := range i.LifetimeActions {
		acts = append(acts, lifetimeActionsFromGenerated(a))
	}
	var attribs *RotationPolicyAttributes
	if i.Attributes != nil {
		attribs = &RotationPolicyAttributes{
			ExpiresIn: i.Attributes.ExpiryTime,
			CreatedOn: i.Attributes.Created,
			UpdatedOn: i.Attributes.Updated,
		}
	}
	return GetKeyRotationPolicyResponse{
		RotationPolicy: RotationPolicy{
			ID:              i.ID,
			LifetimeActions: acts,
			Attributes:      attribs,
		},
	}
}

// GetKeyRotationPolicy gets the specified key's rotation policy. Pass nil for options to accept default values.
func (c *Client) GetKeyRotationPolicy(ctx context.Context, keyName string, options *GetKeyRotationPolicyOptions) (GetKeyRotationPolicyResponse, error) {
	if options == nil {
		options = &GetKeyRotationPolicyOptions{}
	}

	resp, err := c.kvClient.GetKeyRotationPolicy(
		ctx,
		c.vaultURL,
		keyName,
		options.toGenerated(),
	)
	if err != nil {
		return GetKeyRotationPolicyResponse{}, err
	}

	return getKeyRotationPolicyResponseFromGenerated(resp), nil
}

// ReleaseKeyOptions contains optional parameters for Client.ReleaseKey.
type ReleaseKeyOptions struct {
	// Version is the version of the key to release
	Version *string

	// Algorithm is the encryption algorithm used to protected exported key material.
	Algorithm *ExportEncryptionAlg

	// Nonce is client-provided nonce for freshness.
	Nonce *string
}

// ReleaseKeyResponse is returned by ReleaseKey.
type ReleaseKeyResponse struct {
	// Value is a signed token containing the released key.
	Value *string
}

// ReleaseKey is applicable to all key types. The target key must be exportable. Pass nil for options to accept default values.
func (c *Client) ReleaseKey(ctx context.Context, name string, targetAttestationToken string, options *ReleaseKeyOptions) (ReleaseKeyResponse, error) {
	if options == nil {
		options = &ReleaseKeyOptions{}
	}
	version := ""
	if options.Version != nil {
		version = *options.Version
	}
	resp, err := c.kvClient.Release(
		ctx,
		c.vaultURL,
		name,
		version,
		generated.KeyReleaseParameters{
			TargetAttestationToken: &targetAttestationToken,
			Enc:                    (*generated.KeyEncryptionAlgorithm)(options.Algorithm),
			Nonce:                  options.Nonce,
		},
		&generated.KeyVaultClientReleaseOptions{},
	)

	if err != nil {
		return ReleaseKeyResponse{}, err
	}

	return ReleaseKeyResponse{
		Value: resp.Value,
	}, err
}

// UpdateKeyRotationPolicyOptions contains optional parameters for UpdateKeyRotationPolicy.
type UpdateKeyRotationPolicyOptions struct {
	// placeholder for future optional parameters
}

// UpdateKeyRotationPolicyResponse is returned by UpdateKeyRotationPolicy.
type UpdateKeyRotationPolicyResponse struct {
	RotationPolicy
}

func updateKeyRotationPolicyResponseFromGenerated(i generated.KeyVaultClientUpdateKeyRotationPolicyResponse) UpdateKeyRotationPolicyResponse {
	var acts []*LifetimeActions
	for _, a := range i.LifetimeActions {
		acts = append(acts, lifetimeActionsFromGenerated(a))
	}
	var attribs *RotationPolicyAttributes
	if i.Attributes != nil {
		attribs = &RotationPolicyAttributes{
			ExpiresIn: i.Attributes.ExpiryTime,
			CreatedOn: i.Attributes.Created,
			UpdatedOn: i.Attributes.Updated,
		}
	}
	return UpdateKeyRotationPolicyResponse{
		RotationPolicy: RotationPolicy{
			ID:              i.ID,
			LifetimeActions: acts,
			Attributes:      attribs,
		},
	}
}

// UpdateKeyRotationPolicy updates the key's rotation policy. Pass nil for options to accept default values.
func (c *Client) UpdateKeyRotationPolicy(ctx context.Context, keyName string, policy RotationPolicy, options *UpdateKeyRotationPolicyOptions) (UpdateKeyRotationPolicyResponse, error) {
	resp, err := c.kvClient.UpdateKeyRotationPolicy(
		ctx,
		c.vaultURL,
		keyName,
		policy.toGenerated(),
		&generated.KeyVaultClientUpdateKeyRotationPolicyOptions{},
	)

	if err != nil {
		return UpdateKeyRotationPolicyResponse{}, err
	}

	return updateKeyRotationPolicyResponseFromGenerated(resp), nil
}
