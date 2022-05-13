//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azkeys/internal/generated"
	shared "github.com/Azure/azure-sdk-for-go/sdk/keyvault/internal"
)

// Client is the struct for interacting with a KeyVault Keys instance
type Client struct {
	kvClient *generated.KeyVaultClient
	vaultURL string
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
		vaultURL: vaultUrl,
	}, nil
}

// VaultURL returns a string of the vault URL
func (c *Client) VaultURL() string {
	return c.vaultURL
}

// CreateKeyOptions contains the optional parameters for the KeyVaultClient.CreateKey method.
type CreateKeyOptions struct {
	// Elliptic curve name. For valid values, see PossibleCurveNameValues.
	Curve *CurveName `json:"crv,omitempty"`

	// The properties of a key managed by the key vault service.
	Properties *Properties  `json:"attributes,omitempty"`
	Operations []*Operation `json:"key_ops,omitempty"`

	// The policy rules under which the key can be exported.
	ReleasePolicy *ReleasePolicy `json:"release_policy,omitempty"`

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	Size *int32 `json:"key_size,omitempty"`

	// The public exponent for a RSA key.
	PublicExponent *int32 `json:"public_exponent,omitempty"`

	// Tags contain application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
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
		Tags:           convertToGeneratedMap(c.Tags),
		ReleasePolicy:  c.ReleasePolicy.toGenerated(),
	}
}

// CreateKeyResponse contains the response from method KeyVaultClient.CreateKey.
type CreateKeyResponse struct {
	Key
}

// creates CreateKeyResponse from generated.KeyVaultClient.CreateKeyResponse
func createKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) CreateKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return CreateKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// CreateKey can be used to create any key type in Azure Key Vault.  If the named key already exists,
// Azure Key Vault creates a new version of the key. It requires the keys/create permission. Pass nil to use the default options.
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

// CreateECKeyOptions contains the optional parameters for the KeyVaultClient.CreateECKey method
type CreateECKeyOptions struct {
	// Elliptic curve name. For valid values, see PossibleCurveNameValues.
	CurveName *CurveName `json:"crv,omitempty"`

	// Tags contain application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`

	// Whether to create an EC key with HSM protection
	HardwareProtected *bool

	// The properties of a key managed by the key vault service.
	Properties *Properties  `json:"attributes,omitempty"`
	Operations []*Operation `json:"key_ops,omitempty"`

	// The policy rules under which the key can be exported.
	ReleasePolicy *ReleasePolicy `json:"release_policy,omitempty"`
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
		Curve:         (*generated.JSONWebKeyCurveName)(c.CurveName),
		Tags:          convertToGeneratedMap(c.Tags),
		KeyOps:        keyOps,
		ReleasePolicy: c.ReleasePolicy.toGenerated(),
		KeyAttributes: c.Properties.toGenerated(),
	}
}

// CreateECKeyResponse contains the response from method Client.CreateECKey.
type CreateECKeyResponse struct {
	Key
}

// convert the generated.KeyVaultClientCreateKeyResponse to CreateECKeyResponse
func createECKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) CreateECKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return CreateECKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// CreateECKey can be used to create a new elliptic key curve in Azure Key Vault. If the
// named key already exists, Azure Key Vault creates a new version of the key.
// EC keys can be created in Standard or Premium SKU vaults, ECHSM can be created in Premium SKU vaults or Managed HSMs.
// It requires the keys/create permission. Pass nil to use the default options.
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

// CreateOctKeyOptions contains the optional parameters for the Client.CreateOCTKey method
type CreateOctKeyOptions struct {
	// Hardware Protected OCT Key
	HardwareProtected *bool

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	Size *int32 `json:"key_size,omitempty"`

	// The properties of a key managed by the key vault service.
	Properties *Properties  `json:"attributes,omitempty"`
	Operations []*Operation `json:"key_ops,omitempty"`

	// The policy rules under which the key can be exported.
	ReleasePolicy *ReleasePolicy `json:"release_policy,omitempty"`

	// Tags contain application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
}

