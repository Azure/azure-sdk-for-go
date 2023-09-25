//go:build go1.18
// +build go1.18

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package armedgeorder

import "time"

// AddressDetails - Address details for an order item.
type AddressDetails struct {
	// REQUIRED; Customer address and contact details. It should be address resource
	ForwardAddress *AddressProperties

	// READ-ONLY; Return shipping address
	ReturnAddress *AddressProperties
}

// AddressProperties - Address Properties
type AddressProperties struct {
	// REQUIRED; Contact details for the address
	ContactDetails *ContactDetails

	// Shipping details for the address
	ShippingAddress *ShippingAddress

	// READ-ONLY; Status of address validation
	AddressValidationStatus *AddressValidationStatus
}

// AddressResource - Address Resource.
type AddressResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// REQUIRED; Properties of an address.
	Properties *AddressProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Represents resource creation and update time
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// AddressResourceList - Address Resource Collection
type AddressResourceList struct {
	// Link for the next set of job resources.
	NextLink *string

	// READ-ONLY; List of address resources.
	Value []*AddressResource
}

// AddressUpdateParameter - The Address update parameters
type AddressUpdateParameter struct {
	// Properties of a address to be updated.
	Properties *AddressUpdateProperties

	// The list of key value pairs that describe the resource. These tags can be used in viewing and grouping this resource (across
// resource groups).
	Tags map[string]*string
}

// AddressUpdateProperties - Address Properties
type AddressUpdateProperties struct {
	// Contact details for the address
	ContactDetails *ContactDetails

	// Shipping details for the address
	ShippingAddress *ShippingAddress
}

// AvailabilityInformation - Availability information of a product system.
type AvailabilityInformation struct {
	// READ-ONLY; Current availability stage of the product. Availability stage
	AvailabilityStage *AvailabilityStage

	// READ-ONLY; Reason why the product is disabled.
	DisabledReason *DisabledReason

	// READ-ONLY; Message for why the product is disabled.
	DisabledReasonMessage *string
}

// BasicInformation - Basic information for any product system
type BasicInformation struct {
	// READ-ONLY; Availability information of the product system.
	AvailabilityInformation *AvailabilityInformation

	// READ-ONLY; Cost information for the product system.
	CostInformation *CostInformation

	// READ-ONLY; Description related to the product system.
	Description *Description

	// READ-ONLY; Display Name for the product system.
	DisplayName *string

	// READ-ONLY; Hierarchy information of a product.
	HierarchyInformation *HierarchyInformation

	// READ-ONLY; Image information for the product system.
	ImageInformation []*ImageInformation
}

// BillingMeterDetails - Holds billing meter details for each type of billing
type BillingMeterDetails struct {
	// READ-ONLY; Frequency of recurrence
	Frequency *string

	// READ-ONLY; Represents MeterDetails
	MeterDetails MeterDetailsClassification

	// READ-ONLY; Represents Metering type (eg one-time or recurrent)
	MeteringType *MeteringType

	// READ-ONLY; Represents Billing type name
	Name *string
}

// CancellationReason - Reason for cancellation.
type CancellationReason struct {
	// REQUIRED; Reason for cancellation.
	Reason *string
}

// CommonProperties - Represents common properties across product hierarchy
type CommonProperties struct {
	// READ-ONLY; Availability information of the product system.
	AvailabilityInformation *AvailabilityInformation

	// READ-ONLY; Cost information for the product system.
	CostInformation *CostInformation

	// READ-ONLY; Description related to the product system.
	Description *Description

	// READ-ONLY; Display Name for the product system.
	DisplayName *string

	// READ-ONLY; list of filters supported for a product
	FilterableProperties []*FilterableProperty

	// READ-ONLY; Hierarchy information of a product.
	HierarchyInformation *HierarchyInformation

	// READ-ONLY; Image information for the product system.
	ImageInformation []*ImageInformation
}

// Configuration object.
type Configuration struct {
	// READ-ONLY; Properties of configuration
	Properties *ConfigurationProperties
}

