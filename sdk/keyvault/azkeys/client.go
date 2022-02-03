//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.

package azkeys

import (
	"context"
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
		attribs = c.KeyAttributes.toGenerated()
	}

	var ops []*generated.JSONWebKeyOperation
	for _, o := range c.KeyOps {
		ops = append(ops, (*generated.JSONWebKeyOperation)(o))
	}

	return generated.KeyCreateParameters{
		Kty:            keyType.toGenerated(),
		Curve:          (*generated.JSONWebKeyCurveName)(c.Curve),
		KeyAttributes:  attribs,
		KeyOps:         ops,
		KeySize:        c.KeySize,
		PublicExponent: c.PublicExponent,
		Tags:           convertToGeneratedMap(c.Tags),
	}
}

// KeyVaultClientCreateKeyResponse contains the response from method KeyVaultClient.CreateKey.
type CreateKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// creates CreateKeyResponse from generated.KeyVaultClient.CreateKeyResponse
func createKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) CreateKeyResponse {
	return CreateKeyResponse{
		RawResponse: g.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(g.Attributes),
			Key:        jsonWebKeyFromGenerated(g.Key),
			Tags:       convertGeneratedMap(g.Tags),
			Managed:    g.Managed,
		},
	}
}

// CreateKey - The create key operation can be used to create any key type in Azure Key Vault.
// If the named key already exists, Azure Key Vault creates
// a new version of the key. It requires the keys/create  permission.
func (c *Client) CreateKey(ctx context.Context, name string, keyType KeyType, options *CreateKeyOptions) (CreateKeyResponse, error) {
	if options == nil {
		options = &CreateKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), options.toGenerated())
	if err != nil {
		return CreateKeyResponse{}, err
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
		Kty:   keyType.toGenerated(),
		Curve: (*generated.JSONWebKeyCurveName)(c.CurveName),
		Tags:  convertToGeneratedMap(c.Tags),
	}
}

// CreateECKeyResponse contains the response from method Client.CreateECKey.
type CreateECKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert the generated.KeyVaultClientCreateKeyResponse to CreateECKeyResponse
func createECKeyResponseFromGenerated(g generated.KeyVaultClientCreateKeyResponse) CreateECKeyResponse {
	return CreateECKeyResponse{
		RawResponse: g.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(g.Attributes),
			Key:        jsonWebKeyFromGenerated(g.Key),
			Tags:       convertGeneratedMap(g.Tags),
			Managed:    g.Managed,
		},
	}
}

// CreateKey - The create key operation can be used to create a new elliptic key curve in Azure Key Vault.
// If the named key already exists, Azure Key Vault creates
// a new version of the key. It requires the keys/create  permission.
func (c *Client) CreateECKey(ctx context.Context, name string, options *CreateECKeyOptions) (CreateECKeyResponse, error) {
	keyType := EC

	if options != nil && options.HardwareProtected {
		keyType = ECHSM
	} else if options == nil {
		options = &CreateECKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &generated.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return CreateECKeyResponse{}, err
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
		Kty:     keyType.toGenerated(),
		KeySize: c.KeySize,
		Tags:    convertToGeneratedMap(c.Tags),
	}
}