// conver the CreateOCTKeyOptions to generated.KeyCreateParameters
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
		Tags:          convertToGeneratedMap(c.Tags),
		ReleasePolicy: c.ReleasePolicy.toGenerated(),
		KeyAttributes: c.Properties.toGenerated(),
		KeyOps:        keyOps,
	}
}

// CreateOctKeyResponse contains the response from method Client.CreateOCTKey.
type CreateOctKeyResponse struct {
	Key
}

// convert generated response to CreateOCTKeyResponse
func createOctKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) CreateOctKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return CreateOctKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// CreateOctKey can be used to create a new octet sequence (symmetric) key in Azure Key Vault.
// If the named key already exists, Azure Key Vault creates a new version of the key.
// An oct-HSM key can only be created with a Managed HSM vault.
// It requires the keys/create permission. Pass nil to use the default options.
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

// CreateRSAKeyOptions contains the optional parameters for the Client.CreateRSAKey method.
type CreateRSAKeyOptions struct {
	// Hardware Protected OCT Key
	HardwareProtected *bool

	// The key size in bits. For example: 2048, 3072, or 4096 for RSA.
	Size *int32 `json:"key_size,omitempty"`

	// The public exponent for a RSA key.
	PublicExponent *int32 `json:"public_exponent,omitempty"`

	// Tags contain application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`

	// The properties of a key managed by the key vault service.
	Properties *Properties  `json:"attributes,omitempty"`
	Operations []*Operation `json:"key_ops,omitempty"`

	// The policy rules under which the key can be exported.
	ReleasePolicy *ReleasePolicy `json:"release_policy,omitempty"`
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
		Tags:           convertToGeneratedMap(c.Tags),
		KeyAttributes:  c.Properties.toGenerated(),
		KeyOps:         keyOps,
		ReleasePolicy:  c.ReleasePolicy.toGenerated(),
	}
}

// CreateRSAKeyResponse contains the response from method Client.CreateRSAKey.
type CreateRSAKeyResponse struct {
	Key
}

// convert internal response to CreateRSAKeyResponse
func createRSAKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) CreateRSAKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return CreateRSAKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// CreateRSAKey can be used to create a new RSA key in Azure Key Vault. If the named key already
// exists, Azure Key Vault creates a new version of the key. RSA keys can be created in Standard or
// Premium SKU vaults, RSAHSM can be created in Premium SKU vaults or Managed HSMs.
// It requires the keys/create permission. Pass nil to use the default options.
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

// ListPropertiesOfKeysOptions contains the optional parameters for the Client.ListKeys method
type ListPropertiesOfKeysOptions struct {
	// placeholder for future optional parameters
}

