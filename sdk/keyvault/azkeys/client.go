//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"errors"
	"net/http"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/convert"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/models"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/responses"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Client is the struct for interacting with a KeyVault Keys instance
type Client struct {
	kvClient *generated.KeyVaultClient
	vaultUrl string
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
		Logging:          c.Logging,
		Retry:            c.Retry,
		Telemetry:        c.Telemetry,
		Transport:        c.Transport,
		PerCallPolicies:  c.PerCallPolicies,
		PerRetryPolicies: c.PerRetryPolicies,
	}
}

// NewClient returns a pointer to a Client object affinitized to a vaultUrl.
func NewClient(vaultUrl string, credential azcore.TokenCredential, options *ClientOptions) (*Client, error) {
	if options == nil {
		options = &ClientOptions{}
	}

	genOptions := options.toConnectionOptions()

	genOptions.PerRetryPolicies = append(
		genOptions.PerRetryPolicies,
		shared.NewKeyVaultChallengePolicy(credential),
	)

	pl := runtime.NewPipeline(generated.ModuleName, generated.ModuleVersion, runtime.PipelineOptions{}, genOptions)
	return &Client{
		kvClient: generated.NewKeyVaultClient(pl),
		vaultUrl: vaultUrl,
	}, nil
}

// CreateKeyOptions contains the optional parameters for the KeyVaultClient.CreateKey method.
type CreateKeyOptions struct {
	// Elliptic curve name. For valid values, see JsonWebKeyCurveName.
	Curve *JSONWebKeyCurveName `json:"crv,omitempty"`

	// The attributes of a key managed by the key vault service.
	KeyAttributes *KeyAttributes         `json:"attributes,omitempty"`
	KeyOps        []*JSONWebKeyOperation `json:"key_ops,omitempty"`

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`

	// The public exponent for a RSA key.
	PublicExponent *int32 `json:"public_exponent,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// convert CreateKeyOptions to *generated.KeyVaultClientCreateKeyOptions
func (c *CreateKeyOptions) toGenerated() *generated.KeyVaultClientCreateKeyOptions {
	return &generated.KeyVaultClientCreateKeyOptions{}
}

// convert CreateKeyOptions to generated.KeyCreateParameters
func (c *CreateKeyOptions) toKeyCreateParameters(keyType KeyType) generated.KeyCreateParameters {
	var attribs *generated.KeyAttributes
	if c.KeyAttributes != nil {
		attribs = convert.KeyAttributesToGenerated(c.KeyAttributes)
	}

	var ops []*generated.JSONWebKeyOperation
	for _, o := range c.KeyOps {
		ops = append(ops, (*generated.JSONWebKeyOperation)(o))
	}

	return generated.KeyCreateParameters{
		Kty:            convert.KeyTypeToGenerated(keyType),
		Curve:          (*generated.JSONWebKeyCurveName)(c.Curve),
		KeyAttributes:  attribs,
		KeyOps:         ops,
		KeySize:        c.KeySize,
		PublicExponent: c.PublicExponent,
		Tags:           convert.ToGeneratedMap(c.Tags),
	}
}

// creates CreateKeyResponse from generated.KeyVaultClient.CreateKeyResponse
func createKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) responses.CreateKey {
	return responses.CreateKey{
		KeyBundle: models.KeyBundle{
			Attributes: convert.KeyAttributesFromGenerated(g.Attributes),
			Key:        convert.JSONWebKeyFromGenerated(g.Key),
			Tags:       convert.FromGeneratedMap(g.Tags),
			Managed:    g.Managed,
		},
	}
}

// CreateKey - The create key operation can be used to create any key type in Azure Key Vault.
// If the named key already exists, Azure Key Vault creates
// a new version of the key. It requires the keys/create  permission.
func (c *Client) CreateKey(ctx context.Context, name string, keyType KeyType, options *CreateKeyOptions) (responses.CreateKey, error) {
	if options == nil {
		options = &CreateKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), options.toGenerated())
	if err != nil {
		return responses.CreateKey{}, err
	}

	return createKeyResponseFromGenerated(resp), nil
}

// CreateECKeyOptions contains the optional parameters for the KeyVaultClient.CreateECKey method
type CreateECKeyOptions struct {
	// Elliptic curve name. For valid values, see JsonWebKeyCurveName.
	CurveName *JSONWebKeyCurveName `json:"crv,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`

	// Whether to create an EC key with HSM protection
	HardwareProtected bool
}

