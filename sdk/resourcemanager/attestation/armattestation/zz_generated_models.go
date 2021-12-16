//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armattestation

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
	"time"
)

// AttestationProvider - Attestation service response message.
type AttestationProvider struct {
	TrackedResource
	// Describes Attestation service status.
	Properties *StatusResult `json:"properties,omitempty"`

	// READ-ONLY; The system metadata relating to this resource
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type AttestationProvider.
func (a AttestationProvider) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	a.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "properties", a.Properties)
	populate(objectMap, "systemData", a.SystemData)
	return json.Marshal(objectMap)
}

// AttestationProviderListResult - Attestation Providers List.
type AttestationProviderListResult struct {
	// Attestation Provider array.
	Value []*AttestationProvider `json:"value,omitempty"`

	// READ-ONLY; The system metadata relating to this resource
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type AttestationProviderListResult.
func (a AttestationProviderListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "systemData", a.SystemData)
	populate(objectMap, "value", a.Value)
	return json.Marshal(objectMap)
}

// AttestationProvidersCreateOptions contains the optional parameters for the AttestationProviders.Create method.
type AttestationProvidersCreateOptions struct {
	// placeholder for future optional parameters
}

// AttestationProvidersDeleteOptions contains the optional parameters for the AttestationProviders.Delete method.
type AttestationProvidersDeleteOptions struct {
	// placeholder for future optional parameters
}

// AttestationProvidersGetDefaultByLocationOptions contains the optional parameters for the AttestationProviders.GetDefaultByLocation method.
type AttestationProvidersGetDefaultByLocationOptions struct {
	// placeholder for future optional parameters
}

// AttestationProvidersGetOptions contains the optional parameters for the AttestationProviders.Get method.
type AttestationProvidersGetOptions struct {
	// placeholder for future optional parameters
}

// AttestationProvidersListByResourceGroupOptions contains the optional parameters for the AttestationProviders.ListByResourceGroup method.
type AttestationProvidersListByResourceGroupOptions struct {
	// placeholder for future optional parameters
}

// AttestationProvidersListDefaultOptions contains the optional parameters for the AttestationProviders.ListDefault method.
type AttestationProvidersListDefaultOptions struct {
	// placeholder for future optional parameters
}

// AttestationProvidersListOptions contains the optional parameters for the AttestationProviders.List method.
type AttestationProvidersListOptions struct {
	// placeholder for future optional parameters
}

// AttestationProvidersUpdateOptions contains the optional parameters for the AttestationProviders.Update method.
type AttestationProvidersUpdateOptions struct {
	// placeholder for future optional parameters
}

// AttestationServiceCreationParams - Parameters for creating an attestation provider
type AttestationServiceCreationParams struct {
	// REQUIRED; The supported Azure location where the attestation provider should be created.
	Location *string `json:"location,omitempty"`

	// REQUIRED; Properties of the attestation provider
	Properties *AttestationServiceCreationSpecificParams `json:"properties,omitempty"`

	// The tags that will be assigned to the attestation provider.
	Tags map[string]*string `json:"tags,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type AttestationServiceCreationParams.
func (a AttestationServiceCreationParams) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "location", a.Location)
	populate(objectMap, "properties", a.Properties)
	populate(objectMap, "tags", a.Tags)
	return json.Marshal(objectMap)
}

// AttestationServiceCreationSpecificParams - Client supplied parameters used to create a new attestation provider.
type AttestationServiceCreationSpecificParams struct {
	// JSON Web Key Set defining a set of X.509 Certificates that will represent the parent certificate for the signing certificate used for policy operations
	PolicySigningCertificates *JSONWebKeySet `json:"policySigningCertificates,omitempty"`
}

// AttestationServicePatchParams - Parameters for patching an attestation provider
type AttestationServicePatchParams struct {
	// The tags that will be assigned to the attestation provider.
	Tags map[string]*string `json:"tags,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type AttestationServicePatchParams.
func (a AttestationServicePatchParams) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "tags", a.Tags)
	return json.Marshal(objectMap)
}