// ConfigurationFilters - Configuration filters
type ConfigurationFilters struct {
	// REQUIRED; Product hierarchy information
	HierarchyInformation *HierarchyInformation

	// Filters specific to product
	FilterableProperty []*FilterableProperty
}

// ConfigurationProperties - Properties of configuration
type ConfigurationProperties struct {
	// READ-ONLY; Availability information of the product system.
	AvailabilityInformation *AvailabilityInformation

	// READ-ONLY; Cost information for the product system.
	CostInformation *CostInformation

	// READ-ONLY; Description related to the product system.
	Description *Description

	// READ-ONLY; Dimensions of the configuration
	Dimensions *Dimensions

	// READ-ONLY; Display Name for the product system.
	DisplayName *string

	// READ-ONLY; list of filters supported for a product
	FilterableProperties []*FilterableProperty

	// READ-ONLY; Hierarchy information of a product.
	HierarchyInformation *HierarchyInformation

	// READ-ONLY; Image information for the product system.
	ImageInformation []*ImageInformation

	// READ-ONLY; Specifications of the configuration
	Specifications []*Specification
}

// Configurations - The list of configurations.
type Configurations struct {
	// Link for the next set of configurations.
	NextLink *string

	// READ-ONLY; List of configurations.
	Value []*Configuration
}

// ConfigurationsRequest - Configuration request object.
type ConfigurationsRequest struct {
	// REQUIRED; Holds details about product hierarchy information and filterable property.
	ConfigurationFilters []*ConfigurationFilters

	// Customer subscription properties. Clients can display available products to unregistered customers by explicitly passing
// subscription details
	CustomerSubscriptionDetails *CustomerSubscriptionDetails
}

// ContactDetails - Contact Details.
type ContactDetails struct {
	// REQUIRED; Contact name of the person.
	ContactName *string

	// REQUIRED; List of Email-ids to be notified about job progress.
	EmailList []*string

	// REQUIRED; Phone number of the contact person.
	Phone *string

	// Mobile number of the contact person.
	Mobile *string

	// Phone extension number of the contact person.
	PhoneExtension *string
}

// CostInformation - Cost information for the product system
type CostInformation struct {
	// READ-ONLY; Default url to display billing information
	BillingInfoURL *string

	// READ-ONLY; Details on the various billing aspects for the product system.
	BillingMeterDetails []*BillingMeterDetails
}

// CustomerSubscriptionDetails - Holds Customer subscription details. Clients can display available products to unregistered
// customers by explicitly passing subscription details
type CustomerSubscriptionDetails struct {
	// REQUIRED; Quota ID of a subscription
	QuotaID *string

	// Location placement Id of a subscription
	LocationPlacementID *string

	// List of registered feature flags for subscription
	RegisteredFeatures []*CustomerSubscriptionRegisteredFeatures
}

// CustomerSubscriptionRegisteredFeatures - Represents subscription registered features
type CustomerSubscriptionRegisteredFeatures struct {
	// Name of subscription registered feature
	Name *string

	// State of subscription registered feature
	State *string
}

// Description related properties of a product system.
type Description struct {
	// READ-ONLY; Attributes for the product system.
	Attributes []*string

	// READ-ONLY; Type of description.
	DescriptionType *DescriptionType

	// READ-ONLY; Keywords for the product system.
	Keywords []*string

	// READ-ONLY; Links for the product system.
	Links []*Link

	// READ-ONLY; Long description of the product system.
	LongDescription *string

	// READ-ONLY; Short description of the product system.
	ShortDescription *string
}

// DeviceDetails - Device details.
type DeviceDetails struct {
	// READ-ONLY; Management Resource Id
	ManagementResourceID *string

	// READ-ONLY; Management Resource Tenant ID
	ManagementResourceTenantID *string

	// READ-ONLY; device serial number
	SerialNumber *string
}