// CreateOCTKeyResponse contains the response from method Client.CreateOCTKey.
type CreateOCTKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert generated response to CreateOCTKeyResponse
func createOCTKeyResponseFromGenerated(i generated.KeyVaultClientCreateKeyResponse) CreateOCTKeyResponse {
	return CreateOCTKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       convertGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// CreateKey - The create key operation can be used to create a new octet sequence (symmetric) key in Azure Key Vault.
// If the named key already exists, Azure Key Vault creates
// a new version of the key. It requires the keys/create permission.
func (c *Client) CreateOCTKey(ctx context.Context, name string, options *CreateOCTKeyOptions) (CreateOCTKeyResponse, error) {
	keyType := Oct

	if options != nil && options.HardwareProtected {
		keyType = OctHSM
	} else if options == nil {
		options = &CreateOCTKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &generated.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return CreateOCTKeyResponse{}, err
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
		Kty:            k.toGenerated(),
		KeySize:        c.KeySize,
		PublicExponent: c.PublicExponent,
		Tags:           convertToGeneratedMap(c.Tags),
	}
}

// CreateRSAKeyResponse contains the response from method Client.CreateRSAKey.
type CreateRSAKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert internal response to CreateRSAKeyResponse
func createRSAKeyResponseFromGenerated(i generated.KeyVaultClientCreateKeyResponse) CreateRSAKeyResponse {
	return CreateRSAKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       convertGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// CreateKey - The create key operation can be used to create a new RSA key in Azure Key Vault.
// If the named key already exists, Azure Key Vault creates
// a new version of the key. It requires the keys/create  permission.
func (c *Client) CreateRSAKey(ctx context.Context, name string, options *CreateRSAKeyOptions) (CreateRSAKeyResponse, error) {
	keyType := RSA

	if options != nil && options.HardwareProtected {
		keyType = RSAHSM
	} else if options == nil {
		options = &CreateRSAKeyOptions{}
	}

	resp, err := c.kvClient.CreateKey(ctx, c.vaultUrl, name, options.toKeyCreateParameters(keyType), &generated.KeyVaultClientCreateKeyOptions{})
	if err != nil {
		return CreateRSAKeyResponse{}, err
	}

	return createRSAKeyResponseFromGenerated(resp), nil
}

// ListKeysPager is a Pager for the Client.ListSecrets operation
type ListKeysPager interface {
	// PageResponse returns the current ListKeysPage
	PageResponse() ListKeysPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

// listKeysPager implements the ListKeysPager interface
type listKeysPager struct {
	genPager generated.KeyVaultClientGetKeysPager
}

// PageResponse returns the results from the page most recently fetched from the service
func (l *listKeysPager) PageResponse() ListKeysPage {
	return listKeysPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil
func (l *listKeysPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *listKeysPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListKeysOptions contains the optional parameters for the Client.ListKeys method
type ListKeysOptions struct {
	MaxResults *int32
}

// convert ListKeysOptions to generated options
func (l ListKeysOptions) toGenerated() *generated.KeyVaultClientGetKeysOptions {
	return &generated.KeyVaultClientGetKeysOptions{Maxresults: l.MaxResults}
}

// ListKeysPage contains the current page of results for the Client.ListSecrets operation
type ListKeysPage struct {
	// READ-ONLY; The URL to get the next set of keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of keys in the key vault along with a link to the next page of keys.
	Keys []*KeyItem `json:"value,omitempty" azure:"ro"`
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert internal Response to ListKeysPage
func listKeysPageFromGenerated(i generated.KeyVaultClientGetKeysResponse) ListKeysPage {
	var keys []*KeyItem
	for _, k := range i.Value {
		keys = append(keys, keyItemFromGenerated(k))
	}
	return ListKeysPage{
		RawResponse: i.RawResponse,
		NextLink:    i.NextLink,
		Keys:        keys,
	}
}

// ListKeys retrieves a list of the keys in the Key Vault as JSON Web Key structures that contain the
// public part of a stored key. The LIST operation is applicable to all key types, however only the
// base key identifier, attributes, and tags are provided in the response. Individual versions of a
// key are not listed in the response. This operation requires the keys/list permission.
func (c *Client) ListKeys(options *ListKeysOptions) ListKeysPager {
	if options == nil {
		options = &ListKeysOptions{}
	}
	p := c.kvClient.GetKeys(c.vaultUrl, options.toGenerated())

	return &listKeysPager{
		genPager: *p,
	}
}

// GetKeyOptions contains the options for the Client.GetKey method
type GetKeyOptions struct {
	Version string
}

// GetKeyResponse contains the response for the Client.GetResponse method
type GetKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert internal response to GetKeyResponse
func getKeyResponseFromGenerated(i generated.KeyVaultClientGetKeyResponse) GetKeyResponse {
	return GetKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       convertGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// GetKey - The get key operation is applicable to all key types. If the requested key is symmetric, then
// no key material is released in the response. This operation requires the keys/get permission.
func (c *Client) GetKey(ctx context.Context, keyName string, options *GetKeyOptions) (GetKeyResponse, error) {
	if options == nil {
		options = &GetKeyOptions{}
	}

	resp, err := c.kvClient.GetKey(ctx, c.vaultUrl, keyName, options.Version, &generated.KeyVaultClientGetKeyOptions{})
	if err != nil {
		return GetKeyResponse{}, err
	}

	return getKeyResponseFromGenerated(resp), err
}

// GetDeletedKeyOptions contains the optional parameters for the Client.GetDeletedKey method
type GetDeletedKeyOptions struct{}

// convert the GetDeletedKeyOptions to the internal representation
func (g GetDeletedKeyOptions) toGenerated() *generated.KeyVaultClientGetDeletedKeyOptions {
	return &generated.KeyVaultClientGetDeletedKeyOptions{}
}

// GetDeletedKeyResponse contains the response from a Client.GetDeletedKey
type GetDeletedKeyResponse struct {
	DeletedKeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert generated response to GetDeletedKeyResponse
func getDeletedKeyResponseFromGenerated(i generated.KeyVaultClientGetDeletedKeyResponse) GetDeletedKeyResponse {
	return GetDeletedKeyResponse{
		RawResponse: i.RawResponse,
		DeletedKeyBundle: DeletedKeyBundle{
			KeyBundle: KeyBundle{
				Attributes: keyAttributesFromGenerated(i.Attributes),
				Key:        jsonWebKeyFromGenerated(i.Key),
				Tags:       convertGeneratedMap(i.Tags),
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
func (c *Client) GetDeletedKey(ctx context.Context, keyName string, options *GetDeletedKeyOptions) (GetDeletedKeyResponse, error) {
	if options == nil {
		options = &GetDeletedKeyOptions{}
	}

	resp, err := c.kvClient.GetDeletedKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return GetDeletedKeyResponse{}, err
	}

	return getDeletedKeyResponseFromGenerated(resp), nil
}

// PurgeDeletedKeyOptions is the struct for any future options for Client.PurgeDeletedKey.
type PurgeDeletedKeyOptions struct{}

// convert options to internal options
func (p *PurgeDeletedKeyOptions) toGenerated() *generated.KeyVaultClientPurgeDeletedKeyOptions {
	return &generated.KeyVaultClientPurgeDeletedKeyOptions{}
}

// PurgeDeletedKeyResponse contains the response from method Client.PurgeDeletedKey.
type PurgeDeletedKeyResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// Converts the generated response to the publicly exposed version.
func purgeDeletedKeyResponseFromGenerated(i generated.KeyVaultClientPurgeDeletedKeyResponse) PurgeDeletedKeyResponse {
	return PurgeDeletedKeyResponse{
		RawResponse: i.RawResponse,
	}
}

// PurgeDeletedKey deletes the specified key. The purge deleted key operation removes the key permanently, without the possibility of recovery.
// This operation can only be enabled on a soft-delete enabled vault. This operation requires the key/purge permission.
func (c *Client) PurgeDeletedKey(ctx context.Context, keyName string, options *PurgeDeletedKeyOptions) (PurgeDeletedKeyResponse, error) {
	if options == nil {
		options = &PurgeDeletedKeyOptions{}
	}
	resp, err := c.kvClient.PurgeDeletedKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	return purgeDeletedKeyResponseFromGenerated(resp), err
}

// DeletedKeyResponse contains the response for a Client.BeginDeleteKey operation.
type DeleteKeyResponse struct {
	DeletedKeyBundle
	// RawResponse holds the underlying HTTP response
	RawResponse *http.Response
}

// convert interal response to DeleteKeyResponse
func deleteKeyResponseFromGenerated(i *generated.KeyVaultClientDeleteKeyResponse) *DeleteKeyResponse {
	if i == nil {
		return nil
	}
	return &DeleteKeyResponse{
		RawResponse: i.RawResponse,
	}
}

// BeginDeleteKeyOptions contains the optional parameters for the Client.BeginDeleteKey method.
type BeginDeleteKeyOptions struct{}

// convert public options to generated options struct
func (b *BeginDeleteKeyOptions) toGenerated() *generated.KeyVaultClientDeleteKeyOptions {
	return &generated.KeyVaultClientDeleteKeyOptions{}
}

// DeleteKeyPoller is the interface for the Client.DeleteKey operation.
type DeleteKeyPoller interface {
	// Done returns true if the LRO has reached a terminal state
	Done() bool

	// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
	// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
	// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
	Poll(context.Context) (*http.Response, error)

	// FinalResponse returns the final response after the operations has finished
	FinalResponse(context.Context) (DeleteKeyResponse, error)
}

// The poller returned by the Client.StartDeleteKey operation
type startDeleteKeyPoller struct {
	keyName        string // This is the key to Poll for in GetDeletedKey
	vaultUrl       string
	client         *generated.KeyVaultClient
	deleteResponse generated.KeyVaultClientDeleteKeyResponse
	lastResponse   generated.KeyVaultClientGetDeletedKeyResponse
	RawResponse    *http.Response
}

// Done returns true if the LRO has reached a terminal state
func (s *startDeleteKeyPoller) Done() bool {
	return s.lastResponse.RawResponse != nil
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.(
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.)
func (s *startDeleteKeyPoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := s.client.GetDeletedKey(ctx, s.vaultUrl, s.keyName, nil)
	if err == nil {
		// Service recognizes DeletedKey, operation is done
		s.lastResponse = resp
		return resp.RawResponse, nil
	}

	var httpResponseErr *azcore.ResponseError
	if errors.As(err, &httpResponseErr) {
		if httpResponseErr.StatusCode == http.StatusNotFound {
			// This is the expected result
			return s.deleteResponse.RawResponse, nil
		}
	}
	return s.deleteResponse.RawResponse, err
}

// FinalResponse returns the final response after the operations has finished
func (s *startDeleteKeyPoller) FinalResponse(ctx context.Context) (DeleteKeyResponse, error) {
	return *deleteKeyResponseFromGenerated(&s.deleteResponse), nil
}

// pollUntilDone continually calls the Poll operation until the operation is completed. In between each
// Poll is a wait determined by the t parameter.
func (s *startDeleteKeyPoller) pollUntilDone(ctx context.Context, t time.Duration) (DeleteKeyResponse, error) {
	for {
		resp, err := s.Poll(ctx)
		if err != nil {
			return DeleteKeyResponse{}, err
		}
		s.RawResponse = resp
		if s.Done() {
			break
		}
		time.Sleep(t)
	}
	return DeleteKeyResponse{}, nil
}

// DeleteKeyPollerResponse contains the response from the Client.BeginDeleteKey method
type DeleteKeyPollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (DeleteKeyResponse, error)

	// Poller contains an initialized WidgetPoller
	Poller DeleteKeyPoller

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// BeginDeleteKey deletes a key from the keyvault. Delete cannot be applied to an individual version of a key. This operation
// requires the key/delete permission. This response contains a Poller struct that can be used to Poll for a response, or the
// response PollUntilDone function can be used to poll until completion.
func (c *Client) BeginDeleteKey(ctx context.Context, keyName string, options *BeginDeleteKeyOptions) (DeleteKeyPollerResponse, error) {
	if options == nil {
		options = &BeginDeleteKeyOptions{}
	}
	resp, err := c.kvClient.DeleteKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return DeleteKeyPollerResponse{}, err
	}

	getResp, err := c.kvClient.GetDeletedKey(ctx, c.vaultUrl, keyName, nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != http.StatusNotFound {
			return DeleteKeyPollerResponse{}, err
		}
	}

	s := &startDeleteKeyPoller{
		vaultUrl:       c.vaultUrl,
		keyName:        keyName,
		client:         c.kvClient,
		deleteResponse: resp,
		lastResponse:   getResp,
	}

	return DeleteKeyPollerResponse{
		Poller:        s,
		RawResponse:   resp.RawResponse,
		PollUntilDone: s.pollUntilDone,
	}, nil
}

// BackupKeyOptions contains the optional parameters for the Client.BackupKey method
type BackupKeyOptions struct{}

// convert Options to generated version
func (b BackupKeyOptions) toGenerated() *generated.KeyVaultClientBackupKeyOptions {
	return &generated.KeyVaultClientBackupKeyOptions{}
}

// BackupKeyResponse contains the response from the Client.BackupKey method
type BackupKeyResponse struct {
	// READ-ONLY; The backup blob containing the backed up key.
	Value []byte `json:"value,omitempty" azure:"ro"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert internal reponse to BackupKeyResponse
func backupKeyResponseFromGenerated(i generated.KeyVaultClientBackupKeyResponse) BackupKeyResponse {
	return BackupKeyResponse{
		RawResponse: i.RawResponse,
		Value:       i.Value,
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
func (c *Client) BackupKey(ctx context.Context, keyName string, options *BackupKeyOptions) (BackupKeyResponse, error) {
	if options == nil {
		options = &BackupKeyOptions{}
	}

	resp, err := c.kvClient.BackupKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return BackupKeyResponse{}, err
	}

	return backupKeyResponseFromGenerated(resp), nil
}

// RecoverDeletedKeyPoller is the interface for the Client.RecoverDeletedKey operation
type RecoverDeletedKeyPoller interface {
	// Done returns true if the LRO has reached a terminal state
	Done() bool

	// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
	// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
	// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
	Poll(context.Context) (*http.Response, error)

	// FinalResponse returns the final response after the operations has finished
	FinalResponse(context.Context) (RecoverDeletedKeyResponse, error)
}

// beginRecoverPoller implements the RecoverDeletedKeyPoller interface
type beginRecoverPoller struct {
	keyName         string
	vaultUrl        string
	client          *generated.KeyVaultClient
	recoverResponse generated.KeyVaultClientRecoverDeletedKeyResponse
	lastResponse    generated.KeyVaultClientGetKeyResponse
	RawResponse     *http.Response
}

// Done returns true when the polling operation is completed
func (b *beginRecoverPoller) Done() bool {
	return b.RawResponse.StatusCode == http.StatusOK
}

// Poll fetches the latest state of the LRO. It returns an HTTP response or error.
// If the LRO has completed successfully, the poller's state is updated and the HTTP response is returned.
// If the LRO has completed with failure or was cancelled, the poller's state is updated and the error is returned.
func (b *beginRecoverPoller) Poll(ctx context.Context) (*http.Response, error) {
	resp, err := b.client.GetKey(ctx, b.vaultUrl, b.keyName, "", nil)
	b.lastResponse = resp
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		return httpErr.RawResponse, err
	}
	return resp.RawResponse, nil
}

// FinalResponse returns the final response after the operations has finished
func (b *beginRecoverPoller) FinalResponse(ctx context.Context) (RecoverDeletedKeyResponse, error) {
	return recoverDeletedKeyResponseFromGenerated(b.recoverResponse), nil
}

// pollUntilDone is the method for the Response.PollUntilDone struct
func (b *beginRecoverPoller) pollUntilDone(ctx context.Context, t time.Duration) (RecoverDeletedKeyResponse, error) {
	for {
		resp, err := b.Poll(ctx)
		if err != nil {
			b.RawResponse = resp
		}
		if b.Done() {
			break
		}
		b.RawResponse = resp
		time.Sleep(t)
	}
	return recoverDeletedKeyResponseFromGenerated(b.recoverResponse), nil
}

// BeginRecoverDeletedKeyOptions contains the optional parameters for the Client.BeginRecoverDeletedKey operation
type BeginRecoverDeletedKeyOptions struct{}

// Convert the publicly exposed options object to the generated version
func (b BeginRecoverDeletedKeyOptions) toGenerated() *generated.KeyVaultClientRecoverDeletedKeyOptions {
	return &generated.KeyVaultClientRecoverDeletedKeyOptions{}
}

// RecoverDeletedKeyResponse is the response object for the Client.RecoverDeletedKey operation.
type RecoverDeletedKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// change recover deleted key reponse to the generated version.
func recoverDeletedKeyResponseFromGenerated(i generated.KeyVaultClientRecoverDeletedKeyResponse) RecoverDeletedKeyResponse {
	return RecoverDeletedKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       convertGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// RecoverDeletedKeyPollerResponse contains the response of the Client.BeginRecoverDeletedKey operations
type RecoverDeletedKeyPollerResponse struct {
	// PollUntilDone will poll the service endpoint until a terminal state is reached or an error occurs
	PollUntilDone func(context.Context, time.Duration) (RecoverDeletedKeyResponse, error)

	// Poller contains an initialized RecoverDeletedKeyPoller
	Poller RecoverDeletedKeyPoller

	// RawResponse cotains the underlying HTTP response
	RawResponse *http.Response
}

// BeginRecoverDeletedKey recovers the deleted key in the specified vault to the latest version.
// This operation can only be performed on a soft-delete enabled vault. This operation requires the keys/recover permission.
func (c *Client) BeginRecoverDeletedKey(ctx context.Context, keyName string, options *BeginRecoverDeletedKeyOptions) (RecoverDeletedKeyPollerResponse, error) {
	if options == nil {
		options = &BeginRecoverDeletedKeyOptions{}
	}
	resp, err := c.kvClient.RecoverDeletedKey(ctx, c.vaultUrl, keyName, options.toGenerated())
	if err != nil {
		return RecoverDeletedKeyPollerResponse{}, err
	}

	getResp, err := c.kvClient.GetKey(ctx, c.vaultUrl, keyName, "", nil)
	var httpErr *azcore.ResponseError
	if errors.As(err, &httpErr) {
		if httpErr.StatusCode != http.StatusNotFound {
			return RecoverDeletedKeyPollerResponse{}, err
		}
	}

	b := &beginRecoverPoller{
		lastResponse:    getResp,
		keyName:         keyName,
		client:          c.kvClient,
		vaultUrl:        c.vaultUrl,
		recoverResponse: resp,
		RawResponse:     getResp.RawResponse,
	}

	return RecoverDeletedKeyPollerResponse{
		PollUntilDone: b.pollUntilDone,
		Poller:        b,
		RawResponse:   getResp.RawResponse,
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
		attribs = u.KeyAttributes.toGenerated()
	}

	var ops []*generated.JSONWebKeyOperation
	for _, o := range u.KeyOps {
		ops = append(ops, (*generated.JSONWebKeyOperation)(o))
	}

	return generated.KeyUpdateParameters{
		KeyOps:        ops,
		KeyAttributes: attribs,
		Tags:          convertToGeneratedMap(u.Tags),
	}
}

// convert options to generated options
func (u UpdateKeyPropertiesOptions) toGeneratedOptions() *generated.KeyVaultClientUpdateKeyOptions {
	return &generated.KeyVaultClientUpdateKeyOptions{}
}

// UpdateKeyPropertiesResponse contains the response for the Client.UpdateKeyProperties method
type UpdateKeyPropertiesResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert the internal response to UpdateKeyPropertiesResponse
func updateKeyPropertiesFromGenerated(i generated.KeyVaultClientUpdateKeyResponse) UpdateKeyPropertiesResponse {
	return UpdateKeyPropertiesResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       convertGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// UpdateKey - In order to perform this operation, the key must already exist in the Key Vault.
// Note: The cryptographic material of a key itself cannot be changed. This operation requires the keys/update permission.
func (c *Client) UpdateKeyProperties(ctx context.Context, keyName string, options *UpdateKeyPropertiesOptions) (UpdateKeyPropertiesResponse, error) {
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
		return UpdateKeyPropertiesResponse{}, err
	}

	return updateKeyPropertiesFromGenerated(resp), nil
}

// ListDeletedKeys is the interface for the Client.ListDeletedKeys operation
type ListDeletedKeysPager interface {
	// PageResponse returns the current ListDeletedKeysPage
	PageResponse() ListDeletedKeysPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

// listDeletedKeysPager is the pager returned by Client.ListDeletedKeys
type listDeletedKeysPager struct {
	genPager *generated.KeyVaultClientGetDeletedKeysPager
}

// PageResponse returns the current page of results
func (l *listDeletedKeysPager) PageResponse() ListDeletedKeysPage {
	resp := l.genPager.PageResponse()

	var values []*DeletedKeyItem
	for _, d := range resp.Value {
		values = append(values, deletedKeyItemFromGenerated(d))
	}

	return ListDeletedKeysPage{
		RawResponse: resp.RawResponse,
		NextLink:    resp.NextLink,
		DeletedKeys: values,
	}
}

// Err returns an error if the last operation resulted in an error.
func (l *listDeletedKeysPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next page of results.
func (l *listDeletedKeysPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
}

// ListDeletedKeysPage holds the data for a single page.
type ListDeletedKeysPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; The URL to get the next set of deleted keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of the deleted keys in the vault along with a link to the next page of deleted keys
	DeletedKeys []*DeletedKeyItem `json:"value,omitempty" azure:"ro"`
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
func (c *Client) ListDeletedKeys(options *ListDeletedKeysOptions) ListDeletedKeysPager {
	if options == nil {
		options = &ListDeletedKeysOptions{}
	}

	return &listDeletedKeysPager{
		genPager: c.kvClient.GetDeletedKeys(c.vaultUrl, options.toGenerated()),
	}
}

// ListKeyVersionsPager is a Pager for Client.ListKeyVersions results
type ListKeyVersionsPager interface {
	// PageResponse returns the current ListKeyVersionsPage
	PageResponse() ListKeyVersionsPage

	// Err returns true if there is another page of data available, false if not
	Err() error

	// NextPage returns true if there is another page of data available, false if not
	NextPage(context.Context) bool
}

// listKeyVersionsPager implements the ListKeyVersionsPager interface
type listKeyVersionsPager struct {
	genPager *generated.KeyVaultClientGetKeyVersionsPager
}

// PageResponse returns the results from the page most recently fetched from the service.
func (l *listKeyVersionsPager) PageResponse() ListKeyVersionsPage {
	return listKeyVersionsPageFromGenerated(l.genPager.PageResponse())
}

// Err returns an error value if the most recent call to NextPage was not successful, else nil.
func (l *listKeyVersionsPager) Err() error {
	return l.genPager.Err()
}

// NextPage fetches the next available page of results from the service. If the fetched page
// contains results, the return value is true, else false. Results fetched from the service
// can be evaluated by calling PageResponse on this Pager.
func (l *listKeyVersionsPager) NextPage(ctx context.Context) bool {
	return l.genPager.NextPage(ctx)
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

// ListKeyVersionsPage contains the current page from a ListKeyVersionsPager.PageResponse method
type ListKeyVersionsPage struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; The URL to get the next set of keys.
	NextLink *string `json:"nextLink,omitempty" azure:"ro"`

	// READ-ONLY; A response message containing a list of keys in the key vault along with a link to the next page of keys.
	Keys []KeyItem `json:"value,omitempty" azure:"ro"`
}

// create ListKeysPage from generated pager
func listKeyVersionsPageFromGenerated(i generated.KeyVaultClientGetKeyVersionsResponse) ListKeyVersionsPage {
	var keys []KeyItem
	for _, s := range i.Value {
		if s != nil {
			keys = append(keys, *keyItemFromGenerated(s))
		}
	}
	return ListKeyVersionsPage{
		RawResponse: i.RawResponse,
		NextLink:    i.NextLink,
		Keys:        keys,
	}
}

// ListKeyVersions lists all versions of the specified key. The full key identifer and
// attributes are provided in the response. No values are returned for the keys. This operation
// requires the keys/list permission.
func (c *Client) ListKeyVersions(keyName string, options *ListKeyVersionsOptions) ListKeyVersionsPager {
	if options == nil {
		options = &ListKeyVersionsOptions{}
	}

	return &listKeyVersionsPager{
		genPager: c.kvClient.GetKeyVersions(
			c.vaultUrl,
			keyName,
			options.toGenerated(),
		),
	}
}

// RestoreKeyBackupOptions contains the optional parameters for the Client.RestoreKey method.
type RestoreKeyBackupOptions struct{}

func (r RestoreKeyBackupOptions) toGenerated() *generated.KeyVaultClientRestoreKeyOptions {
	return &generated.KeyVaultClientRestoreKeyOptions{}
}

// RestoreKeyBackupResponse contains the response object for the Client.RestoreKeyBackup operation.
type RestoreKeyBackupResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// converts the generated response to the publicly exposed version.
func restoreKeyBackupResponseFromGenerated(i generated.KeyVaultClientRestoreKeyResponse) RestoreKeyBackupResponse {
	return RestoreKeyBackupResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       convertGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// RestoreKeyBackup restores a backed up key, and all its versions, to a vault. This operation requires the keys/restore permission.
// The backup parameter is a blob of the key to restore, this can be received from the Client.BackupKey function.
func (c *Client) RestoreKeyBackup(ctx context.Context, keyBackup []byte, options *RestoreKeyBackupOptions) (RestoreKeyBackupResponse, error) {
	if options == nil {
		options = &RestoreKeyBackupOptions{}
	}

	resp, err := c.kvClient.RestoreKey(ctx, c.vaultUrl, generated.KeyRestoreParameters{KeyBundleBackup: keyBackup}, options.toGenerated())
	if err != nil {
		return RestoreKeyBackupResponse{}, err
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
		attribs = i.KeyAttributes.toGenerated()
	}
	return generated.KeyImportParameters{
		Key:           key.toGenerated(),
		Hsm:           i.Hsm,
		KeyAttributes: attribs,
		Tags:          convertToGeneratedMap(i.Tags),
	}
}

// ImportKeyResponse contains the response of the Client.ImportKey method
type ImportKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// convert the generated response to the ImportKeyResponse
func importKeyResponseFromGenerated(i generated.KeyVaultClientImportKeyResponse) ImportKeyResponse {
	return ImportKeyResponse{
		RawResponse: i.RawResponse,
		KeyBundle: KeyBundle{
			Attributes: keyAttributesFromGenerated(i.Attributes),
			Key:        jsonWebKeyFromGenerated(i.Key),
			Tags:       convertGeneratedMap(i.Tags),
			Managed:    i.Managed,
		},
	}
}

// ImportKey - The import key operation may be used to import any key type into an Azure Key Vault.
// If the named key already exists, Azure Key Vault creates a new version of the key. This operation
// requires the  keys/import permission.
func (c *Client) ImportKey(ctx context.Context, keyName string, key JSONWebKey, options *ImportKeyOptions) (ImportKeyResponse, error) {
	if options == nil {
		options = &ImportKeyOptions{}
	}

	resp, err := c.kvClient.ImportKey(ctx, c.vaultUrl, keyName, options.toImportKeyParameters(key), &generated.KeyVaultClientImportKeyOptions{})
	if err != nil {
		return ImportKeyResponse{}, err
	}

	return importKeyResponseFromGenerated(resp), nil
}

// GetRandomBytesOptions contains the optional parameters for the Client.GetRandomBytes function.
type GetRandomBytesOptions struct{}

func (g GetRandomBytesOptions) toGenerated() *generated.KeyVaultClientGetRandomBytesOptions {
	return &generated.KeyVaultClientGetRandomBytesOptions{}
}

// GetRandomBytesResponse is the response struct for the Client.GetRandomBytes function.
type GetRandomBytesResponse struct {
	// The bytes encoded as a base64url string.
	Value []byte `json:"value,omitempty"`

	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

// GetRandomBytes gets the requested number of bytes containing random values from a managed HSM.
// If the operation fails it returns the *KeyVaultError error type.
func (c *Client) GetRandomBytes(ctx context.Context, count *int32, options *GetRandomBytesOptions) (GetRandomBytesResponse, error) {
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
		return GetRandomBytesResponse{}, err
	}

	return GetRandomBytesResponse{
		Value:       resp.Value,
		RawResponse: resp.RawResponse,
	}, nil
}

type RotateKeyOptions struct{}

func (r RotateKeyOptions) toGenerated() *generated.KeyVaultClientRotateKeyOptions {
	return &generated.KeyVaultClientRotateKeyOptions{}
}

type RotateKeyResponse struct {
	KeyBundle
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func (c *Client) RotateKey(ctx context.Context, name string, options *RotateKeyOptions) (RotateKeyResponse, error) {
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
		return RotateKeyResponse{}, err
	}

	return RotateKeyResponse{
		RawResponse: resp.RawResponse,
		KeyBundle: KeyBundle{
			Attributes:    keyAttributesFromGenerated(resp.Attributes),
			Key:           jsonWebKeyFromGenerated(resp.Key),
			ReleasePolicy: keyReleasePolicyFromGenerated(resp.ReleasePolicy),
			Tags:          convertGeneratedMap(resp.Tags),
			Managed:       resp.Managed,
		},
	}, nil
}

// GetKeyRotationPolicyOptions contains the optional parameters for the Client.GetKeyRotationPolicy function
type GetKeyRotationPolicyOptions struct{}

func (g GetKeyRotationPolicyOptions) toGenerated() *generated.KeyVaultClientGetKeyRotationPolicyOptions {
	return &generated.KeyVaultClientGetKeyRotationPolicyOptions{}
}

// GetKeyRotationPolicyResponse contains the response struct for the Client.GetKeyRotationPolicy function
type GetKeyRotationPolicyResponse struct {
	KeyRotationPolicy
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func getKeyRotationPolicyResponseFromGenerated(i generated.KeyVaultClientGetKeyRotationPolicyResponse) GetKeyRotationPolicyResponse {
	var acts []*LifetimeActions
	for _, a := range i.LifetimeActions {
		acts = append(acts, lifetimeActionsFromGenerated(a))
	}
	var attribs *KeyRotationPolicyAttributes
	if i.Attributes != nil {
		attribs = &KeyRotationPolicyAttributes{
			ExpiryTime: i.Attributes.ExpiryTime,
			Created:    i.Attributes.Created,
			Updated:    i.Attributes.Updated,
		}
	}
	return GetKeyRotationPolicyResponse{
		RawResponse: i.RawResponse,
		KeyRotationPolicy: KeyRotationPolicy{
			ID:              i.ID,
			LifetimeActions: acts,
			Attributes:      attribs,
		},
	}
}

// The GetKeyRotationPolicy operation returns the specified key policy resources in the specified key vault. This operation requires
// the keys/get permission.
func (c *Client) GetKeyRotationPolicy(ctx context.Context, name string, options *GetKeyRotationPolicyOptions) (GetKeyRotationPolicyResponse, error) {
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
		return GetKeyRotationPolicyResponse{}, err
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

// ReleaseKeyResponse contains the response of Client.ReleaseKey
type ReleaseKeyResponse struct {
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response

	// READ-ONLY; A signed object containing the released key.
	Value *string `json:"value,omitempty" azure:"ro"`
}

func (c *Client) ReleaseKey(ctx context.Context, name string, target string, options *ReleaseKeyOptions) (ReleaseKeyResponse, error) {
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
		return ReleaseKeyResponse{}, err
	}

	return ReleaseKeyResponse{
		RawResponse: resp.RawResponse,
		Value:       resp.Value,
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
		attribs = u.Attributes.toGenerated()
	}
	var la []*generated.LifetimeActions
	for _, l := range u.LifetimeActions {
		if l == nil {
			la = append(la, nil)
		} else {
			la = append(la, l.toGenerated())
		}
	}

	return generated.KeyRotationPolicy{
		ID:              u.ID,
		LifetimeActions: la,
		Attributes:      attribs,
	}
}

// UpdateKeyRotationPolicyResponse contains the response for the Client.UpdateKeyRotationPolicy function
type UpdateKeyRotationPolicyResponse struct {
	KeyRotationPolicy
	// RawResponse contains the underlying HTTP response.
	RawResponse *http.Response
}

func updateKeyRotationPolicyResponseFromGenerated(i generated.KeyVaultClientUpdateKeyRotationPolicyResponse) UpdateKeyRotationPolicyResponse {
	var acts []*LifetimeActions
	for _, a := range i.LifetimeActions {
		acts = append(acts, lifetimeActionsFromGenerated(a))
	}
	var attribs *KeyRotationPolicyAttributes
	if i.Attributes != nil {
		attribs = &KeyRotationPolicyAttributes{
			ExpiryTime: i.Attributes.ExpiryTime,
			Created:    i.Attributes.Created,
			Updated:    i.Attributes.Updated,
		}
	}
	return UpdateKeyRotationPolicyResponse{
		RawResponse: i.RawResponse,
		KeyRotationPolicy: KeyRotationPolicy{
			ID:              i.ID,
			LifetimeActions: acts,
			Attributes:      attribs,
		},
	}
}

func (c *Client) UpdateKeyRotationPolicy(ctx context.Context, name string, options *UpdateKeyRotationPolicyOptions) (UpdateKeyRotationPolicyResponse, error) {
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
		return UpdateKeyRotationPolicyResponse{}, err
	}

	return updateKeyRotationPolicyResponseFromGenerated(resp), nil
}