// CloudError - An error response from Attestation.
// Implements the error and azcore.HTTPResponse interfaces.
type CloudError struct {
	raw string
	// An error response from Attestation.
	InnerError *CloudErrorBody `json:"error,omitempty"`
}

// Error implements the error interface for type CloudError.
// The contents of the error text are not contractual and subject to change.
func (e CloudError) Error() string {
	return e.raw
}

// CloudErrorBody - An error response from Attestation.
type CloudErrorBody struct {
	// An identifier for the error. Codes are invariant and are intended to be consumed programmatically.
	Code *string `json:"code,omitempty"`

	// A message describing the error, intended to be suitable for displaying in a user interface.
	Message *string `json:"message,omitempty"`
}

type JSONWebKey struct {
	// REQUIRED; The "kty" (key type) parameter identifies the cryptographic algorithm family used with the key, such as "RSA" or "EC". "kty" values should
	// either be registered in the IANA "JSON Web Key Types"
	// registry established by [JWA] or be a value that contains a Collision- Resistant Name. The "kty" value is a case-sensitive string.
	Kty *string `json:"kty,omitempty"`

	// The "alg" (algorithm) parameter identifies the algorithm intended for use with the key. The values used should either be registered in the IANA "JSON
	// Web Signature and Encryption Algorithms" registry
	// established by [JWA] or be a value that contains a Collision- Resistant Name.
	Alg *string `json:"alg,omitempty"`

	// The "crv" (curve) parameter identifies the curve type
	Crv *string `json:"crv,omitempty"`

	// RSA private exponent or ECC private key
	D *string `json:"d,omitempty"`

	// RSA Private Key Parameter
	Dp *string `json:"dp,omitempty"`

	// RSA Private Key Parameter
	Dq *string `json:"dq,omitempty"`

	// RSA public exponent, in Base64
	E *string `json:"e,omitempty"`

	// Symmetric key
	K *string `json:"k,omitempty"`

	// The "kid" (key ID) parameter is used to match a specific key. This is used, for instance, to choose among a set of keys within a JWK Set during key rollover.
	// The structure of the "kid" value is
	// unspecified. When "kid" values are used within a JWK Set, different keys within the JWK Set SHOULD use distinct "kid" values. (One example in which different
	// keys might use the same "kid" value is if
	// they have different "kty" (key type) values but are considered to be equivalent alternatives by the application using them.) The "kid" value is a case-sensitive
	// string.
	Kid *string `json:"kid,omitempty"`

	// RSA modulus, in Base64
	N *string `json:"n,omitempty"`

	// RSA secret prime
	P *string `json:"p,omitempty"`

	// RSA secret prime, with p < q
	Q *string `json:"q,omitempty"`

	// RSA Private Key Parameter
	Qi *string `json:"qi,omitempty"`

	// Use ("public key use") identifies the intended use of the public key. The "use" parameter is employed to indicate whether a public key is used for encrypting
	// data or verifying the signature on data.
	// Values are commonly "sig" (signature) or "enc" (encryption).
	Use *string `json:"use,omitempty"`

	// X coordinate for the Elliptic Curve point
	X *string `json:"x,omitempty"`

	// The "x5c" (X.509 certificate chain) parameter contains a chain of one or more PKIX certificates [RFC5280]. The certificate chain is represented as a
	// JSON array of certificate value strings. Each
	// string in the array is a base64-encoded (Section 4 of [RFC4648] -- not base64url-encoded) DER [ITU.X690.1994] PKIX certificate value. The PKIX certificate
	// containing the key value MUST be the first
	// certificate.
	X5C []*string `json:"x5c,omitempty"`

	// Y coordinate for the Elliptic Curve point
	Y *string `json:"y,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type JSONWebKey.
func (j JSONWebKey) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "alg", j.Alg)
	populate(objectMap, "crv", j.Crv)
	populate(objectMap, "d", j.D)
	populate(objectMap, "dp", j.Dp)
	populate(objectMap, "dq", j.Dq)
	populate(objectMap, "e", j.E)
	populate(objectMap, "k", j.K)
	populate(objectMap, "kid", j.Kid)
	populate(objectMap, "kty", j.Kty)
	populate(objectMap, "n", j.N)
	populate(objectMap, "p", j.P)
	populate(objectMap, "q", j.Q)
	populate(objectMap, "qi", j.Qi)
	populate(objectMap, "use", j.Use)
	populate(objectMap, "x", j.X)
	populate(objectMap, "x5c", j.X5C)
	populate(objectMap, "y", j.Y)
	return json.Marshal(objectMap)
}

type JSONWebKeySet struct {
	// The value of the "keys" parameter is an array of JWK values. By default, the order of the JWK values within the array does not imply an order of preference
	// among them, although applications of JWK
	// Sets can choose to assign a meaning to the order for their purposes, if desired.
	Keys []*JSONWebKey `json:"keys,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type JSONWebKeySet.
func (j JSONWebKeySet) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "keys", j.Keys)
	return json.Marshal(objectMap)
}