// Dimensions of a configuration.
type Dimensions struct {
	// READ-ONLY; Depth of the device.
	Depth *float64

	// READ-ONLY; Height of the device.
	Height *float64

	// READ-ONLY; Length of the device.
	Length *float64

	// READ-ONLY; Unit for the dimensions of length, height and width.
	LengthHeightUnit *LengthHeightUnit

	// READ-ONLY; Weight of the device.
	Weight *float64

	// READ-ONLY; Unit for the dimensions of weight.
	WeightUnit *WeightMeasurementUnit

	// READ-ONLY; Width of the device.
	Width *float64
}

// DisplayInfo - Describes product display information
type DisplayInfo struct {
	// READ-ONLY; Configuration display name
	ConfigurationDisplayName *string

	// READ-ONLY; Product family display name
	ProductFamilyDisplayName *string
}

// EncryptionPreferences - Preferences related to the double encryption
type EncryptionPreferences struct {
	// Double encryption status as entered by the customer. It is compulsory to give this parameter if the 'Deny' or 'Disabled'
// policy is configured.
	DoubleEncryptionStatus *DoubleEncryptionStatus
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info any

	// READ-ONLY; The additional info type.
	Type *string
}

// ErrorDetail - The error detail.
type ErrorDetail struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo

	// READ-ONLY; The error code.
	Code *string

	// READ-ONLY; The error details.
	Details []*ErrorDetail

	// READ-ONLY; The error message.
	Message *string

	// READ-ONLY; The error target.
	Target *string
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations.
// (This also follows the OData error response format.).
type ErrorResponse struct {
	// The error object.
	Error *ErrorDetail
}

// FilterableProperty - Different types of filters supported and its values.
type FilterableProperty struct {
	// REQUIRED; Values to be filtered.
	SupportedValues []*string

	// REQUIRED; Type of product filter.
	Type *SupportedFilterTypes
}

// ForwardShippingDetails - Forward shipment details.
type ForwardShippingDetails struct {
	// READ-ONLY; Carrier Name for display purpose. Not to be used for any processing.
	CarrierDisplayName *string

	// READ-ONLY; Name of the carrier.
	CarrierName *string

	// READ-ONLY; TrackingId of the package
	TrackingID *string

	// READ-ONLY; TrackingUrl of the package.
	TrackingURL *string
}

// HierarchyInformation - Holds details about product hierarchy information
type HierarchyInformation struct {
	// Represents configuration name that uniquely identifies configuration
	ConfigurationName *string

	// Represents product family name that uniquely identifies product family
	ProductFamilyName *string

	// Represents product line name that uniquely identifies product line
	ProductLineName *string

	// Represents product name that uniquely identifies product
	ProductName *string
}

// ImageInformation - Image for the product
type ImageInformation struct {
	// READ-ONLY; Type of the image
	ImageType *ImageType

	// READ-ONLY; Url of the image
	ImageURL *string
}

// Link - Returns link related to the product
type Link struct {
	// READ-ONLY; Type of link
	LinkType *LinkType

	// READ-ONLY; Url of the link
	LinkURL *string
}

// ManagementResourcePreferences - Management resource preference to link device
type ManagementResourcePreferences struct {
	// Customer preferred Management resource ARM ID
	PreferredManagementResourceID *string
}

// MeterDetails - Holds details about billing type and its meter guids
type MeterDetails struct {
	// REQUIRED; Represents billing type.
	BillingType *BillingType

	// READ-ONLY; Charging type.
	ChargingType *ChargingType

	// READ-ONLY; Billing unit applicable for Pav2 billing
	Multiplier *float64
}

// GetMeterDetails implements the MeterDetailsClassification interface for type MeterDetails.
func (m *MeterDetails) GetMeterDetails() *MeterDetails { return m }

// NotificationPreference - Notification preference for a job stage.
type NotificationPreference struct {
	// REQUIRED; Notification is required or not.
	SendNotification *bool

	// REQUIRED; Name of the stage.
	StageName *NotificationStageName
}