// convert CreateECKeyOptions to generated.KeyCreateParameters
func (c *CreateECKeyOptions) toKeyCreateParameters(keyType KeyType) generated.KeyCreateParameters {
	return generated.KeyCreateParameters{
		Kty:   convert.KeyTypeToGenerated(keyType),
		Curve: (*generated.JSONWebKeyCurveName)(c.CurveName),
		Tags:  convert.ToGeneratedMap(c.Tags),
	}
}

// convert the generated.KeyVaultClientCreateKeyResponse to CreateECKeyResponse
func createECKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) responses.CreateECKey {
	return responses.CreateECKey{
		KeyBundle: models.KeyBundle{
			Attributes: convert.KeyAttributesFromGenerated(g.Attributes),
			Key:        convert.JSONWebKeyFromGenerated(g.Key),
			Tags:       convert.FromGeneratedMap(g.Tags),
			Managed:    g.Managed,
		},
	}
}

// CreateKey - The create key operation can be used to create a new elliptic key curve in Azure Key Vault.
// If the named key already exists, Azure Key Vault creates
// a new version of the key. It requires the keys/create  permission.
func (c *Client) CreateECKey(ctx context.Context, name string, options *CreateECKeyOptions) (responses.CreateECKey, error) {
	keyType := KeyTypeValues.EC()

	if options != nil && options.HardwareProtected {
		keyType = KeyTypeValues.ECHSM()
	} else if options == nil {
		options = &CreateECKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &generated.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return responses.CreateECKey{}, err
	}

	return createECKeyResponseFromGenerated(resp), nil
}