// OperationList - List of supported operations.
type OperationList struct {
	// List of supported operations.
	Value []*OperationsDefinition `json:"value,omitempty"`

	// READ-ONLY; The system metadata relating to this resource
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type OperationList.
func (o OperationList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "systemData", o.SystemData)
	populate(objectMap, "value", o.Value)
	return json.Marshal(objectMap)
}

// OperationsDefinition - Definition object with the name and properties of an operation.
type OperationsDefinition struct {
	// Display object with properties of the operation.
	Display *OperationsDisplayDefinition `json:"display,omitempty"`

	// Name of the operation.
	Name *string `json:"name,omitempty"`
}

// OperationsDisplayDefinition - Display object with properties of the operation.
type OperationsDisplayDefinition struct {
	// Description of the operation.
	Description *string `json:"description,omitempty"`

	// Short description of the operation.
	Operation *string `json:"operation,omitempty"`

	// Resource provider of the operation.
	Provider *string `json:"provider,omitempty"`

	// Resource for the operation.
	Resource *string `json:"resource,omitempty"`
}

// OperationsListOptions contains the optional parameters for the Operations.List method.
type OperationsListOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpoint - The Private Endpoint resource.
type PrivateEndpoint struct {
	// READ-ONLY; The ARM identifier for Private Endpoint
	ID *string `json:"id,omitempty" azure:"ro"`
}

// PrivateEndpointConnection - The Private Endpoint Connection resource.
type PrivateEndpointConnection struct {
	Resource
	// Resource properties.
	Properties *PrivateEndpointConnectionProperties `json:"properties,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type PrivateEndpointConnection.
func (p PrivateEndpointConnection) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	p.Resource.marshalInternal(objectMap)
	populate(objectMap, "properties", p.Properties)
	return json.Marshal(objectMap)
}

// PrivateEndpointConnectionListResult - List of private endpoint connection associated with the specified storage account
type PrivateEndpointConnectionListResult struct {
	// Array of private endpoint connections
	Value []*PrivateEndpointConnection `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type PrivateEndpointConnectionListResult.
func (p PrivateEndpointConnectionListResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "value", p.Value)
	return json.Marshal(objectMap)
}