// Operation - Details of a REST API operation, returned from the Resource Provider Operations API
type Operation struct {
	// Localized display information for this particular operation.
	Display *OperationDisplay

	// READ-ONLY; Enum. Indicates the action type. "Internal" refers to actions that are for internal only APIs.
	ActionType *ActionType

	// READ-ONLY; Whether the operation applies to data-plane. This is "true" for data-plane operations and "false" for ARM/control-plane
// operations.
	IsDataAction *bool

	// READ-ONLY; The name of the operation, as per Resource-Based Access Control (RBAC). Examples: "Microsoft.Compute/virtualMachines/write",
// "Microsoft.Compute/virtualMachines/capture/action"
	Name *string

	// READ-ONLY; The intended executor of the operation; as in Resource Based Access Control (RBAC) and audit logs UX. Default
// value is "user,system"
	Origin *Origin
}

// OperationDisplay - Localized display information for this particular operation.
type OperationDisplay struct {
	// READ-ONLY; The short, localized friendly description of the operation; suitable for tool tips and detailed views.
	Description *string

	// READ-ONLY; The concise, localized friendly name for the operation; suitable for dropdowns. E.g. "Create or Update Virtual
// Machine", "Restart Virtual Machine".
	Operation *string

	// READ-ONLY; The localized friendly form of the resource provider name, e.g. "Microsoft Monitoring Insights" or "Microsoft
// Compute".
	Provider *string

	// READ-ONLY; The localized friendly name of the resource type related to this operation. E.g. "Virtual Machines" or "Job
// Schedule Collections".
	Resource *string
}

// OperationListResult - A list of REST API operations supported by an Azure Resource Provider. It contains an URL link to
// get the next set of results.
type OperationListResult struct {
	// READ-ONLY; URL to get the next set of operation list results (if there are any).
	NextLink *string

	// READ-ONLY; List of operations supported by the resource provider
	Value []*Operation
}

// OrderItemDetails - Order item details
type OrderItemDetails struct {
	// REQUIRED; Order item type.
	OrderItemType *OrderItemType

	// REQUIRED; Unique identifier for configuration.
	ProductDetails *ProductDetails

	// Additional notification email list
	NotificationEmailList []*string

	// Customer notification Preferences
	Preferences *Preferences

	// READ-ONLY; Cancellation reason.
	CancellationReason *string

	// READ-ONLY; Describes whether the order item is cancellable or not.
	CancellationStatus *OrderItemCancellationEnum

	// READ-ONLY; Current Order item Status
	CurrentStage *StageDetails

	// READ-ONLY; Describes whether the order item is deletable or not.
	DeletionStatus *ActionStatusEnum

	// READ-ONLY; Top level error for the job.
	Error *ErrorDetail

	// READ-ONLY; Forward Package Shipping details
	ForwardShippingDetails *ForwardShippingDetails

	// READ-ONLY; Parent RP details - this returns only the first or default parent RP from the entire list
	ManagementRpDetails *ResourceProviderDetails

	// READ-ONLY; List of parent RP details supported for configuration.
	ManagementRpDetailsList []*ResourceProviderDetails

	// READ-ONLY; Order item status history
	OrderItemStageHistory []*StageDetails

	// READ-ONLY; Return reason.
	ReturnReason *string

	// READ-ONLY; Describes whether the order item is returnable or not.
	ReturnStatus *OrderItemReturnEnum

	// READ-ONLY; Reverse Package Shipping details
	ReverseShippingDetails *ReverseShippingDetails
}

// OrderItemProperties - Represents order item details.
type OrderItemProperties struct {
	// REQUIRED; Represents shipping and return address for order item
	AddressDetails *AddressDetails

	// REQUIRED; Id of the order to which order item belongs to
	OrderID *string

	// REQUIRED; Represents order item details.
	OrderItemDetails *OrderItemDetails

	// READ-ONLY; Start time of order item
	StartTime *time.Time
}