// ListPropertiesOfKeysResponse contains the current page of results for the Client.ListSecrets operation
type ListPropertiesOfKeysResponse struct {
	// READ-ONLY; The URL to get the next set of keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of keys in the key vault along with a link to the next page of keys.
	Keys []*KeyItem `json:"value,omitempty" azure:"ro"`
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

// ListPropertiesOfKeys retrieves a list of the keys in the Key Vault as JSON Web Key structures that contain the
// public part of a stored key. The LIST operation is applicable to all key types, however only the
// base key identifier, attributes, and tags are provided in the response. Individual versions of a
// key are not listed in the response. This operation requires the keys/list permission.
func (c *Client) ListPropertiesOfKeys(options *ListPropertiesOfKeysOptions) *runtime.Pager[ListPropertiesOfKeysResponse] {
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
			resp, err := c.kvClient.Pl.Do(req)
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

// GetKeyResponse contains the response for the Client.GetResponse method
type GetKeyResponse struct {
	Key
}

// convert internal response to GetKeyResponse
func getKeyResponseFromGenerated(g generated.KeyVaultClientGetKeyResponse) GetKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return GetKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// GetKey is used to retrieve the content for any single Key. If the requested key is symmetric, then
// no key material is released in the response. This operation requires the keys/get permission.
// Pass nil to use the default options.
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

// GetDeletedKeyOptions contains the optional parameters for the Client.GetDeletedKey method
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
			Properties:         keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
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
// the keys/get permission. Pass nil to use the default options.
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

// PurgeDeletedKeyOptions is the struct for any future options for Client.PurgeDeletedKey.
type PurgeDeletedKeyOptions struct {
	// placeholder for future optional parameters
}

// convert options to internal options
func (p *PurgeDeletedKeyOptions) toGenerated() *generated.KeyVaultClientPurgeDeletedKeyOptions {
	return &generated.KeyVaultClientPurgeDeletedKeyOptions{}
}

// PurgeDeletedKeyResponse contains the response from method Client.PurgeDeletedKey.
type PurgeDeletedKeyResponse struct {
	// placeholder for future response values
}

// Converts the generated response to the publicly exposed version.
func purgeDeletedKeyResponseFromGenerated(i generated.KeyVaultClientPurgeDeletedKeyResponse) PurgeDeletedKeyResponse {
	return PurgeDeletedKeyResponse{}
}

// PurgeDeletedKey deletes the specified key. The purge deleted key operation removes the key permanently, without the possibility of recovery.
// This operation can only be enabled on a soft-delete enabled vault. This operation requires the key/purge permission.
// Pass nil to use the default options.
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

// convert interal response to DeleteKeyResponse
func deleteKeyResponseFromGenerated(g generated.KeyVaultClientDeleteKeyResponse) DeleteKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return DeleteKeyResponse{
		DeletedKey: DeletedKey{
			Properties:         keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
			Key:                jsonWebKeyFromGenerated(g.Key),
			RecoveryID:         g.RecoveryID,
			ReleasePolicy:      keyReleasePolicyFromGenerated(g.ReleasePolicy),
			DeletedOn:          g.DeletedDate,
			ScheduledPurgeDate: g.ScheduledPurgeDate,
		},
	}
}

// BeginDeleteKeyOptions contains the optional parameters for the Client.BeginDeleteKey method.
type BeginDeleteKeyOptions struct {
	// ResumeToken is a string to rehydrate a poller for an operation that has already begun.
	ResumeToken *string
}

// convert public options to generated options struct
func (b *BeginDeleteKeyOptions) toGenerated() *generated.KeyVaultClientDeleteKeyOptions {
	return &generated.KeyVaultClientDeleteKeyOptions{}
}

// DeleteKeyPoller is the poller returned by the Client.StartDeleteKey operation
type DeleteKeyPoller struct {
	keyName           string // This is the key to Poll for in GetDeletedKey
	vaultUrl          string
	client            *generated.KeyVaultClient
	deleteResponse    generated.KeyVaultClientDeleteKeyResponse
	deleteRawResponse *http.Response
	lastResponse      generated.KeyVaultClientGetDeletedKeyResponse
	resumeToken       string
}

// Done returns true if the LRO has reached a terminal state
func (s *DeleteKeyPoller) Done() bool {
	if s.deleteRawResponse == nil {
		return false
	}
	return s.deleteRawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.(
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
func (s *DeleteKeyPoller) Poll(ctx context.Context) (*http.Response, error) {
	var deleteRawResponse *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &deleteRawResponse)
	resp, err := s.client.GetDeletedKey(ctx, s.vaultUrl, s.keyName, nil)
	s.deleteRawResponse = deleteRawResponse
	if deleteRawResponse.StatusCode == http.StatusOK {
		// Service recognizes DeletedKey, operation is done
		s.lastResponse = resp
		return s.deleteRawResponse, nil
	}

	var httpResponseErr *azcore.ResponseError
	if errors.As(err, &httpResponseErr) {
		if httpResponseErr.StatusCode == http.StatusNotFound {
			// This is the expected result
			return s.deleteRawResponse, nil
		}
	}
	return s.deleteRawResponse, err
}

// FinalResponse returns the final response after the operations has finished
func (s *DeleteKeyPoller) FinalResponse(ctx context.Context) (DeleteKeyResponse, error) {
	return deleteKeyResponseFromGenerated(s.deleteResponse), nil
}

// PollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
func (s *DeleteKeyPoller) PollUntilDone(ctx context.Context, t time.Duration) (DeleteKeyResponse, error) {
	for {
		_, err := s.Poll(ctx)
		if err != nil {
			return DeleteKeyResponse{}, err
		}
		if s.Done() {
			break
		}
		time.Sleep(t)
	}

	vaultURL, name, version := shared.ParseID(s.deleteResponse.Key.Kid)
	return DeleteKeyResponse{
		DeletedKey: DeletedKey{
			RecoveryID:         s.deleteResponse.RecoveryID,
			DeletedOn:          s.deleteResponse.DeletedDate,
			ScheduledPurgeDate: s.deleteResponse.ScheduledPurgeDate,
			ReleasePolicy:      keyReleasePolicyFromGenerated(s.deleteResponse.ReleasePolicy),
			Properties:         keyPropertiesFromGenerated(s.deleteResponse.Attributes, s.deleteResponse.Key.Kid, name, version, s.deleteResponse.Managed, vaultURL, s.deleteResponse.Tags),
			Key:                jsonWebKeyFromGenerated(s.deleteResponse.Key),
		},
	}, nil
}

// ResumeToken returns a token for resuming polling at a later time
func (s *DeleteKeyPoller) ResumeToken() (string, error) {
	return s.resumeToken, nil
}

// BeginDeleteKey deletes a key from the keyvault. Delete cannot be applied to an individual version of a key. This operation
// requires the key/delete permission. This response contains a Poller struct that can be used to Poll for a response, or the
// PollUntilDone function can be used to poll until completion. Pass nil to use the default options.
func (c *Client) BeginDeleteKey(ctx context.Context, name string, options *BeginDeleteKeyOptions) (*DeleteKeyPoller, error) {
	if options == nil {
		options = &BeginDeleteKeyOptions{}
	}
	var resumeToken string
	var delResp generated.KeyVaultClientDeleteKeyResponse
	var err error
	if options.ResumeToken == nil {
		delResp, err = c.kvClient.DeleteKey(ctx, c.vaultURL, name, options.toGenerated())
		if err != nil {
			return nil, err
		}

		resumeTokenMarshalled, err := json.Marshal(delResp)
		if err != nil {
			return nil, err
		}
		resumeToken = string(resumeTokenMarshalled)
	} else {
		resumeToken = *options.ResumeToken
		err = json.Unmarshal([]byte(resumeToken), &delResp)
		if err != nil {
			return nil, err
		}
	}

	getResp, err := c.kvClient.GetDeletedKey(ctx, c.vaultURL, name, nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	return &DeleteKeyPoller{
		vaultUrl:       c.vaultURL,
		keyName:        name,
		client:         c.kvClient,
		deleteResponse: delResp,
		lastResponse:   getResp,
		resumeToken:    resumeToken,
	}, nil
}

// BackupKeyOptions contains the optional parameters for the Client.BackupKey method
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
	Value []byte `json:"value,omitempty" azure:"ro"`
}

// convert internal reponse to BackupKeyResponse
func backupKeyResponseFromGenerated(i generated.KeyVaultClientBackupKeyResponse) BackupKeyResponse {
	return BackupKeyResponse{
		Value: i.Value,
	}
}

// BackupKey exports a key from Azure Key Vault in a protected form.
// Note that this operation does NOT return key material in a form that can be used outside the
// Azure Key Vault system, the returned key material is either protected to a Azure Key Vault
// HSM or to Azure Key Vault itself. The intent of this operation is to allow a client to GENERATE
// a key in one Azure Key Vault instance, BACKUP the key, and then RESTORE it into another Azure
// Key Vault instance. The BACKUP operation may be used to export, in protected form, any key
// type from Azure Key Vault. Individual versions of a key cannot be backed up. BACKUP / RESTORE
// can be performed within geographical boundaries only; meaning that a BACKUP from one geographical
// area cannot be restored to another geographical area. For example, a backup from the US
// geographical area cannot be restored in an EU geographical area. This operation requires the key/backup permission.
// Pass nil to use the default options.
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

// RecoverDeletedKeyPoller implements the RecoverDeletedKeyPoller interface
type RecoverDeletedKeyPoller struct {
	keyName         string
	vaultUrl        string
	client          *generated.KeyVaultClient
	recoverResponse generated.KeyVaultClientRecoverDeletedKeyResponse
	lastResponse    generated.KeyVaultClientGetKeyResponse
	lastRawResponse *http.Response
	finished        bool
	resumeToken     string
}

// Done returns true when the polling operation is completed
func (p *RecoverDeletedKeyPoller) Done() bool {
	if p.lastRawResponse == nil {
		return false
	}
	return p.finished
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
func (p *RecoverDeletedKeyPoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := p.client.GetKey(ctx, p.vaultUrl, p.keyName, "", nil)
	if err == nil {
		// Polling is finished
		p.finished = true
		return nil, nil
	}
	p.lastResponse = resp
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		p.lastRawResponse = httpErr.RawResponse
		if httpErr.StatusCode == http.StatusOK || httpErr.StatusCode == http.StatusNotFound {
			return httpErr.RawResponse, nil
		} else {
			return httpErr.RawResponse, err
		}
	}
	return p.lastRawResponse, err
}

// FinalResponse returns the final response after the operations has finished
func (p *RecoverDeletedKeyPoller) FinalResponse(ctx context.Context) (RecoverDeletedKeyResponse, error) {
	return recoverDeletedKeyResponseFromGenerated(p.recoverResponse), nil
}

// PollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
func (p *RecoverDeletedKeyPoller) PollUntilDone(ctx context.Context, t time.Duration) (RecoverDeletedKeyResponse, error) {
	for {
		_, err := p.Poll(ctx)
		if err != nil {
			return RecoverDeletedKeyResponse{}, err
		}
		if p.Done() {
			break
		}
		time.Sleep(t)
	}
	return recoverDeletedKeyResponseFromGenerated(p.recoverResponse), nil
}

// BeginRecoverDeletedKeyOptions contains the optional parameters for the Client.BeginRecoverDeletedKey operation
type BeginRecoverDeletedKeyOptions struct {
	// ResumeToken returns a string for creating a new poller to begin polling again
	ResumeToken *string
}

// Convert the publicly exposed options object to the generated version
func (b BeginRecoverDeletedKeyOptions) toGenerated() *generated.KeyVaultClientRecoverDeletedKeyOptions {
	return &generated.KeyVaultClientRecoverDeletedKeyOptions{}
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
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// ResumeToken returns a token for resuming polling at a later time
func (r *RecoverDeletedKeyPoller) ResumeToken() (string, error) {
	return r.resumeToken, nil
}

// BeginRecoverDeletedKey recovers the deleted key in the specified vault to the latest version.
// This operation can only be performed on a soft-delete enabled vault. This operation requires
// the keys/recover permission. Pass nil to use the default options.
func (c *Client) BeginRecoverDeletedKey(ctx context.Context, name string, options *BeginRecoverDeletedKeyOptions) (*RecoverDeletedKeyPoller, error) {
	if options == nil {
		options = &BeginRecoverDeletedKeyOptions{}
	}
	var resumeToken string
	var recoverResp generated.KeyVaultClientRecoverDeletedKeyResponse
	var err error
	if options.ResumeToken == nil {
		recoverResp, err = c.kvClient.RecoverDeletedKey(ctx, c.vaultURL, name, options.toGenerated())
		if err != nil {
			return nil, err
		}

		marshalled, err := json.Marshal(recoverResp)
		if err != nil {
			return nil, err
		}
		resumeToken = string(marshalled)
	} else {
		resumeToken = *options.ResumeToken
		err = json.Unmarshal([]byte(resumeToken), &recoverResp)
		if err != nil {
			return nil, err
		}
	}

	var getRawResp *http.Response
	ctx = runtime.WithCaptureResponse(ctx, &getRawResp)
	getResp, err := c.kvClient.GetKey(ctx, c.vaultURL, name, "", nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != http.StatusNotFound {
			return nil, err
		}
	}

	return &RecoverDeletedKeyPoller{
		lastResponse:    getResp,
		keyName:         name,
		client:          c.kvClient,
		vaultUrl:        c.vaultURL,
		recoverResponse: recoverResp,
		lastRawResponse: getRawResp,
		resumeToken:     resumeToken,
	}, nil
}

// UpdateKeyPropertiesOptions contains the optional parameters for the Client.UpdateKeyProperties method
type UpdateKeyPropertiesOptions struct {
	// placeholder for future optional parameters
}

// UpdateKeyPropertiesResponse contains the response for the Client.UpdateKeyProperties method
type UpdateKeyPropertiesResponse struct {
	Key
}

// convert the internal response to UpdateKeyPropertiesResponse
func updateKeyPropertiesFromGenerated(g generated.KeyVaultClientUpdateKeyResponse) UpdateKeyPropertiesResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return UpdateKeyPropertiesResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// UpdateKeyProperties updates the parameters of a key, but cannot be used to update the cryptographic material
// of a key itself. In order to perform this operation, the key must already exist in the Key Vault.
// This operation requires the keys/update permission. Pass nil to use the default options.
func (c *Client) UpdateKeyProperties(ctx context.Context, key Key, options *UpdateKeyPropertiesOptions) (UpdateKeyPropertiesResponse, error) {
	name, version := "", ""
	if key.Properties != nil && key.Properties.Name != nil {
		name = *key.Properties.Name
	}
	if key.Properties != nil && key.Properties.Version != nil {
		version = *key.Properties.Version
	}
	resp, err := c.kvClient.UpdateKey(
		ctx,
		c.vaultURL,
		name,
		version,
		key.toKeyUpdateParameters(),
		&generated.KeyVaultClientUpdateKeyOptions{},
	)
	if err != nil {
		return UpdateKeyPropertiesResponse{}, err
	}

	return updateKeyPropertiesFromGenerated(resp), nil
}

// ListDeletedKeysResponse holds the data for a single page.
type ListDeletedKeysResponse struct {
	// READ-ONLY; The URL to get the next set of deleted keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of the deleted keys in the vault along with a link to the next page of deleted keys
	DeletedKeys []*DeletedKeyItem `json:"value,omitempty" azure:"ro"`
}

// ListDeletedKeysOptions contains the optional parameters for the Client.ListDeletedKeys operation.
type ListDeletedKeysOptions struct {
	// placeholder for future optional parameters
}

// Convert publicly exposed options to the generated version.a
func (l *ListDeletedKeysOptions) toGenerated() *generated.KeyVaultClientGetDeletedKeysOptions {
	return &generated.KeyVaultClientGetDeletedKeysOptions{}
}

// ListDeletedKeys retrieves a list of the public part of deleted keys. This operation includes deletion-specific information.
// The ListDeleted operation is applicable for vaults enabled for soft-delete. While the operation can be invoked on any vault, it will return
// an error if invoked on a non soft-delete enabled vault. This operation requires the keys/list permission.
// If the operation fails it returns an *azcore.ResponseError type. Pass nil to use the default options.
func (c *Client) ListDeletedKeys(options *ListDeletedKeysOptions) *runtime.Pager[ListDeletedKeysResponse] {
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
			resp, err := c.kvClient.Pl.Do(req)
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

// ListPropertiesOfKeyVersionsOptions contains the options for the ListKeyVersions operations
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

// ListPropertiesOfKeyVersionsResponse contains the current page from a ListKeyVersionsPager.PageResponse method
type ListPropertiesOfKeyVersionsResponse struct {
	// READ-ONLY; The URL to get the next set of keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of keys in the key vault along with a link to the next page of keys.
	Keys []KeyItem `json:"value,omitempty" azure:"ro"`
}

// create ListKeysPage from generated pager
func listKeyVersionsPageFromGenerated(i generated.KeyVaultClientGetKeyVersionsResponse) ListPropertiesOfKeyVersionsResponse {
	var keys []KeyItem
	for _, s := range i.Value {
		if s != nil {
			keys = append(keys, *keyItemFromGenerated(s))
		}
	}
	return ListPropertiesOfKeyVersionsResponse{
		NextLink: i.NextLink,
		Keys:     keys,
	}
}

// ListPropertiesOfKeyVersions lists all versions of the specified key. The full key identifer and
// attributes are provided in the response. No values are returned for the keys. This operation
// requires the keys/list permission.
func (c *Client) ListPropertiesOfKeyVersions(keyName string, options *ListPropertiesOfKeyVersionsOptions) *runtime.Pager[ListPropertiesOfKeyVersionsResponse] {
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
			resp, err := c.kvClient.Pl.Do(req)
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

// RestoreKeyBackupOptions contains the optional parameters for the Client.RestoreKey method.
type RestoreKeyBackupOptions struct {
	// placeholder for future optional parameters
}

func (r RestoreKeyBackupOptions) toGenerated() *generated.KeyVaultClientRestoreKeyOptions {
	return &generated.KeyVaultClientRestoreKeyOptions{}
}

// RestoreKeyBackupResponse contains the response object for the Client.RestoreKeyBackup operation.
type RestoreKeyBackupResponse struct {
	Key
}

// converts the generated response to the publicly exposed version.
func restoreKeyBackupResponseFromGenerated(g generated.KeyVaultClientRestoreKeyResponse) RestoreKeyBackupResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return RestoreKeyBackupResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// RestoreKeyBackup restores a backed up key, and all its versions, to a vault. This operation requires the keys/restore permission.
// The backup parameter is a blob of the key to restore, this can be received from the Client.BackupKey function.
// Pass nil to use the default options.
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

// ImportKeyOptions contains the optional parameters for the Client.ImportKeyOptions parameter
type ImportKeyOptions struct {
	// Whether to import as a hardware key (HSM) or software key.
	HardwareProtected *bool `json:"Hsm,omitempty"`

	// The key management attributes.
	Properties *Properties `json:"attributes,omitempty"`

	// Tags contain application specific metadata in the form of key-value pairs.
	Tags map[string]string `json:"tags,omitempty"`
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
		Tags:          convertToGeneratedMap(i.Tags),
	}
}

// ImportKeyResponse contains the response of the Client.ImportKey method
type ImportKeyResponse struct {
	Key
}

// convert the generated response to the ImportKeyResponse
func importKeyResponseFromGenerated(g generated.KeyVaultClientImportKeyResponse) ImportKeyResponse {
	vaultURL, name, version := shared.ParseID(g.Key.Kid)
	return ImportKeyResponse{
		Key: Key{
			Properties: keyPropertiesFromGenerated(g.Attributes, g.Key.Kid, name, version, g.Managed, vaultURL, g.Tags),
			JSONWebKey: jsonWebKeyFromGenerated(g.Key),
			ID:         g.Key.Kid,
			Name:       name,
		},
	}
}

// ImportKey may be used to import any key type into an Azure Key Vault. If the named key already exists,
// Azure Key Vault creates a new version of the key. This operation requires the keys/import permission.
// Pass nil to use the default options.
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

// GetRandomBytesOptions contains the optional parameters for the Client.GetRandomBytes function.
type GetRandomBytesOptions struct {
	// placeholder for future optional parameters
}

func (g GetRandomBytesOptions) toGenerated() *generated.KeyVaultClientGetRandomBytesOptions {
	return &generated.KeyVaultClientGetRandomBytesOptions{}
}

// GetRandomBytesResponse is the response struct for the Client.GetRandomBytes function.
type GetRandomBytesResponse struct {
	// The bytes encoded as a base64url string.
	Value []byte `json:"value,omitempty"`
}

// GetRandomBytes gets the requested number of bytes containing random values from a managed HSM.
// If the operation fails it returns the *KeyVaultError error type. Pass nil to use the default options.
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

// RotateKeyOptions contains the optional parameters for the Client.RotateKey function
type RotateKeyOptions struct {
	// placeholder for future optional parameters
}

func (r RotateKeyOptions) toGenerated() *generated.KeyVaultClientRotateKeyOptions {
	return &generated.KeyVaultClientRotateKeyOptions{}
}

// RotateKeyResponse contains response fields for Client.RotateKey
type RotateKeyResponse struct {
	Key
}

// RotateKey will rotate the key based on the key policy. It requires the keys/rotate permission.
// The system will generate a new version in the specified key.
// Pass nil to use the default options.
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
			Properties:    keyPropertiesFromGenerated(resp.Attributes, resp.Key.Kid, name, version, resp.Managed, vaultURL, resp.Tags),
			JSONWebKey:    jsonWebKeyFromGenerated(resp.Key),
			ReleasePolicy: keyReleasePolicyFromGenerated(resp.ReleasePolicy),
			ID:            resp.Key.Kid,
			Name:          name,
		},
	}, nil
}

// GetKeyRotationPolicyOptions contains the optional parameters for the Client.GetKeyRotationPolicy function
type GetKeyRotationPolicyOptions struct {
	// placeholder for future optional parameters
}

func (g GetKeyRotationPolicyOptions) toGenerated() *generated.KeyVaultClientGetKeyRotationPolicyOptions {
	return &generated.KeyVaultClientGetKeyRotationPolicyOptions{}
}

// GetKeyRotationPolicyResponse contains the response struct for the Client.GetKeyRotationPolicy function
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

// GetKeyRotationPolicy returns the specified key policy resources in the specified key vault. This operation requires
// the keys/get permission. Pass nil to use the default options.
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

// ReleaseKeyOptions contains optional parameters for Client.ReleaseKey
type ReleaseKeyOptions struct {
	// The version of the key to release
	Version string

	// The encryption algorithm to use to protected the exported key material
	Enc *ExportEncryptionAlg `json:"enc,omitempty"`

	// A client provided nonce for freshness.
	Nonce *string `json:"nonce,omitempty"`
}

// ReleaseKeyResponse contains the response of Client.ReleaseKey
type ReleaseKeyResponse struct {
	// READ-ONLY; A signed object containing the released key.
	Value *string `json:"value,omitempty" azure:"ro"`
}

// ReleaseKey is applicable to all key types. The target key must be marked exportable. This operation requires the keys/release permission.
// Pass nil to use the default options.
func (c *Client) ReleaseKey(ctx context.Context, name string, targetAttestationToken string, options *ReleaseKeyOptions) (ReleaseKeyResponse, error) {
	if options == nil {
		options = &ReleaseKeyOptions{}
	}

	resp, err := c.kvClient.Release(
		ctx,
		c.vaultURL,
		name,
		options.Version,
		generated.KeyReleaseParameters{
			TargetAttestationToken: &targetAttestationToken,
			Enc:                    (*generated.KeyEncryptionAlgorithm)(options.Enc),
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

// UpdateKeyRotationPolicyOptions contains the optional parameters for the Client.UpdateKeyRotationPolicy function
type UpdateKeyRotationPolicyOptions struct {
	// placeholder for future optional parameters
}

// UpdateKeyRotationPolicyResponse contains the response for the Client.UpdateKeyRotationPolicy function
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

// UpdateKeyRotationPolicy sets specified members in the key policy.
// This operation requires the keys/update permission.
// Pass nil to use the default options.
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