// PrivateEndpointConnectionProperties - Properties of the PrivateEndpointConnectProperties.
type PrivateEndpointConnectionProperties struct {
	// REQUIRED; A collection of information about the state of the connection between service consumer and provider.
	PrivateLinkServiceConnectionState *PrivateLinkServiceConnectionState `json:"privateLinkServiceConnectionState,omitempty"`

	// The resource of private end point.
	PrivateEndpoint *PrivateEndpoint `json:"privateEndpoint,omitempty"`

	// READ-ONLY; The provisioning state of the private endpoint connection resource.
	ProvisioningState *PrivateEndpointConnectionProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// PrivateEndpointConnectionsCreateOptions contains the optional parameters for the PrivateEndpointConnections.Create method.
type PrivateEndpointConnectionsCreateOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsDeleteOptions contains the optional parameters for the PrivateEndpointConnections.Delete method.
type PrivateEndpointConnectionsDeleteOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsGetOptions contains the optional parameters for the PrivateEndpointConnections.Get method.
type PrivateEndpointConnectionsGetOptions struct {
	// placeholder for future optional parameters
}

// PrivateEndpointConnectionsListOptions contains the optional parameters for the PrivateEndpointConnections.List method.
type PrivateEndpointConnectionsListOptions struct {
	// placeholder for future optional parameters
}

// PrivateLinkServiceConnectionState - A collection of information about the state of the connection between service consumer and provider.
type PrivateLinkServiceConnectionState struct {
	// A message indicating if changes on the service provider require any updates on the consumer.
	ActionsRequired *string `json:"actionsRequired,omitempty"`

	// The reason for approval/rejection of the connection.
	Description *string `json:"description,omitempty"`

	// Indicates whether the connection has been Approved/Rejected/Removed by the owner of the service.
	Status *PrivateEndpointServiceConnectionStatus `json:"status,omitempty"`
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	r.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (r Resource) marshalInternal(objectMap map[string]interface{}) {
	populate(objectMap, "id", r.ID)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "type", r.Type)
}

// StatusResult - Status of attestation service.
type StatusResult struct {
	// Gets the uri of attestation service
	AttestURI *string `json:"attestUri,omitempty"`

	// Status of attestation service.
	Status *AttestationServiceStatus `json:"status,omitempty"`

	// Trust model for the attestation provider.
	TrustModel *string `json:"trustModel,omitempty"`

	// READ-ONLY; List of private endpoint connections associated with the attestation provider.
	PrivateEndpointConnections []*PrivateEndpointConnection `json:"privateEndpointConnections,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type StatusResult.
func (s StatusResult) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "attestUri", s.AttestURI)
	populate(objectMap, "privateEndpointConnections", s.PrivateEndpointConnections)
	populate(objectMap, "status", s.Status)
	populate(objectMap, "trustModel", s.TrustModel)
	return json.Marshal(objectMap)
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// The identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`

	// The type of identity that created the resource.
	CreatedByType *CreatedByType `json:"createdByType,omitempty"`

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`

	// The identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type SystemData.
func (s SystemData) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populateTimeRFC3339(objectMap, "createdAt", s.CreatedAt)
	populate(objectMap, "createdBy", s.CreatedBy)
	populate(objectMap, "createdByType", s.CreatedByType)
	populateTimeRFC3339(objectMap, "lastModifiedAt", s.LastModifiedAt)
	populate(objectMap, "lastModifiedBy", s.LastModifiedBy)
	populate(objectMap, "lastModifiedByType", s.LastModifiedByType)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type SystemData.
func (s *SystemData) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "createdAt":
			err = unpopulateTimeRFC3339(val, &s.CreatedAt)
			delete(rawMsg, key)
		case "createdBy":
			err = unpopulate(val, &s.CreatedBy)
			delete(rawMsg, key)
		case "createdByType":
			err = unpopulate(val, &s.CreatedByType)
			delete(rawMsg, key)
		case "lastModifiedAt":
			err = unpopulateTimeRFC3339(val, &s.LastModifiedAt)
			delete(rawMsg, key)
		case "lastModifiedBy":
			err = unpopulate(val, &s.LastModifiedBy)
			delete(rawMsg, key)
		case "lastModifiedByType":
			err = unpopulate(val, &s.LastModifiedByType)
			delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags' and a 'location'
type TrackedResource struct {
	Resource
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type TrackedResource.
func (t TrackedResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	t.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (t TrackedResource) marshalInternal(objectMap map[string]interface{}) {
	t.Resource.marshalInternal(objectMap)
	populate(objectMap, "location", t.Location)
	populate(objectMap, "tags", t.Tags)
}

func populate(m map[string]interface{}, k string, v interface{}) {
	if v == nil {
		return
	} else if azcore.IsNullValue(v) {
		m[k] = nil
	} else if !reflect.ValueOf(v).IsNil() {
		m[k] = v
	}
}

func unpopulate(data json.RawMessage, v interface{}) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, v)
}