// OrderItemResource - Represents order item contract
type OrderItemResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// REQUIRED; Order item properties
	Properties *OrderItemProperties

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Represents resource creation and update time
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// OrderItemResourceList - List of orderItems.
type OrderItemResourceList struct {
	// Link for the next set of order item resources.
	NextLink *string

	// READ-ONLY; List of order item resources.
	Value []*OrderItemResource
}

// OrderItemUpdateParameter - Updates order item parameters.
type OrderItemUpdateParameter struct {
	// Order item update properties
	Properties *OrderItemUpdateProperties

	// The list of key value pairs that describe the resource. These tags can be used in viewing and grouping this resource (across
// resource groups).
	Tags map[string]*string
}

// OrderItemUpdateProperties - Order item update properties.
type OrderItemUpdateProperties struct {
	// Updates forward shipping address and contact details.
	ForwardAddress *AddressProperties

	// Additional notification email list.
	NotificationEmailList []*string

	// Customer preference.
	Preferences *Preferences
}

// OrderProperties - Represents order details.
type OrderProperties struct {
	// READ-ONLY; Order current status.
	CurrentStage *StageDetails

	// READ-ONLY; List of order item ARM Ids which are part of an order.
	OrderItemIDs []*string

	// READ-ONLY; Order status history.
	OrderStageHistory []*StageDetails
}