// CreateOCTKeyOptions contains the optional parameters for the Client.CreateOCTKey method
type CreateOCTKeyOptions struct {
	// Hardware Protected OCT Key
	HardwareProtected bool

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// conver the CreateOCTKeyOptions to generated.KeyCreateParameters
func (c *CreateOCTKeyOptions) toKeyCreateParameters(keyType KeyType) generated.KeyCreateParameters {
	return generated.KeyCreateParameters{
		Kty:     convert.KeyTypeToGenerated(keyType),
		KeySize: c.KeySize,
		Tags:    convert.ToGeneratedMap(c.Tags),
	}
}

// convert generated response to CreateOCTKeyResponse
func createOCTKeyResponseFromGenerated(i generated.KeyVaultClientCreateKeyResponse) responses.CreateOCTKey {
	return responses.CreateOCTKey{
		KeyBundle: models.KeyBundle{
			Attributes: convert.KeyAttributesFromGenerated(i.Attributes),
			Key:        convert.JSONWebKeyFromGenerated(i.Key),
			Tags:       convert.FromGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// CreateKey - The create key operation can be used to create a new octet sequence (symmetric) key in Azure Key Vault.
// If the named key already exists, Azure Key Vault creates
// a new version of the key. It requires the keys/create permission.
func (c *Client) CreateOCTKey(ctx context.Context, name string, options *CreateOCTKeyOptions) (responses.CreateOCTKey, error) {
	keyType := KeyTypeValues.Oct()

	if options != nil && options.HardwareProtected {
		keyType = KeyTypeValues.OctHSM()
	} else if options == nil {
		options = &CreateOCTKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &generated.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return responses.CreateOCTKey{}, err
	}

	return createOCTKeyResponseFromGenerated(resp), nil
}

// CreateRSAKeyOptions contains the optional parameters for the Client.CreateRSAKey method.
type CreateRSAKeyOptions struct {
	// Hardware Protected OCT Key
	HardwareProtected bool

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	KeySize *int32 `json:"key_size,omitempty"`

	// The public exponent for a RSA key.
	PublicExponent *int32 `json:"public_exponent,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// convert CreateRSAKeyOptions to generated.KeyCreateParameters
func (c CreateRSAKeyOptions) toKeyCreateParameters(k KeyType) generated.KeyCreateParameters {
	return generated.KeyCreateParameters{
		Kty:            convert.KeyTypeToGenerated(k),
		KeySize:        c.KeySize,
		PublicExponent: c.PublicExponent,
		Tags:           convert.ToGeneratedMap(c.Tags),
	}
}

// convert internal response to CreateRSAKeyResponse
func createRSAKeyResponseFromGenerated(i generated.KeyVaultClientCreateKeyResponse) responses.CreateRSAKey {
	return responses.CreateRSAKey{
		KeyBundle: models.KeyBundle{
			Attributes: convert.KeyAttributesFromGenerated(i.Attributes),
			Key:        convert.JSONWebKeyFromGenerated(i.Key),
			Tags:       convert.FromGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// CreateKey - The create key operation can be used to create a new RSA key in Azure Key Vault.
// If the named key already exists, Azure Key Vault creates
// a new version of the key. It requires the keys/create  permission.
func (c *Client) CreateRSAKey(ctx context.Context, name string, options *CreateRSAKeyOptions) (responses.CreateRSAKey, error) {
	keyType := KeyTypeValues.RSA()

	if options != nil && options.HardwareProtected {
		keyType = KeyTypeValues.RSAHSM()
	} else if options == nil {
		options = &CreateRSAKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &generated.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return responses.CreateRSAKey{}, err
	}

	return createRSAKeyResponseFromGenerated(resp), nil
}

// ListKeysOptions contains the optional parameters for the Client.ListKeys method
type ListKeysOptions struct {
	MaxResults *int32
}

// convert ListKeysOptions to generated options
func (l ListKeysOptions) toGenerated() *generated.KeyVaultClientGetKeysOptions {
	return &generated.KeyVaultClientGetKeysOptions{Maxresults: l.MaxResults}
}

// ListKeys retrieves a list of the keys in the Key Vault as JSON Web Key structures that contain the
// public part of a stored key. The LIST operation is applicable to all key types, however only the
// base key identifier, attributes, and tags are provided in the response. Individual versions of a
// key are not listed in the response. This operation requires the keys/list permission.
func (c *Client) ListKeys(options *ListKeysOptions) *responses.ListKeysPager {
	if options == nil {
		options = &ListKeysOptions{}
	}
	return responses.NewListKeysPager(c.kvClient.GetKeys(c.vaultUrl, options.toGenerated()))
}

// GetKeyOptions contains the options for the Client.GetKey method
type GetKeyOptions struct {
	Version string
}

// convert internal response to GetKeyResponse
func getKeyResponseFromGenerated(i generated.KeyVaultClientGetKeyResponse) responses.GetKey {
	return responses.GetKey{
		KeyBundle: models.KeyBundle{
			Attributes: convert.KeyAttributesFromGenerated(i.Attributes),
			Key:        convert.JSONWebKeyFromGenerated(i.Key),
			Tags:       convert.FromGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// GetKey - The get key operation is applicable to all key types. If the requested key is symmetric, then
// no key material is released in the response. This operation requires the keys/get permission.
func (c *Client) GetKey(ctx context.Context, keyName string, options *GetKeyOptions) (responses.GetKey, error) {
	if options == nil {
		options = &GetKeyOptions{}
	}

	resp, err := c.kvClient.GetKey(ctx, c.vaultUrl, keyName, options.Version, &generated.KeyVaultClientGetKeyOptions{})
	if err != nil {
		return responses.GetKey{}, err
	}

	return getKeyResponseFromGenerated(resp), err
}

// GetDeletedKeyOptions contains the optional parameters for the Client.GetDeletedKey method
type GetDeletedKeyOptions struct{}

// convert the GetDeletedKeyOptions to the internal representation
func (g GetDeletedKeyOptions) toGenerated() *generated.KeyVaultClientGetDeletedKeyOptions {
	return &generated.KeyVaultClientGetDeletedKeyOptions{}
}

// convert generated response to GetDeletedKeyResponse
func getDeletedKeyResponseFromGenerated(i generated.KeyVaultClientGetDeletedKeyResponse) responses.GetDeletedKey {
	return responses.GetDeletedKey{
		DeletedKeyBundle: models.DeletedKeyBundle{
			KeyBundle: models.KeyBundle{
				Attributes: convert.KeyAttributesFromGenerated(i.Attributes),
				Key:        convert.JSONWebKeyFromGenerated(i.Key),
				Tags:       convert.FromGeneratedMap(i.Tags),
				Managed:    i.Managed,
			},
			RecoveryID:         i.RecoveryID,
			DeletedDate:        i.DeletedDate,
			ScheduledPurgeDate: i.ScheduledPurgeDate,
		},
	}
}

// GetDeletedKey - The Get Deleted Key operation is applicable for soft-delete enabled vaults.
// While the operation can be invoked on any vault, it will return an error if invoked on a non
// soft-delete enabled vault. This operation requires the keys/get permission.
func (c *Client) GetDeletedKey(ctx context.Context, keyName string, options *GetDeletedKeyOptions) (responses.GetDeletedKey, error) {
	if options == nil {
		options = &GetDeletedKeyOptions{}
	}

	resp, err := c.kvClient.GetDeletedKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return responses.GetDeletedKey{}, err
	}

	return getDeletedKeyResponseFromGenerated(resp), nil
}

// PurgeDeletedKeyOptions is the struct for any future options for Client.PurgeDeletedKey.
type PurgeDeletedKeyOptions struct{}

// convert options to internal options
func (p *PurgeDeletedKeyOptions) toGenerated() *generated.KeyVaultClientPurgeDeletedKeyOptions {
	return &generated.KeyVaultClientPurgeDeletedKeyOptions{}
}

// Converts the generated response to the publicly exposed version.
func purgeDeletedKeyResponseFromGenerated(i generated.KeyVaultClientPurgeDeletedKeyResponse) responses.PurgeDeletedKey {
	return responses.PurgeDeletedKey{}
}

// PurgeDeletedKey deletes the specified key. The purge deleted key operation removes the key permanently, without the possibility of recovery.
// This operation can only be enabled on a soft-delete enabled vault. This operation requires the key/purge permission.
func (c *Client) PurgeDeletedKey(ctx context.Context, keyName string, options *PurgeDeletedKeyOptions) (responses.PurgeDeletedKey, error) {
	if options == nil {
		options = &PurgeDeletedKeyOptions{}
	}
	resp, err := c.kvClient.PurgeDeletedKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	return purgeDeletedKeyResponseFromGenerated(resp), err
}

// BeginDeleteKeyOptions contains the optional parameters for the Client.BeginDeleteKey method.
type BeginDeleteKeyOptions struct{}

// convert public options to generated options struct
func (b *BeginDeleteKeyOptions) toGenerated() *generated.KeyVaultClientDeleteKeyOptions {
	return &generated.KeyVaultClientDeleteKeyOptions{}
}

// BeginDeleteKey deletes a key from the keyvault. Delete cannot be applied to an individual version of a key. This operation
// requires the key/delete permission. This response contains a Poller struct that can be used to Poll for a response, or the
// response PollUntilDone function can be used to poll until completion.
func (c *Client) BeginDeleteKey(ctx context.Context, keyName string, options *BeginDeleteKeyOptions) (responses.BeginDeleteKey, error) {
	if options == nil {
		options = &BeginDeleteKeyOptions{}
	}
	resp, err := c.kvClient.DeleteKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return responses.BeginDeleteKey{}, err
	}

	getResp, err := c.kvClient.GetDeletedKey(ctx, c.vaultUrl, keyName, nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != http.StatusNotFound {
			return responses.BeginDeleteKey{}, err
		}
	}

	s := responses.NewDeleteKeyPoller(responses.NewDeleteKeyPollerParams{
		VaultUrl:       c.vaultUrl,
		KeyName:        keyName,
		Client:         c.kvClient,
		DeleteResponse: resp,
		LastResponse:   getResp,
	})

	return responses.BeginDeleteKey{
		Poller: s,
	}, nil
}

// BackupKeyOptions contains the optional parameters for the Client.BackupKey method
type BackupKeyOptions struct{}

// convert Options to generated version
func (b BackupKeyOptions) toGenerated() *generated.KeyVaultClientBackupKeyOptions {
	return &generated.KeyVaultClientBackupKeyOptions{}
}

// convert internal reponse to BackupKeyResponse
func backupKeyResponseFromGenerated(i generated.KeyVaultClientBackupKeyResponse) responses.BackupKey {
	return responses.BackupKey{
		Value: i.Value,
	}
}

// BackupKey - The Key Backup operation exports a key from Azure Key Vault in a protected form.
// Note that this operation does NOT return key material in a form that can be used outside the
// Azure Key Vault system, the returned key material is either protected to a Azure Key Vault
// HSM or to Azure Key Vault itself. The intent of this operation is to allow a client to GENERATE
// a key in one Azure Key Vault instance, BACKUP the key, and then RESTORE it into another Azure
// Key Vault instance. The BACKUP operation may be used to export, in protected form, any key
// type from Azure Key Vault. Individual versions of a key cannot be backed up. BACKUP / RESTORE
// can be performed within geographical boundaries only; meaning that a BACKUP from one geographical
// area cannot be restored to another geographical area. For example, a backup from the US
// geographical area cannot be restored in an EU geographical area. This operation requires the key/backup permission.
func (c *Client) BackupKey(ctx context.Context, keyName string, options *BackupKeyOptions) (responses.BackupKey, error) {
	if options == nil {
		options = &BackupKeyOptions{}
	}

	resp, err := c.kvClient.BackupKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return responses.BackupKey{}, err
	}

	return backupKeyResponseFromGenerated(resp), nil
}

// BeginRecoverDeletedKeyOptions contains the optional parameters for the Client.BeginRecoverDeletedKey operation
type BeginRecoverDeletedKeyOptions struct{}

// Convert the publicly exposed options object to the generated version
func (b BeginRecoverDeletedKeyOptions) toGenerated() *generated.KeyVaultClientRecoverDeletedKeyOptions {
	return &generated.KeyVaultClientRecoverDeletedKeyOptions{}
}

// BeginRecoverDeletedKey recovers the deleted key in the specified vault to the latest version.
// This operation can only be performed on a soft-delete enabled vault. This operation requires the keys/recover permission.
func (c *Client) BeginRecoverDeletedKey(ctx context.Context, keyName string, options *BeginRecoverDeletedKeyOptions) (responses.BeginRecoverDeletedKey, error) {
	if options == nil {
		options = &BeginRecoverDeletedKeyOptions{}
	}
	resp, err := c.kvClient.RecoverDeletedKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return responses.BeginRecoverDeletedKey{}, err
	}

	getResp, err := c.kvClient.GetKey(ctx, c.vaultUrl, keyName, "", nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != http.StatusNotFound {
			return responses.BeginRecoverDeletedKey{}, err
		}
	}

	b := responses.NewRecoverDeletedKeyPoller(responses.NewRecoverDeletedKeyPollerParams{
		LastResponse:    getResp,
		KeyName:         keyName,
		Client:          c.kvClient,
		VaultUrl:        c.vaultUrl,
		RecoverResponse: resp,
	})

	return responses.BeginRecoverDeletedKey{
		Poller: b,
	}, nil
}

// UpdateKeyPropertiesOptions contains the optional parameters for the Client.UpdateKeyProperties method
type UpdateKeyPropertiesOptions struct {
	// The version of a key to update
	Version string

	// The attributes of a key managed by the key vault service.
	KeyAttributes *KeyAttributes `json:"attributes,omitempty"`

	// Json web key operations. For more information on possible key operations, see JsonWebKeyOperation.
	KeyOps []*JSONWebKeyOperation `json:"key_ops,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// convert the options to generated.KeyUpdateParameters struct
func (u UpdateKeyPropertiesOptions) toKeyUpdateParameters() generated.KeyUpdateParameters {
	var attribs *generated.KeyAttributes
	if u.KeyAttributes != nil {
		attribs = convert.KeyAttributesToGenerated(u.KeyAttributes)
	}

	var ops []*generated.JSONWebKeyOperation
	for _, o := range u.KeyOps {
		ops = append(ops, (*generated.JSONWebKeyOperation)(o))
	}

	return generated.KeyUpdateParameters{
		KeyOps:        ops,
		KeyAttributes: attribs,
		Tags:          convert.ToGeneratedMap(u.Tags),
	}
}

// convert options to generated options
func (u UpdateKeyPropertiesOptions) toGeneratedOptions() *generated.KeyVaultClientUpdateKeyOptions {
	return &generated.KeyVaultClientUpdateKeyOptions{}
}

// convert the internal response to UpdateKeyPropertiesResponse
func updateKeyPropertiesFromGenerated(i generated.KeyVaultClientUpdateKeyResponse) responses.UpdateKeyProperties {
	return responses.UpdateKeyProperties{
		KeyBundle: models.KeyBundle{
			Attributes: convert.KeyAttributesFromGenerated(i.Attributes),
			Key:        convert.JSONWebKeyFromGenerated(i.Key),
			Tags:       convert.FromGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// UpdateKey - In order to perform this operation, the key must already exist in the Key Vault.
// Note: The cryptographic material of a key itself cannot be changed. This operation requires the keys/update permission.
func (c *Client) UpdateKeyProperties(ctx context.Context, keyName string, options *UpdateKeyPropertiesOptions) (responses.UpdateKeyProperties, error) {
	if options == nil {
		options = &UpdateKeyPropertiesOptions{}
	}
	resp, err := c.kvClient.UpdateKey(
		ctx,
		c.vaultUrl,
		keyName,
		options.Version,
		options.toKeyUpdateParameters(),
		options.toGeneratedOptions(),
	)
	if err != nil {
		return responses.UpdateKeyProperties{}, err
	}

	return updateKeyPropertiesFromGenerated(resp), nil
}

// ListDeletedKeysOptions contains the optional parameters for the Client.ListDeletedKeys operation.
type ListDeletedKeysOptions struct {
	// Maximum number of results to return in a page. If not specified the service will return up to 25 results.
	MaxResults *int32
}

// Convert publicly exposed options to the generated version.a
func (l *ListDeletedKeysOptions) toGenerated() *generated.KeyVaultClientGetDeletedKeysOptions {
	return &generated.KeyVaultClientGetDeletedKeysOptions{
		Maxresults: l.MaxResults,
	}
}

// ListDeletedKeys lists all versions of the specified key. The full key identifier and attributes are provided
// in the response. No values are returned for the keys. This operation requires the keys/list permission.
func (c *Client) ListDeletedKeys(options *ListDeletedKeysOptions) *responses.ListDeletedKeysPager {
	if options == nil {
		options = &ListDeletedKeysOptions{}
	}
	return responses.NewListDeletedKeysPager(c.kvClient.GetDeletedKeys(c.vaultUrl, options.toGenerated()))
}

// ListKeyVersionsOptions contains the options for the ListKeyVersions operations
type ListKeyVersionsOptions struct {
	// Maximum number of results to return in a page. If not specified the service will return up to 25 results.
	MaxResults *int32
}

// convert the public ListKeyVersionsOptions to the generated version
func (l *ListKeyVersionsOptions) toGenerated() *generated.KeyVaultClientGetKeyVersionsOptions {
	if l == nil {
		return &generated.KeyVaultClientGetKeyVersionsOptions{}
	}
	return &generated.KeyVaultClientGetKeyVersionsOptions{
		Maxresults: l.MaxResults,
	}
}

// ListKeyVersions lists all versions of the specified key. The full key identifer and
// attributes are provided in the response. No values are returned for the keys. This operation
// requires the keys/list permission.
func (c *Client) ListKeyVersions(keyName string, options *ListKeyVersionsOptions) *responses.ListKeyVersionsPager {
	if options == nil {
		options = &ListKeyVersionsOptions{}
	}
	return responses.NewListKeyVersionsPager(c.kvClient.GetKeyVersions(c.vaultUrl, keyName, options.toGenerated()))
}

// RestoreKeyBackupOptions contains the optional parameters for the Client.RestoreKey method.
type RestoreKeyBackupOptions struct{}

func (r RestoreKeyBackupOptions) toGenerated() *generated.KeyVaultClientRestoreKeyOptions {
	return &generated.KeyVaultClientRestoreKeyOptions{}
}

// converts the generated response to the publicly exposed version.
func restoreKeyBackupResponseFromGenerated(i generated.KeyVaultClientRestoreKeyResponse) responses.RestoreKeyBackup {
	return responses.RestoreKeyBackup{
		KeyBundle: models.KeyBundle{
			Attributes: convert.KeyAttributesFromGenerated(i.Attributes),
			Key:        convert.JSONWebKeyFromGenerated(i.Key),
			Tags:       convert.FromGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// RestoreKeyBackup restores a backed up key, and all its versions, to a vault. This operation requires the keys/restore permission.
// The backup parameter is a blob of the key to restore, this can be received from the Client.BackupKey function.
func (c *Client) RestoreKeyBackup(ctx context.Context, keyBackup []byte, options *RestoreKeyBackupOptions) (responses.RestoreKeyBackup, error) {
	if options == nil {
		options = &RestoreKeyBackupOptions{}
	}

	resp, err := c.kvClient.RestoreKey(ctx, c.vaultUrl, generated.KeyRestoreParameters{KeyBundleBackup: keyBackup}, options.toGenerated())
	if err != nil {
		return responses.RestoreKeyBackup{}, err
	}

	return restoreKeyBackupResponseFromGenerated(resp), nil
}

// ImportKeyOptions contains the optional parameters for the Client.ImportKeyOptions parameter
type ImportKeyOptions struct {
	// Whether to import as a hardware key (HSM) or software key.
	Hsm *bool `json:"Hsm,omitempty"`

	// The key management attributes.
	KeyAttributes *KeyAttributes `json:"attributes,omitempty"`

	// Application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

func (i ImportKeyOptions) toImportKeyParameters(key JSONWebKey) generated.KeyImportParameters {
	var attribs *generated.KeyAttributes
	if i.KeyAttributes != nil {
		attribs = convert.KeyAttributesToGenerated(i.KeyAttributes)
	}
	return generated.KeyImportParameters{
		Key:           convert.JSONWebKeyToGenerated(key),
		Hsm:           i.Hsm,
		KeyAttributes: attribs,
		Tags:          convert.ToGeneratedMap(i.Tags),
	}
}

// convert the generated response to the ImportKeyResponse
func importKeyResponseFromGenerated(i generated.KeyVaultClientImportKeyResponse) responses.ImportKey {
	return responses.ImportKey{
		KeyBundle: models.KeyBundle{
			Attributes: convert.KeyAttributesFromGenerated(i.Attributes),
			Key:        convert.JSONWebKeyFromGenerated(i.Key),
			Tags:       convert.FromGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// ImportKey - The import key operation may be used to import any key type into an Azure Key Vault.
// If the named key already exists, Azure Key Vault creates a new version of the key. This operation
// requires the  keys/import permission.
func (c *Client) ImportKey(ctx context.Context, keyName string, key JSONWebKey, options *ImportKeyOptions) (responses.ImportKey, error) {
	if options == nil {
		options = &ImportKeyOptions{}
	}

	resp, err := c.kvClient.ImportKey(ctx, c.vaultUrl, keyName, options.toImportKeyParameters(key), &generated.KeyVaultClientImportKeyOptions{})
	if err != nil {
		return responses.ImportKey{}, err
	}

	return importKeyResponseFromGenerated(resp), nil
}

// GetRandomBytesOptions contains the optional parameters for the Client.GetRandomBytes function.
type GetRandomBytesOptions struct{}

func (g GetRandomBytesOptions) toGenerated() *generated.KeyVaultClientGetRandomBytesOptions {
	return &generated.KeyVaultClientGetRandomBytesOptions{}
}

// GetRandomBytes gets the requested number of bytes containing random values from a managed HSM.
// If the operation fails it returns the *KeyVaultError error type.
func (c *Client) GetRandomBytes(ctx context.Context, count *int32, options *GetRandomBytesOptions) (responses.GetRandomBytes, error) {
	if options == nil {
		options = &GetRandomBytesOptions{}
	}

	resp, err := c.kvClient.GetRandomBytes(
		ctx,
		c.vaultUrl,
		generated.GetRandomBytesRequest{Count: count},
		options.toGenerated(),
	)

	if err != nil {
		return responses.GetRandomBytes{}, err
	}

	return responses.GetRandomBytes{
		Value: resp.Value,
	}, nil
}

type RotateKeyOptions struct{}

func (r RotateKeyOptions) toGenerated() *generated.KeyVaultClientRotateKeyOptions {
	return &generated.KeyVaultClientRotateKeyOptions{}
}

func (c *Client) RotateKey(ctx context.Context, name string, options *RotateKeyOptions) (responses.RotateKey, error) {
	if options == nil {
		options = &RotateKeyOptions{}
	}

	resp, err := c.kvClient.RotateKey(
		ctx,
		c.vaultUrl,
		name,
		options.toGenerated(),
	)
	if err != nil {
		return responses.RotateKey{}, err
	}

	return responses.RotateKey{
		KeyBundle: models.KeyBundle{
			Attributes:    convert.KeyAttributesFromGenerated(resp.Attributes),
			Key:           convert.JSONWebKeyFromGenerated(resp.Key),
			ReleasePolicy: convert.KeyReleasePolicyFromGenerated(resp.ReleasePolicy),
			Tags:          convert.FromGeneratedMap(resp.Tags),
			Managed:       resp.Managed,
		},
	}, nil
}

// GetKeyRotationPolicyOptions contains the optional parameters for the Client.GetKeyRotationPolicy function
type GetKeyRotationPolicyOptions struct{}

func (g GetKeyRotationPolicyOptions) toGenerated() *generated.KeyVaultClientGetKeyRotationPolicyOptions {
	return &generated.KeyVaultClientGetKeyRotationPolicyOptions{}
}

func getKeyRotationPolicyResponseFromGenerated(i generated.KeyVaultClientGetKeyRotationPolicyResponse) responses.GetKeyRotationPolicy {
	var acts []*LifetimeActions
	for _, a := range i.LifetimeActions {
		acts = append(acts, convert.LifetimeActionsFromGenerated(a))
	}
	var attribs *KeyRotationPolicyAttributes
	if i.Attributes != nil {
		attribs = &KeyRotationPolicyAttributes{
			ExpiryTime: i.Attributes.ExpiryTime,
			Created:    i.Attributes.Created,
			Updated:    i.Attributes.Updated,
		}
	}
	return responses.GetKeyRotationPolicy{
		KeyRotationPolicy: models.KeyRotationPolicy{
			ID:              i.ID,
			LifetimeActions: acts,
			Attributes:      attribs,
		},
	}
}

// The GetKeyRotationPolicy operation returns the specified key policy resources in the specified key vault. This operation requires
// the keys/get permission.
func (c *Client) GetKeyRotationPolicy(ctx context.Context, name string, options *GetKeyRotationPolicyOptions) (responses.GetKeyRotationPolicy, error) {
	if options == nil {
		options = &GetKeyRotationPolicyOptions{}
	}

	resp, err := c.kvClient.GetKeyRotationPolicy(
		ctx,
		c.vaultUrl,
		name,
		options.toGenerated(),
	)
	if err != nil {
		return responses.GetKeyRotationPolicy{}, err
	}

	return getKeyRotationPolicyResponseFromGenerated(resp), nil
}

type ReleaseKeyOptions struct {
	// The version of the key to release
	Version string

	// The encryption algorithm to use to protected the exported key material
	Enc *KeyEncryptionAlgorithm `json:"enc,omitempty"`

	// A client provided nonce for freshness.
	Nonce *string `json:"nonce,omitempty"`
}

func (c *Client) ReleaseKey(ctx context.Context, name string, target string, options *ReleaseKeyOptions) (responses.ReleaseKey, error) {
	if options == nil {
		options = &ReleaseKeyOptions{}
	}

	resp, err := c.kvClient.Release(
		ctx,
		c.vaultUrl,
		name,
		options.Version,
		generated.KeyReleaseParameters{
			Target: &target,
			Enc:    (*generated.KeyEncryptionAlgorithm)(options.Enc),
			Nonce:  options.Nonce,
		},
		&generated.KeyVaultClientReleaseOptions{},
	)

	if err != nil {
		return responses.ReleaseKey{}, err
	}

	return responses.ReleaseKey{
		Value: resp.Value,
	}, err
}

// UpdateKeyRotationPolicyOptions contains the optional parameters for the Client.UpdateKeyRotationPolicy function
type UpdateKeyRotationPolicyOptions struct {
	// The key rotation policy attributes.
	Attributes *KeyRotationPolicyAttributes `json:"attributes,omitempty"`

	// Actions that will be performed by Key Vault over the lifetime of a key. For preview, lifetimeActions can only have two items at maximum: one for rotate,
	// one for notify. Notification time would be
	// default to 30 days before expiry and it is not configurable.
	LifetimeActions []*LifetimeActions `json:"lifetimeActions,omitempty"`

	// READ-ONLY; The key policy id.
	ID *string `json:"id,omitempty" azure:"ro"`
}

func (u UpdateKeyRotationPolicyOptions) toGenerated() generated.KeyRotationPolicy {
	var attribs *generated.KeyRotationPolicyAttributes
	if u.Attributes != nil {
		attribs = convert.KeyRotationPolicyAttributesToGenerated(u.Attributes)
	}
	var la []*generated.LifetimeActions
	for _, l := range u.LifetimeActions {
		if l == nil {
			la = append(la, nil)
		} else {
			la = append(la, convert.LifetimeActionsToGenerated(l))
		}
	}

	return generated.KeyRotationPolicy{
		ID:              u.ID,
		LifetimeActions: la,
		Attributes:      attribs,
	}
}

func updateKeyRotationPolicyResponseFromGenerated(i generated.KeyVaultClientUpdateKeyRotationPolicyResponse) responses.UpdateKeyRotationPolicy {
	var acts []*LifetimeActions
	for _, a := range i.LifetimeActions {
		acts = append(acts, convert.LifetimeActionsFromGenerated(a))
	}
	var attribs *KeyRotationPolicyAttributes
	if i.Attributes != nil {
		attribs = &KeyRotationPolicyAttributes{
			ExpiryTime: i.Attributes.ExpiryTime,
			Created:    i.Attributes.Created,
			Updated:    i.Attributes.Updated,
		}
	}
	return responses.UpdateKeyRotationPolicy{
		KeyRotationPolicy: models.KeyRotationPolicy{
			ID:              i.ID,
			LifetimeActions: acts,
			Attributes:      attribs,
		},
	}
}

func (c *Client) UpdateKeyRotationPolicy(ctx context.Context, name string, options *UpdateKeyRotationPolicyOptions) (responses.UpdateKeyRotationPolicy, error) {
	if options == nil {
		options = &UpdateKeyRotationPolicyOptions{}
	}

	resp, err := c.kvClient.UpdateKeyRotationPolicy(
		ctx,
		c.vaultUrl,
		name,
		options.toGenerated(),
		&generated.KeyVaultClientUpdateKeyRotationPolicyOptions{},
	)

	if err != nil {
		return responses.UpdateKeyRotationPolicy{}, err
	}

	return updateKeyRotationPolicyResponseFromGenerated(resp), nil
}