// OrderResource - Specifies the properties or parameters for an order. Order is a grouping of one or more order items.
type OrderResource struct {
	// REQUIRED; Order properties
	Properties *OrderProperties

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; Represents resource creation and update time
	SystemData *SystemData

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// OrderResourceList - List of orders.
type OrderResourceList struct {
	// Link for the next set of order resources.
	NextLink *string

	// READ-ONLY; List of order resources.
	Value []*OrderResource
}

// Pav2MeterDetails - Billing type PAV2 meter details
type Pav2MeterDetails struct {
	// REQUIRED; Represents billing type.
	BillingType *BillingType

	// READ-ONLY; Charging type.
	ChargingType *ChargingType

	// READ-ONLY; Validation status of requested data center and transport.
	MeterGUID *string

	// READ-ONLY; Billing unit applicable for Pav2 billing
	Multiplier *float64
}

// GetMeterDetails implements the MeterDetailsClassification interface for type Pav2MeterDetails.
func (p *Pav2MeterDetails) GetMeterDetails() *MeterDetails {
	return &MeterDetails{
		BillingType: p.BillingType,
		ChargingType: p.ChargingType,
		Multiplier: p.Multiplier,
	}
}

// Preferences related to the order
type Preferences struct {
	// Preferences related to the Encryption.
	EncryptionPreferences *EncryptionPreferences

	// Preferences related to the Management resource.
	ManagementResourcePreferences *ManagementResourcePreferences

	// Notification preferences.
	NotificationPreferences []*NotificationPreference

	// Preferences related to the shipment logistics of the order.
	TransportPreferences *TransportPreferences
}

// Product - List of Products
type Product struct {
	// READ-ONLY; Properties of product
	Properties *ProductProperties
}

// ProductDetails - Represents product details
type ProductDetails struct {
	// REQUIRED; Hierarchy of the product which uniquely identifies the product
	HierarchyInformation *HierarchyInformation

	// Display details of the product
	DisplayInfo *DisplayInfo

	// READ-ONLY; Quantity of the product
	Count *int32

	// READ-ONLY; list of device details
	DeviceDetails []*DeviceDetails

	// READ-ONLY; Double encryption status of the configuration. Read-only field.
	ProductDoubleEncryptionStatus *DoubleEncryptionStatus
}

// ProductFamilies - The list of product families.
type ProductFamilies struct {
	// Link for the next set of product families.
	NextLink *string

	// READ-ONLY; List of product families.
	Value []*ProductFamily
}

// ProductFamiliesMetadata - Holds details about product family metadata
type ProductFamiliesMetadata struct {
	// READ-ONLY; Link for the next set of product families.
	NextLink *string

	// READ-ONLY; List of product family metadata details.
	Value []*ProductFamiliesMetadataDetails
}

// ProductFamiliesMetadataDetails - Product families metadata details.
type ProductFamiliesMetadataDetails struct {
	// READ-ONLY; Product family properties
	Properties *ProductFamilyProperties
}

// ProductFamiliesRequest - The filters for showing the product families.
type ProductFamiliesRequest struct {
	// REQUIRED; Dictionary of filterable properties on product family.
	FilterableProperties map[string][]*FilterableProperty

	// Customer subscription properties. Clients can display available products to unregistered customers by explicitly passing
// subscription details
	CustomerSubscriptionDetails *CustomerSubscriptionDetails
}

// ProductFamily - Product Family
type ProductFamily struct {
	// READ-ONLY; Properties of product family
	Properties *ProductFamilyProperties
}

// ProductFamilyProperties - Properties of product family
type ProductFamilyProperties struct {
	// Contains details related to resource provider
	ResourceProviderDetails []*ResourceProviderDetails

	// READ-ONLY; Availability information of the product system.
	AvailabilityInformation *AvailabilityInformation

	// READ-ONLY; Cost information for the product system.
	CostInformation *CostInformation

	// READ-ONLY; Description related to the product system.
	Description *Description

	// READ-ONLY; Display Name for the product system.
	DisplayName *string

	// READ-ONLY; list of filters supported for a product
	FilterableProperties []*FilterableProperty

	// READ-ONLY; Hierarchy information of a product.
	HierarchyInformation *HierarchyInformation

	// READ-ONLY; Image information for the product system.
	ImageInformation []*ImageInformation

	// READ-ONLY; List of product lines supported in the product family
	ProductLines []*ProductLine
}

// ProductLine - Product line
type ProductLine struct {
	// READ-ONLY; Properties of product line
	Properties *ProductLineProperties
}

// ProductLineProperties - Properties of product line
type ProductLineProperties struct {
	// READ-ONLY; Availability information of the product system.
	AvailabilityInformation *AvailabilityInformation

	// READ-ONLY; Cost information for the product system.
	CostInformation *CostInformation

	// READ-ONLY; Description related to the product system.
	Description *Description

	// READ-ONLY; Display Name for the product system.
	DisplayName *string

	// READ-ONLY; list of filters supported for a product
	FilterableProperties []*FilterableProperty

	// READ-ONLY; Hierarchy information of a product.
	HierarchyInformation *HierarchyInformation

	// READ-ONLY; Image information for the product system.
	ImageInformation []*ImageInformation

	// READ-ONLY; List of products in the product line
	Products []*Product
}

// ProductProperties - Properties of products
type ProductProperties struct {
	// READ-ONLY; Availability information of the product system.
	AvailabilityInformation *AvailabilityInformation

	// READ-ONLY; List of configurations for the product
	Configurations []*Configuration

	// READ-ONLY; Cost information for the product system.
	CostInformation *CostInformation

	// READ-ONLY; Description related to the product system.
	Description *Description

	// READ-ONLY; Display Name for the product system.
	DisplayName *string

	// READ-ONLY; list of filters supported for a product
	FilterableProperties []*FilterableProperty

	// READ-ONLY; Hierarchy information of a product.
	HierarchyInformation *HierarchyInformation

	// READ-ONLY; Image information for the product system.
	ImageInformation []*ImageInformation
}

// ProxyResource - The resource model definition for a Azure Resource Manager proxy resource. It will not have tags and a
// location
type ProxyResource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// PurchaseMeterDetails - Billing type Purchase meter details
type PurchaseMeterDetails struct {
	// REQUIRED; Represents billing type.
	BillingType *BillingType

	// READ-ONLY; Charging type.
	ChargingType *ChargingType

	// READ-ONLY; Billing unit applicable for Pav2 billing
	Multiplier *float64

	// READ-ONLY; Product Id
	ProductID *string

	// READ-ONLY; Sku Id
	SKUID *string

	// READ-ONLY; Term Id
	TermID *string
}

// GetMeterDetails implements the MeterDetailsClassification interface for type PurchaseMeterDetails.
func (p *PurchaseMeterDetails) GetMeterDetails() *MeterDetails {
	return &MeterDetails{
		BillingType: p.BillingType,
		ChargingType: p.ChargingType,
		Multiplier: p.Multiplier,
	}
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// ResourceIdentity - Msi identity details of the resource
type ResourceIdentity struct {
	// Identity type
	Type *string

	// READ-ONLY; Service Principal Id backing the Msi
	PrincipalID *string

	// READ-ONLY; Home Tenant Id
	TenantID *string
}

// ResourceProviderDetails - Management RP details
type ResourceProviderDetails struct {
	// READ-ONLY; Resource provider namespace
	ResourceProviderNamespace *string
}

// ReturnOrderItemDetails - Return order item request body
type ReturnOrderItemDetails struct {
	// REQUIRED; Return Reason.
	ReturnReason *string

	// customer return address.
	ReturnAddress *AddressProperties

	// Service tag (located on the bottom-right corner of the device)
	ServiceTag *string

	// Shipping Box required
	ShippingBoxRequired *bool
}

// ReverseShippingDetails - Reverse shipment details.
type ReverseShippingDetails struct {
	// READ-ONLY; Carrier Name for display purpose. Not to be used for any processing.
	CarrierDisplayName *string

	// READ-ONLY; Name of the carrier.
	CarrierName *string

	// READ-ONLY; SAS key to download the reverse shipment label of the package.
	SasKeyForLabel *string

	// READ-ONLY; TrackingId of the package
	TrackingID *string

	// READ-ONLY; TrackingUrl of the package.
	TrackingURL *string
}

// ShippingAddress - Shipping address where customer wishes to receive the device.
type ShippingAddress struct {
	// REQUIRED; Name of the Country.
	Country *string

	// REQUIRED; Street Address line 1.
	StreetAddress1 *string

	// Type of address.
	AddressType *AddressType

	// Name of the City.
	City *string

	// Name of the company.
	CompanyName *string

	// Postal code.
	PostalCode *string

	// Name of the State or Province.
	StateOrProvince *string

	// Street Address line 2.
	StreetAddress2 *string

	// Street Address line 3.
	StreetAddress3 *string

	// Extended Zip Code.
	ZipExtendedCode *string
}

// ShippingDetails - Package shipping details
type ShippingDetails struct {
	// READ-ONLY; Carrier Name for display purpose. Not to be used for any processing.
	CarrierDisplayName *string

	// READ-ONLY; Name of the carrier.
	CarrierName *string

	// READ-ONLY; TrackingId of the package
	TrackingID *string

	// READ-ONLY; TrackingUrl of the package.
	TrackingURL *string
}

// Specifications of the configurations
type Specification struct {
	// READ-ONLY; Name of the specification
	Name *string

	// READ-ONLY; Value of the specification
	Value *string
}

// StageDetails - Resource stage details.
type StageDetails struct {
	// READ-ONLY; Display name of the resource stage.
	DisplayName *string

	// READ-ONLY; Stage name
	StageName *StageName

	// READ-ONLY; Stage status.
	StageStatus *StageStatus

	// READ-ONLY; Stage start time
	StartTime *time.Time
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time

	// The identity that created the resource.
	CreatedBy *string

	// The type of identity that created the resource.
	CreatedByType *CreatedByType

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time

	// The identity that last modified the resource.
	LastModifiedBy *string

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags'
// and a 'location'
type TrackedResource struct {
	// REQUIRED; The geo-location where the resource lives
	Location *string

	// Resource tags.
	Tags map[string]*string

	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string

	// READ-ONLY; The name of the resource
	Name *string

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string
}

// TransportPreferences - Preferences related to the shipment logistics of the sku
type TransportPreferences struct {
	// REQUIRED; Indicates Shipment Logistics type that the customer preferred.
	PreferredShipmentType *TransportShipmentTypes
}

